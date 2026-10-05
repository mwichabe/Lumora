package database

import (
	"gorm.io/gorm"

	"lumora/backend/models"
)

// ===== Japanese course (JLPT N5 → N1 and beyond, mapped onto CEFR A1 → C2) =====
//
// Structured on the JLPT (Japanese-Language Proficiency Test). The app's level
// system is CEFR-based, so each JLPT level is one CEFR-tagged unit, using the
// common approximation:
//
//	N5 = A1 · N4 = A2 · N3 = B1 · N2 = B2 · N1 = C1 · beyond N1 = C2
//
// The JLPT stops at N1, so C2 ("Beyond N1") covers what near-native learners
// still have to master: classical Japanese, literature, dialects, pragmatics
// and ceremonial keigo. The Final exam spans everything.
//
// N5 opens with the foundations every learner needs first — hiragana, then
// katakana, then the sound system and pitch accent — before survival topics.
// Kanji are introduced in N5 and integrated with vocabulary at every level.
// Readings: every vocabulary item shows its kana reading; romaji is shown
// through N5 only, so learners stop leaning on it early.
//
// Every exercise carries its own answer options (mc / listen / match), so none
// rely on the server's generic distractor fallback, which is Spanish. Japanese
// lessons are answered by choosing (a learner can't be assumed to have a
// Japanese keyboard set up), except the free-writing tasks from N4 up.

// jv builds a Japanese vocabulary item: the word (as written) is what's spoken
// aloud; its reading rides along in the translation so learners always see it.
func jv(word, reading, english, example, exampleReading, exampleEnglish, speaker string) models.VocabItem {
	return vw(word, reading+" · "+english, example, exampleReading+" — "+exampleEnglish, speaker)
}

func seedJapanese(db *gorm.DB) {
	seedJapaneseN5(db)
	seedJapaneseN4(db)
	seedJapaneseN3(db)
	seedJapaneseN2(db)
	seedJapaneseN1(db)
	seedJapaneseBeyondN1(db)
	seedJapaneseListening(db)
	seedJapaneseReading(db)
}

// ─────────────── N5 · A1 — Foundations & Survival Japanese ───────────────

func seedJapaneseN5(db *gorm.DB) {
	const ja = "ja"
	const u = "A1 · JLPT N5 — Foundations & Survival Japanese"
	finch := "Professor Finch"

	// ── Stage 0: hiragana ──
	s := addSkillL(db, ja, u, "Hiragana I", "The first script: vowels and the K, S and T rows.", "Languages", "#6C3FC5", 1, 0)
	l := addLesson(db, s, "Vowels & the K-row", 1, 15,
		char(finch, "Japanese uses three scripts together. Hiragana comes first: 46 basic characters, each one a sound (a mora). Start with the five vowels — あ a, い i, う u, え e, お o — then add k: か ka, き ki, く ku, け ke, こ ko."),
		listen("Listen and choose the character", "あ", "あ (a)", "あ (a)", "い (i)", "う (u)", "お (o)"),
		listen("Listen and choose the character", "け", "け (ke)", "か (ka)", "き (ki)", "け (ke)", "こ (ko)"),
		mc("How do you read this?", "いえ", "ie (house)", "ie (house)", "ue (above)", "ai (love)", "ao (blue)"),
		mc("How do you read this?", "かお", "kao (face)", "kao (face)", "koe (voice)", "kaki (persimmon)", "eki (station)"),
		mc("Which hiragana is 'ku'?", "ku", "く", "く", "へ", "し", "つ"),
		match("Match the reading", "こえ", "koe (voice)", "koe (voice)", "kao (face)", "kiku (to listen)", "ike (pond)"),
		speak("Mira", "あいうえお、かきくけこ"),
	)
	addVocab(db, l,
		jv("あい", "ai", "love", "あいはだいじです。", "Ai wa daiji desu.", "Love is important.", "Mira"),
		jv("いえ", "ie", "house, home", "いえにかえります。", "Ie ni kaerimasu.", "I'm going home.", "Cora"),
		jv("かお", "kao", "face", "かおをあらいます。", "Kao o araimasu.", "I wash my face.", "Cora"),
		jv("いく", "iku", "to go", "いきましょう！", "Ikimashō!", "Let's go!", "Blaze"),
		jv("あおい", "aoi", "blue", "そらはあおいです。", "Sora wa aoi desu.", "The sky is blue.", "Mira"),
	)
	l = addLesson(db, s, "The S-row & T-row", 2, 15,
		char(finch, "S-row: さ sa, し shi (not 'si'), す su, せ se, そ so. T-row: た ta, ち chi, つ tsu, て te, と to. Note the irregular readings: し shi, ち chi and つ tsu."),
		listen("Listen and choose the character", "し", "し (shi)", "さ (sa)", "し (shi)", "す (su)", "そ (so)"),
		listen("Listen and choose the character", "つ", "つ (tsu)", "た (ta)", "ち (chi)", "つ (tsu)", "と (to)"),
		mc("How do you read this?", "すし", "sushi", "sushi", "sashi", "soshi", "susu"),
		mc("How do you read this?", "ちかてつ", "chikatetsu (subway)", "chikatetsu (subway)", "tsukatechi", "chikatesu", "shikatetsu"),
		mc("Which hiragana is 'chi'?", "chi", "ち", "ち", "さ", "き", "た"),
		mc("How do you read this?", "あした", "ashita (tomorrow)", "ashita (tomorrow)", "asita", "ashito", "oshita"),
		speak("Blaze", "さしすせそ、たちつてと"),
	)
	addVocab(db, l,
		jv("すし", "sushi", "sushi", "すしがすきです。", "Sushi ga suki desu.", "I like sushi.", "Cora"),
		jv("あした", "ashita", "tomorrow", "またあした。", "Mata ashita.", "See you tomorrow.", "Lumora"),
		jv("ちかてつ", "chikatetsu", "subway", "ちかてつでいきます。", "Chikatetsu de ikimasu.", "I'll go by subway.", "Riko"),
		jv("くつ", "kutsu", "shoes", "くつをぬいでください。", "Kutsu o nuide kudasai.", "Please take off your shoes.", finch),
		jv("せかい", "sekai", "world", "せかいはひろいです。", "Sekai wa hiroi desu.", "The world is wide.", "Zephyr"),
	)

	s = addSkillL(db, ja, u, "Hiragana II", "The N, H, M, Y, R and W rows, and ん.", "Languages", "#00C2A8", 2, 8)
	l = addLesson(db, s, "N, H & M Rows", 1, 15,
		char(finch, "N-row: な ni ぬ ね の. H-row: は hi ふ (fu) へ ほ. M-row: ま み む め も. Careful: as the topic particle, は is read 'wa', and as the direction particle へ is read 'e'."),
		listen("Listen and choose the character", "ふ", "ふ (fu)", "は (ha)", "ひ (hi)", "ふ (fu)", "ほ (ho)"),
		listen("Listen and choose the character", "ね", "ね (ne)", "ぬ (nu)", "ね (ne)", "れ (re)", "わ (wa)"),
		mc("How do you read this?", "ねこ", "neko (cat)", "neko (cat)", "nuko", "reko", "neku"),
		mc("How do you read this?", "はな", "hana (flower / nose)", "hana (flower / nose)", "wana", "hama", "nana"),
		mc("How is は read in わたしは?", "わたしは (as for me)", "wa", "wa", "ha", "ho", "ba"),
		mc("How do you read this?", "みみ", "mimi (ear)", "mimi (ear)", "nini", "mumu", "meme"),
		speak("Mira", "なにぬねの、はひふへほ、まみむめも"),
	)
	addVocab(db, l,
		jv("ねこ", "neko", "cat", "ねこがいます。", "Neko ga imasu.", "There's a cat.", "Pip"),
		jv("はな", "hana", "flower; nose", "はながきれいです。", "Hana ga kirei desu.", "The flowers are pretty.", "Mira"),
		jv("みみ", "mimi", "ear", "みみがいたいです。", "Mimi ga itai desu.", "My ear hurts.", "Nana"),
		jv("ふね", "fune", "boat, ship", "ふねにのります。", "Fune ni norimasu.", "I'll get on the boat.", "Zephyr"),
		jv("むし", "mushi", "insect", "むしがこわいです。", "Mushi ga kowai desu.", "I'm scared of insects.", "Pip"),
	)
	l = addLesson(db, s, "Y, R & W Rows and ん", 2, 15,
		char(finch, "Y-row has only three: や ゆ よ. R-row: ら り る れ ろ — the Japanese r is a light flap, between r, l and d. W-row: わ and を (read 'o', used only as the object particle). Finally ん n, the only consonant that stands alone."),
		listen("Listen and choose the character", "ろ", "ろ (ro)", "る (ru)", "ろ (ro)", "れ (re)", "ら (ra)"),
		listen("Listen and choose the character", "ゆ", "ゆ (yu)", "や (ya)", "ゆ (yu)", "よ (yo)", "わ (wa)"),
		mc("How do you read this?", "ほん", "hon (book)", "hon (book)", "hono", "han", "ho"),
		mc("How do you read this?", "りんご", "ringo (apple)", "ringo (apple)", "rinko", "lingo", "rengo"),
		mc("What is を used for?", "パンをたべます。", "Marking the object of a verb", "Marking the object of a verb", "Marking the topic", "Marking a place", "Asking a question"),
		mc("How do you read this?", "よる", "yoru (night)", "yoru (night)", "yora", "yuru", "yoro"),
		speak("Blaze", "やゆよ、らりるれろ、わをん"),
	)
	addVocab(db, l,
		jv("ほん", "hon", "book", "ほんをよみます。", "Hon o yomimasu.", "I read a book.", finch),
		jv("よる", "yoru", "night", "よるはしずかです。", "Yoru wa shizuka desu.", "The night is quiet.", "Nana"),
		jv("りんご", "ringo", "apple", "りんごをたべます。", "Ringo o tabemasu.", "I eat an apple.", "Cora"),
		jv("わたし", "watashi", "I, me", "わたしはがくせいです。", "Watashi wa gakusei desu.", "I am a student.", "Lumora"),
		jv("えん", "en", "yen", "ひゃくえんです。", "Hyaku en desu.", "It's 100 yen.", "Cora"),
	)

	s = addSkillL(db, ja, u, "Hiragana III", "Dakuten, combined sounds, small っ and long vowels.", "Languages", "#F5A623", 3, 16)
	l = addLesson(db, s, "Voiced Sounds & Combinations", 1, 15,
		char(finch, "Two dots (゛, dakuten) voice a sound: か ka → が ga, さ sa → ざ za, た ta → だ da, は ha → ば ba. A circle (゜) makes p: ぱ pa. A small ゃ ゅ ょ merges with the sound before it: き + ゃ = きゃ kya, し + ょ = しょ sho."),
		mc("Add dakuten to か", "か + ゛", "が (ga)", "が (ga)", "ぱ (pa)", "ざ (za)", "だ (da)"),
		mc("Which one is 'pa'?", "pa", "ぱ", "ぱ", "ば", "は", "ほ"),
		mc("How do you read this?", "きゃく", "kyaku (guest)", "kyaku (guest)", "kiyaku", "kyoku", "kaku"),
		mc("How do you read this?", "しゅくだい", "shukudai (homework)", "shukudai (homework)", "shiyukudai", "sukudai", "shokudai"),
		listen("Listen and choose", "ぎゅうにゅう", "ぎゅうにゅう (gyūnyū, milk)", "ぎゅうにゅう (gyūnyū, milk)", "きゅうにゅう", "ぎょうにょう", "ぐうにゅう"),
		speak("Cora", "がぎぐげご、ぱぴぷぺぽ、きゃきゅきょ"),
	)
	addVocab(db, l,
		jv("がっこう", "gakkō", "school", "がっこうにいきます。", "Gakkō ni ikimasu.", "I go to school.", "Pip"),
		jv("ぎゅうにゅう", "gyūnyū", "milk", "ぎゅうにゅうをのみます。", "Gyūnyū o nomimasu.", "I drink milk.", "Cora"),
		jv("しゅくだい", "shukudai", "homework", "しゅくだいがあります。", "Shukudai ga arimasu.", "I have homework.", finch),
		jv("でんしゃ", "densha", "train", "でんしゃはべんりです。", "Densha wa benri desu.", "Trains are convenient.", "Riko"),
	)
	l = addLesson(db, s, "Small っ & Long Vowels", 2, 15,
		char(finch, "A small っ doubles the next consonant with a tiny pause: きて kite (come) vs きって kitte (stamp). Long vowels last twice as long and change meaning: おばさん obasan (aunt) vs おばあさん obāsan (grandmother). Length matters!"),
		mc("Which word means 'stamp'?", "kitte", "きって", "きって", "きて", "きいて", "きっと"),
		mc("Which word means 'grandmother'?", "obāsan", "おばあさん", "おばあさん", "おばさん", "おばさあん", "おぱあさん"),
		listen("Listen and choose", "おばさん", "おばさん (aunt)", "おばさん (aunt)", "おばあさん (grandmother)", "おじさん (uncle)", "おじいさん (grandfather)"),
		listen("Listen and choose", "おじいさん", "おじいさん (grandfather)", "おばさん (aunt)", "おばあさん (grandmother)", "おじさん (uncle)", "おじいさん (grandfather)"),
		mc("How do you read this?", "ちょっと", "chotto (a little)", "chotto (a little)", "choto", "chōtto", "chitto"),
		speak("Mira", "きって、ちょっと、おばあさん"),
	)
	addVocab(db, l,
		jv("きって", "kitte", "postage stamp", "きってをかいます。", "Kitte o kaimasu.", "I'll buy stamps.", finch),
		jv("ちょっと", "chotto", "a little; a moment", "ちょっとまってください。", "Chotto matte kudasai.", "Please wait a moment.", "Lumora"),
		jv("おばあさん", "obāsan", "grandmother; old woman", "おばあさんはげんきです。", "Obāsan wa genki desu.", "Grandma is well.", "Nana"),
		jv("おじさん", "ojisan", "uncle; middle-aged man", "おじさんはやさしいです。", "Ojisan wa yasashii desu.", "My uncle is kind.", "Zephyr"),
	)

	s = addSkillL(db, ja, u, "Katakana", "The script for loanwords, names and emphasis.", "Languages", "#FF5C5C", 4, 24)
	l = addLesson(db, s, "Katakana Basics", 1, 15,
		char(finch, "Katakana writes the same sounds as hiragana, in sharper shapes. It's used for loanwords (コーヒー kōhī, coffee), foreign names, onomatopoeia and emphasis. Watch the look-alikes: シ shi vs ツ tsu, ソ so vs ン n."),
		mc("How do you read this?", "テレビ", "terebi (TV)", "terebi (TV)", "tarebi", "terabi", "toribi"),
		mc("How do you read this?", "アイス", "aisu (ice cream)", "aisu (ice cream)", "aishi", "eisu", "oisu"),
		mc("Which katakana is 'shi'?", "shi — strokes slant upward from the left", "シ", "シ", "ツ", "ソ", "ン"),
		mc("Which katakana is 'n'?", "n — stroke slants upward", "ン", "ン", "ソ", "シ", "ツ"),
		match("Match the loanword", "カメラ", "kamera (camera)", "kamera (camera)", "kamara (comma)", "kemura", "kanera"),
		listen("Listen and choose", "ノート", "ノート (nōto, notebook)", "ノート (nōto, notebook)", "ナート", "ノーテ", "メート"),
		speak("Blaze", "テレビとカメラとノート"),
	)
	addVocab(db, l,
		jv("テレビ", "terebi", "television", "テレビをみます。", "Terebi o mimasu.", "I watch TV.", "Pip"),
		jv("カメラ", "kamera", "camera", "あたらしいカメラです。", "Atarashii kamera desu.", "It's a new camera.", "Zephyr"),
		jv("ノート", "nōto", "notebook", "ノートにかきます。", "Nōto ni kakimasu.", "I'll write it in my notebook.", finch),
		jv("ホテル", "hoteru", "hotel", "ホテルはどこですか。", "Hoteru wa doko desu ka.", "Where is the hotel?", "Riko"),
	)
	l = addLesson(db, s, "Long Vowels & Loanwords", 2, 15,
		char(finch, "In katakana a long vowel is a dash: ー. コーヒー kōhī (coffee), ケーキ kēki (cake). Loanwords are adapted to Japanese sounds, and new combinations exist for foreign sounds: ティ ti (パーティー party), ファ fa (ファン fan), ウィ wi."),
		mc("How do you read this?", "コーヒー", "kōhī (coffee)", "kōhī (coffee)", "kohi", "kōhi", "kohī"),
		mc("What is this loanword?", "パーティー", "party", "party", "parting", "pasta", "patio"),
		mc("What is this loanword?", "コンビニ", "convenience store", "convenience store", "combination", "computer", "concert"),
		mc("What is this loanword?", "スマホ", "smartphone", "smartphone", "smart home", "small house", "smoke"),
		mc("What does ー do in ケーキ?", "ケーキ (kēki)", "Makes the vowel long", "Makes the vowel long", "Doubles the consonant", "Makes it a question", "Marks a pause"),
		speak("Cora", "コーヒーとケーキをください。"),
	)
	addVocab(db, l,
		jv("コーヒー", "kōhī", "coffee", "コーヒーをのみます。", "Kōhī o nomimasu.", "I drink coffee.", "Cora"),
		jv("ケーキ", "kēki", "cake", "ケーキがおいしいです。", "Kēki ga oishii desu.", "The cake is delicious.", "Cora"),
		jv("コンビニ", "konbini", "convenience store", "コンビニはあさまであいています。", "Konbini wa asa made aite imasu.", "The convenience store is open till morning.", "Riko"),
		jv("パーティー", "pātī", "party", "パーティーにいきます。", "Pātī ni ikimasu.", "I'm going to a party.", "Blaze"),
	)

	s = addSkillL(db, ja, u, "Sounds & Pitch Accent", "Morae, devoiced vowels and pitch.", "Music", "#17A3DD", 5, 32)
	l = addLesson(db, s, "Rhythm & Pitch", 1, 15,
		char(finch, "Japanese is mora-timed: every kana beat takes the same time — にほん is ni-ho-n, three beats. Vowels i and u often whisper away: です sounds like 'dess'. And pitch, not stress, distinguishes words: あめ with a HIGH-low pitch is 'rain' (雨); low-HIGH is 'candy' (飴)."),
		mc("How many morae (beats) are in にほん?", "に・ほ・ん", "3", "3", "2", "4", "1"),
		mc("How is です usually pronounced?", "です", "dess (the u is whispered)", "dess (the u is whispered)", "de-su with a strong u", "dez", "deshu"),
		mc("Pitch accent: what does あめ mean with HIGH-low pitch?", "Tokyo pitch: Á-me", "rain (雨)", "rain (雨)", "candy (飴)", "sky (天)", "both"),
		mc("What makes 箸 (chopsticks) and 橋 (bridge) sound different?", "both are はし (hashi)", "Pitch: HÁ-shi vs ha-SHÍ", "Pitch: HÁ-shi vs ha-SHÍ", "Vowel length", "A small っ", "They sound identical"),
		mc("How many morae are in がっこう (school)?", "が・っ・こ・う", "4", "4", "2", "3", "5"),
		speak("Mira", "あめがふっています。"),
	)
	addVocab(db, l,
		jv("雨", "あめ ame (HIGH-low)", "rain", "雨がふっています。", "あめがふっています。Ame ga futte imasu.", "It's raining.", "Riko"),
		jv("飴", "あめ ame (low-HIGH)", "candy", "飴をどうぞ。", "あめをどうぞ。Ame o dōzo.", "Have a candy.", "Pip"),
		jv("箸", "はし hashi (HIGH-low)", "chopsticks", "箸でたべます。", "はしでたべます。Hashi de tabemasu.", "I eat with chopsticks.", "Cora"),
		jv("橋", "はし hashi (low-HIGH)", "bridge", "橋をわたります。", "はしをわたります。Hashi o watarimasu.", "I cross the bridge.", "Zephyr"),
	)

	// ── Survival Japanese ──
	s = addSkillL(db, ja, u, "Greetings & Courtesy", "Hello, thank you, excuse me — and when to bow.", "Hand", "#00C2A8", 6, 40)
	l = addLesson(db, s, "Everyday Greetings", 1, 15,
		char("Lumora", "こんにちは！Greetings change with the time of day: おはようございます in the morning, こんにちは in the daytime, こんばんは in the evening. すみません is the all-purpose 'excuse me / sorry', and a small bow goes a long way."),
		mc("What do you say in the morning (politely)?", "Good morning", "おはようございます", "おはようございます", "こんばんは", "おやすみなさい", "さようなら"),
		mc("What does this mean?", "ありがとうございます", "Thank you very much", "Thank you very much", "Excuse me", "Good night", "You're welcome"),
		mc("What does this mean?", "すみません", "Excuse me / I'm sorry", "Excuse me / I'm sorry", "Thank you", "Goodbye", "Please"),
		mc("What do you say before eating?", "Said before a meal", "いただきます", "いただきます", "ごちそうさまでした", "おかえりなさい", "いってきます"),
		mc("What do you say after eating?", "Said after a meal", "ごちそうさまでした", "ごちそうさまでした", "いただきます", "ただいま", "はじめまして"),
		listen("Listen and choose the meaning", "おやすみなさい", "Good night", "Good night", "Good evening", "Goodbye", "Welcome home"),
		speak("Lumora", "おはようございます。ありがとうございます。"),
	)
	addVocab(db, l,
		jv("こんにちは", "konnichiwa", "hello, good afternoon", "こんにちは、せんせい。", "Konnichiwa, sensei.", "Hello, teacher.", "Lumora"),
		jv("おはようございます", "ohayō gozaimasu", "good morning (polite)", "おはようございます！", "Ohayō gozaimasu!", "Good morning!", "Cora"),
		jv("ありがとうございます", "arigatō gozaimasu", "thank you very much", "どうもありがとうございます。", "Dōmo arigatō gozaimasu.", "Thank you so much.", "Mira"),
		jv("すみません", "sumimasen", "excuse me; sorry", "すみません、えきはどこですか。", "Sumimasen, eki wa doko desu ka.", "Excuse me, where is the station?", "Riko"),
		jv("いただきます", "itadakimasu", "said before eating", "いただきます！", "Itadakimasu!", "Let's eat! (thankfully)", "Cora"),
	)
	l = addLesson(db, s, "Home & Leaving Phrases", 2, 15,
		char("Nana", "Japanese has set phrases for coming and going. Leaving home: いってきます ('I'm off'), answered with いってらっしゃい. Coming home: ただいま ('I'm back'), answered with おかえりなさい. Meeting someone new: はじめまして… よろしくおねがいします."),
		mc("What do you say when you leave home?", "I'm off!", "いってきます", "いってきます", "ただいま", "おかえりなさい", "いってらっしゃい"),
		mc("Someone says ただいま. You reply:", "ただいま！", "おかえりなさい", "おかえりなさい", "いってらっしゃい", "いただきます", "はじめまして"),
		mc("What does よろしくおねがいします express?", "said after introducing yourself", "Please treat me well / I look forward to working with you", "Please treat me well / I look forward to working with you", "Goodbye forever", "Excuse me for interrupting", "Happy birthday"),
		mc("What does どういたしまして mean?", "reply to ありがとう", "You're welcome", "You're welcome", "Thank you", "Excuse me", "See you"),
		speak("Nana", "ただいま！おかえりなさい。"),
	)
	addVocab(db, l,
		jv("いってきます", "ittekimasu", "I'm off (leaving home)", "いってきます！", "Ittekimasu!", "I'm heading out!", "Pip"),
		jv("ただいま", "tadaima", "I'm home", "ただいま！", "Tadaima!", "I'm back!", "Pip"),
		jv("おかえりなさい", "okaerinasai", "welcome home", "おかえりなさい。", "Okaerinasai.", "Welcome back.", "Nana"),
		jv("どういたしまして", "dō itashimashite", "you're welcome", "いいえ、どういたしまして。", "Iie, dō itashimashite.", "Not at all, you're welcome.", "Mira"),
	)

	s = addSkillL(db, ja, u, "Self-Introduction", "はじめまして: name, country, job — with です and の.", "Users", "#9B51E0", 7, 50)
	l = addLesson(db, s, "Hajimemashite", 1, 15,
		char("Lumora", "A polite self-introduction: はじめまして。わたしは ルモラ です。ケニア から きました。がくせい です。よろしく おねがいします。は marks the topic, です is the polite 'is', and の links nouns: にほんごの せんせい = a Japanese teacher."),
		mc("What does this mean?", "わたしはアンナです。", "I am Anna.", "I am Anna.", "I like Anna.", "This is Anna's.", "Anna is here."),
		mc("Complete: I'm from Kenya.", "ケニア＿きました。", "から", "から", "まで", "を", "が"),
		mc("What does の do here?", "にほんごのほん", "Links nouns: a Japanese (language) book", "Links nouns: a Japanese (language) book", "Makes a question", "Marks the subject", "Means 'and'"),
		mc("Make it a question", "がくせいです → ?", "がくせいですか。", "がくせいですか。", "がくせいですね。", "がくせいでした。", "がくせいじゃないです。"),
		mc("Make it negative", "せんせいです (I'm a teacher)", "せんせいじゃないです", "せんせいじゃないです", "せんせいでした", "せんせいですか", "せんせいません"),
		listen("Listen and choose", "はじめまして", "Nice to meet you (first time)", "Nice to meet you (first time)", "Long time no see", "See you later", "Good evening"),
		speak("Lumora", "はじめまして。わたしはルモラです。よろしくおねがいします。"),
	)
	addVocab(db, l,
		jv("はじめまして", "hajimemashite", "nice to meet you", "はじめまして、ミラです。", "Hajimemashite, Mira desu.", "Nice to meet you, I'm Mira.", "Mira"),
		jv("がくせい", "gakusei", "student", "わたしはがくせいです。", "Watashi wa gakusei desu.", "I'm a student.", "Pip"),
		jv("せんせい", "sensei", "teacher", "たなかせんせいです。", "Tanaka sensei desu.", "This is Professor Tanaka.", finch),
		jv("〜から きました", "〜kara kimashita", "I'm from ~", "ケニアからきました。", "Kenia kara kimashita.", "I'm from Kenya.", "Lumora"),
		jv("〜じん", "〜jin", "person from ~ (nationality)", "わたしはケニアじんです。", "Watashi wa Keniajin desu.", "I'm Kenyan.", "Lumora"),
	)

	s = addSkillL(db, ja, u, "Numbers, Counters & Money", "Count to 10,000, use counters, ask how much.", "Hash", "#FF5C5C", 8, 60)
	l = addLesson(db, s, "Numbers & Prices", 1, 15,
		char(finch, "いち、に、さん、よん（し）、ご、ろく、なな（しち）、はち、きゅう（く）、じゅう. Bigger numbers stack: じゅうご 15, にじゅう 20, ひゃく 100, せん 1,000, いちまん 10,000. Watch the sound changes: さんびゃく 300, はっぴゃく 800, さんぜん 3,000."),
		mc("How do you say 300?", "300", "さんびゃく", "さんびゃく", "さんひゃく", "さんぴゃく", "みひゃく"),
		mc("How do you say 800?", "800", "はっぴゃく", "はっぴゃく", "はちひゃく", "はちびゃく", "やひゃく"),
		mc("What number is this?", "にせんごひゃく", "2,500", "2,500", "25,000", "250", "2,050"),
		mc("How do you ask the price?", "How much is it?", "いくらですか。", "いくらですか。", "なんですか。", "どこですか。", "いつですか。"),
		listen("Listen and choose the price", "せんにひゃくえん", "¥1,200", "¥1,200", "¥2,100", "¥12,000", "¥1,020"),
		speak("Cora", "これはいくらですか。せんえんです。"),
	)
	addVocab(db, l,
		jv("いくら", "ikura", "how much", "これはいくらですか。", "Kore wa ikura desu ka.", "How much is this?", "Cora"),
		jv("百", "ひゃく hyaku", "hundred", "百えんです。", "ひゃくえんです。Hyaku en desu.", "It's 100 yen.", finch),
		jv("千", "せん sen", "thousand", "千えんさつ。", "せんえんさつ。Sen en satsu.", "A 1,000-yen note.", finch),
		jv("万", "まん man", "ten thousand", "一万えんです。", "いちまんえんです。Ichiman en desu.", "It's 10,000 yen.", finch),
	)
	l = addLesson(db, s, "Counters", 2, 15,
		char(finch, "Japanese counts things with counters that depend on shape or kind. The general set ひとつ, ふたつ, みっつ… works for most objects. 〜にん counts people (ひとり, ふたり, さんにん), 〜まい flat things (paper, tickets), 〜ほん long things (pens, bottles), 〜ひき small animals."),
		mc("How do you say 'two people'?", "2 people", "ふたり", "ふたり", "ににん", "ふたつ", "にほん"),
		mc("Which counter for tickets (flat things)?", "きっぷを2＿ください。", "まい", "まい", "ほん", "ひき", "にん"),
		mc("Which counter for bottles (long things)?", "ビールを3＿", "ぼん（ほん）", "ぼん（ほん）", "まい", "ひき", "にん"),
		mc("Ordering three items with the general counter", "Three, please.", "みっつください。", "みっつください。", "さんください。", "さんまいください。", "さんにんください。"),
		mc("How do you say 'one person'?", "1 person", "ひとり", "ひとり", "いちにん", "ひとつ", "いっぽん"),
		speak("Cora", "コーヒーをふたつください。"),
	)
	addVocab(db, l,
		jv("ひとつ", "hitotsu", "one (thing)", "りんごをひとつください。", "Ringo o hitotsu kudasai.", "One apple, please.", "Cora"),
		jv("ふたり", "futari", "two people", "ふたりです。", "Futari desu.", "Two people (a table for two).", "Riko"),
		jv("〜まい", "〜mai", "counter for flat things", "きっぷをにまいください。", "Kippu o nimai kudasai.", "Two tickets, please.", "Riko"),
		jv("〜ほん", "〜hon / bon / pon", "counter for long things", "ペンがいっぽんあります。", "Pen ga ippon arimasu.", "There's one pen.", finch),
	)

	s = addSkillL(db, ja, u, "Time, Days & Dates", "Clock time, days of the week and when things happen.", "Clock", "#17A3DD", 9, 70)
	l = addLesson(db, s, "What Time Is It?", 1, 15,
		char(finch, "Time: 〜じ for o'clock, 〜ふん/ぷん for minutes. いま なんじ ですか — what time is it now? Watch: よじ 4:00, しちじ 7:00, くじ 9:00, and はん for half past. Days end in 〜ようび: げつようび Monday … にちようび Sunday."),
		mc("How do you say 4 o'clock?", "4:00", "よじ", "よじ", "しじ", "よんじ", "よっじ"),
		mc("How do you say 9:30?", "9:30", "くじはん", "くじはん", "きゅうじはん", "くじさんじゅう", "きゅうじはんぷん"),
		mc("What day is this?", "きんようび", "Friday", "Friday", "Monday", "Wednesday", "Sunday"),
		mc("What does this ask?", "いまなんじですか。", "What time is it now?", "What time is it now?", "What day is it?", "Where are you now?", "How much is it?"),
		listen("Listen and choose", "げつようびのしちじ", "Monday at 7:00", "Monday at 7:00", "Monday at 4:00", "Sunday at 7:00", "Thursday at 7:00"),
		speak("Riko", "いま、さんじはんです。"),
	)
	addVocab(db, l,
		jv("いま", "ima", "now", "いまなんじですか。", "Ima nanji desu ka.", "What time is it now?", "Riko"),
		jv("〜じ", "〜ji", "o'clock", "しちじにおきます。", "Shichiji ni okimasu.", "I get up at seven.", "Riko"),
		jv("はん", "han", "half (past)", "くじはんです。", "Kuji han desu.", "It's 9:30.", "Riko"),
		jv("〜ようび", "〜yōbi", "day of the week", "なんようびですか。", "Nan'yōbi desu ka.", "What day is it?", finch),
		jv("きょう", "kyō", "today", "きょうはげつようびです。", "Kyō wa getsuyōbi desu.", "Today is Monday.", "Cora"),
	)

	// ── Core grammar ──
	s = addSkillL(db, ja, u, "Core Particles", "は・が・を・に・で・へ・と・も — the glue of Japanese.", "Link", "#6C3FC5", 10, 80)
	l = addLesson(db, s, "Topic, Subject & Object", 1, 18,
		char(finch, "Particles follow the word they mark. は (wa) sets the topic: 'as for…'. が marks the subject, often new information or with likes/abilities (すしがすきです). を marks the direct object (パンをたべます). も means 'also' and replaces は/が/を."),
		mc("Choose the particle", "わたし＿がくせいです。(As for me, I'm a student.)", "は", "は", "を", "で", "へ"),
		mc("Choose the particle", "みず＿のみます。(I drink water.)", "を", "を", "は", "に", "で"),
		mc("Choose the particle", "すし＿すきです。(I like sushi.)", "が", "が", "を", "で", "へ"),
		mc("Choose the particle", "わたし＿がくせいです。(I'm a student too.)", "も", "も", "が", "を", "と"),
		mc("Which sentence answers 'Who came?'", "だれがきましたか。", "たなかさんがきました。", "たなかさんがきました。", "たなかさんはきました。", "たなかさんをきました。", "たなかさんできました。"),
		speak("Mira", "わたしはコーヒーがすきです。"),
	)
	addVocab(db, l,
		jv("は", "wa", "topic particle (as for…)", "これはほんです。", "Kore wa hon desu.", "This is a book.", finch),
		jv("が", "ga", "subject particle", "ねこがすきです。", "Neko ga suki desu.", "I like cats.", finch),
		jv("を", "o", "object particle", "テレビをみます。", "Terebi o mimasu.", "I watch TV.", finch),
		jv("も", "mo", "also, too", "わたしもいきます。", "Watashi mo ikimasu.", "I'll go too.", "Blaze"),
	)
	l = addLesson(db, s, "Place, Time, Means & Company", 2, 18,
		char(finch, "に marks a destination, a point in time, or where something exists: がっこうにいきます, しちじにおきます. で marks where an action happens or the means: としょかんでべんきょうします, バスでいきます. へ (e) marks direction. と means 'with' or 'and' between nouns."),
		mc("Choose the particle", "しちじ＿おきます。(I get up at 7.)", "に", "に", "で", "を", "が"),
		mc("Choose the particle", "としょかん＿べんきょうします。(I study at the library.)", "で", "で", "に", "を", "へ"),
		mc("Choose the particle", "バス＿いきます。(I go by bus.)", "で", "で", "に", "と", "を"),
		mc("Choose the particle", "ともだち＿えいがをみます。(I watch a film with a friend.)", "と", "と", "で", "を", "は"),
		mc("Choose the particle", "にほん＿いきます。(I'm going to Japan.)", "に／へ", "に／へ", "で", "を", "が"),
		speak("Blaze", "あした、ともだちとでんしゃできょうとへいきます。"),
	)
	addVocab(db, l,
		jv("に", "ni", "to; at (time); in (existence)", "がっこうにいきます。", "Gakkō ni ikimasu.", "I go to school.", finch),
		jv("で", "de", "at (action place); by (means)", "うちでたべます。", "Uchi de tabemasu.", "I eat at home.", finch),
		jv("と", "to", "with; and", "ははとはなします。", "Haha to hanashimasu.", "I talk with my mother.", "Nana"),
		jv("としょかん", "toshokan", "library", "としょかんでべんきょうします。", "Toshokan de benkyō shimasu.", "I study at the library.", "Mira"),
	)

	s = addSkillL(db, ja, u, "Verbs: the ます Form", "Polite present, past, negative and 'let's'.", "CheckCircle", "#00C2A8", 11, 92)
	l = addLesson(db, s, "Present & Negative", 1, 18,
		char(finch, "Polite verbs end in ます. There's no future tense — たべます means 'I eat' or 'I will eat'. Negative: ません (たべません). Questions add か. Japanese word order is subject–object–VERB: the verb always comes last."),
		mc("What does this mean?", "あしたテニスをします。", "I'll play tennis tomorrow.", "I'll play tennis tomorrow.", "I played tennis yesterday.", "I don't play tennis.", "Let's play tennis."),
		mc("Make it negative", "のみます (drink)", "のみません", "のみません", "のみました", "のみましょう", "のまないです"),
		mc("Put it in order", "I read a book. (ほん / を / よみます / わたしは)", "わたしはほんをよみます。", "わたしはほんをよみます。", "わたしはよみますほんを。", "ほんをわたしはよみますを。", "よみますわたしはほんを。"),
		mc("What does this mean?", "おさけをのみません。", "I don't drink alcohol.", "I don't drink alcohol.", "I drank alcohol.", "Let's drink.", "I want to drink."),
		mc("Choose the verb", "まいにち、にほんごを＿。(I study Japanese every day.)", "べんきょうします", "べんきょうします", "べんきょうでした", "べんきょうです", "べんきょうましょう"),
		speak("Cora", "まいにちにほんごをべんきょうします。"),
	)
	addVocab(db, l,
		jv("たべます", "tabemasu", "to eat (polite)", "あさごはんをたべます。", "Asagohan o tabemasu.", "I eat breakfast.", "Cora"),
		jv("のみます", "nomimasu", "to drink (polite)", "おちゃをのみます。", "Ocha o nomimasu.", "I drink tea.", "Nana"),
		jv("します", "shimasu", "to do (polite)", "しゅくだいをします。", "Shukudai o shimasu.", "I do my homework.", "Pip"),
		jv("べんきょう", "benkyō", "study", "にほんごのべんきょう。", "Nihongo no benkyō.", "Studying Japanese.", finch),
		jv("まいにち", "mainichi", "every day", "まいにちはしります。", "Mainichi hashirimasu.", "I run every day.", "Blaze"),
	)
	l = addLesson(db, s, "Past & 'Let's'", 2, 18,
		char(finch, "Past: ました (たべました, I ate). Past negative: ませんでした (たべませんでした). Invitations: 〜ませんか (won't you…?) and 〜ましょう (let's…)."),
		mc("Put it in the past", "いきます (go)", "いきました", "いきました", "いきません", "いきましょう", "いきませんか"),
		mc("Past negative", "みます (watch)", "みませんでした", "みませんでした", "みましたでした", "みません", "みないでした"),
		mc("What does this mean?", "いっしょにたべませんか。", "Would you like to eat together?", "Would you like to eat together?", "We didn't eat together.", "Don't eat together.", "I ate alone."),
		mc("What does this mean?", "きのうえいがをみました。", "I watched a film yesterday.", "I watched a film yesterday.", "I'll watch a film tomorrow.", "I don't watch films.", "Let's watch a film."),
		mc("Accept an invitation", "えいがをみませんか。", "ええ、みましょう。", "ええ、みましょう。", "いいえ、みました。", "ええ、みません。", "はい、みませんでした。"),
		speak("Blaze", "きのうはどこへいきましたか。"),
	)
	addVocab(db, l,
		jv("きのう", "kinō", "yesterday", "きのうはやすみでした。", "Kinō wa yasumi deshita.", "Yesterday was a day off.", "Cora"),
		jv("いっしょに", "issho ni", "together", "いっしょにいきましょう。", "Issho ni ikimashō.", "Let's go together.", "Blaze"),
		jv("えいが", "eiga", "film, movie", "えいがをみました。", "Eiga o mimashita.", "I watched a film.", "Zephyr"),
		jv("みます", "mimasu", "to see, watch (polite)", "テレビをみます。", "Terebi o mimasu.", "I watch TV.", "Pip"),
	)

	s = addSkillL(db, ja, u, "い- and な-Adjectives", "Describe things — and conjugate the adjectives.", "Sparkles", "#F5A623", 12, 104)
	l = addLesson(db, s, "Two Kinds of Adjectives", 1, 18,
		char(finch, "い-adjectives end in い and conjugate themselves: たかい (expensive) → たかくない (not) → たかかった (was) → たかくなかった. な-adjectives take な before a noun (しずかなまち, a quiet town) and conjugate with です: しずかじゃないです, しずかでした. いい is irregular: よくない, よかった."),
		mc("Negative of たかい", "たかい (expensive)", "たかくない", "たかくない", "たかいじゃない", "たかないです", "たかくなかった"),
		mc("Past of おいしい", "おいしい (tasty)", "おいしかった", "おいしかった", "おいしいでした", "おいしくた", "おいしだった"),
		mc("Past of いい (good)", "いい", "よかった", "よかった", "いかった", "いいでした", "よくった"),
		mc("Choose the right form", "＿まち (a quiet town)", "しずかなまち", "しずかなまち", "しずかいまち", "しずかのまち", "しずかくまち"),
		mc("Negative of きれい (pretty — a な-adjective!)", "きれいです", "きれいじゃないです", "きれいじゃないです", "きれくないです", "きれいくないです", "きれいないです"),
		speak("Mira", "このラーメンはとてもおいしかったです。"),
	)
	addVocab(db, l,
		jv("たかい", "takai", "expensive; tall", "このかばんはたかいです。", "Kono kaban wa takai desu.", "This bag is expensive.", "Cora"),
		jv("やすい", "yasui", "cheap", "このみせはやすいです。", "Kono mise wa yasui desu.", "This shop is cheap.", "Cora"),
		jv("おいしい", "oishii", "delicious", "すしはおいしいです。", "Sushi wa oishii desu.", "Sushi is delicious.", "Cora"),
		jv("しずか（な）", "shizuka (na)", "quiet", "しずかなへやです。", "Shizuka na heya desu.", "It's a quiet room.", "Nana"),
		jv("きれい（な）", "kirei (na)", "pretty; clean", "きれいなはなですね。", "Kirei na hana desu ne.", "What pretty flowers.", "Mira"),
	)

	s = addSkillL(db, ja, u, "Existence & Location", "ある / いる, こ・そ・あ・ど words, where things are.", "Compass", "#17A3DD", 13, 116)
	l = addLesson(db, s, "Where Is It?", 1, 18,
		char(finch, "'There is' has two verbs: あります for things and plants, います for people and animals. Location words come in sets: ここ (here), そこ (there, near you), あそこ (over there), どこ (where?). Likewise これ/それ/あれ/どれ for this/that/that over there/which."),
		mc("Choose the verb", "こうえんにいぬが＿。(There's a dog in the park.)", "います", "います", "あります", "です", "します"),
		mc("Choose the verb", "つくえのうえにほんが＿。(There's a book on the desk.)", "あります", "あります", "います", "します", "いきます"),
		mc("What does this ask?", "トイレはどこですか。", "Where is the toilet?", "Where is the toilet?", "Is this the toilet?", "Which toilet?", "Is there a toilet?"),
		mc("'That one over there (far from both of us)'", "that over there", "あれ", "あれ", "それ", "これ", "どれ"),
		mc("What does this mean?", "えきのまえにぎんこうがあります。", "There's a bank in front of the station.", "There's a bank in front of the station.", "The station is behind the bank.", "The bank is inside the station.", "There's no bank near the station."),
		speak("Riko", "すみません、えきはどこですか。"),
	)
	addVocab(db, l,
		jv("あります", "arimasu", "there is (things)", "コンビニがあります。", "Konbini ga arimasu.", "There's a convenience store.", "Riko"),
		jv("います", "imasu", "there is (living beings)", "ねこがいます。", "Neko ga imasu.", "There's a cat.", "Pip"),
		jv("どこ", "doko", "where", "えきはどこですか。", "Eki wa doko desu ka.", "Where is the station?", "Riko"),
		jv("まえ", "mae", "in front; before", "えきのまえでまちます。", "Eki no mae de machimasu.", "I'll wait in front of the station.", "Zephyr"),
		jv("えき", "eki", "station", "えきはちかいです。", "Eki wa chikai desu.", "The station is close.", "Riko"),
	)

	// ── First kanji ──
	s = addSkillL(db, ja, u, "First Kanji (N5)", "Numbers, days, people and nature — about 100 kanji at N5.", "PenLine", "#FF5C5C", 14, 128)
	l = addLesson(db, s, "Numbers & the Days of the Week", 1, 18,
		char(finch, "Kanji carry meaning; most have a Chinese-derived on'yomi and a native kun'yomi reading. The week is built from nature: 月 moon (Monday 月曜日), 火 fire, 水 water, 木 tree, 金 gold, 土 earth, 日 sun (Sunday 日曜日). Numbers: 一 二 三 四 五 六 七 八 九 十."),
		mc("What does 水 mean?", "水", "water", "water", "fire", "tree", "gold"),
		mc("Which day is 木曜日?", "木曜日", "Thursday", "Thursday", "Tuesday", "Friday", "Saturday"),
		mc("How do you read 三?", "三", "さん", "さん", "し", "みっ", "ご"),
		mc("What does 火 mean?", "火 (as in 火曜日 Tuesday)", "fire", "fire", "water", "sun", "earth"),
		mc("Which kanji means 'sun / day'?", "sun, day", "日", "日", "月", "目", "白"),
		mc("What day is 金曜日?", "金曜日", "Friday", "Friday", "Monday", "Wednesday", "Sunday"),
		speak("Mira", "月曜日から金曜日まではたらきます。"),
	)
	addVocab(db, l,
		jv("日", "ひ・にち hi / nichi", "sun; day", "日曜日はやすみです。", "にちようびはやすみです。", "Sunday is a day off.", finch),
		jv("月", "つき・げつ tsuki / getsu", "moon; month", "月曜日です。", "げつようびです。", "It's Monday.", finch),
		jv("水", "みず・すい mizu / sui", "water", "水をください。", "みずをください。", "Water, please.", "Cora"),
		jv("火", "ひ・か hi / ka", "fire", "火曜日にあいましょう。", "かようびにあいましょう。", "Let's meet on Tuesday.", finch),
		jv("木", "き・もく ki / moku", "tree", "大きい木があります。", "おおきいきがあります。", "There's a big tree.", "Zephyr"),
	)
	l = addLesson(db, s, "People, Nature & Size", 2, 18,
		char(finch, "Many kanji began as pictures: 人 a person, 山 a mountain, 川 a river, 口 a mouth, 目 an eye. 大 (big) is a person with arms spread; 小 (small) and 上 (up) / 下 (down) point the way. Combine them: 大人 おとな, adult; 日本 にほん, Japan ('sun origin')."),
		mc("What does 山 mean?", "山", "mountain", "mountain", "river", "field", "tree"),
		mc("What does 大 mean?", "大きい (おおきい)", "big", "big", "small", "up", "person"),
		mc("How is 日本 read?", "日本", "にほん", "にほん", "ひもと", "にちほん", "じつほん"),
		mc("What does 上 mean?", "上", "up, above", "up, above", "down, below", "middle", "inside"),
		mc("What does 川 mean?", "川", "river", "river", "mountain", "rice field", "rain"),
		speak("Zephyr", "山の上に大きい木があります。"),
	)
	addVocab(db, l,
		jv("人", "ひと・じん・にん hito / jin / nin", "person", "あの人はだれですか。", "あのひとはだれですか。", "Who is that person?", finch),
		jv("山", "やま・さん yama / san", "mountain", "富士山はたかいです。", "ふじさんはたかいです。", "Mount Fuji is tall.", "Zephyr"),
		jv("大きい", "おおきい ōkii", "big", "大きいいえです。", "おおきいいえです。", "It's a big house.", "Blaze"),
		jv("小さい", "ちいさい chiisai", "small", "小さいねこです。", "ちいさいねこです。", "It's a small cat.", "Pip"),
		jv("日本", "にほん nihon", "Japan", "日本へいきたいです。", "にほんへいきたいです。", "I want to go to Japan.", "Lumora"),
	)

	s = addSkillL(db, ja, u, "Food, Shopping & Ordering", "Order at a restaurant, shop, and pay.", "Coffee", "#00C2A8", 15, 142)
	l = addLesson(db, s, "At the Restaurant", 1, 18,
		char("Cora", "いらっしゃいませ！Order with 〜をください or the softer 〜をおねがいします. Ask for the bill with おかいけいおねがいします. In Japan there's no tipping — great service is simply expected."),
		mc("Order one coffee", "Coffee, please.", "コーヒーをひとつおねがいします。", "コーヒーをひとつおねがいします。", "コーヒーがひとつです。", "コーヒーをのみましょう。", "コーヒーはいくらでしたか。"),
		mc("What does the staff mean?", "いらっしゃいませ", "Welcome (to our shop)", "Welcome (to our shop)", "Goodbye", "Thank you for waiting", "Excuse me"),
		mc("How do you ask for the bill?", "The bill, please.", "おかいけいおねがいします。", "おかいけいおねがいします。", "メニューをください。", "いただきます。", "ごちそうさまでした。"),
		mc("What does this mean?", "これはなんですか。", "What is this?", "What is this?", "How much is this?", "Where is this?", "Is this delicious?"),
		listen("Listen and choose", "みずをください", "Water, please.", "Water, please.", "Tea, please.", "The menu, please.", "Rice, please."),
		speak("Cora", "すみません、ラーメンをひとつおねがいします。"),
	)
	addVocab(db, l,
		jv("〜をください", "〜o kudasai", "please give me ~", "メニューをください。", "Menyū o kudasai.", "The menu, please.", "Cora"),
		jv("おねがいします", "onegaishimasu", "please (request)", "おかいけいおねがいします。", "Okaikei onegaishimasu.", "The bill, please.", "Cora"),
		jv("ごはん", "gohan", "rice; meal", "ごはんをたべましょう。", "Gohan o tabemashō.", "Let's eat.", "Nana"),
		jv("みせ", "mise", "shop", "あのみせはやすいです。", "Ano mise wa yasui desu.", "That shop is cheap.", "Cora"),
	)

	s = addSkillL(db, ja, u, "The て-Form", "Requests, ongoing actions and linking verbs.", "Layers", "#6C3FC5", 16, 156)
	l = addLesson(db, s, "Making the て-Form", 1, 20,
		char(finch, "The て-form joins and requests. Group 2 (る-verbs): drop る, add て — たべる → たべて. Group 1 depends on the ending: う/つ/る → って (かう → かって), む/ぶ/ぬ → んで (よむ → よんで), く → いて (かく → かいて), ぐ → いで, す → して. Irregular: する → して, くる → きて, and いく → いって."),
		mc("て-form of たべる", "たべる (eat)", "たべて", "たべて", "たべって", "たべいて", "たべんで"),
		mc("て-form of よむ", "よむ (read)", "よんで", "よんで", "よみて", "よって", "よいて"),
		mc("て-form of かく", "かく (write)", "かいて", "かいて", "かって", "かきて", "かんで"),
		mc("て-form of いく (irregular!)", "いく (go)", "いって", "いって", "いいて", "いきて", "いんで"),
		mc("て-form of はなす", "はなす (speak)", "はなして", "はなして", "はなって", "はないて", "はなんで"),
		mc("て-form of くる", "くる (come)", "きて", "きて", "くって", "くて", "こて"),
		speak("Blaze", "たべて、のんで、よんで、かいて。"),
	)
	addVocab(db, l,
		jv("たべる", "taberu", "to eat (dictionary form)", "なにをたべる？", "Nani o taberu?", "What will you eat?", "Cora"),
		jv("よむ", "yomu", "to read", "まんがをよむ。", "Manga o yomu.", "I read manga.", "Pip"),
		jv("かく", "kaku", "to write", "てがみをかく。", "Tegami o kaku.", "I write a letter.", finch),
		jv("はなす", "hanasu", "to speak", "にほんごをはなす。", "Nihongo o hanasu.", "I speak Japanese.", "Lumora"),
	)
	l = addLesson(db, s, "Please… & -ing", 2, 20,
		char(finch, "〜てください makes a polite request: まってください (please wait). 〜ています describes an action in progress or a state: いまたべています (I'm eating now), けっこんしています (I'm married). Linking: あさおきて、シャワーをあびて、でかけます — I get up, shower and go out."),
		mc("Please speak slowly.", "ゆっくり＿ください。", "はなして", "はなして", "はなし", "はなす", "はなさない"),
		mc("What does this mean?", "いま、べんきょうしています。", "I'm studying now.", "I'm studying now.", "I studied before.", "I'll study now.", "Please study now."),
		mc("What does this mean?", "あにはとうきょうにすんでいます。", "My older brother lives in Tokyo.", "My older brother lives in Tokyo.", "My brother will move to Tokyo.", "My brother lived in Tokyo.", "My brother visits Tokyo."),
		mc("Please write your name here.", "ここになまえを＿ください。", "かいて", "かいて", "かきて", "かって", "かく"),
		mc("What does this mean?", "うちへかえって、ねました。", "I went home and slept.", "I went home and slept.", "I'll go home and sleep.", "I slept, then went home.", "Please go home and sleep."),
		speak("Mira", "すみません、もういちどいってください。"),
	)
	addVocab(db, l,
		jv("まってください", "matte kudasai", "please wait", "ちょっとまってください。", "Chotto matte kudasai.", "Please wait a moment.", "Lumora"),
		jv("ゆっくり", "yukkuri", "slowly", "ゆっくりはなしてください。", "Yukkuri hanashite kudasai.", "Please speak slowly.", "Lumora"),
		jv("すんでいます", "sunde imasu", "live (reside)", "おおさかにすんでいます。", "Ōsaka ni sunde imasu.", "I live in Osaka.", "Riko"),
		jv("もういちど", "mō ichido", "once more", "もういちどおねがいします。", "Mō ichido onegaishimasu.", "Once more, please.", "Mira"),
	)
}
