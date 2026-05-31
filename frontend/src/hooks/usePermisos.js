import { useAuth } from '../context/AuthContext';

export function usePermisos() {
  const { user } = useAuth();

  const puede = (permiso) => {
    if (!user) return false;
    if (user.global) return true;
    if (!user.permisos) return false;
    return user.permisos.includes(permiso) || user.permisos.includes('*');
  };

  const esAdmin = () => puede('identidad:usuario:crear');
  const esSysAdmin = () => user?.global === true;

  return { puede, esAdmin, esSysAdmin };
}
