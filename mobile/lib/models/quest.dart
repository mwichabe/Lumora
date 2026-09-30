import '../core/json_utils.dart';

class Quest {
  final int id;
  final String title;
  final String description;
  final String icon;
  final int xpReward;
  final int target;

  const Quest({
    required this.id,
    required this.title,
    required this.description,
    required this.icon,
    required this.xpReward,
    required this.target,
  });

  factory Quest.fromJson(Map<String, dynamic> j) => Quest(
        id: asInt(j['id']),
        title: asString(j['title']),
        description: asString(j['description']),
        icon: asString(j['icon'], '⭐'),
        xpReward: asInt(j['xpReward']),
        target: asInt(j['target'], 1),
      );
}

class UserQuest {
  final int id;
  final int questId;
  final String date;
  final int progress;
  final bool completed;
  final Quest? quest;

  const UserQuest({
    required this.id,
    required this.questId,
    required this.date,
    required this.progress,
    required this.completed,
    required this.quest,
  });

  factory UserQuest.fromJson(Map<String, dynamic> j) => UserQuest(
        id: asInt(j['id']),
        questId: asInt(j['questId']),
        date: asString(j['date']),
        progress: asInt(j['progress']),
        completed: asBool(j['completed']),
        quest: j['quest'] == null ? null : Quest.fromJson(asMap(j['quest'])),
      );
}
