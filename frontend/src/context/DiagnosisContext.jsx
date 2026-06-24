import { createContext, useContext, useState, useCallback, useEffect } from 'react';
import { useAuth } from './AuthContext';

const DiagnosisContext = createContext(null);

export function DiagnosisProvider({ children }) {
  const { user } = useAuth();
  const [imagePreview, setImagePreview] = useState(null);
  const [imageFile, setImageFile] = useState(null);
  const [results, setResults] = useState(null);

  const clearDiagnosis = useCallback(() => {
    setImagePreview(null);
    setImageFile(null);
    setResults(null);
  }, []);

  useEffect(() => {
    if (!user) clearDiagnosis();
  }, [user, clearDiagnosis]);

  const setImage = useCallback((file, preview) => {
    setImageFile(file);
    setImagePreview(preview);
    setResults(null);
  }, []);

  return (
    <DiagnosisContext.Provider value={{
      imagePreview,
      imageFile,
      results,
      setImage,
      setResults,
      clearDiagnosis,
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
