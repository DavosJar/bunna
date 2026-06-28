package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/dto"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/facades"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/middleware"
	presentation "github.com/davosjar/bunna/services/identidad/internal/shared/presentation"
)

type ObtenerInvitacionInput struct {
	Token string `path:"token" doc:"Token de invitación"`
}

type ObtenerInvitacionOutput struct {
	Body presentation.ApiResponse[dto.ObtenerInvitacionResponse]
}

type ObtenerInvitacionHandler struct {
	facade facades.InvitacionFacade
}

func NewObtenerInvitacionHandler(facade facades.InvitacionFacade) *ObtenerInvitacionHandler {
	return &ObtenerInvitacionHandler{facade: facade}
}

func (h *ObtenerInvitacionHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "obtener-invitacion",
		Method:      http.MethodGet,
		Path:        "/api/v1/invitaciones/{token}",
		Summary:     "Obtener información de invitación",
		Description: "Endpoint público. Obtiene los datos de una invitación por su token, sin requerir autenticación.",
		Tags:        []string{"Invitaciones"},
	}, h.handle)
}

func (h *ObtenerInvitacionHandler) handle(ctx context.Context, input *ObtenerInvitacionInput) (*ObtenerInvitacionOutput, error) {
	// SIN JWT — endpoint público
	resp, err := h.facade.ObtenerInvitacion(ctx, input.Token)
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &ObtenerInvitacionOutput{}
	out.Body = presentation.NewApiResponse(dto.ObtenerInvitacionResponse{
		ID:           resp.ID,
		TenantID:     resp.TenantID,
		TenantNombre: resp.TenantNombre,
		RolID:        resp.RolID,
		RolNombre:    resp.RolNombre,
		Email:        resp.Email,
		Estado:       resp.Estado,
		Expiracion:   resp.Expiracion,
	})
	return out, nil
}

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
	Token string `path:"token" doc:"Token de invitación"`
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
		Path:          "/api/v1/invitaciones/{token}/aceptar",
		Summary:       "Aceptar invitación",
		Description:   "Endpoint público. Acepta una invitación usando solo el token, sin requerir JWT.",
		Tags:          []string{"Invitaciones"},
		DefaultStatus: http.StatusOK,
	}, h.handle)
}

func (h *AceptarInvitacionHandler) handle(ctx context.Context, input *AceptarInvitacionInput) (*AceptarInvitacionOutput, error) {
	// SIN JWT — endpoint público. El token de invitación es la autorización.
	resp, err := h.facade.AceptarInvitacion(ctx, facades.ComandoAceptarInvitacion{
		Token: input.Token,
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

// ─────────────────────────────────────────────────────────────────────────────
// GET /api/v1/invitaciones (admin, con JWT) — Listar invitaciones del tenant
// ─────────────────────────────────────────────────────────────────────────────

type ListarInvitacionesInput struct {
	Pagina       int    `query:"pagina" doc:"Número de página (1-based)" example:"1"`
	TamanoPagina int    `query:"tamano_pagina" doc:"Tamaño de página" example:"20"`
	Estado       string `query:"estado" doc:"Filtrar por estado: pendiente, aceptada, expirada (vacío = todas)"`
}

type ListarInvitacionesOutput struct {
	Body presentation.ApiResponse[dto.ListarInvitacionesResponse]
}

type ListarInvitacionesHandler struct {
	facade facades.InvitacionFacade
}

func NewListarInvitacionesHandler(facade facades.InvitacionFacade) *ListarInvitacionesHandler {
	return &ListarInvitacionesHandler{facade: facade}
}

func (h *ListarInvitacionesHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "listar-invitaciones",
		Method:      http.MethodGet,
		Path:        "/api/v1/invitaciones",
		Summary:     "Listar invitaciones del tenant",
		Description: "Lista las invitaciones enviadas desde el tenant activo, con filtro opcional por estado.",
		Tags:        []string{"Invitaciones"},
	}, h.handle)
}

func (h *ListarInvitacionesHandler) handle(ctx context.Context, input *ListarInvitacionesInput) (*ListarInvitacionesOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}
	tenantID := middleware.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		return nil, huma.Error400BadRequest("tenant no encontrado en el token")
	}

	pagina := input.Pagina
	if pagina < 1 {
		pagina = 1
	}
	tamanoPagina := input.TamanoPagina
	if tamanoPagina < 1 || tamanoPagina > 100 {
		tamanoPagina = 20
	}

	resp, err := h.facade.ListarInvitaciones(ctx, facades.ComandoListarInvitaciones{
		TenantID:     tenantID,
		Pagina:       pagina,
		TamanoPagina: tamanoPagina,
		Estado:       input.Estado,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	items := make([]dto.InvitacionItem, len(resp.Invitaciones))
	for i, inv := range resp.Invitaciones {
		items[i] = dto.InvitacionItem{
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

	out := &ListarInvitacionesOutput{}
	out.Body = presentation.NewApiResponse(dto.ListarInvitacionesResponse{
		Invitaciones: items,
		Total:        resp.Total,
	})
	return out, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// DELETE /api/v1/invitaciones/{id} (admin, con JWT) — Eliminar invitación
// ─────────────────────────────────────────────────────────────────────────────

type EliminarInvitacionInput struct {
	ID string `path:"id" doc:"ID de la invitación a eliminar" example:"01926b1e-dead-beef-cafe-000000000001"`
}

type EliminarInvitacionOutput struct {
	Body presentation.ApiResponse[dto.EliminarInvitacionResponse]
}

type EliminarInvitacionHandler struct {
	facade facades.InvitacionFacade
}

func NewEliminarInvitacionHandler(facade facades.InvitacionFacade) *EliminarInvitacionHandler {
	return &EliminarInvitacionHandler{facade: facade}
}

func (h *EliminarInvitacionHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "eliminar-invitacion",
		Method:        http.MethodDelete,
		Path:          "/api/v1/invitaciones/{id}",
		Summary:       "Eliminar invitación",
		Description:   "Elimina definitivamente una invitación pendiente. Solo para invitaciones del mismo tenant que no hayan sido aceptadas.",
		Tags:          []string{"Invitaciones"},
		DefaultStatus: http.StatusOK,
	}, h.handle)
}

func (h *EliminarInvitacionHandler) handle(ctx context.Context, input *EliminarInvitacionInput) (*EliminarInvitacionOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}
	tenantID := middleware.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		return nil, huma.Error400BadRequest("tenant no encontrado en el token")
	}

	resp, err := h.facade.EliminarInvitacion(ctx, facades.ComandoEliminarInvitacion{
		InvitacionID: input.ID,
		TenantID:     tenantID,
		EjecutorID:   ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &EliminarInvitacionOutput{}
	out.Body = presentation.NewApiResponse(dto.EliminarInvitacionResponse{
		Mensaje: resp.Mensaje,
	})
	return out, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// POST /api/v1/invitaciones/{id}/reenviar (admin, con JWT) — Reenviar invitación
// ─────────────────────────────────────────────────────────────────────────────

type ReenviarInvitacionInput struct {
	ID string `path:"token" doc:"ID de la invitación a reenviar" example:"01926b1e-dead-beef-cafe-000000000001"`
}

type ReenviarInvitacionOutput struct {
	Body presentation.ApiResponse[dto.ReenviarInvitacionResponse]
}

type ReenviarInvitacionHandler struct {
	facade facades.InvitacionFacade
}

func NewReenviarInvitacionHandler(facade facades.InvitacionFacade) *ReenviarInvitacionHandler {
	return &ReenviarInvitacionHandler{facade: facade}
}

func (h *ReenviarInvitacionHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "reenviar-invitacion",
		Method:        http.MethodPost,
		Path:          "/api/v1/invitaciones/{token}/reenviar",
		Summary:       "Reenviar invitación",
		Description:   "Reenvía el email de invitación generando un nuevo token. Solo para invitaciones pendientes del mismo tenant.",
		Tags:          []string{"Invitaciones"},
		DefaultStatus: http.StatusOK,
	}, h.handle)
}

func (h *ReenviarInvitacionHandler) handle(ctx context.Context, input *ReenviarInvitacionInput) (*ReenviarInvitacionOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}
	tenantID := middleware.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		return nil, huma.Error400BadRequest("tenant no encontrado en el token")
	}

	resp, err := h.facade.ReenviarInvitacion(ctx, facades.ComandoReenviarInvitacion{
		InvitacionID: input.ID,
		TenantID:     tenantID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &ReenviarInvitacionOutput{}
	out.Body = presentation.NewApiResponse(dto.ReenviarInvitacionResponse{
		Mensaje: resp.Mensaje,
	})
	return out, nil
}
