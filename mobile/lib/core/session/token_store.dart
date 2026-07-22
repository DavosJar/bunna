import 'dart:convert';

import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import 'auth_tokens.dart';

/// Fuente única de verdad de los tokens: copia en memoria (lectura síncrona,
/// usada por los interceptores en cada request sin I/O de disco) respaldada
/// por Flutter Secure Storage para sobrevivir el reinicio de la app.
final class TokenStore {
  TokenStore({FlutterSecureStorage? storage})
    : _storage = storage ?? const FlutterSecureStorage();

  static const _storageKey = 'bunna_auth_tokens';

  final FlutterSecureStorage _storage;
  AuthTokens? _current;

  /// Lectura síncrona desde memoria — segura para usar dentro de
  /// interceptores de Dio.
  AuthTokens? get current => _current;

  /// Carga los tokens persistidos a memoria. Debe llamarse una vez al
  /// arrancar la app, antes de cualquier request.
  Future<void> loadFromDisk() async {
    final raw = await _storage.read(key: _storageKey);
    if (raw == null) {
      _current = null;
      return;
    }
    try {
      _current = AuthTokens.fromJson(
        jsonDecode(raw) as Map<String, dynamic>,
      );
    } catch (_) {
      // Payload corrupto o de un formato anterior — tratar como sesión nueva.
      _current = null;
      await _storage.delete(key: _storageKey);
    }
  }

  /// Guarda AMBOS tokens de forma atómica: primero memoria (lectura
  /// inmediata para el siguiente request), luego disco.
  Future<void> save(AuthTokens tokens) async {
    _current = tokens;
    await _storage.write(
      key: _storageKey,
      value: jsonEncode(tokens.toJson()),
    );
  }

  Future<void> clear() async {
    _current = null;
    await _storage.delete(key: _storageKey);
  }
}
