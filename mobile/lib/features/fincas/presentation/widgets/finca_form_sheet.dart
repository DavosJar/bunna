import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/error/app_exception.dart';
import '../../../../core/widgets/app_error_banner.dart';
import '../fincas_controller.dart';

/// Bottom sheet para registrar una finca. Devuelve `true` si se creó.
class FincaFormSheet extends ConsumerStatefulWidget {
  const FincaFormSheet({super.key});

  static Future<bool?> show(BuildContext context) => showModalBottomSheet<bool>(
    context: context,
    isScrollControlled: true,
    builder: (_) => const FincaFormSheet(),
  );

  @override
  ConsumerState<FincaFormSheet> createState() => _FincaFormSheetState();
}

class _FincaFormSheetState extends ConsumerState<FincaFormSheet> {
  final _formKey = GlobalKey<FormState>();
  final _nombre = TextEditingController();
  final _ubicacion = TextEditingController();
  final _descripcion = TextEditingController();
  bool _saving = false;
  AppException? _error;

  @override
  void dispose() {
    _nombre.dispose();
    _ubicacion.dispose();
    _descripcion.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;
    setState(() {
      _saving = true;
      _error = null;
    });
    try {
      await ref.read(fincasControllerProvider.notifier).registrar(
        nombre: _nombre.text.trim(),
        ubicacion: _ubicacion.text.trim(),
        descripcion: _descripcion.text.trim(),
      );
      if (mounted) Navigator.of(context).pop(true);
    } on AppException catch (e) {
      setState(() => _error = e);
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final bottomInset = MediaQuery.of(context).viewInsets.bottom;
    return Padding(
      padding: EdgeInsets.fromLTRB(20, 20, 20, 20 + bottomInset),
      child: Form(
        key: _formKey,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text('Nueva finca', style: Theme.of(context).textTheme.titleLarge),
            const SizedBox(height: 16),
            if (_error != null) ...[
              AppErrorBanner(_error!),
              const SizedBox(height: 12),
            ],
            TextFormField(
              controller: _nombre,
              decoration: const InputDecoration(labelText: 'Nombre'),
              textInputAction: TextInputAction.next,
              validator: (v) =>
                  (v == null || v.trim().isEmpty) ? 'Ingresa un nombre' : null,
            ),
            const SizedBox(height: 12),
            TextFormField(
              controller: _ubicacion,
              decoration: const InputDecoration(labelText: 'Ubicación'),
              textInputAction: TextInputAction.next,
              validator: (v) =>
                  (v == null || v.trim().isEmpty) ? 'Ingresa la ubicación' : null,
            ),
            const SizedBox(height: 12),
            TextFormField(
              controller: _descripcion,
              decoration: const InputDecoration(
                labelText: 'Descripción (opcional)',
              ),
              maxLines: 2,
            ),
            const SizedBox(height: 20),
            FilledButton(
              onPressed: _saving ? null : _submit,
              child: _saving
                  ? const SizedBox(
                      height: 20,
                      width: 20,
                      child: CircularProgressIndicator(strokeWidth: 2),
                    )
                  : const Text('Registrar finca'),
            ),
          ],
        ),
      ),
    );
  }
}
