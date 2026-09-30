import 'package:flutter/material.dart';
import '../core/theme/colors.dart';
import '../core/theme/radii.dart';

/// Amber-gradient XP progress bar (frontend/components/widgets.tsx XPBar).
class XPBar extends StatelessWidget {
  final int value;
  final int max;
  const XPBar({super.key, required this.value, required this.max});

  @override
  Widget build(BuildContext context) {
    final pct = max > 0 ? (value / max).clamp(0, 1).toDouble() : 0.0;
    return ClipRRect(
      borderRadius: BorderRadius.circular(LumoraRadii.full),
      child: Container(
        height: 12,
        color: LumoraColors.gray100,
        child: Align(
          alignment: Alignment.centerLeft,
          child: FractionallySizedBox(
            widthFactor: pct,
            child: TweenAnimationBuilder<double>(
              tween: Tween(begin: 0, end: 1),
              duration: const Duration(milliseconds: 500),
              builder: (context, v, _) => Container(
                decoration: const BoxDecoration(
                  gradient: LinearGradient(colors: [LumoraColors.amber, Color(0xFFFF9A00)]),
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}

/// Row of 5 hearts with an optional countdown to the next regenerated one.
class HeartIndicator extends StatelessWidget {
  final int hearts;
  final int secondsToNext;
  final bool light;
  const HeartIndicator({super.key, required this.hearts, this.secondsToNext = 0, this.light = false});

  String _fmt(int s) {
    final m = s ~/ 60;
    final r = s % 60;
    return '$m:${r.toString().padLeft(2, '0')}';
  }

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        for (var i = 0; i < 5; i++)
          Icon(
            i < hearts ? Icons.favorite : Icons.favorite_border,
            size: 20,
            color: i < hearts ? LumoraColors.coral : LumoraColors.coral.withValues(alpha: 0.3),
          ),
        if (hearts < 5 && secondsToNext > 0) ...[
          const SizedBox(width: 4),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
            decoration: BoxDecoration(
              color: LumoraColors.coral.withValues(alpha: 0.1),
              borderRadius: BorderRadius.circular(LumoraRadii.full),
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Icon(Icons.access_time, size: 11, color: LumoraColors.coral),
                const SizedBox(width: 2),
                Text(_fmt(secondsToNext),
                    style: const TextStyle(
                        fontSize: 10, fontWeight: FontWeight.w800, color: LumoraColors.coral)),
              ],
            ),
          ),
        ],
      ],
    );
  }
}

class StreakFlame extends StatelessWidget {
  final int streak;
  final bool light;
  const StreakFlame({super.key, required this.streak, this.light = false});

  @override
  Widget build(BuildContext context) {
    final active = streak > 0;
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(Icons.local_fire_department,
            size: 22, color: active ? LumoraColors.amber : (light ? Colors.white38 : LumoraColors.gray300)),
        const SizedBox(width: 4),
        Text('$streak',
            style: TextStyle(
                fontWeight: FontWeight.w800, color: light ? Colors.white : LumoraColors.ink)),
      ],
    );
  }
}

class GemCounter extends StatelessWidget {
  final int gems;
  const GemCounter({super.key, required this.gems});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
      decoration: BoxDecoration(color: LumoraColors.teal, borderRadius: BorderRadius.circular(LumoraRadii.full)),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          const Icon(Icons.diamond, size: 16, color: Colors.white),
          const SizedBox(width: 4),
          Text('$gems', style: const TextStyle(color: Colors.white, fontWeight: FontWeight.w700, fontSize: 12)),
        ],
      ),
    );
  }
}

class SpeechBubble extends StatelessWidget {
  final Widget child;
  const SpeechBubble({super.key, required this.child});

  @override
  Widget build(BuildContext context) {
    return Stack(
      clipBehavior: Clip.none,
      children: [
        Container(
          constraints: const BoxConstraints(maxWidth: 280),
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.circular(28),
            border: Border.all(color: LumoraColors.gray100),
            boxShadow: const [BoxShadow(color: Color(0x14000000), blurRadius: 8, offset: Offset(0, 2))],
          ),
          child: child,
        ),
        Positioned(
          bottom: -8,
          left: 28,
          child: Transform.rotate(
            angle: 0.785398,
            child: Container(
              width: 16,
              height: 16,
              decoration: BoxDecoration(
                color: Colors.white,
                border: Border.all(color: LumoraColors.gray100),
              ),
            ),
          ),
        ),
      ],
    );
  }
}
