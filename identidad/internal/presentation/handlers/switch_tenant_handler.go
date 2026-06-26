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

// SwitchTenantInput es el input del endpoint POST /api/v1/auth/switch-tenant.
type SwitchTenantInput struct {
	Body dto.SwitchTenantRequest
}

// SwitchTenantOutput es el output del endpoint POST /api/v1/auth/switch-tenant.
type SwitchTenantOutput struct {
	Body presentation.ApiResponse[dto.SwitchTenantResponse]
}

// SwitchTenantHandler maneja el cambio de tenant activo.
type SwitchTenantHandler struct {
	facade facades.AuthFacade
}

// NewSwitchTenantHandler construye el handler con su facade.
func NewSwitchTenantHandler(facade facades.AuthFacade) *SwitchTenantHandler {
	return &SwitchTenantHandler{facade: facade}
}

// Register registra el endpoint en la API Huma.
func (h *SwitchTenantHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "switch-tenant",
		Method:      http.MethodPost,
		Path:        "/api/v1/auth/switch-tenant",
		Summary:     "Cambiar de tenant",
		Description: "Cambia el tenant activo del usuario autenticado y retorna nuevos tokens JWT con el tenant y rol actualizados.",
		Tags:        []string{"Autenticación"},
	}, h.handle)
}

func (h *SwitchTenantHandler) handle(ctx context.Context, input *SwitchTenantInput) (*SwitchTenantOutput, error) {
	usuarioID := middleware.GetUsuarioIDFromCtx(ctx)
	sesionID := middleware.GetSesionIDFromCtx(ctx)
	if usuarioID == "" || sesionID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	resp, err := h.facade.SwitchTenant(ctx, facades.ComandoSwitchTenant{
		UsuarioID: usuarioID,
		SesionID:  sesionID,
		TenantID:  input.Body.TenantID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &SwitchTenantOutput{}
	out.Body = presentation.NewApiResponse(dto.SwitchTenantResponse{
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
