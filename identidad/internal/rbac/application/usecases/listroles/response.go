package listroles

type RolDTO struct {
	ID          string
	Nombre      string
	Descripcion string
	EsSistema   bool
	Permisos    []string
}

type RespuestaListarRoles struct {
	Roles  []RolDTO
	Total  int
	Pagina int
}
