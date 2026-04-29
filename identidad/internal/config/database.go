package config

import (
	"fmt"
	"log"

	seguridad_postgres "github.com/davosjar/bunna/services/identidad/internal/seguridad/infrastructure/persistence/postgres"
	usuarios_postgres "github.com/davosjar/bunna/services/identidad/internal/usuarios/infrastructure/persistence/postgres"
	postgresdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

func RunMigrations(db *gorm.DB) error {
	// Enable uuid-ossp extension for PostgreSQL 18
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return fmt.Errorf("failed to create uuid-ossp extension: %w", err)
	}

	// Migrate usuario
	if err := db.AutoMigrate(&usuarios_postgres.UsuarioModel{}); err != nil {
		return fmt.Errorf("failed to migrate usuario: %w", err)
	}

	// Migrate credenciales
	if err := db.AutoMigrate(&seguridad_postgres.CredencialesModel{}); err != nil {
		return fmt.Errorf("failed to migrate credenciales: %w", err)
	}

	return nil
}
