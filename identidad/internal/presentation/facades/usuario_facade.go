package facades

import (
	"context"
	"time"

	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	uc_createuser "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/createuser"
	uc_deleteuser "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/deleteuser"
	uc_expeluser "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/expeluser"
	uc_listusers "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/listusers"
	uc_updatemyprofile "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/updatemyprofile"
	uc_updateuser "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/updateuser"
	uc_viewmyprofile "github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/viewmyprofile"
)

type ComandoCrearUsuario struct {
	Correo     string
	Nombre     string
	Apellido   string
	Password   string
	EjecutorID string
}

type RespuestaCrearUsuario struct {
	ID       string
	Correo   string
	Nombre   string
	Apellido string
	Activo   bool
	CreadoEn string
}

type ComandoListarUsuarios struct {
	Filtros    []shared_domain.CriterioFiltro
	Paginacion shared_domain.Paginacion
	TenantID   string
	EjecutorID string
}

type RespuestaListarUsuarios struct {
	Usuarios []uc_listusers.UsuarioDTO
	Total    int
	Pagina   int
	Tamano   int
}

type ComandoModificarUsuario struct {
	UsuarioID  string
	Nombre     string
	Apellido   string
	EjecutorID string
}

type RespuestaModificarUsuario struct {
	ID           string
	Correo       string
	Nombre       string
	Apellido     string
	ModificadoEn string
}

type ComandoDarDeBajaUsuario struct {
	UsuarioID  string
	Motivo     string
	EjecutorID string
}

type RespuestaDarDeBajaUsuario struct {
	UsuarioID string
	Estado    string
	BajaEn    string
}

type ComandoExpulsarUsuario struct {
	UsuarioID  string
	EjecutorID string
}

type RespuestaExpulsarUsuario struct {
	UsuarioID         string
	Estado            string
	SesionesRevocadas int
	ExpulsadoEn       string
}

type ComandoVerMiPerfil struct {
	EjecutorID string
}

type RespuestaVerMiPerfil struct {
	ID       string
	Correo   string
	Nombre   string
	Apellido string
	Telefono string
	Estado   string
	CreadoEn string
}

type ComandoModificarMiPerfil struct {
	EjecutorID string
	Nombre     string
	Apellido   string
	Telefono   string
}

type RespuestaModificarMiPerfil struct {
	ID           string
	Correo       string
	Nombre       string
	Apellido     string
	ModificadoEn string
}

type UsuarioFacade interface {
	CrearUsuario(ctx context.Context, cmd ComandoCrearUsuario) (*RespuestaCrearUsuario, error)
	ListarUsuarios(ctx context.Context, cmd ComandoListarUsuarios) (*RespuestaListarUsuarios, error)
	ModificarUsuario(ctx context.Context, cmd ComandoModificarUsuario) (*RespuestaModificarUsuario, error)
	DarDeBajaUsuario(ctx context.Context, cmd ComandoDarDeBajaUsuario) (*RespuestaDarDeBajaUsuario, error)
	ExpulsarUsuario(ctx context.Context, cmd ComandoExpulsarUsuario) (*RespuestaExpulsarUsuario, error)
	VerMiPerfil(ctx context.Context, cmd ComandoVerMiPerfil) (*RespuestaVerMiPerfil, error)
	ModificarMiPerfil(ctx context.Context, cmd ComandoModificarMiPerfil) (*RespuestaModificarMiPerfil, error)
}

type usuarioFacadeImpl struct {
	crearUsuario      *uc_createuser.CrearUsuarioCasoDeUso
	listarUsuarios    *uc_listusers.ListarUsuariosCasoDeUso
	modificarUsuario  *uc_updateuser.ModificarUsuarioCasoDeUso
	darDeBajaUsuario  *uc_deleteuser.DarDeBajaUsuarioCasoDeUso
	expulsarUsuario   *uc_expeluser.ExpulsarUsuarioCasoDeUso
	verMiPerfil       *uc_viewmyprofile.VerMiPerfilCasoDeUso
	modificarMiPerfil *uc_updatemyprofile.ModificarMiPerfilCasoDeUso
}

func NewUsuarioFacade(
	crearUsuario *uc_createuser.CrearUsuarioCasoDeUso,
	listarUsuarios *uc_listusers.ListarUsuariosCasoDeUso,
	modificarUsuario *uc_updateuser.ModificarUsuarioCasoDeUso,
	darDeBajaUsuario *uc_deleteuser.DarDeBajaUsuarioCasoDeUso,
	expulsarUsuario *uc_expeluser.ExpulsarUsuarioCasoDeUso,
	verMiPerfil *uc_viewmyprofile.VerMiPerfilCasoDeUso,
	modificarMiPerfil *uc_updatemyprofile.ModificarMiPerfilCasoDeUso,
) UsuarioFacade {
	return &usuarioFacadeImpl{
		crearUsuario:      crearUsuario,
		listarUsuarios:    listarUsuarios,
		modificarUsuario:  modificarUsuario,
		darDeBajaUsuario:  darDeBajaUsuario,
		expulsarUsuario:   expulsarUsuario,
		verMiPerfil:       verMiPerfil,
		modificarMiPerfil: modificarMiPerfil,
	}
}

func (f *usuarioFacadeImpl) CrearUsuario(ctx context.Context, cmd ComandoCrearUsuario) (*RespuestaCrearUsuario, error) {
	resp, err := f.crearUsuario.Ejecutar(ctx, &uc_createuser.ComandoCrearUsuario{
		Correo:     cmd.Correo,
		Nombre:     cmd.Nombre,
		Apellido:   cmd.Apellido,
		Password:   cmd.Password,
		CreatedBy:  cmd.EjecutorID,
		EjecutorID: cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaCrearUsuario{
		ID:       resp.ID,
		Correo:   resp.Correo,
		Nombre:   resp.Nombre,
		Apellido: resp.Apellido,
		Activo:   resp.Activo,
		CreadoEn: resp.CreadoEn,
	}, nil
}

func (f *usuarioFacadeImpl) ListarUsuarios(ctx context.Context, cmd ComandoListarUsuarios) (*RespuestaListarUsuarios, error) {
	resp, err := f.listarUsuarios.Ejecutar(ctx, &uc_listusers.ComandoListarUsuarios{
		Filtros:    cmd.Filtros,
		Paginacion: cmd.Paginacion,
		TenantID:   cmd.TenantID,
		EjecutorID: cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaListarUsuarios{
		Usuarios: resp.Usuarios,
		Total:    resp.Total,
		Pagina:   resp.Pagina,
		Tamano:   resp.Tamano,
	}, nil
}

func (f *usuarioFacadeImpl) ModificarUsuario(ctx context.Context, cmd ComandoModificarUsuario) (*RespuestaModificarUsuario, error) {
	resp, err := f.modificarUsuario.Ejecutar(ctx, &uc_updateuser.ComandoModificarUsuario{
		UsuarioID:  cmd.UsuarioID,
		Nombre:     cmd.Nombre,
		Apellido:   cmd.Apellido,
		Telefono:   "",
		TenantID:   "",
		EjecutorID: cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaModificarUsuario{
		ID:           resp.ID,
		Correo:       resp.Correo,
		Nombre:       resp.Nombre,
		Apellido:     resp.Apellido,
		ModificadoEn: resp.ModificadoEn,
	}, nil
}

func (f *usuarioFacadeImpl) DarDeBajaUsuario(ctx context.Context, cmd ComandoDarDeBajaUsuario) (*RespuestaDarDeBajaUsuario, error) {
	resp, err := f.darDeBajaUsuario.Ejecutar(ctx, &uc_deleteuser.ComandoDarDeBajaUsuario{
		UsuarioID:  cmd.UsuarioID,
		Motivo:     cmd.Motivo,
		TenantID:   "",
		EjecutorID: cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaDarDeBajaUsuario{
		UsuarioID: resp.UsuarioID,
		Estado:    resp.Estado,
		BajaEn:    resp.BajaEn,
	}, nil
}

func (f *usuarioFacadeImpl) ExpulsarUsuario(ctx context.Context, cmd ComandoExpulsarUsuario) (*RespuestaExpulsarUsuario, error) {
	resp, err := f.expulsarUsuario.Ejecutar(ctx, &uc_expeluser.ComandoExpulsarUsuario{
		UsuarioID:  cmd.UsuarioID,
		TenantID:   "",
		EjecutorID: cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaExpulsarUsuario{
		UsuarioID:         resp.UsuarioID,
		Estado:            resp.Estado,
		SesionesRevocadas: resp.SesionesRevocadas,
		ExpulsadoEn:       resp.ExpulsadoEn,
	}, nil
}

func (f *usuarioFacadeImpl) VerMiPerfil(ctx context.Context, cmd ComandoVerMiPerfil) (*RespuestaVerMiPerfil, error) {
	t := time.Now()
	resp, err := f.verMiPerfil.Ejecutar(ctx, &uc_viewmyprofile.ComandoVerMiPerfil{
		EjecutorID: cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	if resp.CreadoEn == "" {
		return &RespuestaVerMiPerfil{
			ID:       resp.ID,
			Correo:   resp.Correo,
			Nombre:   resp.Nombre,
			Apellido: resp.Apellido,
			Telefono: resp.Telefono,
			Estado:   resp.Estado,
			CreadoEn: t.Format(time.RFC3339),
		}, nil
	}
	return &RespuestaVerMiPerfil{
		ID:       resp.ID,
		Correo:   resp.Correo,
		Nombre:   resp.Nombre,
		Apellido: resp.Apellido,
		Telefono: resp.Telefono,
		Estado:   resp.Estado,
		CreadoEn: resp.CreadoEn,
	}, nil
}

func (f *usuarioFacadeImpl) ModificarMiPerfil(ctx context.Context, cmd ComandoModificarMiPerfil) (*RespuestaModificarMiPerfil, error) {
	resp, err := f.modificarMiPerfil.Ejecutar(ctx, &uc_updatemyprofile.ComandoModificarMiPerfil{
		EjecutorID: cmd.EjecutorID,
		Nombre:     cmd.Nombre,
		Apellido:   cmd.Apellido,
		Telefono:   cmd.Telefono,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaModificarMiPerfil{
		ID:           resp.ID,
		Correo:       resp.Correo,
		Nombre:       resp.Nombre,
		Apellido:     resp.Apellido,
		ModificadoEn: resp.ModificadoEn,
	}, nil
}
