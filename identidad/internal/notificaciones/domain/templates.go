package notificaciones

import "strings"

var templates = map[TipoEmail]struct {
	Asunto string
	Cuerpo string
}{
	TipoVerificacionCorreo: {
		Asunto: "Verifica tu dirección de correo electrónico — CafeScan",
		Cuerpo: `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="margin:0;padding:0;background-color:#f4f4f4;font-family:Arial,sans-serif;">
  <table width="100%" cellpadding="0" cellspacing="0">
    <tr>
      <td align="center" style="padding:40px 0;">
        <table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,0.08);">
          <tr>
            <td style="background:#1a3a2a;padding:32px 40px;text-align:center;">
              <h1 style="color:#ffffff;margin:0;font-size:24px;font-weight:700;">CafeScan</h1>
              <p style="color:#a8c5b0;margin:4px 0 0;font-size:13px;">Diagnóstico de Nitrógeno</p>
            </td>
          </tr>
          <tr>
            <td style="padding:40px;">
              <h2 style="color:#1a3a2a;font-size:22px;margin:0 0 12px;">Hola {{nombre}},</h2>
              <p style="color:#4a5568;font-size:15px;line-height:1.6;margin:0 0 24px;">
                Gracias por registrarte en CafeScan. Para activar tu cuenta, verifica tu dirección de correo electrónico haciendo clic en el botón de abajo.
              </p>
              <table width="100%" cellpadding="0" cellspacing="0">
                <tr>
                  <td align="center" style="padding:8px 0 32px;">
                    <a href="{{url_verificacion}}" style="display:inline-block;background:#2d6a4f;color:#ffffff;text-decoration:none;font-size:16px;font-weight:600;padding:14px 40px;border-radius:8px;">
                      Verificar correo electrónico
                    </a>
                  </td>
                </tr>
              </table>
              <p style="color:#718096;font-size:13px;line-height:1.6;margin:0 0 8px;">
                Este enlace expirará en <strong>{{expiracion_horas}} horas</strong>.
              </p>
              <p style="color:#718096;font-size:13px;line-height:1.6;margin:0;">
                Si no solicitaste esta verificación, puedes ignorar este mensaje.
              </p>
            </td>
          </tr>
          <tr>
            <td style="background:#f7fafc;padding:20px 40px;text-align:center;border-top:1px solid #e2e8f0;">
              <p style="color:#a0aec0;font-size:12px;margin:0;">
                El equipo de CafeScan
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`,
	},
	TipoRecuperacionContrasena: {
		Asunto: "Recuperación de contraseña — CafeScan",
		Cuerpo: `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="margin:0;padding:0;background-color:#f4f4f4;font-family:Arial,sans-serif;">
  <table width="100%" cellpadding="0" cellspacing="0">
    <tr>
      <td align="center" style="padding:40px 0;">
        <table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,0.08);">
          <tr>
            <td style="background:#1a3a2a;padding:32px 40px;text-align:center;">
              <h1 style="color:#ffffff;margin:0;font-size:24px;font-weight:700;">CafeScan</h1>
              <p style="color:#a8c5b0;margin:4px 0 0;font-size:13px;">Diagnóstico de Nitrógeno</p>
            </td>
          </tr>
          <tr>
            <td style="padding:40px;">
              <h2 style="color:#1a3a2a;font-size:22px;margin:0 0 12px;">Hola {{nombre}},</h2>
              <p style="color:#4a5568;font-size:15px;line-height:1.6;margin:0 0 24px;">
                Recibimos una solicitud para restablecer tu contraseña. Haz clic en el botón de abajo para continuar.
              </p>
              <table width="100%" cellpadding="0" cellspacing="0">
                <tr>
                  <td align="center" style="padding:8px 0 32px;">
                    <a href="{{url_recuperacion}}" style="display:inline-block;background:#2d6a4f;color:#ffffff;text-decoration:none;font-size:16px;font-weight:600;padding:14px 40px;border-radius:8px;">
                      Restablecer contraseña
                    </a>
                  </td>
                </tr>
              </table>
              <p style="color:#718096;font-size:13px;line-height:1.6;margin:0 0 8px;">
                Este enlace expirará en <strong>{{expiracion_horas}} horas</strong>.
              </p>
              <p style="color:#718096;font-size:13px;line-height:1.6;margin:0;">
                Si no solicitaste este cambio, puedes ignorar este mensaje.
              </p>
            </td>
          </tr>
          <tr>
            <td style="background:#f7fafc;padding:20px 40px;text-align:center;border-top:1px solid #e2e8f0;">
              <p style="color:#a0aec0;font-size:12px;margin:0;">
                El equipo de CafeScan
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`,
	},
	TipoInvitacion: {
		Asunto: "Has sido invitado a colaborar en Bunna",
		Cuerpo: `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="margin:0;padding:0;background-color:#f4f4f4;font-family:Arial,sans-serif;">
  <table width="100%" cellpadding="0" cellspacing="0">
    <tr>
      <td align="center" style="padding:40px 0;">
        <table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,0.08);">
          <tr>
            <td style="background:#1a3a2a;padding:32px 40px;text-align:center;">
              <h1 style="color:#ffffff;margin:0;font-size:24px;font-weight:700;">Bunna</h1>
              <p style="color:#a8c5b0;margin:4px 0 0;font-size:13px;">Gestión de cafetales</p>
            </td>
          </tr>
          <tr>
            <td style="padding:40px;">
              <h2 style="color:#1a3a2a;font-size:22px;margin:0 0 12px;">¡Hola!</h2>
              <p style="color:#4a5568;font-size:15px;line-height:1.6;margin:0 0 24px;">
                <strong>{{nombre_tenant}}</strong> te ha invitado a colaborar en su equipo. Haz clic en el botón de abajo para aceptar la invitación.
              </p>
              <table width="100%" cellpadding="0" cellspacing="0">
                <tr>
                  <td align="center" style="padding:8px 0 32px;">
                    <a href="{{url_invitacion}}" style="display:inline-block;background:#2d6a4f;color:#ffffff;text-decoration:none;font-size:16px;font-weight:600;padding:14px 40px;border-radius:8px;">
                      Aceptar invitación
                    </a>
                  </td>
                </tr>
              </table>
              <p style="color:#718096;font-size:13px;line-height:1.6;margin:0 0 8px;">
                Este enlace expirará en <strong>{{expiracion_horas}} horas</strong>.
              </p>
              <p style="color:#718096;font-size:13px;line-height:1.6;margin:0;">
                Si no esperabas esta invitación, puedes ignorar este mensaje.
              </p>
            </td>
          </tr>
          <tr>
            <td style="background:#f7fafc;padding:20px 40px;text-align:center;border-top:1px solid #e2e8f0;">
              <p style="color:#a0aec0;font-size:12px;margin:0;">
                El equipo de Bunna
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`,
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
