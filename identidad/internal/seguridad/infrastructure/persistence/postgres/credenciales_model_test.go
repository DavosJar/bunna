package postgres

import (
	"testing"
	"time"

	"github.com/davosjar/bunna/services/identidad/internal/seguridad/domain"
)

func TestCredencialesModelTableName(t *testing.T) {
	// Arrange
	model := &CredencialesModel{}

	// Act
	tableName := model.TableName()

	// Assert
	if tableName != "credenciales_usuarios" {
		t.Errorf("TableName() = %q, want %q", tableName, "credenciales_usuarios")
	}
}

func TestCredencialesModelToDomain(t *testing.T) {
	tests := []struct {
		name     string
		model    *CredencialesModel
		expected func(*domain.CredencialesUsuario) bool
	}{
		{
			name: "when converting model with active credential, returns domain with correct values",
			model: &CredencialesModel{
				UsuarioID:        "user123",
				PasswordHash:     "hash_value",
				Activo:           true,
				CorreoVerificado: true,
				IntentosFallidos: 0,
				BloqueadoHasta:   time.Time{},
			},
			expected: func(c *domain.CredencialesUsuario) bool {
				return c.UsuarioID() == "user123" &&
					c.PasswordHash() == "hash_value" &&
					c.Activo() == true &&
					c.CorreoVerificado() == true &&
					c.IntentosFallidos() == 0 &&
					c.BloqueadoHasta().IsZero()
			},
		},
		{
			name: "when converting model with blocked credential, returns domain with block time",
			model: &CredencialesModel{
				UsuarioID:        "user456",
				PasswordHash:     "blocked_hash",
				Activo:           false,
				CorreoVerificado: false,
				IntentosFallidos: 5,
				BloqueadoHasta:   time.Now().Add(15 * time.Minute),
			},
			expected: func(c *domain.CredencialesUsuario) bool {
				return c.UsuarioID() == "user456" &&
					c.PasswordHash() == "blocked_hash" &&
					c.Activo() == false &&
					c.CorreoVerificado() == false &&
					c.IntentosFallidos() == 5 &&
					!c.BloqueadoHasta().IsZero()
			},
		},
		{
			name: "when converting model with zero values, returns domain with those values",
			model: &CredencialesModel{
				UsuarioID:        "user789",
				PasswordHash:     "",
				Activo:           false,
				CorreoVerificado: false,
				IntentosFallidos: 0,
				BloqueadoHasta:   time.Time{},
			},
			expected: func(c *domain.CredencialesUsuario) bool {
				return c.UsuarioID() == "user789" &&
					c.PasswordHash() == "" &&
					c.Activo() == false &&
					c.CorreoVerificado() == false &&
					c.IntentosFallidos() == 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.model.ToDomain()

			// Assert
			if !tt.expected(result) {
				t.Errorf("ToDomain() returned unexpected values")
			}
		})
	}
}

func TestCredencialesFromDomain(t *testing.T) {
	tests := []struct {
		name     string
		domain   *domain.CredencialesUsuario
		expected func(*CredencialesModel) bool
	}{
		{
			name:   "when converting active domain credential, returns model with correct values",
			domain: domain.NuevaCredencialesUsuario("user123", "hash_value"),
			expected: func(m *CredencialesModel) bool {
				return m.UsuarioID == "user123" &&
					m.PasswordHash == "hash_value" &&
					m.Activo == true &&
					m.CorreoVerificado == false &&
					m.IntentosFallidos == 0 &&
					m.BloqueadoHasta.IsZero()
			},
		},
		{
			name: "when converting domain credential with failed attempts, returns model with attempt count",
			domain: func() *domain.CredencialesUsuario {
				c := domain.NuevaCredencialesUsuario("user456", "hash_locked")
				ahora := time.Now()
				c.MarcarIntentoFallido(ahora)
				c.MarcarIntentoFallido(ahora)
				c.MarcarIntentoFallido(ahora)
				c.MarcarIntentoFallido(ahora)
				c.MarcarIntentoFallido(ahora)
				return c
			}(),
			expected: func(m *CredencialesModel) bool {
				return m.UsuarioID == "user456" &&
					m.PasswordHash == "hash_locked" &&
					m.IntentosFallidos == 5 &&
					!m.BloqueadoHasta.IsZero()
			},
		},
		{
			name: "when converting domain credential with verified email, returns model with email verified",
			domain: func() *domain.CredencialesUsuario {
				c := domain.NuevaCredencialesUsuario("user789", "hash_verified")
				c.VerificarCorreo()
				return c
			}(),
			expected: func(m *CredencialesModel) bool {
				return m.UsuarioID == "user789" &&
					m.CorreoVerificado == true &&
					m.Activo == true
			},
		},
		{
			name: "when converting domain credential and then back to domain, preserves all values",
			domain: func() *domain.CredencialesUsuario {
				c := domain.NuevaCredencialesUsuarioDesdeBD(
					"user999",
					"complex_hash",
					false,
					true,
					2,
					time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
				)
				return c
			}(),
			expected: func(m *CredencialesModel) bool {
				return m.UsuarioID == "user999" &&
					m.PasswordHash == "complex_hash" &&
					m.Activo == false &&
					m.CorreoVerificado == true &&
					m.IntentosFallidos == 2 &&
					m.BloqueadoHasta.Equal(time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result, err := CredencialesFromDomain(tt.domain)

			// Assert
			if err != nil {
				t.Fatalf("CredencialesFromDomain() error = %v", err)
			}
			if !tt.expected(result) {
				t.Errorf("CredencialesFromDomain() returned unexpected model")
			}
		})
	}
}

func TestCredencialesModelRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		cred *domain.CredencialesUsuario
	}{
		{
			name: "roundtrip preserves active credential data",
			cred: domain.NuevaCredencialesUsuario("user1", "hash1"),
		},
		{
			name: "roundtrip preserves blocked credential data",
			cred: domain.NuevaCredencialesUsuarioDesdeBD(
				"user2",
				"hash2",
				false,
				false,
				5,
				time.Now().Add(15*time.Minute),
			),
		},
		{
			name: "roundtrip preserves credential with verified email",
			cred: func() *domain.CredencialesUsuario {
				c := domain.NuevaCredencialesUsuario("user3", "hash3")
				c.VerificarCorreo()
				return c
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act - convert domain to model and back
			model, err := CredencialesFromDomain(tt.cred)
			if err != nil {
				t.Fatalf("CredencialesFromDomain() error = %v", err)
			}
			recovered := model.ToDomain()

			// Assert - verify all fields match
			if recovered.UsuarioID() != tt.cred.UsuarioID() {
				t.Errorf("UsuarioID roundtrip failed: got %q, want %q",
					recovered.UsuarioID(), tt.cred.UsuarioID())
			}
			if recovered.PasswordHash() != tt.cred.PasswordHash() {
				t.Errorf("PasswordHash roundtrip failed")
			}
			if recovered.Activo() != tt.cred.Activo() {
				t.Errorf("Activo roundtrip failed: got %v, want %v",
					recovered.Activo(), tt.cred.Activo())
			}
			if recovered.CorreoVerificado() != tt.cred.CorreoVerificado() {
				t.Errorf("CorreoVerificado roundtrip failed")
			}
			if recovered.IntentosFallidos() != tt.cred.IntentosFallidos() {
				t.Errorf("IntentosFallidos roundtrip failed")
			}
		})
	}
}
