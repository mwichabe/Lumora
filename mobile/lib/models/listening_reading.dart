import '../core/json_utils.dart';

class ListeningMatch {
  final int id;
  final int sessionId;
  final int orderIndex;
  final String word;
  final String translation;

  const ListeningMatch({
    required this.id,
    required this.sessionId,
    required this.orderIndex,
    required this.word,
    required this.translation,
  });

  factory ListeningMatch.fromJson(Map<String, dynamic> j) => ListeningMatch(
        id: asInt(j['id']),
        sessionId: asInt(j['sessionId']),
        orderIndex: asInt(j['orderIndex']),
        word: asString(j['word']),
        translation: asString(j['translation']),
      );
}

class ListeningLine {
  final int id;
  final int sessionId;
  final int orderIndex;
  final String character;
  final String text;
  final String translation;

  const ListeningLine({
    required this.id,
    required this.sessionId,
    required this.orderIndex,
    required this.character,
    required this.text,
    required this.translation,
  });

  factory ListeningLine.fromJson(Map<String, dynamic> j) => ListeningLine(
        id: asInt(j['id']),
        sessionId: asInt(j['sessionId']),
        orderIndex: asInt(j['orderIndex']),
        character: asString(j['character']),
        text: asString(j['text']),
        translation: asString(j['translation']),
      );
}

class ListeningQuestion {
  final int id;
  final int sessionId;
  final int orderIndex;
  final String prompt;
  final String question;
  final List<String>? options;
  final String correctAnswer;

  const ListeningQuestion({
    required this.id,
    required this.sessionId,
    required this.orderIndex,
    required this.prompt,
    required this.question,
    required this.options,
    required this.correctAnswer,
  });

  factory ListeningQuestion.fromJson(Map<String, dynamic> j) => ListeningQuestion(
        id: asInt(j['id']),
        sessionId: asInt(j['sessionId']),
        orderIndex: asInt(j['orderIndex']),
        prompt: asString(j['prompt']),
        question: asString(j['question']),
        options: j['options'] == null ? null : asStringList(j['options']),
        correctAnswer: asString(j['correctAnswer']),
      );
}

class ListeningSession {
  final int id;
  final String language;
  final String unit;
  final String title;
  final String description;
  final int orderIndex;
  final int xpReward;
  final List<ListeningMatch> matches;
  final List<ListeningLine> lines;
  final List<ListeningQuestion> questions;

  const ListeningSession({
    required this.id,
    required this.language,
    required this.unit,
    required this.title,
    required this.description,
    required this.orderIndex,
    required this.xpReward,
    required this.matches,
    required this.lines,
    required this.questions,
  });

  factory ListeningSession.fromJson(Map<String, dynamic> j) => ListeningSession(
        id: asInt(j['id']),
        language: asString(j['language']),
        unit: asString(j['unit']),
        title: asString(j['title']),
        description: asString(j['description']),
        orderIndex: asInt(j['orderIndex']),
        xpReward: asInt(j['xpReward']),
        matches: asList(j['matches'], (e) => ListeningMatch.fromJson(asMap(e))),
        lines: asList(j['lines'], (e) => ListeningLine.fromJson(asMap(e))),
        questions: asList(j['questions'], (e) => ListeningQuestion.fromJson(asMap(e))),
      );
}

class ReadingLine {
  final int id;
  final int sessionId;
  final int orderIndex;
  final String text;
  final String translation;

  const ReadingLine({
    required this.id,
    required this.sessionId,
    required this.orderIndex,
    required this.text,
    required this.translation,
  });

  factory ReadingLine.fromJson(Map<String, dynamic> j) => ReadingLine(
        id: asInt(j['id']),
        sessionId: asInt(j['sessionId']),
        orderIndex: asInt(j['orderIndex']),
        text: asString(j['text']),
        translation: asString(j['translation']),
      );
}

class ReadingQuestion {
  final int id;
  final int sessionId;
  final int orderIndex;
  final String prompt;
  final String question;
  final List<String>? options;
  final String correctAnswer;

  const ReadingQuestion({
    required this.id,
    required this.sessionId,
    required this.orderIndex,
    required this.prompt,
    required this.question,
    required this.options,
    required this.correctAnswer,
  });

  factory ReadingQuestion.fromJson(Map<String, dynamic> j) => ReadingQuestion(
        id: asInt(j['id']),
        sessionId: asInt(j['sessionId']),
        orderIndex: asInt(j['orderIndex']),
        prompt: asString(j['prompt']),
        question: asString(j['question']),
        options: j['options'] == null ? null : asStringList(j['options']),
        correctAnswer: asString(j['correctAnswer']),
      );
}

class ReadingSession {
  final int id;
  final String language;
  final String unit;
  final String title;
  final String description;
  final int orderIndex;
  final int xpReward;
  final List<ReadingLine> lines;
  final List<ReadingQuestion> questions;

  const ReadingSession({
    required this.id,
    required this.language,
    required this.unit,
    required this.title,
    required this.description,
    required this.orderIndex,
    required this.xpReward,
    required this.lines,
    required this.questions,
  });

  factory ReadingSession.fromJson(Map<String, dynamic> j) => ReadingSession(
        id: asInt(j['id']),
        language: asString(j['language']),
        unit: asString(j['unit']),
        title: asString(j['title']),
        description: asString(j['description']),
        orderIndex: asInt(j['orderIndex']),
        xpReward: asInt(j['xpReward']),
        lines: asList(j['lines'], (e) => ReadingLine.fromJson(asMap(e))),
        questions: asList(j['questions'], (e) => ReadingQuestion.fromJson(asMap(e))),
      );
}
