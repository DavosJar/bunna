package createuser_test

import (
	"context"
	"errors"
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/rbac/domain"
	seguridad "github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
	shareddomain "github.com/davosjar/bunna/services/identidad/internal/shared/domain"
	"github.com/davosjar/bunna/services/identidad/internal/usuarios/application/usecases/createuser"
	usuariodomain "github.com/davosjar/bunna/services/identidad/internal/usuarios/domain/usuario"
)

type mockUsuarioRepoCreate struct {
	crearFunc      func(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error)
	obtenerPorID   func(ctx context.Context, id string) (*usuariodomain.Usuario, error)
}

func (m *mockUsuarioRepoCreate) Crear(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
	return m.crearFunc(ctx, u)
}
func (m *mockUsuarioRepoCreate) Actualizar(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
	return u, nil
}
func (m *mockUsuarioRepoCreate) Eliminar(ctx context.Context, id string) error { return nil }
func (m *mockUsuarioRepoCreate) ObtenerPorID(ctx context.Context, id string) (*usuariodomain.Usuario, error) {
	if m.obtenerPorID != nil {
		return m.obtenerPorID(ctx, id)
	}
	return nil, nil
}
func (m *mockUsuarioRepoCreate) Listar(ctx context.Context, _ usuariodomain.EspecificacionUsuario, _ shareddomain.Paginacion) ([]*usuariodomain.Usuario, error) {
	return nil, nil
}

type mockCredencialesRepoCreate struct {
	crearFunc func(ctx context.Context, c *seguridad.CredencialesUsuario) (*seguridad.CredencialesUsuario, error)
}

func (m *mockCredencialesRepoCreate) Crear(ctx context.Context, c *seguridad.CredencialesUsuario) (*seguridad.CredencialesUsuario, error) {
	return m.crearFunc(ctx, c)
}
func (m *mockCredencialesRepoCreate) Actualizar(ctx context.Context, c *seguridad.CredencialesUsuario) (*seguridad.CredencialesUsuario, error) {
	return c, nil
}
func (m *mockCredencialesRepoCreate) ObtenerPorUsuarioID(ctx context.Context, usuarioID string) (*seguridad.CredencialesUsuario, error) {
	return nil, nil
}
func (m *mockCredencialesRepoCreate) Eliminar(ctx context.Context, usuarioID string) error { return nil }
func (m *mockCredencialesRepoCreate) Find(ctx context.Context, _ seguridad.EspecificacionCredenciales, _ shareddomain.Paginacion) ([]*seguridad.CredencialesUsuario, error) {
	return nil, nil
}

type mockEncriptacionCreate struct{}

func (m *mockEncriptacionCreate) Hashear(password string) (string, error) {
	return "hashed:" + password, nil
}
func (m *mockEncriptacionCreate) Verificar(password, hash string) bool {
	return "hashed:"+password == hash
}

type mockAuthSvcCreate struct {
	ok  bool
	err error
}

func (m *mockAuthSvcCreate) TienePermiso(ctx context.Context, usuarioID, tenantID, codigoPermiso string) (bool, error) {
	return m.ok, m.err
}

type mockGeneradorIDCreate struct {
	id  string
	err error
}

func (m *mockGeneradorIDCreate) NextID(ctx context.Context) (string, error) {
	return m.id, m.err
}

func TestCrearUsuarioExitoso(t *testing.T) {
	userRepo := &mockUsuarioRepoCreate{
		crearFunc: func(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
			return u, nil
		},
	}
	credRepo := &mockCredencialesRepoCreate{
		crearFunc: func(ctx context.Context, c *seguridad.CredencialesUsuario) (*seguridad.CredencialesUsuario, error) {
			return c, nil
		},
	}
	uc := createuser.NewCrearUsuarioCasoDeUso(
		userRepo, credRepo, &mockEncriptacionCreate{},
		&mockAuthSvcCreate{ok: true}, &mockGeneradorIDCreate{id: "new-id"},
	)
	resp, err := uc.Ejecutar(context.Background(), &createuser.ComandoCrearUsuario{
		Correo: "test@gmail.com", Nombre: "Juan", Apellido: "Perez",
		Password: "Abcdef1!", EjecutorID: "admin-id",
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if resp.ID != "new-id" {
		t.Errorf("ID incorrecto: %s", resp.ID)
	}
	if resp.Correo != "test@gmail.com" {
		t.Errorf("Correo incorrecto: %s", resp.Correo)
	}
	if resp.Nombre != "Juan" {
		t.Errorf("Nombre incorrecto: %s", resp.Nombre)
	}
	if !resp.Activo {
		t.Error("Activo debe ser true")
	}
	if resp.CreadoEn == "" {
		t.Error("CreadoEn no debe estar vacío")
	}
}

func TestCrearUsuarioPermisoDenegado(t *testing.T) {
	uc := createuser.NewCrearUsuarioCasoDeUso(
		&mockUsuarioRepoCreate{}, &mockCredencialesRepoCreate{},
		&mockEncriptacionCreate{}, &mockAuthSvcCreate{ok: false},
		&mockGeneradorIDCreate{},
	)
	_, err := uc.Ejecutar(context.Background(), &createuser.ComandoCrearUsuario{
		Correo: "test@gmail.com", Nombre: "Juan", Password: "Abcdef1!",
	})
	if !errors.Is(err, rbac.ErrPermisoDenegado) {
		t.Errorf("esperaba ErrPermisoDenegado, got %v", err)
	}
}

func TestCrearUsuarioAuthError(t *testing.T) {
	uc := createuser.NewCrearUsuarioCasoDeUso(
		&mockUsuarioRepoCreate{}, &mockCredencialesRepoCreate{},
		&mockEncriptacionCreate{}, &mockAuthSvcCreate{err: errors.New("fallo")},
		&mockGeneradorIDCreate{},
	)
	_, err := uc.Ejecutar(context.Background(), &createuser.ComandoCrearUsuario{
		Correo: "test@gmail.com", Nombre: "Juan", Password: "Abcdef1!",
	})
	if err == nil {
		t.Fatal("esperaba error")
	}
}

func TestCrearUsuarioCorreoVacio(t *testing.T) {
	uc := createuser.NewCrearUsuarioCasoDeUso(
		&mockUsuarioRepoCreate{}, &mockCredencialesRepoCreate{},
		&mockEncriptacionCreate{}, &mockAuthSvcCreate{ok: true},
		&mockGeneradorIDCreate{},
	)
	_, err := uc.Ejecutar(context.Background(), &createuser.ComandoCrearUsuario{
		Correo: "", Nombre: "Juan", Password: "Abcdef1!",
	})
	if err == nil || err.Error() != "correo no puede estar vacío" {
		t.Errorf("esperaba error de correo vacío, got %v", err)
	}
}

func TestCrearUsuarioCorreoInvalido(t *testing.T) {
	uc := createuser.NewCrearUsuarioCasoDeUso(
		&mockUsuarioRepoCreate{}, &mockCredencialesRepoCreate{},
		&mockEncriptacionCreate{}, &mockAuthSvcCreate{ok: true},
		&mockGeneradorIDCreate{},
	)
	_, err := uc.Ejecutar(context.Background(), &createuser.ComandoCrearUsuario{
		Correo: "invalido", Nombre: "Juan", Password: "Abcdef1!",
	})
	if err == nil {
		t.Fatal("esperaba error de formato de correo")
	}
}

func TestCrearUsuarioNombreVacio(t *testing.T) {
	uc := createuser.NewCrearUsuarioCasoDeUso(
		&mockUsuarioRepoCreate{}, &mockCredencialesRepoCreate{},
		&mockEncriptacionCreate{}, &mockAuthSvcCreate{ok: true},
		&mockGeneradorIDCreate{},
	)
	_, err := uc.Ejecutar(context.Background(), &createuser.ComandoCrearUsuario{
		Correo: "test@gmail.com", Nombre: "", Password: "Abcdef1!",
	})
	if err == nil || err.Error() != "nombre no puede estar vacío" {
		t.Errorf("esperaba error de nombre vacío, got %v", err)
	}
}

func TestCrearUsuarioPasswordVacio(t *testing.T) {
	uc := createuser.NewCrearUsuarioCasoDeUso(
		&mockUsuarioRepoCreate{}, &mockCredencialesRepoCreate{},
		&mockEncriptacionCreate{}, &mockAuthSvcCreate{ok: true},
		&mockGeneradorIDCreate{},
	)
	_, err := uc.Ejecutar(context.Background(), &createuser.ComandoCrearUsuario{
		Correo: "test@gmail.com", Nombre: "Juan", Password: "",
	})
	if err == nil || err.Error() != "password no puede estar vacío" {
		t.Errorf("esperaba error de password vacío, got %v", err)
	}
}

func TestCrearUsuarioPasswordInvalido(t *testing.T) {
	uc := createuser.NewCrearUsuarioCasoDeUso(
		&mockUsuarioRepoCreate{}, &mockCredencialesRepoCreate{},
		&mockEncriptacionCreate{}, &mockAuthSvcCreate{ok: true},
		&mockGeneradorIDCreate{},
	)
	_, err := uc.Ejecutar(context.Background(), &createuser.ComandoCrearUsuario{
		Correo: "test@gmail.com", Nombre: "Juan", Password: "abc",
	})
	if err == nil {
		t.Fatal("esperaba error de formato de password")
	}
}

func TestCrearUsuarioIDGenError(t *testing.T) {
	uc := createuser.NewCrearUsuarioCasoDeUso(
		&mockUsuarioRepoCreate{}, &mockCredencialesRepoCreate{},
		&mockEncriptacionCreate{}, &mockAuthSvcCreate{ok: true},
		&mockGeneradorIDCreate{err: errors.New("gen error")},
	)
	_, err := uc.Ejecutar(context.Background(), &createuser.ComandoCrearUsuario{
		Correo: "test@gmail.com", Nombre: "Juan", Password: "Abcdef1!",
	})
	if err == nil {
		t.Fatal("esperaba error al generar ID")
	}
}

func TestCrearUsuarioRepoError(t *testing.T) {
	userRepo := &mockUsuarioRepoCreate{
		crearFunc: func(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
			return nil, errors.New("db error")
		},
	}
	uc := createuser.NewCrearUsuarioCasoDeUso(
		userRepo, &mockCredencialesRepoCreate{},
		&mockEncriptacionCreate{}, &mockAuthSvcCreate{ok: true},
		&mockGeneradorIDCreate{id: "new-id"},
	)
	_, err := uc.Ejecutar(context.Background(), &createuser.ComandoCrearUsuario{
		Correo: "test@gmail.com", Nombre: "Juan", Password: "Abcdef1!",
	})
	if err == nil {
		t.Fatal("esperaba error de repositorio")
	}
}

func TestCrearUsuarioHashError(t *testing.T) {
	userRepo := &mockUsuarioRepoCreate{
		crearFunc: func(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
			return u, nil
		},
	}
	encSvc := &mockEncriptacionFail{}
	uc := createuser.NewCrearUsuarioCasoDeUso(
		userRepo, &mockCredencialesRepoCreate{},
		encSvc, &mockAuthSvcCreate{ok: true},
		&mockGeneradorIDCreate{id: "new-id"},
	)
	_, err := uc.Ejecutar(context.Background(), &createuser.ComandoCrearUsuario{
		Correo: "test@gmail.com", Nombre: "Juan", Password: "Abcdef1!",
	})
	if err == nil {
		t.Fatal("esperaba error al hashear")
	}
}

type mockEncriptacionFail struct{}

func (m *mockEncriptacionFail) Hashear(password string) (string, error) {
	return "", errors.New("hash error")
}
func (m *mockEncriptacionFail) Verificar(password, hash string) bool { return false }

func TestCrearUsuarioCredRepoError(t *testing.T) {
	userRepo := &mockUsuarioRepoCreate{
		crearFunc: func(ctx context.Context, u *usuariodomain.Usuario) (*usuariodomain.Usuario, error) {
			return u, nil
		},
	}
	credRepo := &mockCredencialesRepoCreate{
		crearFunc: func(ctx context.Context, c *seguridad.CredencialesUsuario) (*seguridad.CredencialesUsuario, error) {
			return nil, errors.New("db error")
		},
	}
	uc := createuser.NewCrearUsuarioCasoDeUso(
		userRepo, credRepo, &mockEncriptacionCreate{},
		&mockAuthSvcCreate{ok: true}, &mockGeneradorIDCreate{id: "new-id"},
	)
	_, err := uc.Ejecutar(context.Background(), &createuser.ComandoCrearUsuario{
		Correo: "test@gmail.com", Nombre: "Juan", Password: "Abcdef1!",
	})
	if err == nil {
		t.Fatal("esperaba error de repositorio de credenciales")
	}
}
