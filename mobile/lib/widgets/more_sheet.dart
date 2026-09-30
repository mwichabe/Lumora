import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../core/theme/colors.dart';
import '../core/theme/radii.dart';
import '../providers/notifications_provider.dart';
import '../providers/chat_provider.dart';

void showMoreSheet(BuildContext context) {
  showModalBottomSheet(
    context: context,
    backgroundColor: Colors.white,
    shape: const RoundedRectangleBorder(
      borderRadius: BorderRadius.vertical(top: Radius.circular(28)),
    ),
    builder: (_) => const _MoreSheetContent(),
  );
}

class _MoreSheetContent extends ConsumerWidget {
  const _MoreSheetContent();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final unreadNotifs = ref.watch(unreadNotificationsProvider);
    final unreadChat = ref.watch(chatUnreadProvider);

    final items = [
      (route: '/ideas', label: 'Ideas', hint: 'Propose, vote, discuss', icon: Icons.lightbulb_rounded, tint: LumoraColors.amber, badge: 0),
      (route: '/chat', label: 'Messages', hint: 'Chat with other learners', icon: Icons.chat_bubble_rounded, tint: LumoraColors.purple, badge: unreadChat.value ?? 0),
      (route: '/notifications', label: 'Notifications', hint: 'League results, milestones', icon: Icons.notifications_rounded, tint: LumoraColors.teal, badge: unreadNotifs.value ?? 0),
      (route: '/profile', label: 'Profile', hint: 'Certificates, settings', icon: Icons.person_rounded, tint: const Color(0xFF17A3DD), badge: 0),
    ];

    return SafeArea(
      child: Padding(
        padding: const EdgeInsets.fromLTRB(12, 8, 12, 12),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(width: 40, height: 4, margin: const EdgeInsets.only(bottom: 12),
                decoration: BoxDecoration(color: LumoraColors.gray100, borderRadius: BorderRadius.circular(2))),
            Row(
              children: [
                const Padding(
                  padding: EdgeInsets.only(left: 8),
                  child: Text('More', style: TextStyle(fontSize: 18, fontWeight: FontWeight.w800, color: LumoraColors.ink)),
                ),
                const Spacer(),
                IconButton(icon: const Icon(Icons.close, size: 20), onPressed: () => Navigator.pop(context)),
              ],
            ),
            for (final item in items)
              ListTile(
                onTap: () {
                  Navigator.pop(context);
                  context.push(item.route);
                },
                leading: Container(
                  width: 44, height: 44,
                  decoration: BoxDecoration(color: item.tint.withValues(alpha: 0.1), borderRadius: BorderRadius.circular(LumoraRadii.lg)),
                  child: Icon(item.icon, color: item.tint, size: 21),
                ),
                title: Row(
                  children: [
                    Text(item.label, style: const TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.ink)),
                    if (item.badge > 0) ...[
                      const SizedBox(width: 8),
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 1),
                        decoration: BoxDecoration(color: LumoraColors.coral, borderRadius: BorderRadius.circular(LumoraRadii.full)),
                        child: Text(item.badge > 9 ? '9+' : '${item.badge}',
                            style: const TextStyle(color: Colors.white, fontSize: 11, fontWeight: FontWeight.w800)),
                      ),
                    ],
                  ],
                ),
                subtitle: Text(item.hint, style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
                trailing: const Icon(Icons.chevron_right, color: LumoraColors.gray300),
              ),
          ],
        ),
      ),
    );
  }
}
