import 'package:google_fonts/google_fonts.dart';
import 'package:flutter/material.dart';
import 'colors.dart';

/// Type scale from the Design Spec (tailwind.config.ts `fontSize`).
/// Font family: Nunito, weights 400/600/700/800.
class LumoraText {
  LumoraText._();

  static TextStyle _base(double size, double height, FontWeight weight, Color color) {
    return GoogleFonts.nunito(
      fontSize: size,
      height: height / size,
      fontWeight: weight,
      color: color,
    );
  }

  static TextStyle displayXl({Color color = LumoraColors.ink}) =>
      _base(32, 40, FontWeight.w800, color);
  static TextStyle displayLg({Color color = LumoraColors.ink}) =>
      _base(28, 36, FontWeight.w800, color);
  static TextStyle headingXl({Color color = LumoraColors.ink}) =>
      _base(24, 32, FontWeight.w800, color);
  static TextStyle headingLg({Color color = LumoraColors.ink}) =>
      _base(20, 28, FontWeight.w800, color);
  static TextStyle headingMd({Color color = LumoraColors.ink}) =>
      _base(18, 26, FontWeight.w800, color);
  static TextStyle headingSm({Color color = LumoraColors.ink}) =>
      _base(16, 24, FontWeight.w800, color);
  static TextStyle bodyLg({Color color = LumoraColors.ink}) =>
      _base(16, 24, FontWeight.w400, color);
  static TextStyle bodyMd({Color color = LumoraColors.ink}) =>
      _base(14, 22, FontWeight.w400, color);
  static TextStyle bodySm({Color color = LumoraColors.slatey}) =>
      _base(12, 18, FontWeight.w400, color);
  static TextStyle labelLg({Color color = LumoraColors.ink}) =>
      _base(14, 20, FontWeight.w700, color);
  static TextStyle labelMd({Color color = LumoraColors.slatey}) =>
      _base(12, 16, FontWeight.w700, color);
  static TextStyle labelSm({Color color = LumoraColors.slatey}) =>
      _base(10, 14, FontWeight.w700, color);
}
