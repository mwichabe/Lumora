package utils

import (
	"sort"
	"strings"
	"unicode"
)

// Offline language detection.
//
// Every message posted anywhere in the app runs through this, so it has to be
// free and instant — no network call, no model. Two stages:
//
//  1. Script detection. Cyrillic, Arabic, CJK, Devanagari and friends are
//     unambiguous from the characters alone; nothing else is needed.
//  2. For Latin-script text, a stop-word vote. Function words ("the", "que",
//     "und") are the highest-signal, highest-frequency tokens in any language
//     and they barely overlap between the languages this app teaches.
//
// The output feeds two decisions: whether to show a "translated from X" badge,
// and whether to spend an API call translating at all. Being wrong in the
// direction of "this is English" is the cheap mistake — it means no
// translation — so the thresholds below are deliberately conservative.

// Lang is a detection result. Confidence runs 0..1; anything below
// MinConfidence should be treated as unknown rather than acted on.
type Lang struct {
	Code       string  // ISO 639-1, or "" when undetermined
	Name       string  // English display name
	Confidence float64 // 0..1
}

// MinConfidence is the bar for acting on a detection. Short messages are
// genuinely ambiguous — "no" is a word in English, Spanish and Italian — and a
// wrong translation is worse than none.
const MinConfidence = 0.55

// Below this many recognisable words, a stop-word vote is noise. "ok",
// "haha", "👍" and bare names should never be translated.
const minWordsForLatin = 3

var languageNames = map[string]string{
	"en": "English", "es": "Spanish", "fr": "French", "de": "German",
	"pt": "Portuguese", "it": "Italian", "nl": "Dutch", "sw": "Swahili",
	"ru": "Russian", "uk": "Ukrainian", "ar": "Arabic", "he": "Hebrew",
	"hi": "Hindi", "zh": "Chinese", "ja": "Japanese", "ko": "Korean",
	"el": "Greek", "th": "Thai",
}

// LanguageName maps a code to its English name, falling back to the code.
func LanguageName(code string) string {
	if n, ok := languageNames[code]; ok {
		return n
	}
	return code
}

// Stop words per language. Kept to high-frequency function words that are
// distinctive: shared tokens like Spanish/Italian "no" or Portuguese/Spanish
// "de" are deliberately absent from one side or scored by overlap below.
var stopWords = map[string][]string{
	"en": {"the", "and", "is", "are", "was", "were", "this", "that", "with", "for",
		"you", "your", "have", "has", "not", "but", "what", "when", "how", "would",
		"could", "should", "there", "their", "about", "from", "they", "been", "will"},
	"es": {"el", "la", "los", "las", "que", "y", "es", "son", "para", "con",
		"por", "una", "pero", "como", "más", "muy", "este", "esta", "cuando", "porque",
		"también", "todo", "hacer", "puede", "tiene", "estoy", "gracias", "hola"},
	"fr": {"le", "la", "les", "des", "une", "est", "sont", "que", "qui", "pour",
		"avec", "dans", "pas", "mais", "vous", "nous", "être", "avoir", "cette", "comme",
		"plus", "très", "bien", "merci", "bonjour", "aussi", "faire", "peut"},
	"de": {"der", "die", "das", "und", "ist", "sind", "nicht", "mit", "für", "auf",
		"ein", "eine", "auch", "aber", "wenn", "wie", "was", "haben", "werden", "kann",
		"sehr", "schon", "noch", "hallo", "danke", "über", "durch", "diese"},
	"pt": {"o", "os", "as", "um", "uma", "não", "para", "com", "por", "mas",
		"como", "mais", "muito", "este", "esta", "quando", "porque", "também", "fazer", "pode",
		"tem", "obrigado", "olá", "você", "isso", "ser", "está"},
	"it": {"il", "lo", "gli", "un", "una", "che", "per", "con", "non", "ma",
		"come", "più", "molto", "questo", "questa", "quando", "perché", "anche", "fare", "può",
		"sono", "grazie", "ciao", "essere", "della", "nella"},
	"nl": {"de", "het", "een", "en", "van", "is", "zijn", "niet", "met", "voor",
		"maar", "ook", "als", "dat", "deze", "kan", "heeft", "wordt", "hoe", "wat",
		"bedankt", "hallo", "heel", "naar"},
	"sw": {"na", "ya", "wa", "kwa", "ni", "katika", "hii", "hiyo", "kwamba", "lakini",
		"sana", "kama", "pia", "hapa", "kuna", "asante", "habari", "nzuri", "mimi", "wewe",
		"tunaweza", "nini", "kwanini", "yangu", "yako"},
}

// Characters that only appear in some Latin-script languages. A single "ß" or
// "ñ" is often the whole answer on a short message where the stop-word vote is
// too thin to be trusted.
var diacriticHints = map[rune][]string{
	'ñ': {"es"}, '¿': {"es"}, '¡': {"es"},
	'ß': {"de"}, 'ä': {"de"}, 'ö': {"de", "sv"}, 'ü': {"de"},
	'ç': {"fr", "pt"}, 'œ': {"fr"}, 'è': {"fr", "it"}, 'ê': {"fr"}, 'û': {"fr"},
	'ã': {"pt"}, 'õ': {"pt"},
	'ì': {"it"}, 'ò': {"it"},
}

// Detect identifies the language of a piece of text.
//
// Returns a zero Lang when the text is too short, too symbolic, or genuinely
// ambiguous. Callers should treat that as "leave it alone".
func Detect(text string) Lang {
	trimmed := strings.TrimSpace(text)
	if len([]rune(trimmed)) < 3 {
		return Lang{}
	}

	if l := detectByScript(trimmed); l.Code != "" {
		return l
	}
	return detectLatin(trimmed)
}

// detectByScript resolves the non-Latin writing systems, where the characters
// themselves identify the language with near-certainty.
func detectByScript(text string) Lang {
	counts := map[string]int{}
	letters := 0

	for _, r := range text {
		if !unicode.IsLetter(r) {
			continue
		}
		letters++
		switch {
		case unicode.Is(unicode.Han, r):
			counts["zh"]++
		case unicode.Is(unicode.Hiragana, r), unicode.Is(unicode.Katakana, r):
			// Japanese mixes kana with Han; kana is the giveaway, so it
			// outweighs any Han characters in the same message.
			counts["ja"] += 3
		case unicode.Is(unicode.Hangul, r):
			counts["ko"]++
		case unicode.Is(unicode.Cyrillic, r):
			// Cyrillic is scored as one script here; which language it is gets
			// decided below by the letters exclusive to each. Counting
			// Ukrainian-only letters against all other Cyrillic would never
			// work — shared letters outnumber them in every real sentence.
			counts["cyrillic"]++
		case unicode.Is(unicode.Arabic, r):
			counts["ar"]++
		case unicode.Is(unicode.Hebrew, r):
			counts["he"]++
		case unicode.Is(unicode.Devanagari, r):
			counts["hi"]++
		case unicode.Is(unicode.Greek, r):
			counts["el"]++
		case unicode.Is(unicode.Thai, r):
			counts["th"]++
		}
	}

	if letters == 0 {
		return Lang{}
	}
	best, bestN := "", 0
	for code, n := range counts {
		if n > bestN {
			best, bestN = code, n
		}
	}
	// A stray Greek letter in an English maths message shouldn't flip the
	// whole message; require the script to dominate.
	if best == "" || float64(bestN)/float64(letters) < 0.4 {
		return Lang{}
	}
	if best == "cyrillic" {
		return Lang{Code: cyrillicLanguage(text), Name: LanguageName(cyrillicLanguage(text)), Confidence: 0.95}
	}
	return Lang{Code: best, Name: LanguageName(best), Confidence: 0.97}
}

// cyrillicLanguage separates Ukrainian from Russian by the letters exclusive to
// each alphabet. Presence is the signal, not frequency: і/ї/є/ґ do not exist in
// Russian and ы/э/ъ do not exist in Ukrainian, so a single one is decisive
// where a count would be drowned out by the letters they share.
func cyrillicLanguage(text string) string {
	ukrainianOnly := 0
	russianOnly := 0
	for _, r := range text {
		switch {
		case strings.ContainsRune("іїєґІЇЄҐ", r):
			ukrainianOnly++
		case strings.ContainsRune("ыэъЫЭЪ", r):
			russianOnly++
		}
	}
	if ukrainianOnly > russianOnly {
		return "uk"
	}
	// Russian is the default for unmarked Cyrillic — it's by far the most
	// likely, and a Ukrainian sentence of any length contains an і.
	return "ru"
}

// detectLatin scores Latin-script text by stop-word overlap, nudged by
// language-specific diacritics.
func detectLatin(text string) Lang {
	words := tokenize(text)
	if len(words) < minWordsForLatin {
		return Lang{}
	}

	scores := map[string]float64{}
	for code, list := range stopWords {
		set := make(map[string]bool, len(list))
		for _, w := range list {
			set[w] = true
		}
		hits := 0
		for _, w := range words {
			if set[w] {
				hits++
			}
		}
		scores[code] = float64(hits) / float64(len(words))
	}

	// Diacritics are worth a lot on short text, where a couple of stop words
	// either way is within the noise.
	for _, r := range strings.ToLower(text) {
		for _, code := range diacriticHints[r] {
			if _, tracked := scores[code]; tracked {
				scores[code] += 0.12
			}
		}
	}

	type scored struct {
		code string
		s    float64
	}
	ranked := make([]scored, 0, len(scores))
	for code, s := range scores {
		ranked = append(ranked, scored{code, s})
	}
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].s != ranked[j].s {
			return ranked[i].s > ranked[j].s
		}
		return ranked[i].code < ranked[j].code // stable for equal scores
	})

	if len(ranked) == 0 || ranked[0].s == 0 {
		// No function words matched anything. Unknown rather than a guess —
		// this is usually a name, a URL, or shorthand.
		return Lang{}
	}

	top := ranked[0]
	runnerUp := 0.0
	if len(ranked) > 1 {
		runnerUp = ranked[1].s
	}

	// Confidence combines how much of the text matched with how clearly the
	// winner beat the runner-up. Two closely-scored Romance languages should
	// come out uncertain, not confidently wrong.
	margin := top.s - runnerUp
	confidence := top.s*2 + margin*3
	if confidence > 0.99 {
		confidence = 0.99
	}
	return Lang{Code: top.code, Name: LanguageName(top.code), Confidence: confidence}
}

// tokenize lowercases and splits on anything that isn't a letter or an
// apostrophe, keeping accented characters intact.
func tokenize(text string) []string {
	fields := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && r != '\''
	})
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		if f != "" {
			out = append(out, f)
		}
	}
	return out
}

// NeedsTranslation reports whether text should be offered in English —
// confidently detected, and confidently not English already.
func NeedsTranslation(text string) (Lang, bool) {
	l := Detect(text)
	if l.Code == "" || l.Code == "en" || l.Confidence < MinConfidence {
		return l, false
	}
	return l, true
}
