import 'package:confetti/confetti.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../widgets/fox_mascot.dart';
import '../../widgets/lumora_button.dart';

class LessonResult {
  final int xp;
  final int accuracy;
  final bool firstClear;
  const LessonResult({required this.xp, required this.accuracy, required this.firstClear});
}

class LessonCompleteScreen extends StatefulWidget {
  final int id;

  /// Set by LessonScreen right before navigating here — mirrors the web
  /// app's sessionStorage handoff (frontend/app/lesson/[id]/complete).
  static LessonResult? pendingResult;

  const LessonCompleteScreen({super.key, required this.id});

  @override
  State<LessonCompleteScreen> createState() => _LessonCompleteScreenState();
}

class _LessonCompleteScreenState extends State<LessonCompleteScreen> {
  late final ConfettiController _confetti;
  late final LessonResult _result;

  @override
  void initState() {
    super.initState();
    _result = LessonCompleteScreen.pendingResult ?? const LessonResult(xp: 10, accuracy: 100, firstClear: true);
    LessonCompleteScreen.pendingResult = null;
    _confetti = ConfettiController(duration: const Duration(seconds: 2))..play();
  }

  @override
  void dispose() {
    _confetti.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        decoration: const BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
            colors: [LumoraColors.purpleDark, LumoraColors.purple, LumoraColors.purpleDark],
          ),
        ),
        child: Stack(
          alignment: Alignment.topCenter,
          children: [
            Align(
              alignment: Alignment.topCenter,
              child: ConfettiWidget(
                confettiController: _confetti,
                blastDirectionality: BlastDirectionality.explosive,
                shouldLoop: false,
                numberOfParticles: 30,
                colors: const [LumoraColors.amber, LumoraColors.teal, LumoraColors.coral, Colors.white],
              ),
            ),
            SafeArea(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 24),
                child: Column(
                  children: [
                    Expanded(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          const FoxMascot(size: 150, glow: true, bounce: true),
                          const SizedBox(height: 24),
                          const Text('Lesson Complete!', style: TextStyle(color: Colors.white, fontSize: 26, fontWeight: FontWeight.w800)),
                          const SizedBox(height: 8),
                          Text(
                            _result.firstClear ? 'Brilliant work — Lumora is proud of you!' : 'Nicely done — practice makes fluent!',
                            textAlign: TextAlign.center,
                            style: const TextStyle(color: LumoraColors.purpleLight, fontSize: 16),
                          ),
                          const SizedBox(height: 32),
                          Row(children: [
                            Expanded(child: _StatCard(label: 'XP Earned', value: '+${_result.xp}', accent: LumoraColors.amber)),
                            const SizedBox(width: 16),
                            Expanded(child: _StatCard(label: 'Accuracy', value: '${_result.accuracy}%', accent: LumoraColors.teal)),
                          ]),
                        ],
                      ),
                    ),
                    LumoraButton(
                      label: 'Continue',
                      full: true,
                      variant: LumoraButtonVariant.secondary,
                      onPressed: () => context.go('/home'),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _StatCard extends StatelessWidget {
  final String label;
  final String value;
  final Color accent;
  const _StatCard({required this.label, required this.value, required this.accent});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 20),
      decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.1), borderRadius: BorderRadius.circular(LumoraRadii.xl)),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(value, style: TextStyle(color: accent, fontSize: 26, fontWeight: FontWeight.w800)),
          const SizedBox(height: 4),
          Text(label.toUpperCase(), style: const TextStyle(color: LumoraColors.purpleLight, fontSize: 10, fontWeight: FontWeight.w700, letterSpacing: 0.5)),
        ],
      ),
    );
  }
}
