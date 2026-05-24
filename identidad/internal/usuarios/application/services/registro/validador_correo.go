package registro

import (
	"fmt"
	"net"
)

// validarDominioCorreo verifica que el dominio del correo tenga registros MX
func validarDominioCorreo(correo string) error {
	// Extraer dominio
	dominio := ""
	for i := len(correo) - 1; i >= 0; i-- {
		if correo[i] == '@' {
			dominio = correo[i+1:]
			break
		}
	}
	if dominio == "" {
		return fmt.Errorf("dominio inválido")
	}

	// Consultar registros MX
	mxRecords, err := net.LookupMX(dominio)
	if err != nil || len(mxRecords) == 0 {
		// Intentar con registros A como fallback
		_, err = net.LookupHost(dominio)
		if err != nil {
			return fmt.Errorf("el dominio '%s' no existe o no acepta emails", dominio)
		}
	}

	return nil
}