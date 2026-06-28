package register_test

import (
	"context"
	"errors"
	"testing"

	rbac "github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	"github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/davosjar/bunna/services/identidad/internal/tenants/domain/tenant"
	"github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/register"
	usuariodomain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

// ── Mock UnitOfWork ────────────────────────────────────────────────────────────

type mockUnitOfWork struct {
	usuarioRepo          usuariodomain.UsuarioRepositorio
	credencialesRepo     domain.CredencialesRepositorio
	tenantRepo           tenant.TenantRepositorio
	membresiaRepo        tenant.MembresiaRepositorio
	rolRepo              rbac.RolRepositorio
	usuarioTenantRolRepo rbac.UsuarioTenantRolRepositorio
	rolPermisoRepo       rbac.RolPermisoRepositorio
}

func (m *mockUnitOfWork) Transaccional(ctx context.Context, fn func(tx register.UnitOfWork) error) error {
	return fn(m)
}

func (m *mockUnitOfWork) UsuarioRepository() usuariodomain.UsuarioRepositorio {
	return m.usuarioRepo
}

func (m *mockUnitOfWork) CredencialesRepository() domain.CredencialesRepositorio {
	return m.credencialesRepo
}

func (m *mockUnitOfWork) TenantRepository() tenant.TenantRepositorio {
	return m.tenantRepo
}

func (m *mockUnitOfWork) MembresiaRepository() tenant.MembresiaRepositorio {
	return m.membresiaRepo
}

func (m *mockUnitOfWork) RolRepository() rbac.RolRepositorio {
	return m.rolRepo
}

func (m *mockUnitOfWork) UsuarioTenantRolRepository() rbac.UsuarioTenantRolRepositorio {
	return m.usuarioTenantRolRepo
}

func (m *mockUnitOfWork) RolPermisoRepository() rbac.RolPermisoRepositorio {
	return m.rolPermisoRepo
}

// ── Other Mocks ────────────────────────────────────────────────────────────────

type mockUsuarioRepo struct {
	crearFunc func(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error)
}

func (m *mockUsuarioRepo) Crear(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
	return m.crearFunc(ctx, u)
}
func (m *mockUsuarioRepo) Actualizar(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
	return u, nil
}
func (m *mockUsuarioRepo) Eliminar(ctx context.Context, id string) error { return nil }
func (m *mockUsuarioRepo) ObtenerPorID(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
	return nil, nil
}
func (m *mockUsuarioRepo) ObtenerPorCorreo(ctx context.Context, correo string) (*usuariodomain.Usuario, error) {
	return nil, nil
}
func (m *mockUsuarioRepo) Listar(ctx context.Context, _ usuariodomain.EspecificacionUsuario, _ shareddomain.Paginacion) ([]*usuariodomain.Usuario, error) {
	return nil, nil
}

type mockCredencialesRepo struct {
	crearFunc func(ctx context.Context, c *domain.CredencialesUsuario) (*domain.CredencialesUsuario, error)
}

func (m *mockCredencialesRepo) Crear(ctx context.Context, c *domain.CredencialesUsuario) (*domain.CredencialesUsuario, error) {
	return m.crearFunc(ctx, c)
}
func (m *mockCredencialesRepo) Actualizar(ctx context.Context, c *domain.CredencialesUsuario) (*domain.CredencialesUsuario, error) {
	return c, nil
}
func (m *mockCredencialesRepo) ObtenerPorUsuarioID(ctx context.Context, usuarioID string) (*domain.CredencialesUsuario, error) {
	return nil, nil
}
func (m *mockCredencialesRepo) Eliminar(ctx context.Context, usuarioID string) error { return nil }
func (m *mockCredencialesRepo) Find(ctx context.Context, _ domain.EspecificacionCredenciales, _ shareddomain.Paginacion) ([]*domain.CredencialesUsuario, error) {
	return nil, nil
}

type mockEncriptacion struct {
	hash string
}

func (m *mockEncriptacion) Hashear(password string) (string, error) {
	return m.hash, nil
}
func (m *mockEncriptacion) Verificar(password, hash string) bool {
	return m.hash == hash
}

type mockGeneradorID struct {
	ids []string
	idx int
}

func (m *mockGeneradorID) NextID(ctx context.Context) (string, error) {
	id := m.ids[m.idx]
	m.idx++
	return id, nil
}

type mockTenantRepo struct {
	crearFunc      func(ctx context.Context, t *tenant.Tenant) (*tenant.Tenant, error)
	obtenerSlugErr error
}

func (m *mockTenantRepo) Crear(ctx context.Context, t *tenant.Tenant) (*tenant.Tenant, error) {
	return m.crearFunc(ctx, t)
}
func (m *mockTenantRepo) ObtenerPorID(ctx context.Context, id string) (*tenant.Tenant, error) {
	return nil, nil
}
func (m *mockTenantRepo) ObtenerPorSlug(ctx context.Context, slug string) (*tenant.Tenant, error) {
	if m.obtenerSlugErr != nil {
		return nil, m.obtenerSlugErr
	}
	return nil, tenant.ErrTenantNoEncontrado
}
func (m *mockTenantRepo) Actualizar(ctx context.Context, t *tenant.Tenant) (*tenant.Tenant, error) {
	return t, nil
}
func (m *mockTenantRepo) Listar(ctx context.Context) ([]*tenant.Tenant, error) { return nil, nil }
func (m *mockTenantRepo) ListarPorUsuario(ctx context.Context, usuarioID string) ([]*tenant.Tenant, error) {
	return nil, nil
}

type mockMembresiaRepo struct {
	crearFunc func(ctx context.Context, m *tenant.Membresia) error
}

func (m *mockMembresiaRepo) Crear(ctx context.Context, memb *tenant.Membresia) error {
	return m.crearFunc(ctx, memb)
}
func (m *mockMembresiaRepo) Eliminar(ctx context.Context, usuarioID, tenantID string) error {
	return nil
}
func (m *mockMembresiaRepo) ExisteMiembro(ctx context.Context, usuarioID, tenantID string) (bool, error) {
	return false, nil
}
func (m *mockMembresiaRepo) ListarUsuariosPorTenant(ctx context.Context, tenantID string) ([]string, error) {
	return nil, nil
}
func (m *mockMembresiaRepo) ListarTenantsPorUsuario(ctx context.Context, usuarioID string) ([]string, error) {
	return nil, nil
}

type mockRolRepo struct {
	obtenerPorNombreFunc func(ctx context.Context, nombre string) (*rbac.RolDB, error)
}

func (m *mockRolRepo) ObtenerPorNombre(ctx context.Context, nombre string) (*rbac.RolDB, error) {
	return m.obtenerPorNombreFunc(ctx, nombre)
}
func (m *mockRolRepo) ObtenerPorID(ctx context.Context, id string) (*rbac.RolDB, error) {
	return &rbac.RolDB{ID: id}, nil
}
func (m *mockRolRepo) Listar(ctx context.Context, _ rbac.EspecificacionRol, _ shareddomain.Paginacion) ([]*rbac.RolDB, error) {
	return nil, nil
}
func (m *mockRolRepo) Crear(ctx context.Context, _ *rbac.RolDB) error { return nil }
func (m *mockRolRepo) ActualizarDescripcion(ctx context.Context, _, _ string) error { return nil }

type mockUsuarioTenantRolRepo struct {
	crearFunc func(ctx context.Context, usuarioID, tenantID, rolID string) error
}

func (m *mockUsuarioTenantRolRepo) Crear(ctx context.Context, usuarioID, tenantID, rolID string) error {
	return m.crearFunc(ctx, usuarioID, tenantID, rolID)
}
func (m *mockUsuarioTenantRolRepo) Eliminar(ctx context.Context, _, _ string, _ string) error { return nil }
func (m *mockUsuarioTenantRolRepo) ListarRolesPorUsuarioEnTenant(ctx context.Context, _, _ string) ([]*rbac.RolDB, error) {
	return nil, nil
}
func (m *mockUsuarioTenantRolRepo) TieneRolEnTenant(ctx context.Context, _, _, _ string) (bool, error) {
	return false, nil
}

type mockRolPermisoRepo struct {
	listarPorRolYTenantFunc func(ctx context.Context, rolID, tenantID string) ([]*rbac.PermisoDB, error)
}

func (m *mockRolPermisoRepo) AsignarPermiso(ctx context.Context, rolID, permisoID, tenantID, asignadoPor string) error {
	return nil
}
func (m *mockRolPermisoRepo) EliminarPermiso(ctx context.Context, rolID, permisoID, tenantID string) error {
	return nil
}
func (m *mockRolPermisoRepo) ListarPorRolYTenant(ctx context.Context, rolID, tenantID string) ([]*rbac.PermisoDB, error) {
	if m.listarPorRolYTenantFunc != nil {
		return m.listarPorRolYTenantFunc(ctx, rolID, tenantID)
	}
	return nil, nil
}

type mockRolPublisher struct {
	publicarRolActualizadoFunc func(ctx context.Context, rolID, tenantID string, permisos []string) error
}

func (m *mockRolPublisher) PublicarRolActualizado(ctx context.Context, rolID, tenantID string, permisos []string) error {
	if m.publicarRolActualizadoFunc != nil {
		return m.publicarRolActualizadoFunc(ctx, rolID, tenantID, permisos)
	}
	return nil
}

// ── Helpers ────────────────────────────────────────────────────────────────────

func newUCSuccess() *register.RegistrarUsuarioCasoDeUso {
	return register.NewRegistrarUsuarioCasoDeUso(
		&mockUnitOfWork{
			usuarioRepo: &mockUsuarioRepo{
				crearFunc: func(_ context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
					return u, nil
				},
			},
			credencialesRepo: &mockCredencialesRepo{
				crearFunc: func(_ context.Context, c *domain.CredencialesUsuario) (*domain.CredencialesUsuario, error) {
					return c, nil
				},
			},
			tenantRepo: &mockTenantRepo{
				crearFunc: func(_ context.Context, t *tenant.Tenant) (*tenant.Tenant, error) {
					return t, nil
				},
			},
			membresiaRepo: &mockMembresiaRepo{
				crearFunc: func(_ context.Context, _ *tenant.Membresia) error {
					return nil
				},
			},
			rolRepo: &mockRolRepo{
				obtenerPorNombreFunc: func(_ context.Context, nombre string) (*rbac.RolDB, error) {
					return &rbac.RolDB{ID: "rol-admin-id", Nombre: nombre}, nil
				},
			},
			usuarioTenantRolRepo: &mockUsuarioTenantRolRepo{
				crearFunc: func(_ context.Context, _, _, _ string) error {
					return nil
				},
			},
			rolPermisoRepo: &mockRolPermisoRepo{},
		},
		&mockEncriptacion{hash: "$2a$10$hashedpassword"},
		&mockGeneradorID{ids: []string{"user-id-1", "tenant-id-1"}},
		&mockRolPublisher{},
	)
}

func newUCWithRepos(
	usuarioRepo *mockUsuarioRepo,
	credencialesRepo *mockCredencialesRepo,
	encriptacion *mockEncriptacion,
	generadorID *mockGeneradorID,
	tenantRepo *mockTenantRepo,
	membresiaRepo *mockMembresiaRepo,
	rolRepo *mockRolRepo,
	usuarioTenantRolRepo *mockUsuarioTenantRolRepo,
	rolPermisoRepo *mockRolPermisoRepo,
	rolPublisher *mockRolPublisher,
) *register.RegistrarUsuarioCasoDeUso {
	return register.NewRegistrarUsuarioCasoDeUso(
		&mockUnitOfWork{
			usuarioRepo:          usuarioRepo,
			credencialesRepo:     credencialesRepo,
			tenantRepo:           tenantRepo,
			membresiaRepo:        membresiaRepo,
			rolRepo:              rolRepo,
			usuarioTenantRolRepo: usuarioTenantRolRepo,
			rolPermisoRepo:       rolPermisoRepo,
		},
		encriptacion,
		generadorID,
		rolPublisher,
	)
}

// ── Tests ─────────────────────────────────────────────────────────────────────

func TestRegistrarUsuarioExitoso(t *testing.T) {
	uc := newUCSuccess()

	cmd := &register.ComandoRegistrarUsuario{
		Correo:   "test@example.com",
		Password: "Secreto1!",
		Nombre:   "Juan",
		Apellido: "García",
		Telefono: "6001234567",
	}

	respuesta, err := uc.Ejecutar(context.Background(), cmd)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if respuesta.UsuarioID != "user-id-1" {
		t.Errorf("UsuarioID esperado 'user-id-1', obtenido '%s'", respuesta.UsuarioID)
	}
	if respuesta.TenantID != "tenant-id-1" {
		t.Errorf("TenantID esperado 'tenant-id-1', obtenido '%s'", respuesta.TenantID)
	}
	if respuesta.Correo != "test@example.com" {
		t.Errorf("Correo esperado 'test@example.com', obtenido '%s'", respuesta.Correo)
	}
	if respuesta.Estado != string(usuariodomain.NO_VERIFICADO) {
		t.Errorf("Estado esperado '%s', obtenido '%s'", usuariodomain.NO_VERIFICADO, respuesta.Estado)
	}
	if respuesta.CreadoEn.IsZero() {
		t.Error("CreadoEn no debería ser zero")
	}
}

func TestRegistrarUsuario_CreaTenantConNombreYSlug(t *testing.T) {
	var tenantCreado *tenant.Tenant

	uc := newUCWithRepos(
		&mockUsuarioRepo{
			crearFunc: func(_ context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
				return u, nil
			},
		},
		&mockCredencialesRepo{
			crearFunc: func(_ context.Context, c *domain.CredencialesUsuario) (*domain.CredencialesUsuario, error) {
				return c, nil
			},
		},
		&mockEncriptacion{hash: "$2a$10$hash"},
		&mockGeneradorID{ids: []string{"uid", "tid"}},
		&mockTenantRepo{
			crearFunc: func(_ context.Context, t *tenant.Tenant) (*tenant.Tenant, error) {
				tenantCreado = t
				return t, nil
			},
		},
		&mockMembresiaRepo{
			crearFunc: func(_ context.Context, _ *tenant.Membresia) error {
				return nil
			},
		},
		&mockRolRepo{
			obtenerPorNombreFunc: func(_ context.Context, nombre string) (*rbac.RolDB, error) {
				return &rbac.RolDB{ID: "rol-admin-id", Nombre: nombre}, nil
			},
		},
		&mockUsuarioTenantRolRepo{
			crearFunc: func(_ context.Context, _, _, _ string) error {
				return nil
			},
		},
		&mockRolPermisoRepo{},
		&mockRolPublisher{},
	)

	cmd := &register.ComandoRegistrarUsuario{
		Correo: "ana@example.com", Password: "Secreto1!",
		Nombre: "Ana", Apellido: "López",
	}

	_, err := uc.Ejecutar(context.Background(), cmd)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if tenantCreado == nil {
		t.Fatal("no se creó un tenant")
	}
	if tenantCreado.Nombre() != "Ana López" {
		t.Errorf("Nombre de tenant esperado 'Ana López', obtenido '%s'", tenantCreado.Nombre())
	}
	if tenantCreado.Slug() != "ana-lopez" {
		t.Errorf("Slug de tenant esperado 'ana-lopez', obtenido '%s'", tenantCreado.Slug())
	}
	if !tenantCreado.EstaActivo() {
		t.Error("Tenant debería estar activo")
	}
}

func TestRegistrarUsuario_CreaMembresia(t *testing.T) {
	var membresiaCreada *tenant.Membresia

	uc := newUCWithRepos(
		&mockUsuarioRepo{
			crearFunc: func(_ context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
				return u, nil
			},
		},
		&mockCredencialesRepo{
			crearFunc: func(_ context.Context, c *domain.CredencialesUsuario) (*domain.CredencialesUsuario, error) {
				return c, nil
			},
		},
		&mockEncriptacion{hash: "hash"},
		&mockGeneradorID{ids: []string{"uid", "tid"}},
		&mockTenantRepo{
			crearFunc: func(_ context.Context, t *tenant.Tenant) (*tenant.Tenant, error) {
				return t, nil
			},
		},
		&mockMembresiaRepo{
			crearFunc: func(_ context.Context, m *tenant.Membresia) error {
				membresiaCreada = m
				return nil
			},
		},
		&mockRolRepo{
			obtenerPorNombreFunc: func(_ context.Context, nombre string) (*rbac.RolDB, error) {
				return &rbac.RolDB{ID: "rol-admin-id", Nombre: nombre}, nil
			},
		},
		&mockUsuarioTenantRolRepo{
			crearFunc: func(_ context.Context, _, _, _ string) error {
				return nil
			},
		},
		&mockRolPermisoRepo{},
		&mockRolPublisher{},
	)

	_, err := uc.Ejecutar(context.Background(), &register.ComandoRegistrarUsuario{
		Correo: "test@example.com", Password: "Secreto1!",
		Nombre: "Test", Apellido: "User",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if membresiaCreada == nil {
		t.Fatal("no se creó una membresía")
	}
	if membresiaCreada.UsuarioID() != "uid" {
		t.Errorf("UsuarioID en membresía esperado 'uid', obtenido '%s'", membresiaCreada.UsuarioID())
	}
	if membresiaCreada.TenantID() != "tid" {
		t.Errorf("TenantID en membresía esperado 'tid', obtenido '%s'", membresiaCreada.TenantID())
	}
}

func TestRegistrarUsuarioCorreoVacio(t *testing.T) {
	uc := newUCSuccess()

	_, err := uc.Ejecutar(context.Background(), &register.ComandoRegistrarUsuario{
		Correo: "", Password: "pass", Nombre: "Juan",
	})
	if err == nil {
		t.Fatal("esperaba error por correo vacío")
	}
}

func TestRegistrarUsuarioPasswordVacio(t *testing.T) {
	uc := newUCSuccess()

	_, err := uc.Ejecutar(context.Background(), &register.ComandoRegistrarUsuario{
		Correo: "test@example.com", Password: "", Nombre: "Juan",
	})
	if err == nil {
		t.Fatal("esperaba error por password vacío")
	}
}

func TestRegistrarUsuarioNombreVacio(t *testing.T) {
	uc := newUCSuccess()

	_, err := uc.Ejecutar(context.Background(), &register.ComandoRegistrarUsuario{
		Correo: "test@example.com", Password: "Secreto1!", Nombre: "",
	})
	if err == nil {
		t.Fatal("esperaba error por nombre vacío")
	}
}

func TestRegistrarUsuarioEmailInvalido(t *testing.T) {
	uc := newUCSuccess()

	_, err := uc.Ejecutar(context.Background(), &register.ComandoRegistrarUsuario{
		Correo: "notanemail", Password: "Secreto1!", Nombre: "Test",
	})
	if err == nil {
		t.Fatal("esperaba error por email inválido")
	}
}

func TestRegistrarUsuarioCorreoDuplicado(t *testing.T) {
	errDuplicado := errors.New("el correo ya está registrado")
	uc := newUCWithRepos(
		&mockUsuarioRepo{
			crearFunc: func(_ context.Context, _ *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
				return nil, errDuplicado
			},
		},
		&mockCredencialesRepo{
			crearFunc: func(_ context.Context, c *domain.CredencialesUsuario) (*domain.CredencialesUsuario, error) {
				return c, nil
			},
		},
		&mockEncriptacion{hash: "hash"},
		&mockGeneradorID{ids: []string{"uid", "tid"}},
		&mockTenantRepo{
			crearFunc: func(_ context.Context, t *tenant.Tenant) (*tenant.Tenant, error) {
				return t, nil
			},
		},
		&mockMembresiaRepo{
			crearFunc: func(_ context.Context, _ *tenant.Membresia) error {
				return nil
			},
		},
		&mockRolRepo{
			obtenerPorNombreFunc: func(_ context.Context, nombre string) (*rbac.RolDB, error) {
				return &rbac.RolDB{ID: "rol-admin-id", Nombre: nombre}, nil
			},
		},
		&mockUsuarioTenantRolRepo{
			crearFunc: func(_ context.Context, _, _, _ string) error {
				return nil
			},
		},
		&mockRolPermisoRepo{},
		&mockRolPublisher{},
	)

	_, err := uc.Ejecutar(context.Background(), &register.ComandoRegistrarUsuario{
		Correo: "dup@example.com", Password: "Secreto1!",
		Nombre: "Dup", Apellido: "User",
	})
	if err == nil {
		t.Fatal("esperaba error por correo duplicado")
	}
	if !errors.Is(err, errDuplicado) {
		t.Errorf("esperaba error '%v', obtenido '%v'", errDuplicado, err)
	}
}

func TestRegistrarUsuarioErrorPersistiendoCredenciales(t *testing.T) {
	errCred := errors.New("error en BD de credenciales")
	uc := newUCWithRepos(
		&mockUsuarioRepo{
			crearFunc: func(_ context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
				return u, nil
			},
		},
		&mockCredencialesRepo{
			crearFunc: func(_ context.Context, _ *domain.CredencialesUsuario) (*domain.CredencialesUsuario, error) {
				return nil, errCred
			},
		},
		&mockEncriptacion{hash: "hash"},
		&mockGeneradorID{ids: []string{"uid", "tid"}},
		&mockTenantRepo{
			crearFunc: func(_ context.Context, t *tenant.Tenant) (*tenant.Tenant, error) {
				return t, nil
			},
		},
		&mockMembresiaRepo{
			crearFunc: func(_ context.Context, _ *tenant.Membresia) error {
				return nil
			},
		},
		&mockRolRepo{
			obtenerPorNombreFunc: func(_ context.Context, nombre string) (*rbac.RolDB, error) {
				return &rbac.RolDB{ID: "rol-admin-id", Nombre: nombre}, nil
			},
		},
		&mockUsuarioTenantRolRepo{
			crearFunc: func(_ context.Context, _, _, _ string) error {
				return nil
			},
		},
		&mockRolPermisoRepo{},
		&mockRolPublisher{},
	)

	_, err := uc.Ejecutar(context.Background(), &register.ComandoRegistrarUsuario{
		Correo: "test@example.com", Password: "Secreto1!",
		Nombre: "Test", Apellido: "User",
	})
	if err == nil {
		t.Fatal("esperaba error al persistir credenciales")
	}
}

func TestRegistrarUsuarioErrorPersistiendoTenant(t *testing.T) {
	errTenant := errors.New("slug duplicado")
	uc := newUCWithRepos(
		&mockUsuarioRepo{
			crearFunc: func(_ context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
				return u, nil
			},
		},
		&mockCredencialesRepo{
			crearFunc: func(_ context.Context, c *domain.CredencialesUsuario) (*domain.CredencialesUsuario, error) {
				return c, nil
			},
		},
		&mockEncriptacion{hash: "hash"},
		&mockGeneradorID{ids: []string{"uid", "tid"}},
		&mockTenantRepo{
			crearFunc: func(_ context.Context, _ *tenant.Tenant) (*tenant.Tenant, error) {
				return nil, errTenant
			},
		},
		&mockMembresiaRepo{
			crearFunc: func(_ context.Context, _ *tenant.Membresia) error {
				return nil
			},
		},
		&mockRolRepo{
			obtenerPorNombreFunc: func(_ context.Context, nombre string) (*rbac.RolDB, error) {
				return &rbac.RolDB{ID: "rol-admin-id", Nombre: nombre}, nil
			},
		},
		&mockUsuarioTenantRolRepo{
			crearFunc: func(_ context.Context, _, _, _ string) error {
				return nil
			},
		},
		&mockRolPermisoRepo{},
		&mockRolPublisher{},
	)

	_, err := uc.Ejecutar(context.Background(), &register.ComandoRegistrarUsuario{
		Correo: "test@example.com", Password: "Secreto1!",
		Nombre: "Test", Apellido: "User",
	})
	if err == nil {
		t.Fatal("esperaba error al persistir tenant")
	}
}

func TestRegistrarUsuarioErrorPersistiendoMembresia(t *testing.T) {
	errMemb := errors.New("error en membresía")
	uc := newUCWithRepos(
		&mockUsuarioRepo{
			crearFunc: func(_ context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
				return u, nil
			},
		},
		&mockCredencialesRepo{
			crearFunc: func(_ context.Context, c *domain.CredencialesUsuario) (*domain.CredencialesUsuario, error) {
				return c, nil
			},
		},
		&mockEncriptacion{hash: "hash"},
		&mockGeneradorID{ids: []string{"uid", "tid"}},
		&mockTenantRepo{
			crearFunc: func(_ context.Context, t *tenant.Tenant) (*tenant.Tenant, error) {
				return t, nil
			},
		},
		&mockMembresiaRepo{
			crearFunc: func(_ context.Context, _ *tenant.Membresia) error {
				return errMemb
			},
		},
		&mockRolRepo{
			obtenerPorNombreFunc: func(_ context.Context, nombre string) (*rbac.RolDB, error) {
				return &rbac.RolDB{ID: "rol-admin-id", Nombre: nombre}, nil
			},
		},
		&mockUsuarioTenantRolRepo{
			crearFunc: func(_ context.Context, _, _, _ string) error {
				return nil
			},
		},
		&mockRolPermisoRepo{},
		&mockRolPublisher{},
	)

	_, err := uc.Ejecutar(context.Background(), &register.ComandoRegistrarUsuario{
		Correo: "test@example.com", Password: "Secreto1!",
		Nombre: "Test", Apellido: "User",
	})
	if err == nil {
		t.Fatal("esperaba error al persistir membresía")
	}
}

func TestRegistrarUsuarioIDsInyectados(t *testing.T) {
	uc := newUCWithRepos(
		&mockUsuarioRepo{
			crearFunc: func(_ context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
				return u, nil
			},
		},
		&mockCredencialesRepo{
			crearFunc: func(_ context.Context, c *domain.CredencialesUsuario) (*domain.CredencialesUsuario, error) {
				return c, nil
			},
		},
		&mockEncriptacion{hash: "hash"},
		&mockGeneradorID{ids: []string{"custom-user-id", "custom-tenant-id"}},
		&mockTenantRepo{
			crearFunc: func(_ context.Context, t *tenant.Tenant) (*tenant.Tenant, error) {
				return t, nil
			},
		},
		&mockMembresiaRepo{
			crearFunc: func(_ context.Context, _ *tenant.Membresia) error {
				return nil
			},
		},
		&mockRolRepo{
			obtenerPorNombreFunc: func(_ context.Context, nombre string) (*rbac.RolDB, error) {
				return &rbac.RolDB{ID: "rol-admin-id", Nombre: nombre}, nil
			},
		},
		&mockUsuarioTenantRolRepo{
			crearFunc: func(_ context.Context, _, _, _ string) error {
				return nil
			},
		},
		&mockRolPermisoRepo{},
		&mockRolPublisher{},
	)

	resp, err := uc.Ejecutar(context.Background(), &register.ComandoRegistrarUsuario{
		Correo: "test@example.com", Password: "Secreto1!",
		Nombre: "Test", Apellido: "User",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.UsuarioID != "custom-user-id" {
		t.Errorf("UsuarioID esperado 'custom-user-id', obtenido '%s'", resp.UsuarioID)
	}
	if resp.TenantID != "custom-tenant-id" {
		t.Errorf("TenantID esperado 'custom-tenant-id', obtenido '%s'", resp.TenantID)
	}
}

func TestRegistrarUsuarioPasswordHasheado(t *testing.T) {
	var hashRecibido string
	uc := newUCWithRepos(
		&mockUsuarioRepo{
			crearFunc: func(_ context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
				return u, nil
			},
		},
		&mockCredencialesRepo{
			crearFunc: func(_ context.Context, c *domain.CredencialesUsuario) (*domain.CredencialesUsuario, error) {
				hashRecibido = c.PasswordHash()
				return c, nil
			},
		},
		&mockEncriptacion{hash: "bcrypt_hash_123"},
		&mockGeneradorID{ids: []string{"uid", "tid"}},
		&mockTenantRepo{
			crearFunc: func(_ context.Context, t *tenant.Tenant) (*tenant.Tenant, error) {
				return t, nil
			},
		},
		&mockMembresiaRepo{
			crearFunc: func(_ context.Context, _ *tenant.Membresia) error {
				return nil
			},
		},
		&mockRolRepo{
			obtenerPorNombreFunc: func(_ context.Context, nombre string) (*rbac.RolDB, error) {
				return &rbac.RolDB{ID: "rol-admin-id", Nombre: nombre}, nil
			},
		},
		&mockUsuarioTenantRolRepo{
			crearFunc: func(_ context.Context, _, _, _ string) error {
				return nil
			},
		},
		&mockRolPermisoRepo{},
		&mockRolPublisher{},
	)

	_, err := uc.Ejecutar(context.Background(), &register.ComandoRegistrarUsuario{
		Correo: "test@example.com", Password: "Secreto1!",
		Nombre: "Test", Apellido: "User",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if hashRecibido != "bcrypt_hash_123" {
		t.Errorf("hash esperado 'bcrypt_hash_123', obtenido '%s'", hashRecibido)
	}
}

func TestRegistrarUsuarioSlugGenerado(t *testing.T) {
	tests := []struct {
		nombre   string
		apellido string
		slug     string
	}{
		{"Ana", "López", "ana-lopez"},
		{"José", "Martínez", "jose-martinez"},
		{"Carlos", "Ramos", "carlos-ramos"},
		{"  Spaces  ", "  Here  ", "spaces-here"},
		{"UpperCase", "User", "uppercase-user"},
		{"Special!@#", "Name$%^", "special-name"},
	}

	for _, tt := range tests {
		t.Run(tt.nombre+"_"+tt.apellido, func(t *testing.T) {
			var tenantCreado *tenant.Tenant
			uc := newUCWithRepos(
				&mockUsuarioRepo{
					crearFunc: func(_ context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
						return u, nil
					},
				},
				&mockCredencialesRepo{
					crearFunc: func(_ context.Context, c *domain.CredencialesUsuario) (*domain.CredencialesUsuario, error) {
						return c, nil
					},
				},
				&mockEncriptacion{hash: "hash"},
				&mockGeneradorID{ids: []string{"uid", "tid"}},
				&mockTenantRepo{
					crearFunc: func(_ context.Context, t *tenant.Tenant) (*tenant.Tenant, error) {
						tenantCreado = t
						return t, nil
					},
				},
				&mockMembresiaRepo{
					crearFunc: func(_ context.Context, _ *tenant.Membresia) error {
						return nil
					},
				},
				&mockRolRepo{
					obtenerPorNombreFunc: func(_ context.Context, nombre string) (*rbac.RolDB, error) {
						return &rbac.RolDB{ID: "rol-admin-id", Nombre: nombre}, nil
					},
				},
				&mockUsuarioTenantRolRepo{
					crearFunc: func(_ context.Context, _, _, _ string) error {
						return nil
					},
				},
				&mockRolPermisoRepo{},
				&mockRolPublisher{},
			)

			_, err := uc.Ejecutar(context.Background(), &register.ComandoRegistrarUsuario{
				Correo: "test@example.com", Password: "Secreto1!",
				Nombre: tt.nombre, Apellido: tt.apellido,
			})
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}
			if tenantCreado.Slug() != tt.slug {
				t.Errorf("slug esperado '%s', obtenido '%s'", tt.slug, tenantCreado.Slug())
			}
		})
	}
}

func TestRegistrarUsuarioCamposCorreoYTelefonoOpcional(t *testing.T) {
	uc := newUCSuccess()

	resp, err := uc.Ejecutar(context.Background(), &register.ComandoRegistrarUsuario{
		Correo: "test@domain.com", Password: "Secreto1!",
		Nombre: "Test", Apellido: "User", Telefono: "0987654321",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.Correo != "test@domain.com" {
		t.Errorf("Correo esperado 'test@domain.com', obtenido '%s'", resp.Correo)
	}
}

func TestRegistrarUsuario_AsignaRolAdministrador(t *testing.T) {
	var (
		usuarioIDAsignado string
		tenantIDAsignado  string
		rolIDAsignado     string
	)

	uc := newUCWithRepos(
		&mockUsuarioRepo{
			crearFunc: func(_ context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
				return u, nil
			},
		},
		&mockCredencialesRepo{
			crearFunc: func(_ context.Context, c *domain.CredencialesUsuario) (*domain.CredencialesUsuario, error) {
				return c, nil
			},
		},
		&mockEncriptacion{hash: "hash"},
		&mockGeneradorID{ids: []string{"uid", "tid"}},
		&mockTenantRepo{
			crearFunc: func(_ context.Context, t *tenant.Tenant) (*tenant.Tenant, error) {
				return t, nil
			},
		},
		&mockMembresiaRepo{
			crearFunc: func(_ context.Context, _ *tenant.Membresia) error {
				return nil
			},
		},
		&mockRolRepo{
			obtenerPorNombreFunc: func(_ context.Context, nombre string) (*rbac.RolDB, error) {
				return &rbac.RolDB{ID: "rol-admin-id", Nombre: nombre}, nil
			},
		},
		&mockUsuarioTenantRolRepo{
			crearFunc: func(_ context.Context, usuarioID, tenantID, rolID string) error {
				usuarioIDAsignado = usuarioID
				tenantIDAsignado = tenantID
				rolIDAsignado = rolID
				return nil
			},
		},
		&mockRolPermisoRepo{},
		&mockRolPublisher{},
	)

	_, err := uc.Ejecutar(context.Background(), &register.ComandoRegistrarUsuario{
		Correo: "test@example.com", Password: "Secreto1!",
		Nombre: "Test", Apellido: "User",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if usuarioIDAsignado != "uid" {
		t.Errorf("usuarioID en rol esperado 'uid', obtenido '%s'", usuarioIDAsignado)
	}
	if tenantIDAsignado != "tid" {
		t.Errorf("tenantID en rol esperado 'tid', obtenido '%s'", tenantIDAsignado)
	}
	if rolIDAsignado != "rol-admin-id" {
		t.Errorf("rolID esperado 'rol-admin-id', obtenido '%s'", rolIDAsignado)
	}
}
