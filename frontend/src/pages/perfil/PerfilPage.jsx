import { useState, useEffect, useMemo } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { getMiPerfil, updateMiPerfil, updateMiPassword, reenviarVerificacion } from '../../services/identidadApi';
import { validarPassword } from '../../services/validacionPassword';
import { getGlobalFarmStats } from '../../utils/farmAnalytics';
import StatCard from '../../components/ui/StatCard';
import { usePermisos } from '../../hooks/usePermisos';
import { IconFarm, IconSample, IconCheck, IconClock } from '../../components/icons/Icons';
import RoleAccessCard from '../../components/perfil/RoleAccessCard';
import Layout from '../../components/layout/Layout';
import '../../components/ui/StatCard.css';
import './Perfil.css';

export default function PerfilPage() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const { puedeConfigurarTenant, rolLabel } = usePermisos();
  const farmStats = useMemo(() => getGlobalFarmStats(user?.id), [user?.id]);

  const [perfil, setPerfil] = useState(null);
  const [nombre, setNombre] = useState('');
  const [apellido, setApellido] = useState('');
  const [passwordActual, setPasswordActual] = useState('');
  const [nuevaPassword, setNuevaPassword] = useState('');
  const [loadingPerfil, setLoadingPerfil] = useState(false);
  const [loadingPassword, setLoadingPassword] = useState(false);
  const [msgPerfil, setMsgPerfil] = useState(null);
  const [msgPassword, setMsgPassword] = useState(null);
  const [msgSesion, setMsgSesion] = useState(null);
  const [msgVerificacion, setMsgVerificacion] = useState(null);

  useEffect(() => {
    getMiPerfil().then((d) => {
      setPerfil(d);
      setNombre(d.nombre);
      setApellido(d.apellido);
    }).catch(() => {});
  }, []);

  const handleUpdatePerfil = async (e) => {
    e.preventDefault();
    setLoadingPerfil(true);
    setMsgPerfil(null);
    try {
      await updateMiPerfil({ nombre, apellido });
      setMsgPerfil({ tipo: 'exito', texto: 'Perfil actualizado correctamente.' });
    } catch {
      setMsgPerfil({ tipo: 'error', texto: 'Error al actualizar el perfil.' });
    } finally {
      setLoadingPerfil(false);
    }
  };

  const handleUpdatePassword = async (e) => {
    e.preventDefault();
    const validacion = validarPassword(nuevaPassword);
    if (!validacion.valida) {
      setMsgPassword({ tipo: 'error', texto: validacion.errores.join('. ') });
      return;
    }
    setLoadingPassword(true);
    setMsgPassword(null);
    try {
      await updateMiPassword({ password_actual: passwordActual, nueva_password: nuevaPassword });
      setMsgPassword({ tipo: 'exito', texto: 'Contraseña actualizada correctamente.' });
      setPasswordActual('');
      setNuevaPassword('');
    } catch {
      setMsgPassword({ tipo: 'error', texto: 'Contraseña actual incorrecta.' });
    } finally {
      setLoadingPassword(false);
    }
  };

  const initials = perfil
    ? `${perfil.nombre?.[0] || ''}${perfil.apellido?.[0] || ''}`.toUpperCase()
    : '??';

  return (
    <Layout title="Mi Perfil" subtitle={`Cuenta de ${rolLabel} — gestiona tu información personal.`}>
      <div className="stat-grid" style={{ marginBottom: '1.5rem' }}>
        <StatCard icon={<IconFarm />} label="Mis fincas" value={farmStats.fincas} accent="green" />
        <StatCard icon={<IconSample />} label="Análisis realizados" value={farmStats.muestras + farmStats.diagnosticos} accent="blue" />
        <StatCard icon={<IconCheck />} label="Aceptados" value={farmStats.aceptados} accent="green" />
        <StatCard icon={<IconClock />} label="Pendientes" value={farmStats.pendientes} accent="amber" />
      </div>

      <RoleAccessCard />

      {/* Info personal */}
      <div className="perfil-card">
        <h2 className="perfil-card__title">Información personal</h2>
        <div className="perfil-avatar">
          <div className="perfil-avatar__circle">{initials}</div>
          <div className="perfil-avatar__info">
            <p className="perfil-avatar__name">{perfil ? `${perfil.nombre} ${perfil.apellido}` : '...'}</p>
            <p>{perfil?.correo}</p>
            <p>
              Estado: {perfil?.estado}
              {perfil?.correo_verificado != null && (
                <span className={`perfil-verif-badge ${perfil.correo_verificado ? 'perfil-verif-badge--ok' : 'perfil-verif-badge--pending'}`}>
                  {perfil.correo_verificado ? 'Correo verificado' : 'Correo pendiente'}
                </span>
              )}
            </p>
          </div>
        </div>

        {msgPerfil && (
          <div className={msgPerfil.tipo === 'exito' ? 'perfil-success' : 'perfil-error'}>
            {msgPerfil.texto}
          </div>
        )}

        <form onSubmit={handleUpdatePerfil}>
          <div className="form-row">
            <div className="form-group">
              <label className="form-label" htmlFor="p-nombre">Nombre</label>
              <input id="p-nombre" className="form-input" type="text" value={nombre} onChange={(e) => setNombre(e.target.value)} required disabled={loadingPerfil} />
            </div>
            <div className="form-group">
              <label className="form-label" htmlFor="p-apellido">Apellido</label>
              <input id="p-apellido" className="form-input" type="text" value={apellido} onChange={(e) => setApellido(e.target.value)} required disabled={loadingPerfil} />
            </div>
          </div>
          <button type="submit" className="btn-perfil btn-perfil--primary" disabled={loadingPerfil}>
            {loadingPerfil ? 'Guardando...' : 'Guardar cambios'}
          </button>
        </form>
      </div>

      {/* Cambiar contraseña */}
      <div className="perfil-card">
        <h2 className="perfil-card__title">Cambiar contraseña</h2>

        {msgPassword && (
          <div className={msgPassword.tipo === 'exito' ? 'perfil-success' : 'perfil-error'}>
            {msgPassword.texto}
          </div>
        )}

        <form onSubmit={handleUpdatePassword}>
          <div className="form-group">
            <label className="form-label" htmlFor="p-actual">Contraseña actual</label>
            <input id="p-actual" className="form-input" type="password" placeholder="••••••••" value={passwordActual} onChange={(e) => setPasswordActual(e.target.value)} required disabled={loadingPassword} />
          </div>
          <div className="form-group">
            <label className="form-label" htmlFor="p-nueva">Nueva contraseña</label>
            <input id="p-nueva" className="form-input" type="password" placeholder="Mínimo 8 caracteres" value={nuevaPassword} onChange={(e) => setNuevaPassword(e.target.value)} required disabled={loadingPassword} />
          </div>
          <button type="submit" className="btn-perfil btn-perfil--primary" disabled={loadingPassword}>
            {loadingPassword ? 'Actualizando...' : 'Actualizar contraseña'}
          </button>
        </form>
      </div>

      {puedeConfigurarTenant() && (
      <div className="perfil-card">
        <h2 className="perfil-card__title">Organización / Tenant</h2>
        <p style={{ marginBottom: '1rem', color: 'var(--color-gray-600)', fontSize: '0.9rem' }}>
          Configura el nombre y slug de tu organización.
        </p>
        <Link to="/finca-config" className="btn-perfil btn-perfil--primary" style={{ display: 'inline-block', textAlign: 'center', textDecoration: 'none' }}>
          Configurar tenant
        </Link>
      </div>
      )}

      <div className="perfil-card">
        <h2 className="perfil-card__title">Verificación de correo</h2>
        {msgVerificacion && <div className={msgVerificacion.tipo === 'exito' ? 'perfil-success' : 'perfil-error'}>{msgVerificacion.texto}</div>}
        <button type="button" className="btn-perfil btn-perfil--primary" onClick={async () => {
          try {
            await reenviarVerificacion({ correo: user?.correo });
            setMsgVerificacion({ tipo: 'exito', texto: 'Correo de verificación reenviado.' });
          } catch {
            setMsgVerificacion({ tipo: 'error', texto: 'No se pudo reenviar el correo.' });
          }
        }}>Reenviar verificación</button>
      </div>

      <div className="perfil-card">
        <h2 className="perfil-card__title">Sesiones</h2>
        {msgSesion && <div className={msgSesion.tipo === 'exito' ? 'perfil-success' : 'perfil-error'}>{msgSesion.texto}</div>}
        <button type="button" className="btn-perfil btn-perfil--primary" onClick={async () => {
          await logout({ allSessions: true });
          navigate('/login');
        }}>Cerrar todas las sesiones</button>
      </div>
    </Layout>
  );
}