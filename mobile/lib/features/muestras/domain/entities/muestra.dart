import 'package:freezed_annotation/freezed_annotation.dart';

part 'muestra.freezed.dart';

@freezed
abstract class Muestra with _$Muestra {
  const factory Muestra({
    required String id,
    required String fincaId,
    required String loteId,
    required double latitud,
    required double longitud,
    required DateTime createdAt,
  }) = _Muestra;
}
