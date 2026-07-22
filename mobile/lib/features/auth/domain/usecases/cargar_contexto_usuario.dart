import '../../../../core/error/app_exception.dart';
import '../auth_repository.dart';
import '../entities/mis_tenants.dart';
import '../entities/perfil.dart';
import '../entities/permiso.dart';

typedef ContextoUsuario = ({
  Perfil perfil,
  MisTenants tenants,
  List<Permiso> permisos,
});

const _tenantsVacios = MisTenants(tenants: [], propioId: '');

/// Orquesta el bootstrap post-login: perfil → tenants → permisos.
final class CargarContextoUsuario {
  const CargarContextoUsuario(this._repo);

  final AuthRepository _repo;

  /// Secuencial a propósito, NUNCA `Future.wait`: identidad limita a
  /// 10 req/1s por IP y esto se dispara justo después del login, cuando ya
  /// se gastó una request en `/auth/login` (ver ARQUITECTURA.md §0,
  /// "Bootstrap post-login").
  ///
  /// `tenants`/`permisos` degradan a vacío en vez de abortar el login si
  /// fallan — igual que `fetchMisTenants`/`fetchMisPermisos` en
  /// `AuthContext.jsx` del frontend web. Caso real confirmado contra el
  /// backend: usuarios `sys_admin` no tienen tenant, su JWT no lleva claim
  /// `rol`, y `GET /mis-permisos` trata `rol` vacío como "sin token" → 401,
  /// aunque el mismo token acabe de validar `GET /mi-perfil` con éxito. No
  /// es una sesión inválida (si lo fuera, `getMiPerfil()` ya habría fallado
  /// primero) — es una particularidad del backend que no se puede corregir
  /// desde el cliente. `sys_admin` de todas formas tiene acceso total vía
  /// `puede()` (bypass explícito por rol), así que una lista de permisos
  /// vacía no le quita funcionalidad.
  Future<ContextoUsuario> call() async {
    final perfil = await _repo.getMiPerfil();
    final tenants = await _conFallback(_repo.getMisTenants, _tenantsVacios);
    final permisos = await _conFallback(
      _repo.getMisPermisos,
      const <Permiso>[],
    );
    return (perfil: perfil, tenants: tenants, permisos: permisos);
  }

  Future<T> _conFallback<T>(Future<T> Function() call, T fallback) async {
    try {
      return await call();
    } on AppException {
      return fallback;
    }
  }
}
