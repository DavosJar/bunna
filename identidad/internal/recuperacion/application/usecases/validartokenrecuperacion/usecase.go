package validartokenrecuperacion

import (
	"context"
	"time"

	dominio "github.com/davosjar/bunna/services/identidad/internal/recuperacion/domain"
)

type ValidarTokenRecuperacionCasoDeUso struct {
	tokenRepo dominio.TokenRecuperacionRepositorio
}

func NewValidarTokenRecuperacionCasoDeUso(
	tokenRepo dominio.TokenRecuperacionRepositorio,
) *ValidarTokenRecuperacionCasoDeUso {
	return &ValidarTokenRecuperacionCasoDeUso{tokenRepo: tokenRepo}
}

func (uc *ValidarTokenRecuperacionCasoDeUso) Ejecutar(ctx context.Context, cmd *ComandoValidarTokenRecuperacion) (*RespuestaValidarTokenRecuperacion, error) {
	if cmd.Token == "" {
		return nil, dominio.ErrEnlaceInvalido
	}

	hash := dominio.HashearToken(cmd.Token)
	token, err := uc.tokenRepo.ObtenerPorHash(ctx, hash)
	if err != nil {
		return nil, dominio.ErrEnlaceInvalido
	}

	if err := token.EsValido(time.Now()); err != nil {
		return nil, err
	}

	return &RespuestaValidarTokenRecuperacion{
		UsuarioID: token.UsuarioID(),
		Valido:    true,
	}, nil
}
