import { useState, useEffect } from 'react';
import { useSearchParams, Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { getInvitacionInfo, aceptarInvitacion } from '../../services/identidadApi';
import LogoIcon from '../../components/LogoIcon';
import { IconCheckCircle, IconAlert } from '../../components/icons/Icons';
import './Auth.css';

export default function AceptarInvitacionPage() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token') || '';
  const { user, fetchMisTenants, fetchMisPermisos } = useAuth();
  const navigate = useNavigate();

  const [cargando, setCargando] = useState(true);
  const [invitacion, setInvitacion] = useState(null);
  const [errorCarga, setErrorCarga] = useState('');

  const [aceptando, setAceptando] = useState(false);
  const [estado, setEstado] = useState(null); // 'exito' | 'error'
  const [msgError, setMsgError] = useState('');

  // 1. Cargar info de la invitación (endpoint público)
  useEffect(() => {
    if (!token) {
      setErrorCarga('No se proporcionó un token de invitación.');
      setCargando(false);
      return;
    }
    getInvitacionInfo(token)
      .then(data => {
        setInvitacion(data);
        setCargando(false);
      })
      .catch(() => {
        setErrorCarga('El enlace de invitación no es válido o ha expirado.');
        setCargando(false);
      });
  }, [token]);

  // 2. Aceptar invitación (solo con clic)
  const handleAceptar = async () => {
    setAceptando(true);
    setEstado(null);
    setMsgError('');
    try {
      await aceptarInvitacion(token);
      await Promise.all([fetchMisTenants(), fetchMisPermisos()]);
      setEstado('exito');
    } catch (err) {
      const detail = err.response?.data?.detail || '';
      setEstado('error');
      setMsgError(detail || 'No se pudo aceptar la invitación.');
    } finally {
      setAceptando(false);
    }
  };

  const emailCoincide = user && invitacion && user.email?.toLowerCase() === invitacion.email?.toLowerCase();
  const returnUrl = encodeURIComponent(`/aceptar-invitacion?token=${token}`);

  if (cargando) {
    return (
      <div className="auth-layout">
        <div className="auth-hero">
          <img src="/coffee-bg.png" alt="Café" className="auth-hero__bg" />
          <div className="auth-hero__overlay" />
          <div className="auth-hero__logo">
            <LogoIcon className="auth-hero__logo-icon" imgClassName="auth-hero__logo-img" />
            <span className="auth-hero__logo-text">Bunna</span>
          </div>
          <div className="auth-hero__content">
            <span className="auth-hero__tag">Invitación a finca</span>
            <h1 className="auth-hero__title">Tienes una invitación</h1>
            <p className="auth-hero__desc">Estamos verificando tu enlace de invitación…</p>
          </div>
        </div>
        <div className="auth-form-panel">
          <div className="auth-form-container" style={{ textAlign: 'center' }}>
            <div
              className="btn-spinner"
              style={{
                margin: '0 auto 1.5rem',
                width: 48,
                height: 48,
                borderWidth: 4,
                borderColor: 'var(--color-gray-200)',
                borderTopColor: 'var(--color-green-600)',
              }}
            />
            <h2 className="auth-form__title">Cargando invitación</h2>
            <p className="auth-form__description">
              Estamos obteniendo los detalles de tu invitación…
            </p>
          </div>
        </div>
      </div>
    );
  }

  if (errorCarga) {
    return (
      <div className="auth-layout">
        <div className="auth-hero">
          <img src="/coffee-bg.png" alt="Café" className="auth-hero__bg" />
          <div className="auth-hero__overlay" />
          <div className="auth-hero__logo">
            <LogoIcon className="auth-hero__logo-icon" imgClassName="auth-hero__logo-img" />
            <span className="auth-hero__logo-text">Bunna</span>
          </div>
          <div className="auth-hero__content">
            <span className="auth-hero__tag">Invitación a finca</span>
            <h1 className="auth-hero__title">Algo salió mal</h1>
            <p className="auth-hero__desc">No pudimos procesar la invitación.</p>
          </div>
        </div>
        <div className="auth-form-panel">
          <div className="auth-form-container" style={{ textAlign: 'center' }}>
            <div className="auth-status-icon auth-status-icon--error">
              <IconAlert size={52} />
            </div>
            <h2 className="auth-form__title">Invitación no válida</h2>
            <p
              className="auth-form__description"
              style={{
                marginBottom: '1.5rem',
                padding: '0.85rem 1rem',
                background: '#fef2f2',
                borderRadius: '0.75rem',
                border: '1px solid #fecaca',
                color: '#991b1b',
                textAlign: 'left',
                fontSize: '0.875rem',
                lineHeight: 1.5,
              }}
            >
              {errorCarga}
            </p>
            <Link to="/fincas" className="btn-primary" style={{ display: 'flex' }}>
              Ir a Mis Fincas →
            </Link>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="auth-layout">
      {/* Panel izquierdo — hero */}
      <div className="auth-hero">
        <img src="/coffee-bg.png" alt="Café" className="auth-hero__bg" />
        <div className="auth-hero__overlay" />
        <div className="auth-hero__logo">
          <LogoIcon className="auth-hero__logo-icon" imgClassName="auth-hero__logo-img" />
          <span className="auth-hero__logo-text">Bunna</span>
        </div>
        <div className="auth-hero__content">
          <span className="auth-hero__tag">Invitación a finca</span>
          <h1 className="auth-hero__title">
            {estado === 'exito' ? '¡Bienvenido al equipo!' : 'Tienes una invitación'}
          </h1>
          <p className="auth-hero__desc">
            {estado === 'exito'
              ? 'Ya tienes acceso a la finca con tu rol asignado.'
              : 'Acepta la invitación para unirte a la finca.'}
          </p>
        </div>
      </div>

      {/* Panel derecho */}
      <div className="auth-form-panel">
        <div className="auth-form-container" style={{ textAlign: 'center' }}>
          {estado === 'exito' ? (
            // Mensaje de éxito
            <>
              <div className="auth-status-icon auth-status-icon--success">
                <IconCheckCircle size={52} />
              </div>
              <h2 className="auth-form__title">¡Invitación aceptada!</h2>
              <p className="auth-form__description">
                Ya formas parte del equipo con el rol asignado.
              </p>
              <Link to="/fincas" className="btn-primary" style={{ display: 'flex' }}>
                Ir a Mis Fincas →
              </Link>
            </>
          ) : estado === 'error' ? (
            // Mensaje de error con acción
            <>
              <div className="auth-status-icon auth-status-icon--error">
                <IconAlert size={52} />
              </div>
              <h2 className="auth-form__title">No se pudo aceptar</h2>
              <p style={{ color: '#991b1b', marginBottom: '1.5rem' }}>{msgError}</p>
              {msgError.includes('crear una cuenta') && (
                <Link to={`/register?returnUrl=${returnUrl}`} className="btn-primary" style={{ display: 'flex' }}>
                  Crear cuenta →
                </Link>
              )}
              {msgError.includes('expirado') && (
                <Link to="/fincas" className="btn-primary" style={{ display: 'flex' }}>
                  Ir a Mis Fincas →
                </Link>
              )}
              {(msgError.includes('ya aceptaste') || msgError.includes('Ya formas parte')) && (
                <Link to="/fincas" className="btn-primary" style={{ display: 'flex' }}>
                  Ir a Mis Fincas →
                </Link>
              )}
            </>
          ) : invitacion ? (
            // Mostrar detalles de la invitación
            <>
              <div className="auth-status-icon auth-status-icon--info" style={{ background: '#f0fdf4' }}>
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="#16a34a" strokeWidth="2">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/>
                </svg>
              </div>

              <h2 className="auth-form__title">Invitación a finca</h2>

              <div style={{
                background: '#f9fafb', borderRadius: '0.75rem', padding: '1rem',
                marginBottom: '1.5rem', textAlign: 'left', fontSize: '0.9rem',
              }}>
                <p style={{ margin: '0 0 0.5rem' }}>
                  <strong>Tenant:</strong> {invitacion.tenant_nombre || invitacion.tenant_id}
                </p>
                <p style={{ margin: '0 0 0.5rem' }}>
                  <strong>Rol:</strong> {invitacion.rol_nombre || invitacion.rol_id}
                </p>
                <p style={{ margin: 0, color: '#6b7280', fontSize: '0.82rem' }}>
                  <strong>Email:</strong> {invitacion.email}
                </p>
              </div>

              {user && emailCoincide && (
                <button
                  className="btn-primary"
                  onClick={handleAceptar}
                  disabled={aceptando}
                  style={{ width: '100%' }}
                >
                  {aceptando ? 'Aceptando…' : 'Aceptar invitación'}
                </button>
              )}

              {user && !emailCoincide && (
                <div style={{
                  padding: '1rem', background: '#fef2f2', borderRadius: '0.75rem',
                  color: '#991b1b', fontSize: '0.85rem',
                }}>
                  Esta invitación es para <strong>{invitacion.email}</strong>, pero tienes
                  sesión con <strong>{user.email}</strong>.
                  Cierra sesión e inicia con la cuenta correcta.
                </div>
              )}

              {!user && (
                <>
                  <Link
                    to={`/register?returnUrl=${returnUrl}`}
                    className="btn-primary"
                    style={{ display: 'flex', marginBottom: '0.75rem' }}
                  >
                    Crear cuenta para aceptar →
                  </Link>
                  <Link
                    to={`/login?returnUrl=${returnUrl}`}
                    className="auth-footer__link"
                    style={{ display: 'inline-block' }}
                  >
                    Ya tengo cuenta — Iniciar sesión
                  </Link>
                </>
              )}
            </>
          ) : null}
        </div>
      </div>
    </div>
  );
}
