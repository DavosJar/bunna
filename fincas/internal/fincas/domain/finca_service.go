package domain

// FincaService es un servicio de dominio que coordina operaciones entre Finca y Lote
type FincaService struct{}

func NewFincaService() *FincaService {
	return &FincaService{}
}

// EliminarFincaConLotes valida si se puede eliminar una finca según sus lotes
func (s *FincaService) EliminarFincaConLotes(finca *Finca, cantidadLotes int, confirmado bool) error {
	if finca.TieneLotes(cantidadLotes) && !confirmado {
		return ErrFincaConLotes(cantidadLotes)
	}
	return finca.CambiarEstado(FincaPendienteEliminar)
}
