import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './Auth.css';

export default function RegisterPage() {
  const [nombre, setNombre] = useState('');
  const [apellido, setApellido] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const { register } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = (e) => {
    e.preventDefault();
    register({ nombre: `${nombre} ${apellido}`, email });
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
          <span className="auth-hero__tag">Únete a Bunna</span>
          <h1 className="auth-hero__title">
            Cuida tu cafetal con tecnología.
          </h1>
          <p className="auth-hero__desc">
            Diagnóstico inteligente de nitrógeno en hojas de café.
            Analiza fotos de tus plantas y recibe resultados al instante.
          </p>
        </div>
      </div>

      {/* Right Form Panel */}
      <div className="auth-form-panel">
        <div className="auth-form-container">
          <p className="auth-form__subtitle">Empieza ahora</p>
          <h2 className="auth-form__title">Crear cuenta</h2>
          <p className="auth-form__description">
            Regístrate como caficultor para diagnosticar tus plantas.
          </p>

          <form onSubmit={handleSubmit} id="register-form">
            <div className="form-row">
              <div className="form-group">
                <label className="form-label" htmlFor="reg-nombre">
                  Nombre
                </label>
                <input
                  id="reg-nombre"
                  className="form-input"
                  type="text"
                  placeholder="Juan"
                  value={nombre}
                  onChange={(e) => setNombre(e.target.value)}
                />
              </div>
              <div className="form-group">
                <label className="form-label" htmlFor="reg-apellido">
                  Apellido
                </label>
                <input
                  id="reg-apellido"
                  className="form-input"
                  type="text"
                  placeholder="Pérez"
                  value={apellido}
                  onChange={(e) => setApellido(e.target.value)}
                />
              </div>
            </div>

            <div className="form-group">
              <label className="form-label" htmlFor="reg-email">
                Correo electrónico
              </label>
              <input
                id="reg-email"
                className="form-input"
                type="email"
                placeholder="tu@correo.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                autoComplete="email"
              />
            </div>

            <div className="form-group">
              <label className="form-label" htmlFor="reg-password">
                Contraseña
              </label>
              <div className="form-input-wrapper">
                <input
                  id="reg-password"
                  className="form-input"
                  type={showPassword ? 'text' : 'password'}
                  placeholder="Mínimo 8 caracteres"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  autoComplete="new-password"
                />
                <button
                  type="button"
                  className="form-input-icon"
                  onClick={() => setShowPassword(!showPassword)}
                  aria-label={showPassword ? 'Ocultar contraseña' : 'Mostrar contraseña'}
                >
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

            <div className="form-group" style={{ marginTop: 'var(--space-2)' }}>
              <div style={{
                display: 'inline-flex',
                alignItems: 'center',
                gap: 'var(--space-2)',
                padding: '0.4rem 0.8rem',
                background: 'var(--color-green-50)',
                borderRadius: 'var(--radius-full)',
                fontSize: '0.78rem',
                fontWeight: 500,
                color: 'var(--color-green-800)',
              }}>
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                  <circle cx="8.5" cy="7" r="4"/>
                </svg>
                Rol: Caficultor
              </div>
            </div>

            <button type="submit" className="btn-primary" id="register-submit">
              Crear cuenta
              <span className="btn-primary__arrow">→</span>
            </button>
          </form>

          <p className="auth-footer">
            ¿Ya tienes cuenta?{' '}
            <Link to="/login" className="auth-footer__link">
              Iniciar sesión
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
