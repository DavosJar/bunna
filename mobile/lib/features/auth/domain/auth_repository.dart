import 'entities/auth_session.dart';
import 'entities/mis_tenants.dart';
import 'entities/perfil.dart';
import 'entities/permiso.dart';

/// Todos los métodos lanzan `AppException` (nunca DTOs ni `DioException`) —
/// ver ARQUITECTURA.md §2 y §3.
abstract interface class AuthRepository {
  /// POST /api/v1/identidad/auth/login — persiste tokens en TokenStore.
  Future<AuthSession> login({
    required String correo,
    required String password,
  });

  /// Arranque: lee tokens de secure storage; si el access expiró intenta un
  /// refresh vía coordinator. Devuelve `null` si no hay sesión recuperable.
  Future<AuthSession?> restoreSession();

  /// POST /auth/logout | /auth/logout/all (best-effort) y limpia TokenStore
  /// SIEMPRE, incluso si la llamada al servidor falla.
  Future<void> logout({bool todasLasSesiones = false});

  /// POST /auth/switch-tenant — reemplaza AMBOS tokens (rotación).
  Future<AuthSession> switchTenant(String tenantId);

  /// GET /mi-perfil
  Future<Perfil> getMiPerfil();

  /// GET /tenants/mis-tenants
  Future<MisTenants> getMisTenants();

  /// GET /mis-permisos
  Future<List<Permiso>> getMisPermisos();

  // Diferido a fases posteriores: register, recuperación y verificación de correo.
}
