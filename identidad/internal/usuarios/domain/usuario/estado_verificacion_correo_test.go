package usuario

import (
	"testing"
)

func TestEstadoVerificacionCorreoValorPendiente(t *testing.T) {
	estado := PENDIENTE_VERIFICACION
	if estado != "PENDIENTE_VERIFICACION" {
		t.Errorf("Expected 'PENDIENTE_VERIFICACION', got '%s'", estado)
	}
}

func TestEstadoVerificacionCorreoValorVerificado(t *testing.T) {
	estado := VERIFICADO
	if estado != "VERIFICADO" {
		t.Errorf("Expected 'VERIFICADO', got '%s'", estado)
	}
}

func TestPuedeTransicionarDeDesdeFormatos(t *testing.T) {
	tests := []struct {
		nombre  string
		origen  EstadoVerificacionCorreo
		destino EstadoVerificacionCorreo
		valido  bool
	}{
		{
			nombre:  "PENDIENTE a VERIFICADO",
			origen:  PENDIENTE_VERIFICACION,
			destino: VERIFICADO,
			valido:  true,
		},
		{
			nombre:  "PENDIENTE a ENLACE_EXPIRADO",
			origen:  PENDIENTE_VERIFICACION,
			destino: ENLACE_EXPIRADO,
			valido:  true,
		},
		{
			nombre:  "PENDIENTE a REENVIO_SOLICITADO",
			origen:  PENDIENTE_VERIFICACION,
			destino: REENVIO_SOLICITADO,
			valido:  true,
		},
		{
			nombre:  "ENLACE_EXPIRADO a REENVIO_SOLICITADO",
			origen:  ENLACE_EXPIRADO,
			destino: REENVIO_SOLICITADO,
			valido:  true,
		},
		{
			nombre:  "ENLACE_EXPIRADO a VERIFICADO (inválido)",
			origen:  ENLACE_EXPIRADO,
			destino: VERIFICADO,
			valido:  false,
		},
		{
			nombre:  "REENVIO_SOLICITADO a VERIFICADO",
			origen:  REENVIO_SOLICITADO,
			destino: VERIFICADO,
			valido:  true,
		},
		{
			nombre:  "REENVIO_SOLICITADO a ENLACE_EXPIRADO",
			origen:  REENVIO_SOLICITADO,
			destino: ENLACE_EXPIRADO,
			valido:  true,
		},
		{
			nombre:  "VERIFICADO es terminal (no puede ir a PENDIENTE)",
			origen:  VERIFICADO,
			destino: PENDIENTE_VERIFICACION,
			valido:  false,
		},
		{
			nombre:  "VERIFICADO es terminal (no puede ir a ENLACE_EXPIRADO)",
			origen:  VERIFICADO,
			destino: ENLACE_EXPIRADO,
			valido:  false,
		},
		{
			nombre:  "VERIFICADO es terminal (no puede ir a REENVIO_SOLICITADO)",
			origen:  VERIFICADO,
			destino: REENVIO_SOLICITADO,
			valido:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			resultado := tt.origen.PuedeTransicionarA(tt.destino)
			if resultado != tt.valido {
				t.Errorf("PuedeTransicionarA(%s -> %s): esperado %v, obtuvo %v",
					tt.origen, tt.destino, tt.valido, resultado)
			}
		})
	}
}

func TestPuedeTransicionarConEstadoInvalido(t *testing.T) {
	estadoInvalido := EstadoVerificacionCorreo("ESTADO_INEXISTENTE")
	resultado := estadoInvalido.PuedeTransicionarA(VERIFICADO)

	if resultado != false {
		t.Errorf("Expected false para estado inválido, obtuvo %v", resultado)
	}
}

func TestTransicionesVerificacionMapaCompleto(t *testing.T) {
	// El mapa debe tener exactamente 4 estados (sin VERIFICACION_FALLIDA)
	if len(TransicionesVerificacion) != 4 {
		t.Fatalf("Se esperaban 4 estados en TransicionesVerificacion, obtuvo %d", len(TransicionesVerificacion))
	}

	estadosEsperados := []EstadoVerificacionCorreo{
		PENDIENTE_VERIFICACION,
		ENLACE_EXPIRADO,
		REENVIO_SOLICITADO,
		VERIFICADO,
	}

	for _, estado := range estadosEsperados {
		if _, existe := TransicionesVerificacion[estado]; !existe {
			t.Errorf("Estado %s no encontrado en TransicionesVerificacion", estado)
		}
	}
}

func TestVerificadoEsTerminal(t *testing.T) {
	transiciones, existe := TransicionesVerificacion[VERIFICADO]

	if !existe {
		t.Fatal("VERIFICADO no encontrado en mapa")
	}

	if len(transiciones) != 0 {
		t.Errorf("VERIFICADO debería ser terminal (0 transiciones), obtuvo %d", len(transiciones))
	}

	estadosVerificacion := []EstadoVerificacionCorreo{
		PENDIENTE_VERIFICACION,
		ENLACE_EXPIRADO,
		REENVIO_SOLICITADO,
	}

	for _, estado := range estadosVerificacion {
		if VERIFICADO.PuedeTransicionarA(estado) {
			t.Errorf("VERIFICADO no debería poder transicionar a %s", estado)
		}
	}
}