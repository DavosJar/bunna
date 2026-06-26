import { Navigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { usePermisos } from '../../hooks/usePermisos';
import { getRutaInicio } from '../../utils/roleAccess';

export function PrivateRoute({ children }) {
  const { user } = useAuth();
  return user ? children : <Navigate to="/login" replace />;
}

export function PublicRoute({ children }) {
  const { user } = useAuth();
  const { permisos } = usePermisos();
  if (user) return <Navigate to={getRutaInicio(user, permisos)} replace />;
  return children;
}

export function AdminRoute({ children }) {
  const { user } = useAuth();
  const { puedeAccederAdmin, rutaInicio } = usePermisos();
  if (!user) return <Navigate to="/login" replace />;
  if (!puedeAccederAdmin()) return <Navigate to={rutaInicio} replace />;
  return children;
}

export function TenantConfigRoute({ children }) {
  const { user } = useAuth();
  const { puedeConfigurarTenant, rutaInicio } = usePermisos();
  if (!user) return <Navigate to="/login" replace />;
  if (!puedeConfigurarTenant()) return <Navigate to={rutaInicio} replace />;
  return children;
}
