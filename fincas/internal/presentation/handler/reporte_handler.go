package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/davosjar/bunna/services/fincas/internal/application"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/facades"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/middleware"
	shared "github.com/davosjar/bunna/services/fincas/shared/presentation"
)

type ReporteHandler struct {
	facade facades.ReportesFacade
}

func NewReporteHandler(facade facades.ReportesFacade) *ReporteHandler {
	return &ReporteHandler{facade: facade}
}

func (h *ReporteHandler) GenerarPorLote(c *gin.Context) {
	loteID := c.Param("loteID")
	if loteID == "" {
		loteID = c.Param("id")
	}

	auth := middleware.GetAuthContext(c)
	if auth == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}
	if !auth.TienePermiso(application.PermisoGenerarReporte) {
		c.JSON(http.StatusForbidden, gin.H{"error": "permiso denegado"})
		return
	}

	resp, err := h.facade.GenerarPorLote(c.Request.Context(), auth, loteID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	response := shared.NewResponse(*resp, map[string]shared.Link{
		"self":     {Href: c.Request.URL.Path, Method: "GET"},
		"muestras": {Href: "/lotes/" + loteID + "/muestras", Method: "GET"},
	})
	c.JSON(http.StatusOK, response)
}
