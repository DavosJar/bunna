import 'dart:convert';

/// Claims relevantes del access token JWT, leídos SIN verificar la firma
/// (solo para uso de UI/estado local — la validación real ocurre en los
/// backends). Equivalente a `parseJWT`/`isTokenExpired` del frontend web.
final class JwtClaims {
  const JwtClaims({
    required this.usuarioId,
    required this.sesionId,
    required this.tenantId,
    required this.rol,
    required this.expiresAt,
  });

  final String usuarioId; // claim `sub`
  final String? sesionId; // claim `sid`
  final String? tenantId; // claim `tenant_id`
  final String? rol; // claim `rol`
  final DateTime? expiresAt; // claim `exp`

  bool get isExpired =>
      expiresAt == null || DateTime.now().isAfter(expiresAt!);

  /// Decodifica el payload de un JWT. Devuelve `null` si el token está mal
  /// formado o no trae `sub` — nunca lanza.
  static JwtClaims? tryDecode(String token) {
    try {
      final parts = token.split('.');
      if (parts.length != 3) return null;

      final payloadJson = utf8.decode(
        base64Url.decode(base64Url.normalize(parts[1])),
      );
      final payload = jsonDecode(payloadJson) as Map<String, dynamic>;

      final sub = payload['sub'] as String?;
      if (sub == null || sub.isEmpty) return null;

      DateTime? expiresAt;
      final exp = payload['exp'];
      if (exp is int) {
        expiresAt = DateTime.fromMillisecondsSinceEpoch(
          exp * 1000,
          isUtc: true,
        ).toLocal();
      }

      return JwtClaims(
        usuarioId: sub,
        sesionId: payload['sid'] as String?,
        tenantId: payload['tenant_id'] as String?,
        rol: payload['rol'] as String?,
        expiresAt: expiresAt,
      );
    } catch (_) {
      return null;
    }
  }
}
