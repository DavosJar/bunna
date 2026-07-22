import 'package:dio/dio.dart';

import '../../../core/network/api_envelope.dart';
import '../../../core/network/dio_exception_x.dart';
import 'dtos/mis_tenants_dto.dart';
import 'dtos/perfil_dto.dart';
import 'dtos/permiso_dto.dart';
import 'dtos/token_session_dto.dart';

/// Datasource remoto de auth. Devuelve DTOs (nunca entidades de domain) y
/// lanza `AppException` (nunca `DioException`) — el repositorio es quien
/// traduce DTO → entidad.
///
/// No es `final`: `AuthRepositoryImpl` la consume directo (no hay una
/// interfaz separada para el datasource, solo para el repositorio), y
/// dejarla abierta permite mockearla con mocktail en los tests.
class AuthApi {
  AuthApi(this._dio);

  final Dio _dio;

  static const _base = '/api/v1/identidad';

  Future<TokenSessionDto> login({
    required String correo,
    required String password,
  }) => _post('$_base/auth/login', {'correo': correo, 'password': password})
      .then(TokenSessionDto.fromJson);

  Future<TokenSessionDto> switchTenant(String tenantId) =>
      _post('$_base/auth/switch-tenant', {'tenant_id': tenantId})
          .then(TokenSessionDto.fromJson);

  Future<void> logout() => _postVoid('$_base/auth/logout');

  Future<void> logoutAll() => _postVoid('$_base/auth/logout/all');

  Future<PerfilDto> getMiPerfil() =>
      _get('$_base/mi-perfil').then(PerfilDto.fromJson);

  Future<MisTenantsDto> getMisTenants() =>
      _get('$_base/tenants/mis-tenants').then(MisTenantsDto.fromJson);

  Future<MisPermisosDto> getMisPermisos() =>
      _get('$_base/mis-permisos').then(MisPermisosDto.fromJson);

  Future<Map<String, dynamic>> _get(String path) async {
    try {
      final response = await _dio.get<Map<String, dynamic>>(path);
      return ApiEnvelope.unwrap(
        response.data!,
        (json) => json as Map<String, dynamic>,
      );
    } on DioException catch (e) {
      throw e.toAppException();
    }
  }

  Future<Map<String, dynamic>> _post(
    String path,
    Map<String, dynamic> body,
  ) async {
    try {
      final response = await _dio.post<Map<String, dynamic>>(
        path,
        data: body,
      );
      return ApiEnvelope.unwrap(
        response.data!,
        (json) => json as Map<String, dynamic>,
      );
    } on DioException catch (e) {
      throw e.toAppException();
    }
  }

  Future<void> _postVoid(String path) async {
    try {
      await _dio.post<Map<String, dynamic>>(path);
    } on DioException catch (e) {
      throw e.toAppException();
    }
  }
}
