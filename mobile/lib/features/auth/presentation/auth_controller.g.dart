// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'auth_controller.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning
/// Estado global de sesión — `keepAlive` porque vive durante toda la vida de
/// la app, no por pantalla.

@ProviderFor(AuthController)
const authControllerProvider = AuthControllerProvider._();

/// Estado global de sesión — `keepAlive` porque vive durante toda la vida de
/// la app, no por pantalla.
final class AuthControllerProvider
    extends $AsyncNotifierProvider<AuthController, AuthState> {
  /// Estado global de sesión — `keepAlive` porque vive durante toda la vida de
  /// la app, no por pantalla.
  const AuthControllerProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'authControllerProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$authControllerHash();

  @$internal
  @override
  AuthController create() => AuthController();
}

String _$authControllerHash() => r'4108cd4a040a7c0b9c3e5d5ff90e3584d3a61906';

/// Estado global de sesión — `keepAlive` porque vive durante toda la vida de
/// la app, no por pantalla.

abstract class _$AuthController extends $AsyncNotifier<AuthState> {
  FutureOr<AuthState> build();
  @$mustCallSuper
  @override
  void runBuild() {
    final created = build();
    final ref = this.ref as $Ref<AsyncValue<AuthState>, AuthState>;
    final element =
        ref.element
            as $ClassProviderElement<
              AnyNotifier<AsyncValue<AuthState>, AuthState>,
              AsyncValue<AuthState>,
              Object?,
              Object?
            >;
    element.handleValue(ref, created);
  }
}
