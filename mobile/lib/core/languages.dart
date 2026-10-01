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
  LanguageMeta(code: 'ja', name: 'Japanese', nativeName: '日本語', flag: '🇯🇵', available: false),
  LanguageMeta(code: 'zh', name: 'Mandarin', nativeName: '中文', flag: '🇨🇳', available: true),
  LanguageMeta(code: 'ar', name: 'Arabic', nativeName: 'العربية', flag: '🇸🇦', available: false),
  LanguageMeta(code: 'sw', name: 'Swahili', nativeName: 'Kiswahili', flag: '🇰🇪', available: false),
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

/// How to show an exam / certificate level for a language: Mandarin learners
/// see the HSK level they know ("HSK 1 · A1"); everyone else sees CEFR.
/// Mirrors frontend/lib/languages.ts `levelDisplay`.
String levelDisplay(String level, String? lang) {
  final hsk = _kHsk[level];
  if (lang == 'zh' && hsk != null) return level == 'FINAL' ? '$hsk (Advanced)' : '$hsk · $level';
  return level == 'FINAL' ? 'Final Mastery' : level;
}
