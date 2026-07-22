import 'package:freezed_annotation/freezed_annotation.dart';

part 'finca.freezed.dart';

@freezed
abstract class Finca with _$Finca {
  const factory Finca({
    required String id,
    required String nombre,
    required String ubicacion,
    required String descripcion,
    required String estado,
    required DateTime createdAt,
  }) = _Finca;

  const Finca._();

  bool get estaActiva => estado.toUpperCase() == 'ACTIVA';
}
