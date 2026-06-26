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

export async function desactivarFinca(fincaID, { confirmar = true } = {}) {
  const res = await client.post(`/fincas/${fincaID}/desactivar`, { confirmar });
  return res.data.data;
}

// ── Lotes ──────────────────────────────────────────────────
export async function agregarLote(fincaID, { nombre, area, descripcion }) {
  const res = await client.post(`/fincas/${fincaID}/lotes`, { nombre, area, descripcion });
  return res.data.data;
}

export async function eliminarLote(loteID) {
  const res = await client.post(`/lotes/${loteID}/eliminar`);
  return res.data.data;
}

// ── Muestras ───────────────────────────────────────────────
export async function tomarMuestra(fincaID, loteID, { latitud, longitud }) {
  const res = await client.post(`/fincas/${fincaID}/lotes/${loteID}/muestras`, { latitud, longitud });
  return res.data.data;
}

export async function listarMuestras(fincaID, loteID) {
  const res = await client.get(`/fincas/${fincaID}/lotes/${loteID}/muestras`);
  return res.data.data || [];
}

// ── Diagnósticos ───────────────────────────────────────────
export async function solicitarDiagnosticoManual(muestraID, { imageURL }) {
  const res = await client.post(`/muestras/${muestraID}/diagnosticos/manual`, { imageURL });
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
