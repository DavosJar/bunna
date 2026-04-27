"""
Esquemas/Interfaces de los datos del flujo.
Define la estructura que fluye entre servicios.
"""

from pydantic import BaseModel
from typing import List, Optional
from datetime import datetime


# ==================== DIAGNÓSTICO ====================

class SintomaDiagnostico(BaseModel):
    """Un síntoma detectado."""
    nombre: str
    descripcion: Optional[str] = None


class Diagnostico(BaseModel):
    """Resultado del diagnóstico."""
    enfermedad: str
    severidad: str  # "baja", "media", "alta"
    confianza: float  # 0-1
    sintomas: List[SintomaDiagnostico]
    area_afectada: str  # "30%", "60%", etc
    tipo_dano: str
    recomendaciones_iniciales: List[str]


# ==================== PRESCRIPCIÓN ====================

class Medicamento(BaseModel):
    """Un medicamento/tratamiento."""
    nombre: str
    dosis: str
    frecuencia: str
    duracion: Optional[str] = None
    notas: Optional[str] = None


class Prescripcion(BaseModel):
    """Prescripción basada en diagnóstico."""
    diagnostico_id: str  # UUID del diagnóstico
    tratamientos: List[Medicamento]
    urgencia: str  # "baja", "media", "alta", "inmediata"
    observaciones: str
    seguimiento_dias: int
    contraindicaciones: Optional[List[str]] = None
    creado_por: Optional[str] = "Sistema"


# ==================== RESPUESTA FINAL ====================

class RespuestaDiagnosticoPrescripcion(BaseModel):
    """Respuesta completa con diagnóstico + prescripción."""
    id: str  # UUID
    diagnostico: Diagnostico
    prescripcion: Prescripcion
    timestamp: datetime
    
    class Config:
        json_schema_extra = {
            "example": {
                "id": "550e8400-e29b-41d4-a716-446655440000",
                "diagnostico": {
                    "enfermedad": "Clorósis",
                    "severidad": "alta",
                    "confianza": 0.92,
                    "sintomas": [
                        {"nombre": "Amarillamiento", "descripcion": "Hojas amarillas"},
                        {"nombre": "Decoloración", "descripcion": "Pérdida de color"}
                    ],
                    "area_afectada": "60%",
                    "tipo_dano": "Fotosintético",
                    "recomendaciones_iniciales": ["Aumentar riego", "Mejorar drenaje"]
                },
                "prescripcion": {
                    "diagnostico_id": "550e8400-e29b-41d4-a716-446655440000",
                    "tratamientos": [
                        {
                            "nombre": "Sulfato de hierro",
                            "dosis": "5 gramos/litro",
                            "frecuencia": "cada 3 días",
                            "duracion": "2 semanas"
                        }
                    ],
                    "urgencia": "inmediata",
                    "observaciones": "Aplicar en las primeras horas del día",
                    "seguimiento_dias": 7,
                    "contraindicaciones": ["No mezclar con calcio"]
                },
                "timestamp": "2026-04-27T10:30:00"
            }
        }
