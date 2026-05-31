package domain

import (
	"context"

	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

type DiagnosticoRepositorio interface {
	Crear(ctx context.Context, diagnostico *Diagnostico) error
	ObtenerPorID(ctx context.Context, id string) (*Diagnostico, error)
	ListarPorFinca(ctx context.Context, fincaID string) ([]Diagnostico, error)
	Actualizar(ctx context.Context, diagnostico *Diagnostico) error
	Eliminar(ctx context.Context, id string) error
	Buscar(ctx context.Context, especificacion EspecificacionDiagnostico, paginacion shared.Paginacion) ([]Diagnostico, error)
}

type MuestraRepositorio interface {
	Crear(ctx context.Context, muestra *Muestra) error
	ObtenerPorID(ctx context.Context, id string) (*Muestra, error)
	ListarPorDiagnostico(ctx context.Context, diagnosticoID string) ([]Muestra, error)
	Actualizar(ctx context.Context, muestra *Muestra) error
	Eliminar(ctx context.Context, id string) error
	Buscar(ctx context.Context, especificacion EspecificacionMuestra, paginacion shared.Paginacion) ([]Muestra, error)
}

type CandidatoReentrenamientoRepositorio interface {
	Crear(ctx context.Context, candidato *CandidatoReentrenamiento) error
	ObtenerPorDiagnosticoID(ctx context.Context, diagnosticoID string) (*CandidatoReentrenamiento, error)
	ListarPendientes(ctx context.Context, limite int) ([]CandidatoReentrenamiento, error)
}

