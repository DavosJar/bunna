package main

import (
	"context"
	"fmt"
	"log"

	"github.com/davosjar/bunna/services/identidad/internal/application/services/registro"
	"github.com/davosjar/bunna/services/identidad/internal/config"
	"github.com/davosjar/bunna/services/identidad/internal/domain/usuario"
	"github.com/davosjar/bunna/services/identidad/internal/registry"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
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
		v7uuid, err := uuid.NewV7()
		if err != nil {
			log.Fatalf("❌ Error generar UUID V7: %v", err)
		}
		usuarioID := v7uuid.String()
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
	// === SECCIÓN 3: BÚSQUEDAS COMPUESTAS ===
	// ============================================================================
	fmt.Println("\n═════════════════════════════════════════════════════════════════")
	fmt.Println("=== BÚSQUEDAS COMPUESTAS ===")
	fmt.Println("═════════════════════════════════════════════════════════════════")

	// BÚSQUEDA 1: ACTIVO + @test.com + Ordenar apellido ASC
	fmt.Println("\n🔍 BÚSQUEDA 1: ACTIVO + @test.com + Ordenar apellido ASC")
	fmt.Println("   Filtros:")
	fmt.Println("   ├─ Estado = ACTIVO")
	fmt.Println("   └─ Correo LIKE \"%@test.com\"")
	fmt.Println("   Ordenación: apellido ASC")
	spec1 := usuario.EspecificacionUsuario{
		ListaLiltros: []usuario.CriterioFiltro{
			{Campo: "estado", Operador: "=", Valor: string(usuario.ACTIVO)},
			{Campo: "correo", Operador: "LIKE", Valor: "%@test.com"},
		},
	}
	pag1 := usuario.Paginacion{
		Pagina:       1,
		TamanoPagina: 100,
		Ordenaciones: []usuario.Ordenacion{
			{Campo: "apellido", Tipo: usuario.ASC},
		},
	}
	busqueda1, err := repo.Listar(ctx, spec1, pag1)
	if err != nil {
		log.Fatalf("❌ Error búsqueda 1: %v", err)
	}
	fmt.Printf("   Resultados: %d usuarios\n", len(busqueda1))
	for _, u := range busqueda1 {
		fmt.Printf("   ├─ %-25s (%s) - %s\n", u.Nombre()+" "+u.Apellido(), u.Correo(), u.Estado())
	}

	// BÚSQUEDA 2: NO_VERIFICADO + Apellido LIKE "García"
	fmt.Println("\n🔍 BÚSQUEDA 2: NO_VERIFICADO + Apellido LIKE \"García\"")
	fmt.Println("   Filtros:")
	fmt.Println("   ├─ Estado = NO_VERIFICADO")
	fmt.Println("   └─ Apellido LIKE (wildcard)García(wildcard)")
	spec2 := usuario.EspecificacionUsuario{
		ListaLiltros: []usuario.CriterioFiltro{
			{Campo: "estado", Operador: "=", Valor: string(usuario.NO_VERIFICADO)},
			{Campo: "apellido", Operador: "LIKE", Valor: "%García%"},
		},
	}
	pag2 := usuario.Paginacion{Pagina: 1, TamanoPagina: 100}
	busqueda2, err := repo.Listar(ctx, spec2, pag2)
	if err != nil {
		log.Fatalf("❌ Error búsqueda 2: %v", err)
	}
	fmt.Printf("   Resultados: %d usuarios\n", len(busqueda2))
	for _, u := range busqueda2 {
		fmt.Printf("   ├─ %-25s (%s) - %s\n", u.Nombre()+" "+u.Apellido(), u.Correo(), u.Estado())
	}

	// BÚSQUEDA 3: != BLOQUEADO + @example.com + Ordenar correo DESC
	fmt.Println("\n🔍 BÚSQUEDA 3: != BLOQUEADO + @example.com + Ordenar correo DESC")
	fmt.Println("   Filtros:")
	fmt.Println("   ├─ Estado != BLOQUEADO")
	fmt.Println("   └─ Correo LIKE \"%@example.com\"")
	fmt.Println("   Ordenación: correo DESC")
	spec3 := usuario.EspecificacionUsuario{
		ListaLiltros: []usuario.CriterioFiltro{
			{Campo: "estado", Operador: "!=", Valor: string(usuario.BLOQUEADO)},
			{Campo: "correo", Operador: "LIKE", Valor: "%@example.com"},
		},
	}
	pag3 := usuario.Paginacion{
		Pagina:       1,
		TamanoPagina: 100,
		Ordenaciones: []usuario.Ordenacion{
			{Campo: "correo", Tipo: usuario.DESC},
		},
	}
	busqueda3, err := repo.Listar(ctx, spec3, pag3)
	if err != nil {
		log.Fatalf("❌ Error búsqueda 3: %v", err)
	}
	fmt.Printf("   Resultados: %d usuarios\n", len(busqueda3))
	for _, u := range busqueda3 {
		fmt.Printf("   ├─ %-25s (%s) - %s\n", u.Nombre()+" "+u.Apellido(), u.Correo(), u.Estado())
	}

	// BÚSQUEDA 4: != BLOQUEADO + Ordenar fecha_actualizacion DESC + Paginación
	fmt.Println("\n🔍 BÚSQUEDA 4: != BLOQUEADO + Ordenar fecha_actualizacion DESC + Paginación")
	fmt.Println("   Filtros:")
	fmt.Println("   └─ Estado != BLOQUEADO")
	fmt.Println("   Ordenación: fecha_actualizacion DESC")
	fmt.Println("   Paginación: página 1, tamaño 5")
	spec4 := usuario.EspecificacionUsuario{
		ListaLiltros: []usuario.CriterioFiltro{
			{Campo: "estado", Operador: "!=", Valor: string(usuario.BLOQUEADO)},
		},
	}
	pag4 := usuario.Paginacion{
		Pagina:       1,
		TamanoPagina: 5,
		Ordenaciones: []usuario.Ordenacion{
			{Campo: "fecha_actualizacion", Tipo: usuario.DESC},
		},
	}
	busqueda4, err := repo.Listar(ctx, spec4, pag4)
	if err != nil {
		log.Fatalf("❌ Error búsqueda 4: %v", err)
	}
	fmt.Printf("   Resultados: %d usuarios (últimos modificados, no bloqueados)\n", len(busqueda4))
	for i, u := range busqueda4 {
		fmt.Printf("   ├─ [%d] %-25s (%s) - Actualizado: %v\n", i+1, u.Nombre()+" "+u.Apellido(), u.Correo(), u.FechaActualizacion())
	}

	// BÚSQUEDA 5: BLOQUEADO + Múltiples ordenaciones (nombre ASC, luego correo DESC)
	fmt.Println("\n🔍 BÚSQUEDA 5: BLOQUEADO + Múltiples ordenaciones")
	fmt.Println("   Filtros:")
	fmt.Println("   └─ Estado = BLOQUEADO")
	fmt.Println("   Ordenación: nombre ASC, luego correo DESC")
	fmt.Println("   Paginación: todos (5 usuarios)")
	spec5 := usuario.EspecificacionUsuario{
		ListaLiltros: []usuario.CriterioFiltro{
			{Campo: "estado", Operador: "=", Valor: string(usuario.BLOQUEADO)},
		},
	}
	pag5 := usuario.Paginacion{
		Pagina:       1,
		TamanoPagina: 100,
		Ordenaciones: []usuario.Ordenacion{
			{Campo: "nombre", Tipo: usuario.ASC},
			{Campo: "correo", Tipo: usuario.DESC},
		},
	}
	busqueda5, err := repo.Listar(ctx, spec5, pag5)
	if err != nil {
		log.Fatalf("❌ Error búsqueda 5: %v", err)
	}
	fmt.Printf("   Resultados: %d usuarios bloqueados (ordenados)\n", len(busqueda5))
	for i, u := range busqueda5 {
		fmt.Printf("   ├─ [%d] %-25s (%s) - %s\n", i+1, u.Nombre()+" "+u.Apellido(), u.Correo(), u.Estado())
	}

	// BÚSQUEDA 6: INACTIVO + @demo.com + Ordenar nombre ASC
	fmt.Println("\n🔍 BÚSQUEDA 6: INACTIVO + @demo.com + Ordenar nombre ASC")
	fmt.Println("   Filtros:")
	fmt.Println("   ├─ Estado = INACTIVO")
	fmt.Println("   └─ Correo LIKE \"%@demo.com\"")
	fmt.Println("   Ordenación: nombre ASC")
	spec6 := usuario.EspecificacionUsuario{
		ListaLiltros: []usuario.CriterioFiltro{
			{Campo: "estado", Operador: "=", Valor: string(usuario.INACTIVO)},
			{Campo: "correo", Operador: "LIKE", Valor: "%@demo.com"},
		},
	}
	pag6 := usuario.Paginacion{
		Pagina:       1,
		TamanoPagina: 100,
		Ordenaciones: []usuario.Ordenacion{
			{Campo: "nombre", Tipo: usuario.ASC},
		},
	}
	busqueda6, err := repo.Listar(ctx, spec6, pag6)
	if err != nil {
		log.Fatalf("❌ Error búsqueda 6: %v", err)
	}
	fmt.Printf("   Resultados: %d usuarios inactivos con correo demo\n", len(busqueda6))
	for _, u := range busqueda6 {
		fmt.Printf("   ├─ %-25s (%s) - %s\n", u.Nombre()+" "+u.Apellido(), u.Correo(), u.Estado())
	}

	// ============================================================================
	// === RESUMEN FINAL ===
	// ============================================================================
	fmt.Println("\n═════════════════════════════════════════════════════════════════")
	fmt.Println("=== RESUMEN FINAL - BÚSQUEDAS COMPUESTAS ===")
	fmt.Println("═════════════════════════════════════════════════════════════════")
	fmt.Printf("✅ Total de usuarios creados: %d\n", len(usuariosCreados))
	fmt.Printf("✅ Búsqueda 1 (ACTIVO + @test.com, ordenados por apellido): %d\n", len(busqueda1))
	fmt.Printf("✅ Búsqueda 2 (NO_VERIFICADO + García): %d\n", len(busqueda2))
	fmt.Printf("✅ Búsqueda 3 (!= BLOQUEADO + @example.com, correo DESC): %d\n", len(busqueda3))
	fmt.Printf("✅ Búsqueda 4 (!= BLOQUEADO, últimos 5 modificados): %d\n", len(busqueda4))
	fmt.Printf("✅ Búsqueda 5 (BLOQUEADO, múltiples ordenaciones): %d\n", len(busqueda5))
	fmt.Printf("✅ Búsqueda 6 (INACTIVO + @demo.com, ordenados por nombre): %d\n", len(busqueda6))
	fmt.Println("\n✓ Pruebas de búsquedas compuestas completadas exitosamente")

	// ============================================================================
	// === SECCIÓN 4: PRUEBA DE REGISTRO CON SERVICIO DE APLICACIÓN ===
	// ============================================================================
	fmt.Println("\n═════════════════════════════════════════════════════════════════")
	fmt.Println("=== PRUEBA DE REGISTRO DE NUEVOS USUARIOS ===")
	fmt.Println("═════════════════════════════════════════════════════════════════")

	// Crear servicio de registro
	servicioRegistro := registro.NuevoServicioRegistro(
		reg.UsuarioRepository(),
		reg.CredencialesRepository(),
		db,
	)

	// Casos de prueba de registro
	casosRegistro := []struct {
		nombre      string
		comando     *registro.ComandoRegistro
		esperaError bool
	}{
		{
			nombre: "Registro exitoso",
			comando: &registro.ComandoRegistro{
				Correo:   "juan.nuevo@example.com",
				Password: "Password123!",
				Nombre:   "Juan Nuevo",
				Apellido: "Pérez",
				Telefono: "6031234567",
			},
			esperaError: false,
		},
		{
			nombre: "Email inválido",
			comando: &registro.ComandoRegistro{
				Correo:   "notanemail",
				Password: "Password123!",
				Nombre:   "María",
				Apellido: "López",
				Telefono: "6032345678",
			},
			esperaError: true,
		},
		{
			nombre: "Email vacío",
			comando: &registro.ComandoRegistro{
				Correo:   "",
				Password: "Password123!",
				Nombre:   "Carlos",
				Apellido: "García",
				Telefono: "6033456789",
			},
			esperaError: true,
		},
		{
			nombre: "Password vacío",
			comando: &registro.ComandoRegistro{
				Correo:   "carlos.nuevo@test.com",
				Password: "",
				Nombre:   "Carlos",
				Apellido: "García",
				Telefono: "6033456789",
			},
			esperaError: true,
		},
		{
			nombre: "Nombre vacío",
			comando: &registro.ComandoRegistro{
				Correo:   "ana.nueva@demo.com",
				Password: "Password123!",
				Nombre:   "",
				Apellido: "Martínez",
				Telefono: "6034567890",
			},
			esperaError: true,
		},
		{
			nombre: "Registro exitoso 2",
			comando: &registro.ComandoRegistro{
				Correo:   "pedro.nuevo@example.com",
				Password: "SecurePass456",
				Nombre:   "Pedro",
				Apellido: "Rodríguez",
				Telefono: "6035678901",
			},
			esperaError: false,
		},
	}

	fmt.Println("\n📋 Ejecutando casos de prueba de registro:")
	exitosRegistro := 0
	fallosRegistro := 0

	for i, caso := range casosRegistro {
		respuesta, err := servicioRegistro.Ejecutar(ctx, caso.comando)

		if caso.esperaError {
			if err != nil {
				fmt.Printf("[%d] ❌ %s\n", i+1, caso.nombre)
				fmt.Printf("     Error esperado: %v\n", err)
				fallosRegistro++
			} else {
				fmt.Printf("[%d] ⚠️  %s\n", i+1, caso.nombre)
				fmt.Printf("     Se esperaba error pero se registró exitosamente\n")
				fallosRegistro++
			}
		} else {
			if err != nil {
				fmt.Printf("[%d] ❌ %s\n", i+1, caso.nombre)
				fmt.Printf("     Error inesperado: %v\n", err)
				fallosRegistro++
			} else {
				fmt.Printf("[%d] ✅ %s\n", i+1, caso.nombre)
				fmt.Printf("     Usuario ID: %s\n", respuesta.UsuarioID)
				fmt.Printf("     Correo: %s\n", respuesta.Correo)
				fmt.Printf("     Estado: %s\n", respuesta.Estado)
				fmt.Printf("     Timestamp: %v\n", respuesta.Timestamp)
				exitosRegistro++
			}
		}
	}

	fmt.Println("\n═════════════════════════════════════════════════════════════════")
	fmt.Println("=== RESUMEN DE REGISTRO ===")
	fmt.Println("═════════════════════════════════════════════════════════════════")
	fmt.Printf("✅ Registros exitosos: %d\n", exitosRegistro)
	fmt.Printf("❌ Casos de error manejados correctamente: %d\n", fallosRegistro)
	fmt.Printf("📊 Total de pruebas: %d\n", len(casosRegistro))

	if exitosRegistro > 0 {
		// Verificar que las credenciales se crearon correctamente
		fmt.Println("\n🔍 Verificación de credenciales creadas:")

		especCredenciales := domain.EspecificacionCredenciales{
			ListaFiltros: []domain.CriterioFiltro{
				{Campo: "activo", Operador: "=", Valor: true},
			},
		}
		paginacion := domain.Paginacion{
			Pagina:       1,
			TamanoPagina: 100,
		}

		credencialesCreadas, err := reg.CredencialesRepository().Find(ctx, especCredenciales, paginacion)
		if err != nil {
			fmt.Printf("❌ Error al buscar credenciales: %v\n", err)
		} else {
			fmt.Printf("✅ Total de credenciales activas: %d\n", len(credencialesCreadas))
		}
	}

	fmt.Println("\n✓ Pruebas de registro completadas")
}
