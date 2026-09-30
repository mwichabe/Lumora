import 'package:flutter/material.dart';

import '../core/theme/colors.dart';
import '../core/theme/radii.dart';
import '../core/theme/shadows.dart';

/// Shared visual language for the auth flow (welcome, forgot/reset password):
/// a light cream background, a white rounded card, and pill inputs with a
/// leading icon — no gradients, purple only as an accent.
InputDecoration authFieldDecoration({required IconData icon, required String hint, Widget? suffixIcon}) {
  return InputDecoration(
    hintText: hint,
    prefixIcon: Icon(icon, size: 20, color: LumoraColors.gray500),
    suffixIcon: suffixIcon,
    filled: true,
    fillColor: LumoraColors.gray50,
    contentPadding: const EdgeInsets.symmetric(vertical: 14),
    border: OutlineInputBorder(borderRadius: BorderRadius.circular(LumoraRadii.lg), borderSide: BorderSide.none),
    enabledBorder: OutlineInputBorder(borderRadius: BorderRadius.circular(LumoraRadii.lg), borderSide: BorderSide.none),
    focusedBorder: OutlineInputBorder(borderRadius: BorderRadius.circular(LumoraRadii.lg), borderSide: const BorderSide(color: LumoraColors.purple, width: 2)),
    errorBorder: OutlineInputBorder(borderRadius: BorderRadius.circular(LumoraRadii.lg), borderSide: const BorderSide(color: LumoraColors.coral, width: 1.5)),
    focusedErrorBorder: OutlineInputBorder(borderRadius: BorderRadius.circular(LumoraRadii.lg), borderSide: const BorderSide(color: LumoraColors.coral, width: 2)),
  );
}

/// The rounded white card auth screens sit their content in.
class AuthCard extends StatelessWidget {
  final Widget child;
  const AuthCard({super.key, required this.child});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(28),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(LumoraRadii.xl2),
        boxShadow: LumoraShadows.cardLg,
      ),
      child: child,
    );
  }
}

/// A tinted circular icon badge used at the top of auth cards (check-email,
/// success, etc.) — a solid tint, never a gradient.
class AuthIconBadge extends StatelessWidget {
  final IconData icon;
  final Color color;
  const AuthIconBadge({super.key, required this.icon, this.color = LumoraColors.purple});

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 56,
      height: 56,
      decoration: BoxDecoration(color: color.withValues(alpha: 0.12), shape: BoxShape.circle),
      child: Icon(icon, color: color, size: 26),
    );
  }
}

class AuthFieldLabel extends StatelessWidget {
  final String text;
  const AuthFieldLabel(this.text, {super.key});
  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 6, left: 2),
      child: Text(text, style: const TextStyle(fontSize: 12, fontWeight: FontWeight.w700, color: LumoraColors.slatey)),
    );
  }
}

/// The small circular back button used on auth-flow screens.
class AuthBackButton extends StatelessWidget {
  final VoidCallback onTap;
  const AuthBackButton({super.key, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return Material(
      color: LumoraColors.gray50,
      shape: const CircleBorder(),
      child: InkWell(
        customBorder: const CircleBorder(),
        onTap: onTap,
        child: const SizedBox(width: 36, height: 36, child: Icon(Icons.arrow_back_rounded, size: 18, color: LumoraColors.ink)),
      ),
    );
  }
}
