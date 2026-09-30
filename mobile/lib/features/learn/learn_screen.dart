import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/languages.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../core/theme/shadows.dart';
import '../../models/lesson.dart';
import '../../models/listening_reading.dart';
import '../../providers/auth_provider.dart';
import '../../providers/learn_provider.dart';
import '../../widgets/language_switcher.dart';
import '../../widgets/lumora_button.dart';
import '../../widgets/skill_icon.dart';

class _Unit {
  final String name;
  final List<Skill> skills;
  _Unit(this.name) : skills = [];
}

class LearnScreen extends ConsumerStatefulWidget {
  const LearnScreen({super.key});

  @override
  ConsumerState<LearnScreen> createState() => _LearnScreenState();
}

class _LearnScreenState extends ConsumerState<LearnScreen> {
  bool _roadmap = false;

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(learnProvider);
    final user = ref.watch(authProvider).user;

    final totalLessons = state.skills.fold<int>(0, (n, s) => n + s.lessonCount);
    final doneLessons = state.skills.fold<int>(0, (n, s) => n + s.completedCount);
    final pct = totalLessons > 0 ? doneLessons / totalLessons : 0.0;

    final units = <_Unit>[];
    for (final s in state.skills) {
      final name = s.unit.isEmpty ? 'Course' : s.unit;
      var u = units.where((x) => x.name == name).firstOrNull;
      if (u == null) {
        u = _Unit(name);
        units.add(u);
      }
      u.skills.add(s);
    }

    return Scaffold(
      backgroundColor: LumoraColors.cream,
      body: CustomScrollView(
        slivers: [
          SliverToBoxAdapter(
            child: Container(
              color: LumoraColors.purple,
              padding: const EdgeInsets.fromLTRB(20, 56, 20, 20),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text('${languageName(user?.targetLanguage ?? "es")} COURSE'.toUpperCase(),
                                style: const TextStyle(color: Colors.white60, fontSize: 11, fontWeight: FontWeight.w800)),
                            const SizedBox(height: 2),
                            const Text('Your path', style: TextStyle(color: Colors.white, fontSize: 24, fontWeight: FontWeight.w800)),
                          ],
                        ),
                      ),
                      LanguageSwitcher(onChanged: () => ref.read(learnProvider.notifier).load()),
                    ],
                  ),
                  const SizedBox(height: 16),
                  Row(
                    children: [
                      Expanded(
                        child: ClipRRect(
                          borderRadius: BorderRadius.circular(LumoraRadii.full),
                          child: LinearProgressIndicator(
                            value: pct,
                            minHeight: 8,
                            backgroundColor: Colors.white24,
                            valueColor: const AlwaysStoppedAnimation(LumoraColors.amber),
                          ),
                        ),
                      ),
                      const SizedBox(width: 10),
                      Text('$doneLessons/${totalLessons == 0 ? "—" : totalLessons} lessons',
                          style: const TextStyle(color: Colors.white, fontWeight: FontWeight.w800, fontSize: 12)),
                    ],
                  ),
                  const SizedBox(height: 4),
                  Text('${user?.xp ?? 0} XP · ${user?.levelName.isNotEmpty == true ? user!.levelName : "Spark"}',
                      style: const TextStyle(color: Colors.white70, fontSize: 12)),
                  const SizedBox(height: 14),
                  Container(
                    padding: const EdgeInsets.all(4),
                    decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.15), borderRadius: BorderRadius.circular(LumoraRadii.full)),
                    child: Row(children: [
                      Expanded(child: _ViewTab(label: 'Course', icon: Icons.menu_book_rounded, active: !_roadmap, onTap: () => setState(() => _roadmap = false))),
                      Expanded(child: _ViewTab(label: 'Roadmap', icon: Icons.map_rounded, active: _roadmap, onTap: () => setState(() => _roadmap = true))),
                    ]),
                  ),
                ],
              ),
            ),
          ),
          // Only fall back to the error / spinner when there's nothing to show:
          // a reload keeps the current course on screen until the new one lands.
          if (state.error && state.skills.isEmpty)
            SliverFillRemaining(
              child: Center(
                child: Padding(
                  padding: const EdgeInsets.all(24),
                  child: Column(mainAxisSize: MainAxisSize.min, children: [
                    const Text("Couldn't load your course", style: TextStyle(fontWeight: FontWeight.w800)),
                    const SizedBox(height: 12),
                    LumoraButton(label: 'Retry', onPressed: () => ref.read(learnProvider.notifier).load()),
                  ]),
                ),
              ),
            )
          else if (state.loading && state.skills.isEmpty)
            const SliverFillRemaining(child: Center(child: CircularProgressIndicator()))
          else if (state.skills.isEmpty)
            const SliverFillRemaining(
              child: Center(child: Text('No lessons yet. Your course is being prepared.', style: TextStyle(color: LumoraColors.slatey))),
            )
          else if (_roadmap)
            SliverPadding(
              padding: const EdgeInsets.all(20),
              sliver: SliverToBoxAdapter(child: _RoadmapView(skills: state.skills)),
            )
          else
            SliverPadding(
              padding: const EdgeInsets.fromLTRB(16, 20, 16, 32),
              sliver: SliverList(
                delegate: SliverChildListDelegate([
                  for (final u in units) _UnitSection(unit: u, listening: state.listening, reading: state.reading),
                ]),
              ),
            ),
        ],
      ),
    );
  }
}

extension _FirstOrNull<T> on Iterable<T> {
  T? get firstOrNull => isEmpty ? null : first;
}

class _ViewTab extends StatelessWidget {
  final String label;
  final IconData icon;
  final bool active;
  final VoidCallback onTap;
  const _ViewTab({required this.label, required this.icon, required this.active, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return Material(
      color: active ? Colors.white : Colors.transparent,
      borderRadius: BorderRadius.circular(LumoraRadii.full),
      child: InkWell(
        borderRadius: BorderRadius.circular(LumoraRadii.full),
        onTap: onTap,
        child: Container(
          padding: const EdgeInsets.symmetric(vertical: 8),
          child: Row(mainAxisAlignment: MainAxisAlignment.center, children: [
            Icon(icon, size: 16, color: active ? LumoraColors.purple : Colors.white70),
            const SizedBox(width: 6),
            Text(label, style: TextStyle(color: active ? LumoraColors.purple : Colors.white70, fontWeight: FontWeight.w800, fontSize: 13)),
          ]),
        ),
      ),
    );
  }
}

class _UnitSection extends StatelessWidget {
  final _Unit unit;
  final List<ListeningSession> listening;
  final List<ReadingSession> reading;
  const _UnitSection({required this.unit, required this.listening, required this.reading});

  @override
  Widget build(BuildContext context) {
    final unitUnlocked = unit.skills.any((s) => s.unlocked);
    final unitComplete = unit.skills.isNotEmpty && unit.skills.every((s) => s.completed);

    return Padding(
      padding: const EdgeInsets.only(bottom: 28),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.only(left: 4, bottom: 10),
            child: Text(unit.name.toUpperCase(),
                style: const TextStyle(fontSize: 12, fontWeight: FontWeight.w800, color: LumoraColors.gray500, letterSpacing: 0.5)),
          ),
          for (final s in unit.skills) Padding(padding: const EdgeInsets.only(bottom: 10), child: _SkillCard(skill: s)),
          for (final r in reading.where((r) => r.unit == unit.name))
            Padding(padding: const EdgeInsets.only(bottom: 10), child: _ReadingCard(session: r, locked: !unitUnlocked, unit: unit.name)),
          for (final l in listening.where((l) => l.unit == unit.name))
            Padding(padding: const EdgeInsets.only(bottom: 10), child: _ListeningCard(session: l, locked: !unitUnlocked, unit: unit.name)),
          if (unitComplete) ...[
            _ReviewMistakesCard(),
            const SizedBox(height: 10),
            _DoExamCard(),
          ],
        ],
      ),
    );
  }
}

Color _hex(String s) {
  var h = s.replaceAll('#', '');
  if (h.length == 6) h = 'FF$h';
  final v = int.tryParse(h, radix: 16);
  return v == null ? LumoraColors.purple : Color(v);
}

class _SkillCard extends StatelessWidget {
  final Skill skill;
  const _SkillCard({required this.skill});

  @override
  Widget build(BuildContext context) {
    final color = _hex(skill.color);
    if (!skill.unlocked) {
      return Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), border: Border.all(color: LumoraColors.gray100), boxShadow: LumoraShadows.card),
        child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Row(children: [
            Container(width: 48, height: 48, decoration: BoxDecoration(color: LumoraColors.gray100, borderRadius: BorderRadius.circular(LumoraRadii.md)),
                child: const Icon(Icons.lock_rounded, color: LumoraColors.gray500)),
            const SizedBox(width: 12),
            Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Text(skill.title, style: const TextStyle(fontWeight: FontWeight.w800)),
              Text(skill.description, maxLines: 1, overflow: TextOverflow.ellipsis, style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
            ])),
          ]),
          const SizedBox(height: 10),
          Container(
            width: double.infinity,
            padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 8),
            decoration: BoxDecoration(color: LumoraColors.gray50, borderRadius: BorderRadius.circular(LumoraRadii.sm)),
            child: const Text('Complete the previous skill to unlock this one', style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600, color: LumoraColors.slatey)),
          ),
        ]),
      );
    }

    final pct = skill.lessonCount > 0 ? skill.completedCount / skill.lessonCount : 0.0;

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            Container(width: 48, height: 48, decoration: BoxDecoration(color: color, borderRadius: BorderRadius.circular(LumoraRadii.md)),
                child: Center(child: SkillIcon(name: skill.icon, size: 24, color: Colors.white))),
            const SizedBox(width: 12),
            Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Text(skill.title, style: const TextStyle(fontWeight: FontWeight.w800)),
              Text(skill.description, maxLines: 1, overflow: TextOverflow.ellipsis, style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
            ])),
            if (skill.completed)
              Container(width: 28, height: 28, decoration: const BoxDecoration(color: LumoraColors.teal, shape: BoxShape.circle),
                  child: const Icon(Icons.check, size: 16, color: Colors.white)),
          ]),
          const SizedBox(height: 10),
          Row(children: [
            Expanded(
              child: ClipRRect(
                borderRadius: BorderRadius.circular(LumoraRadii.full),
                child: LinearProgressIndicator(value: pct, minHeight: 6, backgroundColor: LumoraColors.gray100, valueColor: AlwaysStoppedAnimation(color)),
              ),
            ),
            const SizedBox(width: 8),
            Text('${skill.completedCount}/${skill.lessonCount}', style: const TextStyle(fontSize: 12, fontWeight: FontWeight.w800, color: LumoraColors.slatey)),
          ]),
          const Divider(height: 20),
          for (var i = 0; i < skill.lessons.length; i++) _LessonRow(lesson: skill.lessons[i], index: i, completedCount: skill.completedCount, color: color),
        ],
      ),
    );
  }
}

class _LessonRow extends StatelessWidget {
  final Lesson lesson;
  final int index;
  final int completedCount;
  final Color color;
  const _LessonRow({required this.lesson, required this.index, required this.completedCount, required this.color});

  @override
  Widget build(BuildContext context) {
    final done = index < completedCount;
    final current = index == completedCount;
    return InkWell(
      onTap: () => context.push('/lesson/${lesson.id}'),
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 8),
        child: Row(children: [
          Container(
            width: 32, height: 32,
            decoration: BoxDecoration(
              color: done ? LumoraColors.teal : (current ? color : Colors.transparent),
              shape: BoxShape.circle,
              border: (!done && !current) ? Border.all(color: LumoraColors.gray100, width: 2) : null,
            ),
            child: Center(
              child: done
                  ? const Icon(Icons.check, size: 16, color: Colors.white)
                  : current
                      ? const Icon(Icons.play_arrow_rounded, size: 16, color: Colors.white)
                      : Text('${index + 1}', style: const TextStyle(fontSize: 12, fontWeight: FontWeight.w800, color: LumoraColors.gray300)),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Text(lesson.title, maxLines: 1, overflow: TextOverflow.ellipsis, style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 14)),
              Text('+${lesson.xpReward} XP${current ? " · Start now" : done ? " · Completed" : ""}',
                  style: const TextStyle(color: LumoraColors.slatey, fontSize: 11)),
            ]),
          ),
          const Icon(Icons.chevron_right, color: LumoraColors.gray300, size: 18),
        ]),
      ),
    );
  }
}

class _ReadingCard extends StatelessWidget {
  final ReadingSession session;
  final bool locked;
  final String unit;
  const _ReadingCard({required this.session, required this.locked, required this.unit});

  @override
  Widget build(BuildContext context) {
    if (locked) return _LockedCard(title: session.title, kind: 'Reading', unit: unit);
    return Material(
      color: LumoraColors.tealLight.withValues(alpha: 0.5),
      borderRadius: BorderRadius.circular(LumoraRadii.xl),
      child: InkWell(
        borderRadius: BorderRadius.circular(LumoraRadii.xl),
        onTap: () => context.push('/reading/${session.id}'),
        child: Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(borderRadius: BorderRadius.circular(LumoraRadii.xl), border: Border.all(color: LumoraColors.teal.withValues(alpha: 0.2))),
          child: Row(children: [
            Container(width: 48, height: 48, decoration: BoxDecoration(color: LumoraColors.teal, borderRadius: BorderRadius.circular(LumoraRadii.md)),
                child: const Icon(Icons.menu_book_rounded, color: Colors.white)),
            const SizedBox(width: 12),
            Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Row(children: [
                Flexible(child: Text(session.title, overflow: TextOverflow.ellipsis, style: const TextStyle(fontWeight: FontWeight.w800))),
                const SizedBox(width: 6),
                _Badge('Reading', LumoraColors.teal),
              ]),
              Text(session.description, maxLines: 1, overflow: TextOverflow.ellipsis, style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
            ])),
            _Badge('+${session.xpReward} XP', LumoraColors.amber),
          ]),
        ),
      ),
    );
  }
}

class _ListeningCard extends StatelessWidget {
  final ListeningSession session;
  final bool locked;
  final String unit;
  const _ListeningCard({required this.session, required this.locked, required this.unit});

  @override
  Widget build(BuildContext context) {
    if (locked) return _LockedCard(title: session.title, kind: 'Listening', unit: unit);
    return Material(
      color: LumoraColors.purpleLight.withValues(alpha: 0.5),
      borderRadius: BorderRadius.circular(LumoraRadii.xl),
      child: InkWell(
        borderRadius: BorderRadius.circular(LumoraRadii.xl),
        onTap: () => context.push('/listening/${session.id}'),
        child: Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(borderRadius: BorderRadius.circular(LumoraRadii.xl), border: Border.all(color: LumoraColors.purple.withValues(alpha: 0.15))),
          child: Row(children: [
            Container(width: 48, height: 48, decoration: BoxDecoration(color: LumoraColors.purple, borderRadius: BorderRadius.circular(LumoraRadii.md)),
                child: const Icon(Icons.headphones_rounded, color: Colors.white)),
            const SizedBox(width: 12),
            Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Row(children: [
                Flexible(child: Text(session.title, overflow: TextOverflow.ellipsis, style: const TextStyle(fontWeight: FontWeight.w800))),
                const SizedBox(width: 6),
                _Badge('Listening', LumoraColors.purple),
              ]),
              Text(session.description, maxLines: 1, overflow: TextOverflow.ellipsis, style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
            ])),
            _Badge('+${session.xpReward} XP', LumoraColors.amber),
          ]),
        ),
      ),
    );
  }
}

class _LockedCard extends StatelessWidget {
  final String title;
  final String kind;
  final String unit;
  const _LockedCard({required this.title, required this.kind, required this.unit});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), border: Border.all(color: LumoraColors.gray100)),
      child: Row(children: [
        Container(width: 48, height: 48, decoration: BoxDecoration(color: LumoraColors.gray100, borderRadius: BorderRadius.circular(LumoraRadii.md)),
            child: const Icon(Icons.lock_rounded, color: LumoraColors.gray500)),
        const SizedBox(width: 12),
        Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Row(children: [
            Flexible(child: Text(title, overflow: TextOverflow.ellipsis, style: const TextStyle(fontWeight: FontWeight.w800))),
            const SizedBox(width: 6),
            _Badge(kind, LumoraColors.gray500),
          ]),
          Text('Reach the $unit unit to unlock', style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
        ])),
      ]),
    );
  }
}

class _Badge extends StatelessWidget {
  final String text;
  final Color color;
  const _Badge(this.text, this.color);
  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.full)),
      child: Text(text, style: TextStyle(color: color, fontSize: 10, fontWeight: FontWeight.w800)),
    );
  }
}

class _ReviewMistakesCard extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Material(
      color: LumoraColors.amberLight.withValues(alpha: 0.6),
      borderRadius: BorderRadius.circular(LumoraRadii.xl),
      child: InkWell(
        borderRadius: BorderRadius.circular(LumoraRadii.xl),
        onTap: () => context.push('/practice/run?mode=mistakes'),
        child: Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(borderRadius: BorderRadius.circular(LumoraRadii.xl), border: Border.all(color: LumoraColors.amber.withValues(alpha: 0.3))),
          child: Row(children: [
            Container(width: 48, height: 48, decoration: BoxDecoration(color: LumoraColors.amber, borderRadius: BorderRadius.circular(LumoraRadii.md)),
                child: const Icon(Icons.refresh_rounded, color: LumoraColors.ink)),
            const SizedBox(width: 12),
            Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Row(children: [
                const Flexible(child: Text('Review your mistakes', style: TextStyle(fontWeight: FontWeight.w800))),
                const SizedBox(width: 6),
                _Badge('Unit done', LumoraColors.amber),
              ]),
              const Text('Turn the words you missed into a quick refresher.', style: TextStyle(color: LumoraColors.slatey, fontSize: 12)),
            ])),
            const Icon(Icons.chevron_right, color: LumoraColors.gray300),
          ]),
        ),
      ),
    );
  }
}

class _DoExamCard extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Material(
      color: LumoraColors.purple,
      borderRadius: BorderRadius.circular(LumoraRadii.xl),
      child: InkWell(
        borderRadius: BorderRadius.circular(LumoraRadii.xl),
        onTap: () => context.push('/exam'),
        child: Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.float),
          child: Row(children: [
            Container(width: 48, height: 48, decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.15), borderRadius: BorderRadius.circular(LumoraRadii.md)),
                child: const Icon(Icons.school_rounded, color: Colors.white)),
            const SizedBox(width: 12),
            const Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Text('Ready? Take the exam', style: TextStyle(color: Colors.white, fontWeight: FontWeight.w800)),
              Text('Test all four skills and earn your certificate.', style: TextStyle(color: Colors.white70, fontSize: 12)),
            ])),
            const Icon(Icons.chevron_right, color: Colors.white70),
          ]),
        ),
      ),
    );
  }
}

/// A simplified "galaxy map": skills laid out on a winding path over a
/// starfield backdrop (frontend/components/RoadmapView.tsx).
class _RoadmapView extends StatelessWidget {
  final List<Skill> skills;
  const _RoadmapView({required this.skills});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 24, horizontal: 12),
      decoration: BoxDecoration(color: LumoraColors.space, borderRadius: BorderRadius.circular(LumoraRadii.xl2)),
      child: Column(
        children: [
          for (var i = 0; i < skills.length; i++)
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 10),
              child: Align(
                alignment: i.isEven ? Alignment.centerLeft : Alignment.centerRight,
                child: Padding(
                  padding: EdgeInsets.only(left: i.isEven ? 16 : 0, right: i.isEven ? 0 : 16),
                  child: GestureDetector(
                    onTap: skills[i].unlocked && skills[i].lessons.isNotEmpty
                        ? () => context.push('/lesson/${skills[i].lessons.first.id}')
                        : null,
                    child: Column(children: [
                      Container(
                        width: 64, height: 64,
                        decoration: BoxDecoration(
                          shape: BoxShape.circle,
                          color: skills[i].unlocked ? _hex(skills[i].color) : Colors.white12,
                          boxShadow: skills[i].unlocked ? LumoraShadows.skill : null,
                        ),
                        child: Center(
                          child: skills[i].unlocked
                              ? SkillIcon(name: skills[i].icon, size: 28, color: Colors.white)
                              : const Icon(Icons.lock_rounded, color: Colors.white38),
                        ),
                      ),
                      const SizedBox(height: 6),
                      Text(skills[i].title, style: const TextStyle(color: Colors.white, fontSize: 11, fontWeight: FontWeight.w700)),
                    ]),
                  ),
                ),
              ),
            ),
        ],
      ),
    );
  }
}
