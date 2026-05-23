package notificaciones

import "context"

// EmailServicio define la interfaz para envío de emails
type EmailServicio interface {
	Enviar(ctx context.Context, destinatario, asunto, cuerpo string) error
	EnviarTemplate(ctx context.Context, destinatario string, tipo TipoEmail, datos map[string]string) error
}