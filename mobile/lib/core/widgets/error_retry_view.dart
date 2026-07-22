import 'package:flutter/material.dart';

import '../error/app_exception.dart';
import '../error/error_ui.dart';

class ErrorRetryView extends StatelessWidget {
  const ErrorRetryView({super.key, required this.error, this.onRetry});

  final Object error;
  final VoidCallback? onRetry;

  AppException? get _appError => error is AppException ? error as AppException : null;

  String get _message =>
      _appError?.mensajeUsuario ?? 'Ocurrió un error inesperado';

  @override
  Widget build(BuildContext context) {
    final permiso = _appError?.esPermiso ?? false;
    final acento = permiso
        ? Theme.of(context).colorScheme.tertiary
        : Theme.of(context).colorScheme.error;
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              _appError?.iconoUsuario ?? Icons.error_outline,
              size: 40,
              color: acento,
            ),
            const SizedBox(height: 12),
            if (_appError != null) ...[
              Text(
                _appError!.tituloUsuario,
                textAlign: TextAlign.center,
                style: Theme.of(context).textTheme.titleMedium
                    ?.copyWith(fontWeight: FontWeight.w700),
              ),
              const SizedBox(height: 4),
            ],
            Text(
              _message,
              textAlign: TextAlign.center,
              style: Theme.of(context).textTheme.bodyMedium,
            ),
            if (onRetry != null) ...[
              const SizedBox(height: 16),
              OutlinedButton(
                onPressed: onRetry,
                child: const Text('Reintentar'),
              ),
            ],
          ],
        ),
      ),
    );
  }
}
