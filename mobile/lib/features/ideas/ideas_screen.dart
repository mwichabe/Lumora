import 'dart:async';

import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../core/network/api_client.dart';
import '../../core/network/api_exception.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../core/theme/shadows.dart';
import '../../models/idea.dart';
import '../../widgets/avatar.dart';
import '../../widgets/lumora_button.dart';

const _kStatusFilters = ['', 'draft', 'under_review', 'approved', 'in_progress', 'completed'];
const _kSorts = [
  ('hot', 'Hot', Icons.local_fire_department_rounded),
  ('top', 'Top', Icons.trending_up_rounded),
  ('new', 'New', Icons.schedule_rounded),
  ('controversial', 'Debated', Icons.forum_rounded),
];

/// The board's count key for each status filter.
const _kCountKeys = {'': 'all', 'draft': 'draft', 'under_review': 'underReview', 'approved': 'approved', 'in_progress': 'inProgress', 'completed': 'completed'};

Color _statusColor(IdeaStatus status) => switch (status) {
      IdeaStatus.draft => LumoraColors.gray500,
      IdeaStatus.underReview => const Color(0xFFD48806),
      IdeaStatus.approved => LumoraColors.teal,
      IdeaStatus.inProgress => LumoraColors.purple,
      IdeaStatus.completed => const Color(0xFF0B9E6E),
      IdeaStatus.archived => LumoraColors.gray500,
    };

IconData _statusIcon(IdeaStatus status) => switch (status) {
      IdeaStatus.draft => Icons.edit_note_rounded,
      IdeaStatus.underReview => Icons.hourglass_top_rounded,
      IdeaStatus.approved => Icons.thumb_up_alt_rounded,
      IdeaStatus.inProgress => Icons.construction_rounded,
      IdeaStatus.completed => Icons.check_circle_rounded,
      IdeaStatus.archived => Icons.inventory_2_rounded,
    };

class IdeasScreen extends StatefulWidget {
  final int? ideaId;
  const IdeasScreen({super.key, this.ideaId});

  @override
  State<IdeasScreen> createState() => _IdeasScreenState();
}

class _IdeasScreenState extends State<IdeasScreen> {
  IdeaBoard? _board;
  bool _loading = true;
  String _status = '';
  String _sort = 'hot';
  String _q = '';
  String? _error;
  final _search = TextEditingController();

  @override
  void dispose() {
    _search.dispose();
    super.dispose();
  }

  @override
  void initState() {
    super.initState();
    _load();
    if (widget.ideaId != null) {
      WidgetsBinding.instance.addPostFrameCallback((_) => context.push('/ideas/${widget.ideaId}'));
    }
  }

  Future<void> _load() async {
    setState(() { _loading = true; _error = null; });
    try {
      final board = await ApiClient.instance.ideas(status: _status, sort: _sort, q: _q);
      if (mounted) setState(() => _board = board);
    } on ApiException catch (e) {
      // Say so, rather than showing an empty board as if there were no ideas.
      if (mounted) setState(() => _error = e.message);
    } catch (_) {
      if (mounted) setState(() => _error = 'Something went wrong. Please try again.');
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _vote(Idea idea, int value) async {
    final v = idea.myVote == value ? 0 : value;
    try {
      await ApiClient.instance.voteIdea(idea.id, v);
      _load();
    } catch (_) {}
  }

  Future<void> _star(Idea idea) async {
    try {
      await ApiClient.instance.starIdea(idea.id);
      _load();
    } catch (_) {}
  }

  void _openCreate() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.white,
      shape: const RoundedRectangleBorder(borderRadius: BorderRadius.vertical(top: Radius.circular(24))),
      builder: (_) => _NewIdeaSheet(onCreated: _load),
    );
  }

  @override
  Widget build(BuildContext context) {
    final board = _board;
    return Scaffold(
      backgroundColor: LumoraColors.cream,
      appBar: AppBar(
        backgroundColor: Colors.white,
        foregroundColor: LumoraColors.ink,
        surfaceTintColor: Colors.white,
        elevation: 0,
        leading: BackButton(onPressed: () => context.canPop() ? context.pop() : context.go('/home')),
        title: const Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Text('Ideas', style: TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.ink, fontSize: 20)),
          Text('Shape what Lumora builds next', style: TextStyle(color: LumoraColors.slatey, fontSize: 12)),
        ]),
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: _openCreate,
        backgroundColor: LumoraColors.purple,
        foregroundColor: Colors.white,
        shape: const StadiumBorder(),
        icon: const Icon(Icons.lightbulb_rounded),
        label: const Text('New idea', style: TextStyle(fontWeight: FontWeight.w800)),
      ),
      body: Column(
        children: [
          Container(
            color: Colors.white,
            padding: const EdgeInsets.fromLTRB(16, 4, 16, 12),
            child: Column(children: [
              TextField(
                controller: _search,
                textInputAction: TextInputAction.search,
                onSubmitted: (v) { _q = v.trim(); _load(); },
                decoration: InputDecoration(
                  hintText: 'Search ideas',
                  prefixIcon: const Icon(Icons.search_rounded, color: LumoraColors.gray500),
                  suffixIcon: _q.isEmpty
                      ? null
                      : IconButton(
                          tooltip: 'Clear search',
                          icon: const Icon(Icons.close_rounded, color: LumoraColors.gray500),
                          onPressed: () { _search.clear(); _q = ''; _load(); },
                        ),
                ),
              ),
              const SizedBox(height: 12),
              SizedBox(
                height: 36,
                child: ListView(
                  scrollDirection: Axis.horizontal,
                  children: [
                    for (final s in _kStatusFilters)
                      Padding(
                        padding: const EdgeInsets.only(right: 8),
                        child: _FilterChip(
                          label: s.isEmpty ? 'All' : ideaStatusLabel(ideaStatusFromString(s)),
                          count: board?.counts[_kCountKeys[s]],
                          selected: _status == s,
                          onTap: () { setState(() => _status = s); _load(); },
                        ),
                      ),
                  ],
                ),
              ),
            ]),
          ),
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 12, 8, 4),
            child: Row(children: [
              Expanded(
                child: Text(
                  board == null ? '' : '${board.ideas.length} ${board.ideas.length == 1 ? "idea" : "ideas"}',
                  style: const TextStyle(color: LumoraColors.slatey, fontWeight: FontWeight.w700, fontSize: 13),
                ),
              ),
              PopupMenuButton<String>(
                tooltip: 'Sort ideas',
                color: Colors.white,
                initialValue: _sort,
                onSelected: (v) { setState(() => _sort = v); _load(); },
                itemBuilder: (_) => [
                  for (final s in _kSorts)
                    PopupMenuItem(
                      value: s.$1,
                      child: Row(children: [
                        Icon(s.$3, size: 18, color: _sort == s.$1 ? LumoraColors.purple : LumoraColors.slatey),
                        const SizedBox(width: 10),
                        Text(s.$2, style: TextStyle(fontWeight: _sort == s.$1 ? FontWeight.w800 : FontWeight.w600)),
                      ]),
                    ),
                ],
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 6),
                  child: Row(mainAxisSize: MainAxisSize.min, children: [
                    Icon(_kSorts.firstWhere((s) => s.$1 == _sort).$3, size: 16, color: LumoraColors.purple),
                    const SizedBox(width: 4),
                    Text(_kSorts.firstWhere((s) => s.$1 == _sort).$2,
                        style: const TextStyle(color: LumoraColors.purple, fontWeight: FontWeight.w800, fontSize: 13)),
                    const Icon(Icons.expand_more_rounded, size: 18, color: LumoraColors.purple),
                  ]),
                ),
              ),
            ]),
          ),
          Expanded(child: _body(board)),
        ],
      ),
    );
  }

  Widget _body(IdeaBoard? board) {
    if (_loading && board == null) return const Center(child: CircularProgressIndicator());
    if (_error != null && board == null) {
      return _EmptyState(
        icon: Icons.wifi_off_rounded,
        title: "Couldn't load ideas",
        message: _error!,
        action: LumoraButton(label: 'Try again', onPressed: _load),
      );
    }
    if (board == null || board.ideas.isEmpty) {
      return _EmptyState(
        icon: Icons.lightbulb_outline_rounded,
        title: _q.isNotEmpty || _status.isNotEmpty ? 'No ideas match' : 'No ideas yet',
        message: _q.isNotEmpty || _status.isNotEmpty
            ? 'Try another filter or search term.'
            : 'Have a feature in mind? Be the first to propose it.',
        action: LumoraButton(label: 'Propose an idea', onPressed: _openCreate),
      );
    }
    return RefreshIndicator(
      onRefresh: _load,
      child: ListView.builder(
        // Bottom padding keeps the last card clear of the "New idea" button.
        padding: const EdgeInsets.fromLTRB(16, 4, 16, 96),
        itemCount: board.ideas.length,
        itemBuilder: (context, i) {
          final idea = board.ideas[i];
          return Padding(
            padding: const EdgeInsets.only(bottom: 12),
            child: _IdeaCard(idea: idea, onVote: (v) => _vote(idea, v), onStar: () => _star(idea)),
          );
        },
      ),
    );
  }
}

class _FilterChip extends StatelessWidget {
  final String label;
  final int? count;
  final bool selected;
  final VoidCallback onTap;
  const _FilterChip({required this.label, required this.count, required this.selected, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return Material(
      color: selected ? LumoraColors.purple : LumoraColors.gray50,
      shape: const StadiumBorder(),
      child: InkWell(
        customBorder: const StadiumBorder(),
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 14),
          child: Row(mainAxisSize: MainAxisSize.min, children: [
            Text(label, style: TextStyle(color: selected ? Colors.white : LumoraColors.ink, fontWeight: FontWeight.w700, fontSize: 13)),
            if (count != null && count! > 0) ...[
              const SizedBox(width: 6),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 1),
                decoration: BoxDecoration(
                  color: selected ? Colors.white.withValues(alpha: 0.25) : Colors.white,
                  borderRadius: BorderRadius.circular(LumoraRadii.full),
                ),
                child: Text('$count', style: TextStyle(color: selected ? Colors.white : LumoraColors.slatey, fontWeight: FontWeight.w800, fontSize: 11)),
              ),
            ],
          ]),
        ),
      ),
    );
  }
}

class _IdeaCard extends StatelessWidget {
  final Idea idea;
  final void Function(int) onVote;
  final VoidCallback onStar;
  const _IdeaCard({required this.idea, required this.onVote, required this.onStar});

  @override
  Widget build(BuildContext context) {
    final statusColor = _statusColor(idea.status);
    return Material(
      color: Colors.white,
      borderRadius: BorderRadius.circular(LumoraRadii.xl),
      child: InkWell(
        borderRadius: BorderRadius.circular(LumoraRadii.xl),
        onTap: () => context.push('/ideas/${idea.id}'),
        child: Container(
          decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
          clipBehavior: Clip.antiAlias,
          child: IntrinsicHeight(
            child: Row(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
              // Status stripe down the left edge.
              Container(width: 5, color: statusColor),
              Expanded(
                child: Padding(
                  padding: const EdgeInsets.fromLTRB(14, 14, 14, 12),
                  child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                    Row(children: [
                      _StatusBadge(status: idea.status),
                      const Spacer(),
                      InkResponse(
                        onTap: onStar,
                        radius: 20,
                        child: Icon(idea.starred ? Icons.star_rounded : Icons.star_outline_rounded,
                            size: 22, color: idea.starred ? LumoraColors.amber : LumoraColors.gray300),
                      ),
                    ]),
                    const SizedBox(height: 8),
                    Text(idea.title, maxLines: 2, overflow: TextOverflow.ellipsis,
                        style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 16, color: LumoraColors.ink, height: 1.25)),
                    if (idea.description.isNotEmpty) ...[
                      const SizedBox(height: 4),
                      Text(idea.description, maxLines: 2, overflow: TextOverflow.ellipsis,
                          style: const TextStyle(color: LumoraColors.slatey, fontSize: 13, height: 1.35)),
                    ],
                    if (idea.tags.isNotEmpty) ...[
                      const SizedBox(height: 8),
                      Wrap(spacing: 6, runSpacing: 6, children: [
                        for (final t in idea.tags.take(3))
                          Container(
                            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                            decoration: BoxDecoration(color: LumoraColors.purpleLight, borderRadius: BorderRadius.circular(LumoraRadii.full)),
                            child: Text('#$t', style: const TextStyle(color: LumoraColors.purple, fontSize: 11, fontWeight: FontWeight.w700)),
                          ),
                      ]),
                    ],
                    const SizedBox(height: 12),
                    Row(children: [
                      _VotePill(score: idea.score, myVote: idea.myVote, onVote: onVote),
                      const SizedBox(width: 10),
                      const Icon(Icons.chat_bubble_outline_rounded, size: 15, color: LumoraColors.gray500),
                      const SizedBox(width: 4),
                      Text('${idea.messageCount}', style: const TextStyle(color: LumoraColors.slatey, fontSize: 12, fontWeight: FontWeight.w700)),
                      const Spacer(),
                      LumoraAvatar(name: idea.owner.name, avatarColor: idea.owner.avatarColor, avatarUrl: idea.owner.avatarUrl, size: 20),
                      const SizedBox(width: 6),
                      ConstrainedBox(
                        constraints: const BoxConstraints(maxWidth: 110),
                        child: Text(idea.owner.name, maxLines: 1, overflow: TextOverflow.ellipsis,
                            style: const TextStyle(color: LumoraColors.slatey, fontSize: 12, fontWeight: FontWeight.w600)),
                      ),
                    ]),
                  ]),
                ),
              ),
            ]),
          ),
        ),
      ),
    );
  }
}

class _StatusBadge extends StatelessWidget {
  final IdeaStatus status;
  const _StatusBadge({required this.status});

  @override
  Widget build(BuildContext context) {
    final c = _statusColor(status);
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
      decoration: BoxDecoration(color: c.withValues(alpha: 0.12), borderRadius: BorderRadius.circular(LumoraRadii.full)),
      child: Row(mainAxisSize: MainAxisSize.min, children: [
        Icon(_statusIcon(status), size: 12, color: c),
        const SizedBox(width: 4),
        Text(ideaStatusLabel(status), style: TextStyle(color: c, fontSize: 11, fontWeight: FontWeight.w800)),
      ]),
    );
  }
}

/// Up/down voting in one pill: the score in the middle, your vote highlighted.
class _VotePill extends StatelessWidget {
  final int score;
  final int myVote;
  final void Function(int) onVote;
  const _VotePill({required this.score, required this.myVote, required this.onVote});

  @override
  Widget build(BuildContext context) {
    final color = myVote == 1 ? LumoraColors.teal : myVote == -1 ? LumoraColors.coral : LumoraColors.slatey;
    return Container(
      decoration: BoxDecoration(
        color: myVote == 0 ? LumoraColors.gray50 : color.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(LumoraRadii.full),
      ),
      child: Row(mainAxisSize: MainAxisSize.min, children: [
        _voteButton(Icons.arrow_upward_rounded, 'Upvote', myVote == 1 ? LumoraColors.teal : LumoraColors.gray500, () => onVote(1)),
        Text('$score', style: TextStyle(fontWeight: FontWeight.w800, fontSize: 13, color: myVote == 0 ? LumoraColors.ink : color)),
        _voteButton(Icons.arrow_downward_rounded, 'Downvote', myVote == -1 ? LumoraColors.coral : LumoraColors.gray500, () => onVote(-1)),
      ]),
    );
  }

  Widget _voteButton(IconData icon, String tooltip, Color color, VoidCallback onTap) => IconButton(
        tooltip: tooltip,
        onPressed: onTap,
        icon: Icon(icon, size: 18, color: color),
        padding: EdgeInsets.zero,
        visualDensity: VisualDensity.compact,
        constraints: const BoxConstraints(minWidth: 34, minHeight: 34),
      );
}

class _EmptyState extends StatelessWidget {
  final IconData icon;
  final String title;
  final String message;
  final Widget action;
  const _EmptyState({required this.icon, required this.title, required this.message, required this.action});

  @override
  Widget build(BuildContext context) {
    return Center(
      child: SingleChildScrollView(
        padding: const EdgeInsets.all(32),
        child: Column(mainAxisSize: MainAxisSize.min, children: [
          Container(
            width: 72,
            height: 72,
            decoration: const BoxDecoration(color: LumoraColors.purpleLight, shape: BoxShape.circle),
            child: Icon(icon, size: 34, color: LumoraColors.purple),
          ),
          const SizedBox(height: 16),
          Text(title, textAlign: TextAlign.center, style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w800, color: LumoraColors.ink)),
          const SizedBox(height: 6),
          Text(message, textAlign: TextAlign.center, style: const TextStyle(color: LumoraColors.slatey, height: 1.4)),
          const SizedBox(height: 20),
          action,
        ]),
      ),
    );
  }
}

class _NewIdeaSheet extends StatefulWidget {
  final VoidCallback onCreated;
  const _NewIdeaSheet({required this.onCreated});
  @override
  State<_NewIdeaSheet> createState() => _NewIdeaSheetState();
}

class _NewIdeaSheetState extends State<_NewIdeaSheet> {
  final _title = TextEditingController();
  final _desc = TextEditingController();
  bool _busy = false;
  String? _error;
  List<SimilarIdea> _similar = [];
  Timer? _debounce;

  @override
  void dispose() {
    _debounce?.cancel();
    super.dispose();
  }

  void _onTitleChanged(String v) {
    _debounce?.cancel();
    if (v.trim().length < 4) {
      setState(() => _similar = []);
      return;
    }
    _debounce = Timer(const Duration(milliseconds: 400), () async {
      try {
        final similar = await ApiClient.instance.similarIdeas(v.trim());
        if (mounted) setState(() => _similar = similar);
      } catch (_) {}
    });
  }

  Future<void> _submit() async {
    if (_title.text.trim().isEmpty) return;
    setState(() { _busy = true; _error = null; });
    try {
      await ApiClient.instance.createIdea(title: _title.text.trim(), description: _desc.text.trim());
      widget.onCreated();
      if (mounted) Navigator.pop(context);
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
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text('New idea', style: TextStyle(fontSize: 18, fontWeight: FontWeight.w800)),
            const SizedBox(height: 16),
            TextField(controller: _title, onChanged: _onTitleChanged, decoration: const InputDecoration(labelText: 'Title')),
            if (_similar.isNotEmpty)
              Container(
                margin: const EdgeInsets.only(top: 8),
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(color: LumoraColors.amberLight, borderRadius: BorderRadius.circular(LumoraRadii.md)),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text('Similar ideas already exist:', style: TextStyle(fontWeight: FontWeight.w800, fontSize: 12)),
                    for (final s in _similar.take(3))
                      Padding(
                        padding: const EdgeInsets.only(top: 4),
                        child: Text('• ${s.title} (${(s.similarity * 100).round()}% match)', style: const TextStyle(fontSize: 12, color: LumoraColors.slatey)),
                      ),
                  ],
                ),
              ),
            const SizedBox(height: 12),
            TextField(controller: _desc, maxLines: 4, decoration: const InputDecoration(labelText: 'Description (optional)')),
            if (_error != null) Padding(padding: const EdgeInsets.only(top: 8), child: Text(_error!, style: const TextStyle(color: LumoraColors.coral))),
            const SizedBox(height: 16),
            LumoraButton(label: 'Post idea', full: true, loading: _busy, onPressed: _submit),
          ],
        ),
      ),
    );
  }
}
