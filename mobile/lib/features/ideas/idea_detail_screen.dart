import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../core/network/api_client.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../core/theme/shadows.dart';
import '../../models/idea.dart';
import '../../widgets/avatar.dart';
import '../../widgets/lumora_button.dart';

class IdeaDetailScreen extends StatefulWidget {
  final int id;
  const IdeaDetailScreen({super.key, required this.id});

  @override
  State<IdeaDetailScreen> createState() => _IdeaDetailScreenState();
}

class _IdeaDetailScreenState extends State<IdeaDetailScreen> {
  IdeaDetail? _detail;
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() => _loading = true);
    try {
      final d = await ApiClient.instance.idea(widget.id);
      if (mounted) setState(() => _detail = d);
    } catch (_) {
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _changeStatus(IdeaStatus status) async {
    try {
      await ApiClient.instance.updateIdea(widget.id, status: status);
      _load();
    } catch (_) {}
  }

  Future<void> _addTask() async {
    final controller = TextEditingController();
    final result = await showDialog<String>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('New task'),
        content: TextField(controller: controller, decoration: const InputDecoration(hintText: 'Task title')),
        actions: [
          TextButton(onPressed: () => Navigator.pop(context), child: const Text('Cancel')),
          TextButton(onPressed: () => Navigator.pop(context, controller.text.trim()), child: const Text('Add')),
        ],
      ),
    );
    if (result != null && result.isNotEmpty) {
      try {
        await ApiClient.instance.createIdeaTask(widget.id, title: result);
        _load();
      } catch (_) {}
    }
  }

  Future<void> _toggleTask(IdeaTask task) async {
    final next = task.status == 'done' ? 'todo' : 'done';
    try {
      await ApiClient.instance.updateIdeaTask(task.id, status: next);
      _load();
    } catch (_) {}
  }

  Future<void> _archive() async {
    final controller = TextEditingController();
    final reason = await showDialog<String>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Archive idea'),
        content: TextField(controller: controller, decoration: const InputDecoration(hintText: 'Reason')),
        actions: [
          TextButton(onPressed: () => Navigator.pop(context), child: const Text('Cancel')),
          TextButton(onPressed: () => Navigator.pop(context, controller.text.trim()), child: const Text('Archive')),
        ],
      ),
    );
    if (reason != null && reason.isNotEmpty) {
      try {
        await ApiClient.instance.archiveIdea(widget.id, reason);
        _load();
      } catch (_) {}
    }
  }

  @override
  Widget build(BuildContext context) {
    final detail = _detail;
    return Scaffold(
      backgroundColor: LumoraColors.cream,
      appBar: AppBar(
        backgroundColor: Colors.white,
        foregroundColor: LumoraColors.ink,
        elevation: 0,
        title: const Text('Idea', style: TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.ink)),
        actions: [
          if (detail?.canEdit == true)
            PopupMenuButton<String>(
              onSelected: (v) {
                if (v == 'archive') _archive();
                if (v == 'restore') ApiClient.instance.restoreIdea(widget.id).then((_) => _load());
              },
              itemBuilder: (context) => [
                if (!detail!.idea.archived) const PopupMenuItem(value: 'archive', child: Text('Archive')),
                if (detail.idea.archived) const PopupMenuItem(value: 'restore', child: Text('Restore')),
              ],
            ),
        ],
      ),
      body: _loading || detail == null
          ? const Center(child: CircularProgressIndicator())
          : RefreshIndicator(
              onRefresh: _load,
              child: ListView(
                padding: const EdgeInsets.all(16),
                children: [
                  Container(
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(detail.idea.title, style: const TextStyle(fontSize: 20, fontWeight: FontWeight.w800)),
                        const SizedBox(height: 8),
                        Row(children: [
                          LumoraAvatar(name: detail.idea.owner.name, avatarColor: detail.idea.owner.avatarColor, avatarUrl: detail.idea.owner.avatarUrl, size: 24),
                          const SizedBox(width: 8),
                          Text(detail.idea.owner.name, style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
                        ]),
                        if (detail.idea.description.isNotEmpty) ...[
                          const SizedBox(height: 12),
                          Text(detail.idea.description),
                        ],
                        if (detail.idea.tags.isNotEmpty) ...[
                          const SizedBox(height: 10),
                          Wrap(spacing: 6, children: [for (final t in detail.idea.tags) Chip(label: Text(t), visualDensity: VisualDensity.compact)]),
                        ],
                      ],
                    ),
                  ),
                  const SizedBox(height: 16),
                  const Text('STATUS', style: TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.gray500)),
                  const SizedBox(height: 8),
                  Wrap(spacing: 8, runSpacing: 8, children: [
                    for (final s in detail.statusFlow)
                      ChoiceChip(label: Text(ideaStatusLabel(s)), selected: detail.idea.status == s, onSelected: (_) => _changeStatus(s)),
                  ]),
                  const SizedBox(height: 20),
                  Row(children: [
                    const Text('TASKS', style: TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.gray500)),
                    const Spacer(),
                    TextButton.icon(onPressed: _addTask, icon: const Icon(Icons.add, size: 16), label: const Text('Add')),
                  ]),
                  for (final t in detail.tasks)
                    Card(
                      margin: const EdgeInsets.only(bottom: 6),
                      child: CheckboxListTile(
                        value: t.status == 'done',
                        onChanged: (_) => _toggleTask(t),
                        title: Text(t.title, style: TextStyle(decoration: t.status == 'done' ? TextDecoration.lineThrough : null)),
                        controlAffinity: ListTileControlAffinity.leading,
                      ),
                    ),
                  if (detail.tasks.isEmpty) const Padding(padding: EdgeInsets.symmetric(vertical: 8), child: Text('No tasks yet.', style: TextStyle(color: LumoraColors.slatey))),
                  const SizedBox(height: 20),
                  LumoraButton(label: 'Open discussion (${detail.idea.messageCount})', full: true, onPressed: () => context.push('/ideas/${widget.id}/thread')),
                  if (detail.history.isNotEmpty) ...[
                    const SizedBox(height: 24),
                    const Text('HISTORY', style: TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.gray500)),
                    const SizedBox(height: 8),
                    for (final ev in detail.history.take(10))
                      Padding(
                        padding: const EdgeInsets.symmetric(vertical: 4),
                        child: Text('${ev.actor.name} ${ev.kind} ${ev.field.isNotEmpty ? ev.field : ""}', style: const TextStyle(fontSize: 12, color: LumoraColors.slatey)),
                      ),
                  ],
                ],
              ),
            ),
    );
  }
}
