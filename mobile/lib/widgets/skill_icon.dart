import 'package:flutter/material.dart';

/// Maps the lucide icon names stored on a Skill to Material icons, mirroring
/// frontend/components/SkillIcon.tsx (kept icon-based, not emoji).
const _kSkillIcons = <String, IconData>{
  'Hand': Icons.pan_tool_rounded,
  'Sparkles': Icons.auto_awesome_rounded,
  'Coffee': Icons.coffee_rounded,
  'ShoppingBag': Icons.shopping_bag_rounded,
  'Plane': Icons.flight_rounded,
  'Users': Icons.groups_rounded,
  'MessageCircle': Icons.chat_bubble_rounded,
  'Heart': Icons.favorite_rounded,
  'Compass': Icons.explore_rounded,
  'Hash': Icons.tag_rounded,
  'BookOpen': Icons.menu_book_rounded,
  'GraduationCap': Icons.school_rounded,
  'PenLine': Icons.edit_rounded,
  'Quote': Icons.format_quote_rounded,
  'Languages': Icons.translate_rounded,
  'Layers': Icons.layers_rounded,
  'Clock': Icons.access_time_rounded,
  'Link2': Icons.link_rounded,
  // Mandarin course
  'Music': Icons.music_note_rounded,
  'Waves': Icons.waves_rounded,
  'CheckCircle': Icons.check_circle_rounded,
  'Globe': Icons.public_rounded,
  'Scale': Icons.balance_rounded,
  'Link': Icons.link_rounded,
  'Shield': Icons.shield_rounded,
  'HeartPulse': Icons.monitor_heart_rounded,
  'Briefcase': Icons.work_rounded,
  'Home': Icons.home_rounded,
  'Smartphone': Icons.smartphone_rounded,
  'Landmark': Icons.account_balance_rounded,
  'FlaskConical': Icons.science_rounded,
  'Brain': Icons.psychology_rounded,
  'Cpu': Icons.memory_rounded,
  'Leaf': Icons.eco_rounded,
  'Mic': Icons.mic_rounded,
};

class SkillIcon extends StatelessWidget {
  final String? name;
  final double size;
  final Color? color;
  const SkillIcon({super.key, this.name, this.size = 24, this.color});

  @override
  Widget build(BuildContext context) {
    final icon = (name != null ? _kSkillIcons[name] : null) ?? Icons.menu_book_rounded;
    return Icon(icon, size: size, color: color);
  }
}
