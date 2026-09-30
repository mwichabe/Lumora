/// The companion cast — colour + short bio, mirroring frontend/lib/characters.ts
/// (minus the DiceBear portrait URLs; the mobile app renders an initial badge
/// instead of fetching remote portraits for every character mention).
class CharacterInfo {
  final String name;
  final String gender;
  final String color;
  final String quote;
  const CharacterInfo({required this.name, required this.gender, required this.color, this.quote = ''});
}

const _kCharacters = <String, CharacterInfo>{
  'Lumora': CharacterInfo(name: 'Lumora', gender: 'female', color: '#6C3FC5', quote: "Every word you learn is a little light. Let's make you shine!"),
  'Professor Finch': CharacterInfo(name: 'Professor Finch', gender: 'male', color: '#8B6F47', quote: 'Grammar is not a cage — it is the architecture of meaning.'),
  'Cora': CharacterInfo(name: 'Cora', gender: 'female', color: '#00C2A8', quote: "Infinite vocabulary, infinite puns. Let's get word-y!"),
  'Blaze': CharacterInfo(name: 'Blaze', gender: 'male', color: '#FF5C5C', quote: "MORE FIRE! Say it louder — you've absolutely got this!"),
  'Mira': CharacterInfo(name: 'Mira', gender: 'female', color: '#9090A0', quote: 'Close your eyes and let the sounds tell their story.'),
  'Riko': CharacterInfo(name: 'Riko', gender: 'male', color: '#F5A623', quote: 'Think you can beat me? ...Fine, maybe you can. A little.'),
  'Zephyr': CharacterInfo(name: 'Zephyr', gender: 'male', color: '#17A3DD', quote: 'Words are wind — give them shape and they carry far.'),
  'Nana': CharacterInfo(name: 'Nana', gender: 'female', color: '#06AECE', quote: 'Slow and steady. Every journey takes the time it needs.'),
  'Pip': CharacterInfo(name: 'Pip', gender: 'male', color: '#F5A623', quote: "Quick-quick! So many quests, so little time — let's go!"),
};

CharacterInfo characterInfo(String? name) {
  if (name != null && _kCharacters.containsKey(name)) return _kCharacters[name]!;
  final safe = (name != null && name.trim().isNotEmpty) ? name : 'Lumora';
  return CharacterInfo(name: safe, gender: 'female', color: '#6C3FC5');
}
