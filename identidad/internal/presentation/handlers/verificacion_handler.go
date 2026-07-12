package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/dto"
	"github.com/davosjar/bunna/services/identidad/internal/presentation/facades"
	presentation "github.com/davosjar/bunna/services/identidad/internal/shared/presentation"
)

// ── Solicitar Verificación ─────────────────────────────────────────────────────

type SolicitarVerificacionInput struct {
	Body struct {
		Email string `json:"correo" doc:"Correo electrónico del usuario"`
	}
}

type SolicitarVerificacionOutput struct {
	Body presentation.ApiResponse[dto.SolicitarVerificacionResponse]
}

type SolicitarVerificacionHandler struct {
	facade facades.VerificacionFacade
}

func NewSolicitarVerificacionHandler(facade facades.VerificacionFacade) *SolicitarVerificacionHandler {
	return &SolicitarVerificacionHandler{facade: facade}
}

func (h *SolicitarVerificacionHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "solicitar-verificacion",
		Method:      http.MethodPost,
		Path:        "/api/v1/identidad/verificacion/solicitar",
		Summary:     "Solicitar verificación de correo",
		Description: "Solicita un enlace de verificación para el correo electrónico indicado. Endpoint público (no requiere JWT).",
		Tags:        []string{"Verificación"},
	}, h.handle)
}

func (h *SolicitarVerificacionHandler) handle(ctx context.Context, input *SolicitarVerificacionInput) (*SolicitarVerificacionOutput, error) {
	email := input.Body.Email
	if email == "" {
		return nil, huma.Error422UnprocessableEntity("email es requerido")
	}

	resp, err := h.facade.SolicitarVerificacion(ctx, facades.ComandoSolicitarVerificacion{
		Email: email,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &SolicitarVerificacionOutput{}
	out.Body = presentation.NewApiResponse(dto.SolicitarVerificacionResponse{Mensaje: resp.Mensaje})
	return out, nil
}

// ── Confirmar Verificación ─────────────────────────────────────────────────────

type ConfirmarVerificacionInput struct {
	Body dto.ConfirmarVerificacionRequest
}

type ConfirmarVerificacionOutput struct {
	Body presentation.ApiResponse[dto.ConfirmarVerificacionResponse]
}

type ConfirmarVerificacionHandler struct {
	facade facades.VerificacionFacade
}

func NewConfirmarVerificacionHandler(facade facades.VerificacionFacade) *ConfirmarVerificacionHandler {
	return &ConfirmarVerificacionHandler{facade: facade}
}

func (h *ConfirmarVerificacionHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "confirmar-verificacion",
		Method:      http.MethodPost,
		Path:        "/api/v1/identidad/verificacion/confirmar",
		Summary:     "Confirmar verificación de correo",
		Description: "Confirma la verificación del correo electrónico usando el token recibido.",
		Tags:        []string{"Verificación"},
	}, h.handle)
}

func (h *ConfirmarVerificacionHandler) handle(ctx context.Context, input *ConfirmarVerificacionInput) (*ConfirmarVerificacionOutput, error) {
	resp, err := h.facade.ConfirmarVerificacion(ctx, facades.ComandoConfirmarVerificacion{
		Token: input.Body.Token,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &ConfirmarVerificacionOutput{}
	out.Body = presentation.NewApiResponse(dto.ConfirmarVerificacionResponse{Mensaje: resp.Mensaje})
	return out, nil
}

// ── Reenviar Verificación ─────────────────────────────────────────────────────

type ReenviarVerificacionInput struct {
	Body struct {
		Email string `json:"correo" doc:"Correo electrónico del usuario"`
	}
}

type ReenviarVerificacionOutput struct {
	Body presentation.ApiResponse[dto.ReenviarVerificacionResponse]
}

type ReenviarVerificacionHandler struct {
	facade facades.VerificacionFacade
}

func NewReenviarVerificacionHandler(facade facades.VerificacionFacade) *ReenviarVerificacionHandler {
	return &ReenviarVerificacionHandler{facade: facade}
}

func (h *ReenviarVerificacionHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "reenviar-verificacion",
		Method:      http.MethodPost,
		Path:        "/api/v1/identidad/verificacion/reenviar",
		Summary:     "Reenviar verificación de correo",
		Description: "Reenvía el enlace de verificación al correo indicado. Endpoint público (no requiere JWT).",
		Tags:        []string{"Verificación"},
	}, h.handle)
}

func (h *ReenviarVerificacionHandler) handle(ctx context.Context, input *ReenviarVerificacionInput) (*ReenviarVerificacionOutput, error) {
	email := input.Body.Email
	if email == "" {
		return nil, huma.Error422UnprocessableEntity("email es requerido")
	}

	resp, err := h.facade.ReenviarVerificacion(ctx, facades.ComandoReenviarVerificacion{
		Email: email,
	})
	if err != nil {
		return nil, presentation.MapearError(err)
	}

	out := &ReenviarVerificacionOutput{}
	out.Body = presentation.NewApiResponse(dto.ReenviarVerificacionResponse{Mensaje: resp.Mensaje})
	return out, nil
}
