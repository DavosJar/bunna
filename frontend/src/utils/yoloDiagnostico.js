/** Indica clorosis / deficiencia de nitrógeno según resultado YOLO */
export function tieneClorosisFromYolo(yolo) {
  if (!yolo) return null;
  const detections = yolo.detections || [];
  if (detections.some((d) => d.class_name === 'deficiencia_nitrogeno')) return true;
  if (detections.length === 0) return null;
  const level = yolo.feedback?.level;
  if (level === 'low' || level === 'medium') return true;
  if (level === 'high') return false;
  return null;
}

/** Convierte data URL base64 a File para re-enviar a YOLO */
export function base64ToFile(dataUrl, filename = 'imagen.jpg') {
  const [header, data] = dataUrl.split(',');
  const mime = header?.match(/:(.*?);/)?.[1] || 'image/jpeg';
  const binary = atob(data);
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) bytes[i] = binary.charCodeAt(i);
  return new File([bytes], filename, { type: mime });
}

/** Construye reporte local alineado con el formato del backend fincas */
export function buildReporteLocal({ finca, lote, muestras, diagnosticos }) {
  const activas = muestras.filter((m) => m.estado !== 'ELIMINADA');
  const diagMap = Object.fromEntries(diagnosticos.map((d) => [d.muestraID, d]));

  const filas = activas.map((m) => {
    const d = diagMap[m.id];
    const clorosis = d ? tieneClorosisFromYolo(d.yolo) : null;
    return {
      id: m.id,
      latitud: m.latitud,
      longitud: m.longitud,
      createdAt: m.createdAt,
      nombreArchivo: m.nombreArchivo,
      diagnosticoID: d?.diagnosticoID || d?.id || null,
      estadoDiagnostico: d?.estado || 'SIN_DIAGNOSTICO',
      tieneClorosis: clorosis,
      yoloLabel: d?.yolo?.feedback?.label,
      yoloLevel: d?.yolo?.feedback?.level,
      recomendacion: d?.yolo?.feedback?.recommendation,
      numDetections: d?.yolo?.num_detections,
      avgConfidence: d?.yolo?.avg_confidence,
      yoloImage: d?.yolo?.image,
    };
  });

  const conDiag = filas.filter((f) => f.estadoDiagnostico !== 'SIN_DIAGNOSTICO');
  const aceptados = conDiag.filter((f) => f.estadoDiagnostico === 'ACEPTADO');
  const pendientes = conDiag.filter((f) => f.estadoDiagnostico === 'PENDIENTE');
  const conClorosis = aceptados.filter((f) => f.tieneClorosis === true);
  const porcentajeAfectado = aceptados.length
    ? (conClorosis.length / aceptados.length) * 100
    : 0;

  return {
    finca: finca ? { nombre: finca.nombre, ubicacion: finca.ubicacion, descripcion: finca.descripcion } : null,
    nombre: lote.nombre,
    areaTotal: lote.area,
    descripcion: lote.descripcion,
    _local: true,
    generadoEn: new Date().toISOString(),
    metricas: {
      totalMuestras: activas.length,
      diagnosticosAceptados: aceptados.length,
      diagnosticosPendientes: pendientes.length,
      porcentajeAfectado,
    },
    muestras: filas,
  };
}
