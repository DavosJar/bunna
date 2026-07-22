// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'dio_providers.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning
/// Singleton de proceso: memoria + secure storage de los tokens JWT.
/// `main()` llama a `loadFromDisk()` antes de `runApp` (ver ARQUITECTURA.md
/// §1, `main.dart`).

@ProviderFor(tokenStore)
const tokenStoreProvider = TokenStoreProvider._();

/// Singleton de proceso: memoria + secure storage de los tokens JWT.
/// `main()` llama a `loadFromDisk()` antes de `runApp` (ver ARQUITECTURA.md
/// §1, `main.dart`).

final class TokenStoreProvider
    extends $FunctionalProvider<TokenStore, TokenStore, TokenStore>
    with $Provider<TokenStore> {
  /// Singleton de proceso: memoria + secure storage de los tokens JWT.
  /// `main()` llama a `loadFromDisk()` antes de `runApp` (ver ARQUITECTURA.md
  /// §1, `main.dart`).
  const TokenStoreProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'tokenStoreProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$tokenStoreHash();

  @$internal
  @override
  $ProviderElement<TokenStore> $createElement($ProviderPointer pointer) =>
      $ProviderElement(pointer);

  @override
  TokenStore create(Ref ref) {
    return tokenStore(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(TokenStore value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<TokenStore>(value),
    );
  }
}

String _$tokenStoreHash() => r'1a3723cb9c87a80308bd518f02c3e17b42cdd4e8';

@ProviderFor(sessionEvents)
const sessionEventsProvider = SessionEventsProvider._();

final class SessionEventsProvider
    extends $FunctionalProvider<SessionEvents, SessionEvents, SessionEvents>
    with $Provider<SessionEvents> {
  const SessionEventsProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'sessionEventsProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$sessionEventsHash();

  @$internal
  @override
  $ProviderElement<SessionEvents> $createElement($ProviderPointer pointer) =>
      $ProviderElement(pointer);

  @override
  SessionEvents create(Ref ref) {
    return sessionEvents(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(SessionEvents value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<SessionEvents>(value),
    );
  }
}

String _$sessionEventsHash() => r'e11d40f19d2b43c3e06c460fc54ac3882ba22132';

/// Single-flight de refresh compartido por AMBOS clientes Dio — es lo que
/// evita el doble refresh cuando identidad y fincas dan 401 a la vez.

@ProviderFor(tokenRefreshCoordinator)
const tokenRefreshCoordinatorProvider = TokenRefreshCoordinatorProvider._();

/// Single-flight de refresh compartido por AMBOS clientes Dio — es lo que
/// evita el doble refresh cuando identidad y fincas dan 401 a la vez.

final class TokenRefreshCoordinatorProvider
    extends
        $FunctionalProvider<
          TokenRefreshCoordinator,
          TokenRefreshCoordinator,
          TokenRefreshCoordinator
        >
    with $Provider<TokenRefreshCoordinator> {
  /// Single-flight de refresh compartido por AMBOS clientes Dio — es lo que
  /// evita el doble refresh cuando identidad y fincas dan 401 a la vez.
  const TokenRefreshCoordinatorProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'tokenRefreshCoordinatorProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$tokenRefreshCoordinatorHash();

  @$internal
  @override
  $ProviderElement<TokenRefreshCoordinator> $createElement(
    $ProviderPointer pointer,
  ) => $ProviderElement(pointer);

  @override
  TokenRefreshCoordinator create(Ref ref) {
    return tokenRefreshCoordinator(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(TokenRefreshCoordinator value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<TokenRefreshCoordinator>(value),
    );
  }
}

String _$tokenRefreshCoordinatorHash() =>
    r'1796103683c87ef79d08bce689bdd5c730a3ad25';

@ProviderFor(identidadDio)
const identidadDioProvider = IdentidadDioProvider._();

final class IdentidadDioProvider extends $FunctionalProvider<Dio, Dio, Dio>
    with $Provider<Dio> {
  const IdentidadDioProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'identidadDioProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$identidadDioHash();

  @$internal
  @override
  $ProviderElement<Dio> $createElement($ProviderPointer pointer) =>
      $ProviderElement(pointer);

  @override
  Dio create(Ref ref) {
    return identidadDio(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(Dio value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<Dio>(value),
    );
  }
}

String _$identidadDioHash() => r'2c734ee3a36179e8e5aeab7567b3afb9405084ef';

@ProviderFor(fincasDio)
const fincasDioProvider = FincasDioProvider._();

final class FincasDioProvider extends $FunctionalProvider<Dio, Dio, Dio>
    with $Provider<Dio> {
  const FincasDioProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'fincasDioProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$fincasDioHash();

  @$internal
  @override
  $ProviderElement<Dio> $createElement($ProviderPointer pointer) =>
      $ProviderElement(pointer);

  @override
  Dio create(Ref ref) {
    return fincasDio(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(Dio value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<Dio>(value),
    );
  }
}

String _$fincasDioHash() => r'e076d72b7ae81750c41e6935df0fa1833cbbcefb';

/// YOLO no recibe JWT (ver ARQUITECTURA.md §4) y no tiene contrato de error
/// estable, por eso reutiliza el mapper tolerante de fincas en vez de tener
/// uno propio.

@ProviderFor(yoloDio)
const yoloDioProvider = YoloDioProvider._();

/// YOLO no recibe JWT (ver ARQUITECTURA.md §4) y no tiene contrato de error
/// estable, por eso reutiliza el mapper tolerante de fincas en vez de tener
/// uno propio.

final class YoloDioProvider extends $FunctionalProvider<Dio, Dio, Dio>
    with $Provider<Dio> {
  /// YOLO no recibe JWT (ver ARQUITECTURA.md §4) y no tiene contrato de error
  /// estable, por eso reutiliza el mapper tolerante de fincas en vez de tener
  /// uno propio.
  const YoloDioProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'yoloDioProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$yoloDioHash();

  @$internal
  @override
  $ProviderElement<Dio> $createElement($ProviderPointer pointer) =>
      $ProviderElement(pointer);

  @override
  Dio create(Ref ref) {
    return yoloDio(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(Dio value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<Dio>(value),
    );
  }
}

String _$yoloDioHash() => r'98004b0ef905006b14a6ae9d7872aad08460e606';
