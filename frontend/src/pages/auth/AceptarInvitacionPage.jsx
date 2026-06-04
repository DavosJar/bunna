import { useState, useEffect } from 'react';
import { useSearchParams, Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { aceptarInvitacion } from '../../services/identidadApi';
import './Auth.css';

export default function AceptarInvitacionPage() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token') || '';
  const { user } = useAuth();
  const navigate = useNavigate();
  const [estado, setEstado] = useState('cargando');
  const [errorMsg, setErrorMsg] = useState('');

  useEffect(() => {
    if (!token) {
      setEstado('error');
      setErrorMsg('No se proporcionó un token de invitación.');
      return;
    }
    if (!user) {
      navigate(`/login?returnUrl=${encodeURIComponent(`/aceptar-invitacion?token=${token}`)}`);
      return;
    }
    aceptarInvitacion({ token })
      .then(() => setEstado('exito'))
      .catch((err) => {
        setEstado('error');
        const detail = err.response?.data?.detail || '';
        if (detail.includes('expirado') || detail.includes('expirado')) {
          setErrorMsg('Este enlace de invitación ha expirado.');
        } else if (detail.includes('aceptada')) {
          setErrorMsg('Esta invitación ya fue aceptada.');
        } else {
          setErrorMsg(detail || 'El enlace de invitación no es válido.');
        }
      });
  }, [token, user, navigate]);

  if (!user) return null;

  return (
    <div className="auth-layout">
      <div className="auth-hero">
        <img src="/coffee-bg.png" alt="Café" className="auth-hero__bg" />
        <div className="auth-hero__overlay" />
        <div className="auth-hero__logo">
          <div className="auth-hero__logo-icon">☕</div>
          <span className="auth-hero__logo-text">Bunna</span>
        </div>
        <div className="auth-hero__content">
          <span className="auth-hero__tag">Invitación</span>
          <h1 className="auth-hero__title">Aceptando tu invitación.</h1>
        </div>
      </div>

      <div className="auth-form-panel">
        <div className="auth-form-container" style={{ textAlign: 'center' }}>
          {estado === 'cargando' && (
            <>
              <div className="btn-spinner" style={{ margin: '0 auto 1.5rem', width: 40, height: 40, borderWidth: 4, borderColor: 'var(--color-gray-200)', borderTopColor: 'var(--color-green-600)' }} />
              <p style={{ color: 'var(--color-gray-500)' }}>Aceptando invitación...</p>
            </>
          )}

          {estado === 'exito' && (
            <>
              <div style={{ fontSize: '4rem', marginBottom: '1rem' }}>🎉</div>
              <h2 className="auth-form__title">¡Invitación aceptada!</h2>
              <p className="auth-form__description">Ya formas parte del equipo. Puedes cambiar de finca desde tu perfil.</p>
              <Link to="/dashboard" className="btn-primary" style={{ display: 'flex', marginTop: '2rem' }}>
                Ir al dashboard <span className="btn-primary__arrow">→</span>
              </Link>
            </>
          )}

          {estado === 'error' && (
            <>
              <div style={{ fontSize: '4rem', marginBottom: '1rem' }}>❌</div>
              <h2 className="auth-form__title">Enlace inválido</h2>
              <p className="auth-form__description">{errorMsg}</p>
              <Link to="/dashboard" className="auth-footer__link" style={{ display: 'inline-block', marginTop: '1.5rem' }}>
                Volver al dashboard
              </Link>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
