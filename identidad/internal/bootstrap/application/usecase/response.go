package usecase

import "time"

// RespuestaCrearPrimerSysAdmin es el resultado de Ejecutar.
//
//   - Si YaExistia == true: no se creó nada; ExistenteID contiene el ID del
//     sys_admin preexistente. El resto de los campos quedan vacíos.
//   - Si YaExistia == false: se creó el sys_admin y los campos describen al
//     nuevo usuario (Estado="ACTIVO", Verificado=true).
type RespuestaCrearPrimerSysAdmin struct {
	UsuarioID   string
	Nombre      string
	Apellido    string
	Correo      string
	Estado      string // "ACTIVO"
	Verificado  bool   // true
	CreadoEn    time.Time
	YaExistia   bool   // true → no se creó nada
	ExistenteID string // ID del sys_admin preexistente si YaExistia
}
