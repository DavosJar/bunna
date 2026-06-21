package notificaciones

import (
	"strings"
	"testing"
)

func TestRenderizarTemplateVerificacion(t *testing.T) {
	asunto, cuerpo, err := RenderizarTemplate(TipoVerificacionCorreo, map[string]string{
		"nombre":           "Juan",
		"url_verificacion": "https://bunna.app/verificar?token=abc-123",
		"expiracion_horas": "24",
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if asunto == "" {
		t.Error("Expected asunto no vacío")
	}
	if !strings.Contains(cuerpo, "Juan") {
		t.Error("Expected cuerpo con nombre 'Juan'")
	}
	if !strings.Contains(cuerpo, "bunna.app/verificar") {
		t.Error("Expected cuerpo con url_verificacion")
	}
	if !strings.Contains(cuerpo, "24") {
		t.Error("Expected cuerpo con expiracion '24'")
	}
}

func TestRenderizarTemplateRecuperacion(t *testing.T) {
	asunto, cuerpo, err := RenderizarTemplate(TipoRecuperacionContrasena, map[string]string{
		"nombre":           "Maria",
		"url_recuperacion": "https://bunna.app/recuperar?token=xyz-456",
		"expiracion_horas": "2",
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if asunto == "" {
		t.Error("Expected asunto no vacío")
	}
	if !strings.Contains(cuerpo, "Maria") {
		t.Error("Expected cuerpo con nombre 'Maria'")
	}
	if !strings.Contains(cuerpo, "bunna.app/recuperar") {
		t.Error("Expected cuerpo con url_recuperacion")
	}
}

func TestRenderizarTemplateNoEncontrado(t *testing.T) {
	_, _, err := RenderizarTemplate("TIPO_INEXISTENTE", map[string]string{})
	if err != ErrTemplateNoEncontrado {
		t.Errorf("Expected ErrTemplateNoEncontrado, got %v", err)
	}
}

func TestRenderizarTemplateSinVariables(t *testing.T) {
	_, cuerpo, err := RenderizarTemplate(TipoVerificacionCorreo, map[string]string{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !strings.Contains(cuerpo, "{{nombre}}") {
		t.Error("Expected marcadores sin reemplazar cuando no se pasan datos")
	}
}

func TestTiposEmailConstantes(t *testing.T) {
	if TipoVerificacionCorreo == "" {
		t.Error("TipoVerificacionCorreo no debe estar vacío")
	}
	if TipoRecuperacionContrasena == "" {
		t.Error("TipoRecuperacionContrasena no debe estar vacío")
	}
	if TipoVerificacionCorreo == TipoRecuperacionContrasena {
		t.Error("Los tipos de email deben ser distintos")
	}
}
