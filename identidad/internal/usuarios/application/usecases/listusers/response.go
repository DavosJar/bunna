package listusers

type UsuarioDTO struct {
	ID       string
	Correo   string
	Nombre   string
	Apellido string
	Estado   string
	CreadoEn string
}

type RespuestaListarUsuarios struct {
	Usuarios []UsuarioDTO
	Total    int
	Pagina   int
	Tamano   int
}
