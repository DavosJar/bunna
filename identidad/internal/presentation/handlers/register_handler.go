package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/dto"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/facades"
	presentation "github.com/davosjar/bunna/services/identidad/internal/shared/presentation"
)

// RegisterInput es el input del endpoint POST /api/v1/identidad/auth/register.
type RegisterInput struct {
	Body dto.RegisterRequest
}

// RegisterOutput es el output del endpoint POST /api/v1/identidad/auth/register.
type RegisterOutput struct {
	Body presentation.ApiResponse[dto.RegisterResponse]
}

// RegisterHandler maneja el registro de nuevos usuarios.
type RegisterHandler struct {
	facade facades.AuthFacade
}

// NewRegisterHandler construye el handler con su facade.
func NewRegisterHandler(facade facades.AuthFacade) *RegisterHandler {
	return &RegisterHandler{facade: facade}
}

// Register registra el endpoint en la API Huma.
func (h *RegisterHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "register-usuario",
		Method:        http.MethodPost,
		Path:          "/api/v1/identidad/auth/register",
		Summary:       "Registro de usuario",
		Description:   "Crea un nuevo usuario con sus credenciales. Devuelve el ID y estado inicial.",
		Tags:          []string{"Autenticación"},
		DefaultStatus: http.StatusCreated,
	}, h.handle)
}

func (h *RegisterHandler) handle(ctx context.Context, input *RegisterInput) (*RegisterOutput, error) {
	telefono := ""
	if input.Body.Telefono != nil {
		telefono = *input.Body.Telefono
	}

	resp, err := h.facade.Registrar(ctx, facades.ComandoRegistro{
		Nombre:   input.Body.Nombre,
		Apellido: input.Body.Apellido,
		Correo:   input.Body.Correo,
		Password: input.Body.Password,
		Telefono: telefono,
	})
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	// CON-PRES-005: links HATEOAS los construye el handler
	links := map[string]presentation.Link{
		"self": {
			Href:   "/api/v1/identidad/usuarios/" + resp.UsuarioID,
			Method: http.MethodGet,
		},
	}

	out := &RegisterOutput{}
	out.Body = presentation.NewApiResponseWithLinks(dto.RegisterResponse{
		UsuarioID: resp.UsuarioID,
		Correo:    resp.Correo,
		Estado:    resp.Estado,
	}, links)

	return out, nil
}
