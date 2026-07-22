// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'lotes_providers.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning

@ProviderFor(lotesApi)
const lotesApiProvider = LotesApiProvider._();

final class LotesApiProvider
    extends $FunctionalProvider<LotesApi, LotesApi, LotesApi>
    with $Provider<LotesApi> {
  const LotesApiProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'lotesApiProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$lotesApiHash();

  @$internal
  @override
  $ProviderElement<LotesApi> $createElement($ProviderPointer pointer) =>
      $ProviderElement(pointer);

  @override
  LotesApi create(Ref ref) {
    return lotesApi(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(LotesApi value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<LotesApi>(value),
    );
  }
}

String _$lotesApiHash() => r'96d645aaed1c777d061c0e9232210104dfed4a1b';

@ProviderFor(lotesRepository)
const lotesRepositoryProvider = LotesRepositoryProvider._();

final class LotesRepositoryProvider
    extends
        $FunctionalProvider<LotesRepository, LotesRepository, LotesRepository>
    with $Provider<LotesRepository> {
  const LotesRepositoryProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'lotesRepositoryProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$lotesRepositoryHash();

  @$internal
  @override
  $ProviderElement<LotesRepository> $createElement($ProviderPointer pointer) =>
      $ProviderElement(pointer);

  @override
  LotesRepository create(Ref ref) {
    return lotesRepository(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(LotesRepository value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<LotesRepository>(value),
    );
  }
}

String _$lotesRepositoryHash() => r'4d2b00eeb3f8af0d21ec0db4262515a11756a6ed';
