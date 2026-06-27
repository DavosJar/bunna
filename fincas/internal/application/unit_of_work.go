package application

import (
	"context"

	diagnosticodomain "github.com/davosjar/bunna/services/fincas/internal/diagnostico/domain"
)

// UnitOfWorkDiagnostico abstrae las transacciones que abarcan
// DiagnosticoRepositorio, MuestraRepositorio y CandidatoReentrenamientoRepositorio.
// La implementación concreta vive en infraestructura.
type UnitOfWorkDiagnostico interface {
	Transaccional(ctx context.Context, fn func(UnitOfWorkDiagnostico) error) error
	DiagnosticoRepo() diagnosticodomain.DiagnosticoRepositorio
	MuestraRepo() diagnosticodomain.MuestraRepositorio
	CandidatoRepo() diagnosticodomain.CandidatoReentrenamientoRepositorio
}
