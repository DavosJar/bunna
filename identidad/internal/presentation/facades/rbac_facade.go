package facades

import (
	"context"

	uc_assignpermissiontorole "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/assignpermissiontorole"
	uc_assignrole "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/assignrole"
	uc_createrole "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/createrole"
	uc_deleterole "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/deleterole"
	uc_listarmispermisos "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/listarmispermisos"
	uc_listpermisos "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/listpermisos"
	uc_listroles "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/listroles"
	uc_revokepermissionfromrole "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/revokepermissionfromrole"
	uc_revokerole "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/revokerole"
	uc_updaterole "github.com/davosjar/bunna/services/identidad/internal/rbac/application/usecases/updaterole"
	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type ComandoListarRoles struct {
	Paginacion shared_domain.Paginacion
	EjecutorID string
}

type RespuestaListarRoles struct {
	Roles  []uc_listroles.RolDTO
	Total  int
	Pagina int
}

type RespuestaListarPermisos struct {
	Permisos []uc_listpermisos.PermisoItem
	Total    int
}

type RespuestaListarMisPermisos struct {
	Permisos []string
}

type ComandoCrearRol struct {
	Nombre      string
	Descripcion string
	Permisos    []string
	EjecutorID  string
}

type RespuestaCrearRol struct {
	ID          string
	Nombre      string
	Descripcion string
	EsSistema   bool
	CreadoEn    string
}

type ComandoModificarRol struct {
	RolID       string
	Nombre      string
	Descripcion string
	EjecutorID  string
}

type RespuestaModificarRol struct {
	ID           string
	Nombre       string
	Descripcion  string
	ModificadoEn string
}

type ComandoEliminarRol struct {
	RolID      string
	EjecutorID string
}

type RespuestaEliminarRol struct {
	RolID       string
	EliminadoEn string
}

type ComandoAsignarRol struct {
	UsuarioID  string
	RolID      string
	TenantID   string
	EjecutorID string
}

type RespuestaAsignarRol struct {
	UsuarioID  string
	RolID      string
	TenantID   string
	AsignadoEn string
}

type ComandoRevocarRol struct {
	UsuarioID  string
	RolID      string
	TenantID   string
	EjecutorID string
}

type RespuestaRevocarRol struct {
	UsuarioID  string
	RolID      string
	TenantID   string
	RevocadoEn string
}

type ComandoAsignarPermisoARol struct {
	RolID         string
	PermisoCodigo string
	TenantID      string
	EjecutorID    string
	AsignadoPor   string
}

type RespuestaAsignarPermisoARol struct {
	RolID         string
	PermisoCodigo string
	AsignadoEn    string
}

type ComandoRevocarPermisoDeRol struct {
	RolID         string
	PermisoCodigo string
	TenantID      string
	EjecutorID    string
}

type RespuestaRevocarPermisoDeRol struct {
	RolID         string
	PermisoCodigo string
	RevocadoEn    string
}

type RbacFacade interface {
	ListarRoles(ctx context.Context, cmd ComandoListarRoles) (*RespuestaListarRoles, error)
	ListarPermisos(ctx context.Context, ejecutorID string) (*RespuestaListarPermisos, error)
	ListarMisPermisos(ctx context.Context, rol, tenantID string) (*RespuestaListarMisPermisos, error)
	CrearRol(ctx context.Context, cmd ComandoCrearRol) (*RespuestaCrearRol, error)
	ModificarRol(ctx context.Context, cmd ComandoModificarRol) (*RespuestaModificarRol, error)
	EliminarRol(ctx context.Context, cmd ComandoEliminarRol) (*RespuestaEliminarRol, error)
	AsignarRol(ctx context.Context, cmd ComandoAsignarRol) (*RespuestaAsignarRol, error)
	RevocarRol(ctx context.Context, cmd ComandoRevocarRol) (*RespuestaRevocarRol, error)
	AsignarPermisoARol(ctx context.Context, cmd ComandoAsignarPermisoARol) (*RespuestaAsignarPermisoARol, error)
	RevocarPermisoDeRol(ctx context.Context, cmd ComandoRevocarPermisoDeRol) (*RespuestaRevocarPermisoDeRol, error)
}

type rbacFacadeImpl struct {
	listarRoles         *uc_listroles.ListarRolesCasoDeUso
	listarPermisos      *uc_listpermisos.ListarPermisosCasoDeUso
	listarMisPermisos   *uc_listarmispermisos.ListarMisPermisosCasoDeUso
	crearRol            *uc_createrole.CrearRolCasoDeUso
	modificarRol        *uc_updaterole.ModificarRolCasoDeUso
	eliminarRol         *uc_deleterole.EliminarRolCasoDeUso
	asignarRol          *uc_assignrole.AsignarRolCasoDeUso
	revocarRol          *uc_revokerole.RevocarRolCasoDeUso
	asignarPermisoARol  *uc_assignpermissiontorole.AsignarPermisoARolCasoDeUso
	revocarPermisoDeRol *uc_revokepermissionfromrole.RevocarPermisoDeRolCasoDeUso
}

func NewRbacFacade(
	listarRoles *uc_listroles.ListarRolesCasoDeUso,
	listarPermisos *uc_listpermisos.ListarPermisosCasoDeUso,
	listarMisPermisos *uc_listarmispermisos.ListarMisPermisosCasoDeUso,
	crearRol *uc_createrole.CrearRolCasoDeUso,
	modificarRol *uc_updaterole.ModificarRolCasoDeUso,
	eliminarRol *uc_deleterole.EliminarRolCasoDeUso,
	asignarRol *uc_assignrole.AsignarRolCasoDeUso,
	revocarRol *uc_revokerole.RevocarRolCasoDeUso,
	asignarPermisoARol *uc_assignpermissiontorole.AsignarPermisoARolCasoDeUso,
	revocarPermisoDeRol *uc_revokepermissionfromrole.RevocarPermisoDeRolCasoDeUso,
) RbacFacade {
	return &rbacFacadeImpl{
		listarRoles:         listarRoles,
		listarPermisos:      listarPermisos,
		listarMisPermisos:   listarMisPermisos,
		crearRol:            crearRol,
		modificarRol:        modificarRol,
		eliminarRol:         eliminarRol,
		asignarRol:          asignarRol,
		revocarRol:          revocarRol,
		asignarPermisoARol:  asignarPermisoARol,
		revocarPermisoDeRol: revocarPermisoDeRol,
	}
}

func (f *rbacFacadeImpl) ListarPermisos(ctx context.Context, ejecutorID string) (*RespuestaListarPermisos, error) {
	resp, err := f.listarPermisos.Ejecutar(ctx, ejecutorID)
	if err != nil {
		return nil, err
	}
	return &RespuestaListarPermisos{
		Permisos: resp.Permisos,
		Total:    resp.Total,
	}, nil
}

func (f *rbacFacadeImpl) ListarMisPermisos(ctx context.Context, rol, tenantID string) (*RespuestaListarMisPermisos, error) {
	codigos, err := f.listarMisPermisos.Ejecutar(ctx, rol, tenantID)
	if err != nil {
		return nil, err
	}
	return &RespuestaListarMisPermisos{
		Permisos: codigos,
	}, nil
}

func (f *rbacFacadeImpl) ListarRoles(ctx context.Context, cmd ComandoListarRoles) (*RespuestaListarRoles, error) {
	resp, err := f.listarRoles.Ejecutar(ctx, &uc_listroles.ComandoListarRoles{
		Paginacion: cmd.Paginacion,
		TenantID:   "",
		EjecutorID: cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaListarRoles{
		Roles:  resp.Roles,
		Total:  resp.Total,
		Pagina: resp.Pagina,
	}, nil
}

func (f *rbacFacadeImpl) CrearRol(ctx context.Context, cmd ComandoCrearRol) (*RespuestaCrearRol, error) {
	resp, err := f.crearRol.Ejecutar(ctx, &uc_createrole.ComandoCrearRol{
		Nombre:      cmd.Nombre,
		Descripcion: cmd.Descripcion,
		Permisos:    cmd.Permisos,
		TenantID:    "",
		EjecutorID:  cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaCrearRol{
		ID:          resp.ID,
		Nombre:      resp.Nombre,
		Descripcion: resp.Descripcion,
		EsSistema:   resp.EsSistema,
		CreadoEn:    resp.CreadoEn,
	}, nil
}

func (f *rbacFacadeImpl) ModificarRol(ctx context.Context, cmd ComandoModificarRol) (*RespuestaModificarRol, error) {
	resp, err := f.modificarRol.Ejecutar(ctx, &uc_updaterole.ComandoModificarRol{
		RolID:       cmd.RolID,
		Nombre:      cmd.Nombre,
		Descripcion: cmd.Descripcion,
		TenantID:    "",
		EjecutorID:  cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaModificarRol{
		ID:           resp.ID,
		Nombre:       resp.Nombre,
		Descripcion:  resp.Descripcion,
		ModificadoEn: resp.ModificadoEn,
	}, nil
}

func (f *rbacFacadeImpl) EliminarRol(ctx context.Context, cmd ComandoEliminarRol) (*RespuestaEliminarRol, error) {
	resp, err := f.eliminarRol.Ejecutar(ctx, &uc_deleterole.ComandoEliminarRol{
		RolID:      cmd.RolID,
		TenantID:   "",
		EjecutorID: cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaEliminarRol{
		RolID:       resp.RolID,
		EliminadoEn: resp.EliminadoEn,
	}, nil
}

func (f *rbacFacadeImpl) AsignarRol(ctx context.Context, cmd ComandoAsignarRol) (*RespuestaAsignarRol, error) {
	resp, err := f.asignarRol.Ejecutar(ctx, &uc_assignrole.ComandoAsignarRol{
		UsuarioID:  cmd.UsuarioID,
		RolID:      cmd.RolID,
		TenantID:   cmd.TenantID,
		EjecutorID: cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaAsignarRol{
		UsuarioID:  resp.UsuarioID,
		RolID:      resp.RolID,
		TenantID:   resp.TenantID,
		AsignadoEn: resp.AsignadoEn,
	}, nil
}

func (f *rbacFacadeImpl) RevocarRol(ctx context.Context, cmd ComandoRevocarRol) (*RespuestaRevocarRol, error) {
	resp, err := f.revocarRol.Ejecutar(ctx, &uc_revokerole.ComandoRevocarRol{
		UsuarioID:  cmd.UsuarioID,
		RolID:      cmd.RolID,
		TenantID:   cmd.TenantID,
		EjecutorID: cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaRevocarRol{
		UsuarioID:  resp.UsuarioID,
		RolID:      resp.RolID,
		TenantID:   resp.TenantID,
		RevocadoEn: resp.RevocadoEn,
	}, nil
}

func (f *rbacFacadeImpl) AsignarPermisoARol(ctx context.Context, cmd ComandoAsignarPermisoARol) (*RespuestaAsignarPermisoARol, error) {
	resp, err := f.asignarPermisoARol.Ejecutar(ctx, &uc_assignpermissiontorole.ComandoAsignarPermisoARol{
		RolID:         cmd.RolID,
		PermisoCodigo: cmd.PermisoCodigo,
		TenantID:      cmd.TenantID,
		EjecutorID:    cmd.EjecutorID,
		AsignadoPor:   cmd.AsignadoPor,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaAsignarPermisoARol{
		RolID:         resp.RolID,
		PermisoCodigo: resp.PermisoCodigo,
		AsignadoEn:    resp.AsignadoEn,
	}, nil
}

func (f *rbacFacadeImpl) RevocarPermisoDeRol(ctx context.Context, cmd ComandoRevocarPermisoDeRol) (*RespuestaRevocarPermisoDeRol, error) {
	resp, err := f.revocarPermisoDeRol.Ejecutar(ctx, &uc_revokepermissionfromrole.ComandoRevocarPermisoDeRol{
		RolID:         cmd.RolID,
		PermisoCodigo: cmd.PermisoCodigo,
		TenantID:      cmd.TenantID,
		EjecutorID:    cmd.EjecutorID,
	})
	if err != nil {
		return nil, err
	}
	return &RespuestaRevocarPermisoDeRol{
		RolID:         resp.RolID,
		PermisoCodigo: resp.PermisoCodigo,
		RevocadoEn:    resp.RevocadoEn,
	}, nil
}
