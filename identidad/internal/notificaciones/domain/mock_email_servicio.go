package notificaciones

import "context"

type EmailCapturado struct {
	Destinatario string
	Asunto       string
	Cuerpo       string
	Tipo         TipoEmail
	Datos        map[string]string
}

type MockEmailServicio struct {
	EmailsEnviados []EmailCapturado
	ErrorEnviar    error
	ErrorTemplate  error
}

func (m *MockEmailServicio) Enviar(ctx context.Context, destinatario, asunto, cuerpo string) error {
	if m.ErrorEnviar != nil {
		return m.ErrorEnviar
	}
	m.EmailsEnviados = append(m.EmailsEnviados, EmailCapturado{
		Destinatario: destinatario,
		Asunto:       asunto,
		Cuerpo:       cuerpo,
	})
	return nil
}

func (m *MockEmailServicio) EnviarTemplate(ctx context.Context, destinatario string, tipo TipoEmail, datos map[string]string) error {
	if m.ErrorTemplate != nil {
		return m.ErrorTemplate
	}
	m.EmailsEnviados = append(m.EmailsEnviados, EmailCapturado{
		Destinatario: destinatario,
		Tipo:         tipo,
		Datos:        datos,
	})
	return nil
}

func (m *MockEmailServicio) UltimoEmail() *EmailCapturado {
	if len(m.EmailsEnviados) == 0 {
		return nil
	}
	return &m.EmailsEnviados[len(m.EmailsEnviados)-1]
}

func (m *MockEmailServicio) Limpiar() {
	m.EmailsEnviados = nil
}
