import 'dart:async';

import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../../../core/network/dio_providers.dart';
import '../../../core/session/session_events.dart';
import '../domain/entities/auth_session.dart';
import '../domain/entities/mis_tenants.dart';
import '../domain/entities/perfil.dart';
import '../domain/entities/permiso.dart';
import 'auth_usecases_providers.dart';

part 'auth_controller.freezed.dart';
part 'auth_controller.g.dart';

/// `unknown` representa el estado antes de que `restoreSession()` resuelva
/// (splash). El propio `AsyncNotifier` ya expone eso como `AsyncLoading`
/// mientras `build()` corre — `unknown` queda como valor explícito para que
/// el redirect de GoRouter tenga un caso exhaustivo y legible, sin depender
/// implícitamente de `AsyncLoading` (ver `app_router.dart`).
@freezed
sealed class AuthState with _$AuthState {
  const factory AuthState.unknown() = AuthUnknown;
  const factory AuthState.unauthenticated() = AuthUnauthenticated;
  const factory AuthState.authenticated({
    required AuthSession session,
    required Perfil perfil,
    required MisTenants tenants,
    required List<Permiso> permisos,
  }) = AuthAuthenticated;
}

/// Estado global de sesión — `keepAlive` porque vive durante toda la vida de
/// la app, no por pantalla.
@Riverpod(keepAlive: true)
class AuthController extends _$AuthController {
  StreamSubscription<SessionEventType>? _sessionSub;

  @override
  Future<AuthState> build() async {
    final events = ref.watch(sessionEventsProvider);
    unawaited(_sessionSub?.cancel());
    _sessionSub = events.stream.listen(_onSessionEvent);
    ref.onDispose(() => _sessionSub?.cancel());

    final restore = ref.read(restoreSessionUseCaseProvider);
    final session = await restore();
    if (session == null) return const AuthState.unauthenticated();

    return _cargarContexto(session);
  }

  void _onSessionEvent(SessionEventType event) {
    if (event == SessionEventType.expired) {
      state = const AsyncData(AuthState.unauthenticated());
    }
  }

  Future<AuthState> _cargarContexto(AuthSession session) async {
    final cargarContexto = ref.read(cargarContextoUsuarioUseCaseProvider);
    final contexto = await cargarContexto();
    return AuthState.authenticated(
      session: session,
      perfil: contexto.perfil,
      tenants: contexto.tenants,
      permisos: contexto.permisos,
    );
  }

  /// Lanza `AppException` en caso de error — la página de login la captura
  /// para mostrarla inline, en vez de que el estado global pase por
  /// `AsyncLoading`/`AsyncError` (eso confundiría al redirect del router
  /// durante lo que es, en esencia, una interacción de un formulario).
  Future<void> login({required String correo, required String password}) async {
    final login = ref.read(loginUseCaseProvider);
    final session = await login(correo: correo, password: password);
    state = AsyncData(await _cargarContexto(session));
  }

  /// Best-effort en el repositorio: nunca lanza, siempre deja al usuario
  /// deslogueado localmente.
  Future<void> logout({bool todasLasSesiones = false}) async {
    final logout = ref.read(logoutUseCaseProvider);
    await logout(todasLasSesiones: todasLasSesiones);
    state = const AsyncData(AuthState.unauthenticated());
  }

  /// Igual que `login`: lanza en caso de error para que quien lo invoque
  /// (selector de tenant en Perfil) lo maneje localmente.
  Future<void> switchTenant(String tenantId) async {
    final switchTenant = ref.read(switchTenantUseCaseProvider);
    final session = await switchTenant(tenantId);
    state = AsyncData(await _cargarContexto(session));
  }
}
