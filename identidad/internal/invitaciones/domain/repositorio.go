package invitaciones

import "context"

type InvitacionRepositorio interface {
	Crear(ctx context.Context, invitacion *Invitacion) error
	ObtenerPorTokenHash(ctx context.Context, tokenHash string) (*Invitacion, error)
	MarcarAceptada(ctx context.Context, id string) error
	ObtenerPorID(ctx context.Context, id string) (*Invitacion, error)
}
