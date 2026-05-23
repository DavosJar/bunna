package usuario

type EstadoVerificacionCorreo string

const (
	PENDIENTE_VERIFICACION EstadoVerificacionCorreo = "PENDIENTE_VERIFICACION"
	VERIFICADO             EstadoVerificacionCorreo = "VERIFICADO"
	ENLACE_EXPIRADO        EstadoVerificacionCorreo = "ENLACE_EXPIRADO"
	REENVIO_SOLICITADO     EstadoVerificacionCorreo = "REENVIO_SOLICITADO"
)

var TransicionesVerificacion = map[EstadoVerificacionCorreo][]EstadoVerificacionCorreo{
	PENDIENTE_VERIFICACION: {VERIFICADO, ENLACE_EXPIRADO, REENVIO_SOLICITADO},
	ENLACE_EXPIRADO:        {REENVIO_SOLICITADO},
	REENVIO_SOLICITADO:     {VERIFICADO, ENLACE_EXPIRADO},
	VERIFICADO:             {}, // terminal
}

func (e EstadoVerificacionCorreo) PuedeTransicionarA(destino EstadoVerificacionCorreo) bool {
	destinosPermitidos, existe := TransicionesVerificacion[e]
	if !existe {
		return false
	}
	for _, d := range destinosPermitidos {
		if d == destino {
			return true
		}
	}
	return false
}