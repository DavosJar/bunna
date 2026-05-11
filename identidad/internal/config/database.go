package config

import (
	"fmt"
	"log"

	seguridad_postgres "github.com/davosjar/bunna/services/identidad/internal/seguridad/infrastructure/persistence/postgres"
	sesiones_postgres "github.com/davosjar/bunna/services/identidad/internal/sesiones/infrastructure/persistence/postgres"
	usuarios_postgres "github.com/davosjar/bunna/services/identidad/internal/usuarios/infrastructure/persistence/postgres"
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

	// Índice único para emails (case-insensitive via LOWER).
	// CREATE UNIQUE INDEX IF NOT EXISTS es soportado desde PostgreSQL 9.5+.
	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_usuarios_correo_unique ON usuarios (correo)")

	return nil
}