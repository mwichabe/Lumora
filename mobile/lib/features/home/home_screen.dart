import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../core/theme/shadows.dart';
import '../../models/home_data.dart';
import '../../providers/home_provider.dart';
import '../../providers/learn_provider.dart';
import '../../widgets/avatar.dart';
import '../../widgets/fox_mascot.dart';
import '../../widgets/header_bells.dart';
import '../../widgets/lumora_button.dart';
import '../../widgets/skill_icon.dart';
import '../../widgets/stat_widgets.dart';

String _greeting() {
  final h = DateTime.now().hour;
  if (h < 12) return 'Good morning';
  if (h < 18) return 'Good afternoon';
  return 'Late night study sesh?';
}

const _steps = [
  (icon: Icons.menu_book_rounded, title: 'Take a lesson', desc: 'Short, playful exercises'),
  (icon: Icons.bolt_rounded, title: 'Earn XP', desc: 'Hit your daily goal'),
  (icon: Icons.local_fire_department_rounded, title: 'Build a streak', desc: 'Practise a little daily'),
];

class HomeScreen extends ConsumerWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final homeAsync = ref.watch(homeProvider);

    return Scaffold(
      backgroundColor: LumoraColors.cream,
      body: RefreshIndicator(
        onRefresh: () => ref.read(homeProvider.notifier).refresh(),
        child: CustomScrollView(
          slivers: [
            SliverToBoxAdapter(child: _Header(data: homeAsync.valueOrNull)),
            SliverPadding(
              padding: const EdgeInsets.fromLTRB(20, 20, 20, 32),
              sliver: SliverList(
                delegate: SliverChildListDelegate([
                  homeAsync.when(
                    data: (data) => _StartHero(data: data),
                    loading: () => const _HeroSkeleton(),
                    error: (_, _) => _ErrorCard(onRetry: () => ref.read(homeProvider.notifier).refresh()),
                  ),
                  if (homeAsync.hasValue && homeAsync.value!.user.xp == 0) ...[
                    const SizedBox(height: 24),
                    const _FirstTimerGuide(),
                  ],
                  const SizedBox(height: 24),
                  if (homeAsync.hasValue) _DailyGoalCard(user: homeAsync.value!.user),
                  const SizedBox(height: 24),
                  _RoadmapEntry(),
                  const SizedBox(height: 24),
                  const _SectionLabel('Daily quests'),
                  const SizedBox(height: 8),
                  Row(
                    children: const [
                      Text('🦔', style: TextStyle(fontSize: 22)),
                      SizedBox(width: 8),
                      Expanded(
                        child: Text('Pip: "Your quests await!! Let\'s GO!!"',
                            style: TextStyle(color: LumoraColors.slatey, fontStyle: FontStyle.italic, fontSize: 13)),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  homeAsync.when(
                    data: (data) => _QuestList(data: data),
                    loading: () => const SizedBox(
                        height: 64, child: Center(child: CircularProgressIndicator())),
                    error: (_, _) => const SizedBox.shrink(),
                  ),
                ]),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _Header extends StatelessWidget {
  final HomeData? data;
  const _Header({required this.data});

  @override
  Widget build(BuildContext context) {
    final user = data?.user;
    return Container(
      color: LumoraColors.purple,
      padding: const EdgeInsets.fromLTRB(20, 56, 20, 20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                decoration: BoxDecoration(shape: BoxShape.circle, border: Border.all(color: Colors.white54, width: 2)),
                child: LumoraAvatar(
                  name: user?.name ?? '',
                  avatarColor: user?.avatarColor ?? '#6C3FC5',
                  avatarUrl: user?.avatarUrl,
                  size: 40,
                ),
              ),
              const SizedBox(width: 10),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(_greeting(), style: const TextStyle(color: Colors.white70, fontSize: 12)),
                    Text(user?.name.isNotEmpty == true ? user!.name : 'friend',
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                        style: const TextStyle(color: Colors.white, fontWeight: FontWeight.w800, fontSize: 16)),
                  ],
                ),
              ),
              const IdeasBell(),
              const SizedBox(width: 6),
              const ChatBell(),
              const SizedBox(width: 6),
              const NotificationBell(),
            ],
          ),
          const SizedBox(height: 12),
          Row(children: [
            GemCounter(gems: user?.gems ?? 0),
            const SizedBox(width: 8),
            StreakFlame(streak: user?.streak ?? 0, light: true),
          ]),
          if ((user?.streak ?? 0) > 0) ...[
            const SizedBox(height: 14),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
              decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.15), borderRadius: BorderRadius.circular(LumoraRadii.lg)),
              child: Row(
                children: [
                  const Text('🔥', style: TextStyle(fontSize: 20)),
                  const SizedBox(width: 8),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text('${user!.streak} day streak!', style: const TextStyle(color: Colors.white, fontWeight: FontWeight.w800)),
                      const Text('Keep it going!', style: TextStyle(color: Colors.white70, fontSize: 12)),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ],
      ),
    );
  }
}

class _StartHero extends StatelessWidget {
  final HomeData data;
  const _StartHero({required this.data});

  @override
  Widget build(BuildContext context) {
    final lesson = data.nextLesson;
    final skill = data.nextSkill;
    final isNewUser = data.user.xp == 0;

    if (lesson == null) {
      return Container(
        padding: const EdgeInsets.all(20),
        decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl2), boxShadow: LumoraShadows.cardLg),
        child: Column(
          children: [
            const FoxMascot(size: 88, glow: true),
            const SizedBox(height: 12),
            const Text("You're all caught up! 🎉", textAlign: TextAlign.center, style: TextStyle(fontSize: 18, fontWeight: FontWeight.w800)),
            const SizedBox(height: 4),
            const Text('Great work. Unlock new skills on your learning path.',
                textAlign: TextAlign.center, style: TextStyle(color: LumoraColors.slatey)),
            const SizedBox(height: 16),
            Consumer(
              builder: (context, ref, _) => LumoraButton(
                label: 'Explore the map',
                onPressed: () {
                  ref.read(learnViewProvider.notifier).state = LearnView.roadmap;
                  context.go('/learn');
                },
              ),
            ),
          ],
        ),
      );
    }

    final label = isNewUser ? 'START HERE' : 'CONTINUE WHERE YOU LEFT OFF';
    final cta = isNewUser ? 'Start your first lesson' : 'Continue lesson';

    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(color: LumoraColors.purple, borderRadius: BorderRadius.circular(LumoraRadii.xl2), boxShadow: LumoraShadows.float),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(label, style: const TextStyle(color: Colors.white70, fontSize: 11, fontWeight: FontWeight.w800, letterSpacing: 0.5)),
          const SizedBox(height: 12),
          Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Container(
                width: 64, height: 64,
                decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.15), borderRadius: BorderRadius.circular(LumoraRadii.lg)),
                child: Center(child: SkillIcon(name: skill?.icon, size: 28, color: Colors.white)),
              ),
              const SizedBox(width: 14),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(lesson.title, maxLines: 1, overflow: TextOverflow.ellipsis,
                        style: const TextStyle(color: Colors.white, fontSize: 20, fontWeight: FontWeight.w800)),
                    Text(skill?.title ?? 'Your next lesson', maxLines: 1, overflow: TextOverflow.ellipsis,
                        style: const TextStyle(color: Colors.white70)),
                    const SizedBox(height: 6),
                    Wrap(spacing: 8, runSpacing: 4, children: [
                      _Chip('Lesson ${lesson.orderIndex}', Colors.white.withValues(alpha: 0.15), Colors.white),
                      _Chip('+${lesson.xpReward} XP', LumoraColors.amber, LumoraColors.ink),
                      _Chip('~5 min', Colors.white.withValues(alpha: 0.15), Colors.white),
                    ]),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 18),
          SizedBox(
            width: double.infinity,
            child: ElevatedButton.icon(
              onPressed: () => context.push('/lesson/${lesson.id}'),
              style: ElevatedButton.styleFrom(backgroundColor: Colors.white, foregroundColor: LumoraColors.purple, minimumSize: const Size.fromHeight(48)),
              icon: const Icon(Icons.play_arrow_rounded),
              label: Text(cta, style: const TextStyle(fontWeight: FontWeight.w800)),
            ),
          ),
        ],
      ),
    );
  }
}

class _Chip extends StatelessWidget {
  final String label;
  final Color bg;
  final Color fg;
  const _Chip(this.label, this.bg, this.fg);
  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 3),
      decoration: BoxDecoration(color: bg, borderRadius: BorderRadius.circular(LumoraRadii.full)),
      child: Text(label, style: TextStyle(color: fg, fontSize: 11, fontWeight: FontWeight.w800)),
    );
  }
}

class _HeroSkeleton extends StatelessWidget {
  const _HeroSkeleton();
  @override
  Widget build(BuildContext context) {
    return Container(
      height: 180,
      decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl2), boxShadow: LumoraShadows.cardLg),
      child: const Center(child: CircularProgressIndicator()),
    );
  }
}

class _ErrorCard extends StatelessWidget {
  final VoidCallback onRetry;
  const _ErrorCard({required this.onRetry});
  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl2), boxShadow: LumoraShadows.cardLg),
      child: Column(children: [
        const Text("Couldn't load your dashboard", style: TextStyle(fontWeight: FontWeight.w800, fontSize: 16)),
        const SizedBox(height: 4),
        const Text('Check your connection and try again.', style: TextStyle(color: LumoraColors.slatey)),
        const SizedBox(height: 12),
        LumoraButton(label: 'Retry', onPressed: onRetry),
      ]),
    );
  }
}

class _FirstTimerGuide extends StatelessWidget {
  const _FirstTimerGuide();
  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const _SectionLabel('How Lumora works'),
        const SizedBox(height: 8),
        for (var i = 0; i < _steps.length; i++)
          Padding(
            padding: const EdgeInsets.only(bottom: 8),
            child: Container(
              padding: const EdgeInsets.all(14),
              decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.lg), boxShadow: LumoraShadows.card),
              child: Row(children: [
                Stack(clipBehavior: Clip.none, children: [
                  Container(
                    width: 40, height: 40,
                    decoration: BoxDecoration(color: LumoraColors.purpleLight, borderRadius: BorderRadius.circular(LumoraRadii.md)),
                    child: Icon(_steps[i].icon, color: LumoraColors.purple, size: 20),
                  ),
                  Positioned(
                    top: -4, left: -4,
                    child: Container(
                      width: 18, height: 18,
                      decoration: const BoxDecoration(color: LumoraColors.purple, shape: BoxShape.circle),
                      child: Center(child: Text('${i + 1}', style: const TextStyle(color: Colors.white, fontSize: 10, fontWeight: FontWeight.w800))),
                    ),
                  ),
                ]),
                const SizedBox(width: 12),
                Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                  Text(_steps[i].title, style: const TextStyle(fontWeight: FontWeight.w800)),
                  Text(_steps[i].desc, style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
                ]),
              ]),
            ),
          ),
      ],
    );
  }
}

class _DailyGoalCard extends StatelessWidget {
  final dynamic user;
  const _DailyGoalCard({required this.user});
  @override
  Widget build(BuildContext context) {
    final reached = user.xpToday >= user.dailyGoalXp;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const _SectionLabel('Your daily goal'),
        const SizedBox(height: 8),
        Container(
          padding: const EdgeInsets.all(18),
          decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.lg), boxShadow: LumoraShadows.card),
          child: Column(children: [
            XPBar(value: user.xpToday, max: user.dailyGoalXp),
            const SizedBox(height: 8),
            Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
              Text('${user.xpToday}/${user.dailyGoalXp} XP today', style: const TextStyle(fontWeight: FontWeight.w800)),
              Text(reached ? 'Goal reached! 🎉' : '${(user.dailyGoalXp - user.xpToday).clamp(0, 999)} XP to go',
                  style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
            ]),
          ]),
        ),
      ],
    );
  }
}

class _RoadmapEntry extends ConsumerWidget {
  const _RoadmapEntry();
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const _SectionLabel('Your learning path'),
        const SizedBox(height: 8),
        Material(
          color: Colors.white,
          borderRadius: BorderRadius.circular(LumoraRadii.lg),
          child: InkWell(
            borderRadius: BorderRadius.circular(LumoraRadii.lg),
            onTap: () {
              // Open the Learn tab on its Roadmap view, not the course list.
              ref.read(learnViewProvider.notifier).state = LearnView.roadmap;
              context.go('/learn');
            },
            child: Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.lg), boxShadow: LumoraShadows.card),
              child: Row(children: [
                Container(
                  width: 48, height: 48,
                  decoration: BoxDecoration(color: LumoraColors.tealLight, borderRadius: BorderRadius.circular(LumoraRadii.md)),
                  child: const Icon(Icons.map_rounded, color: LumoraColors.teal),
                ),
                const SizedBox(width: 12),
                const Expanded(
                  child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                    Text('Explore the galaxy map', style: TextStyle(fontWeight: FontWeight.w800)),
                    Text('See every skill and what to learn next.', style: TextStyle(color: LumoraColors.slatey, fontSize: 12)),
                  ]),
                ),
                const Icon(Icons.chevron_right, color: LumoraColors.gray300),
              ]),
            ),
          ),
        ),
      ],
    );
  }
}

class _QuestList extends StatelessWidget {
  final HomeData data;
  const _QuestList({required this.data});

  @override
  Widget build(BuildContext context) {
    if (data.quests.isEmpty) {
      return Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.md), boxShadow: LumoraShadows.quest),
        child: const Text("Complete a lesson to unlock today's quests.", style: TextStyle(color: LumoraColors.slatey)),
      );
    }
    return Column(
      children: [
        for (final q in data.quests)
          Padding(
            key: ValueKey(q.questId),
            padding: const EdgeInsets.only(bottom: 8),
            child: Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.md), boxShadow: LumoraShadows.quest),
              child: Row(children: [
                Text(q.quest?.icon ?? '⭐', style: const TextStyle(fontSize: 20)),
                const SizedBox(width: 10),
                Expanded(
                  child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                    Text(q.quest?.title ?? '', style: const TextStyle(fontWeight: FontWeight.w700)),
                    const SizedBox(height: 4),
                    ClipRRect(
                      borderRadius: BorderRadius.circular(LumoraRadii.full),
                      // Animates from the old value when a refresh lands, so
                      // progress made in a lesson visibly fills the bar.
                      child: TweenAnimationBuilder<double>(
                        tween: Tween(end: (q.progress / (q.quest?.target ?? 1)).clamp(0, 1).toDouble()),
                        duration: const Duration(milliseconds: 600),
                        curve: Curves.easeOutCubic,
                        builder: (context, value, _) => LinearProgressIndicator(
                          value: value,
                          minHeight: 6,
                          backgroundColor: LumoraColors.gray100,
                          valueColor: const AlwaysStoppedAnimation(LumoraColors.amber),
                        ),
                      ),
                    ),
                  ]),
                ),
                const SizedBox(width: 8),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                  decoration: BoxDecoration(color: LumoraColors.amber.withValues(alpha: 0.15), borderRadius: BorderRadius.circular(LumoraRadii.full)),
                  child: Text('+${q.quest?.xpReward ?? 0}', style: const TextStyle(color: LumoraColors.amber, fontSize: 11, fontWeight: FontWeight.w800)),
                ),
                const SizedBox(width: 8),
                Container(
                  width: 24, height: 24,
                  decoration: BoxDecoration(
                    color: q.completed ? LumoraColors.teal : Colors.transparent,
                    shape: BoxShape.circle,
                    border: q.completed ? null : Border.all(color: LumoraColors.gray100, width: 2),
                  ),
                  child: q.completed ? const Icon(Icons.check, size: 14, color: Colors.white) : null,
                ),
              ]),
            ),
          ),
      ],
    );
  }
}

class _SectionLabel extends StatelessWidget {
  final String text;
  const _SectionLabel(this.text);
  @override
  Widget build(BuildContext context) {
    return Text(text.toUpperCase(),
        style: const TextStyle(fontSize: 11, fontWeight: FontWeight.w700, color: LumoraColors.gray500, letterSpacing: 0.5));
  }
}
