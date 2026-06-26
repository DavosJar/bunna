import { useState, useEffect } from 'react';
import { useSearchParams, Link } from 'react-router-dom';
import { confirmarVerificacion } from '../../services/identidadApi';
import LogoIcon from '../../components/LogoIcon';
import { IconCheckCircle, IconAlert } from '../../components/icons/Icons';
import './Auth.css';

export default function VerificacionCorreoPage() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token') || '';
  const [estado, setEstado] = useState('cargando'); // cargando | exito | error

  useEffect(() => {
    if (!token) { setEstado('error'); return; }
    confirmarVerificacion({ token })
      .then(() => setEstado('exito'))
      .catch(() => setEstado('error'));
  }, [token]);

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
          <span className="auth-hero__tag">Verificación</span>
          <h1 className="auth-hero__title">Confirmando tu correo electrónico.</h1>
        </div>
      </div>

      <div className="auth-form-panel">
        <div className="auth-form-container" style={{ textAlign: 'center' }}>
          {estado === 'cargando' && (
            <>
              <div className="btn-spinner" style={{ margin: '0 auto 1.5rem', width: 40, height: 40, borderWidth: 4, borderColor: 'var(--color-gray-200)', borderTopColor: 'var(--color-green-600)' }} />
              <p style={{ color: 'var(--color-gray-500)' }}>Verificando tu correo...</p>
            </>
          )}

          {estado === 'exito' && (
            <>
              <div className="auth-status-icon auth-status-icon--success"><IconCheckCircle size={48} /></div>
              <h2 className="auth-form__title">¡Correo verificado!</h2>
              <p className="auth-form__description">Tu cuenta está activa. Ya puedes usar Bunna.</p>
              <Link to="/dashboard" className="btn-primary" style={{ display: 'flex', marginTop: '2rem' }}>
                Ir al dashboard <span className="btn-primary__arrow">→</span>
              </Link>
            </>
          )}

          {estado === 'error' && (
            <>
              <div className="auth-status-icon auth-status-icon--error"><IconAlert size={48} /></div>
              <h2 className="auth-form__title">Enlace inválido</h2>
              <p className="auth-form__description">El enlace de verificación expiró o no es válido.</p>
              <Link to="/login" className="auth-footer__link" style={{ display: 'inline-block', marginTop: '1.5rem' }}>
                Volver al inicio de sesión
              </Link>
            </>
          )}
        </div>
      </div>
    </div>
  );
}