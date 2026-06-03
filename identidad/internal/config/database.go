package config

import (
	"fmt"
	"log"

	rbac_postgres "github.com/davosjar/bunna/services/identidad/internal/rbac/infrastructure/persistence/postgres"
	recuperacion_postgres "github.com/davosjar/bunna/services/identidad/internal/recuperacion/infrastructure/persistence/postgres"
	seguridad_postgres "github.com/davosjar/bunna/services/identidad/internal/seguridad/infrastructure/persistence/postgres"
	sesiones_postgres "github.com/davosjar/bunna/services/identidad/internal/sesiones/infrastructure/persistence/postgres"
	tenant_postgres "github.com/davosjar/bunna/services/identidad/internal/tenants/infrastructure/persistence/postgres"
	usuarios_postgres "github.com/davosjar/bunna/services/identidad/internal/usuarios/infrastructure/persistence/postgres"
	verificacion_postgres "github.com/davosjar/bunna/services/identidad/internal/verificacion/infrastructure/persistence/postgres"
	postgresdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB inicializa la conexión a PostgreSQL y ejecuta las migraciones automáticas.
func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgresdriver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")

	if err := RunMigrations(db); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return db, nil
}

// RunMigrations ejecuta todas las migraciones automáticas de GORM.
// El orden importa: usuarios antes que credenciales (FK), sesiones al final.
func RunMigrations(db *gorm.DB) error {
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return fmt.Errorf("failed to create uuid-ossp extension: %w", err)
	}

	if err := db.AutoMigrate(&usuarios_postgres.UsuarioModel{}); err != nil {
		return fmt.Errorf("failed to migrate usuario: %w", err)
	}

	if err := db.AutoMigrate(&seguridad_postgres.CredencialesModel{}); err != nil {
		return fmt.Errorf("failed to migrate credenciales: %w", err)
	}

	if err := db.AutoMigrate(&sesiones_postgres.SesionModel{}); err != nil {
		return fmt.Errorf("failed to migrate sesiones: %w", err)
	}

	if err := db.AutoMigrate(&seguridad_postgres.IntentoIPModel{}); err != nil {
		return fmt.Errorf("failed to migrate intentos_por_ip: %w", err)
	}

	if err := db.AutoMigrate(&seguridad_postgres.RateLimitIPModel{}); err != nil {
		return fmt.Errorf("failed to migrate rate_limit_ip: %w", err)
	}

	if err := db.AutoMigrate(&rbac_postgres.RolModel{}); err != nil {
		return fmt.Errorf("failed to migrate roles: %w", err)
	}
	if err := db.AutoMigrate(&rbac_postgres.PermisoModel{}); err != nil {
		return fmt.Errorf("failed to migrate permisos: %w", err)
	}
	if err := db.AutoMigrate(&rbac_postgres.RolPermisoModel{}); err != nil {
		return fmt.Errorf("failed to migrate rol_permisos: %w", err)
	}
	if err := db.AutoMigrate(&rbac_postgres.UsuarioRolModel{}); err != nil {
		return fmt.Errorf("failed to migrate usuario_roles: %w", err)
	}
	if err := db.AutoMigrate(&rbac_postgres.UsuarioTenantRolModel{}); err != nil {
		return fmt.Errorf("failed to migrate usuario_tenant_roles: %w", err)
	}

	if err := db.AutoMigrate(&tenant_postgres.TenantModel{}); err != nil {
		return fmt.Errorf("failed to migrate tenants: %w", err)
	}
	if err := db.AutoMigrate(&tenant_postgres.MembresiaModel{}); err != nil {
		return fmt.Errorf("failed to migrate usuario_tenants: %w", err)
	}

	if err := db.AutoMigrate(&verificacion_postgres.VerificacionUsuarioModel{}); err != nil {
		return fmt.Errorf("failed to migrate verificacion columns: %w", err)
	}

	if err := db.AutoMigrate(&recuperacion_postgres.TokenRecuperacionModel{}); err != nil {
		return fmt.Errorf("failed to migrate tokens_recuperacion: %w", err)
	}

	// Migrar rol_permisos: agregar id, tenant_id, asignado_por y unique index
	if err := db.Exec(`
		ALTER TABLE rol_permisos ADD COLUMN IF NOT EXISTS id varchar(36);
		ALTER TABLE rol_permisos ADD COLUMN IF NOT EXISTS tenant_id varchar(36) NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000';
		ALTER TABLE rol_permisos ADD COLUMN IF NOT EXISTS asignado_por varchar(36) NULL;
	`).Error; err != nil {
		return fmt.Errorf("failed to migrate rol_permisos columns: %w", err)
	}

	// Poblar IDs para filas existentes (sistema, ya que no tenían tenant_id antes)
	if err := db.Exec(`
		UPDATE rol_permisos SET id = gen_random_uuid()::text WHERE id IS NULL;
	`).Error; err != nil {
		return fmt.Errorf("failed to populate rol_permisos ids: %w", err)
	}

	// Hacer id NOT NULL y PK
	if err := db.Exec(`
		ALTER TABLE rol_permisos ALTER COLUMN id SET NOT NULL;
		ALTER TABLE rol_permisos DROP CONSTRAINT IF EXISTS rol_permisos_pkey;
		ALTER TABLE rol_permisos ADD PRIMARY KEY (id);
	`).Error; err != nil {
		return fmt.Errorf("failed to set rol_permisos primary key: %w", err)
	}

	// Unique index
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_rp_rol_perm_tenant ON rol_permisos (rol_id, permiso_id, tenant_id);
	`).Error; err != nil {
		return fmt.Errorf("failed to create unique index on rol_permisos: %w", err)
	}

	// Índice único para emails (case-insensitive via LOWER).
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_usuarios_correo_unique ON usuarios (correo)")

	return nil
}
