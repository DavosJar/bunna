// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'quick_scan_controller.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning
/// Escaneo YOLO independiente (tab central "Escanear"): analiza una foto sin
/// vincularla a ninguna muestra. Para vincular a una muestra concreta se usa
/// el flujo dentro de un lote.

@ProviderFor(QuickScanController)
const quickScanControllerProvider = QuickScanControllerProvider._();

/// Escaneo YOLO independiente (tab central "Escanear"): analiza una foto sin
/// vincularla a ninguna muestra. Para vincular a una muestra concreta se usa
/// el flujo dentro de un lote.
final class QuickScanControllerProvider
    extends $NotifierProvider<QuickScanController, QuickScanState> {
  /// Escaneo YOLO independiente (tab central "Escanear"): analiza una foto sin
  /// vincularla a ninguna muestra. Para vincular a una muestra concreta se usa
  /// el flujo dentro de un lote.
  const QuickScanControllerProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'quickScanControllerProvider',
        isAutoDispose: true,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$quickScanControllerHash();

  @$internal
  @override
  QuickScanController create() => QuickScanController();

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(QuickScanState value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<QuickScanState>(value),
    );
  }
}

String _$quickScanControllerHash() =>
    r'2e6944dafc3471870356422df64db495694ba909';

/// Escaneo YOLO independiente (tab central "Escanear"): analiza una foto sin
/// vincularla a ninguna muestra. Para vincular a una muestra concreta se usa
/// el flujo dentro de un lote.

abstract class _$QuickScanController extends $Notifier<QuickScanState> {
  QuickScanState build();
  @$mustCallSuper
  @override
  void runBuild() {
    final created = build();
    final ref = this.ref as $Ref<QuickScanState, QuickScanState>;
    final element =
        ref.element
            as $ClassProviderElement<
              AnyNotifier<QuickScanState, QuickScanState>,
              QuickScanState,
              Object?,
              Object?
            >;
    element.handleValue(ref, created);
  }
}
