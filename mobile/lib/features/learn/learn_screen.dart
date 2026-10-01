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
              padding: const EdgeInsets.fromLTRB(20, 8, 20, 32),
              // A lazy list: only the units near the viewport are built, which
              // matters for the long courses.
              sliver: SliverList.builder(
                itemCount: units.length + 1,
                itemBuilder: (context, i) {
                  if (i == units.length) return const _RoadmapFinish();
                  var before = 0;
                  for (var j = 0; j < i; j++) {
                    before += units[j].skills.length;
                  }
                  return _RoadmapUnit(unit: units[i], startIndex: before);
                },
              ),
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
          decoration: BoxDecoration(color: LumoraColors.purple, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.float),
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

/// One unit of the course "journey" (frontend/components/RoadmapView.tsx): a
/// heading, then the unit's skills as nodes zig-zagging down a dashed path,
/// coloured by status — completed, current or locked.
///
/// Every node is centred and then nudged left or right, and its label is
/// width-limited and wraps, so nothing can run off the edge of a phone however
/// long a skill's title is.
class _RoadmapUnit extends StatelessWidget {
  final _Unit unit;

  /// How many nodes come before this unit, so the zig-zag carries on across
  /// unit boundaries instead of restarting on the same side.
  final int startIndex;
  const _RoadmapUnit({required this.unit, required this.startIndex});

  static const _labelWidth = 140.0;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 12),
            child: Row(children: [
              const Expanded(child: Divider(color: LumoraColors.gray100, thickness: 1)),
              const SizedBox(width: 12),
              Flexible(
                flex: 4,
                child: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 5),
                  decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.full), boxShadow: LumoraShadows.card),
                  child: Text(
                    unit.name.toUpperCase(),
                    textAlign: TextAlign.center,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                    style: const TextStyle(color: LumoraColors.purple, fontSize: 11, fontWeight: FontWeight.w800, letterSpacing: 0.6),
                  ),
                ),
              ),
              const SizedBox(width: 12),
              const Expanded(child: Divider(color: LumoraColors.gray100, thickness: 1)),
            ]),
          ),
          LayoutBuilder(builder: (context, constraints) {
            // How far a node may swing from the centre line while its label
            // still fits inside the available width.
            final swing = ((constraints.maxWidth - _labelWidth) / 2).clamp(0.0, 56.0);
            return Stack(
              alignment: Alignment.topCenter,
              children: [
                const Positioned(top: 24, bottom: 24, child: _DashedLine()),
                Column(
                  children: [
                    for (var i = 0; i < unit.skills.length; i++)
                      Padding(
                        padding: EdgeInsets.only(bottom: i == unit.skills.length - 1 ? 0 : 24),
                        child: Transform.translate(
                          offset: Offset((startIndex + i).isEven ? -swing : swing, 0),
                          child: _RoadmapNode(skill: unit.skills[i]),
                        ),
                      ),
                  ],
                ),
              ],
            );
          }),
        ],
      ),
    );
  }
}

class _RoadmapNode extends StatelessWidget {
  final Skill skill;
  const _RoadmapNode({required this.skill});

  @override
  Widget build(BuildContext context) {
    final locked = !skill.unlocked;
    final current = !locked && !skill.completed;
    return SizedBox(
      width: _RoadmapUnit._labelWidth,
      child: Column(
        children: [
          GestureDetector(
            onTap: !locked && skill.lessons.isNotEmpty ? () => context.push('/lesson/${skill.lessons.first.id}') : null,
            child: Stack(
              clipBehavior: Clip.none,
              children: [
                Container(
                  width: 80,
                  height: 80,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    color: locked ? const Color(0xFFE5E5EC) : _hex(skill.color),
                    boxShadow: locked ? null : LumoraShadows.skill,
                    // The skill to do next gets an amber ring.
                    border: current ? Border.all(color: LumoraColors.amber.withValues(alpha: 0.7), width: 4) : null,
                  ),
                  child: Center(
                    child: locked
                        ? const Icon(Icons.lock_rounded, color: LumoraColors.gray500, size: 26)
                        : SkillIcon(name: skill.icon, size: 28, color: Colors.white),
                  ),
                ),
                if (skill.completed)
                  Positioned(
                    top: -4,
                    right: -4,
                    child: Container(
                      width: 28,
                      height: 28,
                      decoration: BoxDecoration(
                        color: LumoraColors.teal,
                        shape: BoxShape.circle,
                        border: Border.all(color: LumoraColors.cream, width: 2),
                      ),
                      child: const Icon(Icons.check_rounded, color: Colors.white, size: 16),
                    ),
                  ),
              ],
            ),
          ),
          const SizedBox(height: 8),
          Text(
            skill.title,
            textAlign: TextAlign.center,
            maxLines: 2,
            overflow: TextOverflow.ellipsis,
            style: const TextStyle(color: LumoraColors.ink, fontSize: 12, fontWeight: FontWeight.w800, height: 1.25),
          ),
          const SizedBox(height: 2),
          Text(
            locked ? 'Finish the previous skill' : '${skill.completedCount}/${skill.lessonCount} lessons',
            textAlign: TextAlign.center,
            style: TextStyle(color: locked ? LumoraColors.amber : LumoraColors.slatey, fontSize: 11, fontWeight: FontWeight.w600),
          ),
        ],
      ),
    );
  }
}

/// The end of the path.
class _RoadmapFinish extends StatelessWidget {
  const _RoadmapFinish();

  @override
  Widget build(BuildContext context) {
    return const Padding(
      padding: EdgeInsets.only(top: 8),
      child: Column(children: [
        Text('🏁', style: TextStyle(fontSize: 30)),
        Text('Fluency', style: TextStyle(color: LumoraColors.slatey, fontSize: 12, fontWeight: FontWeight.w800)),
      ]),
    );
  }
}

/// The dashed centre line the nodes hang off.
class _DashedLine extends StatelessWidget {
  const _DashedLine();

  @override
  Widget build(BuildContext context) => const CustomPaint(size: Size(2, double.infinity), painter: _DashedLinePainter());
}

class _DashedLinePainter extends CustomPainter {
  const _DashedLinePainter();

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = LumoraColors.gray300
      ..strokeWidth = 2;
    for (double y = 0; y < size.height; y += 10) {
      canvas.drawLine(Offset(size.width / 2, y), Offset(size.width / 2, (y + 5).clamp(0, size.height)), paint);
    }
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}
