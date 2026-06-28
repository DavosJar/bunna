package facades

import (
	"context"

	decorator "github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry/decorator"
	uc_aceptar "github.com/davosjar/bunna/services/identidad/internal/invitaciones/application/usecases/aceptarinvitacion"
	uc_crear "github.com/davosjar/bunna/services/identidad/internal/invitaciones/application/usecases/crearinvitacion"
	uc_eliminar "github.com/davosjar/bunna/services/identidad/internal/invitaciones/application/usecases/eliminarinvitacion"
	uc_listar "github.com/davosjar/bunna/services/identidad/internal/invitaciones/application/usecases/listarinvitaciones"
	uc_obtener "github.com/davosjar/bunna/services/identidad/internal/invitaciones/application/usecases/obtenerinvitacion"
	uc_reenviar "github.com/davosjar/bunna/services/identidad/internal/invitaciones/application/usecases/reenviarinvitacion"
)

type ComandoCrearInvitacion struct {
	TenantID  string
	RolID     string
	Correo    string
	CreadoPor string
}

type RespuestaCrearInvitacion struct {
	Mensaje string
}

type ComandoAceptarInvitacion struct {
	Token string
}

type RespuestaAceptarInvitacion struct {
	TenantID string
	RolID    string
}

type RespuestaObtenerInvitacion struct {
	ID           string
	TenantID     string
	TenantNombre string
	RolID        string
	RolNombre    string
	Email        string
	Estado       string
	Expiracion   string
}

type ComandoListarInvitaciones struct {
	TenantID     string
	Pagina       int
	TamanoPagina int
	Estado       string
}

type RespuestaListarInvitacionesItem struct {
	ID            string
	Email         string
	Nombre        string
	RolID         string
	RolNombre     string
	Estado        string
	FechaCreacion string
	Expiracion    string
}

type RespuestaListarInvitaciones struct {
	Invitaciones []RespuestaListarInvitacionesItem
	Total        int
}

type ComandoReenviarInvitacion struct {
	InvitacionID string
	TenantID     string
}

type RespuestaReenviarInvitacion struct {
	Mensaje string
}

type ComandoEliminarInvitacion struct {
	InvitacionID string
	TenantID     string
	EjecutorID   string
}

type RespuestaEliminarInvitacion struct {
	Mensaje string
}

type InvitacionFacade interface {
	CrearInvitacion(ctx context.Context, cmd ComandoCrearInvitacion) (*RespuestaCrearInvitacion, error)
	AceptarInvitacion(ctx context.Context, cmd ComandoAceptarInvitacion) (*RespuestaAceptarInvitacion, error)
	ObtenerInvitacion(ctx context.Context, token string) (*RespuestaObtenerInvitacion, error)
	ListarInvitaciones(ctx context.Context, cmd ComandoListarInvitaciones) (*RespuestaListarInvitaciones, error)
	ReenviarInvitacion(ctx context.Context, cmd ComandoReenviarInvitacion) (*RespuestaReenviarInvitacion, error)
	EliminarInvitacion(ctx context.Context, cmd ComandoEliminarInvitacion) (*RespuestaEliminarInvitacion, error)
}

type invitacionFacadeImpl struct {
	crearUseCase    decorator.UseCase[*uc_crear.ComandoCrearInvitacion, *uc_crear.RespuestaCrearInvitacion]
	aceptarUseCase  decorator.UseCase[*uc_aceptar.ComandoAceptarInvitacion, *uc_aceptar.RespuestaAceptarInvitacion]
	obtenerUseCase  decorator.UseCase[*uc_obtener.ComandoObtenerInvitacion, *uc_obtener.RespuestaObtenerInvitacion]
	listarUseCase   decorator.UseCase[*uc_listar.ComandoListarInvitaciones, *uc_listar.RespuestaListarInvitaciones]
	reenviarUseCase decorator.UseCase[*uc_reenviar.ComandoReenviarInvitacion, *uc_reenviar.RespuestaReenviarInvitacion]
	eliminarUseCase decorator.UseCase[*uc_eliminar.ComandoEliminarInvitacion, *uc_eliminar.RespuestaEliminarInvitacion]
}

func NewInvitacionFacade(
	crearUseCase decorator.UseCase[*uc_crear.ComandoCrearInvitacion, *uc_crear.RespuestaCrearInvitacion],
	aceptarUseCase decorator.UseCase[*uc_aceptar.ComandoAceptarInvitacion, *uc_aceptar.RespuestaAceptarInvitacion],
	obtenerUseCase decorator.UseCase[*uc_obtener.ComandoObtenerInvitacion, *uc_obtener.RespuestaObtenerInvitacion],
	listarUseCase decorator.UseCase[*uc_listar.ComandoListarInvitaciones, *uc_listar.RespuestaListarInvitaciones],
	reenviarUseCase decorator.UseCase[*uc_reenviar.ComandoReenviarInvitacion, *uc_reenviar.RespuestaReenviarInvitacion],
	eliminarUseCase decorator.UseCase[*uc_eliminar.ComandoEliminarInvitacion, *uc_eliminar.RespuestaEliminarInvitacion],
) InvitacionFacade {
	return &invitacionFacadeImpl{
		crearUseCase:    crearUseCase,
		aceptarUseCase:  aceptarUseCase,
		obtenerUseCase:  obtenerUseCase,
		listarUseCase:   listarUseCase,
		reenviarUseCase: reenviarUseCase,
		eliminarUseCase: eliminarUseCase,
	}
}

func (f *invitacionFacadeImpl) CrearInvitacion(ctx context.Context, cmd ComandoCrearInvitacion) (*RespuestaCrearInvitacion, error) {
	_, err := f.crearUseCase.Ejecutar(ctx, &uc_crear.ComandoCrearInvitacion{
		TenantID:  cmd.TenantID,
		RolID:     cmd.RolID,
		Correo:    cmd.Correo,
		CreadoPor: cmd.CreadoPor,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaCrearInvitacion{
		Mensaje: "Invitación enviada exitosamente",
	}, nil
}

func (f *invitacionFacadeImpl) AceptarInvitacion(ctx context.Context, cmd ComandoAceptarInvitacion) (*RespuestaAceptarInvitacion, error) {
	resp, err := f.aceptarUseCase.Ejecutar(ctx, &uc_aceptar.ComandoAceptarInvitacion{
		Token: cmd.Token,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaAceptarInvitacion{
		TenantID: resp.TenantID,
		RolID:    resp.RolID,
	}, nil
}

func (f *invitacionFacadeImpl) ListarInvitaciones(ctx context.Context, cmd ComandoListarInvitaciones) (*RespuestaListarInvitaciones, error) {
	resp, err := f.listarUseCase.Ejecutar(ctx, &uc_listar.ComandoListarInvitaciones{
		TenantID:     cmd.TenantID,
		Pagina:       cmd.Pagina,
		TamanoPagina: cmd.TamanoPagina,
		Estado:       cmd.Estado,
	})
	if err != nil {
		return nil, err
	}

	items := make([]RespuestaListarInvitacionesItem, len(resp.Invitaciones))
	for i, inv := range resp.Invitaciones {
		items[i] = RespuestaListarInvitacionesItem{
			ID:            inv.ID,
			Email:         inv.Email,
			Nombre:        inv.Nombre,
			RolID:         inv.RolID,
			RolNombre:     inv.RolNombre,
			Estado:        inv.Estado,
			FechaCreacion: inv.FechaCreacion,
			Expiracion:    inv.Expiracion,
		}
	}

	return &RespuestaListarInvitaciones{
		Invitaciones: items,
		Total:        resp.Total,
	}, nil
}

func (f *invitacionFacadeImpl) ReenviarInvitacion(ctx context.Context, cmd ComandoReenviarInvitacion) (*RespuestaReenviarInvitacion, error) {
	_, err := f.reenviarUseCase.Ejecutar(ctx, &uc_reenviar.ComandoReenviarInvitacion{
		InvitacionID: cmd.InvitacionID,
		TenantID:     cmd.TenantID,
	})
	if err != nil {
		return nil, err
	}

	return &RespuestaReenviarInvitacion{
		Mensaje: "Invitación reenviada exitosamente",
	}, nil
}

func (f *invitacionFacadeImpl) EliminarInvitacion(ctx context.Context, cmd ComandoEliminarInvitacion) (*RespuestaEliminarInvitacion, error) {
	resp, err := f.eliminarUseCase.Ejecutar(ctx, &uc_eliminar.ComandoEliminarInvitacion{
		InvitacionID: cmd.InvitacionID,
		TenantID:     cmd.TenantID,
		EjecutorID:   cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}

	return &RespuestaEliminarInvitacion{
		Mensaje: resp.Mensaje,
	}, nil
}

func (f *invitacionFacadeImpl) ObtenerInvitacion(ctx context.Context, token string) (*RespuestaObtenerInvitacion, error) {
	resp, err := f.obtenerUseCase.Ejecutar(ctx, &uc_obtener.ComandoObtenerInvitacion{
		Token: token,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaObtenerInvitacion{
		ID:           resp.ID,
		TenantID:     resp.TenantID,
		TenantNombre: resp.TenantNombre,
		RolID:        resp.RolID,
		RolNombre:    resp.RolNombre,
		Email:        resp.Email,
		Estado:       resp.Estado,
		Expiracion:   resp.Expiracion,
	}, nil
}
