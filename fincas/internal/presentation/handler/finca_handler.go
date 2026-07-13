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

type FincaHandler struct {
	facade facades.FincasFacade
}

func NewFincaHandler(facade facades.FincasFacade) *FincaHandler {
	return &FincaHandler{facade: facade}
}

func (h *FincaHandler) Registrar(c *gin.Context) {
	var req dto.RegistrarFincaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cuerpo inválido", "detalle": err.Error()})
		return
	}

	auth := middleware.GetAuthContext(c)
	if auth == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}
	if !auth.TienePermiso(application.PermisoCrearFinca) {
		c.JSON(http.StatusForbidden, gin.H{"error": "permiso denegado"})
		return
	}

	resp, err := h.facade.Registrar(c.Request.Context(), auth, req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	response := shared.NewResponse(*resp, map[string]shared.Link{
		"self":        {Href: c.Request.URL.Path, Method: "POST"},
		"desactivar":  {Href: "/fincas/" + resp.ID + "/desactivar", Method: "POST"},
	})
	c.JSON(http.StatusCreated, response)
}

func (h *FincaHandler) Desactivar(c *gin.Context) {
	fincaID := c.Param("id")

	var req dto.DesactivarFincaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cuerpo inválido", "detalle": err.Error()})
		return
	}

	auth := middleware.GetAuthContext(c)
	if auth == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}
	if !auth.TienePermiso(application.PermisoDesactivarFinca) {
		c.JSON(http.StatusForbidden, gin.H{"error": "permiso denegado"})
		return
	}

	resp, err := h.facade.Desactivar(c.Request.Context(), auth, fincaID, req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	response := shared.NewResponse(*resp, map[string]shared.Link{
		"self": {Href: c.Request.URL.Path, Method: "POST"},
	})
	c.JSON(http.StatusOK, response)
}

func (h *FincaHandler) Listar(c *gin.Context) {
	auth := middleware.GetAuthContext(c)
	if auth == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	resp, err := h.facade.Listar(c.Request.Context(), auth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}
