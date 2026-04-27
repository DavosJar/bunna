from fastapi import APIRouter, File, UploadFile
import uuid
from datetime import datetime

from servicios.diagnostico import servicio_diagnostico
from servicios.prescripcion import servicio_prescripcion
from schemas import RespuestaDiagnosticoPrescripcion

router = APIRouter()


@router.post("/diagnostico", response_model=RespuestaDiagnosticoPrescripcion)
async def diagnostico(archivo: UploadFile = File(...)):
    """
    Endpoint completo: Recibe imagen y retorna diagnóstico + prescripción.
    """
    # Leer archivo
    contenido = await archivo.read()
    
    # Generar ID único
    diagnosis_id = str(uuid.uuid4())
    
    # Servicio 1: Diagnosticar (implementación actual: mock)
    diagnostico_resultado = servicio_diagnostico.diagnosticar(contenido)
    
    # Servicio 2: Prescribir (implementación actual: mock)
    prescripcion_resultado = servicio_prescripcion.generar_prescripcion(
        diagnostico_id=diagnosis_id,
        enfermedad=diagnostico_resultado.enfermedad,
        severidad=diagnostico_resultado.severidad
    )
    
    # Respuesta completa
    respuesta = RespuestaDiagnosticoPrescripcion(
        id=diagnosis_id,
        diagnostico=diagnostico_resultado,
        prescripcion=prescripcion_resultado,
        timestamp=datetime.now()
    )
    
    return respuesta
