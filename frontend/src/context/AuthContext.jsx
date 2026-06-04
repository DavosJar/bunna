import { createContext, useContext, useState, useEffect, useCallback, useMemo } from 'react';
import { loginUsuario, registrarUsuario, parseJWT, isTokenExpired, switchTenantAPI } from '../services/authApi';
import { getMisTenants, getMisPermisos } from '../services/identidadApi';

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

export function AuthProvider({ children }) {
  const [user, setUser] = useState(loadStoredUser);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [availableTenants, setAvailableTenants] = useState([]);
  const [ownTenantId, setOwnTenantId] = useState(null);
  const [permisos, setPermisos] = useState([]);

  // Derived: current tenant object
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

  // Cargar tenants del usuario al montar si ya está autenticado
  const fetchMisTenants = useCallback(async () => {
    try {
      const data = await getMisTenants();
      setAvailableTenants(data.tenants || []);
      setOwnTenantId(data.propio_id || null);
    } catch {
      // Silently fail — the user might not be fully authenticated yet
      setAvailableTenants([]);
      setOwnTenantId(null);
    }
  }, []);

  const fetchMisPermisos = useCallback(async () => {
    try {
      const data = await getMisPermisos();
      setPermisos(data);
    } catch {
      setPermisos([]);
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
      const payload = parseJWT(data.access_token);

      const userData = {
        id: data.usuario_id,
        email: correo,
        nombre: correo.split('@')[0],
        tenantID: data.tenant_id,
        ownTenantID: data.tenant_id,
        rol: data.rol || 'caficultor',
        sessionId: payload?.sid,
      };

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
  }, [fetchMisTenants]);

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
  }, [user, fetchMisTenants]);

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

  const logout = useCallback(() => {
    clearStorage();
    setUser(null);
    setError(null);
  }, []);

  const getAccessToken = useCallback(() => {
    const token = localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN);
    if (!token || isTokenExpired(token)) {
      logout();
      return null;
    }
    return token;
  }, [logout]);

  return (
    <AuthContext.Provider value={{ user, permisos, loading, error, login, register, logout, getAccessToken, setError, availableTenants, ownTenantId, currentTenant, switchTenant, fetchMisTenants, fetchMisPermisos }}>
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