import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { solicitarRecuperacion } from '../../services/identidadApi';
import LogoIcon from '../../components/LogoIcon';
import { useAuth } from '../../context/AuthContext';
import { IconMail } from '../../components/icons/Icons';
import './Auth.css';

export default function ForgotPasswordPage() {
  const [correo, setCorreo] = useState('');
  const [loading, setLoading] = useState(false);
  const [enviado, setEnviado] = useState(false);
  const [error, setError] = useState('');
  const { setError: setAuthError } = useAuth();

  useEffect(() => { setAuthError(null); }, [setAuthError]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      await solicitarRecuperacion({ correo });
      setEnviado(true);
    } catch {
      setError('Ocurrió un error. Intenta de nuevo.');
    } finally {
      setLoading(false);
    }
  };

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
          <span className="auth-hero__tag">Recuperar acceso</span>
          <h1 className="auth-hero__title">Te enviamos un enlace a tu correo.</h1>
        </div>
      </div>

      <div className="auth-form-panel">
        <div className="auth-form-container">
          <p className="auth-form__subtitle">Seguridad</p>
          <h2 className="auth-form__title">¿Olvidaste tu contraseña?</h2>
          <p className="auth-form__description">
            Ingresa tu correo y te enviamos un enlace para restablecer tu contraseña.
          </p>

          {error && (
            <div className="auth-error">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <circle cx="12" cy="12" r="10" /><line x1="15" y1="9" x2="9" y2="15" /><line x1="9" y1="9" x2="15" y2="15" />
              </svg>
              {error}
            </div>
          )}

          {enviado ? (
            <div style={{ textAlign: 'center', padding: '2rem 0' }}>
              <div className="auth-status-icon auth-status-icon--success"><IconMail size={40} /></div>
              <p style={{ fontWeight: 600, color: 'var(--color-gray-800)', marginBottom: '0.5rem' }}>
                Revisa tu correo
              </p>
              <p style={{ fontSize: '0.9rem', color: 'var(--color-gray-500)' }}>
                Si el correo está registrado, recibirás el enlace en unos minutos.
              </p>
              <Link to="/login" style={{ display: 'inline-block', marginTop: '1.5rem', color: 'var(--color-green-700)', fontWeight: 600 }}>
                Volver al inicio de sesión
              </Link>
            </div>
          ) : (
            <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label className="form-label" htmlFor="forgot-correo">Correo electrónico</label>
                <input
                  id="forgot-correo"
                  className="form-input"
                  type="email"
                  placeholder="tu@correo.com"
                  value={correo}
                  onChange={(e) => setCorreo(e.target.value)}
                  required
                  disabled={loading}
                />
              </div>

              <button
                type="submit"
                className={`btn-primary ${loading ? 'btn-primary--loading' : ''}`}
                disabled={loading}
              >
                {loading ? <><div className="btn-spinner" />Enviando...</> : <>Enviar enlace <span className="btn-primary__arrow">→</span></>}
              </button>
            </form>
          )}

          <p className="auth-footer">
            <Link to="/login" className="auth-footer__link">← Volver al inicio de sesión</Link>
          </p>
        </div>
      </div>
    </div>
  );
}