package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/dto"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/facades"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/middleware"
	presentation "github.com/davosjar/bunna/services/identidad/shared/presentation"
)

type CrearInvitacionInput struct {
	Body dto.CrearInvitacionRequest
}

type CrearInvitacionOutput struct {
	Body presentation.ApiResponse[dto.CrearInvitacionResponse]
}

type CrearInvitacionHandler struct {
	facade facades.InvitacionFacade
}

func NewCrearInvitacionHandler(facade facades.InvitacionFacade) *CrearInvitacionHandler {
	return &CrearInvitacionHandler{facade: facade}
}

func (h *CrearInvitacionHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "crear-invitacion",
		Method:        http.MethodPost,
		Path:          "/api/v1/invitaciones",
		Summary:       "Crear invitación",
		Description:   "Crea una invitación para que un usuario se una al tenant del ejecutor.",
		Tags:          []string{"Invitaciones"},
		DefaultStatus: http.StatusCreated,
	}, h.handle)
}

func (h *CrearInvitacionHandler) handle(ctx context.Context, input *CrearInvitacionInput) (*CrearInvitacionOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}
	tenantID := middleware.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		return nil, huma.Error400BadRequest("tenant no encontrado en el token")
	}

	resp, err := h.facade.CrearInvitacion(ctx, facades.ComandoCrearInvitacion{
		TenantID:  tenantID,
		RolID:     input.Body.RolID,
		Correo:    input.Body.Correo,
		CreadoPor: ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &CrearInvitacionOutput{}
	out.Body = presentation.NewApiResponse(dto.CrearInvitacionResponse{
		Mensaje: resp.Mensaje,
	})
	return out, nil
}

type AceptarInvitacionInput struct {
	Body dto.AceptarInvitacionRequest
}

type AceptarInvitacionOutput struct {
	Body presentation.ApiResponse[dto.AceptarInvitacionResponse]
}

type AceptarInvitacionHandler struct {
	facade facades.InvitacionFacade
}

func NewAceptarInvitacionHandler(facade facades.InvitacionFacade) *AceptarInvitacionHandler {
	return &AceptarInvitacionHandler{facade: facade}
}

func (h *AceptarInvitacionHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "aceptar-invitacion",
		Method:        http.MethodPost,
		Path:          "/api/v1/invitaciones/aceptar",
		Summary:       "Aceptar invitación",
		Description:   "Acepta una invitación para unirse a un tenant con un rol específico.",
		Tags:          []string{"Invitaciones"},
		DefaultStatus: http.StatusOK,
	}, h.handle)
}

func (h *AceptarInvitacionHandler) handle(ctx context.Context, input *AceptarInvitacionInput) (*AceptarInvitacionOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	resp, err := h.facade.AceptarInvitacion(ctx, facades.ComandoAceptarInvitacion{
		Token:     input.Body.Token,
		UsuarioID: ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &AceptarInvitacionOutput{}
	out.Body = presentation.NewApiResponse(dto.AceptarInvitacionResponse{
		TenantID: resp.TenantID,
		RolID:    resp.RolID,
	})
	return out, nil
}
