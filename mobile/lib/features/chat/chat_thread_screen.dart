import 'dart:async';

import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:image_picker/image_picker.dart';

import '../../core/env.dart';
import '../../core/network/api_client.dart';
import '../../core/network/api_exception.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../models/chat.dart';
import '../../providers/auth_provider.dart';
import '../../widgets/avatar.dart';
import '../../widgets/fox_mascot.dart';

class ChatThreadScreen extends ConsumerStatefulWidget {
  final int userId;
  const ChatThreadScreen({super.key, required this.userId});

  @override
  ConsumerState<ChatThreadScreen> createState() => _ChatThreadScreenState();
}

class _ChatThreadScreenState extends ConsumerState<ChatThreadScreen> {
  ChatUser? _other;
  List<ChatMessage> _messages = [];
  bool _loaded = false;
  bool _sending = false;
  String? _error;
  int? _editingId;
  final _textController = TextEditingController();
  final _scrollController = ScrollController();
  Timer? _poll;

  @override
  void initState() {
    super.initState();
    _load();
    _poll = Timer.periodic(const Duration(seconds: 4), (_) => _load());
  }

  @override
  void dispose() {
    _poll?.cancel();
    _textController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  Future<void> _load() async {
    try {
      final (messages, user) = await ApiClient.instance.chatMessages(widget.userId);
      if (!mounted) return;
      final grew = messages.length != _messages.length;
      setState(() {
        _other = user;
        _messages = messages;
        _loaded = true;
      });
      if (grew) {
        WidgetsBinding.instance.addPostFrameCallback((_) {
          if (_scrollController.hasClients) {
            _scrollController.animateTo(_scrollController.position.maxScrollExtent, duration: const Duration(milliseconds: 250), curve: Curves.easeOut);
          }
        });
      }
    } catch (_) {
      if (mounted) setState(() => _loaded = true);
    }
  }

  Future<void> _send() async {
    final body = _textController.text.trim();
    if (body.isEmpty || _sending) return;
    setState(() {
      _sending = true;
      _textController.clear();
      _error = null;
    });
    try {
      final msg = await ApiClient.instance.sendChatMessage(widget.userId, body);
      setState(() => _messages = [..._messages, msg]);
    } on ApiException catch (e) {
      _textController.text = body;
      setState(() => _error = e.message);
    } finally {
      if (mounted) setState(() => _sending = false);
    }
  }

  Future<void> _sendImage() async {
    final picker = ImagePicker();
    final file = await picker.pickImage(source: ImageSource.gallery, imageQuality: 85);
    if (file == null) return;
    final bytes = await file.readAsBytes();
    setState(() {
      _sending = true;
      _error = null;
    });
    try {
      final caption = _textController.text.trim();
      final msg = await ApiClient.instance.sendChatImage(widget.userId, UploadFile(bytes, file.name), caption);
      setState(() {
        _messages = [..._messages, msg];
        _textController.clear();
      });
    } on ApiException catch (e) {
      setState(() => _error = e.message);
    } finally {
      if (mounted) setState(() => _sending = false);
    }
  }

  Future<void> _saveEdit(int messageId, String body) async {
    try {
      final msg = await ApiClient.instance.editChatMessage(messageId, body);
      setState(() {
        _messages = [for (final m in _messages) m.id == messageId ? msg : m];
        _editingId = null;
      });
    } on ApiException catch (e) {
      setState(() => _error = e.message);
    }
  }

  Future<void> _delete(int messageId) async {
    try {
      final msg = await ApiClient.instance.deleteChatMessage(messageId);
      setState(() => _messages = [for (final m in _messages) m.id == messageId ? msg : m]);
    } on ApiException catch (e) {
      setState(() => _error = e.message);
    }
  }

  Future<void> _translate(int messageId) async {
    try {
      final t = await ApiClient.instance.translateChatMessage(messageId);
      if (t == null) return;
      setState(() => _messages = [
            for (final m in _messages)
              if (m.id == messageId)
                ChatMessage(id: m.id, senderId: m.senderId, recipientId: m.recipientId, kind: m.kind, body: m.body,
                    url: m.url, fileName: m.fileName, width: m.width, height: m.height, read: m.read, mine: m.mine,
                    edited: m.edited, deleted: m.deleted, canEdit: m.canEdit, createdAt: m.createdAt, translation: t)
              else
                m,
          ]);
    } on ApiException catch (e) {
      setState(() => _error = e.message);
    }
  }

  @override
  Widget build(BuildContext context) {
    final myId = ref.watch(authProvider).user?.id;
    final other = _other;

    return Scaffold(
      backgroundColor: LumoraColors.cream,
      appBar: AppBar(
        backgroundColor: Colors.white,
        foregroundColor: LumoraColors.ink,
        elevation: 0,
        leading: BackButton(onPressed: () => context.canPop() ? context.pop() : context.go('/chat')),
        titleSpacing: 0,
        title: other == null
            ? null
            : Row(children: [
                LumoraAvatar(name: other.name, avatarColor: other.avatarColor, avatarUrl: other.avatarUrl, size: 36),
                const SizedBox(width: 10),
                Expanded(
                  child: Column(crossAxisAlignment: CrossAxisAlignment.start, mainAxisSize: MainAxisSize.min, children: [
                    Text(other.name, maxLines: 1, overflow: TextOverflow.ellipsis, style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 15)),
                    Text(other.levelName.isEmpty ? 'Learner' : other.levelName, style: const TextStyle(color: LumoraColors.slatey, fontSize: 11)),
                  ]),
                ),
              ]),
      ),
      body: Column(
        children: [
          Expanded(
            child: !_loaded
                ? const Center(child: FoxMascot(size: 90, glow: true))
                : _messages.isEmpty
                    ? Center(
                        child: Column(mainAxisSize: MainAxisSize.min, children: [
                          Text('Say hi to ${other?.name ?? "your new friend"} 👋', style: const TextStyle(fontWeight: FontWeight.w700)),
                          const Text('Practising together makes it stick.', style: TextStyle(color: LumoraColors.slatey)),
                        ]),
                      )
                    : ListView.builder(
                        controller: _scrollController,
                        padding: const EdgeInsets.all(16),
                        itemCount: _messages.length,
                        itemBuilder: (context, i) {
                          final m = _messages[i];
                          return _Bubble(
                            message: m,
                            mine: m.senderId == myId,
                            editing: _editingId == m.id,
                            onStartEdit: () => setState(() => _editingId = m.id),
                            onCancelEdit: () => setState(() => _editingId = null),
                            onSaveEdit: (body) => _saveEdit(m.id, body),
                            onDelete: () => _delete(m.id),
                            onTranslate: () => _translate(m.id),
                          );
                        },
                      ),
          ),
          if (_error != null)
            Container(
              margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
              decoration: BoxDecoration(color: LumoraColors.coralLight, borderRadius: BorderRadius.circular(LumoraRadii.md)),
              child: Text(_error!, style: const TextStyle(color: LumoraColors.coral, fontSize: 12)),
            ),
          SafeArea(
            top: false,
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
              decoration: const BoxDecoration(color: Colors.white, border: Border(top: BorderSide(color: LumoraColors.gray100))),
              child: Row(children: [
                IconButton(icon: const Icon(Icons.image_outlined, color: LumoraColors.slatey), onPressed: _sending ? null : _sendImage),
                Expanded(
                  child: TextField(
                    controller: _textController,
                    onSubmitted: (_) => _send(),
                    decoration: InputDecoration(
                      hintText: 'Type a message…',
                      filled: true,
                      fillColor: LumoraColors.gray50,
                      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
                      border: OutlineInputBorder(borderRadius: BorderRadius.circular(LumoraRadii.full), borderSide: BorderSide.none),
                    ),
                  ),
                ),
                const SizedBox(width: 8),
                InkWell(
                  onTap: _sending ? null : _send,
                  child: Container(
                    width: 44, height: 44,
                    decoration: BoxDecoration(color: LumoraColors.purple, shape: BoxShape.circle),
                    child: _sending
                        ? const Padding(padding: EdgeInsets.all(12), child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white))
                        : const Icon(Icons.send_rounded, color: Colors.white, size: 18),
                  ),
                ),
              ]),
            ),
          ),
        ],
      ),
    );
  }
}

class _Bubble extends StatefulWidget {
  final ChatMessage message;
  final bool mine;
  final bool editing;
  final VoidCallback onStartEdit;
  final VoidCallback onCancelEdit;
  final ValueChanged<String> onSaveEdit;
  final VoidCallback onDelete;
  final VoidCallback onTranslate;

  const _Bubble({
    required this.message, required this.mine, required this.editing, required this.onStartEdit,
    required this.onCancelEdit, required this.onSaveEdit, required this.onDelete, required this.onTranslate,
  });

  @override
  State<_Bubble> createState() => _BubbleState();
}

class _BubbleState extends State<_Bubble> {
  late final TextEditingController _draft = TextEditingController(text: widget.message.body);

  @override
  Widget build(BuildContext context) {
    final m = widget.message;
    final align = widget.mine ? Alignment.centerRight : Alignment.centerLeft;

    if (m.deleted) {
      return Align(
        alignment: align,
        child: Container(
          margin: const EdgeInsets.symmetric(vertical: 4),
          padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 8),
          decoration: BoxDecoration(color: LumoraColors.gray100, borderRadius: BorderRadius.circular(LumoraRadii.lg)),
          child: const Row(mainAxisSize: MainAxisSize.min, children: [
            Icon(Icons.delete_outline, size: 14, color: LumoraColors.gray500),
            SizedBox(width: 6),
            Text('Message deleted', style: TextStyle(fontStyle: FontStyle.italic, color: LumoraColors.gray500, fontSize: 13)),
          ]),
        ),
      );
    }

    if (widget.editing) {
      return Align(
        alignment: Alignment.centerRight,
        child: Container(
          width: MediaQuery.of(context).size.width * 0.8,
          margin: const EdgeInsets.symmetric(vertical: 4),
          padding: const EdgeInsets.all(8),
          decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.lg), boxShadow: const [BoxShadow(color: Color(0x14000000), blurRadius: 8)]),
          child: Column(children: [
            TextField(controller: _draft, maxLines: 2, autofocus: true, decoration: const InputDecoration(filled: true, fillColor: LumoraColors.gray50)),
            Row(mainAxisAlignment: MainAxisAlignment.end, children: [
              TextButton(onPressed: widget.onCancelEdit, child: const Text('Cancel')),
              TextButton(onPressed: () => widget.onSaveEdit(_draft.text.trim()), child: const Text('Save')),
            ]),
          ]),
        ),
      );
    }

    return Align(
      alignment: align,
      child: GestureDetector(
        onLongPress: widget.mine ? () => _showActions(context) : null,
        child: Container(
          constraints: BoxConstraints(maxWidth: MediaQuery.of(context).size.width * 0.78),
          margin: const EdgeInsets.symmetric(vertical: 4),
          padding: m.kind == 'image' ? const EdgeInsets.all(4) : const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
          decoration: BoxDecoration(
            color: widget.mine ? LumoraColors.purple : Colors.white,
            borderRadius: BorderRadius.only(
              topLeft: const Radius.circular(18), topRight: const Radius.circular(18),
              bottomLeft: Radius.circular(widget.mine ? 18 : 4),
              bottomRight: Radius.circular(widget.mine ? 4 : 18),
            ),
            boxShadow: widget.mine ? null : const [BoxShadow(color: Color(0x14000000), blurRadius: 6)],
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisSize: MainAxisSize.min,
            children: [
              if (m.kind == 'image')
                ClipRRect(
                  borderRadius: BorderRadius.circular(14),
                  child: CachedNetworkImage(imageUrl: mediaUrl(m.url), fit: BoxFit.contain, height: 220,
                      placeholder: (_, _) => const SizedBox(height: 120, child: Center(child: CircularProgressIndicator()))),
                ),
              if (m.body.isNotEmpty)
                Padding(
                  padding: m.kind == 'image' ? const EdgeInsets.fromLTRB(8, 6, 8, 4) : EdgeInsets.zero,
                  child: Text(m.body, style: TextStyle(color: widget.mine ? Colors.white : LumoraColors.ink, fontSize: 15)),
                ),
              if (m.translation != null)
                Padding(
                  padding: const EdgeInsets.only(top: 4),
                  child: m.translation!.pending
                      ? Text('Translating…', style: TextStyle(color: widget.mine ? Colors.white70 : LumoraColors.slatey, fontStyle: FontStyle.italic, fontSize: 12))
                      : Text('${m.translation!.langName}: ${m.translation!.text}',
                          style: TextStyle(color: widget.mine ? Colors.white70 : LumoraColors.slatey, fontStyle: FontStyle.italic, fontSize: 12)),
                ),
              if (m.edited)
                Padding(
                  padding: const EdgeInsets.only(top: 2),
                  child: Text('edited', style: TextStyle(color: widget.mine ? Colors.white60 : LumoraColors.gray500, fontSize: 10, fontStyle: FontStyle.italic)),
                ),
            ],
          ),
        ),
      ),
    );
  }

  void _showActions(BuildContext context) {
    showModalBottomSheet(
      context: context,
      builder: (_) => SafeArea(
        child: Column(mainAxisSize: MainAxisSize.min, children: [
          if (widget.message.canEdit)
            ListTile(leading: const Icon(Icons.edit), title: const Text('Edit'), onTap: () { Navigator.pop(context); widget.onStartEdit(); })
          else
            ListTile(
              leading: const Icon(Icons.info_outline, color: LumoraColors.gray500),
              title: Text(widget.message.kind == 'image' ? "Photos can't be edited." : 'Edits close after 24 hours.', style: const TextStyle(color: LumoraColors.gray500)),
            ),
          ListTile(leading: const Icon(Icons.delete, color: LumoraColors.coral), title: const Text('Delete', style: TextStyle(color: LumoraColors.coral)),
              onTap: () { Navigator.pop(context); widget.onDelete(); }),
        ]),
      ),
    );
  }
}
