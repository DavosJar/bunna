package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/dto"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/facades"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/middleware"
	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	presentation "github.com/davosjar/bunna/services/identidad/internal/shared/presentation"
)

// ── Cambiar Mi Contraseña ──────────────────────────────────────────────────────

type CambiarMiPasswordInput struct {
	Body dto.CambiarMiPasswordRequest
}

type CambiarMiPasswordOutput struct {
	Body presentation.ApiResponse[dto.CambiarMiPasswordResponse]
}

type CambiarMiPasswordHandler struct {
	facade facades.SeguridadFacade
}

func NewCambiarMiPasswordHandler(facade facades.SeguridadFacade) *CambiarMiPasswordHandler {
	return &CambiarMiPasswordHandler{facade: facade}
}

func (h *CambiarMiPasswordHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "cambiar-mi-password",
		Method:      http.MethodPut,
		Path:        "/api/v1/mi-password",
		Summary:     "Cambiar mi contraseña",
		Description: "Cambia la contraseña del usuario autenticado.",
		Tags:        []string{"Mi Perfil"},
	}, h.handle)
}

func (h *CambiarMiPasswordHandler) handle(ctx context.Context, input *CambiarMiPasswordInput) (*CambiarMiPasswordOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	resp, err := h.facade.CambiarMiPassword(ctx, facades.ComandoCambiarMiPassword{
		EjecutorID:     ejecutorID,
		PasswordActual: input.Body.PasswordActual,
		NuevaPassword:  input.Body.NuevaPassword,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &CambiarMiPasswordOutput{}
	out.Body = presentation.NewApiResponse(dto.CambiarMiPasswordResponse{
		ModificadoEn: resp.ModificadoEn,
	})
	return out, nil
}

// ── Resetear Contraseña ────────────────────────────────────────────────────────

type ResetearPasswordInput struct {
	UsuarioID string `path:"usuarioID" doc:"ID del usuario"`
	Body      dto.ResetearPasswordRequest
}

type ResetearPasswordOutput struct {
	Body presentation.ApiResponse[dto.ResetearPasswordResponse]
}

type ResetearPasswordHandler struct {
	facade facades.SeguridadFacade
}

func NewResetearPasswordHandler(facade facades.SeguridadFacade) *ResetearPasswordHandler {
	return &ResetearPasswordHandler{facade: facade}
}

func (h *ResetearPasswordHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "resetear-password",
		Method:      http.MethodPost,
		Path:        "/api/v1/usuarios/{usuarioID}/reset-password",
		Summary:     "Resetear contraseña",
		Description: "Resetea la contraseña de un usuario (requiere permisos administrativos).",
		Tags:        []string{"Seguridad"},
	}, h.handle)
}

func (h *ResetearPasswordHandler) handle(ctx context.Context, input *ResetearPasswordInput) (*ResetearPasswordOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	resp, err := h.facade.ResetearPassword(ctx, facades.ComandoResetearPassword{
		UsuarioID:     input.UsuarioID,
		NuevaPassword: input.Body.NuevaPassword,
		EjecutorID:    ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &ResetearPasswordOutput{}
	out.Body = presentation.NewApiResponse(dto.ResetearPasswordResponse{
		UsuarioID:    resp.UsuarioID,
		ModificadoEn: resp.ModificadoEn,
	})
	return out, nil
}

// ── Desbloquear Cuenta ─────────────────────────────────────────────────────────

type DesbloquearCuentaInput struct {
	UsuarioID string `path:"usuarioID" doc:"ID del usuario a desbloquear"`
}

type DesbloquearCuentaOutput struct {
	Body presentation.ApiResponse[dto.DesbloquearCuentaResponse]
}

type DesbloquearCuentaHandler struct {
	facade facades.SeguridadFacade
}

func NewDesbloquearCuentaHandler(facade facades.SeguridadFacade) *DesbloquearCuentaHandler {
	return &DesbloquearCuentaHandler{facade: facade}
}

func (h *DesbloquearCuentaHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "desbloquear-cuenta",
		Method:      http.MethodPost,
		Path:        "/api/v1/usuarios/{usuarioID}/unlock",
		Summary:     "Desbloquear cuenta",
		Description: "Desbloquea la cuenta de un usuario bloqueada por intentos fallidos.",
		Tags:        []string{"Seguridad"},
	}, h.handle)
}

func (h *DesbloquearCuentaHandler) handle(ctx context.Context, input *DesbloquearCuentaInput) (*DesbloquearCuentaOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	resp, err := h.facade.DesbloquearCuenta(ctx, facades.ComandoDesbloquearCuenta{
		UsuarioID:  input.UsuarioID,
		EjecutorID: ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &DesbloquearCuentaOutput{}
	out.Body = presentation.NewApiResponse(dto.DesbloquearCuentaResponse{
		UsuarioID:      resp.UsuarioID,
		DesbloqueadoEn: resp.DesbloqueadoEn,
	})
	return out, nil
}

// ── Listar IPs Bloqueadas ──────────────────────────────────────────────────────

type ListarIPsBloqueadasInput struct {
	Pagina       int `query:"pagina"       doc:"Número de página (1-based)" example:"1"`
	TamanoPagina int `query:"tamano"       doc:"Elementos por página"       example:"20"`
}

type ListarIPsBloqueadasOutput struct {
	Body presentation.ApiResponse[dto.ListarIPsBloqueadasResponse]
}

type ListarIPsBloqueadasHandler struct {
	facade facades.SeguridadFacade
}

func NewListarIPsBloqueadasHandler(facade facades.SeguridadFacade) *ListarIPsBloqueadasHandler {
	return &ListarIPsBloqueadasHandler{facade: facade}
}

func (h *ListarIPsBloqueadasHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "listar-ips-bloqueadas",
		Method:      http.MethodGet,
		Path:        "/api/v1/ips-bloqueadas",
		Summary:     "Listar IPs bloqueadas",
		Description: "Lista las direcciones IP bloqueadas temporalmente por exceso de intentos.",
		Tags:        []string{"Seguridad"},
	}, h.handle)
}

func (h *ListarIPsBloqueadasHandler) handle(ctx context.Context, input *ListarIPsBloqueadasInput) (*ListarIPsBloqueadasOutput, error) {
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

	resp, err := h.facade.ListarIPsBloqueadas(ctx, facades.ComandoListarIPsBloqueadas{
		Paginacion: shared_domain.Paginacion{Pagina: pagina, TamanoPagina: tamano},
		EjecutorID: ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	items := make([]dto.IPBloqueadaItem, len(resp.IPs))
	for i, ip := range resp.IPs {
		items[i] = dto.IPBloqueadaItem{
			IP: ip.IP, Intentos: ip.Intentos,
			BloqueadoHasta: ip.BloqueadoHasta,
		}
	}

	out := &ListarIPsBloqueadasOutput{}
	out.Body = presentation.NewApiResponse(dto.ListarIPsBloqueadasResponse{
		IPs: items, Total: resp.Total, Pagina: resp.Pagina,
	})
	return out, nil
}

// ── Desbloquear IP ─────────────────────────────────────────────────────────────

type DesbloquearIPInput struct {
	IP string `path:"ip" doc:"Dirección IP a desbloquear"`
}

type DesbloquearIPOutput struct {
	Body presentation.ApiResponse[dto.DesbloquearIPResponse]
}

type DesbloquearIPHandler struct {
	facade facades.SeguridadFacade
}

func NewDesbloquearIPHandler(facade facades.SeguridadFacade) *DesbloquearIPHandler {
	return &DesbloquearIPHandler{facade: facade}
}

func (h *DesbloquearIPHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "desbloquear-ip",
		Method:      http.MethodDelete,
		Path:        "/api/v1/ips-bloqueadas/{ip}",
		Summary:     "Desbloquear IP",
		Description: "Elimina el bloqueo de una dirección IP.",
		Tags:        []string{"Seguridad"},
	}, h.handle)
}

func (h *DesbloquearIPHandler) handle(ctx context.Context, input *DesbloquearIPInput) (*DesbloquearIPOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	resp, err := h.facade.DesbloquearIP(ctx, facades.ComandoDesbloquearIP{
		IP:         input.IP,
		EjecutorID: ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &DesbloquearIPOutput{}
	out.Body = presentation.NewApiResponse(dto.DesbloquearIPResponse{
		IP: resp.IP, DesbloqueadoEn: resp.DesbloqueadoEn,
	})
	return out, nil
}

// ── Consultar Credenciales ─────────────────────────────────────────────────────

type ConsultarCredencialesInput struct {
	UsuarioID string `path:"usuarioID" doc:"ID del usuario"`
}

type ConsultarCredencialesOutput struct {
	Body presentation.ApiResponse[dto.ConsultarCredencialesResponse]
}

type ConsultarCredencialesHandler struct {
	facade facades.SeguridadFacade
}

func NewConsultarCredencialesHandler(facade facades.SeguridadFacade) *ConsultarCredencialesHandler {
	return &ConsultarCredencialesHandler{facade: facade}
}

func (h *ConsultarCredencialesHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "consultar-credenciales",
		Method:      http.MethodGet,
		Path:        "/api/v1/credenciales/{usuarioID}",
		Summary:     "Consultar credenciales",
		Description: "Obtiene el estado de las credenciales de un usuario (bloqueo, intentos, verificación).",
		Tags:        []string{"Seguridad"},
	}, h.handle)
}

func (h *ConsultarCredencialesHandler) handle(ctx context.Context, input *ConsultarCredencialesInput) (*ConsultarCredencialesOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	resp, err := h.facade.ConsultarCredenciales(ctx, facades.ComandoConsultarCredenciales{
		UsuarioID:  input.UsuarioID,
		EjecutorID: ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &ConsultarCredencialesOutput{}
	out.Body = presentation.NewApiResponse(dto.ConsultarCredencialesResponse{
		UsuarioID:        resp.UsuarioID,
		Activo:           resp.Activo,
		CorreoVerificado: resp.CorreoVerificado,
		IntentosFallidos: resp.IntentosFallidos,
		BloqueadoHasta:   resp.BloqueadoHasta,
	})
	return out, nil
}
