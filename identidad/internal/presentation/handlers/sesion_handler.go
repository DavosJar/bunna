package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/dto"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/facades"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/middleware"
	shared_domain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	presentation "github.com/davosjar/bunna/services/identidad/internal/shared/presentation"
)

// ── Listar Sesiones ────────────────────────────────────────────────────────────

type ListarSesionesInput struct {
	Pagina       int `query:"pagina"       doc:"Número de página (1-based)" example:"1"`
	TamanoPagina int `query:"tamano"       doc:"Elementos por página"       example:"20"`
}

type ListarSesionesOutput struct {
	Body presentation.ApiResponse[dto.ListarSesionesResponse]
}

type ListarSesionesHandler struct {
	facade facades.SesionFacade
}

func NewListarSesionesHandler(facade facades.SesionFacade) *ListarSesionesHandler {
	return &ListarSesionesHandler{facade: facade}
}

func (h *ListarSesionesHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "listar-sesiones",
		Method:      http.MethodGet,
		Path:        "/api/v1/sesiones",
		Summary:     "Listar sesiones",
		Description: "Lista las sesiones activas del sistema con paginación.",
		Tags:        []string{"Sesiones"},
	}, h.handle)
}

func (h *ListarSesionesHandler) handle(ctx context.Context, input *ListarSesionesInput) (*ListarSesionesOutput, error) {
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

	resp, err := h.facade.ListarSesiones(ctx, facades.ComandoListarSesiones{
		Paginacion: shared_domain.Paginacion{Pagina: pagina, TamanoPagina: tamano},
		EjecutorID: ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	items := make([]dto.SesionItem, len(resp.Sesiones))
	for i, s := range resp.Sesiones {
		items[i] = dto.SesionItem{
			ID: s.ID, UsuarioID: s.UsuarioID, IPOrigen: s.IPOrigen,
			Estado: s.Estado, UltimaActividad: s.UltimaActividad.Format(time.RFC3339),
		}
	}

	out := &ListarSesionesOutput{}
	out.Body = presentation.NewApiResponse(dto.ListarSesionesResponse{
		Sesiones: items, Total: resp.Total, Pagina: resp.Pagina,
	})
	return out, nil
}

// ── Forzar Cierre de Sesión ────────────────────────────────────────────────────

type ForzarCierreSesionInput struct {
	SesionID string `path:"sesionID" doc:"ID de la sesión a cerrar"`
}

type ForzarCierreSesionOutput struct {
	Body presentation.ApiResponse[dto.ForzarCierreSesionResponse]
}

type ForzarCierreSesionHandler struct {
	facade facades.SesionFacade
}

func NewForzarCierreSesionHandler(facade facades.SesionFacade) *ForzarCierreSesionHandler {
	return &ForzarCierreSesionHandler{facade: facade}
}

func (h *ForzarCierreSesionHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "forzar-cierre-sesion",
		Method:      http.MethodDelete,
		Path:        "/api/v1/sesiones/{sesionID}",
		Summary:     "Forzar cierre de sesión",
		Description: " fuerza el cierre de una sesión específica (requiere permisos administrativos).",
		Tags:        []string{"Sesiones"},
	}, h.handle)
}

func (h *ForzarCierreSesionHandler) handle(ctx context.Context, input *ForzarCierreSesionInput) (*ForzarCierreSesionOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	resp, err := h.facade.ForzarCierreSesion(ctx, facades.ComandoForzarCierreSesion{
		SesionID:   input.SesionID,
		EjecutorID: ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &ForzarCierreSesionOutput{}
	out.Body = presentation.NewApiResponse(dto.ForzarCierreSesionResponse{
		SesionID: resp.SesionID, Estado: resp.Estado, RevocadoEn: resp.RevocadoEn,
	})
	return out, nil
}
