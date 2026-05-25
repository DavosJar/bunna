import { useState } from 'react';
import { useNavigate, Link, useLocation } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import './Auth.css';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const { login, loading, error } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const mensajeBienvenida = location.state?.mensaje;

  const [correoNoVerificado, setCorreoNoVerificado] = useState(false);
  const [reenvioExitoso, setReenvioExitoso] = useState(false);
  const [reenvioError, setReenvioError] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setCorreoNoVerificado(false);
    const result = await login(email, password);
    if (result.success) {
      navigate('/dashboard');
    } else if (result.error?.includes('verificar tu correo')) {
      setCorreoNoVerificado(true);
    }
  };

  const handleReenviarVerificacion = async () => {
    setReenvioExitoso(false);
    setReenvioError(false);
    try {
      await fetch('/api/v1/verificacion/solicitar', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
      });
      setReenvioExitoso(true);
      setTimeout(() => setReenvioExitoso(false), 5000);
    } catch {
      setReenvioError(true);
      setTimeout(() => setReenvioError(false), 5000);
    }
  };

  return (
    <div className="auth-layout">
      <div className="auth-hero">
        <img src="/coffee-bg.png" alt="Plantación de café" className="auth-hero__bg" />
        <div className="auth-hero__overlay" />
        <div className="auth-hero__logo">
          <div className="auth-hero__logo-icon">☕</div>
          <span className="auth-hero__logo-text">Bunna</span>
        </div>
        <div className="auth-hero__content">
          <span className="auth-hero__tag">Diagnóstico de Nitrógeno</span>
          <h1 className="auth-hero__title">Entra a tu finca desde cualquier dispositivo.</h1>
        </div>
      </div>

      <div className="auth-form-panel">
        <div className="auth-form-container">
          <p className="auth-form__subtitle">Bienvenido de vuelta</p>
          <h2 className="auth-form__title">Iniciar sesión</h2>
          <p className="auth-form__description">Accede con tu correo y contraseña.</p>

          {mensajeBienvenida && (
            <div style={{
              display: 'flex',
              alignItems: 'flex-start',
              gap: '0.75rem',
              padding: '0.85rem 1rem',
              marginBottom: '1.5rem',
              background: '#f0fdf4',
              border: '1px solid #bbf7d0',
              borderRadius: '0.75rem',
              color: '#166534',
              fontSize: '0.85rem',
              fontWeight: 500,
              lineHeight: 1.5,
            }}>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" style={{ flexShrink: 0, marginTop: 2 }}>
                <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07A19.5 19.5 0 0 1 4.99 12a19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 3.9 1.18h3a2 2 0 0 1 2 1.72c.127.96.361 1.903.7 2.81a2 2 0 0 1-.45 2.11L8.09 8.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0 1 22 16.92z"/>
              </svg>
              {mensajeBienvenida}
            </div>
          )}

          {correoNoVerificado && (
            <div style={{
              padding: '1rem',
              marginBottom: '1rem',
              background: '#fffbeb',
              border: '1px solid #fcd34d',
              borderRadius: '0.75rem',
              color: '#92400e',
              fontSize: '0.875rem',
            }}>
              <p style={{ margin: '0 0 0.75rem', fontWeight: 600 }}>
                Debes verificar tu correo electrónico antes de iniciar sesión.
              </p>
              <p style={{ margin: '0 0 0.75rem', color: '#78350f' }}>
                Revisa tu bandeja de entrada en <strong>{email.replace(/(.{2}).*(@.*)/, '$1***$2')}</strong>
              </p>
              <button
                type="button"
                onClick={handleReenviarVerificacion}
                style={{
                  background: '#d97706',
                  color: '#fff',
                  border: 'none',
                  borderRadius: '0.5rem',
                  padding: '0.5rem 1rem',
                  fontSize: '0.875rem',
                  fontWeight: 600,
                  cursor: 'pointer',
                  width: '100%',
                }}
              >
                Reenviar email de verificación
              </button>
              {reenvioExitoso && (
                <p style={{ margin: '0.5rem 0 0', color: '#065f46', fontSize: '0.8rem', fontWeight: 500, textAlign: 'center' }}>
                  Email reenviado. Revisa tu bandeja de entrada.
                </p>
              )}
              {reenvioError && (
                <p style={{ margin: '0.5rem 0 0', color: '#991b1b', fontSize: '0.8rem', fontWeight: 500, textAlign: 'center' }}>
                  No se pudo reenviar. Intenta registrarte de nuevo.
                </p>
              )}
            </div>
          )}

          {error && !correoNoVerificado && (
            <div className="auth-error" id="login-error">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <circle cx="12" cy="12" r="10"/>
                <line x1="15" y1="9" x2="9" y2="15"/>
                <line x1="9" y1="9" x2="15" y2="15"/>
              </svg>
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} id="login-form">
            <div className="form-group">
              <label className="form-label" htmlFor="login-email">Correo</label>
              <div className="form-input-wrapper">
                <input
                  id="login-email"
                  className="form-input"
                  type="email"
                  placeholder="tu@correo.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  autoComplete="email"
                  required
                  disabled={loading}
                />
              </div>
            </div>

            <div className="form-group">
              <label className="form-label" htmlFor="login-password">
                Contraseña
                <Link to="/forgot-password" className="form-label__link">¿Olvidaste?</Link>
              </label>
              <div className="form-input-wrapper">
                <input
                  id="login-password"
                  className="form-input"
                  type={showPassword ? 'text' : 'password'}
                  placeholder="••••••••"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  autoComplete="current-password"
                  required
                  disabled={loading}
                />
                <button type="button" className="form-input-icon" onClick={() => setShowPassword(!showPassword)} aria-label={showPassword ? 'Ocultar contraseña' : 'Mostrar contraseña'}>
                  {showPassword ? (
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                      <line x1="1" y1="1" x2="23" y2="23"/>
                    </svg>
                  ) : (
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                      <circle cx="12" cy="12" r="3"/>
                    </svg>
                  )}
                </button>
              </div>
            </div>

            <button type="submit" className={`btn-primary ${loading ? 'btn-primary--loading' : ''}`} id="login-submit" disabled={loading}>
              {loading ? (
                <><div className="btn-spinner" />Entrando...</>
              ) : (
                <>Entrar <span className="btn-primary__arrow">→</span></>
              )}
            </button>
          </form>

          <p className="auth-footer">
            ¿No tienes cuenta?{' '}
            <Link to="/register" className="auth-footer__link">Crear cuenta</Link>
          </p>
        </div>
      </div>
    </div>
  );
}