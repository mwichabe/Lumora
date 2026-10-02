import 'dart:async';
import 'dart:io' show Platform;
import 'dart:math';

import 'package:flutter/foundation.dart';
import 'package:flutter_tts/flutter_tts.dart';
import 'package:speech_to_text/speech_to_text.dart' as stt;

/// A character's voice. [rate] is on the web scale (frontend/lib/voices.ts):
/// 1.0 = the engine's normal speed. It's converted for flutter_tts — and
/// slowed to a learner's pace — in [Voices._ttsRate].
class VoiceProfile {
  final double pitch;
  final double rate;
  final bool female;
  const VoiceProfile({required this.pitch, required this.rate, required this.female});
}

const _kLangLocale = <String, String>{
  'es': 'es-ES', 'de': 'de-DE', 'fr': 'fr-FR', 'en': 'en-US', 'it': 'it-IT',
  'pt': 'pt-PT', 'ja': 'ja-JP', 'zh': 'zh-CN', 'ar': 'ar-SA', 'sw': 'sw-KE',
};

const _kLangName = <String, String>{
  'es': 'Spanish', 'de': 'German', 'fr': 'French', 'en': 'English', 'it': 'Italian',
  'pt': 'Portuguese', 'ja': 'Japanese', 'zh': 'Chinese (Mandarin)', 'ar': 'Arabic', 'sw': 'Swahili',
};

/// Each character gets a distinct pitch/rate (and, where the phone has
/// several voices for the language, a distinct voice) so they stay
/// recognisable within any spoken language — CHARACTER_VOICES in voices.ts.
const _kCharacterVoices = <String, VoiceProfile>{
  'Lumora': VoiceProfile(pitch: 1.05, rate: 1.0, female: true),
  'Cora': VoiceProfile(pitch: 1.28, rate: 1.22, female: true),
  'Professor Finch': VoiceProfile(pitch: 0.8, rate: 0.88, female: false),
  'Blaze': VoiceProfile(pitch: 1.1, rate: 1.2, female: false),
  'Mira': VoiceProfile(pitch: 0.96, rate: 0.86, female: true),
  'Riko': VoiceProfile(pitch: 0.9, rate: 1.08, female: false),
  'Zephyr': VoiceProfile(pitch: 0.86, rate: 0.94, female: false),
  'Nana': VoiceProfile(pitch: 0.84, rate: 0.74, female: true),
  'Pip': VoiceProfile(pitch: 1.2, rate: 1.35, female: false),
};
const _kDefaultVoice = VoiceProfile(pitch: 1, rate: 1.0, female: true);

// Voice-name fragments that give away a system voice's gender (as voices.ts).
const _kFemaleHints = [
  'female', 'mujer', 'femme', 'frau', 'mónica', 'monica', 'paulina', 'helena', 'elvira',
  'marisol', 'lucia', 'laura', 'katja', 'hedda', 'marlene', 'anna', 'petra', 'amelie',
  'amélie', 'audrey', 'julie', 'celine', 'denise', 'samantha', 'karen', 'moira', 'tessa',
  'fiona', 'alice', 'federica', 'joana', 'luciana', 'tingting', 'meijia', 'sinji', 'kyoko', 'zuzana',
];
const _kMaleHints = [
  'male', 'hombre', 'homme', 'mann', 'jorge', 'diego', 'juan', 'carlos', 'pablo', 'conrad',
  'stefan', 'hans', 'yannick', 'thomas', 'paul', 'henri', 'nicolas', 'daniel', 'alex',
  'fred', 'david', 'luca', 'otoya', 'maged',
];

/// A system voice as flutter_tts reports it.
class _SysVoice {
  final String name;
  final String locale; // normalised: lower-case, "-" separated
  final Map<String, String> raw;
  final int quality;
  _SysVoice(this.raw)
      : name = raw['name'] ?? '',
        locale = (raw['locale'] ?? '').toLowerCase().replaceAll('_', '-'),
        quality = _qualityScore(raw);

  static int _qualityScore(Map<String, String> v) {
    var s = 0;
    switch ((v['quality'] ?? '').toLowerCase()) {
      case 'very high' || 'premium':
        s += 6;
      case 'high' || 'enhanced':
        s += 4;
      case 'normal' || 'default':
        s += 2;
    }
    // Voices that stream from the network stall or fail offline; prefer local.
    if (v['network_required'] == '1') s -= 3;
    if (v['features']?.contains('notInstalled') == true) s -= 10;
    return s;
  }

  bool? get female {
    final n = name.toLowerCase();
    if (_kFemaleHints.any(n.contains)) return true;
    if (_kMaleHints.any(n.contains)) return false;
    return null;
  }
}

/// Wraps flutter_tts + speech_to_text as a drop-in for the web app's
/// Web-Speech-API wrapper (frontend/lib/voices.ts): per-character voice
/// profiles, sequenced dialogue playback, and simple pronunciation scoring.
///
/// Pronunciation must match the CONTENT language: German text is read by a
/// German voice (so "ja" sounds like "ya"), French by a French one, and so on.
/// The app keeps [setSpeechLanguage] in step with the learner's course; a voice
/// from another language is never used as a stand-in.
class Voices {
  Voices._();
  static final Voices instance = Voices._();

  final FlutterTts _tts = FlutterTts();
  final stt.SpeechToText _stt = stt.SpeechToText();
  String _lang = 'es';
  bool _sttInitialized = false;

  Future<void>? _ready;
  List<_SysVoice>? _voices;
  final Map<String, Map<String, String>?> _voiceChoice = {}; // "lang|character" → voice
  final Set<String> _warnedMissing = {};

  // Replaying the same phrase soon after slows it right down — the natural
  // "say that again, slower" gesture.
  String? _lastText;
  DateTime _lastSpokenAt = DateTime(0);

  /// Called once per language when the phone has no voice for it, with a
  /// message telling the learner how to add one. Wired to a SnackBar in main.
  void Function(String message)? onMissingVoice;

  void setSpeechLanguage(String? code) {
    if (code == null || code.isEmpty) return;
    _lang = code.split(RegExp('[-_]')).first.toLowerCase();
  }

  String get locale => _kLangLocale[_lang] ?? _lang;

  Future<void> stopSpeaking() => _tts.stop();

  Future<void> _init() => _ready ??= () async {
        // speak() completes when the line has finished, so dialogue can be
        // sequenced without guessing durations.
        await _tts.awaitSpeakCompletion(true);
        if (!kIsWeb && Platform.isIOS) {
          // Play through the silent switch, like any language-learning audio.
          await _tts.setSharedInstance(true);
          await _tts.setIosAudioCategory(IosTextToSpeechAudioCategory.playback,
              [IosTextToSpeechAudioCategoryOptions.duckOthers]);
        }
        try {
          final raw = await _tts.getVoices;
          _voices = [
            for (final v in (raw as List? ?? const []))
              if (v is Map) _SysVoice(v.map((k, val) => MapEntry('$k', '$val'))),
          ];
        } catch (_) {
          _voices = const [];
        }
      }();

  /// The best installed voice for [lang], varied by character where there's a
  /// choice. Null when the phone has no voice for the language at all.
  Map<String, String>? _pickVoice(String lang, String? character, VoiceProfile profile) {
    final cacheKey = '$lang|${character ?? ''}';
    if (_voiceChoice.containsKey(cacheKey)) return _voiceChoice[cacheKey];

    final want = locale.toLowerCase();
    final pool = (_voices ?? const <_SysVoice>[]).where((v) => v.locale.split('-').first == lang).toList();
    Map<String, String>? choice;
    if (pool.isNotEmpty) {
      int score(_SysVoice v) {
        var s = v.quality;
        if (v.locale == want) s += 3; // the course's own region (de-DE, not de-CH)
        final g = v.female;
        if (g != null) s += g == profile.female ? 2 : -2;
        return s;
      }

      pool.sort((a, b) => score(b).compareTo(score(a)));
      // Spread characters over the best few voices so they sound different.
      final top = pool.where((v) => score(v) >= score(pool.first) - 2).toList();
      final pick = top[(character ?? '').hashCode.abs() % top.length];
      choice = pick.raw;
    }
    _voiceChoice[cacheKey] = choice;
    return choice;
  }

  /// flutter_tts takes 0.5 as normal speed on both Android and iOS. Learners
  /// need clear, unhurried speech, so every character is slowed and none is
  /// ever faster than 90% of normal; a quick replay drops to ~65%.
  double _ttsRate(VoiceProfile profile, {required bool replay}) {
    if (replay) return 0.32;
    return (0.5 * profile.rate * 0.82).clamp(0.3, 0.45);
  }

  Future<void> speakAs(String? character, String text, {String? lang}) async {
    if (text.trim().isEmpty) return;
    await _init();
    final code = (lang ?? _lang).split(RegExp('[-_]')).first.toLowerCase();
    final profile = _kCharacterVoices[character] ?? _kDefaultVoice;

    final now = DateTime.now();
    final replay = text == _lastText && now.difference(_lastSpokenAt) < const Duration(seconds: 10);
    _lastText = text;
    _lastSpokenAt = now;

    await _tts.stop();
    final loc = _kLangLocale[code] ?? code;
    await _tts.setLanguage(loc);
    final voice = _pickVoice(code, character, profile);
    if (voice != null) {
      await _tts.setVoice({
        'name': voice['name'] ?? '',
        'locale': voice['locale'] ?? loc,
        if (voice['identifier'] != null) 'identifier': voice['identifier']!,
      });
    } else if (_voices != null && _voices!.isNotEmpty && _warnedMissing.add(code)) {
      // The phone would read this with another language's voice — wrong
      // accent. Say why, once, rather than mislead silently.
      final name = _kLangName[code] ?? code;
      onMissingVoice?.call(
        !kIsWeb && Platform.isIOS
            ? 'No $name voice is installed. Add one in Settings → Accessibility → Spoken Content → Voices for correct pronunciation.'
            : 'No $name voice is installed. Add it in Settings → Text-to-speech → Install voice data for correct pronunciation.',
      );
    }
    await _tts.setPitch(profile.pitch.clamp(0.5, 2.0));
    await _tts.setSpeechRate(_ttsRate(profile, replay: replay));
    await _tts.setVolume(1.0);
    await _tts.speak(text);
  }

  Future<void> speakSequence(
    List<({String? character, String text})> lines, {
    void Function(int index)? onLine,
    bool Function()? shouldContinue,
  }) async {
    _lastText = null; // a conversation never counts as a replay
    for (var i = 0; i < lines.length; i++) {
      if (shouldContinue != null && !shouldContinue()) return;
      onLine?.call(i);
      // speak() resolves when the line ends (awaitSpeakCompletion). The cap
      // scales with the line so a long sentence at learner pace isn't cut off.
      final cap = Duration(milliseconds: 6000 + lines[i].text.length * 220);
      await speakAs(lines[i].character, lines[i].text).timeout(cap, onTimeout: () {});
      _lastText = null;
      if (shouldContinue != null && !shouldContinue()) return;
      // A beat between speakers so each line can sink in.
      await Future.delayed(const Duration(milliseconds: 700));
    }
  }

  Future<bool> speechRecognitionSupported() async {
    if (_sttInitialized) return true;
    _sttInitialized = await _stt.initialize();
    return _sttInitialized;
  }

  Future<String> recognizeSpeech({String? lang}) async {
    final ok = await speechRecognitionSupported();
    if (!ok) throw Exception('unsupported');
    final completer = Completer<String>();
    await _stt.listen(
      onResult: (result) {
        if (result.finalResult && !completer.isCompleted) {
          completer.complete(result.recognizedWords);
        }
      },
      localeId: (lang != null ? (_kLangLocale[lang] ?? lang) : locale).replaceAll('-', '_'),
    );
    final text = await completer.future.timeout(
      const Duration(seconds: 8),
      onTimeout: () {
        _stt.stop();
        return '';
      },
    );
    return text;
  }
}

String _normalize(String s) {
  var out = s.toLowerCase();
  const accents = {
    'á': 'a', 'à': 'a', 'ä': 'a', 'â': 'a',
    'é': 'e', 'è': 'e', 'ë': 'e', 'ê': 'e',
    'í': 'i', 'ì': 'i', 'ï': 'i', 'î': 'i',
    'ó': 'o', 'ò': 'o', 'ö': 'o', 'ô': 'o',
    'ú': 'u', 'ù': 'u', 'ü': 'u', 'û': 'u',
    'ñ': 'n', 'ç': 'c',
  };
  accents.forEach((k, v) => out = out.replaceAll(k, v));
  out = out.replaceAll(RegExp(r'[.,!¡¿?"，。！？、；：“”‘’（）《》…—]'), '').trim();
  return out;
}

int _editDistance(String a, String b) {
  final m = a.length, n = b.length;
  final dp = List.generate(m + 1, (_) => List.filled(n + 1, 0));
  for (var i = 0; i <= m; i++) dp[i][0] = i;
  for (var j = 0; j <= n; j++) dp[0][j] = j;
  for (var i = 1; i <= m; i++) {
    for (var j = 1; j <= n; j++) {
      dp[i][j] = [
        dp[i - 1][j] + 1,
        dp[i][j - 1] + 1,
        dp[i - 1][j - 1] + (a[i - 1] == b[j - 1] ? 0 : 1),
      ].reduce(min);
    }
  }
  return dp[m][n];
}

final _han = RegExp(r'\p{Script=Han}', unicode: true);
final _hanOrRun = RegExp(r'\p{Script=Han}|[^\p{Script=Han}]+', unicode: true);

/// Splits text into the units we compare and count: words for space-separated
/// languages, and individual characters for Chinese, which has no spaces —
/// otherwise a whole Chinese sentence is a single "word". Mirrors voices.ts.
List<String> _units(String s) {
  final out = <String>[];
  for (final chunk in s.split(RegExp(r'\s+')).where((w) => w.isNotEmpty)) {
    if (!_han.hasMatch(chunk)) {
      out.add(chunk);
      continue;
    }
    out.addAll(_hanOrRun.allMatches(chunk).map((m) => m.group(0)!));
  }
  return out;
}

/// Words written, counting each Chinese character as one (for writing tasks).
int countWords(String text) => _units(_normalize(text)).length;

/// 0–100 pronunciation score combining word overlap and character similarity,
/// ported 1:1 from frontend/lib/voices.ts `scorePronunciation`.
int scorePronunciation(String target, String said) {
  final t = _normalize(target);
  final s = _normalize(said);
  if (s.isEmpty) return 0;
  if (t == s) return 100;

  final tWords = _units(t);
  final sWords = _units(s).toSet();
  final hit = tWords.where((w) => sWords.contains(w)).length;
  final wordScore = tWords.isNotEmpty ? hit / tWords.length : 0.0;

  final charScore = 1 - _editDistance(t, s) / [t.length, s.length, 1].reduce(max);

  return ((wordScore * 0.6 + charScore * 0.4) * 100).round();
}
