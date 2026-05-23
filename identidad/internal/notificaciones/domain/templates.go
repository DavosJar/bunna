package notificaciones

import "strings"

var templates = map[TipoEmail]struct {
	Asunto string
	Cuerpo string
}{
	TipoVerificacionCorreo: {
		Asunto: "Verifica tu dirección de correo electrónico",
		Cuerpo: `Hola {{nombre}},

Gracias por registrarte. Para verificar tu dirección de correo electrónico,
usa el siguiente token:

{{token}}

Este enlace expirará en {{expiracion_horas}} horas.

Si no solicitaste esta verificación, puedes ignorar este mensaje.

Saludos,
El equipo de CafeScan`,
	},
	TipoRecuperacionContrasena: {
		Asunto: "Recuperación de contraseña",
		Cuerpo: `Hola {{nombre}},

Recibimos una solicitud para restablecer tu contraseña.
Usa el siguiente token para continuar:

{{token}}

Este enlace expirará en {{expiracion_horas}} horas.

Si no solicitaste este cambio, puedes ignorar este mensaje.

Saludos,
El equipo de CafeScan`,
	},
}

func RenderizarTemplate(tipo TipoEmail, datos map[string]string) (asunto, cuerpo string, err error) {
	tmpl, existe := templates[tipo]
	if !existe {
		return "", "", ErrTemplateNoEncontrado
	}
	asunto = tmpl.Asunto
	cuerpo = tmpl.Cuerpo
	for clave, valor := range datos {
		marcador := "{{" + clave + "}}"
		cuerpo = strings.ReplaceAll(cuerpo, marcador, valor)
	}
	return asunto, cuerpo, nil
}
