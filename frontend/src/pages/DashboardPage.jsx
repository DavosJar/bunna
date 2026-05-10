import { useState, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './Dashboard.css';

/* Mock results data */
const MOCK_RESULTS = [
  {
    level: 'high',
    label: 'Alto',
    percentage: 82,
    confidence: '94.2%',
    detections: 12,
    recommendation:
      'El nivel de nitrógeno es óptimo. Mantén el plan de fertilización actual y realiza un nuevo análisis en 30 días para monitorear la evolución.',
  },
  {
    level: 'medium',
    label: 'Medio',
    percentage: 54,
    confidence: '87.5%',
    detections: 8,
    recommendation:
      'Se recomienda aplicar urea (46-0-0) a razón de 150g por planta. Repite el análisis en 15 días para evaluar la respuesta del cultivo.',
  },
  {
    level: 'low',
    label: 'Bajo',
    percentage: 28,
    confidence: '91.8%',
    detections: 15,
    recommendation:
      'Nivel crítico de nitrógeno detectado. Aplica fertilizante nitrogenado de forma urgente (200g/planta de urea). Consulta con un agrónomo para un plan de recuperación.',
  },
];

export default function DashboardPage() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const fileInputRef = useRef(null);

  const [imagePreview, setImagePreview] = useState(null);
  const [analyzing, setAnalyzing] = useState(false);
  const [results, setResults] = useState(null);
  const [dragActive, setDragActive] = useState(false);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const handleFileChange = (e) => {
    const file = e.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (ev) => {
        setImagePreview(ev.target.result);
        setResults(null);
      };
      reader.readAsDataURL(file);
    }
  };

  const handleDrag = (e) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true);
    } else {
      setDragActive(false);
    }
  };

  const handleDrop = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
    const file = e.dataTransfer.files?.[0];
    if (file && file.type.startsWith('image/')) {
      const reader = new FileReader();
      reader.onload = (ev) => {
        setImagePreview(ev.target.result);
        setResults(null);
      };
      reader.readAsDataURL(file);
    }
  };

  const handleAnalyze = () => {
    setAnalyzing(true);
    // Simulate API call with random result
    setTimeout(() => {
      const randomResult = MOCK_RESULTS[Math.floor(Math.random() * MOCK_RESULTS.length)];
      setResults(randomResult);
      setAnalyzing(false);
    }, 2000);
  };

  const handleRemoveImage = () => {
    setImagePreview(null);
    setResults(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const initials = user?.nombre
    ? user.nombre.split(' ').map((w) => w[0]).join('').toUpperCase().slice(0, 2)
    : '??';

  return (
    <div className="dashboard">
      {/* Navbar */}
      <nav className="navbar" id="dashboard-navbar">
        <div className="navbar__brand">
          <div className="navbar__logo">☕</div>
          <span className="navbar__name">Bunna</span>
        </div>
        <div className="navbar__right">
          <div className="navbar__user">
            <div className="navbar__avatar">{initials}</div>
            <div className="navbar__user-info">
              <span className="navbar__user-name">{user?.nombre || 'Usuario'}</span>
              <span className="navbar__user-role">Caficultor</span>
            </div>
          </div>
          <button
            className="navbar__logout"
            onClick={handleLogout}
            id="logout-btn"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
              <polyline points="16 17 21 12 16 7"/>
              <line x1="21" y1="12" x2="9" y2="12"/>
            </svg>
            Salir
          </button>
        </div>
      </nav>

      {/* Main Content */}
      <main className="dashboard__main">
        <div className="dashboard__header">
          <h1 className="dashboard__greeting">
            Hola, {user?.nombre || 'Caficultor'} 👋
          </h1>
          <p className="dashboard__greeting-sub">
            Sube una foto de tus hojas de café para analizar el nivel de nitrógeno.
          </p>
        </div>

        <div className="upload-section">
          {/* Upload Card */}
          <div className="upload-card" id="upload-card">
            <h2 className="upload-card__title">Subir imagen</h2>
            <p className="upload-card__desc">
              Arrastra una foto o selecciónala desde tu dispositivo.
            </p>

            {!imagePreview ? (
              <div
                className={`dropzone ${dragActive ? 'dropzone--active' : ''}`}
                onDragEnter={handleDrag}
                onDragLeave={handleDrag}
                onDragOver={handleDrag}
                onDrop={handleDrop}
              >
                <input
                  ref={fileInputRef}
                  type="file"
                  accept="image/*"
                  className="dropzone__input"
                  onChange={handleFileChange}
                  id="file-upload-input"
                />
                <div className="dropzone__icon">
                  <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                    <polyline points="17 8 12 3 7 8"/>
                    <line x1="12" y1="3" x2="12" y2="15"/>
                  </svg>
                </div>
                <p className="dropzone__text">
                  <span>Haz clic aquí</span> o arrastra tu imagen
                </p>
                <p className="dropzone__hint">
                  PNG, JPG o WEBP — máximo 10MB
                </p>
              </div>
            ) : (
              <div className="preview">
                <img
                  src={imagePreview}
                  alt="Vista previa de la hoja"
                  className="preview__img"
                />
                <button
                  className="preview__remove"
                  onClick={handleRemoveImage}
                  aria-label="Eliminar imagen"
                >
                  ✕
                </button>
              </div>
            )}

            <button
              className={`btn-analyze ${analyzing ? 'btn-analyze--loading' : ''}`}
              onClick={handleAnalyze}
              disabled={!imagePreview || analyzing}
              id="analyze-btn"
            >
              {analyzing ? (
                <>
                  <div className="spinner" />
                  Analizando...
                </>
              ) : (
                <>
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <circle cx="11" cy="11" r="8"/>
                    <line x1="21" y1="21" x2="16.65" y2="16.65"/>
                  </svg>
                  Analizar imagen
                </>
              )}
            </button>
          </div>

          {/* Results Card */}
          <div className="results-card" id="results-card">
            <h2 className="results-card__title">Resultados</h2>
            <p className="results-card__desc">
              Diagnóstico de nitrógeno basado en IA.
            </p>

            {!results ? (
              <div className="results-empty">
                <div className="results-empty__icon">🔬</div>
                <p className="results-empty__text">Sin resultados aún</p>
                <p className="results-empty__hint">
                  Sube una imagen y presiona "Analizar" para obtener el diagnóstico.
                </p>
              </div>
            ) : (
              <div className="results-content">
                {/* Badge */}
                <div className={`results-badge results-badge--${results.level}`}>
                  <span className="results-badge__dot" />
                  Nitrógeno {results.label}
                </div>

                {/* Gauge */}
                <div className="results-gauge">
                  <div className="results-gauge__label">
                    <span>Nivel de nitrógeno</span>
                    <span className="results-gauge__value">{results.percentage}%</span>
                  </div>
                  <div className="results-gauge__bar">
                    <div
                      className={`results-gauge__fill results-gauge__fill--${results.level}`}
                      style={{ width: `${results.percentage}%` }}
                    />
                  </div>
                </div>

                {/* Stats */}
                <div className="results-stats">
                  <div className="results-stat">
                    <div className="results-stat__label">Confianza</div>
                    <div className="results-stat__value">{results.confidence}</div>
                  </div>
                  <div className="results-stat">
                    <div className="results-stat__label">Detecciones</div>
                    <div className="results-stat__value">{results.detections}</div>
                  </div>
                </div>

                {/* Recommendation */}
                <div className="results-rec">
                  <div className="results-rec__title">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
                    </svg>
                    Recomendación
                  </div>
                  <p className="results-rec__text">{results.recommendation}</p>
                </div>
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  );
}
