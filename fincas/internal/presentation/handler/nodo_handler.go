package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/davosjar/bunna/services/fincas/internal/application"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/dto"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/facades"
	"github.com/davosjar/bunna/services/fincas/internal/presentation/middleware"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
	pshared "github.com/davosjar/bunna/services/fincas/shared/presentation"
)

type NodoHandler struct {
	facade facades.NodosFacade
}

func NewNodoHandler(facade facades.NodosFacade) *NodoHandler {
	return &NodoHandler{facade: facade}
}

func (h *NodoHandler) Registrar(c *gin.Context) {
	var req dto.RegistrarNodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cuerpo inválido", "detalle": err.Error()})
		return
	}

	auth := middleware.GetAuthContext(c)
	if auth == nil {
		auth = &application.AuthContext{}
	}

	resp, err := h.facade.Registrar(c.Request.Context(), auth, req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	response := pshared.NewResponse(*resp, map[string]pshared.Link{
		"self": {Href: c.Request.URL.Path, Method: "POST"},
	})
	c.JSON(http.StatusCreated, response)
}

func (h *NodoHandler) Listar(c *gin.Context) {
	auth := middleware.GetAuthContext(c)
	if auth == nil {
		auth = &application.AuthContext{}
	}

	pagina, _ := strconv.Atoi(c.DefaultQuery("pagina", "1"))
	tamano, _ := strconv.Atoi(c.DefaultQuery("tamanoPagina", "10"))

	paginacion := shared.Paginacion{
		Pagina:       pagina,
		TamanoPagina: tamano,
	}

	resp, err := h.facade.Listar(c.Request.Context(), auth, paginacion)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	response := pshared.NewResponse(*resp, map[string]pshared.Link{
		"self": {Href: c.Request.URL.Path, Method: "GET"},
	})
	c.JSON(http.StatusOK, response)
}

func (h *NodoHandler) Obtener(c *gin.Context) {
	nodoID := c.Param("id")

	auth := middleware.GetAuthContext(c)
	if auth == nil {
		auth = &application.AuthContext{}
	}

	resp, err := h.facade.Obtener(c.Request.Context(), auth, nodoID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	response := pshared.NewResponse(*resp, map[string]pshared.Link{
		"self": {Href: c.Request.URL.Path, Method: "GET"},
	})
	c.JSON(http.StatusOK, response)
}

func (h *NodoHandler) Editar(c *gin.Context) {
	nodoID := c.Param("id")

	var req dto.EditarNodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cuerpo inválido", "detalle": err.Error()})
		return
	}

	auth := middleware.GetAuthContext(c)
	if auth == nil {
		auth = &application.AuthContext{}
	}

	resp, err := h.facade.Editar(c.Request.Context(), auth, nodoID, req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	response := pshared.NewResponse(*resp, map[string]pshared.Link{
		"self": {Href: c.Request.URL.Path, Method: "PUT"},
	})
	c.JSON(http.StatusOK, response)
}

func (h *NodoHandler) Desactivar(c *gin.Context) {
	nodoID := c.Param("id")

	var req dto.DesactivarNodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cuerpo inválido", "detalle": err.Error()})
		return
	}

	auth := middleware.GetAuthContext(c)
	if auth == nil {
		auth = &application.AuthContext{}
	}

	resp, err := h.facade.Desactivar(c.Request.Context(), auth, nodoID, req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	response := pshared.NewResponse(*resp, map[string]pshared.Link{
		"self": {Href: c.Request.URL.Path, Method: "POST"},
	})
	c.JSON(http.StatusOK, response)
}

func (h *NodoHandler) Validar(c *gin.Context) {
	nodeKey := c.Query("nodeKey")
	if nodeKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parámetro nodeKey es requerido"})
		return
	}

	resp, err := h.facade.Validar(c.Request.Context(), nodeKey)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *NodoHandler) RegistrarInferencia(c *gin.Context) {
	var req dto.RegistrarInferenciaDesdeNodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("ERROR BINDING INFERENCIA: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "cuerpo inválido", "detalle": err.Error()})
		return
	}
	fmt.Printf("YOLO CALLBACK REQ: %+v\n", req)

	resp, err := h.facade.RegistrarInferencia(c.Request.Context(), req)
	if err != nil {
		fmt.Printf("ERROR FACADE INFERENCIA: %v\n", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	response := pshared.NewResponse(*resp, map[string]pshared.Link{
		"self": {Href: c.Request.URL.Path, Method: "POST"},
	})
	c.JSON(http.StatusCreated, response)
}
