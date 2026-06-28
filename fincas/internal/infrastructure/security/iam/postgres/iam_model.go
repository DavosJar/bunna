package postgres

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// IamRolPermisosModel almacena la copia local de los permisos de un rol para un tenant.
type IamRolPermisosModel struct {
	RolID    string `gorm:"primaryKey;column:rol_id;type:varchar(36)"`
	TenantID string `gorm:"primaryKey;column:tenant_id;type:varchar(36)"`
	Permisos string `gorm:"column:permisos;type:text"` // JSON array de codigos de permisos
}

func (IamRolPermisosModel) TableName() string { return "iam_rol_permisos" }

type IAMRepositorio struct {
	db *gorm.DB
}

func NewIAMRepositorio(db *gorm.DB) *IAMRepositorio {
	return &IAMRepositorio{db: db}
}

// UpsertPermisos actualiza la lista de permisos de un rol para un tenant.
func (r *IAMRepositorio) UpsertPermisos(ctx context.Context, rolID, tenantID string, permisos []string) error {
	b, err := json.Marshal(permisos)
	if err != nil {
		return err
	}

	model := IamRolPermisosModel{
		RolID:    rolID,
		TenantID: tenantID,
		Permisos: string(b),
	}

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "rol_id"}, {Name: "tenant_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"permisos"}),
	}).Create(&model).Error
}

// ObtenerPermisos devuelve la lista de permisos de un rol para un tenant específico.
// No hace fallback: si no existe la combinación exacta (rol_id, tenant_id), retorna lista vacía.
func (r *IAMRepositorio) ObtenerPermisos(ctx context.Context, rolID, tenantID string) ([]string, error) {
	var model IamRolPermisosModel
	err := r.db.WithContext(ctx).
		Where("rol_id = ? AND tenant_id = ?", rolID, tenantID).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []string{}, nil
		}
		return nil, err
	}

	var permisos []string
	if err := json.Unmarshal([]byte(model.Permisos), &permisos); err != nil {
		return nil, err
	}

	return permisos, nil
}
