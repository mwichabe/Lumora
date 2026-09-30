import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../core/network/api_client.dart';
import 'auth_provider.dart';

/// Polls the global unread chat count every 15s, matching frontend/lib/chat.ts.
final chatUnreadProvider = StreamProvider.autoDispose<int>((ref) async* {
  final auth = ref.watch(authProvider);
  if (!auth.isAuthenticated) {
    yield 0;
    return;
  }
  while (true) {
    try {
      yield await ApiClient.instance.chatUnread();
    } catch (_) {
      // keep previous value on transient failure
    }
    await Future.delayed(const Duration(seconds: 15));
  }
});
