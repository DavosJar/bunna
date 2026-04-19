package usuario

import (
	"context"
	"errors"
	"testing"
)

type mockRepositorio struct {
	usuarios []*Usuario
}

func (m *mockRepositorio) Crear(ctx context.Context, usuario *Usuario) (*Usuario, error) {
	m.usuarios = append(m.usuarios, usuario)
	return usuario, nil
}

func (m *mockRepositorio) Actualizar(ctx context.Context, usuario *Usuario) (*Usuario, error) {
	for i, u := range m.usuarios {
		if u.ID() == usuario.ID() {
			m.usuarios[i] = usuario
			return usuario, nil
		}
	}
	return nil, errors.New("usuario no encontrado")
}

func (m *mockRepositorio) Eliminar(ctx context.Context, id string) error {
	for i, u := range m.usuarios {
		if u.ID() == id {
			m.usuarios = append(m.usuarios[:i], m.usuarios[i+1:]...)
			return nil
		}
	}
	return errors.New("usuario no encontrado")
}

func (m *mockRepositorio) ObtenerPorID(ctx context.Context, id string) (*Usuario, error) {
	for _, u := range m.usuarios {
		if u.ID() == id {
			return u, nil
		}
	}
	return nil, errors.New("usuario no encontrado")
}

func (m *mockRepositorio) Listar(ctx context.Context, especificacion EspecificacionUsuario, paginacion Paginacion) ([]*Usuario, error) {
	resultado := m.filtrar(especificacion)
	resultado = m.ordenar(resultado, paginacion.Ordenaciones)
	return m.paginar(resultado, paginacion), nil
}

func (m *mockRepositorio) filtrar(especificacion EspecificacionUsuario) []*Usuario {
	if len(especificacion.ListaLiltros) == 0 {
		return m.usuarios
	}

	var resultado []*Usuario
	for _, u := range m.usuarios {
		if cumpleFiltros(u, especificacion.ListaLiltros) {
			resultado = append(resultado, u)
		}
	}
	return resultado
}

func cumpleFiltros(u *Usuario, filtros []CriterioFiltro) bool {
	for _, f := range filtros {
		switch f.Campo {
		case "nombre":
			if !compararString(u.Nombre(), f.Operador, f.Valor) {
				return false
			}
		case "apellido":
			if !compararString(u.Apellido(), f.Operador, f.Valor) {
				return false
			}
		case "correo":
			if !compararString(u.Correo(), f.Operador, f.Valor) {
				return false
			}
		case "telefono":
			if !compararString(u.Telefono(), f.Operador, f.Valor) {
				return false
			}
		case "estado":
			if !compararString(string(u.Estado()), f.Operador, f.Valor) {
				return false
			}
		}
	}
	return true
}

func compararString(valor string, operador string, esperado any) bool {
	esperadoStr, ok := esperado.(string)
	if !ok {
		return false
	}
	switch operador {
	case "=":
		return valor == esperadoStr
	case "!=":
		return valor != esperadoStr
	case "LIKE":
		return contains(valor, esperadoStr)
	}
	return false
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func ptrString(s string) *string {
	return &s
}

func (m *mockRepositorio) ordenar(usuarios []*Usuario, ordenaciones []Ordenacion) []*Usuario {
	if len(ordenaciones) == 0 {
		return usuarios
	}

	resultado := make([]*Usuario, len(usuarios))
	copy(resultado, usuarios)

	for _, ord := range ordenaciones {
		switch ord.Campo {
		case "nombre":
			sortByString(resultado, func(u *Usuario) string { return u.Nombre() }, ord.Tipo == DESC)
		case "apellido":
			sortByString(resultado, func(u *Usuario) string { return u.Apellido() }, ord.Tipo == DESC)
		case "correo":
			sortByString(resultado, func(u *Usuario) string { return u.Correo() }, ord.Tipo == DESC)
		}
	}
	return resultado
}

func sortByString(usuarios []*Usuario, getter func(*Usuario) string, desc bool) {
	for i := 0; i < len(usuarios)-1; i++ {
		for j := i + 1; j < len(usuarios); j++ {
			if desc {
				if getter(usuarios[i]) < getter(usuarios[j]) {
					usuarios[i], usuarios[j] = usuarios[j], usuarios[i]
				}
			} else {
				if getter(usuarios[i]) > getter(usuarios[j]) {
					usuarios[i], usuarios[j] = usuarios[j], usuarios[i]
				}
			}
		}
	}
}

func (m *mockRepositorio) paginar(usuarios []*Usuario, paginacion Paginacion) []*Usuario {
	if paginacion.TamanoPagina <= 0 {
		return usuarios
	}

	inicio := (paginacion.Pagina - 1) * paginacion.TamanoPagina
	if inicio >= len(usuarios) {
		return []*Usuario{}
	}

	fin := inicio + paginacion.TamanoPagina
	if fin > len(usuarios) {
		fin = len(usuarios)
	}

	return usuarios[inicio:fin]
}

func crearUsuariosPrueba() []*Usuario {
	u1, _ := NuevoUsuario("", "Ana", "García", "ana@test.com", "+34111111111")
	u2, _ := NuevoUsuario("", "Carlos", "López", "carlos@test.com", "+34222222222")
	u3, _ := NuevoUsuario("", "Beatriz", "Martínez", "beatriz@test.com", "+34333333333")
	u4, _ := NuevoUsuario("", "David", "García", "david@test.com", "+34444444444")
	u5, _ := NuevoUsuario("", "Elena", "Sánchez", "elena@test.com", "+34555555555")
	u1.CambiarEstado(ACTIVO)
	u2.CambiarEstado(ACTIVO)
	u3.CambiarEstado(INACTIVO)
	u4.CambiarEstado(ACTIVO)
	u5.CambiarEstado(BLOQUEADO)
	return []*Usuario{u1, u2, u3, u4, u5}
}

func TestListarSinFiltros(t *testing.T) {
	mock := &mockRepositorio{usuarios: crearUsuariosPrueba()}
	repo := UsuarioRepositorio(mock)

	resultado, err := repo.Listar(context.Background(), EspecificacionUsuario{}, Paginacion{Pagina: 1, TamanoPagina: 10})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(resultado) != 5 {
		t.Errorf("Expected 5 usuarios, got %d", len(resultado))
	}
}

func TestListarConFiltroIgualdad(t *testing.T) {
	mock := &mockRepositorio{usuarios: crearUsuariosPrueba()}
	repo := UsuarioRepositorio(mock)

	especificacion := EspecificacionUsuario{
		ListaLiltros: []CriterioFiltro{
			{Campo: "nombre", Operador: "=", Valor: "Ana"},
		},
	}

	resultado, err := repo.Listar(context.Background(), especificacion, Paginacion{Pagina: 1, TamanoPagina: 10})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(resultado) != 1 {
		t.Errorf("Expected 1 usuario, got %d", len(resultado))
	}
	if len(resultado) > 0 && resultado[0].Nombre() != "Ana" {
		t.Errorf("Expected nombre 'Ana', got '%s'", resultado[0].Nombre())
	}
}

func TestListarConFiltroEstado(t *testing.T) {
	mock := &mockRepositorio{usuarios: crearUsuariosPrueba()}
	repo := UsuarioRepositorio(mock)

	especificacion := EspecificacionUsuario{
		ListaLiltros: []CriterioFiltro{
			{Campo: "estado", Operador: "=", Valor: "ACTIVO"},
		},
	}

	resultado, err := repo.Listar(context.Background(), especificacion, Paginacion{Pagina: 1, TamanoPagina: 10})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(resultado) != 3 {
		t.Errorf("Expected 3 usuarios activos, got %d", len(resultado))
	}
}

func TestListarConFiltroLike(t *testing.T) {
	mock := &mockRepositorio{usuarios: crearUsuariosPrueba()}
	repo := UsuarioRepositorio(mock)

	especificacion := EspecificacionUsuario{
		ListaLiltros: []CriterioFiltro{
			{Campo: "correo", Operador: "LIKE", Valor: "test"},
		},
	}

	resultado, err := repo.Listar(context.Background(), especificacion, Paginacion{Pagina: 1, TamanoPagina: 10})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(resultado) != 5 {
		t.Errorf("Expected 5 usuarios con 'test' en correo, got %d", len(resultado))
	}
}

func TestListarConFiltroApellido(t *testing.T) {
	mock := &mockRepositorio{usuarios: crearUsuariosPrueba()}
	repo := UsuarioRepositorio(mock)

	especificacion := EspecificacionUsuario{
		ListaLiltros: []CriterioFiltro{
			{Campo: "apellido", Operador: "=", Valor: "García"},
		},
	}

	resultado, err := repo.Listar(context.Background(), especificacion, Paginacion{Pagina: 1, TamanoPagina: 10})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(resultado) != 2 {
		t.Errorf("Expected 2 usuarios con apellido García, got %d", len(resultado))
	}
}

func TestListarConMultipleFiltros(t *testing.T) {
	mock := &mockRepositorio{usuarios: crearUsuariosPrueba()}
	repo := UsuarioRepositorio(mock)

	especificacion := EspecificacionUsuario{
		ListaLiltros: []CriterioFiltro{
			{Campo: "estado", Operador: "=", Valor: "ACTIVO"},
			{Campo: "apellido", Operador: "=", Valor: "García"},
		},
	}

	resultado, err := repo.Listar(context.Background(), especificacion, Paginacion{Pagina: 1, TamanoPagina: 10})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(resultado) != 1 {
		t.Errorf("Expected 1 usuario activo con apellido García, got %d", len(resultado))
	}
}

func TestListarConPaginacionPrimeraPagina(t *testing.T) {
	mock := &mockRepositorio{usuarios: crearUsuariosPrueba()}
	repo := UsuarioRepositorio(mock)

	resultado, err := repo.Listar(context.Background(), EspecificacionUsuario{}, Paginacion{Pagina: 1, TamanoPagina: 2})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(resultado) != 2 {
		t.Errorf("Expected 2 usuarios en página 1, got %d", len(resultado))
	}
}

func TestListarConPaginacionSegundaPagina(t *testing.T) {
	mock := &mockRepositorio{usuarios: crearUsuariosPrueba()}
	repo := UsuarioRepositorio(mock)

	resultado, err := repo.Listar(context.Background(), EspecificacionUsuario{}, Paginacion{Pagina: 2, TamanoPagina: 2})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(resultado) != 2 {
		t.Errorf("Expected 2 usuarios en página 2, got %d", len(resultado))
	}
}

func TestListarConPaginacionUltimaPaginaParcial(t *testing.T) {
	mock := &mockRepositorio{usuarios: crearUsuariosPrueba()}
	repo := UsuarioRepositorio(mock)

	resultado, err := repo.Listar(context.Background(), EspecificacionUsuario{}, Paginacion{Pagina: 3, TamanoPagina: 2})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(resultado) != 1 {
		t.Errorf("Expected 1 usuario en página 3, got %d", len(resultado))
	}
}

func TestListarConPaginacionPaginaVacia(t *testing.T) {
	mock := &mockRepositorio{usuarios: crearUsuariosPrueba()}
	repo := UsuarioRepositorio(mock)

	resultado, err := repo.Listar(context.Background(), EspecificacionUsuario{}, Paginacion{Pagina: 10, TamanoPagina: 2})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(resultado) != 0 {
		t.Errorf("Expected 0 usuarios en página 10, got %d", len(resultado))
	}
}

func TestListarConOrdenacionASC(t *testing.T) {
	mock := &mockRepositorio{usuarios: crearUsuariosPrueba()}
	repo := UsuarioRepositorio(mock)

	paginacion := Paginacion{
		Pagina:       1,
		TamanoPagina: 10,
		Ordenaciones: []Ordenacion{{Campo: "nombre", Tipo: ASC}},
	}

	resultado, err := repo.Listar(context.Background(), EspecificacionUsuario{}, paginacion)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(resultado) < 2 {
		t.Fatalf("Expected at least 2 usuarios, got %d", len(resultado))
	}

	nombres := []string{"Ana", "Beatriz", "Carlos", "David", "Elena"}
	for i, esperado := range nombres {
		if resultado[i].Nombre() != esperado {
			t.Errorf("Position %d: expected '%s', got '%s'", i, esperado, resultado[i].Nombre())
		}
	}
}

func TestListarConOrdenacionDESC(t *testing.T) {
	mock := &mockRepositorio{usuarios: crearUsuariosPrueba()}
	repo := UsuarioRepositorio(mock)

	paginacion := Paginacion{
		Pagina:       1,
		TamanoPagina: 10,
		Ordenaciones: []Ordenacion{{Campo: "nombre", Tipo: DESC}},
	}

	resultado, err := repo.Listar(context.Background(), EspecificacionUsuario{}, paginacion)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(resultado) < 2 {
		t.Fatalf("Expected at least 2 usuarios, got %d", len(resultado))
	}

	nombres := []string{"Elena", "David", "Carlos", "Beatriz", "Ana"}
	for i, esperado := range nombres {
		if resultado[i].Nombre() != esperado {
			t.Errorf("Position %d: expected '%s', got '%s'", i, esperado, resultado[i].Nombre())
		}
	}
}

func TestListarConFiltroYPaginacion(t *testing.T) {
	mock := &mockRepositorio{usuarios: crearUsuariosPrueba()}
	repo := UsuarioRepositorio(mock)

	especificacion := EspecificacionUsuario{
		ListaLiltros: []CriterioFiltro{
			{Campo: "estado", Operador: "=", Valor: "ACTIVO"},
		},
	}
	paginacion := Paginacion{Pagina: 1, TamanoPagina: 2}

	resultado, err := repo.Listar(context.Background(), especificacion, paginacion)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(resultado) != 2 {
		t.Errorf("Expected 2 usuarios activos en página 1, got %d", len(resultado))
	}
}

func TestListarConFiltroOrdenacionYPaginacion(t *testing.T) {
	mock := &mockRepositorio{usuarios: crearUsuariosPrueba()}
	repo := UsuarioRepositorio(mock)

	especificacion := EspecificacionUsuario{
		ListaLiltros: []CriterioFiltro{
			{Campo: "estado", Operador: "=", Valor: "ACTIVO"},
		},
	}
	paginacion := Paginacion{
		Pagina:       1,
		TamanoPagina: 10,
		Ordenaciones: []Ordenacion{{Campo: "nombre", Tipo: ASC}},
	}

	resultado, err := repo.Listar(context.Background(), especificacion, paginacion)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(resultado) != 3 {
		t.Errorf("Expected 3 usuarios activos, got %d", len(resultado))
	}

	nombres := []string{"Ana", "Carlos", "David"}
	for i, esperado := range nombres {
		if resultado[i].Nombre() != esperado {
			t.Errorf("Position %d: expected '%s', got '%s'", i, esperado, resultado[i].Nombre())
		}
	}
}

func TestListarConFiltroNoEncontrado(t *testing.T) {
	mock := &mockRepositorio{usuarios: crearUsuariosPrueba()}
	repo := UsuarioRepositorio(mock)

	especificacion := EspecificacionUsuario{
		ListaLiltros: []CriterioFiltro{
			{Campo: "nombre", Operador: "=", Valor: "NoExiste"},
		},
	}

	resultado, err := repo.Listar(context.Background(), especificacion, Paginacion{Pagina: 1, TamanoPagina: 10})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(resultado) != 0 {
		t.Errorf("Expected 0 usuarios, got %d", len(resultado))
	}
}

func TestListarConOperadorDesigualdad(t *testing.T) {
	mock := &mockRepositorio{usuarios: crearUsuariosPrueba()}
	repo := UsuarioRepositorio(mock)

	especificacion := EspecificacionUsuario{
		ListaLiltros: []CriterioFiltro{
			{Campo: "estado", Operador: "!=", Valor: "ACTIVO"},
		},
	}

	resultado, err := repo.Listar(context.Background(), especificacion, Paginacion{Pagina: 1, TamanoPagina: 10})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(resultado) != 2 {
		t.Errorf("Expected 2 usuarios no activos, got %d", len(resultado))
	}
}
