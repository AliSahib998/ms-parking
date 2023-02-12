package handler

import (
	"github.com/go-chi/chi"
	"net/http"
)

// HealthHandler just checking service is alive
type HealthHandler struct {
}

func NewHealthHandler(router *chi.Mux) *HealthHandler {
	h := &HealthHandler{}
	router.Get("/health", h.Health)
	return h
}

func (*HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
