import 'package:freezed_annotation/freezed_annotation.dart';

part 'perfil.freezed.dart';

@freezed
abstract class Perfil with _$Perfil {
  const factory Perfil({
    required String id,
    required String correo,
    required String nombre,
    required String apellido,
    required String telefono,
    required String estado,
    required DateTime creadoEn,
  }) = _Perfil;
}
