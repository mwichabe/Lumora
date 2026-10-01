import 'package:flutter/material.dart';
import '../core/theme/colors.dart';
import '../core/theme/radii.dart';

enum LumoraButtonVariant { primary, secondary, outline, ghost, danger }

/// Pill buttons matching frontend/components/Button.tsx's five variants.
class LumoraButton extends StatelessWidget {
  final String label;
  final VoidCallback? onPressed;
  final LumoraButtonVariant variant;
  final bool full;
  final bool loading;
  final Widget? icon;

  const LumoraButton({
    super.key,
    required this.label,
    required this.onPressed,
    this.variant = LumoraButtonVariant.primary,
    this.full = false,
    this.loading = false,
    this.icon,
  });

  @override
  Widget build(BuildContext context) {
    final disabled = onPressed == null || loading;
    final child = loading
        ? SizedBox(
            height: 20,
            width: 20,
            child: CircularProgressIndicator(
              strokeWidth: 2.5,
              valueColor: AlwaysStoppedAnimation(_fg().withValues(alpha: 0.9)),
            ),
          )
        : Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              if (icon != null) ...[icon!, const SizedBox(width: 8)],
              // Shrinks with an ellipsis rather than overflowing the pill when
              // a label is long for the screen (narrow phones, large text).
              Flexible(
                child: Text(label, maxLines: 1, overflow: TextOverflow.ellipsis,
                    style: TextStyle(fontWeight: FontWeight.w800, fontSize: 16, color: _fg())),
              ),
            ],
          );

    final button = AnimatedScale(
      duration: const Duration(milliseconds: 120),
      scale: 1.0,
      child: Material(
        color: _bg(),
        shape: variant == LumoraButtonVariant.ghost
            ? RoundedRectangleBorder(borderRadius: BorderRadius.circular(LumoraRadii.md))
            : StadiumBorder(
                side: variant == LumoraButtonVariant.outline
                    ? const BorderSide(color: LumoraColors.purple, width: 2)
                    : BorderSide.none,
              ),
        child: InkWell(
          onTap: disabled ? null : onPressed,
          customBorder: variant == LumoraButtonVariant.ghost
              ? RoundedRectangleBorder(borderRadius: BorderRadius.circular(LumoraRadii.md))
              : const StadiumBorder(),
          child: Container(
            height: variant == LumoraButtonVariant.ghost ? 44 : 52,
            padding: const EdgeInsets.symmetric(horizontal: 24),
            alignment: Alignment.center,
            width: full ? double.infinity : null,
            child: Opacity(opacity: disabled && !loading ? 0.4 : 1, child: child),
          ),
        ),
      ),
    );

    return full ? SizedBox(width: double.infinity, child: button) : button;
  }

  Color _bg() {
    switch (variant) {
      case LumoraButtonVariant.primary:
        return LumoraColors.purple;
      case LumoraButtonVariant.secondary:
        return LumoraColors.amber;
      case LumoraButtonVariant.outline:
        return Colors.transparent;
      case LumoraButtonVariant.ghost:
        return Colors.transparent;
      case LumoraButtonVariant.danger:
        return LumoraColors.coral;
    }
  }

  Color _fg() {
    switch (variant) {
      case LumoraButtonVariant.primary:
      case LumoraButtonVariant.danger:
        return Colors.white;
      case LumoraButtonVariant.secondary:
        return LumoraColors.ink;
      case LumoraButtonVariant.outline:
      case LumoraButtonVariant.ghost:
        return LumoraColors.purple;
    }
  }
}
