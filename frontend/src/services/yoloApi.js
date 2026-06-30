import axios from 'axios';

// En dev usa /yolo (proxy Vite → API remota). En prod usa la URL directa de la VM por defecto.
const API_BASE = import.meta.env.VITE_YOLO_API_URL || (import.meta.env.PROD ? 'https://i2u6hsbhf1.execute-api.us-east-2.amazonaws.com' : '/yolo');

const yoloClient = axios.create({
  baseURL: API_BASE,
  timeout: 30000,
});

/**
 * Verifica que la API YOLO esté activa.
 * GET /health
 */
export async function verificarSaludYOLO() {
  try {
    // Como el API Gateway no rutea /health a YOLO, le hacemos un GET a su ruta principal.
    // YOLO espera un POST, así que devolverá 405 Method Not Allowed, lo que confirma que está vivo.
    await yoloClient.get('/api/v1/diagnostico');
    return { status: "ok" };
  } catch (error) {
    if (error.response && error.response.status === 405) {
      return { status: "ok" }; // 405 significa que YOLO respondió
    }
    throw error;
  }
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

/**
 * Re-ejecuta diagnóstico sobre un archivo ya subido al servidor YOLO.
 * GET /api/v1/diagnostico/{image_name}
 */
export async function diagnosticarPorNombre(imageName) {
  const response = await yoloClient.get(`/api/v1/diagnostico/${encodeURIComponent(imageName)}`);
  return response.data;
}

/** Re-analiza desde data URL base64 (re-sube la imagen vía POST) */
export async function diagnosticarDesdeBase64(dataUrl, filename = 'imagen.jpg') {
  const res = await fetch(dataUrl);
  const blob = await res.blob();
  const file = new File([blob], filename, { type: blob.type || 'image/jpeg' });
  return diagnosticar(file);
}
