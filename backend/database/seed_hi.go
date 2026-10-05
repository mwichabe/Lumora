package database

import (
	"gorm.io/gorm"

	"lumora/backend/models"
)

// ===== Hindi (हिन्दी) course, A1 → C2 =====
//
// Standard Hindi in Devanagari, CEFR-aligned. Following the roadmap's first
// rule — learn the real script from day one — A1 opens with three Devanagari
// skills (vowels, consonants with aspiration and retroflexes, then matras and
// conjuncts) before anything else. Gender is taught with every noun and
// practised through agreement immediately; postpositions and the oblique case
// come early; आप (polite you) is used from the first lesson.
//
// Romanization is a crutch, not a path: vocabulary always shows a simplified
// transliteration next to the Devanagari (as a pronunciation aid), but
// exercises drop it from B1 on so learners read the script itself.
//
// Like the other non-Latin courses, lessons are answered by choosing — a
// learner can't be assumed to have a Devanagari keyboard set up — with typed
// free-writing tasks from A2 up. Every exercise carries its own options.

// hv builds a Hindi vocabulary item: the Devanagari is what's spoken aloud;
// the transliteration rides along in the translation.
func hv(word, roman, english, example, exampleRoman, exampleEnglish, speaker string) models.VocabItem {
	return vw(word, roman+" · "+english, example, exampleRoman+" — "+exampleEnglish, speaker)
}

func seedHindi(db *gorm.DB) {
	seedHindiA1(db)
	seedHindiA2(db)
	seedHindiB1(db)
	seedHindiB2(db)
	seedHindiC1(db)
	seedHindiC2(db)
	seedHindiListening(db)
	seedHindiReading(db)
}

// ───────────────────── A1 — शुरुआत: Foundations & Survival Hindi ─────────────────────

func seedHindiA1(db *gorm.DB) {
	const hi = "hi"
	const u = "A1 · शुरुआत — Foundations & Survival Hindi"
	finch := "Professor Finch"

	s := addSkillL(db, hi, u, "Devanagari I: Vowels", "The script from day one: independent vowels and the shirorekha line.", "Languages", "#6C3FC5", 1, 0)
	l := addLesson(db, s, "स्वर — Vowels", 1, 15,
		char(finch, "नमस्ते! Hindi is written in Devanagari, left to right, with a line (shirorekha) along the top. Vowels come in short–long pairs: अ a, आ ā, इ i, ई ī, उ u, ऊ ū, then ए e, ऐ ai, ओ o, औ au. Length matters: the long vowels are held twice as long. ं (anusvara) adds a nasal: अं an."),
		listen("Listen and choose the vowel", "आ", "आ (ā — long)", "आ (ā — long)", "अ (a — short)", "ओ (o)", "औ (au)"),
		listen("Listen and choose the vowel", "ई", "ई (ī — long)", "ई (ī — long)", "इ (i — short)", "ए (e)", "ऐ (ai)"),
		mc("Which is the short 'u'?", "u (as in 'put')", "उ", "उ", "ऊ", "ओ", "अ"),
		mc("How is ऐ pronounced?", "ऐ", "ai (like 'a' in 'bat', drawn out)", "ai (like 'a' in 'bat', drawn out)", "ī", "o", "u"),
		mc("What is the line along the top of Devanagari called?", "शिरोरेखा", "shirorekha", "shirorekha", "matra", "virama", "danda"),
		match("Match the pair", "इ — ई", "short i — long ī", "short i — long ī", "a — ā", "u — ū", "e — ai"),
		speak("Mira", "अ आ इ ई उ ऊ ए ऐ ओ औ"),
	)
	addVocab(db, l,
		hv("आम", "ām", "mango (m)", "यह आम मीठा है।", "yah ām mīṭhā hai.", "This mango is sweet.", "Cora"),
		hv("ईख", "īkh", "sugarcane (f)", "ईख मीठी होती है।", "īkh mīṭhī hotī hai.", "Sugarcane is sweet.", "Pip"),
		hv("ऊन", "ūn", "wool (m)", "यह ऊन गरम है।", "yah ūn garam hai.", "This wool is warm.", "Nana"),
		hv("एक", "ek", "one", "एक आम दीजिए।", "ek ām dījie.", "Please give me one mango.", "Cora"),
		hv("और", "aur", "and", "चाय और पानी।", "chāy aur pānī.", "Tea and water.", "Lumora"),
	)

	s = addSkillL(db, hi, u, "Devanagari II: Consonants", "Consonants by place — with aspiration and retroflexes.", "Languages", "#00C2A8", 2, 10)
	l = addLesson(db, s, "व्यंजन — Consonants", 1, 15,
		char(finch, "Every consonant carries a built-in 'a': क is ka. They're grouped by where the mouth makes them: क ख ग घ ङ (throat), च छ ज झ ञ (palate), ट ठ ड ढ ण (retroflex — tongue curled back), त थ द ध न (teeth), प फ ब भ म (lips), then य र ल व श ष स ह. Aspiration changes meaning: the second of each pair has a puff of air — कल kal (tomorrow) vs खल khal (villain)."),
		mc("Which consonant is aspirated?", "puff of air", "ख (kha)", "ख (kha)", "क (ka)", "ग (ga)", "न (na)"),
		mc("Which is the retroflex 't' (tongue curled back)?", "retroflex", "ट (ṭa)", "ट (ṭa)", "त (ta)", "थ (tha)", "द (da)"),
		listen("Listen and choose", "फल", "फल (phal — fruit)", "फल (phal — fruit)", "पल (pal — moment)", "बल (bal — strength)", "मल (mal — dirt)"),
		mc("What vowel does every bare consonant carry?", "क = k + ?", "a", "a", "i", "u", "none"),
		mc("Which letter is 'ma'?", "lips group", "म", "म", "न", "ब", "भ"),
		listen("Listen and choose", "घर", "घर (ghar — house)", "घर (ghar — house)", "गर (gar)", "कर (kar — do)", "खर (khar)"),
		speak("Blaze", "क ख ग घ, त थ द ध, ट ठ ड ढ"),
	)
	addVocab(db, l,
		hv("घर", "ghar", "house, home (m)", "मेरा घर बड़ा है।", "merā ghar baṛā hai.", "My house is big.", "Nana"),
		hv("फल", "phal", "fruit (m)", "फल ताज़ा है।", "phal tāzā hai.", "The fruit is fresh.", "Cora"),
		hv("कल", "kal", "tomorrow; yesterday", "कल मिलेंगे।", "kal milenge.", "See you tomorrow.", "Lumora"),
		hv("नमक", "namak", "salt (m)", "थोड़ा नमक दीजिए।", "thoṛā namak dījie.", "Please give a little salt.", "Cora"),
	)
	l = addLesson(db, s, "Tricky Sounds", 2, 15,
		char(finch, "Some sounds need special care. ड़ (ṛa) is a flapped retroflex: पेड़ peṛ (tree), लड़का laṛkā (boy). A dot below (nuqta) marks sounds from Persian, Arabic and English: ज़ za (ज़रूर, certainly), फ़ fa (फ़ोन). श sha and ष ṣa both sound like 'sh' today; स is 's'."),
		mc("Which word means 'boy'?", "laṛkā", "लड़का", "लड़का", "लडका", "लरका", "लकड़ा"),
		mc("What does the dot (nuqta) in ज़ mark?", "ज़रूर", "A sound from Persian/Arabic/English (z)", "A sound from Persian/Arabic/English (z)", "A long vowel", "Aspiration", "A silent letter"),
		mc("How is फ़ in फ़ोन pronounced?", "फ़ोन", "fa (as in 'phone')", "fa (as in 'phone')", "pha (aspirated p)", "ba", "va"),
		listen("Listen and choose", "पेड़", "पेड़ (peṛ — tree)", "पेड़ (peṛ — tree)", "पेट (peṭ — stomach)", "पैर (pair — foot)", "पेन (pen)"),
		mc("Which letter is 'sa' (plain s)?", "s", "स", "स", "श", "ष", "ह"),
		speak("Pip", "लड़का पेड़ के नीचे है।"),
	)
	addVocab(db, l,
		hv("लड़का", "laṛkā", "boy (m)", "लड़का स्कूल जाता है।", "laṛkā skūl jātā hai.", "The boy goes to school.", "Pip"),
		hv("पेड़", "peṛ", "tree (m)", "यह पेड़ बहुत ऊँचा है।", "yah peṛ bahut ū̃chā hai.", "This tree is very tall.", "Zephyr"),
		hv("ज़रूर", "zarūr", "certainly, of course", "ज़रूर आइए!", "zarūr āie!", "Do come!", "Lumora"),
		hv("फ़ोन", "fon", "phone (m)", "मेरा फ़ोन कहाँ है?", "merā fon kahā̃ hai?", "Where is my phone?", "Blaze"),
	)

	s = addSkillL(db, hi, u, "Devanagari III: Matras & Conjuncts", "Vowel signs on consonants, nasalization and joined letters.", "PenLine", "#F5A623", 3, 20)
	l = addLesson(db, s, "मात्राएँ — Vowel Signs", 1, 15,
		char(finch, "After a consonant, vowels become signs (matras): का kā, कि ki (written before the letter!), की kī, कु ku, कू kū, के ke, कै kai, को ko, कौ kau. Nasal vowels take ँ (chandrabindu) or ं: माँ mā̃ (mother), हाँ hā̃ (yes), में mẽ (in)."),
		mc("How do you read कि?", "क + ि", "ki", "ki", "kī", "ke", "ku"),
		mc("How do you read की?", "क + ी", "kī", "kī", "ki", "kai", "ko"),
		mc("Which spelling is 'kitāb' (book)?", "kitāb", "किताब", "किताब", "कीताब", "कितब", "केताब"),
		mc("What does ँ in माँ do?", "माँ", "Nasalizes the vowel", "Nasalizes the vowel", "Makes it short", "Doubles the consonant", "Marks a question"),
		listen("Listen and choose", "पानी", "पानी (pānī — water)", "पानी (pānī — water)", "पनी", "पाना", "पीनी"),
		mc("How do you read को?", "क + ो", "ko", "ko", "kau", "ku", "kā"),
		speak("Cora", "मुझे पानी और किताब दीजिए।"),
	)
	addVocab(db, l,
		hv("किताब", "kitāb", "book (f)", "यह किताब अच्छी है।", "yah kitāb acchī hai.", "This book is good.", finch),
		hv("पानी", "pānī", "water (m)", "पानी ठंडा है।", "pānī ṭhanḍā hai.", "The water is cold.", "Cora"),
		hv("माँ", "mā̃", "mother (f)", "मेरी माँ डॉक्टर हैं।", "merī mā̃ ḍākṭar haĩ.", "My mother is a doctor.", "Nana"),
		hv("हाँ / नहीं", "hā̃ / nahī̃", "yes / no", "हाँ, ठीक है।", "hā̃, ṭhīk hai.", "Yes, that's fine.", "Lumora"),
	)
	l = addLesson(db, s, "संयुक्त अक्षर — Conjuncts", 2, 15,
		char(finch, "When two consonants meet with no vowel between them, the first loses its vertical stroke or tucks in: प्यार pyār (love), नमस्ते namaste (स्त = s+t), हिन्दी hindī. Some combine unexpectedly: क्ष kṣa (क्षमा, forgiveness), त्र tra (मित्र, friend), ज्ञ gya (ज्ञान, knowledge), श्र shra (श्री). The ् (virama) shows a consonant with no vowel. The full stop is ।(danda)."),
		mc("Which conjunct is in नमस्ते?", "namaste", "स्त (s+t)", "स्त (s+t)", "क्ष", "त्र", "ज्ञ"),
		mc("How is ज्ञ usually pronounced in Hindi?", "ज्ञान (knowledge)", "gya", "gya", "jna", "ksha", "tra"),
		mc("What does त्र combine?", "मित्र (friend)", "t + r", "t + r", "k + ṣ", "s + t", "p + y"),
		mc("What is the Hindi full stop called?", "।", "पूर्ण विराम (danda)", "पूर्ण विराम (danda)", "शिरोरेखा", "मात्रा", "नुक़्ता"),
		listen("Listen and choose", "प्यार", "प्यार (pyār — love)", "प्यार (pyār — love)", "पार (pār — across)", "पयार", "प्रार"),
		speak("Mira", "मेरा मित्र हिन्दी सीखता है।"),
	)
	addVocab(db, l,
		hv("प्यार", "pyār", "love (m)", "माँ का प्यार।", "mā̃ kā pyār.", "A mother's love.", "Mira"),
		hv("मित्र / दोस्त", "mitra / dost", "friend (m)", "वह मेरा दोस्त है।", "vah merā dost hai.", "He is my friend.", "Blaze"),
		hv("ज्ञान", "gyān", "knowledge (m)", "ज्ञान ही शक्ति है।", "gyān hī shakti hai.", "Knowledge is power.", finch),
		hv("हिन्दी / हिंदी", "hindī", "Hindi (f)", "मैं हिंदी सीख रहा हूँ।", "mãi hindī sīkh rahā hū̃.", "I'm learning Hindi.", "Lumora"),
	)

	s = addSkillL(db, hi, u, "Greetings & Politeness", "नमस्ते, धन्यवाद, कृपया — and जी for respect.", "Hand", "#FF5C5C", 4, 32)
	l = addLesson(db, s, "नमस्ते!", 1, 15,
		char("Lumora", "नमस्ते (namaste), with palms together, works at any time of day for hello and goodbye. आप कैसे हैं? (how are you — to a man) / आप कैसी हैं? (to a woman) — मैं ठीक हूँ (I'm fine). धन्यवाद or शुक्रिया = thank you. कृपया = please. माफ़ कीजिए = excuse me / sorry. Add जी after names and answers for respect: हाँ जी, शर्मा जी."),
		mc("How do you thank someone?", "Thank you", "धन्यवाद", "धन्यवाद", "नमस्ते", "माफ़ कीजिए", "कृपया"),
		mc("How are you? (to a woman, polite)", "Asking a woman", "आप कैसी हैं?", "आप कैसी हैं?", "आप कैसे हैं?", "तुम कैसा है?", "आप कौन हैं?"),
		mc("Reply: I'm fine.", "आप कैसे हैं?", "मैं ठीक हूँ।", "मैं ठीक हूँ।", "मैं ठीक है।", "आप ठीक हैं।", "ठीक हो मैं।"),
		mc("What does जी add to 'हाँ जी'?", "हाँ जी", "Respect", "Respect", "A question", "Negation", "Plural"),
		mc("Excuse me / sorry", "Apologising", "माफ़ कीजिए", "माफ़ कीजिए", "शुक्रिया", "फिर मिलेंगे", "नमस्ते"),
		listen("Listen and choose", "फिर मिलेंगे", "See you again", "See you again", "Good morning", "Thank you", "Welcome"),
		speak("Lumora", "नमस्ते! आप कैसे हैं? मैं ठीक हूँ, धन्यवाद।"),
	)
	addVocab(db, l,
		hv("नमस्ते", "namaste", "hello / goodbye", "नमस्ते, शर्मा जी!", "namaste, sharmā jī!", "Hello, Mr Sharma!", "Lumora"),
		hv("धन्यवाद / शुक्रिया", "dhanyavād / shukriyā", "thank you", "बहुत-बहुत धन्यवाद।", "bahut-bahut dhanyavād.", "Thank you very much.", "Mira"),
		hv("कृपया", "kṛpayā", "please", "कृपया बैठिए।", "kṛpayā baiṭhie.", "Please sit down.", finch),
		hv("माफ़ कीजिए", "māf kījie", "excuse me; sorry", "माफ़ कीजिए, स्टेशन कहाँ है?", "māf kījie, sṭeshan kahā̃ hai?", "Excuse me, where is the station?", "Riko"),
		hv("ठीक", "ṭhīk", "fine, OK", "सब ठीक है।", "sab ṭhīk hai.", "Everything's fine.", "Blaze"),
	)

	s = addSkillL(db, hi, u, "Pronouns & होना (to be)", "मैं हूँ, आप हैं, वह है — and आप / तुम / तू.", "Users", "#9B51E0", 5, 42)
	l = addLesson(db, s, "मैं हूँ", 1, 15,
		char(finch, "होना (to be), present: मैं हूँ (I am), तू है (you, intimate), तुम हो (you, familiar), आप हैं (you, polite), यह/वह है (he/she/it — this/that), हम हैं (we), ये/वे हैं (they). Three 'you's: आप for strangers, elders and respect; तुम for friends; तू only for very close people or children. Use आप by default!"),
		mc("We are students.", "हम ___ छात्र", "हम छात्र हैं।", "हम छात्र हैं।", "हम छात्र है।", "हम छात्र हूँ।", "हम छात्र हो।"),
		mc("You are my teacher. (polite)", "आप", "आप मेरे अध्यापक हैं।", "आप मेरे अध्यापक हैं।", "आप मेरा अध्यापक है।", "तू मेरे अध्यापक हैं।", "आप मेरे अध्यापक हो।"),
		mc("Which 'you' should you use with a stranger?", "Register", "आप", "आप", "तुम", "तू", "वह"),
		mc("Which form goes with तुम?", "तुम कहाँ ___?", "हो", "हो", "हूँ", "है", "हैं"),
		mc("What does वह mean?", "वह", "he / she / that", "he / she / that", "we", "I", "they only"),
		listen("Listen and choose", "मैं भारतीय नहीं हूँ", "I am not Indian.", "I am not Indian.", "I am Indian.", "You are not Indian.", "We are Indian."),
		speak("Lumora", "मैं केन्या से हूँ। आप कहाँ से हैं?"),
	)
	addVocab(db, l,
		hv("मैं / हम", "mãi / ham", "I / we", "हम दोस्त हैं।", "ham dost haĩ.", "We are friends.", finch),
		hv("आप / तुम / तू", "āp / tum / tū", "you (polite / familiar / intimate)", "आप कैसे हैं?", "āp kaise haĩ?", "How are you?", "Lumora"),
		hv("वह / वे", "vah / ve", "he, she, that / they", "वह मेरी बहन है।", "vah merī bahan hai.", "She is my sister.", "Nana"),
		hv("छात्र / छात्रा", "chātra / chātrā", "student (m / f)", "मैं छात्रा हूँ।", "mãi chātrā hū̃.", "I'm a (female) student.", "Mira"),
	)

	s = addSkillL(db, hi, u, "Introducing Yourself", "मेरा नाम… है, मैं… से हूँ.", "Users", "#17A3DD", 6, 52)
	l = addLesson(db, s, "मेरा नाम…", 1, 15,
		char("Lumora", "मेरा नाम लुमोरा है (my name is Lumora). आपका नाम क्या है? (what's your name?). मैं केन्या से हूँ (I'm from Kenya) — से is a postposition meaning 'from', and in Hindi it comes AFTER the noun. आपसे मिलकर ख़ुशी हुई (pleased to meet you)."),
		mc("What's your name? (polite)", "Asking", "आपका नाम क्या है?", "आपका नाम क्या है?", "तेरा नाम कौन है?", "आप नाम क्या हैं?", "नाम आपका कहाँ है?"),
		mc("I'm from Kenya.", "से after the noun", "मैं केन्या से हूँ।", "मैं केन्या से हूँ।", "मैं से केन्या हूँ।", "मैं केन्या हूँ से।", "मैं केन्या में हूँ।"),
		mc("Pleased to meet you.", "Introductions", "आपसे मिलकर ख़ुशी हुई।", "आपसे मिलकर ख़ुशी हुई।", "आपका नाम ख़ुशी है।", "मैं ख़ुश नहीं हूँ।", "फिर मिलकर।"),
		mc("Where does the verb go in a Hindi sentence?", "Word order: SOV", "At the end", "At the end", "At the start", "After the subject", "Anywhere"),
		mc("What does 'मेरा नाम रवि है' mean?", "मेरा नाम रवि है।", "My name is Ravi.", "My name is Ravi.", "I know Ravi.", "Ravi is my friend.", "Is your name Ravi?"),
		speak("Lumora", "नमस्ते! मेरा नाम लुमोरा है। मैं केन्या से हूँ। आपसे मिलकर ख़ुशी हुई।"),
	)
	addVocab(db, l,
		hv("नाम", "nām", "name (m)", "मेरा नाम अमित है।", "merā nām amit hai.", "My name is Amit.", "Lumora"),
		hv("क्या", "kyā", "what; (question marker)", "यह क्या है?", "yah kyā hai?", "What is this?", finch),
		hv("से", "se", "from; with (postposition)", "मैं दिल्ली से हूँ।", "mãi dillī se hū̃.", "I'm from Delhi.", "Riko"),
		hv("ख़ुशी", "khushī", "happiness, pleasure (f)", "मुझे बहुत ख़ुशी हुई।", "mujhe bahut khushī huī.", "I was very happy.", "Mira"),
	)

	s = addSkillL(db, hi, u, "Numbers & Money", "एक से सौ तक, हज़ार, लाख, करोड़ — and कितने का है?", "Hash", "#00C2A8", 7, 62)
	l = addLesson(db, s, "एक, दो, तीन…", 1, 15,
		char(finch, "एक, दो, तीन, चार, पाँच, छह, सात, आठ, नौ, दस. Every number up to 100 has its own word — learn them gradually: ग्यारह (11), बीस (20), पचास (50), सौ (100), हज़ार (1,000). India counts big numbers in लाख (100,000) and करोड़ (10 million). Prices: यह कितने का है? (how much is this?) — सौ रुपये (100 rupees)."),
		mc("What number is सात?", "सात", "7", "7", "6", "8", "9"),
		mc("How do you say 5?", "5", "पाँच", "पाँच", "पचास", "चार", "छह"),
		mc("What is एक लाख?", "एक लाख", "100,000", "100,000", "1,000", "10,000,000", "1,000,000"),
		mc("How much is this?", "At a shop", "यह कितने का है?", "यह कितने का है?", "यह कहाँ है?", "यह क्या है?", "यह कौन है?"),
		listen("Listen and choose the price", "पचास रुपये", "₹50", "₹50", "₹15", "₹500", "₹5"),
		mc("What number is दस?", "दस", "10", "10", "2", "100", "12"),
		speak("Cora", "यह कितने का है? सौ रुपये।"),
	)
	addVocab(db, l,
		hv("एक, दो, तीन", "ek, do, tīn", "one, two, three", "दो चाय दीजिए।", "do chāy dījie.", "Two teas, please.", "Cora"),
		hv("दस", "das", "ten", "दस रुपये।", "das rupaye.", "Ten rupees.", "Cora"),
		hv("सौ", "sau", "hundred", "सौ लोग आए।", "sau log āe.", "A hundred people came.", finch),
		hv("रुपया / रुपये", "rupayā / rupaye", "rupee / rupees (m)", "पाँच सौ रुपये।", "pā̃ch sau rupaye.", "Five hundred rupees.", "Cora"),
		hv("कितना / कितने", "kitnā / kitne", "how much / how many", "कितने लोग हैं?", "kitne log haĩ?", "How many people are there?", "Riko"),
	)

	s = addSkillL(db, hi, u, "Gender & Adjectives", "लड़का / लड़की, अच्छा / अच्छी / अच्छे — learn gender with every noun.", "Layers", "#F5A623", 8, 72)
	l = addLesson(db, s, "Masculine & Feminine", 1, 18,
		char(finch, "Every noun is masculine or feminine. Many masculine nouns end in -आ (लड़का boy, कमरा room), many feminine in -ई (लड़की girl, रोटी bread) — but always check: किताब (book) is feminine, घर (house) masculine. Plurals: लड़का → लड़के, कमरा → कमरे; लड़की → लड़कियाँ; किताब → किताबें; masculine nouns not in -आ don't change: घर → घर."),
		mc("Plural of लड़का (boy)", "लड़का", "लड़के", "लड़के", "लड़कियाँ", "लड़काएँ", "लड़कों"),
		mc("Plural of लड़की (girl)", "लड़की", "लड़कियाँ", "लड़कियाँ", "लड़के", "लड़कीयें", "लड़कियों"),
		mc("Plural of किताब (book, f)", "किताब", "किताबें", "किताबें", "किताबे", "किताबों", "किताब"),
		mc("What gender is किताब?", "किताब", "Feminine", "Feminine", "Masculine", "Neuter", "Both"),
		mc("Plural of घर (house, m)", "घर", "घर", "घर", "घरें", "घरे", "घरों"),
		match("Match the noun to its gender", "कमरा (room)", "masculine", "masculine", "feminine", "neuter", "plural only"),
		speak("Pip", "एक लड़का, दो लड़के; एक लड़की, दो लड़कियाँ।"),
	)
	addVocab(db, l,
		hv("लड़की", "laṛkī", "girl (f)", "लड़की गाना गाती है।", "laṛkī gānā gātī hai.", "The girl sings a song.", "Mira"),
		hv("कमरा", "kamrā", "room (m)", "कमरा साफ़ है।", "kamrā sāf hai.", "The room is clean.", "Nana"),
		hv("रोटी", "roṭī", "flatbread (f)", "रोटी गरम है।", "roṭī garam hai.", "The roti is hot.", "Cora"),
		hv("आदमी / औरत", "ādmī / aurat", "man (m) / woman (f)", "वह औरत डॉक्टर है।", "vah aurat ḍākṭar hai.", "That woman is a doctor.", finch),
	)
	l = addLesson(db, s, "Adjective Agreement", 2, 18,
		char(finch, "Adjectives ending in -आ agree: अच्छा लड़का (good boy), अच्छी लड़की (good girl), अच्छे लड़के (good boys) — feminine is -ई in both singular and plural: अच्छी लड़कियाँ. Adjectives not ending in -आ never change: सुंदर लड़का, सुंदर लड़की, लाल किताबें. Adjectives come before the noun."),
		mc("a big house (घर, m)", "बड़ा / बड़ी / बड़े", "बड़ा घर", "बड़ा घर", "बड़ी घर", "बड़े घर (one house)", "घर बड़ा एक"),
		mc("good girls", "अच्छा", "अच्छी लड़कियाँ", "अच्छी लड़कियाँ", "अच्छे लड़कियाँ", "अच्छा लड़कियाँ", "अच्छीं लड़कियाँ"),
		mc("small rooms (कमरे, m pl)", "छोटा", "छोटे कमरे", "छोटे कमरे", "छोटा कमरे", "छोटी कमरे", "कमरे छोटा"),
		mc("a beautiful girl", "सुंदर (doesn't change)", "सुंदर लड़की", "सुंदर लड़की", "सुंदरी लड़की", "सुंदरा लड़की", "सुंदरे लड़की"),
		mc("a new book (किताब, f)", "नया", "नई किताब", "नई किताब", "नया किताब", "नए किताब", "किताब नया"),
		speak("Mira", "यह नई किताब बहुत अच्छी है।"),
	)
	addVocab(db, l,
		hv("अच्छा / अच्छी", "acchā / acchī", "good (m / f)", "अच्छा दिन!", "acchā din!", "A good day!", "Lumora"),
		hv("बड़ा / छोटा", "baṛā / choṭā", "big / small", "छोटा भाई।", "choṭā bhāī.", "Younger brother.", "Pip"),
		hv("नया / नई", "nayā / naī", "new (m / f)", "नई गाड़ी।", "naī gāṛī.", "A new car.", "Blaze"),
		hv("सुंदर", "sundar", "beautiful (invariable)", "सुंदर शहर।", "sundar shahar.", "A beautiful city.", "Zephyr"),
	)

	s = addSkillL(db, hi, u, "Postpositions", "में, पर, से, को, का/की/के — and the oblique case.", "Link", "#6C3FC5", 9, 86)
	l = addLesson(db, s, "In, On, From, Of", 1, 18,
		char(finch, "Postpositions follow the noun: कमरे में (in the room), मेज़ पर (on the table), दिल्ली से (from Delhi), राम को (to Ram). Before a postposition, masculine -आ nouns change to -ए (the oblique): कमरा → कमरे में, लड़का → लड़के को. 'Of' agrees with the thing owned: राम का घर (m), राम की किताब (f), राम के दोस्त (m pl)."),
		mc("in the room", "कमरा + में", "कमरे में", "कमरे में", "कमरा में", "में कमरा", "कमरी में"),
		mc("on the table", "मेज़ + पर", "मेज़ पर", "मेज़ पर", "पर मेज़", "मेज़े पर", "मेज़ में ऊपर"),
		mc("Ravi's book (किताब, f)", "का / की / के", "रवि की किताब", "रवि की किताब", "रवि का किताब", "रवि के किताब", "किताब रवि की"),
		mc("Ravi's brother (भाई, m)", "का / की / के", "रवि का भाई", "रवि का भाई", "रवि की भाई", "रवि के भाई (one brother)", "भाई का रवि"),
		mc("What does 'मैं बस से जाता हूँ' mean?", "से = by / from", "I go by bus.", "I go by bus.", "I go to the bus.", "I come from the bus.", "I don't take the bus."),
		match("Match", "दिल्ली से", "from Delhi", "from Delhi", "in Delhi", "to Delhi", "of Delhi"),
		speak("Riko", "किताब मेज़ पर है और फ़ोन कमरे में है।"),
	)
	addVocab(db, l,
		hv("में", "mẽ", "in", "मैं घर में हूँ।", "mãi ghar mẽ hū̃.", "I'm at home.", finch),
		hv("पर", "par", "on, at", "चाबी मेज़ पर है।", "chābī mez par hai.", "The key is on the table.", "Riko"),
		hv("को", "ko", "to; (object marker)", "माँ को फ़ोन करो।", "mā̃ ko fon karo.", "Call Mum.", "Nana"),
		hv("का / की / के", "kā / kī / ke", "of, 's (agrees with the owned thing)", "सीता का घर।", "sītā kā ghar.", "Sita's house.", finch),
		hv("मेज़", "mez", "table (f)", "मेज़ बड़ी है।", "mez baṛī hai.", "The table is big.", "Cora"),
	)

	s = addSkillL(db, hi, u, "Present Habitual", "मैं खाता हूँ / खाती हूँ — verbs agree with the speaker's gender.", "MessageCircle", "#00C2A8", 10, 98)
	l = addLesson(db, s, "What I Do", 1, 18,
		char(finch, "Habits and general facts: verb stem + ता/ती/ते + होना. The verb agrees with the subject's gender: a man says मैं हिंदी बोलता हूँ, a woman मैं हिंदी बोलती हूँ. वह चाय पीता है (he drinks tea) / पीती है (she). हम खाते हैं. Negative: नहीं before the verb, and the auxiliary is often dropped — मैं मांस नहीं खाता."),
		mc("I speak Hindi. (said by a woman)", "बोलना", "मैं हिंदी बोलती हूँ।", "मैं हिंदी बोलती हूँ।", "मैं हिंदी बोलता हूँ।", "मैं हिंदी बोलते हैं।", "मैं बोलती हिंदी हूँ।"),
		mc("He drinks tea.", "पीना", "वह चाय पीता है।", "वह चाय पीता है।", "वह चाय पीती है।", "वह चाय पीते हैं।", "वह पीता चाय है।"),
		mc("We eat rice.", "खाना", "हम चावल खाते हैं।", "हम चावल खाते हैं।", "हम चावल खाता हैं।", "हम चावल खाती है।", "हम खाते चावल हैं।"),
		mc("I don't eat meat. (man)", "Negative", "मैं मांस नहीं खाता।", "मैं मांस नहीं खाता।", "मैं नहीं मांस खाता हूँ।", "मैं मांस खाता नहीं हूँ नहीं।", "नहीं मैं मांस खाती।"),
		mc("Where do you live? (polite)", "रहना", "आप कहाँ रहते हैं?", "आप कहाँ रहते हैं?", "आप कहाँ रहता है?", "आप कहाँ रहती हूँ?", "तुम कहाँ रहते हैं?"),
		listen("Listen and choose", "वह स्कूल जाती है", "She goes to school.", "She goes to school.", "He goes to school.", "They go to school.", "I go to school."),
		speak("Cora", "मैं रोज़ सुबह चाय पीती हूँ।"),
	)
	addVocab(db, l,
		hv("बोलना", "bolnā", "to speak", "आप हिंदी बोलते हैं?", "āp hindī bolte haĩ?", "Do you speak Hindi?", "Lumora"),
		hv("खाना", "khānā", "to eat; food (m)", "हम साथ खाना खाते हैं।", "ham sāth khānā khāte haĩ.", "We eat together.", "Cora"),
		hv("पीना", "pīnā", "to drink", "वह कॉफ़ी पीती है।", "vah kôfī pītī hai.", "She drinks coffee.", "Nana"),
		hv("रहना", "rahnā", "to live, stay", "मैं मुंबई में रहता हूँ।", "mãi mumbaī mẽ rahtā hū̃.", "I live in Mumbai.", "Riko"),
		hv("रोज़", "roz", "every day", "वह रोज़ दौड़ती है।", "vah roz dauṛtī hai.", "She runs every day.", "Blaze"),
	)

	s = addSkillL(db, hi, u, "Present Continuous", "मैं पढ़ रहा हूँ / पढ़ रही हूँ.", "Clock", "#FF5C5C", 11, 110)
	l = addLesson(db, s, "Right Now", 1, 18,
		char(finch, "For what's happening now: stem + रहा/रही/रहे + होना. मैं पढ़ रहा हूँ (I'm reading — man), मैं पढ़ रही हूँ (woman), बच्चे खेल रहे हैं (the children are playing), वह खाना बना रही है (she is cooking)."),
		mc("The children are playing.", "खेलना", "बच्चे खेल रहे हैं।", "बच्चे खेल रहे हैं।", "बच्चे खेल रहा है।", "बच्चे खेलते रहे हैं।", "बच्चे खेल रही है।"),
		mc("I'm learning Hindi. (woman)", "सीखना", "मैं हिंदी सीख रही हूँ।", "मैं हिंदी सीख रही हूँ।", "मैं हिंदी सीख रहा हूँ।", "मैं हिंदी सीखती रही।", "मैं हिंदी सीख रहे हैं।"),
		mc("What are you doing? (polite)", "करना", "आप क्या कर रहे हैं?", "आप क्या कर रहे हैं?", "आप क्या करते रहे?", "आप क्या कर रहा है?", "क्या आप करते हैं क्या?"),
		mc("It's raining.", "Weather", "बारिश हो रही है।", "बारिश हो रही है।", "बारिश हो रहा है।", "बारिश होती रहा।", "बारिश है रही।"),
		mc("She is cooking.", "खाना बनाना", "वह खाना बना रही है।", "वह खाना बना रही है।", "वह खाना बना रहा है।", "वह खाना बनाती रहे।", "वह बना खाना रही।"),
		speak("Blaze", "मैं अभी काम कर रहा हूँ, बाद में बात करते हैं।"),
	)
	addVocab(db, l,
		hv("पढ़ना", "paṛhnā", "to read, study", "मैं किताब पढ़ रहा हूँ।", "mãi kitāb paṛh rahā hū̃.", "I'm reading a book.", finch),
		hv("खेलना", "khelnā", "to play", "बच्चे खेल रहे हैं।", "bacce khel rahe haĩ.", "The children are playing.", "Pip"),
		hv("बनाना", "banānā", "to make", "माँ खाना बना रही हैं।", "mā̃ khānā banā rahī haĩ.", "Mum is cooking.", "Nana"),
		hv("अभी", "abhī", "right now", "मैं अभी आ रहा हूँ।", "mãi abhī ā rahā hū̃.", "I'm coming right now.", "Blaze"),
		hv("बारिश", "bārish", "rain (f)", "बारिश हो रही है।", "bārish ho rahī hai.", "It's raining.", "Zephyr"),
	)

	s = addSkillL(db, hi, u, "Family & Possessives", "मेरा / मेरी / मेरे, माँ, पिता जी, भाई, बहन.", "Heart", "#17A3DD", 12, 122)
	l = addLesson(db, s, "मेरा परिवार", 1, 18,
		char("Nana", "Possessives agree with the thing owned, like का/की/के: मेरा भाई (my brother), मेरी बहन (my sister), मेरे माता-पिता (my parents). Others: तुम्हारा (your, familiar), आपका (your, polite), उसका (his/her), हमारा (our), उनका (their / his-her respectful). Elders get जी: पिता जी, दादी जी."),
		mc("my sister", "बहन (f)", "मेरी बहन", "मेरी बहन", "मेरा बहन", "मेरे बहन", "बहन मेरी का"),
		mc("your father (polite)", "पिता (m)", "आपके पिता जी", "आपके पिता जी", "आपकी पिता जी", "आपका पिता जी (respect plural needed)", "तेरा पिता"),
		mc("our house", "घर (m)", "हमारा घर", "हमारा घर", "हमारी घर", "हमारे घर (one house)", "हम का घर"),
		mc("Who is दादी?", "दादी जी", "Father's mother (grandmother)", "Father's mother (grandmother)", "Mother's mother", "Aunt", "Sister"),
		mc("my parents", "माता-पिता (pl)", "मेरे माता-पिता", "मेरे माता-पिता", "मेरा माता-पिता", "मेरी माता-पिता", "माता-पिता मेरा"),
		speak("Nana", "यह मेरी बहन है और ये मेरे दादा जी हैं।"),
	)
	addVocab(db, l,
		hv("परिवार", "parivār", "family (m)", "मेरा परिवार बड़ा है।", "merā parivār baṛā hai.", "My family is big.", "Nana"),
		hv("भाई / बहन", "bhāī / bahan", "brother (m) / sister (f)", "मेरे दो भाई हैं।", "mere do bhāī haĩ.", "I have two brothers.", "Pip"),
		hv("पिता जी / माता जी", "pitā jī / mātā jī", "father / mother (respectful)", "पिता जी घर पर हैं।", "pitā jī ghar par haĩ.", "Father is at home.", "Zephyr"),
		hv("दादा / दादी", "dādā / dādī", "paternal grandfather / grandmother", "दादी कहानी सुनाती हैं।", "dādī kahānī sunātī haĩ.", "Grandma tells stories.", "Nana"),
		hv("मेरा / मेरी / मेरे", "merā / merī / mere", "my (agrees with the owned thing)", "मेरी माँ।", "merī mā̃.", "My mother.", "Lumora"),
	)

	s = addSkillL(db, hi, u, "Food & the Market", "मुझे… चाहिए, … दीजिए — and bargaining.", "Coffee", "#00C2A8", 13, 134)
	l = addLesson(db, s, "बाज़ार में", 1, 18,
		char("Cora", "To ask for something: मुझे पानी चाहिए (I need/want water — 'to me water is wanted'). To order: एक चाय दीजिए (please give one tea). Staples: रोटी, चावल (rice), दाल (lentils), सब्ज़ी (vegetables), चाय. At the market: यह बहुत महँगा है! थोड़ा कम कीजिए (it's too expensive — lower it a little)."),
		mc("I need water.", "चाहिए", "मुझे पानी चाहिए।", "मुझे पानी चाहिए।", "मैं पानी चाहिए हूँ।", "मेरा पानी चाहिए।", "पानी मैं चाहता चाहिए।"),
		mc("Please give two teas.", "Ordering", "दो चाय दीजिए।", "दो चाय दीजिए।", "दो चाय दो।", "दो चाय लीजिए।", "चाय दो दीजिए।"),
		mc("This is very expensive!", "Bargaining", "यह बहुत महँगा है!", "यह बहुत महँगा है!", "यह बहुत सस्ता है!", "यह बहुत अच्छा है!", "यह बहुत छोटा है!"),
		mc("What is दाल?", "दाल-चावल", "Lentils", "Lentils", "Bread", "Tea", "Sweets"),
		mc("Lower the price a little, please.", "Bargaining", "थोड़ा कम कीजिए।", "थोड़ा कम कीजिए।", "थोड़ा ज़्यादा दीजिए।", "बहुत महँगा कीजिए।", "दाम बढ़ाइए।"),
		listen("Listen and choose", "खाना बहुत स्वादिष्ट है", "The food is very tasty.", "The food is very tasty.", "The food is very hot.", "The food is expensive.", "I don't like the food."),
		speak("Cora", "भैया, ये आम कितने के हैं? थोड़ा कम कीजिए!"),
	)
	addVocab(db, l,
		hv("चाहिए", "chāhie", "is needed / wanted (with मुझे)", "मुझे मदद चाहिए।", "mujhe madad chāhie.", "I need help.", "Lumora"),
		hv("दीजिए", "dījie", "please give", "पानी दीजिए।", "pānī dījie.", "Please give me water.", "Cora"),
		hv("महँगा / सस्ता", "mahãgā / sastā", "expensive / cheap", "यह सस्ता है।", "yah sastā hai.", "This is cheap.", "Cora"),
		hv("चावल / दाल", "chāval / dāl", "rice (m) / lentils (f)", "दाल-चावल मेरा पसंदीदा खाना है।", "dāl-chāval merā pasandīdā khānā hai.", "Dal-rice is my favourite food.", "Cora"),
		hv("स्वादिष्ट", "svādiṣṭ", "tasty, delicious", "खाना स्वादिष्ट है।", "khānā svādiṣṭ hai.", "The food is delicious.", "Pip"),
		hv("भैया / दीदी", "bhaiyā / dīdī", "brother / sister (friendly address)", "भैया, एक चाय!", "bhaiyā, ek chāy!", "Brother, one tea!", "Blaze"),
	)

	s = addSkillL(db, hi, u, "Questions", "क्या, कौन, कहाँ, कब, क्यों, कैसे, कितना.", "Quote", "#F5A623", 14, 146)
	l = addLesson(db, s, "Asking Questions", 1, 18,
		char(finch, "Question words usually sit just before the verb: आप कहाँ जा रहे हैं? (where are you going?). क्या (what), कौन (who), कहाँ (where), कब (when), क्यों (why), कैसे (how), कितना/कितनी/कितने (how much/many — agrees). At the start of a sentence, क्या turns it into a yes/no question: क्या आप हिंदी बोलते हैं?"),
		mc("Where are you going? (polite)", "कहाँ", "आप कहाँ जा रहे हैं?", "आप कहाँ जा रहे हैं?", "आप कब जा रहे हैं?", "आप क्यों जा रहे हैं?", "आप कौन जा रहे हैं?"),
		mc("Why are you sad?", "क्यों", "आप उदास क्यों हैं?", "आप उदास क्यों हैं?", "आप उदास कहाँ हैं?", "आप उदास कब हैं?", "आप उदास कौन हैं?"),
		mc("Who is that?", "कौन", "वह कौन है?", "वह कौन है?", "वह क्या है?", "वह कहाँ है?", "वह कैसे है?"),
		mc("What does क्या do at the start here?", "क्या आप चाय पीते हैं?", "Makes a yes/no question", "Makes a yes/no question", "Means 'what'", "Makes it negative", "Makes it past tense"),
		mc("When will the train come?", "कब", "ट्रेन कब आएगी?", "ट्रेन कब आएगी?", "ट्रेन कहाँ आएगी?", "ट्रेन क्यों आएगी?", "ट्रेन कैसे आएगी?"),
		speak("Riko", "माफ़ कीजिए, बस अड्डा कहाँ है? और बस कब आएगी?"),
	)
	addVocab(db, l,
		hv("कहाँ", "kahā̃", "where", "आप कहाँ रहते हैं?", "āp kahā̃ rahte haĩ?", "Where do you live?", "Riko"),
		hv("कब", "kab", "when", "आप कब आएँगे?", "āp kab āẽge?", "When will you come?", "Nana"),
		hv("क्यों", "kyõ", "why", "आप क्यों हँस रहे हैं?", "āp kyõ hãs rahe haĩ?", "Why are you laughing?", "Pip"),
		hv("कैसे", "kaise", "how", "यह कैसे बनता है?", "yah kaise bantā hai?", "How is this made?", "Cora"),
		hv("कौन", "kaun", "who", "दरवाज़े पर कौन है?", "darvāze par kaun hai?", "Who's at the door?", finch),
	)
}

// ───────────────────── A2 — रोज़मर्रा की ज़िंदगी: Everyday Life ─────────────────────

func seedHindiA2(db *gorm.DB) {
	const hi = "hi"
	const u = "A2 · रोज़मर्रा की ज़िंदगी — Everyday Life"
	finch := "Professor Finch"

	s := addSkillL(db, hi, u, "The Future Tense", "मैं जाऊँगा / जाऊँगी — गा, गी, गे.", "Compass", "#6C3FC5", 15, 170)
	l := addLesson(db, s, "कल…", 1, 20,
		char(finch, "Future = stem + ending + गा/गी/गे (by gender): मैं जाऊँगा / जाऊँगी, तुम जाओगे, वह जाएगा / जाएगी, हम जाएँगे, आप जाएँगे, वे जाएँगे. होना: होगा, होगी, होंगे — also used for guesses: वह घर पर होगा (he's probably at home)."),
		mc("I will go tomorrow. (man)", "जाना", "मैं कल जाऊँगा।", "मैं कल जाऊँगा।", "मैं कल जाऊँगी।", "मैं कल गया।", "मैं कल जाता हूँ।"),
		mc("She will come.", "आना", "वह आएगी।", "वह आएगी।", "वह आएगा।", "वह आएँगे।", "वह आई।"),
		mc("We will eat together.", "खाना", "हम साथ खाएँगे।", "हम साथ खाएँगे।", "हम साथ खाएगा।", "हम साथ खाया।", "हम साथ खाती हैं।"),
		mc("What does 'वह दफ़्तर में होगा' suggest?", "Future of probability", "He's probably at the office", "He's probably at the office", "He will go to the office next year", "He was at the office", "He must not be at the office"),
		mc("Will you come? (polite)", "आना", "क्या आप आएँगे?", "क्या आप आएँगे?", "क्या आप आएगा?", "क्या आप आओगी?", "क्या आप आया?"),
		speak("Zephyr", "अगले साल मैं भारत जाऊँगा और ताजमहल देखूँगा।"),
	)
	addVocab(db, l,
		hv("जाना", "jānā", "to go", "हम दिल्ली जाएँगे।", "ham dillī jāẽge.", "We'll go to Delhi.", "Zephyr"),
		hv("आना", "ānā", "to come", "वह कल आएगी।", "vah kal āegī.", "She'll come tomorrow.", "Mira"),
		hv("अगला साल", "aglā sāl", "next year (m)", "अगले साल मिलेंगे।", "agle sāl milẽge.", "See you next year.", "Lumora"),
		hv("देखना", "dekhnā", "to see, watch", "मैं फ़िल्म देखूँगा।", "mãi film dekhū̃gā.", "I'll watch a film.", "Blaze"),
	)

	s = addSkillL(db, hi, u, "Simple Past (Intransitive)", "मैं गया / गई, वह आया / आई.", "Clock", "#00C2A8", 16, 185)
	l = addLesson(db, s, "Went, Came, Slept", 1, 20,
		char(finch, "For intransitive verbs (no object — going, coming, sleeping), the simple past is the stem + आ/ई/ए/ईं agreeing with the subject: वह आया (he came), वह आई (she came), वे आए, वे आईं (they — women). जाना is irregular: गया / गई / गए. होना: हुआ / हुई / हुए."),
		mc("She went to the market.", "जाना", "वह बाज़ार गई।", "वह बाज़ार गई।", "वह बाज़ार गया।", "वह बाज़ार जाई।", "वह बाज़ार जाती।"),
		mc("We arrived late.", "पहुँचना", "हम देर से पहुँचे।", "हम देर से पहुँचे।", "हम देर से पहुँचा।", "हम देर से पहुँची।", "हम देर से पहुँचते।"),
		mc("The child slept.", "सोना (boy)", "बच्चा सोया।", "बच्चा सोया।", "बच्चा सोई।", "बच्चा सोए।", "बच्चा सोता।"),
		mc("What happened?", "होना", "क्या हुआ?", "क्या हुआ?", "क्या होया?", "क्या हुई?", "क्या होगा कल?"),
		mc("They came. (women)", "आना", "वे आईं।", "वे आईं।", "वे आए।", "वे आया।", "वे आई।"),
		speak("Mira", "मैं कल रात देर से घर पहुँची।"),
	)
	addVocab(db, l,
		hv("गया / गई", "gayā / gaī", "went (m / f)", "वह घर गया।", "vah ghar gayā.", "He went home.", finch),
		hv("पहुँचना", "pahũchnā", "to arrive, reach", "ट्रेन समय पर पहुँची।", "ṭren samay par pahũchī.", "The train arrived on time.", "Riko"),
		hv("सोना", "sonā", "to sleep; gold (m)", "बच्चा सो गया।", "bacchā so gayā.", "The child fell asleep.", "Nana"),
		hv("देर", "der", "lateness, delay (f)", "देर हो गई।", "der ho gaī.", "It got late.", "Blaze"),
	)

	s = addSkillL(db, hi, u, "The ने Construction", "मैंने रोटी खाई — the past of transitive verbs agrees with the object.", "Shield", "#F5A623", 17, 200)
	l = addLesson(db, s, "मैंने खाना खाया", 1, 22,
		char(finch, "The trickiest rule in Hindi: in the past tense of transitive verbs (eat, read, see, do), the subject takes ने and the verb agrees with the OBJECT, not the subject. मैंने खाना खाया (I ate food — खाना m), मैंने रोटी खाई (I ate roti — f), उसने किताबें पढ़ीं (she read books — f pl). Irregular pasts: करना → किया, देना → दिया, लेना → लिया, पीना → पिया. If the object takes को, the verb defaults to masculine singular: मैंने सीता को देखा."),
		mc("I ate roti. (रोटी, f)", "ने + object agreement", "मैंने रोटी खाई।", "मैंने रोटी खाई।", "मैंने रोटी खाया।", "मैं रोटी खाई।", "मैंने रोटी खाए।"),
		mc("She read a book. (किताब, f)", "पढ़ना", "उसने किताब पढ़ी।", "उसने किताब पढ़ी।", "उसने किताब पढ़ा।", "वह किताब पढ़ी।", "उसने किताब पढ़ीं।"),
		mc("Ravi drank tea. (चाय, f)", "पीना → पिया/पी", "रवि ने चाय पी।", "रवि ने चाय पी।", "रवि ने चाय पिया।", "रवि चाय पिया।", "रवि ने चाय पीया।"),
		mc("What did you do? (polite)", "करना → किया", "आपने क्या किया?", "आपने क्या किया?", "आप क्या किया?", "आपने क्या करा?", "आपने क्या की?"),
		mc("I saw Sita. (object with को)", "Default agreement", "मैंने सीता को देखा।", "मैंने सीता को देखा।", "मैंने सीता को देखी।", "मैं सीता को देखा।", "मैंने सीता देखी को।"),
		mc("Does ने apply to जाना (to go)?", "Intransitive", "No — मैं गया, not मैंने गया", "No — मैं गया, not मैंने गया", "Yes, always", "Only for women", "Only in questions"),
		speak("Cora", "मैंने आज सुबह चाय पी और दो रोटियाँ खाईं।"),
	)
	addVocab(db, l,
		hv("मैंने", "mãine", "I (subject with ने, past transitive)", "मैंने फ़िल्म देखी।", "mãine film dekhī.", "I watched a film.", finch),
		hv("किया", "kiyā", "did (करना)", "आपने बहुत अच्छा किया।", "āpne bahut acchā kiyā.", "You did very well.", "Lumora"),
		hv("दिया / लिया", "diyā / liyā", "gave / took", "उसने मुझे तोहफ़ा दिया।", "usne mujhe tohfā diyā.", "He gave me a present.", "Mira"),
		hv("चिट्ठी", "chiṭṭhī", "letter (f)", "मैंने चिट्ठी लिखी।", "mãine chiṭṭhī likhī.", "I wrote a letter.", "Nana"),
	)

	s = addSkillL(db, hi, u, "Past Habits & Ongoing Past", "मैं जाता था, वह खेल रही थी.", "BookOpen", "#FF5C5C", 18, 215)
	l = addLesson(db, s, "Used To & Was Doing", 1, 20,
		char(finch, "Past of होना: था (m sg), थी (f sg), थे (m pl), थीं (f pl). Combine it: past habit — मैं बचपन में क्रिकेट खेलता था (I used to play cricket as a child); past continuous — वह खाना बना रही थी (she was cooking). Simple 'was': मैं बीमार था (I was ill)."),
		mc("I used to live in Delhi. (man)", "Past habitual", "मैं दिल्ली में रहता था।", "मैं दिल्ली में रहता था।", "मैं दिल्ली में रहा था।", "मैं दिल्ली में रहता थी।", "मैं दिल्ली में रहूँगा।"),
		mc("She was cooking.", "Past continuous", "वह खाना बना रही थी।", "वह खाना बना रही थी।", "वह खाना बना रहा था।", "वह खाना बनाती है।", "उसने खाना बना रही थी।"),
		mc("We were tired.", "था / थी / थे", "हम थके हुए थे।", "हम थके हुए थे।", "हम थका हुआ था।", "हम थके हुए थी।", "हम थके हुए हैं।"),
		mc("Past of 'है' for a woman", "वह बीमार ___।", "थी", "थी", "था", "थे", "थीं"),
		mc("What does 'बचपन में' mean?", "बचपन में मैं गाँव में रहता था।", "In childhood", "In childhood", "Yesterday", "In the village", "Every day"),
		speak("Nana", "बचपन में हम हर गर्मी नानी के घर जाते थे।"),
	)
	addVocab(db, l,
		hv("था / थी / थे", "thā / thī / the", "was / were", "मौसम अच्छा था।", "mausam acchā thā.", "The weather was good.", finch),
		hv("बचपन", "bacpan", "childhood (m)", "बचपन में मैं बहुत खेलता था।", "bacpan mẽ mãi bahut kheltā thā.", "As a child I played a lot.", "Nana"),
		hv("गाँव", "gā̃v", "village (m)", "मेरे दादा गाँव में रहते थे।", "mere dādā gā̃v mẽ rahte the.", "My grandfather lived in a village.", "Zephyr"),
		hv("नानी", "nānī", "maternal grandmother (f)", "नानी कहानियाँ सुनाती थीं।", "nānī kahāniyā̃ sunātī thī̃.", "Grandma used to tell stories.", "Nana"),
	)

	s = addSkillL(db, hi, u, "Requests & Commands", "आइए, बैठिए, आओ, मत जाओ.", "Hand", "#17A3DD", 19, 230)
	l = addLesson(db, s, "Please Come In", 1, 20,
		char(finch, "Imperatives follow the three 'you's. आप (polite): stem + इए — आइए (come), बैठिए (sit), लीजिए (take), दीजिए (give), कीजिए (do). तुम: stem + ओ — आओ, बैठो, खाओ, करो. तू: bare stem — आ, बैठ. Negative: मत or न before the verb — मत जाओ (don't go), कृपया न छुएँ (please don't touch). ज़रा softens: ज़रा सुनिए (excuse me / listen a moment)."),
		mc("Please sit down. (polite)", "बैठना", "बैठिए।", "बैठिए।", "बैठो।", "बैठ।", "बैठता।"),
		mc("Come here! (to a friend)", "आना, तुम", "यहाँ आओ!", "यहाँ आओ!", "यहाँ आइए!", "यहाँ आ!", "यहाँ आता!"),
		mc("Don't go! (to a friend)", "Negative", "मत जाओ!", "मत जाओ!", "नहीं जाओ मत!", "जाओ नहीं!", "मत जाता!"),
		mc("Please do this. (polite)", "करना", "यह कीजिए।", "यह कीजिए।", "यह करिए मत।", "यह करो जी।", "यह करा।"),
		mc("What does 'ज़रा सुनिए' do?", "ज़रा सुनिए!", "Politely gets someone's attention", "Politely gets someone's attention", "Insults someone", "Says goodbye", "Asks the time"),
		speak("Lumora", "आइए, बैठिए, चाय लीजिए!"),
	)
	addVocab(db, l,
		hv("आइए", "āie", "please come", "अंदर आइए।", "andar āie.", "Please come in.", "Lumora"),
		hv("बैठिए", "baiṭhie", "please sit", "यहाँ बैठिए।", "yahā̃ baiṭhie.", "Please sit here.", "Nana"),
		hv("लीजिए", "lījie", "please take", "मिठाई लीजिए।", "miṭhāī lījie.", "Please have a sweet.", "Cora"),
		hv("मत", "mat", "don't", "शोर मत करो।", "shor mat karo.", "Don't make noise.", finch),
		hv("ज़रा", "zarā", "a little; just (softener)", "ज़रा रुकिए।", "zarā rukie.", "Wait a moment, please.", "Riko"),
	)

	s = addSkillL(db, hi, u, "Liking, Needing & Having", "मुझे… पसंद है, मेरे पास… है, मुझे भूख लगी है.", "Heart", "#9B51E0", 20, 245)
	l = addLesson(db, s, "मुझे पसंद है", 1, 20,
		char(finch, "Many feelings use the dative 'to me' (मुझे) as the subject: मुझे चाय पसंद है (I like tea), मुझे पता है (I know), मुझे भूख लगी है (I'm hungry), मुझे बुख़ार है (I have a fever). Owning things uses के पास: मेरे पास एक गाड़ी है (I have a car); for relatives just use the possessive: मेरे दो भाई हैं."),
		mc("I like mangoes.", "पसंद", "मुझे आम पसंद हैं।", "मुझे आम पसंद हैं।", "मैं आम पसंद हूँ।", "मेरा आम पसंद है।", "मुझे आम पसंद करते।"),
		mc("I have a car.", "के पास", "मेरे पास एक गाड़ी है।", "मेरे पास एक गाड़ी है।", "मुझे एक गाड़ी है।", "मैं एक गाड़ी हूँ।", "मेरी एक गाड़ी पास है।"),
		mc("I'm hungry.", "भूख", "मुझे भूख लगी है।", "मुझे भूख लगी है।", "मैं भूख हूँ।", "मेरा भूख है।", "मुझे भूखा लगा।"),
		mc("I don't know.", "पता", "मुझे पता नहीं।", "मुझे पता नहीं।", "मैं पता नहीं हूँ।", "मेरा पता नहीं।", "पता मुझे नहीं हूँ।"),
		mc("Do you like Hindi films? (polite)", "आपको", "क्या आपको हिंदी फ़िल्में पसंद हैं?", "क्या आपको हिंदी फ़िल्में पसंद हैं?", "क्या आप हिंदी फ़िल्में पसंद हैं?", "क्या आपका हिंदी फ़िल्म पसंद?", "आप फ़िल्में पसंद करता?"),
		speak("Pip", "मुझे बहुत भूख लगी है! मुझे समोसे पसंद हैं।"),
	)
	addVocab(db, l,
		hv("पसंद", "pasand", "liked (f) — मुझे… पसंद है", "मुझे संगीत पसंद है।", "mujhe sangīt pasand hai.", "I like music.", "Mira"),
		hv("के पास", "ke pās", "near; (have) — मेरे पास", "मेरे पास समय नहीं है।", "mere pās samay nahī̃ hai.", "I don't have time.", "Blaze"),
		hv("भूख / प्यास", "bhūkh / pyās", "hunger / thirst (f)", "मुझे प्यास लगी है।", "mujhe pyās lagī hai.", "I'm thirsty.", "Pip"),
		hv("पता", "patā", "knowledge (m) — मुझे पता है", "क्या आपको पता है?", "kyā āpko patā hai?", "Do you know?", finch),
		hv("गाड़ी", "gāṛī", "car; vehicle (f)", "उसकी गाड़ी नई है।", "uskī gāṛī naī hai.", "His car is new.", "Blaze"),
	)

	s = addSkillL(db, hi, u, "Directions & Travel", "सीधे जाइए, बाएँ मुड़िए — trains, autos and tickets.", "Plane", "#00C2A8", 21, 260)
	l = addLesson(db, s, "रास्ता", 1, 20,
		char("Riko", "Directions: सीधे जाइए (go straight), बाएँ मुड़िए (turn left), दाएँ मुड़िए (turn right), पास (near), दूर (far), के सामने (opposite), के पीछे (behind). Travel: रेलवे स्टेशन, टिकट, ऑटो (auto-rickshaw), प्लेटफ़ॉर्म. Ask: यह ट्रेन दिल्ली जाती है? and agree the auto fare first: कितना लेंगे?"),
		mc("Turn left.", "Directions", "बाएँ मुड़िए।", "बाएँ मुड़िए।", "दाएँ मुड़िए।", "सीधे जाइए।", "पीछे जाइए।"),
		mc("Where is the railway station?", "Asking the way", "रेलवे स्टेशन कहाँ है?", "रेलवे स्टेशन कहाँ है?", "रेलवे स्टेशन कब है?", "रेलवे स्टेशन क्या है?", "रेलवे स्टेशन कौन है?"),
		mc("How much will you charge? (to an auto driver)", "Fare", "कितना लेंगे?", "कितना लेंगे?", "कितना देंगे?", "कहाँ लेंगे?", "क्या लेंगे?"),
		mc("opposite the bank", "के सामने", "बैंक के सामने", "बैंक के सामने", "बैंक के पीछे", "बैंक में", "सामने बैंक का"),
		mc("Is it far?", "दूर", "क्या यह दूर है?", "क्या यह दूर है?", "क्या यह पास है?", "क्या यह दाएँ है?", "क्या यह सीधा है?"),
		speak("Riko", "भैया, कनॉट प्लेस चलेंगे? कितना लेंगे?"),
	)
	addVocab(db, l,
		hv("सीधे", "sīdhe", "straight ahead", "सीधे जाइए।", "sīdhe jāie.", "Go straight.", "Riko"),
		hv("बाएँ / दाएँ", "bāẽ / dāẽ", "left / right", "अगले चौराहे पर दाएँ मुड़िए।", "agle caurāhe par dāẽ muṛie.", "Turn right at the next crossing.", "Riko"),
		hv("टिकट", "ṭikaṭ", "ticket (m)", "दो टिकट दीजिए।", "do ṭikaṭ dījie.", "Two tickets, please.", "Riko"),
		hv("ऑटो", "ôṭo", "auto-rickshaw (m)", "चलो, ऑटो लेते हैं।", "calo, ôṭo lete haĩ.", "Come on, let's take an auto.", "Blaze"),
		hv("दूर / पास", "dūr / pās", "far / near", "स्टेशन पास है।", "sṭeshan pās hai.", "The station is near.", "Zephyr"),
	)

	s = addSkillL(db, hi, u, "Weather & Telling Time", "मौसम; साढ़े, सवा, पौने, डेढ़, ढाई.", "Waves", "#F5A623", 22, 275)
	l = addLesson(db, s, "कितने बजे हैं?", 1, 22,
		char(finch, "Time: कितने बजे हैं? — तीन बजे हैं (it's three o'clock). Hindi has special fraction words: सवा तीन (3:15, 'quarter-plus'), साढ़े तीन (3:30), पौने तीन (2:45, 'quarter-less'), and two irregulars: डेढ़ (1:30) and ढाई (2:30). Weather: मौसम कैसा है? — गर्मी है (it's hot), सर्दी है (cold), बारिश हो रही है; monsoon = मानसून."),
		mc("What time is साढ़े चार?", "साढ़े चार बजे", "4:30", "4:30", "4:15", "3:45", "5:30"),
		mc("How do you say 1:30?", "1:30", "डेढ़ बजे", "डेढ़ बजे", "साढ़े एक बजे", "ढाई बजे", "सवा एक बजे"),
		mc("What time is पौने पाँच?", "पौने पाँच", "4:45", "4:45", "5:15", "5:45", "4:15"),
		mc("How do you say 2:30?", "2:30", "ढाई बजे", "ढाई बजे", "डेढ़ बजे", "साढ़े दो बजे", "पौने तीन बजे"),
		mc("It's very hot today.", "Weather", "आज बहुत गर्मी है।", "आज बहुत गर्मी है।", "आज बहुत सर्दी है।", "आज बारिश हो रही है।", "आज मौसम ठंडा है।"),
		write("Write 4–5 sentences in Hindi (about 30 words, in Devanagari) about the weather where you live and what you do on a rainy day.",
			"मैं नैरोबी में रहती हूँ। सुबह थोड़ी सर्दी होती है, लेकिन दोपहर में धूप निकलती है। अप्रैल में बहुत बारिश होती है। बारिश के दिन मैं घर पर चाय पीती हूँ और किताब पढ़ती हूँ।"),
	)
	addVocab(db, l,
		hv("बजे", "baje", "o'clock", "पाँच बजे मिलते हैं।", "pā̃ch baje milte haĩ.", "Let's meet at five.", "Riko"),
		hv("साढ़े / सवा / पौने", "sāṛhe / savā / paune", "half past / quarter past / quarter to", "सवा दस बजे।", "savā das baje.", "At 10:15.", finch),
		hv("डेढ़ / ढाई", "ḍeṛh / ḍhāī", "one and a half / two and a half", "ढाई घंटे लगेंगे।", "ḍhāī ghanṭe lagẽge.", "It'll take two and a half hours.", "Zephyr"),
		hv("मौसम", "mausam", "weather (m)", "आज मौसम सुहावना है।", "āj mausam suhāvnā hai.", "The weather is pleasant today.", "Zephyr"),
		hv("गर्मी / सर्दी", "garmī / sardī", "heat, summer / cold, winter (f)", "दिल्ली में गर्मी बहुत होती है।", "dillī mẽ garmī bahut hotī hai.", "Delhi gets very hot.", "Nana"),
	)
}
