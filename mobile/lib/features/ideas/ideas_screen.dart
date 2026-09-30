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
const _kSorts = [('hot', 'Hot'), ('top', 'Top'), ('new', 'New'), ('controversial', 'Controversial')];

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

  @override
  void initState() {
    super.initState();
    _load();
    if (widget.ideaId != null) {
      WidgetsBinding.instance.addPostFrameCallback((_) => context.push('/ideas/${widget.ideaId}'));
    }
  }

  Future<void> _load() async {
    setState(() => _loading = true);
    try {
      final board = await ApiClient.instance.ideas(status: _status, sort: _sort, q: _q);
      if (mounted) setState(() => _board = board);
    } catch (_) {
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
        elevation: 0,
        leading: BackButton(onPressed: () => context.canPop() ? context.pop() : context.go('/home')),
        title: const Text('Ideas', style: TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.ink)),
        actions: [IconButton(icon: const Icon(Icons.add), onPressed: _openCreate)],
      ),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 12, 16, 8),
            child: TextField(
              onSubmitted: (v) { _q = v; _load(); },
              decoration: const InputDecoration(hintText: 'Search ideas…', prefixIcon: Icon(Icons.search)),
            ),
          ),
          SizedBox(
            height: 40,
            child: ListView(
              scrollDirection: Axis.horizontal,
              padding: const EdgeInsets.symmetric(horizontal: 16),
              children: [
                for (final s in _kStatusFilters)
                  Padding(
                    padding: const EdgeInsets.only(right: 8),
                    child: ChoiceChip(
                      label: Text(s.isEmpty ? 'All' : ideaStatusLabel(ideaStatusFromString(s))),
                      selected: _status == s,
                      onSelected: (_) { setState(() => _status = s); _load(); },
                    ),
                  ),
                Container(width: 1, height: 24, color: LumoraColors.gray100, margin: const EdgeInsets.symmetric(horizontal: 4)),
                for (final s in _kSorts)
                  Padding(
                    padding: const EdgeInsets.only(right: 8),
                    child: ChoiceChip(
                      label: Text(s.$2),
                      selected: _sort == s.$1,
                      onSelected: (_) { setState(() => _sort = s.$1); _load(); },
                    ),
                  ),
              ],
            ),
          ),
          const SizedBox(height: 8),
          Expanded(
            child: _loading
                ? const Center(child: CircularProgressIndicator())
                : board == null || board.ideas.isEmpty
                    ? const Center(child: Text('No ideas yet. Be the first to propose one!', style: TextStyle(color: LumoraColors.slatey)))
                    : RefreshIndicator(
                        onRefresh: _load,
                        child: ListView.builder(
                          padding: const EdgeInsets.fromLTRB(16, 0, 16, 24),
                          itemCount: board.ideas.length,
                          itemBuilder: (context, i) {
                            final idea = board.ideas[i];
                            return Padding(
                              padding: const EdgeInsets.only(bottom: 10),
                              child: _IdeaCard(idea: idea, onVote: (v) => _vote(idea, v), onStar: () => _star(idea)),
                            );
                          },
                        ),
                      ),
          ),
        ],
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
    return Material(
      color: Colors.white,
      borderRadius: BorderRadius.circular(LumoraRadii.xl),
      child: InkWell(
        borderRadius: BorderRadius.circular(LumoraRadii.xl),
        onTap: () => context.push('/ideas/${idea.id}'),
        child: Container(
          padding: const EdgeInsets.all(14),
          decoration: BoxDecoration(borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Column(children: [
                IconButton(
                  icon: Icon(Icons.arrow_upward, size: 18, color: idea.myVote == 1 ? LumoraColors.teal : LumoraColors.gray300),
                  onPressed: () => onVote(1),
                  padding: EdgeInsets.zero,
                  constraints: const BoxConstraints(minWidth: 28, minHeight: 28),
                ),
                Text('${idea.score}', style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 13)),
                IconButton(
                  icon: Icon(Icons.arrow_downward, size: 18, color: idea.myVote == -1 ? LumoraColors.coral : LumoraColors.gray300),
                  onPressed: () => onVote(-1),
                  padding: EdgeInsets.zero,
                  constraints: const BoxConstraints(minWidth: 28, minHeight: 28),
                ),
              ]),
              const SizedBox(width: 10),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(children: [
                      Expanded(child: Text(idea.title, maxLines: 1, overflow: TextOverflow.ellipsis, style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 15))),
                      InkWell(onTap: onStar, child: Icon(idea.starred ? Icons.star : Icons.star_border, size: 18, color: idea.starred ? LumoraColors.amber : LumoraColors.gray300)),
                    ]),
                    if (idea.description.isNotEmpty)
                      Text(idea.description, maxLines: 2, overflow: TextOverflow.ellipsis, style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
                    const SizedBox(height: 6),
                    Row(children: [
                      _statusBadge(idea.status),
                      const SizedBox(width: 8),
                      LumoraAvatar(name: idea.owner.name, avatarColor: idea.owner.avatarColor, avatarUrl: idea.owner.avatarUrl, size: 18),
                      const SizedBox(width: 4),
                      Expanded(child: Text(idea.owner.name, maxLines: 1, overflow: TextOverflow.ellipsis, style: const TextStyle(color: LumoraColors.slatey, fontSize: 11))),
                      const Icon(Icons.chat_bubble_outline, size: 12, color: LumoraColors.gray500),
                      const SizedBox(width: 2),
                      Text('${idea.messageCount}', style: const TextStyle(color: LumoraColors.gray500, fontSize: 11)),
                    ]),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _statusBadge(IdeaStatus status) {
    final colors = {
      IdeaStatus.draft: LumoraColors.gray500, IdeaStatus.underReview: LumoraColors.amber,
      IdeaStatus.approved: LumoraColors.teal, IdeaStatus.inProgress: LumoraColors.purple,
      IdeaStatus.completed: LumoraColors.teal, IdeaStatus.archived: LumoraColors.gray500,
    };
    final c = colors[status] ?? LumoraColors.gray500;
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 1),
      decoration: BoxDecoration(color: c.withValues(alpha: 0.12), borderRadius: BorderRadius.circular(LumoraRadii.full)),
      child: Text(ideaStatusLabel(status), style: TextStyle(color: c, fontSize: 10, fontWeight: FontWeight.w800)),
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
