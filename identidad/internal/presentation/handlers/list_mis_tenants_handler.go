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

// ListarMisTenantsOutput es el output del endpoint GET /api/v1/tenants/mis-tenants.
type ListarMisTenantsOutput struct {
	Body presentation.ApiResponse[dto.ListarMisTenantsResponse]
}

// ListarMisTenantsHandler maneja la consulta de tenants del usuario autenticado.
type ListarMisTenantsHandler struct {
	facade facades.TenantFacade
}

// NewListarMisTenantsHandler construye el handler con su facade.
func NewListarMisTenantsHandler(facade facades.TenantFacade) *ListarMisTenantsHandler {
	return &ListarMisTenantsHandler{facade: facade}
}

// Register registra el endpoint en la API Huma.
func (h *ListarMisTenantsHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "listar-mis-tenants",
		Method:      http.MethodGet,
		Path:        "/api/v1/tenants/mis-tenants",
		Summary:     "Listar mis tenants",
		Description: "Retorna la lista de tenants a los que pertenece el usuario autenticado, con su rol en cada uno.",
		Tags:        []string{"Tenants"},
	}, h.handle)
}

func (h *ListarMisTenantsHandler) handle(ctx context.Context, input *struct{}) (*ListarMisTenantsOutput, error) {
	usuarioID := middleware.GetUsuarioIDFromCtx(ctx)
	if usuarioID == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	resp, err := h.facade.ListarMisTenants(ctx, usuarioID)
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	items := make([]dto.TenantConRolDTO, len(resp.Tenants))
	for i, t := range resp.Tenants {
		items[i] = dto.TenantConRolDTO{
			ID:       t.ID,
			Nombre:   t.Nombre,
			Slug:     t.Slug,
			Rol:      t.Rol,
			EsPropio: t.EsPropio,
		}
	}

	out := &ListarMisTenantsOutput{}
	out.Body = presentation.NewApiResponse(dto.ListarMisTenantsResponse{
		Tenants:  items,
		PropioID: resp.PropioID,
	})
	return out, nil
}
