class LanguageMeta {
  final String code;
  final String name;
  final String nativeName;
  final String flag;
  final bool available;
  const LanguageMeta({
    required this.code,
    required this.name,
    required this.nativeName,
    required this.flag,
    required this.available,
  });
}

/// Mirrors frontend/lib/languages.ts — only es/de/fr have a full seeded course;
/// the rest are selectable but route into the same starter content for now.
const kLanguages = <LanguageMeta>[
  LanguageMeta(code: 'es', name: 'Spanish', nativeName: 'Español', flag: '🇪🇸', available: true),
  LanguageMeta(code: 'de', name: 'German', nativeName: 'Deutsch', flag: '🇩🇪', available: true),
  LanguageMeta(code: 'fr', name: 'French', nativeName: 'Français', flag: '🇫🇷', available: true),
  LanguageMeta(code: 'ja', name: 'Japanese', nativeName: '日本語', flag: '🇯🇵', available: true),
  LanguageMeta(code: 'zh', name: 'Mandarin', nativeName: '中文', flag: '🇨🇳', available: true),
  LanguageMeta(code: 'ar', name: 'Arabic', nativeName: 'العربية', flag: '🇸🇦', available: false),
  LanguageMeta(code: 'sw', name: 'Swahili', nativeName: 'Kiswahili', flag: '🇰🇪', available: true),
  LanguageMeta(code: 'pt', name: 'Portuguese', nativeName: 'Português', flag: '🇵🇹', available: false),
  LanguageMeta(code: 'it', name: 'Italian', nativeName: 'Italiano', flag: '🇮🇹', available: false),
  LanguageMeta(code: 'ko', name: 'Korean', nativeName: '한국어', flag: '🇰🇷', available: false),
  LanguageMeta(code: 'hi', name: 'Hindi', nativeName: 'हिन्दी', flag: '🇮🇳', available: false),
];

LanguageMeta languageMeta(String code) =>
    kLanguages.firstWhere((l) => l.code == code, orElse: () => kLanguages.first);

String languageName(String code) => languageMeta(code).name;

/// Languages with a full course today — the only ones it's honest to advertise.
final kAvailableLanguages = kLanguages.where((l) => l.available).toList();

/// HSK equivalent of each CEFR level, for Mandarin (HSK 3.0).
const _kHsk = {'A1': 'HSK 1', 'A2': 'HSK 2', 'B1': 'HSK 3', 'B2': 'HSK 4', 'C1': 'HSK 5', 'C2': 'HSK 6', 'FINAL': 'HSK 7–9'};

/// JLPT equivalent of each CEFR level, for Japanese. The JLPT stops at N1, so
/// C2 is "Beyond N1" (classical Japanese, literature, dialects, ceremony).
const _kJlpt = {'A1': 'JLPT N5', 'A2': 'JLPT N4', 'B1': 'JLPT N3', 'B2': 'JLPT N2', 'C1': 'JLPT N1', 'C2': 'Beyond N1', 'FINAL': 'JLPT N1+'};

const _kNativeScales = {'zh': _kHsk, 'ja': _kJlpt};

/// Whether a language's levels are shown on its own exam scale (HSK, JLPT).
bool hasNativeLevels(String? lang) => _kNativeScales.containsKey(lang);

/// Languages written without spaces, whose writing is counted in characters.
bool countsCharacters(String? lang) => lang == 'zh' || lang == 'ja';

/// How to show an exam / certificate level for a language: Mandarin and
/// Japanese learners see the level they know ("HSK 1 · A1", "JLPT N5 · A1");
/// everyone else sees CEFR. Mirrors frontend/lib/languages.ts `levelDisplay`.
String levelDisplay(String level, String? lang) {
  final native = _kNativeScales[lang]?[level];
  if (native != null) return level == 'FINAL' ? '$native (Advanced)' : '$native · $level';
  return level == 'FINAL' ? 'Final Mastery' : level;
}
