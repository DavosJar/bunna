// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'diagnostico_flow_controller.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning
/// Orquesta el flujo de diagnóstico de UNA muestra: foto → análisis YOLO →
/// vincular como diagnóstico en el backend → aceptar/rechazar. Family por
/// `muestraId`. Imperativo (Notifier), no AsyncNotifier: cada paso lo dispara
/// el usuario y el estado avanza por fases.

@ProviderFor(DiagnosticoFlowController)
const diagnosticoFlowControllerProvider = DiagnosticoFlowControllerFamily._();

/// Orquesta el flujo de diagnóstico de UNA muestra: foto → análisis YOLO →
/// vincular como diagnóstico en el backend → aceptar/rechazar. Family por
/// `muestraId`. Imperativo (Notifier), no AsyncNotifier: cada paso lo dispara
/// el usuario y el estado avanza por fases.
final class DiagnosticoFlowControllerProvider
    extends $NotifierProvider<DiagnosticoFlowController, DiagnosticoFlowState> {
  /// Orquesta el flujo de diagnóstico de UNA muestra: foto → análisis YOLO →
  /// vincular como diagnóstico en el backend → aceptar/rechazar. Family por
  /// `muestraId`. Imperativo (Notifier), no AsyncNotifier: cada paso lo dispara
  /// el usuario y el estado avanza por fases.
  const DiagnosticoFlowControllerProvider._({
    required DiagnosticoFlowControllerFamily super.from,
    required String super.argument,
  }) : super(
         retry: null,
         name: r'diagnosticoFlowControllerProvider',
         isAutoDispose: true,
         dependencies: null,
         $allTransitiveDependencies: null,
       );

  @override
  String debugGetCreateSourceHash() => _$diagnosticoFlowControllerHash();

  @override
  String toString() {
    return r'diagnosticoFlowControllerProvider'
        ''
        '($argument)';
  }

  @$internal
  @override
  DiagnosticoFlowController create() => DiagnosticoFlowController();

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(DiagnosticoFlowState value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<DiagnosticoFlowState>(value),
    );
  }

  @override
  bool operator ==(Object other) {
    return other is DiagnosticoFlowControllerProvider &&
        other.argument == argument;
  }

  @override
  int get hashCode {
    return argument.hashCode;
  }
}

String _$diagnosticoFlowControllerHash() =>
    r'c403e396eb497c736e923d83aeb0a4abbf827840';

/// Orquesta el flujo de diagnóstico de UNA muestra: foto → análisis YOLO →
/// vincular como diagnóstico en el backend → aceptar/rechazar. Family por
/// `muestraId`. Imperativo (Notifier), no AsyncNotifier: cada paso lo dispara
/// el usuario y el estado avanza por fases.

final class DiagnosticoFlowControllerFamily extends $Family
    with
        $ClassFamilyOverride<
          DiagnosticoFlowController,
          DiagnosticoFlowState,
          DiagnosticoFlowState,
          DiagnosticoFlowState,
          String
        > {
  const DiagnosticoFlowControllerFamily._()
    : super(
        retry: null,
        name: r'diagnosticoFlowControllerProvider',
        dependencies: null,
        $allTransitiveDependencies: null,
        isAutoDispose: true,
      );

  /// Orquesta el flujo de diagnóstico de UNA muestra: foto → análisis YOLO →
  /// vincular como diagnóstico en el backend → aceptar/rechazar. Family por
  /// `muestraId`. Imperativo (Notifier), no AsyncNotifier: cada paso lo dispara
  /// el usuario y el estado avanza por fases.

  DiagnosticoFlowControllerProvider call(String muestraId) =>
      DiagnosticoFlowControllerProvider._(argument: muestraId, from: this);

  @override
  String toString() => r'diagnosticoFlowControllerProvider';
}

/// Orquesta el flujo de diagnóstico de UNA muestra: foto → análisis YOLO →
/// vincular como diagnóstico en el backend → aceptar/rechazar. Family por
/// `muestraId`. Imperativo (Notifier), no AsyncNotifier: cada paso lo dispara
/// el usuario y el estado avanza por fases.

abstract class _$DiagnosticoFlowController
    extends $Notifier<DiagnosticoFlowState> {
  late final _$args = ref.$arg as String;
  String get muestraId => _$args;

  DiagnosticoFlowState build(String muestraId);
  @$mustCallSuper
  @override
  void runBuild() {
    final created = build(_$args);
    final ref = this.ref as $Ref<DiagnosticoFlowState, DiagnosticoFlowState>;
    final element =
        ref.element
            as $ClassProviderElement<
              AnyNotifier<DiagnosticoFlowState, DiagnosticoFlowState>,
              DiagnosticoFlowState,
              Object?,
              Object?
            >;
    element.handleValue(ref, created);
  }
}
