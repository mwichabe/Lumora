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
  LanguageMeta(code: 'zh', name: 'Chinese', nativeName: '中文', flag: '🇨🇳', available: false),
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
