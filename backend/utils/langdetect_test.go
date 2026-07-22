package utils

import "testing"

// The detector gates every translation call, so two failure modes matter and
// they are not symmetric. Missing a foreign message costs a translation nobody
// gets. Misreading English as foreign costs an API call *and* shows a
// pointless "translated from" banner on a message that was already English.
// The tests below weight the second case more heavily.

func TestDetectsScriptLanguages(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{"Здравствуйте, как дела сегодня?", "ru"},
		{"Привіт, як справи? Це українська мова їжа", "uk"},
		{"مرحبا كيف حالك اليوم", "ar"},
		{"こんにちは、元気ですか", "ja"},
		{"안녕하세요 반갑습니다", "ko"},
		{"你好，今天天气很好", "zh"},
		{"नमस्ते आप कैसे हैं", "hi"},
		{"Γειά σου τι κάνεις σήμερα", "el"},
	}
	for _, c := range cases {
		t.Run(c.want, func(t *testing.T) {
			got := Detect(c.text)
			if got.Code != c.want {
				t.Errorf("Detect(%q) = %q (%.2f), want %q", c.text, got.Code, got.Confidence, c.want)
			}
			if got.Confidence < MinConfidence {
				t.Errorf("confidence %.2f below the acting threshold", got.Confidence)
			}
		})
	}
}

// Japanese mixes kana with Han characters. Kana is the giveaway and must
// outweigh the Han, or every Japanese message is mislabelled Chinese.
func TestJapaneseIsNotMistakenForChinese(t *testing.T) {
	got := Detect("今日は天気がいいですね")
	if got.Code != "ja" {
		t.Errorf("got %q, want ja — kana should outweigh the Han characters", got.Code)
	}
}

func TestDetectsLatinLanguages(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{"The weather is very nice today and I think we should go outside", "en"},
		{"Hola, ¿cómo estás? Espero que todo esté muy bien contigo", "es"},
		{"Bonjour, comment allez-vous aujourd'hui? J'espère que tout va bien", "fr"},
		{"Guten Tag, wie geht es Ihnen heute? Ich hoffe, alles ist gut", "de"},
		{"Olá, como você está hoje? Espero que esteja tudo bem com você", "pt"},
		{"Ciao, come stai oggi? Spero che vada tutto molto bene", "it"},
		{"Habari yako leo? Naomba tuweze kuzungumza kwa sababu ni muhimu sana", "sw"},
	}
	for _, c := range cases {
		t.Run(c.want, func(t *testing.T) {
			got := Detect(c.text)
			if got.Code != c.want {
				t.Errorf("Detect(%q) = %q (%.2f), want %q", c.text, got.Code, got.Confidence, c.want)
			}
		})
	}
}

// The expensive mistake: translating something that was already English.
func TestEnglishIsNeverQueuedForTranslation(t *testing.T) {
	english := []string{
		"The weather is very nice today and I think we should go outside",
		"I have not been able to find the answer to that question yet",
		"Could you explain how this works when you have a moment?",
		"That was a great lesson, thanks for all of your help with it",
		"What do you think about the new update that they released?",
	}
	for _, text := range english {
		if lang, needs := NeedsTranslation(text); needs {
			t.Errorf("English text queued for translation as %q (%.2f): %q",
				lang.Code, lang.Confidence, text)
		}
	}
}

// Short, symbolic, or ambiguous input must come back unknown rather than
// guessed — a wrong guess here is a wrong banner on someone's message.
func TestAmbiguousInputIsLeftAlone(t *testing.T) {
	cases := []string{
		"ok",
		"haha",
		"👍",
		"👍👍👍",
		"lol",
		"???",
		"Kwame",              // a name
		"https://lumora.app", // a bare URL
		"42",
		"",
		"   ",
	}
	for _, text := range cases {
		t.Run(text, func(t *testing.T) {
			if _, needs := NeedsTranslation(text); needs {
				t.Errorf("queued %q for translation; short/symbolic input should be left alone", text)
			}
		})
	}
}

// Diacritics carry the decision on short messages where the stop-word vote is
// too thin on its own.
func TestDiacriticsDisambiguateShortText(t *testing.T) {
	cases := []struct {
		text string
		want string
	}{
		{"¿Dónde está el baño?", "es"},
		{"Ich möchte gern die Straße sehen", "de"},
	}
	for _, c := range cases {
		t.Run(c.want, func(t *testing.T) {
			if got := Detect(c.text); got.Code != c.want {
				t.Errorf("Detect(%q) = %q, want %q", c.text, got.Code, c.want)
			}
		})
	}
}

func TestNeedsTranslationGatesCorrectly(t *testing.T) {
	if lang, needs := NeedsTranslation("Hola, ¿cómo estás? Espero que todo esté muy bien"); !needs {
		t.Errorf("Spanish was not queued for translation (detected %q, %.2f)",
			lang.Code, lang.Confidence)
	}
	if _, needs := NeedsTranslation("Where is the library from here please"); needs {
		t.Error("English was queued for translation")
	}
}

func TestLanguageNameFallsBackToCode(t *testing.T) {
	if got := LanguageName("es"); got != "Spanish" {
		t.Errorf("LanguageName(es) = %q, want Spanish", got)
	}
	if got := LanguageName("xx"); got != "xx" {
		t.Errorf("LanguageName(xx) = %q, want the code back", got)
	}
}

// A mostly-English message with one stray Greek letter (π, common in maths
// talk) must stay English — the script check requires dominance, not presence.
func TestStrayForeignCharacterDoesNotFlipTheMessage(t *testing.T) {
	got := Detect("The area of the circle is π times the radius squared, right?")
	if got.Code != "en" {
		t.Errorf("got %q, want en — one Greek letter shouldn't flip the message", got.Code)
	}
}

// Mixed-language input is real (code-switching). Whatever the detector picks,
// it must not crash or claim high confidence in a language that isn't there.
func TestMixedLanguageDoesNotPanic(t *testing.T) {
	for _, text := range []string{
		"Hola! I am learning español and it is muy divertido",
		"Das ist very interessant, thanks",
		"안녕 hello bonjour",
	} {
		got := Detect(text)
		if got.Confidence > 1 || got.Confidence < 0 {
			t.Errorf("Detect(%q) confidence %.2f out of range", text, got.Confidence)
		}
	}
}
