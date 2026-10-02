import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';

import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../models/notification.dart';
import '../../providers/notifications_provider.dart';

const _kKindLabel = {
  'welcome': 'Welcome', 'welcome_back': 'Welcome back', 'milestone': 'Milestone', 'exam': 'Exam',
  'hearts': 'Hearts', 'payment': 'Payment', 'tip': 'Tip', 'feature': 'New feature',
  'language': 'Languages', 'streak': 'Streak', 'league': 'League', 'chat': 'Message', 'message': 'Message',
};

String _timeAgo(String iso) {
  final t = DateTime.tryParse(iso);
  if (t == null) return '';
  final s = DateTime.now().difference(t).inSeconds.clamp(0, 1 << 31);
  if (s < 60) return 'just now';
  final m = s ~/ 60;
  if (m < 60) return '${m}m ago';
  final h = m ~/ 60;
  if (h < 24) return '${h}h ago';
  return '${h ~/ 24}d ago';
}

class NotificationsScreen extends ConsumerWidget {
  const NotificationsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(notificationsProvider);
    final newItems = state.items.where((n) => !n.read).toList();
    final earlier = state.items.where((n) => n.read).toList();

    return Scaffold(
      backgroundColor: Colors.white,
      appBar: AppBar(
        backgroundColor: Colors.white,
        foregroundColor: LumoraColors.ink,
        elevation: 0,
        leading: BackButton(onPressed: () => context.canPop() ? context.pop() : context.go('/home')),
        title: const Text('Notifications', style: TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.ink)),
        actions: [
          if (newItems.isNotEmpty)
            TextButton(
              onPressed: () => ref.read(notificationsProvider.notifier).markAllRead(),
              child: const Text('Mark all as read'),
            ),
        ],
      ),
      body: state.loading
          ? const Center(child: CircularProgressIndicator())
          : RefreshIndicator(
              onRefresh: () => ref.read(notificationsProvider.notifier).load(),
              child: state.items.isEmpty
                  ? ListView(
                      physics: const AlwaysScrollableScrollPhysics(),
                      children: const [SizedBox(height: 120), _EmptyState()],
                    )
                  : ListView(
                      physics: const AlwaysScrollableScrollPhysics(),
                      children: [
                        if (newItems.isNotEmpty) _Section(label: 'New', items: newItems, ref: ref),
                        if (earlier.isNotEmpty) _Section(label: 'Earlier', items: earlier, ref: ref),
                      ],
                    ),
            ),
    );
  }
}

class _EmptyState extends StatelessWidget {
  const _EmptyState();
  @override
  Widget build(BuildContext context) {
    return const Center(
      child: Padding(
        padding: EdgeInsets.all(32),
        child: Column(mainAxisSize: MainAxisSize.min, children: [
          Icon(Icons.notifications_off_outlined, size: 48, color: LumoraColors.gray300),
          SizedBox(height: 12),
          Text("You're all caught up", style: TextStyle(fontWeight: FontWeight.w800)),
          SizedBox(height: 4),
          Text('New milestones, league results, and tips will show up here.', textAlign: TextAlign.center, style: TextStyle(color: LumoraColors.slatey)),
        ]),
      ),
    );
  }
}

class _Section extends StatelessWidget {
  final String label;
  final List<AppNotification> items;
  final WidgetRef ref;
  const _Section({required this.label, required this.items, required this.ref});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 16, 16, 4),
          child: Text(label.toUpperCase(), style: const TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.gray500)),
        ),
        for (final n in items) _NotificationTile(n: n, ref: ref),
      ],
    );
  }
}

class _NotificationTile extends StatelessWidget {
  final AppNotification n;
  final WidgetRef ref;
  const _NotificationTile({required this.n, required this.ref});

  @override
  Widget build(BuildContext context) {
    return Dismissible(
      key: ValueKey(n.id),
      direction: DismissDirection.endToStart,
      background: Container(color: LumoraColors.coral, alignment: Alignment.centerRight, padding: const EdgeInsets.only(right: 20), child: const Icon(Icons.delete, color: Colors.white)),
      onDismissed: (_) => ref.read(notificationsProvider.notifier).delete(n.id),
      child: ListTile(
        onTap: () {
          if (!n.read) ref.read(notificationsProvider.notifier).markRead(n.id);
          showModalBottomSheet(
            context: context,
            backgroundColor: Colors.white,
            shape: const RoundedRectangleBorder(borderRadius: BorderRadius.vertical(top: Radius.circular(24))),
            builder: (sheetContext) => _NotificationDetail(
              n: n,
              onFollowLink: n.link.isEmpty
                  ? null
                  : () {
                      Navigator.pop(sheetContext);
                      _followLink(context, n.link);
                    },
            ),
          );
        },
        leading: Container(
          width: 44, height: 44,
          decoration: BoxDecoration(color: _hex(n.tint).withValues(alpha: 0.12), borderRadius: BorderRadius.circular(LumoraRadii.lg)),
          child: Center(child: Text(n.emoji, style: const TextStyle(fontSize: 20))),
        ),
        title: Text(n.title, maxLines: 1, overflow: TextOverflow.ellipsis, style: TextStyle(fontWeight: n.read ? FontWeight.w600 : FontWeight.w800)),
        subtitle: Text(n.body, maxLines: 2, overflow: TextOverflow.ellipsis, style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
        trailing: Column(mainAxisAlignment: MainAxisAlignment.center, crossAxisAlignment: CrossAxisAlignment.end, children: [
          Text(_timeAgo(n.createdAt), style: const TextStyle(fontSize: 10, color: LumoraColors.gray500)),
          if (!n.read) Container(margin: const EdgeInsets.only(top: 4), width: 8, height: 8, decoration: const BoxDecoration(color: LumoraColors.coral, shape: BoxShape.circle)),
        ]),
      ),
    );
  }
}

/// Tab destinations live in the bottom-nav shell and must be switched to with
/// `go`; anything else opens on top of the notifications screen.
const _kTabRoutes = ['/home', '/learn', '/practice', '/leaderboard'];

void _followLink(BuildContext context, String link) {
  final path = Uri.parse(link).path;
  if (_kTabRoutes.contains(path)) {
    context.go(link);
  } else {
    context.push(link);
  }
}

class _NotificationDetail extends StatelessWidget {
  final AppNotification n;
  final VoidCallback? onFollowLink;
  const _NotificationDetail({required this.n, this.onFollowLink});
  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Text(n.emoji, style: const TextStyle(fontSize: 28)),
              const SizedBox(width: 12),
              Expanded(child: Text(n.title, style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w800))),
            ]),
            const SizedBox(height: 4),
            Text(_kKindLabel[n.kind] ?? 'Notification', style: const TextStyle(color: LumoraColors.slatey, fontSize: 11)),
            const SizedBox(height: 12),
            Text(n.body, style: const TextStyle(fontSize: 15)),
            const SizedBox(height: 12),
            Text(_fullDate(n.createdAt), style: const TextStyle(color: LumoraColors.gray500, fontSize: 11)),
            if (onFollowLink != null) ...[
              const SizedBox(height: 16),
              SizedBox(
                width: double.infinity,
                child: FilledButton(
                  onPressed: onFollowLink,
                  style: FilledButton.styleFrom(backgroundColor: LumoraColors.purple),
                  child: const Text('Open'),
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }

  String _fullDate(String iso) {
    final d = DateTime.tryParse(iso);
    if (d == null) return '';
    return DateFormat.yMMMd().add_jm().format(d);
  }
}

Color _hex(String s) {
  var h = s.replaceAll('#', '');
  if (h.length == 6) h = 'FF$h';
  final v = int.tryParse(h, radix: 16);
  return v == null ? LumoraColors.purple : Color(v);
}
