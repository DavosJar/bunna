import 'package:flutter/material.dart';

import '../../app/theme/app_colors.dart';

/// Tile de icono redondeado — replica el `.stat-card__icon` del frontend web:
/// cuadrado con esquinas redondeadas (no círculo), fondo verde-100 y glifo
/// verde-700. Es el patrón visual recurrente en tarjetas y encabezados.
///
/// Conserva el nombre `CircleIcon` por compatibilidad con las pantallas ya
/// construidas, pero el estilo real de marca es el tile redondeado.
class CircleIcon extends StatelessWidget {
  const CircleIcon(
    this.icon, {
    super.key,
    this.size = 44,
    this.background = AppColors.green100,
    this.foreground = AppColors.green700,
  });

  final IconData icon;
  final double size;
  final Color background;
  final Color foreground;

  @override
  Widget build(BuildContext context) {
    // Radio proporcional (~0.27x) equivalente a --radius-lg sobre el tile de
    // 44px del frontend, escalado para tamaños mayores.
    final radius = (size * 0.27).clamp(10.0, 18.0);
    return Container(
      width: size,
      height: size,
      decoration: BoxDecoration(
        color: background,
        borderRadius: BorderRadius.circular(radius),
      ),
      child: Icon(icon, color: foreground, size: size * 0.5),
    );
  }
}
