import 'dart:math';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/network/api_client.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../core/theme/shadows.dart';
import '../../core/voices.dart';
import '../../models/lesson.dart';
import '../../models/listening_reading.dart';
import '../../providers/auth_provider.dart';
import '../../widgets/fox_mascot.dart';
import '../../widgets/lumora_button.dart';
import '../../widgets/mistakes_review.dart';

const _kModeTitle = {
  'mix': 'Daily Mix', 'quiz': 'Vocabulary Quiz', 'listen': 'Listening Drill',
  'speak': 'Speaking Practice', 'mistakes': 'Review Mistakes',
  'listening': 'Listening Comprehension', 'reading': 'Reading Comprehension',
};

final _rand = Random();

/// Back to the Practice tab. Popping (rather than navigating) hands control
/// back to the tab, which then reloads its counts — mistakes fixed or made in
/// this drill show up straight away.
void _backToPractice(BuildContext context) {
  if (context.canPop()) {
    context.pop();
  } else {
    context.go('/practice');
  }
}
List<T> _shuffle<T>(List<T> a) => [...a]..shuffle(_rand);
List<T> _sample<T>(List<T> a, int n) => _shuffle(a).take(n).toList();

class _Drill {
  final String kind; // choose | speak
  final String prompt;
  final String question;
  final String correct;
  final List<String> options;
  final String? speaker;
  final int? mistakeId;

  /// Play the question as audio instead of showing it. Set explicitly by the
  /// listening drill — never guessed from the prompt, because a saved mistake
  /// keeps its original prompt ("Listen and choose…", "Listening
  /// Comprehension"), and guessing hid the question when it came up for review.
  final bool listen;
  const _Drill({required this.kind, required this.prompt, required this.question, required this.correct, this.options = const [], this.speaker, this.mistakeId, this.listen = false});
}

List<_Drill> _buildDrills(String mode, List<VocabItem> vocab, List<Mistake> mistakes) {
  final translations = vocab.map((v) => v.translation).toSet().toList();
  List<String> options(String correct, List<String> pool) =>
      _shuffle([correct, ..._sample(pool.where((x) => x.isNotEmpty && x != correct).toList(), 3)]);

  _Drill quizDrill(VocabItem v, bool listen) => _Drill(
        kind: 'choose', prompt: listen ? 'Listen and choose the meaning' : 'What does this mean?',
        question: v.word, correct: v.translation, options: options(v.translation, translations), speaker: v.speaker, listen: listen,
      );
  _Drill speakDrill(VocabItem v) => _Drill(
        kind: 'speak', prompt: 'Say it out loud',
        question: v.example.isNotEmpty ? v.example : v.word, correct: v.example.isNotEmpty ? v.example : v.word, speaker: v.speaker,
      );

  if (mode == 'mistakes') {
    final answers = {...translations, ...vocab.map((v) => v.word), ...mistakes.map((m) => m.correctAnswer)}.toList();
    return mistakes.take(12).map((m) => _Drill(
          kind: 'choose', prompt: m.prompt.isNotEmpty ? m.prompt : 'Choose the correct answer',
          question: m.question, correct: m.correctAnswer, options: options(m.correctAnswer, answers), mistakeId: m.id,
        )).toList();
  }
  if (mode == 'speak') return _sample(vocab, min(8, vocab.length)).map(speakDrill).toList();
  if (mode == 'mix') {
    final picks = _sample(vocab, min(10, vocab.length));
    final drills = <_Drill>[
      ...picks.take(4).map((v) => quizDrill(v, false)),
      ...picks.skip(4).take(3).map((v) => quizDrill(v, true)),
      ...picks.skip(7).map(speakDrill),
    ];
    return _shuffle(drills);
  }
  return _sample(vocab, min(8, vocab.length)).map((v) => quizDrill(v, mode == 'listen')).toList();
}

class PracticeRunScreen extends ConsumerStatefulWidget {
  final String mode;
  const PracticeRunScreen({super.key, required this.mode});

  @override
  ConsumerState<PracticeRunScreen> createState() => _PracticeRunScreenState();
}

class _PracticeRunScreenState extends ConsumerState<PracticeRunScreen> {
  bool get _isSession => widget.mode == 'listening' || widget.mode == 'reading';

  bool _loading = true;
  List<_Drill> _drills = [];
  int _idx = 0;
  int _correct = 0;
  final List<int> _resolved = [];
  final List<ReviewItem> _misses = [];
  bool _reviewing = false;
  bool _done = false;
  int _xp = 0;

  // session runner state
  List<ListeningSession> _listeningSessions = [];
  List<ReadingSession> _readingSessions = [];
  int _sessionIdx = 0;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() => _loading = true);
    try {
      if (_isSession) {
        if (widget.mode == 'listening') {
          _listeningSessions = await ApiClient.instance.practiceListening();
        } else {
          _readingSessions = await ApiClient.instance.practiceReading();
        }
      } else {
        final r = await ApiClient.instance.practice();
        _drills = _buildDrills(widget.mode, r.vocab, r.mistakes);
      }
    } catch (_) {
      _drills = [];
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  Future<void> _finish(int finalCorrect, List<int> resolvedIds) async {
    final earned = min(50, 5 + finalCorrect * 3);
    setState(() {
      _xp = earned;
      _reviewing = false;
      _done = true;
    });
    try {
      final (_, user) = await ApiClient.instance.completePractice(earned);
      ref.read(authProvider.notifier).setUser(user);
    } catch (_) {}
    if (resolvedIds.isNotEmpty) {
      ApiClient.instance.resolveMistakes(resolvedIds).catchError((_) => false);
    }
  }

  void _next(bool wasCorrect) {
    final d = _drills[_idx];
    final nc = _correct + (wasCorrect ? 1 : 0);
    if (wasCorrect && d.mistakeId != null) _resolved.add(d.mistakeId!);

    if (!wasCorrect) {
      if (d.kind == 'choose' && d.mistakeId == null) {
        ApiClient.instance.recordMistake(prompt: d.prompt, question: d.question, correctAnswer: d.correct).catchError((_) {});
      }
      _misses.add(ReviewItem(prompt: d.prompt, question: d.question, correctAnswer: d.correct,
          playText: d.mistakeId == null ? d.question : null, speaker: d.speaker));
    }

    setState(() {
      _correct = nc;
      if (_idx + 1 < _drills.length) {
        _idx++;
      } else if (_misses.isNotEmpty) {
        _reviewing = true;
      } else {
        _finish(nc, _resolved);
      }
    });
  }

  @override
  void dispose() {
    Voices.instance.stopSpeaking();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final title = _kModeTitle[widget.mode] ?? 'Practice';
    return Scaffold(
      backgroundColor: LumoraColors.cream,
      body: SafeArea(
        child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 8, 16, 8),
              child: Row(children: [
                IconButton(icon: const Icon(Icons.close, color: LumoraColors.gray500), onPressed: () {
                  Voices.instance.stopSpeaking();
                  _backToPractice(context);
                }),
                Expanded(child: Text(title, style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 16))),
                if (!_isSession && !_loading && _drills.isNotEmpty && !_done)
                  Padding(
                    padding: const EdgeInsets.only(right: 8),
                    child: SizedBox(
                      width: 100,
                      child: ClipRRect(
                        borderRadius: BorderRadius.circular(LumoraRadii.full),
                        child: LinearProgressIndicator(value: _idx / _drills.length, minHeight: 6, backgroundColor: LumoraColors.gray100,
                            valueColor: const AlwaysStoppedAnimation(LumoraColors.purple)),
                      ),
                    ),
                  ),
              ]),
            ),
            Expanded(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20),
                child: _loading
                    ? const Center(child: FoxMascot(size: 110, glow: true))
                    : _isSession
                        ? _buildSessionRunner()
                        : _buildDrillRunner(),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildDrillRunner() {
    if (_reviewing) {
      return MistakesReview(items: _misses, finishLabel: 'Finish practice', onDone: () => _finish(_correct, _resolved));
    }
    if (_done) {
      return _CompleteView(correct: _correct, total: _drills.length, xp: _xp, onBack: () => _backToPractice(context));
    }
    if (_drills.isEmpty) {
      return _EmptyView(
        title: widget.mode == 'mistakes' ? 'No mistakes to review 🎉' : 'Nothing to practise yet',
        subtitle: widget.mode == 'mistakes'
            ? 'Answer some lessons — anything you miss shows up here.'
            : 'Start a lesson to unlock practice for this language.',
      );
    }
    final drill = _drills[_idx];
    return drill.kind == 'speak' ? _SpeakDrillView(key: ValueKey(_idx), drill: drill, onNext: _next) : _ChooseDrillView(key: ValueKey(_idx), drill: drill, onNext: _next);
  }

  Widget _buildSessionRunner() {
    final listening = widget.mode == 'listening';
    final count = listening ? _listeningSessions.length : _readingSessions.length;
    if (count == 0) {
      return _EmptyView(
        title: 'Nothing unlocked yet',
        subtitle: 'Complete more lessons to unlock ${listening ? "listening conversations" : "reading passages"} for practice.',
      );
    }
    final i = _sessionIdx % count;
    if (listening) {
      final s = _listeningSessions[i];
      return _ComprehensionView(
        key: ValueKey('listening-${s.id}-$_sessionIdx'),
        listening: true,
        title: s.title,
        conversation: [for (final l in s.lines) (character: l.character, text: l.text)],
        passage: const [],
        questions: [for (final q in s.questions) (question: q.question, options: q.options ?? [], correctAnswer: q.correctAnswer)],
        onScored: _finishSession,
        onAnother: () => setState(() => _sessionIdx++),
      );
    }
    final s = _readingSessions[i];
    return _ComprehensionView(
      key: ValueKey('reading-${s.id}-$_sessionIdx'),
      listening: false,
      title: s.title,
      conversation: const [],
      passage: [for (final l in s.lines) l.text],
      questions: [for (final q in s.questions) (question: q.question, options: q.options ?? [], correctAnswer: q.correctAnswer)],
      onScored: _finishSession,
      onAnother: () => setState(() => _sessionIdx++),
    );
  }

  /// Awards the session's XP; returns what was earned.
  Future<int> _finishSession(int correct, int total) async {
    final earned = min(45, 8 + correct * 4);
    try {
      final (_, user) = await ApiClient.instance.completePractice(earned);
      ref.read(authProvider.notifier).setUser(user);
    } catch (_) {}
    return earned;
  }
}

class _EmptyView extends StatelessWidget {
  final String title;
  final String subtitle;
  const _EmptyView({required this.title, required this.subtitle});
  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(mainAxisSize: MainAxisSize.min, children: [
        const FoxMascot(size: 110, glow: true),
        const SizedBox(height: 16),
        Text(title, style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 16)),
        const SizedBox(height: 4),
        Text(subtitle, textAlign: TextAlign.center, style: const TextStyle(color: LumoraColors.slatey)),
        const SizedBox(height: 20),
        LumoraButton(label: 'Go to lessons', variant: LumoraButtonVariant.outline, onPressed: () => context.go('/learn')),
      ]),
    );
  }
}

class _CompleteView extends StatelessWidget {
  final int correct;
  final int total;
  final int xp;
  final VoidCallback onBack;
  const _CompleteView({required this.correct, required this.total, required this.xp, required this.onBack});
  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(mainAxisSize: MainAxisSize.min, children: [
        const FoxMascot(size: 130, glow: true, bounce: true),
        const SizedBox(height: 16),
        const Text('Practice complete!', style: TextStyle(fontSize: 22, fontWeight: FontWeight.w800)),
        const SizedBox(height: 4),
        Text('You got $correct/$total right.', style: const TextStyle(color: LumoraColors.slatey)),
        const SizedBox(height: 16),
        Container(padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8), decoration: BoxDecoration(color: LumoraColors.amber, borderRadius: BorderRadius.circular(LumoraRadii.full)),
            child: Text('+$xp XP', style: const TextStyle(fontWeight: FontWeight.w800))),
        const SizedBox(height: 24),
        LumoraButton(label: 'Back to practice', full: true, onPressed: onBack),
      ]),
    );
  }
}

class _ChooseDrillView extends StatefulWidget {
  final _Drill drill;
  final ValueChanged<bool> onNext;
  const _ChooseDrillView({super.key, required this.drill, required this.onNext});
  @override
  State<_ChooseDrillView> createState() => _ChooseDrillViewState();
}

class _ChooseDrillViewState extends State<_ChooseDrillView> {
  String? _answer;
  bool? _correct;

  @override
  void initState() {
    super.initState();
    if (widget.drill.listen) Voices.instance.speakAs(widget.drill.speaker, widget.drill.question);
  }

  void _check() {
    if (_answer == null) return;
    setState(() => _correct = _answer == widget.drill.correct);
  }

  @override
  Widget build(BuildContext context) {
    final d = widget.drill;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        // Long questions and sentence-length answers scroll instead of being
        // cut off; the Check / Continue bar stays pinned below.
        Expanded(
          child: SingleChildScrollView(
            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Text(d.prompt.toUpperCase(), style: const TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.gray500)),
              const SizedBox(height: 12),
              if (d.listen)
                Center(
                  child: Column(children: [
                    InkWell(
                      customBorder: const CircleBorder(),
                      onTap: () => Voices.instance.speakAs(d.speaker, d.question),
                      child: Container(width: 72, height: 72, decoration: const BoxDecoration(color: LumoraColors.purple, shape: BoxShape.circle),
                          child: const Icon(Icons.volume_up_rounded, color: Colors.white, size: 30)),
                    ),
                    const SizedBox(height: 8),
                    const Text('Tap to listen again', style: TextStyle(color: LumoraColors.slatey, fontSize: 12)),
                  ]),
                )
              else
                Text(d.question, style: const TextStyle(fontSize: 22, fontWeight: FontWeight.w800, height: 1.3)),
              const SizedBox(height: 20),
              for (final opt in d.options)
                Padding(
                  padding: const EdgeInsets.only(bottom: 8),
                  child: _OptionTile(label: opt, selected: _answer == opt, isCorrect: opt == d.correct, feedback: _correct != null,
                      onTap: _correct == null ? () => setState(() => _answer = opt) : null),
                ),
            ]),
          ),
        ),
        Padding(
          padding: const EdgeInsets.symmetric(vertical: 16),
          child: _correct != null
              ? _Feedback(
                  correct: _correct!,
                  answer: d.correct,
                  onContinue: () => widget.onNext(_correct!),
                )
              : LumoraButton(label: 'Check', full: true, onPressed: _answer == null ? null : _check),
        ),
      ],
    );
  }
}

/// The result bar under a checked answer.
class _Feedback extends StatelessWidget {
  final bool correct;
  final String answer;
  final VoidCallback onContinue;
  const _Feedback({required this.correct, required this.answer, required this.onContinue});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(color: correct ? LumoraColors.tealLight : LumoraColors.coralLight, borderRadius: BorderRadius.circular(LumoraRadii.xl)),
      child: Column(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
        Text(correct ? 'Correct!' : 'Answer: $answer',
            style: TextStyle(fontWeight: FontWeight.w800, color: correct ? LumoraColors.teal : LumoraColors.coral)),
        const SizedBox(height: 10),
        LumoraButton(label: 'Continue', full: true, variant: correct ? LumoraButtonVariant.primary : LumoraButtonVariant.danger, onPressed: onContinue),
      ]),
    );
  }
}

class _OptionTile extends StatelessWidget {
  final String label;
  final bool selected;
  final bool isCorrect;
  final bool feedback;
  final VoidCallback? onTap;
  const _OptionTile({required this.label, required this.selected, required this.isCorrect, required this.feedback, required this.onTap});

  @override
  Widget build(BuildContext context) {
    Color border = LumoraColors.gray100, bg = Colors.white;
    if (feedback && isCorrect) { border = LumoraColors.teal; bg = LumoraColors.tealLight; }
    else if (feedback && selected && !isCorrect) { border = LumoraColors.coral; bg = LumoraColors.coralLight; }
    else if (selected) { border = LumoraColors.purple; bg = LumoraColors.purpleLight; }
    return Material(
      color: bg,
      borderRadius: BorderRadius.circular(LumoraRadii.md),
      child: InkWell(
        borderRadius: BorderRadius.circular(LumoraRadii.md),
        onTap: onTap,
        child: Container(
          // At least 56 tall, but free to grow: answers can be whole sentences.
          constraints: const BoxConstraints(minHeight: 56),
          alignment: Alignment.centerLeft,
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
          decoration: BoxDecoration(border: Border.all(color: border, width: 2), borderRadius: BorderRadius.circular(LumoraRadii.md)),
          child: Text(label, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w700, height: 1.3)),
        ),
      ),
    );
  }
}

class _SpeakDrillView extends StatefulWidget {
  final _Drill drill;
  final ValueChanged<bool> onNext;
  const _SpeakDrillView({super.key, required this.drill, required this.onNext});
  @override
  State<_SpeakDrillView> createState() => _SpeakDrillViewState();
}

class _SpeakDrillViewState extends State<_SpeakDrillView> {
  bool _listening = false;
  int? _score;
  String? _heard;

  Future<void> _record() async {
    if (_listening) return;
    setState(() { _listening = true; _score = null; _heard = null; });
    try {
      final said = await Voices.instance.recognizeSpeech();
      setState(() { _heard = said; _score = scorePronunciation(widget.drill.correct, said); });
    } catch (_) {
      setState(() => _heard = '');
    } finally {
      if (mounted) setState(() => _listening = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        const SizedBox(height: 20),
        const Text('SAY IT OUT LOUD', style: TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.coral)),
        const SizedBox(height: 12),
        Text('"${widget.drill.question}"', textAlign: TextAlign.center, style: const TextStyle(fontSize: 20, fontWeight: FontWeight.w800, color: LumoraColors.purple)),
        const SizedBox(height: 20),
        Row(mainAxisAlignment: MainAxisAlignment.center, children: [
          InkWell(onTap: () => Voices.instance.speakAs(widget.drill.speaker, widget.drill.question),
              child: Container(width: 48, height: 48, decoration: const BoxDecoration(color: LumoraColors.purple, shape: BoxShape.circle), child: const Icon(Icons.volume_up_rounded, color: Colors.white))),
          const SizedBox(width: 16),
          InkWell(onTap: _record,
              child: Container(width: 64, height: 64, decoration: const BoxDecoration(color: LumoraColors.coral, shape: BoxShape.circle), child: Icon(_listening ? Icons.hearing : Icons.mic, color: Colors.white, size: 28))),
        ]),
        if (_score != null) ...[
          const SizedBox(height: 16),
          Text('$_score%', style: const TextStyle(fontSize: 28, fontWeight: FontWeight.w800)),
          Text(_heard != null && _heard!.isNotEmpty ? 'You said: "$_heard"' : "Didn't catch that", style: const TextStyle(color: LumoraColors.slatey)),
        ],
        const Spacer(),
        Padding(padding: const EdgeInsets.only(bottom: 16), child: LumoraButton(label: "I said it!", full: true, onPressed: () => widget.onNext(_score != null && _score! >= 50))),
      ],
    );
  }
}

/// One reading passage or listening conversation, then its questions —
/// all on one page like the web: read (or play) it, answer every question,
/// submit. Missed questions are saved for Review Mistakes and recapped before
/// the score.
class _ComprehensionView extends StatefulWidget {
  final bool listening;
  final String title;
  final List<({String? character, String text})> conversation;
  final List<String> passage;
  final List<({String question, List<String> options, String correctAnswer})> questions;
  final Future<int> Function(int correct, int total) onScored;
  final VoidCallback onAnother;

  const _ComprehensionView({
    super.key,
    required this.listening,
    required this.title,
    required this.conversation,
    required this.passage,
    required this.questions,
    required this.onScored,
    required this.onAnother,
  });

  @override
  State<_ComprehensionView> createState() => _ComprehensionViewState();
}

class _ComprehensionViewState extends State<_ComprehensionView> {
  final Map<int, String> _answers = {};
  bool _playing = false;
  bool _played = false;
  int _speakingLine = -1;
  bool _reviewing = false;
  bool _submitted = false;
  int _score = 0;
  int _xp = 0;
  List<ReviewItem> _missed = [];

  String get _modeTitle => widget.listening ? 'Listening Comprehension' : 'Reading Comprehension';

  @override
  void dispose() {
    Voices.instance.stopSpeaking();
    super.dispose();
  }

  Future<void> _play() async {
    if (_playing) {
      await Voices.instance.stopSpeaking();
      setState(() { _playing = false; _speakingLine = -1; });
      return;
    }
    setState(() { _playing = true; _played = true; });
    await Voices.instance.speakSequence(
      widget.conversation,
      onLine: (i) { if (mounted) setState(() => _speakingLine = i); },
      shouldContinue: () => mounted && _playing,
    );
    if (mounted) setState(() { _playing = false; _speakingLine = -1; });
  }

  Future<void> _submit() async {
    Voices.instance.stopSpeaking();
    final qs = widget.questions;
    var correct = 0;
    final missed = <ReviewItem>[];
    for (var i = 0; i < qs.length; i++) {
      final q = qs[i];
      if (_answers[i] == q.correctAnswer) {
        correct++;
      } else {
        missed.add(ReviewItem(prompt: _modeTitle, question: q.question, correctAnswer: q.correctAnswer));
        // Feeds Practice → Review Mistakes, as on the web.
        ApiClient.instance
            .recordMistake(prompt: _modeTitle, question: q.question, correctAnswer: q.correctAnswer)
            .catchError((_) {});
      }
    }
    setState(() {
      _playing = false;
      _score = qs.isEmpty ? 0 : (correct / qs.length * 100).round();
      _missed = missed;
      _reviewing = missed.isNotEmpty; // study the misses before the score
      _submitted = true;
    });
    final xp = await widget.onScored(correct, qs.length);
    if (mounted) setState(() => _xp = xp);
  }

  @override
  Widget build(BuildContext context) {
    final qs = widget.questions;
    if (qs.isEmpty) {
      return const _EmptyView(title: 'Nothing to practise yet', subtitle: 'Try another mode.');
    }
    if (_reviewing) {
      return MistakesReview(items: _missed, finishLabel: 'See my score', onDone: () => setState(() => _reviewing = false));
    }
    if (_submitted) return _scoreView();

    final allAnswered = _answers.length == qs.length;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Expanded(
          child: ListView(
            padding: const EdgeInsets.only(top: 4, bottom: 8),
            children: [
              Text(widget.listening ? 'LISTEN AND ANSWER' : 'READ AND ANSWER',
                  style: const TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.gray500)),
              const SizedBox(height: 4),
              Text(widget.title, style: const TextStyle(fontSize: 20, fontWeight: FontWeight.w800, color: LumoraColors.ink)),
              const SizedBox(height: 14),
              if (widget.listening) _conversationCard() else _passageCard(),
              const SizedBox(height: 20),
              for (var i = 0; i < qs.length; i++) ...[
                Text('${i + 1}. ${qs[i].question}',
                    style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w800, color: LumoraColors.ink, height: 1.35)),
                const SizedBox(height: 8),
                for (final opt in qs[i].options)
                  Padding(
                    padding: const EdgeInsets.only(bottom: 8),
                    child: _OptionTile(
                      label: opt,
                      selected: _answers[i] == opt,
                      isCorrect: false,
                      feedback: false,
                      onTap: () => setState(() => _answers[i] = opt),
                    ),
                  ),
                const SizedBox(height: 12),
              ],
            ],
          ),
        ),
        Padding(
          padding: const EdgeInsets.symmetric(vertical: 16),
          child: LumoraButton(
            label: allAnswered ? 'Submit answers' : 'Answer all ${qs.length} questions',
            full: true,
            onPressed: allAnswered ? _submit : null,
          ),
        ),
      ],
    );
  }

  Widget _passageCard() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        border: Border.all(color: LumoraColors.gray100),
        borderRadius: BorderRadius.circular(LumoraRadii.xl),
        boxShadow: LumoraShadows.card,
      ),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        for (final (i, line) in widget.passage.indexed)
          Padding(
            padding: EdgeInsets.only(bottom: i == widget.passage.length - 1 ? 0 : 10),
            child: Text(line, style: const TextStyle(fontSize: 16, height: 1.55, color: LumoraColors.ink)),
          ),
      ]),
    );
  }

  Widget _conversationCard() {
    final speakers = {for (final l in widget.conversation) l.character ?? ''}..remove('');
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        border: Border.all(color: LumoraColors.gray100),
        borderRadius: BorderRadius.circular(LumoraRadii.xl),
        boxShadow: LumoraShadows.card,
      ),
      child: Column(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
        if (speakers.isNotEmpty)
          Wrap(spacing: 6, runSpacing: 6, children: [
            for (final name in speakers)
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                decoration: BoxDecoration(color: LumoraColors.purpleLight, borderRadius: BorderRadius.circular(LumoraRadii.full)),
                child: Row(mainAxisSize: MainAxisSize.min, children: [
                  const Icon(Icons.person_rounded, size: 14, color: LumoraColors.purple),
                  const SizedBox(width: 4),
                  Text(name, style: const TextStyle(color: LumoraColors.purple, fontWeight: FontWeight.w700, fontSize: 12)),
                ]),
              ),
          ]),
        const SizedBox(height: 12),
        LumoraButton(
          label: _playing ? 'Stop' : (_played ? 'Play again' : 'Play the conversation'),
          full: true,
          icon: Icon(_playing ? Icons.stop_rounded : Icons.volume_up_rounded, color: Colors.white, size: 20),
          onPressed: _play,
        ),
        const SizedBox(height: 8),
        Text(
          _playing && _speakingLine >= 0
              ? 'Line ${_speakingLine + 1} of ${widget.conversation.length}'
              : 'Replay as many times as you need.',
          textAlign: TextAlign.center,
          style: const TextStyle(color: LumoraColors.slatey, fontSize: 12),
        ),
      ]),
    );
  }

  Widget _scoreView() {
    final color = _score >= 80 ? LumoraColors.teal : _score >= 50 ? LumoraColors.amber : LumoraColors.coral;
    return Center(
      child: SingleChildScrollView(
        child: Column(mainAxisSize: MainAxisSize.min, children: [
          const FoxMascot(size: 130, glow: true, bounce: true),
          const SizedBox(height: 12),
          Text(_score >= 60 ? 'Nicely done!' : 'Keep going!', style: const TextStyle(fontSize: 22, fontWeight: FontWeight.w800)),
          Text('$_score%', style: TextStyle(fontSize: 32, fontWeight: FontWeight.w800, color: color)),
          Text('on "${widget.title}"', textAlign: TextAlign.center, style: const TextStyle(color: LumoraColors.slatey)),
          if (_xp > 0) ...[
            const SizedBox(height: 12),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              decoration: BoxDecoration(color: LumoraColors.amber, borderRadius: BorderRadius.circular(LumoraRadii.full)),
              child: Text('+$_xp XP', style: const TextStyle(fontWeight: FontWeight.w800)),
            ),
          ],
          const SizedBox(height: 20),
          LumoraButton(label: 'Practice another', full: true, onPressed: widget.onAnother),
          const SizedBox(height: 8),
          LumoraButton(label: 'Back to practice', full: true, variant: LumoraButtonVariant.outline, onPressed: () => _backToPractice(context)),
        ]),
      ),
    );
  }
}
