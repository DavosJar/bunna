import { createContext, useContext, useState, useEffect, useCallback, useMemo } from 'react';
import { loginUsuario, registrarUsuario, parseJWT, isTokenExpired, switchTenantAPI, logoutUsuario, logoutAllSesiones } from '../services/authApi';
import { getMisTenants, getMisPermisos, getMiPerfil } from '../services/identidadApi';

const AuthContext = createContext(null);

const STORAGE_KEYS = {
  ACCESS_TOKEN: 'bunna_access_token',
  REFRESH_TOKEN: 'bunna_refresh_token',
  USER: 'bunna_user',
};

function loadStoredUser() {
  try {
    const accessToken = localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN);
    const userJson = localStorage.getItem(STORAGE_KEYS.USER);
    if (!accessToken || !userJson) return null;
    if (isTokenExpired(accessToken)) {
      clearStorage();
      return null;
    }
    return JSON.parse(userJson);
  } catch {
    clearStorage();
    return null;
  }
}

function clearStorage() {
  localStorage.removeItem(STORAGE_KEYS.ACCESS_TOKEN);
  localStorage.removeItem(STORAGE_KEYS.REFRESH_TOKEN);
  localStorage.removeItem(STORAGE_KEYS.USER);
}

function saveSession(accessToken, refreshToken, userData) {
  localStorage.setItem(STORAGE_KEYS.ACCESS_TOKEN, accessToken);
  localStorage.setItem(STORAGE_KEYS.REFRESH_TOKEN, refreshToken);
  localStorage.setItem(STORAGE_KEYS.USER, JSON.stringify(userData));
}

async function buildUserFromLogin(data, correo) {
  const payload = parseJWT(data.access_token);
  let nombre = correo.split('@')[0];
  let apellido = '';
  try {
    const perfil = await getMiPerfil();
    nombre = perfil.nombre || nombre;
    apellido = perfil.apellido || '';
  } catch { /* usar fallback */ }

  return {
    id: data.usuario_id,
    email: correo,
    nombre,
    apellido,
    tenantID: data.tenant_id,
    ownTenantID: data.tenant_id,
    rol: data.rol || 'caficultor',
    sessionId: payload?.sid,
  };
}

export function AuthProvider({ children }) {
  const [user, setUser] = useState(loadStoredUser);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [availableTenants, setAvailableTenants] = useState([]);
  const [ownTenantId, setOwnTenantId] = useState(null);
  const [permisos, setPermisos] = useState([]);
  const [degraded, setDegraded] = useState(false);

  const currentTenant = useMemo(
    () => availableTenants?.find(t => t.id === user?.tenantID) || null,
    [availableTenants, user?.tenantID]
  );

  useEffect(() => {
    if (error) {
      const t = setTimeout(() => setError(null), 5000);
      return () => clearTimeout(t);
    }
  }, [error]);

  const fetchMisTenants = useCallback(async () => {
    try {
      const data = await getMisTenants();
      setDegraded(false);
      setAvailableTenants(data.tenants || []);
      const propioId = data.propio_id || null;
      setOwnTenantId(propioId);
      if (propioId) {
        setUser((prev) => (prev && prev.ownTenantID !== propioId ? { ...prev, ownTenantID: propioId } : prev));
      }
    } catch {
      setDegraded(true);
      // No vaciar availableTenants — getMisTenants() ya intentó cache
    }
  }, []);

  const fetchMisPermisos = useCallback(async () => {
    try {
      const data = await getMisPermisos();
      setDegraded(false);
      setPermisos(data.map(p => p.codigo));
    } catch {
      setDegraded(true);
      // No vaciar permisos — getMisPermisos() ya intentó cache
    }
  }, []);

  useEffect(() => {
    if (user) {
      fetchMisTenants();
      fetchMisPermisos();
    }
  }, [user, fetchMisTenants, fetchMisPermisos]);

  const login = useCallback(async (correo, password) => {
    setLoading(true);
    setError(null);
    try {
      const res = await loginUsuario({ correo, password });
      const data = res.data;
      const userData = await buildUserFromLogin(data, correo);

      saveSession(data.access_token, data.refresh_token, userData);
      setUser(userData);
      await fetchMisTenants();
      await fetchMisPermisos();
      return { success: true, user: userData };
    } catch (err) {
      const msg = err.response?.data?.detail || 'Error al iniciar sesión';
      setError(msg);
      return { success: false, error: msg };
    } finally {
      setLoading(false);
    }
  }, [fetchMisTenants, fetchMisPermisos]);

  const switchTenant = useCallback(async (tenantId) => {
    setLoading(true);
    setError(null);
    try {
      const res = await switchTenantAPI(tenantId);
      const data = res.data;
      const payload = parseJWT(data.access_token);

      const updatedUser = {
        ...user,
        tenantID: data.tenant_id,
        rol: data.rol,
        sessionId: payload?.sid,
      };

      saveSession(data.access_token, data.refresh_token, updatedUser);
      setUser(updatedUser);
      await fetchMisTenants();
      await fetchMisPermisos();
      return { success: true };
    } catch (err) {
      const msg = err.response?.data?.detail || 'Error al cambiar de tenant';
      setError(msg);
      return { success: false, error: msg };
    } finally {
      setLoading(false);
    }
  }, [user, fetchMisTenants, fetchMisPermisos]);

  const register = useCallback(async ({ nombre, apellido, correo, password, telefono }) => {
    setLoading(true);
    setError(null);
    try {
      await registrarUsuario({ nombre, apellido, correo, password, telefono });
      return { success: true, correo };
    } catch (err) {
      const detail = err.response?.data?.detail || '';
      let msg = 'Error al registrarse';
      if (detail.includes('duplicate') || detail.includes('23505')) {
        msg = 'Este correo ya está registrado';
      } else if (detail) {
        msg = detail;
      }
      setError(msg);
      return { success: false, error: msg };
    } finally {
      setLoading(false);
    }
  }, []);

  const logout = useCallback(async ({ allSessions = false } = {}) => {
    const userId = user?.id;
    try {
      if (allSessions) await logoutAllSesiones();
      else await logoutUsuario();
    } catch { /* limpiar local aunque falle el servidor */ }
    clearStorage();
    setUser(null);
    setError(null);
    setAvailableTenants([]);
    setPermisos([]);
    setDegraded(false);
  }, [user]);

  const getAccessToken = useCallback(() => {
    const token = localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN);
    if (!token || isTokenExpired(token)) {
      clearStorage();
      setUser(null);
      return null;
    }
    return token;
  }, []);

  return (
    <AuthContext.Provider value={{
      user, permisos, loading, error, login, register, logout, getAccessToken,
      setError, availableTenants, ownTenantId, currentTenant, switchTenant,
      fetchMisTenants, fetchMisPermisos, degraded,
    }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth debe usarse dentro de un AuthProvider');
  }
  return context;
}
