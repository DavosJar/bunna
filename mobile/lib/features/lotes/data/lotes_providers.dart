import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../../../core/network/dio_providers.dart';
import '../domain/lotes_repository.dart';
import 'lotes_api.dart';
import 'lotes_repository_impl.dart';

part 'lotes_providers.g.dart';

@Riverpod(keepAlive: true)
LotesApi lotesApi(Ref ref) => LotesApi(ref.watch(fincasDioProvider));

@Riverpod(keepAlive: true)
LotesRepository lotesRepository(Ref ref) =>
    LotesRepositoryImpl(ref.watch(lotesApiProvider));
