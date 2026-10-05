import '../core/json_utils.dart';

enum ExerciseType { translate, multipleChoice, listen, match, speak, fill, write, character, unknown }

ExerciseType exerciseTypeFromString(String s) {
  switch (s) {
    case 'translate':
      return ExerciseType.translate;
    case 'multiple_choice':
      return ExerciseType.multipleChoice;
    case 'listen':
      return ExerciseType.listen;
    case 'match':
      return ExerciseType.match;
    case 'speak':
      return ExerciseType.speak;
    case 'fill':
      return ExerciseType.fill;
    case 'write':
      return ExerciseType.write;
    case 'character':
      return ExerciseType.character;
    default:
      return ExerciseType.unknown;
  }
}

class Exercise {
  final int id;
  final int lessonId;
  final ExerciseType type;
  final int orderIndex;
  final String prompt;
  final String question;
  final List<String>? options;
  final String correctAnswer;
  final String character;

  const Exercise({
    required this.id,
    required this.lessonId,
    required this.type,
    required this.orderIndex,
    required this.prompt,
    required this.question,
    required this.options,
    required this.correctAnswer,
    required this.character,
  });

  factory Exercise.fromJson(Map<String, dynamic> j) => Exercise(
        id: asInt(j['id']),
        lessonId: asInt(j['lessonId']),
        type: exerciseTypeFromString(asString(j['type'])),
        orderIndex: asInt(j['orderIndex']),
        prompt: asString(j['prompt']),
        question: asString(j['question']),
        options: j['options'] == null ? null : asStringList(j['options']),
        correctAnswer: asString(j['correctAnswer']),
        character: asString(j['character']),
      );
}

class VocabItem {
  final int id;
  final int lessonId;
  final int orderIndex;
  final String word;
  final String translation;
  final String example;
  final String exampleTranslation;
  final String speaker;

  const VocabItem({
    required this.id,
    required this.lessonId,
    required this.orderIndex,
    required this.word,
    required this.translation,
    required this.example,
    required this.exampleTranslation,
    required this.speaker,
  });

  factory VocabItem.fromJson(Map<String, dynamic> j) => VocabItem(
        id: asInt(j['id']),
        lessonId: asInt(j['lessonId']),
        orderIndex: asInt(j['orderIndex']),
        word: asString(j['word']),
        translation: asString(j['translation']),
        example: asString(j['example']),
        exampleTranslation: asString(j['exampleTranslation']),
        speaker: asString(j['speaker']),
      );
}

class Mistake {
  final int id;
  final int userId;
  final String language;
  final String prompt;
  final String question;
  final String correctAnswer;

  const Mistake({
    required this.id,
    required this.userId,
    required this.language,
    required this.prompt,
    required this.question,
    required this.correctAnswer,
  });

  factory Mistake.fromJson(Map<String, dynamic> j) => Mistake(
        id: asInt(j['id']),
        userId: asInt(j['userId']),
        language: asString(j['language']),
        prompt: asString(j['prompt']),
        question: asString(j['question']),
        correctAnswer: asString(j['correctAnswer']),
      );
}

class Lesson {
  final int id;
  final int skillId;
  final String title;
  final int orderIndex;
  final int xpReward;
  final List<VocabItem> vocab;
  final List<Exercise> exercises;

  const Lesson({
    required this.id,
    required this.skillId,
    required this.title,
    required this.orderIndex,
    required this.xpReward,
    required this.vocab,
    required this.exercises,
  });

  factory Lesson.fromJson(Map<String, dynamic> j) => Lesson(
        id: asInt(j['id']),
        skillId: asInt(j['skillId']),
        title: asString(j['title']),
        orderIndex: asInt(j['orderIndex']),
        xpReward: asInt(j['xpReward']),
        vocab: asList(j['vocab'], (e) => VocabItem.fromJson(asMap(e))),
        exercises: asList(j['exercises'], (e) => Exercise.fromJson(asMap(e))),
      );
}

class Skill {
  final int id;
  final String language;
  final String unit;
  final String title;
  final String description;
  final String icon;
  final String color;
  final int orderIndex;
  final int requiredXp;
  final List<Lesson> lessons;
  final bool unlocked;
  final bool completed;
  final int lessonCount;
  final int completedCount;

  const Skill({
    required this.id,
    required this.language,
    required this.unit,
    required this.title,
    required this.description,
    required this.icon,
    required this.color,
    required this.orderIndex,
    required this.requiredXp,
    required this.lessons,
    required this.unlocked,
    required this.completed,
    required this.lessonCount,
    required this.completedCount,
  });

  factory Skill.fromJson(Map<String, dynamic> j) => Skill(
        id: asInt(j['id']),
        language: asString(j['language']),
        unit: asString(j['unit']),
        title: asString(j['title']),
        description: asString(j['description']),
        icon: asString(j['icon']),
        color: asString(j['color'], '#6C3FC5'),
        orderIndex: asInt(j['orderIndex']),
        requiredXp: asInt(j['requiredXp']),
        lessons: asList(j['lessons'], (e) => Lesson.fromJson(asMap(e))),
        unlocked: asBool(j['unlocked']),
        completed: asBool(j['completed']),
        lessonCount: asInt(j['lessonCount']),
        completedCount: asInt(j['completedCount']),
      );
}

/// A mistake the writing coach found, with its fix.
class WritingCorrection {
  final String original;
  final String correction;
  final String explanation;
  const WritingCorrection({required this.original, required this.correction, required this.explanation});
  factory WritingCorrection.fromJson(Map<String, dynamic> j) => WritingCorrection(
        original: asString(j['original']),
        correction: asString(j['correction']),
        explanation: asString(j['explanation']),
      );
}

/// The server's verdict on a free-writing answer.
class WritingCheck {
  /// False when rule checks failed — fix [problems] and try again.
  final bool ok;
  final List<String> problems;
  final int wordCount;
  final int minWords;
  final bool acceptable;
  final String summary;
  final List<WritingCorrection> corrections;

  const WritingCheck({
    required this.ok,
    required this.problems,
    required this.wordCount,
    required this.minWords,
    required this.acceptable,
    required this.summary,
    required this.corrections,
  });

  factory WritingCheck.fromJson(Map<String, dynamic> j) => WritingCheck(
        ok: asBool(j['ok']),
        problems: asList(j['problems'], (e) => asString(asMap(e)['message'])),
        wordCount: asInt(j['wordCount']),
        minWords: asInt(j['minWords']),
        acceptable: j['acceptable'] == null ? true : asBool(j['acceptable']),
        summary: asString(j['summary']),
        corrections: asList(j['corrections'], (e) => WritingCorrection.fromJson(asMap(e))),
      );
}
