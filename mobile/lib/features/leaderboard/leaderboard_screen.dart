import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../core/theme/shadows.dart';
import '../../models/league.dart';
import '../../providers/league_provider.dart';
import '../../widgets/avatar.dart';
import '../../widgets/lumora_button.dart';

Color _hex(String s) {
  var h = s.replaceAll('#', '');
  if (h.length == 6) h = 'FF$h';
  final v = int.tryParse(h, radix: 16);
  return v == null ? LumoraColors.purple : Color(v);
}

String _fmtRemaining(int seconds) {
  final d = seconds ~/ 86400;
  final h = (seconds % 86400) ~/ 3600;
  final m = (seconds % 3600) ~/ 60;
  if (d > 0) return '${d}d ${h}h';
  if (h > 0) return '${h}h ${m}m';
  return '${m}m';
}

class LeaderboardScreen extends ConsumerWidget {
  const LeaderboardScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(leagueProvider);
    final data = state.data;
    final tint = data != null ? _hex(data.tier.tint) : LumoraColors.purple;

    return Scaffold(
      backgroundColor: LumoraColors.cream,
      body: SafeArea(
        child: RefreshIndicator(
          onRefresh: () => ref.read(leagueProvider.notifier).refresh(),
          child: Stack(
            children: [
              ListView(
                padding: EdgeInsets.zero,
                children: [
                  _Header(data: data, loading: state.loading, tint: tint),
                  if (data != null)
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 16),
                      child: Column(
                        children: [
                          if (!data.casual && data.stage != null) _TournamentBanner(stage: data.stage!),
                          if (!state.loading) _TierRail(data: data),
                          if (data.groupGoal != null && data.joined) _GroupGoal(goal: data.groupGoal!, tint: tint),
                          if (data.me?.flagged == true)
                            Container(
                              margin: const EdgeInsets.only(top: 16),
                              padding: const EdgeInsets.all(14),
                              decoration: BoxDecoration(
                                color: LumoraColors.coralLight,
                                border: Border.all(color: LumoraColors.coral.withValues(alpha: 0.4), width: 2),
                                borderRadius: BorderRadius.circular(LumoraRadii.xl),
                              ),
                              child: Text.rich(TextSpan(children: [
                                const TextSpan(text: 'Your week is under review. ', style: TextStyle(fontWeight: FontWeight.w800)),
                                TextSpan(text: '${data.me!.flagReason}. You keep every point of XP, but promotion and chest rewards are paused for this season.'),
                              ])),
                            ),
                          if (state.loading)
                            const Padding(padding: EdgeInsets.only(top: 40), child: Center(child: CircularProgressIndicator()))
                          else if (data.casual)
                            _CasualCard(busy: state.busy, onLeave: () => ref.read(leagueProvider.notifier).toggleCasual(false))
                          else if (!data.joined)
                            const _NotJoinedCard()
                          else
                            _Standings(data: data, onReport: (r) => ref.read(leagueProvider.notifier).report(r)),
                          if (!state.loading && !data.casual)
                            Padding(
                              padding: const EdgeInsets.symmetric(vertical: 24),
                              child: Center(
                                child: TextButton(
                                  onPressed: state.busy ? null : () => ref.read(leagueProvider.notifier).toggleCasual(true),
                                  child: const Text('Switch to casual mode — learn without the leaderboard',
                                      style: TextStyle(color: LumoraColors.slatey, decoration: TextDecoration.underline)),
                                ),
                              ),
                            ),
                        ],
                      ),
                    ),
                  const SizedBox(height: 24),
                ],
              ),
              if (state.result != null)
                _LeagueCeremonyOverlay(result: state.result!, onClose: () => ref.read(leagueProvider.notifier).closeCeremony()),
            ],
          ),
        ),
      ),
    );
  }
}

class _Header extends StatefulWidget {
  final LeagueStandings? data;
  final bool loading;
  final Color tint;
  const _Header({required this.data, required this.loading, required this.tint});

  @override
  State<_Header> createState() => _HeaderState();
}

class _HeaderState extends State<_Header> {
  int _remaining = 0;
  Timer? _timer;

  @override
  void initState() {
    super.initState();
    _remaining = widget.data?.secondsRemaining ?? 0;
    _timer = Timer.periodic(const Duration(seconds: 1), (_) {
      if (mounted && _remaining > 0) setState(() => _remaining--);
    });
  }

  @override
  void didUpdateWidget(covariant _Header oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (widget.data?.secondsRemaining != null) _remaining = widget.data!.secondsRemaining;
  }

  @override
  void dispose() {
    _timer?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final data = widget.data;
    return Container(
      padding: const EdgeInsets.fromLTRB(24, 24, 24, 28),
      decoration: BoxDecoration(
        borderRadius: const BorderRadius.vertical(bottom: Radius.circular(32)),
        color: widget.tint, // solid league colour — no gradients on this screen
      ),
      child: Column(
        children: [
          Container(width: 64, height: 64, decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.2), borderRadius: BorderRadius.circular(LumoraRadii.lg)),
              child: const Icon(Icons.emoji_events_rounded, color: Colors.white, size: 32)),
          const SizedBox(height: 12),
          Text(data != null ? '${data.tier.name} League' : 'League', textAlign: TextAlign.center,
              style: const TextStyle(color: Colors.white, fontSize: 22, fontWeight: FontWeight.w800)),
          if (!widget.loading && data != null) ...[
            const SizedBox(height: 4),
            Text(
              data.casual
                  ? "Casual mode — you're learning without a leaderboard"
                  : data.joined
                      ? "You're #${data.userRank} of ${data.rows.length} · top ${data.promoteTop ?? 0} promote"
                      : "Complete one lesson to enter this week's race",
              textAlign: TextAlign.center,
              style: const TextStyle(color: Colors.white70, fontSize: 13),
            ),
          ],
          if (!widget.loading && data != null && !data.casual) ...[
            const SizedBox(height: 12),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 6),
              decoration: BoxDecoration(color: Colors.black.withValues(alpha: 0.2), borderRadius: BorderRadius.circular(LumoraRadii.full)),
              child: Row(mainAxisSize: MainAxisSize.min, children: [
                const Icon(Icons.timer_outlined, size: 14, color: Colors.white),
                const SizedBox(width: 6),
                Text('${_fmtRemaining(_remaining)} left', style: const TextStyle(color: Colors.white, fontWeight: FontWeight.w700, fontSize: 12)),
              ]),
            ),
          ],
          if (!widget.loading && data != null) ...[
            const SizedBox(height: 12),
            Wrap(spacing: 8, alignment: WrapAlignment.center, children: [
              if (data.you.fairPlay) _pill(Icons.verified_user, 'Fair play'),
              if (data.you.trophies > 0) _pill(Icons.emoji_events, '${data.you.trophies}'),
              _pill(null, 'Best: ${data.you.best}'),
            ]),
          ],
        ],
      ),
    );
  }

  Widget _pill(IconData? icon, String label) => Container(
        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
        decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.15), borderRadius: BorderRadius.circular(LumoraRadii.full)),
        child: Row(mainAxisSize: MainAxisSize.min, children: [
          if (icon != null) ...[Icon(icon, size: 12, color: Colors.white), const SizedBox(width: 4)],
          Text(label, style: const TextStyle(color: Colors.white, fontSize: 11, fontWeight: FontWeight.w800)),
        ]),
      );
}

class _TierRail extends StatelessWidget {
  final LeagueStandings data;
  const _TierRail({required this.data});

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 84,
      child: ListView(
        padding: const EdgeInsets.symmetric(vertical: 10),
        scrollDirection: Axis.horizontal,
        children: [
          for (final t in data.tiers)
            Container(
              width: 68,
              margin: const EdgeInsets.only(right: 6),
              padding: const EdgeInsets.symmetric(vertical: 6),
              decoration: t.index == data.tier.index
                  ? BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.lg), boxShadow: LumoraShadows.card)
                  : null,
              child: Column(children: [
                Container(
                  width: 36, height: 36,
                  decoration: BoxDecoration(
                    color: (t.index <= data.you.bestTier || t.index == data.tier.index) ? _hex(t.tint) : LumoraColors.gray100,
                    borderRadius: BorderRadius.circular(LumoraRadii.md),
                  ),
                  child: Icon(
                    (t.index <= data.you.bestTier || t.index == data.tier.index) ? Icons.workspace_premium_rounded : Icons.lock_rounded,
                    size: 16, color: (t.index <= data.you.bestTier || t.index == data.tier.index) ? Colors.white : LumoraColors.gray500,
                  ),
                ),
                const SizedBox(height: 4),
                Text(t.name, maxLines: 1, overflow: TextOverflow.ellipsis, style: const TextStyle(fontSize: 10, fontWeight: FontWeight.w800)),
                Text(t.index == 9 ? 'cup' : 'top ${t.promoteTop}', style: const TextStyle(fontSize: 9, color: LumoraColors.gray500)),
              ]),
            ),
        ],
      ),
    );
  }
}

class _TournamentBanner extends StatelessWidget {
  final String stage;
  const _TournamentBanner({required this.stage});

  @override
  Widget build(BuildContext context) {
    final label = stage == 'quarterfinal' ? 'Quarterfinal' : stage == 'semifinal' ? 'Semifinal' : 'Final';
    return Container(
      margin: const EdgeInsets.only(top: 12),
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        color: LumoraColors.purple,
        borderRadius: BorderRadius.circular(LumoraRadii.xl),
      ),
      child: Row(children: [
        const Icon(Icons.sports_martial_arts_rounded, color: Colors.white),
        const SizedBox(width: 12),
        Expanded(
          child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Text('Diamond Tournament · $label', style: const TextStyle(color: Colors.white, fontWeight: FontWeight.w800)),
            Text(stage == 'final' ? 'Top 3 take the trophy. Everyone returns to Diamond next week.' : 'Top 10 of 30 advance to the next round.',
                style: const TextStyle(color: Colors.white70, fontSize: 12)),
          ]),
        ),
      ]),
    );
  }
}

class _GroupGoal extends StatelessWidget {
  final GroupGoal goal;
  final Color tint;
  const _GroupGoal({required this.goal, required this.tint});

  @override
  Widget build(BuildContext context) {
    final pct = goal.target > 0 ? (goal.current / goal.target).clamp(0, 1).toDouble() : 0.0;
    return Container(
      margin: const EdgeInsets.only(top: 12),
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Row(children: [
          Icon(Icons.groups_rounded, size: 16, color: tint),
          const SizedBox(width: 6),
          const Text('Pod group goal', style: TextStyle(fontWeight: FontWeight.w800, fontSize: 13)),
          const Spacer(),
          if (goal.hit) const Icon(Icons.check_circle, size: 16, color: LumoraColors.teal),
        ]),
        const SizedBox(height: 8),
        ClipRRect(
          borderRadius: BorderRadius.circular(LumoraRadii.full),
          child: LinearProgressIndicator(value: pct, minHeight: 8, backgroundColor: LumoraColors.gray100, valueColor: AlwaysStoppedAnimation(tint)),
        ),
        const SizedBox(height: 4),
        Text('${goal.current}/${goal.target} · +${goal.bonus} gems for everyone', style: const TextStyle(fontSize: 11, color: LumoraColors.slatey)),
      ]),
    );
  }
}

class _CasualCard extends StatelessWidget {
  final bool busy;
  final VoidCallback onLeave;
  const _CasualCard({required this.busy, required this.onLeave});
  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(top: 16),
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
      child: Column(children: [
        const Icon(Icons.spa_rounded, size: 32, color: LumoraColors.teal),
        const SizedBox(height: 8),
        const Text("You're in casual mode", style: TextStyle(fontWeight: FontWeight.w800, fontSize: 16)),
        const SizedBox(height: 4),
        const Text('You keep every feature and all your XP — you\'re just never placed in a weekly pod.',
            textAlign: TextAlign.center, style: TextStyle(color: LumoraColors.slatey, fontSize: 13)),
        const SizedBox(height: 16),
        LumoraButton(label: 'Join the league', full: true, loading: busy, onPressed: onLeave),
      ]),
    );
  }
}

class _NotJoinedCard extends StatelessWidget {
  const _NotJoinedCard();
  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(top: 16),
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
      child: Column(children: const [
        Icon(Icons.flag_circle_rounded, size: 32, color: LumoraColors.purple),
        SizedBox(height: 8),
        Text('Ready when you are', style: TextStyle(fontWeight: FontWeight.w800, fontSize: 16)),
        SizedBox(height: 4),
        Text('Complete one lesson this week to join this pod and start earning league points.',
            textAlign: TextAlign.center, style: TextStyle(color: LumoraColors.slatey, fontSize: 13)),
      ]),
    );
  }
}

class _Standings extends StatelessWidget {
  final LeagueStandings data;
  final void Function(LeagueRow) onReport;
  const _Standings({required this.data, required this.onReport});

  @override
  Widget build(BuildContext context) {
    final promote = data.promoteTop ?? 0;
    final demote = data.demoteBottom ?? 0;
    final demoteFrom = demote > 0 ? data.rows.length - demote : -1;

    return Padding(
      padding: const EdgeInsets.only(top: 16),
      child: Column(
        children: [
          for (var i = 0; i < data.rows.length; i++) ...[
            if (i + 1 == promote + 1 && promote > 0)
              _ZoneDivider(promote: true, label: 'Promotion zone — top $promote move up to ${data.tiers[(data.tier.index + 1).clamp(0, data.tiers.length - 1)].name}'),
            if (demoteFrom > 0 && i + 1 == demoteFrom + 1)
              _ZoneDivider(promote: false, label: 'Demotion zone — bottom $demote move down'),
            Padding(padding: const EdgeInsets.only(bottom: 8), child: _Row(row: data.rows[i], tint: _hex(data.tier.tint), onReport: onReport)),
          ],
        ],
      ),
    );
  }
}

class _ZoneDivider extends StatelessWidget {
  final bool promote;
  final String label;
  const _ZoneDivider({required this.promote, required this.label});
  @override
  Widget build(BuildContext context) {
    final color = promote ? LumoraColors.teal : LumoraColors.coral;
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(children: [
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
          decoration: BoxDecoration(color: color, borderRadius: BorderRadius.circular(LumoraRadii.full)),
          child: Row(mainAxisSize: MainAxisSize.min, children: [
            Icon(promote ? Icons.trending_up : Icons.trending_down, size: 12, color: Colors.white),
            const SizedBox(width: 4),
            Text(promote ? 'Promotion' : 'Demotion', style: const TextStyle(color: Colors.white, fontSize: 10, fontWeight: FontWeight.w800)),
          ]),
        ),
        const SizedBox(width: 8),
        Expanded(child: Text(label, style: const TextStyle(fontSize: 11, fontWeight: FontWeight.w700, color: LumoraColors.slatey))),
      ]),
    );
  }
}

class _Row extends StatelessWidget {
  final LeagueRow row;
  final Color tint;
  final void Function(LeagueRow) onReport;
  const _Row({required this.row, required this.tint, required this.onReport});

  @override
  Widget build(BuildContext context) {
    final ringColor = row.isUser ? LumoraColors.purple : row.zone == LeagueZone.promote ? LumoraColors.teal.withValues(alpha: 0.4) : row.zone == LeagueZone.demote ? LumoraColors.coral.withValues(alpha: 0.4) : Colors.transparent;
    return GestureDetector(
      onLongPress: row.isUser || row.reported ? null : () => onReport(row),
      child: Container(
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          color: row.isUser ? LumoraColors.purpleLight : Colors.white,
          borderRadius: BorderRadius.circular(LumoraRadii.xl),
          border: Border.all(color: ringColor, width: 2),
          boxShadow: LumoraShadows.card,
        ),
        child: Row(children: [
          SizedBox(width: 28, child: Center(child: _RankBadge(rank: row.rank))),
          const SizedBox(width: 4),
          LumoraAvatar(name: row.name, avatarColor: row.avatarColor, avatarUrl: row.avatarUrl, size: 44),
          const SizedBox(width: 10),
          Expanded(
            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Row(children: [
                Flexible(child: Text(row.name, overflow: TextOverflow.ellipsis, style: const TextStyle(fontWeight: FontWeight.w800))),
                if (row.isUser) const Text(' (you)', style: TextStyle(color: LumoraColors.purple, fontSize: 11)),
                if (row.fairPlay && !row.isUser) const Padding(padding: EdgeInsets.only(left: 4), child: Icon(Icons.verified_user, size: 12, color: LumoraColors.teal)),
                if (row.flagged) const Padding(padding: EdgeInsets.only(left: 4), child: Icon(Icons.flag, size: 12, color: LumoraColors.coral)),
              ]),
              Row(children: [
                const Icon(Icons.local_fire_department, size: 12, color: LumoraColors.coral),
                Text(' ${row.streak}', style: const TextStyle(fontSize: 11, color: LumoraColors.slatey)),
                if (row.accuracy > 0) Text('  · ${row.accuracy}% acc', style: const TextStyle(fontSize: 11, color: LumoraColors.slatey)),
              ]),
            ]),
          ),
          Text('${row.points}', style: TextStyle(fontWeight: FontWeight.w800, fontSize: 16, color: tint)),
        ]),
      ),
    );
  }
}

class _RankBadge extends StatelessWidget {
  final int rank;
  const _RankBadge({required this.rank});
  @override
  Widget build(BuildContext context) {
    if (rank <= 3) {
      final color = rank == 1 ? const Color(0xFFFFD700) : rank == 2 ? const Color(0xFFC0C0C0) : const Color(0xFFCD7F32);
      return Icon(Icons.emoji_events, color: color, size: 20);
    }
    return Text('$rank', style: const TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.gray500));
  }
}

class _LeagueCeremonyOverlay extends StatelessWidget {
  final LeagueResult result;
  final VoidCallback onClose;
  const _LeagueCeremonyOverlay({required this.result, required this.onClose});

  @override
  Widget build(BuildContext context) {
    final resultLabel = {
      'promoted': 'Promoted!', 'held': 'Held your rank', 'demoted': 'Demoted',
      'champion': 'Champion!', 'eliminated': 'Eliminated', 'qualified': 'Qualified!', 'advanced': 'Advanced!',
    }[result.result] ?? result.result;

    return Positioned.fill(
      child: Container(
        color: Colors.black54,
        child: Center(
          child: Container(
            margin: const EdgeInsets.all(24),
            padding: const EdgeInsets.all(24),
            decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl2)),
            child: Column(mainAxisSize: MainAxisSize.min, children: [
              Icon(
                result.result == 'promoted' || result.result == 'champion' || result.result == 'qualified' || result.result == 'advanced'
                    ? Icons.emoji_events : result.result == 'demoted' || result.result == 'eliminated' ? Icons.trending_down : Icons.check_circle,
                size: 56,
                color: result.result == 'demoted' || result.result == 'eliminated' ? LumoraColors.coral : LumoraColors.amber,
              ),
              const SizedBox(height: 12),
              Text(resultLabel, style: const TextStyle(fontSize: 22, fontWeight: FontWeight.w800)),
              const SizedBox(height: 4),
              Text('${result.from.name} → ${result.to.name}', style: const TextStyle(color: LumoraColors.slatey)),
              const SizedBox(height: 16),
              Row(mainAxisAlignment: MainAxisAlignment.spaceEvenly, children: [
                _stat('Rank', '#${result.rank}'),
                _stat('Points', '${result.points}'),
                _stat('Gems', '+${result.gems}'),
              ]),
              if (result.podium.isNotEmpty) ...[
                const SizedBox(height: 16),
                for (final p in result.podium.take(3))
                  Padding(
                    padding: const EdgeInsets.symmetric(vertical: 2),
                    child: Row(children: [
                      SizedBox(width: 30, child: Text('#${p.rank}', style: const TextStyle(fontWeight: FontWeight.w800))),
                      const SizedBox(width: 8),
                      Expanded(child: Text(p.name, style: TextStyle(fontWeight: p.isUser ? FontWeight.w800 : FontWeight.w500))),
                      Text('${p.points}', style: const TextStyle(color: LumoraColors.slatey)),
                    ]),
                  ),
              ],
              const SizedBox(height: 20),
              LumoraButton(label: 'Continue', full: true, onPressed: onClose),
            ]),
          ),
        ),
      ),
    );
  }

  Widget _stat(String label, String value) => Column(children: [
        Text(value, style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 18)),
        Text(label, style: const TextStyle(color: LumoraColors.slatey, fontSize: 11)),
      ]);
}
