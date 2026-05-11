import axios from 'axios';

const API_BASE = import.meta.env.VITE_YOLO_API_URL || '/yolo';

const yoloClient = axios.create({
  baseURL: API_BASE,
  timeout: 60000,
});

/**
 * Envía una imagen al backend YOLO y devuelve los resultados reales.
 * Respuesta incluye: detections, feedback, avg_confidence, image (base64)
 */
export async function diagnosticar(imagen) {
  const formData = new FormData();
  formData.append('archivo', imagen);

  const response = await yoloClient.post('/api/v1/diagnostico', formData);
  return response.data;
}
