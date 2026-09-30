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

enum _Phase { read, quiz, done }

class ReadingScreen extends ConsumerStatefulWidget {
  final int id;
  const ReadingScreen({super.key, required this.id});

  @override
  ConsumerState<ReadingScreen> createState() => _ReadingScreenState();
}

class _ReadingScreenState extends ConsumerState<ReadingScreen> {
  ReadingSession? _session;
  _Phase _phase = _Phase.read;
  bool _revealed = false;

  @override
  void initState() {
    super.initState();
    ApiClient.instance.readingSession(widget.id).then((s) {
      if (mounted) setState(() => _session = s);
    });
  }

  @override
  void dispose() {
    Voices.instance.stopSpeaking();
    super.dispose();
  }

  Future<void> _complete() async {
    try {
      final (_, user) = await ApiClient.instance.completeReading(widget.id);
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
                      Text('READING · ${session.unit}',
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
                  _Phase.read => _buildRead(session),
                  _Phase.quiz => ComprehensionQuiz(
                      questions: [
                        for (final q in session.questions)
                          (question: q.question, options: q.options ?? [], correctAnswer: q.correctAnswer),
                      ],
                      onDone: _complete,
                    ),
                  _Phase.done => SessionDone(
                      title: 'Reading complete!',
                      subtitle: 'You read real text. ¡Muy bien!',
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

  Widget _buildRead(ReadingSession session) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(session.description, style: const TextStyle(color: LumoraColors.slatey)),
        const SizedBox(height: 8),
        Row(children: const [
          CircleAvatar(radius: 12, backgroundColor: LumoraColors.purple, child: Text('M', style: TextStyle(color: Colors.white, fontSize: 11))),
          SizedBox(width: 8),
          Text('Mira reads with you', style: TextStyle(color: LumoraColors.slatey, fontSize: 12)),
        ]),
        const SizedBox(height: 12),
        Expanded(
          child: ListView.separated(
            itemCount: session.lines.length,
            separatorBuilder: (_, _) => const SizedBox(height: 10),
            itemBuilder: (context, i) {
              final l = session.lines[i];
              return Container(
                padding: const EdgeInsets.all(14),
                decoration: BoxDecoration(color: Colors.white, border: Border.all(color: LumoraColors.gray100), borderRadius: BorderRadius.circular(LumoraRadii.xl)),
                child: Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
                  Expanded(
                    child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                      Text(l.text, style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 15, height: 1.4)),
                      if (_revealed) Text(l.translation, style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
                    ]),
                  ),
                  IconButton(icon: const Icon(Icons.volume_up_rounded, size: 18, color: LumoraColors.purple),
                      onPressed: () => Voices.instance.speakAs('Mira', l.text)),
                ]),
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
          child: LumoraButton(label: "I've read it — quiz me", full: true, onPressed: () => setState(() => _phase = _Phase.quiz)),
        ),
      ],
    );
  }
}
