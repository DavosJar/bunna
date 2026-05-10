import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './Auth.css';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = (e) => {
    e.preventDefault();
    login(email);
    navigate('/dashboard');
  };

  return (
    <div className="auth-layout">
      {/* Left Hero Panel */}
      <div className="auth-hero">
        <img
          src="/coffee-bg.png"
          alt="Plantación de café"
          className="auth-hero__bg"
        />
        <div className="auth-hero__overlay" />

        <div className="auth-hero__logo">
          <div className="auth-hero__logo-icon">☕</div>
          <span className="auth-hero__logo-text">Bunna</span>
        </div>

        <div className="auth-hero__content">
          <span className="auth-hero__tag">Diagnóstico de Nitrógeno</span>
          <h1 className="auth-hero__title">
            Entra a tu finca desde cualquier dispositivo.
          </h1>
        </div>
      </div>

      {/* Right Form Panel */}
      <div className="auth-form-panel">
        <div className="auth-form-container">
          <p className="auth-form__subtitle">Bienvenida de vuelta</p>
          <h2 className="auth-form__title">Iniciar sesión</h2>
          <p className="auth-form__description">
            Accede con tu correo y contraseña.
          </p>

          <form onSubmit={handleSubmit} id="login-form">
            <div className="form-group">
              <label className="form-label" htmlFor="login-email">
                Correo
              </label>
              <div className="form-input-wrapper">
                <input
                  id="login-email"
                  className="form-input"
                  type="email"
                  placeholder="tu@correo.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  autoComplete="email"
                />
              </div>
            </div>

            <div className="form-group">
              <label className="form-label" htmlFor="login-password">
                Contraseña
                <span className="form-label__link">¿Olvidaste?</span>
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
                />
                <button
                  type="button"
                  className="form-input-icon"
                  onClick={() => setShowPassword(!showPassword)}
                  aria-label={showPassword ? 'Ocultar contraseña' : 'Mostrar contraseña'}
                >
                  {showPassword ? (
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24" />
                      <line x1="1" y1="1" x2="23" y2="23" />
                    </svg>
                  ) : (
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
                      <circle cx="12" cy="12" r="3" />
                    </svg>
                  )}
                </button>
              </div>
            </div>

            <button type="submit" className="btn-primary" id="login-submit">
              Entrar
              <span className="btn-primary__arrow">→</span>
            </button>
          </form>

          <p className="auth-footer">
            ¿No tienes cuenta?{' '}
            <Link to="/register" className="auth-footer__link">
              Crear cuenta
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
