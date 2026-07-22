import 'package:freezed_annotation/freezed_annotation.dart';

part 'lote.freezed.dart';

@freezed
abstract class Lote with _$Lote {
  const factory Lote({
    required String id,
    required String fincaId,
    required String nombre,
    required double area,
    required String descripcion,
    required String estado,
    required DateTime createdAt,
  }) = _Lote;

  const Lote._();

  bool get estaActivo => estado.toUpperCase() == 'ACTIVO';
}
