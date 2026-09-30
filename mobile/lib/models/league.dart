import '../core/json_utils.dart';

class LeaderRow {
  final int id;
  final String name;
  final int xp;
  final int streak;
  final String avatarColor;
  final String avatarUrl;
  final String language;
  final bool isUser;
  final int rank;

  const LeaderRow({
    required this.id,
    required this.name,
    required this.xp,
    required this.streak,
    required this.avatarColor,
    required this.avatarUrl,
    required this.language,
    required this.isUser,
    required this.rank,
  });

  factory LeaderRow.fromJson(Map<String, dynamic> j) => LeaderRow(
        id: asInt(j['id']),
        name: asString(j['name']),
        xp: asInt(j['xp']),
        streak: asInt(j['streak']),
        avatarColor: asString(j['avatarColor'], '#6C3FC5'),
        avatarUrl: asString(j['avatarUrl']),
        language: asString(j['language']),
        isUser: asBool(j['isUser']),
        rank: asInt(j['rank']),
      );
}

class LeagueTier {
  final int index;
  final String name;
  final String tint;
  final int promoteTop;
  final int demoteBottom;
  final int goldGems;
  final int groupBonus;

  const LeagueTier({
    required this.index,
    required this.name,
    required this.tint,
    required this.promoteTop,
    required this.demoteBottom,
    required this.goldGems,
    required this.groupBonus,
  });

  factory LeagueTier.fromJson(Map<String, dynamic> j) => LeagueTier(
        index: asInt(j['index']),
        name: asString(j['name']),
        tint: asString(j['tint'], '#6C3FC5'),
        promoteTop: asInt(j['promoteTop']),
        demoteBottom: asInt(j['demoteBottom']),
        goldGems: asInt(j['goldGems']),
        groupBonus: asInt(j['groupBonus']),
      );
}

enum LeagueZone { promote, hold, demote }

LeagueZone leagueZoneFromString(String s) {
  switch (s) {
    case 'promote':
      return LeagueZone.promote;
    case 'demote':
      return LeagueZone.demote;
    default:
      return LeagueZone.hold;
  }
}

class LeagueRow {
  final int id;
  final String name;
  final int points;
  final int rawXp;
  final int streak;
  final String avatarColor;
  final String avatarUrl;
  final String language;
  final int rank;
  final LeagueZone zone;
  final bool isUser;
  final int accuracy;
  final int perfectRuns;
  final bool fairPlay;
  final bool flagged;
  final bool reported;

  const LeagueRow({
    required this.id,
    required this.name,
    required this.points,
    required this.rawXp,
    required this.streak,
    required this.avatarColor,
    required this.avatarUrl,
    required this.language,
    required this.rank,
    required this.zone,
    required this.isUser,
    required this.accuracy,
    required this.perfectRuns,
    required this.fairPlay,
    required this.flagged,
    required this.reported,
  });

  factory LeagueRow.fromJson(Map<String, dynamic> j) => LeagueRow(
        id: asInt(j['id']),
        name: asString(j['name']),
        points: asInt(j['points']),
        rawXp: asInt(j['rawXp']),
        streak: asInt(j['streak']),
        avatarColor: asString(j['avatarColor'], '#6C3FC5'),
        avatarUrl: asString(j['avatarUrl']),
        language: asString(j['language']),
        rank: asInt(j['rank']),
        zone: leagueZoneFromString(asString(j['zone'])),
        isUser: asBool(j['isUser']),
        accuracy: asInt(j['accuracy']),
        perfectRuns: asInt(j['perfectRuns']),
        fairPlay: asBool(j['fairPlay']),
        flagged: asBool(j['flagged']),
        reported: asBool(j['reported']),
      );
}

class GroupGoal {
  final int target;
  final int current;
  final bool hit;
  final int bonus;
  const GroupGoal({required this.target, required this.current, required this.hit, required this.bonus});
  factory GroupGoal.fromJson(Map<String, dynamic> j) => GroupGoal(
        target: asInt(j['target']),
        current: asInt(j['current']),
        hit: asBool(j['hit']),
        bonus: asInt(j['bonus']),
      );
}

class LeagueMe {
  final int points;
  final int rawXp;
  final int activities;
  final int perfectRuns;
  final bool flagged;
  final String flagReason;
  const LeagueMe({
    required this.points,
    required this.rawXp,
    required this.activities,
    required this.perfectRuns,
    required this.flagged,
    required this.flagReason,
  });
  factory LeagueMe.fromJson(Map<String, dynamic> j) => LeagueMe(
        points: asInt(j['points']),
        rawXp: asInt(j['rawXp']),
        activities: asInt(j['activities']),
        perfectRuns: asInt(j['perfectRuns']),
        flagged: asBool(j['flagged']),
        flagReason: asString(j['flagReason']),
      );
}

class LeagueYou {
  final int integrity;
  final bool fairPlay;
  final int trophies;
  final String best;
  final int bestTier;
  final bool casual;
  const LeagueYou({
    required this.integrity,
    required this.fairPlay,
    required this.trophies,
    required this.best,
    required this.bestTier,
    required this.casual,
  });
  factory LeagueYou.fromJson(Map<String, dynamic> j) => LeagueYou(
        integrity: asInt(j['integrity'], 100),
        fairPlay: asBool(j['fairPlay'], true),
        trophies: asInt(j['trophies']),
        best: asString(j['best']),
        bestTier: asInt(j['bestTier']),
        casual: asBool(j['casual']),
      );
}

class LeagueStandings {
  final String seasonId;
  final String endsAt;
  final int secondsRemaining;
  final List<LeagueTier> tiers;
  final LeagueTier tier;
  final int podSize;
  final bool casual;
  final bool joined;
  final List<LeagueRow> rows;
  final int? userRank;
  final int? promoteTop;
  final int? demoteBottom;
  final String? stage;
  final GroupGoal? groupGoal;
  final LeagueMe? me;
  final LeagueYou you;

  const LeagueStandings({
    required this.seasonId,
    required this.endsAt,
    required this.secondsRemaining,
    required this.tiers,
    required this.tier,
    required this.podSize,
    required this.casual,
    required this.joined,
    required this.rows,
    required this.userRank,
    required this.promoteTop,
    required this.demoteBottom,
    required this.stage,
    required this.groupGoal,
    required this.me,
    required this.you,
  });

  factory LeagueStandings.fromJson(Map<String, dynamic> j) => LeagueStandings(
        seasonId: asString(j['seasonId']),
        endsAt: asString(j['endsAt']),
        secondsRemaining: asInt(j['secondsRemaining']),
        tiers: asList(j['tiers'], (e) => LeagueTier.fromJson(asMap(e))),
        tier: LeagueTier.fromJson(asMap(j['tier'])),
        podSize: asInt(j['podSize'], 30),
        casual: asBool(j['casual']),
        joined: asBool(j['joined']),
        rows: asList(j['rows'], (e) => LeagueRow.fromJson(asMap(e))),
        userRank: j['userRank'] == null ? null : asInt(j['userRank']),
        promoteTop: j['promoteTop'] == null ? null : asInt(j['promoteTop']),
        demoteBottom: j['demoteBottom'] == null ? null : asInt(j['demoteBottom']),
        stage: j['stage'] == null ? null : asString(j['stage']),
        groupGoal: j['groupGoal'] == null ? null : GroupGoal.fromJson(asMap(j['groupGoal'])),
        me: j['me'] == null ? null : LeagueMe.fromJson(asMap(j['me'])),
        you: LeagueYou.fromJson(asMap(j['you'])),
      );
}

class LeaguePodiumEntry {
  final int rank;
  final String name;
  final int points;
  final String avatarColor;
  final String avatarUrl;
  final bool isUser;
  const LeaguePodiumEntry({
    required this.rank,
    required this.name,
    required this.points,
    required this.avatarColor,
    required this.avatarUrl,
    required this.isUser,
  });
  factory LeaguePodiumEntry.fromJson(Map<String, dynamic> j) => LeaguePodiumEntry(
        rank: asInt(j['rank']),
        name: asString(j['name']),
        points: asInt(j['points']),
        avatarColor: asString(j['avatarColor'], '#6C3FC5'),
        avatarUrl: asString(j['avatarUrl']),
        isUser: asBool(j['isUser']),
      );
}

class LeagueResult {
  final String seasonId;
  final String result;
  final int rank;
  final int podSize;
  final int points;
  final int rawXp;
  final int activities;
  final int accuracy;
  final int perfectRuns;
  final int gems;
  final bool groupGoalHit;
  final bool flagged;
  final String flagReason;
  final String stage;
  final LeagueTier from;
  final LeagueTier to;
  final List<LeaguePodiumEntry> podium;
  final int trophies;

  const LeagueResult({
    required this.seasonId,
    required this.result,
    required this.rank,
    required this.podSize,
    required this.points,
    required this.rawXp,
    required this.activities,
    required this.accuracy,
    required this.perfectRuns,
    required this.gems,
    required this.groupGoalHit,
    required this.flagged,
    required this.flagReason,
    required this.stage,
    required this.from,
    required this.to,
    required this.podium,
    required this.trophies,
  });

  factory LeagueResult.fromJson(Map<String, dynamic> j) => LeagueResult(
        seasonId: asString(j['seasonId']),
        result: asString(j['result']),
        rank: asInt(j['rank']),
        podSize: asInt(j['podSize']),
        points: asInt(j['points']),
        rawXp: asInt(j['rawXp']),
        activities: asInt(j['activities']),
        accuracy: asInt(j['accuracy']),
        perfectRuns: asInt(j['perfectRuns']),
        gems: asInt(j['gems']),
        groupGoalHit: asBool(j['groupGoalHit']),
        flagged: asBool(j['flagged']),
        flagReason: asString(j['flagReason']),
        stage: asString(j['stage']),
        from: LeagueTier.fromJson(asMap(j['from'])),
        to: LeagueTier.fromJson(asMap(j['to'])),
        podium: asList(j['podium'], (e) => LeaguePodiumEntry.fromJson(asMap(e))),
        trophies: asInt(j['trophies']),
      );
}

class LeagueHistoryEntry {
  final String seasonId;
  final LeagueTier tier;
  final LeagueTier to;
  final int rank;
  final int points;
  final String result;
  final int gems;
  final String stage;

  const LeagueHistoryEntry({
    required this.seasonId,
    required this.tier,
    required this.to,
    required this.rank,
    required this.points,
    required this.result,
    required this.gems,
    required this.stage,
  });

  factory LeagueHistoryEntry.fromJson(Map<String, dynamic> j) => LeagueHistoryEntry(
        seasonId: asString(j['seasonId']),
        tier: LeagueTier.fromJson(asMap(j['tier'])),
        to: LeagueTier.fromJson(asMap(j['to'])),
        rank: asInt(j['rank']),
        points: asInt(j['points']),
        result: asString(j['result']),
        gems: asInt(j['gems']),
        stage: asString(j['stage']),
      );
}
