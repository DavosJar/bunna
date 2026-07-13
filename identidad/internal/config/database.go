package config

import (
	"fmt"
	"log"

	invitaciones_postgres "github.com/davosjar/bunna/services/identidad/internal/invitaciones/infrastructure/persistence/postgres"
	rbac_postgres "github.com/davosjar/bunna/services/identidad/internal/rbac/infrastructure/persistence/postgres"
	recuperacion_postgres "github.com/davosjar/bunna/services/identidad/internal/recuperacion/infrastructure/persistence/postgres"
	seguridad_postgres "github.com/davosjar/bunna/services/identidad/internal/seguridad/infrastructure/persistence/postgres"
	sesiones_postgres "github.com/davosjar/bunna/services/identidad/internal/sesiones/infrastructure/persistence/postgres"
	"github.com/davosjar/bunna/services/identidad/internal/shared/infrastructure/outbox"
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
		return nil, fmt.Errorf("error al conectar con la base de datos: %w", err)
	}

	log.Println("Database connection established")

	if err := RunMigrations(db); err != nil {
		return nil, fmt.Errorf("error al ejecutar migraciones: %w", err)
	}

	return db, nil
}

// RunMigrations ejecuta todas las migraciones automáticas de GORM.
// El orden importa: usuarios antes que credenciales (FK), sesiones al final.
func RunMigrations(db *gorm.DB) error {
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return fmt.Errorf("error al crear extensión uuid-ossp: %w", err)
	}

	if err := db.AutoMigrate(&usuarios_postgres.UsuarioModel{}); err != nil {
		return fmt.Errorf("error al migrar usuarios: %w", err)
	}

	if err := db.AutoMigrate(&seguridad_postgres.CredencialesModel{}); err != nil {
		return fmt.Errorf("error al migrar credenciales: %w", err)
	}

	if err := db.AutoMigrate(&sesiones_postgres.SesionModel{}); err != nil {
		return fmt.Errorf("error al migrar sesiones: %w", err)
	}

	if err := db.AutoMigrate(&seguridad_postgres.IntentoIPModel{}); err != nil {
		return fmt.Errorf("error al migrar intentos por IP: %w", err)
	}

	if err := db.AutoMigrate(&seguridad_postgres.RateLimitIPModel{}); err != nil {
		return fmt.Errorf("error al migrar rate limit IP: %w", err)
	}

	if err := db.AutoMigrate(&rbac_postgres.RolModel{}); err != nil {
		return fmt.Errorf("error al migrar roles: %w", err)
	}

	// Limpiar roles huérfanos de sistema (creados en registros fallidos, no toca roles personalizados)
	db.Exec(`DELETE FROM roles WHERE id IN (
		SELECT r.id FROM roles r
		LEFT JOIN usuario_tenant_roles utr ON r.id = utr.rol_id
		WHERE utr.rol_id IS NULL AND r.es_sistema = true
	)`)

	// Migrar unique index de roles: de (nombre) a (nombre, tenant_id)
	if err := db.Exec(`DROP INDEX IF EXISTS idx_roles_nombre`).Error; err != nil {
		return fmt.Errorf("error al dropear old index idx_roles_nombre: %w", err)
	}
	if err := db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_nombre_tenant ON roles (nombre, tenant_id)`).Error; err != nil {
		return fmt.Errorf("error al crear unique index idx_roles_nombre_tenant: %w", err)
	}

	if err := db.AutoMigrate(&rbac_postgres.PermisoModel{}); err != nil {
		return fmt.Errorf("error al migrar permisos: %w", err)
	}
	if err := db.AutoMigrate(&rbac_postgres.RolPermisoModel{}); err != nil {
		return fmt.Errorf("error al migrar rol_permisos: %w", err)
	}
	if err := db.AutoMigrate(&rbac_postgres.UsuarioRolModel{}); err != nil {
		return fmt.Errorf("error al migrar usuario_roles: %w", err)
	}
	if err := db.AutoMigrate(&rbac_postgres.UsuarioTenantRolModel{}); err != nil {
		return fmt.Errorf("error al migrar usuario_tenant_roles: %w", err)
	}

	if err := db.AutoMigrate(&tenant_postgres.TenantModel{}); err != nil {
		return fmt.Errorf("error al migrar tenants: %w", err)
	}
	if err := db.AutoMigrate(&tenant_postgres.MembresiaModel{}); err != nil {
		return fmt.Errorf("error al migrar membresías: %w", err)
	}

	if err := db.AutoMigrate(&verificacion_postgres.VerificacionUsuarioModel{}); err != nil {
		return fmt.Errorf("error al migrar columnas de verificación: %w", err)
	}

	if err := db.AutoMigrate(&recuperacion_postgres.TokenRecuperacionModel{}); err != nil {
		return fmt.Errorf("error al migrar tokens de recuperación: %w", err)
	}

	if err := db.AutoMigrate(&invitaciones_postgres.InvitacionModel{}); err != nil {
		return fmt.Errorf("error al migrar invitaciones: %w", err)
	}

	// Migrar event_outbox (Outbox Pattern)
	if err := db.AutoMigrate(&outbox.EventoOutbox{}); err != nil {
		return fmt.Errorf("error al migrar event_outbox: %w", err)
	}

	// Migrar rol_permisos: agregar id, tenant_id, asignado_por y unique index
	if err := db.Exec(`
		ALTER TABLE rol_permisos ADD COLUMN IF NOT EXISTS id varchar(36);
		ALTER TABLE rol_permisos ADD COLUMN IF NOT EXISTS tenant_id varchar(36) NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000';
		ALTER TABLE rol_permisos ADD COLUMN IF NOT EXISTS asignado_por varchar(36) NULL;
	`).Error; err != nil {
		return fmt.Errorf("error al migrar columnas de rol_permisos: %w", err)
	}

	// Poblar IDs para filas existentes (sistema, ya que no tenían tenant_id antes)
	if err := db.Exec(`
		UPDATE rol_permisos SET id = gen_random_uuid()::text WHERE id IS NULL;
	`).Error; err != nil {
		return fmt.Errorf("error al poblar IDs de rol_permisos: %w", err)
	}

	// Hacer id NOT NULL y PK
	if err := db.Exec(`
		ALTER TABLE rol_permisos ALTER COLUMN id SET NOT NULL;
		ALTER TABLE rol_permisos DROP CONSTRAINT IF EXISTS rol_permisos_pkey;
		ALTER TABLE rol_permisos ADD PRIMARY KEY (id);
	`).Error; err != nil {
		return fmt.Errorf("error al establecer PK de rol_permisos: %w", err)
	}

	// Unique index
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_rp_rol_perm_tenant ON rol_permisos (rol_id, permiso_id, tenant_id);
	`).Error; err != nil {
		return fmt.Errorf("error al crear índice único en rol_permisos: %w", err)
	}

	// Índice único para emails (case-insensitive via LOWER).
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_usuarios_correo_unique ON usuarios (correo)")

	return nil
}
