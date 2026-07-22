import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../../../core/network/dio_providers.dart';
import '../domain/auth_repository.dart';
import 'auth_api.dart';
import 'auth_repository_impl.dart';

part 'auth_providers.g.dart';

@Riverpod(keepAlive: true)
AuthApi authApi(Ref ref) => AuthApi(ref.watch(identidadDioProvider));

/// Expone el contrato `AuthRepository`, no la implementación — domain y
/// presentation nunca dependen de `AuthRepositoryImpl` directamente.
@Riverpod(keepAlive: true)
AuthRepository authRepository(Ref ref) {
  return AuthRepositoryImpl(
    api: ref.watch(authApiProvider),
    tokenStore: ref.watch(tokenStoreProvider),
    refreshCoordinator: ref.watch(tokenRefreshCoordinatorProvider),
  );
}
