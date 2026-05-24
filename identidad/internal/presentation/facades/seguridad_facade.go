package facades

import (
	"context"

	uc_changemypassword "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/changemypassword"
	uc_listblockedips "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/listblockedips"
	uc_resetpassword "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/resetpassword"
	uc_unblockip "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/unblockip"
	uc_unlockaccount "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/unlockaccount"
	uc_viewcredentials "github.com/davosjar/bunna/services/identidad/internal/seguridad/application/usecases/viewcredentials"
	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type ComandoCambiarMiPassword struct {
	EjecutorID     string
	PasswordActual string
	NuevaPassword  string
}

type RespuestaCambiarMiPassword struct {
	ModificadoEn string
}

type ComandoResetearPassword struct {
	UsuarioID     string
	NuevaPassword string
	EjecutorID    string
}

type RespuestaResetearPassword struct {
	UsuarioID    string
	ModificadoEn string
}

type ComandoDesbloquearCuenta struct {
	UsuarioID  string
	EjecutorID string
}

type RespuestaDesbloquearCuenta struct {
	UsuarioID      string
	DesbloqueadoEn string
}

type ComandoListarIPsBloqueadas struct {
	Paginacion shared_domain.Paginacion
	EjecutorID string
}

type RespuestaListarIPsBloqueadas struct {
	IPs    []uc_listblockedips.IPBloqueadaDTO
	Total  int
	Pagina int
}

type ComandoDesbloquearIP struct {
	IP         string
	EjecutorID string
}

type RespuestaDesbloquearIP struct {
	IP             string
	DesbloqueadoEn string
}

type ComandoConsultarCredenciales struct {
	UsuarioID  string
	EjecutorID string
}

type RespuestaConsultarCredenciales struct {
	UsuarioID        string
	Activo           bool
	CorreoVerificado bool
	IntentosFallidos int
	BloqueadoHasta   string
}

type SeguridadFacade interface {
	CambiarMiPassword(ctx context.Context, cmd ComandoCambiarMiPassword) (*RespuestaCambiarMiPassword, error)
	ResetearPassword(ctx context.Context, cmd ComandoResetearPassword) (*RespuestaResetearPassword, error)
	DesbloquearCuenta(ctx context.Context, cmd ComandoDesbloquearCuenta) (*RespuestaDesbloquearCuenta, error)
	ListarIPsBloqueadas(ctx context.Context, cmd ComandoListarIPsBloqueadas) (*RespuestaListarIPsBloqueadas, error)
	DesbloquearIP(ctx context.Context, cmd ComandoDesbloquearIP) (*RespuestaDesbloquearIP, error)
	ConsultarCredenciales(ctx context.Context, cmd ComandoConsultarCredenciales) (*RespuestaConsultarCredenciales, error)
}

type seguridadFacadeImpl struct {
	cambiarMiPassword     *uc_changemypassword.CambiarMiContrasenaCasoDeUso
	resetearPassword      *uc_resetpassword.ResetearContrasenaCasoDeUso
	desbloquearCuenta     *uc_unlockaccount.DesbloquearCuentaCasoDeUso
	listarIPsBloqueadas   *uc_listblockedips.ListarIPsBloqueadasCasoDeUso
	desbloquearIP         *uc_unblockip.DesbloquearIPCasoDeUso
	consultarCredenciales *uc_viewcredentials.ConsultarCredencialesCasoDeUso
}

func NewSeguridadFacade(
	cambiarMiPassword *uc_changemypassword.CambiarMiContrasenaCasoDeUso,
	resetearPassword *uc_resetpassword.ResetearContrasenaCasoDeUso,
	desbloquearCuenta *uc_unlockaccount.DesbloquearCuentaCasoDeUso,
	listarIPsBloqueadas *uc_listblockedips.ListarIPsBloqueadasCasoDeUso,
	desbloquearIP *uc_unblockip.DesbloquearIPCasoDeUso,
	consultarCredenciales *uc_viewcredentials.ConsultarCredencialesCasoDeUso,
) SeguridadFacade {
	return &seguridadFacadeImpl{
		cambiarMiPassword:     cambiarMiPassword,
		resetearPassword:      resetearPassword,
		desbloquearCuenta:     desbloquearCuenta,
		listarIPsBloqueadas:   listarIPsBloqueadas,
		desbloquearIP:         desbloquearIP,
		consultarCredenciales: consultarCredenciales,
	}
}

func (f *seguridadFacadeImpl) CambiarMiPassword(ctx context.Context, cmd ComandoCambiarMiPassword) (*RespuestaCambiarMiPassword, error) {
	resp, err := f.cambiarMiPassword.Ejecutar(ctx, &uc_changemypassword.ComandoCambiarMiContrasena{
		EjecutorID:     cmd.EjecutorID,
		PasswordActual: cmd.PasswordActual,
		NuevaPassword:  cmd.NuevaPassword,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaCambiarMiPassword{
		ModificadoEn: resp.ModificadoEn,
	}, nil
}

func (f *seguridadFacadeImpl) ResetearPassword(ctx context.Context, cmd ComandoResetearPassword) (*RespuestaResetearPassword, error) {
	resp, err := f.resetearPassword.Ejecutar(ctx, &uc_resetpassword.ComandoResetearContrasena{
		UsuarioID:     cmd.UsuarioID,
		NuevaPassword: cmd.NuevaPassword,
		TenantID:      "",
		EjecutorID:    cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaResetearPassword{
		UsuarioID:    resp.UsuarioID,
		ModificadoEn: resp.ModificadoEn,
	}, nil
}

func (f *seguridadFacadeImpl) DesbloquearCuenta(ctx context.Context, cmd ComandoDesbloquearCuenta) (*RespuestaDesbloquearCuenta, error) {
	resp, err := f.desbloquearCuenta.Ejecutar(ctx, &uc_unlockaccount.ComandoDesbloquearCuenta{
		UsuarioID:  cmd.UsuarioID,
		TenantID:   "",
		EjecutorID: cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaDesbloquearCuenta{
		UsuarioID:      resp.UsuarioID,
		DesbloqueadoEn: resp.DesbloqueadoEn,
	}, nil
}

func (f *seguridadFacadeImpl) ListarIPsBloqueadas(ctx context.Context, cmd ComandoListarIPsBloqueadas) (*RespuestaListarIPsBloqueadas, error) {
	resp, err := f.listarIPsBloqueadas.Ejecutar(ctx, &uc_listblockedips.ComandoListarIPsBloqueadas{
		Paginacion: cmd.Paginacion,
		TenantID:   "",
		EjecutorID: cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaListarIPsBloqueadas{
		IPs:    resp.IPs,
		Total:  resp.Total,
		Pagina: resp.Pagina,
	}, nil
}

func (f *seguridadFacadeImpl) DesbloquearIP(ctx context.Context, cmd ComandoDesbloquearIP) (*RespuestaDesbloquearIP, error) {
	resp, err := f.desbloquearIP.Ejecutar(ctx, &uc_unblockip.ComandoDesbloquearIP{
		IP:         cmd.IP,
		TenantID:   "",
		EjecutorID: cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaDesbloquearIP{
		IP:             resp.IP,
		DesbloqueadoEn: resp.DesbloqueadoEn,
	}, nil
}

func (f *seguridadFacadeImpl) ConsultarCredenciales(ctx context.Context, cmd ComandoConsultarCredenciales) (*RespuestaConsultarCredenciales, error) {
	resp, err := f.consultarCredenciales.Ejecutar(ctx, &uc_viewcredentials.ComandoConsultarCredenciales{
		UsuarioID:  cmd.UsuarioID,
		TenantID:   "",
		EjecutorID: cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaConsultarCredenciales{
		UsuarioID:        resp.UsuarioID,
		Activo:           resp.Activo,
		CorreoVerificado: resp.CorreoVerificado,
		IntentosFallidos: resp.IntentosFallidos,
		BloqueadoHasta:   resp.BloqueadoHasta,
	}, nil
}
