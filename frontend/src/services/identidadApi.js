import axios from 'axios';
import { showToast } from './toastService';

const API_BASE = import.meta.env.VITE_API_URL || '/api';
const AUTH_KEYS = ['bunna_access_token', 'bunna_refresh_token', 'bunna_user'];

function clearAuthStorage() {
  AUTH_KEYS.forEach((k) => localStorage.removeItem(k));
}

const client = axios.create({
  baseURL: API_BASE,
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
});

// Interceptor para agregar token automáticamente
client.interceptors.request.use((config) => {
  const token = localStorage.getItem('bunna_access_token');
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

let isRefreshing = false;
let failedQueue = [];

const processQueue = (error, token = null) => {
  failedQueue.forEach(prom => {
    if (error) prom.reject(error);
    else prom.resolve(token);
  });
  failedQueue = [];
};

// Interceptor para refrescar token en 401
client.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // Evitar bucles: no intentar refresh si es la propia petición de refresh
    if (originalRequest?.url?.includes('/auth/refresh')) {
      return Promise.reject(error);
    }

    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        // Si ya se está refrescando, encolar la solicitud
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        }).then(token => {
          originalRequest.headers.Authorization = `Bearer ${token}`;
          return client(originalRequest);
        });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      const refreshToken = localStorage.getItem('bunna_refresh_token');
      if (!refreshToken) {
        isRefreshing = false;
        clearAuthStorage();
        window.location.href = '/login';
        return Promise.reject(error);
      }

      try {
        const res = await client.post('/api/v1/identidad/auth/refresh', { refresh_token: refreshToken });
        const data = res.data?.data || res.data;

        localStorage.setItem('bunna_access_token', data.access_token);
        localStorage.setItem('bunna_refresh_token', data.refresh_token);

        processQueue(null, data.access_token);

        originalRequest.headers.Authorization = `Bearer ${data.access_token}`;
        return client(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError, null);
        clearAuthStorage();
        window.location.href = '/login';
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    // Show toast for rate limit errors
    if (error.response?.status === 429) {
      const msg = error.response?.data?.error || 'Demasiadas solicitudes. Intente más tarde.';
      showToast(msg, 'error');
    }

    return Promise.reject(error);
  }
);

// ── Mi Perfil ──────────────────────────────────────────────
export async function getMiPerfil() {
  const res = await client.get('/api/v1/identidad/mi-perfil');
  return res.data.data;
}

export async function updateMiPerfil({ nombre, apellido }) {
  const res = await client.put('/api/v1/identidad/mi-perfil', { nombre, apellido });
  return res.data.data;
}

export async function updateMiPassword({ password_actual, nueva_password }) {
  const res = await client.put('/api/v1/identidad/mi-password', { password_actual, nueva_password });
  return res.data.data;
}

// ── Tenants ────────────────────────────────────────────────
export async function getMisTenants() {
  const res = await client.get('/api/v1/identidad/tenants/mis-tenants');
  return res.data.data;
}

export async function configurarTenant(tenantID, { nombre, slug }) {
  const res = await client.put(`/api/v1/identidad/tenants/${tenantID}`, { nombre, slug });
  return res.data.data;
}

// ── Credenciales (admin) ───────────────────────────────────
export async function getCredenciales(usuarioID) {
  const res = await client.get(`/api/v1/identidad/credenciales/${usuarioID}`);
  return res.data.data;
}

// ── Recuperación de contraseña ─────────────────────────────
export async function solicitarRecuperacion({ correo }) {
  const res = await client.post('/api/v1/identidad/recuperacion/solicitar', { correo });
  return res.data.data;
}

export async function validarTokenRecuperacion({ token }) {
  const res = await client.post('/api/v1/identidad/recuperacion/validar', { token });
  return res.data.data;
}

export async function confirmarRecuperacion({ token, nueva_password }) {
  const res = await client.post('/api/v1/identidad/recuperacion/confirmar', { token, nueva_password });
  return res.data.data;
}

// ── Verificación de correo ─────────────────────────────────
export async function solicitarVerificacion() {
  const res = await client.post('/api/v1/identidad/verificacion/solicitar');
  return res.data.data;
}

export async function confirmarVerificacion({ token }) {
  const res = await client.post('/api/v1/identidad/verificacion/confirmar', { token });
  return res.data.data;
}

export async function reenviarVerificacion({ correo }) {
  const res = await client.post('/api/v1/identidad/verificacion/reenviar', { correo });
  return res.data.data;
}

// ── Usuarios (admin) ───────────────────────────────────────
export async function getUsuarios({ pagina = 1, tamano = 20, correo = '', estado = '' } = {}) {
  const params = { pagina, tamano };
  if (correo) params.correo = correo;
  if (estado) params.estado = estado;
  const res = await client.get('/api/v1/identidad/usuarios', { params });
  return res.data.data;
}

export async function createUsuario({ correo, nombre, apellido, password }) {
  const res = await client.post('/api/v1/identidad/usuarios', { correo, nombre, apellido, password });
  return res.data.data;
}

export async function updateUsuario(usuarioID, { nombre, apellido }) {
  const res = await client.put(`/api/v1/identidad/usuarios/${usuarioID}`, { nombre, apellido });
  return res.data.data;
}

export async function bajaUsuario(usuarioID, motivo = '') {
  const res = await client.delete(`/api/v1/identidad/usuarios/${usuarioID}`, { data: { motivo } });
  return res.data.data;
}

export async function expulsarUsuario(usuarioID) {
  const res = await client.post(`/api/v1/identidad/usuarios/${usuarioID}/expulsar`);
  return res.data.data;
}

export async function unlockUsuario(usuarioID) {
  const res = await client.post(`/api/v1/identidad/usuarios/${usuarioID}/unlock`);
  return res.data.data;
}

export async function resetPasswordUsuario(usuarioID, nueva_password) {
  const res = await client.post(`/api/v1/identidad/usuarios/${usuarioID}/reset-password`, { nueva_password });
  return res.data.data;
}

// ── Roles (admin) ──────────────────────────────────────────
export async function getRoles({ pagina = 1, tamano = 20 } = {}) {
  const res = await client.get('/api/v1/identidad/roles', { params: { pagina, tamano } });
  return res.data.data;
}

export async function createRol({ nombre, descripcion, permisos = [] }) {
  const res = await client.post('/api/v1/identidad/roles', { nombre, descripcion, permisos });
  return res.data.data;
}

export async function updateRol(rolID, { nombre, descripcion }) {
  const res = await client.put(`/api/v1/identidad/roles/${rolID}`, { nombre, descripcion });
  return res.data.data;
}

export async function deleteRol(rolID) {
  const res = await client.delete(`/api/v1/identidad/roles/${rolID}`);
  return res.data.data;
}

export async function asignarRol(usuarioID, { rol_id, tenant_id = '' }) {
  const res = await client.post(`/api/v1/identidad/usuarios/${usuarioID}/roles`, { rol_id, tenant_id });
  return res.data.data;
}

export async function revocarRol(usuarioID, rolID) {
  const res = await client.delete(`/api/v1/identidad/usuarios/${usuarioID}/roles/${rolID}`);
  return res.data.data;
}

// ── Sesiones (admin) ───────────────────────────────────────
export async function getSesiones({ pagina = 1, tamano = 20 } = {}) {
  const res = await client.get('/api/v1/identidad/sesiones', { params: { pagina, tamano } });
  return res.data.data;
}

export async function cerrarSesion(sesionID) {
  const res = await client.delete(`/api/v1/identidad/sesiones/${sesionID}`);
  return res.data.data;
}

// ── IPs bloqueadas (admin) ─────────────────────────────────
export async function getIPsBloqueadas({ pagina = 1, tamano = 20 } = {}) {
  const res = await client.get('/api/v1/identidad/ips-bloqueadas', { params: { pagina, tamano } });
  return res.data.data;
}

export async function desbloquearIP(ip) {
  const res = await client.delete(`/api/v1/identidad/ips-bloqueadas/${ip}`);
  return res.data.data;
}
// ── Invitaciones ──────────────────────────────────────────────
export async function invitarUsuario({ correo, rol_id }) {
  const res = await client.post('/api/v1/identidad/invitaciones', { correo, rol_id });
  return res.data.data;
}

export async function getInvitacionInfo(token) {
  const res = await client.get(`/api/v1/identidad/invitaciones/${token}`);
  return res.data.data;
}

export async function aceptarInvitacion(token) {
  const res = await client.post(`/api/v1/identidad/invitaciones/${token}/aceptar`);
  return res.data.data;
}

export async function getInvitaciones({ pagina = 1, tamano = 20, estado = '' } = {}) {
  const params = { pagina, tamano };
  if (estado) params.estado = estado;
  try {
    const res = await client.get('/api/v1/identidad/invitaciones', { params });
    return res.data.data;
  } catch {
    return { invitaciones: [], total: 0 };
  }
}

export async function reenviarInvitacion(invitacionID) {
  try {
    const res = await client.post(`/api/v1/identidad/invitaciones/${invitacionID}/reenviar`);
    return res.data.data;
  } catch {
    throw new Error('Error al reenviar invitación');
  }
}

export async function eliminarInvitacion(invitacionID) {
  try {
    const res = await client.delete(`/api/v1/identidad/invitaciones/${invitacionID}`);
    return res.data.data;
  } catch {
    throw new Error('Error al eliminar invitación');
  }
}

export async function getRolesUsuario(usuarioID) {
  const res = await client.get(`/api/v1/identidad/usuarios/${usuarioID}/roles`);
  return res.data?.data?.roles || [];
}

// ── Permisos de roles ──────────────────────────────────────────
export async function getPermisos() {
  const res = await client.get('/api/v1/identidad/permisos');
  const data = res.data?.data || {};
  return data.permisos || data || [];
}

export async function getMisPermisos() {
  const res = await client.get('/api/v1/identidad/mis-permisos');
  return res.data?.data?.permisos || [];
}

export async function asignarPermisoARol(rolID, permisoCodigo) {
  const res = await client.post(`/api/v1/identidad/roles/${rolID}/permisos`, { permiso_codigo: permisoCodigo });
  return res.data?.data || {};
}

export async function revocarPermisoDeRol(rolID, permisoCodigo) {
  const res = await client.delete(`/api/v1/identidad/roles/${rolID}/permisos/${permisoCodigo}`);
  return res.data?.data || {};
}
