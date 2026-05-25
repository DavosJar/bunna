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

type LoteHandler struct {
	facade facades.LotesFacade
}

func NewLoteHandler(facade facades.LotesFacade) *LoteHandler {
	return &LoteHandler{facade: facade}
}

func (h *LoteHandler) Agregar(c *gin.Context) {
	fincaID := c.Param("fincaID")

	var req dto.AgregarLoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cuerpo inválido", "detalle": err.Error()})
		return
	}

	auth := middleware.GetAuthContext(c)
	if auth == nil {
		auth = &application.AuthContext{}
	}

	resp, err := h.facade.Agregar(c.Request.Context(), auth, fincaID, req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	response := shared.NewResponse(*resp, map[string]shared.Link{
		"self":     {Href: c.Request.URL.Path, Method: "POST"},
		"eliminar": {Href: "/lotes/" + resp.ID + "/eliminar", Method: "POST"},
		"muestras": {Href: "/lotes/" + resp.ID + "/muestras", Method: "GET"},
	})
	c.JSON(http.StatusCreated, response)
}

func (h *LoteHandler) Eliminar(c *gin.Context) {
	loteID := c.Param("id")

	auth := middleware.GetAuthContext(c)
	if auth == nil {
		auth = &application.AuthContext{}
	}

	resp, err := h.facade.Eliminar(c.Request.Context(), auth, loteID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	response := shared.NewResponse(*resp, map[string]shared.Link{
		"self": {Href: c.Request.URL.Path, Method: "POST"},
	})
	c.JSON(http.StatusOK, response)
}
