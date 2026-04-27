"""
Interfaces/Contratos que deben cumplir los servicios.
Las implementaciones (mock, real) heredan de estas.
"""

from abc import ABC, abstractmethod
from schemas import Diagnostico, Prescripcion


class IServicioDiagnostico(ABC):
    """Interfaz que cualquier servicio de diagnóstico debe cumplir."""
    
    @abstractmethod
    def diagnosticar(self, contenido_imagen: bytes) -> Diagnostico:
        """Recibe imagen, retorna Diagnostico."""
        pass


class IServicioPrescripcion(ABC):
    """Interfaz que cualquier servicio de prescripción debe cumplir."""
    
    @abstractmethod
    def generar_prescripcion(self, diagnostico_id: str, enfermedad: str, severidad: str) -> Prescripcion:
        """Recibe datos diagnóstico, retorna Prescripcion."""
        pass
