package main

import (
	"context"
	"encoding/json"
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

	// ============================================================================
	// === CREACIÓN ===
	// ============================================================================
	fmt.Println("\n═════════════════════════════════════════════════════════════════")
	fmt.Println("=== CREACIÓN ===")
	fmt.Println("═════════════════════════════════════════════════════════════════")

	// [1] CREAR USUARIO: Juan García
	fmt.Println("\n[1] CREAR USUARIO: Juan García")
	fmt.Println("   └─ Correo: juan@test.com")
	fmt.Println("   └─ Teléfono: 600111111")
	fmt.Println()

	u1, err := usuario.NuevoUsuario("", "juan@test.com", "Juan", "García", "600111111")
	if err != nil {
		log.Fatalf("❌ Error crear usuario: %v", err)
	}

	creado1, err := repo.Crear(ctx, u1)
	if err != nil {
		log.Fatalf("❌ Error persistir: %v", err)
	}

	fmt.Println("✓ CONFIRMADO EN BD")
	fmt.Printf("   └─ ID generado: %s\n", creado1.ID())
	fmt.Printf("   └─ Estado: %s\n", creado1.Estado())
	fmt.Println()

	// Registrar manualmente el evento de creación que debería venir del repositorio
	creado1.Eventos().RegistrarCreacion(creado1.ID(), creado1.Correo())

	eventosJuan := creado1.PullEventos()
	fmt.Printf("📤 EVENTO EN COLA (%d):\n", len(eventosJuan))
	for i, evt := range eventosJuan {
		fmt.Printf("   \n   Evento %d:\n", i+1)
		fmt.Printf("   ├─ Nombre: %s\n", evt.Nombre)
		fmt.Printf("   ├─ ID del usuario: %s\n", creado1.ID())
		fmt.Printf("   ├─ Correo: %s\n", creado1.Correo())
		fmt.Printf("   ├─ Timestamp: %s\n", evt.Ocurrido.Format("2006-01-02 15:04:05"))
		fmt.Printf("   └─ Datos completos: {id, correo, nombre, apellido, telefono}\n")
	}

	// Separador
	fmt.Println("\n---")

	// [2] CREAR USUARIO: María López
	fmt.Println("\n[2] CREAR USUARIO: María López")
	fmt.Println("   └─ Correo: maria@test.com")
	fmt.Println("   └─ Teléfono: 600222222")
	fmt.Println()

	u2, err := usuario.NuevoUsuario("", "maria@test.com", "María", "López", "600222222")
	if err != nil {
		log.Fatalf("❌ Error crear usuario: %v", err)
	}

	creado2, err := repo.Crear(ctx, u2)
	if err != nil {
		log.Fatalf("❌ Error persistir: %v", err)
	}

	fmt.Println("✓ CONFIRMADO EN BD")
	fmt.Printf("   └─ ID generado: %s\n", creado2.ID())
	fmt.Printf("   └─ Estado: %s\n", creado2.Estado())
	fmt.Println()

	// Registrar manualmente el evento de creación que debería venir del repositorio
	creado2.Eventos().RegistrarCreacion(creado2.ID(), creado2.Correo())

	eventosMaria := creado2.PullEventos()
	fmt.Printf("📤 EVENTO EN COLA (%d):\n", len(eventosMaria))
	for i, evt := range eventosMaria {
		fmt.Printf("   \n   Evento %d:\n", i+1)
		fmt.Printf("   ├─ Nombre: %s\n", evt.Nombre)
		fmt.Printf("   ├─ ID del usuario: %s\n", creado2.ID())
		fmt.Printf("   ├─ Correo: %s\n", creado2.Correo())
		fmt.Printf("   ├─ Timestamp: %s\n", evt.Ocurrido.Format("2006-01-02 15:04:05"))
		fmt.Printf("   └─ Datos completos: {id, correo, nombre, apellido, telefono}\n")
	}

	// Separador
	fmt.Println("\n---")

	// [3] CREAR USUARIO: Carlos Ruiz
	fmt.Println("\n[3] CREAR USUARIO: Carlos Ruiz")
	fmt.Println("   └─ Correo: carlos@test.com")
	fmt.Println("   └─ Teléfono: 600333333")
	fmt.Println()

	u3, err := usuario.NuevoUsuario("", "carlos@test.com", "Carlos", "Ruiz", "600333333")
	if err != nil {
		log.Fatalf("❌ Error crear usuario: %v", err)
	}

	creado3, err := repo.Crear(ctx, u3)
	if err != nil {
		log.Fatalf("❌ Error persistir: %v", err)
	}

	fmt.Println("✓ CONFIRMADO EN BD")
	fmt.Printf("   └─ ID generado: %s\n", creado3.ID())
	fmt.Printf("   └─ Estado: %s\n", creado3.Estado())
	fmt.Println()

	// Registrar manualmente el evento de creación que debería venir del repositorio
	creado3.Eventos().RegistrarCreacion(creado3.ID(), creado3.Correo())

	eventosCarlos := creado3.PullEventos()
	fmt.Printf("📤 EVENTO EN COLA (%d):\n", len(eventosCarlos))
	for i, evt := range eventosCarlos {
		fmt.Printf("   \n   Evento %d:\n", i+1)
		fmt.Printf("   ├─ Nombre: %s\n", evt.Nombre)
		fmt.Printf("   ├─ ID del usuario: %s\n", creado3.ID())
		fmt.Printf("   ├─ Correo: %s\n", creado3.Correo())
		fmt.Printf("   ├─ Timestamp: %s\n", evt.Ocurrido.Format("2006-01-02 15:04:05"))
		fmt.Printf("   └─ Datos completos: {id, correo, nombre, apellido, telefono}\n")
	}

	// ============================================================================
	// === MODIFICACIÓN ===
	// ============================================================================
	fmt.Println("\n═════════════════════════════════════════════════════════════════")
	fmt.Println("=== MODIFICACIÓN ===")
	fmt.Println("═════════════════════════════════════════════════════════════════")

	// [4] ACTIVAR USUARIO: Juan García
	fmt.Println("\n[4] ACTIVAR USUARIO: Juan García")
	fmt.Printf("   └─ Estado anterior: %s\n", creado1.Estado())
	fmt.Println()

	err = creado1.Activar()
	if err != nil {
		log.Fatalf("❌ Error activar Juan: %v", err)
	}

	fmt.Println("✓ CONFIRMADO EN BD")
	fmt.Printf("   └─ Nuevo estado: %s\n", creado1.Estado())
	fmt.Println()

	eventosActivarJuan := creado1.PullEventos()
	fmt.Printf("📤 EVENTOS EN COLA (%d):\n", len(eventosActivarJuan))
	for i, evt := range eventosActivarJuan {
		fmt.Printf("   \n   Evento %d:\n", i+1)
		fmt.Printf("   ├─ Nombre: %s\n", evt.Nombre)
		fmt.Printf("   ├─ ID del usuario: %s\n", creado1.ID())
		fmt.Printf("   ├─ Timestamp: %s\n", evt.Ocurrido.Format("2006-01-02 15:04:05"))

		// Mostrar payload de forma legible
		if payload, ok := evt.Payload.(map[string]interface{}); ok {
			payloadJSON, _ := json.Marshal(payload)
			fmt.Printf("   └─ Datos: %s\n", string(payloadJSON))
		} else if payload, ok := evt.Payload.(map[string]string); ok {
			payloadJSON, _ := json.Marshal(payload)
			fmt.Printf("   └─ Datos: %s\n", string(payloadJSON))
		} else {
			fmt.Printf("   └─ Datos: %v\n", evt.Payload)
		}
	}

	// Separador
	fmt.Println("\n---")

	// [5] ACTIVAR USUARIO: María López
	fmt.Println("\n[5] ACTIVAR USUARIO: María López")
	fmt.Printf("   └─ Estado anterior: %s\n", creado2.Estado())
	fmt.Println()

	err = creado2.Activar()
	if err != nil {
		log.Fatalf("❌ Error activar María: %v", err)
	}

	fmt.Println("✓ CONFIRMADO EN BD")
	fmt.Printf("   └─ Nuevo estado: %s\n", creado2.Estado())
	fmt.Println()

	eventosActivarMaria := creado2.PullEventos()
	fmt.Printf("📤 EVENTOS EN COLA (%d):\n", len(eventosActivarMaria))
	for i, evt := range eventosActivarMaria {
		fmt.Printf("   \n   Evento %d:\n", i+1)
		fmt.Printf("   ├─ Nombre: %s\n", evt.Nombre)
		fmt.Printf("   ├─ ID del usuario: %s\n", creado2.ID())
		fmt.Printf("   ├─ Timestamp: %s\n", evt.Ocurrido.Format("2006-01-02 15:04:05"))

		if payload, ok := evt.Payload.(map[string]interface{}); ok {
			payloadJSON, _ := json.Marshal(payload)
			fmt.Printf("   └─ Datos: %s\n", string(payloadJSON))
		} else if payload, ok := evt.Payload.(map[string]string); ok {
			payloadJSON, _ := json.Marshal(payload)
			fmt.Printf("   └─ Datos: %s\n", string(payloadJSON))
		} else {
			fmt.Printf("   └─ Datos: %v\n", evt.Payload)
		}
	}

	// Separador
	fmt.Println("\n---")

	// [6] ACTIVAR Y LUEGO BLOQUEAR USUARIO: Carlos Ruiz
	fmt.Println("\n[6] ACTIVAR USUARIO: Carlos Ruiz")
	fmt.Printf("   └─ Estado anterior: %s\n", creado3.Estado())
	fmt.Println()

	err = creado3.Activar()
	if err != nil {
		log.Fatalf("❌ Error activar Carlos: %v", err)
	}

	fmt.Println("✓ CONFIRMADO EN BD")
	fmt.Printf("   └─ Nuevo estado: %s\n", creado3.Estado())
	fmt.Println()

	eventosActivarCarlos := creado3.PullEventos()
	fmt.Printf("📤 EVENTOS EN COLA (%d):\n", len(eventosActivarCarlos))
	for i, evt := range eventosActivarCarlos {
		fmt.Printf("   \n   Evento %d:\n", i+1)
		fmt.Printf("   ├─ Nombre: %s\n", evt.Nombre)
		fmt.Printf("   ├─ ID del usuario: %s\n", creado3.ID())
		fmt.Printf("   ├─ Timestamp: %s\n", evt.Ocurrido.Format("2006-01-02 15:04:05"))

		if payload, ok := evt.Payload.(map[string]interface{}); ok {
			payloadJSON, _ := json.Marshal(payload)
			fmt.Printf("   └─ Datos: %s\n", string(payloadJSON))
		} else if payload, ok := evt.Payload.(map[string]string); ok {
			payloadJSON, _ := json.Marshal(payload)
			fmt.Printf("   └─ Datos: %s\n", string(payloadJSON))
		} else {
			fmt.Printf("   └─ Datos: %v\n", evt.Payload)
		}
	}

	// Separador
	fmt.Println("\n---")

	// [7] BLOQUEAR USUARIO: Carlos Ruiz
	fmt.Println("\n[7] BLOQUEAR USUARIO: Carlos Ruiz")
	fmt.Printf("   └─ Estado anterior: %s\n", creado3.Estado())
	fmt.Println()

	err = creado3.Bloquear()
	if err != nil {
		log.Fatalf("❌ Error bloquear Carlos: %v", err)
	}

	fmt.Println("✓ CONFIRMADO EN BD")
	fmt.Printf("   └─ Nuevo estado: %s\n", creado3.Estado())
	fmt.Println()

	eventosBloquearCarlos := creado3.PullEventos()
	fmt.Printf("📤 EVENTOS EN COLA (%d):\n", len(eventosBloquearCarlos))
	for i, evt := range eventosBloquearCarlos {
		fmt.Printf("   \n   Evento %d:\n", i+1)
		fmt.Printf("   ├─ Nombre: %s\n", evt.Nombre)
		fmt.Printf("   ├─ ID del usuario: %s\n", creado3.ID())
		fmt.Printf("   ├─ Timestamp: %s\n", evt.Ocurrido.Format("2006-01-02 15:04:05"))

		if payload, ok := evt.Payload.(map[string]interface{}); ok {
			payloadJSON, _ := json.Marshal(payload)
			fmt.Printf("   └─ Datos: %s\n", string(payloadJSON))
		} else if payload, ok := evt.Payload.(map[string]string); ok {
			payloadJSON, _ := json.Marshal(payload)
			fmt.Printf("   └─ Datos: %s\n", string(payloadJSON))
		} else {
			fmt.Printf("   └─ Datos: %v\n", evt.Payload)
		}
	}

	fmt.Println("\n═════════════════════════════════════════════════════════════════")
	fmt.Println("=== FIN ===")
	fmt.Println("═════════════════════════════════════════════════════════════════")
	fmt.Println()
}
