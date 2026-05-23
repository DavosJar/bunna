package notificaciones

import "errors"

var (
	ErrDestinatarioInvalido  = errors.New("destinatario de email inválido")
	ErrAsuntoRequerido       = errors.New("asunto del email requerido")
	ErrCuerpoRequerido       = errors.New("cuerpo del email requerido")
	ErrTemplateNoEncontrado  = errors.New("template de email no encontrado")
	ErrSMTPConexionFallida   = errors.New("error al conectar con servidor SMTP")
)