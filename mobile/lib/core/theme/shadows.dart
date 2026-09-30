import 'package:flutter/material.dart';

/// Box-shadow tokens from the Design Spec (tailwind.config.ts `boxShadow`).
class LumoraShadows {
  LumoraShadows._();

  static const card = [
    BoxShadow(color: Color(0x14000000), blurRadius: 8, offset: Offset(0, 2)),
  ];

  static const cardLg = [
    BoxShadow(color: Color(0x14000000), blurRadius: 16, offset: Offset(0, 4)),
  ];

  static const skill = [
    BoxShadow(color: Color(0x266C3FC5), blurRadius: 16, offset: Offset(0, 4)),
  ];

  static const quest = [
    BoxShadow(color: Color(0x26F5A623), blurRadius: 8, offset: Offset(0, 2)),
  ];

  static const float = [
    BoxShadow(color: Color(0x406C3FC5), blurRadius: 24, offset: Offset(0, 8)),
  ];
}
