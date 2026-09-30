import 'dart:async';

import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';

import '../../core/network/api_client.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../models/idea.dart';
import '../../widgets/avatar.dart';

const _kReactionEmoji = ['👍', '💡', '🔥', '❓', '🎉', '👀'];

class IdeaThreadScreen extends StatefulWidget {
  final int id;
  const IdeaThreadScreen({super.key, required this.id});

  @override
  State<IdeaThreadScreen> createState() => _IdeaThreadScreenState();
}

class _IdeaThreadScreenState extends State<IdeaThreadScreen> {
  IdeaThread? _thread;
  bool _loading = true;
  final _textController = TextEditingController();
  bool _sending = false;
  Timer? _poll;

  @override
  void initState() {
    super.initState();
    _load();
    _poll = Timer.periodic(const Duration(seconds: 6), (_) => _load(silent: true));
  }

  @override
  void dispose() {
    _poll?.cancel();
    _textController.dispose();
    super.dispose();
  }

  Future<void> _load({bool silent = false}) async {
    if (!silent) setState(() => _loading = true);
    try {
      final t = await ApiClient.instance.ideaMessages(widget.id);
      if (mounted) setState(() => _thread = t);
    } catch (_) {
    } finally {
      if (mounted && !silent) setState(() => _loading = false);
    }
  }

  Future<void> _send() async {
    final body = _textController.text.trim();
    if (body.isEmpty || _sending) return;
    setState(() {
      _sending = true;
      _textController.clear();
    });
    try {
      await ApiClient.instance.postIdeaMessage(widget.id, body: body);
      _load();
    } catch (_) {
    } finally {
      if (mounted) setState(() => _sending = false);
    }
  }

  Future<void> _sendImage() async {
    final picker = ImagePicker();
    final file = await picker.pickImage(source: ImageSource.gallery, imageQuality: 85);
    if (file == null) return;
    final bytes = await file.readAsBytes();
    setState(() => _sending = true);
    try {
      await ApiClient.instance.postIdeaAttachment(widget.id, UploadFile(bytes, file.name), kind: 'image');
      _load();
    } catch (_) {
    } finally {
      if (mounted) setState(() => _sending = false);
    }
  }

  Future<void> _react(int messageId, String emoji) async {
    try {
      await ApiClient.instance.reactToIdeaMessage(messageId, emoji);
      _load(silent: true);
    } catch (_) {}
  }

  Future<void> _startBrainstorm() async {
    try {
      await ApiClient.instance.startBrainstorm(widget.id, 10, 'Quick ideas');
      _load();
    } catch (_) {}
  }

  Future<void> _stopBrainstorm() async {
    try {
      await ApiClient.instance.stopBrainstorm(widget.id);
      _load();
    } catch (_) {}
  }

  Future<void> _showSummary() async {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.white,
      shape: const RoundedRectangleBorder(borderRadius: BorderRadius.vertical(top: Radius.circular(24))),
      builder: (context) => FutureBuilder<ThreadSummary>(
        future: ApiClient.instance.ideaSummary(widget.id),
        builder: (context, snapshot) {
          return SafeArea(
            child: Padding(
              padding: const EdgeInsets.all(24),
              child: !snapshot.hasData
                  ? const SizedBox(height: 120, child: Center(child: CircularProgressIndicator()))
                  : Column(
                      mainAxisSize: MainAxisSize.min,
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Text('Thread summary', style: TextStyle(fontSize: 18, fontWeight: FontWeight.w800)),
                        const SizedBox(height: 12),
                        if (!snapshot.data!.generated)
                          const Text('The thread is too thin to summarise yet.', style: TextStyle(color: LumoraColors.slatey))
                        else ...[
                          Text(snapshot.data!.gist),
                          if (snapshot.data!.keyPoints.isNotEmpty) ...[
                            const SizedBox(height: 16),
                            const Text('KEY POINTS', style: TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.gray500)),
                            for (final p in snapshot.data!.keyPoints)
                              Padding(padding: const EdgeInsets.only(top: 6), child: Text('• ${p.text} — ${p.author.name}', style: const TextStyle(fontSize: 13))),
                          ],
                          if (snapshot.data!.questions.isNotEmpty) ...[
                            const SizedBox(height: 16),
                            const Text('UNANSWERED QUESTIONS', style: TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.gray500)),
                            for (final q in snapshot.data!.questions)
                              Padding(padding: const EdgeInsets.only(top: 6), child: Text('• ${q.text}', style: const TextStyle(fontSize: 13))),
                          ],
                        ],
                      ],
                    ),
            ),
          );
        },
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final thread = _thread;
    return Scaffold(
      backgroundColor: LumoraColors.cream,
      appBar: AppBar(
        backgroundColor: Colors.white,
        foregroundColor: LumoraColors.ink,
        elevation: 0,
        title: Text(thread?.idea.title ?? 'Discussion', maxLines: 1, overflow: TextOverflow.ellipsis, style: const TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.ink)),
        actions: [
          IconButton(
            icon: const Icon(Icons.auto_awesome_outlined),
            tooltip: 'Summarize thread',
            onPressed: _showSummary,
          ),
          IconButton(
            icon: Icon(thread?.brainstorm != null ? Icons.stop_circle_outlined : Icons.psychology_outlined),
            tooltip: thread?.brainstorm != null ? 'Stop silent brainstorm' : 'Start silent brainstorm',
            onPressed: thread?.brainstorm != null ? _stopBrainstorm : _startBrainstorm,
          ),
        ],
      ),
      body: Column(
        children: [
          if (thread?.brainstorm != null)
            Container(
              width: double.infinity,
              color: LumoraColors.purpleLight,
              padding: const EdgeInsets.all(10),
              child: Text('🤫 Silent brainstorm running: "${thread!.brainstorm!.topic}" — posts stay anonymous until it ends.',
                  style: const TextStyle(fontSize: 12, color: LumoraColors.purple, fontWeight: FontWeight.w700)),
            ),
          Expanded(
            child: _loading
                ? const Center(child: CircularProgressIndicator())
                : thread == null || thread.messages.isEmpty
                    ? const Center(child: Text('No messages yet. Start the discussion!', style: TextStyle(color: LumoraColors.slatey)))
                    : ListView.builder(
                        padding: const EdgeInsets.all(16),
                        itemCount: thread.messages.length,
                        itemBuilder: (context, i) => _MessageTile(message: thread.messages[i], onReact: _react),
                      ),
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
                      hintText: 'Add to the discussion…',
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
                  child: Container(width: 44, height: 44, decoration: const BoxDecoration(color: LumoraColors.purple, shape: BoxShape.circle),
                      child: const Icon(Icons.send_rounded, color: Colors.white, size: 18)),
                ),
              ]),
            ),
          ),
        ],
      ),
    );
  }
}

class _MessageTile extends StatelessWidget {
  final IdeaMessage message;
  final Future<void> Function(int, String) onReact;
  const _MessageTile({required this.message, required this.onReact});

  @override
  Widget build(BuildContext context) {
    if (message.deleted) {
      return const Padding(
        padding: EdgeInsets.symmetric(vertical: 6),
        child: Text('Message deleted', style: TextStyle(fontStyle: FontStyle.italic, color: LumoraColors.gray500, fontSize: 12)),
      );
    }
    final authorName = message.anonymous ? 'Anonymous' : (message.author?.name ?? 'Unknown');
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          message.anonymous
              ? const CircleAvatar(radius: 16, backgroundColor: LumoraColors.gray300, child: Icon(Icons.person, size: 16, color: Colors.white))
              : LumoraAvatar(name: authorName, avatarColor: message.author?.avatarColor ?? '#6C3FC5', avatarUrl: message.author?.avatarUrl, size: 32),
          const SizedBox(width: 10),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(children: [
                  Text(authorName, style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 13)),
                  if (message.edited) const Padding(padding: EdgeInsets.only(left: 6), child: Text('edited', style: TextStyle(fontSize: 10, color: LumoraColors.gray500, fontStyle: FontStyle.italic))),
                ]),
                const SizedBox(height: 2),
                if (message.body.isNotEmpty) Text(message.body, style: const TextStyle(fontSize: 14)),
                if (message.translation != null && !message.translation!.pending)
                  Padding(
                    padding: const EdgeInsets.only(top: 2),
                    child: Text('${message.translation!.langName}: ${message.translation!.text}', style: const TextStyle(fontSize: 12, color: LumoraColors.slatey, fontStyle: FontStyle.italic)),
                  ),
                const SizedBox(height: 6),
                Wrap(spacing: 6, children: [
                  for (final r in message.reactions)
                    ActionChip(
                      label: Text('${r.emoji} ${r.count}', style: const TextStyle(fontSize: 11)),
                      backgroundColor: r.mine ? LumoraColors.purpleLight : LumoraColors.gray50,
                      onPressed: () => onReact(message.id, r.emoji),
                    ),
                  InkWell(
                    onTap: () => _showEmojiPicker(context),
                    child: Container(
                      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                      decoration: BoxDecoration(color: LumoraColors.gray50, borderRadius: BorderRadius.circular(LumoraRadii.full)),
                      child: const Icon(Icons.add_reaction_outlined, size: 14, color: LumoraColors.gray500),
                    ),
                  ),
                ]),
                if (message.replyCount > 0)
                  Padding(padding: const EdgeInsets.only(top: 4), child: Text('${message.replyCount} replies', style: const TextStyle(fontSize: 11, color: LumoraColors.purple, fontWeight: FontWeight.w700))),
              ],
            ),
          ),
        ],
      ),
    );
  }

  void _showEmojiPicker(BuildContext context) {
    showModalBottomSheet(
      context: context,
      builder: (_) => SafeArea(
        child: Wrap(
          children: [
            for (final e in _kReactionEmoji)
              InkWell(
                onTap: () { Navigator.pop(context); onReact(message.id, e); },
                child: Padding(padding: const EdgeInsets.all(16), child: Text(e, style: const TextStyle(fontSize: 28))),
              ),
          ],
        ),
      ),
    );
  }
}
