import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../core/network/api_client.dart';
import '../../core/network/api_exception.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../core/theme/shadows.dart';
import '../../models/idea.dart';
import '../../widgets/avatar.dart';
import '../../widgets/confirm_dialog.dart';
import '../../widgets/lumora_button.dart';
import 'idea_status_style.dart';

/// Everything about one idea that isn't conversation — where it stands in the
/// workflow, how it got there, and what to do with it next. Mirrors
/// frontend/components/ideas/IdeaDetailsPanel.tsx.
class IdeaDetailScreen extends StatefulWidget {
  final int id;
  const IdeaDetailScreen({super.key, required this.id});

  @override
  State<IdeaDetailScreen> createState() => _IdeaDetailScreenState();
}

class _IdeaDetailScreenState extends State<IdeaDetailScreen> {
  IdeaDetail? _detail;
  bool _loading = true;
  bool _busy = false;
  String? _loadError;

  @override
  void initState() {
    super.initState();
    _load(spinner: true);
  }

  Future<void> _load({bool spinner = false}) async {
    if (spinner) setState(() => _loading = true);
    try {
      final d = await ApiClient.instance.idea(widget.id);
      if (mounted) setState(() { _detail = d; _loadError = null; });
    } on ApiException catch (e) {
      if (mounted) setState(() => _loadError = e.message);
    } catch (_) {
      if (mounted) setState(() => _loadError = "Couldn't load this idea.");
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  /// Runs an action, reloads, and surfaces the server's reason when it's
  /// refused (e.g. "You can't approve your own idea").
  Future<void> _act(Future<void> Function() fn, {String? done}) async {
    if (_busy) return;
    setState(() => _busy = true);
    try {
      await fn();
      await _load();
      if (done != null) _toast(done);
    } on ApiException catch (e) {
      _toast(e.message, error: true);
    } catch (_) {
      _toast("That didn't work — try again.", error: true);
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  void _toast(String msg, {bool error = false}) {
    if (!mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      content: Text(msg),
      backgroundColor: error ? LumoraColors.coral : LumoraColors.ink,
    ));
  }

  Future<String?> _prompt({
    required String title,
    String? message,
    String hint = '',
    required String confirm,
    bool required = false,
    TextInputType? keyboard,
  }) {
    final controller = TextEditingController();
    return showDialog<String>(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setLocal) => AlertDialog(
          title: Text(title),
          content: Column(mainAxisSize: MainAxisSize.min, crossAxisAlignment: CrossAxisAlignment.start, children: [
            if (message != null) ...[
              Text(message, style: const TextStyle(color: LumoraColors.slatey, fontSize: 13)),
              const SizedBox(height: 8),
            ],
            TextField(
              controller: controller,
              autofocus: true,
              keyboardType: keyboard,
              maxLines: keyboard == TextInputType.number ? 1 : 3,
              minLines: 1,
              onChanged: (_) => setLocal(() {}),
              decoration: InputDecoration(hintText: hint),
            ),
          ]),
          actions: [
            TextButton(onPressed: () => Navigator.pop(context), child: const Text('Cancel')),
            TextButton(
              onPressed: required && controller.text.trim().isEmpty
                  ? null
                  : () => Navigator.pop(context, controller.text.trim()),
              child: Text(confirm),
            ),
          ],
        ),
      ),
    );
  }

  // --- actions ---------------------------------------------------------------

  Future<void> _move(IdeaTransition t) async {
    String? note;
    if (!t.primary) {
      // Backward or sideways moves ask for an optional note, so everyone
      // following the idea learns why.
      note = await _prompt(title: t.label, message: t.hint, hint: 'Add a note (optional)', confirm: t.label);
      if (note == null) return;
    }
    await _act(
      () => ApiClient.instance.updateIdea(widget.id, status: t.to, note: note),
      done: 'Moved to ${ideaStatusLabel(t.to)}',
    );
  }

  Future<void> _vote(int value) async {
    final idea = _detail!.idea;
    await _act(() => ApiClient.instance.voteIdea(widget.id, idea.myVote == value ? 0 : value));
  }

  Future<void> _addTask() async {
    final title = await _prompt(
      title: 'New task',
      hint: 'Task title (defaults to the idea title)',
      confirm: 'Add',
    );
    if (title == null) return;
    await _act(
      () => ApiClient.instance.createIdeaTask(widget.id, title: title.isEmpty ? null : title),
      done: 'Task added',
    );
  }

  Future<void> _toggleTask(IdeaTask task) =>
      _act(() => ApiClient.instance.updateIdeaTask(task.id, status: task.status == 'done' ? 'todo' : 'done'));

  Future<void> _archive() async {
    final reason = await _prompt(
      title: 'Archive idea',
      message: "Required — it's what tells everyone who contributed what happened to their idea.",
      hint: 'Duplicate of #12 / out of scope…',
      confirm: 'Archive',
      required: true,
    );
    if (reason == null || reason.isEmpty) return;
    await _act(() => ApiClient.instance.archiveIdea(widget.id, reason), done: 'Idea archived');
  }

  Future<void> _merge() async {
    final raw = await _prompt(
      title: 'Merge into another idea',
      message: 'Votes and the whole thread move across. Anyone who backed both is only counted once.',
      hint: 'Surviving idea number, e.g. 12',
      confirm: 'Merge',
      required: true,
      keyboard: TextInputType.number,
    );
    final target = int.tryParse((raw ?? '').replaceAll('#', ''));
    if (target == null) return;
    await _act(() => ApiClient.instance.mergeIdea(widget.id, target), done: 'Merged into #$target');
  }

  Future<void> _delete() async {
    final detail = _detail!;
    final ok = await showLumoraConfirmDialog(
      context,
      title: 'Delete this idea?',
      message: detail.idea.messageCount > 0
          ? 'This permanently destroys the idea, its ${detail.idea.messageCount} message(s), votes and history. Archiving keeps the discussion readable instead.'
          : "This permanently deletes the idea. It can't be undone.",
      confirmLabel: 'Delete',
      danger: true,
    );
    if (!ok) return;
    try {
      await ApiClient.instance.deleteIdea(widget.id);
      if (mounted) context.pop(true); // true = deleted, so the board says so
    } on ApiException catch (e) {
      _toast(e.message, error: true);
    }
  }

  Future<void> _edit() async {
    final saved = await showModalBottomSheet<bool>(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.white,
      shape: const RoundedRectangleBorder(borderRadius: BorderRadius.vertical(top: Radius.circular(24))),
      builder: (_) => _EditIdeaSheet(idea: _detail!.idea),
    );
    if (saved == true) {
      await _load();
      _toast('Idea updated');
    }
  }

  // --- build -----------------------------------------------------------------

  @override
  Widget build(BuildContext context) {
    final detail = _detail;
    final idea = detail?.idea;
    final closed = idea != null && (idea.archived || idea.mergedIntoId != null);

    return Scaffold(
      backgroundColor: LumoraColors.cream,
      appBar: AppBar(
        backgroundColor: Colors.white,
        foregroundColor: LumoraColors.ink,
        elevation: 0,
        title: Text(idea == null ? 'Idea' : 'Idea #${idea.id}',
            style: const TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.ink)),
        actions: [
          if (idea != null)
            IconButton(
              tooltip: idea.starred ? 'Unstar' : 'Star',
              icon: Icon(idea.starred ? Icons.star_rounded : Icons.star_outline_rounded,
                  color: idea.starred ? LumoraColors.amber : LumoraColors.gray500),
              onPressed: () => _act(() => ApiClient.instance.starIdea(widget.id)),
            ),
          if (detail != null)
            PopupMenuButton<String>(
              onSelected: (v) {
                switch (v) {
                  case 'edit': _edit();
                  case 'merge': _merge();
                  case 'archive': _archive();
                  case 'restore': _act(() => ApiClient.instance.restoreIdea(widget.id), done: 'Idea restored');
                  case 'delete': _delete();
                }
              },
              itemBuilder: (context) => [
                if (detail.canEdit && !closed) const PopupMenuItem(value: 'edit', child: Text('Edit')),
                if (!closed) const PopupMenuItem(value: 'merge', child: Text('Merge into another idea')),
                if (!closed) const PopupMenuItem(value: 'archive', child: Text('Archive')),
                if (idea!.archived && idea.mergedIntoId == null)
                  const PopupMenuItem(value: 'restore', child: Text('Restore to the board')),
                if (detail.canEdit)
                  const PopupMenuItem(value: 'delete', child: Text('Delete', style: TextStyle(color: LumoraColors.coral))),
              ],
            ),
        ],
      ),
      body: _loading
          ? const Center(child: CircularProgressIndicator())
          : detail == null
              ? _ErrorState(message: _loadError ?? "Couldn't load this idea.", onRetry: () => _load(spinner: true))
              : RefreshIndicator(
                  onRefresh: _load,
                  child: ListView(
                    padding: const EdgeInsets.all(16),
                    children: [
                      _HeaderCard(detail: detail, busy: _busy, onVote: _vote),
                      if (closed) ...[
                        const SizedBox(height: 12),
                        _ClosedBanner(idea: detail.idea),
                      ],
                      const SizedBox(height: 16),
                      _StatusCard(detail: detail, busy: _busy, onMove: _move),
                      const SizedBox(height: 16),
                      _TasksCard(detail: detail, busy: _busy, onAdd: _addTask, onToggle: _toggleTask),
                      const SizedBox(height: 16),
                      LumoraButton(
                        label: 'Open discussion (${detail.idea.messageCount})',
                        full: true,
                        onPressed: () async {
                          await context.push('/ideas/${widget.id}/thread');
                          if (mounted) _load();
                        },
                      ),
                      if (detail.participants.isNotEmpty) ...[
                        const SizedBox(height: 20),
                        const _SectionLabel('PEOPLE'),
                        const SizedBox(height: 8),
                        Wrap(spacing: 6, runSpacing: 6, children: [
                          for (final p in detail.participants)
                            Chip(
                              avatar: LumoraAvatar(name: p.name, avatarColor: p.avatarColor, avatarUrl: p.avatarUrl, size: 20),
                              label: Text(p.name, style: const TextStyle(fontSize: 12)),
                              visualDensity: VisualDensity.compact,
                            ),
                        ]),
                      ],
                      if (detail.similar.isNotEmpty) ...[
                        const SizedBox(height: 20),
                        const _SectionLabel('POSSIBLY RELATED'),
                        const SizedBox(height: 8),
                        for (final s in detail.similar)
                          Card(
                            margin: const EdgeInsets.only(bottom: 6),
                            child: ListTile(
                              dense: true,
                              title: Text('#${s.id} ${s.title}', maxLines: 1, overflow: TextOverflow.ellipsis,
                                  style: const TextStyle(fontWeight: FontWeight.w700)),
                              subtitle: Text('${(s.similarity * 100).round()}% overlap · ${s.score} votes'),
                              trailing: IdeaStatusBadge(status: s.status),
                              onTap: () => context.push('/ideas/${s.id}'),
                            ),
                          ),
                      ],
                      if (detail.mergedIn.isNotEmpty) ...[
                        const SizedBox(height: 20),
                        const _SectionLabel('MERGED IN'),
                        const SizedBox(height: 8),
                        for (final m in detail.mergedIn)
                          Text('#${m.id} ${m.title}', style: const TextStyle(color: LumoraColors.slatey, fontSize: 13)),
                      ],
                      if (detail.history.isNotEmpty) ...[
                        const SizedBox(height: 20),
                        const _SectionLabel('HISTORY'),
                        const SizedBox(height: 8),
                        for (final ev in detail.history.take(20)) _HistoryRow(event: ev),
                      ],
                      const SizedBox(height: 24),
                      if (detail.canEdit)
                        LumoraButton(
                          label: 'Delete idea',
                          full: true,
                          variant: LumoraButtonVariant.danger,
                          onPressed: _busy ? null : _delete,
                        )
                      else
                        Text(
                          'Only ${detail.idea.owner.name.isEmpty ? "the person who posted it" : detail.idea.owner.name} can delete this idea. '
                          "Archive it from the ⋮ menu if it's no longer relevant.",
                          textAlign: TextAlign.center,
                          style: const TextStyle(color: LumoraColors.gray500, fontSize: 12),
                        ),
                    ],
                  ),
                ),
    );
  }
}

// --- sections ----------------------------------------------------------------

class _SectionLabel extends StatelessWidget {
  final String text;
  const _SectionLabel(this.text);
  @override
  Widget build(BuildContext context) =>
      Text(text, style: const TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.gray500));
}

BoxDecoration _card() => BoxDecoration(
      color: Colors.white,
      borderRadius: BorderRadius.circular(LumoraRadii.xl),
      boxShadow: LumoraShadows.card,
    );

class _HeaderCard extends StatelessWidget {
  final IdeaDetail detail;
  final bool busy;
  final void Function(int) onVote;
  const _HeaderCard({required this.detail, required this.busy, required this.onVote});

  @override
  Widget build(BuildContext context) {
    final idea = detail.idea;
    final total = idea.upvotes + idea.downvotes;
    final closed = idea.archived || idea.mergedIntoId != null;
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: _card(),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          IdeaStatusBadge(status: idea.status),
          const SizedBox(height: 10),
          Text(idea.title, style: const TextStyle(fontSize: 20, fontWeight: FontWeight.w800)),
          const SizedBox(height: 8),
          Row(children: [
            LumoraAvatar(name: idea.owner.name, avatarColor: idea.owner.avatarColor, avatarUrl: idea.owner.avatarUrl, size: 24),
            const SizedBox(width: 8),
            Expanded(
              child: Text('${idea.owner.name} · ${_ago(idea.createdAt)}',
                  style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
            ),
          ]),
          if (idea.description.isNotEmpty) ...[
            const SizedBox(height: 12),
            Text(idea.description, style: const TextStyle(height: 1.4)),
          ],
          if (idea.tags.isNotEmpty) ...[
            const SizedBox(height: 10),
            Wrap(spacing: 6, runSpacing: 6, children: [
              for (final t in idea.tags)
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                  decoration: BoxDecoration(color: LumoraColors.purpleLight, borderRadius: BorderRadius.circular(LumoraRadii.full)),
                  child: Text('#$t', style: const TextStyle(color: LumoraColors.purple, fontSize: 11, fontWeight: FontWeight.w700)),
                ),
            ]),
          ],
          const SizedBox(height: 14),
          Row(children: [
            _VoteButton(
              icon: Icons.arrow_upward_rounded,
              label: '${idea.upvotes}',
              active: idea.myVote == 1,
              color: LumoraColors.teal,
              onTap: busy || closed ? null : () => onVote(1),
            ),
            const SizedBox(width: 8),
            _VoteButton(
              icon: Icons.arrow_downward_rounded,
              label: '${idea.downvotes}',
              active: idea.myVote == -1,
              color: LumoraColors.coral,
              onTap: busy || closed ? null : () => onVote(-1),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Text('Score ${idea.score}', style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 13)),
                const SizedBox(height: 4),
                ClipRRect(
                  borderRadius: BorderRadius.circular(4),
                  child: LinearProgressIndicator(
                    value: total == 0 ? 0 : idea.upvotes / total,
                    minHeight: 6,
                    backgroundColor: total == 0 ? LumoraColors.gray100 : LumoraColors.coralLight,
                    color: LumoraColors.teal,
                  ),
                ),
              ]),
            ),
          ]),
        ],
      ),
    );
  }
}

class _VoteButton extends StatelessWidget {
  final IconData icon;
  final String label;
  final bool active;
  final Color color;
  final VoidCallback? onTap;
  const _VoteButton({required this.icon, required this.label, required this.active, required this.color, this.onTap});

  @override
  Widget build(BuildContext context) {
    return Material(
      color: active ? color.withValues(alpha: 0.14) : LumoraColors.gray50,
      borderRadius: BorderRadius.circular(LumoraRadii.full),
      child: InkWell(
        borderRadius: BorderRadius.circular(LumoraRadii.full),
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
          child: Row(mainAxisSize: MainAxisSize.min, children: [
            Icon(icon, size: 18, color: active ? color : LumoraColors.slatey),
            const SizedBox(width: 4),
            Text(label, style: TextStyle(fontWeight: FontWeight.w800, color: active ? color : LumoraColors.slatey)),
          ]),
        ),
      ),
    );
  }
}

class _ClosedBanner extends StatelessWidget {
  final Idea idea;
  const _ClosedBanner({required this.idea});

  @override
  Widget build(BuildContext context) {
    final merged = idea.mergedIntoId != null;
    return Container(
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(color: LumoraColors.gray50, borderRadius: BorderRadius.circular(LumoraRadii.lg)),
      child: Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Icon(merged ? Icons.merge_rounded : Icons.inventory_2_rounded, color: LumoraColors.slatey),
        const SizedBox(width: 10),
        Expanded(
          child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Text(merged ? 'Merged' : 'Archived', style: const TextStyle(fontWeight: FontWeight.w800)),
            if (idea.archiveReason.isNotEmpty)
              Text(idea.archiveReason, style: const TextStyle(color: LumoraColors.slatey, fontSize: 13)),
            if (merged)
              TextButton(
                style: TextButton.styleFrom(padding: EdgeInsets.zero, visualDensity: VisualDensity.compact),
                onPressed: () => context.push('/ideas/${idea.mergedIntoId}'),
                child: Text('Open idea #${idea.mergedIntoId}'),
              ),
          ]),
        ),
      ]),
    );
  }
}

/// The workflow ladder plus the moves available from here. The ladder is
/// read-only — status changes only through the actions, which the server
/// validates, so draft → review → approved means the same thing everywhere.
class _StatusCard extends StatelessWidget {
  final IdeaDetail detail;
  final bool busy;
  final void Function(IdeaTransition) onMove;
  const _StatusCard({required this.detail, required this.busy, required this.onMove});

  @override
  Widget build(BuildContext context) {
    final steps = detail.statusFlow.where((s) => s != IdeaStatus.archived).toList();
    final current = detail.idea.status;
    final currentIndex = steps.indexOf(current);

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: _card(),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const _SectionLabel('STATUS'),
          const SizedBox(height: 10),
          for (var i = 0; i < steps.length; i++)
            _StepRow(
              status: steps[i],
              active: steps[i] == current,
              done: currentIndex > i,
              last: i == steps.length - 1,
            ),
          if (detail.nextStep.isNotEmpty) ...[
            const SizedBox(height: 8),
            Text(detail.nextStep, style: const TextStyle(color: LumoraColors.slatey, fontSize: 13, height: 1.35)),
          ],
          for (final t in detail.transitions) ...[
            const SizedBox(height: 10),
            LumoraButton(
              label: t.label,
              full: true,
              variant: t.primary ? LumoraButtonVariant.primary : LumoraButtonVariant.outline,
              onPressed: t.allowed && !busy ? () => onMove(t) : null,
            ),
            if (!t.allowed && t.reason.isNotEmpty)
              Padding(
                padding: const EdgeInsets.only(top: 4, left: 4, right: 4),
                child: Text(t.reason, style: const TextStyle(color: LumoraColors.gray500, fontSize: 12)),
              ),
          ],
        ],
      ),
    );
  }
}

class _StepRow extends StatelessWidget {
  final IdeaStatus status;
  final bool active;
  final bool done;
  final bool last;
  const _StepRow({required this.status, required this.active, required this.done, required this.last});

  @override
  Widget build(BuildContext context) {
    final c = ideaStatusColor(status);
    final reached = active || done;
    return IntrinsicHeight(
      child: Row(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
        SizedBox(
          width: 22,
          child: Column(children: [
            Container(
              width: 18, height: 18,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                color: done ? c : (active ? c.withValues(alpha: 0.15) : Colors.transparent),
                border: Border.all(color: reached ? c : LumoraColors.gray300, width: 2),
              ),
              child: done ? const Icon(Icons.check, size: 11, color: Colors.white) : null,
            ),
            if (!last)
              Expanded(child: Container(width: 2, color: done ? c.withValues(alpha: 0.5) : LumoraColors.gray100)),
          ]),
        ),
        const SizedBox(width: 10),
        Expanded(
          child: Padding(
            padding: const EdgeInsets.only(bottom: 12),
            child: Text(
              ideaStatusLabel(status),
              style: TextStyle(
                fontWeight: active ? FontWeight.w800 : FontWeight.w600,
                color: active ? c : (done ? LumoraColors.slatey : LumoraColors.gray500),
              ),
            ),
          ),
        ),
      ]),
    );
  }
}

class _TasksCard extends StatelessWidget {
  final IdeaDetail detail;
  final bool busy;
  final VoidCallback onAdd;
  final void Function(IdeaTask) onToggle;
  const _TasksCard({required this.detail, required this.busy, required this.onAdd, required this.onToggle});

  @override
  Widget build(BuildContext context) {
    final idea = detail.idea;
    final closed = idea.archived || idea.mergedIntoId != null;
    // Only approved work becomes tasks — the server enforces the same rule.
    final canAdd = !closed &&
        (idea.status == IdeaStatus.approved || idea.status == IdeaStatus.inProgress || idea.status == IdeaStatus.completed);

    return Container(
      padding: const EdgeInsets.fromLTRB(16, 12, 8, 12),
      decoration: _card(),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            const _SectionLabel('LINKED TASKS'),
            const Spacer(),
            if (canAdd)
              TextButton.icon(
                onPressed: busy ? null : onAdd,
                icon: const Icon(Icons.add, size: 16),
                label: Text(idea.status == IdeaStatus.approved ? 'Convert to task' : 'Add'),
              ),
          ]),
          if (detail.tasks.isEmpty)
            Padding(
              padding: const EdgeInsets.only(top: 4, bottom: 4, right: 8),
              child: Text(
                canAdd
                    ? 'No tasks yet. Adding one starts the work.'
                    : 'Tasks open up once the idea is approved.',
                style: const TextStyle(color: LumoraColors.slatey, fontSize: 13),
              ),
            ),
          for (final t in detail.tasks)
            CheckboxListTile(
              value: t.status == 'done',
              onChanged: busy || closed ? null : (_) => onToggle(t),
              dense: true,
              contentPadding: EdgeInsets.zero,
              title: Text(t.title,
                  style: TextStyle(
                    decoration: t.status == 'done' ? TextDecoration.lineThrough : null,
                    color: t.status == 'done' ? LumoraColors.gray500 : LumoraColors.ink,
                  )),
              subtitle: t.sprint.isNotEmpty ? Text(t.sprint) : null,
              controlAffinity: ListTileControlAffinity.leading,
            ),
          if (detail.tasks.isNotEmpty && idea.status == IdeaStatus.inProgress)
            Padding(
              padding: const EdgeInsets.only(top: 4, right: 8),
              child: Text(
                detail.openTasks == 0
                    ? 'All tasks done.'
                    : '${detail.openTasks} open — the idea completes when every task is done.',
                style: const TextStyle(color: LumoraColors.gray500, fontSize: 12),
              ),
            ),
        ],
      ),
    );
  }
}

class _HistoryRow extends StatelessWidget {
  final IdeaEvent event;
  const _HistoryRow({required this.event});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 5),
      child: Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Container(
          margin: const EdgeInsets.only(top: 6, right: 10),
          width: 7, height: 7,
          decoration: const BoxDecoration(color: LumoraColors.purple, shape: BoxShape.circle),
        ),
        Expanded(
          child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Text.rich(TextSpan(children: [
              TextSpan(text: event.actor.name.isEmpty ? 'Someone' : event.actor.name,
                  style: const TextStyle(fontWeight: FontWeight.w800)),
              TextSpan(text: ' ${_describeEvent(event)}'),
            ]), style: const TextStyle(fontSize: 13, color: LumoraColors.ink)),
            if (event.note.isNotEmpty)
              Text(event.note, style: const TextStyle(fontSize: 12, fontStyle: FontStyle.italic, color: LumoraColors.slatey)),
            Text(_ago(event.at), style: const TextStyle(fontSize: 11, color: LumoraColors.gray500)),
          ]),
        ),
      ]),
    );
  }
}

String _label(String status) => ideaStatusLabel(ideaStatusFromString(status));

String _describeEvent(IdeaEvent e) => switch (e.kind) {
      'created' => 'posted this idea',
      'status' => 'moved it from ${_label(e.from)} to ${_label(e.to)}',
      'vote_threshold' => '— community votes moved it to ${_label(e.to)}',
      'edited' => 'edited the ${e.field}',
      'tagged' => 'changed the tags to ${e.to.isEmpty ? 'none' : e.to}',
      'merged' => 'merged "${e.from}" into "${e.to}"',
      'archived' => 'archived it',
      'restored' => 'restored it to ${_label(e.to)}',
      'task' => 'converted it to a task: "${e.to}"',
      'brainstorm' => 'started a ${e.to} silent brainstorm',
      'brainstorm_ended' => 'ended the silent brainstorm early',
      _ => e.kind,
    };

String _ago(String iso) {
  final t = DateTime.tryParse(iso);
  if (t == null) return '';
  final s = DateTime.now().difference(t).inSeconds;
  if (s < 60) return 'just now';
  if (s < 3600) return '${s ~/ 60}m ago';
  if (s < 86400) return '${s ~/ 3600}h ago';
  return '${s ~/ 86400}d ago';
}

class _ErrorState extends StatelessWidget {
  final String message;
  final VoidCallback onRetry;
  const _ErrorState({required this.message, required this.onRetry});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(mainAxisSize: MainAxisSize.min, children: [
          const Icon(Icons.lightbulb_outline_rounded, size: 48, color: LumoraColors.gray300),
          const SizedBox(height: 12),
          Text(message, textAlign: TextAlign.center, style: const TextStyle(color: LumoraColors.slatey)),
          const SizedBox(height: 16),
          LumoraButton(label: 'Try again', onPressed: onRetry),
        ]),
      ),
    );
  }
}

/// The owner's edit form: title, description and tags.
class _EditIdeaSheet extends StatefulWidget {
  final Idea idea;
  const _EditIdeaSheet({required this.idea});
  @override
  State<_EditIdeaSheet> createState() => _EditIdeaSheetState();
}

class _EditIdeaSheetState extends State<_EditIdeaSheet> {
  late final _title = TextEditingController(text: widget.idea.title);
  late final _desc = TextEditingController(text: widget.idea.description);
  late final _tags = TextEditingController(text: widget.idea.tags.join(', '));
  bool _busy = false;
  String? _error;

  @override
  void dispose() {
    _title.dispose();
    _desc.dispose();
    _tags.dispose();
    super.dispose();
  }

  Future<void> _save() async {
    if (_title.text.trim().isEmpty) {
      setState(() => _error = 'An idea needs a title');
      return;
    }
    setState(() { _busy = true; _error = null; });
    try {
      await ApiClient.instance.updateIdea(
        widget.idea.id,
        title: _title.text.trim(),
        description: _desc.text.trim(),
        tags: _tags.text.split(',').map((t) => t.trim()).where((t) => t.isNotEmpty).toList(),
      );
      if (mounted) Navigator.pop(context, true);
    } on ApiException catch (e) {
      setState(() => _error = e.message);
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Padding(
        padding: EdgeInsets.only(left: 20, right: 20, top: 20, bottom: 20 + MediaQuery.of(context).viewInsets.bottom),
        child: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text('Edit idea', style: TextStyle(fontSize: 18, fontWeight: FontWeight.w800)),
              const SizedBox(height: 16),
              TextField(controller: _title, decoration: const InputDecoration(labelText: 'Title')),
              const SizedBox(height: 12),
              TextField(controller: _desc, maxLines: 4, decoration: const InputDecoration(labelText: 'Description')),
              const SizedBox(height: 12),
              TextField(controller: _tags, decoration: const InputDecoration(labelText: 'Tags (comma separated)')),
              if (_error != null)
                Padding(padding: const EdgeInsets.only(top: 8), child: Text(_error!, style: const TextStyle(color: LumoraColors.coral))),
              const SizedBox(height: 16),
              LumoraButton(label: 'Save changes', full: true, loading: _busy, onPressed: _save),
            ],
          ),
        ),
      ),
    );
  }
}
