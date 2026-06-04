package application

import (
	"fmt"
	"unicode"
)

// ValidarFormatoPassword valida que una contraseña cumpla con los requisitos de seguridad:
// - Mínimo 8 caracteres
// - Al menos 1 mayúscula
// - Al menos 1 minúscula
// - Al menos 1 número
// - Al menos 1 carácter no alfanumérico
func ValidarFormatoPassword(password, campo string) error {
	if len(password) < 8 {
		return fmt.Errorf("%s debe tener al menos 8 caracteres", campo)
	}
	var mayuscula, minuscula, numero, noAlfanumerico bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			mayuscula = true
		case unicode.IsLower(c):
			minuscula = true
		case unicode.IsDigit(c):
			numero = true
		case !unicode.IsLetter(c) && !unicode.IsDigit(c):
			noAlfanumerico = true
		}
	}
	if !mayuscula {
		return fmt.Errorf("%s debe contener al menos una mayúscula", campo)
	}
	if !minuscula {
		return fmt.Errorf("%s debe contener al menos una minúscula", campo)
	}
	if !numero {
		return fmt.Errorf("%s debe contener al menos un número", campo)
	}
	if !noAlfanumerico {
		return fmt.Errorf("%s debe contener al menos un carácter no alfanumérico", campo)
	}
	return nil
}
