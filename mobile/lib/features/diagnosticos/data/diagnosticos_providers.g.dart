// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'diagnosticos_providers.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning

@ProviderFor(yoloApi)
const yoloApiProvider = YoloApiProvider._();

final class YoloApiProvider
    extends $FunctionalProvider<YoloApi, YoloApi, YoloApi>
    with $Provider<YoloApi> {
  const YoloApiProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'yoloApiProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$yoloApiHash();

  @$internal
  @override
  $ProviderElement<YoloApi> $createElement($ProviderPointer pointer) =>
      $ProviderElement(pointer);

  @override
  YoloApi create(Ref ref) {
    return yoloApi(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(YoloApi value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<YoloApi>(value),
    );
  }
}

String _$yoloApiHash() => r'e2a82c757aab2db112f353da64b8afd15c6e370c';

@ProviderFor(diagnosticosApi)
const diagnosticosApiProvider = DiagnosticosApiProvider._();

final class DiagnosticosApiProvider
    extends
        $FunctionalProvider<DiagnosticosApi, DiagnosticosApi, DiagnosticosApi>
    with $Provider<DiagnosticosApi> {
  const DiagnosticosApiProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'diagnosticosApiProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$diagnosticosApiHash();

  @$internal
  @override
  $ProviderElement<DiagnosticosApi> $createElement($ProviderPointer pointer) =>
      $ProviderElement(pointer);

  @override
  DiagnosticosApi create(Ref ref) {
    return diagnosticosApi(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(DiagnosticosApi value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<DiagnosticosApi>(value),
    );
  }
}

String _$diagnosticosApiHash() => r'2ac0ee46421bebc8596c7bcb4454d8a285d01437';

/// Se expone como `DiagnosticosRepositoryImpl` (no la interfaz) porque la UI
/// necesita el método extra `analizarBytes` para el caso web/bytes en
/// memoria, que no está en el contrato `DiagnosticosRepository`.

@ProviderFor(diagnosticosRepository)
const diagnosticosRepositoryProvider = DiagnosticosRepositoryProvider._();

/// Se expone como `DiagnosticosRepositoryImpl` (no la interfaz) porque la UI
/// necesita el método extra `analizarBytes` para el caso web/bytes en
/// memoria, que no está en el contrato `DiagnosticosRepository`.

final class DiagnosticosRepositoryProvider
    extends
        $FunctionalProvider<
          DiagnosticosRepositoryImpl,
          DiagnosticosRepositoryImpl,
          DiagnosticosRepositoryImpl
        >
    with $Provider<DiagnosticosRepositoryImpl> {
  /// Se expone como `DiagnosticosRepositoryImpl` (no la interfaz) porque la UI
  /// necesita el método extra `analizarBytes` para el caso web/bytes en
  /// memoria, que no está en el contrato `DiagnosticosRepository`.
  const DiagnosticosRepositoryProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'diagnosticosRepositoryProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$diagnosticosRepositoryHash();

  @$internal
  @override
  $ProviderElement<DiagnosticosRepositoryImpl> $createElement(
    $ProviderPointer pointer,
  ) => $ProviderElement(pointer);

  @override
  DiagnosticosRepositoryImpl create(Ref ref) {
    return diagnosticosRepository(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(DiagnosticosRepositoryImpl value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<DiagnosticosRepositoryImpl>(value),
    );
  }
}

String _$diagnosticosRepositoryHash() =>
    r'9aa541e758c7e049ceea29ee7dbe6d42c0e6c6e6';

/// Alias al contrato para código que solo necesita la interfaz.

@ProviderFor(diagnosticosRepositoryContract)
const diagnosticosRepositoryContractProvider =
    DiagnosticosRepositoryContractProvider._();

/// Alias al contrato para código que solo necesita la interfaz.

final class DiagnosticosRepositoryContractProvider
    extends
        $FunctionalProvider<
          DiagnosticosRepository,
          DiagnosticosRepository,
          DiagnosticosRepository
        >
    with $Provider<DiagnosticosRepository> {
  /// Alias al contrato para código que solo necesita la interfaz.
  const DiagnosticosRepositoryContractProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'diagnosticosRepositoryContractProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$diagnosticosRepositoryContractHash();

  @$internal
  @override
  $ProviderElement<DiagnosticosRepository> $createElement(
    $ProviderPointer pointer,
  ) => $ProviderElement(pointer);

  @override
  DiagnosticosRepository create(Ref ref) {
    return diagnosticosRepositoryContract(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(DiagnosticosRepository value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<DiagnosticosRepository>(value),
    );
  }
}

String _$diagnosticosRepositoryContractHash() =>
    r'3394b9e213b282dca3bc9cd4d0f9561f90a5b3cc';
