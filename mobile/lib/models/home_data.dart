import '../core/json_utils.dart';
import 'lesson.dart';
import 'quest.dart';
import 'user.dart';

class HomeData {
  final User user;
  final Lesson? nextLesson;
  final Skill? nextSkill;
  final List<UserQuest> quests;

  const HomeData({
    required this.user,
    required this.nextLesson,
    required this.nextSkill,
    required this.quests,
  });

  factory HomeData.fromJson(Map<String, dynamic> j) => HomeData(
        user: User.fromJson(asMap(j['user'])),
        nextLesson: j['nextLesson'] == null ? null : Lesson.fromJson(asMap(j['nextLesson'])),
        nextSkill: j['nextSkill'] == null ? null : Skill.fromJson(asMap(j['nextSkill'])),
        quests: asList(j['quests'], (e) => UserQuest.fromJson(asMap(e))),
      );
}
