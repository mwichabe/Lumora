import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/network/api_client.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../core/voices.dart';
import '../../models/listening_reading.dart';
import '../../providers/auth_provider.dart';
import '../../widgets/fox_mascot.dart';
import '../../widgets/lumora_button.dart';
import '../../widgets/quiz_phase.dart';

enum _Phase { match, listen, quiz, done }

class ListeningScreen extends ConsumerStatefulWidget {
  final int id;
  const ListeningScreen({super.key, required this.id});

  @override
  ConsumerState<ListeningScreen> createState() => _ListeningScreenState();
}

class _ListeningScreenState extends ConsumerState<ListeningScreen> {
  ListeningSession? _session;
  _Phase _phase = _Phase.match;

  @override
  void initState() {
    super.initState();
    ApiClient.instance.listeningSession(widget.id).then((s) {
      if (!mounted) return;
      setState(() {
        _session = s;
        _phase = s.matches.isNotEmpty ? _Phase.match : _Phase.listen;
      });
    });
  }

  @override
  void dispose() {
    Voices.instance.stopSpeaking();
    super.dispose();
  }

  Future<void> _complete() async {
    try {
      final (_, user) = await ApiClient.instance.completeListening(widget.id);
      ref.read(authProvider.notifier).setUser(user);
    } catch (_) {}
    setState(() => _phase = _Phase.done);
  }

  @override
  Widget build(BuildContext context) {
    final session = _session;
    if (session == null) {
      return const Scaffold(body: Center(child: FoxMascot(size: 120, glow: true)));
    }

    return Scaffold(
      backgroundColor: LumoraColors.cream,
      body: SafeArea(
        child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 8, 16, 8),
              child: Row(children: [
                IconButton(
                  icon: const Icon(Icons.close, color: LumoraColors.gray500),
                  onPressed: () {
                    Voices.instance.stopSpeaking();
                    context.go('/learn');
                  },
                ),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text('LISTENING · ${session.unit}',
                          style: const TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.gray500)),
                      Text(session.title, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w800)),
                    ],
                  ),
                ),
              ]),
            ),
            Expanded(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20),
                child: switch (_phase) {
                  _Phase.match => _MatchPhase(pairs: session.matches, onDone: () => setState(() => _phase = _Phase.listen)),
                  _Phase.listen => _ListenPhase(session: session, onDone: () => setState(() => _phase = _Phase.quiz)),
                  _Phase.quiz => ComprehensionQuiz(
                      questions: [
                        for (final q in session.questions)
                          (question: q.question, options: q.options ?? [], correctAnswer: q.correctAnswer),
                      ],
                      onDone: _complete,
                    ),
                  _Phase.done => SessionDone(
                      title: 'Session complete!',
                      subtitle: 'Your ear is getting sharper. ¡Bien hecho!',
                      xp: session.xpReward,
                      onContinue: () => context.go('/learn'),
                    ),
                },
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _MatchPhase extends StatefulWidget {
  final List<ListeningMatch> pairs;
  final VoidCallback onDone;
  const _MatchPhase({required this.pairs, required this.onDone});

  @override
  State<_MatchPhase> createState() => _MatchPhaseState();
}

class _MatchPhaseState extends State<_MatchPhase> {
  late final List<ListeningMatch> _shuffled;
  int? _selected;
  final Set<int> _matched = {};
  int? _wrong;

  @override
  void initState() {
    super.initState();
    _shuffled = [...widget.pairs]..sort((a, b) => ((a.id * 7) % 5).compareTo((b.id * 7) % 5));
  }

  void _tapRight(int rightId, String word) {
    if (_selected == null || _matched.contains(rightId)) return;
    final left = widget.pairs.where((p) => p.id == _selected).firstOrNull;
    if (left != null && left.word == word) {
      setState(() {
        _matched.add(rightId);
        _selected = null;
      });
      Voices.instance.speakAs(null, word);
    } else {
      setState(() => _wrong = rightId);
      Future.delayed(const Duration(milliseconds: 500), () {
        if (mounted) setState(() => _wrong = null);
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final allDone = _matched.length == widget.pairs.length && widget.pairs.isNotEmpty;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text('WARM-UP · MATCH THE WORDS', style: TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.gray500)),
        const SizedBox(height: 4),
        const Text("These words appear in the conversation you're about to hear.", style: TextStyle(color: LumoraColors.slatey)),
        const SizedBox(height: 16),
        Expanded(
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Expanded(
                child: Column(children: [
                  for (final p in widget.pairs) _MatchTile(
                    label: p.translation,
                    done: _matched.contains(p.id),
                    active: _selected == p.id,
                    onTap: _matched.contains(p.id) ? null : () => setState(() => _selected = p.id),
                  ),
                ]),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(children: [
                  for (final p in _shuffled) _MatchTile(
                    label: p.word,
                    done: _matched.contains(p.id),
                    wrong: _wrong == p.id,
                    onTap: _matched.contains(p.id) ? null : () => _tapRight(p.id, p.word),
                  ),
                ]),
              ),
            ],
          ),
        ),
        Padding(
          padding: const EdgeInsets.symmetric(vertical: 12),
          child: LumoraButton(label: allDone ? 'Listen to the conversation' : 'Match all the words', full: true, onPressed: allDone ? widget.onDone : null),
        ),
      ],
    );
  }
}

extension _FirstOrNullExt<T> on Iterable<T> {
  T? get firstOrNull => isEmpty ? null : first;
}

class _MatchTile extends StatelessWidget {
  final String label;
  final bool done;
  final bool active;
  final bool wrong;
  final VoidCallback? onTap;
  const _MatchTile({required this.label, this.done = false, this.active = false, this.wrong = false, this.onTap});

  @override
  Widget build(BuildContext context) {
    Color border = LumoraColors.gray100, bg = Colors.white, fg = LumoraColors.ink;
    if (done) { border = LumoraColors.teal; bg = LumoraColors.tealLight; fg = LumoraColors.teal; }
    else if (wrong) { border = LumoraColors.coral; bg = LumoraColors.coralLight; }
    else if (active) { border = LumoraColors.purple; bg = LumoraColors.purpleLight; }

    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Material(
        color: bg,
        borderRadius: BorderRadius.circular(LumoraRadii.lg),
        child: InkWell(
          borderRadius: BorderRadius.circular(LumoraRadii.lg),
          onTap: onTap,
          child: Container(
            height: 52,
            padding: const EdgeInsets.symmetric(horizontal: 12),
            decoration: BoxDecoration(border: Border.all(color: border, width: 2), borderRadius: BorderRadius.circular(LumoraRadii.lg)),
            child: Row(children: [
              Expanded(child: Text(label, maxLines: 1, overflow: TextOverflow.ellipsis, style: TextStyle(fontWeight: FontWeight.w700, color: fg))),
              if (done) const Icon(Icons.check, size: 16, color: LumoraColors.teal),
            ]),
          ),
        ),
      ),
    );
  }
}

class _ListenPhase extends StatefulWidget {
  final ListeningSession session;
  final VoidCallback onDone;
  const _ListenPhase({required this.session, required this.onDone});

  @override
  State<_ListenPhase> createState() => _ListenPhaseState();
}

class _ListenPhaseState extends State<_ListenPhase> {
  int _active = -1;
  bool _playing = false;
  bool _revealed = false;

  Future<void> _playAll() async {
    if (_playing) {
      setState(() => _playing = false);
      Voices.instance.stopSpeaking();
      return;
    }
    setState(() => _playing = true);
    await Voices.instance.speakSequence(
      [for (final l in widget.session.lines) (character: l.character, text: l.text)],
      onLine: (i) { if (mounted) setState(() => _active = i); },
      shouldContinue: () => _playing,
    );
    if (mounted) setState(() { _playing = false; _active = -1; });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(widget.session.description, style: const TextStyle(color: LumoraColors.slatey)),
        const SizedBox(height: 12),
        LumoraButton(
          label: _playing ? 'Pause' : 'Play conversation',
          full: true,
          icon: Icon(_playing ? Icons.pause : Icons.play_arrow, color: Colors.white, size: 20),
          onPressed: _playAll,
        ),
        const SizedBox(height: 16),
        Expanded(
          child: ListView.separated(
            itemCount: widget.session.lines.length,
            separatorBuilder: (_, _) => const SizedBox(height: 10),
            itemBuilder: (context, i) {
              final l = widget.session.lines[i];
              final isActive = i == _active;
              return Opacity(
                opacity: isActive || _active == -1 ? 1 : 0.55,
                child: Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: isActive ? LumoraColors.purpleLight : Colors.white,
                    border: Border.all(color: isActive ? LumoraColors.purple : LumoraColors.gray100),
                    borderRadius: BorderRadius.circular(LumoraRadii.xl),
                  ),
                  child: Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
                    CircleAvatar(radius: 20, backgroundColor: LumoraColors.purple, child: Text(l.character.isNotEmpty ? l.character[0] : '?', style: const TextStyle(color: Colors.white))),
                    const SizedBox(width: 10),
                    Expanded(
                      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                        Text(l.character, style: const TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.slatey)),
                        Text(l.text, style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 15)),
                        if (_revealed) Text(l.translation, style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
                      ]),
                    ),
                    IconButton(icon: const Icon(Icons.volume_up_rounded, size: 18, color: LumoraColors.purple),
                        onPressed: () => Voices.instance.speakAs(l.character, l.text)),
                  ]),
                ),
              );
            },
          ),
        ),
        Center(
          child: TextButton(
            onPressed: () => setState(() => _revealed = !_revealed),
            child: Text(_revealed ? 'Hide translations' : 'Show translations', style: const TextStyle(color: LumoraColors.teal, fontWeight: FontWeight.w700)),
          ),
        ),
        Padding(
          padding: const EdgeInsets.only(top: 8, bottom: 12),
          child: LumoraButton(label: "I'm ready — quiz me", full: true, onPressed: widget.onDone),
        ),
      ],
    );
  }
}
