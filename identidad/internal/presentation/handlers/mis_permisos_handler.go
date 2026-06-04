package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/facades"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/middleware"
	presentation "github.com/davosjar/bunna/services/identidad/shared/presentation"
)

type misPermisoItem struct {
	Codigo      string `json:"codigo" doc:"Código del permiso"`
	Nombre      string `json:"nombre" doc:"Nombre del permiso"`
	Descripcion string `json:"descripcion" doc:"Descripción del permiso"`
	Modulo      string `json:"modulo" doc:"Módulo del permiso"`
}

type listarMisPermisosData struct {
	Permisos []misPermisoItem `json:"permisos" doc:"Lista de permisos del usuario autenticado"`
}

type ListarMisPermisosOutput struct {
	Body presentation.ApiResponse[listarMisPermisosData]
}

type ListarMisPermisosHandler struct {
	facade facades.RbacFacade
}

func NewListarMisPermisosHandler(facade facades.RbacFacade) *ListarMisPermisosHandler {
	return &ListarMisPermisosHandler{facade: facade}
}

func (h *ListarMisPermisosHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "listar-mis-permisos",
		Method:      http.MethodGet,
		Path:        "/api/v1/mis-permisos",
		Summary:     "Listar mis permisos",
		Description: "Retorna los códigos de permiso del usuario autenticado según su rol en el tenant activo.",
		Tags:        []string{"Permisos"},
	}, h.handle)
}

func (h *ListarMisPermisosHandler) handle(ctx context.Context, input *struct{}) (*ListarMisPermisosOutput, error) {
	rol := middleware.GetRolFromCtx(ctx)
	if rol == "" {
		return nil, huma.Error401Unauthorized("token requerido")
	}

	tenantID := middleware.GetTenantIDFromCtx(ctx)

	resp, err := h.facade.ListarMisPermisos(ctx, rol, tenantID)
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	usuarioID := middleware.GetUsuarioIDFromCtx(ctx)
	log.Printf("[MisPermisos] usuario=%s rol=%s tenant=%s permisos=%v", usuarioID, rol, tenantID, resp.Permisos)

	items := make([]misPermisoItem, len(resp.Permisos))
	for i, p := range resp.Permisos {
		items[i] = misPermisoItem{
			Codigo:      p.Codigo,
			Nombre:      p.Nombre,
			Descripcion: p.Descripcion,
			Modulo:      p.Modulo,
		}
	}

	out := &ListarMisPermisosOutput{}
	out.Body = presentation.NewApiResponse(listarMisPermisosData{
		Permisos: items,
	})
	return out, nil
}
