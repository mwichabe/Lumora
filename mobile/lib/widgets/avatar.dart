import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import '../core/env.dart';

Color colorFromHex(String hex, {Color fallback = const Color(0xFF6C3FC5)}) {
  var h = hex.replaceAll('#', '');
  if (h.length == 6) h = 'FF$h';
  if (h.length != 8) return fallback;
  final v = int.tryParse(h, radix: 16);
  return v == null ? fallback : Color(v);
}

/// Circular avatar: shows the uploaded photo when present, falling back to a
/// coloured initial circle (frontend/components/Avatar.tsx).
class LumoraAvatar extends StatelessWidget {
  final String? avatarUrl;
  final String name;
  final String avatarColor;
  final double size;

  const LumoraAvatar({
    super.key,
    required this.name,
    this.avatarUrl,
    this.avatarColor = '#6C3FC5',
    this.size = 44,
  });

  @override
  Widget build(BuildContext context) {
    final bg = colorFromHex(avatarColor);
    final initial = name.isNotEmpty ? name.trim()[0].toUpperCase() : '?';

    final fallback = Container(
      width: size,
      height: size,
      alignment: Alignment.center,
      decoration: BoxDecoration(color: bg, shape: BoxShape.circle),
      child: Text(
        initial,
        style: TextStyle(color: Colors.white, fontWeight: FontWeight.w800, fontSize: size * 0.42),
      ),
    );

    if (avatarUrl == null || avatarUrl!.isEmpty) return fallback;

    return ClipOval(
      child: CachedNetworkImage(
        imageUrl: mediaUrl(avatarUrl),
        width: size,
        height: size,
        fit: BoxFit.cover,
        placeholder: (_, _) => fallback,
        errorWidget: (_, _, _) => fallback,
      ),
    );
  }
}
