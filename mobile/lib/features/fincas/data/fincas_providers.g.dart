// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'fincas_providers.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning

@ProviderFor(fincasApi)
const fincasApiProvider = FincasApiProvider._();

final class FincasApiProvider
    extends $FunctionalProvider<FincasApi, FincasApi, FincasApi>
    with $Provider<FincasApi> {
  const FincasApiProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'fincasApiProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$fincasApiHash();

  @$internal
  @override
  $ProviderElement<FincasApi> $createElement($ProviderPointer pointer) =>
      $ProviderElement(pointer);

  @override
  FincasApi create(Ref ref) {
    return fincasApi(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(FincasApi value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<FincasApi>(value),
    );
  }
}

String _$fincasApiHash() => r'14f9ec9cc5f55daceecfd888f2888389a2a1484c';

@ProviderFor(fincasRepository)
const fincasRepositoryProvider = FincasRepositoryProvider._();

final class FincasRepositoryProvider
    extends
        $FunctionalProvider<
          FincasRepository,
          FincasRepository,
          FincasRepository
        >
    with $Provider<FincasRepository> {
  const FincasRepositoryProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'fincasRepositoryProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$fincasRepositoryHash();

  @$internal
  @override
  $ProviderElement<FincasRepository> $createElement($ProviderPointer pointer) =>
      $ProviderElement(pointer);

  @override
  FincasRepository create(Ref ref) {
    return fincasRepository(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(FincasRepository value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<FincasRepository>(value),
    );
  }
}

String _$fincasRepositoryHash() => r'19ee33b009fe4df01df4a2d6f39b893391b368d0';
