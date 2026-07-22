// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'muestras_controller.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning
/// Muestras de un lote (family por fincaId + loteId).

@ProviderFor(MuestrasController)
const muestrasControllerProvider = MuestrasControllerFamily._();

/// Muestras de un lote (family por fincaId + loteId).
final class MuestrasControllerProvider
    extends $AsyncNotifierProvider<MuestrasController, List<Muestra>> {
  /// Muestras de un lote (family por fincaId + loteId).
  const MuestrasControllerProvider._({
    required MuestrasControllerFamily super.from,
    required (String, String) super.argument,
  }) : super(
         retry: null,
         name: r'muestrasControllerProvider',
         isAutoDispose: true,
         dependencies: null,
         $allTransitiveDependencies: null,
       );

  @override
  String debugGetCreateSourceHash() => _$muestrasControllerHash();

  @override
  String toString() {
    return r'muestrasControllerProvider'
        ''
        '$argument';
  }

  @$internal
  @override
  MuestrasController create() => MuestrasController();

  @override
  bool operator ==(Object other) {
    return other is MuestrasControllerProvider && other.argument == argument;
  }

  @override
  int get hashCode {
    return argument.hashCode;
  }
}

String _$muestrasControllerHash() =>
    r'8dd575a0a323db95327073c6ed76ca5f32a262e7';

/// Muestras de un lote (family por fincaId + loteId).

final class MuestrasControllerFamily extends $Family
    with
        $ClassFamilyOverride<
          MuestrasController,
          AsyncValue<List<Muestra>>,
          List<Muestra>,
          FutureOr<List<Muestra>>,
          (String, String)
        > {
  const MuestrasControllerFamily._()
    : super(
        retry: null,
        name: r'muestrasControllerProvider',
        dependencies: null,
        $allTransitiveDependencies: null,
        isAutoDispose: true,
      );

  /// Muestras de un lote (family por fincaId + loteId).

  MuestrasControllerProvider call(String fincaId, String loteId) =>
      MuestrasControllerProvider._(argument: (fincaId, loteId), from: this);

  @override
  String toString() => r'muestrasControllerProvider';
}

/// Muestras de un lote (family por fincaId + loteId).

abstract class _$MuestrasController extends $AsyncNotifier<List<Muestra>> {
  late final _$args = ref.$arg as (String, String);
  String get fincaId => _$args.$1;
  String get loteId => _$args.$2;

  FutureOr<List<Muestra>> build(String fincaId, String loteId);
  @$mustCallSuper
  @override
  void runBuild() {
    final created = build(_$args.$1, _$args.$2);
    final ref = this.ref as $Ref<AsyncValue<List<Muestra>>, List<Muestra>>;
    final element =
        ref.element
            as $ClassProviderElement<
              AnyNotifier<AsyncValue<List<Muestra>>, List<Muestra>>,
              AsyncValue<List<Muestra>>,
              Object?,
              Object?
            >;
    element.handleValue(ref, created);
  }
}
