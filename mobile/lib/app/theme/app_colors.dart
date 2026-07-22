import 'package:flutter/material.dart';

/// Tokens de color EXACTOS del frontend web (`frontend/src/index.css`,
/// sección `:root`). No son aproximaciones: son los mismos hex que usa la app
/// web en producción, para que móvil y web se vean como el mismo producto.
///
/// Estilo real: SaaS agro premium claro — fondo crema, tarjetas blancas con
/// borde gris fino y sombra sutil, verde bosque como acento, títulos serif
/// (Playfair Display) y cuerpo Inter.
abstract final class AppColors {
  // ── Verde (escala completa del frontend) ──────────────────────────────
  static const green900 = Color(0xFF1A3A1A);
  static const green800 = Color(0xFF1E4D2B);
  static const green700 = Color(0xFF2D6A3F);
  static const green600 = Color(0xFF3A7D4F);
  static const green500 = Color(0xFF4A9163);
  static const green400 = Color(0xFF6BB07E);
  static const green300 = Color(0xFF8EC99E);
  static const green200 = Color(0xFFB5DFC2);
  static const green100 = Color(0xFFDDF0E3);
  static const green50 = Color(0xFFF0F9F2);

  // ── Tierra / café ─────────────────────────────────────────────────────
  static const earth900 = Color(0xFF3D2B1F);
  static const earth800 = Color(0xFF5C3D2E);
  static const earth700 = Color(0xFF7A5040);
  static const earth600 = Color(0xFF96644F);
  static const earth500 = Color(0xFFB07D63);
  static const earth100 = Color(0xFFF2E8DE);
  static const earth50 = Color(0xFFFAF5F0);

  // ── Crema / neutros de fondo ──────────────────────────────────────────
  static const cream = Color(0xFFFAF7F2);
  static const creamDark = Color(0xFFF0EBE3);
  static const white = Color(0xFFFFFFFF);

  // ── Grises (escala del frontend) ──────────────────────────────────────
  static const gray50 = Color(0xFFF9FAFB);
  static const gray100 = Color(0xFFF3F4F6);
  static const gray200 = Color(0xFFE5E7EB);
  static const gray300 = Color(0xFFD1D5DB);
  static const gray400 = Color(0xFF9CA3AF);
  static const gray500 = Color(0xFF6B7280);
  static const gray600 = Color(0xFF4B5563);
  static const gray700 = Color(0xFF374151);
  static const gray800 = Color(0xFF1F2937);
  static const gray900 = Color(0xFF111827);

  // ── Ámbar (acento secundario / advertencia) ──────────────────────────
  static const amber500 = Color(0xFFD4A038);
  static const amber600 = Color(0xFFB8891F);

  // ── Semánticos (tomados de Toast.css / Auth.css del frontend) ─────────
  // Error
  static const errorBg = Color(0xFFFEF2F2);
  static const errorBorder = Color(0xFFFECACA);
  static const errorText = Color(0xFF991B1B);
  static const errorIcon = Color(0xFFDC2626);
  // Éxito
  static const successBg = Color(0xFFF0FDF4);
  static const successBorder = Color(0xFFBBF7D0);
  static const successText = Color(0xFF166534);
  // Info
  static const infoBg = Color(0xFFEFF6FF);
  static const infoBorder = Color(0xFFBFDBFE);
  static const infoText = Color(0xFF1E40AF);
  // Advertencia (chip ámbar del frontend: bg #fef3c7 / texto #b45309)
  static const warningBg = Color(0xFFFEF3C7);
  static const warningText = Color(0xFFB45309);

  // ────────────────────────────────────────────────────────────────────
  // Alias semánticos usados por las pantallas ya construidas. Apuntan a los
  // tokens exactos de arriba para que toda la app adopte la paleta real.
  // ────────────────────────────────────────────────────────────────────
  static const verde = green700; // acento primario sólido
  static const verdeVivo = green600;
  static const verdeOscuro = green800; // AppBar/marca, nav activo
  static const verdeOscuroProfundo = green900;
  static const verdeSuave = green100; // fondo de tiles de icono
  static const cafe = earth700;
  static const cafeClaro = earth500;
  static const crema = cream; // fondo de pantallas
  static const cremaClaro = gray50; // fondo de inputs
  static const textoOscuro = gray900; // títulos / texto fuerte
  static const textoTenue = gray500; // texto secundario
  static const blanco = white;
  static const exito = green600;
  static const advertencia = amber500;
  static const error = errorIcon;
  static const superficie = white; // tarjetas
  static const superficieTenue = gray200; // bordes / divisores
}
