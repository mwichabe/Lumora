import 'package:flutter/material.dart';

/// Colour tokens lifted 1:1 from the Lumora Figma Design Spec
/// (frontend/tailwind.config.ts).
class LumoraColors {
  LumoraColors._();

  // Brand
  static const purple = Color(0xFF6C3FC5);
  static const purpleDark = Color(0xFF3A1F8A);
  static const purpleLight = Color(0xFFEDE7F6);

  static const amber = Color(0xFFF5A623);
  static const amberLight = Color(0xFFFFF8E7);

  static const teal = Color(0xFF00C2A8);
  static const tealLight = Color(0xFFE0FAF7);

  static const coral = Color(0xFFFF5C5C);
  static const coralLight = Color(0xFFFFE8E8);

  // Neutrals
  static const cream = Color(0xFFFAFAF7);
  static const gray50 = Color(0xFFF5F5F5);
  static const gray100 = Color(0xFFEBEBEB);
  static const gray300 = Color(0xFFCCCCCC);
  static const gray500 = Color(0xFF9090A0);

  static const ink = Color(0xFF1A1A2E); // primary text
  static const slatey = Color(0xFF4A4A6A); // secondary text
  static const space = Color(0xFF0F0F24); // galaxy map deep background

  // Backdrop used behind the centred app frame on wide screens / splash
  static const outerBg = Color(0xFFECEAF3);

  // Misc one-offs used across the web app
  static const stampInk = Color(0xFFB23A6B);
}
