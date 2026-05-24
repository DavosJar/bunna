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

// ── Crear Usuario ──────────────────────────────────────────────────────────────

type CrearUsuarioInput struct {
	Body dto.CrearUsuarioRequest
}

type CrearUsuarioOutput struct {
	Body presentation.ApiResponse[dto.CrearUsuarioResponse]
}

type CrearUsuarioHandler struct {
	facade facades.UsuarioFacade
}

func NewCrearUsuarioHandler(facade facades.UsuarioFacade) *CrearUsuarioHandler {
	return &CrearUsuarioHandler{facade: facade}
}

func (h *CrearUsuarioHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "crear-usuario",
		Method:        http.MethodPost,
		Path:          "/api/v1/usuarios",
		Summary:       "Crear usuario",
		Description:   "Crea un nuevo usuario en el sistema (requiere permisos de administración).",
		Tags:          []string{"Usuarios"},
		DefaultStatus: http.StatusCreated,
	}, h.handle)
}

func (h *CrearUsuarioHandler) handle(ctx context.Context, input *CrearUsuarioInput) (*CrearUsuarioOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	resp, err := h.facade.CrearUsuario(ctx, facades.ComandoCrearUsuario{
		Correo:     input.Body.Correo,
		Nombre:     input.Body.Nombre,
		Apellido:   input.Body.Apellido,
		Password:   input.Body.Password,
		EjecutorID: ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &CrearUsuarioOutput{}
	out.Body = presentation.NewApiResponse(dto.CrearUsuarioResponse{
		ID:       resp.ID,
		Correo:   resp.Correo,
		Nombre:   resp.Nombre,
		Apellido: resp.Apellido,
		Activo:   resp.Activo,
		CreadoEn: resp.CreadoEn,
	})
	return out, nil
}

// ── Listar Usuarios ────────────────────────────────────────────────────────────

type ListarUsuariosInput struct {
	Pagina       int    `query:"pagina"       doc:"Número de página (1-based)" example:"1"`
	TamanoPagina int    `query:"tamano"       doc:"Elementos por página"       example:"20"`
	Correo       string `query:"correo"       doc:"Filtrar por correo"         example:""`
	Estado       string `query:"estado"       doc:"Filtrar por estado"         example:"ACTIVO"`
}

type ListarUsuariosOutput struct {
	Body presentation.ApiResponse[dto.ListarUsuariosResponse]
}

type ListarUsuariosHandler struct {
	facade facades.UsuarioFacade
}

func NewListarUsuariosHandler(facade facades.UsuarioFacade) *ListarUsuariosHandler {
	return &ListarUsuariosHandler{facade: facade}
}

func (h *ListarUsuariosHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "listar-usuarios",
		Method:      http.MethodGet,
		Path:        "/api/v1/usuarios",
		Summary:     "Listar usuarios",
		Description: "Lista usuarios del sistema con filtros y paginación.",
		Tags:        []string{"Usuarios"},
	}, h.handle)
}

func (h *ListarUsuariosHandler) handle(ctx context.Context, input *ListarUsuariosInput) (*ListarUsuariosOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	pagina := input.Pagina
	if pagina < 1 {
		pagina = 1
	}
	tamano := input.TamanoPagina
	if tamano < 1 || tamano > 100 {
		tamano = 20
	}

	var filtros []shared_domain.CriterioFiltro
	if input.Correo != "" {
		filtros = append(filtros, shared_domain.CriterioFiltro{Campo: "correo", Operador: "=", Valor: input.Correo})
	}
	if input.Estado != "" {
		filtros = append(filtros, shared_domain.CriterioFiltro{Campo: "estado", Operador: "=", Valor: input.Estado})
	}

	resp, err := h.facade.ListarUsuarios(ctx, facades.ComandoListarUsuarios{
		Filtros:    filtros,
		Paginacion: shared_domain.Paginacion{Pagina: pagina, TamanoPagina: tamano},
		EjecutorID: ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	items := make([]dto.UsuarioItem, len(resp.Usuarios))
	for i, u := range resp.Usuarios {
		items[i] = dto.UsuarioItem{
			ID: u.ID, Correo: u.Correo, Nombre: u.Nombre,
			Apellido: u.Apellido, Estado: u.Estado, CreadoEn: u.CreadoEn,
		}
	}

	out := &ListarUsuariosOutput{}
	out.Body = presentation.NewApiResponse(dto.ListarUsuariosResponse{
		Usuarios: items,
		Total:    resp.Total,
		Pagina:   resp.Pagina,
		Tamano:   resp.Tamano,
	})
	return out, nil
}

// ── Modificar Usuario ─────────────────────────────────────────────────────────

type ModificarUsuarioInput struct {
	UsuarioID string `path:"usuarioID" doc:"ID del usuario a modificar"`
	Body      dto.ModificarUsuarioRequest
}

type ModificarUsuarioOutput struct {
	Body presentation.ApiResponse[dto.ModificarUsuarioResponse]
}

type ModificarUsuarioHandler struct {
	facade facades.UsuarioFacade
}

func NewModificarUsuarioHandler(facade facades.UsuarioFacade) *ModificarUsuarioHandler {
	return &ModificarUsuarioHandler{facade: facade}
}

func (h *ModificarUsuarioHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "modificar-usuario",
		Method:      http.MethodPut,
		Path:        "/api/v1/usuarios/{usuarioID}",
		Summary:     "Modificar usuario",
		Description: "Modifica los datos de un usuario existente.",
		Tags:        []string{"Usuarios"},
	}, h.handle)
}

func (h *ModificarUsuarioHandler) handle(ctx context.Context, input *ModificarUsuarioInput) (*ModificarUsuarioOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	resp, err := h.facade.ModificarUsuario(ctx, facades.ComandoModificarUsuario{
		UsuarioID:  input.UsuarioID,
		Nombre:     input.Body.Nombre,
		Apellido:   input.Body.Apellido,
		EjecutorID: ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &ModificarUsuarioOutput{}
	out.Body = presentation.NewApiResponse(dto.ModificarUsuarioResponse{
		ID: resp.ID, Correo: resp.Correo,
		Nombre: resp.Nombre, Apellido: resp.Apellido,
		ModificadoEn: resp.ModificadoEn,
	})
	return out, nil
}

// ── Dar de Baja Usuario ────────────────────────────────────────────────────────

type DarDeBajaUsuarioInput struct {
	UsuarioID string `path:"usuarioID" doc:"ID del usuario a dar de baja"`
	Body      dto.DarDeBajaUsuarioRequest
}

type DarDeBajaUsuarioOutput struct {
	Body presentation.ApiResponse[dto.DarDeBajaUsuarioResponse]
}

type DarDeBajaUsuarioHandler struct {
	facade facades.UsuarioFacade
}

func NewDarDeBajaUsuarioHandler(facade facades.UsuarioFacade) *DarDeBajaUsuarioHandler {
	return &DarDeBajaUsuarioHandler{facade: facade}
}

func (h *DarDeBajaUsuarioHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "dar-de-baja-usuario",
		Method:      http.MethodDelete,
		Path:        "/api/v1/usuarios/{usuarioID}",
		Summary:     "Dar de baja usuario",
		Description: "Desactiva un usuario del sistema (baja lógica).",
		Tags:        []string{"Usuarios"},
	}, h.handle)
}

func (h *DarDeBajaUsuarioHandler) handle(ctx context.Context, input *DarDeBajaUsuarioInput) (*DarDeBajaUsuarioOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	motivo := ""
	if input.Body.Motivo != "" {
		motivo = input.Body.Motivo
	}

	resp, err := h.facade.DarDeBajaUsuario(ctx, facades.ComandoDarDeBajaUsuario{
		UsuarioID:  input.UsuarioID,
		Motivo:     motivo,
		EjecutorID: ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &DarDeBajaUsuarioOutput{}
	out.Body = presentation.NewApiResponse(dto.DarDeBajaUsuarioResponse{
		UsuarioID: resp.UsuarioID,
		Estado:    resp.Estado,
		BajaEn:    resp.BajaEn,
	})
	return out, nil
}

// ── Expulsar Usuario ───────────────────────────────────────────────────────────

type ExpulsarUsuarioInput struct {
	UsuarioID string `path:"usuarioID" doc:"ID del usuario a expulsar"`
}

type ExpulsarUsuarioOutput struct {
	Body presentation.ApiResponse[dto.ExpulsarUsuarioResponse]
}

type ExpulsarUsuarioHandler struct {
	facade facades.UsuarioFacade
}

func NewExpulsarUsuarioHandler(facade facades.UsuarioFacade) *ExpulsarUsuarioHandler {
	return &ExpulsarUsuarioHandler{facade: facade}
}

func (h *ExpulsarUsuarioHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "expulsar-usuario",
		Method:      http.MethodPost,
		Path:        "/api/v1/usuarios/{usuarioID}/expulsar",
		Summary:     "Expulsar usuario",
		Description: "Expulsa a un usuario del sistema, desactivándolo e invalidando todas sus sesiones.",
		Tags:        []string{"Usuarios"},
	}, h.handle)
}

func (h *ExpulsarUsuarioHandler) handle(ctx context.Context, input *ExpulsarUsuarioInput) (*ExpulsarUsuarioOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	resp, err := h.facade.ExpulsarUsuario(ctx, facades.ComandoExpulsarUsuario{
		UsuarioID:  input.UsuarioID,
		EjecutorID: ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &ExpulsarUsuarioOutput{}
	out.Body = presentation.NewApiResponse(dto.ExpulsarUsuarioResponse{
		UsuarioID:         resp.UsuarioID,
		Estado:            resp.Estado,
		SesionesRevocadas: resp.SesionesRevocadas,
		ExpulsadoEn:       resp.ExpulsadoEn,
	})
	return out, nil
}

// ── Ver Mi Perfil ──────────────────────────────────────────────────────────────

type VerMiPerfilOutput struct {
	Body presentation.ApiResponse[dto.VerMiPerfilResponse]
}

type VerMiPerfilHandler struct {
	facade facades.UsuarioFacade
}

func NewVerMiPerfilHandler(facade facades.UsuarioFacade) *VerMiPerfilHandler {
	return &VerMiPerfilHandler{facade: facade}
}

func (h *VerMiPerfilHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "ver-mi-perfil",
		Method:      http.MethodGet,
		Path:        "/api/v1/mi-perfil",
		Summary:     "Ver mi perfil",
		Description: "Obtiene los datos del perfil del usuario autenticado.",
		Tags:        []string{"Mi Perfil"},
	}, h.handle)
}

func (h *VerMiPerfilHandler) handle(ctx context.Context, input *struct{}) (*VerMiPerfilOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	resp, err := h.facade.VerMiPerfil(ctx, facades.ComandoVerMiPerfil{
		EjecutorID: ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &VerMiPerfilOutput{}
	out.Body = presentation.NewApiResponse(dto.VerMiPerfilResponse{
		ID: resp.ID, Correo: resp.Correo, Nombre: resp.Nombre,
		Apellido: resp.Apellido, Telefono: resp.Telefono,
		Estado: resp.Estado, CreadoEn: resp.CreadoEn,
	})
	return out, nil
}

// ── Modificar Mi Perfil ────────────────────────────────────────────────────────

type ModificarMiPerfilInput struct {
	Body dto.ModificarMiPerfilRequest
}

type ModificarMiPerfilOutput struct {
	Body presentation.ApiResponse[dto.ModificarMiPerfilResponse]
}

type ModificarMiPerfilHandler struct {
	facade facades.UsuarioFacade
}

func NewModificarMiPerfilHandler(facade facades.UsuarioFacade) *ModificarMiPerfilHandler {
	return &ModificarMiPerfilHandler{facade: facade}
}

func (h *ModificarMiPerfilHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "modificar-mi-perfil",
		Method:      http.MethodPut,
		Path:        "/api/v1/mi-perfil",
		Summary:     "Modificar mi perfil",
		Description: "Actualiza los datos del perfil del usuario autenticado.",
		Tags:        []string{"Mi Perfil"},
	}, h.handle)
}

func (h *ModificarMiPerfilHandler) handle(ctx context.Context, input *ModificarMiPerfilInput) (*ModificarMiPerfilOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	resp, err := h.facade.ModificarMiPerfil(ctx, facades.ComandoModificarMiPerfil{
		EjecutorID: ejecutorID,
		Nombre:     input.Body.Nombre,
		Apellido:   input.Body.Apellido,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &ModificarMiPerfilOutput{}
	out.Body = presentation.NewApiResponse(dto.ModificarMiPerfilResponse{
		ID: resp.ID, Correo: resp.Correo,
		Nombre: resp.Nombre, Apellido: resp.Apellido,
		ModificadoEn: resp.ModificadoEn,
	})
	return out, nil
}
