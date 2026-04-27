"""
Servicio de Diagnóstico (MOCK por ahora).
Implementa la interfaz IServicioDiagnostico.
"""

from interfaces import IServicioDiagnostico
from schemas import Diagnostico, SintomaDiagnostico


class ServicioDiagnosticoMock(IServicioDiagnostico):
    """Mock implementation - Retorna datos ficticios."""
    
    def diagnosticar(self, contenido_imagen: bytes) -> Diagnostico:
        """Mock: Retorna diagnóstico de prueba."""
        diagnostico = Diagnostico(
            enfermedad="Clorósis",
            severidad="alta",
            confianza=0.92,
            sintomas=[
                SintomaDiagnostico(
                    nombre="Amarillamiento",
                    descripcion="Hojas presentan color amarillo intenso"
                ),
                SintomaDiagnostico(
                    nombre="Decoloración",
                    descripcion="Pérdida progresiva del pigmento verde"
                )
            ],
            area_afectada="60%",
            tipo_dano="Fotosintético",
            recomendaciones_iniciales=[
                "Aumentar frecuencia de riego",
                "Mejorar drenaje del suelo",
                "Revisar pH del terreno"
            ]
        )
        return diagnostico


# DESPUÉS: Reemplazar con implementación real
# class ServicioDiagnosticoYOLO(IServicioDiagnostico):
#     def diagnosticar(self, contenido_imagen: bytes) -> Diagnostico:
#         imagen = preprocesar_imagen(contenido_imagen)
#         resultados = modelo_yolo.predict(imagen)
#         # ... procesar resultados ...
#         return Diagnostico(...)


# Instancia por defecto (usar mock ahora)
servicio_diagnostico = ServicioDiagnosticoMock()
