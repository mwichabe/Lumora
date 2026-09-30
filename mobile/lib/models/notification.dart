import '../core/json_utils.dart';

class AppNotification {
  final int id;
  final String kind;
  final String emoji;
  final String tint;
  final String title;
  final String body;
  final String link;
  final bool read;
  final String createdAt;

  const AppNotification({
    required this.id,
    required this.kind,
    required this.emoji,
    required this.tint,
    required this.title,
    required this.body,
    required this.link,
    required this.read,
    required this.createdAt,
  });

  factory AppNotification.fromJson(Map<String, dynamic> j) => AppNotification(
        id: asInt(j['id']),
        kind: asString(j['kind']),
        emoji: asString(j['emoji'], '🔔'),
        tint: asString(j['tint'], '#6C3FC5'),
        title: asString(j['title']),
        body: asString(j['body']),
        link: asString(j['link']),
        read: asBool(j['read']),
        createdAt: asString(j['createdAt']),
      );
}
