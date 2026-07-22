import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../../../core/network/dio_providers.dart';
import '../domain/fincas_repository.dart';
import 'fincas_api.dart';
import 'fincas_repository_impl.dart';

part 'fincas_providers.g.dart';

@Riverpod(keepAlive: true)
FincasApi fincasApi(Ref ref) => FincasApi(ref.watch(fincasDioProvider));

@Riverpod(keepAlive: true)
FincasRepository fincasRepository(Ref ref) =>
    FincasRepositoryImpl(ref.watch(fincasApiProvider));
