// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'muestras_providers.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint, type=warning

@ProviderFor(muestrasApi)
const muestrasApiProvider = MuestrasApiProvider._();

final class MuestrasApiProvider
    extends $FunctionalProvider<MuestrasApi, MuestrasApi, MuestrasApi>
    with $Provider<MuestrasApi> {
  const MuestrasApiProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'muestrasApiProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$muestrasApiHash();

  @$internal
  @override
  $ProviderElement<MuestrasApi> $createElement($ProviderPointer pointer) =>
      $ProviderElement(pointer);

  @override
  MuestrasApi create(Ref ref) {
    return muestrasApi(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(MuestrasApi value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<MuestrasApi>(value),
    );
  }
}

String _$muestrasApiHash() => r'700453f449370342ebcf8d48a4aa7aaffc67c569';

@ProviderFor(muestrasRepository)
const muestrasRepositoryProvider = MuestrasRepositoryProvider._();

final class MuestrasRepositoryProvider
    extends
        $FunctionalProvider<
          MuestrasRepository,
          MuestrasRepository,
          MuestrasRepository
        >
    with $Provider<MuestrasRepository> {
  const MuestrasRepositoryProvider._()
    : super(
        from: null,
        argument: null,
        retry: null,
        name: r'muestrasRepositoryProvider',
        isAutoDispose: false,
        dependencies: null,
        $allTransitiveDependencies: null,
      );

  @override
  String debugGetCreateSourceHash() => _$muestrasRepositoryHash();

  @$internal
  @override
  $ProviderElement<MuestrasRepository> $createElement(
    $ProviderPointer pointer,
  ) => $ProviderElement(pointer);

  @override
  MuestrasRepository create(Ref ref) {
    return muestrasRepository(ref);
  }

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(MuestrasRepository value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<MuestrasRepository>(value),
    );
  }
}

String _$muestrasRepositoryHash() =>
    r'9f4f1c8833515f58a4c7e985dd563b8c436795bd';
