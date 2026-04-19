package main

import (
	"context"
	"fmt"
	"log"

	"github.com/davosjar/bunna/services/identidad/internal/config"
	"github.com/davosjar/bunna/services/identidad/internal/domain/usuario"
	"github.com/davosjar/bunna/services/identidad/internal/registry"
	"github.com/google/uuid"
)

// Estructura para datos de usuario
type datosUsuario struct {
	nombre   string
	apellido string
	correo   string
	telefono string
	estado   usuario.EstadoUsuario
}

func main() {
	// ============================================================================
	// === SECCIÓN 1: INICIALIZACIÓN ===
	// ============================================================================
	fmt.Println("\n═════════════════════════════════════════════════════════════════")
	fmt.Println("=== INICIALIZACIÓN ===")
	fmt.Println("═════════════════════════════════════════════════════════════════")

	// Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Config: %v", err)
	}
	fmt.Println("✓ Configuración cargada")

	// Inicializar base de datos
	db, err := config.InitDB(cfg.GetDSN())
	if err != nil {
		log.Fatalf("❌ DB: %v", err)
	}
	fmt.Println("✓ Base de datos conectada")

	// Limpiar y migrar
	db.Migrator().DropTable(&usuario.Usuario{})
	if err := config.RunMigrations(db); err != nil {
		log.Fatalf("❌ Migrations: %v", err)
	}
	fmt.Println("✓ Tabla usuarios limpiada y migrada")

	// Crear registry e inyectar repositorio
	reg := registry.NewRegistry(db)
	repo := reg.UsuarioRepository()
	ctx := context.Background()

	// ============================================================================
	// === SECCIÓN 2: CREAR 30 USUARIOS ===
	// ============================================================================
	fmt.Println("\n═════════════════════════════════════════════════════════════════")
	fmt.Println("=== CREACIÓN DE 30 USUARIOS ===")
	fmt.Println("═════════════════════════════════════════════════════════════════")

	// Datos de los 30 usuarios (ACTIVO, NO_VERIFICADO, BLOQUEADO, INACTIVO)
	usuariosData := []datosUsuario{
		// ACTIVO (10)
		{nombre: "Juan", apellido: "García", correo: "juan@test.com", telefono: "6001111111", estado: usuario.ACTIVO},
		{nombre: "María", apellido: "López", correo: "maria@test.com", telefono: "6002222222", estado: usuario.ACTIVO},
		{nombre: "Carlos", apellido: "Ruiz", correo: "carlos@test.com", telefono: "6003333333", estado: usuario.ACTIVO},
		{nombre: "Ana", apellido: "Martínez", correo: "ana@example.com", telefono: "6004444444", estado: usuario.ACTIVO},
		{nombre: "Pedro", apellido: "Sánchez", correo: "pedro@demo.com", telefono: "6005555555", estado: usuario.ACTIVO},
		{nombre: "Isabel", apellido: "Fernández", correo: "isabel@test.com", telefono: "6006666666", estado: usuario.ACTIVO},
		{nombre: "Diego", apellido: "Morales", correo: "diego@example.com", telefono: "6007777777", estado: usuario.ACTIVO},
		{nombre: "Rosa", apellido: "García", correo: "rosa@demo.com", telefono: "6008888888", estado: usuario.ACTIVO},
		{nombre: "Luis", apellido: "Rodríguez", correo: "luis@test.com", telefono: "6009999999", estado: usuario.ACTIVO},
		{nombre: "Miguel", apellido: "Díaz", correo: "miguel@example.com", telefono: "6010101010", estado: usuario.ACTIVO},
		// NO_VERIFICADO (10)
		{nombre: "Elena", apellido: "Gutiérrez", correo: "elena@test.com", telefono: "6011111111", estado: usuario.NO_VERIFICADO},
		{nombre: "Francisco", apellido: "López", correo: "francisco@example.com", telefono: "6012121212", estado: usuario.NO_VERIFICADO},
		{nombre: "Carmen", apellido: "García", correo: "carmen@demo.com", telefono: "6013131313", estado: usuario.NO_VERIFICADO},
		{nombre: "Manuel", apellido: "Sánchez", correo: "manuel@test.com", telefono: "6014141414", estado: usuario.NO_VERIFICADO},
		{nombre: "Javier", apellido: "Pérez", correo: "javier@example.com", telefono: "6015151515", estado: usuario.NO_VERIFICADO},
		{nombre: "Cristina", apellido: "Flores", correo: "cristina@demo.com", telefono: "6016161616", estado: usuario.NO_VERIFICADO},
		{nombre: "Antonio", apellido: "Ramírez", correo: "antonio@test.com", telefono: "6017171717", estado: usuario.NO_VERIFICADO},
		{nombre: "Margarita", apellido: "Vega", correo: "margarita@example.com", telefono: "6018181818", estado: usuario.NO_VERIFICADO},
		{nombre: "Rafael", apellido: "Ortiz", correo: "rafael@demo.com", telefono: "6019191919", estado: usuario.NO_VERIFICADO},
		{nombre: "Silvia", apellido: "Ramos", correo: "silvia@test.com", telefono: "6020202020", estado: usuario.NO_VERIFICADO},
		// BLOQUEADO (5)
		{nombre: "Ángel", apellido: "Torres", correo: "angel@example.com", telefono: "6021212121", estado: usuario.BLOQUEADO},
		{nombre: "Beatriz", apellido: "Iglesias", correo: "beatriz@demo.com", telefono: "6022222222", estado: usuario.BLOQUEADO},
		{nombre: "Enrique", apellido: "Núñez", correo: "enrique@test.com", telefono: "6023232323", estado: usuario.BLOQUEADO},
		{nombre: "Dolores", apellido: "Campos", correo: "dolores@example.com", telefono: "6024242424", estado: usuario.BLOQUEADO},
		{nombre: "Gerardo", apellido: "Herrera", correo: "gerardo@demo.com", telefono: "6025252525", estado: usuario.BLOQUEADO},
		// INACTIVO (5)
		{nombre: "Herminia", apellido: "Romero", correo: "herminia@test.com", telefono: "6026262626", estado: usuario.INACTIVO},
		{nombre: "Ignacio", apellido: "Vargas", correo: "ignacio@example.com", telefono: "6027272727", estado: usuario.INACTIVO},
		{nombre: "Julia", apellido: "Salazar", correo: "julia@demo.com", telefono: "6028282828", estado: usuario.INACTIVO},
		{nombre: "Karim", apellido: "Malik", correo: "karim@test.com", telefono: "6029292929", estado: usuario.INACTIVO},
		{nombre: "Lucía", apellido: "Estrada", correo: "lucia@example.com", telefono: "6030303030", estado: usuario.INACTIVO},
	}

	var usuariosCreados []*usuario.Usuario
	for i, datos := range usuariosData {
		// Crear usuario (se crea siempre en NO_VERIFICADO según NuevoUsuario)
		// Generar UUID para el usuario
		usuarioID := uuid.New().String()
		u, err := usuario.NuevoUsuario(usuarioID, datos.correo, datos.nombre, datos.apellido, datos.telefono)
		if err != nil {
			log.Fatalf("❌ Error crear usuario %d: %v", i+1, err)
		}

		// Cambiar estado si es necesario (siguiendo las transiciones permitidas)
		if datos.estado != usuario.NO_VERIFICADO {
			// Primero pasar a ACTIVO si es necesario
			if err := u.Activar(); err != nil {
				log.Fatalf("❌ Error activar usuario %d: %v", i+1, err)
			}
			// Si el estado final es diferente, hacer la transición
			if datos.estado != usuario.ACTIVO {
				if err := u.CambiarEstado(datos.estado); err != nil {
					log.Fatalf("❌ Error cambiar estado usuario %d: %v", i+1, err)
				}
			}
		}

		// Persistir
		creado, err := repo.Crear(ctx, u)
		if err != nil {
			log.Fatalf("❌ Error persistir usuario %d: %v", i+1, err)
		}

		usuariosCreados = append(usuariosCreados, creado)

		// Mostrar resumen
		fmt.Printf("[%2d] %-15s %-20s → %-15s → %s\n",
			i+1,
			datos.nombre+" "+datos.apellido,
			datos.correo,
			datos.estado,
			creado.ID())
	}

	fmt.Printf("\n✓ 30 usuarios creados exitosamente\n")

	// ============================================================================
	// === SECCIÓN 3: PRUEBAS DE FILTROS ===
	// ============================================================================
	fmt.Println("\n═════════════════════════════════════════════════════════════════")
	fmt.Println("=== PRUEBAS DE FILTROS ===")
	fmt.Println("═════════════════════════════════════════════════════════════════")

	// FILTRO 1: Estado = ACTIVO
	fmt.Println("📌 FILTRO 1: Estado = ACTIVO")
	spec1 := usuario.EspecificacionUsuario{
		ListaLiltros: []usuario.CriterioFiltro{
			{Campo: "estado", Operador: "=", Valor: string(usuario.ACTIVO)},
		},
	}
	pag1 := usuario.Paginacion{Pagina: 1, TamanoPagina: 100}
	activos, err := repo.Listar(ctx, spec1, pag1)
	if err != nil {
		log.Fatalf("❌ Error filtro ACTIVO: %v", err)
	}
	fmt.Printf("   Resultados: %d usuarios\n", len(activos))
	for _, u := range activos {
		fmt.Printf("   ├─ %-20s (%s)\n", u.Nombre()+" "+u.Apellido(), u.Correo())
	}

	// FILTRO 2: Estado = NO_VERIFICADO
	fmt.Println("\n📌 FILTRO 2: Estado = NO_VERIFICADO")
	spec2 := usuario.EspecificacionUsuario{
		ListaLiltros: []usuario.CriterioFiltro{
			{Campo: "estado", Operador: "=", Valor: string(usuario.NO_VERIFICADO)},
		},
	}
	pag2 := usuario.Paginacion{Pagina: 1, TamanoPagina: 100}
	noVerificados, err := repo.Listar(ctx, spec2, pag2)
	if err != nil {
		log.Fatalf("❌ Error filtro NO_VERIFICADO: %v", err)
	}
	fmt.Printf("   Resultados: %d usuarios\n", len(noVerificados))
	for _, u := range noVerificados {
		fmt.Printf("   ├─ %-20s (%s)\n", u.Nombre()+" "+u.Apellido(), u.Correo())
	}

	// FILTRO 3: Correo LIKE "%@example.com"
	fmt.Println("\n📌 FILTRO 3: Correo LIKE \"%@example.com\"")
	spec3 := usuario.EspecificacionUsuario{
		ListaLiltros: []usuario.CriterioFiltro{
			{Campo: "correo", Operador: "LIKE", Valor: "%@example.com"},
		},
	}
	pag3 := usuario.Paginacion{Pagina: 1, TamanoPagina: 100}
	example, err := repo.Listar(ctx, spec3, pag3)
	if err != nil {
		log.Fatalf("❌ Error filtro example.com: %v", err)
	}
	fmt.Printf("   Resultados: %d usuarios\n", len(example))
	for _, u := range example {
		fmt.Printf("   ├─ %-20s (%s)\n", u.Nombre()+" "+u.Apellido(), u.Correo())
	}

	// FILTRO 4: Correo LIKE "%test%"
	fmt.Printf("📌 FILTRO 4: Correo LIKE \"%%test%%\"\n")
	spec4 := usuario.EspecificacionUsuario{
		ListaLiltros: []usuario.CriterioFiltro{
			{Campo: "correo", Operador: "LIKE", Valor: "%test%"},
		},
	}
	pag4 := usuario.Paginacion{Pagina: 1, TamanoPagina: 100}
	testDom, err := repo.Listar(ctx, spec4, pag4)
	if err != nil {
		log.Fatalf("❌ Error filtro test domain: %v", err)
	}
	fmt.Printf("   Resultados: %d usuarios\n", len(testDom))
	for _, u := range testDom {
		fmt.Printf("   ├─ %-20s (%s)\n", u.Nombre()+" "+u.Apellido(), u.Correo())
	}

	// FILTRO 5: Estado = NO_VERIFICADO + Correo LIKE "%test%"
	fmt.Printf("📌 FILTRO 5: Estado = NO_VERIFICADO + Correo LIKE \"%%test%%\"\n")
	spec5 := usuario.EspecificacionUsuario{
		ListaLiltros: []usuario.CriterioFiltro{
			{Campo: "estado", Operador: "=", Valor: string(usuario.NO_VERIFICADO)},
			{Campo: "correo", Operador: "LIKE", Valor: "%test%"},
		},
	}
	pag5 := usuario.Paginacion{Pagina: 1, TamanoPagina: 100}
	combinado, err := repo.Listar(ctx, spec5, pag5)
	if err != nil {
		log.Fatalf("❌ Error filtro combinado: %v", err)
	}
	fmt.Printf("   Resultados: %d usuarios\n", len(combinado))
	for _, u := range combinado {
		fmt.Printf("   ├─ %-20s (%s)\n", u.Nombre()+" "+u.Apellido(), u.Correo())
	}

	// FILTRO 6: Estado = BLOQUEADO
	fmt.Println("\n📌 FILTRO 6: Estado = BLOQUEADO")
	spec6 := usuario.EspecificacionUsuario{
		ListaLiltros: []usuario.CriterioFiltro{
			{Campo: "estado", Operador: "=", Valor: string(usuario.BLOQUEADO)},
		},
	}
	pag6 := usuario.Paginacion{Pagina: 1, TamanoPagina: 100}
	bloqueados, err := repo.Listar(ctx, spec6, pag6)
	if err != nil {
		log.Fatalf("❌ Error filtro BLOQUEADO: %v", err)
	}
	fmt.Printf("   Resultados: %d usuarios\n", len(bloqueados))
	for _, u := range bloqueados {
		fmt.Printf("   ├─ %-20s (%s)\n", u.Nombre()+" "+u.Apellido(), u.Correo())
	}

	// FILTRO 7: Nombre LIKE "%García%"
	fmt.Printf("📌 FILTRO 7: Nombre LIKE \"%%García%%\"\n")
	spec7 := usuario.EspecificacionUsuario{
		ListaLiltros: []usuario.CriterioFiltro{
			{Campo: "nombre", Operador: "LIKE", Valor: "%García%"},
		},
	}
	pag7 := usuario.Paginacion{Pagina: 1, TamanoPagina: 100}
	garcia, err := repo.Listar(ctx, spec7, pag7)
	if err != nil {
		log.Fatalf("❌ Error filtro García: %v", err)
	}
	fmt.Printf("   Resultados: %d usuarios\n", len(garcia))
	for _, u := range garcia {
		fmt.Printf("   ├─ %-20s (%s)\n", u.Nombre()+" "+u.Apellido(), u.Correo())
	}

	// ============================================================================
	// === SECCIÓN 4: PRUEBAS DE PAGINACIÓN ===
	// ============================================================================
	fmt.Println("\n═════════════════════════════════════════════════════════════════")
	fmt.Println("=== PRUEBAS DE PAGINACIÓN ===")
	fmt.Println("═════════════════════════════════════════════════════════════════")

	// PÁGINA 1 (10 usuarios, tamaño 10)
	fmt.Println("📄 PÁGINA 1 (tamaño 10)")
	specTodos := usuario.EspecificacionUsuario{ListaLiltros: []usuario.CriterioFiltro{}}
	pagPag1 := usuario.Paginacion{
		Pagina:       1,
		TamanoPagina: 10,
		Ordenaciones: []usuario.Ordenacion{
			{Campo: "nombre", Tipo: usuario.ASC},
		},
	}
	pag1Usuarios, err := repo.Listar(ctx, specTodos, pagPag1)
	if err != nil {
		log.Fatalf("❌ Error página 1: %v", err)
	}
	fmt.Printf("   Mostrando 1-10 de 30\n")
	for i, u := range pag1Usuarios {
		fmt.Printf("   ├─ [%2d] %-20s (%s) - %s\n", i+1, u.Nombre()+" "+u.Apellido(), u.Correo(), u.Estado())
	}

	// PÁGINA 2 (10 usuarios, tamaño 10)
	fmt.Println("\n📄 PÁGINA 2 (tamaño 10)")
	pagPag2 := usuario.Paginacion{
		Pagina:       2,
		TamanoPagina: 10,
		Ordenaciones: []usuario.Ordenacion{
			{Campo: "nombre", Tipo: usuario.ASC},
		},
	}
	pag2Usuarios, err := repo.Listar(ctx, specTodos, pagPag2)
	if err != nil {
		log.Fatalf("❌ Error página 2: %v", err)
	}
	fmt.Printf("   Mostrando 11-20 de 30\n")
	for i, u := range pag2Usuarios {
		fmt.Printf("   ├─ [%2d] %-20s (%s) - %s\n", 11+i, u.Nombre()+" "+u.Apellido(), u.Correo(), u.Estado())
	}

	// PÁGINA 3 (10 usuarios, tamaño 10)
	fmt.Println("\n📄 PÁGINA 3 (tamaño 10)")
	pagPag3 := usuario.Paginacion{
		Pagina:       3,
		TamanoPagina: 10,
		Ordenaciones: []usuario.Ordenacion{
			{Campo: "nombre", Tipo: usuario.ASC},
		},
	}
	pag3Usuarios, err := repo.Listar(ctx, specTodos, pagPag3)
	if err != nil {
		log.Fatalf("❌ Error página 3: %v", err)
	}
	fmt.Printf("   Mostrando 21-30 de 30\n")
	for i, u := range pag3Usuarios {
		fmt.Printf("   ├─ [%2d] %-20s (%s) - %s\n", 21+i, u.Nombre()+" "+u.Apellido(), u.Correo(), u.Estado())
	}

	// PAGINACIÓN CON ORDENACIÓN ASC
	fmt.Println("\n📊 PAGINACIÓN CON ORDENACIÓN: Nombre ASC (Página 1, tamaño 5)")
	pagAscNombre := usuario.Paginacion{
		Pagina:       1,
		TamanoPagina: 5,
		Ordenaciones: []usuario.Ordenacion{
			{Campo: "nombre", Tipo: usuario.ASC},
		},
	}
	ascNombre, err := repo.Listar(ctx, specTodos, pagAscNombre)
	if err != nil {
		log.Fatalf("❌ Error ordenación ASC nombre: %v", err)
	}
	for i, u := range ascNombre {
		fmt.Printf("   ├─ [%d] %-20s (%s)\n", i+1, u.Nombre()+" "+u.Apellido(), u.Correo())
	}

	// PAGINACIÓN CON ORDENACIÓN DESC
	fmt.Println("\n📊 PAGINACIÓN CON ORDENACIÓN: Correo DESC (Página 1, tamaño 5)")
	pagDescCorreo := usuario.Paginacion{
		Pagina:       1,
		TamanoPagina: 5,
		Ordenaciones: []usuario.Ordenacion{
			{Campo: "correo", Tipo: usuario.DESC},
		},
	}
	descCorreo, err := repo.Listar(ctx, specTodos, pagDescCorreo)
	if err != nil {
		log.Fatalf("❌ Error ordenación DESC correo: %v", err)
	}
	for i, u := range descCorreo {
		fmt.Printf("   ├─ [%d] %-20s (%s)\n", i+1, u.Nombre()+" "+u.Apellido(), u.Correo())
	}

	// Paginación con filtro combinado
	fmt.Println("\n📊 PAGINACIÓN CON FILTRO Y ORDENACIÓN: Estado=ACTIVO, Correo DESC (Página 1, tamaño 5)")
	specActivoDesc := usuario.EspecificacionUsuario{
		ListaLiltros: []usuario.CriterioFiltro{
			{Campo: "estado", Operador: "=", Valor: string(usuario.ACTIVO)},
		},
	}
	pagActivoDesc := usuario.Paginacion{
		Pagina:       1,
		TamanoPagina: 5,
		Ordenaciones: []usuario.Ordenacion{
			{Campo: "correo", Tipo: usuario.DESC},
		},
	}
	activoDesc, err := repo.Listar(ctx, specActivoDesc, pagActivoDesc)
	if err != nil {
		log.Fatalf("❌ Error filtro ACTIVO + ordenación DESC: %v", err)
	}
	for i, u := range activoDesc {
		fmt.Printf("   ├─ [%d] %-20s (%s)\n", i+1, u.Nombre()+" "+u.Apellido(), u.Correo())
	}

	// ============================================================================
	// === RESUMEN FINAL ===
	// ============================================================================
	fmt.Println("\n═════════════════════════════════════════════════════════════════")
	fmt.Println("=== RESUMEN FINAL ===")
	fmt.Println("═════════════════════════════════════════════════════════════════")
	fmt.Printf("✅ Total de usuarios creados: %d\n", len(usuariosCreados))
	fmt.Printf("✅ ACTIVOS: %d\n", len(activos))
	fmt.Printf("✅ NO_VERIFICADOS: %d\n", len(noVerificados))
	fmt.Printf("✅ BLOQUEADOS: %d\n", len(bloqueados))
	fmt.Printf("✅ INACTIVOS: %d (no se filtraron explícitamente)\n", 5)
	fmt.Printf("✅ Correos en @example.com: %d\n", len(example))
	fmt.Printf("✅ Correos con 'test': %d\n", len(testDom))
	fmt.Printf("✅ García (nombre): %d\n", len(garcia))
	fmt.Println("\n✓ Pruebas completadas exitosamente")
}
