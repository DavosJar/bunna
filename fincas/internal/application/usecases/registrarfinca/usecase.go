package registrarfinca

import (
	"context"
	"time"

	"github.com/davosjar/bunna/services/fincas/internal/application"
	"github.com/davosjar/bunna/services/fincas/internal/fincas/domain"
	shared "github.com/davosjar/bunna/services/fincas/internal/shared/domain"
)

const permisoRequerido = "CREAR_FINCA"

// FincaCreada es el evento publicado en RabbitMQ tras el registro exitoso.
type FincaCreada struct {
	EventID    string    `json:"event_id"`
	FincaID    string    `json:"finca_id"`
	Nombre     string    `json:"nombre"`
	UsuarioID  string    `json:"usuario_id"`
	TenantID   string    `json:"tenant_id,omitempty"`
	OcurredAt  time.Time `json:"ocurred_at"`
}

// UseCase orquesta el registro de una nueva finca.
type UseCase struct {
	fincaRepo  domain.FincaRepositorio
	generador  shared.GeneradorID
	publisher  application.EventPublisher
}

func NewUseCase(
	fincaRepo domain.FincaRepositorio,
	generador shared.GeneradorID,
	publisher application.EventPublisher,
) *UseCase {
	return &UseCase{
		fincaRepo: fincaRepo,
		generador: generador,
		publisher: publisher,
	}
}

func (uc *UseCase) Ejecutar(ctx context.Context, auth *application.AuthContext, cmd Command) (*Salida, error) {
	// 1. Validar permisos
	if !auth.TienePermiso(permisoRequerido) {
		return nil, application.ErrForbidden
	}

	// 2. Validar campos del comando
	if err := cmd.Validar(); err != nil {
		return nil, err
	}

	// 3. Generar ID
	id, err := uc.generador.NextID(ctx)
	if err != nil {
		return nil, err
	}

	// 4. Construir entidad Finca vía constructor de dominio
	finca := domain.NuevaFinca(cmd.Nombre, cmd.Ubicacion, cmd.Descripcion, auth.UsuarioID)

	// 5. Asignar ID y TenantID. El dominio no expone setters, por lo que
	// reconstruimos desde persistencia con los valores correctos.
	var tenantID *string
	if auth.TenantID != "" {
		tenantID = &auth.TenantID
	}
	now := time.Now()
	finca = domain.NewFincaFromPersistence(
		id, cmd.Nombre, cmd.Ubicacion, cmd.Descripcion, auth.UsuarioID,
		tenantID, domain.FincaActiva, now, now,
	)

	// 6. Persistir
	if err := uc.fincaRepo.Crear(ctx, finca); err != nil {
		return nil, err
	}

	// 7. Publicar evento
	eventID, _ := uc.generador.NextID(ctx)
	evento := FincaCreada{
		EventID:   eventID,
		FincaID:   id,
		Nombre:    cmd.Nombre,
		UsuarioID: auth.UsuarioID,
		TenantID:  auth.TenantID,
		OcurredAt: now,
	}
	_ = uc.publisher.Publish(ctx, "fincas.v1.finca.creada", evento)

	return &Salida{
		ID:          id,
		Nombre:      finca.Nombre(),
		Ubicacion:   finca.Ubicacion(),
		Descripcion: finca.Descripcion(),
		Estado:      string(finca.Estado()),
		CreatedAt:   finca.CreatedAt(),
	}, nil
}
