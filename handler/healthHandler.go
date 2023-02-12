package handler

import (
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

// HealthHandler just checking service is alive
type HealthHandler struct {
}

func NewHealthHandler(router *chi.Mux) *HealthHandler {
	h := &HealthHandler{}
	router.Get("/health", h.Health)
	router.Get("/swagger/*", httpSwagger.Handler())
	return h
}

// Health godoc
// @Summary Health endpoint for kubernetes health and readiness check
// @Tags health-handler
// @Success 200 {} http.Response
// @Router /health [get]
func (*HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
