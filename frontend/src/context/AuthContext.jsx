import { createContext, useContext, useState, useEffect, useCallback } from 'react';
import { loginUsuario, registrarUsuario, parseJWT, isTokenExpired } from '../services/authApi';

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

  useEffect(() => {
    if (error) {
      const t = setTimeout(() => setError(null), 5000);
      return () => clearTimeout(t);
    }
  }, [error]);

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
        rol: 'caficultor',
        sessionId: payload?.sid,
      };

      saveSession(data.access_token, data.refresh_token, userData);
      setUser(userData);
      return { success: true };
    } catch (err) {
      const msg = err.response?.data?.detail || 'Error al iniciar sesión';
      setError(msg);
      return { success: false, error: msg };
    } finally {
      setLoading(false);
    }
  }, []);

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
    <AuthContext.Provider value={{ user, loading, error, login, register, logout, getAccessToken, setError }}>
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