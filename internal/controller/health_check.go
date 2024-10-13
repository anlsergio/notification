package controller

import (
	"github.com/gorilla/mux"
	"net/http"
)

// NewHealthCheck creates a new HealthCheck controller instance.
func NewHealthCheck() *HealthCheck {
	return &HealthCheck{}
}

// HealthCheck is the health check controller.
// It defines routes and handlers to serve Liveness and Readiness probes
// using the "z" suffix convention: https://kubernetes.io/docs/reference/using-api/health-checks/
type HealthCheck struct{}

// SetRouter returns the router r with all the necessary routes for the
// HealthCheck controller setup.
func (h HealthCheck) SetRouter(r *mux.Router) {
	r.HandleFunc("/healthz", h.checkHealth).Methods(http.MethodGet)
	r.HandleFunc("/readyz", h.checkReady).Methods(http.MethodGet)
}

func (h HealthCheck) checkHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h HealthCheck) checkReady(w http.ResponseWriter, r *http.Request) {
	// TODO: check if dependency servers are ready, such as DB servers, Messaging brokers, ...
	w.WriteHeader(http.StatusOK)
}
