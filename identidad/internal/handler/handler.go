package handler

import (
	"net/http"

	"github.com/davosjar/bunna/services/identidad/internal/registry"
)

type Handler struct {
	registry *registry.Registry
}

func NewHandler(reg *registry.Registry) *Handler {
	return &Handler{registry: reg}
}

// Health check endpoint
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// Register handlers to router
func (h *Handler) RegisterRoutes() {
	http.HandleFunc("/health", h.Health)
}
