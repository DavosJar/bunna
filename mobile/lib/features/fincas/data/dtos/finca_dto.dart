import 'package:freezed_annotation/freezed_annotation.dart';

import '../../domain/entities/finca.dart';

part 'finca_dto.freezed.dart';
part 'finca_dto.g.dart';

/// `FincaResponse` de fincas (familia camelCase; `createdAt` coincide con
/// lowerCamel Dart). Cotejado contra
/// `fincas/internal/presentation/dto/finca.go`.
@freezed
abstract class FincaDto with _$FincaDto {
  const factory FincaDto({
    required String id,
    required String nombre,
    required String ubicacion,
    required String descripcion,
    required String estado,
    required DateTime createdAt,
  }) = _FincaDto;

  const FincaDto._();

  factory FincaDto.fromJson(Map<String, dynamic> json) =>
      _$FincaDtoFromJson(json);

  Finca toDomain() => Finca(
    id: id,
    nombre: nombre,
    ubicacion: ubicacion,
    descripcion: descripcion,
    estado: estado,
    createdAt: createdAt,
  );
}
