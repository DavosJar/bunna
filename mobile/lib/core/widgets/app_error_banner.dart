import 'package:flutter/material.dart';

import '../../app/theme/app_colors.dart';
import '../error/app_exception.dart';
import '../error/error_ui.dart';

/// Banner de alerta para mostrar un [AppException] dentro de un formulario o
/// pantalla de acción. Replica el estilo de `.auth-error` / `.toast--error`
/// del frontend web (bloque suave con borde e icono).
///
/// Distingue el caso de permiso (403) — que no es un fallo del sistema sino
/// una restricción esperada — con el tono ámbar del frontend (`#fef3c7` /
/// `#b45309`), en vez del rojo de error.
class AppErrorBanner extends StatelessWidget {
  const AppErrorBanner(this.error, {super.key});

  final AppException error;

  @override
  Widget build(BuildContext context) {
    final permiso = error.esPermiso;

    final bg = permiso ? AppColors.warningBg : AppColors.errorBg;
    final border = permiso
        ? AppColors.amber500.withValues(alpha: 0.45)
        : AppColors.errorBorder;
    final titleColor = permiso ? AppColors.warningText : AppColors.errorText;
    final iconColor = permiso ? AppColors.amber600 : AppColors.errorIcon;

    return Container(
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        color: bg,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: border),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(error.iconoUsuario, color: iconColor, size: 20),
          const SizedBox(width: 11),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  error.tituloUsuario,
                  style: TextStyle(
                    fontWeight: FontWeight.w700,
                    fontSize: 14,
                    color: titleColor,
                  ),
                ),
                const SizedBox(height: 3),
                Text(
                  error.mensajeUsuario,
                  style: TextStyle(
                    color: titleColor.withValues(alpha: 0.85),
                    fontSize: 13,
                    height: 1.35,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
