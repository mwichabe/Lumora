import 'dart:async';

import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../core/network/api_client.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../models/chat.dart';
import '../../widgets/avatar.dart';

String _timeAgo(String iso) {
  final t = DateTime.tryParse(iso);
  if (t == null) return '';
  final s = DateTime.now().difference(t).inSeconds.clamp(0, 1 << 31);
  if (s < 60) return 'now';
  final m = s ~/ 60;
  if (m < 60) return '${m}m';
  final h = m ~/ 60;
  if (h < 24) return '${h}h';
  return '${h ~/ 24}d';
}

class ChatListScreen extends StatefulWidget {
  const ChatListScreen({super.key});

  @override
  State<ChatListScreen> createState() => _ChatListScreenState();
}

class _ChatListScreenState extends State<ChatListScreen> {
  List<ChatThread>? _threads;
  Timer? _poll;

  @override
  void initState() {
    super.initState();
    _load();
    _poll = Timer.periodic(const Duration(seconds: 8), (_) => _load());
  }

  @override
  void dispose() {
    _poll?.cancel();
    super.dispose();
  }

  Future<void> _load() async {
    try {
      final threads = await ApiClient.instance.chatThreads();
      if (mounted) setState(() => _threads = threads);
    } catch (_) {
      if (mounted && _threads == null) setState(() => _threads = []);
    }
  }

  void _openPicker() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.white,
      shape: const RoundedRectangleBorder(borderRadius: BorderRadius.vertical(top: Radius.circular(24))),
      builder: (_) => const _ContactPicker(),
    );
  }

  @override
  Widget build(BuildContext context) {
    final threads = _threads;
    return Scaffold(
      backgroundColor: LumoraColors.cream,
      appBar: AppBar(
        backgroundColor: Colors.white,
        foregroundColor: LumoraColors.ink,
        elevation: 0,
        leading: BackButton(onPressed: () => context.canPop() ? context.pop() : context.go('/home')),
        title: const Text('Messages', style: TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.ink)),
        actions: [
          Padding(
            padding: const EdgeInsets.only(right: 12),
            child: TextButton.icon(
              onPressed: _openPicker,
              icon: const Icon(Icons.add_comment_rounded, size: 16),
              label: const Text('New'),
              style: TextButton.styleFrom(backgroundColor: LumoraColors.purple, foregroundColor: Colors.white, shape: const StadiumBorder()),
            ),
          ),
        ],
      ),
      body: threads == null
          ? const Center(child: CircularProgressIndicator())
          : threads.isEmpty
              ? _EmptyState(onStart: _openPicker)
              : ListView.separated(
                  itemCount: threads.length,
                  separatorBuilder: (_, _) => const Divider(height: 1),
                  itemBuilder: (context, i) {
                    final t = threads[i];
                    return ListTile(
                      onTap: () => context.push('/chat/${t.user.id}'),
                      leading: LumoraAvatar(name: t.user.name, avatarColor: t.user.avatarColor, avatarUrl: t.user.avatarUrl, size: 48),
                      title: Row(children: [
                        Expanded(child: Text(t.user.name, maxLines: 1, overflow: TextOverflow.ellipsis, style: const TextStyle(fontWeight: FontWeight.w800))),
                        Text(_timeAgo(t.lastAt), style: const TextStyle(color: LumoraColors.gray500, fontSize: 11)),
                      ]),
                      subtitle: Text(t.lastMessage, maxLines: 1, overflow: TextOverflow.ellipsis,
                          style: TextStyle(color: t.unread > 0 ? LumoraColors.ink : LumoraColors.slatey, fontWeight: t.unread > 0 ? FontWeight.w700 : FontWeight.w400)),
                      trailing: t.unread > 0
                          ? Container(
                              padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                              decoration: BoxDecoration(color: LumoraColors.coral, borderRadius: BorderRadius.circular(LumoraRadii.full)),
                              child: Text(t.unread > 9 ? '9+' : '${t.unread}', style: const TextStyle(color: Colors.white, fontSize: 11, fontWeight: FontWeight.w800)),
                            )
                          : null,
                    );
                  },
                ),
    );
  }
}

class _EmptyState extends StatefulWidget {
  final VoidCallback onStart;
  const _EmptyState({required this.onStart});
  @override
  State<_EmptyState> createState() => _EmptyStateState();
}

class _EmptyStateState extends State<_EmptyState> {
  List<ChatUser>? _contacts;

  @override
  void initState() {
    super.initState();
    ApiClient.instance.chatContacts().then((c) {
      if (mounted) setState(() => _contacts = c);
    }).catchError((_) {
      if (mounted) setState(() => _contacts = []);
    });
  }

  @override
  Widget build(BuildContext context) {
    final shown = (_contacts ?? []).take(5).toList();
    final extra = ((_contacts?.length ?? 0) - shown.length).clamp(0, 1 << 31);

    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(mainAxisSize: MainAxisSize.min, children: [
          if (shown.isNotEmpty)
            SizedBox(
              height: 56,
              child: Stack(
                children: [
                  for (var i = 0; i < shown.length; i++)
                    Positioned(
                      left: i * 42.0,
                      child: GestureDetector(
                        onTap: () => context.push('/chat/${shown[i].id}'),
                        child: Container(
                          decoration: BoxDecoration(shape: BoxShape.circle, border: Border.all(color: LumoraColors.cream, width: 3)),
                          child: LumoraAvatar(name: shown[i].name, avatarColor: shown[i].avatarColor, avatarUrl: shown[i].avatarUrl, size: 56),
                        ),
                      ),
                    ),
                  if (extra > 0)
                    Positioned(
                      left: shown.length * 42.0,
                      child: Container(
                        width: 56, height: 56,
                        decoration: BoxDecoration(color: LumoraColors.purpleLight, shape: BoxShape.circle, border: Border.all(color: LumoraColors.cream, width: 3)),
                        child: Center(child: Text('+$extra', style: const TextStyle(color: LumoraColors.purple, fontWeight: FontWeight.w800))),
                      ),
                    ),
                ],
              ),
            )
          else
            Container(width: 64, height: 64, decoration: BoxDecoration(color: LumoraColors.gray50, shape: BoxShape.circle), child: const Icon(Icons.forum_outlined, color: LumoraColors.gray300, size: 28)),
          const SizedBox(height: 20),
          Text(shown.isNotEmpty ? 'No conversations yet' : 'Nobody else is here yet', style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 16)),
          const SizedBox(height: 4),
          Text(
            shown.isNotEmpty
                ? '${shown.length + extra} learners are learning alongside you. Tap a face above, or start a new conversation.'
                : "As soon as other learners join, they'll show up here and you can say hello.",
            textAlign: TextAlign.center,
            style: const TextStyle(color: LumoraColors.slatey),
          ),
        ]),
      ),
    );
  }
}

class _ContactPicker extends StatefulWidget {
  const _ContactPicker();
  @override
  State<_ContactPicker> createState() => _ContactPickerState();
}

class _ContactPickerState extends State<_ContactPicker> {
  List<ChatUser>? _contacts;
  String _q = '';

  @override
  void initState() {
    super.initState();
    ApiClient.instance.chatContacts().then((c) {
      if (mounted) setState(() => _contacts = c);
    }).catchError((_) {
      if (mounted) setState(() => _contacts = []);
    });
  }

  @override
  Widget build(BuildContext context) {
    final filtered = (_contacts ?? []).where((c) => c.name.toLowerCase().contains(_q.toLowerCase())).toList();
    return SafeArea(
      child: Padding(
        padding: EdgeInsets.only(bottom: MediaQuery.of(context).viewInsets.bottom),
        child: SizedBox(
          height: MediaQuery.of(context).size.height * 0.7,
          child: Column(
            children: [
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 16, 12, 8),
                child: Row(children: [
                  const Expanded(child: Text('New message', style: TextStyle(fontSize: 18, fontWeight: FontWeight.w800))),
                  IconButton(icon: const Icon(Icons.close), onPressed: () => Navigator.pop(context)),
                ]),
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16),
                child: TextField(
                  onChanged: (v) => setState(() => _q = v),
                  decoration: const InputDecoration(hintText: 'Search learners…', prefixIcon: Icon(Icons.search)),
                ),
              ),
              const SizedBox(height: 8),
              Expanded(
                child: _contacts == null
                    ? const Center(child: CircularProgressIndicator())
                    : filtered.isEmpty
                        ? const Center(child: Text('No learners found.', style: TextStyle(color: LumoraColors.slatey)))
                        : ListView.builder(
                            itemCount: filtered.length,
                            itemBuilder: (context, i) {
                              final c = filtered[i];
                              return ListTile(
                                onTap: () {
                                  Navigator.pop(context);
                                  context.push('/chat/${c.id}');
                                },
                                leading: LumoraAvatar(name: c.name, avatarColor: c.avatarColor, avatarUrl: c.avatarUrl, size: 40),
                                title: Text(c.name, style: const TextStyle(fontWeight: FontWeight.w800)),
                                subtitle: Text(c.levelName, style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
                              );
                            },
                          ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
