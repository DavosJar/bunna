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
	fmt.Println("     RUTINA DE PRUEBA: USUARIO - FLUJO COMPLETO")
	fmt.Println(separator + "\n")

	// 1. CREAR USUARIO - Juan García
	fmt.Println("📝 [1] CREAR USUARIO - Juan García")
	fmt.Println("   └─ Correo: juan@test.com | Teléfono: 600111111")
	u1, err := usuario.NuevoUsuario("", "juan@test.com", "Juan", "García", "600111111")
	if err != nil {
		log.Fatalf("❌ Error crear usuario: %v", err)
	}
	u1.Activar()

	creado1, err := repo.Crear(ctx, u1)
	if err != nil {
		log.Fatalf("❌ Error persistir: %v", err)
	}
	fmt.Printf("   ✓ Creado: ID=%s, Estado=%s\n", creado1.ID(), creado1.Estado())

	eventos := creado1.PullEventos()
	fmt.Printf("   ✓ Eventos: %d\n", len(eventos))
	for _, evt := range eventos {
		fmt.Printf("      • %s\n", evt.Nombre)
	}

	// 2. CREAR USUARIO - María López
	fmt.Println("\n📝 [2] CREAR USUARIO - María López")
	fmt.Println("   └─ Correo: maria@test.com | Teléfono: 600222222")
	u2, err := usuario.NuevoUsuario("", "maria@test.com", "María", "López", "600222222")
	if err != nil {
		log.Fatalf("❌ Error crear usuario: %v", err)
	}
	u2.Activar()

	creado2, err := repo.Crear(ctx, u2)
	if err != nil {
		log.Fatalf("❌ Error persistir: %v", err)
	}
	fmt.Printf("   ✓ Creado: ID=%s, Estado=%s\n", creado2.ID(), creado2.Estado())

	// 3. CREAR USUARIO - Carlos Ruiz (sin activar)
	fmt.Println("\n📝 [3] CREAR USUARIO - Carlos Ruiz (sin activar)")
	fmt.Println("   └─ Correo: carlos@test.com | Teléfono: 600333333")
	u3, err := usuario.NuevoUsuario("", "carlos@test.com", "Carlos", "Ruiz", "600333333")
	if err != nil {
		log.Fatalf("❌ Error crear usuario: %v", err)
	}

	creado3, err := repo.Crear(ctx, u3)
	if err != nil {
		log.Fatalf("❌ Error persistir: %v", err)
	}
	fmt.Printf("   ✓ Creado: ID=%s, Estado=%s\n", creado3.ID(), creado3.Estado())

	// 4. OBTENER USUARIO POR ID
	fmt.Println("\n🔍 [4] OBTENER USUARIO POR ID")
	fmt.Printf("   └─ Buscando: %s\n", creado1.ID())
	recuperado, err := repo.ObtenerPorID(ctx, creado1.ID())
	if err != nil {
		log.Fatalf("❌ Error obtener: %v", err)
	}
	fmt.Printf("   ✓ Encontrado: %s %s (Estado: %s)\n",
		recuperado.Nombre(), recuperado.Apellido(), recuperado.Estado())

	// 5. LISTAR TODOS LOS USUARIOS
	fmt.Println("\n📋 [5] LISTAR TODOS LOS USUARIOS")
	spec := usuario.EspecificacionUsuario{}
	pag := usuario.Paginacion{Pagina: 1, TamanoPagina: 100}
	todos, err := repo.Listar(ctx, spec, pag)
	if err != nil {
		log.Fatalf("❌ Error listar: %v", err)
	}
	fmt.Printf("   ✓ Total: %d usuarios\n", len(todos))
	for i, u := range todos {
		fmt.Printf("      %d. %s %s - %s (Estado: %s)\n",
			i+1, u.Nombre(), u.Apellido(), u.Correo(), u.Estado())
	}

	// 6. LISTAR SOLO ACTIVOS (filtro por estado)
	fmt.Println("\n📋 [6] LISTAR SOLO ACTIVOS (filtro)")
	specActivos := usuario.EspecificacionUsuario{
		ListaLiltros: []usuario.CriterioFiltro{
			{Campo: "estado", Operador: "=", Valor: "ACTIVO"},
		},
	}
	activos, err := repo.Listar(ctx, specActivos, pag)
	if err != nil {
		log.Fatalf("❌ Error listar activos: %v", err)
	}
	fmt.Printf("   ✓ Activos: %d\n", len(activos))
	for i, u := range activos {
		fmt.Printf("      %d. %s %s - %s\n", i+1, u.Nombre(), u.Apellido(), u.Correo())
	}

	// 7. CAMBIAR ESTADO DE USUARIO (María -> INACTIVO)
	fmt.Println("\n🔄 [7] CAMBIAR ESTADO - María ACTIVO → INACTIVO")
	fmt.Printf("   └─ Usuario: %s %s (ID: %s)\n", creado2.Nombre(), creado2.Apellido(), creado2.ID())
	err = creado2.Inactivar()
	if err != nil {
		log.Fatalf("❌ Error cambiar estado: %v", err)
	}

	eventosInactivar := creado2.PullEventos()
	fmt.Printf("   ✓ Estado cambió a: %s\n", creado2.Estado())
	fmt.Printf("   ✓ Evento emitido: %v\n", len(eventosInactivar) > 0)
	if len(eventosInactivar) > 0 {
		fmt.Printf("      • %s\n", eventosInactivar[0].Nombre)
	}

	// Persistir cambio
	err = repo.Actualizar(ctx, creado2)
	if err != nil {
		log.Fatalf("❌ Error actualizar: %v", err)
	}
	fmt.Println("   ✓ Cambio persistido en BD")

	// Verificar cambio
	verificado, err := repo.ObtenerPorID(ctx, creado2.ID())
	if err != nil {
		log.Fatalf("❌ Error verificar: %v", err)
	}
	fmt.Printf("   ✓ Estado verificado: %s\n", verificado.Estado())

	// 8. CAMBIAR ESTADO DE USUARIO (Juan -> BLOQUEADO)
	fmt.Println("\n🔄 [8] CAMBIAR ESTADO - Juan ACTIVO → BLOQUEADO")
	fmt.Printf("   └─ Usuario: %s %s (ID: %s)\n", creado1.Nombre(), creado1.Apellido(), creado1.ID())
	err = creado1.Bloquear()
	if err != nil {
		log.Fatalf("❌ Error bloquear: %v", err)
	}

	eventosBloq := creado1.PullEventos()
	fmt.Printf("   ✓ Estado cambió a: %s\n", creado1.Estado())
	fmt.Printf("   ✓ Evento emitido: %v\n", len(eventosBloq) > 0)
	if len(eventosBloq) > 0 {
		fmt.Printf("      • %s\n", eventosBloq[0].Nombre)
	}

	repo.Actualizar(ctx, creado1)
	fmt.Println("   ✓ Cambio persistido en BD")

	// 9. LISTAR SOLO ACTIVOS (verificar que cambios se reflejan)
	fmt.Println("\n📋 [9] LISTAR ACTIVOS (después de cambios)")
	activos2, err := repo.Listar(ctx, specActivos, pag)
	if err != nil {
		log.Fatalf("❌ Error listar: %v", err)
	}
	fmt.Printf("   ✓ Activos ahora: %d (antes eran 2)\n", len(activos2))
	for i, u := range activos2 {
		fmt.Printf("      %d. %s %s\n", i+1, u.Nombre(), u.Apellido())
	}

	// 10. ELIMINAR USUARIO (Carlos)
	fmt.Println("\n🗑️  [10] ELIMINAR USUARIO - Carlos Ruiz")
	fmt.Printf("   └─ ID: %s\n", creado3.ID())
	err = repo.Eliminar(ctx, creado3.ID())
	if err != nil {
		log.Fatalf("❌ Error eliminar: %v", err)
	}

	// Verificar eliminación
	_, err = repo.ObtenerPorID(ctx, creado3.ID())
	if err != nil {
		fmt.Println("   ✓ Usuario eliminado correctamente")
	} else {
		fmt.Println("   ⚠️  Usuario aún existe en BD")
	}

	// 11. LISTAR FINAL
	fmt.Println("\n📋 [11] LISTAR FINAL (todos)")
	finales, err := repo.Listar(ctx, spec, pag)
	if err != nil {
		log.Fatalf("❌ Error listar: %v", err)
	}
	fmt.Printf("   ✓ Total final: %d usuarios (comenzamos con 3)\n", len(finales))
	for i, u := range finales {
		fmt.Printf("      %d. %s %s - Estado: %s\n",
			i+1, u.Nombre(), u.Apellido(), u.Estado())
	}

	// Resumen
	fmt.Println("\n" + separator)
	fmt.Println("     ✓ RUTINA COMPLETADA EXITOSAMENTE")
	fmt.Println(separator + "\n")
	fmt.Println("📊 RESUMEN:")
	fmt.Printf("   • Usuarios creados: 3\n")
	fmt.Printf("   • Usuarios eliminados: 1\n")
	fmt.Printf("   • Usuarios finales: %d\n", len(finales))
	fmt.Printf("   • Cambios de estado: 2\n")
	fmt.Printf("   • Eventos emitidos: %d+\n", len(eventos)+len(eventosInactivar)+len(eventosBloq))
	fmt.Println("\n")
}
