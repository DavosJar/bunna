package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/davosjar/bunna/services/identidad/internal/config"
	seguridad_postgres "github.com/davosjar/bunna/services/identidad/internal/seguridad/infrastructure/persistence/postgres"
	bcrypt_lib "github.com/davosjar/bunna/services/identidad/internal/seguridad/infrastructure/security/bcrypt"
	sesiones_login "github.com/davosjar/bunna/services/identidad/internal/sesiones/application/usecases/login"
	sesiones_postgres "github.com/davosjar/bunna/services/identidad/internal/sesiones/infrastructure/persistence/postgres"
	sesiones_jwt "github.com/davosjar/bunna/services/identidad/internal/sesiones/infrastructure/security/jwt"
	shared_idgenerator "github.com/davosjar/bunna/services/identidad/internal/shared/infrastructure/idgenerator"
	usuarios_postgres "github.com/davosjar/bunna/services/identidad/internal/usuarios/infrastructure/persistence/postgres"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	// ── 1. Set up env vars needed by config ──
	os.Setenv("JWT_SECRET", "supersecreto_dev_cambiar_en_prod")

	// ── 2. Build DSN ──
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s client_encoding=UTF8",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "identidad_user"),
		getEnv("DB_PASSWORD", "identidad_pass_dev"),
		getEnv("DB_NAME", "identidad_db"),
		getEnv("DB_SSLMODE", "disable"),
	)

	fmt.Println("📦 Connecting to database and running migrations...")

	// ── 3. Connect to DB and run migrations ──
	db, err := config.InitDB(dsn)
	if err != nil {
		log.Fatalf("❌ Failed to init DB: %v", err)
	}
	fmt.Println("✅ Database connected and migrations applied")

	// ── 4. Clean up previous test data ──
	db.Exec("DELETE FROM credenciales_usuarios WHERE usuario_id IN (SELECT id::text FROM usuarios WHERE correo = 'test_block@test.com')")
	db.Exec("DELETE FROM usuarios WHERE correo = 'test_block@test.com'")

	// ── 5. Create test user ──
	correctPassword := "correct_password"
	hash, err := bcrypt.GenerateFromPassword([]byte(correctPassword), 12)
	if err != nil {
		log.Fatalf("❌ Failed to hash password: %v", err)
	}

	userID := uuid.Must(uuid.NewV7())
	ahora := time.Now()

	usuario := usuarios_postgres.UsuarioModel{
		ID:                       userID,
		Nombre:                   "Test",
		Apellido:                 "Block",
		Correo:                   "test_block@test.com",
		Telefono:                 "",
		Estado:                   "ACTIVO",
		EstadoVerificacionCorreo: "VERIFICADO",
		FechaCreacion:            ahora,
		FechaActualizacion:       ahora,
	}

	credenciales := seguridad_postgres.CredencialesModel{
		UsuarioID:        userID.String(),
		PasswordHash:     string(hash),
		Activo:           true,
		CorreoVerificado: true,
		IntentosFallidos: 0,
		BloqueadoHasta:   time.Time{},
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&usuario).Error; err != nil {
			return fmt.Errorf("error creating user: %w", err)
		}
		if err := tx.Create(&credenciales).Error; err != nil {
			return fmt.Errorf("error creating credenciales: %w", err)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("❌ Failed to create test user: %v", err)
	}
	fmt.Println("✅ Test user created: test_block@test.com")

	// ── 6. Wire up real dependencies ──
	encriptacion := bcrypt_lib.NewBcryptEncriptacion(12)
	generadorID := shared_idgenerator.NewUUIDv7Generator()

	usuarioRepo := usuarios_postgres.NewUsuarioRepositorio(db)
	credencialesRepo := seguridad_postgres.NewCredencialesRepositorio(db)
	sesionRepo := sesiones_postgres.NewSesionRepositorio(db)

	tokenSvc := sesiones_jwt.NewJWTTokenServicio(sesiones_jwt.ConfigJWT{
		Secret:            "supersecreto_dev_cambiar_en_prod",
		Issuer:            "ServicioIdentidad",
		ExpiracionAccess:  15 * time.Minute,
		ExpiracionRefresh: 24 * time.Hour,
	})

	sesionUoW := sesiones_postgres.NewSesionUnitOfWork(
		db,
		sesionRepo,
		credencialesRepo,
		usuarioRepo,
		encriptacion,
		tokenSvc,
		generadorID,
	)

	// Without IP bloqueo / rate limiter services (IPOrigen="" skips them)
	useCase := sesiones_login.NewIniciarSesionCasoDeUso(sesionUoW, nil, nil, sesiones_login.ConfigLogin{
		CuentaMaxIntentos:     5,
		CuentaBloqueoDuracion: 15 * time.Minute,
	})

	// ── 7. Call Ejecutar with wrong password 6 times ──
	ctx := context.Background()
	wrongPassword := "wrong_password"

	fmt.Println("\n═══════════════════════════════════════════")
	fmt.Println("  Calling Ejecutar with wrong password ×6")
	fmt.Println("═══════════════════════════════════════════")

	var lastErr error
	for i := 1; i <= 6; i++ {
		cmd := sesiones_login.ComandoIniciarSesion{
			Email:    "test_block@test.com",
			Password: wrongPassword,
			IPOrigen: "",
		}

		_, err := useCase.Ejecutar(ctx, cmd)
		lastErr = err
		if err != nil {
			fmt.Printf("  Attempt %d: ❌ %v\n", i, err)
		} else {
			fmt.Printf("  Attempt %d: ✅ SUCCESS (unexpected)\n", i)
		}
	}

	// ── 8. Assertions ──
	fmt.Println("\n═══════════════════════════════════════════")
	fmt.Println("  Assertions")
	fmt.Println("═══════════════════════════════════════════")

	// Assert 1: The 6th call error should be ErrCuentaBloqueada
	fmt.Printf("\n🔍 6th call error: %v\n", lastErr)
	if errors.Is(lastErr, sesiones_login.ErrCuentaBloqueada) {
		fmt.Println("   ✅ Error is ErrCuentaBloqueada (cuenta temporalmente bloqueada)")
	} else {
		fmt.Printf("   ❌ Expected ErrCuentaBloqueada, got: %v\n", lastErr)
	}

	// Assert 2: Database state
	var finalCred seguridad_postgres.CredencialesModel
	result := db.Where("usuario_id = ?", userID.String()).First(&finalCred)
	if result.Error != nil {
		log.Fatalf("❌ Failed to query credenciales: %v", result.Error)
	}

	fmt.Printf("\n📊 Database state:\n")
	fmt.Printf("   intentos_fallidos = %d (expected 5)\n", finalCred.IntentosFallidos)
	fmt.Printf("   bloqueado_hasta   = %v\n", finalCred.BloqueadoHasta)

	ahoraVerif := time.Now()
	if finalCred.IntentosFallidos == 5 {
		fmt.Println("   ✅ intentos_fallidos = 5")
	} else {
		fmt.Printf("   ❌ Expected intentos_fallidos = 5, got %d\n", finalCred.IntentosFallidos)
	}

	if !finalCred.BloqueadoHasta.IsZero() && finalCred.BloqueadoHasta.After(ahoraVerif) {
		fmt.Println("   ✅ bloqueado_hasta is set and in the future")
	} else if !finalCred.BloqueadoHasta.IsZero() {
		fmt.Println("   ❌ bloqueado_hasta is in the past")
	} else {
		fmt.Println("   ❌ bloqueado_hasta is zero (not blocked)")
	}

	// ── Summary ──
	allPassed := errors.Is(lastErr, sesiones_login.ErrCuentaBloqueada) &&
		finalCred.IntentosFallidos == 5 &&
		!finalCred.BloqueadoHasta.IsZero() &&
		finalCred.BloqueadoHasta.After(ahoraVerif)

	fmt.Println("\n═══════════════════════════════════════════")
	if allPassed {
		fmt.Println("  ✅✅✅ ALL ASSERTIONS PASSED")
	} else {
		fmt.Println("  ❌❌❌ SOME ASSERTIONS FAILED")
	}
	fmt.Println("═══════════════════════════════════════════")
}
