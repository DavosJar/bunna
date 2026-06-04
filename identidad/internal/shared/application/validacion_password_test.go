package application_test

import (
	"testing"

	"github.com/davosjar/bunna/services/identidad/internal/shared/application"
)

func TestValidarFormatoPassword_Valido(t *testing.T) {
	err := application.ValidarFormatoPassword("Secreto1!", "password")
	if err != nil {
		t.Errorf("password válido no debería dar error: %v", err)
	}
}

func TestValidarFormatoPassword_MuyCorto(t *testing.T) {
	err := application.ValidarFormatoPassword("Secre1!", "password")
	if err == nil {
		t.Fatal("password muy corto debería dar error")
	}
	if err.Error() != "password debe tener al menos 8 caracteres" {
		t.Errorf("mensaje de error incorrecto: %v", err)
	}
}

func TestValidarFormatoPassword_SinMayuscula(t *testing.T) {
	err := application.ValidarFormatoPassword("secreto1!", "password")
	if err == nil {
		t.Fatal("password sin mayúscula debería dar error")
	}
	if err.Error() != "password debe contener al menos una mayúscula" {
		t.Errorf("mensaje de error incorrecto: %v", err)
	}
}

func TestValidarFormatoPassword_SinMinuscula(t *testing.T) {
	err := application.ValidarFormatoPassword("SECRETO1!", "password")
	if err == nil {
		t.Fatal("password sin minúscula debería dar error")
	}
	if err.Error() != "password debe contener al menos una minúscula" {
		t.Errorf("mensaje de error incorrecto: %v", err)
	}
}

func TestValidarFormatoPassword_SinNumero(t *testing.T) {
	err := application.ValidarFormatoPassword("Secreto!", "password")
	if err == nil {
		t.Fatal("password sin número debería dar error")
	}
	if err.Error() != "password debe contener al menos un número" {
		t.Errorf("mensaje de error incorrecto: %v", err)
	}
}

func TestValidarFormatoPassword_SinNoAlfanumerico(t *testing.T) {
	err := application.ValidarFormatoPassword("Secreto1", "password")
	if err == nil {
		t.Fatal("password sin carácter no alfanumérico debería dar error")
	}
	if err.Error() != "password debe contener al menos un carácter no alfanumérico" {
		t.Errorf("mensaje de error incorrecto: %v", err)
	}
}

func TestValidarFormatoPassword_Vacio(t *testing.T) {
	err := application.ValidarFormatoPassword("", "password")
	if err == nil {
		t.Fatal("password vacío debería dar error")
	}
	if err.Error() != "password debe tener al menos 8 caracteres" {
		t.Errorf("mensaje de error incorrecto: %v", err)
	}
}

func TestValidarFormatoPassword_CampoPersonalizado(t *testing.T) {
	err := application.ValidarFormatoPassword("abc", "nueva_password")
	if err == nil {
		t.Fatal("debería dar error")
	}
	if err.Error() != "nueva_password debe tener al menos 8 caracteres" {
		t.Errorf("mensaje de error con campo personalizado incorrecto: %v", err)
	}
}
