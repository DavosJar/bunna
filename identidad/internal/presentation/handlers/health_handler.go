package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

// HealthOutput es la respuesta del health check.
type HealthOutput struct {
	Body struct {
		Status string `json:"status" doc:"Estado del servicio" example:"ok"`
	}
}

// RegisterHealthHandler registra el endpoint GET /health en la API Huma.
func RegisterHealthHandler(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "health-check",
		Method:      "GET",
		Path:        "/api/v1/identidad/health",
		Summary:     "Health check",
		Description: "Verifica que el servicio está activo y respondiendo.",
		Tags:        []string{"Sistema"},
	}, func(ctx context.Context, input *struct{}) (*HealthOutput, error) {
		resp := &HealthOutput{}
		resp.Body.Status = "ok"
		return resp, nil
	})
}
