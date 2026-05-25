package dto

import "time"

// RegistrarFincaRequest es el body para crear una finca.
type RegistrarFincaRequest struct {
	Nombre      string `json:"nombre"`
	Ubicacion   string `json:"ubicacion"`
	Descripcion string `json:"descripcion"`
}

// DesactivarFincaRequest es el body para desactivar una finca.
type DesactivarFincaRequest struct {
	Confirmar bool `json:"confirmar"`
}

// FincaResponse es la respuesta de una finca (registrar, obtener, listar).
type FincaResponse struct {
	ID          string    `json:"id"`
	Nombre      string    `json:"nombre"`
	Ubicacion   string    `json:"ubicacion"`
	Descripcion string    `json:"descripcion"`
	Estado      string    `json:"estado"`
	CreatedAt   time.Time `json:"createdAt"`
}
