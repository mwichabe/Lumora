import 'package:flutter/material.dart';

/// Lumora the fennec fox — a simple vector mascot so the app needs no bundled
/// image assets (frontend renders it as inline SVG for the same reason).
class FoxMascot extends StatefulWidget {
  final double size;
  final bool bounce;
  final bool glow;

  const FoxMascot({super.key, this.size = 120, this.bounce = false, this.glow = false});

  @override
  State<FoxMascot> createState() => _FoxMascotState();
}

class _FoxMascotState extends State<FoxMascot> with SingleTickerProviderStateMixin {
  late final AnimationController _c;

  @override
  void initState() {
    super.initState();
    _c = AnimationController(vsync: this, duration: const Duration(seconds: 3))
      ..repeat(reverse: true);
  }

  @override
  void dispose() {
    _c.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final child = SizedBox(
      width: widget.size,
      height: widget.size,
      child: CustomPaint(painter: _FoxPainter(glow: widget.glow)),
    );
    if (!widget.bounce) return child;
    return AnimatedBuilder(
      animation: _c,
      builder: (context, c) => Transform.translate(
        offset: Offset(0, -6 * _c.value),
        child: c,
      ),
      child: child,
    );
  }
}

class _FoxPainter extends CustomPainter {
  final bool glow;
  _FoxPainter({required this.glow});

  @override
  void paint(Canvas canvas, Size size) {
    final w = size.width, h = size.height;
    final center = Offset(w / 2, h / 2);

    if (glow) {
      final glowPaint = Paint()
        ..color = const Color(0xFFF5A623).withValues(alpha: 0.35)
        ..maskFilter = const MaskFilter.blur(BlurStyle.normal, 30);
      canvas.drawCircle(center, w * 0.42, glowPaint);
    }

    final fur = Paint()..color = const Color(0xFFF5A623);
    final furDark = Paint()..color = const Color(0xFFE08E1A);
    final cream = Paint()..color = const Color(0xFFFFF8E7);
    final ink = Paint()..color = const Color(0xFF1A1A2E);

    // Ears
    final earL = Path()
      ..moveTo(w * 0.28, h * 0.30)
      ..lineTo(w * 0.14, h * 0.02)
      ..lineTo(w * 0.42, h * 0.20)
      ..close();
    final earR = Path()
      ..moveTo(w * 0.72, h * 0.30)
      ..lineTo(w * 0.86, h * 0.02)
      ..lineTo(w * 0.58, h * 0.20)
      ..close();
    canvas.drawPath(earL, fur);
    canvas.drawPath(earR, fur);
    final innerEarL = Path()
      ..moveTo(w * 0.30, h * 0.26)
      ..lineTo(w * 0.23, h * 0.10)
      ..lineTo(w * 0.38, h * 0.21)
      ..close();
    final innerEarR = Path()
      ..moveTo(w * 0.70, h * 0.26)
      ..lineTo(w * 0.77, h * 0.10)
      ..lineTo(w * 0.62, h * 0.21)
      ..close();
    canvas.drawPath(innerEarL, cream);
    canvas.drawPath(innerEarR, cream);

    // Head
    canvas.drawOval(Rect.fromCenter(center: Offset(w * 0.5, h * 0.52), width: w * 0.62, height: h * 0.56), fur);

    // Muzzle
    final muzzle = Path()
      ..moveTo(w * 0.34, h * 0.55)
      ..quadraticBezierTo(w * 0.5, h * 0.78, w * 0.66, h * 0.55)
      ..quadraticBezierTo(w * 0.5, h * 0.68, w * 0.34, h * 0.55)
      ..close();
    canvas.drawPath(muzzle, cream);

    // Nose
    canvas.drawCircle(Offset(w * 0.5, h * 0.60), w * 0.035, ink);

    // Eyes
    canvas.drawCircle(Offset(w * 0.40, h * 0.47), w * 0.045, ink);
    canvas.drawCircle(Offset(w * 0.60, h * 0.47), w * 0.045, ink);
    canvas.drawCircle(Offset(w * 0.415, h * 0.458), w * 0.014, Paint()..color = Colors.white);
    canvas.drawCircle(Offset(w * 0.615, h * 0.458), w * 0.014, Paint()..color = Colors.white);

    // Cheeks
    canvas.drawCircle(
      Offset(w * 0.30, h * 0.58),
      w * 0.05,
      Paint()..color = const Color(0xFFFF9A9A).withValues(alpha: 0.55),
    );
    canvas.drawCircle(
      Offset(w * 0.70, h * 0.58),
      w * 0.05,
      Paint()..color = const Color(0xFFFF9A9A).withValues(alpha: 0.55),
    );

    // Head shading
    canvas.drawArc(
      Rect.fromCenter(center: Offset(w * 0.5, h * 0.52), width: w * 0.62, height: h * 0.56),
      0.3,
      1.6,
      false,
      furDark..style = PaintingStyle.stroke..strokeWidth = 2,
    );
  }

  @override
  bool shouldRepaint(covariant _FoxPainter oldDelegate) => oldDelegate.glow != glow;
}
