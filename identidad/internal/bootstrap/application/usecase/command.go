package usecase

// ComandoCrearPrimerSysAdmin transporta los datos necesarios para
// crear el primer sys_admin del sistema.
//
// Es construido y validado por el caller (CLI u otro). El caso de uso
// aplica una validación baseline (ver usecase.go) antes de ejecutar.
type ComandoCrearPrimerSysAdmin struct {
	Nombre   string
	Apellido string
	Correo   string
	Password string
}

// ToLog retorna una representación segura del comando para logs/telemetría.
// El password NUNCA se incluye.
func (c ComandoCrearPrimerSysAdmin) ToLog() map[string]any {
	return map[string]any{
		"correo":   c.Correo,
		"nombre":   c.Nombre,
		"apellido": c.Apellido,
	}
}
