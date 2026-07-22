import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/core/data/cambio_estado_dto.dart';
import 'package:mobile/features/diagnosticos/data/dtos/diagnostico_dto.dart';
import 'package:mobile/features/fincas/data/dtos/finca_dto.dart';
import 'package:mobile/features/lotes/data/dtos/lote_dto.dart';
import 'package:mobile/features/muestras/data/dtos/muestra_dto.dart';

/// Verifica el parsing de los DTOs de fincas con las formas JSON EXACTAS de
/// los structs Go (`fincas/internal/presentation/dto/*.go`) — el naming mixto
/// (sufijos ID/URL en mayúscula vs camelCase) es la parte de mayor riesgo.
void main() {
  test('FincaDto: camelCase + createdAt', () {
    final dto = FincaDto.fromJson(const {
      'id': 'f1',
      'nombre': 'Finca El Cafetal',
      'ubicacion': 'Loja',
      'descripcion': 'demo',
      'estado': 'ACTIVA',
      'createdAt': '2026-07-22T03:14:42Z',
    });
    final finca = dto.toDomain();
    expect(finca.nombre, 'Finca El Cafetal');
    expect(finca.estaActiva, isTrue);
    expect(finca.createdAt.toUtc().year, 2026);
  });

  test('LoteDto: fincaID → fincaId', () {
    final dto = LoteDto.fromJson(const {
      'id': 'l1',
      'fincaID': 'f1',
      'nombre': 'Lote Norte',
      'area': 2.5,
      'descripcion': '',
      'estado': 'ACTIVO',
      'createdAt': '2026-07-22T03:14:42Z',
    });
    final lote = dto.toDomain();
    expect(lote.fincaId, 'f1');
    expect(lote.area, 2.5);
    expect(lote.estaActivo, isTrue);
  });

  test('LoteDto: area entera se parsea como double', () {
    final dto = LoteDto.fromJson(const {
      'id': 'l1',
      'fincaID': 'f1',
      'nombre': 'Lote',
      'area': 3, // el backend puede serializar sin decimales
      'descripcion': '',
      'estado': 'ACTIVO',
      'createdAt': '2026-07-22T03:14:42Z',
    });
    expect(dto.toDomain().area, 3.0);
  });

  test('MuestraDto: fincaID/loteID → fincaId/loteId', () {
    final dto = MuestraDto.fromJson(const {
      'id': 'm1',
      'fincaID': 'f1',
      'loteID': 'l1',
      'latitud': -3.99313,
      'longitud': -79.20422,
      'createdAt': '2026-07-22T03:14:42Z',
    });
    final m = dto.toDomain();
    expect(m.fincaId, 'f1');
    expect(m.loteId, 'l1');
    expect(m.latitud, closeTo(-3.99313, 1e-9));
  });

  test('DiagnosticoDto: muestraID/imageURL → muestraId/imageUrl', () {
    final dto = DiagnosticoDto.fromJson(const {
      'id': 'd1',
      'muestraID': 'm1',
      'nombre': 'diag',
      'estado': 'PENDIENTE',
      'tieneClorosis': true,
      'confianza': 0.87,
      'imageURL': 'data:image/jpeg;base64,abc',
      'imageBase64': '',
      'procesadoAt': '2026-07-22T03:14:42Z',
      'createdAt': '2026-07-22T03:14:42Z',
    });
    final d = dto.toDomain();
    expect(d.muestraId, 'm1');
    expect(d.tieneClorosis, isTrue);
    expect(d.confianza, 0.87);
    expect(d.imageUrl, 'data:image/jpeg;base64,abc');
  });

  test('CambioEstadoDto: motivo opcional ausente', () {
    final dto = CambioEstadoDto.fromJson(const {
      'id': 'x1',
      'estado': 'INACTIVA',
      'updatedAt': '2026-07-22T03:14:42Z',
    });
    final c = dto.toDomain();
    expect(c.estado, 'INACTIVA');
    expect(c.motivo, isNull);
  });
}
