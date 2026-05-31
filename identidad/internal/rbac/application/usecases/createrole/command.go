package createrole

type ComandoCrearRol struct {
	Nombre      string
	Descripcion string
	Permisos    []string
	TenantID    string
	EjecutorID  string
}
