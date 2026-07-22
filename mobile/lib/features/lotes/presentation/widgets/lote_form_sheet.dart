import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/error/app_exception.dart';
import '../../../../core/widgets/app_error_banner.dart';
import '../lotes_controller.dart';

class LoteFormSheet extends ConsumerStatefulWidget {
  const LoteFormSheet({super.key, required this.fincaId});

  final String fincaId;

  static Future<bool?> show(BuildContext context, String fincaId) =>
      showModalBottomSheet<bool>(
        context: context,
        isScrollControlled: true,
        builder: (_) => LoteFormSheet(fincaId: fincaId),
      );

  @override
  ConsumerState<LoteFormSheet> createState() => _LoteFormSheetState();
}

class _LoteFormSheetState extends ConsumerState<LoteFormSheet> {
  final _formKey = GlobalKey<FormState>();
  final _nombre = TextEditingController();
  final _area = TextEditingController();
  final _descripcion = TextEditingController();
  bool _saving = false;
  AppException? _error;

  @override
  void dispose() {
    _nombre.dispose();
    _area.dispose();
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
      await ref.read(lotesControllerProvider(widget.fincaId).notifier).agregar(
        nombre: _nombre.text.trim(),
        area: double.parse(_area.text.trim().replaceAll(',', '.')),
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
            Text('Nuevo lote', style: Theme.of(context).textTheme.titleLarge),
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
              controller: _area,
              decoration: const InputDecoration(
                labelText: 'Área (hectáreas)',
                suffixText: 'ha',
              ),
              keyboardType: const TextInputType.numberWithOptions(decimal: true),
              textInputAction: TextInputAction.next,
              validator: (v) {
                if (v == null || v.trim().isEmpty) return 'Ingresa el área';
                final parsed = double.tryParse(v.trim().replaceAll(',', '.'));
                if (parsed == null || parsed <= 0) return 'Área inválida';
                return null;
              },
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
                  : const Text('Agregar lote'),
            ),
          ],
        ),
      ),
    );
  }
}
