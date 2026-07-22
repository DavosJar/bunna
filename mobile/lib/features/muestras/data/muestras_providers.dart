import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../../../core/network/dio_providers.dart';
import '../domain/muestras_repository.dart';
import 'muestras_api.dart';
import 'muestras_repository_impl.dart';

part 'muestras_providers.g.dart';

@Riverpod(keepAlive: true)
MuestrasApi muestrasApi(Ref ref) => MuestrasApi(ref.watch(fincasDioProvider));

@Riverpod(keepAlive: true)
MuestrasRepository muestrasRepository(Ref ref) =>
    MuestrasRepositoryImpl(ref.watch(muestrasApiProvider));
