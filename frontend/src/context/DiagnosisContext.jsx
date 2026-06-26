import { createContext, useContext, useState, useCallback, useEffect } from 'react';
import { useAuth } from './AuthContext';
import {
  loadDiagnosticoHistorial,
  saveDiagnosticoHistorial,
  updateHistorialEntry,
} from '../services/localStore';

const DiagnosisContext = createContext(null);

export function DiagnosisProvider({ children }) {
  const { user } = useAuth();
  const userId = user?.id;

  const [imagePreview, setImagePreview] = useState(null);
  const [imageFile, setImageFile] = useState(null);
  const [results, setResults] = useState(null);
  const [historial, setHistorial] = useState([]);
  const [activeHistorialId, setActiveHistorialId] = useState(null);

  const clearDiagnosis = useCallback(() => {
    setImagePreview(null);
    setImageFile(null);
    setResults(null);
    setActiveHistorialId(null);
  }, []);

  useEffect(() => {
    if (!userId) {
      clearDiagnosis();
      setHistorial([]);
      return;
    }
    setHistorial(loadDiagnosticoHistorial(userId));
  }, [userId, clearDiagnosis]);

  const setImage = useCallback((file, preview) => {
    setImageFile(file);
    setImagePreview(preview);
    setResults(null);
    setActiveHistorialId(null);
  }, []);

  const setResultsAndHistorial = useCallback((data, meta = {}) => {
    setResults(data);
    if (!userId || !data?.feedback) return;

    const entry = {
      id: meta.historialId || crypto.randomUUID(),
      feedback: data.feedback,
      num_detections: data.num_detections,
      avg_confidence: data.avg_confidence,
      detections: data.detections || [],
      image: data.image,
      originalPreview: meta.originalPreview || imagePreview,
      nombreArchivo: meta.nombreArchivo || imageFile?.name || 'imagen',
      vinculado: meta.vinculado || null,
      fecha: meta.fecha || new Date().toISOString(),
    };

    if (meta.historialId) {
      const updated = updateHistorialEntry(userId, meta.historialId, entry);
      setHistorial(updated);
      setActiveHistorialId(meta.historialId);
    } else {
      const updated = saveDiagnosticoHistorial(userId, entry);
      setHistorial(updated);
      setActiveHistorialId(updated[0]?.id || entry.id);
    }
  }, [userId, imagePreview, imageFile]);

  const restoreFromHistorial = useCallback((entry) => {
    setImagePreview(entry.originalPreview || entry.image);
    setImageFile(null);
    setResults({
      feedback: entry.feedback,
      num_detections: entry.num_detections,
      avg_confidence: entry.avg_confidence,
      detections: entry.detections || [],
      image: entry.image,
    });
    setActiveHistorialId(entry.id);
  }, []);

  const marcarHistorialVinculado = useCallback((historialId, vinculado) => {
    if (!userId) return;
    const updated = updateHistorialEntry(userId, historialId, { vinculado });
    setHistorial(updated);
  }, [userId]);

  const getActiveHistorialEntry = useCallback(() => {
    if (!activeHistorialId) return null;
    return historial.find((h) => h.id === activeHistorialId) || null;
  }, [activeHistorialId, historial]);

  return (
    <DiagnosisContext.Provider value={{
      imagePreview,
      imageFile,
      results,
      historial,
      activeHistorialId,
      setImage,
      setResults,
      setResultsAndHistorial,
      clearDiagnosis,
      restoreFromHistorial,
      marcarHistorialVinculado,
      getActiveHistorialEntry,
    }}>
      {children}
    </DiagnosisContext.Provider>
  );
}

export function useDiagnosis() {
  const context = useContext(DiagnosisContext);
  if (!context) {
    throw new Error('useDiagnosis debe usarse dentro de un DiagnosisProvider');
  }
  return context;
}
