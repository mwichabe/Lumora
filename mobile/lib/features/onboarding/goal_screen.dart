import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/network/api_client.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../providers/auth_provider.dart';
import '../../widgets/fox_mascot.dart';
import '../../widgets/lumora_button.dart';
import '../../widgets/stat_widgets.dart';
import 'language_screen.dart';

class _GoalOption {
  final String label;
  final String minutes;
  final int xp;
  final String reason;
  const _GoalOption(this.label, this.minutes, this.xp, this.reason);
}

const _goals = [
  _GoalOption('Casual', '5 min/day', 10, 'casual'),
  _GoalOption('Regular', '10 min/day', 20, 'regular'),
  _GoalOption('Serious', '15 min/day', 30, 'serious'),
  _GoalOption('Intense', '20+ min/day', 50, 'intense'),
];

/// frontend/app/onboarding/goal — picks a daily XP goal, then calls
/// POST /api/auth/setup to finish onboarding.
class OnboardingGoalScreen extends ConsumerStatefulWidget {
  const OnboardingGoalScreen({super.key});

  @override
  ConsumerState<OnboardingGoalScreen> createState() => _OnboardingGoalScreenState();
}

class _OnboardingGoalScreenState extends ConsumerState<OnboardingGoalScreen> {
  int _selected = 1;
  bool _loading = false;

  Future<void> _finish() async {
    final language = OnboardingLanguageScreen.pendingChoice;
    if (language == null) {
      if (mounted) context.go('/onboarding/language');
      return;
    }
    setState(() => _loading = true);
    try {
      final goal = _goals[_selected];
      final user = await ApiClient.instance.setup(language, goal.xp, goal.reason);
      ref.read(authProvider.notifier).setUser(user);
      if (mounted) context.go('/home');
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: LumoraColors.cream,
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              const SizedBox(height: 8),
              const Center(child: FoxMascot(size: 88)),
              const SizedBox(height: 12),
              const Center(
                child: SpeechBubble(
                  child: Text('How much time can you commit each day?', textAlign: TextAlign.center),
                ),
              ),
              const SizedBox(height: 28),
              const Text('Set your daily goal', style: TextStyle(fontSize: 22, fontWeight: FontWeight.w800, color: LumoraColors.ink)),
              const SizedBox(height: 16),
              for (var i = 0; i < _goals.length; i++) ...[
                _GoalCard(
                  goal: _goals[i],
                  selected: _selected == i,
                  onTap: () => setState(() => _selected = i),
                ),
                const SizedBox(height: 10),
              ],
              const SizedBox(height: 8),
              LumoraButton(label: 'Continue', full: true, loading: _loading, onPressed: _loading ? null : _finish),
            ],
          ),
        ),
      ),
    );
  }
}

class _GoalCard extends StatelessWidget {
  final _GoalOption goal;
  final bool selected;
  final VoidCallback onTap;
  const _GoalCard({required this.goal, required this.selected, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return Material(
      color: selected ? LumoraColors.purpleLight : Colors.white,
      borderRadius: BorderRadius.circular(LumoraRadii.lg),
      child: InkWell(
        borderRadius: BorderRadius.circular(LumoraRadii.lg),
        onTap: onTap,
        child: Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            border: Border.all(color: selected ? LumoraColors.purple : LumoraColors.gray100, width: selected ? 2 : 1),
            borderRadius: BorderRadius.circular(LumoraRadii.lg),
          ),
          child: Row(
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(goal.label, style: const TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.ink)),
                    Text(goal.minutes, style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
                  ],
                ),
              ),
              Text('${goal.xp} XP', style: const TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.amber)),
              const SizedBox(width: 8),
              Icon(selected ? Icons.check_circle_rounded : Icons.circle_outlined,
                  color: selected ? LumoraColors.purple : LumoraColors.gray300),
            ],
          ),
        ),
      ),
    );
  }
}
