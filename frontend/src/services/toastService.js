// Toast service — allows any module (including axios interceptors) to trigger toasts
// without depending on React context. React components subscribe via onToast().

let handler = null;

/**
 * Show a toast notification.
 * @param {string} message - The message to display
 * @param {'error'|'success'|'info'} type - Toast type
 * @param {number} duration - Auto-dismiss duration in ms (default 10000)
 */
export function showToast(message, type = 'error', duration = 10000) {
  if (handler) {
    handler({ id: Date.now() + Math.random(), message, type, duration });
  }
}

/**
 * Register a handler for toast events (called by ToastContainer on mount).
 * Returns an unsubscribe function.
 * @param {(toast: object) => void} fn
 * @returns {() => void}
 */
export function onToast(fn) {
  handler = fn;
  return () => {
    handler = null;
  };
}
