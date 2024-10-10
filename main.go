package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"notification/internal/domain"
	"notification/internal/infra"
	"notification/internal/repository"
	"notification/internal/service"
	"time"
)

func main() {
	// TODO: temporary integration. Create an appropriate User Interface for triggering
	// the application.
	cacheService := infra.NewRedisCache()
	rules := domain.RateLimitRules{
		domain.Status: domain.RateLimitRule{
			MaxCount:   2,
			Expiration: time.Minute * 1,
		},
		domain.News: domain.RateLimitRule{
			MaxCount:   1,
			Expiration: time.Hour * 24,
		},
		domain.Marketing: domain.RateLimitRule{
			MaxCount:   3,
			Expiration: time.Hour * 1,
		},
	}
	rateLimitHandler := service.NewCacheRateLimitHandler(cacheService, rules)
	mailClient := infra.NewSMTPMailer("localhost:1025", "no-reply@example.com")
	userRepo := repository.NewInMemoryUserRepository()
	notificationSvc := service.NewEmailNotificationSender(rateLimitHandler, mailClient, userRepo)

	// Create some test users
	user1 := domain.User{
		ID:       "123-abc",
		Name:     "John",
		LastName: "Doe",
		Email:    "john@example.com",
	}
	_ = userRepo.Save(user1)

	user2 := domain.User{
		ID:       "456-bbb",
		Name:     "Jane",
		LastName: "Doe",
		Email:    "jane@example.com",
	}
	_ = userRepo.Save(user2)

	var counter int
	for counter <= 2 {
		counter++
		err := notificationSvc.Send(context.Background(),
			user1.ID,
			fmt.Sprintf("Hey %s! This is the email #%d", user1.Name, counter), domain.Status)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		err = notificationSvc.Send(context.Background(),
			user2.ID,
			fmt.Sprintf("Hey %s! This is the email #%d", user2.Name, counter), domain.Status)
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}
}
