import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../core/theme/colors.dart';
import '../providers/notifications_provider.dart';
import '../providers/chat_provider.dart';

/// Circular translucent icon button for coloured headers, with a coral unread
/// badge. Mirrors NotificationBell.tsx / ChatBell.tsx / IdeasBell.tsx.
class _HeaderBell extends StatelessWidget {
  final IconData icon;
  final int badge;
  final VoidCallback onTap;
  const _HeaderBell({required this.icon, required this.badge, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return Stack(
      clipBehavior: Clip.none,
      children: [
        Material(
          color: Colors.white.withValues(alpha: 0.15),
          shape: const CircleBorder(),
          child: InkWell(
            customBorder: const CircleBorder(),
            onTap: onTap,
            child: SizedBox(
              width: 38,
              height: 38,
              child: Icon(icon, color: Colors.white, size: 19),
            ),
          ),
        ),
        if (badge > 0)
          Positioned(
            top: -2,
            right: -2,
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 1),
              constraints: const BoxConstraints(minWidth: 16),
              decoration: const BoxDecoration(color: LumoraColors.coral, shape: BoxShape.circle),
              child: Text(
                badge > 9 ? '9+' : '$badge',
                textAlign: TextAlign.center,
                style: const TextStyle(color: Colors.white, fontSize: 10, fontWeight: FontWeight.w800),
              ),
            ),
          ),
      ],
    );
  }
}

class NotificationBell extends ConsumerWidget {
  const NotificationBell({super.key});
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final unread = ref.watch(unreadNotificationsProvider).value ?? 0;
    return _HeaderBell(icon: Icons.notifications_rounded, badge: unread, onTap: () => context.push('/notifications'));
  }
}

class ChatBell extends ConsumerWidget {
  const ChatBell({super.key});
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final unread = ref.watch(chatUnreadProvider).value ?? 0;
    return _HeaderBell(icon: Icons.chat_bubble_rounded, badge: unread, onTap: () => context.push('/chat'));
  }
}

class IdeasBell extends StatelessWidget {
  const IdeasBell({super.key});
  @override
  Widget build(BuildContext context) {
    return _HeaderBell(icon: Icons.lightbulb_rounded, badge: 0, onTap: () => context.push('/ideas'));
  }
}
