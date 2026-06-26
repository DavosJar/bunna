import { useState, useRef, useCallback, useEffect, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { useDiagnosis } from '../../context/DiagnosisContext';
import { diagnosticar, diagnosticarDesdeBase64 } from '../../services/yoloApi';
import VincularDiagnosticoPanel from '../../components/fincas/VincularDiagnosticoPanel';
import DashboardOverview from '../../components/dashboard/DashboardOverview';
import ServiceStatusBar from '../../components/ui/ServiceStatusBar';
import { IconFlask, IconLeaf, IconAlert, IconX, IconCheckCircle, IconLink } from '../../components/icons/Icons';
import { getGlobalFarmStats } from '../../utils/farmAnalytics';
import Layout from '../../components/layout/Layout';
import '../../components/fincas/VincularDiagnosticoPanel.css';
import '../../components/dashboard/DashboardOverview.css';
import '../../components/ui/StatCard.css';
import './Dashboard.css';

export default function DashboardPage() {
  const { user } = useAuth();
  const navigate = useNavigate();
  const {
    imagePreview, imageFile, results, historial, activeHistorialId,
    setImage, setResultsAndHistorial, clearDiagnosis,
    restoreFromHistorial, getActiveHistorialEntry,
  } = useDiagnosis();
  const fileInputRef = useRef(null);

  const [analyzing, setAnalyzing] = useState(false);
  const [rediagnosingId, setRediagnosingId] = useState(null);
  const [tab, setTab] = useState('resumen');
  const [dragActive, setDragActive] = useState(false);
  const [lightboxOpen, setLightboxOpen] = useState(false);
  const [zoom, setZoom] = useState(1);
  const [pan, setPan] = useState({ x: 0, y: 0 });
  const [isPanning, setIsPanning] = useState(false);
  const panStart = useRef({ x: 0, y: 0 });
  const panOffset = useRef({ x: 0, y: 0 });

  const openLightbox = () => { setZoom(1); setPan({ x: 0, y: 0 }); setLightboxOpen(true); };
  const closeLightbox = () => setLightboxOpen(false);

  const handleWheel = useCallback((e) => {
    e.preventDefault();
    setZoom((z) => Math.min(Math.max(z + (e.deltaY > 0 ? -0.15 : 0.15), 0.5), 5));
  }, []);

  const handlePointerDown = useCallback((e) => {
    if (e.button !== 0) return;
    setIsPanning(true);
    panStart.current = { x: e.clientX, y: e.clientY };
    panOffset.current = { ...pan };
    e.currentTarget.setPointerCapture(e.pointerId);
  }, [pan]);

  const handlePointerMove = useCallback((e) => {
    if (!isPanning) return;
    setPan({
      x: panOffset.current.x + (e.clientX - panStart.current.x),
      y: panOffset.current.y + (e.clientY - panStart.current.y),
    });
  }, [isPanning]);

  const handlePointerUp = useCallback(() => setIsPanning(false), []);

  const farmStats = useMemo(() => getGlobalFarmStats(user?.id), [user?.id, historial.length, results]);

  useEffect(() => {
    const onKey = (e) => { if (e.key === 'Escape') closeLightbox(); };
    if (lightboxOpen) {
      document.addEventListener('keydown', onKey);
      document.body.style.overflow = 'hidden';
    }
    return () => {
      document.removeEventListener('keydown', onKey);
      document.body.style.overflow = '';
    };
  }, [lightboxOpen]);

  const loadImage = (file) => {
    const reader = new FileReader();
    reader.onload = (ev) => setImage(file, ev.target.result);
    reader.readAsDataURL(file);
  };

  const handleFileChange = (e) => {
    const file = e.target.files?.[0];
    if (file) loadImage(file);
  };

  const handleDrag = (e) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') setDragActive(true);
    else setDragActive(false);
  };

  const handleDrop = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
    const file = e.dataTransfer.files?.[0];
    if (file && (file.type === 'image/jpeg' || file.type === 'image/png')) {
      loadImage(file);
    }
  };

  const applyYoloResult = (data, meta = {}) => {
    setResultsAndHistorial(data, {
      ...meta,
      originalPreview: meta.originalPreview || imagePreview,
      nombreArchivo: meta.nombreArchivo || imageFile?.name,
    });
    setTab('analisis');
  };

  const handleAnalyze = async () => {
    if (!imageFile) return;
    setAnalyzing(true);
    try {
      const data = await diagnosticar(imageFile);
      applyYoloResult(data, { nombreArchivo: imageFile.name, originalPreview: imagePreview });
    } catch (err) {
      const serverMsg = err.response?.data?.detail || err.response?.data?.message;
      const isNetwork = !err.response && err.message;
      const recommendation = serverMsg
        || (isNetwork ? `Error de conexión: ${err.message}` : 'No se pudo analizar la imagen. Intenta de nuevo.');
      applyYoloResult({
        feedback: { level: 'unknown', label: 'Error', percentage: 0, recommendation },
        num_detections: 0, detections: [], avg_confidence: 0, image: null,
      });
    } finally {
      setAnalyzing(false);
    }
  };

  const handleRediagnosticar = async (entry) => {
    const source = entry.originalPreview || entry.image;
    if (!source) return;
    setRediagnosingId(entry.id);
    try {
      const data = await diagnosticarDesdeBase64(source, entry.nombreArchivo || 'imagen.jpg');
      restoreFromHistorial(entry);
      applyYoloResult(data, {
        historialId: entry.id,
        nombreArchivo: entry.nombreArchivo,
        originalPreview: entry.originalPreview || source,
        vinculado: entry.vinculado,
      });
    } catch (err) {
      const serverMsg = err.response?.data?.detail || err.message;
      alert(`No se pudo re-analizar: ${serverMsg}`);
    } finally {
      setRediagnosingId(null);
    }
  };

  const handleVerHistorial = (entry) => {
    restoreFromHistorial(entry);
    setTab('analisis');
  };

  const handleRemoveImage = () => {
    clearDiagnosis();
    if (fileInputRef.current) fileInputRef.current.value = '';
  };

  const feedback = results?.feedback;
  const confidenceStr = results?.avg_confidence ? `${(results.avg_confidence * 100).toFixed(1)}%` : '-';
  const activeEntry = getActiveHistorialEntry();

  return (
    <Layout
      title={`Hola, ${user?.nombre || 'Caficultor'}`}
      subtitle="Panel de control — diagnóstico de nitrógeno con inteligencia artificial."
    >
      <ServiceStatusBar compact />

      <div className="dash-tabs">
        <button type="button" className={`dash-tabs__btn ${tab === 'resumen' ? 'dash-tabs__btn--active' : ''}`} onClick={() => setTab('resumen')}>
          Resumen
        </button>
        <button type="button" className={`dash-tabs__btn ${tab === 'analisis' ? 'dash-tabs__btn--active' : ''}`} onClick={() => setTab('analisis')}>
          Análisis
        </button>
        <button type="button" className={`dash-tabs__btn ${tab === 'historial' ? 'dash-tabs__btn--active' : ''}`} onClick={() => setTab('historial')}>
          Historial ({historial.length})
        </button>
      </div>

      {tab === 'resumen' && (
        <DashboardOverview stats={farmStats} onNavigate={setTab} />
      )}

      {tab === 'historial' ? (
        <div className="results-card">
          <h2 className="results-card__title">Historial de diagnósticos</h2>
          {historial.length === 0 ? <p className="results-empty__hint">Sin análisis previos.</p> : (
            <div style={{ display: 'grid', gap: '1rem' }}>
              {historial.map((h) => (
                <div key={h.id} className={`historial-card ${h.vinculado ? 'historial-card--vinculado' : ''}`}>
                  <p style={{ fontSize: '0.8rem', color: 'var(--color-gray-500)' }}>
                    {new Date(h.fecha).toLocaleString()} — {h.nombreArchivo}
                  </p>
                  <p><strong>{h.feedback?.label}</strong> · {h.num_detections} detecciones · {h.avg_confidence ? `${(h.avg_confidence * 100).toFixed(1)}%` : '—'}</p>
                  {h.vinculado && (
                    <p className="historial-vinculado">
                      <IconLink size={14} /> Vinculado: {h.vinculado.fincaNombre} → {h.vinculado.loteNombre}
                    </p>
                  )}
                  {h.image && <img src={h.image} alt="" style={{ maxWidth: '100%', marginTop: '0.5rem', borderRadius: '0.5rem' }} />}
                  <div className="historial-actions">
                    <button type="button" className="btn-admin btn-admin--primary" onClick={() => handleVerHistorial(h)}>Ver resultados</button>
                    <button
                      type="button"
                      className="btn-admin"
                      onClick={() => handleRediagnosticar(h)}
                      disabled={rediagnosingId === h.id || !h.originalPreview && !h.image}
                    >
                      {rediagnosingId === h.id ? 'Re-analizando...' : 'Re-analizar'}
                    </button>
                    {h.vinculado ? (
                      <button
                        type="button"
                        className="btn-admin"
                        onClick={() => navigate('/fincas', { state: { fincaId: h.vinculado.fincaId, loteId: h.vinculado.loteId } })}
                      >
                        Ver en finca
                      </button>
                    ) : (
                      <button type="button" className="btn-admin" onClick={() => { handleVerHistorial(h); setTab('analisis'); }}>
                        Vincular a finca
                      </button>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      ) : tab === 'analisis' ? (
      <div className="upload-section">
        <div className="upload-card" id="upload-card">
          <h2 className="upload-card__title">Subir imagen</h2>
          <p className="upload-card__desc">Arrastra una foto o selecciónala desde tu dispositivo.</p>

          {!imagePreview ? (
            <div
              className={`dropzone ${dragActive ? 'dropzone--active' : ''}`}
              onDragEnter={handleDrag}
              onDragLeave={handleDrag}
              onDragOver={handleDrag}
              onDrop={handleDrop}
            >
              <input ref={fileInputRef} type="file" accept="image/jpeg,image/png,.jpg,.jpeg,.png" className="dropzone__input" onChange={handleFileChange} id="file-upload-input" />
              <div className="dropzone__icon">
                <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                  <polyline points="17 8 12 3 7 8"/>
                  <line x1="12" y1="3" x2="12" y2="15"/>
                </svg>
              </div>
              <p className="dropzone__text"><span>Haz clic aquí</span> o arrastra tu imagen</p>
              <p className="dropzone__hint">PNG o JPG — máximo 10MB</p>
            </div>
          ) : (
            <div className="preview">
              <img src={imagePreview} alt="Vista previa de la hoja" className="preview__img" />
              <button className="preview__remove" onClick={handleRemoveImage} aria-label="Eliminar imagen">
                <IconX size={16} />
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
              <><div className="spinner" />Analizando...</>
            ) : (
              <>
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
                </svg>
                Analizar imagen
              </>
            )}
          </button>
        </div>

        <div className="results-card" id="results-card">
          <h2 className="results-card__title">Resultados</h2>
          <p className="results-card__desc">Diagnóstico de nitrógeno basado en IA.</p>

          {!results ? (
            <div className="results-empty">
              <div className="results-empty__icon"><IconFlask size={40} /></div>
              <p className="results-empty__text">Sin resultados aún</p>
              <p className="results-empty__hint">Sube una imagen y presiona "Analizar" para obtener el diagnóstico.</p>
            </div>
          ) : (
            <div className="results-content">
              {results.image && (
                <div className="results-image">
                  <p className="results-image__label">Imagen procesada por YOLO</p>
                  <img src={results.image} alt="Detección YOLO" className="results-image__img results-image__img--clickable" onClick={openLightbox} title="Clic para ampliar" />
                </div>
              )}
              {feedback && (
                <div className={`results-badge results-badge--${feedback.level}`}>
                  <span className="results-badge__dot" />
                  Nitrógeno {feedback.label}
                </div>
              )}
              {feedback && (
                <div className="results-gauge">
                  <div className="results-gauge__label">
                    <span>Nivel de nitrógeno</span>
                    <span className="results-gauge__value">{feedback.percentage}%</span>
                  </div>
                  <div className="results-gauge__bar">
                    <div className={`results-gauge__fill results-gauge__fill--${feedback.level}`} style={{ width: `${feedback.percentage}%` }} />
                  </div>
                </div>
              )}
              <div className="results-stats">
                <div className="results-stat">
                  <div className="results-stat__label">Confianza</div>
                  <div className="results-stat__value">{confidenceStr}</div>
                </div>
                <div className="results-stat">
                  <div className="results-stat__label">Detecciones</div>
                  <div className="results-stat__value">{results.num_detections}</div>
                </div>
              </div>
              {results.detections && results.detections.length > 0 && (
                <div className="results-detections">
                  <p className="results-detections__title">Detalle de detecciones</p>
                  <div className="results-detections__list">
                    {results.detections.map((det, i) => {
                      const esSana = det.class_name === 'hoja_sana';
                      const label = esSana
                        ? 'Hoja sana'
                        : det.class_name === 'deficiencia_nitrogeno'
                          ? 'Deficiencia de nitrógeno'
                          : det.class_name;
                      return (
                        <div key={i} className={`detection-chip detection-chip--${esSana ? 'sana' : 'deficiente'}`}>
                          <span className="detection-chip__name">
                            {esSana ? <IconLeaf size={14} /> : <IconAlert size={14} />}
                            {label}
                          </span>
                          <span className="detection-chip__conf">{(det.confidence * 100).toFixed(1)}%</span>
                        </div>
                      );
                    })}
                  </div>
                </div>
              )}
              {feedback && (
                <div className="results-rec">
                  <div className="results-rec__title">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
                    </svg>
                    Recomendación
                  </div>
                  <p className="results-rec__text">{feedback.recommendation}</p>
                </div>
              )}

              {feedback && feedback.label !== 'Error' && !activeEntry?.vinculado && (
                <VincularDiagnosticoPanel
                  yoloResults={results}
                  historialId={activeHistorialId}
                  onVinculado={(v) => navigate('/fincas', { state: { fincaId: v.fincaId, loteId: v.loteId } })}
                />
              )}
              {activeEntry?.vinculado && (
                <div className="vincular-panel" style={{ marginTop: '1rem' }}>
                  <p className="historial-vinculado" style={{ margin: 0, display: 'flex', alignItems: 'center', gap: '0.35rem' }}>
                    <IconCheckCircle size={16} /> Vinculado a {activeEntry.vinculado.fincaNombre} → {activeEntry.vinculado.loteNombre}
                  </p>
                  <button
                    type="button"
                    className="btn-admin btn-admin--primary"
                    style={{ marginTop: '0.5rem' }}
                    onClick={() => navigate('/fincas', {
                      state: { fincaId: activeEntry.vinculado.fincaId, loteId: activeEntry.vinculado.loteId },
                    })}
                  >
                    Ver en Mis Fincas
                  </button>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
      ) : null}

      {lightboxOpen && results?.image && (
        <div className="lightbox" onClick={closeLightbox}>
          <button className="lightbox__close" onClick={closeLightbox} aria-label="Cerrar">
            <IconX size={20} />
          </button>
          <div className="lightbox__zoom-info">{Math.round(zoom * 100)}%</div>
          <div
            className="lightbox__container"
            onClick={(e) => e.stopPropagation()}
            onWheel={handleWheel}
            onPointerDown={handlePointerDown}
            onPointerMove={handlePointerMove}
            onPointerUp={handlePointerUp}
            style={{ cursor: isPanning ? 'grabbing' : 'grab' }}
          >
            <img
              src={results.image}
              alt="Detección YOLO ampliada"
              className="lightbox__img"
              draggable={false}
              style={{ transform: `translate(${pan.x}px, ${pan.y}px) scale(${zoom})` }}
            />
          </div>
        </div>
      )}
    </Layout>
  );
}
