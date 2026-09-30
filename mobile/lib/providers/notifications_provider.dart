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

class NotificationsController extends Notifier<NotificationsState> {
  @override
  NotificationsState build() {
    load();
    return NotificationsState.initial();
  }

  Future<void> load() async {
    try {
      final (items, unread) = await ApiClient.instance.notifications();
      state = state.copyWith(items: items, unread: unread, loading: false);
    } catch (_) {
      state = state.copyWith(loading: false);
    }
  }

  Future<void> markAllRead() async {
    await ApiClient.instance.markNotificationsRead();
    state = state.copyWith(
      items: state.items.map((n) => AppNotification(
            id: n.id, kind: n.kind, emoji: n.emoji, tint: n.tint, title: n.title,
            body: n.body, link: n.link, read: true, createdAt: n.createdAt,
          )).toList(),
      unread: 0,
    );
  }

  Future<void> markRead(int id) async {
    final unread = await ApiClient.instance.markNotificationRead(id);
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
  }

  Future<void> delete(int id) async {
    final unread = await ApiClient.instance.deleteNotification(id);
    state = state.copyWith(items: state.items.where((n) => n.id != id).toList(), unread: unread);
  }
}

final notificationsProvider = NotifierProvider<NotificationsController, NotificationsState>(
  NotificationsController.new,
);
