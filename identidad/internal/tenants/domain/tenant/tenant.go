package tenant

import (
	"regexp"
	"time"
)

var regexSlug = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type Tenant struct {
	id                 string
	nombre             string
	slug               string
	activo             bool
	fechaCreacion      time.Time
	fechaActualizacion time.Time
}

// NuevoTenant crea un tenant con validación completa.
func NuevoTenant(id, nombre, slug string) (*Tenant, error) {
	if nombre == "" {
		return nil, ErrNombreRequerido
	}
	if len(nombre) > 200 {
		return nil, ErrNombreRequerido
	}
	if slug == "" {
		return nil, ErrSlugRequerido
	}
	if len(slug) > 100 || !regexSlug.MatchString(slug) {
		return nil, ErrSlugInvalido
	}

	ahora := time.Now()
	return &Tenant{
		id:                 id,
		nombre:             nombre,
		slug:               slug,
		activo:             true,
		fechaCreacion:      ahora,
		fechaActualizacion: ahora,
	}, nil
}

// NuevoTenantDesdeBD reconstruye un tenant desde persistencia sin validar.
func NuevoTenantDesdeBD(id, nombre, slug string, activo bool, fechaCreacion, fechaActualizacion time.Time) *Tenant {
	return &Tenant{
		id:                 id,
		nombre:             nombre,
		slug:               slug,
		activo:             activo,
		fechaCreacion:      fechaCreacion,
		fechaActualizacion: fechaActualizacion,
	}
}

// Activar cambia el estado a activo.
func (t *Tenant) Activar() error {
	if t.activo {
		return ErrTenantYaActivo
	}
	t.activo = true
	t.fechaActualizacion = time.Now()
	return nil
}

// Desactivar cambia el estado a inactivo.
func (t *Tenant) Desactivar() error {
	if !t.activo {
		return ErrTenantYaInactivo
	}
	t.activo = false
	t.fechaActualizacion = time.Now()
	return nil
}

// EstaActivo retorna si el tenant está activo.
func (t *Tenant) EstaActivo() bool { return t.activo }

// Getters
func (t *Tenant) ID() string                   { return t.id }
func (t *Tenant) Nombre() string               { return t.nombre }
func (t *Tenant) Slug() string                 { return t.slug }
func (t *Tenant) FechaCreacion() time.Time     { return t.fechaCreacion }
func (t *Tenant) FechaActualizacion() time.Time { return t.fechaActualizacion }