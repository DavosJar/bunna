package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/davosjar/bunna/services/fincas/internal/application"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/dto"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/facades"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/middleware"
	shared "github.com/davosjar/bunna/services/fincas/shared/presentation"
)

type MuestraHandler struct {
	facade facades.MuestrasFacade
}

func NewMuestraHandler(facade facades.MuestrasFacade) *MuestraHandler {
	return &MuestraHandler{facade: facade}
}

func (h *MuestraHandler) Tomar(c *gin.Context) {
	fincaID := c.Param("id")
	loteID := c.Param("loteID")
	if c.FullPath() == "/api/v1/lotes/:id/muestras" {
		loteID = c.Param("id")
		fincaID = ""
	}

	var req dto.TomarMuestraRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cuerpo inválido", "detalle": err.Error()})
		return
	}

	auth := middleware.GetAuthContext(c)
	if auth == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}
	if !auth.TienePermiso(application.PermisoCrearMuestra) {
		c.JSON(http.StatusForbidden, gin.H{"error": "permiso denegado"})
		return
	}

	resp, err := h.facade.Tomar(c.Request.Context(), auth, fincaID, loteID, req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	response := shared.NewResponse(*resp, map[string]shared.Link{
		"self":               {Href: c.Request.URL.Path, Method: "POST"},
		"diagnostico_manual": {Href: "/muestras/" + resp.ID + "/diagnosticos/manual", Method: "POST"},
	})
	c.JSON(http.StatusCreated, response)
}

func (h *MuestraHandler) ListarPorLote(c *gin.Context) {
	fincaID := c.Param("id")
	loteID := c.Param("loteID")
	if c.FullPath() == "/api/v1/lotes/:id/muestras" {
		loteID = c.Param("id")
		fincaID = ""
	}

	auth := middleware.GetAuthContext(c)
	if auth == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}
	if !auth.TienePermiso(application.PermisoVerMuestras) {
		c.JSON(http.StatusForbidden, gin.H{"error": "permiso denegado"})
		return
	}

	items, err := h.facade.ListarPorLote(c.Request.Context(), auth, fincaID, loteID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if items == nil {
		items = []dto.MuestraItemResponse{}
	}

	response := shared.NewResponse(items, map[string]shared.Link{
		"self":  {Href: c.Request.URL.Path, Method: "GET"},
		"tomar": {Href: "/lotes/" + loteID + "/muestras", Method: "POST"},
	})
	c.JSON(http.StatusOK, response)
}
