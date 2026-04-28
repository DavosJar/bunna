package registro

import (
	"context"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/config"
	"github.com/davosjar/bunna/services/identidad/internal/domain/usuario"
	"github.com/davosjar/bunna/services/identidad/internal/registry"
	"gorm.io/gorm"
)

// setupTestDB configura la BD de prueba
func setupTestDB(t *testing.T) *gorm.DB {
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	db, err := config.InitDB(cfg.GetDSN())
	if err != nil {
		t.Fatalf("Error initializing database: %v", err)
	}

	// Limpiar y migrar
	db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	if err := config.RunMigrations(db); err != nil {
		t.Fatalf("Error running migrations: %v", err)
	}

	return db
}

// TestServicioRegistroExitoso verifica que se creen correctamente usuario y credenciales
func TestServicioRegistroExitoso(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	// Crear registry y servicio
	reg := registry.NewRegistry(db)
	usuarioRepo := reg.UsuarioRepository()
	credencialesRepo := reg.CredencialesRepository()

	servicio := NuevoServicioRegistro(usuarioRepo, credencialesRepo, db)

	ctx := context.Background()
	comando := &ComandoRegistro{
		Correo:   "test@example.com",
		Password: "password123",
		Nombre:   "Juan",
		Apellido: "García",
		Telefono: "6001234567",
	}

	// Ejecutar servicio
	respuesta, err := servicio.Ejecutar(ctx, comando)
	if err != nil {
		t.Fatalf("Error ejecutando servicio de registro: %v", err)
	}

	// Verificaciones
	if respuesta.UsuarioID == "" {
		t.Error("UsuarioID no debe estar vacío")
	}
	if respuesta.Correo != comando.Correo {
		t.Errorf("Correo esperado %s, obtenido %s", comando.Correo, respuesta.Correo)
	}
	if respuesta.Estado != string(usuario.NO_VERIFICADO) {
		t.Errorf("Estado esperado %s, obtenido %s", usuario.NO_VERIFICADO, respuesta.Estado)
	}

	// Verificar que el usuario existe en BD
	usuarioObtener, err := usuarioRepo.ObtenerPorID(ctx, respuesta.UsuarioID)
	if err != nil {
		t.Fatalf("Error obteniendo usuario: %v", err)
	}
	if usuarioObtener.Correo() != comando.Correo {
		t.Errorf("Usuario en BD tiene correo diferente")
	}

	// Verificar que las credenciales existen en BD
	credenciales, err := credencialesRepo.ObtenerPorUsuarioID(ctx, respuesta.UsuarioID)
	if err != nil {
		t.Fatalf("Error obteniendo credenciales: %v", err)
	}

	// Verificar hash de password
	hashEsperado := hashPassword(comando.Password)
	if credenciales.PasswordHash() != hashEsperado {
		t.Errorf("Hash de password no coincide")
	}
	if !credenciales.Activo() {
		t.Error("Credenciales deberían estar activas")
	}
	if credenciales.CorreoVerificado() {
		t.Error("Correo no debería estar verificado inicialmente")
	}
}

// TestServicioRegistroValidacionCorreoVacio verifica validación de correo
func TestServicioRegistroValidacionCorreoVacio(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	reg := registry.NewRegistry(db)
	servicio := NuevoServicioRegistro(
		reg.UsuarioRepository(),
		reg.CredencialesRepository(),
		db,
	)

	ctx := context.Background()
	comando := &ComandoRegistro{
		Correo:   "", // Vacío
		Password: "password123",
		Nombre:   "Juan",
		Apellido: "García",
		Telefono: "6001234567",
	}

	respuesta, err := servicio.Ejecutar(ctx, comando)
	if err == nil {
		t.Error("Debería haber error por correo vacío")
	}
	if respuesta != nil {
		t.Error("Respuesta debería ser nil cuando hay error de validación")
	}
}

// TestServicioRegistroValidacionPasswordVacio verifica validación de password
func TestServicioRegistroValidacionPasswordVacio(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	reg := registry.NewRegistry(db)
	servicio := NuevoServicioRegistro(
		reg.UsuarioRepository(),
		reg.CredencialesRepository(),
		db,
	)

	ctx := context.Background()
	comando := &ComandoRegistro{
		Correo:   "test@example.com",
		Password: "", // Vacío
		Nombre:   "Juan",
		Apellido: "García",
		Telefono: "6001234567",
	}

	respuesta, err := servicio.Ejecutar(ctx, comando)
	if err == nil {
		t.Error("Debería haber error por password vacío")
	}
	if respuesta != nil {
		t.Error("Respuesta debería ser nil cuando hay error de validación")
	}
}

// TestServicioRegistroValidacionNombreVacio verifica validación de nombre
func TestServicioRegistroValidacionNombreVacio(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	reg := registry.NewRegistry(db)
	servicio := NuevoServicioRegistro(
		reg.UsuarioRepository(),
		reg.CredencialesRepository(),
		db,
	)

	ctx := context.Background()
	comando := &ComandoRegistro{
		Correo:   "test@example.com",
		Password: "password123",
		Nombre:   "", // Vacío
		Apellido: "García",
		Telefono: "6001234567",
	}

	respuesta, err := servicio.Ejecutar(ctx, comando)
	if err == nil {
		t.Error("Debería haber error por nombre vacío")
	}
	if respuesta != nil {
		t.Error("Respuesta debería ser nil cuando hay error de validación")
	}
}

// TestServicioRegistroRollbackPorErrorCredenciales verifica rollback si falla creación de credenciales
func TestServicioRegistroRollbackPorErrorCredenciales(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	reg := registry.NewRegistry(db)
	servicio := NuevoServicioRegistro(
		reg.UsuarioRepository(),
		reg.CredencialesRepository(),
		db,
	)

	ctx := context.Background()

	// Crear un usuario exitosamente
	comando1 := &ComandoRegistro{
		Correo:   "test1@example.com",
		Password: "password123",
		Nombre:   "Juan",
		Apellido: "García",
		Telefono: "6001234567",
	}

	respuesta1, err := servicio.Ejecutar(ctx, comando1)
	if err != nil {
		t.Fatalf("Error ejecutando primer registro: %v", err)
	}

	// Verificar que el usuario existe
	usuarioObtener, err := reg.UsuarioRepository().ObtenerPorID(ctx, respuesta1.UsuarioID)
	if err != nil {
		t.Fatalf("Error obteniendo usuario: %v", err)
	}
	if usuarioObtener.Correo() != comando1.Correo {
		t.Error("El primer usuario no se creó correctamente")
	}

	// Crear otro usuario exitosamente
	comando2 := &ComandoRegistro{
		Correo:   "test2@example.com",
		Password: "password456",
		Nombre:   "Pedro",
		Apellido: "López",
		Telefono: "6009876543",
	}

	respuesta2, err := servicio.Ejecutar(ctx, comando2)
	if err != nil {
		t.Fatalf("Error ejecutando segundo registro: %v", err)
	}

	// Verificar que ambos usuarios existen con credenciales
	credenciales1, err := reg.CredencialesRepository().ObtenerPorUsuarioID(ctx, respuesta1.UsuarioID)
	if err != nil {
		t.Fatalf("Error obteniendo credenciales del usuario 1: %v", err)
	}

	credenciales2, err := reg.CredencialesRepository().ObtenerPorUsuarioID(ctx, respuesta2.UsuarioID)
	if err != nil {
		t.Fatalf("Error obteniendo credenciales del usuario 2: %v", err)
	}

	// Verificar que ambos tienen los hashes correctos
	if credenciales1.PasswordHash() != hashPassword(comando1.Password) {
		t.Error("Hash de password del usuario 1 no coincide")
	}
	if credenciales2.PasswordHash() != hashPassword(comando2.Password) {
		t.Error("Hash de password del usuario 2 no coincide")
	}
}

// TestServicioRegistroMultiplesUsuarios verifica que se puedan crear múltiples usuarios
func TestServicioRegistroMultiplesUsuarios(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	reg := registry.NewRegistry(db)
	servicio := NuevoServicioRegistro(
		reg.UsuarioRepository(),
		reg.CredencialesRepository(),
		db,
	)

	ctx := context.Background()

	// Crear 3 usuarios
	usuarios := []*ComandoRegistro{
		{
			Correo:   "user1@example.com",
			Password: "pass1",
			Nombre:   "Usuario",
			Apellido: "Uno",
			Telefono: "6001111111",
		},
		{
			Correo:   "user2@example.com",
			Password: "pass2",
			Nombre:   "Usuario",
			Apellido: "Dos",
			Telefono: "6002222222",
		},
		{
			Correo:   "user3@example.com",
			Password: "pass3",
			Nombre:   "Usuario",
			Apellido: "Tres",
			Telefono: "6003333333",
		},
	}

	var usuariosCreados []*DtoRespuestaRegistro
	for _, cmd := range usuarios {
		respuesta, err := servicio.Ejecutar(ctx, cmd)
		if err != nil {
			t.Fatalf("Error creando usuario: %v", err)
		}
		usuariosCreados = append(usuariosCreados, respuesta)
	}

	if len(usuariosCreados) != 3 {
		t.Errorf("Se esperaban 3 usuarios, se crearon %d", len(usuariosCreados))
	}

	// Verificar que todos los usuarios existen en BD
	for i, respuesta := range usuariosCreados {
		usuarioObtener, err := reg.UsuarioRepository().ObtenerPorID(ctx, respuesta.UsuarioID)
		if err != nil {
			t.Fatalf("Error obteniendo usuario %d: %v", i+1, err)
		}
		if usuarioObtener.Correo() != usuarios[i].Correo {
			t.Errorf("Usuario %d tiene correo diferente", i+1)
		}
	}
}

// TestServicioRegistroValidacionEmailInvalido verifica validación de email con formato inválido
func TestServicioRegistroValidacionEmailInvalido(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	reg := registry.NewRegistry(db)
	servicio := NuevoServicioRegistro(
		reg.UsuarioRepository(),
		reg.CredencialesRepository(),
		db,
	)

	ctx := context.Background()
	comando := &ComandoRegistro{
		Correo:   "notanemail", // Email inválido
		Password: "password123",
		Nombre:   "Juan",
		Apellido: "García",
		Telefono: "6001234567",
	}

	respuesta, err := servicio.Ejecutar(ctx, comando)
	if err == nil {
		t.Error("Debería haber error por email con formato inválido")
	}
	if respuesta != nil {
		t.Error("Respuesta debería ser nil cuando hay error de validación")
	}
	if err != nil && err.Error() != "formato de correo inválido: mail: missing '@'" {
		// Verificar que el error contiene el texto esperado
		if !containsSubstring(err.Error(), "formato de correo inválido") {
			t.Errorf("Error esperado contiene 'formato de correo inválido', obtenido: %v", err)
		}
	}
}

// TestServicioRegistroValidacionEmailSinArroba verifica validación de email sin @
func TestServicioRegistroValidacionEmailSinArroba(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	reg := registry.NewRegistry(db)
	servicio := NuevoServicioRegistro(
		reg.UsuarioRepository(),
		reg.CredencialesRepository(),
		db,
	)

	ctx := context.Background()
	comando := &ComandoRegistro{
		Correo:   "user.com", // Email sin @
		Password: "password123",
		Nombre:   "Juan",
		Apellido: "García",
		Telefono: "6001234567",
	}

	respuesta, err := servicio.Ejecutar(ctx, comando)
	if err == nil {
		t.Error("Debería haber error por email sin @")
	}
	if respuesta != nil {
		t.Error("Respuesta debería ser nil cuando hay error de validación")
	}
	if err != nil && !containsSubstring(err.Error(), "formato de correo inválido") {
		t.Errorf("Error esperado contiene 'formato de correo inválido', obtenido: %v", err)
	}
}

// containsSubstring es una función auxiliar para verificar si una cadena contiene una subcadena
func containsSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
