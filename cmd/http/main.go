package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"notification/internal/config"
	"notification/internal/controller"
	"notification/internal/infra"
	"notification/internal/repository"
	"notification/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Notification API
// @version 1.0
// @description A really nice description

// @contact.name Notification API Support
// @contact.email foo@bar.com

// @host notification
// @BasePath /
func main() {
	// Load the application configuration params
	cfg := config.NewAppConfig()

	// set the controller handlers injecting the dependency
	// in the router
	r := mux.NewRouter()

	// TODO: Health Check controller set up

	// Notification resource controller set up
	cacheService := infra.NewRedisCache()
	rateLimitRulesRepo := repository.NewInMemoryRateLimitRuleRepository()
	rateLimitHandler := service.NewCacheRateLimitHandler(cacheService, rateLimitRulesRepo)
	smtpAddress := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)
	mailClient := infra.NewSMTPMailer(smtpAddress, cfg.MailFrom)
	userRepo := repository.NewInMemoryUserRepository()
	notificationSvc := service.NewEmailNotificationSender(rateLimitHandler, mailClient, userRepo)

	notificationController := controller.NewNotification(notificationSvc)
	notificationController.SetRouter(r)

	// TODO: Set the Swagger endpoint to render the OpenAPI specs.

	// start the HTTP server
	log.Printf("Starting server on port %d", cfg.ServerPort)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: r,
	}

	// Start up the HTTP server in a Go routine
	// to not block the execution so that the Signal listener can
	// take it from there.
	go func() {
		// When the server exits, make sure the error states that the server
		// was closed normally, meaning there's no unexpected error.
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Server is shutting down")
	}()

	// Listen to OS termination signals to allow for a graceful shutdown
	// (especially important in Kubernetes runtimes)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	// define the main context with cancel to release
	// associated resources upon shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call Shutdown for a graceful shutdown.
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server graceful shutdown complete.")
}
