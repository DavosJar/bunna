package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/dto"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/facades"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/middleware"
	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	presentation "github.com/davosjar/bunna/services/identidad/shared/presentation"
)

// ── Listar Roles ───────────────────────────────────────────────────────────────

type ListarRolesInput struct {
	Pagina       int `query:"pagina"       doc:"Número de página (1-based)" example:"1"`
	TamanoPagina int `query:"tamano"       doc:"Elementos por página"       example:"20"`
}

type ListarRolesOutput struct {
	Body presentation.ApiResponse[dto.ListarRolesResponse]
}

type ListarRolesHandler struct {
	facade facades.RbacFacade
}

func NewListarRolesHandler(facade facades.RbacFacade) *ListarRolesHandler {
	return &ListarRolesHandler{facade: facade}
}

func (h *ListarRolesHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "listar-roles",
		Method:      http.MethodGet,
		Path:        "/api/v1/roles",
		Summary:     "Listar roles",
		Description: "Lista los roles del sistema con paginación.",
		Tags:        []string{"Roles"},
	}, h.handle)
}

func (h *ListarRolesHandler) handle(ctx context.Context, input *ListarRolesInput) (*ListarRolesOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	tenantID := middleware.GetTenantIDFromCtx(ctx)

	pagina := input.Pagina
	if pagina < 1 {
		pagina = 1
	}
	tamano := input.TamanoPagina
	if tamano < 1 || tamano > 100 {
		tamano = 20
	}

	resp, err := h.facade.ListarRoles(ctx, facades.ComandoListarRoles{
		Paginacion: shared_domain.Paginacion{Pagina: pagina, TamanoPagina: tamano},
		TenantID:   tenantID,
		EjecutorID: ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	items := make([]dto.RolItem, len(resp.Roles))
	for i, r := range resp.Roles {
		items[i] = dto.RolItem{
			ID: r.ID, Nombre: r.Nombre, Descripcion: r.Descripcion,
			EsSistema: r.EsSistema, Permisos: r.Permisos,
		}
	}

	out := &ListarRolesOutput{}
	out.Body = presentation.NewApiResponse(dto.ListarRolesResponse{
		Roles: items, Total: resp.Total, Pagina: resp.Pagina,
	})
	return out, nil
}

// ── Crear Rol ──────────────────────────────────────────────────────────────────

type CrearRolInput struct {
	Body dto.CrearRolRequest
}

type CrearRolOutput struct {
	Body presentation.ApiResponse[dto.CrearRolResponse]
}

type CrearRolHandler struct {
	facade facades.RbacFacade
}

func NewCrearRolHandler(facade facades.RbacFacade) *CrearRolHandler {
	return &CrearRolHandler{facade: facade}
}

func (h *CrearRolHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "crear-rol",
		Method:        http.MethodPost,
		Path:          "/api/v1/roles",
		Summary:       "Crear rol",
		Description:   "Crea un nuevo rol en el sistema con permisos opcionales.",
		Tags:          []string{"Roles"},
		DefaultStatus: http.StatusCreated,
	}, h.handle)
}

func (h *CrearRolHandler) handle(ctx context.Context, input *CrearRolInput) (*CrearRolOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	tenantID := middleware.GetTenantIDFromCtx(ctx)

	resp, err := h.facade.CrearRol(ctx, facades.ComandoCrearRol{
		Nombre:      input.Body.Nombre,
		Descripcion: input.Body.Descripcion,
		Permisos:    input.Body.Permisos,
		TenantID:    tenantID,
		EjecutorID:  ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &CrearRolOutput{}
	out.Body = presentation.NewApiResponse(dto.CrearRolResponse{
		ID: resp.ID, Nombre: resp.Nombre, Descripcion: resp.Descripcion,
		EsSistema: resp.EsSistema, CreadoEn: resp.CreadoEn,
	})
	return out, nil
}

// ── Modificar Rol ──────────────────────────────────────────────────────────────

type ModificarRolInput struct {
	RolID string `path:"rolID" doc:"ID del rol a modificar"`
	Body  dto.ModificarRolRequest
}

type ModificarRolOutput struct {
	Body presentation.ApiResponse[dto.ModificarRolResponse]
}

type ModificarRolHandler struct {
	facade facades.RbacFacade
}

func NewModificarRolHandler(facade facades.RbacFacade) *ModificarRolHandler {
	return &ModificarRolHandler{facade: facade}
}

func (h *ModificarRolHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "modificar-rol",
		Method:      http.MethodPut,
		Path:        "/api/v1/roles/{rolID}",
		Summary:     "Modificar rol",
		Description: "Actualiza el nombre y descripción de un rol.",
		Tags:        []string{"Roles"},
	}, h.handle)
}

func (h *ModificarRolHandler) handle(ctx context.Context, input *ModificarRolInput) (*ModificarRolOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	tenantID := middleware.GetTenantIDFromCtx(ctx)

	resp, err := h.facade.ModificarRol(ctx, facades.ComandoModificarRol{
		RolID:       input.RolID,
		Nombre:      input.Body.Nombre,
		Descripcion: input.Body.Descripcion,
		TenantID:    tenantID,
		EjecutorID:  ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &ModificarRolOutput{}
	out.Body = presentation.NewApiResponse(dto.ModificarRolResponse{
		ID: resp.ID, Nombre: resp.Nombre, Descripcion: resp.Descripcion,
		ModificadoEn: resp.ModificadoEn,
	})
	return out, nil
}

// ── Eliminar Rol ───────────────────────────────────────────────────────────────

type EliminarRolInput struct {
	RolID string `path:"rolID" doc:"ID del rol a eliminar"`
}

type EliminarRolOutput struct {
	Body presentation.ApiResponse[dto.EliminarRolResponse]
}

type EliminarRolHandler struct {
	facade facades.RbacFacade
}

func NewEliminarRolHandler(facade facades.RbacFacade) *EliminarRolHandler {
	return &EliminarRolHandler{facade: facade}
}

func (h *EliminarRolHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "eliminar-rol",
		Method:      http.MethodDelete,
		Path:        "/api/v1/roles/{rolID}",
		Summary:     "Eliminar rol",
		Description: "Elimina un rol del sistema (no se pueden eliminar roles de sistema).",
		Tags:        []string{"Roles"},
	}, h.handle)
}

func (h *EliminarRolHandler) handle(ctx context.Context, input *EliminarRolInput) (*EliminarRolOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	tenantID := middleware.GetTenantIDFromCtx(ctx)

	resp, err := h.facade.EliminarRol(ctx, facades.ComandoEliminarRol{
		RolID:      input.RolID,
		TenantID:   tenantID,
		EjecutorID: ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &EliminarRolOutput{}
	out.Body = presentation.NewApiResponse(dto.EliminarRolResponse{
		RolID: resp.RolID, EliminadoEn: resp.EliminadoEn,
	})
	return out, nil
}

// ── Asignar Rol a Usuario ─────────────────────────────────────────────────────

type AsignarRolInput struct {
	UsuarioID string `path:"usuarioID" doc:"ID del usuario"`
	Body      dto.AsignarRolRequest
}

type AsignarRolOutput struct {
	Body presentation.ApiResponse[dto.AsignarRolResponse]
}

type AsignarRolHandler struct {
	facade facades.RbacFacade
}

func NewAsignarRolHandler(facade facades.RbacFacade) *AsignarRolHandler {
	return &AsignarRolHandler{facade: facade}
}

func (h *AsignarRolHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "asignar-rol",
		Method:        http.MethodPost,
		Path:          "/api/v1/usuarios/{usuarioID}/roles",
		Summary:       "Asignar rol a usuario",
		Description:   "Asigna un rol a un usuario, opcionalmente en un tenant específico.",
		Tags:          []string{"Roles"},
		DefaultStatus: http.StatusCreated,
	}, h.handle)
}

func (h *AsignarRolHandler) handle(ctx context.Context, input *AsignarRolInput) (*AsignarRolOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	resp, err := h.facade.AsignarRol(ctx, facades.ComandoAsignarRol{
		UsuarioID:  input.UsuarioID,
		RolID:      input.Body.RolID,
		TenantID:   input.Body.TenantID,
		EjecutorID: ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &AsignarRolOutput{}
	out.Body = presentation.NewApiResponse(dto.AsignarRolResponse{
		UsuarioID: resp.UsuarioID, RolID: resp.RolID,
		TenantID: resp.TenantID, AsignadoEn: resp.AsignadoEn,
	})
	return out, nil
}

// ── Revocar Rol de Usuario ─────────────────────────────────────────────────────

type RevocarRolInput struct {
	UsuarioID string `path:"usuarioID" doc:"ID del usuario"`
	RolID     string `path:"rolID"     doc:"ID del rol a revocar"`
}

type RevocarRolOutput struct {
	Body presentation.ApiResponse[dto.RevocarRolResponse]
}

type RevocarRolHandler struct {
	facade facades.RbacFacade
}

func NewRevocarRolHandler(facade facades.RbacFacade) *RevocarRolHandler {
	return &RevocarRolHandler{facade: facade}
}

func (h *RevocarRolHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "revocar-rol",
		Method:      http.MethodDelete,
		Path:        "/api/v1/usuarios/{usuarioID}/roles/{rolID}",
		Summary:     "Revocar rol de usuario",
		Description: "Revoca un rol asignado a un usuario.",
		Tags:        []string{"Roles"},
	}, h.handle)
}

func (h *RevocarRolHandler) handle(ctx context.Context, input *RevocarRolInput) (*RevocarRolOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	tenantID := middleware.GetTenantIDFromCtx(ctx)

	resp, err := h.facade.RevocarRol(ctx, facades.ComandoRevocarRol{
		UsuarioID:  input.UsuarioID,
		RolID:      input.RolID,
		TenantID:   tenantID,
		EjecutorID: ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &RevocarRolOutput{}
	out.Body = presentation.NewApiResponse(dto.RevocarRolResponse{
		UsuarioID: resp.UsuarioID, RolID: resp.RolID,
		TenantID: resp.TenantID, RevocadoEn: resp.RevocadoEn,
	})
	return out, nil
}

// ── Asignar Permiso a Rol ──────────────────────────────────────────────────────

type AsignarPermisoARolInput struct {
	RolID string `path:"rolID" doc:"ID del rol"`
	Body  dto.AsignarPermisoRequest
}

type AsignarPermisoARolOutput struct {
	Body presentation.ApiResponse[dto.AsignarPermisoResponse]
}

type AsignarPermisoARolHandler struct {
	facade facades.RbacFacade
}

func NewAsignarPermisoARolHandler(facade facades.RbacFacade) *AsignarPermisoARolHandler {
	return &AsignarPermisoARolHandler{facade: facade}
}

func (h *AsignarPermisoARolHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "asignar-permiso-a-rol",
		Method:        http.MethodPost,
		Path:          "/api/v1/roles/{rolID}/permisos",
		Summary:       "Asignar permiso a rol",
		Description:   "Asigna un permiso a un rol específico.",
		Tags:          []string{"Roles"},
		DefaultStatus: http.StatusCreated,
	}, h.handle)
}

func (h *AsignarPermisoARolHandler) handle(ctx context.Context, input *AsignarPermisoARolInput) (*AsignarPermisoARolOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	tenantID := middleware.GetTenantIDFromCtx(ctx)

	resp, err := h.facade.AsignarPermisoARol(ctx, facades.ComandoAsignarPermisoARol{
		RolID:         input.RolID,
		PermisoCodigo: input.Body.PermisoCodigo,
		TenantID:      tenantID,
		EjecutorID:    ejecutorID,
		AsignadoPor:   ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &AsignarPermisoARolOutput{}
	out.Body = presentation.NewApiResponse(dto.AsignarPermisoResponse{
		RolID: resp.RolID, PermisoCodigo: resp.PermisoCodigo,
		AsignadoEn: resp.AsignadoEn,
	})
	return out, nil
}

// ── Revocar Permiso de Rol ─────────────────────────────────────────────────────

type RevocarPermisoDeRolInput struct {
	RolID         string `path:"rolID"         doc:"ID del rol"`
	PermisoCodigo string `path:"codigo"        doc:"Código del permiso a revocar"`
}

type RevocarPermisoDeRolOutput struct {
	Body presentation.ApiResponse[dto.RevocarPermisoResponse]
}

type RevocarPermisoDeRolHandler struct {
	facade facades.RbacFacade
}

func NewRevocarPermisoDeRolHandler(facade facades.RbacFacade) *RevocarPermisoDeRolHandler {
	return &RevocarPermisoDeRolHandler{facade: facade}
}

func (h *RevocarPermisoDeRolHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "revocar-permiso-de-rol",
		Method:      http.MethodDelete,
		Path:        "/api/v1/roles/{rolID}/permisos/{codigo}",
		Summary:     "Revocar permiso de rol",
		Description: "Revoca un permiso previamente asignado a un rol.",
		Tags:        []string{"Roles"},
	}, h.handle)
}

func (h *RevocarPermisoDeRolHandler) handle(ctx context.Context, input *RevocarPermisoDeRolInput) (*RevocarPermisoDeRolOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	tenantID := middleware.GetTenantIDFromCtx(ctx)

	resp, err := h.facade.RevocarPermisoDeRol(ctx, facades.ComandoRevocarPermisoDeRol{
		RolID:         input.RolID,
		PermisoCodigo: input.PermisoCodigo,
		TenantID:      tenantID,
		EjecutorID:    ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &RevocarPermisoDeRolOutput{}
	out.Body = presentation.NewApiResponse(dto.RevocarPermisoResponse{
		RolID: resp.RolID, PermisoCodigo: resp.PermisoCodigo,
		RevocadoEn: resp.RevocadoEn,
	})
	return out, nil
}

// ── Listar Permisos ────────────────────────────────────────────────────────────
type ListarPermisosOutput struct {
	Body presentation.ApiResponse[dto.ListarPermisosResponse]
}

type ListarPermisosHandler struct {
	facade facades.RbacFacade
}

func NewListarPermisosHandler(facade facades.RbacFacade) *ListarPermisosHandler {
	return &ListarPermisosHandler{facade: facade}
}

func (h *ListarPermisosHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "listar-permisos",
		Method:      http.MethodGet,
		Path:        "/api/v1/permisos",
		Summary:     "Listar permisos",
		Description: "Lista todos los permisos disponibles en el sistema.",
		Tags:        []string{"Roles"},
	}, h.handle)
}

func (h *ListarPermisosHandler) handle(ctx context.Context, input *struct{}) (*ListarPermisosOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	tenantID := middleware.GetTenantIDFromCtx(ctx)

	resp, err := h.facade.ListarPermisos(ctx, ejecutorID, tenantID)
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	items := make([]dto.PermisoItem, len(resp.Permisos))
	for i, p := range resp.Permisos {
		items[i] = dto.PermisoItem{
			ID:          p.ID,
			Codigo:      p.Codigo,
			Nombre:      p.Nombre,
			Descripcion: p.Descripcion,
			Modulo:      p.Modulo,
		}
	}
	out := &ListarPermisosOutput{}
	out.Body = presentation.NewApiResponse(dto.ListarPermisosResponse{
		Permisos: items,
		Total:    resp.Total,
	})
	return out, nil
}
