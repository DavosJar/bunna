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

type ConfigurarTenantInput struct {
	TenantID string `path:"tenantID" doc:"ID del tenant a configurar"`
	Body     dto.ConfigurarTenantRequest
}

type ConfigurarTenantOutput struct {
	Body presentation.ApiResponse[dto.ConfigurarTenantResponse]
}

type ConfigurarTenantHandler struct {
	facade facades.TenantFacade
}

func NewConfigurarTenantHandler(facade facades.TenantFacade) *ConfigurarTenantHandler {
	return &ConfigurarTenantHandler{facade: facade}
}

func (h *ConfigurarTenantHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "configurar-tenant",
		Method:      http.MethodPut,
		Path:        "/api/v1/tenants/{tenantID}",
		Summary:     "Configurar tenant",
		Description: "Actualiza la configuración de un tenant.",
		Tags:        []string{"Tenants"},
	}, h.handle)
}

func (h *ConfigurarTenantHandler) handle(ctx context.Context, input *ConfigurarTenantInput) (*ConfigurarTenantOutput, error) {
	ejecutorID := middleware.GetUsuarioIDFromCtx(ctx)
	if ejecutorID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	resp, err := h.facade.ConfigurarTenant(ctx, facades.ComandoConfigurarTenant{
		TenantID:   input.TenantID,
		Nombre:     input.Body.Nombre,
		Slug:       input.Body.Slug,
		EjecutorID: ejecutorID,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &ConfigurarTenantOutput{}
	out.Body = presentation.NewApiResponse(dto.ConfigurarTenantResponse{
		TenantID: resp.TenantID, Nombre: resp.Nombre,
		Slug: resp.Slug, ModificadoEn: resp.ModificadoEn,
	})
	return out, nil
}
