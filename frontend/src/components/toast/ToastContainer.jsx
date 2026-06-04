import { useState, useEffect, useRef } from 'react';
import { onToast } from '../../services/toastService';
import './Toast.css';

const MAX_TOASTS = 1;

export default function ToastContainer() {
  const [toasts, setToasts] = useState([]);
  const timeoutRef = useRef(null);

  useEffect(() => {
    const unsubscribe = onToast((toast) => {
      // Clear previous timeout so it doesn't try to remove a replaced toast
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }

      setToasts(prev => {
        const next = [...prev, toast];
        return next.slice(-MAX_TOASTS);
      });

      // Auto-remove after duration
      timeoutRef.current = setTimeout(() => {
        setToasts(prev => prev.filter(t => t.id !== toast.id));
      }, toast.duration);
    });

    return () => {
      unsubscribe();
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  if (toasts.length === 0) return null;

  return (
    <div className="toast-container" aria-live="polite">
      {toasts.map(toast => (
        <div key={toast.id} className={`toast toast--${toast.type}`} role="alert">
          <span className="toast__message">{toast.message}</span>
          <button
            className="toast__close"
            onClick={() => setToasts(prev => prev.filter(t => t.id !== toast.id))}
            aria-label="Cerrar"
          >
            ✕
          </button>
        </div>
      ))}
    </div>
  );
}
