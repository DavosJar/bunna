package postgres

import (
	"context"
	"errors"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// --- RolRepositorio ---

type rolRepositorio struct{ db *gorm.DB }

func NewRolRepositorio(db *gorm.DB) rbac.RolRepositorio {
	return &rolRepositorio{db: db}
}

func (r *rolRepositorio) ObtenerPorNombre(ctx context.Context, nombre string) (*rbac.RolDB, error) {
	var m RolModel
	if err := r.db.WithContext(ctx).First(&m, "nombre = ?", nombre).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rbac.ErrRolNoEncontrado
		}
		return nil, err
	}
	return toRolDB(&m), nil
}

func (r *rolRepositorio) ObtenerPorID(ctx context.Context, id string) (*rbac.RolDB, error) {
	var m RolModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rbac.ErrRolNoEncontrado
		}
		return nil, err
	}
	return toRolDB(&m), nil
}

func (r *rolRepositorio) Listar(ctx context.Context, spec rbac.EspecificacionRol, pag shareddomain.Paginacion) ([]*rbac.RolDB, error) {
	query := r.db.WithContext(ctx).Model(&RolModel{})

	mapeoColumnas := map[string]string{
		"nombre":    "nombre",
		"esSistema": "es_sistema",
	}

	for _, filtro := range spec.ListaFiltros {
		if !rbac.ColumnasPermitidasRol[filtro.Campo] {
			continue
		}

		columnaDB, ok := mapeoColumnas[filtro.Campo]
		if !ok {
			continue
		}

		switch filtro.Operador {
		case "=":
			query = query.Where(columnaDB+" = ?", filtro.Valor)
		case "!=":
			query = query.Where(columnaDB+" != ?", filtro.Valor)
		case "LIKE":
			query = query.Where(columnaDB+" LIKE ?", filtro.Valor)
		}
	}

	// Filtro por tenant: si TenantID está especificado y no es vacío,
	// mostrar roles del tenant + roles de sistema (tenant_id vacío)
	if spec.TenantID != "" {
		query = query.Where("tenant_id = ? OR tenant_id = '' OR es_sistema = ?", spec.TenantID, true)
	}

	for _, ord := range pag.Ordenaciones {
		if !rbac.ColumnasPermitidasRol[ord.Campo] {
			continue
		}

		columnaDB, ok := mapeoColumnas[ord.Campo]
		if !ok {
			continue
		}

		orden := "ASC"
		if ord.Tipo == shareddomain.DESC {
			orden = "DESC"
		}
		query = query.Order(columnaDB + " " + orden)
	}

	offset := (pag.Pagina - 1) * pag.TamanoPagina
	if pag.Pagina < 1 {
		offset = 0
	}
	if pag.TamanoPagina < 1 {
		pag.TamanoPagina = 10
	}

	var models []RolModel
	result := query.Offset(offset).Limit(pag.TamanoPagina).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	roles := make([]*rbac.RolDB, len(models))
	for i, m := range models {
		roles[i] = toRolDB(&m)
	}
	return roles, nil
}

func (r *rolRepositorio) Crear(ctx context.Context, rol *rbac.RolDB) error {
	m := &RolModel{
		ID:          rol.ID,
		Nombre:      rol.Nombre,
		Descripcion: rol.Descripcion,
		EsSistema:   rol.EsSistema,
		TenantID:    rol.TenantID,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *rolRepositorio) ActualizarDescripcion(ctx context.Context, id, descripcion string) error {
	return r.db.WithContext(ctx).Model(&RolModel{}).
		Where("id = ?", id).
		Update("descripcion", descripcion).Error
}

func (r *rolRepositorio) Eliminar(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete role-permission assignments
		if err := tx.Where("rol_id = ?", id).Delete(&RolPermisoModel{}).Error; err != nil {
			return err
		}
		// Delete user-role assignments (global)
		if err := tx.Where("rol_id = ?", id).Delete(&UsuarioRolModel{}).Error; err != nil {
			return err
		}
		// Delete user-tenant-role assignments
		if err := tx.Where("rol_id = ?", id).Delete(&UsuarioTenantRolModel{}).Error; err != nil {
			return err
		}
		// Delete the role itself
		if err := tx.Delete(&RolModel{}, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}

func toRolDB(m *RolModel) *rbac.RolDB {
	return &rbac.RolDB{
		ID:          m.ID,
		Nombre:      m.Nombre,
		Descripcion: m.Descripcion,
		EsSistema:   m.EsSistema,
		TenantID:    m.TenantID,
	}
}

// --- PermisoRepositorio ---

type permisoRepositorio struct{ db *gorm.DB }

func NewPermisoRepositorio(db *gorm.DB) rbac.PermisoRepositorio {
	return &permisoRepositorio{db: db}
}

func (r *permisoRepositorio) ObtenerPorCodigo(ctx context.Context, codigo string) (*rbac.PermisoDB, error) {
	var m PermisoModel
	if err := r.db.WithContext(ctx).First(&m, "codigo = ?", codigo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toPermisoDB(&m), nil
}

func (r *permisoRepositorio) Listar(ctx context.Context) ([]*rbac.PermisoDB, error) {
	var models []PermisoModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	permisos := make([]*rbac.PermisoDB, len(models))
	for i, m := range models {
		permisos[i] = toPermisoDB(&m)
	}
	return permisos, nil
}

func (r *permisoRepositorio) Crear(ctx context.Context, permiso *rbac.PermisoDB) error {
	m := &PermisoModel{
		ID:          permiso.ID,
		Codigo:      permiso.Codigo,
		Nombre:      permiso.Nombre,
		Descripcion: permiso.Descripcion,
		Modulo:      permiso.Modulo,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *permisoRepositorio) ActualizarNombreDescripcion(ctx context.Context, id, nombre, descripcion string) error {
	return r.db.WithContext(ctx).Model(&PermisoModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"nombre": nombre, "descripcion": descripcion}).Error
}

func (r *permisoRepositorio) ListarPorRol(ctx context.Context, rolID string, tenantID string) ([]*rbac.PermisoDB, error) {
	var models []PermisoModel
	query := r.db.WithContext(ctx).
		Joins("JOIN rol_permisos ON rol_permisos.permiso_id = permisos.id").
		Where("rol_permisos.rol_id = ?", rolID)

	// Si hay tenantID, filtrar por tenant + permisos de sistema
	if tenantID != "" {
		query = query.Where("(rol_permisos.tenant_id = ? OR rol_permisos.tenant_id = ?)", tenantID, rbac.TenantIDSistema)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	permisos := make([]*rbac.PermisoDB, len(models))
	for i, m := range models {
		permisos[i] = toPermisoDB(&m)
	}
	return permisos, nil
}

func toPermisoDB(m *PermisoModel) *rbac.PermisoDB {
	return &rbac.PermisoDB{
		ID:          m.ID,
		Codigo:      m.Codigo,
		Nombre:      m.Nombre,
		Descripcion: m.Descripcion,
		Modulo:      m.Modulo,
	}
}

// --- RolPermisoRepositorio ---

type rolPermisoRepositorio struct{ db *gorm.DB }

func NewRolPermisoRepositorio(db *gorm.DB) rbac.RolPermisoRepositorio {
	return &rolPermisoRepositorio{db: db}
}

func (r *rolPermisoRepositorio) AsignarPermiso(ctx context.Context, rolID, permisoID, tenantID, asignadoPor string) error {
	m := &RolPermisoModel{
		ID:          uuid.New().String(),
		RolID:       rolID,
		PermisoID:   permisoID,
		TenantID:    tenantID,
		AsignadoPor: &asignadoPor,
	}
	// Upsert: ON CONFLICT (rol_id, permiso_id, tenant_id) DO NOTHING
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "rol_id"}, {Name: "permiso_id"}, {Name: "tenant_id"}},
		DoNothing: true,
	}).Create(m).Error
}

func (r *rolPermisoRepositorio) EliminarPermiso(ctx context.Context, rolID, permisoID, tenantID string) error {
	return r.db.WithContext(ctx).
		Where("rol_id = ? AND permiso_id = ? AND tenant_id = ?", rolID, permisoID, tenantID).
		Delete(&RolPermisoModel{}).Error
}

func (r *rolPermisoRepositorio) ListarPorRolYTenant(ctx context.Context, rolID, tenantID string) ([]*rbac.PermisoDB, error) {
	var models []PermisoModel
	if err := r.db.WithContext(ctx).
		Joins("JOIN rol_permisos ON rol_permisos.permiso_id = permisos.id").
		Where("rol_permisos.rol_id = ? AND (rol_permisos.tenant_id = ? OR rol_permisos.tenant_id = ?)", rolID, tenantID, rbac.TenantIDSistema).
		Find(&models).Error; err != nil {
		return nil, err
	}
	permisos := make([]*rbac.PermisoDB, len(models))
	for i, m := range models {
		permisos[i] = toPermisoDB(&m)
	}
	return permisos, nil
}

// --- UsuarioRolRepositorio ---

type usuarioRolRepositorio struct{ db *gorm.DB }

func NewUsuarioRolRepositorio(db *gorm.DB) rbac.UsuarioRolRepositorio {
	return &usuarioRolRepositorio{db: db}
}

func (r *usuarioRolRepositorio) Crear(ctx context.Context, usuarioID, rolID string) error {
	m := &UsuarioRolModel{UsuarioID: usuarioID, RolID: rolID}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *usuarioRolRepositorio) Eliminar(ctx context.Context, usuarioID, rolID string) error {
	return r.db.WithContext(ctx).
		Where("usuario_id = ? AND rol_id = ?", usuarioID, rolID).
		Delete(&UsuarioRolModel{}).Error
}

func (r *usuarioRolRepositorio) ListarRolesPorUsuario(ctx context.Context, usuarioID string) ([]*rbac.RolDB, error) {
	var models []RolModel
	if err := r.db.WithContext(ctx).
		Joins("JOIN usuario_roles ON usuario_roles.rol_id = roles.id").
		Where("usuario_roles.usuario_id = ?", usuarioID).
		Find(&models).Error; err != nil {
		return nil, err
	}
	roles := make([]*rbac.RolDB, len(models))
	for i, m := range models {
		roles[i] = toRolDB(&m)
	}
	return roles, nil
}

func (r *usuarioRolRepositorio) TieneRol(ctx context.Context, usuarioID, rolNombre string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&UsuarioRolModel{}).
		Joins("JOIN roles ON roles.id = usuario_roles.rol_id").
		Where("usuario_roles.usuario_id = ? AND roles.nombre = ?", usuarioID, rolNombre).
		Count(&count).Error
	return count > 0, err
}

// --- UsuarioTenantRolRepositorio ---

type usuarioTenantRolRepositorio struct{ db *gorm.DB }

func NewUsuarioTenantRolRepositorio(db *gorm.DB) rbac.UsuarioTenantRolRepositorio {
	return &usuarioTenantRolRepositorio{db: db}
}

func (r *usuarioTenantRolRepositorio) Crear(ctx context.Context, usuarioID, tenantID, rolID string) error {
	m := &UsuarioTenantRolModel{UsuarioID: usuarioID, TenantID: tenantID, RolID: rolID}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *usuarioTenantRolRepositorio) Eliminar(ctx context.Context, usuarioID, tenantID, rolID string) error {
	return r.db.WithContext(ctx).
		Where("usuario_id = ? AND tenant_id = ? AND rol_id = ?", usuarioID, tenantID, rolID).
		Delete(&UsuarioTenantRolModel{}).Error
}

func (r *usuarioTenantRolRepositorio) ListarRolesPorUsuarioEnTenant(ctx context.Context, usuarioID, tenantID string) ([]*rbac.RolDB, error) {
	var models []RolModel
	if err := r.db.WithContext(ctx).
		Joins("JOIN usuario_tenant_roles ON usuario_tenant_roles.rol_id = roles.id").
		Where("usuario_tenant_roles.usuario_id = ? AND usuario_tenant_roles.tenant_id = ?", usuarioID, tenantID).
		Find(&models).Error; err != nil {
		return nil, err
	}
	roles := make([]*rbac.RolDB, len(models))
	for i, m := range models {
		roles[i] = toRolDB(&m)
	}
	return roles, nil
}

func (r *usuarioTenantRolRepositorio) TieneRolEnTenant(ctx context.Context, usuarioID, tenantID, rolNombre string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&UsuarioTenantRolModel{}).
		Joins("JOIN roles ON roles.id = usuario_tenant_roles.rol_id").
		Where("usuario_tenant_roles.usuario_id = ? AND usuario_tenant_roles.tenant_id = ? AND roles.nombre = ?", usuarioID, tenantID, rolNombre).
		Count(&count).Error
	return count > 0, err
}
