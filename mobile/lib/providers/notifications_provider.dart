import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../core/network/api_client.dart';
import '../models/notification.dart';
import 'auth_provider.dart';

/// Polls the unread notification count every 60s, matching frontend/lib/notifications.ts.
final unreadNotificationsProvider = StreamProvider.autoDispose<int>((ref) async* {
  final auth = ref.watch(authProvider);
  if (!auth.isAuthenticated) {
    yield 0;
    return;
  }
  while (true) {
    try {
      final (_, unread) = await ApiClient.instance.notifications();
      yield unread;
    } catch (_) {
      // keep the previous value on a transient failure
    }
    await Future.delayed(const Duration(seconds: 60));
  }
});

class NotificationsState {
  final List<AppNotification> items;
  final int unread;
  final bool loading;
  const NotificationsState({required this.items, required this.unread, required this.loading});

  factory NotificationsState.initial() => const NotificationsState(items: [], unread: 0, loading: true);

  NotificationsState copyWith({List<AppNotification>? items, int? unread, bool? loading}) =>
      NotificationsState(
        items: items ?? this.items,
        unread: unread ?? this.unread,
        loading: loading ?? this.loading,
      );
}

/// The notifications list. Auto-disposed so every visit to the screen fetches a
/// fresh list (the badge polls the server on its own, so a list cached from an
/// earlier visit would show nothing new while the badge says otherwise).
class NotificationsController extends AutoDisposeNotifier<NotificationsState> {
  @override
  NotificationsState build() {
    // Rebuild (and drop the previous account's items) when the session changes.
    ref.watch(authProvider.select((a) => a.user?.id));
    Future.microtask(load);
    return NotificationsState.initial();
  }

  Future<void> load() async {
    try {
      final (items, unread) = await ApiClient.instance.notifications();
      state = state.copyWith(items: items, unread: unread, loading: false);
      _syncBadge();
    } catch (_) {
      state = state.copyWith(loading: false);
    }
  }

  /// Re-polls the unread badge right away instead of waiting for its next tick.
  void _syncBadge() => ref.invalidate(unreadNotificationsProvider);

  Future<void> markAllRead() async {
    try {
      await ApiClient.instance.markNotificationsRead();
    } catch (_) {
      return;
    }
    state = state.copyWith(
      items: state.items.map((n) => AppNotification(
            id: n.id, kind: n.kind, emoji: n.emoji, tint: n.tint, title: n.title,
            body: n.body, link: n.link, read: true, createdAt: n.createdAt,
          )).toList(),
      unread: 0,
    );
    _syncBadge();
  }

  Future<void> markRead(int id) async {
    final int unread;
    try {
      unread = await ApiClient.instance.markNotificationRead(id);
    } catch (_) {
      return;
    }
    state = state.copyWith(
      items: [
        for (final n in state.items)
          if (n.id == id)
            AppNotification(id: n.id, kind: n.kind, emoji: n.emoji, tint: n.tint, title: n.title,
                body: n.body, link: n.link, read: true, createdAt: n.createdAt)
          else
            n,
      ],
      unread: unread,
    );
    _syncBadge();
  }

  Future<void> delete(int id) async {
    final int unread;
    try {
      unread = await ApiClient.instance.deleteNotification(id);
    } catch (_) {
      return load(); // the row was swiped away locally; resync with the server
    }
    state = state.copyWith(items: state.items.where((n) => n.id != id).toList(), unread: unread);
    _syncBadge();
  }
}

final notificationsProvider = NotifierProvider.autoDispose<NotificationsController, NotificationsState>(
  NotificationsController.new,
);
