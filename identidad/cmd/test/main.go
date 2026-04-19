package main

import (
	"context"
	"fmt"
	"log"

	"github.com/davosjar/bunna/services/identidad/internal/config"
	"github.com/davosjar/bunna/services/identidad/internal/domain/usuario"
	"github.com/davosjar/bunna/services/identidad/internal/registry"
)

func main() {
	// Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Config: %v", err)
	}

	// Inicializar base de datos
	db, err := config.InitDB(cfg.GetDSN())
	if err != nil {
		log.Fatalf("❌ DB: %v", err)
	}

	// Limpiar y migrar
	db.Migrator().DropTable(&usuario.Usuario{})
	if err := config.RunMigrations(db); err != nil {
		log.Fatalf("❌ Migrations: %v", err)
	}

	// Crear registry e inyectar repositorio
	reg := registry.NewRegistry(db)
	repo := reg.UsuarioRepository()
	ctx := context.Background()

	fmt.Println("\n=== CREACIÓN ===")
	fmt.Println()

	// [1] Crear Juan
	fmt.Println("[1] Creando Juan...")
	u1, err := usuario.NuevoUsuario("", "juan@test.com", "Juan", "García", "600111111")
	if err != nil {
		log.Fatalf("❌ Error crear usuario: %v", err)
	}

	creado1, err := repo.Crear(ctx, u1)
	if err != nil {
		log.Fatalf("❌ Error persistir: %v", err)
	}
	fmt.Printf("✓ Usuario creado con ID: %s\n", creado1.ID())

	eventosJuan := creado1.PullEventos()
	fmt.Println("\n📤 COLA DE EVENTOS:")
	for _, evt := range eventosJuan {
		fmt.Printf("   • %s (%s)\n", evt.Nombre, creado1.Correo())
	}

	// [2] Crear María
	fmt.Println("\n[2] Creando María...")
	u2, err := usuario.NuevoUsuario("", "maria@test.com", "María", "López", "600222222")
	if err != nil {
		log.Fatalf("❌ Error crear usuario: %v", err)
	}

	creado2, err := repo.Crear(ctx, u2)
	if err != nil {
		log.Fatalf("❌ Error persistir: %v", err)
	}
	fmt.Printf("✓ Usuario creado con ID: %s\n", creado2.ID())

	eventosMaria := creado2.PullEventos()
	fmt.Println("\n📤 COLA DE EVENTOS:")
	for _, evt := range eventosMaria {
		fmt.Printf("   • %s (%s)\n", evt.Nombre, creado2.Correo())
	}

	// [3] Crear Carlos
	fmt.Println("\n[3] Creando Carlos...")
	u3, err := usuario.NuevoUsuario("", "carlos@test.com", "Carlos", "Ruiz", "600333333")
	if err != nil {
		log.Fatalf("❌ Error crear usuario: %v", err)
	}

	creado3, err := repo.Crear(ctx, u3)
	if err != nil {
		log.Fatalf("❌ Error persistir: %v", err)
	}
	fmt.Printf("✓ Usuario creado con ID: %s\n", creado3.ID())

	eventosCarlos := creado3.PullEventos()
	fmt.Println("\n📤 COLA DE EVENTOS:")
	for _, evt := range eventosCarlos {
		fmt.Printf("   • %s (%s)\n", evt.Nombre, creado3.Correo())
	}

	fmt.Println("\n=== MODIFICACIÓN ===")
	fmt.Println()

	// [4] Activar Juan
	fmt.Println("[4] Activando Juan...")
	err = creado1.Activar()
	if err != nil {
		log.Fatalf("❌ Error activar Juan: %v", err)
	}
	fmt.Printf("✓ Juan activado (estado: %s)\n", creado1.Estado())

	eventosActivarJuan := creado1.PullEventos()
	fmt.Println("\n📤 COLA DE EVENTOS:")
	for _, evt := range eventosActivarJuan {
		fmt.Printf("   • %s\n", evt.Nombre)
	}

	// [5] Activar María
	fmt.Println("\n[5] Activando María...")
	err = creado2.Activar()
	if err != nil {
		log.Fatalf("❌ Error activar María: %v", err)
	}
	fmt.Printf("✓ María activada (estado: %s)\n", creado2.Estado())

	eventosActivarMaria := creado2.PullEventos()
	fmt.Println("\n📤 COLA DE EVENTOS:")
	for _, evt := range eventosActivarMaria {
		fmt.Printf("   • %s\n", evt.Nombre)
	}

	// [6] Bloquear Carlos
	fmt.Println("\n[6] Bloqueando Carlos...")
	err = creado3.Bloquear()
	if err != nil {
		log.Fatalf("❌ Error bloquear Carlos: %v", err)
	}
	fmt.Printf("✓ Carlos bloqueado (estado: %s)\n", creado3.Estado())

	eventosBloquearCarlos := creado3.PullEventos()
	fmt.Println("\n📤 COLA DE EVENTOS:")
	for _, evt := range eventosBloquearCarlos {
		fmt.Printf("   • %s\n", evt.Nombre)
	}

	fmt.Println("\n=== FIN ===")

}
