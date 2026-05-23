package postgres

import (
	"context"
	"errors"

	"github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
	"gorm.io/gorm"
)

type tenantRepositorio struct {
	db *gorm.DB
}

func NewTenantRepositorio(db *gorm.DB) tenant.TenantRepositorio {
	return &tenantRepositorio{db: db}
}

func (r *tenantRepositorio) Crear(ctx context.Context, t *tenant.Tenant) (*tenant.Tenant, error) {
	model := TenantFromDomain(t)
	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, tenant.ErrSlugDuplicado
		}
		return nil, result.Error
	}
	return model.ToDomain(), nil
}

func (r *tenantRepositorio) ObtenerPorID(ctx context.Context, id string) (*tenant.Tenant, error) {
	var model TenantModel
	result := r.db.WithContext(ctx).First(&model, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, tenant.ErrTenantNoEncontrado
		}
		return nil, result.Error
	}
	return model.ToDomain(), nil
}

func (r *tenantRepositorio) ObtenerPorSlug(ctx context.Context, slug string) (*tenant.Tenant, error) {
	var model TenantModel
	result := r.db.WithContext(ctx).First(&model, "slug = ?", slug)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, tenant.ErrTenantNoEncontrado
		}
		return nil, result.Error
	}
	return model.ToDomain(), nil
}

func (r *tenantRepositorio) Actualizar(ctx context.Context, t *tenant.Tenant) (*tenant.Tenant, error) {
	result := r.db.WithContext(ctx).Model(&TenantModel{}).
		Where("id = ?", t.ID()).
		Updates(map[string]interface{}{
			"nombre":     t.Nombre(),
			"activo":     t.EstaActivo(),
			"updated_at": t.FechaActualizacion(),
		})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, tenant.ErrTenantNoEncontrado
	}
	return r.ObtenerPorID(ctx, t.ID())
}

func (r *tenantRepositorio) Listar(ctx context.Context) ([]*tenant.Tenant, error) {
	var models []TenantModel
	result := r.db.WithContext(ctx).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}
	tenants := make([]*tenant.Tenant, len(models))
	for i, m := range models {
		tenants[i] = m.ToDomain()
	}
	return tenants, nil
}

func (r *tenantRepositorio) ListarPorUsuario(ctx context.Context, usuarioID string) ([]*tenant.Tenant, error) {
	var models []TenantModel
	result := r.db.WithContext(ctx).
		Joins("JOIN usuario_tenants ON usuario_tenants.tenant_id = tenants.id").
		Where("usuario_tenants.usuario_id = ?", usuarioID).
		Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}
	tenants := make([]*tenant.Tenant, len(models))
	for i, m := range models {
		tenants[i] = m.ToDomain()
	}
	return tenants, nil
}

// membresiaRepositorio

type membresiaRepositorio struct {
	db *gorm.DB
}

func NewMembresiaRepositorio(db *gorm.DB) tenant.MembresiaRepositorio {
	return &membresiaRepositorio{db: db}
}

func (r *membresiaRepositorio) Crear(ctx context.Context, m *tenant.Membresia) error {
	model := &MembresiaModel{
		UsuarioID: m.UsuarioID(),
		TenantID:  m.TenantID(),
	}
	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return tenant.ErrUsuarioYaMiembro
		}
		return result.Error
	}
	return nil
}

func (r *membresiaRepositorio) Eliminar(ctx context.Context, usuarioID, tenantID string) error {
	result := r.db.WithContext(ctx).
		Where("usuario_id = ? AND tenant_id = ?", usuarioID, tenantID).
		Delete(&MembresiaModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return tenant.ErrUsuarioNoEsMiembro
	}
	return nil
}

func (r *membresiaRepositorio) ExisteMiembro(ctx context.Context, usuarioID, tenantID string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&MembresiaModel{}).
		Where("usuario_id = ? AND tenant_id = ?", usuarioID, tenantID).
		Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (r *membresiaRepositorio) ListarUsuariosPorTenant(ctx context.Context, tenantID string) ([]string, error) {
	var models []MembresiaModel
	result := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}
	ids := make([]string, len(models))
	for i, m := range models {
		ids[i] = m.UsuarioID
	}
	return ids, nil
}