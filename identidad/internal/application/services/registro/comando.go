package registro

// ComandoRegistro es el DTO de entrada para el caso de uso de registro
type ComandoRegistro struct {
	Correo   string
	Password string // Password en plano (no hasheado)
	Nombre   string
	Apellido string
	Telefono string
}
