import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/characters.dart';
import '../../core/grading.dart';
import '../../core/network/api_exception.dart';
import '../../core/network/api_client.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../core/theme/shadows.dart';
import '../../core/voices.dart';
import '../../models/lesson.dart';
import '../../providers/hearts_provider.dart';
import '../../providers/auth_provider.dart';
import '../../providers/learn_provider.dart';
import '../../providers/home_provider.dart';
import '../payments/checkout_screen.dart';
import '../../widgets/fox_mascot.dart';
import '../../widgets/lumora_button.dart';
import '../../widgets/mistakes_review.dart';
import '../../widgets/out_of_hearts_modal.dart';
import '../../widgets/stat_widgets.dart';
import 'lesson_complete_screen.dart';

enum _Phase { vocab, practice, review }
enum _Feedback { none, correct, incorrect }

class LessonScreen extends ConsumerStatefulWidget {
  final int id;
  const LessonScreen({super.key, required this.id});

  @override
  ConsumerState<LessonScreen> createState() => _LessonScreenState();
}

class _LessonScreenState extends ConsumerState<LessonScreen> {
  Lesson? _lesson;
  _Phase _phase = _Phase.vocab;
  int _idx = 0;
  final List<ReviewItem> _misses = [];
  String _answer = '';
  _Feedback _feedback = _Feedback.none;
  int _correctCount = 0;
  int _gradedCount = 0;
  bool _submitting = false;
  bool _outOfHeartsShown = false;
  // What the feedback bar explains beyond right/wrong.
  _Detail _detail = const _Detail();
  // Problems with a free-writing answer that must be fixed before it counts.
  List<String> _writeProblems = const [];
  bool _checking = false;

  @override
  void initState() {
    super.initState();
    ApiClient.instance.lesson(widget.id).then((lesson) {
      if (!mounted) return;
      setState(() {
        _lesson = lesson;
        _phase = lesson.vocab.isNotEmpty ? _Phase.vocab : _Phase.practice;
      });
    });
  }

  @override
  void dispose() {
    Voices.instance.stopSpeaking();
    super.dispose();
  }

  void _maybeShowOutOfHearts(HeartsState hearts) {
    if (_outOfHeartsShown || hearts.status == null || hearts.hearts > 0) return;
    _outOfHeartsShown = true;
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted) return;
      showOutOfHeartsModal(
        context,
        status: hearts.status!,
        secondsToNext: hearts.secondsToNext,
        note: 'You ran out of hearts, so this lesson has ended. Refill to try again now, or wait for a heart and restart it.',
        closeLabel: 'Back to lessons',
        onBuy: _refillHearts,
        onClose: () {
          Voices.instance.stopSpeaking();
          Navigator.pop(context);
          context.go('/learn');
        },
      );
    });
  }

  /// Pays for a refill in-app. The lesson has already ended by the time the
  /// sheet is up, so a successful refill sends the learner back to the course
  /// to restart it with full hearts.
  Future<void> _refillHearts() async {
    final messenger = ScaffoldMessenger.of(context);
    final paid = await payWithPaystack(context, product: 'hearts');
    // Re-read either way: a payment that settles late is applied by the server.
    await ref.read(heartsProvider.notifier).reload();
    if (!mounted || !paid) return;
    Voices.instance.stopSpeaking();
    Navigator.pop(context); // the out-of-hearts sheet
    context.go('/learn');
    messenger.showSnackBar(const SnackBar(content: Text('Hearts refilled ❤️ — jump back in!')));
  }

  String _normalise(String s) => s.trim().toLowerCase().replaceAll(RegExp(r'[.,!¡¿?]'), '');

  Exercise? get _ex {
    final exercises = _lesson?.exercises ?? [];
    return _idx < exercises.length ? exercises[_idx] : null;
  }

  bool _isTyped(Exercise ex) =>
      [ExerciseType.translate, ExerciseType.fill].contains(ex.type) && (ex.options == null || ex.options!.isEmpty);

  Future<void> _check() async {
    final ex = _ex;
    if (ex == null) return;
    final isNarrative = ex.type == ExerciseType.character;
    final isSpeak = ex.type == ExerciseType.speak;
    if (isNarrative || isSpeak) {
      _advance();
      return;
    }
    if (ex.type == ExerciseType.write) return _checkWriting(ex);
    if (_isTyped(ex)) return _checkTyped(ex);
    if (_normalise(_answer) == _normalise(ex.correctAnswer)) {
      _markCorrect();
    } else {
      _markWrong(ex);
    }
  }

  void _markCorrect([_Detail d = const _Detail()]) {
    setState(() {
      _gradedCount++;
      _correctCount++;
      _detail = d;
      _feedback = _Feedback.correct;
    });
  }

  /// Typed answers: point out exactly what's wrong, accept accent slips and
  /// one-letter typos, and ask the server about valid alternative wordings.
  Future<void> _checkTyped(Exercise ex) async {
    final g = gradeTyped(_answer, ex.correctAnswer);
    if (g.verdict != Verdict.incorrect) {
      return _markCorrect(_Detail(note: g.note.isEmpty ? null : g.note, lines: issueLines(g.issues)));
    }
    if (ex.correctAnswer.trim().split(RegExp(r'\s+')).length >= 2) {
      setState(() => _checking = true);
      try {
        final second = await ApiClient.instance
            .checkAnswer(lessonId: ex.lessonId, question: ex.question, expected: ex.correctAnswer, answer: _answer)
            .timeout(const Duration(seconds: 9));
        if (second.available && second.correct) {
          if (!mounted) return;
          setState(() => _checking = false);
          return _markCorrect(_Detail(
            note: second.explanation.isEmpty ? "That's another correct way to say it." : second.explanation,
            lines: ['Also correct: “${ex.correctAnswer}”'],
          ));
        }
      } catch (_) {
        // No second opinion — keep the local verdict.
      }
      if (!mounted) return;
      setState(() => _checking = false);
    }
    _markWrong(ex, _Detail(note: g.note, lines: issueLines(g.issues)));
  }

  /// Free writing: rule checks (no copying the example, long enough, right
  /// language) must pass before it counts; then corrections are shown.
  Future<void> _checkWriting(Exercise ex) async {
    if (copiesExample(_answer, ex.correctAnswer)) {
      setState(() => _writeProblems = const ['An example cannot be used here — write your own answer in your own words.']);
      return;
    }
    setState(() {
      _writeProblems = const [];
      _checking = true;
    });
    try {
      final r = await ApiClient.instance.checkWriting(ex.id, _answer);
      if (!mounted) return;
      if (!r.ok) {
        setState(() => _writeProblems = r.problems);
        return;
      }
      final d = _Detail(note: r.summary.isEmpty ? null : r.summary, corrections: r.corrections, hideAnswer: true);
      r.acceptable ? _markCorrect(d) : _markWrong(ex, d);
    } on ApiException catch (_) {
      if (mounted) _markCorrect(const _Detail(note: "Saved. Detailed feedback isn't available right now.", hideAnswer: true));
    } catch (_) {
      // Offline: the local checks passed, so accept it.
      if (mounted) _markCorrect(const _Detail(note: "Saved. Detailed feedback isn't available right now.", hideAnswer: true));
    } finally {
      if (mounted) setState(() => _checking = false);
    }
  }

  void _markWrong(Exercise ex, [_Detail d = const _Detail()]) {
    setState(() {
      _gradedCount++;
      _detail = d;
      _feedback = _Feedback.incorrect;
      _misses.add(ReviewItem(
        prompt: ex.prompt,
        question: ex.question,
        correctAnswer: ex.correctAnswer,
        playText: ex.type == ExerciseType.listen ? ex.question : null,
        speaker: ex.character,
      ));
    });
    ApiClient.instance.recordMistake(prompt: ex.prompt, question: ex.question, correctAnswer: ex.correctAnswer).catchError((_) {});
    ref.read(heartsProvider.notifier).lose().then((s) {
      if (s != null && s.hearts <= 0) _maybeShowOutOfHearts(HeartsState(status: s, secondsToNext: s.secondsToNext));
    });
  }

  void _advance() {
    final total = _lesson?.exercises.length ?? 0;
    setState(() {
      _feedback = _Feedback.none;
      _detail = const _Detail();
      _writeProblems = const [];
      _answer = '';
      if (_idx + 1 < total) {
        _idx++;
      } else if (_misses.isNotEmpty) {
        _phase = _Phase.review;
      } else {
        _finish();
      }
    });
  }

  Future<void> _finish() async {
    setState(() => _submitting = true);
    final accuracy = _gradedCount > 0 ? ((_correctCount / _gradedCount) * 100).round() : 100;
    try {
      final (xpEarned, resultAccuracy, user, firstClear) = await ApiClient.instance.completeLesson(widget.id, accuracy);
      LessonCompleteScreen.pendingResult = LessonResult(xp: xpEarned, accuracy: resultAccuracy, firstClear: firstClear);
      // Start refreshing Home and the course now, while the celebration plays,
      // so "continue where you left off" and the quest bars have moved on by
      // the time the learner taps Continue.
      ref.read(homeProvider.notifier).refresh();
      ref.read(learnProvider.notifier).load();
      ref.read(authProvider.notifier).setUser(user);
    } catch (_) {
      LessonCompleteScreen.pendingResult = LessonResult(xp: _lesson?.exercises.length ?? 0, accuracy: accuracy, firstClear: true);
    } finally {
      if (mounted) context.pushReplacement('/lesson/${widget.id}/complete');
    }
  }

  @override
  Widget build(BuildContext context) {
    final hearts = ref.watch(heartsProvider);
    ref.listen(heartsProvider, (_, next) => _maybeShowOutOfHearts(next));

    if (_lesson == null) {
      return const Scaffold(body: Center(child: FoxMascot(size: 120, glow: true)));
    }

    final total = _lesson!.exercises.length;
    final progress = total > 0 ? _idx / total : 0.0;

    return Scaffold(
      backgroundColor: LumoraColors.cream,
      body: SafeArea(
        child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 8, 16, 8),
              child: Row(children: [
                IconButton(icon: const Icon(Icons.close, color: LumoraColors.gray500), onPressed: () => context.go('/home')),
                Expanded(
                  child: ClipRRect(
                    borderRadius: BorderRadius.circular(LumoraRadii.full),
                    child: LinearProgressIndicator(value: progress, minHeight: 8, backgroundColor: LumoraColors.gray100,
                        valueColor: const AlwaysStoppedAnimation(LumoraColors.purple)),
                  ),
                ),
                const SizedBox(width: 10),
                HeartIndicator(hearts: hearts.hearts, secondsToNext: hearts.secondsToNext),
              ]),
            ),
            Expanded(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20),
                child: switch (_phase) {
                  _Phase.vocab => _VocabPhase(vocab: _lesson!.vocab, onDone: () => setState(() => _phase = _Phase.practice)),
                  _Phase.review => MistakesReview(items: _misses, finishLabel: _submitting ? 'Finishing…' : 'Finish lesson', onDone: _finish),
                  _Phase.practice => _buildPractice(),
                },
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildPractice() {
    final ex = _ex;
    if (ex == null) return const SizedBox.shrink();
    final isNarrative = ex.type == ExerciseType.character;
    final isSpeak = ex.type == ExerciseType.speak;
    final isWrite = ex.type == ExerciseType.write;
    final hasOptions = ex.options != null && ex.options!.isNotEmpty;
    final needsChoice = !isNarrative && !isSpeak && !isWrite &&
        ([ExerciseType.multipleChoice, ExerciseType.listen, ExerciseType.match].contains(ex.type) || hasOptions);
    final needsTyping = [ExerciseType.translate, ExerciseType.fill].contains(ex.type) && !hasOptions;
    final wordCount = countWords(_answer);

    final canCheck = _feedback != _Feedback.none || _checking
        ? false
        : (isNarrative || isSpeak)
            ? true
            : _answer.trim().isNotEmpty;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        if (!isNarrative)
          Text(ex.prompt.toUpperCase(), style: const TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.gray500)),
        const SizedBox(height: 8),
        if (isNarrative) _NarrativeCard(ex: ex) else _QuestionArea(ex: ex),
        if (!isNarrative) ...[
          const SizedBox(height: 16),
          if (needsChoice)
            _ChoiceList(ex: ex, answer: _answer, feedback: _feedback, onSelect: (v) {
              if (_feedback == _Feedback.none) setState(() => _answer = v);
            }),
          if (needsTyping)
            TextField(
              key: ValueKey('typed-$_idx'), // a fresh field per exercise
              autofocus: true,
              enabled: _feedback == _Feedback.none,
              autocorrect: false,
              enableSuggestions: false,
              textInputAction: TextInputAction.done,
              onSubmitted: (_) {
                if (canCheck) _check();
              },
              onChanged: (v) => setState(() => _answer = v),
              decoration: const InputDecoration(hintText: 'Type your answer…'),
            ),
          if (isSpeak) _SpeakControl(phrase: ex.question, disabled: _feedback != _Feedback.none),
          if (isWrite)
            _WriteControl(
              key: ValueKey('write-$_idx'),
              onChange: (v) => setState(() {
                _answer = v;
                if (_writeProblems.isNotEmpty) _writeProblems = const [];
              }),
              words: wordCount,
              example: ex.correctAnswer,
              problems: _writeProblems,
              enabled: _feedback == _Feedback.none,
            ),
        ],
        const Spacer(),
        Padding(
          padding: const EdgeInsets.only(bottom: 24, top: 12),
          child: _feedback != _Feedback.none
              ? _FeedbackBar(correct: _feedback == _Feedback.correct, correctAnswer: ex.correctAnswer, detail: _detail, onContinue: _advance)
              : LumoraButton(
                  label: isNarrative ? 'Continue' : isSpeak ? "I said it!" : isWrite ? 'Check my writing' : 'Check',
                  full: true,
                  loading: _submitting || _checking,
                  onPressed: canCheck ? _check : null,
                ),
        ),
      ],
    );
  }
}

class _VocabPhase extends StatefulWidget {
  final List<VocabItem> vocab;
  final VoidCallback onDone;
  const _VocabPhase({required this.vocab, required this.onDone});

  @override
  State<_VocabPhase> createState() => _VocabPhaseState();
}

class _VocabPhaseState extends State<_VocabPhase> {
  int _i = 0;

  @override
  void initState() {
    super.initState();
    _speak();
  }

  void _speak() {
    final item = widget.vocab[_i];
    Voices.instance.speakAs(item.speaker.isEmpty ? 'Lumora' : item.speaker, item.word);
  }

  @override
  Widget build(BuildContext context) {
    final item = widget.vocab[_i];
    final last = _i == widget.vocab.length - 1;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text('NEW WORDS · ${_i + 1}/${widget.vocab.length}',
            style: const TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.gray500)),
        Expanded(
          child: Center(
            child: Container(
              width: double.infinity,
              padding: const EdgeInsets.all(24),
              decoration: BoxDecoration(color: Colors.white, border: Border.all(color: LumoraColors.gray100),
                  borderRadius: BorderRadius.circular(LumoraRadii.xl2), boxShadow: LumoraShadows.cardLg),
              child: Column(mainAxisSize: MainAxisSize.min, children: [
                Text(item.word, textAlign: TextAlign.center, style: const TextStyle(fontSize: 28, fontWeight: FontWeight.w800)),
                const SizedBox(height: 4),
                Text(item.translation, style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w800, color: LumoraColors.purple)),
                const SizedBox(height: 16),
                InkWell(
                  onTap: _speak,
                  child: Container(width: 48, height: 48, decoration: const BoxDecoration(color: LumoraColors.purple, shape: BoxShape.circle),
                      child: const Icon(Icons.volume_up_rounded, color: Colors.white)),
                ),
                if (item.example.isNotEmpty) ...[
                  const SizedBox(height: 20),
                  InkWell(
                    onTap: () => Voices.instance.speakAs(item.speaker, item.example),
                    child: Container(
                      width: double.infinity,
                      padding: const EdgeInsets.all(12),
                      decoration: BoxDecoration(color: LumoraColors.gray50, borderRadius: BorderRadius.circular(LumoraRadii.lg)),
                      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                        Text(item.example, style: const TextStyle(fontWeight: FontWeight.w700)),
                        Text(item.exampleTranslation, style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
                      ]),
                    ),
                  ),
                ],
              ]),
            ),
          ),
        ),
        Padding(
          padding: const EdgeInsets.symmetric(vertical: 12),
          child: Row(mainAxisAlignment: MainAxisAlignment.center, children: [
            for (var n = 0; n < widget.vocab.length; n++)
              AnimatedContainer(
                duration: const Duration(milliseconds: 200),
                margin: const EdgeInsets.symmetric(horizontal: 2),
                width: n == _i ? 20 : 6,
                height: 6,
                decoration: BoxDecoration(color: n == _i ? LumoraColors.purple : LumoraColors.gray100, borderRadius: BorderRadius.circular(3)),
              ),
          ]),
        ),
        Padding(
          padding: const EdgeInsets.only(bottom: 24),
          child: Row(children: [
            if (_i > 0) ...[
              Expanded(child: LumoraButton(label: 'Back', variant: LumoraButtonVariant.outline, full: true, onPressed: () => setState(() { _i--; _speak(); }))),
              const SizedBox(width: 12),
            ],
            Expanded(
              child: LumoraButton(
                label: last ? 'Start practice' : 'Next word',
                full: true,
                onPressed: () {
                  if (last) {
                    widget.onDone();
                  } else {
                    setState(() { _i++; _speak(); });
                  }
                },
              ),
            ),
          ]),
        ),
      ],
    );
  }
}

class _NarrativeCard extends StatefulWidget {
  final Exercise ex;
  const _NarrativeCard({required this.ex});
  @override
  State<_NarrativeCard> createState() => _NarrativeCardState();
}

class _NarrativeCardState extends State<_NarrativeCard> {
  @override
  void initState() {
    super.initState();
    if (widget.ex.question.isNotEmpty) Voices.instance.speakAs(widget.ex.character, widget.ex.question);
  }

  @override
  Widget build(BuildContext context) {
    final info = characterInfo(widget.ex.character);
    final color = _hexColor(info.color);
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 24),
      child: Column(children: [
        Container(
          padding: const EdgeInsets.all(3),
          decoration: BoxDecoration(shape: BoxShape.circle, border: Border.all(color: color, width: 3)),
          child: CircleAvatar(radius: 46, backgroundColor: color, child: Text(info.name.isNotEmpty ? info.name[0] : '?', style: const TextStyle(color: Colors.white, fontSize: 32, fontWeight: FontWeight.w800))),
        ),
        const SizedBox(height: 8),
        Text(info.name, style: const TextStyle(fontWeight: FontWeight.w800)),
        const SizedBox(height: 12),
        SpeechBubble(child: Text(widget.ex.question, style: const TextStyle(fontSize: 16))),
        const SizedBox(height: 12),
        TextButton.icon(
          onPressed: () => Voices.instance.speakAs(widget.ex.character, widget.ex.question),
          icon: const Icon(Icons.volume_up_rounded, size: 16),
          label: const Text('Replay'),
        ),
      ]),
    );
  }
}

class _QuestionArea extends StatelessWidget {
  final Exercise ex;
  const _QuestionArea({required this.ex});
  @override
  Widget build(BuildContext context) {
    if (ex.type == ExerciseType.listen) {
      return Padding(
        padding: const EdgeInsets.symmetric(vertical: 16),
        child: Column(children: [
          InkWell(
            onTap: () => Voices.instance.speakAs('Mira', ex.question),
            child: Container(width: 64, height: 64, decoration: const BoxDecoration(color: LumoraColors.purple, shape: BoxShape.circle),
                child: const Icon(Icons.volume_up_rounded, color: Colors.white, size: 28)),
          ),
          const SizedBox(height: 8),
          const Text('Tap to listen, then choose the meaning', style: TextStyle(color: LumoraColors.slatey, fontSize: 12)),
        ]),
      );
    }
    return Text(ex.question, style: const TextStyle(fontSize: 20, fontWeight: FontWeight.w700));
  }
}

class _ChoiceList extends StatelessWidget {
  final Exercise ex;
  final String answer;
  final _Feedback feedback;
  final ValueChanged<String> onSelect;
  const _ChoiceList({required this.ex, required this.answer, required this.feedback, required this.onSelect});

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        for (final opt in ex.options ?? [])
          Padding(
            padding: const EdgeInsets.only(bottom: 8),
            child: _ChoiceButton(
              label: opt,
              selected: answer == opt,
              isCorrect: opt == ex.correctAnswer,
              feedback: feedback,
              onTap: () => onSelect(opt),
            ),
          ),
      ],
    );
  }
}

class _ChoiceButton extends StatelessWidget {
  final String label;
  final bool selected;
  final bool isCorrect;
  final _Feedback feedback;
  final VoidCallback onTap;
  const _ChoiceButton({required this.label, required this.selected, required this.isCorrect, required this.feedback, required this.onTap});

  @override
  Widget build(BuildContext context) {
    Color border = LumoraColors.gray100;
    Color bg = Colors.white;
    if (feedback != _Feedback.none && isCorrect) {
      border = LumoraColors.teal;
      bg = LumoraColors.tealLight;
    } else if (feedback != _Feedback.none && selected && !isCorrect) {
      border = LumoraColors.coral;
      bg = LumoraColors.coralLight;
    } else if (selected) {
      border = LumoraColors.purple;
      bg = LumoraColors.purpleLight;
    }
    return Material(
      color: bg,
      borderRadius: BorderRadius.circular(LumoraRadii.md),
      child: InkWell(
        borderRadius: BorderRadius.circular(LumoraRadii.md),
        onTap: onTap,
        child: Container(
          height: 56,
          padding: const EdgeInsets.symmetric(horizontal: 16),
          alignment: Alignment.centerLeft,
          decoration: BoxDecoration(border: Border.all(color: border, width: 2), borderRadius: BorderRadius.circular(LumoraRadii.md)),
          child: Text(label, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w700)),
        ),
      ),
    );
  }
}

class _WriteControl extends StatefulWidget {
  final ValueChanged<String> onChange;
  final int words;
  final String example;
  /// Why the answer can't be accepted yet (copied example, too short…).
  final List<String> problems;
  final bool enabled;
  const _WriteControl({
    super.key,
    required this.onChange,
    required this.words,
    required this.example,
    required this.problems,
    required this.enabled,
  });

  @override
  State<_WriteControl> createState() => _WriteControlState();
}

class _WriteControlState extends State<_WriteControl> {
  bool _showExample = false;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        TextField(
          maxLines: 6,
          enabled: widget.enabled,
          onChanged: widget.onChange,
          decoration: InputDecoration(
            hintText: 'Write your answer here…',
            enabledBorder: widget.problems.isEmpty
                ? null
                : OutlineInputBorder(borderSide: const BorderSide(color: LumoraColors.coral, width: 1.5), borderRadius: BorderRadius.circular(LumoraRadii.lg)),
          ),
        ),
        if (widget.problems.isNotEmpty)
          Container(
            width: double.infinity,
            margin: const EdgeInsets.only(top: 8),
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(color: LumoraColors.coralLight, borderRadius: BorderRadius.circular(LumoraRadii.lg)),
            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              for (final p in widget.problems)
                Padding(
                  padding: const EdgeInsets.symmetric(vertical: 2),
                  child: Text(p, style: const TextStyle(color: LumoraColors.coral, fontWeight: FontWeight.w700, fontSize: 13)),
                ),
            ]),
          ),
        const SizedBox(height: 8),
        Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
          Text('${widget.words} words', style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
          if (widget.example.isNotEmpty)
            TextButton(onPressed: () => setState(() => _showExample = !_showExample),
                child: Text(_showExample ? 'Hide example' : 'Show example', style: const TextStyle(color: LumoraColors.teal))),
        ]),
        if (_showExample && widget.example.isNotEmpty)
          Container(
            width: double.infinity,
            margin: const EdgeInsets.only(top: 4),
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(color: LumoraColors.gray50, borderRadius: BorderRadius.circular(LumoraRadii.lg)),
            child: Text(widget.example, style: const TextStyle(color: LumoraColors.slatey)),
          ),
      ],
    );
  }
}

class _SpeakControl extends StatefulWidget {
  final String phrase;
  final bool disabled;
  const _SpeakControl({required this.phrase, required this.disabled});

  @override
  State<_SpeakControl> createState() => _SpeakControlState();
}

class _SpeakControlState extends State<_SpeakControl> {
  bool _listening = false;
  String? _heard;
  int? _score;

  Future<void> _record() async {
    if (_listening || widget.disabled) return;
    setState(() {
      _listening = true;
      _heard = null;
      _score = null;
    });
    try {
      final said = await Voices.instance.recognizeSpeech();
      setState(() {
        _heard = said;
        _score = scorePronunciation(widget.phrase, said);
      });
    } catch (_) {
      setState(() => _heard = '');
    } finally {
      if (mounted) setState(() => _listening = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final color = _score == null ? LumoraColors.ink : _score! >= 80 ? LumoraColors.teal : _score! >= 50 ? LumoraColors.amber : LumoraColors.coral;
    final label = _score == null ? '' : _score! >= 80 ? 'Excellent!' : _score! >= 50 ? 'Good try!' : 'Keep practising';

    return Column(children: [
      const Text('SPEAKING PRACTICE', style: TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.coral)),
      const SizedBox(height: 8),
      Text('"${widget.phrase}"', style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w800, color: LumoraColors.purple)),
      const SizedBox(height: 16),
      Row(mainAxisAlignment: MainAxisAlignment.center, children: [
        InkWell(
          onTap: widget.disabled ? null : () => Voices.instance.speakAs('Lumora', widget.phrase),
          child: Container(width: 48, height: 48, decoration: const BoxDecoration(color: LumoraColors.purple, shape: BoxShape.circle),
              child: const Icon(Icons.volume_up_rounded, color: Colors.white)),
        ),
        const SizedBox(width: 16),
        InkWell(
          onTap: _record,
          child: Container(width: 64, height: 64, decoration: const BoxDecoration(color: LumoraColors.coral, shape: BoxShape.circle),
              child: Icon(_listening ? Icons.hearing : Icons.mic, color: Colors.white, size: 28)),
        ),
      ]),
      const SizedBox(height: 8),
      Text(_listening ? 'Listening… say the phrase' : 'Tap the mic and say it out loud', style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
      if (_score != null) ...[
        const SizedBox(height: 12),
        Container(
          width: double.infinity,
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(color: LumoraColors.gray50, borderRadius: BorderRadius.circular(LumoraRadii.lg)),
          child: Column(children: [
            Text('$_score%', style: TextStyle(fontSize: 28, fontWeight: FontWeight.w800, color: color)),
            Text('$label · FLUENCY', style: TextStyle(color: color, fontSize: 11, fontWeight: FontWeight.w800)),
            const SizedBox(height: 4),
            Text(_heard != null && _heard!.isNotEmpty ? 'You said: "$_heard"' : "Didn't catch that — try again.",
                style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
          ]),
        ),
      ],
    ]);
  }
}

/// What the feedback bar explains beyond right/wrong.
class _Detail {
  /// Headline under "Nailed it!" / "Not quite" (e.g. "Watch the accents").
  final String? note;
  /// Pointed-out errors in a typed answer, one per line.
  final List<String> lines;
  /// Writing-coach corrections on free writing.
  final List<WritingCorrection> corrections;
  /// Hide "Answer: …" (free writing has no single answer).
  final bool hideAnswer;
  const _Detail({this.note, this.lines = const [], this.corrections = const [], this.hideAnswer = false});
}

class _FeedbackBar extends StatelessWidget {
  final bool correct;
  final String correctAnswer;
  final _Detail detail;
  final VoidCallback onContinue;
  const _FeedbackBar({required this.correct, required this.correctAnswer, required this.detail, required this.onContinue});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: correct ? LumoraColors.tealLight : LumoraColors.coralLight, borderRadius: BorderRadius.circular(LumoraRadii.lg)),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            Container(width: 32, height: 32, decoration: BoxDecoration(color: correct ? LumoraColors.teal : LumoraColors.coral, shape: BoxShape.circle),
                child: Icon(correct ? Icons.check : Icons.close, color: Colors.white, size: 18)),
            const SizedBox(width: 10),
            Expanded(
              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Text(correct ? 'Nailed it!' : 'Not quite', style: TextStyle(fontWeight: FontWeight.w800, color: correct ? LumoraColors.teal : LumoraColors.coral)),
                if (!correct && !detail.hideAnswer)
                  Text.rich(TextSpan(text: 'Answer: ', children: [TextSpan(text: correctAnswer, style: const TextStyle(fontWeight: FontWeight.w800))])),
              ]),
            ),
          ]),
          if (detail.note != null) ...[
            const SizedBox(height: 8),
            Text(detail.note!, style: const TextStyle(color: LumoraColors.ink, fontSize: 13)),
          ],
          for (final l in detail.lines)
            Padding(
              padding: const EdgeInsets.only(top: 2),
              child: Text('• $l', style: const TextStyle(color: LumoraColors.ink, fontSize: 13)),
            ),
          if (detail.corrections.isNotEmpty)
            ConstrainedBox(
              constraints: const BoxConstraints(maxHeight: 220),
              child: ListView(
                shrinkWrap: true,
                padding: const EdgeInsets.only(top: 8),
                children: [
                  for (final c in detail.corrections)
                    Container(
                      margin: const EdgeInsets.only(bottom: 6),
                      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
                      decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.7), borderRadius: BorderRadius.circular(LumoraRadii.md)),
                      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                        Text.rich(TextSpan(children: [
                          TextSpan(text: c.original, style: const TextStyle(color: LumoraColors.coral, decoration: TextDecoration.lineThrough)),
                          const TextSpan(text: '  →  '),
                          TextSpan(text: c.correction, style: const TextStyle(color: LumoraColors.teal, fontWeight: FontWeight.w800)),
                        ]), style: const TextStyle(fontSize: 13)),
                        Text(c.explanation, style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
                      ]),
                    ),
                ],
              ),
            ),
          const SizedBox(height: 12),
          LumoraButton(label: correct ? 'Continue' : 'Got it', full: true, variant: correct ? LumoraButtonVariant.primary : LumoraButtonVariant.danger, onPressed: onContinue),
        ],
      ),
    );
  }
}

Color _hexColor(String hex) {
  var h = hex.replaceAll('#', '');
  if (h.length == 6) h = 'FF$h';
  final v = int.tryParse(h, radix: 16);
  return v == null ? LumoraColors.purple : Color(v);
}
