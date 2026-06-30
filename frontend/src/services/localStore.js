import { generateUUID } from '../utils/uuid';

const PREFIX = 'bunna_local_';

function key(userId, suffix) {
  return `${PREFIX}${userId}_${suffix}`;
}

export function loadFincasLocal(userId) {
  try {
    const raw = localStorage.getItem(key(userId, 'fincas'));
    return raw ? JSON.parse(raw) : [];
  } catch {
    return [];
  }
}

export function saveFincaLocal(userId, finca) {
  const list = loadFincasLocal(userId);
  const idx = list.findIndex((f) => f.id === finca.id);
  if (idx >= 0) list[idx] = { ...list[idx], ...finca };
  else list.push(finca);
  localStorage.setItem(key(userId, 'fincas'), JSON.stringify(list));
  return list;
}

export function updateFincaLocal(userId, fincaId, patch) {
  const list = loadFincasLocal(userId).map((f) => (f.id === fincaId ? { ...f, ...patch } : f));
  localStorage.setItem(key(userId, 'fincas'), JSON.stringify(list));
  return list;
}

export function loadLotesLocal(userId, fincaId) {
  try {
    const raw = localStorage.getItem(key(userId, `lotes_${fincaId}`));
    return raw ? JSON.parse(raw) : [];
  } catch {
    return [];
  }
}

export function saveLoteLocal(userId, fincaId, lote) {
  const list = loadLotesLocal(userId, fincaId);
  const idx = list.findIndex((l) => l.id === lote.id);
  if (idx >= 0) list[idx] = { ...list[idx], ...lote };
  else list.push(lote);
  localStorage.setItem(key(userId, `lotes_${fincaId}`), JSON.stringify(list));
  return list;
}

export function updateLoteLocal(userId, fincaId, loteId, patch) {
  const list = loadLotesLocal(userId, fincaId).map((l) => (l.id === loteId ? { ...l, ...patch } : l));
  localStorage.setItem(key(userId, `lotes_${fincaId}`), JSON.stringify(list));
  return list;
}

export function loadDiagnosticoHistorial(userId) {
  try {
    const raw = localStorage.getItem(key(userId, 'diagnostico_historial'));
    return raw ? JSON.parse(raw) : [];
  } catch {
    return [];
  }
}

export function saveDiagnosticoHistorial(userId, entry) {
  const list = loadDiagnosticoHistorial(userId);
  list.unshift({ ...entry, id: entry.id || generateUUID(), fecha: entry.fecha || new Date().toISOString() });
  
  // Para evitar QuotaExceededError en localStorage (límite de 5MB):
  // Solo guardamos la imagen base64 de los 2 más recientes. Los demás los guardamos sin imagen.
  const trimmed = list.slice(0, 15).map((item, idx) => {
    if (idx >= 2 && item.image) {
      const { image, ...rest } = item;
      return rest;
    }
    return item;
  });

  try {
    localStorage.setItem(key(userId, 'diagnostico_historial'), JSON.stringify(trimmed));
  } catch (err) {
    // Si aún así se llena, vaciamos agresivamente y dejamos solo 1
    if (err.name === 'QuotaExceededError') {
      const ultraTrimmed = trimmed.slice(0, 1);
      localStorage.setItem(key(userId, 'diagnostico_historial'), JSON.stringify(ultraTrimmed));
      return ultraTrimmed;
    }
  }
  return trimmed;
}

export function loadMuestrasLocal(userId, fincaId, loteId) {
  try {
    const raw = localStorage.getItem(key(userId, `muestras_${fincaId}_${loteId}`));
    return raw ? JSON.parse(raw) : [];
  } catch {
    return [];
  }
}

export function saveMuestraLocal(userId, fincaId, loteId, muestra) {
  const list = loadMuestrasLocal(userId, fincaId, loteId);
  const idx = list.findIndex((m) => m.id === muestra.id);
  if (idx >= 0) list[idx] = { ...list[idx], ...muestra };
  else list.push(muestra);
  localStorage.setItem(key(userId, `muestras_${fincaId}_${loteId}`), JSON.stringify(list));
  return list;
}

export function loadDiagnosticosLocal(userId, fincaId, loteId) {
  try {
    const raw = localStorage.getItem(key(userId, `diagnosticos_${fincaId}_${loteId}`));
    return raw ? JSON.parse(raw) : [];
  } catch {
    return [];
  }
}

export function saveDiagnosticoLocal(userId, fincaId, loteId, diagnostico) {
  const list = loadDiagnosticosLocal(userId, fincaId, loteId);
  const idx = list.findIndex((d) => d.id === diagnostico.id);
  if (idx >= 0) list[idx] = { ...list[idx], ...diagnostico };
  else list.push(diagnostico);
  localStorage.setItem(key(userId, `diagnosticos_${fincaId}_${loteId}`), JSON.stringify(list));
  return list;
}

export function updateDiagnosticoLocal(userId, fincaId, loteId, diagnosticoId, patch) {
  const list = loadDiagnosticosLocal(userId, fincaId, loteId).map((d) =>
    d.id === diagnosticoId ? { ...d, ...patch } : d
  );
  localStorage.setItem(key(userId, `diagnosticos_${fincaId}_${loteId}`), JSON.stringify(list));
  return list;
}

export function getDiagnosticoPorMuestra(userId, fincaId, loteId, muestraId) {
  return loadDiagnosticosLocal(userId, fincaId, loteId).find((d) => d.muestraID === muestraId) || null;
}

export function updateHistorialEntry(userId, entryId, patch) {
  const list = loadDiagnosticoHistorial(userId).map((e) =>
    e.id === entryId ? { ...e, ...patch } : e
  );
  localStorage.setItem(key(userId, 'diagnostico_historial'), JSON.stringify(list));
  return list;
}

export function saveWorkflowState(userId, { fincaSel, loteSel, muestraSel, activeTab }) {
  try {
    localStorage.setItem(
      key(userId, 'workflow'),
      JSON.stringify({ fincaSel, loteSel, muestraSel, activeTab, updatedAt: new Date().toISOString() }),
    );
  } catch { /* quota */ }
}

export function loadWorkflowState(userId) {
  try {
    const raw = localStorage.getItem(key(userId, 'workflow'));
    return raw ? JSON.parse(raw) : null;
  } catch {
    return null;
  }
}

export function clearLocalDataForUser(userId) {
  const keys = Object.keys(localStorage).filter((k) => k.startsWith(`${PREFIX}${userId}_`));
  keys.forEach((k) => localStorage.removeItem(k));
}

// ── Nodos (Cámaras) ──────────────────────────────────────────
export function loadNodosLocal(userId, fincaId) {
  try {
    const raw = localStorage.getItem(key(userId, `nodos_${fincaId}`));
    return raw ? JSON.parse(raw) : [];
  } catch {
    return [];
  }
}

export function saveNodoLocal(userId, fincaId, nodo) {
  const list = loadNodosLocal(userId, fincaId);
  const idx = list.findIndex((n) => n.id === nodo.id);
  if (idx >= 0) list[idx] = { ...list[idx], ...nodo };
  else list.push(nodo);
  localStorage.setItem(key(userId, `nodos_${fincaId}`), JSON.stringify(list));
  return list;
}

export function updateNodoLocal(userId, fincaId, nodoId, patch) {
  const list = loadNodosLocal(userId, fincaId).map((n) => (n.id === nodoId ? { ...n, ...patch } : n));
  localStorage.setItem(key(userId, `nodos_${fincaId}`), JSON.stringify(list));
  return list;
}
