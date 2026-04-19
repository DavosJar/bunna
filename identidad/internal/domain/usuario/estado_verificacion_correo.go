package usuario

type EstadoVerificacionCorreo string

const (
	PENDIENTE_VERIFICACION EstadoVerificacionCorreo = "PENDIENTE_VERIFICACION"
	VERIFICADO             EstadoVerificacionCorreo = "VERIFICADO"
	ENLACE_EXPIRADO        EstadoVerificacionCorreo = "ENLACE_EXPIRADO"
	REENVIO_SOLICITADO     EstadoVerificacionCorreo = "REENVIO_SOLICITADO"
	VERIFICACION_FALLIDA   EstadoVerificacionCorreo = "VERIFICACION_FALLIDA"
)

var TransicionesVerificacion = map[EstadoVerificacionCorreo][]EstadoVerificacionCorreo{
	PENDIENTE_VERIFICACION: {VERIFICADO, ENLACE_EXPIRADO, VERIFICACION_FALLIDA, REENVIO_SOLICITADO},
	ENLACE_EXPIRADO:        {REENVIO_SOLICITADO},
	VERIFICACION_FALLIDA:   {REENVIO_SOLICITADO},
	REENVIO_SOLICITADO:     {VERIFICADO, ENLACE_EXPIRADO, VERIFICACION_FALLIDA},
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
