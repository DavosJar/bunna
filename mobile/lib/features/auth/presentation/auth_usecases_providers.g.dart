// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'auth_usecases_providers.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning
/// Único lugar donde se conectan los casos de uso (domain, sin Riverpod) con
/// `authRepositoryProvider` (data). Vive en `presentation/` a propósito: es
/// la capa de composición, la única a la que le está permitido conocer
/// `domain` y la instancia concreta de `data` a la vez.

@ProviderFor(loginUseCase)
const loginUseCaseProvider = LoginUseCaseProvider._();

/// Único lugar donde se conectan los casos de uso (domain, sin Riverpod) con
/// `authRepositoryProvider` (data). Vive en `presentation/` a propósito: es
/// la capa de composición, la única a la que le está permitido conocer
/// `domain` y la instancia concreta de `data` a la vez.

final class LoginUseCaseProvider
    extends $FunctionalProvider<Login, Login, Login>
    with $Provider<Login> {
  /// Único lugar donde se conectan los casos de uso (domain, sin Riverpod) con
  /// `authRepositoryProvider` (data). Vive en `presentation/` a propósito: es
  /// la capa de composición, la única a la que le está permitido conocer
  /// `domain` y la instancia concreta de `data` a la vez.
  const LoginUseCaseProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'loginUseCaseProvider',
        isAutoDispose: true,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$loginUseCaseHash();

  @$internal
  @override
  $ProviderElement<Login> $createElement($ProviderPointer pointer) =>
      $ProviderElement(pointer);

  @override
  Login create(Ref ref) {
    return loginUseCase(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(Login value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<Login>(value),
    );
  }
}

String _$loginUseCaseHash() => r'3dcd67871dd409e48fec25dccaec07b22901b35d';

@ProviderFor(restoreSessionUseCase)
const restoreSessionUseCaseProvider = RestoreSessionUseCaseProvider._();

final class RestoreSessionUseCaseProvider
    extends $FunctionalProvider<RestoreSession, RestoreSession, RestoreSession>
    with $Provider<RestoreSession> {
  const RestoreSessionUseCaseProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'restoreSessionUseCaseProvider',
        isAutoDispose: true,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$restoreSessionUseCaseHash();

  @$internal
  @override
  $ProviderElement<RestoreSession> $createElement($ProviderPointer pointer) =>
      $ProviderElement(pointer);

  @override
  RestoreSession create(Ref ref) {
    return restoreSessionUseCase(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(RestoreSession value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<RestoreSession>(value),
    );
  }
}

String _$restoreSessionUseCaseHash() =>
    r'1ba05d8a62afeb7037d4a05eae28f333513189a9';

@ProviderFor(logoutUseCase)
const logoutUseCaseProvider = LogoutUseCaseProvider._();

final class LogoutUseCaseProvider
    extends $FunctionalProvider<Logout, Logout, Logout>
    with $Provider<Logout> {
  const LogoutUseCaseProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'logoutUseCaseProvider',
        isAutoDispose: true,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$logoutUseCaseHash();

  @$internal
  @override
  $ProviderElement<Logout> $createElement($ProviderPointer pointer) =>
      $ProviderElement(pointer);

  @override
  Logout create(Ref ref) {
    return logoutUseCase(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(Logout value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<Logout>(value),
    );
  }
}

String _$logoutUseCaseHash() => r'9dad12162a83228303546bc3426be5e3d03b3d81';

@ProviderFor(switchTenantUseCase)
const switchTenantUseCaseProvider = SwitchTenantUseCaseProvider._();

final class SwitchTenantUseCaseProvider
    extends $FunctionalProvider<SwitchTenant, SwitchTenant, SwitchTenant>
    with $Provider<SwitchTenant> {
  const SwitchTenantUseCaseProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'switchTenantUseCaseProvider',
        isAutoDispose: true,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$switchTenantUseCaseHash();

  @$internal
  @override
  $ProviderElement<SwitchTenant> $createElement($ProviderPointer pointer) =>
      $ProviderElement(pointer);

  @override
  SwitchTenant create(Ref ref) {
    return switchTenantUseCase(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(SwitchTenant value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<SwitchTenant>(value),
    );
  }
}

String _$switchTenantUseCaseHash() =>
    r'31aaaa511b50cfcc50366257cb37ccd35f9c6460';

@ProviderFor(cargarContextoUsuarioUseCase)
const cargarContextoUsuarioUseCaseProvider =
    CargarContextoUsuarioUseCaseProvider._();

final class CargarContextoUsuarioUseCaseProvider
    extends
        $FunctionalProvider<
          CargarContextoUsuario,
          CargarContextoUsuario,
          CargarContextoUsuario
        >
    with $Provider<CargarContextoUsuario> {
  const CargarContextoUsuarioUseCaseProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'cargarContextoUsuarioUseCaseProvider',
        isAutoDispose: true,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$cargarContextoUsuarioUseCaseHash();

  @$internal
  @override
  $ProviderElement<CargarContextoUsuario> $createElement(
    $ProviderPointer pointer,
  ) => $ProviderElement(pointer);

  @override
  CargarContextoUsuario create(Ref ref) {
    return cargarContextoUsuarioUseCase(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(CargarContextoUsuario value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<CargarContextoUsuario>(value),
    );
  }
}

String _$cargarContextoUsuarioUseCaseHash() =>
    r'd42b6675c047d86c72c18ce4aef40ae7137f12b5';
