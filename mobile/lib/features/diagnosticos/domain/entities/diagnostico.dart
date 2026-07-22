import 'package:freezed_annotation/freezed_annotation.dart';

part 'diagnostico.freezed.dart';

@freezed
abstract class Diagnostico with _$Diagnostico {
  const factory Diagnostico({
    required String id,
    required String muestraId,
    required String estado,
    required bool tieneClorosis,
    required double confianza,
    String? imageUrl,
    String? imageBase64,
  }) = _Diagnostico;
}
