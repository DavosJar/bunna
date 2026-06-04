/**
 * Valida formato de contraseña según reglas de seguridad.
 * @param {string} password
 * @returns {{ valida: boolean, errores: string[] }}
 */
export function validarPassword(password) {
  const errores = [];

  if (!password) {
    errores.push('La contraseña es requerida');
    return { valida: false, errores };
  }

  if (password.length < 8) {
    errores.push('Debe tener al menos 8 caracteres');
  }

  let mayuscula = false, minuscula = false, numero = false, noAlfanumerico = false;

  for (const c of password) {
    if (c >= 'A' && c <= 'Z') mayuscula = true;
    else if (c >= 'a' && c <= 'z') minuscula = true;
    else if (c >= '0' && c <= '9') numero = true;
    else if (!((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9'))) noAlfanumerico = true;
  }

  if (!mayuscula) errores.push('Debe contener al menos una mayúscula');
  if (!minuscula) errores.push('Debe contener al menos una minúscula');
  if (!numero) errores.push('Debe contener al menos un número');
  if (!noAlfanumerico) errores.push('Debe contener al menos un carácter especial (@, #, $, etc.)');

  return { valida: errores.length === 0, errores };
}
