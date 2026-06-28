package invitaciones

import (
	"context"

	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
)

type InvitacionRepositorio interface {
	Crear(ctx context.Context, invitacion *Invitacion) error
	ObtenerPorTokenHash(ctx context.Context, tokenHash string) (*Invitacion, error)
	MarcarAceptada(ctx context.Context, id string) error
	ObtenerPorID(ctx context.Context, id string) (*Invitacion, error)
	ListarPorTenant(ctx context.Context, tenantID string, paginacion shareddomain.Paginacion, estado string) ([]*Invitacion, int, error)
	ActualizarTokenHash(ctx context.Context, id string, tokenHash string) error
	Eliminar(ctx context.Context, id string) error
}
