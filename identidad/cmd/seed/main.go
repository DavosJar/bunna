// DEPRECATO: usar cmd/init_superuser. Este script se mantiene solo por
// compatibilidad histórica.
//
// Script para crear un superusuario con cuenta totalmente activa y verificada.
// Uso: go run ./cmd/seed
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Modelos mínimos para inserción directa
type UsuarioModel struct {
	ID                       uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id"`
	Nombre                   string    `gorm:"column:nombre"`
	Apellido                 string    `gorm:"column:apellido"`
	Correo                   string    `gorm:"column:correo"`
	Telefono                 string    `gorm:"column:telefono"`
	Estado                   string    `gorm:"column:estado"`
	EstadoVerificacionCorreo string    `gorm:"column:estado_verificacion_correo"`
	FechaCreacion            time.Time `gorm:"column:fecha_creacion;autoCreateTime:milli"`
	FechaActualizacion       time.Time `gorm:"column:fecha_actualizacion;autoUpdateTime:milli"`
}

func (UsuarioModel) TableName() string { return "usuarios" }

type CredencialesModel struct {
	UsuarioID        string    `gorm:"column:usuario_id;primaryKey;type:varchar(36)"`
	PasswordHash     string    `gorm:"column:password_hash;type:varchar(255)"`
	Activo           bool      `gorm:"column:activo;type:boolean;default:true"`
	CorreoVerificado bool      `gorm:"column:correo_verificado;type:boolean;default:false"`
	IntentosFallidos int       `gorm:"column:intentos_fallidos;type:int;default:0"`
	BloqueadoHasta   time.Time `gorm:"column:bloqueado_hasta;type:timestamp;default:null"`
}

func (CredencialesModel) TableName() string { return "credenciales_usuarios" }

// ─── Configuración del superusuario ───
const (
	superNombre   = "Admin"
	superApellido = "Bunna"
	superCorreo   = "admin@bunna.com"
	superPassword = "admin1234"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s client_encoding=UTF8",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "identidad_user"),
		getEnv("DB_PASSWORD", "identidad_pass_dev"),
		getEnv("DB_NAME", "identidad_db"),
		getEnv("DB_SSLMODE", "disable"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ No se pudo conectar a la BD: %v", err)
	}

	// Verificar si ya existe
	var count int64
	db.Model(&UsuarioModel{}).Where("correo = ?", superCorreo).Count(&count)
	if count > 0 {
		log.Println("⚠️  El superusuario ya existe, no se creará de nuevo.")
		log.Printf("   📧 Correo:     %s", superCorreo)
		log.Printf("   🔑 Contraseña: %s", superPassword)
		os.Exit(0)
	}

	// Hashear password
	hash, err := bcrypt.GenerateFromPassword([]byte(superPassword), 12)
	if err != nil {
		log.Fatalf("❌ Error al hashear password: %v", err)
	}

	userID := uuid.Must(uuid.NewV7())
	ahora := time.Now()

	// Crear usuario con estado ACTIVO y verificado
	usuario := UsuarioModel{
		ID:                       userID,
		Nombre:                   superNombre,
		Apellido:                 superApellido,
		Correo:                   superCorreo,
		Telefono:                 "",
		Estado:                   "ACTIVO",
		EstadoVerificacionCorreo: "VERIFICADO",
		FechaCreacion:            ahora,
		FechaActualizacion:       ahora,
	}

	credenciales := CredencialesModel{
		UsuarioID:        userID.String(),
		PasswordHash:     string(hash),
		Activo:           true,
		CorreoVerificado: true,
		IntentosFallidos: 0,
		BloqueadoHasta:   time.Time{},
	}

	// Transacción para consistencia
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&usuario).Error; err != nil {
			return fmt.Errorf("error al crear usuario: %w", err)
		}
		if err := tx.Create(&credenciales).Error; err != nil {
			return fmt.Errorf("error al crear credenciales: %w", err)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("❌ Error al crear superusuario: %v", err)
	}

	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════╗")
	fmt.Println("║       ☕ Superusuario creado exitosamente        ║")
	fmt.Println("╠══════════════════════════════════════════════════╣")
	fmt.Printf("║  ID:          %s  ║\n", userID.String())
	fmt.Printf("║  Nombre:      %-36s  ║\n", superNombre+" "+superApellido)
	fmt.Printf("║  Correo:      %-36s  ║\n", superCorreo)
	fmt.Printf("║  Contraseña:  %-36s  ║\n", superPassword)
	fmt.Println("║  Estado:      ACTIVO ✅                           ║")
	fmt.Println("║  Verificado:  SÍ ✅                               ║")
	fmt.Println("╚══════════════════════════════════════════════════╝")
	fmt.Println()
}
