// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'permisos.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning
/// Espejo de `usePermisos`/`roleAccess.js` del frontend web: `sys_admin`
/// pasa cualquier permiso; el resto se decide por la lista real de
/// `/mis-permisos`. Devuelve `false` mientras no hay sesión autenticada.

@ProviderFor(puede)
const puedeProvider = PuedeFamily._();

/// Espejo de `usePermisos`/`roleAccess.js` del frontend web: `sys_admin`
/// pasa cualquier permiso; el resto se decide por la lista real de
/// `/mis-permisos`. Devuelve `false` mientras no hay sesión autenticada.

final class PuedeProvider extends $FunctionalProvider<bool, bool, bool>
    with $Provider<bool> {
  /// Espejo de `usePermisos`/`roleAccess.js` del frontend web: `sys_admin`
  /// pasa cualquier permiso; el resto se decide por la lista real de
  /// `/mis-permisos`. Devuelve `false` mientras no hay sesión autenticada.
  const PuedeProvider._({
    required PuedeFamily super.from,
    required String super.argument,
  }) : super(
         retry: null,
         name: r'puedeProvider',
         isAutoDispose: true,
         dependencies: null,
         $allTransitiveDependencies: null,
       );

  @override
  String debugGetCreateSourceHash() => _$puedeHash();

  @override
  String toString() {
    return r'puedeProvider'
        ''
        '($argument)';
  }

  @$internal
  @override
  $ProviderElement<bool> $createElement($ProviderPointer pointer) =>
      $ProviderElement(pointer);

  @override
  bool create(Ref ref) {
    final argument = this.argument as String;
    return puede(ref, argument);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(bool value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<bool>(value),
    );
  }

  @override
  bool operator ==(Object other) {
    return other is PuedeProvider && other.argument == argument;
  }

  @override
  int get hashCode {
    return argument.hashCode;
  }
}

String _$puedeHash() => r'832a616e8e18968047ba118a0928e54542361fe7';

/// Espejo de `usePermisos`/`roleAccess.js` del frontend web: `sys_admin`
/// pasa cualquier permiso; el resto se decide por la lista real de
/// `/mis-permisos`. Devuelve `false` mientras no hay sesión autenticada.

final class PuedeFamily extends $Family
    with $FunctionalFamilyOverride<bool, String> {
  const PuedeFamily._()
    : super(
        retry: null,
        name: r'puedeProvider',
        dependencies: null,
        $allTransitiveDependencies: null,
        isAutoDispose: true,
      );

  /// Espejo de `usePermisos`/`roleAccess.js` del frontend web: `sys_admin`
  /// pasa cualquier permiso; el resto se decide por la lista real de
  /// `/mis-permisos`. Devuelve `false` mientras no hay sesión autenticada.

  PuedeProvider call(String codigo) =>
      PuedeProvider._(argument: codigo, from: this);

  @override
  String toString() => r'puedeProvider';
}
