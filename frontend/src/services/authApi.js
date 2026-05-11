import axios from 'axios';

const API_BASE = import.meta.env.VITE_AUTH_API_URL || '/api';

const authClient = axios.create({
  baseURL: API_BASE,
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
});

/**
 * Registrar un nuevo usuario.
 * POST /api/v1/auth/register
 */
export async function registrarUsuario({ nombre, apellido, correo, password, telefono }) {
  const body = { nombre, apellido, correo, password };
  if (telefono) body.telefono = telefono;

  const response = await authClient.post('/v1/auth/register', body);
  return response.data;
}

/**
 * Iniciar sesión.
 * POST /api/v1/auth/login
 */
export async function loginUsuario({ correo, password }) {
  const response = await authClient.post('/v1/auth/login', { correo, password });
  return response.data;
}

/**
 * Parsear el payload de un JWT sin verificar firma (solo para UI).
 */
export function parseJWT(token) {
  try {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const payload = JSON.parse(atob(base64));
    return payload;
  } catch {
    return null;
  }
}

/**
 * Verificar si un token JWT ha expirado.
 */
export function isTokenExpired(token) {
  const payload = parseJWT(token);
  if (!payload?.exp) return true;
  return Date.now() >= payload.exp * 1000;
}
