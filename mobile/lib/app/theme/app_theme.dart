import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

import 'app_colors.dart';

/// Tema de CafeScan replicado EXACTAMENTE del frontend web
/// (`frontend/src/index.css` + hojas de componentes). Look real: SaaS agro
/// premium claro — fondo crema, tarjetas blancas con borde gris fino, verde
/// bosque como acento, títulos en serif (Playfair Display) y cuerpo en Inter.
///
/// Nota: las fuentes se cargan vía `google_fonts` (mismo par que la web:
/// Inter + Playfair Display). Si no hay red, degradan al sans del sistema sin
/// romper el layout.
abstract final class AppTheme {
  // Radios exactos del frontend (--radius-*).
  static const double _radiusLg = 12; // 0.75rem
  static const double _radiusXl = 16; // 1rem
  static const double _radius2xl = 24; // 1.5rem

  static ThemeData get light {
    final colorScheme = const ColorScheme.light(
      primary: AppColors.green700,
      onPrimary: AppColors.white,
      primaryContainer: AppColors.green100,
      onPrimaryContainer: AppColors.green800,
      secondary: AppColors.earth700,
      onSecondary: AppColors.white,
      secondaryContainer: AppColors.earth100,
      onSecondaryContainer: AppColors.earth800,
      surface: AppColors.white,
      onSurface: AppColors.gray900,
      surfaceContainerHighest: AppColors.gray100,
      onSurfaceVariant: AppColors.gray600,
      outline: AppColors.gray300,
      outlineVariant: AppColors.gray200,
      error: AppColors.errorIcon,
      onError: AppColors.white,
      errorContainer: AppColors.errorBg,
      onErrorContainer: AppColors.errorText,
    );

    final base = ThemeData(
      useMaterial3: true,
      colorScheme: colorScheme,
      scaffoldBackgroundColor: AppColors.cream,
      splashFactory: InkRipple.splashFactory,
    );

    return base.copyWith(
      textTheme: _buildTextTheme(base.textTheme),

      appBarTheme: const AppBarTheme(
        backgroundColor: AppColors.white,
        foregroundColor: AppColors.gray900,
        surfaceTintColor: Colors.transparent,
        elevation: 0,
        scrolledUnderElevation: 0,
        centerTitle: false,
        // Hairline inferior gris (como el borde del sidebar/navbar web).
        shape: Border(
          bottom: BorderSide(color: AppColors.gray200),
        ),
        iconTheme: IconThemeData(color: AppColors.gray700),
      ),

      // Tarjetas: blanco, borde gris fino, radio xl, sin sombra dura
      // (el frontend usa shadow-sm apenas perceptible → borde + plano).
      cardTheme: CardThemeData(
        color: AppColors.white,
        elevation: 0,
        shadowColor: Colors.black.withValues(alpha: 0.04),
        surfaceTintColor: Colors.transparent,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(_radiusXl),
          side: const BorderSide(color: AppColors.gray200),
        ),
        margin: EdgeInsets.zero,
        clipBehavior: Clip.antiAlias,
      ),

      // Botón primario: verde sólido (el gradiente 800→700 de la web se lee
      // casi igual a este verde a escala de botón), radio xl, peso 600.
      filledButtonTheme: FilledButtonThemeData(
        style: FilledButton.styleFrom(
          backgroundColor: AppColors.green700,
          foregroundColor: AppColors.white,
          disabledBackgroundColor: AppColors.green700.withValues(alpha: 0.5),
          disabledForegroundColor: AppColors.white,
          minimumSize: const Size.fromHeight(50),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(_radiusXl),
          ),
          textStyle: GoogleFonts.inter(
            fontSize: 15,
            fontWeight: FontWeight.w600,
            letterSpacing: 0.3,
          ),
        ),
      ),

      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: AppColors.green700,
          foregroundColor: AppColors.white,
          elevation: 0,
          minimumSize: const Size.fromHeight(50),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(_radiusXl),
          ),
          textStyle: GoogleFonts.inter(
            fontSize: 15,
            fontWeight: FontWeight.w600,
            letterSpacing: 0.3,
          ),
        ),
      ),

      outlinedButtonTheme: OutlinedButtonThemeData(
        style: OutlinedButton.styleFrom(
          foregroundColor: AppColors.green700,
          side: const BorderSide(color: AppColors.gray300, width: 1.5),
          minimumSize: const Size.fromHeight(50),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(_radiusXl),
          ),
          textStyle: GoogleFonts.inter(
            fontSize: 15,
            fontWeight: FontWeight.w600,
            letterSpacing: 0.3,
          ),
        ),
      ),

      textButtonTheme: TextButtonThemeData(
        style: TextButton.styleFrom(
          foregroundColor: AppColors.green700,
          textStyle: GoogleFonts.inter(
            fontSize: 14,
            fontWeight: FontWeight.w600,
          ),
        ),
      ),

      floatingActionButtonTheme: FloatingActionButtonThemeData(
        backgroundColor: AppColors.green700,
        foregroundColor: AppColors.white,
        elevation: 2,
        highlightElevation: 4,
        extendedTextStyle: GoogleFonts.inter(
          fontSize: 15,
          fontWeight: FontWeight.w600,
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(_radiusXl),
        ),
      ),

      // Inputs: fondo gray-50, borde 1.5px gray-200, radio lg, foco verde-500.
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: AppColors.gray50,
        contentPadding:
            const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
        hintStyle: const TextStyle(color: AppColors.gray400),
        labelStyle: const TextStyle(color: AppColors.gray600),
        floatingLabelStyle: const TextStyle(
          color: AppColors.green600,
          fontWeight: FontWeight.w600,
        ),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(_radiusLg),
          borderSide: const BorderSide(color: AppColors.gray200, width: 1.5),
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(_radiusLg),
          borderSide: const BorderSide(color: AppColors.gray200, width: 1.5),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(_radiusLg),
          borderSide: const BorderSide(color: AppColors.green500, width: 1.5),
        ),
        errorBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(_radiusLg),
          borderSide: const BorderSide(color: AppColors.errorBorder, width: 1.5),
        ),
        focusedErrorBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(_radiusLg),
          borderSide: const BorderSide(color: AppColors.errorIcon, width: 1.5),
        ),
        errorStyle: const TextStyle(
          color: AppColors.errorIcon,
          fontWeight: FontWeight.w500,
        ),
      ),

      // Chips / badges de estado: pill verde-100 con texto verde-800.
      chipTheme: base.chipTheme.copyWith(
        backgroundColor: AppColors.green100,
        labelStyle: GoogleFonts.inter(
          color: AppColors.green800,
          fontSize: 12,
          fontWeight: FontWeight.w600,
        ),
        shape: const StadiumBorder(),
        side: BorderSide.none,
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      ),

      dividerTheme: const DividerThemeData(
        color: AppColors.gray200,
        space: 1,
        thickness: 1,
      ),

      snackBarTheme: SnackBarThemeData(
        behavior: SnackBarBehavior.floating,
        backgroundColor: AppColors.green800,
        contentTextStyle: GoogleFonts.inter(
          color: AppColors.white,
          fontSize: 14,
          fontWeight: FontWeight.w500,
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(_radiusLg),
        ),
      ),

      dialogTheme: DialogThemeData(
        backgroundColor: AppColors.white,
        surfaceTintColor: Colors.transparent,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(_radius2xl),
        ),
        titleTextStyle: GoogleFonts.playfairDisplay(
          fontSize: 20,
          fontWeight: FontWeight.w700,
          color: AppColors.gray900,
        ),
      ),

      bottomSheetTheme: const BottomSheetThemeData(
        backgroundColor: AppColors.white,
        surfaceTintColor: Colors.transparent,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.vertical(top: Radius.circular(_radius2xl)),
        ),
      ),

      progressIndicatorTheme: const ProgressIndicatorThemeData(
        color: AppColors.green600,
      ),

      listTileTheme: const ListTileThemeData(
        iconColor: AppColors.gray500,
        textColor: AppColors.gray900,
      ),
    );
  }

  /// Tipografía exacta de la web: Playfair Display (serif) para títulos y
  /// display; Inter para todo lo demás. Colores gray-900 (fuerte) / gray-800
  /// (cuerpo).
  static TextTheme _buildTextTheme(TextTheme base) {
    final inter = GoogleFonts.interTextTheme(base);
    TextStyle serif(double size, {FontWeight weight = FontWeight.w700}) =>
        GoogleFonts.playfairDisplay(
          fontSize: size,
          fontWeight: weight,
          color: AppColors.gray900,
          height: 1.2,
        );

    return inter.copyWith(
      // Display + headline + titleLarge → serif (títulos de página / tarjeta).
      displayLarge: serif(40, weight: FontWeight.w800),
      displayMedium: serif(32, weight: FontWeight.w800),
      displaySmall: serif(28),
      headlineLarge: serif(26),
      headlineMedium: serif(24),
      headlineSmall: serif(20),
      titleLarge: serif(20),
      // titleMedium/Small y cuerpo → Inter.
      titleMedium: GoogleFonts.inter(
        fontSize: 16,
        fontWeight: FontWeight.w600,
        color: AppColors.gray900,
        height: 1.3,
      ),
      titleSmall: GoogleFonts.inter(
        fontSize: 14,
        fontWeight: FontWeight.w600,
        color: AppColors.gray800,
      ),
      bodyLarge: GoogleFonts.inter(
        fontSize: 15,
        color: AppColors.gray800,
        height: 1.5,
      ),
      bodyMedium: GoogleFonts.inter(
        fontSize: 14,
        color: AppColors.gray700,
        height: 1.5,
      ),
      bodySmall: GoogleFonts.inter(
        fontSize: 12.5,
        color: AppColors.gray500,
        height: 1.4,
      ),
      labelLarge: GoogleFonts.inter(
        fontSize: 14,
        fontWeight: FontWeight.w600,
        color: AppColors.gray900,
      ),
    );
  }
}
