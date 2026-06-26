import axios from 'axios';

const API_BASE = import.meta.env.VITE_API_URL || '/api';

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

  const response = await authClient.post('/api/v1/auth/register', body);
  return response.data;
}

/**
 * Iniciar sesión.
 * POST /api/v1/auth/login
 */
export async function loginUsuario({ correo, password }) {
  const response = await authClient.post('/api/v1/auth/login', { correo, password });
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

/**
 * Validar formato básico de correo electrónico.
 */
export function validateEmail(email) {
  const errors = [];

  if (!email?.trim()) {
    errors.push("El correo es requerido");
    return { valid: false, errors };
  }

  const trimmed = email.trim();

  if (trimmed.length > 254)
    errors.push("El correo excede la longitud máxima permitida");

  if (!trimmed.includes("@"))
    errors.push("El correo debe contener '@'");

  if (!/^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/.test(trimmed))
    errors.push("El formato del correo no es válido");

  return {
    valid: errors.length === 0,
    errors,
  };
}

/**
 * Cambiar de tenant activo.
 * POST /api/v1/auth/switch-tenant
 */
export async function switchTenantAPI(tenantId) {
  const token = localStorage.getItem('bunna_access_token');
  const response = await authClient.post(
    '/api/v1/auth/switch-tenant',
    { tenant_id: tenantId },
    { headers: { Authorization: `Bearer ${token}` } }
  );
  return response.data;
}

/**
 * Refrescar access token usando refresh token.
 * POST /api/v1/auth/refresh
 */
export async function refreshTokenAPI(refreshToken) {
  const response = await authClient.post('/api/v1/auth/refresh', { refresh_token: refreshToken });
  return response.data;
}

/**
 * Cerrar sesión actual en el servidor.
 * POST /api/v1/auth/logout
 */
export async function logoutUsuario() {
  const token = localStorage.getItem('bunna_access_token');
  const response = await authClient.post('/api/v1/auth/logout', {}, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return response.data;
}

/**
 * Cerrar todas las sesiones del usuario.
 * POST /api/v1/auth/logout/all
 */
export async function logoutAllSesiones() {
  const token = localStorage.getItem('bunna_access_token');
  const response = await authClient.post('/api/v1/auth/logout/all', {}, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return response.data;
}
