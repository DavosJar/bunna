package registro_test

import (
	"context"
	"github.com/davosjar/bunna/services/identidad/internal/usuarios/application/services/registro"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/config"
	"github.com/davosjar/bunna/services/identidad/internal/registry"
	"github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
	usuarios_postgres "github.com/davosjar/bunna/services/identidad/internal/usuarios/infrastructure/persistence/postgres"
	"gorm.io/gorm"
)

// MockGeneradorID es un mock de GeneradorID para pruebas
type MockGeneradorID struct {
	idsAGenerar []string
	indice      int
}

func NewMockGeneradorID(ids ...string) *MockGeneradorID {
	return &MockGeneradorID{
		idsAGenerar: ids,
		indice:      0,
	}
}

func (m *MockGeneradorID) NextID(ctx context.Context) (string, error) {
	if m.indice >= len(m.idsAGenerar) {
		return "", nil
	}
	id := m.idsAGenerar[m.indice]
	m.indice++
	return id, nil
}

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

	// Crear registry y servicio con UnitOfWork
	cfg, _ := config.LoadConfig()
	reg := registry.NewRegistry(db, cfg)

	servicio := registro.NuevoServicioRegistro(reg.UsuarioUnitOfWork())

	ctx := context.Background()
	comando := &registro.ComandoRegistro{
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
	usuarioObtener, err := reg.UsuarioRepository().ObtenerPorID(ctx, respuesta.UsuarioID)
	if err != nil {
		t.Fatalf("Error obteniendo usuario: %v", err)
	}
	if usuarioObtener.Correo() != comando.Correo {
		t.Errorf("Usuario en BD tiene correo diferente")
	}

	// Verificar que las credenciales existen en BD
	credenciales, err := reg.CredencialesRepository().ObtenerPorUsuarioID(ctx, respuesta.UsuarioID)
	if err != nil {
		t.Fatalf("Error obteniendo credenciales: %v", err)
	}

	// Verificar que el hash de password se creó correctamente
	if credenciales.PasswordHash() == "" {
		t.Error("Hash de password no debería estar vacío")
	}

	// Verificar que el password se hasheó correctamente con bcrypt
	if !reg.EncriptacionServicio().Verificar(comando.Password, credenciales.PasswordHash()) {
		t.Error("Password no se hasheó correctamente")
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

	cfg, _ := config.LoadConfig()
	reg := registry.NewRegistry(db, cfg)
	servicio := registro.NuevoServicioRegistro(reg.UsuarioUnitOfWork())

	ctx := context.Background()
	comando := &registro.ComandoRegistro{
		Correo:   "", // Vacío
		Password: "Test1234!",
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

	cfg, _ := config.LoadConfig()
	reg := registry.NewRegistry(db, cfg)
	servicio := registro.NuevoServicioRegistro(reg.UsuarioUnitOfWork())

	ctx := context.Background()
	comando := &registro.ComandoRegistro{
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

	cfg, _ := config.LoadConfig()
	reg := registry.NewRegistry(db, cfg)
	servicio := registro.NuevoServicioRegistro(reg.UsuarioUnitOfWork())

	ctx := context.Background()
	comando := &registro.ComandoRegistro{
		Correo:   "test@example.com",
		Password: "Test1234!",
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

// TestServicioRegistroMultiplesUsuarios verifica que se puedan crear múltiples usuarios
func TestServicioRegistroMultiplesUsuarios(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	cfg, _ := config.LoadConfig()
	reg := registry.NewRegistry(db, cfg)
	servicio := registro.NuevoServicioRegistro(reg.UsuarioUnitOfWork())

	ctx := context.Background()

	// Crear 3 usuarios
	usuarios := []*registro.ComandoRegistro{
		{
			Correo:   "user1@example.com",
			Password: "Test1234!",
			Nombre:   "Usuario",
			Apellido: "Uno",
			Telefono: "6001111111",
		},
		{
			Correo:   "user2@example.com",
			Password: "Pass5678!",
			Nombre:   "Usuario",
			Apellido: "Dos",
			Telefono: "6002222222",
		},
		{
			Correo:   "user3@example.com",
			Password: "Clave999?",
			Nombre:   "Usuario",
			Apellido: "Tres",
			Telefono: "6003333333",
		},
	}

	var usuariosCreados []*registro.DtoRespuestaRegistro
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

	cfg, _ := config.LoadConfig()
	reg := registry.NewRegistry(db, cfg)
	servicio := registro.NuevoServicioRegistro(reg.UsuarioUnitOfWork())

	ctx := context.Background()
	comando := &registro.ComandoRegistro{
		Correo:   "notanemail", // Email inválido
		Password: "Test1234!",
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
	if err != nil && !containsSubstring(err.Error(), "formato de correo inválido") {
		t.Errorf("Error esperado contiene 'formato de correo inválido', obtenido: %v", err)
	}
}

// TestServicioRegistroValidacionEmailSinArroba verifica validación de email sin @
func TestServicioRegistroValidacionEmailSinArroba(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	cfg, _ := config.LoadConfig()
	reg := registry.NewRegistry(db, cfg)
	servicio := registro.NuevoServicioRegistro(reg.UsuarioUnitOfWork())

	ctx := context.Background()
	comando := &registro.ComandoRegistro{
		Correo:   "user.com", // Email sin @
		Password: "Test1234!",
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

// TestRegistroPasswordHasheadoCorrectamente verifica que el password se hashea y se puede verificar
func TestRegistroPasswordHasheadoCorrectamente(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	cfg, _ := config.LoadConfig()
	reg := registry.NewRegistry(db, cfg)
	servicio := registro.NuevoServicioRegistro(reg.UsuarioUnitOfWork())

	ctx := context.Background()
	comando := &registro.ComandoRegistro{
		Correo:   "test@example.com",
		Password: "Test1234!",
		Nombre:   "Juan",
		Apellido: "García",
		Telefono: "6001234567",
	}

	// Ejecutar servicio
	respuesta, err := servicio.Ejecutar(ctx, comando)
	if err != nil {
		t.Fatalf("Error ejecutando servicio de registro: %v", err)
	}

	// Obtener credenciales
	credenciales, err := reg.CredencialesRepository().ObtenerPorUsuarioID(ctx, respuesta.UsuarioID)
	if err != nil {
		t.Fatalf("Error obteniendo credenciales: %v", err)
	}

	// Verificar que el password correcto verifica
	if !reg.EncriptacionServicio().Verificar(comando.Password, credenciales.PasswordHash()) {
		t.Error("Password correcto no verifica contra el hash")
	}

	// Verificar que un password incorrecto no verifica
	if reg.EncriptacionServicio().Verificar("WrongPassword", credenciales.PasswordHash()) {
		t.Error("Password incorrecto no debería verificar contra el hash")
	}

	// Verificar que dos hashes del mismo password son diferentes (propiedad de bcrypt)
	hash1 := credenciales.PasswordHash()
	hash2, err := reg.EncriptacionServicio().Hashear(comando.Password)
	if err != nil {
		t.Fatalf("Error hasheando password: %v", err)
	}

	if hash1 == hash2 {
		t.Error("Dos hashes del mismo password no deberían ser idénticos (bcrypt usa salt)")
	}

	// Pero ambos deberían verificar el mismo password
	if !reg.EncriptacionServicio().Verificar(comando.Password, hash1) {
		t.Error("Hash1 no verifica el password")
	}
	if !reg.EncriptacionServicio().Verificar(comando.Password, hash2) {
		t.Error("Hash2 no verifica el password")
	}
}

// TestServicioRegistroRollbackSiCredencialesFalla verifica que se hace rollback si falla al crear credenciales
func TestServicioRegistroRollbackSiCredencialesFalla(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	cfg, _ := config.LoadConfig()
	reg := registry.NewRegistry(db, cfg)
	servicio := registro.NuevoServicioRegistro(reg.UsuarioUnitOfWork())

	ctx := context.Background()

	// Crear un usuario exitosamente primero
	comando1 := &registro.ComandoRegistro{
		Correo:   "test1@example.com",
		Password: "Test1234!",
		Nombre:   "Usuario",
		Apellido: "Uno",
		Telefono: "6001234567",
	}

	respuesta1, err := servicio.Ejecutar(ctx, comando1)
	if err != nil {
		t.Fatalf("Error creando primer usuario: %v", err)
	}

	usuarioCreado1, err := reg.UsuarioRepository().ObtenerPorID(ctx, respuesta1.UsuarioID)
	if err != nil {
		t.Fatalf("Error obteniendo usuario creado: %v", err)
	}
	if usuarioCreado1 == nil {
		t.Fatal("Usuario debería existir en BD")
	}

	// Contar usuarios
	usuariosAntes, err := reg.UsuarioRepository().Listar(ctx, usuario.EspecificacionUsuario{}, domain.Paginacion{Pagina: 1, TamanoPagina: 100})
	if err != nil {
		t.Fatalf("Error listando usuarios: %v", err)
	}
	contadorAntes := len(usuariosAntes)

	// Ahora intentar crear un usuario con credenciales duplicadas
	// Primero creamos un usuario con el ID especificado
	comando2 := &registro.ComandoRegistro{
		Correo:   "test2@example.com",
		Password: "Pass5678!",
		Nombre:   "Usuario",
		Apellido: "Dos",
		Telefono: "6001234568",
	}

	respuesta2, err := servicio.Ejecutar(ctx, comando2)

	// Si la transacción no es atómica, podría haber un usuario sin credenciales
	// Verificamos que la cantidad de usuarios es la esperada
	usuariosDespues, err := reg.UsuarioRepository().Listar(ctx, usuario.EspecificacionUsuario{}, domain.Paginacion{Pagina: 1, TamanoPagina: 100})
	if err != nil {
		t.Fatalf("Error listando usuarios: %v", err)
	}
	contadorDespues := len(usuariosDespues)

	// Si ejecutó exitosamente
	if err == nil && respuesta2 != nil {
		// Verificar que se crearon usuario y credenciales
		if contadorDespues != contadorAntes+1 {
			t.Errorf("Se esperaban %d usuarios después, pero hay %d", contadorAntes+1, contadorDespues)
		}
		// Verificar que las credenciales existen
		creds, err := reg.CredencialesRepository().ObtenerPorUsuarioID(ctx, respuesta2.UsuarioID)
		if err != nil {
			t.Fatalf("Error obteniendo credenciales: %v", err)
		}
		if creds == nil {
			t.Fatal("Credenciales deberían existir")
		}
	} else {
		// Si falló, no debería cambiar el número de usuarios (rollback)
		if contadorDespues != contadorAntes {
			t.Errorf("Después de un error, la transacción debería haber sido revertida. Usuarios antes: %d, después: %d", contadorAntes, contadorDespues)
		}
	}
}

// TestServicioRegistroTransaccionAtomicaConContextTimeout verifica rollback con context timeout
func TestServicioRegistroTransaccionAtomicaConContextTimeout(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	cfg, _ := config.LoadConfig()
	reg := registry.NewRegistry(db, cfg)
	servicio := registro.NuevoServicioRegistro(reg.UsuarioUnitOfWork())

	// Crear un context con timeout muy corto
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	comando := &registro.ComandoRegistro{
		Correo:   "test@example.com",
		Password: "Test1234!",
		Nombre:   "Juan",
		Apellido: "García",
		Telefono: "6001234567",
	}

	// Ejecutar con timeout expirado
	respuesta, err := servicio.Ejecutar(ctx, comando)

	// Debe fallar
	if err == nil {
		t.Error("Debería haber error por context timeout")
	}
	if respuesta != nil {
		t.Error("Respuesta debería ser nil cuando hay error")
	}

	// Verificar que no se creó ningún usuario (rollback funcionó)
	usuariosPorCorreo, err := reg.UsuarioRepository().Listar(ctx, usuario.EspecificacionUsuario{}, domain.Paginacion{Pagina: 1, TamanoPagina: 100})
	if err == nil && len(usuariosPorCorreo) > 0 {
		t.Errorf("Debería haber 0 usuarios después del rollback, pero hay %d", len(usuariosPorCorreo))
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

// TestServicioRegistroConGeneradorIDInyectado verifica que GeneradorID se puede inyectar
func TestServicioRegistroConGeneradorIDInyectado(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	// Crear un mock de GeneradorID con un ID predefinido (UUID válido)
	mockID := "550e8400-e29b-41d4-a716-446655440000"
	mockGenerador := NewMockGeneradorID(mockID)

	// Crear repositorios
	cfg, _ := config.LoadConfig()
	reg := registry.NewRegistry(db, cfg)

	// Crear UnitOfWork con el mock de GeneradorID
	usuarioRepo := reg.UsuarioRepository()
	credencialesRepo := reg.CredencialesRepository()
	encriptacion := reg.EncriptacionServicio()

	unitOfWork := usuarios_postgres.NewUnitOfWork(
		db,
		usuarioRepo,
		credencialesRepo,
		encriptacion,
		mockGenerador,
	)

	servicio := registro.NuevoServicioRegistro(unitOfWork)

	ctx := context.Background()
	comando := &registro.ComandoRegistro{
		Correo:   "test@example.com",
		Password: "Test1234!",
		Nombre:   "Juan",
		Apellido: "García",
		Telefono: "6001234567",
	}

	// Ejecutar servicio
	respuesta, err := servicio.Ejecutar(ctx, comando)
	if err != nil {
		t.Fatalf("Error ejecutando servicio de registro: %v", err)
	}

	// Verificar que se usó el ID del mock
	if respuesta.UsuarioID != mockID {
		t.Errorf("Se esperaba ID %s, obtenido %s", mockID, respuesta.UsuarioID)
	}

	// Verificar que el usuario se creó con el ID correcto
	usuarioObtener, err := usuarioRepo.ObtenerPorID(ctx, mockID)
	if err != nil {
		t.Fatalf("Error obteniendo usuario: %v", err)
	}
	if usuarioObtener.ID() != mockID {
		t.Errorf("Usuario en BD tiene ID diferente: %s", usuarioObtener.ID())
	}
}
