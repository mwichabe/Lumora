import 'dart:async';
import 'dart:math';

import 'package:flutter_tts/flutter_tts.dart';
import 'package:speech_to_text/speech_to_text.dart' as stt;

class VoiceProfile {
  final double pitch;
  final double rate;
  const VoiceProfile({required this.pitch, required this.rate});
}

const _kLangLocale = <String, String>{
  'es': 'es-ES', 'de': 'de-DE', 'fr': 'fr-FR', 'en': 'en-US', 'it': 'it-IT',
  'pt': 'pt-PT', 'ja': 'ja-JP', 'zh': 'zh-CN', 'ar': 'ar-SA', 'sw': 'sw-KE',
};

/// Each character gets a distinct pitch/rate so they stay recognisable within
/// any spoken language, mirroring frontend/lib/voices.ts CHARACTER_VOICES.
const _kCharacterVoices = <String, VoiceProfile>{
  'Lumora': VoiceProfile(pitch: 1.05, rate: 0.5),
  'Cora': VoiceProfile(pitch: 1.28, rate: 0.58),
  'Professor Finch': VoiceProfile(pitch: 0.8, rate: 0.42),
  'Blaze': VoiceProfile(pitch: 1.1, rate: 0.56),
  'Mira': VoiceProfile(pitch: 0.96, rate: 0.4),
  'Riko': VoiceProfile(pitch: 0.9, rate: 0.5),
  'Zephyr': VoiceProfile(pitch: 0.86, rate: 0.44),
  'Nana': VoiceProfile(pitch: 0.84, rate: 0.36),
  'Pip': VoiceProfile(pitch: 1.2, rate: 0.62),
};
const _kDefaultVoice = VoiceProfile(pitch: 1, rate: 0.5);

/// Wraps flutter_tts + speech_to_text as a drop-in for the web app's
/// Web-Speech-API wrapper (frontend/lib/voices.ts): per-character voice
/// profiles, sequenced dialogue playback, and simple pronunciation scoring.
class Voices {
  Voices._();
  static final Voices instance = Voices._();

  final FlutterTts _tts = FlutterTts();
  final stt.SpeechToText _stt = stt.SpeechToText();
  String _locale = 'es-ES';
  bool _sttInitialized = false;

  void setSpeechLanguage(String? code) {
    if (code == null || code.isEmpty) return;
    _locale = _kLangLocale[code] ?? code;
  }

  String get locale => _locale;

  Future<void> stopSpeaking() => _tts.stop();

  Future<void> speakAs(String? character, String text, {String? lang}) async {
    if (text.isEmpty) return;
    final profile = _kCharacterVoices[character] ?? _kDefaultVoice;
    await _tts.stop();
    await _tts.setLanguage(lang != null ? (_kLangLocale[lang] ?? lang) : _locale);
    await _tts.setPitch(profile.pitch.clamp(0.5, 2.0));
    await _tts.setSpeechRate(profile.rate);
    await _tts.setVolume(1.0);
    await _tts.speak(text);
  }

  Future<void> speakSequence(
    List<({String? character, String text})> lines, {
    void Function(int index)? onLine,
    bool Function()? shouldContinue,
  }) async {
    for (var i = 0; i < lines.length; i++) {
      if (shouldContinue != null && !shouldContinue()) return;
      onLine?.call(i);
      final completer = Completer<void>();
      _tts.setCompletionHandler(() {
        if (!completer.isCompleted) completer.complete();
      });
      await speakAs(lines[i].character, lines[i].text);
      await completer.future.timeout(const Duration(seconds: 12), onTimeout: () {});
      await Future.delayed(const Duration(milliseconds: 250));
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
      localeId: (lang != null ? (_kLangLocale[lang] ?? lang) : _locale).replaceAll('-', '_'),
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
  out = out.replaceAll(RegExp(r'[.,!¡¿?"]'), '').trim();
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

/// 0–100 pronunciation score combining word overlap and character similarity,
/// ported 1:1 from frontend/lib/voices.ts `scorePronunciation`.
int scorePronunciation(String target, String said) {
  final t = _normalize(target);
  final s = _normalize(said);
  if (s.isEmpty) return 0;
  if (t == s) return 100;

  final tWords = t.split(RegExp(r'\s+')).where((w) => w.isNotEmpty).toList();
  final sWords = s.split(RegExp(r'\s+')).where((w) => w.isNotEmpty).toSet();
  final hit = tWords.where((w) => sWords.contains(w)).length;
  final wordScore = tWords.isNotEmpty ? hit / tWords.length : 0.0;

  final charScore = 1 - _editDistance(t, s) / [t.length, s.length, 1].reduce(max);

  return ((wordScore * 0.6 + charScore * 0.4) * 100).round();
}
