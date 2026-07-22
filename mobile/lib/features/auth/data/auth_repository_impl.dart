import '../../../core/error/app_exception.dart';
import '../../../core/session/auth_tokens.dart';
import '../../../core/session/jwt_claims.dart';
import '../../../core/session/token_refresh_coordinator.dart';
import '../../../core/session/token_store.dart';
import 'auth_api.dart';
import 'dtos/token_session_dto.dart';
import '../domain/auth_repository.dart';
import '../domain/entities/auth_session.dart';
import '../domain/entities/mis_tenants.dart';
import '../domain/entities/perfil.dart';
import '../domain/entities/permiso.dart';
import '../domain/entities/tenant_con_rol.dart';

final class AuthRepositoryImpl implements AuthRepository {
  AuthRepositoryImpl({
    required AuthApi api,
    required TokenStore tokenStore,
    required TokenRefreshCoordinator refreshCoordinator,
  }) : _api = api,
       _tokenStore = tokenStore,
       _refreshCoordinator = refreshCoordinator;

  final AuthApi _api;
  final TokenStore _tokenStore;
  final TokenRefreshCoordinator _refreshCoordinator;

  @override
  Future<AuthSession> login({
    required String correo,
    required String password,
  }) async {
    final dto = await _api.login(correo: correo, password: password);
    return _persistirYConstruirSesion(dto);
  }

  @override
  Future<AuthSession> switchTenant(String tenantId) async {
    final dto = await _api.switchTenant(tenantId);
    return _persistirYConstruirSesion(dto);
  }

  Future<AuthSession> _persistirYConstruirSesion(TokenSessionDto dto) async {
    final tokens = AuthTokens(
      accessToken: dto.accessToken,
      refreshToken: dto.refreshToken,
      accessExpiresAt: DateTime.now().add(Duration(seconds: dto.expiresIn)),
    );
    await _tokenStore.save(tokens);

    final claims = JwtClaims.tryDecode(dto.accessToken);
    return AuthSession(
      usuarioId: dto.usuarioId,
      tenantId: dto.tenantId,
      rol: dto.rol,
      sesionId: claims?.sesionId,
    );
  }

  @override
  Future<AuthSession?> restoreSession() async {
    final stored = _tokenStore.current;
    if (stored == null) return null;

    var tokens = stored;
    if (tokens.isAccessExpired) {
      try {
        tokens = await _refreshCoordinator.refresh(
          failedAccessToken: tokens.accessToken,
        );
      } on SessionExpiredException {
        // El coordinator ya limpió el TokenStore y emitió el evento.
        return null;
      } on AppException {
        // Fallo transitorio (red/servidor caído al arrancar): NO borramos la
        // sesión local. Seguimos con el access token expirado; el primer
        // request real disparará el refresh reactivo del RefreshInterceptor.
      }
    }

    final claims = JwtClaims.tryDecode(tokens.accessToken);
    if (claims == null) {
      // Token corrupto o imposible de leer: no hay forma de reconstruir la
      // sesión localmente.
      await _tokenStore.clear();
      return null;
    }

    return AuthSession(
      usuarioId: claims.usuarioId,
      tenantId: claims.tenantId ?? '',
      rol: claims.rol ?? '',
      sesionId: claims.sesionId,
    );
  }

  @override
  Future<void> logout({bool todasLasSesiones = false}) async {
    try {
      if (todasLasSesiones) {
        await _api.logoutAll();
      } else {
        await _api.logout();
      }
    } catch (_) {
      // Best-effort: si el servidor no responde, igual limpiamos localmente.
    } finally {
      await _tokenStore.clear();
    }
  }

  @override
  Future<Perfil> getMiPerfil() async {
    final dto = await _api.getMiPerfil();
    return Perfil(
      id: dto.id,
      correo: dto.correo,
      nombre: dto.nombre,
      apellido: dto.apellido,
      telefono: dto.telefono,
      estado: dto.estado,
      creadoEn: dto.creadoEn,
    );
  }

  @override
  Future<MisTenants> getMisTenants() async {
    final dto = await _api.getMisTenants();
    return MisTenants(
      tenants: dto.tenants
          .map(
            (t) => TenantConRol(
              id: t.id,
              nombre: t.nombre,
              slug: t.slug,
              rol: t.rol,
              esPropio: t.esPropio,
            ),
          )
          .toList(growable: false),
      propioId: dto.propioId,
    );
  }

  @override
  Future<List<Permiso>> getMisPermisos() async {
    final dto = await _api.getMisPermisos();
    return dto.permisos
        .map(
          (p) => Permiso(
            codigo: p.codigo,
            nombre: p.nombre,
            descripcion: p.descripcion,
            modulo: p.modulo,
          ),
        )
        .toList(growable: false);
  }
}
