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

type DiagnosticoHandler struct {
	facade facades.DiagnosticosFacade
}

func NewDiagnosticoHandler(facade facades.DiagnosticosFacade) *DiagnosticoHandler {
	return &DiagnosticoHandler{facade: facade}
}

func (h *DiagnosticoHandler) SolicitarManual(c *gin.Context) {
	muestraID := c.Param("muestraID")

	var req dto.SolicitarDiagnosticoManualRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cuerpo inválido", "detalle": err.Error()})
		return
	}

	auth := middleware.GetAuthContext(c)
	if auth == nil {
		auth = &application.AuthContext{}
	}

	resp, err := h.facade.SolicitarManual(c.Request.Context(), auth, muestraID, req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	response := shared.NewResponse(*resp, map[string]shared.Link{
		"self": {Href: c.Request.URL.Path, Method: "POST"},
	})
	c.JSON(http.StatusCreated, response)
}

func (h *DiagnosticoHandler) Aceptar(c *gin.Context) {
	diagnosticoID := c.Param("id")

	auth := middleware.GetAuthContext(c)
	if auth == nil {
		auth = &application.AuthContext{}
	}

	resp, err := h.facade.Aceptar(c.Request.Context(), auth, diagnosticoID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	response := shared.NewResponse(*resp, map[string]shared.Link{
		"self":     {Href: c.Request.URL.Path, Method: "POST"},
		"rechazar": {Href: "/diagnosticos/" + diagnosticoID + "/rechazar", Method: "POST"},
	})
	c.JSON(http.StatusOK, response)
}

func (h *DiagnosticoHandler) Rechazar(c *gin.Context) {
	diagnosticoID := c.Param("id")

	var req dto.RechazarDiagnosticoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cuerpo inválido", "detalle": err.Error()})
		return
	}

	auth := middleware.GetAuthContext(c)
	if auth == nil {
		auth = &application.AuthContext{}
	}

	resp, err := h.facade.Rechazar(c.Request.Context(), auth, diagnosticoID, req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	response := shared.NewResponse(*resp, map[string]shared.Link{
		"self":    {Href: c.Request.URL.Path, Method: "POST"},
		"aceptar": {Href: "/diagnosticos/" + diagnosticoID + "/aceptar", Method: "POST"},
	})
	c.JSON(http.StatusOK, response)
}
