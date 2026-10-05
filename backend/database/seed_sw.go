package database

import "gorm.io/gorm"

// ===== Swahili (Kiswahili) course, A1 → C2 =====
//
// Standard Swahili (Kiswahili sanifu), the Kiunguja-based standard used in
// schools, media and government across Kenya and Tanzania. Swahili is
// CEFR-assessed, so the units map directly onto the app's levels.
//
// The order follows how learners actually crack the language: sounds and
// greetings first (respect for elders is part of grammar here), then the noun
// classes (ngeli) introduced gradually and always in context — M-/WA- people
// and KI-/VI- things first — alongside the verb template
// subject–tense–(object)–root–ending, which is the other half of the system.
//
// Swahili is written in the Latin alphabet and spelled phonetically, so
// lessons mix choosing with typing from A1 (translate / fill exercises are
// typed or picked per lesson-level rules in controllers/lesson_writing.go).
// Multiple-choice questions always carry their own options.

func seedSwahili(db *gorm.DB) {
	seedSwahiliA1(db)
	seedSwahiliA2(db)
	seedSwahiliB1(db)
	seedSwahiliB2(db)
	seedSwahiliC1(db)
	seedSwahiliC2(db)
	seedSwahiliListening(db)
	seedSwahiliReading(db)
}

// ───────────────────── A1 — Mwanzo: Survival Swahili ─────────────────────

func seedSwahiliA1(db *gorm.DB) {
	const sw = "sw"
	const u = "A1 · Mwanzo — Survival Swahili"
	finch := "Professor Finch"

	s := addSkillL(db, sw, u, "Sounds of Swahili", "Five pure vowels, ng', ny, ch, dh, th, gh — and where the stress falls.", "Music", "#6C3FC5", 1, 0)
	l := addLesson(db, s, "Vowels & Stress", 1, 15,
		char(finch, "Karibu! Swahili is spelled exactly as it sounds. The five vowels never change: a as in 'father', e as in 'bed', i as in 'machine', o as in 'more', u as in 'rule'. Every vowel is pronounced — kaa (crab) has two. Stress falls on the second-to-last syllable: ki-TA-bu, ha-BA-ri, ki-swa-HI-li."),
		mc("Where is the stress in 'kitabu' (book)?", "ki-ta-bu", "ki-TA-bu", "ki-TA-bu", "KI-ta-bu", "ki-ta-BU", "No stress"),
		mc("How is the 'e' in 'wewe' (you) pronounced?", "wewe", "Like 'e' in 'bed'", "Like 'e' in 'bed'", "Like 'ee' in 'see'", "Silent", "Like 'ay' in 'day'"),
		mc("How many syllables are in 'kaa' (crab)?", "ka-a", "2", "2", "1", "3", "4"),
		listen("Listen and choose the word", "habari", "habari (news)", "habari (news)", "harabi", "hebari", "habiri"),
		mc("Where is the stress in 'Kiswahili'?", "ki-swa-hi-li", "ki-swa-HI-li", "ki-swa-HI-li", "KI-swa-hi-li", "ki-SWA-hi-li", "ki-swa-hi-LI"),
		speak("Mira", "Karibu Kiswahili!"),
	)
	addVocab(db, l,
		vw("karibu", "welcome; come in", "Karibu nyumbani!", "Welcome home!", "Lumora"),
		vw("habari", "news; how are things?", "Habari yako?", "How are you? (What's your news?)", "Cora"),
		vw("Kiswahili", "the Swahili language", "Ninajifunza Kiswahili.", "I'm learning Swahili.", "Lumora"),
		vw("kitabu", "book", "Hiki ni kitabu.", "This is a book.", finch),
	)
	l = addLesson(db, s, "Special Consonants", 2, 15,
		char(finch, "A few sounds need care. ng' (with the apostrophe) is the 'ng' of 'singer': ng'ombe, cow. Without it, ng is as in 'finger': ngoma, drum. ny is like 'ny' in 'canyon': nyumba, house. dh is 'th' in 'this' (dhahabu, gold), th is 'th' in 'thin' (thelathini, thirty), gh is a throaty sound like the French r (ghali, expensive). These come mostly from Arabic loanwords."),
		mc("Which word means 'cow'?", "The 'ng' of singer", "ng'ombe", "ng'ombe", "ngombe", "nombe", "gombe"),
		mc("How is 'ny' in 'nyumba' (house) pronounced?", "nyumba", "Like 'ny' in 'canyon'", "Like 'ny' in 'canyon'", "Like 'n' then 'y' separately", "Like 'nee'", "Silent n"),
		mc("'dh' in 'dhahabu' (gold) sounds like…", "dhahabu", "'th' in 'this'", "'th' in 'this'", "'th' in 'thin'", "'d'", "'z'"),
		mc("Which sound starts 'ghali' (expensive)?", "ghali", "A throaty 'gh', like French r", "A throaty 'gh', like French r", "A hard 'g' as in 'go'", "Silent g", "'j' as in 'jam'"),
		listen("Listen and choose the word", "ngoma", "ngoma (drum)", "ngoma (drum)", "ng'ombe (cow)", "nyoka (snake)", "ngazi (ladder)"),
		speak("Blaze", "Ng'ombe, ngoma, nyumba, thelathini."),
	)
	addVocab(db, l,
		vw("ng'ombe", "cow, cattle", "Ng'ombe anakula nyasi.", "The cow is eating grass.", "Zephyr"),
		vw("ngoma", "drum; traditional dance", "Ngoma inapigwa.", "The drum is being played.", "Blaze"),
		vw("nyumba", "house", "Nyumba yangu ni ndogo.", "My house is small.", "Nana"),
		vw("ghali", "expensive", "Simu hii ni ghali.", "This phone is expensive.", "Cora"),
		vw("thelathini", "thirty", "Nina miaka thelathini.", "I'm thirty years old.", finch),
	)

	s = addSkillL(db, sw, u, "Greetings & Respect", "Hujambo, Habari, Shikamoo — greetings come first in Swahili.", "Hand", "#00C2A8", 2, 10)
	l = addLesson(db, s, "Hujambo & Habari", 1, 15,
		char("Lumora", "In East Africa you always greet before anything else — skipping it is rude. Hujambo? (Are you well?) is answered Sijambo (I'm well). To several people: Hamjambo? — Hatujambo. Habari? (Any news?) is answered Nzuri, Salama or Njema — even on a bad day! Add the time: Habari za asubuhi (good morning)."),
		mc("Someone says 'Hujambo?' You answer:", "Hujambo?", "Sijambo", "Sijambo", "Hujambo", "Hatujambo", "Shikamoo"),
		mc("How do you answer 'Habari yako?'", "Habari yako?", "Nzuri", "Nzuri", "Karibu", "Kwaheri", "Asante"),
		mc("Greeting a group: 'Hamjambo?' They answer:", "Hamjambo?", "Hatujambo", "Hatujambo", "Sijambo", "Hujambo", "Nzuri sana"),
		mc("Good morning (literally 'news of the morning?')", "Good morning", "Habari za asubuhi?", "Habari za asubuhi?", "Habari za jioni?", "Lala salama", "Habari za mchana?"),
		tr("Translate", "Thank you very much", "Asante sana"),
		fill("Complete the greeting", "Habari za ___? (Good evening — news of the evening?)", "jioni"),
		speak("Lumora", "Hujambo? Sijambo. Habari za asubuhi?"),
	)
	addVocab(db, l,
		vw("Hujambo?", "How are you? (to one person)", "Hujambo, mama?", "How are you, madam?", "Lumora"),
		vw("Sijambo", "I'm fine (reply to Hujambo)", "Sijambo, asante.", "I'm fine, thanks.", "Mira"),
		vw("nzuri", "good, fine", "Habari nzuri.", "The news is good. (I'm fine.)", "Cora"),
		vw("asante", "thank you", "Asante sana!", "Thank you very much!", "Pip"),
		vw("asubuhi", "morning", "Habari za asubuhi?", "Good morning?", "Nana"),
		vw("jioni", "evening", "Tutaonana jioni.", "We'll see each other in the evening.", "Zephyr"),
	)
	l = addLesson(db, s, "Shikamoo, Mambo & Goodbyes", 2, 15,
		char("Nana", "Respect for elders is built into greetings. A younger person greets an elder with Shikamoo (I hold your feet), and the elder replies Marahaba (I'm delighted). Among friends: Mambo? — Poa! or Safi! Leaving: Kwaheri (goodbye), Kwaherini to several people, Tutaonana (see you), Lala salama (good night). And Pole! (sorry / sympathy) when someone is unwell or tired."),
		mc("A young person greets their grandmother:", "Greeting an elder", "Shikamoo", "Shikamoo", "Mambo", "Marahaba", "Poa"),
		mc("The elder replies to 'Shikamoo' with:", "Shikamoo, bibi!", "Marahaba", "Marahaba", "Sijambo", "Poa", "Kwaheri"),
		mc("Your friend says 'Mambo?' You reply:", "Mambo?", "Poa!", "Poa!", "Marahaba", "Shikamoo", "Kwaherini"),
		mc("What do you say to someone who is sick or tired?", "Sympathy", "Pole!", "Pole!", "Hongera!", "Karibu!", "Poa!"),
		mc("Goodbye to a group of people", "Leaving several friends", "Kwaherini", "Kwaherini", "Kwaheri", "Karibuni", "Hatujambo"),
		tr("Translate", "Good night", "Lala salama"),
		speak("Nana", "Shikamoo, bibi! Marahaba, mjukuu wangu."),
	)
	addVocab(db, l,
		vw("Shikamoo", "respectful greeting to an elder", "Shikamoo, babu!", "Greetings, grandfather!", "Pip"),
		vw("Marahaba", "elder's reply to Shikamoo", "Marahaba, mwanangu.", "Delighted, my child.", "Nana"),
		vw("Mambo? — Poa!", "What's up? — Cool! (casual)", "Mambo vipi? Poa tu!", "What's up? Just fine!", "Blaze"),
		vw("pole", "sorry (sympathy)", "Pole sana kwa ugonjwa.", "So sorry you're ill.", "Nana"),
		vw("kwaheri", "goodbye", "Kwaheri, tutaonana kesho.", "Goodbye, see you tomorrow.", "Lumora"),
	)

	s = addSkillL(db, sw, u, "Introducing Yourself", "Pronouns, ni / si, jina langu ni…, ninatoka…", "Users", "#F5A623", 3, 20)
	l = addLesson(db, s, "Mimi ni…", 1, 15,
		char("Lumora", "Personal pronouns: mimi (I), wewe (you), yeye (he/she — Swahili has no gender!), sisi (we), ninyi (you all), wao (they). 'Is/am/are' is ni, and its negative is si: Mimi ni mwanafunzi (I'm a student), Mimi si mwalimu (I'm not a teacher). Jina langu ni… = my name is… Ninatoka Kenya = I'm from Kenya."),
		mc("What does 'yeye' mean?", "yeye", "he or she", "he or she", "he only", "she only", "they"),
		mc("I am not a teacher.", "Mimi ___ mwalimu.", "si", "si", "ni", "na", "ha"),
		tr("Translate", "My name is Anna", "Jina langu ni Anna"),
		tr("Translate", "I am a student", "Mimi ni mwanafunzi"),
		mc("Where are you from?", "Ask a new friend", "Unatoka wapi?", "Unatoka wapi?", "Jina lako nani?", "Uko wapi?", "Unataka nini?"),
		fill("Complete", "Nina___ Kenya. (I'm from Kenya.)", "toka"),
		speak("Lumora", "Jina langu ni Lumora. Ninatoka Kenya. Nimefurahi kukujua."),
	)
	addVocab(db, l,
		vw("mimi", "I, me", "Mimi ni Mkenya.", "I'm Kenyan.", "Lumora"),
		vw("wewe", "you (one person)", "Wewe ni nani?", "Who are you?", "Cora"),
		vw("jina", "name", "Jina lako nani?", "What's your name?", "Mira"),
		vw("mwanafunzi", "student", "Yeye ni mwanafunzi.", "She is a student.", finch),
		vw("Nimefurahi kukujua", "pleased to meet you", "Nimefurahi kukujua, Anna.", "Pleased to meet you, Anna.", "Lumora"),
	)

	s = addSkillL(db, sw, u, "Numbers & Money", "Count, give your age and ask the price.", "Hash", "#FF5C5C", 4, 30)
	l = addLesson(db, s, "Moja, Mbili, Tatu…", 1, 15,
		char(finch, "moja, mbili, tatu, nne, tano, sita, saba, nane, tisa, kumi. Teens: kumi na moja (11), kumi na mbili (12). Tens: ishirini (20), thelathini, arobaini, hamsini, sitini, sabini, themanini, tisini. Then mia (100), elfu (1,000), laki (100,000). Many of these come from Arabic. Prices: Bei gani? — How much? Shilingi mia mbili — 200 shillings."),
		mc("What number is 'saba'?", "saba", "7", "7", "6", "8", "9"),
		mc("How do you say 15?", "15", "kumi na tano", "kumi na tano", "tano na kumi", "hamsini", "kumi tano"),
		mc("What number is 'hamsini'?", "hamsini", "50", "50", "15", "5", "500"),
		tr("Translate", "How much is it?", "Bei gani"),
		listen("Listen and choose the price", "shilingi mia tatu", "300 shillings", "300 shillings", "30 shillings", "3,000 shillings", "130 shillings"),
		fill("Complete: I'm twenty years old.", "Nina miaka ___.", "ishirini"),
		speak("Cora", "Ndizi hizi ni bei gani? Shilingi mia moja."),
	)
	addVocab(db, l,
		vw("moja, mbili, tatu", "one, two, three", "Nataka chai moja.", "I want one tea.", finch),
		vw("kumi", "ten", "Nina vitabu kumi.", "I have ten books.", finch),
		vw("mia", "hundred", "Shilingi mia tano.", "Five hundred shillings.", "Cora"),
		vw("elfu", "thousand", "Elfu moja tu.", "Just one thousand.", "Cora"),
		vw("bei", "price", "Bei gani?", "How much?", "Cora"),
		vw("miaka", "years", "Ana miaka kumi.", "She is ten years old.", "Pip"),
	)

	s = addSkillL(db, sw, u, "Noun Classes I: People (M-/WA-)", "mtu → watu: the first ngeli and its agreement.", "Users", "#9B51E0", 5, 40)
	l = addLesson(db, s, "One Person, Many People", 1, 18,
		char(finch, "Every Swahili noun belongs to a class (ngeli) that changes its prefix for singular and plural — and everything that agrees with it changes too. The M-/WA- class is for people: mtu → watu (person → people), mtoto → watoto (child → children), mwalimu → walimu (teacher → teachers). Adjectives copy the prefix: mtu mzuri (a good person), watu wazuri (good people)."),
		mc("Plural of 'mtoto' (child)", "mtoto", "watoto", "watoto", "vitoto", "mitoto", "matoto"),
		mc("Plural of 'mwalimu' (teacher)", "mwalimu", "walimu", "walimu", "wamwalimu", "miwalimu", "valimu"),
		mc("Good people", "watu ___", "wazuri", "wazuri", "mzuri", "vizuri", "nzuri"),
		tr("Translate", "The children", "Watoto"),
		tr("Translate", "a good teacher", "mwalimu mzuri"),
		fill("Make it plural", "mtu mrefu (a tall person) → watu ___", "warefu"),
		speak("Mira", "Mwalimu mzuri, walimu wazuri."),
	)
	addVocab(db, l,
		vw("mtu / watu", "person / people", "Watu wengi wako sokoni.", "Many people are at the market.", finch),
		vw("mtoto / watoto", "child / children", "Watoto wanacheza.", "The children are playing.", "Pip"),
		vw("mwalimu / walimu", "teacher / teachers", "Mwalimu wetu ni mzuri.", "Our teacher is good.", finch),
		vw("-zuri", "good, nice, beautiful", "Ni mtu mzuri.", "She's a good person.", "Mira"),
		vw("-refu", "tall, long", "Kaka yangu ni mrefu.", "My brother is tall.", "Blaze"),
	)

	s = addSkillL(db, sw, u, "Noun Classes II: Things (KI-/VI-)", "kitabu → vitabu, and ch-/vy- before vowels.", "Layers", "#17A3DD", 6, 52)
	l = addLesson(db, s, "Things & Tools", 1, 18,
		char(finch, "The KI-/VI- class holds many things, tools and languages: kitabu → vitabu (book), kiti → viti (chair), Kiswahili, Kiingereza. Before a vowel ki- becomes ch- and vi- becomes vy-: chakula → vyakula (food), choo → vyoo (toilet). Agreement copies the prefix: kitabu kizuri, vitabu vizuri."),
		mc("Plural of 'kiti' (chair)", "kiti", "viti", "viti", "watiti", "miti", "makiti"),
		mc("Plural of 'chakula' (food)", "chakula", "vyakula", "vyakula", "vichakula", "chakula", "machakula"),
		mc("A big book", "kitabu ___", "kikubwa", "kikubwa", "mkubwa", "kubwa", "vikubwa"),
		mc("Which class is 'Kiingereza' (English language)?", "Languages start with…", "KI-/VI-", "KI-/VI-", "M-/WA-", "N-/N-", "JI-/MA-"),
		tr("Translate", "These books are good", "Vitabu hivi ni vizuri"),
		fill("Agreement", "kiti ki___ (a small chair, -dogo)", "dogo"),
		speak("Blaze", "Kitabu kikubwa, vitabu vikubwa."),
	)
	addVocab(db, l,
		vw("kitabu / vitabu", "book / books", "Nina vitabu vingi.", "I have many books.", finch),
		vw("kiti / viti", "chair / chairs", "Kaa kitini.", "Sit on the chair.", "Nana"),
		vw("chakula / vyakula", "food / foods", "Chakula ni kitamu.", "The food is delicious.", "Cora"),
		vw("-kubwa", "big", "Ni kitabu kikubwa.", "It's a big book.", "Blaze"),
		vw("-dogo", "small", "Kiti kidogo.", "A small chair.", "Pip"),
	)

	s = addSkillL(db, sw, u, "Present Tense with -na-", "The verb template: ni-na-soma, and the negative si-som-i.", "CheckCircle", "#00C2A8", 7, 64)
	l = addLesson(db, s, "I Am Reading", 1, 18,
		char(finch, "Swahili verbs are built like Lego: subject prefix + tense + root. Subject prefixes for people: ni- (I), u- (you), a- (he/she), tu- (we), m- (you all), wa- (they). Present tense is -na-: ninasoma (I read / am reading), anasoma (she reads), wanasoma (they read). The dictionary form starts with ku-: kusoma, to read."),
		mc("She is reading.", "kusoma (to read)", "anasoma", "anasoma", "ninasoma", "wanasoma", "unasoma"),
		mc("We are cooking.", "kupika (to cook)", "tunapika", "tunapika", "wanapika", "mnapika", "ninapika"),
		mc("What does 'wanacheza' mean?", "kucheza (to play)", "They are playing", "They are playing", "We are playing", "You are playing", "He is playing"),
		tr("Translate", "I am learning Swahili", "Ninajifunza Kiswahili"),
		fill("Complete: You (one person) are speaking.", "U___sema.", "na"),
		mc("Prefix for 'you all'", "___nasoma (you all are reading)", "m", "m", "u", "wa", "tu"),
		speak("Cora", "Ninasoma, unasoma, anasoma, tunasoma."),
	)
	addVocab(db, l,
		vw("kusoma", "to read; to study", "Ninasoma kitabu.", "I'm reading a book.", finch),
		vw("kupika", "to cook", "Mama anapika wali.", "Mum is cooking rice.", "Cora"),
		vw("kucheza", "to play; to dance", "Watoto wanacheza mpira.", "The children are playing football.", "Pip"),
		vw("kujifunza", "to learn", "Tunajifunza Kiswahili.", "We are learning Swahili.", "Lumora"),
		vw("kusema", "to say, speak", "Anasema nini?", "What is he saying?", "Mira"),
	)
	l = addLesson(db, s, "Negatives & Short Verbs", 2, 18,
		char(finch, "Present negative swaps the prefixes — si- (I), hu- (you), ha- (he/she), hatu- (we), ham- (you all), hawa- (they) — drops -na-, and changes the final -a to -i: sisomi (I don't read), hasomi, hawasomi. One-syllable verbs keep their ku- in positive tenses: kula (eat) → ninakula, kunywa (drink) → anakunywa; negative: sili (I don't eat), hanywi."),
		mc("I don't read.", "kusoma", "sisomi", "sisomi", "sinasoma", "hasomi", "sisoma"),
		mc("They don't play.", "kucheza", "hawachezi", "hawachezi", "hawacheza", "hachezi", "wanachezi"),
		mc("I am eating.", "kula (to eat — one syllable)", "ninakula", "ninakula", "ninala", "nanakula", "ninakulia"),
		mc("What does 'hanywi kahawa' mean?", "kunywa (to drink)", "He/she doesn't drink coffee", "He/she doesn't drink coffee", "I don't drink coffee", "They drink coffee", "Don't drink coffee"),
		tr("Translate", "I don't understand", "Sielewi"),
		fill("Negative: We don't cook.", "Hatu___. (kupika)", "piki"),
		speak("Nana", "Samahani, sielewi. Sema polepole, tafadhali."),
	)
	addVocab(db, l,
		vw("kula", "to eat", "Ninakula ugali.", "I'm eating ugali.", "Cora"),
		vw("kunywa", "to drink", "Anakunywa chai.", "She is drinking tea.", "Nana"),
		vw("kuelewa", "to understand", "Sielewi.", "I don't understand.", "Lumora"),
		vw("polepole", "slowly", "Sema polepole, tafadhali.", "Speak slowly, please.", "Lumora"),
		vw("samahani", "excuse me; sorry", "Samahani, choo kiko wapi?", "Excuse me, where is the toilet?", "Riko"),
	)

	s = addSkillL(db, sw, u, "Questions & Having", "nini, nani, wapi, lini — and nina / sina.", "Quote", "#F5A623", 8, 76)
	l = addLesson(db, s, "Asking & Having", 1, 18,
		char(finch, "Question words usually come at the end: Unataka nini? (What do you want?), Yeye ni nani? (Who is she?), Unakwenda wapi? (Where are you going?), Utafika lini? (When will you arrive?), Habari gani? (What news?), Vipi? (How?), ngapi? (how many?). 'To have' is -na: nina (I have), una, ana, tuna, mna, wana; negative sina, huna, hana. Yes/no questions can start with Je."),
		mc("Where are you going?", "Unakwenda ___?", "wapi", "wapi", "nini", "nani", "lini"),
		mc("Who is he?", "Yeye ni ___?", "nani", "nani", "nini", "wapi", "gani"),
		mc("I don't have money.", "kuwa na pesa", "Sina pesa.", "Sina pesa.", "Nina pesa.", "Hana pesa.", "Si pesa."),
		mc("How many children do you have?", "Una watoto ___?", "wangapi", "wangapi", "gani", "wapi", "nini"),
		tr("Translate", "What do you want?", "Unataka nini"),
		tr("Translate", "Do you have water?", "Una maji"),
		speak("Riko", "Samahani, una chenji? Hapana, sina."),
	)
	addVocab(db, l,
		vw("nini", "what", "Hii ni nini?", "What is this?", finch),
		vw("nani", "who", "Ni nani?", "Who is it?", "Riko"),
		vw("wapi", "where", "Unaishi wapi?", "Where do you live?", "Mira"),
		vw("lini", "when", "Utarudi lini?", "When will you come back?", "Nana"),
		vw("nina / sina", "I have / I don't have", "Nina swali.", "I have a question.", "Pip"),
		vw("pesa", "money", "Sina pesa za kutosha.", "I don't have enough money.", "Blaze"),
	)

	s = addSkillL(db, sw, u, "Family & Possessives", "mama yangu, kitabu changu — 'my' agrees too.", "Heart", "#FF5C5C", 9, 88)
	l = addLesson(db, s, "My Family", 1, 18,
		char("Nana", "Possessive stems are -angu (my), -ako (your), -ake (his/her), -etu (our), -enu (your, pl.), -ao (their). They take the noun's class prefix: kitabu changu, vitabu vyangu, mtoto wangu, watoto wangu. Most family words — mama, baba, dada, kaka, rafiki — are in the N-class for possessives: mama yangu, kaka yangu."),
		mc("my book", "kitabu ___", "changu", "changu", "wangu", "yangu", "vyangu"),
		mc("my children", "watoto ___", "wangu", "wangu", "yangu", "changu", "zangu"),
		mc("her mother", "mama ___", "yake", "yake", "wake", "chake", "lake"),
		mc("What does 'dada yangu' mean?", "dada yangu", "my sister", "my sister", "my brother", "my mother", "my aunt"),
		tr("Translate", "my father", "baba yangu"),
		tr("Translate", "our teacher", "mwalimu wetu"),
		speak("Nana", "Huyu ni mama yangu na huyu ni kaka yangu."),
	)
	addVocab(db, l,
		vw("mama", "mother; madam", "Mama yangu ni daktari.", "My mother is a doctor.", "Nana"),
		vw("baba", "father; sir", "Baba yake ni mkulima.", "His father is a farmer.", "Zephyr"),
		vw("kaka / dada", "brother / sister", "Nina kaka mmoja na dada wawili.", "I have one brother and two sisters.", "Pip"),
		vw("bibi / babu", "grandmother / grandfather", "Bibi yangu anaishi kijijini.", "My grandmother lives in the village.", "Nana"),
		vw("familia", "family", "Familia yangu ni kubwa.", "My family is big.", "Lumora"),
		vw("rafiki", "friend", "Yeye ni rafiki yangu.", "She is my friend.", "Mira"),
	)

	s = addSkillL(db, sw, u, "Food & the Market", "Order, shop and bargain: Naomba…, Punguza bei!", "ShoppingBag", "#00C2A8", 10, 100)
	l = addLesson(db, s, "Sokoni — At the Market", 1, 18,
		char("Cora", "Ask politely with Naomba… (I request…) or Nataka… (I want…). Staples: wali (cooked rice), ugali (maize meal), nyama (meat), samaki (fish), maharagwe (beans), chai, maji. At the market, bargaining is normal: Bei gani? … Ghali sana! Punguza kidogo, tafadhali (lower it a little, please)."),
		mc("Polite way to ask for water", "Request", "Naomba maji, tafadhali.", "Naomba maji, tafadhali.", "Maji sasa!", "Wewe maji.", "Nina maji."),
		mc("What is 'ugali'?", "ugali na sukuma wiki", "A stiff maize-meal staple", "A stiff maize-meal staple", "Fried fish", "Rice with spices", "Sweet tea"),
		mc("That's very expensive!", "Bargaining", "Ghali sana!", "Ghali sana!", "Rahisi sana!", "Tamu sana!", "Karibu sana!"),
		mc("Lower the price a little, please.", "Bargaining", "Punguza bei kidogo, tafadhali.", "Punguza bei kidogo, tafadhali.", "Ongeza bei, tafadhali.", "Lipa sasa.", "Bei gani?"),
		tr("Translate", "I want two bananas", "Nataka ndizi mbili"),
		tr("Translate", "The food is delicious", "Chakula ni kitamu"),
		speak("Cora", "Naomba chai na mandazi mawili, tafadhali."),
	)
	addVocab(db, l,
		vw("Naomba…", "May I have… (polite)", "Naomba bili, tafadhali.", "May I have the bill, please.", "Cora"),
		vw("tafadhali", "please", "Saidia, tafadhali.", "Help, please.", "Lumora"),
		vw("ndizi", "banana", "Ndizi tano kwa shilingi mia.", "Five bananas for a hundred shillings.", "Cora"),
		vw("-tamu", "sweet; delicious", "Embe hili ni tamu.", "This mango is sweet.", "Pip"),
		vw("soko", "market", "Ninakwenda sokoni.", "I'm going to the market.", "Cora"),
		vw("maji", "water", "Naomba maji baridi.", "May I have cold water.", "Riko"),
	)

	s = addSkillL(db, sw, u, "Time & Days", "Swahili time starts at dawn — saa moja is 7 a.m.!", "Clock", "#6C3FC5", 11, 112)
	l = addLesson(db, s, "Saa Ngapi?", 1, 18,
		char(finch, "Swahili time counts from sunrise, so it's six hours off from English time: saa moja asubuhi = 7 a.m. (the first hour of daylight), saa sita mchana = noon, saa moja usiku = 7 p.m. Add na nusu (half past), na robo (quarter past), kasoro robo (quarter to). Days: Jumatatu (Mon), Jumanne, Jumatano, Alhamisi, Ijumaa (Fri), Jumamosi (Sat), Jumapili (Sun)."),
		mc("What time is 'saa moja asubuhi'?", "saa moja asubuhi", "7:00 a.m.", "7:00 a.m.", "1:00 a.m.", "1:00 p.m.", "6:00 a.m."),
		mc("What time is 'saa sita mchana'?", "saa sita mchana", "12:00 noon", "12:00 noon", "6:00 p.m.", "6:00 a.m.", "4:00 p.m."),
		mc("How do you say 9:30 a.m.?", "9:30 a.m.", "saa tatu na nusu asubuhi", "saa tatu na nusu asubuhi", "saa tisa na nusu asubuhi", "saa tatu kasoro robo", "saa nne na nusu asubuhi"),
		mc("Which day is Ijumaa?", "Ijumaa", "Friday", "Friday", "Thursday", "Saturday", "Monday"),
		tr("Translate", "What time is it?", "Saa ngapi"),
		fill("Days", "Siku ya kwanza ya wiki ya kazi ni ___ (Monday).", "Jumatatu"),
		speak("Riko", "Sasa ni saa mbili na robo asubuhi."),
	)
	addVocab(db, l,
		vw("saa", "hour; clock; watch", "Saa ngapi sasa?", "What time is it now?", "Riko"),
		vw("nusu / robo", "half / quarter", "Saa nne na nusu.", "Half past ten (10:30).", finch),
		vw("leo / kesho / jana", "today / tomorrow / yesterday", "Kesho ni Jumamosi.", "Tomorrow is Saturday.", "Pip"),
		vw("Jumatatu", "Monday", "Shule inaanza Jumatatu.", "School starts on Monday.", "Pip"),
		vw("Ijumaa", "Friday", "Ijumaa tutasafiri.", "On Friday we'll travel.", "Zephyr"),
		vw("usiku", "night", "Saa mbili usiku.", "8 p.m.", "Nana"),
	)
}

// ───────────────────── A2 — Maisha ya Kila Siku: Everyday Life ─────────────────────

func seedSwahiliA2(db *gorm.DB) {
	const sw = "sw"
	const u = "A2 · Maisha ya Kila Siku — Everyday Life"
	finch := "Professor Finch"

	s := addSkillL(db, sw, u, "Past Tense with -li-", "nilisoma — and the negative sikusoma.", "Clock", "#6C3FC5", 12, 140)
	l := addLesson(db, s, "Yesterday I…", 1, 18,
		char(finch, "Past tense swaps -na- for -li-: nilisoma (I read), alikwenda (she went), tulikula (we ate). The negative uses the negative prefix + -ku-: sikusoma (I didn't read), hakwenda (she didn't go), hatukula (we didn't eat)."),
		mc("I cooked.", "kupika", "nilipika", "nilipika", "ninapika", "nitapika", "nimepika"),
		mc("They didn't come.", "kuja (to come)", "hawakuja", "hawakuja", "hawaji", "hawatakuja", "walikuja"),
		mc("What does 'tulikula wali jana' mean?", "kula", "We ate rice yesterday", "We ate rice yesterday", "We will eat rice tomorrow", "We are eating rice", "We didn't eat rice"),
		tr("Translate", "I went to the market yesterday", "Nilikwenda sokoni jana"),
		fill("Negative: She didn't read.", "Ha___soma.", "ku"),
		tr("Translate", "We didn't see him", "Hatukumwona"),
		speak("Blaze", "Jana nilikwenda mjini na nilinunua viatu."),
	)
	addVocab(db, l,
		vw("kwenda", "to go", "Alikwenda shule.", "She went to school.", "Pip"),
		vw("kuja", "to come", "Walikuja jana.", "They came yesterday.", "Riko"),
		vw("kununua", "to buy", "Nilinunua mkate.", "I bought bread.", "Cora"),
		vw("jana", "yesterday", "Jana kulikuwa na mvua.", "It rained yesterday.", "Zephyr"),
		vw("mjini", "in town", "Ninafanya kazi mjini.", "I work in town.", "Blaze"),
	)

	s = addSkillL(db, sw, u, "Future & Perfect", "-ta- (will), -me- (have done) and -ja- (not yet).", "CheckCircle", "#00C2A8", 13, 155)
	l = addLesson(db, s, "Will Do, Have Done", 1, 18,
		char(finch, "Future: -ta- — nitasoma (I will read); negative si-ta-: sitasoma (I won't read). Perfect: -me- — nimekula (I have eaten), amefika (he has arrived), used constantly for states too: nimechoka (I'm tired), amelala (she's asleep). Its negative is -ja- (not yet): sijala (I haven't eaten yet), hajafika (he hasn't arrived yet)."),
		mc("I will travel tomorrow.", "kusafiri", "Nitasafiri kesho.", "Nitasafiri kesho.", "Nilisafiri kesho.", "Ninasafiri jana.", "Nimesafiri kesho."),
		mc("He hasn't arrived yet.", "kufika", "Hajafika.", "Hajafika.", "Hakufika.", "Hatafika.", "Amefika."),
		mc("What does 'nimechoka' mean?", "kuchoka (to get tired)", "I'm tired", "I'm tired", "I will be tired", "I was tired yesterday", "I'm not tired"),
		mc("We won't go.", "kwenda", "Hatutakwenda.", "Hatutakwenda.", "Hatukwenda.", "Hatujaenda.", "Tutakwenda."),
		tr("Translate", "I have finished", "Nimemaliza"),
		tr("Translate", "I haven't eaten yet", "Sijala"),
		speak("Mira", "Nimechoka sana. Nitalala mapema leo."),
	)
	addVocab(db, l,
		vw("kusafiri", "to travel", "Tutasafiri kwenda Mombasa.", "We'll travel to Mombasa.", "Zephyr"),
		vw("kufika", "to arrive", "Basi limefika.", "The bus has arrived.", "Riko"),
		vw("kuchoka", "to get tired", "Umechoka?", "Are you tired?", "Nana"),
		vw("kumaliza", "to finish", "Bado sijamaliza.", "I haven't finished yet.", finch),
		vw("bado", "still; not yet", "Bado yuko kazini.", "She's still at work.", "Mira"),
	)

	s = addSkillL(db, sw, u, "Noun Classes III: JI-/MA- & N-/N-", "gari → magari, nyumba → nyumba, and their agreement.", "Layers", "#F5A623", 14, 170)
	l = addLesson(db, s, "JI-/MA-", 1, 20,
		char(finch, "The JI-/MA- class often has no singular prefix: gari → magari (car), tunda → matunda (fruit), jina → majina (name), but jicho → macho (eye). Agreement: verbs take li- (singular) and ya- (plural): Gari limeharibika (the car has broken down), Magari yameharibika. Adjectives: tunda zuri, matunda mazuri."),
		mc("Plural of 'gari' (car)", "gari", "magari", "magari", "vigari", "wagari", "migari"),
		mc("Plural of 'jicho' (eye)", "jicho", "macho", "macho", "majicho", "vicho", "micho"),
		mc("The car has broken down.", "Gari ___haribika.", "lime", "lime", "ame", "kime", "ime"),
		mc("good fruits", "matunda ___", "mazuri", "mazuri", "wazuri", "vizuri", "nzuri"),
		tr("Translate", "The cars are big", "Magari ni makubwa"),
		fill("Agreement", "Yai ___mevunjika. (The egg has broken.)", "li"),
		speak("Blaze", "Gari langu limeharibika tena!"),
	)
	addVocab(db, l,
		vw("gari / magari", "car / cars", "Gari lake ni jipya.", "Her car is new.", "Blaze"),
		vw("tunda / matunda", "fruit / fruits", "Matunda haya ni matamu.", "These fruits are sweet.", "Cora"),
		vw("jicho / macho", "eye / eyes", "Ana macho mazuri.", "She has beautiful eyes.", "Mira"),
		vw("yai / mayai", "egg / eggs", "Nataka mayai mawili.", "I want two eggs.", "Cora"),
		vw("-pya", "new", "Gari jipya.", "A new car.", "Blaze"),
	)
	l = addLesson(db, s, "N-/N-", 2, 20,
		char(finch, "N-/N- nouns look the same in singular and plural: nyumba (house/houses), ndizi, kalamu (pen), barua (letter). Agreement shows the number: i- singular, zi- plural — Nyumba imejengwa (the house has been built), Nyumba zimejengwa (the houses have been built); nyumba yangu / nyumba zangu. Adjectives take n-/m- or nothing: nyumba nzuri, nyumba kubwa. Animals and people in this class (rafiki, ng'ombe) take M-/WA- agreement."),
		mc("my houses", "nyumba ___", "zangu", "zangu", "yangu", "vyangu", "wangu"),
		mc("The pens have got lost.", "Kalamu ___mepotea.", "zi", "zi", "i", "vi", "wa"),
		mc("a beautiful house", "nyumba ___", "nzuri", "nzuri", "mzuri", "kizuri", "zuri"),
		mc("My friend has arrived (animate N-noun).", "Rafiki yangu ___mefika.", "a", "a", "i", "zi", "li"),
		tr("Translate", "My house is big", "Nyumba yangu ni kubwa"),
		fill("Agreement", "Barua ___mefika. (The letter has arrived.)", "i"),
		speak("Nana", "Nyumba yetu ni ndogo lakini ni nzuri."),
	)
	addVocab(db, l,
		vw("kalamu", "pen(s)", "Kalamu yangu iko wapi?", "Where is my pen?", finch),
		vw("barua", "letter(s)", "Nimepokea barua yako.", "I've received your letter.", "Mira"),
		vw("simu", "phone(s)", "Simu yangu imeharibika.", "My phone is broken.", "Blaze"),
		vw("chai", "tea", "Chai imeiva.", "The tea is ready.", "Nana"),
		vw("-ingi", "many, much", "Kuna nyumba nyingi hapa.", "There are many houses here.", "Zephyr"),
	)

	s = addSkillL(db, sw, u, "Noun Classes IV: M-/MI-, U- & KU-", "mti → miti, ukuta → kuta, and verbs as nouns.", "Leaf", "#FF5C5C", 15, 185)
	l = addLesson(db, s, "Trees, Walls & Activities", 1, 20,
		char(finch, "M-/MI- is mostly trees, plants and body parts — not people: mti → miti (tree), mkono → mikono (arm/hand), mji → miji (town). Agreement: u- / i- — Mti umeanguka, Miti imeanguka. U- nouns are often long or abstract: ukuta → kuta (wall), ufunguo → funguo (key), uhuru (freedom). KU- turns verbs into nouns: Kusoma ni kuzuri — reading is good."),
		mc("Plural of 'mti' (tree)", "mti", "miti", "miti", "watu", "viti", "mati"),
		mc("The tree has fallen.", "Mti ___meanguka.", "u", "u", "a", "ki", "i"),
		mc("Plural of 'ufunguo' (key)", "ufunguo", "funguo", "funguo", "mafunguo", "vifunguo", "ufunguo"),
		mc("What does 'Kuogelea ni kuzuri' mean?", "kuogelea (to swim)", "Swimming is good", "Swimming is good", "I like swimming", "Don't swim here", "The swimmer is good"),
		tr("Translate", "My hand hurts", "Mkono wangu unauma"),
		fill("Agreement", "Miji mi___ (big towns, -kubwa)", "kubwa"),
		speak("Zephyr", "Miti mirefu, mji mkubwa, ufunguo mdogo."),
	)
	addVocab(db, l,
		vw("mti / miti", "tree / trees", "Mti huu ni mrefu.", "This tree is tall.", "Zephyr"),
		vw("mkono / mikono", "arm, hand / arms", "Osha mikono yako.", "Wash your hands.", "Nana"),
		vw("mji / miji", "town / towns", "Nairobi ni mji mkubwa.", "Nairobi is a big city.", "Riko"),
		vw("ufunguo / funguo", "key / keys", "Nimepoteza funguo zangu.", "I've lost my keys.", "Blaze"),
		vw("uhuru", "freedom, independence", "Siku ya Uhuru.", "Independence Day.", finch),
	)

	s = addSkillL(db, sw, u, "Object Infixes", "Ninakupenda: slot the object into the verb.", "Link", "#17A3DD", 16, 200)
	l = addLesson(db, s, "Ni-na-ku-penda", 1, 20,
		char(finch, "Objects can be slotted in before the verb root: -ni- (me), -ku- (you), -m-/-mw- (him/her), -tu- (us), -wa- (you all / them). Ninakupenda = I love you. Alimwona = she saw him. Things use their class: Nimekinunua (I've bought it — kitabu), Nimeyanunua (I've bought them — matunda). Verbs ending -a with -wa-: Ninawapenda (I love you all / them)."),
		mc("I love you.", "kupenda", "Ninakupenda", "Ninakupenda", "Unanipenda", "Ninampenda", "Tunakupenda"),
		mc("She saw him.", "kuona", "Alimwona", "Alimwona", "Alikuona", "Aliniona", "Walimwona"),
		mc("Help me!", "kusaidia", "Nisaidie!", "Nisaidie!", "Msaidie!", "Tusaidie!", "Wasaidie!"),
		mc("I've bought it (the book).", "kitabu → ki", "Nimekinunua", "Nimekinunua", "Nimemnunua", "Nimeinunua", "Nimelinunua"),
		tr("Translate", "Do you understand me?", "Unanielewa"),
		fill("Insert the object", "Nili___ona jana. (I saw you yesterday.)", "ku"),
		speak("Mira", "Nakupenda sana, rafiki yangu."),
	)
	addVocab(db, l,
		vw("kupenda", "to love, like", "Ninakupenda.", "I love you.", "Mira"),
		vw("kuona", "to see", "Nilimwona sokoni.", "I saw her at the market.", "Cora"),
		vw("kusaidia", "to help", "Unaweza kunisaidia?", "Can you help me?", "Lumora"),
		vw("kupigia simu", "to phone (someone)", "Nitakupigia simu.", "I'll call you.", "Riko"),
	)

	s = addSkillL(db, sw, u, "Commands & Requests", "Imperatives and the subjunctive: Njoo! Nisaidie! Usiende!", "Hand", "#9B51E0", 17, 215)
	l = addLesson(db, s, "Do This, Please", 1, 20,
		char(finch, "The imperative is the bare root: Soma! (read!), Andika! (write!). Plural adds -eni: Someni! Irregular: Njoo! (come!), Nenda! (go!), Kula! (eat!). The subjunctive (ending -e) is for polite requests and wishes: Nisaidie (help me), Ninataka uje (I want you to come), Tuende! (let's go!). Negative: usi- — Usiende! (don't go!)."),
		mc("Come! (to one person)", "kuja", "Njoo!", "Njoo!", "Kuja!", "Uje!", "Njooni!"),
		mc("Let's go!", "kwenda", "Tuende!", "Tuende!", "Tunakwenda!", "Twende kesho!", "Mwende!"),
		mc("Don't go!", "kwenda", "Usiende!", "Usiende!", "Huendi!", "Hutaenda!", "Sienda!"),
		mc("I want you to come tomorrow.", "Ninataka ___ kesho.", "uje", "uje", "unakuja", "utakuja", "ulikuja"),
		tr("Translate", "Write your name, please", "Andika jina lako, tafadhali"),
		fill("Plural command: Read! (to a class)", "Som___!", "eni"),
		speak("Lumora", "Karibuni! Ingieni, mkae, mnywe chai."),
	)
	addVocab(db, l,
		vw("Njoo! / Njooni!", "Come! (one / several)", "Njoo hapa!", "Come here!", "Nana"),
		vw("Nenda!", "Go!", "Nenda nyumbani.", "Go home.", "Zephyr"),
		vw("kuandika", "to write", "Andika barua.", "Write a letter.", finch),
		vw("kuingia", "to enter", "Ingia ndani.", "Come inside.", "Lumora"),
		vw("Twende!", "Let's go!", "Twende sokoni!", "Let's go to the market!", "Blaze"),
	)

	s = addSkillL(db, sw, u, "Directions & Travel", "Matatu, daladala, kulia, kushoto — and -ni for 'at'.", "Plane", "#00C2A8", 18, 230)
	l = addLesson(db, s, "Getting Around", 1, 20,
		char("Riko", "Directions: nenda moja kwa moja (go straight), pinda kulia (turn right), pinda kushoto (turn left), karibu na (near), mbele ya (in front of). Adding -ni makes a place: soko → sokoni (at the market), nyumba → nyumbani (at home), shule → shuleni. 'Where is…?' uses -ko: Hoteli iko wapi? Minibuses are matatu in Kenya and daladala in Tanzania."),
		mc("Turn left.", "Directions", "Pinda kushoto.", "Pinda kushoto.", "Pinda kulia.", "Nenda moja kwa moja.", "Simama hapa."),
		mc("Where is the bus station?", "Kituo cha basi ___ wapi?", "kiko", "kiko", "iko", "yuko", "liko"),
		mc("at school", "shule + -ni", "shuleni", "shuleni", "kwashule", "shulani", "nishule"),
		mc("What is a 'daladala'?", "Tanzania", "A shared minibus", "A shared minibus", "A motorbike taxi", "A train", "A ferry"),
		tr("Translate", "Where is the hotel?", "Hoteli iko wapi"),
		tr("Translate", "Stop here, please", "Simama hapa, tafadhali"),
		speak("Riko", "Samahani, matatu ya kwenda mjini inapatikana wapi?"),
	)
	addVocab(db, l,
		vw("kulia / kushoto", "right / left", "Pinda kulia hapo.", "Turn right there.", "Riko"),
		vw("moja kwa moja", "straight on", "Nenda moja kwa moja.", "Go straight on.", "Riko"),
		vw("kituo", "stop, station", "Kituo cha basi kiko mbele.", "The bus stop is ahead.", "Riko"),
		vw("matatu / daladala", "minibus (Kenya / Tanzania)", "Nilipanda matatu.", "I took a matatu.", "Blaze"),
		vw("boda boda", "motorbike taxi", "Tuchukue boda boda.", "Let's take a boda boda.", "Blaze"),
		vw("karibu na", "near", "Karibu na benki.", "Near the bank.", "Zephyr"),
	)

	s = addSkillL(db, sw, u, "Weather, Seasons & Being", "mvua, jua — and kuwa: nilikuwa, nitakuwa, -ko.", "Waves", "#F5A623", 19, 245)
	l = addLesson(db, s, "The Weather & 'To Be'", 1, 20,
		char("Zephyr", "East Africa has rainy seasons — masika (long rains) and vuli (short rains) — and the hot dry kiangazi. Weather talk: kuna jua (it's sunny), mvua inanyesha (it's raining), kuna baridi (it's cold), kuna joto (it's hot). 'To be' in other tenses uses kuwa: nilikuwa (I was), nitakuwa (I will be). Location uses -ko: niko nyumbani (I'm at home), yuko wapi? (where is she?)."),
		mc("It's raining.", "mvua", "Mvua inanyesha.", "Mvua inanyesha.", "Kuna jua.", "Mvua imekwisha.", "Kuna joto."),
		mc("I was a teacher.", "kuwa", "Nilikuwa mwalimu.", "Nilikuwa mwalimu.", "Nitakuwa mwalimu.", "Ni mwalimu.", "Nimekuwa mwalimu jana."),
		mc("Where are you?", "location", "Uko wapi?", "Uko wapi?", "Una wapi?", "Ni wapi wewe?", "Unaenda wapi?"),
		mc("What is 'masika'?", "masika", "The long rainy season", "The long rainy season", "The hot dry season", "A kind of fruit", "Winter snow"),
		tr("Translate", "It is very hot today", "Leo kuna joto sana"),
		write("Write 4–5 sentences in Swahili (about 30 words) describing the weather where you live and what you do on a rainy day.",
			"Ninaishi Nairobi. Asubuhi kuna baridi kidogo, lakini mchana kuna jua na joto. Mwezi wa nne mvua nyingi inanyesha. Siku ya mvua, ninakaa nyumbani, ninakunywa chai na ninasoma kitabu."),
	)
	addVocab(db, l,
		vw("mvua", "rain", "Mvua kubwa inanyesha.", "Heavy rain is falling.", "Zephyr"),
		vw("jua", "sun", "Leo kuna jua kali.", "There's strong sun today.", "Pip"),
		vw("baridi / joto", "cold / heat", "Kuna baridi asubuhi.", "It's cold in the morning.", "Nana"),
		vw("kuwa", "to be, become", "Nitakuwa daktari.", "I will be a doctor.", "Pip"),
		vw("niko / yuko", "I am (at) / she is (at)", "Niko kazini.", "I'm at work.", "Mira"),
	)
}
