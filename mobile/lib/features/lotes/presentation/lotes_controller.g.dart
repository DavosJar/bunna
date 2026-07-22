// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'lotes_controller.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning
/// Lotes de una finca (family por `fincaId`). `autoDispose`: caché por
/// pantalla, sin persistencia.

@ProviderFor(LotesController)
const lotesControllerProvider = LotesControllerFamily._();

/// Lotes de una finca (family por `fincaId`). `autoDispose`: caché por
/// pantalla, sin persistencia.
final class LotesControllerProvider
    extends $AsyncNotifierProvider<LotesController, List<Lote>> {
  /// Lotes de una finca (family por `fincaId`). `autoDispose`: caché por
  /// pantalla, sin persistencia.
  const LotesControllerProvider._({
    required LotesControllerFamily super.from,
    required String super.argument,
  }) : super(
         retry: null,
         name: r'lotesControllerProvider',
         isAutoDispose: true,
         dependencies: null,
         $allTransitiveDependencies: null,
       );

  @override
  String debugGetCreateSourceHash() => _$lotesControllerHash();

  @override
  String toString() {
    return r'lotesControllerProvider'
        ''
        '($argument)';
  }

  @$internal
  @override
  LotesController create() => LotesController();

  @override
  bool operator ==(Object other) {
    return other is LotesControllerProvider && other.argument == argument;
  }

  @override
  int get hashCode {
    return argument.hashCode;
  }
}

String _$lotesControllerHash() => r'4d1177ad9e3fa5aa8b9f5682ed9571a15b50d687';

/// Lotes de una finca (family por `fincaId`). `autoDispose`: caché por
/// pantalla, sin persistencia.

final class LotesControllerFamily extends $Family
    with
        $ClassFamilyOverride<
          LotesController,
          AsyncValue<List<Lote>>,
          List<Lote>,
          FutureOr<List<Lote>>,
          String
        > {
  const LotesControllerFamily._()
    : super(
        retry: null,
        name: r'lotesControllerProvider',
        dependencies: null,
        $allTransitiveDependencies: null,
        isAutoDispose: true,
      );

  /// Lotes de una finca (family por `fincaId`). `autoDispose`: caché por
  /// pantalla, sin persistencia.

  LotesControllerProvider call(String fincaId) =>
      LotesControllerProvider._(argument: fincaId, from: this);

  @override
  String toString() => r'lotesControllerProvider';
}

/// Lotes de una finca (family por `fincaId`). `autoDispose`: caché por
/// pantalla, sin persistencia.

abstract class _$LotesController extends $AsyncNotifier<List<Lote>> {
  late final _$args = ref.$arg as String;
  String get fincaId => _$args;

  FutureOr<List<Lote>> build(String fincaId);
  @$mustCallSuper
  @override
  void runBuild() {
    final created = build(_$args);
    final ref = this.ref as $Ref<AsyncValue<List<Lote>>, List<Lote>>;
    final element =
        ref.element
            as $ClassProviderElement<
              AnyNotifier<AsyncValue<List<Lote>>, List<Lote>>,
              AsyncValue<List<Lote>>,
              Object?,
              Object?
            >;
    element.handleValue(ref, created);
  }
}
