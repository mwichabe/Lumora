import 'package:flutter/material.dart';

import '../core/theme/colors.dart';
import '../core/theme/radii.dart';
import 'lumora_button.dart';

/// A single-question multiple-choice check, shared by the Listening and
/// Reading session players (their QuizPhase components are identical).
class ComprehensionQuiz extends StatefulWidget {
  final List<({String question, List<String> options, String correctAnswer})> questions;
  final Future<void> Function() onDone;

  const ComprehensionQuiz({super.key, required this.questions, required this.onDone});

  @override
  State<ComprehensionQuiz> createState() => _ComprehensionQuizState();
}

class _ComprehensionQuizState extends State<ComprehensionQuiz> {
  int _qi = 0;
  String? _answer;
  bool? _correct;
  bool _finishing = false;

  @override
  void initState() {
    super.initState();
    if (widget.questions.isEmpty) {
      WidgetsBinding.instance.addPostFrameCallback((_) => widget.onDone());
    }
  }

  void _check() {
    if (_answer == null) return;
    setState(() => _correct = _answer == widget.questions[_qi].correctAnswer);
  }

  Future<void> _next() async {
    setState(() {
      _correct = null;
      _answer = null;
    });
    if (_qi + 1 < widget.questions.length) {
      setState(() => _qi++);
    } else {
      setState(() => _finishing = true);
      await widget.onDone();
    }
  }

  @override
  Widget build(BuildContext context) {
    if (widget.questions.isEmpty) return const SizedBox.shrink();
    final q = widget.questions[_qi];

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text('COMPREHENSION · ${_qi + 1}/${widget.questions.length}',
            style: const TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.gray500)),
        const SizedBox(height: 12),
        Text(q.question, style: const TextStyle(fontSize: 20, fontWeight: FontWeight.w700)),
        const SizedBox(height: 20),
        for (final opt in q.options)
          Padding(
            padding: const EdgeInsets.only(bottom: 8),
            child: _buildOption(opt, q.correctAnswer),
          ),
        const Spacer(),
        Padding(
          padding: const EdgeInsets.only(top: 16, bottom: 8),
          child: _correct != null
              ? Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: _correct! ? LumoraColors.tealLight : LumoraColors.coralLight,
                    borderRadius: BorderRadius.circular(LumoraRadii.lg),
                  ),
                  child: Column(children: [
                    Text(_correct! ? 'Correct!' : 'Answer: ${q.correctAnswer}',
                        style: TextStyle(fontWeight: FontWeight.w800, color: _correct! ? LumoraColors.teal : LumoraColors.coral)),
                    const SizedBox(height: 12),
                    LumoraButton(
                      label: 'Continue',
                      full: true,
                      loading: _finishing,
                      variant: _correct! ? LumoraButtonVariant.primary : LumoraButtonVariant.danger,
                      onPressed: _next,
                    ),
                  ]),
                )
              : LumoraButton(label: 'Check', full: true, onPressed: _answer == null ? null : _check),
        ),
      ],
    );
  }

  Widget _buildOption(String opt, String correctAnswer) {
    final selected = _answer == opt;
    final isCorrect = opt == correctAnswer;
    Color border = LumoraColors.gray100;
    Color bg = Colors.white;
    if (_correct != null && isCorrect) {
      border = LumoraColors.teal;
      bg = LumoraColors.tealLight;
    } else if (_correct != null && selected && !isCorrect) {
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
        onTap: _correct == null ? () => setState(() => _answer = opt) : null,
        child: Container(
          height: 56,
          alignment: Alignment.centerLeft,
          padding: const EdgeInsets.symmetric(horizontal: 16),
          decoration: BoxDecoration(border: Border.all(color: border, width: 2), borderRadius: BorderRadius.circular(LumoraRadii.md)),
          child: Text(opt, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w700)),
        ),
      ),
    );
  }
}

/// Shared "Session complete!" screen for Listening/Reading.
class SessionDone extends StatelessWidget {
  final String title;
  final String subtitle;
  final int xp;
  final VoidCallback onContinue;
  const SessionDone({super.key, required this.title, required this.subtitle, required this.xp, required this.onContinue});

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        const Icon(Icons.celebration_rounded, size: 72, color: LumoraColors.amber),
        const SizedBox(height: 16),
        Text(title, style: const TextStyle(fontSize: 24, fontWeight: FontWeight.w800)),
        const SizedBox(height: 4),
        Text(subtitle, style: const TextStyle(color: LumoraColors.slatey)),
        const SizedBox(height: 20),
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
          decoration: BoxDecoration(color: LumoraColors.amber, borderRadius: BorderRadius.circular(LumoraRadii.full)),
          child: Text('+$xp XP', style: const TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.ink)),
        ),
        const SizedBox(height: 32),
        LumoraButton(label: 'Back to course', full: true, onPressed: onContinue),
      ],
    );
  }
}
