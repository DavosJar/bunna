package facades

import (
	"context"

	uc_listsessions "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/listsessions"
	uc_terminatesession "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/terminatesession"
	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type ComandoListarSesiones struct {
	Paginacion shared_domain.Paginacion
	EjecutorID string
}

type RespuestaListarSesiones struct {
	Sesiones []uc_listsessions.SesionDTO
	Total    int
	Pagina   int
}

type ComandoForzarCierreSesion struct {
	SesionID   string
	EjecutorID string
}

type RespuestaForzarCierreSesion struct {
	SesionID   string
	Estado     string
	RevocadoEn string
}

type SesionFacade interface {
	ListarSesiones(ctx context.Context, cmd ComandoListarSesiones) (*RespuestaListarSesiones, error)
	ForzarCierreSesion(ctx context.Context, cmd ComandoForzarCierreSesion) (*RespuestaForzarCierreSesion, error)
}

type sesionFacadeImpl struct {
	listarSesiones     *uc_listsessions.ListarSesionesCasoDeUso
	forzarCierreSesion *uc_terminatesession.ForzarCierreSesionCasoDeUso
}

func NewSesionFacade(
	listarSesiones *uc_listsessions.ListarSesionesCasoDeUso,
	forzarCierreSesion *uc_terminatesession.ForzarCierreSesionCasoDeUso,
) SesionFacade {
	return &sesionFacadeImpl{
		listarSesiones:     listarSesiones,
		forzarCierreSesion: forzarCierreSesion,
	}
}

func (f *sesionFacadeImpl) ListarSesiones(ctx context.Context, cmd ComandoListarSesiones) (*RespuestaListarSesiones, error) {
	resp, err := f.listarSesiones.Ejecutar(ctx, &uc_listsessions.ComandoListarSesiones{
		Paginacion: cmd.Paginacion,
		EjecutorID: cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaListarSesiones{
		Sesiones: resp.Sesiones,
		Total:    resp.Total,
		Pagina:   resp.Pagina,
	}, nil
}

func (f *sesionFacadeImpl) ForzarCierreSesion(ctx context.Context, cmd ComandoForzarCierreSesion) (*RespuestaForzarCierreSesion, error) {
	resp, err := f.forzarCierreSesion.Ejecutar(ctx, &uc_terminatesession.ComandoForzarCierreSesion{
		SesionID:   cmd.SesionID,
		EjecutorID: cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaForzarCierreSesion{
		SesionID:   resp.SesionID,
		Estado:     resp.Estado,
		RevocadoEn: resp.RevocadoEn,
	}, nil
}
