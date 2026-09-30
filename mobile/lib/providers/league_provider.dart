import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../core/network/api_client.dart';
import '../models/league.dart';

class LeagueState {
  final LeagueStandings? data;
  final LeagueResult? result;
  final bool loading;
  final bool busy;
  const LeagueState({required this.data, required this.result, required this.loading, required this.busy});

  factory LeagueState.initial() => const LeagueState(data: null, result: null, loading: true, busy: false);

  LeagueState copyWith({LeagueStandings? data, LeagueResult? result, bool clearResult = false, bool? loading, bool? busy}) =>
      LeagueState(
        data: data ?? this.data,
        result: clearResult ? null : (result ?? this.result),
        loading: loading ?? this.loading,
        busy: busy ?? this.busy,
      );
}

class LeagueController extends Notifier<LeagueState> {
  @override
  LeagueState build() {
    _load();
    return LeagueState.initial();
  }

  Future<void> _load() async {
    try {
      final data = await ApiClient.instance.league();
      LeagueResult? result;
      try {
        result = await ApiClient.instance.leagueResult();
      } catch (_) {}
      state = state.copyWith(data: data, result: result, loading: false);
    } catch (_) {
      state = state.copyWith(loading: false);
    }
  }

  Future<void> refresh() => _load();

  Future<void> closeCeremony() async {
    final r = state.result;
    if (r != null) {
      try {
        await ApiClient.instance.markLeagueResultSeen(r.seasonId);
      } catch (_) {}
    }
    state = state.copyWith(clearResult: true);
    await _load();
  }

  Future<void> toggleCasual(bool enabled) async {
    state = state.copyWith(busy: true);
    try {
      await ApiClient.instance.setLeagueCasual(enabled);
      await _load();
    } finally {
      state = state.copyWith(busy: false);
    }
  }

  Future<void> report(LeagueRow row) async {
    if (row.reported || row.isUser || state.data == null) return;
    final rows = [
      for (final r in state.data!.rows)
        if (r.id == row.id)
          LeagueRow(id: r.id, name: r.name, points: r.points, rawXp: r.rawXp, streak: r.streak,
              avatarColor: r.avatarColor, avatarUrl: r.avatarUrl, language: r.language, rank: r.rank,
              zone: r.zone, isUser: r.isUser, accuracy: r.accuracy, perfectRuns: r.perfectRuns,
              fairPlay: r.fairPlay, flagged: r.flagged, reported: true)
        else
          r,
    ];
    state = state.copyWith(data: LeagueStandings(
      seasonId: state.data!.seasonId, endsAt: state.data!.endsAt, secondsRemaining: state.data!.secondsRemaining,
      tiers: state.data!.tiers, tier: state.data!.tier, podSize: state.data!.podSize, casual: state.data!.casual,
      joined: state.data!.joined, rows: rows, userRank: state.data!.userRank, promoteTop: state.data!.promoteTop,
      demoteBottom: state.data!.demoteBottom, stage: state.data!.stage, groupGoal: state.data!.groupGoal,
      me: state.data!.me, you: state.data!.you,
    ));
    try {
      await ApiClient.instance.reportLeagueMember(row.id, 'suspicious activity');
    } catch (_) {}
  }
}

final leagueProvider = NotifierProvider<LeagueController, LeagueState>(LeagueController.new);
