package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"time"

	notificaciones "github.com/davosjar/bunna/services/identidad/internal/notificaciones/domain"
)

// ConfigSMTP contiene la configuración del servidor SMTP
type ConfigSMTP struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
	Timeout  time.Duration
	Async    bool
}

// SMTPServicio implementa EmailServicio usando net/smtp
type SMTPServicio struct {
	config ConfigSMTP
}

// NewSMTPServicio crea una nueva instancia del servicio SMTP
func NewSMTPServicio(config ConfigSMTP) notificaciones.EmailServicio {
	if config.From == "" {
		config.From = config.User
	}
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}
	return &SMTPServicio{config: config}
}

// Enviar envía un email con asunto y cuerpo dados
func (s *SMTPServicio) Enviar(ctx context.Context, destinatario, asunto, cuerpo string) error {
	if destinatario == "" {
		return notificaciones.ErrDestinatarioInvalido
	}
	if asunto == "" {
		return notificaciones.ErrAsuntoRequerido
	}
	if cuerpo == "" {
		return notificaciones.ErrCuerpoRequerido
	}

	if s.config.Async {
		go func() {
			if err := s.enviar(destinatario, asunto, cuerpo); err != nil {
				log.Printf("[EmailServicio] Error al enviar email a %s: %v", destinatario, err)
			}
		}()
		return nil
	}

	return s.enviar(destinatario, asunto, cuerpo)
}

// EnviarTemplate renderiza un template y envía el email
func (s *SMTPServicio) EnviarTemplate(ctx context.Context, destinatario string, tipo notificaciones.TipoEmail, datos map[string]string) error {
	asunto, cuerpo, err := notificaciones.RenderizarTemplate(tipo, datos)
	if err != nil {
		return err
	}
	return s.Enviar(ctx, destinatario, asunto, cuerpo)
}

// enviar realiza el envío SMTP real
func (s *SMTPServicio) enviar(destinatario, asunto, cuerpo string) error {
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)

	auth := smtp.PlainAuth("", s.config.User, s.config.Password, s.config.Host)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.config.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		// Intentar sin TLS (STARTTLS en puerto 587)
		client, err := smtp.Dial(addr)
		if err != nil {
			return fmt.Errorf("%w: %v", notificaciones.ErrSMTPConexionFallida, err)
		}
		defer client.Close()

		if err = client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("%w: %v", notificaciones.ErrSMTPConexionFallida, err)
		}

		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("error de autenticación SMTP: %v", err)
		}

		return s.enviarConCliente(client, destinatario, asunto, cuerpo)
	}

	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return fmt.Errorf("%w: %v", notificaciones.ErrSMTPConexionFallida, err)
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("error de autenticación SMTP: %v", err)
	}

	return s.enviarConCliente(client, destinatario, asunto, cuerpo)
}

func (s *SMTPServicio) enviarConCliente(client *smtp.Client, destinatario, asunto, cuerpo string) error {
	if err := client.Mail(s.config.From); err != nil {
		return fmt.Errorf("error al establecer remitente: %v", err)
	}
	if err := client.Rcpt(destinatario); err != nil {
		return fmt.Errorf("error al establecer destinatario: %v", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("error al abrir canal de datos: %v", err)
	}
	defer w.Close()

	mensaje := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		s.config.From, destinatario, asunto, cuerpo,
	)

	if _, err := fmt.Fprint(w, mensaje); err != nil {
		return fmt.Errorf("error al escribir mensaje: %v", err)
	}

	return nil
}
