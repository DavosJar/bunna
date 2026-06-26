import { jsPDF } from 'jspdf';
import autoTable from 'jspdf-autotable';

const MARGIN = 14;
const PAGE_W = 210;
const CONTENT_W = PAGE_W - MARGIN * 2;

function drawBox(doc, x, y, w, h) {
  doc.setDrawColor(0);
  doc.setLineWidth(0.3);
  doc.rect(x, y, w, h);
}

function sectionTitle(doc, text, x, y, w) {
  doc.setFillColor(240, 240, 240);
  doc.rect(x, y, w, 7, 'F');
  doc.setDrawColor(0);
  doc.rect(x, y, w, 7);
  doc.setFont('helvetica', 'bold');
  doc.setFontSize(9);
  doc.setTextColor(0);
  doc.text(text, x + 2, y + 5);
  return y + 7;
}

function labelValue(doc, label, value, x, y, labelW = 45) {
  doc.setFont('helvetica', 'bold');
  doc.setFontSize(8);
  doc.text(label, x + 2, y);
  doc.setFont('helvetica', 'normal');
  const lines = doc.splitTextToSize(String(value ?? '—'), CONTENT_W - labelW - 4);
  doc.text(lines, x + labelW, y);
  return y + Math.max(5, lines.length * 4);
}

/**
 * Genera un PDF técnico estilo ficha técnica con datos del lote y diagnósticos YOLO.
 */
export function generateReportePDF(reporte, { usuario } = {}) {
  const doc = new jsPDF({ unit: 'mm', format: 'a4' });
  let y = MARGIN;

  // ── Encabezado ──
  doc.setFont('helvetica', 'bold');
  doc.setFontSize(18);
  doc.text('FICHA TÉCNICA DE DIAGNÓSTICO', PAGE_W / 2, y, { align: 'center' });
  y += 8;

  doc.setFontSize(10);
  doc.setFont('helvetica', 'normal');
  doc.text('BUNNA — Análisis de Nitrógeno en Café', PAGE_W / 2, y, { align: 'center' });
  y += 5;

  doc.setFontSize(8);
  const fecha = new Date(reporte.generadoEn || Date.now()).toLocaleDateString('es-EC', {
    year: 'numeric', month: 'long', day: 'numeric',
  });
  doc.text(`Fecha: ${fecha}`, MARGIN, y);
  doc.text('1 de 1', PAGE_W - MARGIN, y, { align: 'right' });
  y += 4;
  doc.setDrawColor(0);
  doc.line(MARGIN, y, PAGE_W - MARGIN, y);
  y += 6;

  // ── Identificación ──
  y = sectionTitle(doc, 'IDENTIFICACIÓN DEL ANÁLISIS', MARGIN, y, CONTENT_W);
  const idBoxH = 28;
  drawBox(doc, MARGIN, y, CONTENT_W, idBoxH);
  let innerY = y + 5;
  innerY = labelValue(doc, 'FINCA:', reporte.finca?.nombre || '—', MARGIN, innerY);
  innerY = labelValue(doc, 'UBICACIÓN:', reporte.finca?.ubicacion || '—', MARGIN, innerY);
  innerY = labelValue(doc, 'LOTE:', `${reporte.nombre} (${reporte.areaTotal} ha)`, MARGIN, innerY);
  if (usuario) {
    labelValue(doc, 'RESPONSABLE:', usuario, MARGIN, innerY);
  }
  y += idBoxH + 5;

  // ── Especificaciones / métricas ──
  y = sectionTitle(doc, 'RESULTADOS DEL DIAGNÓSTICO', MARGIN, y, CONTENT_W);
  const m = reporte.metricas || {};
  const specBoxH = 22;
  drawBox(doc, MARGIN, y, CONTENT_W, specBoxH);
  const colW = CONTENT_W / 4;
  const metrics = [
    ['MUESTRAS', String(m.totalMuestras ?? 0)],
    ['ACEPTADOS', String(m.diagnosticosAceptados ?? 0)],
    ['PENDIENTES', String(m.diagnosticosPendientes ?? 0)],
    ['ÁREA AFECTADA', `${(m.porcentajeAfectado ?? 0).toFixed(1)}%`],
  ];
  metrics.forEach(([label, val], i) => {
    const cx = MARGIN + i * colW;
    doc.setFont('helvetica', 'bold');
    doc.setFontSize(7);
    doc.text(label, cx + colW / 2, y + 6, { align: 'center' });
    doc.setFontSize(14);
    doc.text(val, cx + colW / 2, y + 16, { align: 'center' });
    if (i > 0) {
      doc.line(cx, y, cx, y + specBoxH);
    }
  });
  y += specBoxH + 5;

  // ── Tabla de muestras ──
  y = sectionTitle(doc, 'DETALLE DE MUESTRAS Y DIAGNÓSTICOS', MARGIN, y, CONTENT_W);

  const rows = (reporte.muestras || []).map((m) => [
    m.id?.slice(0, 8) || '—',
    m.latitud != null && m.longitud != null
      ? `${Number(m.latitud).toFixed(4)}, ${Number(m.longitud).toFixed(4)}`
      : '—',
    m.createdAt ? new Date(m.createdAt).toLocaleDateString('es-EC') : '—',
    m.estadoDiagnostico || '—',
    m.tieneClorosis == null ? '—' : m.tieneClorosis ? 'Sí' : 'No',
    m.yoloLabel || '—',
    m.recomendacion ? (m.recomendacion.length > 60 ? `${m.recomendacion.slice(0, 57)}…` : m.recomendacion) : '—',
  ]);

  autoTable(doc, {
    startY: y,
    margin: { left: MARGIN, right: MARGIN },
    head: [['ID', 'Coordenadas', 'Fecha', 'Estado', 'Clorosis', 'Nitrógeno', 'Recomendación']],
    body: rows.length ? rows : [['—', '—', '—', 'Sin muestras', '—', '—', '—']],
    styles: { fontSize: 7, cellPadding: 2, lineColor: [0, 0, 0], lineWidth: 0.2 },
    headStyles: { fillColor: [34, 85, 51], textColor: 255, fontStyle: 'bold', fontSize: 7 },
    alternateRowStyles: { fillColor: [248, 250, 252] },
    theme: 'grid',
  });

  y = doc.lastAutoTable.finalY + 6;

  // ── Modo de acción / interpretación ──
  if (y > 240) {
    doc.addPage();
    y = MARGIN;
  }

  y = sectionTitle(doc, 'INTERPRETACIÓN Y RECOMENDACIONES', MARGIN, y, CONTENT_W);
  const interpH = 35;
  drawBox(doc, MARGIN, y, CONTENT_W, interpH);

  const aceptados = (reporte.muestras || []).filter((m) => m.estadoDiagnostico === 'ACEPTADO');
  const conClorosis = aceptados.filter((m) => m.tieneClorosis);
  const pct = m.porcentajeAfectado ?? 0;

  let interpretacion;
  if (!aceptados.length) {
    interpretacion = 'No hay diagnósticos aceptados. Acepte al menos un diagnóstico en la pestaña Diagnósticos para obtener recomendaciones definitivas.';
  } else if (pct >= 50) {
    interpretacion = `Se detectó deficiencia de nitrógeno (clorosis) en ${pct.toFixed(1)}% de las muestras aceptadas. Se recomienda aplicar fertilización nitrogenada según el plan agronómico del lote y monitorear nuevamente en 2-3 semanas.`;
  } else if (pct > 0) {
    interpretacion = `Deficiencia de nitrógeno presente en ${pct.toFixed(1)}% de las muestras. Monitoreo preventivo recomendado; evaluar ajuste de fertilización en zonas afectadas.`;
  } else {
    interpretacion = 'Los niveles de nitrógeno son óptimos en las muestras analizadas. Mantener el plan de fertilización actual y continuar el monitoreo periódico.';
  }

  doc.setFont('helvetica', 'normal');
  doc.setFontSize(8);
  const interpLines = doc.splitTextToSize(interpretacion, CONTENT_W - 6);
  doc.text(interpLines, MARGIN + 3, y + 6);
  y += interpH + 5;

  // ── Pie ──
  if (y > 265) {
    doc.addPage();
    y = MARGIN;
  }

  y = sectionTitle(doc, 'INFORMACIÓN DEL REPORTE', MARGIN, y, CONTENT_W);
  const footerH = 22;
  drawBox(doc, MARGIN, y, CONTENT_W, footerH);
  doc.setFontSize(7);
  doc.setFont('helvetica', 'normal');
  doc.text('Sistema: Bunna — Diagnóstico YOLO de hojas de café', MARGIN + 3, y + 6);
  doc.text('Método: Detección automática de clorosis y niveles de nitrógeno mediante visión por computadora.', MARGIN + 3, y + 11);
  doc.text('Este documento fue generado automáticamente en el navegador. Los datos se almacenan localmente.', MARGIN + 3, y + 16);
  doc.setFont('helvetica', 'bold');
  doc.text('Bunna — Investigación & Desarrollo', MARGIN + 3, y + 21);

  const filename = `ficha-tecnica-${reporte.nombre?.replace(/\s+/g, '-') || 'lote'}-${new Date().toISOString().slice(0, 10)}.pdf`;
  doc.save(filename);
  return filename;
}
