package postgres

import "time"

// RolModel representa la tabla roles
type RolModel struct {
	ID                 string    `gorm:"type:varchar(36);primaryKey;column:id"`
	Nombre             string    `gorm:"column:nombre;uniqueIndex"`
	Descripcion        string    `gorm:"column:descripcion"`
	EsSistema          bool      `gorm:"column:es_sistema;default:false"`
	FechaCreacion      time.Time `gorm:"column:created_at;autoCreateTime"`
	FechaActualizacion time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (RolModel) TableName() string { return "roles" }

// PermisoModel representa la tabla permisos
type PermisoModel struct {
	ID            string    `gorm:"type:varchar(36);primaryKey;column:id"`
	Codigo        string    `gorm:"column:codigo;uniqueIndex"`
	Nombre        string    `gorm:"column:nombre"`
	Descripcion   string    `gorm:"column:descripcion"`
	Modulo        string    `gorm:"column:modulo"`
	FechaCreacion time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (PermisoModel) TableName() string { return "permisos" }

// RolPermisoModel representa la tabla rol_permisos
type RolPermisoModel struct {
	ID            string    `gorm:"type:varchar(36);primaryKey;column:id"`
	RolID         string    `gorm:"type:varchar(36);not null;column:rol_id;uniqueIndex:idx_rp_rol_perm_tenant"`
	PermisoID     string    `gorm:"type:varchar(36);not null;column:permiso_id;uniqueIndex:idx_rp_rol_perm_tenant"`
	TenantID      string    `gorm:"type:varchar(36);not null;column:tenant_id;default:00000000-0000-0000-0000-000000000000;uniqueIndex:idx_rp_rol_perm_tenant"`
	AsignadoPor   *string   `gorm:"type:varchar(36);column:asignado_por"`
	FechaCreacion time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (RolPermisoModel) TableName() string { return "rol_permisos" }

// UsuarioRolModel representa la tabla usuario_roles (roles globales)
type UsuarioRolModel struct {
	UsuarioID     string    `gorm:"type:varchar(36);primaryKey;column:usuario_id"`
	RolID         string    `gorm:"type:varchar(36);primaryKey;column:rol_id"`
	FechaCreacion time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (UsuarioRolModel) TableName() string { return "usuario_roles" }

// UsuarioTenantRolModel representa la tabla usuario_tenant_roles
type UsuarioTenantRolModel struct {
	UsuarioID     string    `gorm:"type:varchar(36);primaryKey;column:usuario_id"`
	TenantID      string    `gorm:"type:varchar(36);primaryKey;column:tenant_id"`
	RolID         string    `gorm:"type:varchar(36);primaryKey;column:rol_id"`
	FechaCreacion time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (UsuarioTenantRolModel) TableName() string { return "usuario_tenant_roles" }
