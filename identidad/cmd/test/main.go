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

	separator := "============================================================"
	fmt.Println("\n" + separator)
	fmt.Println("     RUTINA DE PRUEBA: USUARIO - FLUJO CON EVENTOS TEMPORAL")
	fmt.Println(separator + "\n")

	// ============================================================================
	// FASE 1: CREACIÓN DE USUARIOS SIN MANIPULACIÓN
	// ============================================================================
	fmt.Println("\n=== FASE 1: CREACIÓN DE USUARIOS ===")

	fmt.Println("\n📝 [1] Crear Juan → estado NO_VERIFICADO")
	fmt.Println("   └─ Correo: juan@test.com | Teléfono: 600111111")
	u1, err := usuario.NuevoUsuario("", "juan@test.com", "Juan", "García", "600111111")
	if err != nil {
		log.Fatalf("❌ Error crear usuario: %v", err)
	}
	// CAPTURAR eventos de creación ANTES de persistir
	eventosCreacion1 := u1.PullEventos()

	creado1, err := repo.Crear(ctx, u1)
	if err != nil {
		log.Fatalf("❌ Error persistir: %v", err)
	}
	fmt.Printf("   ✓ Creado: ID=%s, Estado=%s\n", creado1.ID(), creado1.Estado())

	fmt.Println("\n📝 [2] Crear María → estado NO_VERIFICADO")
	fmt.Println("   └─ Correo: maria@test.com | Teléfono: 600222222")
	u2, err := usuario.NuevoUsuario("", "maria@test.com", "María", "López", "600222222")
	if err != nil {
		log.Fatalf("❌ Error crear usuario: %v", err)
	}
	// CAPTURAR eventos de creación ANTES de persistir
	eventosCreacion2 := u2.PullEventos()

	creado2, err := repo.Crear(ctx, u2)
	if err != nil {
		log.Fatalf("❌ Error persistir: %v", err)
	}
	fmt.Printf("   ✓ Creado: ID=%s, Estado=%s\n", creado2.ID(), creado2.Estado())

	fmt.Println("\n📝 [3] Crear Carlos → estado NO_VERIFICADO")
	fmt.Println("   └─ Correo: carlos@test.com | Teléfono: 600333333")
	u3, err := usuario.NuevoUsuario("", "carlos@test.com", "Carlos", "Ruiz", "600333333")
	if err != nil {
		log.Fatalf("❌ Error crear usuario: %v", err)
	}
	// CAPTURAR eventos de creación ANTES de persistir
	eventosCreacion3 := u3.PullEventos()

	creado3, err := repo.Crear(ctx, u3)
	if err != nil {
		log.Fatalf("❌ Error persistir: %v", err)
	}
	fmt.Printf("   ✓ Creado: ID=%s, Estado=%s\n", creado3.ID(), creado3.Estado())

	// ============================================================================
	// FASE 2: MOSTRAR EVENTOS INICIALES DE CREACIÓN
	// ============================================================================
	fmt.Println("\n=== FASE 2: EVENTOS INICIALES DE CREACIÓN ===")

	fmt.Println("\n📊 Juan eventos:")
	fmt.Printf("   Total: %d\n", len(eventosCreacion1))
	for i, evt := range eventosCreacion1 {
		fmt.Printf("   %d. %s\n", i+1, evt.Nombre)
	}

	fmt.Println("\n📊 María eventos:")
	fmt.Printf("   Total: %d\n", len(eventosCreacion2))
	for i, evt := range eventosCreacion2 {
		fmt.Printf("   %d. %s\n", i+1, evt.Nombre)
	}

	fmt.Println("\n📊 Carlos eventos:")
	fmt.Printf("   Total: %d\n", len(eventosCreacion3))
	for i, evt := range eventosCreacion3 {
		fmt.Printf("   %d. %s\n", i+1, evt.Nombre)
	}

	// ============================================================================
	// FASE 3: MANIPULACIÓN DE ESTADOS
	// ============================================================================
	fmt.Println("\n=== FASE 3: MANIPULACIÓN DE ESTADOS ===")

	fmt.Println("\n🔄 [4] Activar Juan → UsuarioActivado")
	err = creado1.Activar()
	if err != nil {
		log.Fatalf("❌ Error activar Juan: %v", err)
	}
	// CAPTURAR eventos de activación ANTES de cualquier otra operación
	eventosActivar1 := creado1.PullEventos()
	fmt.Printf("   ✓ Juan ahora: %s\n", creado1.Estado())

	fmt.Println("\n🔄 [5] Activar María → UsuarioActivado")
	err = creado2.Activar()
	if err != nil {
		log.Fatalf("❌ Error activar María: %v", err)
	}
	// CAPTURAR eventos de activación ANTES de cualquier otra operación
	eventosActivar2 := creado2.PullEventos()
	fmt.Printf("   ✓ María ahora: %s\n", creado2.Estado())

	fmt.Println("\n🔄 [6] Activar Carlos → UsuarioActivado")
	err = creado3.Activar()
	if err != nil {
		log.Fatalf("❌ Error activar Carlos: %v", err)
	}
	// CAPTURAR eventos de activación ANTES de cualquier otra operación
	eventosActivar3 := creado3.PullEventos()
	fmt.Printf("   ✓ Carlos ahora: %s\n", creado3.Estado())

	fmt.Println("\n🔄 [7] Bloquear Carlos → UsuarioBloqueado")
	err = creado3.Bloquear()
	if err != nil {
		log.Fatalf("❌ Error bloquear Carlos: %v", err)
	}
	// CAPTURAR eventos de bloqueo
	eventosBloquear := creado3.PullEventos()
	fmt.Printf("   ✓ Carlos ahora: %s\n", creado3.Estado())

	// ============================================================================
	// FASE 4: EVENTOS DESPUÉS DE CAMBIOS
	// ============================================================================
	fmt.Println("\n=== FASE 4: EVENTOS DESPUÉS DE CAMBIOS ===")

	fmt.Println("\n📊 Juan eventos (post-activación):")
	fmt.Printf("   Total: %d\n", len(eventosActivar1))
	for i, evt := range eventosActivar1 {
		fmt.Printf("   %d. %s\n", i+1, evt.Nombre)
	}

	fmt.Println("\n📊 María eventos (post-activación):")
	fmt.Printf("   Total: %d\n", len(eventosActivar2))
	for i, evt := range eventosActivar2 {
		fmt.Printf("   %d. %s\n", i+1, evt.Nombre)
	}

	fmt.Println("\n📊 Carlos eventos (post-activación + bloqueo):")
	fmt.Printf("   Total de activación: %d\n", len(eventosActivar3))
	for i, evt := range eventosActivar3 {
		fmt.Printf("   %d. %s\n", i+1, evt.Nombre)
	}
	fmt.Printf("   Total de bloqueo: %d\n", len(eventosBloquear))
	for i, evt := range eventosBloquear {
		fmt.Printf("   %d. %s\n", i+1, evt.Nombre)
	}

	// ============================================================================
	// FASE 5: OPERACIONES BD - PERSISTIR Y LISTAR
	// ============================================================================
	fmt.Println("\n=== FASE 5: OPERACIONES BD ===")

	fmt.Println("\n💾 [8] Persistir cambios en BD")
	err = repo.Actualizar(ctx, creado1)
	if err != nil {
		log.Fatalf("❌ Error actualizar Juan: %v", err)
	}
	err = repo.Actualizar(ctx, creado2)
	if err != nil {
		log.Fatalf("❌ Error actualizar María: %v", err)
	}
	err = repo.Actualizar(ctx, creado3)
	if err != nil {
		log.Fatalf("❌ Error actualizar Carlos: %v", err)
	}
	fmt.Println("   ✓ Cambios persistidos")

	fmt.Println("\n📋 [9] Listar todos los usuarios")
	spec := usuario.EspecificacionUsuario{}
	pag := usuario.Paginacion{Pagina: 1, TamanoPagina: 100}
	todos, err := repo.Listar(ctx, spec, pag)
	if err != nil {
		log.Fatalf("❌ Error listar: %v", err)
	}
	fmt.Printf("   Total: %d usuarios\n", len(todos))
	for i, u := range todos {
		fmt.Printf("   %d. %s %s (%s) - Estado: %s\n",
			i+1, u.Nombre(), u.Apellido(), u.Correo(), u.Estado())
	}

	fmt.Println("\n📋 [10] Listar solo ACTIVOS")
	specActivos := usuario.EspecificacionUsuario{
		ListaLiltros: []usuario.CriterioFiltro{
			{Campo: "estado", Operador: "=", Valor: "ACTIVO"},
		},
	}
	activos, err := repo.Listar(ctx, specActivos, pag)
	if err != nil {
		log.Fatalf("❌ Error listar activos: %v", err)
	}
	fmt.Printf("   Activos: %d\n", len(activos))
	for i, u := range activos {
		fmt.Printf("   %d. %s %s\n", i+1, u.Nombre(), u.Apellido())
	}

	fmt.Println("\n🗑️  [11] Eliminar Carlos")
	err = repo.Eliminar(ctx, creado3.ID())
	if err != nil {
		log.Fatalf("❌ Error eliminar: %v", err)
	}
	fmt.Println("   ✓ Carlos eliminado")

	fmt.Println("\n📋 [12] Listar finales (todos)")
	finales, err := repo.Listar(ctx, spec, pag)
	if err != nil {
		log.Fatalf("❌ Error listar: %v", err)
	}
	fmt.Printf("   Total: %d usuarios\n", len(finales))
	for i, u := range finales {
		fmt.Printf("   %d. %s %s - Estado: %s\n",
			i+1, u.Nombre(), u.Apellido(), u.Estado())
	}

	// ============================================================================
	// RESUMEN
	// ============================================================================
	fmt.Println("\n" + separator)
	fmt.Println("     ✓ RUTINA COMPLETADA EXITOSAMENTE")
	fmt.Println(separator + "\n")
	fmt.Println("📊 RESUMEN:")
	fmt.Printf("   • Usuarios creados: 3\n")
	fmt.Printf("   • Usuarios eliminados: 1\n")
	fmt.Printf("   • Usuarios finales: %d\n", len(finales))
	fmt.Printf("   • Eventos de creación: %d (Juan), %d (María), %d (Carlos)\n",
		len(eventosCreacion1), len(eventosCreacion2), len(eventosCreacion3))
	fmt.Printf("   • Eventos de activación: %d (Juan), %d (María), %d (Carlos)\n",
		len(eventosActivar1), len(eventosActivar2), len(eventosActivar3))
	fmt.Printf("   • Eventos de bloqueo: %d (Carlos)\n", len(eventosBloquear))
	fmt.Println()
}
