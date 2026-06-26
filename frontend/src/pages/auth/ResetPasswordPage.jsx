import { useState, useEffect } from 'react';
import { useSearchParams, useNavigate, Link } from 'react-router-dom';
import { validarTokenRecuperacion, confirmarRecuperacion } from '../../services/identidadApi';
import LogoIcon from '../../components/LogoIcon';
import { validarPassword } from '../../services/validacionPassword';
import { IconCheckCircle, IconAlert } from '../../components/icons/Icons';
import './Auth.css';

export default function ResetPasswordPage() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const token = searchParams.get('token') || '';

  const [tokenValido, setTokenValido] = useState(null);
  const [password, setPassword] = useState('');
  const [confirmar, setConfirmar] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [exito, setExito] = useState(false);

  useEffect(() => {
    if (!token) { setTokenValido(false); return; }
    validarTokenRecuperacion({ token })
      .then((d) => setTokenValido(d.valido))
      .catch(() => setTokenValido(false));
  }, [token]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (password !== confirmar) { setError('Las contraseñas no coinciden.'); return; }
    const validacion = validarPassword(password);
    if (!validacion.valida) {
      setError(validacion.errores.join('. '));
      return;
    }
    setLoading(true);
    setError('');
    try {
      await confirmarRecuperacion({ token, nueva_password: password });
      setExito(true);
      setTimeout(() => navigate('/login'), 3000);
    } catch {
      setError('El enlace expiró o es inválido. Solicita uno nuevo.');
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
          <span className="auth-hero__tag">Nueva contraseña</span>
          <h1 className="auth-hero__title">Elige una contraseña segura.</h1>
        </div>
      </div>

      <div className="auth-form-panel">
        <div className="auth-form-container">
          <p className="auth-form__subtitle">Seguridad</p>
          <h2 className="auth-form__title">Restablecer contraseña</h2>

          {tokenValido === null && (
            <p style={{ color: 'var(--color-gray-500)' }}>Validando enlace...</p>
          )}

          {tokenValido === false && (
            <div style={{ textAlign: 'center', padding: '2rem 0' }}>
              <div className="auth-status-icon auth-status-icon--error"><IconAlert size={40} /></div>
              <p style={{ fontWeight: 600, color: 'var(--color-gray-800)', marginBottom: '0.5rem' }}>
                Enlace inválido o expirado
              </p>
              <Link to="/forgot-password" className="auth-footer__link">
                Solicitar nuevo enlace
              </Link>
            </div>
          )}

          {tokenValido === true && !exito && (
            <form onSubmit={handleSubmit}>
              <p className="auth-form__description">Ingresa tu nueva contraseña.</p>

              {error && (
                <div className="auth-error">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <circle cx="12" cy="12" r="10" /><line x1="15" y1="9" x2="9" y2="15" /><line x1="9" y1="9" x2="15" y2="15" />
                  </svg>
                  {error}
                </div>
              )}

              <div className="form-group">
                <label className="form-label" htmlFor="nueva-password">Nueva contraseña</label>
                <div className="form-input-wrapper">
                  <input
                    id="nueva-password"
                    className="form-input"
                    type={showPassword ? 'text' : 'password'}
                    placeholder="Mínimo 8 caracteres"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    disabled={loading}
                  />
                  <button type="button" className="form-input-icon" onClick={() => setShowPassword(!showPassword)}>
                    {showPassword
                      ? <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24" /><line x1="1" y1="1" x2="23" y2="23" /></svg>
                      : <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" /><circle cx="12" cy="12" r="3" /></svg>
                    }
                  </button>
                </div>
              </div>

              <div className="form-group">
                <label className="form-label" htmlFor="confirmar-password">Confirmar contraseña</label>
                <input
                  id="confirmar-password"
                  className="form-input"
                  type={showPassword ? 'text' : 'password'}
                  placeholder="Repite la contraseña"
                  value={confirmar}
                  onChange={(e) => setConfirmar(e.target.value)}
                  required
                  disabled={loading}
                />
              </div>

              <button type="submit" className={`btn-primary ${loading ? 'btn-primary--loading' : ''}`} disabled={loading}>
                {loading ? <><div className="btn-spinner" />Guardando...</> : <>Guardar contraseña <span className="btn-primary__arrow">→</span></>}
              </button>
            </form>
          )}

          {exito && (
            <div style={{ textAlign: 'center', padding: '2rem 0' }}>
              <div className="auth-status-icon auth-status-icon--success"><IconCheckCircle size={40} /></div>
              <p style={{ fontWeight: 600, color: 'var(--color-gray-800)', marginBottom: '0.5rem' }}>
                ¡Contraseña actualizada!
              </p>
              <p style={{ fontSize: '0.9rem', color: 'var(--color-gray-500)' }}>
                Redirigiendo al inicio de sesión...
              </p>
            </div>
          )}

          <p className="auth-footer">
            <Link to="/login" className="auth-footer__link">← Volver al inicio de sesión</Link>
          </p>
        </div>
      </div>
    </div>
  );
}