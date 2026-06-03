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

type RefreshInput struct {
	Body dto.RefreshRequest
}

type RefreshOutput struct {
	Body presentation.ApiResponse[dto.RefreshResponse]
}

type RefreshHandler struct {
	facade facades.AuthFacade
}

func NewRefreshHandler(facade facades.AuthFacade) *RefreshHandler {
	return &RefreshHandler{facade: facade}
}

func (h *RefreshHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "refresh-sesion",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/refresh",
		Summary:     "Renovar sesión",
		Description: "Renueva el access token usando el refresh token. Aplica rotación de tokens.",
		Tags:        []string{"Autenticación"},
	}, h.handle)
}

func (h *RefreshHandler) handle(ctx context.Context, input *RefreshInput) (*RefreshOutput, error) {
	resp, err := h.facade.Refresh(ctx, facades.ComandoRefresh{
		RefreshToken: input.Body.RefreshToken,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &RefreshOutput{}
	out.Body = presentation.NewApiResponse(dto.RefreshResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
		TokenType:    resp.TokenType,
		UsuarioID:    resp.UsuarioID,
		TenantID:     resp.TenantID,
		Rol:          resp.Rol,
	})
	return out, nil
}

// ── Logout ─────────────────────────────────────────────────────────────────────

type LogoutOutput struct {
	Body presentation.ApiResponse[dto.LogoutResponse]
}

type LogoutHandler struct {
	facade facades.AuthFacade
}

func NewLogoutHandler(facade facades.AuthFacade) *LogoutHandler {
	return &LogoutHandler{facade: facade}
}

func (h *LogoutHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "cerrar-sesion",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/logout",
		Summary:     "Cerrar sesión",
		Description: "Cierra la sesión actual del usuario autenticado.",
		Tags:        []string{"Autenticación"},
	}, h.handle)
}

func (h *LogoutHandler) handle(ctx context.Context, input *struct{}) (*LogoutOutput, error) {
	usuarioID := middleware.GetUsuarioIDFromCtx(ctx)
	sesionID := middleware.GetSesionIDFromCtx(ctx)
	if usuarioID == "" || sesionID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	resp, err := h.facade.Logout(ctx, facades.ComandoLogout{
		SesionID:  sesionID,
		UsuarioID: usuarioID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &LogoutOutput{}
	out.Body = presentation.NewApiResponse(dto.LogoutResponse{
		SesionesRevocadas: resp.SesionesRevocadas,
	})
	return out, nil
}

// ── LogoutAll ──────────────────────────────────────────────────────────────────

type LogoutAllOutput struct {
	Body presentation.ApiResponse[dto.LogoutResponse]
}

type LogoutAllHandler struct {
	facade facades.AuthFacade
}

func NewLogoutAllHandler(facade facades.AuthFacade) *LogoutAllHandler {
	return &LogoutAllHandler{facade: facade}
}

func (h *LogoutAllHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "cerrar-todas-sesiones",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/logout/all",
		Summary:     "Cerrar todas las sesiones",
		Description: "Cierra todas las sesiones activas del usuario autenticado.",
		Tags:        []string{"Autenticación"},
	}, h.handle)
}

func (h *LogoutAllHandler) handle(ctx context.Context, input *struct{}) (*LogoutAllOutput, error) {
	usuarioID := middleware.GetUsuarioIDFromCtx(ctx)
	if usuarioID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	resp, err := h.facade.LogoutAll(ctx, facades.ComandoLogoutAll{
		UsuarioID: usuarioID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &LogoutAllOutput{}
	out.Body = presentation.NewApiResponse(dto.LogoutResponse{
		SesionesRevocadas: resp.SesionesRevocadas,
	})
	return out, nil
}
