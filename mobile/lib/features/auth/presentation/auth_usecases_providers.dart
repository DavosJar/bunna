import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../data/auth_providers.dart';
import '../domain/usecases/cargar_contexto_usuario.dart';
import '../domain/usecases/login.dart';
import '../domain/usecases/logout.dart';
import '../domain/usecases/restore_session.dart';
import '../domain/usecases/switch_tenant.dart';

part 'auth_usecases_providers.g.dart';

/// Único lugar donde se conectan los casos de uso (domain, sin Riverpod) con
/// `authRepositoryProvider` (data). Vive en `presentation/` a propósito: es
/// la capa de composición, la única a la que le está permitido conocer
/// `domain` y la instancia concreta de `data` a la vez.
@riverpod
Login loginUseCase(Ref ref) => Login(ref.watch(authRepositoryProvider));

@riverpod
RestoreSession restoreSessionUseCase(Ref ref) =>
    RestoreSession(ref.watch(authRepositoryProvider));

@riverpod
Logout logoutUseCase(Ref ref) => Logout(ref.watch(authRepositoryProvider));

@riverpod
SwitchTenant switchTenantUseCase(Ref ref) =>
    SwitchTenant(ref.watch(authRepositoryProvider));

@riverpod
CargarContextoUsuario cargarContextoUsuarioUseCase(Ref ref) =>
    CargarContextoUsuario(ref.watch(authRepositoryProvider));
