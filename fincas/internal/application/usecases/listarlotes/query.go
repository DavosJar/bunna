package listarlotes

import "fmt"

type Query struct {
	FincaID string
}

func (q *Query) Validar() error {
	if q.FincaID == "" {
		return fmt.Errorf("validación: el fincaID es requerido")
	}
	return nil
}
