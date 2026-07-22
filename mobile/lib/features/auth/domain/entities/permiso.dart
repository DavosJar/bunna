import 'package:freezed_annotation/freezed_annotation.dart';

part 'permiso.freezed.dart';

@freezed
abstract class Permiso with _$Permiso {
  const factory Permiso({
    required String codigo,
    required String nombre,
    required String descripcion,
    required String modulo,
  }) = _Permiso;
}
