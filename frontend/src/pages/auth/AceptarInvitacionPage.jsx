import { useState, useEffect } from 'react';
import { useSearchParams, Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { aceptarInvitacion } from '../../services/identidadApi';
import LogoIcon from '../../components/LogoIcon';
import { IconCheckCircle, IconAlert } from '../../components/icons/Icons';
import './Auth.css';

const ROL_INFO = {
  caficultor:   { label: 'Caficultor',    emoji: '🌱', desc: 'Podrás gestionar fincas, lotes y análisis YOLO.',       color: '#92400e', bg: '#fef3c7' },
  agronomo:     { label: 'Agrónomo',      emoji: '🔬', desc: 'Tendrás acceso a operación de fincas y panel admin limitado.', color: '#1e40af', bg: '#dbeafe' },
  administrador:{ label: 'Administrador', emoji: '⚙️', desc: 'Podrás gestionar usuarios, roles y configuración del tenant.', color: '#166534', bg: '#dcfce7' },
};

export default function AceptarInvitacionPage() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token') || '';
  const { user, fetchMisTenants, fetchMisPermisos } = useAuth();
  const navigate = useNavigate();
  const [estado, setEstado] = useState('cargando'); // cargando | exito | error
  const [errorMsg, setErrorMsg] = useState('');
  const [rolAsignado, setRolAsignado] = useState(null);
  const [tenantID, setTenantID] = useState(null);

  useEffect(() => {
    if (!token) {
      setEstado('error');
      setErrorMsg('No se proporcionó un token de invitación válido.');
      return;
    }

    if (!user) {
      // Redirige al login manteniendo el returnUrl con el token
      const returnUrl = encodeURIComponent(`/aceptar-invitacion?token=${token}`);
      navigate(`/login?returnUrl=${returnUrl}`);
      return;
    }

    // Usuario logueado: intentar aceptar la invitación
    aceptarInvitacion({ token })
      .then(async (data) => {
        setRolAsignado(data?.rol_id || data?.rolID || null);
        setTenantID(data?.tenant_id || data?.tenantID || null);
        // Refrescar tenants y permisos para que el nuevo rol aparezca
        try {
          await Promise.all([fetchMisTenants(), fetchMisPermisos()]);
        } catch { /* no crítico */ }
        setEstado('exito');
      })
      .catch((err) => {
        setEstado('error');
        const detail = err.response?.data?.detail || '';
        if (detail.toLowerCase().includes('expirado')) {
          setErrorMsg('Este enlace de invitación ha expirado. Pide al administrador que te reenvíe la invitación.');
        } else if (detail.toLowerCase().includes('aceptada')) {
          setErrorMsg('Esta invitación ya fue aceptada anteriormente. Ya formas parte de la finca.');
        } else if (detail.toLowerCase().includes('inválido') || detail.toLowerCase().includes('invalido')) {
          setErrorMsg('El enlace de invitación no es válido. Puede haber sido revocado.');
        } else {
          setErrorMsg(detail || 'El enlace de invitación no es válido o ha expirado.');
        }
      });
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [token, user]);

  if (!user) return null;

  const rolInfo = ROL_INFO[rolAsignado] || ROL_INFO[rolAsignado?.toLowerCase()] || null;

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
            {estado === 'exito'
              ? '¡Bienvenido al equipo!'
              : estado === 'error'
                ? 'Algo salió mal'
                : 'Procesando tu invitación…'}
          </h1>
          <p className="auth-hero__desc">
            {estado === 'exito'
              ? 'Ya tienes acceso a la finca con tu rol asignado.'
              : estado === 'error'
                ? 'No pudimos procesar la invitación.'
                : 'Estamos verificando tu enlace y asignando tu rol.'}
          </p>
        </div>
      </div>

      {/* Panel derecho — estado */}
      <div className="auth-form-panel">
        <div className="auth-form-container" style={{ textAlign: 'center' }}>

          {/* ── CARGANDO ── */}
          {estado === 'cargando' && (
            <>
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
              <h2 className="auth-form__title">Verificando invitación</h2>
              <p className="auth-form__description">
                Estamos comprobando el enlace y asignando tu acceso…
              </p>
            </>
          )}

          {/* ── ÉXITO ── */}
          {estado === 'exito' && (
            <>
              {/* Ícono de éxito */}
              <div className="auth-status-icon auth-status-icon--success">
                <IconCheckCircle size={52} />
              </div>

              <h2 className="auth-form__title" style={{ marginBottom: '0.5rem' }}>
                ¡Invitación aceptada!
              </h2>
              <p className="auth-form__description" style={{ marginBottom: '1.5rem' }}>
                Ya formas parte del equipo. A continuación verás tu rol asignado.
              </p>

              {/* Tarjeta del rol asignado */}
              {rolInfo && (
                <div
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: '1rem',
                    padding: '1rem 1.25rem',
                    marginBottom: '1.5rem',
                    background: rolInfo.bg,
                    borderRadius: '0.875rem',
                    border: `1px solid ${rolInfo.color}30`,
                    textAlign: 'left',
                  }}
                >
                  <span style={{ fontSize: '2rem', lineHeight: 1 }}>{rolInfo.emoji}</span>
                  <div>
                    <p style={{ margin: '0 0 0.15rem', fontWeight: 700, color: rolInfo.color, fontSize: '1rem' }}>
                      Tu rol: {rolInfo.label}
                    </p>
                    <p style={{ margin: 0, fontSize: '0.82rem', color: rolInfo.color + 'cc' }}>
                      {rolInfo.desc}
                    </p>
                  </div>
                </div>
              )}

              {/* Si el rol no está en la lista (ID en vez de nombre), mostrar genérico */}
              {!rolInfo && rolAsignado && (
                <div
                  style={{
                    padding: '0.75rem 1rem',
                    marginBottom: '1.5rem',
                    background: '#f0fdf4',
                    borderRadius: '0.75rem',
                    border: '1px solid #bbf7d0',
                    color: '#166534',
                    fontSize: '0.875rem',
                    fontWeight: 500,
                  }}
                >
                  Rol asignado correctamente ✓
                </div>
              )}

              {/* Tip de cambio de tenant */}
              <div
                style={{
                  padding: '0.75rem 1rem',
                  marginBottom: '1.5rem',
                  background: '#f8fafc',
                  borderRadius: '0.75rem',
                  border: '1px solid #e2e8f0',
                  color: '#475569',
                  fontSize: '0.82rem',
                  textAlign: 'left',
                  lineHeight: 1.5,
                }}
              >
                <strong>💡 Tip:</strong> Si tienes acceso a varias fincas, puedes cambiar entre ellas
                desde tu <strong>Perfil</strong> o el selector en la barra superior.
              </div>

              <Link
                to="/fincas"
                className="btn-primary"
                style={{ display: 'flex', marginBottom: '0.75rem' }}
              >
                Ir a Mis Fincas <span className="btn-primary__arrow">→</span>
              </Link>
              <Link
                to="/dashboard"
                className="auth-footer__link"
                style={{ display: 'inline-block' }}
              >
                Ver Panel de Análisis
              </Link>
            </>
          )}

          {/* ── ERROR ── */}
          {estado === 'error' && (
            <>
              <div className="auth-status-icon auth-status-icon--error">
                <IconAlert size={52} />
              </div>
              <h2 className="auth-form__title" style={{ marginBottom: '0.5rem' }}>
                Invitación no válida
              </h2>
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
                {errorMsg}
              </p>

              <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
                <Link to="/fincas" className="btn-primary" style={{ display: 'flex' }}>
                  Ir a Mis Fincas <span className="btn-primary__arrow">→</span>
                </Link>
                <Link to="/dashboard" className="auth-footer__link" style={{ display: 'inline-block' }}>
                  Volver al dashboard
                </Link>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
