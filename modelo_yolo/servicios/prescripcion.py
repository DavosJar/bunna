"""
Servicio de Prescripción (MOCK por ahora).
Implementa la interfaz IServicioPrescripcion.
"""

from interfaces import IServicioPrescripcion
from schemas import Prescripcion, Medicamento


class ServicioPrescripcionMock(IServicioPrescripcion):
    """Mock implementation - Retorna prescripciones ficticias."""
    
    def generar_prescripcion(self, diagnostico_id: str, enfermedad: str, severidad: str) -> Prescripcion:
        """Mock: Genera prescripción según severidad."""
        
        if severidad == "alta":
            tratamientos = [
                Medicamento(
                    nombre="Sulfato de hierro",
                    dosis="5 gramos/litro",
                    frecuencia="cada 3 días",
                    duracion="2 semanas",
                    notas="Aplicar en primeras horas del día"
                ),
                Medicamento(
                    nombre="Quelatante de hierro",
                    dosis="2 gramos/litro",
                    frecuencia="semanal",
                    duracion="1 mes"
                )
            ]
            urgencia = "inmediata"
            seguimiento = 7
        elif severidad == "media":
            tratamientos = [
                Medicamento(
                    nombre="Sulfato de hierro",
                    dosis="3 gramos/litro",
                    frecuencia="cada 5 días",
                    duracion="2 semanas"
                )
            ]
            urgencia = "media"
            seguimiento = 14
        else:  # baja
            tratamientos = [
                Medicamento(
                    nombre="Sulfato de hierro",
                    dosis="2 gramos/litro",
                    frecuencia="semanal",
                    duracion="1 mes"
                )
            ]
            urgencia = "baja"
            seguimiento = 21
        
        prescripcion = Prescripcion(
            diagnostico_id=diagnostico_id,
            tratamientos=tratamientos,
            urgencia=urgencia,
            observaciones="Mantener monitor regular. Consultar especialista si no hay mejora.",
            seguimiento_dias=seguimiento,
            contraindicaciones=["No mezclar con calcio", "Evitar luz solar directa"],
            creado_por="Sistema"
        )
        
        return prescripcion


# DESPUÉS: Reemplazar con implementación real
# class ServicioPrescripcionBD(IServicioPrescripcion):
#     def generar_prescripcion(self, diagnostico_id: str, enfermedad: str, severidad: str) -> Prescripcion:
#         protocolo = consultar_bd(enfermedad, severidad)
#         # ... procesar protocolos ...
#         return Prescripcion(...)


# Instancia por defecto (usar mock ahora)
servicio_prescripcion = ServicioPrescripcionMock()
