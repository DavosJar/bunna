import axios from 'axios';

const API_BASE = import.meta.env.VITE_FINCAS_API_URL || '/fincas';

const client = axios.create({
  baseURL: API_BASE,
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
});

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('bunna_access_token');
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

/** Servicio fincas caído: sin respuesta, o proxy Vite devuelve 502 (ECONNREFUSED en :8082) */
export function isServicioFincasNoDisponible(err) {
  if (!err?.response) return true;
  const status = err.response.status;
  return status === 502 || status === 503 || status === 504;
}

export function parseErrorFincas(err) {
  if (isServicioFincasNoDisponible(err)) {
    return 'El servicio de fincas no está disponible en el puerto 8082.';
  }
  const data = err.response?.data;
  if (typeof data === 'string') return data;
  return data?.error || data?.detalle || data?.detail || `Error ${err.response?.status || ''}`.trim();
}

export async function verificarSaludFincas() {
  const res = await client.get('/health');
  return res.data;
}

export async function fincasApiDisponible() {
  try {
    await verificarSaludFincas();
    return true;
  } catch {
    return false;
  }
}

// ── Fincas ─────────────────────────────────────────────────
export async function registrarFinca({ nombre, ubicacion, descripcion }) {
  const res = await client.post('/fincas', { nombre, ubicacion, descripcion });
  return res.data.data;
}

export async function listarFincas() {
  const res = await client.get('/fincas');
  return res.data.data || [];
}

export async function desactivarFinca(fincaID, { confirmar = true } = {}) {
  const res = await client.post(`/fincas/${fincaID}/desactivar`, { confirmar });
  return res.data.data;
}

// ── Lotes ──────────────────────────────────────────────────
export async function agregarLote(fincaID, { nombre, area, descripcion }) {
  const res = await client.post(`/fincas/${fincaID}/lotes`, { nombre, area, descripcion });
  return res.data.data;
}

export async function listarLotes(fincaID) {
  const res = await client.get(`/fincas/${fincaID}/lotes`);
  return res.data.data || [];
}

export async function eliminarLote(loteID) {
  const res = await client.post(`/lotes/${loteID}/eliminar`);
  return res.data.data;
}

// ── Muestras ───────────────────────────────────────────────
export async function tomarMuestra(fincaID, loteID, { latitud, longitud }) {
  const url = loteID ? `/fincas/${fincaID}/lotes/${loteID}/muestras` : `/fincas/${fincaID}/muestras`;
  const res = await client.post(url, { latitud, longitud });
  return res.data.data;
}

export async function listarMuestras(fincaID, loteID) {
  const url = loteID ? `/fincas/${fincaID}/lotes/${loteID}/muestras` : `/fincas/${fincaID}/muestras`;
  const res = await client.get(url);
  return res.data.data || [];
}

// ── Diagnósticos ───────────────────────────────────────────
export async function solicitarDiagnosticoManual(muestraID, { imageURL }) {
  const res = await client.post(`/muestras/${muestraID}/diagnosticos/manual`, { imageURL });
  return res.data.data;
}

export async function guardarResultadoManual(muestraID, payload) {
  const res = await client.post(`/muestras/${muestraID}/diagnosticos/manual/resultado`, payload);
  return res.data.data;
}

export async function aceptarDiagnostico(diagnosticoID) {
  const res = await client.post(`/diagnosticos/${diagnosticoID}/aceptar`);
  return res.data.data;
}

export async function rechazarDiagnostico(diagnosticoID, { motivo = '' } = {}) {
  const res = await client.post(`/diagnosticos/${diagnosticoID}/rechazar`, { motivo });
  return res.data.data;
}

// ── Reportes ───────────────────────────────────────────────
export async function generarReporteLote(fincaID, loteID) {
  const res = await client.get(`/fincas/${fincaID}/lotes/${loteID}/reporte`);
  return res.data.data;
}

// ── Nodos (Cámaras) ────────────────────────────────────────
export async function registrarNodo(fincaID, { node_key, nombre, lote_id }) {
  const payload = { finca_id: fincaID, node_key, nombre };
  if (lote_id) payload.lote_id = lote_id;
  const res = await client.post('/nodos', payload);
  return res.data.data;
}

export async function listarNodos(fincaID) {
  // El endpoint GET /nodos es paginado y devuelve un objeto { data, total, ... } en res.data.data
  const res = await client.get(`/nodos`);
  const data = res.data?.data?.data || [];
  return data.filter((n) => n.finca_id === fincaID);
}

export async function desactivarNodo(nodoID) {
  const res = await client.post(`/nodos/${nodoID}/desactivar`, { estado: 'INACTIVO' });
  return res.data.data;
}
