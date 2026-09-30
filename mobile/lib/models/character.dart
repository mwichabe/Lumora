import '../core/json_utils.dart';

class CharacterWithFriendship {
  final int id;
  final String name;
  final String species;
  final String role;
  final String personality;
  final String color;
  final String emoji;
  final int friendshipLevel;

  const CharacterWithFriendship({
    required this.id,
    required this.name,
    required this.species,
    required this.role,
    required this.personality,
    required this.color,
    required this.emoji,
    required this.friendshipLevel,
  });

  factory CharacterWithFriendship.fromJson(Map<String, dynamic> j) => CharacterWithFriendship(
        id: asInt(j['id']),
        name: asString(j['name']),
        species: asString(j['species']),
        role: asString(j['role']),
        personality: asString(j['personality']),
        color: asString(j['color'], '#6C3FC5'),
        emoji: asString(j['emoji'], '🦊'),
        friendshipLevel: asInt(j['friendshipLevel']),
      );
}
