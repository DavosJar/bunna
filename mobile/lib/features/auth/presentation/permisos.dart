import 'package:riverpod_annotation/riverpod_annotation.dart';

import 'auth_controller.dart';

part 'permisos.g.dart';

/// Espejo de `usePermisos`/`roleAccess.js` del frontend web: `sys_admin`
/// pasa cualquier permiso; el resto se decide por la lista real de
/// `/mis-permisos`. Devuelve `false` mientras no hay sesión autenticada.
@riverpod
bool puede(Ref ref, String codigo) {
  final authState = ref.watch(authControllerProvider).value;
  if (authState is! AuthAuthenticated) return false;
  if (authState.session.rol == 'sys_admin') return true;
  return authState.permisos.any((p) => p.codigo == codigo);
}
