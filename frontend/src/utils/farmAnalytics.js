import {
  loadFincasLocal,
  loadLotesLocal,
  loadMuestrasLocal,
  loadDiagnosticosLocal,
  loadDiagnosticoHistorial,
} from '../services/localStore';
import { tieneClorosisFromYolo } from './yoloDiagnostico';

/** Estadísticas globales de todas las fincas del usuario */
export function getGlobalFarmStats(userId) {
  if (!userId) {
    return emptyStats();
  }

  const fincas = loadFincasLocal(userId).filter((f) => f.estado !== 'PENDIENTE_ELIMINACION');
  let totalLotes = 0;
  let totalMuestras = 0;
  let totalDiagnosticos = 0;
  let pendientes = 0;
  let aceptados = 0;
  let rechazados = 0;
  let conClorosis = 0;
  const nitrogenCounts = {};
  const timeline = [];
  const recientes = [];
  const fincasDetalle = [];

  fincas.forEach((finca) => {
    const lotes = loadLotesLocal(userId, finca.id).filter((l) => l.estado !== 'ELIMINADO');
    let fincaMuestras = 0;
    let fincaPendientes = 0;

    lotes.forEach((lote) => {
      const muestras = loadMuestrasLocal(userId, finca.id, lote.id);
      const diagnosticos = loadDiagnosticosLocal(userId, finca.id, lote.id);
      totalMuestras += muestras.length;
      totalDiagnosticos += diagnosticos.length;
      fincaMuestras += muestras.length;

      diagnosticos.forEach((d) => {
        if (d.estado === 'PENDIENTE') { pendientes++; fincaPendientes++; }
        if (d.estado === 'ACEPTADO') {
          aceptados++;
          if (tieneClorosisFromYolo(d.yolo)) conClorosis++;
        }
        if (d.estado === 'RECHAZADO') rechazados++;

        const label = d.yolo?.feedback?.label;
        if (label) nitrogenCounts[label] = (nitrogenCounts[label] || 0) + 1;

        const fecha = d.createdAt?.slice(0, 10);
        if (fecha) {
          timeline.push({ fecha, label: label || 'Sin etiqueta', estado: d.estado });
        }

        recientes.push({
          id: d.id,
          tipo: 'diagnostico',
          fecha: d.createdAt,
          titulo: `Nitrógeno ${label || '—'}`,
          subtitulo: `${finca.nombre} → ${lote.nombre}`,
          estado: d.estado,
          fincaId: finca.id,
          loteId: lote.id,
        });
      });
    });

    totalLotes += lotes.length;
    fincasDetalle.push({
      ...finca,
      lotes: lotes.length,
      muestras: fincaMuestras,
      pendientes: fincaPendientes,
    });
  });

  loadDiagnosticoHistorial(userId).forEach((h) => {
    recientes.push({
      id: h.id,
      tipo: 'historial',
      fecha: h.fecha,
      titulo: `Nitrógeno ${h.feedback?.label || '—'}`,
      subtitulo: h.vinculado
        ? `${h.vinculado.fincaNombre} → ${h.vinculado.loteNombre}`
        : (h.nombreArchivo || 'Análisis rápido'),
      estado: h.vinculado ? 'VINCULADO' : 'LOCAL',
    });
  });

  recientes.sort((a, b) => new Date(b.fecha) - new Date(a.fecha));

  const porFecha = {};
  timeline.forEach((t) => {
    porFecha[t.fecha] = (porFecha[t.fecha] || 0) + 1;
  });
  const serieTemporal = Object.entries(porFecha)
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([fecha, total]) => ({ fecha, total }));

  return {
    fincas: fincas.length,
    lotes: totalLotes,
    muestras: totalMuestras,
    diagnosticos: totalDiagnosticos,
    pendientes,
    aceptados,
    rechazados,
    porcentajeAfectado: aceptados ? (conClorosis / aceptados) * 100 : 0,
    nitrogenData: Object.entries(nitrogenCounts).map(([name, value]) => ({ name, value })),
    estadoData: [
      { name: 'Aceptados', value: aceptados, color: '#22c55e' },
      { name: 'Pendientes', value: pendientes, color: '#eab308' },
      { name: 'Rechazados', value: rechazados, color: '#ef4444' },
    ].filter((d) => d.value > 0),
    serieTemporal,
    recientes: recientes.slice(0, 10),
    fincasDetalle,
  };
}

function emptyStats() {
  return {
    fincas: 0, lotes: 0, muestras: 0, diagnosticos: 0,
    pendientes: 0, aceptados: 0, rechazados: 0, porcentajeAfectado: 0,
    nitrogenData: [], estadoData: [], serieTemporal: [], recientes: [], fincasDetalle: [],
  };
}

/** Puntos GPS para mapa de muestras de un lote */
export function getMuestrasMapPoints(reporte) {
  if (!reporte?.muestras) return [];
  return reporte.muestras
    .filter((m) => m.latitud != null && m.longitud != null)
    .map((m) => ({
      id: m.id,
      lat: Number(m.latitud),
      lng: Number(m.longitud),
      label: m.yoloLabel || '—',
      estado: m.estadoDiagnostico,
      clorosis: m.tieneClorosis,
      fecha: m.createdAt,
    }));
}

/** Datos de gráficos para el reporte de un lote */
export function getReportChartData(reporte) {
  if (!reporte) return { nitrogenData: [], estadoData: [], serieTemporal: [] };

  const nitrogenCounts = {};
  const timeline = [];
  let aceptados = 0;
  let pendientes = 0;
  let rechazados = 0;

  (reporte.muestras || []).forEach((m) => {
    if (m.yoloLabel) nitrogenCounts[m.yoloLabel] = (nitrogenCounts[m.yoloLabel] || 0) + 1;
    if (m.estadoDiagnostico === 'ACEPTADO') aceptados++;
    if (m.estadoDiagnostico === 'PENDIENTE') pendientes++;
    if (m.estadoDiagnostico === 'RECHAZADO') rechazados++;
    if (m.createdAt) timeline.push({ fecha: m.createdAt.slice(0, 10) });
  });

  const porFecha = {};
  timeline.forEach((t) => { porFecha[t.fecha] = (porFecha[t.fecha] || 0) + 1; });

  return {
    nitrogenData: Object.entries(nitrogenCounts).map(([name, value]) => ({ name, value })),
    estadoData: [
      { name: 'Aceptados', value: aceptados, color: '#22c55e' },
      { name: 'Pendientes', value: pendientes, color: '#eab308' },
      { name: 'Rechazados', value: rechazados, color: '#ef4444' },
    ].filter((d) => d.value > 0),
    serieTemporal: Object.entries(porFecha)
      .sort(([a], [b]) => a.localeCompare(b))
      .map(([fecha, total]) => ({ fecha, total })),
  };
}
