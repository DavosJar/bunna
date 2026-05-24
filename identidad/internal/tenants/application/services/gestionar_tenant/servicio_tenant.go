package gestionar_tenant

import (
	"context"
	"fmt"

	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
)

// ServicioTenant maneja los casos de uso de gestión de tenants
type ServicioTenant struct {
	tenantRepo    tenant.TenantRepositorio
	membresiaRepo tenant.MembresiaRepositorio
	idGenerator   shareddomain.GeneradorID
}

// NuevoServicioTenant crea una nueva instancia del servicio
func NuevoServicioTenant(
	tenantRepo tenant.TenantRepositorio,
	membresiaRepo tenant.MembresiaRepositorio,
	idGenerator shareddomain.GeneradorID,
) *ServicioTenant {
	return &ServicioTenant{
		tenantRepo:    tenantRepo,
		membresiaRepo: membresiaRepo,
		idGenerator:   idGenerator,
	}
}

// CrearTenant crea un nuevo tenant (solo SYS_ADMIN)
func (s *ServicioTenant) CrearTenant(ctx context.Context, cmd *ComandoCrearTenant) (*DtoTenant, error) {
	nuevoID, err := s.idGenerator.NextID(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al generar ID: %w", err)
	}

	nuevoTenant, err := tenant.NuevoTenant(nuevoID, cmd.Nombre, cmd.Slug)
	if err != nil {
		return nil, err
	}

	creado, err := s.tenantRepo.Crear(ctx, nuevoTenant)
	if err != nil {
		return nil, err
	}

	return toDto(creado), nil
}

// ActivarTenant activa un tenant inactivo (solo SYS_ADMIN)
func (s *ServicioTenant) ActivarTenant(ctx context.Context, cmd *ComandoActivarTenant) error {
	t, err := s.tenantRepo.ObtenerPorID(ctx, cmd.TenantID)
	if err != nil {
		return err
	}

	if err := t.Activar(); err != nil {
		return err
	}

	_, err = s.tenantRepo.Actualizar(ctx, t)
	return err
}

// DesactivarTenant desactiva un tenant activo (solo SYS_ADMIN)
func (s *ServicioTenant) DesactivarTenant(ctx context.Context, cmd *ComandoDesactivarTenant) error {
	t, err := s.tenantRepo.ObtenerPorID(ctx, cmd.TenantID)
	if err != nil {
		return err
	}

	if err := t.Desactivar(); err != nil {
		return err
	}

	_, err = s.tenantRepo.Actualizar(ctx, t)
	return err
}

// ListarTodos retorna todos los tenants (solo SYS_ADMIN)
func (s *ServicioTenant) ListarTodos(ctx context.Context) ([]*DtoTenant, error) {
	tenants, err := s.tenantRepo.Listar(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]*DtoTenant, len(tenants))
	for i, t := range tenants {
		dtos[i] = toDto(t)
	}
	return dtos, nil
}

// ListarMisTenants retorna los tenants del usuario autenticado
func (s *ServicioTenant) ListarMisTenants(ctx context.Context, usuarioID string) ([]*DtoTenant, error) {
	tenants, err := s.tenantRepo.ListarPorUsuario(ctx, usuarioID)
	if err != nil {
		return nil, err
	}

	dtos := make([]*DtoTenant, len(tenants))
	for i, t := range tenants {
		dtos[i] = toDto(t)
	}
	return dtos, nil
}

// ObtenerPorID retorna el detalle de un tenant
func (s *ServicioTenant) ObtenerPorID(ctx context.Context, id string) (*DtoTenant, error) {
	t, err := s.tenantRepo.ObtenerPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDto(t), nil
}

// ObtenerPorSlug retorna el detalle de un tenant por su slug
func (s *ServicioTenant) ObtenerPorSlug(ctx context.Context, slug string) (*DtoTenant, error) {
	t, err := s.tenantRepo.ObtenerPorSlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return toDto(t), nil
}

// AgregarUsuario agrega un usuario como miembro de un tenant
func (s *ServicioTenant) AgregarUsuario(ctx context.Context, cmd *ComandoAgregarUsuario) error {
	t, err := s.tenantRepo.ObtenerPorID(ctx, cmd.TenantID)
	if err != nil {
		return err
	}

	if !t.EstaActivo() {
		return tenant.ErrTenantInactivo
	}

	existe, err := s.membresiaRepo.ExisteMiembro(ctx, cmd.UsuarioID, cmd.TenantID)
	if err != nil {
		return err
	}
	if existe {
		return tenant.ErrUsuarioYaMiembro
	}

	membresia, err := tenant.NuevaMembresia(cmd.UsuarioID, cmd.TenantID)
	if err != nil {
		return err
	}

	return s.membresiaRepo.Crear(ctx, membresia)
}

// RemoverUsuario elimina la membresía de un usuario en un tenant
func (s *ServicioTenant) RemoverUsuario(ctx context.Context, cmd *ComandoRemoverUsuario) error {
	existe, err := s.membresiaRepo.ExisteMiembro(ctx, cmd.UsuarioID, cmd.TenantID)
	if err != nil {
		return err
	}
	if !existe {
		return tenant.ErrUsuarioNoEsMiembro
	}

	return s.membresiaRepo.Eliminar(ctx, cmd.UsuarioID, cmd.TenantID)
}

// ListarUsuariosDeTenant retorna los IDs de usuarios miembros de un tenant
func (s *ServicioTenant) ListarUsuariosDeTenant(ctx context.Context, tenantID string) ([]string, error) {
	return s.membresiaRepo.ListarUsuariosPorTenant(ctx, tenantID)
}

// toDto convierte un Tenant de dominio a DTO
func toDto(t *tenant.Tenant) *DtoTenant {
	return &DtoTenant{
		ID:            t.ID(),
		Nombre:        t.Nombre(),
		Slug:          t.Slug(),
		Activo:        t.EstaActivo(),
		FechaCreacion: t.FechaCreacion(),
	}
}
