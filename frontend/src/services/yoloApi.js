import axios from 'axios';

// En dev usa /yolo (proxy Vite → API remota). En prod puede ser la URL directa.
const API_BASE = import.meta.env.VITE_YOLO_API_URL || '/yolo';

const yoloClient = axios.create({
  baseURL: API_BASE,
  timeout: 30000,
});

/**
 * Verifica que la API YOLO esté activa.
 * GET /health
 */
export async function verificarSaludYOLO() {
  const response = await yoloClient.get('/health');
  return response.data;
}

/**
 * Envía una imagen al backend YOLO y devuelve los resultados de diagnóstico.
 * POST /api/v1/diagnostico — campo multipart: archivo
 */
export async function diagnosticar(imagen) {
  const formData = new FormData();
  formData.append('archivo', imagen);

  // No setear Content-Type manualmente: axios debe incluir el boundary del multipart
  const response = await yoloClient.post('/api/v1/diagnostico', formData);
  return response.data;
}
