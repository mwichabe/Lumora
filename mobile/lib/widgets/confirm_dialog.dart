import 'package:flutter/material.dart';

import '../core/theme/colors.dart';
import '../core/theme/radii.dart';
import '../core/theme/shadows.dart';
import 'lumora_button.dart';

/// A small centered confirmation modal, matching
/// frontend/components/ConfirmDialog.tsx — on-brand rounded card with a
/// title, optional message, and two full-width buttons (outline cancel +
/// primary/danger confirm). Use this instead of the bare Material
/// [AlertDialog] everywhere a yes/no confirmation is needed.
Future<bool> showLumoraConfirmDialog(
  BuildContext context, {
  required String title,
  String? message,
  String confirmLabel = 'Confirm',
  String cancelLabel = 'Cancel',
  bool danger = false,
}) async {
  final result = await showDialog<bool>(
    context: context,
    barrierColor: Colors.black.withValues(alpha: 0.5),
    builder: (context) => _LumoraConfirmDialog(
      title: title,
      message: message,
      confirmLabel: confirmLabel,
      cancelLabel: cancelLabel,
      danger: danger,
    ),
  );
  return result ?? false;
}

class _LumoraConfirmDialog extends StatelessWidget {
  final String title;
  final String? message;
  final String confirmLabel;
  final String cancelLabel;
  final bool danger;

  const _LumoraConfirmDialog({
    required this.title,
    required this.message,
    required this.confirmLabel,
    required this.cancelLabel,
    required this.danger,
  });

  @override
  Widget build(BuildContext context) {
    return Dialog(
      backgroundColor: Colors.transparent,
      insetPadding: const EdgeInsets.all(24),
      child: ConstrainedBox(
        constraints: const BoxConstraints(maxWidth: 320),
        child: Container(
          padding: const EdgeInsets.all(20),
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.circular(LumoraRadii.xl),
            boxShadow: LumoraShadows.cardLg,
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(
                title,
                textAlign: TextAlign.center,
                style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w800, color: LumoraColors.ink),
              ),
              if (message != null) ...[
                const SizedBox(height: 6),
                Text(
                  message!,
                  textAlign: TextAlign.center,
                  style: const TextStyle(fontSize: 14, color: LumoraColors.slatey, height: 1.4),
                ),
              ],
              const SizedBox(height: 20),
              Row(
                children: [
                  Expanded(
                    child: LumoraButton(
                      label: cancelLabel,
                      variant: LumoraButtonVariant.outline,
                      full: true,
                      onPressed: () => Navigator.of(context).pop(false),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: LumoraButton(
                      label: confirmLabel,
                      variant: danger ? LumoraButtonVariant.danger : LumoraButtonVariant.primary,
                      full: true,
                      onPressed: () => Navigator.of(context).pop(true),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}
