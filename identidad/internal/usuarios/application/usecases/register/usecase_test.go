package register_test

import (
	"context"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/config"
	"github.com/davosjar/bunna/services/identidad/internal/registry"
	"github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/register"
	"github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
	usuarios_postgres "github.com/davosjar/bunna/services/identidad/internal/usuarios/infrastructure/persistence/postgres"
	"gorm.io/gorm"
)

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

func setupTestDB(t *testing.T) *gorm.DB {
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	db, err := config.InitDB(cfg.GetDSN())
	if err != nil {
		t.Fatalf("Error initializing database: %v", err)
	}

	db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	if err := config.RunMigrations(db); err != nil {
		t.Fatalf("Error running migrations: %v", err)
	}

	return db
}

func TestRegistrarUsuarioExitoso(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	cfg, _ := config.LoadConfig()
	reg := registry.NewRegistry(db, cfg)

	uc := register.NewRegistrarUsuarioCasoDeUso(reg.UsuarioUnitOfWork())

	ctx := context.Background()
	cmd := &register.ComandoRegistrarUsuario{
		Correo:   "test@example.com",
		Password: "password123",
		Nombre:   "Juan",
		Apellido: "García",
		Telefono: "6001234567",
	}

	respuesta, err := uc.Ejecutar(ctx, cmd)
	if err != nil {
		t.Fatalf("Error ejecutando caso de uso: %v", err)
	}

	if respuesta.UsuarioID == "" {
		t.Error("UsuarioID no debe estar vacío")
	}
	if respuesta.Correo != cmd.Correo {
		t.Errorf("Correo esperado %s, obtenido %s", cmd.Correo, respuesta.Correo)
	}
	if respuesta.Estado != string(usuario.NO_VERIFICADO) {
		t.Errorf("Estado esperado %s, obtenido %s", usuario.NO_VERIFICADO, respuesta.Estado)
	}

	usuarioObtener, err := reg.UsuarioRepository().ObtenerPorID(ctx, respuesta.UsuarioID)
	if err != nil {
		t.Fatalf("Error obteniendo usuario: %v", err)
	}
	if usuarioObtener.Correo() != cmd.Correo {
		t.Errorf("Usuario en BD tiene correo diferente")
	}

	credenciales, err := reg.CredencialesRepository().ObtenerPorUsuarioID(ctx, respuesta.UsuarioID)
	if err != nil {
		t.Fatalf("Error obteniendo credenciales: %v", err)
	}

	if credenciales.PasswordHash() == "" {
		t.Error("Hash de password no debería estar vacío")
	}

	if !reg.EncriptacionServicio().Verificar(cmd.Password, credenciales.PasswordHash()) {
		t.Error("Password no se hasheó correctamente")
	}

	if !credenciales.Activo() {
		t.Error("Credenciales deberían estar activas")
	}
	if credenciales.CorreoVerificado() {
		t.Error("Correo no debería estar verificado inicialmente")
	}
}

func TestRegistrarUsuarioCorreoVacio(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	cfg, _ := config.LoadConfig()
	reg := registry.NewRegistry(db, cfg)
	uc := register.NewRegistrarUsuarioCasoDeUso(reg.UsuarioUnitOfWork())

	ctx := context.Background()
	cmd := &register.ComandoRegistrarUsuario{
		Correo:   "",
		Password: "password123",
		Nombre:   "Juan",
		Apellido: "García",
		Telefono: "6001234567",
	}

	respuesta, err := uc.Ejecutar(ctx, cmd)
	if err == nil {
		t.Error("Debería haber error por correo vacío")
	}
	if respuesta != nil {
		t.Error("Respuesta debería ser nil cuando hay error de validación")
	}
}

func TestRegistrarUsuarioPasswordVacio(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	cfg, _ := config.LoadConfig()
	reg := registry.NewRegistry(db, cfg)
	uc := register.NewRegistrarUsuarioCasoDeUso(reg.UsuarioUnitOfWork())

	ctx := context.Background()
	cmd := &register.ComandoRegistrarUsuario{
		Correo:   "test@example.com",
		Password: "",
		Nombre:   "Juan",
		Apellido: "García",
		Telefono: "6001234567",
	}

	respuesta, err := uc.Ejecutar(ctx, cmd)
	if err == nil {
		t.Error("Debería haber error por password vacío")
	}
	if respuesta != nil {
		t.Error("Respuesta debería ser nil")
	}
}

func TestRegistrarUsuarioNombreVacio(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	cfg, _ := config.LoadConfig()
	reg := registry.NewRegistry(db, cfg)
	uc := register.NewRegistrarUsuarioCasoDeUso(reg.UsuarioUnitOfWork())

	ctx := context.Background()
	cmd := &register.ComandoRegistrarUsuario{
		Correo:   "test@example.com",
		Password: "password123",
		Nombre:   "",
		Apellido: "García",
		Telefono: "6001234567",
	}

	respuesta, err := uc.Ejecutar(ctx, cmd)
	if err == nil {
		t.Error("Debería haber error por nombre vacío")
	}
	if respuesta != nil {
		t.Error("Respuesta debería ser nil")
	}
}

func TestRegistrarUsuarioMultiplesUsuarios(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	cfg, _ := config.LoadConfig()
	reg := registry.NewRegistry(db, cfg)
	uc := register.NewRegistrarUsuarioCasoDeUso(reg.UsuarioUnitOfWork())

	ctx := context.Background()

	usuarios := []*register.ComandoRegistrarUsuario{
		{Correo: "user1@example.com", Password: "pass1", Nombre: "Usuario", Apellido: "Uno", Telefono: "6001111111"},
		{Correo: "user2@example.com", Password: "pass2", Nombre: "Usuario", Apellido: "Dos", Telefono: "6002222222"},
		{Correo: "user3@example.com", Password: "pass3", Nombre: "Usuario", Apellido: "Tres", Telefono: "6003333333"},
	}

	var usuariosCreados []*register.RespuestaRegistrarUsuario
	for _, cmd := range usuarios {
		respuesta, err := uc.Ejecutar(ctx, cmd)
		if err != nil {
			t.Fatalf("Error creando usuario: %v", err)
		}
		usuariosCreados = append(usuariosCreados, respuesta)
	}

	if len(usuariosCreados) != 3 {
		t.Errorf("Se esperaban 3 usuarios, se crearon %d", len(usuariosCreados))
	}

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

func TestRegistrarUsuarioEmailInvalido(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	cfg, _ := config.LoadConfig()
	reg := registry.NewRegistry(db, cfg)
	uc := register.NewRegistrarUsuarioCasoDeUso(reg.UsuarioUnitOfWork())

	ctx := context.Background()
	cmd := &register.ComandoRegistrarUsuario{
		Correo:   "notanemail",
		Password: "password123",
		Nombre:   "Juan",
		Apellido: "García",
		Telefono: "6001234567",
	}

	respuesta, err := uc.Ejecutar(ctx, cmd)
	if err == nil {
		t.Error("Debería haber error por email inválido")
	}
	if respuesta != nil {
		t.Error("Respuesta debería ser nil")
	}
}

func TestRegistrarUsuarioPasswordHasheadoCorrectamente(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	cfg, _ := config.LoadConfig()
	reg := registry.NewRegistry(db, cfg)
	uc := register.NewRegistrarUsuarioCasoDeUso(reg.UsuarioUnitOfWork())

	ctx := context.Background()
	cmd := &register.ComandoRegistrarUsuario{
		Correo:   "test@example.com",
		Password: "MySecurePassword123!",
		Nombre:   "Juan",
		Apellido: "García",
		Telefono: "6001234567",
	}

	respuesta, err := uc.Ejecutar(ctx, cmd)
	if err != nil {
		t.Fatalf("Error ejecutando caso de uso: %v", err)
	}

	credenciales, err := reg.CredencialesRepository().ObtenerPorUsuarioID(ctx, respuesta.UsuarioID)
	if err != nil {
		t.Fatalf("Error obteniendo credenciales: %v", err)
	}

	if !reg.EncriptacionServicio().Verificar(cmd.Password, credenciales.PasswordHash()) {
		t.Error("Password correcto no verifica contra el hash")
	}

	if reg.EncriptacionServicio().Verificar("WrongPassword", credenciales.PasswordHash()) {
		t.Error("Password incorrecto no debería verificar")
	}

	hash1 := credenciales.PasswordHash()
	hash2, err := reg.EncriptacionServicio().Hashear(cmd.Password)
	if err != nil {
		t.Fatalf("Error hasheando password: %v", err)
	}

	if hash1 == hash2 {
		t.Error("Dos hashes del mismo password no deberían ser idénticos (bcrypt usa salt)")
	}

	if !reg.EncriptacionServicio().Verificar(cmd.Password, hash1) {
		t.Error("Hash1 no verifica el password")
	}
	if !reg.EncriptacionServicio().Verificar(cmd.Password, hash2) {
		t.Error("Hash2 no verifica el password")
	}
}

func TestRegistrarUsuarioTransaccionAtomica(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	cfg, _ := config.LoadConfig()
	reg := registry.NewRegistry(db, cfg)
	uc := register.NewRegistrarUsuarioCasoDeUso(reg.UsuarioUnitOfWork())

	ctx := context.Background()

	cmd1 := &register.ComandoRegistrarUsuario{
		Correo: "test1@example.com", Password: "password123",
		Nombre: "Usuario", Apellido: "Uno", Telefono: "6001234567",
	}

	respuesta1, err := uc.Ejecutar(ctx, cmd1)
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

	usuariosAntes, err := reg.UsuarioRepository().Listar(ctx, usuario.EspecificacionUsuario{}, domain.Paginacion{Pagina: 1, TamanoPagina: 100})
	if err != nil {
		t.Fatalf("Error listando usuarios: %v", err)
	}
	contadorAntes := len(usuariosAntes)

	cmd2 := &register.ComandoRegistrarUsuario{
		Correo: "test2@example.com", Password: "password456",
		Nombre: "Usuario", Apellido: "Dos", Telefono: "6001234568",
	}

	respuesta2, err := uc.Ejecutar(ctx, cmd2)

	usuariosDespues, err := reg.UsuarioRepository().Listar(ctx, usuario.EspecificacionUsuario{}, domain.Paginacion{Pagina: 1, TamanoPagina: 100})
	if err != nil {
		t.Fatalf("Error listando usuarios: %v", err)
	}
	contadorDespues := len(usuariosDespues)

	if err == nil && respuesta2 != nil {
		if contadorDespues != contadorAntes+1 {
			t.Errorf("Se esperaban %d usuarios después, pero hay %d", contadorAntes+1, contadorDespues)
		}
	} else {
		if contadorDespues != contadorAntes {
			t.Errorf("Transacción debería haber revertido. Usuarios antes: %d, después: %d", contadorAntes, contadorDespues)
		}
	}
}

func TestRegistrarUsuarioConGeneradorIDInyectado(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		db.Migrator().DropTable("credenciales_usuarios", "usuarios")
	}()

	mockID := "550e8400-e29b-41d4-a716-446655440000"
	mockGenerador := NewMockGeneradorID(mockID)

	cfg, _ := config.LoadConfig()
	reg := registry.NewRegistry(db, cfg)

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

	uc := register.NewRegistrarUsuarioCasoDeUso(unitOfWork)

	ctx := context.Background()
	cmd := &register.ComandoRegistrarUsuario{
		Correo: "test@example.com", Password: "password123",
		Nombre: "Juan", Apellido: "García", Telefono: "6001234567",
	}

	respuesta, err := uc.Ejecutar(ctx, cmd)
	if err != nil {
		t.Fatalf("Error ejecutando caso de uso: %v", err)
	}

	if respuesta.UsuarioID != mockID {
		t.Errorf("Se esperaba ID %s, obtenido %s", mockID, respuesta.UsuarioID)
	}

	usuarioObtener, err := usuarioRepo.ObtenerPorID(ctx, mockID)
	if err != nil {
		t.Fatalf("Error obteniendo usuario: %v", err)
	}
	if usuarioObtener.ID() != mockID {
		t.Errorf("Usuario en BD tiene ID diferente: %s", usuarioObtener.ID())
	}
}
