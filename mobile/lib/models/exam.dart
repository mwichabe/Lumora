import '../core/json_utils.dart';

class Certificate {
  final int id;
  final int userId;
  final String userName;
  final String language;
  final String level;
  final int score;
  final int listening;
  final int reading;
  final int writing;
  final int speaking;
  final String serial;
  final String issuedAt;

  const Certificate({
    required this.id,
    required this.userId,
    required this.userName,
    required this.language,
    required this.level,
    required this.score,
    required this.listening,
    required this.reading,
    required this.writing,
    required this.speaking,
    required this.serial,
    required this.issuedAt,
  });

  factory Certificate.fromJson(Map<String, dynamic> j) => Certificate(
        id: asInt(j['id']),
        userId: asInt(j['userId']),
        userName: asString(j['userName']),
        language: asString(j['language']),
        level: asString(j['level']),
        score: asInt(j['score']),
        listening: asInt(j['listening']),
        reading: asInt(j['reading']),
        writing: asInt(j['writing']),
        speaking: asInt(j['speaking']),
        serial: asString(j['serial']),
        issuedAt: asString(j['issuedAt']),
      );
}

class ExamSectionWeights {
  final int listening;
  final int reading;
  final int writing;
  final int speaking;

  const ExamSectionWeights({
    required this.listening,
    required this.reading,
    required this.writing,
    required this.speaking,
  });

  factory ExamSectionWeights.fromJson(Map<String, dynamic> j) => ExamSectionWeights(
        listening: asInt(j['listening'], 30),
        reading: asInt(j['reading'], 30),
        writing: asInt(j['writing'], 20),
        speaking: asInt(j['speaking'], 20),
      );
}

class ExamSectionScores {
  final int listening;
  final int reading;
  final int writing;
  final int speaking;

  const ExamSectionScores({
    required this.listening,
    required this.reading,
    required this.writing,
    required this.speaking,
  });

  factory ExamSectionScores.fromJson(Map<String, dynamic> j) => ExamSectionScores(
        listening: asInt(j['listening']),
        reading: asInt(j['reading']),
        writing: asInt(j['writing']),
        speaking: asInt(j['speaking']),
      );
}

class ExamResult {
  final bool passed;
  final bool alreadyTaken;
  final int overall;
  final String level;
  final int? passMark;
  final ExamSectionWeights? weights;
  final ExamSectionScores? sections;
  final Certificate? certificate;

  const ExamResult({
    required this.passed,
    required this.alreadyTaken,
    required this.overall,
    required this.level,
    required this.passMark,
    required this.weights,
    required this.sections,
    required this.certificate,
  });

  factory ExamResult.fromJson(Map<String, dynamic> j) => ExamResult(
        passed: asBool(j['passed']),
        alreadyTaken: asBool(j['alreadyTaken']),
        overall: asInt(j['overall']),
        level: asString(j['level']),
        passMark: j['passMark'] == null ? null : asInt(j['passMark']),
        weights: j['weights'] == null ? null : ExamSectionWeights.fromJson(asMap(j['weights'])),
        sections: j['sections'] == null ? null : ExamSectionScores.fromJson(asMap(j['sections'])),
        certificate:
            j['certificate'] == null ? null : Certificate.fromJson(asMap(j['certificate'])),
      );
}

class ExamLevelMeta {
  final int passMark;
  final int durationSeconds;
  const ExamLevelMeta({required this.passMark, required this.durationSeconds});
  factory ExamLevelMeta.fromJson(Map<String, dynamic> j) => ExamLevelMeta(
        passMark: asInt(j['passMark']),
        durationSeconds: asInt(j['durationSeconds']),
      );
}

class ExamMeta {
  final ExamSectionWeights weights;
  final Map<String, ExamLevelMeta> levels;

  const ExamMeta({required this.weights, required this.levels});

  factory ExamMeta.fromJson(Map<String, dynamic> j) => ExamMeta(
        weights: ExamSectionWeights.fromJson(asMap(j['weights'])),
        levels: asMap(j['levels'])
            .map((k, v) => MapEntry(k, ExamLevelMeta.fromJson(asMap(v)))),
      );
}

class PaperQuestion {
  final String question;
  final List<String> options;
  final String correctAnswer;

  const PaperQuestion({required this.question, required this.options, required this.correctAnswer});

  factory PaperQuestion.fromJson(Map<String, dynamic> j) => PaperQuestion(
        question: asString(j['question']),
        options: asStringList(j['options']),
        correctAnswer: asString(j['correctAnswer']),
      );
}

class PaperLine {
  final String character;
  final String text;
  final String translation;

  const PaperLine({required this.character, required this.text, required this.translation});

  factory PaperLine.fromJson(Map<String, dynamic> j) => PaperLine(
        character: asString(j['character']),
        text: asString(j['text']),
        translation: asString(j['translation']),
      );
}

class ExamListeningPaper {
  final String title;
  final List<PaperLine> lines;
  final List<PaperQuestion> questions;
  const ExamListeningPaper({required this.title, required this.lines, required this.questions});
  factory ExamListeningPaper.fromJson(Map<String, dynamic> j) => ExamListeningPaper(
        title: asString(j['title']),
        lines: asList(j['lines'], (e) => PaperLine.fromJson(asMap(e))),
        questions: asList(j['questions'], (e) => PaperQuestion.fromJson(asMap(e))),
      );
}

class ExamReadingPaper {
  final String title;
  final List<String> paragraphs;
  final List<PaperQuestion> questions;
  const ExamReadingPaper({required this.title, required this.paragraphs, required this.questions});
  factory ExamReadingPaper.fromJson(Map<String, dynamic> j) => ExamReadingPaper(
        title: asString(j['title']),
        paragraphs: asStringList(j['paragraphs']),
        questions: asList(j['questions'], (e) => PaperQuestion.fromJson(asMap(e))),
      );
}

class ExamWritingPaper {
  final String prompt;
  final int minWords;
  const ExamWritingPaper({required this.prompt, required this.minWords});
  factory ExamWritingPaper.fromJson(Map<String, dynamic> j) =>
      ExamWritingPaper(prompt: asString(j['prompt']), minWords: asInt(j['minWords']));
}

class ExamSpeakingPaper {
  final String phrase;
  final String speaker;
  final String translation;
  const ExamSpeakingPaper({required this.phrase, required this.speaker, required this.translation});
  factory ExamSpeakingPaper.fromJson(Map<String, dynamic> j) => ExamSpeakingPaper(
        phrase: asString(j['phrase']),
        speaker: asString(j['speaker']),
        translation: asString(j['translation']),
      );
}

class ExamPaper {
  final bool ready;
  final String language;
  final String level;
  final int durationSeconds;
  final int passMark;
  final ExamSectionWeights weights;
  final ExamListeningPaper? listening;
  final ExamReadingPaper? reading;
  final ExamWritingPaper writing;
  final ExamSpeakingPaper speaking;

  const ExamPaper({
    required this.ready,
    required this.language,
    required this.level,
    required this.durationSeconds,
    required this.passMark,
    required this.weights,
    required this.listening,
    required this.reading,
    required this.writing,
    required this.speaking,
  });

  factory ExamPaper.fromJson(Map<String, dynamic> j) => ExamPaper(
        ready: asBool(j['ready']),
        language: asString(j['language']),
        level: asString(j['level']),
        durationSeconds: asInt(j['durationSeconds']),
        passMark: asInt(j['passMark']),
        weights: ExamSectionWeights.fromJson(asMap(j['weights'])),
        listening:
            j['listening'] == null ? null : ExamListeningPaper.fromJson(asMap(j['listening'])),
        reading: j['reading'] == null ? null : ExamReadingPaper.fromJson(asMap(j['reading'])),
        writing: ExamWritingPaper.fromJson(asMap(j['writing'])),
        speaking: ExamSpeakingPaper.fromJson(asMap(j['speaking'])),
      );
}

class CertVerification {
  final bool valid;
  final String? userName;
  final String? language;
  final String? level;
  final int? score;
  final String? serial;
  final String? issuedAt;

  const CertVerification({
    required this.valid,
    this.userName,
    this.language,
    this.level,
    this.score,
    this.serial,
    this.issuedAt,
  });

  factory CertVerification.fromJson(Map<String, dynamic> j) {
    final cert = j['certificate'] as Map<String, dynamic>?;
    return CertVerification(
      valid: asBool(j['valid']),
      userName: cert == null ? null : asString(cert['userName']),
      language: cert == null ? null : asString(cert['language']),
      level: cert == null ? null : asString(cert['level']),
      score: cert == null ? null : asInt(cert['score']),
      serial: cert == null ? null : asString(cert['serial']),
      issuedAt: cert == null ? null : asString(cert['issuedAt']),
    );
  }
}
