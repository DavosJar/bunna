// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'fincas_controller.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning
/// Lista de fincas del tenant activo. `autoDispose` (default de `@riverpod`):
/// la caché vive solo mientras la pantalla está montada — sin persistencia
/// de datos de negocio, se relee del backend en cada entrada. Las mutaciones
/// invalidan el provider para forzar un refetch.

@ProviderFor(FincasController)
const fincasControllerProvider = FincasControllerProvider._();

/// Lista de fincas del tenant activo. `autoDispose` (default de `@riverpod`):
/// la caché vive solo mientras la pantalla está montada — sin persistencia
/// de datos de negocio, se relee del backend en cada entrada. Las mutaciones
/// invalidan el provider para forzar un refetch.
final class FincasControllerProvider
    extends $AsyncNotifierProvider<FincasController, List<Finca>> {
  /// Lista de fincas del tenant activo. `autoDispose` (default de `@riverpod`):
  /// la caché vive solo mientras la pantalla está montada — sin persistencia
  /// de datos de negocio, se relee del backend en cada entrada. Las mutaciones
  /// invalidan el provider para forzar un refetch.
  const FincasControllerProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'fincasControllerProvider',
        isAutoDispose: true,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$fincasControllerHash();

  @$internal
  @override
  FincasController create() => FincasController();
}

String _$fincasControllerHash() => r'8f723b0beb544583d76649ac8173004faf349250';

/// Lista de fincas del tenant activo. `autoDispose` (default de `@riverpod`):
/// la caché vive solo mientras la pantalla está montada — sin persistencia
/// de datos de negocio, se relee del backend en cada entrada. Las mutaciones
/// invalidan el provider para forzar un refetch.

abstract class _$FincasController extends $AsyncNotifier<List<Finca>> {
  FutureOr<List<Finca>> build();
  @$mustCallSuper
  @override
  void runBuild() {
    final created = build();
    final ref = this.ref as $Ref<AsyncValue<List<Finca>>, List<Finca>>;
    final element =
        ref.element
            as $ClassProviderElement<
              AnyNotifier<AsyncValue<List<Finca>>, List<Finca>>,
              AsyncValue<List<Finca>>,
              Object?,
              Object?
            >;
    element.handleValue(ref, created);
  }
}
