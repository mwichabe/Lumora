package database

import (
	"encoding/json"

	"gorm.io/gorm"

	"lumora/backend/models"
)

// Seed inserts starter content. Each section is seeded independently and only
// when missing, so an existing database that is partially populated (e.g. it has
// characters but no skills) still gets the content it's lacking.
func Seed(db *gorm.DB) {
	var characters int64
	db.Model(&models.Character{}).Count(&characters)
	if characters == 0 {
		seedCharacters(db)
	}

	var quests int64
	db.Model(&models.Quest{}).Count(&quests)
	if quests == 0 {
		seedQuests(db)
	}

	var spanish int64
	db.Model(&models.Skill{}).Where("language = ?", "es").Count(&spanish)
	if spanish == 0 {
		seedSpanish(db)
	}

	var listening int64
	db.Model(&models.ListeningSession{}).Where("language = ?", "es").Count(&listening)
	if listening == 0 {
		seedListening(db)
	}

	var reading int64
	db.Model(&models.ReadingSession{}).Where("language = ?", "es").Count(&reading)
	if reading == 0 {
		seedReading(db)
	}

	// Each additional language has its own course (no fallback to Spanish).
	var german int64
	db.Model(&models.Skill{}).Where("language = ?", "de").Count(&german)
	if german == 0 {
		seedGerman(db)
	}
	var french int64
	db.Model(&models.Skill{}).Where("language = ?", "fr").Count(&french)
	if french == 0 {
		seedFrench(db)
	}
}

// --- language-aware builders (used by non-Spanish courses) -------------------

func addSkillL(db *gorm.DB, lang, unit, title, desc, icon, color string, order, reqXP int) uint {
	s := models.Skill{
		Language: lang, Unit: unit, Title: title, Description: desc,
		Icon: icon, Color: color, OrderIndex: order, RequiredXP: reqXP,
	}
	db.Create(&s)
	return s.ID
}

func addListeningL(db *gorm.DB, lang, unit, title, desc string, order, xp int, matches []models.ListeningMatch, lines []models.ListeningLine, qs []models.ListeningQuestion) {
	s := models.ListeningSession{Language: lang, Unit: unit, Title: title, Description: desc, OrderIndex: order, XPReward: xp}
	db.Create(&s)
	for i := range matches {
		matches[i].SessionID = s.ID
		matches[i].OrderIndex = i + 1
	}
	db.Create(&matches)
	for i := range lines {
		lines[i].SessionID = s.ID
		lines[i].OrderIndex = i + 1
	}
	db.Create(&lines)
	for i := range qs {
		qs[i].SessionID = s.ID
		qs[i].OrderIndex = i + 1
	}
	db.Create(&qs)
}

func addReadingL(db *gorm.DB, lang, unit, title, desc string, order, xp int, lines []models.ReadingLine, qs []models.ReadingQuestion) {
	s := models.ReadingSession{Language: lang, Unit: unit, Title: title, Description: desc, OrderIndex: order, XPReward: xp}
	db.Create(&s)
	for i := range lines {
		lines[i].SessionID = s.ID
		lines[i].OrderIndex = i + 1
	}
	db.Create(&lines)
	for i := range qs {
		qs[i].SessionID = s.ID
		qs[i].OrderIndex = i + 1
	}
	db.Create(&qs)
}

// ===== German course (A1 → C2, grammar-focused) ==============================
func seedGerman(db *gorm.DB) {
	const de = "de"
	seedGermanCourse(db)
}

// seedGermanCourse builds a full CEFR A1–C2 German course organised by grammar
// (Grammatik) subtopics, with lengthy conversations and advanced reading.
func seedGermanCourse(db *gorm.DB) {
	const de = "de"
	finch := "Professor Finch"

	// ───────────────────────── A1 · Grundlagen ─────────────────────────
	u1 := "A1 · Grundlagen"
	s := addSkillL(db, de, u1, "Artikel & Nomen", "Genders: der, die, das.", "Hand", "#6C3FC5", 1, 0)
	l := addLesson(db, s, "Bestimmte Artikel", 1, 15,
		char(finch, "Willkommen! Every German noun has a gender — der, die or das. Memorise the article with the noun."),
		mc("Choose the article", "___ Mann (man)", "der", "der", "die", "das", "den"),
		mc("Choose the article", "___ Frau (woman)", "die", "der", "die", "das", "dem"),
		mc("Choose the article", "___ Kind (child)", "das", "der", "die", "das", "den"),
		fill("Fill the blank", "___ Buch ist gut. (the book — neuter)", "Das"),
		speak("Blaze", "der Mann, die Frau, das Kind"),
	)
	addVocab(db, l,
		vw("der Mann", "the man", "Der Mann liest.", "The man reads.", finch),
		vw("die Frau", "the woman", "Die Frau arbeitet.", "The woman works.", "Cora"),
		vw("das Kind", "the child", "Das Kind spielt.", "The child plays.", "Cora"),
		vw("das Buch", "the book", "Das Buch ist neu.", "The book is new.", finch),
	)
	s = addSkillL(db, de, u1, "Präsens: sein & haben", "The two essential verbs.", "Sparkles", "#17A3DD", 2, 10)
	l = addLesson(db, s, "sein und haben", 1, 15,
		char(finch, "Two verbs you'll use constantly: sein (to be) and haben (to have). Learn their forms."),
		fill("Conjugate 'sein'", "Ich ___ müde. (I am tired)", "bin"),
		fill("Conjugate 'haben'", "Du ___ einen Bruder. (you have a brother)", "hast"),
		mc("Conjugate 'sein'", "Wir ___ Studenten. (we are)", "sind", "sind", "seid", "bin", "ist"),
		tr("Translate", "I am happy", "Ich bin glücklich"),
		speak("Blaze", "Ich bin glücklich. Ich habe Zeit."),
	)
	addVocab(db, l,
		vw("sein", "to be", "Ich bin hier.", "I am here.", finch),
		vw("haben", "to have", "Ich habe Zeit.", "I have time.", finch),
		vw("müde", "tired", "Ich bin müde.", "I am tired.", "Cora"),
		vw("glücklich", "happy", "Sie ist glücklich.", "She is happy.", "Cora"),
	)

	// Skill 3 — Personal pronouns
	s = addSkillL(db, de, u1, "Personalpronomen", "ich, du, er … mich, dich.", "Users", "#00C2A8", 3, 20)
	l = addLesson(db, s, "Personalpronomen", 1, 15,
		char(finch, "Subject pronouns: ich, du, er/sie/es, wir, ihr, sie/Sie. Accusative: mich, dich, ihn …"),
		mc("Which means 'we'?", "we", "wir", "wir", "ihr", "sie", "du"),
		fill("Fill in", "___ bist mein Freund. (you, informal)", "Du"),
		mc("Accusative of 'ich'", "Er sieht ___ . (me)", "mich", "mich", "mir", "ich", "dich"),
		tr("Translate", "They are nice", "Sie sind nett"),
		speak("Blaze", "ich, du, er, sie, es, wir, ihr, sie"),
	)
	addVocab(db, l,
		vw("ich", "I", "Ich bin hier.", "I am here.", finch),
		vw("du", "you (informal)", "Du bist nett.", "You are nice.", "Cora"),
		vw("er", "he", "Er kommt.", "He is coming.", "Cora"),
		vw("wir", "we", "Wir lernen.", "We learn.", finch),
	)

	// Skill 4 — Possessives
	s = addSkillL(db, de, u1, "Possessivartikel", "mein, dein, sein …", "Hand", "#F5A623", 4, 30)
	l = addLesson(db, s, "Possessivartikel", 1, 15,
		char(finch, "Possessives — mein, dein, sein, ihr, unser, euer — decline like 'ein'."),
		fill("Fill in", "Das ist ___ Buch. (my)", "mein"),
		mc("'your' (informal)", "___ Buch (your book)", "dein", "dein", "deine", "mein", "sein"),
		fill("Feminine noun → seine", "Das ist ___ Mutter. (his)", "seine"),
		tr("Translate", "This is my friend", "Das ist mein Freund"),
		speak("Blaze", "mein Buch, deine Tasche, sein Hund"),
	)
	addVocab(db, l,
		vw("mein", "my", "mein Buch", "my book", finch),
		vw("dein", "your", "dein Hund", "your dog", "Cora"),
		vw("sein", "his", "seine Mutter", "his mother", "Cora"),
		vw("unser", "our", "unser Haus", "our house", finch),
	)

	// Skill 5 — Regular present tense
	s = addSkillL(db, de, u1, "Verben im Präsens", "Regular conjugation.", "PenLine", "#17A3DD", 5, 42)
	l = addLesson(db, s, "Regelmäßige Verben", 1, 15,
		char(finch, "Regular verbs: stem + endings. spielen → ich spiele, du spielst, er spielt …"),
		fill("spielen → ich", "Ich ___ Fußball. (play)", "spiele"),
		fill("lernen → du", "Du ___ Deutsch. (learn)", "lernst"),
		mc("wohnen → er", "Er ___ in Berlin.", "wohnt", "wohnt", "wohne", "wohnst", "wohnen"),
		tr("Translate", "We learn German", "Wir lernen Deutsch"),
		speak("Blaze", "Ich spiele, du spielst, er spielt."),
	)
	addVocab(db, l,
		vw("spielen", "to play", "Ich spiele gern.", "I like to play.", "Cora"),
		vw("lernen", "to learn", "Ich lerne Deutsch.", "I learn German.", finch),
		vw("wohnen", "to live", "Ich wohne hier.", "I live here.", "Cora"),
		vw("arbeiten", "to work", "Sie arbeitet viel.", "She works a lot.", finch),
	)

	// Skill 6 — sein, haben, werden (irregular)
	s = addSkillL(db, de, u1, "sein, haben, werden", "The key irregular verbs.", "Sparkles", "#6C3FC5", 6, 54)
	l = addLesson(db, s, "Unregelmäßige Verben", 1, 15,
		char(finch, "Three you must know cold: sein (to be), haben (to have), werden (to become)."),
		fill("sein → ich", "Ich ___ müde.", "bin"),
		fill("haben → du", "Du ___ Hunger.", "hast"),
		mc("werden → er", "Er ___ Arzt. (becomes)", "wird", "wird", "werde", "wirst", "werden"),
		tr("Translate", "I am happy", "Ich bin glücklich"),
		speak("Blaze", "Ich bin, du bist, er ist."),
	)
	addVocab(db, l,
		vw("werden", "to become", "Er wird Arzt.", "He becomes a doctor.", finch),
		vw("der Hunger", "hunger", "Ich habe Hunger.", "I'm hungry.", "Cora"),
		vw("der Arzt", "doctor", "Sie ist Ärztin.", "She is a doctor.", "Cora"),
		vw("der Durst", "thirst", "Ich habe Durst.", "I'm thirsty.", finch),
	)

	// Skill 7 — Modal verbs
	s = addSkillL(db, de, u1, "Modalverben", "können, müssen, wollen, mögen.", "Layers", "#FF5C5C", 7, 68)
	l = addLesson(db, s, "Modalverben", 1, 15,
		char(finch, "Modal verbs push the MAIN verb to the end: 'Ich kann Deutsch sprechen.'"),
		mc("Modal 'can'", "Ich ___ Deutsch sprechen.", "kann", "kann", "kannst", "können", "muss"),
		fill("müssen → du", "Du ___ jetzt gehen. (must)", "musst"),
		mc("Where does the main verb go?", "Ich will ein Buch ___ .", "lesen", "lesen", "lese", "liest", "gelesen"),
		tr("Translate", "I would like a coffee", "Ich möchte einen Kaffee"),
		speak("Blaze", "Ich kann Deutsch sprechen."),
	)
	addVocab(db, l,
		vw("können", "can / to be able", "Ich kann schwimmen.", "I can swim.", finch),
		vw("müssen", "must / to have to", "Ich muss gehen.", "I must go.", finch),
		vw("wollen", "to want", "Ich will lernen.", "I want to learn.", "Cora"),
		vw("möchten", "would like", "Ich möchte Tee.", "I'd like tea.", "Cora"),
	)

	// Skill 8 — Separable verbs
	s = addSkillL(db, de, u1, "Trennbare Verben", "aufstehen, einkaufen, anrufen.", "Link2", "#F5A623", 8, 82)
	l = addLesson(db, s, "Trennbare Verben", 1, 15,
		char(finch, "Separable verbs split: 'aufstehen' → 'Ich stehe um 7 Uhr auf.' The prefix goes to the end."),
		mc("Conjugated stem", "Ich ___ um sieben Uhr auf. (aufstehen)", "stehe", "stehe", "aufstehe", "stehst", "steht"),
		fill("Prefix at the end", "Ich rufe dich ___ . (anrufen → call)", "an"),
		tr("Translate", "I am getting up", "Ich stehe auf"),
		speak("Blaze", "Ich stehe um sieben Uhr auf."),
	)
	addVocab(db, l,
		vw("aufstehen", "to get up", "Ich stehe früh auf.", "I get up early.", finch),
		vw("einkaufen", "to shop", "Ich kaufe ein.", "I go shopping.", "Cora"),
		vw("anrufen", "to call", "Ich rufe an.", "I call.", "Cora"),
		vw("fernsehen", "to watch TV", "Ich sehe fern.", "I watch TV.", finch),
	)

	// Skill 9 — Sentence structure & questions
	s = addSkillL(db, de, u1, "Satzbau & Fragen", "Verb position, W-Fragen.", "MessageCircle", "#00C2A8", 9, 96)
	l = addLesson(db, s, "Satzbau & Fragen", 1, 15,
		char(finch, "Verb SECOND in statements, FIRST in yes/no questions; W-questions start with a question word."),
		mc("Yes/no question", "___ du ins Kino? (go)", "Gehst", "Gehst", "Du gehst", "Gehen", "Geht"),
		mc("Question word 'where'", "___ wohnst du?", "Wo", "Wo", "Was", "Wann", "Wer"),
		fill("Question word 'what' (name)", "___ heißt du?", "Wie"),
		tr("Translate", "Where is the station?", "Wo ist der Bahnhof?"),
		speak("Blaze", "Wo wohnst du? Was machst du?"),
	)
	addVocab(db, l,
		vw("wo", "where", "Wo bist du?", "Where are you?", finch),
		vw("was", "what", "Was ist das?", "What is that?", "Cora"),
		vw("wann", "when", "Wann kommst du?", "When do you come?", "Cora"),
		vw("wer", "who", "Wer ist das?", "Who is that?", finch),
	)

	// Skill 10 — Accusative
	s = addSkillL(db, de, u1, "Akkusativ", "The direct object (den/einen).", "Hash", "#17A3DD", 10, 110)
	l = addLesson(db, s, "Akkusativ", 1, 15,
		char(finch, "Accusative = direct object. Only masculine changes: der→den, ein→einen."),
		mc("Accusative (masc.)", "Ich sehe ___ Mann. (the)", "den", "den", "der", "dem", "das"),
		fill("Accusative (masc., indef.)", "Ich habe ___ Bruder. (a)", "einen"),
		mc("Neuter — no change", "Ich kaufe ___ Auto. (a)", "ein", "ein", "einen", "eine", "einem"),
		tr("Translate", "I see the woman", "Ich sehe die Frau"),
		speak("Blaze", "Ich sehe den Mann und kaufe einen Apfel."),
	)
	addVocab(db, l,
		vw("der Hund", "the dog", "Ich habe einen Hund.", "I have a dog.", "Cora"),
		vw("der Apfel", "the apple", "Ich esse einen Apfel.", "I eat an apple.", "Cora"),
		vw("kaufen", "to buy", "Ich kaufe Brot.", "I buy bread.", finch),
		vw("sehen", "to see", "Ich sehe dich.", "I see you.", finch),
	)

	// Skill 11 — Prepositions
	s = addSkillL(db, de, u1, "Präpositionen", "in, auf, mit, zu …", "Compass", "#6C3FC5", 11, 124)
	l = addLesson(db, s, "Präpositionen", 1, 15,
		char(finch, "mit, zu, bei take dative; in/auf/an take accusative for movement, dative for location."),
		mc("Dative preposition", "Ich fahre ___ dem Bus. (with)", "mit", "mit", "für", "ohne", "durch"),
		mc("Location (auf + dem)", "Das Buch ist ___ Tisch. (on the)", "auf dem", "auf dem", "auf den", "in die", "an den"),
		fill("Movement → accusative", "Ich gehe ___ die Schule. (into/to)", "in"),
		tr("Translate", "I go to the doctor", "Ich gehe zum Arzt"),
		speak("Blaze", "Ich fahre mit dem Bus zur Schule."),
	)
	addVocab(db, l,
		vw("mit", "with", "mit dem Bus", "by bus", finch),
		vw("zu", "to", "zum Arzt", "to the doctor", finch),
		vw("der Bus", "the bus", "Der Bus kommt.", "The bus is coming.", "Cora"),
		vw("der Bahnhof", "the station", "Wo ist der Bahnhof?", "Where is the station?", "Cora"),
	)

	// Skill 12 — Negation
	s = addSkillL(db, de, u1, "Negation: nicht & kein", "Saying no.", "Hand", "#FF5C5C", 12, 138)
	l = addLesson(db, s, "Negation", 1, 15,
		char(finch, "Use 'nicht' to negate verbs/adjectives; 'kein' to negate nouns (no / not a)."),
		mc("Negate a noun", "Ich habe ___ Auto. (no car)", "kein", "kein", "nicht", "keine", "nein"),
		fill("Negate an adjective", "Das ist ___ gut. (not good)", "nicht"),
		mc("kein (feminine)", "Ich habe ___ Zeit. (no time)", "keine", "keine", "kein", "nicht", "keinen"),
		tr("Translate", "I don't have a brother", "Ich habe keinen Bruder"),
		speak("Blaze", "Ich habe kein Auto. Das ist nicht gut."),
	)
	addVocab(db, l,
		vw("nicht", "not", "Ich komme nicht.", "I'm not coming.", finch),
		vw("kein", "no / not a", "Ich habe kein Geld.", "I have no money.", finch),
		vw("ja", "yes", "Ja, gern.", "Yes, gladly.", "Cora"),
		vw("nein", "no", "Nein, danke.", "No, thanks.", "Cora"),
	)

	// Skill 13 — Numbers, time, dates
	s = addSkillL(db, de, u1, "Zahlen, Uhrzeit & Datum", "1–100, telling time, dates.", "Clock", "#F5A623", 13, 152)
	l = addLesson(db, s, "Zahlen & Uhrzeit", 1, 15,
		char(finch, "Numbers, clock time and dates. 'Es ist halb drei' = 2:30."),
		mc("Number", "___ = 3", "drei", "drei", "zwei", "vier", "fünf"),
		fill("Telling time", "Es ist ___ Uhr. (ten)", "zehn"),
		mc("'halb drei' means", "halb drei", "2:30", "2:30", "3:30", "3:00", "2:00"),
		tr("Translate", "Today is Monday", "Heute ist Montag"),
		speak("Blaze", "Es ist halb drei. Heute ist Montag."),
	)
	addVocab(db, l,
		vw("die Uhr", "clock / o'clock", "Es ist drei Uhr.", "It's three o'clock.", finch),
		vw("der Tag", "the day", "Schönen Tag!", "Have a nice day!", "Cora"),
		vw("der Montag", "Monday", "Am Montag arbeite ich.", "On Monday I work.", "Cora"),
		vw("heute", "today", "Heute ist schön.", "Today is nice.", finch),
	)

	// Skill 14 — Writing: a short email
	s = addSkillL(db, de, u1, "Schreiben: E-Mail", "A1 writing is all email — learn it well.", "PenLine", "#00C2A8", 14, 166)
	// Lesson 1 — the building blocks of an email (Anrede & Gruß)
	l = addLesson(db, s, "E-Mail: Anrede & Gruß", 1, 15,
		char(finch, "Every German email needs a greeting (Anrede) and a closing (Gruß). Informal: 'Liebe Anna' / 'Viele Grüße'. Formal: 'Sehr geehrte Frau …' / 'Mit freundlichen Grüßen'."),
		mc("Informal greeting to a friend (Anna)", "___ Anna,", "Liebe", "Liebe", "Sehr geehrte", "Hallo Herr", "Tschüss"),
		mc("Formal greeting (Mr Klein)", "___ Herr Klein,", "Sehr geehrter", "Sehr geehrter", "Liebe", "Hallo du", "Viele"),
		mc("Informal closing", "___ , Max", "Viele Grüße", "Viele Grüße", "Sehr geehrte", "Bitte", "Danke"),
		fill("Polite request", "Bitte ___ Sie mir Informationen. (send → schicken)", "schicken"),
		tr("Translate", "Best wishes", "Viele Grüße"),
	)
	addVocab(db, l,
		vw("Liebe / Lieber", "Dear (informal)", "Liebe Anna,", "Dear Anna,", finch),
		vw("Sehr geehrte/r", "Dear (formal)", "Sehr geehrte Frau Klein,", "Dear Ms Klein,", finch),
		vw("Viele Grüße", "Best wishes", "Viele Grüße, Max", "Best wishes, Max", "Cora"),
		vw("Mit freundlichen Grüßen", "Yours sincerely", "Mit freundlichen Grüßen", "Yours sincerely", "Cora"),
	)
	// Lesson 2 — write real A1 emails
	l = addLesson(db, s, "E-Mail schreiben", 2, 20,
		char(finch, "Now write! A good A1 email (~30–40 words) has a greeting, 2–3 sentences covering every point in the task, and a closing."),
		write("Write a short email to a hotel: say when you arrive and ask for a room and the price.",
			"Liebes Hotel-Team,\nich komme vom 10. bis 15. Juli nach München. Haben Sie ein Zimmer frei? Bitte schicken Sie mir Informationen über die Preise.\nViele Grüße,\nMax Müller"),
		write("Write a short email to a friend: invite them to your birthday party on Saturday.",
			"Liebe Anna,\nam Samstag habe ich Geburtstag. Möchtest du zu meiner Party kommen? Wir feiern ab 18 Uhr bei mir zu Hause.\nViele Grüße,\nLukas"),
		write("Write an email to cancel a doctor's appointment on Monday and ask for a new one.",
			"Sehr geehrte Frau Doktor Klein,\nleider kann ich am Montag nicht kommen, ich bin krank. Haben Sie am Mittwoch einen Termin frei? Vielen Dank.\nMit freundlichen Grüßen,\nAnna Rossi"),
		write("Write an email to a friend: suggest meeting at the cinema on Friday at 8 p.m.",
			"Lieber Tom,\nhast du am Freitag Zeit? Wir können um 20 Uhr ins Kino gehen. Treffen wir uns vor dem Kino? Bitte antworte mir.\nViele Grüße,\nLena"),
		write("Write an email to a language school: ask about a German course and the cost.",
			"Sehr geehrte Damen und Herren,\nich möchte Deutsch lernen. Haben Sie einen Kurs für Anfänger? Wann beginnt der Kurs und was kostet er? Vielen Dank.\nMit freundlichen Grüßen,\nKwame Mensah"),
	)
	addVocab(db, l,
		vw("das Zimmer", "the room", "Ein Zimmer, bitte.", "A room, please.", "Cora"),
		vw("der Termin", "the appointment", "Ich habe einen Termin.", "I have an appointment.", "Cora"),
		vw("absagen", "to cancel", "Ich muss absagen.", "I have to cancel.", finch),
		vw("einladen", "to invite", "Ich lade dich ein.", "I invite you.", finch),
	)

	// Skill 15 — Speaking: introduce yourself
	s = addSkillL(db, de, u1, "Sprechen: Vorstellung", "Introduce yourself aloud.", "MessageCircle", "#6C3FC5", 15, 180)
	addLesson(db, s, "Sich vorstellen", 1, 20,
		char(finch, "Speaking! Introduce yourself: name, age, country, city, and why you learn German."),
		mc("How to say 'My name is…'", "My name is …", "Ich heiße …", "Ich heiße …", "Ich komme …", "Ich wohne …", "Ich habe …"),
		speak("Blaze", "Ich heiße Max. Ich bin fünfundzwanzig Jahre alt."),
		speak("Blaze", "Ich komme aus Kenia und ich wohne in Berlin."),
		speak("Blaze", "Ich lerne Deutsch, weil ich in Deutschland arbeiten möchte."),
	)

	addListeningL(db, de, u1, "Der erste Schultag",
		"Professor Finch meets two new students. Listen, then answer.", 1, 20,
		[]models.ListeningMatch{
			lm("Guten Morgen", "Good morning"),
			lm("Wie heißen Sie?", "What's your name? (formal)"),
			lm("Ich komme aus...", "I come from..."),
			lm("Freut mich", "Nice to meet you"),
		},
		[]models.ListeningLine{
			ln(finch, "Guten Morgen! Willkommen im Deutschkurs. Wie heißen Sie?", "Good morning! Welcome to the German course. What's your name?"),
			ln("Lumora", "Guten Morgen, Herr Professor. Ich heiße Lumora.", "Good morning, Professor. My name is Lumora."),
			ln(finch, "Freut mich, Lumora. Und woher kommen Sie?", "Nice to meet you, Lumora. And where are you from?"),
			ln("Lumora", "Ich komme aus Kenia, aber ich wohne jetzt in Berlin.", "I'm from Kenya, but I now live in Berlin."),
			ln("Cora", "Hallo! Ich bin Cora und ich bin sehr aufgeregt!", "Hi! I'm Cora and I'm very excited!"),
			ln(finch, "Wunderbar. Dann fangen wir an!", "Wonderful. Then let's begin!"),
		},
		[]models.ListeningQuestion{
			lq("What is the student's name?", "Lumora", "Lumora", "Cora", "Finch", "Berlin"),
			lq("Where does Lumora live now?", "Berlin", "Berlin", "Kenia", "Köln", "Wien"),
			lq("How does Cora feel?", "Excited", "Excited", "Tired", "Sad", "Bored"),
		},
	)
	addReadingL(db, de, u1, "Sich vorstellen",
		"Read Lumora's introduction, then answer.", 1, 15,
		[]models.ReadingLine{
			rl("Hallo! Ich heiße Lumora und ich bin einundzwanzig Jahre alt.", "Hello! My name is Lumora and I am twenty-one years old."),
			rl("Ich komme aus Kenia, aber ich wohne und studiere in Berlin.", "I'm from Kenya, but I live and study in Berlin."),
			rl("Ich lerne Deutsch, weil ich die Sprache schön finde.", "I'm learning German because I find the language beautiful."),
			rl("In meiner Freizeit lese ich gern und höre Musik.", "In my free time I like to read and listen to music."),
		},
		[]models.ReadingQuestion{
			rq("How old is Lumora?", "21", "21", "18", "25", "30"),
			rq("Why is she learning German?", "She finds it beautiful", "She finds it beautiful", "For work", "For exams", "For travel"),
			rq("What does she do in her free time?", "Read and listen to music", "Read and listen to music", "Cook", "Run", "Paint"),
		},
	)
	addListeningL(db, de, u1, "Im Supermarkt",
		"Blaze goes shopping. Listen, then answer.", 2, 20,
		[]models.ListeningMatch{
			lm("der Supermarkt", "the supermarket"),
			lm("Was kostet das?", "How much is that?"),
			lm("das Brot", "the bread"),
			lm("die Tüte", "the bag"),
		},
		[]models.ListeningLine{
			ln("Blaze", "Entschuldigung, wo finde ich das Brot?", "Excuse me, where do I find the bread?"),
			ln("Cora", "Das Brot ist gleich dort links.", "The bread is right over there on the left."),
			ln("Blaze", "Danke. Was kostet das Wasser?", "Thanks. How much is the water?"),
			ln("Cora", "Eine Flasche kostet einen Euro.", "A bottle costs one euro."),
			ln("Blaze", "Gut, ich nehme zwei Flaschen und ein Brot.", "Good, I'll take two bottles and a bread."),
			ln("Cora", "Möchten Sie eine Tüte?", "Would you like a bag?"),
		},
		[]models.ListeningQuestion{
			lq("Where is the bread?", "On the left", "On the left", "On the right", "At the back", "Upstairs"),
			lq("How much is a bottle of water?", "One euro", "One euro", "Two euros", "Free", "Fifty cents"),
			lq("What does Blaze buy?", "Water and bread", "Water and bread", "Only bread", "Milk", "Apples"),
		},
	)
	addListeningL(db, de, u1, "Wo ist der Bahnhof?",
		"Mira asks for directions. Listen, then answer.", 3, 20,
		[]models.ListeningMatch{
			lm("der Bahnhof", "the station"),
			lm("links", "left"),
			lm("rechts", "right"),
			lm("geradeaus", "straight ahead"),
		},
		[]models.ListeningLine{
			ln("Mira", "Entschuldigung, wo ist der Bahnhof?", "Excuse me, where is the station?"),
			ln("Riko", "Gehen Sie geradeaus und dann links.", "Go straight ahead and then left."),
			ln("Mira", "Ist es weit?", "Is it far?"),
			ln("Riko", "Nein, nur fünf Minuten zu Fuß.", "No, only five minutes on foot."),
			ln("Mira", "Vielen Dank!", "Thank you very much!"),
			ln("Riko", "Gern geschehen!", "You're welcome!"),
		},
		[]models.ListeningQuestion{
			lq("What is Mira looking for?", "The station", "The station", "The hotel", "The bank", "The market"),
			lq("Which way after going straight?", "Left", "Left", "Right", "Back", "Up"),
			lq("How far is it?", "5 minutes", "5 minutes", "15 minutes", "1 hour", "Very far"),
		},
	)
	addListeningL(db, de, u1, "Mein Tag",
		"Lumora describes her daily routine. Listen, then answer.", 4, 20,
		[]models.ListeningMatch{
			lm("aufstehen", "to get up"),
			lm("um sieben Uhr", "at seven o'clock"),
			lm("frühstücken", "to have breakfast"),
			lm("am Abend", "in the evening"),
		},
		[]models.ListeningLine{
			ln("Cora", "Lumora, wann stehst du auf?", "Lumora, when do you get up?"),
			ln("Lumora", "Ich stehe um sieben Uhr auf.", "I get up at seven o'clock."),
			ln("Cora", "Und was machst du dann?", "And what do you do then?"),
			ln("Lumora", "Ich frühstücke und fahre zur Arbeit.", "I have breakfast and go to work."),
			ln("Cora", "Wann kommst du nach Hause?", "When do you come home?"),
			ln("Lumora", "Am Abend, gegen sechs Uhr.", "In the evening, around six."),
		},
		[]models.ListeningQuestion{
			lq("When does Lumora get up?", "Seven o'clock", "Seven o'clock", "Six o'clock", "Eight o'clock", "Nine o'clock"),
			lq("What does she do after waking?", "Breakfast and work", "Breakfast and work", "Sleep", "Study", "Run"),
			lq("When does she come home?", "Around six in the evening", "Around six in the evening", "At noon", "At midnight", "At nine"),
		},
	)
	addReadingL(db, de, u1, "Anna aus Italien",
		"Read about Anna, then answer.", 2, 15,
		[]models.ReadingLine{
			rl("Hallo, ich heiße Anna. Ich bin 25 Jahre alt und komme aus Italien.", "Hello, my name is Anna. I am 25 years old and come from Italy."),
			rl("Ich wohne jetzt in Berlin und lerne Deutsch in einer Schule.", "I now live in Berlin and learn German at a school."),
			rl("Am Morgen trinke ich Kaffee und gehe spazieren.", "In the morning I drink coffee and go for a walk."),
			rl("Ich habe einen Hund. Er heißt Max.", "I have a dog. His name is Max."),
		},
		[]models.ReadingQuestion{
			rq("How old is Anna?", "25", "25", "20", "30", "18"),
			rq("Where is she from?", "Italy", "Italy", "Germany", "Spain", "France"),
			rq("What does she do in the morning?", "Coffee and a walk", "Coffee and a walk", "Works", "Sleeps", "Studies"),
		},
	)
	addReadingL(db, de, u1, "Eine E-Mail",
		"Read Tom's email, then answer.", 3, 15,
		[]models.ReadingLine{
			rl("Liebe Frau Schmidt,", "Dear Ms Schmidt,"),
			rl("ich kann am Montag leider nicht kommen.", "unfortunately I cannot come on Monday."),
			rl("Haben Sie am Dienstag um 10 Uhr Zeit?", "Do you have time on Tuesday at 10 o'clock?"),
			rl("Bitte antworten Sie mir. Viele Grüße, Tom", "Please reply to me. Best wishes, Tom"),
		},
		[]models.ReadingQuestion{
			rq("When can Tom NOT come?", "Monday", "Monday", "Tuesday", "Friday", "Sunday"),
			rq("What time does he suggest?", "Tuesday 10:00", "Tuesday 10:00", "Monday 9:00", "Tuesday 14:00", "Friday 10:00"),
			rq("What does Tom ask for?", "A reply", "A reply", "Money", "A room", "Directions"),
		},
	)
	addReadingL(db, de, u1, "Meine Familie",
		"Read about the family, then answer.", 4, 15,
		[]models.ReadingLine{
			rl("Das ist meine Familie. Wir sind vier Personen.", "This is my family. We are four people."),
			rl("Mein Vater heißt Peter und meine Mutter heißt Maria.", "My father is called Peter and my mother is called Maria."),
			rl("Ich habe einen Bruder. Er ist zehn Jahre alt.", "I have a brother. He is ten years old."),
			rl("Wir haben auch einen Hund und eine Katze.", "We also have a dog and a cat."),
		},
		[]models.ReadingQuestion{
			rq("How many people in the family?", "Four", "Four", "Three", "Five", "Two"),
			rq("How old is the brother?", "Ten", "Ten", "Five", "Twelve", "Eight"),
			rq("What pets do they have?", "A dog and a cat", "A dog and a cat", "A dog only", "Fish", "A bird"),
		},
	)

	// ───────────────────────── A2 · Aufbaustufe ────────────────────────
	u2 := "A2 · Aufbaustufe"

	// Skill 16 — Dative case
	s = addSkillL(db, de, u2, "Der Dativ", "The indirect-object case.", "Hash", "#F5A623", 16, 190)
	l = addLesson(db, s, "Der Dativ", 1, 20,
		char(finch, "Smart tip: learn each new word WITH its article and plural. The dative is the indirect object — der→dem, die→der, das→dem, plural→den (+n)."),
		mc("Dative (masc.)", "Ich gebe ___ Mann ein Buch. (the)", "dem", "dem", "den", "der", "das"),
		mc("Dative (fem.)", "Ich helfe ___ Frau. (the)", "der", "der", "die", "dem", "den"),
		fill("Dative (indef., masc.)", "Ich danke ___ Lehrer. (a → einem)", "einem"),
		tr("Translate", "I give the child a book", "Ich gebe dem Kind ein Buch"),
		speak("Blaze", "Ich gebe dem Mann das Buch."),
	)
	addVocab(db, l,
		vw("geben", "to give", "Ich gebe dir das.", "I give you that.", finch),
		vw("helfen", "to help (+ dat.)", "Ich helfe dir.", "I help you.", finch),
		vw("danken", "to thank (+ dat.)", "Ich danke dir.", "I thank you.", "Cora"),
		vw("gehören", "to belong (+ dat.)", "Das gehört mir.", "That belongs to me.", "Cora"),
	)

	// Skill 17 — Dative pronouns
	s = addSkillL(db, de, u2, "Dativpronomen", "mir, dir, ihm, ihr …", "Users", "#00C2A8", 17, 204)
	l = addLesson(db, s, "Dativpronomen", 1, 20,
		char(finch, "Dative pronouns: mir, dir, ihm, ihr, uns, euch, ihnen/Ihnen."),
		mc("Dative of 'ich'", "Kannst du ___ helfen? (me)", "mir", "mir", "mich", "ich", "dir"),
		fill("Dative of 'du'", "Ich gebe ___ das Buch. (you)", "dir"),
		mc("Dative of 'er'", "Ich schreibe ___ eine E-Mail. (him)", "ihm", "ihm", "ihn", "ihr", "er"),
		tr("Translate", "She gives us the keys", "Sie gibt uns die Schlüssel"),
		speak("Blaze", "Gib mir bitte das Buch."),
	)
	addVocab(db, l,
		vw("mir", "(to) me", "Gib mir das.", "Give me that.", finch),
		vw("dir", "(to) you", "Ich helfe dir.", "I help you.", "Cora"),
		vw("ihm", "(to) him", "Ich danke ihm.", "I thank him.", "Cora"),
		vw("uns", "(to) us", "Hilf uns!", "Help us!", finch),
	)

	// Skill 18 — Perfekt
	s = addSkillL(db, de, u2, "Das Perfekt", "Talking about the past.", "Clock", "#17A3DD", 18, 218)
	l = addLesson(db, s, "Das Perfekt", 1, 20,
		char(finch, "Spoken past = Perfekt: haben/sein + Partizip II. Movement & change of state use 'sein'."),
		fill("Participle", "Ich habe Deutsch ___ . (learned)", "gelernt"),
		mc("Auxiliary", "Wir ___ nach Berlin gefahren.", "sind", "sind", "haben", "seid", "ist"),
		tr("Translate", "I have eaten", "Ich habe gegessen"),
		fill("Participle", "Sie ist nach Hause ___ . (gone)", "gegangen"),
		speak("Blaze", "Ich habe gegessen und bin nach Hause gegangen."),
	)
	addVocab(db, l,
		vw("gelernt", "learned", "Ich habe gelernt.", "I learned.", finch),
		vw("gegessen", "eaten", "Ich habe gegessen.", "I ate.", "Cora"),
		vw("gefahren", "travelled", "Ich bin gefahren.", "I travelled.", "Cora"),
		vw("gegangen", "gone", "Ich bin gegangen.", "I went.", finch),
	)

	// Skill 19 — Präteritum (war, hatte, modals)
	s = addSkillL(db, de, u2, "Präteritum", "war, hatte & modal past.", "Clock", "#6C3FC5", 19, 232)
	l = addLesson(db, s, "Präteritum", 1, 20,
		char(finch, "In writing and with sein/haben/modals, German uses the simple past (Präteritum): war, hatte, konnte, musste."),
		fill("sein → past", "Gestern ___ ich krank. (was)", "war"),
		fill("haben → past", "Wir ___ keine Zeit. (had → hatten)", "hatten"),
		mc("können → past", "Ich ___ nicht kommen. (could)", "konnte", "konnte", "kann", "könnte", "konntest"),
		tr("Translate", "I was at home", "Ich war zu Hause"),
		speak("Blaze", "Gestern war ich müde und hatte keine Zeit."),
	)
	addVocab(db, l,
		vw("war", "was", "Ich war da.", "I was there.", finch),
		vw("hatte", "had", "Ich hatte Glück.", "I was lucky.", finch),
		vw("konnte", "could", "Ich konnte nicht.", "I couldn't.", "Cora"),
		vw("gestern", "yesterday", "Gestern war Montag.", "Yesterday was Monday.", "Cora"),
	)

	// Skill 20 — Reflexive verbs
	s = addSkillL(db, de, u2, "Reflexive Verben", "sich freuen, sich waschen.", "Sparkles", "#FF5C5C", 20, 246)
	l = addLesson(db, s, "Reflexive Verben", 1, 20,
		char(finch, "Reflexive verbs need a reflexive pronoun: mich, dich, sich … 'Ich freue mich.'"),
		mc("Reflexive pronoun (ich)", "Ich freue ___ . (am happy)", "mich", "mich", "mir", "sich", "dich"),
		fill("Reflexive (du)", "Du wäschst ___ . (yourself)", "dich"),
		mc("Reflexive (er)", "Er interessiert ___ für Musik.", "sich", "sich", "ihn", "mich", "dich"),
		tr("Translate", "I am looking forward to the holiday", "Ich freue mich auf den Urlaub"),
		speak("Blaze", "Ich freue mich. Ich wasche mich."),
	)
	addVocab(db, l,
		vw("sich freuen", "to be glad", "Ich freue mich.", "I'm glad.", finch),
		vw("sich waschen", "to wash (oneself)", "Ich wasche mich.", "I wash.", "Cora"),
		vw("sich interessieren", "to be interested", "Ich interessiere mich für …", "I'm interested in …", "Cora"),
		vw("sich treffen", "to meet", "Wir treffen uns.", "We meet.", finch),
	)

	// Skill 21 — Modal verbs: dürfen & sollen
	s = addSkillL(db, de, u2, "dürfen & sollen", "Permission & advice.", "Layers", "#F5A623", 21, 260)
	l = addLesson(db, s, "dürfen & sollen", 1, 20,
		char(finch, "dürfen = be allowed to; sollen = should / be supposed to. Main verb still goes to the end."),
		mc("dürfen (ich)", "___ ich hier rauchen? (may I)", "Darf", "Darf", "Soll", "Muss", "Kann"),
		fill("sollen (du)", "Du ___ mehr lernen. (should)", "sollst"),
		mc("Word order", "Du sollst die Tür ___ .", "schließen", "schließen", "schließt", "geschlossen", "schließe"),
		tr("Translate", "You are not allowed to park here", "Du darfst hier nicht parken"),
		speak("Blaze", "Darf ich hier parken? Du sollst leise sein."),
	)
	addVocab(db, l,
		vw("dürfen", "to be allowed", "Ich darf gehen.", "I may go.", finch),
		vw("sollen", "should", "Du sollst kommen.", "You should come.", finch),
		vw("rauchen", "to smoke", "Rauchen verboten.", "No smoking.", "Cora"),
		vw("parken", "to park", "Parken verboten.", "No parking.", "Cora"),
	)

	// Skill 22 — Separable & inseparable verbs
	s = addSkillL(db, de, u2, "Trennbar & untrennbar", "Prefixes that split — or don't.", "Link2", "#00C2A8", 22, 274)
	l = addLesson(db, s, "Verben mit Vorsilben", 1, 20,
		char(finch, "Separable prefixes (an-, auf-, ein-, mit-) split off. Inseparable prefixes (be-, ver-, ent-, er-) never split."),
		mc("Separable: anrufen", "Ich rufe dich später ___ . (call → an)", "an", "an", "auf", "ein", "ab"),
		fill("Inseparable: bezahlen", "Ich ___ die Rechnung. (pay — no split)", "bezahle"),
		mc("Which is inseparable?", "Choose the inseparable verb", "verstehen", "verstehen", "aufstehen", "einkaufen", "mitkommen"),
		tr("Translate", "I understand the question", "Ich verstehe die Frage"),
		speak("Blaze", "Ich rufe dich an und bezahle die Rechnung."),
	)
	addVocab(db, l,
		vw("anrufen", "to call (sep.)", "Ich rufe an.", "I call.", finch),
		vw("bezahlen", "to pay (insep.)", "Ich bezahle.", "I pay.", "Cora"),
		vw("verstehen", "to understand", "Ich verstehe.", "I understand.", "Cora"),
		vw("vergessen", "to forget", "Ich vergesse nichts.", "I forget nothing.", finch),
	)

	// Skill 23 — Subordinate clauses
	s = addSkillL(db, de, u2, "Nebensätze", "weil, dass, wenn, ob.", "MessageCircle", "#6C3FC5", 23, 288)
	l = addLesson(db, s, "Nebensätze", 1, 20,
		char(finch, "After weil, dass, wenn, ob the conjugated verb goes to the END."),
		mc("Verb at the end", "Ich bleibe zu Hause, weil ich krank ___ .", "bin", "bin", "bist", "ist", "sein"),
		fill("dass-clause", "Ich glaube, dass er recht ___ . (has)", "hat"),
		mc("Indirect question 'whether'", "Ich weiß nicht, ___ er kommt.", "ob", "ob", "weil", "dass", "wenn"),
		tr("Translate", "I hope that you come", "Ich hoffe, dass du kommst"),
		speak("Blaze", "Ich lerne, weil ich die Prüfung bestehen will."),
	)
	addVocab(db, l,
		vw("weil", "because", "…, weil ich müde bin.", "…, because I'm tired.", finch),
		vw("dass", "that", "…, dass es gut ist.", "…, that it's good.", finch),
		vw("ob", "whether", "…, ob er kommt.", "…, whether he comes.", "Cora"),
		vw("glauben", "to believe", "Ich glaube dir.", "I believe you.", "Cora"),
	)

	// Skill 24 — Two-way prepositions
	s = addSkillL(db, de, u2, "Wechselpräpositionen", "Accusative or dative?", "Compass", "#17A3DD", 24, 302)
	l = addLesson(db, s, "Wechselpräpositionen", 1, 20,
		char(finch, "Wohin? (movement) → accusative. Wo? (location) → dative. in, auf, an, über, unter, neben, zwischen."),
		mc("Location → dative", "Das Bild hängt ___ der Wand. (on)", "an", "an", "auf", "in", "über"),
		fill("Movement → accusative", "Ich lege das Buch ___ den Tisch. (onto)", "auf"),
		mc("in + dem = ?", "Ich bin ___ Kino. (in the — location)", "im", "im", "ins", "in", "an"),
		tr("Translate", "I am going into the house", "Ich gehe ins Haus"),
		speak("Blaze", "Das Buch liegt auf dem Tisch."),
	)
	addVocab(db, l,
		vw("an", "at / on (vertical)", "an der Wand", "on the wall", finch),
		vw("über", "over / above", "über dem Tisch", "above the table", "Cora"),
		vw("neben", "next to", "neben dem Haus", "next to the house", "Cora"),
		vw("zwischen", "between", "zwischen uns", "between us", finch),
	)

	// Skill 25 — Adjective declension
	s = addSkillL(db, de, u2, "Adjektivdeklination", "Adjective endings.", "PenLine", "#FF5C5C", 25, 316)
	l = addLesson(db, s, "Adjektivendungen", 1, 20,
		char(finch, "After 'der/die/das' adjectives usually take -e or -en; after 'ein' watch the gender ending."),
		mc("After 'der' (nom.)", "der ___ Mann (big → groß)", "große", "große", "großer", "großen", "groß"),
		fill("After 'ein' (neut. nom.)", "ein ___ Auto (new → neu)", "neues"),
		mc("Accusative (masc.)", "Ich sehe den ___ Hund. (small → klein)", "kleinen", "kleinen", "kleine", "kleiner", "klein"),
		tr("Translate", "a beautiful city", "eine schöne Stadt"),
		speak("Blaze", "Das ist ein schönes Haus."),
	)
	addVocab(db, l,
		vw("groß", "big", "ein großes Haus", "a big house", finch),
		vw("klein", "small", "ein kleiner Hund", "a small dog", "Cora"),
		vw("neu", "new", "ein neues Auto", "a new car", "Cora"),
		vw("schön", "beautiful", "eine schöne Stadt", "a beautiful city", finch),
	)

	// Skill 26 — Comparative & superlative
	s = addSkillL(db, de, u2, "Komparativ & Superlativ", "bigger, the best.", "Sparkles", "#F5A623", 26, 330)
	l = addLesson(db, s, "Vergleichen", 1, 20,
		char(finch, "Comparative: adjective + -er (+ als). Superlative: am + adjective + -sten."),
		fill("Comparative of 'klein'", "Mein Auto ist ___ als deins. (smaller)", "kleiner"),
		mc("Superlative", "Er läuft ___ . (the fastest)", "am schnellsten", "am schnellsten", "schneller", "schnell", "schnellste"),
		mc("Irregular: gut →", "Das ist ___ . (better)", "besser", "besser", "guter", "gutter", "am gut"),
		tr("Translate", "She is taller than me", "Sie ist größer als ich"),
		speak("Blaze", "Berlin ist größer als Bonn."),
	)
	addVocab(db, l,
		vw("größer", "bigger", "größer als …", "bigger than …", finch),
		vw("besser", "better", "besser als …", "better than …", "Cora"),
		vw("am besten", "(the) best", "Das ist am besten.", "That's best.", "Cora"),
		vw("als", "than", "größer als ich", "bigger than me", finch),
	)

	// Skill 27 — Future tense
	s = addSkillL(db, de, u2, "Futur I", "werden + infinitive.", "Clock", "#06AECE", 27, 344)
	l = addLesson(db, s, "Futur I", 1, 20,
		char(finch, "Future / intentions: werden + infinitive at the end. 'Ich werde Deutsch lernen.'"),
		mc("werden (ich)", "Ich ___ morgen kommen.", "werde", "werde", "wird", "wirst", "werden"),
		fill("Infinitive at end", "Wir werden nach Berlin ___ . (drive — fahren)", "fahren"),
		mc("werden (du)", "Du ___ das schaffen.", "wirst", "wirst", "werde", "wird", "werden"),
		tr("Translate", "I will learn German", "Ich werde Deutsch lernen"),
		speak("Blaze", "Ich werde morgen früh aufstehen."),
	)
	addVocab(db, l,
		vw("werden", "will (future)", "Ich werde gehen.", "I will go.", finch),
		vw("morgen", "tomorrow", "Bis morgen!", "See you tomorrow!", "Cora"),
		vw("nächste Woche", "next week", "nächste Woche", "next week", "Cora"),
		vw("der Plan", "the plan", "Ich habe einen Plan.", "I have a plan.", finch),
	)

	// Skill 28 — Writing (A2 email/letter)
	s = addSkillL(db, de, u2, "Schreiben: Brief", "Longer emails (80–100 words).", "PenLine", "#00C2A8", 28, 358)
	l = addLesson(db, s, "Längere E-Mails", 1, 20,
		char(finch, "A2 writing is longer (80–100 words). Link your sentences with und, aber, weil, dann, deshalb."),
		write("Write an email to a friend about your last weekend: where you went, what you did, and how it was (use the Perfekt).",
			"Liebe Sara,\nletztes Wochenende war toll! Ich bin nach Hamburg gefahren und habe meine Freundin besucht. Am Samstag haben wir die Stadt angesehen und sind am Abend ins Kino gegangen. Das Wetter war schön und das Essen hat gut geschmeckt. Wie war dein Wochenende?\nViele Grüße,\nLena"),
		write("Write an email to invite a friend on holiday next month: where, when, and why it will be nice.",
			"Lieber Tom,\nim nächsten Monat fahre ich an die Ostsee. Möchtest du mitkommen? Wir werden am Strand liegen und schwimmen. Ich glaube, es wird sehr schön, weil das Wetter warm ist. Bitte sag mir bis Freitag Bescheid.\nViele Grüße,\nMax"),
		write("Write a complaint email to a shop: the product is broken and you want a refund.",
			"Sehr geehrte Damen und Herren,\nich habe bei Ihnen eine Lampe gekauft, aber sie ist leider kaputt. Ich möchte mein Geld zurück oder eine neue Lampe. Bitte antworten Sie mir bald.\nMit freundlichen Grüßen,\nAnna Rossi"),
	)
	addVocab(db, l,
		vw("deshalb", "therefore", "Es regnet, deshalb bleibe ich.", "It rains, so I stay.", finch),
		vw("dann", "then", "Erst esse ich, dann lerne ich.", "First I eat, then I study.", "Cora"),
		vw("der Strand", "the beach", "am Strand", "at the beach", "Cora"),
		vw("kaputt", "broken", "Es ist kaputt.", "It's broken.", finch),
	)

	// Skill 29 — Speaking (A2 conversation)
	s = addSkillL(db, de, u2, "Sprechen: Gespräch", "Role-plays & opinions.", "MessageCircle", "#6C3FC5", 29, 372)
	addLesson(db, s, "Im Gespräch", 1, 20,
		char(finch, "A2 speaking: describe routines, talk about the past, give opinions with 'Ich finde …, weil …', and handle role-plays."),
		mc("Give an opinion", "How to start an opinion?", "Ich finde …", "Ich finde …", "Ich heiße …", "Ich komme …", "Ich habe …"),
		speak("Blaze", "Am Wochenende bin ich zu Hause geblieben und habe gelesen."),
		speak("Blaze", "Ich finde Berlin toll, weil die Stadt sehr lebendig ist."),
		speak("Blaze", "Entschuldigung, wo finde ich die Brötchen? Was kostet das?"),
	)
	addVocab(db, l,
		vw("Ich finde …", "I think …", "Ich finde das gut.", "I think that's good.", finch),
		vw("meiner Meinung nach", "in my opinion", "Meiner Meinung nach …", "In my opinion …", "Cora"),
		vw("der Termin", "appointment", "Ich habe einen Termin.", "I have an appointment.", "Cora"),
		vw("die Gesundheit", "health", "Gesundheit ist wichtig.", "Health is important.", finch),
	)

	addListeningL(db, de, u2, "Wie war dein Wochenende?",
		"Riko and Cora talk about their weekend. Listen, then answer.", 2, 20,
		[]models.ListeningMatch{
			lm("das Wochenende", "the weekend"),
			lm("Ich bin gefahren", "I drove/went"),
			lm("Ich habe gesehen", "I saw"),
			lm("Es war toll", "It was great"),
		},
		[]models.ListeningLine{
			ln("Riko", "Hey Cora, wie war dein Wochenende?", "Hey Cora, how was your weekend?"),
			ln("Cora", "Es war toll! Ich bin nach Hamburg gefahren.", "It was great! I went to Hamburg."),
			ln("Riko", "Schön! Was hast du dort gemacht?", "Nice! What did you do there?"),
			ln("Cora", "Ich habe Freunde getroffen und das Museum besucht.", "I met friends and visited the museum."),
			ln("Riko", "Ich habe nur gearbeitet und zu viel gegessen.", "I only worked and ate too much."),
			ln("Cora", "Armer Riko! Nächstes Mal kommst du mit.", "Poor Riko! Next time you come along."),
		},
		[]models.ListeningQuestion{
			lq("Where did Cora go?", "Hamburg", "Hamburg", "Berlin", "Munich", "Köln"),
			lq("What did Cora do there?", "Met friends, visited a museum", "Met friends, visited a museum", "Worked", "Slept", "Studied"),
			lq("What did Riko do?", "Worked and ate too much", "Worked and ate too much", "Travelled", "Danced", "Read"),
		},
	)
	addReadingL(db, de, u2, "Ein Brief aus dem Urlaub",
		"Read the postcard, then answer.", 2, 20,
		[]models.ReadingLine{
			rl("Liebe Mira, viele Grüße aus München!", "Dear Mira, greetings from Munich!"),
			rl("Wir sind am Freitag angekommen und haben sofort die Altstadt besucht.", "We arrived on Friday and immediately visited the old town."),
			rl("Gestern haben wir im Englischen Garten ein Picknick gemacht.", "Yesterday we had a picnic in the English Garden."),
			rl("Das Wetter war herrlich und das Essen hat fantastisch geschmeckt.", "The weather was lovely and the food tasted fantastic."),
			rl("Morgen fahren wir weiter nach Salzburg. Bis bald!", "Tomorrow we travel on to Salzburg. See you soon!"),
		},
		[]models.ReadingQuestion{
			rq("When did they arrive?", "On Friday", "On Friday", "On Sunday", "On Monday", "Yesterday"),
			rq("What did they do yesterday?", "A picnic", "A picnic", "A museum", "Shopping", "A concert"),
			rq("Where are they going next?", "Salzburg", "Salzburg", "Berlin", "Vienna", "Home"),
		},
	)
	addListeningL(db, de, u2, "Beim Arzt",
		"Lumora visits the doctor. Listen, then answer.", 3, 20,
		[]models.ListeningMatch{
			lm("Was fehlt Ihnen?", "What's wrong?"),
			lm("Ich habe Kopfschmerzen", "I have a headache"),
			lm("seit gestern", "since yesterday"),
			lm("dreimal am Tag", "three times a day"),
		},
		[]models.ListeningLine{
			ln(finch, "Guten Tag, Frau Lumora. Was fehlt Ihnen?", "Hello, Ms Lumora. What's wrong?"),
			ln("Lumora", "Guten Tag, Herr Doktor. Ich habe Kopfschmerzen und Fieber.", "Hello, Doctor. I have a headache and a fever."),
			ln(finch, "Seit wann haben Sie die Beschwerden?", "Since when have you had the symptoms?"),
			ln("Lumora", "Seit gestern. Ich konnte nicht schlafen.", "Since yesterday. I couldn't sleep."),
			ln(finch, "Nehmen Sie diese Tabletten, dreimal am Tag.", "Take these tablets, three times a day."),
			ln("Lumora", "Vielen Dank. Wann soll ich wiederkommen?", "Thank you. When should I come back?"),
		},
		[]models.ListeningQuestion{
			lq("What's wrong with Lumora?", "Headache and fever", "Headache and fever", "A cold", "A broken arm", "Toothache"),
			lq("Since when?", "Since yesterday", "Since yesterday", "A week", "An hour", "Two days"),
			lq("How often should she take the tablets?", "Three times a day", "Three times a day", "Once", "Twice", "Every hour"),
		},
	)
	addListeningL(db, de, u2, "Durchsage am Bahnhof",
		"A station announcement. Listen, then answer.", 4, 20,
		[]models.ListeningMatch{
			lm("die Durchsage", "the announcement"),
			lm("Gleis", "platform"),
			lm("Verspätung", "delay"),
			lm("der Zug", "the train"),
		},
		[]models.ListeningLine{
			ln("Mira", "Achtung am Gleis drei: Der Zug nach Hamburg hat zehn Minuten Verspätung.", "Attention at platform three: the train to Hamburg is ten minutes late."),
			ln("Mira", "Der Zug nach München fährt heute von Gleis fünf.", "The train to Munich departs today from platform five."),
			ln("Riko", "Oh nein, unser Zug nach Hamburg ist zu spät.", "Oh no, our train to Hamburg is late."),
			ln("Cora", "Kein Problem, wir haben Zeit. Welches Gleis ist es?", "No problem, we have time. Which platform is it?"),
			ln("Riko", "Gleis drei. Wir warten noch zehn Minuten.", "Platform three. We wait ten more minutes."),
		},
		[]models.ListeningQuestion{
			lq("Which train is delayed?", "To Hamburg", "To Hamburg", "To Munich", "To Berlin", "To Köln"),
			lq("How long is the delay?", "10 minutes", "10 minutes", "5 minutes", "1 hour", "20 minutes"),
			lq("From which platform does the Munich train leave?", "Platform 5", "Platform 5", "Platform 3", "Platform 1", "Platform 9"),
		},
	)
	addReadingL(db, de, u2, "Wohnungsanzeige",
		"Read the flat advert, then answer.", 3, 20,
		[]models.ReadingLine{
			rl("Schöne 2-Zimmer-Wohnung in Berlin-Mitte ab 1. September zu vermieten.", "Nice 2-room flat in Berlin-Mitte to let from 1 September."),
			rl("Die Wohnung ist 55 m² groß, hell und ruhig, mit Balkon und Küche.", "The flat is 55 m², bright and quiet, with a balcony and kitchen."),
			rl("Die Miete kostet 800 Euro im Monat plus Nebenkosten.", "The rent is 800 euros a month plus utilities."),
			rl("Interessenten melden sich bitte per E-Mail bei Frau Weber.", "Interested parties please contact Ms Weber by email."),
		},
		[]models.ReadingQuestion{
			rq("How many rooms?", "Two", "Two", "One", "Three", "Four"),
			rq("How much is the rent?", "800 euros", "800 euros", "550 euros", "1000 euros", "80 euros"),
			rq("From when is it available?", "1 September", "1 September", "1 August", "1 October", "Now"),
		},
	)
	addReadingL(db, de, u2, "Eine Einladung",
		"Read the invitation email, then answer.", 4, 20,
		[]models.ReadingLine{
			rl("Liebe Freunde, ich möchte euch zu meiner Geburtstagsfeier einladen.", "Dear friends, I'd like to invite you to my birthday party."),
			rl("Die Feier ist am Samstag, dem 12. Mai, ab 19 Uhr bei mir zu Hause.", "The party is on Saturday, 12 May, from 7 p.m. at my place."),
			rl("Es gibt Essen und Musik. Bitte bringt gute Laune mit!", "There will be food and music. Please bring good spirits!"),
			rl("Sagt mir bitte bis Mittwoch, ob ihr kommen könnt.", "Please tell me by Wednesday whether you can come."),
		},
		[]models.ReadingQuestion{
			rq("What is the event?", "A birthday party", "A birthday party", "A wedding", "A meeting", "A concert"),
			rq("When is it?", "Saturday 12 May, 7 p.m.", "Saturday 12 May, 7 p.m.", "Sunday 10 a.m.", "Friday 8 p.m.", "Monday noon"),
			rq("By when should you reply?", "Wednesday", "Wednesday", "Monday", "Saturday", "Friday"),
		},
	)
	addReadingL(db, de, u2, "Lerntipps für A2",
		"A smart study plan for A2. Read, then answer.", 5, 20,
		[]models.ReadingLine{
			rl("Lerne jeden Tag 30 Minuten — das ist besser als einmal pro Woche drei Stunden.", "Study 30 minutes every day — that's better than three hours once a week."),
			rl("Wiederhole zuerst deine A1-Fehler und konzentriere dich auf Perfekt und Dativ.", "First review your A1 mistakes and focus on the Perfekt and the dative."),
			rl("Übe alle vier Fertigkeiten: Hören, Lesen, Schreiben und Sprechen.", "Practise all four skills: listening, reading, writing and speaking."),
			rl("Lerne neue Wörter in Themen, zum Beispiel Gesundheit, Reisen und Arbeit.", "Learn new words in topics, for example health, travel and work."),
			rl("Mach jede Woche eine kleine Prüfung und sprich so oft wie möglich.", "Do a small test every week and speak as often as possible."),
		},
		[]models.ReadingQuestion{
			rq("How long should you study daily?", "30 minutes", "30 minutes", "3 hours", "5 minutes", "All day"),
			rq("Which grammar should you focus on?", "Perfekt and dative", "Perfekt and dative", "Only articles", "Konjunktiv", "Passive"),
			rq("How should you learn new words?", "In topics", "In topics", "Randomly", "Alphabetically", "Backwards"),
		},
	)

	// ───────────────────────── B1 · Mittelstufe ────────────────────────
	u3 := "B1 · Mittelstufe"

	// 30 — Genitive
	s = addSkillL(db, de, u3, "Der Genitiv", "Possession + genitive prepositions.", "Hash", "#06AECE", 30, 400)
	l = addLesson(db, s, "Genitiv", 1, 20,
		char(finch, "The genitive shows possession: 'das Auto des Mannes', 'die Farbe der Blume'. Masculine/neuter add -(e)s. Genitive prepositions: wegen, trotz, während, statt."),
		mc("Genitive (masc.)", "Das ist das Haus ___ Lehrers. (the teacher's)", "des", "des", "der", "dem", "den"),
		fill("Genitive (fem.)", "Die Tasche ___ Frau. (of the woman → der)", "der"),
		mc("Genitive preposition", "___ des Regens bleiben wir zu Hause. (because of)", "Wegen", "Wegen", "Weil", "Während", "Trotz"),
		tr("Translate", "the title of the book", "der Titel des Buches"),
		speak("Blaze", "Trotz des Wetters gehen wir spazieren."),
	)
	addVocab(db, l,
		vw("wegen", "because of (+gen)", "wegen des Regens", "because of the rain", finch),
		vw("trotz", "despite (+gen)", "trotz des Wetters", "despite the weather", finch),
		vw("während", "during (+gen)", "während des Tages", "during the day", "Cora"),
		vw("der Titel", "the title", "der Titel des Buches", "the title of the book", "Cora"),
	)

	// 31 — Plusquamperfekt
	s = addSkillL(db, de, u3, "Plusquamperfekt", "The past before the past.", "Clock", "#6C3FC5", 31, 430)
	l = addLesson(db, s, "Plusquamperfekt", 1, 20,
		char(finch, "The Plusquamperfekt describes what happened BEFORE another past event: hatte/war + Partizip II. Often with 'nachdem'."),
		fill("haben → past", "Nachdem ich gegessen ___ , ging ich los. (had)", "hatte"),
		mc("sein verb", "Er ___ schon gegangen, als ich kam. (had gone)", "war", "war", "hatte", "ist", "hat"),
		tr("Translate", "I had already eaten", "Ich hatte schon gegessen"),
		fill("Participle", "Sie war müde, weil sie schlecht ___ hatte. (slept)", "geschlafen"),
		speak("Blaze", "Nachdem wir gegessen hatten, gingen wir spazieren."),
	)
	addVocab(db, l,
		vw("nachdem", "after (conj.)", "Nachdem ich aß, …", "After I ate, …", finch),
		vw("schon", "already", "Ich war schon da.", "I was already there.", "Cora"),
		vw("vorher", "before(hand)", "vorher", "beforehand", "Cora"),
		vw("der Moment", "the moment", "in diesem Moment", "at this moment", finch),
	)

	// 32 — Konjunktiv II
	s = addSkillL(db, de, u3, "Konjunktiv II", "Wishes, politeness, hypotheticals.", "Sparkles", "#F5A623", 32, 460)
	l = addLesson(db, s, "Konjunktiv II", 1, 20,
		char(finch, "Konjunktiv II = politeness, wishes and hypotheticals: würde + infinitive, or hätte/wäre/könnte."),
		fill("würde-form", "Ich ___ gern nach Japan reisen. (would)", "würde"),
		mc("hätte", "Wenn ich Zeit ___ , käme ich. (had)", "hätte", "hätte", "habe", "hatte", "hat"),
		mc("Polite request", "___ Sie mir bitte helfen? (could)", "Könnten", "Könnten", "Können", "Konnten", "Kann"),
		tr("Translate", "I would help you", "Ich würde dir helfen"),
		speak("Blaze", "An deiner Stelle würde ich warten."),
	)
	addVocab(db, l,
		vw("würde", "would", "Ich würde gehen.", "I would go.", finch),
		vw("hätte", "would have / had", "Wenn ich Geld hätte …", "If I had money …", finch),
		vw("wäre", "would be", "Das wäre toll.", "That would be great.", "Cora"),
		vw("an deiner Stelle", "if I were you", "An deiner Stelle …", "If I were you …", "Cora"),
	)

	// 33 — Verbs with prepositions
	s = addSkillL(db, de, u3, "Verben mit Präpositionen", "warten auf, denken an …", "Link2", "#00C2A8", 33, 490)
	l = addLesson(db, s, "Verben + Präposition", 1, 20,
		char(finch, "Many verbs take a fixed preposition and case: warten auf (+Akk), denken an (+Akk), sprechen über (+Akk), helfen bei (+Dat)."),
		mc("warten ___", "Ich warte ___ den Bus. (for)", "auf", "auf", "an", "über", "für"),
		fill("denken ___", "Ich denke oft ___ dich. (of/about)", "an"),
		mc("sprechen ___", "Wir sprechen ___ das Wetter. (about)", "über", "über", "auf", "an", "von"),
		tr("Translate", "I am interested in music", "Ich interessiere mich für Musik"),
		speak("Blaze", "Ich freue mich auf das Wochenende."),
	)
	addVocab(db, l,
		vw("warten auf", "to wait for", "Ich warte auf dich.", "I'm waiting for you.", finch),
		vw("denken an", "to think of", "Ich denke an dich.", "I think of you.", finch),
		vw("sich freuen auf", "to look forward to", "Ich freue mich auf …", "I look forward to …", "Cora"),
		vw("sprechen über", "to talk about", "Wir sprechen über …", "We talk about …", "Cora"),
	)

	// 34 — Process passive
	s = addSkillL(db, de, u3, "Vorgangspassiv", "werden + Partizip II.", "Layers", "#FF5C5C", 34, 520)
	l = addLesson(db, s, "Das Passiv", 1, 20,
		char(finch, "Process passive: werden + Partizip II — focus on the action: 'Das Haus wird gebaut.' Past: 'wurde gebaut.'"),
		fill("Passive (singular)", "Das Haus ___ gebaut. (is being built)", "wird"),
		mc("Passive (plural)", "Die Briefe ___ geschrieben.", "werden", "werden", "wird", "ist", "sind"),
		fill("Passive (past)", "Gestern ___ das Haus gebaut. (was built)", "wurde"),
		tr("Translate", "The book is read", "Das Buch wird gelesen"),
		speak("Blaze", "Hier wird Deutsch gesprochen."),
	)
	addVocab(db, l,
		vw("das Passiv", "the passive", "das Vorgangspassiv", "the process passive", finch),
		vw("gebaut", "built", "Es wird gebaut.", "It's being built.", "Cora"),
		vw("hergestellt", "produced", "Es wird hergestellt.", "It's produced.", "Cora"),
		vw("der Vorgang", "the process", "der Vorgang", "the process", finch),
	)

	// 35 — State passive
	s = addSkillL(db, de, u3, "Zustandspassiv", "sein + Partizip II (result).", "Layers", "#17A3DD", 35, 550)
	l = addLesson(db, s, "Zustandspassiv", 1, 20,
		char(finch, "State passive: sein + Partizip II describes a RESULT/state. 'Die Tür wird geschlossen' (action) vs 'Die Tür ist geschlossen' (state)."),
		mc("State (closed)", "Die Tür ist ___ .", "geschlossen", "geschlossen", "schließen", "schließt", "schließe"),
		fill("State (open)", "Das Geschäft ist ___ . (open → geöffnet)", "geöffnet"),
		mc("Which is a state?", "Choose the state passive", "ist geschlossen", "ist geschlossen", "wird geschlossen", "schließt", "geschlossen werden"),
		tr("Translate", "The shop is closed", "Das Geschäft ist geschlossen"),
		speak("Blaze", "Das Fenster ist geöffnet."),
	)
	addVocab(db, l,
		vw("geschlossen", "closed", "Es ist geschlossen.", "It's closed.", finch),
		vw("geöffnet", "open", "Es ist geöffnet.", "It's open.", "Cora"),
		vw("repariert", "repaired", "Es ist repariert.", "It's repaired.", "Cora"),
		vw("fertig", "ready / done", "Es ist fertig.", "It's done.", finch),
	)

	// 36 — Relative clauses
	s = addSkillL(db, de, u3, "Relativsätze", "der/die/das as relative pronouns.", "Link2", "#6C3FC5", 36, 580)
	l = addLesson(db, s, "Relativsätze", 1, 20,
		char(finch, "Relative pronouns match the noun in gender & number but take their case from their role in the clause."),
		mc("masc. nominative", "Der Mann, ___ dort steht, ist mein Lehrer.", "der", "der", "den", "dem", "die"),
		fill("neut. accusative", "Das Buch, ___ ich lese, ist gut.", "das"),
		mc("fem. dative", "Die Frau, ___ ich helfe, ist nett.", "der", "der", "die", "den", "dem"),
		tr("Translate", "the man who works here", "der Mann, der hier arbeitet"),
		speak("Blaze", "Das ist der Freund, der mir geholfen hat."),
	)
	addVocab(db, l,
		vw("der Lehrer", "teacher", "Der Lehrer erklärt.", "The teacher explains.", finch),
		vw("der Nachbar", "neighbour", "mein Nachbar", "my neighbour", "Cora"),
		vw("spannend", "exciting", "ein spannendes Buch", "an exciting book", "Cora"),
		vw("erklären", "to explain", "Ich erkläre es.", "I explain it.", finch),
	)

	// 37 — Advanced conjunctions
	s = addSkillL(db, de, u3, "Konjunktionen", "obwohl, damit, seitdem, als.", "MessageCircle", "#00C2A8", 37, 610)
	l = addLesson(db, s, "Bindewörter", 1, 20,
		char(finch, "Advanced subordinators: obwohl (although), damit (so that), seitdem (since), als (when — single past event), während (while)."),
		mc("although", "___ es regnete, gingen wir. (although)", "Obwohl", "Obwohl", "Weil", "Wenn", "Damit"),
		mc("so that", "Ich lerne, ___ ich die Prüfung bestehe. (so that)", "damit", "damit", "weil", "obwohl", "dass"),
		mc("when (past, once)", "___ ich klein war, wohnte ich in Bonn. (when)", "Als", "Als", "Wenn", "Wann", "Ob"),
		tr("Translate", "since I live here", "seitdem ich hier wohne"),
		speak("Blaze", "Obwohl es regnet, gehe ich spazieren."),
	)
	addVocab(db, l,
		vw("obwohl", "although", "Obwohl es regnet …", "Although it rains …", finch),
		vw("damit", "so that", "…, damit ich lerne.", "…, so that I learn.", finch),
		vw("seitdem", "since (time)", "seitdem ich hier bin", "since I've been here", "Cora"),
		vw("als", "when (past)", "Als ich jung war …", "When I was young …", "Cora"),
	)

	// 38 — Infinitive with zu
	s = addSkillL(db, de, u3, "Infinitiv mit zu", "… zu … and um … zu.", "PenLine", "#F5A623", 38, 640)
	l = addLesson(db, s, "Infinitiv mit zu", 1, 20,
		char(finch, "Use 'zu + infinitive' after many verbs/expressions: 'Ich versuche zu lernen.' Purpose: 'um … zu'."),
		fill("zu + infinitive", "Ich versuche, Deutsch ___ lernen.", "zu"),
		mc("It is important to …", "Es ist wichtig, früh ___ kommen.", "zu", "zu", "für", "um", "dass"),
		mc("in order to", "Ich gehe in die Stadt, ___ einzukaufen.", "um", "um", "zu", "damit", "für"),
		tr("Translate", "I have no time to sleep", "Ich habe keine Zeit zu schlafen"),
		speak("Blaze", "Ich hoffe, dich bald zu sehen."),
	)
	addVocab(db, l,
		vw("versuchen", "to try", "Ich versuche es.", "I try it.", finch),
		vw("hoffen", "to hope", "Ich hoffe es.", "I hope so.", "Cora"),
		vw("vergessen", "to forget", "Vergiss nicht zu …", "Don't forget to …", "Cora"),
		vw("um … zu", "in order to", "um zu lernen", "in order to learn", finch),
	)

	// 39 — Participles as adjectives
	s = addSkillL(db, de, u3, "Partizipien als Adjektive", "the laughing man, the boiled egg.", "Sparkles", "#FF5C5C", 39, 670)
	l = addLesson(db, s, "Partizip als Adjektiv", 1, 20,
		char(finch, "Participles work as adjectives: Partizip I (-end) is active — 'der lachende Mann'; Partizip II is passive/done — 'das gekochte Ei'. They take adjective endings."),
		mc("Partizip I", "der ___ Mann (laughing)", "lachende", "lachende", "lachen", "gelacht", "lacht"),
		mc("Partizip II", "das ___ Ei (boiled)", "gekochte", "gekochte", "kochen", "kocht", "kochend"),
		fill("Partizip II", "die ___ Tür (closed → geschlossene)", "geschlossene"),
		tr("Translate", "the sleeping child", "das schlafende Kind"),
		speak("Blaze", "Das ist ein spannendes Buch."),
	)
	addVocab(db, l,
		vw("lachend", "laughing", "das lachende Kind", "the laughing child", finch),
		vw("kochend", "boiling", "kochendes Wasser", "boiling water", "Cora"),
		vw("gekocht", "cooked / boiled", "ein gekochtes Ei", "a boiled egg", "Cora"),
		vw("spannend", "exciting", "ein spannender Film", "an exciting film", finch),
	)

	// 40 — Modal particles
	s = addSkillL(db, de, u3, "Modalpartikeln", "ja, doch, mal, denn.", "Quote", "#17A3DD", 40, 700)
	l = addLesson(db, s, "Modalpartikeln", 1, 20,
		char(finch, "Modal particles add tone. 'Komm mal her!' (softens), 'Das ist ja toll!' (surprise), 'Warum denn?' (curiosity), 'Komm doch!' (encourages)."),
		mc("softening a request", "Komm ___ her!", "mal", "mal", "ja", "denn", "zu"),
		mc("surprise", "Das ist ___ toll!", "ja", "ja", "mal", "denn", "zu"),
		mc("in questions", "Warum ___ ?", "denn", "denn", "mal", "ja", "doch"),
		tr("Translate (do come in!)", "Do come in", "Komm doch rein"),
		speak("Blaze", "Das ist ja super! Komm doch mal vorbei."),
	)
	addVocab(db, l,
		vw("ja (Partikel)", "(emphasis/surprise)", "Das ist ja schön!", "That's really nice!", finch),
		vw("doch (Partikel)", "(encourage/contrast)", "Komm doch!", "Do come!", finch),
		vw("mal (Partikel)", "(softens)", "Sag mal …", "Say …", "Cora"),
		vw("denn (Partikel)", "(in questions)", "Was machst du denn?", "What are you doing?", "Cora"),
	)

	// 41 — Writing (B1 essay/letter)
	s = addSkillL(db, de, u3, "Schreiben: Aufsatz", "Opinion essays & letters (150–200 words).", "PenLine", "#00C2A8", 41, 730)
	l = addLesson(db, s, "Meinung & Brief", 1, 20,
		char(finch, "B1 writing has structure: introduction, body (reasons & examples), conclusion. Connectors: außerdem, deshalb, obwohl, zum Beispiel, einerseits/andererseits."),
		write("Write your opinion (≈150 words): 'Soll man im Homeoffice arbeiten?' Give advantages, disadvantages, and your own view.",
			"Heutzutage arbeiten viele Menschen im Homeoffice. Meiner Meinung nach hat das Vor- und Nachteile.\nEinerseits spart man Zeit, weil man nicht zur Arbeit fahren muss. Außerdem kann man flexibler arbeiten. Andererseits fehlt der Kontakt zu den Kollegen, und manche Menschen können sich zu Hause schlecht konzentrieren.\nIch persönlich finde das Homeoffice gut, aber nicht jeden Tag. Deshalb wäre eine Mischung aus Büro und Homeoffice ideal."),
		write("Write a formal complaint letter: you bought a phone online, it arrived broken, and you want a replacement or refund.",
			"Sehr geehrte Damen und Herren,\nam 3. Mai habe ich bei Ihnen ein Handy bestellt. Leider war das Gerät bei der Lieferung kaputt. Ich habe es sofort fotografiert.\nIch möchte Sie bitten, mir ein neues Handy zu schicken oder mir das Geld zurückzugeben. Bitte antworten Sie mir bis Ende der Woche.\nMit freundlichen Grüßen,\nKwame Mensah"),
		write("Write an email to a friend about your future plans after the B1 exam (study, work, travel) and ask about theirs.",
			"Liebe Sara,\nbald mache ich die B1-Prüfung und ich bin schon aufgeregt. Danach möchte ich in Deutschland studieren, weil die Universitäten gut sind. Vielleicht reise ich vorher noch ein bisschen.\nUnd du? Was sind deine Pläne nach der Prüfung? Schreib mir bald!\nViele Grüße,\nLena"),
	)
	addVocab(db, l,
		vw("außerdem", "besides / moreover", "Außerdem ist es billig.", "Besides, it's cheap.", finch),
		vw("der Vorteil", "advantage", "ein großer Vorteil", "a big advantage", "Cora"),
		vw("der Nachteil", "disadvantage", "ein Nachteil ist …", "a disadvantage is …", "Cora"),
		vw("einerseits", "on one hand", "einerseits … andererseits", "on one hand … on the other", finch),
	)

	// 42 — Speaking (B1 discussion)
	s = addSkillL(db, de, u3, "Sprechen: Diskussion", "Opinions, pros/cons, planning.", "MessageCircle", "#6C3FC5", 42, 760)
	addLesson(db, s, "Diskutieren", 1, 20,
		char(finch, "B1 speaking: present a topic, justify opinions, weigh pros & cons, and plan together. Use 'Meiner Meinung nach …', 'einerseits … andererseits', 'Ich schlage vor, dass …'."),
		mc("Formal opinion phrase", "How to state an opinion formally?", "Meiner Meinung nach …", "Meiner Meinung nach …", "Ich heiße …", "Ich komme …", "Ich habe …"),
		speak("Blaze", "Meiner Meinung nach ist Umweltschutz sehr wichtig, weil wir nur eine Erde haben."),
		speak("Blaze", "Einerseits ist das Auto bequem, andererseits ist es schlecht für die Umwelt."),
		speak("Blaze", "Ich schlage vor, dass wir das Problem zusammen lösen."),
	)
	addVocab(db, l,
		vw("die Umwelt", "the environment", "Umwelt schützen", "to protect the environment", finch),
		vw("der Vorschlag", "the suggestion", "Ich habe einen Vorschlag.", "I have a suggestion.", "Cora"),
		vw("vorschlagen", "to suggest", "Ich schlage vor …", "I suggest …", "Cora"),
		vw("die Gesellschaft", "society", "in der Gesellschaft", "in society", finch),
	)

	addListeningL(db, de, u3, "Interview: Umwelt",
		"Mira interviews Riko about the environment. Listen, then answer.", 3, 25,
		[]models.ListeningMatch{
			lm("die Umwelt", "the environment"),
			lm("der Klimawandel", "climate change"),
			lm("öffentliche Verkehrsmittel", "public transport"),
			lm("meiner Meinung nach", "in my opinion"),
		},
		[]models.ListeningLine{
			ln("Mira", "Riko, was ist für dich das größte Umweltproblem?", "Riko, what's the biggest environmental problem for you?"),
			ln("Riko", "Meiner Meinung nach ist der Klimawandel das größte Problem.", "In my opinion, climate change is the biggest problem."),
			ln("Mira", "Und was kann man dagegen tun?", "And what can be done about it?"),
			ln("Riko", "Man sollte öfter öffentliche Verkehrsmittel benutzen und weniger fliegen.", "One should use public transport more often and fly less."),
			ln("Mira", "Machst du das selbst auch?", "Do you do that yourself too?"),
			ln("Riko", "Ja, ich fahre meistens mit dem Fahrrad, obwohl es manchmal anstrengend ist.", "Yes, I usually cycle, although it's sometimes tiring."),
		},
		[]models.ListeningQuestion{
			lq("What is the biggest problem for Riko?", "Climate change", "Climate change", "Traffic", "Noise", "Rubbish"),
			lq("What does he suggest?", "Use public transport, fly less", "Use public transport, fly less", "Buy a car", "Eat more", "Work less"),
			lq("How does Riko usually travel?", "By bicycle", "By bicycle", "By car", "By plane", "By taxi"),
		},
	)
	addListeningL(db, de, u3, "Ein Problem im Büro",
		"A workplace conversation. Listen, then answer.", 4, 25,
		[]models.ListeningMatch{
			lm("die Besprechung", "the meeting"),
			lm("der Termin", "the deadline / appointment"),
			lm("Ich kümmere mich darum", "I'll take care of it"),
			lm("kein Problem", "no problem"),
		},
		[]models.ListeningLine{
			ln("Cora", "Riko, der Kunde hat angerufen. Das Projekt ist nicht fertig.", "Riko, the client called. The project isn't finished."),
			ln("Riko", "Das weiß ich. Der Termin ist erst morgen, oder?", "I know. The deadline is only tomorrow, right?"),
			ln("Cora", "Nein, der Termin ist heute Abend. Wir haben ein Problem.", "No, the deadline is this evening. We have a problem."),
			ln("Riko", "Okay, kein Problem. Ich kümmere mich darum und bleibe länger.", "Okay, no problem. I'll take care of it and stay longer."),
			ln("Cora", "Danke! Soll ich dir helfen?", "Thanks! Should I help you?"),
			ln("Riko", "Ja, gern. Zusammen schaffen wir das.", "Yes, gladly. Together we'll manage it."),
		},
		[]models.ListeningQuestion{
			lq("What is the problem?", "The project isn't finished", "The project isn't finished", "Someone is sick", "The computer is broken", "No coffee"),
			lq("When is the deadline?", "This evening", "This evening", "Tomorrow", "Next week", "Friday"),
			lq("What does Riko decide?", "To stay longer and finish it", "To stay longer and finish it", "To go home", "To cancel", "To call the client"),
		},
	)
	addReadingL(db, de, u3, "Artikel: Die Stadt der Zukunft",
		"Read the article, then answer.", 3, 25,
		[]models.ReadingLine{
			rl("Wie sieht die Stadt der Zukunft aus? Experten sind sich einig, dass sie grüner und ruhiger sein muss.", "What does the city of the future look like? Experts agree that it must be greener and quieter."),
			rl("In vielen Städten werden heute neue Radwege gebaut, damit weniger Autos fahren.", "In many cities new cycle paths are being built today, so that fewer cars drive."),
			rl("Außerdem sollen mehr Parks entstehen, weil Grünflächen gut für die Gesundheit sind.", "Moreover, more parks should be created, because green spaces are good for health."),
			rl("Obwohl diese Pläne teuer sind, halten viele Bürger sie für eine gute Investition in die Zukunft.", "Although these plans are expensive, many citizens consider them a good investment in the future."),
		},
		[]models.ReadingQuestion{
			rq("How should the future city be?", "Greener and quieter", "Greener and quieter", "Bigger and louder", "Cheaper", "Faster"),
			rq("Why build cycle paths?", "So fewer cars drive", "So fewer cars drive", "To save money", "For tourists", "For sport"),
			rq("How do citizens view the plans?", "A good investment", "A good investment", "A waste", "Too cheap", "Pointless"),
		},
	)
	addReadingL(db, de, u3, "Lerntipps für B1",
		"A smart study plan for B1. Read, then answer.", 4, 25,
		[]models.ReadingLine{
			rl("Auf dem B1-Niveau musst du die Sprache aktiv produzieren, nicht nur verstehen.", "At B1 level you must actively produce the language, not just understand it."),
			rl("Sprich und schreibe jeden Tag — zum Beispiel ein kurzer Text über deine Meinung.", "Speak and write every day — for example a short text about your opinion."),
			rl("Lerne Wörter in Themen wie Arbeit, Umwelt und Gesundheit, und übe Konnektoren.", "Learn words in topics like work, environment and health, and practise connectors."),
			rl("Mach jede Woche eine vollständige Modellprüfung, damit du das Format gut kennst.", "Do a full model exam every week, so that you know the format well."),
			rl("Wer regelmäßig und mit Plan lernt, besteht die Prüfung leichter.", "Whoever studies regularly and with a plan passes the exam more easily."),
		},
		[]models.ReadingQuestion{
			rq("What must you do at B1?", "Produce the language actively", "Produce the language actively", "Only listen", "Only read", "Memorise lists"),
			rq("How often should you do a model exam?", "Every week", "Every week", "Once a year", "Never", "Daily"),
			rq("Which topics are recommended?", "Work, environment, health", "Work, environment, health", "Only food", "Only travel", "Only sport"),
		},
	)

	// ───────────────────────── B2 · Oberstufe ──────────────────────────
	u4 := "B2 · Oberstufe"

	// 44 — Konjunktiv I (reported speech)
	s = addSkillL(db, de, u4, "Konjunktiv I", "Reported speech.", "Quote", "#FF5C5C", 44, 850)
	l = addLesson(db, s, "Indirekte Rede", 1, 25,
		char(finch, "Formal reported speech uses Konjunktiv I: 'Er sagt, er habe keine Zeit.' For sein → sei."),
		mc("Konjunktiv I of haben", "Er sagt, er ___ keine Zeit.", "habe", "habe", "hat", "hätte", "haben"),
		fill("Konjunktiv I of sein", "Sie sagte, sie ___ krank. (be → sei)", "sei"),
		mc("Reported speech", "Der Minister erklärte, die Lage ___ stabil.", "sei", "sei", "ist", "war", "wäre"),
		tr("Translate", "He says he has no time", "Er sagt, er habe keine Zeit"),
		speak("Blaze", "Sie behauptet, sie sei die Beste."),
	)
	addVocab(db, l,
		vw("die indirekte Rede", "reported speech", "in der indirekten Rede", "in reported speech", finch),
		vw("behaupten", "to claim", "Er behauptet das.", "He claims that.", finch),
		vw("angeblich", "allegedly", "Er ist angeblich krank.", "He is allegedly ill.", "Cora"),
		vw("erklären", "to state / explain", "Sie erklärte …", "She stated …", "Cora"),
	)

	// 45 — Konjunktiv II & conditionals
	s = addSkillL(db, de, u4, "Konditionalsätze", "Real & unreal conditions.", "Sparkles", "#F5A623", 45, 880)
	l = addLesson(db, s, "Konjunktiv II", 1, 25,
		char(finch, "Real conditions use the indicative ('Wenn es regnet, bleibe ich'); unreal use Konjunktiv II ('Wenn ich Zeit hätte, käme ich')."),
		mc("Unreal (sein)", "Wenn ich reich ___ , würde ich reisen.", "wäre", "wäre", "bin", "war", "bist"),
		fill("Unreal past (haben)", "Wenn ich das gewusst ___ , wäre ich gekommen.", "hätte"),
		mc("Polite request", "___ Sie mir bitte helfen?", "Könnten", "Könnten", "Können", "Konnten", "Kann"),
		tr("Translate", "If I had time, I would help", "Wenn ich Zeit hätte, würde ich helfen"),
		speak("Blaze", "Wenn ich das gewusst hätte, wäre ich gekommen."),
	)
	addVocab(db, l,
		vw("die Bedingung", "condition", "unter einer Bedingung", "on one condition", finch),
		vw("falls", "in case / if", "falls es regnet", "in case it rains", "Cora"),
		vw("sonst", "otherwise", "Beeil dich, sonst …", "Hurry, otherwise …", "Cora"),
		vw("beinahe", "almost", "beinahe vergessen", "almost forgotten", finch),
	)

	// 46 — Futur II
	s = addSkillL(db, de, u4, "Futur II", "The future perfect.", "Clock", "#6C3FC5", 46, 910)
	l = addLesson(db, s, "Futur II", 1, 25,
		char(finch, "Futur II = a completed future or an assumption: werden + Partizip II + haben/sein. 'Bis morgen werde ich es gemacht haben.'"),
		mc("werden (ich)", "Bis morgen ___ ich den Text gelesen haben.", "werde", "werde", "wird", "habe", "bin"),
		fill("Auxiliary (motion)", "Er wird angekommen ___ . (sein-verb)", "sein"),
		mc("Assumption", "Sie wird das wohl vergessen ___ .", "haben", "haben", "sein", "werden", "hat"),
		tr("Translate", "I will have read the book", "Ich werde das Buch gelesen haben"),
		speak("Blaze", "Bis Freitag werde ich alles erledigt haben."),
	)
	addVocab(db, l,
		vw("bis morgen", "by tomorrow", "bis morgen", "by tomorrow", finch),
		vw("wohl", "probably", "Das wird wohl so sein.", "That'll probably be so.", "Cora"),
		vw("erledigen", "to get done", "Ich erledige das.", "I'll get it done.", "Cora"),
		vw("vermutlich", "presumably", "vermutlich morgen", "presumably tomorrow", finch),
	)

	// 47 — Passive with modals
	s = addSkillL(db, de, u4, "Passiv mit Modalverben", "must/can be done.", "Layers", "#00C2A8", 47, 940)
	l = addLesson(db, s, "Modalpassiv", 1, 25,
		char(finch, "Passive with a modal: modal + Partizip II + werden. 'Das Problem muss gelöst werden.'"),
		mc("be (passive infinitive)", "Das Problem muss gelöst ___ .", "werden", "werden", "sein", "wird", "worden"),
		fill("Passive with können", "Die Tür kann nicht geöffnet ___ . (be → werden)", "werden"),
		mc("Past modal passive", "Das Haus musste neu gebaut ___ .", "werden", "werden", "wurde", "worden", "sein"),
		tr("Translate", "The letter must be written", "Der Brief muss geschrieben werden"),
		speak("Blaze", "Das muss heute noch gemacht werden."),
	)
	addVocab(db, l,
		vw("lösen", "to solve", "ein Problem lösen", "to solve a problem", finch),
		vw("möglich", "possible", "Das ist möglich.", "That's possible.", "Cora"),
		vw("nötig", "necessary", "Das ist nicht nötig.", "That's not necessary.", "Cora"),
		vw("erledigt", "done / settled", "Es ist erledigt.", "It's done.", finch),
	)

	// 48 — lassen + infinitive
	s = addSkillL(db, de, u4, "lassen + Infinitiv", "Have something done.", "Link2", "#17A3DD", 48, 970)
	l = addLesson(db, s, "lassen", 1, 25,
		char(finch, "'lassen' + infinitive = to have something done, or to let. 'Ich lasse mein Auto reparieren.'"),
		mc("lassen (ich)", "Ich ___ mein Auto reparieren.", "lasse", "lasse", "lässt", "lassen", "gelassen"),
		fill("lassen (sie)", "Sie ___ die Haare schneiden. (has them cut)", "lässt"),
		mc("Meaning of 'Ich lasse dich gehen'", "Choose the meaning", "I let you go", "I let you go", "I make you go", "I go away", "I leave you"),
		tr("Translate", "I'm having a house built", "Ich lasse ein Haus bauen"),
		speak("Blaze", "Ich lasse mein Auto reparieren."),
	)
	addVocab(db, l,
		vw("lassen", "to let / have done", "Ich lasse es machen.", "I have it done.", finch),
		vw("reparieren", "to repair", "reparieren lassen", "to have repaired", "Cora"),
		vw("schneiden", "to cut", "Haare schneiden", "to cut hair", "Cora"),
		vw("der Friseur", "hairdresser", "zum Friseur gehen", "to go to the hairdresser", finch),
	)

	// 49 — Double infinitive
	s = addSkillL(db, de, u4, "Doppelter Infinitiv", "Modals in the perfect.", "Layers", "#F5A623", 49, 1000)
	l = addLesson(db, s, "Doppelter Infinitiv", 1, 25,
		char(finch, "Modals in the Perfekt use a double infinitive: 'Ich habe arbeiten müssen' — not the participle 'gemusst' when there's a main verb."),
		mc("Double infinitive", "Ich habe gestern arbeiten ___ .", "müssen", "müssen", "gemusst", "musste", "muss"),
		fill("können in perfect", "Er hat nicht kommen ___ . (could not)", "können"),
		mc("Word order", "Ich habe das machen ___ .", "müssen", "müssen", "gemusst", "muss", "musste"),
		tr("Translate", "I have had to learn a lot", "Ich habe viel lernen müssen"),
		speak("Blaze", "Ich habe gestern lange arbeiten müssen."),
	)
	addVocab(db, l,
		vw("müssen", "to have to", "Ich habe gehen müssen.", "I had to go.", finch),
		vw("können", "to be able", "Ich habe es sehen können.", "I was able to see it.", "Cora"),
		vw("dürfen", "to be allowed", "Ich habe bleiben dürfen.", "I was allowed to stay.", "Cora"),
		vw("der Infinitiv", "infinitive", "der doppelte Infinitiv", "the double infinitive", finch),
	)

	// 50 — Adjectives without an article
	s = addSkillL(db, de, u4, "Adjektive ohne Artikel", "Strong declension.", "PenLine", "#FF5C5C", 50, 1030)
	l = addLesson(db, s, "Starke Deklination", 1, 25,
		char(finch, "Without an article, the adjective itself carries the case ending (strong declension): 'guter Wein, kaltes Wasser, mit großer Freude.'"),
		mc("Masc. nom.", "___ Wein (good wine — no article)", "guter", "guter", "gute", "gutes", "guten"),
		fill("Neut. nom.", "___ Wasser (cold water → kaltes)", "kaltes"),
		mc("Fem. dative", "mit ___ Freude (great joy)", "großer", "großer", "große", "großen", "großem"),
		tr("Translate", "fresh bread", "frisches Brot"),
		speak("Blaze", "Guter Wein und frisches Brot."),
	)
	addVocab(db, l,
		vw("frisch", "fresh", "frisches Brot", "fresh bread", finch),
		vw("kalt", "cold", "kaltes Wasser", "cold water", "Cora"),
		vw("die Freude", "joy", "mit großer Freude", "with great joy", "Cora"),
		vw("der Wein", "wine", "guter Wein", "good wine", finch),
	)

	// 51 — Relative clauses (advanced)
	s = addSkillL(db, de, u4, "Relativsätze (komplex)", "Prepositions, was, welcher.", "Link2", "#6C3FC5", 51, 1060)
	l = addLesson(db, s, "Komplexe Relativsätze", 1, 25,
		char(finch, "Relative clauses can carry a preposition ('die Frau, mit der …'), use 'welcher' as a variant, and 'was' after alles/etwas/nichts."),
		mc("Preposition + pronoun (fem. dat.)", "Die Frau, mit ___ ich arbeite, …", "der", "der", "die", "den", "dem"),
		fill("was after 'alles'", "Alles, ___ ich weiß, … (that)", "was"),
		mc("welcher (neut.)", "Das Buch, ___ ich lese, … (which)", "welches", "welches", "welcher", "welche", "welchem"),
		tr("Translate", "the house in which I live", "das Haus, in dem ich wohne"),
		speak("Blaze", "Das ist der Kollege, mit dem ich arbeite."),
	)
	addVocab(db, l,
		vw("welcher", "which", "welches Buch", "which book", finch),
		vw("alles", "everything", "alles, was …", "everything that …", "Cora"),
		vw("etwas", "something", "etwas, das …", "something that …", "Cora"),
		vw("nichts", "nothing", "nichts, was …", "nothing that …", finch),
	)

	// 52 — Conjunctions (indem, je…desto, sodass)
	s = addSkillL(db, de, u4, "Konjunktionen (B2)", "indem, je…desto, sodass.", "MessageCircle", "#00C2A8", 52, 1090)
	l = addLesson(db, s, "Satzverbindungen", 1, 25,
		char(finch, "indem = by (means of); je … desto = the … the; sodass = so that (result)."),
		mc("by (means)", "Man lernt, ___ man übt.", "indem", "indem", "weil", "obwohl", "damit"),
		mc("the … the", "Je mehr ich lerne, ___ besser werde ich.", "desto", "desto", "als", "wie", "so"),
		mc("so that (result)", "Er sprach laut, ___ alle ihn hörten.", "sodass", "sodass", "obwohl", "indem", "weil"),
		tr("Translate", "The more, the better", "Je mehr, desto besser"),
		speak("Blaze", "Man lernt eine Sprache, indem man sie spricht."),
	)
	addVocab(db, l,
		vw("indem", "by (…-ing)", "indem man übt", "by practising", finch),
		vw("je … desto", "the … the", "je mehr, desto besser", "the more, the better", finch),
		vw("sodass", "so that", "…, sodass alle es sahen", "…, so that all saw it", "Cora"),
		vw("üben", "to practise", "Ich übe täglich.", "I practise daily.", "Cora"),
	)

	// 53 — Infinitive clauses (ohne/anstatt … zu)
	s = addSkillL(db, de, u4, "ohne / anstatt … zu", "Infinitive clauses.", "PenLine", "#17A3DD", 53, 1120)
	l = addLesson(db, s, "Infinitivsätze", 1, 25,
		char(finch, "Infinitive clauses: 'ohne … zu' (without doing), 'anstatt … zu' (instead of), '(um) … zu' (in order to)."),
		mc("without", "Er ging, ___ zu grüßen.", "ohne", "ohne", "anstatt", "um", "damit"),
		mc("instead of", "___ zu arbeiten, schlief er.", "Anstatt", "Anstatt", "Ohne", "Um", "Damit"),
		fill("in order to", "Ich lerne, ___ die Prüfung zu bestehen.", "um"),
		tr("Translate", "without saying a word", "ohne ein Wort zu sagen"),
		speak("Blaze", "Er ging, ohne sich zu verabschieden."),
	)
	addVocab(db, l,
		vw("ohne … zu", "without (…-ing)", "ohne zu fragen", "without asking", finch),
		vw("anstatt … zu", "instead of", "anstatt zu warten", "instead of waiting", finch),
		vw("grüßen", "to greet", "Er grüßt nicht.", "He doesn't greet.", "Cora"),
		vw("bestehen", "to pass", "die Prüfung bestehen", "to pass the exam", "Cora"),
	)

	// 54 — Nominalisation & compounds
	s = addSkillL(db, de, u4, "Nominalisierung", "Noun style & compounds.", "Hash", "#F5A623", 54, 1150)
	l = addLesson(db, s, "Nominalstil", 1, 25,
		char(finch, "Formal German nominalises verbs ('das Lernen') and builds compound nouns: Umwelt + Schutz → der Umweltschutz. The LAST noun sets the gender."),
		mc("Nominalised verb", "'lernen' as a noun → ___ Lernen", "das", "das", "der", "die", "den"),
		mc("Compound gender", "die Umwelt + der Schutz → ___ Umweltschutz", "der", "der", "die", "das", "den"),
		fill("Nominalise 'lesen'", "das ___ (the reading)", "Lesen"),
		tr("Translate", "environmental protection", "der Umweltschutz"),
		speak("Blaze", "Das Lernen einer Sprache braucht Zeit."),
	)
	addVocab(db, l,
		vw("der Umweltschutz", "environmental protection", "Umweltschutz ist wichtig.", "Environmental protection matters.", finch),
		vw("die Digitalisierung", "digitalisation", "die Digitalisierung", "digitalisation", "Cora"),
		vw("das Berufsleben", "working life", "im Berufsleben", "in working life", "Cora"),
		vw("die Entwicklung", "development", "die Entwicklung", "the development", finch),
	)

	// 55 — Argument connectors
	s = addSkillL(db, de, u4, "Konnektoren", "zudem, hingegen, folglich …", "MessageCircle", "#FF5C5C", 55, 1180)
	l = addLesson(db, s, "Argumentieren", 1, 25,
		char(finch, "Argument connectors: zudem (moreover), hingegen (whereas), folglich (consequently), dennoch (nevertheless), jedoch (however)."),
		mc("moreover", "Es ist teuer. ___ ist es unpraktisch.", "Zudem", "Zudem", "Hingegen", "Dennoch", "Folglich"),
		mc("however", "Ich wollte kommen, ___ war ich krank.", "jedoch", "jedoch", "zudem", "folglich", "indem"),
		mc("consequently", "Er lernte viel; ___ bestand er.", "folglich", "folglich", "hingegen", "zudem", "obwohl"),
		tr("Translate", "Nevertheless, I came", "Dennoch bin ich gekommen"),
		speak("Blaze", "Das Auto ist teuer; zudem ist es schlecht für die Umwelt."),
	)
	addVocab(db, l,
		vw("zudem", "moreover", "Zudem ist es teuer.", "Moreover it's expensive.", finch),
		vw("hingegen", "whereas", "Ich hingegen denke …", "I, on the other hand, think …", "Cora"),
		vw("folglich", "consequently", "folglich …", "consequently …", "Cora"),
		vw("dennoch", "nevertheless", "Dennoch …", "Nevertheless …", finch),
	)

	// 56 — Writing (B2 argumentative)
	s = addSkillL(db, de, u4, "Schreiben: Erörterung", "Argumentative texts & formal letters.", "PenLine", "#00C2A8", 56, 1210)
	l = addLesson(db, s, "Argumentieren & Schreiben", 1, 25,
		char(finch, "B2 writing argues a position: introduction, arguments with examples, a counter-argument, and a conclusion. Mind the register (formal vs informal)."),
		write("Write an argumentative text (≈180 words): 'Sind soziale Medien gut oder schlecht für die Gesellschaft?' Give arguments for and against, then your conclusion.",
			"Soziale Medien sind heute aus unserem Leben kaum wegzudenken. Es stellt sich jedoch die Frage, ob sie der Gesellschaft mehr nutzen oder schaden.\nEinerseits ermöglichen sie schnelle Kommunikation und den Zugang zu Informationen. Zudem können Menschen über Grenzen hinweg in Kontakt bleiben.\nAndererseits führen sie oft zu Stress, Vergleich und der Verbreitung von Falschnachrichten. Folglich leidet manchmal das echte soziale Leben.\nMeiner Ansicht nach überwiegen die Vorteile, wenn man soziale Medien bewusst und in Maßen nutzt."),
		write("Write a formal email applying for a place in a German course: state your level, your goal, and ask about dates and fees.",
			"Sehr geehrte Damen und Herren,\nich interessiere mich für Ihren Deutschkurs auf dem Niveau B2. Mein Ziel ist es, bald an einer deutschen Universität zu studieren.\nKönnten Sie mir bitte mitteilen, wann der nächste Kurs beginnt und wie hoch die Gebühren sind? Außerdem würde ich gern wissen, ob es einen Einstufungstest gibt.\nVielen Dank im Voraus.\nMit freundlichen Grüßen,\nLena Bauer"),
		write("Write a short report (≈120 words) describing the results of a survey about remote work in your company (pros, cons, conclusion).",
			"In unserer Firma wurde eine Umfrage zum Thema Homeoffice durchgeführt. Insgesamt nahmen 80 Mitarbeiter teil.\nDie meisten Befragten sehen Vorteile: Sie sparen Zeit und können flexibler arbeiten. Ein Teil der Mitarbeiter vermisst jedoch den Kontakt zu den Kollegen.\nZusammenfassend lässt sich sagen, dass ein Modell aus Büro und Homeoffice am beliebtesten ist."),
	)
	addVocab(db, l,
		vw("die Gesellschaft", "society", "in der Gesellschaft", "in society", finch),
		vw("der Standpunkt", "standpoint", "mein Standpunkt", "my standpoint", "Cora"),
		vw("begründen", "to justify", "eine Meinung begründen", "to justify an opinion", "Cora"),
		vw("das Fazit", "conclusion", "Als Fazit …", "In conclusion …", finch),
	)

	// 57 — Speaking (B2 discussion/presentation)
	s = addSkillL(db, de, u4, "Sprechen: Präsentation", "Present, argue, react.", "MessageCircle", "#6C3FC5", 57, 1240)
	addLesson(db, s, "Präsentieren & Diskutieren", 1, 25,
		char(finch, "B2 speaking: present a topic with structure, argue and counter, react spontaneously, negotiate. Use 'Ich bin der Ansicht, dass …', 'Das mag stimmen, aber …'."),
		mc("Concede, then object", "Choose the phrase", "Das mag stimmen, aber …", "Das mag stimmen, aber …", "Ich heiße …", "Ich komme …", "Auf Wiedersehen"),
		speak("Blaze", "Ich bin der Ansicht, dass die Digitalisierung mehr Vorteile als Nachteile hat."),
		speak("Blaze", "Das mag stimmen, aber man darf die Risiken nicht vergessen."),
		speak("Blaze", "Zusammenfassend lässt sich sagen, dass beide Seiten gute Argumente haben."),
	)
	addVocab(db, l,
		vw("die Ansicht", "view", "meiner Ansicht nach", "in my view", finch),
		vw("das Argument", "argument", "ein gutes Argument", "a good argument", "Cora"),
		vw("überzeugen", "to convince", "Das überzeugt mich.", "That convinces me.", "Cora"),
		vw("zusammenfassend", "in summary", "Zusammenfassend …", "In summary …", finch),
	)

	addListeningL(db, de, u4, "Nachrichten: Klima",
		"A short news report. Listen, then answer.", 4, 28,
		[]models.ListeningMatch{
			lm("die Nachrichten", "the news"),
			lm("der Klimawandel", "climate change"),
			lm("Maßnahmen", "measures"),
			lm("laut Experten", "according to experts"),
		},
		[]models.ListeningLine{
			ln("Mira", "Guten Abend. In den Nachrichten: Die Regierung hat neue Klimaziele vorgestellt.", "Good evening. In the news: the government has presented new climate targets."),
			ln("Mira", "Laut Experten müssen die Emissionen bis 2030 deutlich gesenkt werden.", "According to experts, emissions must be significantly reduced by 2030."),
			ln("Mira", "Dafür sollen mehr erneuerbare Energien genutzt werden.", "For this, more renewable energy is to be used."),
			ln("Mira", "Kritiker meinen jedoch, die Maßnahmen seien nicht ausreichend.", "Critics, however, say the measures are not sufficient."),
			ln("Mira", "Die Diskussion wird morgen im Parlament fortgesetzt.", "The discussion will continue in parliament tomorrow."),
		},
		[]models.ListeningQuestion{
			lq("What did the government present?", "New climate targets", "New climate targets", "A new law on taxes", "A budget", "A holiday"),
			lq("By when must emissions be reduced?", "By 2030", "By 2030", "By 2025", "By 2040", "By 2050"),
			lq("What do critics say?", "The measures aren't enough", "The measures aren't enough", "It's perfect", "It's too expensive", "It's illegal"),
		},
	)
	addListeningL(db, de, u4, "Diskussion: Digitalisierung",
		"Riko and Zephyr debate technology. Listen, then answer.", 5, 28,
		[]models.ListeningMatch{
			lm("die Digitalisierung", "digitalisation"),
			lm("Vorteile und Nachteile", "advantages and disadvantages"),
			lm("Ich bin der Ansicht", "I am of the view"),
			lm("Das sehe ich anders", "I see that differently"),
		},
		[]models.ListeningLine{
			ln("Riko", "Ich bin der Ansicht, dass die Digitalisierung unser Leben einfacher macht.", "I'm of the view that digitalisation makes our lives easier."),
			ln("Zephyr", "Das sehe ich anders. Viele Menschen verlieren dadurch ihren Arbeitsplatz.", "I see that differently. Many people lose their jobs because of it."),
			ln("Riko", "Das mag stimmen, aber es entstehen auch neue Berufe.", "That may be true, but new jobs are also created."),
			ln("Zephyr", "Zudem werden wir immer abhängiger von der Technik.", "Moreover, we're becoming ever more dependent on technology."),
			ln("Riko", "Einerseits ja, andererseits hilft uns die Technik bei vielen Problemen.", "On one hand yes, on the other technology helps us with many problems."),
			ln("Zephyr", "Zusammenfassend brauchen wir wohl klare Regeln.", "In summary, we probably need clear rules."),
		},
		[]models.ListeningQuestion{
			lq("What is Riko's position?", "Digitalisation makes life easier", "Digitalisation makes life easier", "It's all bad", "It's boring", "It's illegal"),
			lq("What worries Zephyr most?", "Job losses & dependence", "Job losses & dependence", "Cost", "Speed", "Noise"),
			lq("What do they conclude?", "We need clear rules", "We need clear rules", "Ban technology", "Do nothing", "Use more phones"),
		},
	)
	addReadingL(db, de, u4, "Kommentar: Soziale Medien",
		"Read the commentary, then answer.", 4, 28,
		[]models.ReadingLine{
			rl("Soziale Medien haben die Art, wie wir kommunizieren, grundlegend verändert.", "Social media have fundamentally changed the way we communicate."),
			rl("Einerseits ermöglichen sie es, jederzeit mit Menschen auf der ganzen Welt in Kontakt zu bleiben.", "On one hand, they make it possible to stay in touch with people all over the world at any time."),
			rl("Andererseits verbreiten sich Falschnachrichten schneller denn je, sodass die Verunsicherung wächst.", "On the other hand, fake news spreads faster than ever, so that uncertainty grows."),
			rl("Folglich ist ein kritischer und bewusster Umgang mit diesen Plattformen unerlässlich.", "Consequently, a critical and conscious use of these platforms is essential."),
		},
		[]models.ReadingQuestion{
			rq("What have social media changed?", "How we communicate", "How we communicate", "The weather", "Prices", "Laws"),
			rq("What is the downside mentioned?", "Fake news spreads fast", "Fake news spreads fast", "They're expensive", "They're slow", "Too few users"),
			rq("What does the author call essential?", "Critical, conscious use", "Critical, conscious use", "Banning them", "Using them more", "Ignoring them"),
		},
	)
	addReadingL(db, de, u4, "Lerntipps für B2",
		"A smart study plan for B2. Read, then answer.", 5, 28,
		[]models.ReadingLine{
			rl("Auf dem B2-Niveau solltest du vor allem mit authentischen Materialien arbeiten.", "At B2 level you should work above all with authentic materials."),
			rl("Lies Zeitungen, höre Nachrichten und schau Filme mit Untertiteln.", "Read newspapers, listen to the news and watch films with subtitles."),
			rl("Übe integriert: Höre einen Text und schreibe oder sprich danach darüber.", "Practise in an integrated way: listen to a text and then write or speak about it."),
			rl("Lerne Synonyme und feste Wendungen, damit dein Ausdruck präziser wird.", "Learn synonyms and fixed expressions so that your expression becomes more precise."),
			rl("Mach regelmäßig vollständige Modellprüfungen und lass deine Texte korrigieren.", "Do full model exams regularly and have your texts corrected."),
		},
		[]models.ReadingQuestion{
			rq("What should you mainly use at B2?", "Authentic materials", "Authentic materials", "Only flashcards", "Only A1 books", "Nothing"),
			rq("What does 'integrated practice' mean here?", "Listen, then write/speak about it", "Listen, then write/speak about it", "Only grammar drills", "Only reading", "Only tests"),
			rq("Why learn synonyms & fixed phrases?", "For more precise expression", "For more precise expression", "To write less", "To speak slower", "No reason"),
		},
	)

	// ───────────────────────── C1 · Kompetenzstufe ─────────────────────
	u5 := "C1 · Kompetenzstufe"

	// 60 — Konjunktiv I (academic reported speech)
	s = addSkillL(db, de, u5, "Konjunktiv I (Akademisch)", "Reported speech in formal texts.", "Quote", "#17A3DD", 60, 1300)
	l = addLesson(db, s, "Indirekte Rede (C1)", 1, 25,
		char(finch, "In academic/journalistic writing, reported speech uses Konjunktiv I. When the K-I form equals the indicative, switch to Konjunktiv II for clarity (haben → hätten)."),
		mc("K I of haben", "Er betonte, er ___ die Studie gelesen.", "habe", "habe", "hat", "hätte", "haben"),
		mc("Overlap → K II", "Die Forscher sagten, sie ___ keine Beweise. (plural)", "hätten", "hätten", "haben", "habe", "hatten"),
		fill("K I of sein", "Sie erklärte, das Ergebnis ___ eindeutig. (be → sei)", "sei"),
		tr("Translate", "He claims he knows nothing", "Er behauptet, er wisse nichts"),
		speak("Blaze", "Der Minister sagte, er werde zurücktreten."),
	)
	addVocab(db, l,
		vw("betonen", "to emphasise", "Er betonte, dass …", "He emphasised that …", finch),
		vw("die Studie", "study", "laut der Studie", "according to the study", "Cora"),
		vw("das Ergebnis", "result", "das Ergebnis zeigt …", "the result shows …", "Cora"),
		vw("eindeutig", "clear-cut", "ein eindeutiges Ergebnis", "a clear-cut result", finch),
	)

	// 61 — Konjunktiv II (past, regret, mixed)
	s = addSkillL(db, de, u5, "Konjunktiv II (Vergangenheit)", "Regrets & mixed conditionals.", "Sparkles", "#6C3FC5", 61, 1340)
	l = addLesson(db, s, "Irreale Vergangenheit", 1, 25,
		char(finch, "Past Konjunktiv II = hätte/wäre + Partizip II for regrets and unreal past: 'Ich hätte mehr lernen sollen.' Mixed conditionals link past cause and present result."),
		mc("Regret with modal", "Ich ___ mehr lernen sollen. (should have)", "hätte", "hätte", "würde", "wäre", "habe"),
		fill("Unreal past", "Wenn ich das gewusst hätte, ___ ich anders gehandelt.", "hätte"),
		mc("Mixed conditional", "Wenn ich damals studiert hätte, ___ ich jetzt mehr Chancen.", "hätte", "hätte", "habe", "hatte", "würde"),
		tr("Translate", "I wish I had said it", "Ich wünschte, ich hätte es gesagt"),
		speak("Blaze", "Ich hätte das anders machen sollen."),
	)
	addVocab(db, l,
		vw("das Bedauern", "regret", "voller Bedauern", "full of regret", finch),
		vw("handeln", "to act", "anders handeln", "to act differently", "Cora"),
		vw("damals", "back then", "damals war alles anders", "back then everything was different", "Cora"),
		vw("die Chance", "chance / opportunity", "eine Chance haben", "to have a chance", finch),
	)

	// 62 — Extended participial attributes
	s = addSkillL(db, de, u5, "Erweiterte Partizipien", "Packed participial attributes.", "Layers", "#F5A623", 62, 1380)
	l = addLesson(db, s, "Partizipialattribute", 1, 25,
		char(finch, "C1 packs information into extended participial attributes: 'die gestern in der Zeitung gelesene Nachricht' = 'the news read in the paper yesterday'."),
		mc("Partizip II (passive, done)", "die gestern ___ Nachricht (read)", "gelesene", "gelesene", "lesende", "gelesen", "lesen"),
		mc("Partizip I (active, ongoing)", "das ruhig ___ Kind (sleeping)", "schlafende", "schlafende", "geschlafene", "schläft", "schlafen"),
		fill("Partizip II", "der von allen ___ Plan (loved → geliebte)", "geliebte"),
		tr("Translate", "the problem solved by the team", "das vom Team gelöste Problem"),
		speak("Blaze", "die in der Zeitung veröffentlichte Nachricht"),
	)
	addVocab(db, l,
		vw("das Attribut", "attribute", "ein erweitertes Attribut", "an extended attribute", finch),
		vw("veröffentlicht", "published", "ein veröffentlichter Artikel", "a published article", "Cora"),
		vw("gelöst", "solved", "das gelöste Problem", "the solved problem", "Cora"),
		vw("geliebt", "beloved", "der geliebte Ort", "the beloved place", finch),
	)

	// 63 — Nominal style & function-verb phrases
	s = addSkillL(db, de, u5, "Nominalstil & FVG", "Academic noun style.", "Hash", "#00C2A8", 63, 1420)
	l = addLesson(db, s, "Nominalstil", 1, 25,
		char(finch, "Formal German prefers Nominalstil ('nach der Prüfung des Antrags'). Funktionsverbgefüge are fixed verb+noun phrases: 'in Frage stellen', 'zur Verfügung stehen'."),
		mc("FVG: call into question", "Ich ___ das in Frage. (question)", "stelle", "stelle", "setze", "nehme", "mache"),
		fill("FVG: be available", "Die Daten stehen zur ___ . (Verfügung)", "Verfügung"),
		mc("Nominalise 'prüfen'", "→ die ___ (examination)", "Prüfung", "Prüfung", "Prüfen", "Prüfer", "geprüft"),
		tr("Translate", "to make a decision", "eine Entscheidung treffen"),
		speak("Blaze", "Wir stellen den Plan in Frage."),
	)
	addVocab(db, l,
		vw("in Frage stellen", "to call into question", "etwas in Frage stellen", "to question something", finch),
		vw("zur Verfügung stehen", "to be available", "steht zur Verfügung", "is available", "Cora"),
		vw("eine Entscheidung treffen", "to make a decision", "eine Entscheidung treffen", "to make a decision", "Cora"),
		vw("der Antrag", "application", "einen Antrag stellen", "to file an application", finch),
	)

	// 64 — Passive alternatives
	s = addSkillL(db, de, u5, "Passiversatzformen", "man, sich lassen, sein+zu, -bar.", "Layers", "#FF5C5C", 64, 1460)
	l = addLesson(db, s, "Passiversatz", 1, 25,
		char(finch, "Passive alternatives sound more elegant: 'sich lassen' ('lässt sich lösen'), 'sein + zu + Infinitiv' ('ist zu lösen'), and -bar adjectives ('lösbar')."),
		mc("sich lassen", "Das Problem ___ sich leicht lösen.", "lässt", "lässt", "wird", "ist", "kann"),
		mc("sein + zu", "Das Formular ist bis Montag ___ . (to be submitted)", "einzureichen", "einzureichen", "eingereicht", "einreichen", "reicht ein"),
		fill("-bar adjective", "Das Problem ist nicht ___ . (solvable → lösbar)", "lösbar"),
		tr("Translate", "That can be done", "Das lässt sich machen"),
		speak("Blaze", "Das Problem lässt sich lösen."),
	)
	addVocab(db, l,
		vw("lösbar", "solvable", "ein lösbares Problem", "a solvable problem", finch),
		vw("einreichen", "to submit", "einen Antrag einreichen", "to submit an application", "Cora"),
		vw("vermeidbar", "avoidable", "ein vermeidbarer Fehler", "an avoidable mistake", "Cora"),
		vw("machbar", "feasible", "Das ist machbar.", "That's feasible.", finch),
	)

	// 65 — Subjective modal verbs
	s = addSkillL(db, de, u5, "Subjektive Modalverben", "Speculation & deduction.", "Sparkles", "#6C3FC5", 65, 1500)
	l = addLesson(db, s, "Vermutungen ausdrücken", 1, 25,
		char(finch, "Modals express probability/deduction: 'Er dürfte zu Hause sein' (probably), 'Sie müsste angekommen sein' (must, deduction), 'Das könnte stimmen' (might)."),
		mc("Probability", "Er ___ jetzt zu Hause sein.", "dürfte", "dürfte", "darf", "muss", "kann"),
		mc("Deduction", "Sie ___ schon angekommen sein.", "müsste", "müsste", "muss", "dürfte", "könnte"),
		mc("Possibility", "Das ___ stimmen.", "könnte", "könnte", "kann", "muss", "darf"),
		tr("Translate", "He must have forgotten it", "Er muss es vergessen haben"),
		speak("Blaze", "Sie dürfte das längst gewusst haben."),
	)
	addVocab(db, l,
		vw("vermutlich", "presumably", "vermutlich morgen", "presumably tomorrow", finch),
		vw("offenbar", "apparently", "Er ist offenbar krank.", "He's apparently ill.", "Cora"),
		vw("anscheinend", "seemingly", "anscheinend stimmt es", "seemingly it's true", "Cora"),
		vw("die Vermutung", "assumption", "eine Vermutung äußern", "to voice an assumption", finch),
	)

	// 66 — Prepositional adverbs
	s = addSkillL(db, de, u5, "Präpositionaladverbien", "darauf, worauf, womit …", "Link2", "#17A3DD", 66, 1540)
	l = addLesson(db, s, "da(r)- & wo(r)-", 1, 25,
		char(finch, "When a verb's preposition refers to a thing, use da(r)+preposition ('darauf', 'damit'); in questions wo(r)+preposition ('worauf', 'womit')."),
		mc("da-form", "Ich warte ___ . (for it)", "darauf", "darauf", "auf es", "darüber", "davon"),
		mc("wo-form", "___ wartest du? (what for)", "Worauf", "Worauf", "Wofür", "Womit", "Wovon"),
		fill("da-form", "Ich freue mich ___ . (about it → darüber)", "darüber"),
		tr("Translate", "What are you talking about?", "Worüber sprichst du?"),
		speak("Blaze", "Ich denke oft daran."),
	)
	addVocab(db, l,
		vw("darauf", "on/for it", "Ich warte darauf.", "I'm waiting for it.", finch),
		vw("womit", "with what", "Womit kann ich helfen?", "With what can I help?", "Cora"),
		vw("worüber", "about what", "Worüber redest du?", "What are you talking about?", "Cora"),
		vw("davon", "of/from it", "Ich habe davon gehört.", "I've heard of it.", finch),
	)

	// 67 — Genitive relative pronouns
	s = addSkillL(db, de, u5, "Relativsätze: dessen/deren", "Genitive relatives.", "Link2", "#F5A623", 67, 1580)
	l = addLesson(db, s, "dessen & deren", 1, 25,
		char(finch, "Genitive relative pronouns: 'dessen' (masc/neut), 'deren' (fem/plural): 'der Mann, dessen Auto kaputt ist'."),
		mc("dessen (masc.)", "Der Mann, ___ Auto kaputt ist, …", "dessen", "dessen", "deren", "der", "den"),
		mc("deren (fem.)", "Die Frau, ___ Kinder hier spielen, …", "deren", "deren", "dessen", "der", "die"),
		fill("deren (plural)", "die Leute, ___ Namen ich vergaß", "deren"),
		tr("Translate", "the city whose history is famous", "die Stadt, deren Geschichte berühmt ist"),
		speak("Blaze", "Das ist der Autor, dessen Buch ich liebe."),
	)
	addVocab(db, l,
		vw("dessen", "whose (m/n)", "der Mann, dessen …", "the man whose …", finch),
		vw("deren", "whose (f/pl)", "die Frau, deren …", "the woman whose …", finch),
		vw("die Geschichte", "history / story", "eine lange Geschichte", "a long history", "Cora"),
		vw("berühmt", "famous", "ein berühmter Autor", "a famous author", "Cora"),
	)

	// 68 — Advanced subordinate clauses
	s = addSkillL(db, de, u5, "Anspruchsvolle Nebensätze", "obgleich, ohne dass, sodass …", "MessageCircle", "#00C2A8", 68, 1620)
	l = addLesson(db, s, "Komplexe Nebensätze", 1, 25,
		char(finch, "Advanced subordinators: obgleich (although), ohne dass (without … -ing, different subject), sodass (so that), indem (by)."),
		mc("without (diff. subject)", "Er ging, ___ jemand es bemerkte.", "ohne dass", "ohne dass", "ohne zu", "damit", "obwohl"),
		mc("although", "___ er müde war, arbeitete er weiter.", "Obgleich", "Obgleich", "Indem", "Sodass", "Weil"),
		mc("by (means)", "Er überzeugte uns, ___ er Beispiele gab.", "indem", "indem", "ohne dass", "obgleich", "sodass"),
		tr("Translate", "without my noticing", "ohne dass ich es bemerkte"),
		speak("Blaze", "Er half uns, ohne dass wir ihn darum baten."),
	)
	addVocab(db, l,
		vw("obgleich", "although", "obgleich es schwer ist", "although it's hard", finch),
		vw("ohne dass", "without …-ing", "ohne dass er es merkt", "without him noticing", finch),
		vw("bemerken", "to notice", "Ich bemerkte nichts.", "I noticed nothing.", "Cora"),
		vw("überzeugen", "to convince", "Er überzeugte mich.", "He convinced me.", "Cora"),
	)

	// 69 — Discourse markers
	s = addSkillL(db, de, u5, "Diskursmarker", "zwar…aber, allerdings, demzufolge.", "Quote", "#FF5C5C", 69, 1660)
	l = addLesson(db, s, "Diskursmarker", 1, 25,
		char(finch, "Discourse markers structure argument: 'zwar … aber' (admittedly … but), 'allerdings' (however), 'demzufolge' (consequently), 'hingegen' (whereas)."),
		mc("admittedly … but", "Das ist ___ teuer, aber sehr gut.", "zwar", "zwar", "allerdings", "demzufolge", "hingegen"),
		mc("however", "Ich komme, ___ etwas später.", "allerdings", "allerdings", "zwar", "folglich", "indem"),
		mc("consequently", "Die Nachfrage stieg; ___ stiegen die Preise.", "demzufolge", "demzufolge", "zwar", "hingegen", "obgleich"),
		tr("Translate", "admittedly difficult, but possible", "zwar schwierig, aber möglich"),
		speak("Blaze", "Das ist zwar teuer, aber es lohnt sich."),
	)
	addVocab(db, l,
		vw("zwar", "admittedly", "zwar … aber …", "admittedly … but …", finch),
		vw("allerdings", "however", "allerdings ist es spät", "however, it's late", "Cora"),
		vw("demzufolge", "consequently", "demzufolge …", "consequently …", "Cora"),
		vw("ohnehin", "anyway", "Ich gehe ohnehin.", "I'm going anyway.", finch),
	)

	// 70 — Word order for emphasis
	s = addSkillL(db, de, u5, "Wortstellung & Betonung", "Fronting & emphasis.", "PenLine", "#6C3FC5", 70, 1700)
	l = addLesson(db, s, "Betonung im Satz", 1, 25,
		char(finch, "German fronts elements for emphasis — but the verb stays in second position: 'Gestern habe ich ihn gesehen', 'Diesen Film kenne ich'."),
		mc("Front 'gestern'", "___ habe ich ihn gesehen.", "Gestern", "Gestern", "Ich", "Habe", "Gesehen"),
		mc("Object fronted, verb 2nd", "Diesen Film ___ ich gut.", "kenne", "kenne", "ich kenne", "kennen", "gekannt"),
		fill("Front 'das'", "___ verstehe ich nicht.", "Das"),
		tr("Translate", "Him I don't trust", "Ihm vertraue ich nicht"),
		speak("Blaze", "Diesen Fehler mache ich nie wieder."),
	)
	addVocab(db, l,
		vw("betonen", "to stress", "etwas betonen", "to stress something", finch),
		vw("hervorheben", "to highlight", "Ich möchte hervorheben …", "I'd like to highlight …", "Cora"),
		vw("die Wortstellung", "word order", "die freie Wortstellung", "flexible word order", "Cora"),
		vw("vertrauen", "to trust (+dat.)", "Ich vertraue dir.", "I trust you.", finch),
	)

	// 71 — Correlatives & gradation particles
	s = addSkillL(db, de, u5, "nicht nur … sondern auch", "Correlatives & gradation.", "Hash", "#17A3DD", 71, 1740)
	l = addLesson(db, s, "Korrelative", 1, 25,
		char(finch, "Correlative conjunctions and gradation particles: 'nicht nur … sondern auch', 'sowohl … als auch', 'sogar' (even), 'kaum' (hardly), 'lediglich' (merely)."),
		mc("not only … but also", "Er spricht nicht nur Deutsch, ___ auch Französisch.", "sondern", "sondern", "aber", "oder", "denn"),
		mc("both … and", "___ Anna als auch Tom kommen.", "Sowohl", "Sowohl", "Entweder", "Weder", "Nicht"),
		mc("even", "Er hat ___ am Sonntag gearbeitet.", "sogar", "sogar", "kaum", "lediglich", "nur"),
		tr("Translate", "hardly anyone came", "Es kam kaum jemand"),
		speak("Blaze", "Sie ist nicht nur klug, sondern auch fleißig."),
	)
	addVocab(db, l,
		vw("nicht nur … sondern auch", "not only … but also", "nicht nur … sondern auch", "not only … but also", finch),
		vw("sowohl … als auch", "both … and", "sowohl … als auch", "both … and", "Cora"),
		vw("sogar", "even", "sogar am Sonntag", "even on Sunday", "Cora"),
		vw("kaum", "hardly", "kaum jemand", "hardly anyone", finch),
	)

	// 72 — Modal particles (C1)
	s = addSkillL(db, de, u5, "Modalpartikeln (C1)", "eben, halt, schon, wohl.", "Quote", "#F5A623", 72, 1780)
	l = addLesson(db, s, "Feine Partikeln", 1, 25,
		char(finch, "Subtle particles colour speech: 'eben/halt' (resignation), 'schon' (reassurance), 'wohl' (presumably)."),
		mc("resignation", "Das ist ___ so. (just/simply)", "eben", "eben", "wohl", "schon", "denn"),
		mc("reassurance", "Das wird ___ klappen. (surely)", "schon", "schon", "eben", "wohl", "mal"),
		mc("presumably", "Er ist ___ krank. (presumably)", "wohl", "wohl", "eben", "schon", "denn"),
		tr("Translate", "That's just how it is", "Das ist halt so"),
		speak("Blaze", "Das wird schon wieder."),
	)
	addVocab(db, l,
		vw("eben", "just (so)", "Das ist eben so.", "That's just how it is.", finch),
		vw("halt", "just (so)", "Das ist halt so.", "That's just how it is.", finch),
		vw("schon", "(reassurance)", "Das wird schon.", "It'll be fine.", "Cora"),
		vw("wohl", "presumably", "Er ist wohl müde.", "He's presumably tired.", "Cora"),
	)

	// 73 — Writing (C1)
	s = addSkillL(db, de, u5, "Schreiben: Erörterung (C1)", "Academic essays & summaries.", "PenLine", "#00C2A8", 73, 1820)
	l = addLesson(db, s, "Akademisches Schreiben", 1, 25,
		char(finch, "C1 writing is formal/academic: a structured Erörterung and a precise Zusammenfassung. Use Nominalstil, varied connectors and an objective register."),
		write("Write an argumentative essay (≈220 words): 'Sollte ein Studium kostenlos sein?' Present arguments, a counter-position and a reasoned conclusion.",
			"Die Frage, ob ein Studium kostenlos sein sollte, wird seit Langem kontrovers diskutiert.\nBefürworter argumentieren, dass kostenlose Bildung Chancengleichheit schaffe und auch Menschen aus ärmeren Familien ein Studium ermögliche. Zudem profitiere die gesamte Gesellschaft von gut ausgebildeten Fachkräften.\nKritiker wenden hingegen ein, dass kostenlose Studiengänge hohe Kosten für den Staat verursachten und die Qualität darunter leiden könne. Außerdem sei eine moderate Gebühr durchaus zumutbar.\nMeiner Ansicht nach überwiegen die Vorteile: Bildung sollte nicht vom Einkommen abhängen. Eine sinnvolle Lösung wäre, das Studium grundsätzlich kostenlos anzubieten, jedoch gezielt in die Qualität zu investieren."),
		write("Summarise objectively in ≈80 words a text about the advantages and disadvantages of working from home (no personal opinion).",
			"Der Text befasst sich mit den Vor- und Nachteilen des Homeoffice. Als Vorteile werden vor allem Zeitersparnis und Flexibilität genannt, da der Arbeitsweg entfällt. Demgegenüber stehen Nachteile wie der fehlende Kontakt zu Kollegen und mögliche Konzentrationsprobleme zu Hause. Abschließend stellt der Autor fest, dass eine Mischung aus Büro- und Heimarbeit für viele Beschäftigte die beste Lösung darstelle."),
		write("Write a formal letter to a university requesting a deferral of your enrolment, explaining your reasons clearly.",
			"Sehr geehrte Damen und Herren,\nhiermit möchte ich Sie bitten, meine Immatrikulation um ein Semester zu verschieben. Aus gesundheitlichen Gründen ist es mir derzeit leider nicht möglich, das Studium aufzunehmen.\nIch wäre Ihnen sehr dankbar, wenn Sie meinem Antrag entsprechen könnten, und stehe für Rückfragen jederzeit zur Verfügung.\nMit freundlichen Grüßen,\nKwame Mensah"),
	)
	addVocab(db, l,
		vw("die Erörterung", "argumentation", "eine Erörterung schreiben", "to write an argumentation", finch),
		vw("die Zusammenfassung", "summary", "eine kurze Zusammenfassung", "a short summary", "Cora"),
		vw("objektiv", "objective", "objektiv schreiben", "to write objectively", "Cora"),
		vw("die Chancengleichheit", "equal opportunity", "Chancengleichheit schaffen", "to create equal opportunity", finch),
	)

	// 74 — Speaking (C1)
	s = addSkillL(db, de, u5, "Sprechen: Vortrag (C1)", "Structured talks & debate.", "MessageCircle", "#6C3FC5", 74, 1860)
	addLesson(db, s, "Vortrag & Diskussion", 1, 25,
		char(finch, "C1 speaking: deliver a structured talk, discuss abstract ideas, handle counter-questions, and signpost clearly: 'Zunächst …', 'Darüber hinaus …', 'Abschließend …'."),
		mc("Open a structured talk", "How to begin?", "Zunächst möchte ich …", "Zunächst möchte ich …", "Ich heiße …", "Tschüss", "Ich komme aus …"),
		speak("Blaze", "Zunächst möchte ich auf die Vorteile der Globalisierung eingehen."),
		speak("Blaze", "Darüber hinaus spielt die Nachhaltigkeit eine zentrale Rolle."),
		speak("Blaze", "Abschließend lässt sich festhalten, dass die Vorteile überwiegen."),
	)
	addVocab(db, l,
		vw("zunächst", "first(ly)", "Zunächst …", "First …", finch),
		vw("darüber hinaus", "moreover", "Darüber hinaus …", "Moreover …", "Cora"),
		vw("abschließend", "finally", "Abschließend …", "Finally …", "Cora"),
		vw("die Nachhaltigkeit", "sustainability", "die Nachhaltigkeit fördern", "to promote sustainability", finch),
	)

	addListeningL(db, de, u5, "Vortrag: Künstliche Intelligenz",
		"A short lecture excerpt. Listen, then answer.", 6, 32,
		[]models.ListeningMatch{
			lm("die künstliche Intelligenz", "artificial intelligence"),
			lm("Chancen und Risiken", "opportunities and risks"),
			lm("zahlreiche Bereiche", "numerous areas"),
			lm("Verantwortung", "responsibility"),
		},
		[]models.ListeningLine{
			ln(finch, "Sehr geehrte Damen und Herren, mein heutiger Vortrag behandelt die künstliche Intelligenz.", "Ladies and gentlemen, my talk today deals with artificial intelligence."),
			ln(finch, "Zunächst ist festzuhalten, dass KI bereits zahlreiche Bereiche unseres Lebens verändert hat.", "First, it should be noted that AI has already changed numerous areas of our lives."),
			ln(finch, "Einerseits bietet sie enorme Chancen, etwa in der Medizin und der Forschung.", "On one hand, it offers enormous opportunities, for instance in medicine and research."),
			ln(finch, "Andererseits dürfen die Risiken, wie der Verlust von Arbeitsplätzen, nicht ignoriert werden.", "On the other hand, the risks, such as job losses, must not be ignored."),
			ln(finch, "Entscheidend ist daher, dass wir Verantwortung übernehmen und klare Regeln schaffen.", "It is therefore crucial that we take responsibility and create clear rules."),
			ln(finch, "Abschließend lässt sich sagen: Die Technik ist weder gut noch schlecht — es kommt auf uns an.", "In conclusion: the technology is neither good nor bad — it depends on us."),
		},
		[]models.ListeningQuestion{
			lq("What is the lecture about?", "Artificial intelligence", "Artificial intelligence", "History", "Economics", "Sport"),
			lq("Which risk is named?", "Job losses", "Job losses", "Pollution", "Inflation", "Traffic"),
			lq("What does the speaker say is crucial?", "Taking responsibility & clear rules", "Taking responsibility & clear rules", "Banning AI", "Ignoring it", "More speed"),
		},
	)
	addListeningL(db, de, u5, "Eine kleine Debatte",
		"Riko and Zephyr debate online learning. Listen closely, then answer.", 5, 30,
		[]models.ListeningMatch{
			lm("meiner Meinung nach", "in my opinion"),
			lm("einerseits … andererseits", "on one hand … on the other"),
			lm("Das stimmt zwar, aber …", "That's true, but …"),
			lm("zusammenfassend", "to summarise"),
		},
		[]models.ListeningLine{
			ln("Riko", "Meiner Meinung nach ist das Lernen am Computer dem Unterricht klar überlegen.", "In my opinion, learning on the computer is clearly superior to classroom teaching."),
			ln("Zephyr", "Das stimmt zwar teilweise, aber der persönliche Kontakt fehlt völlig.", "That's partly true, but the personal contact is completely missing."),
			ln("Riko", "Einerseits hast du recht, andererseits kann man jederzeit und überall lernen.", "On one hand you're right, on the other you can learn anytime, anywhere."),
			ln("Zephyr", "Flexibilität ist kein Ersatz für echte Gespräche und Diskussionen.", "Flexibility is no substitute for real conversations and discussions."),
			ln("Riko", "Zusammenfassend würde ich sagen: Eine Mischung aus beidem ist ideal.", "To summarise, I'd say: a mix of both is ideal."),
			ln("Zephyr", "Dem kann ich ausnahmsweise zustimmen.", "I can, for once, agree with that."),
		},
		[]models.ListeningQuestion{
			lq("What is Riko's main point?", "Online learning is superior", "Online learning is superior", "Books are best", "Teachers are useless", "Learning is boring"),
			lq("What does Zephyr say is missing online?", "Personal contact", "Personal contact", "Money", "Speed", "Homework"),
			lq("What do they finally agree on?", "A mix of both", "A mix of both", "Only online", "Only classroom", "Nothing"),
		},
	)
	addReadingL(db, de, u5, "Kommentar: Die Kunst des Lesens",
		"Read this advanced commentary, then answer.", 5, 30,
		[]models.ReadingLine{
			rl("In einer Zeit, in der Informationen im Sekundentakt auf uns einströmen, gerät das vertiefte Lesen zunehmend in Vergessenheit.", "In an age in which information streams at us every second, deep reading is increasingly being forgotten."),
			rl("Wer einen langen Text liest, übt nicht nur Geduld, sondern auch die Fähigkeit, komplexe Gedanken zu durchdringen.", "Whoever reads a long text practises not only patience but also the ability to penetrate complex ideas."),
			rl("Trotz aller technischen Bequemlichkeit sollte das Buch deshalb nicht vorschnell abgeschrieben werden.", "Despite all technical convenience, the book should therefore not be written off prematurely."),
			rl("Denn das Lesen formt, so paradox es klingt, gerade jene Aufmerksamkeit, die unsere schnelle Welt am dringendsten benötigt.", "For reading shapes, paradoxical as it sounds, precisely the attention our fast world most urgently needs."),
		},
		[]models.ReadingQuestion{
			rq("What is being forgotten, per the text?", "Deep reading", "Deep reading", "Writing", "Speaking", "Counting"),
			rq("What does long reading train?", "Patience & complex thinking", "Patience & complex thinking", "Typing speed", "Memory only", "Spelling"),
			rq("The author's stance on books is:", "They shouldn't be written off", "They shouldn't be written off", "They are obsolete", "They are overrated", "They are too long"),
		},
	)
	addReadingL(db, de, u5, "Lerntipps für C1",
		"A smart study plan for C1. Read, then answer.", 6, 32,
		[]models.ReadingLine{
			rl("Auf dem C1-Niveau geht es weniger um Regeln als um Stil, Register und Präzision.", "At C1 level it's less about rules and more about style, register and precision."),
			rl("Lies anspruchsvolle Texte aus Qualitätszeitungen und achte bewusst auf Konnektoren und Nominalstil.", "Read demanding texts from quality newspapers and pay conscious attention to connectors and nominal style."),
			rl("Schreibe regelmäßig Erörterungen und lass sie korrigieren, um Fehler systematisch zu beseitigen.", "Write argumentative essays regularly and have them corrected, in order to eliminate mistakes systematically."),
			rl("Sprich über abstrakte Themen und nimm dich dabei auf, um Flüssigkeit und Aussprache zu prüfen.", "Speak about abstract topics and record yourself, in order to check fluency and pronunciation."),
			rl("Wer auf diesem Niveau Fortschritte macht, braucht vor allem hochwertigen Input und ehrliches Feedback.", "Whoever progresses at this level needs above all high-quality input and honest feedback."),
		},
		[]models.ReadingQuestion{
			rq("What matters most at C1?", "Style, register, precision", "Style, register, precision", "Only grammar rules", "Only vocabulary lists", "Nothing"),
			rq("What should you pay attention to when reading?", "Connectors & nominal style", "Connectors & nominal style", "Only spelling", "Page numbers", "Pictures"),
			rq("Why record yourself speaking?", "To check fluency & pronunciation", "To check fluency & pronunciation", "For fun", "To save time", "No reason"),
		},
	)

	// ───────────────────────── C2 · Feinheiten & Idiomatik ─────────────
	u6 := "C2 · Feinheiten"

	// 76 — Konjunktiv I mastery
	s = addSkillL(db, de, u6, "Konjunktiv I (Meisterung)", "Reported speech, full control.", "Quote", "#06AECE", 76, 2000)
	l = addLesson(db, s, "Indirekte Rede (C2)", 1, 30,
		char(finch, "At C2 you wield reported speech with full control — Konjunktiv I throughout; where it coincides with the indicative, switch to Konjunktiv II ('sie kommen' → 'sie kämen')."),
		mc("K I, 3rd sg.", "Er erklärte, er ___ diese These.", "vertrete", "vertrete", "vertritt", "verträte", "vertreten"),
		mc("Overlap → K II", "Die Autoren behaupten, sie ___ recht. (plural)", "hätten", "hätten", "haben", "habe", "hatten"),
		fill("K I of werden", "Sie sagte, sie ___ kommen. (werde)", "werde"),
		tr("Translate", "He claimed he had not known it", "Er behauptete, er habe es nicht gewusst"),
		speak("Blaze", "Die Zeugin sagte, sie habe nichts gesehen."),
	)
	addVocab(db, l,
		vw("vertreten", "to hold (a view)", "eine These vertreten", "to hold a thesis", finch),
		vw("die These", "thesis", "eine gewagte These", "a bold thesis", "Cora"),
		vw("einräumen", "to concede", "Er räumte ein, dass …", "He conceded that …", "Cora"),
		vw("der Zeuge", "witness", "ein wichtiger Zeuge", "an important witness", finch),
	)

	// 77 — Konjunktiv II nuance
	s = addSkillL(db, de, u6, "Konjunktiv II (Nuancen)", "Irony, politeness, caution.", "Sparkles", "#6C3FC5", 77, 2040)
	l = addLesson(db, s, "Feine Nuancen", 1, 30,
		char(finch, "Konjunktiv II conveys nuance: irony ('Das wäre ja noch schöner!'), extreme politeness ('Ich hätte da eine Bitte'), and cautious claims ('Damit wäre alles gesagt')."),
		mc("Extreme politeness", "___ Sie wohl so freundlich? (would you be)", "Wären", "Wären", "Sind", "Waren", "Seid"),
		fill("Cautious claim", "Damit ___ wohl alles gesagt. (would be → wäre)", "wäre"),
		mc("Irony", "Das ___ ja noch schöner! (would be)", "wäre", "wäre", "ist", "war", "sei"),
		tr("Translate", "I'd have one request", "Ich hätte da eine Bitte"),
		speak("Blaze", "Ich hätte da noch eine kleine Anmerkung."),
	)
	addVocab(db, l,
		vw("die Bitte", "request", "eine höfliche Bitte", "a polite request", finch),
		vw("die Anmerkung", "remark", "eine kritische Anmerkung", "a critical remark", "Cora"),
		vw("die Ironie", "irony", "feine Ironie", "subtle irony", "Cora"),
		vw("andeuten", "to hint", "Er deutete an, dass …", "He hinted that …", finch),
	)

	// 78 — Genitive (formal/literary)
	s = addSkillL(db, de, u6, "Genitiv (gehoben)", "Formal & literary genitive.", "Hash", "#F5A623", 78, 2080)
	l = addLesson(db, s, "Gehobener Genitiv", 1, 30,
		char(finch, "Formal/literary German uses the genitive richly: 'angesichts des Problems', and rarer genitive prepositions like mittels, zwecks, hinsichtlich."),
		mc("in view of", "___ des Problems handelten wir schnell.", "Angesichts", "Angesichts", "Wegen", "Trotz", "Mit"),
		mc("by means of", "Er öffnete es ___ eines Schlüssels.", "mittels", "mittels", "mit", "durch", "von"),
		fill("regarding", "___ der Kosten gibt es Bedenken. (hinsichtlich)", "Hinsichtlich"),
		tr("Translate", "for the sake of clarity", "um der Klarheit willen"),
		speak("Blaze", "Angesichts der Lage müssen wir handeln."),
	)
	addVocab(db, l,
		vw("angesichts", "in view of (+gen)", "angesichts der Lage", "in view of the situation", finch),
		vw("mittels", "by means of (+gen)", "mittels eines Tricks", "by means of a trick", "Cora"),
		vw("hinsichtlich", "regarding (+gen)", "hinsichtlich der Frage", "regarding the question", "Cora"),
		vw("zwecks", "for the purpose of", "zwecks Klärung", "for clarification", finch),
	)

	// 79 — Extended attributes (advanced)
	s = addSkillL(db, de, u6, "Erweiterte Attribute", "Heavy participial phrases.", "Layers", "#00C2A8", 79, 2120)
	l = addLesson(db, s, "Komplexe Attribute", 1, 30,
		char(finch, "C2 readers unpack heavy extended attributes: 'die von der Regierung im letzten Jahr beschlossenen Maßnahmen' = 'the measures decided by the government last year'."),
		mc("Partizip II", "die im Parlament ___ Maßnahmen (decided)", "beschlossenen", "beschlossenen", "beschließenden", "beschlossen", "beschließen"),
		mc("Partizip I", "die stetig ___ Zahl (rising)", "steigende", "steigende", "gestiegene", "steigt", "steigen"),
		fill("Partizip II (masc. nom.)", "ein sorgfältig ___ Plan (worked out → ausgearbeiteter)", "ausgearbeiteter"),
		tr("Translate", "the long-awaited decision", "die lang erwartete Entscheidung"),
		speak("Blaze", "die von allen Experten empfohlene Lösung"),
	)
	addVocab(db, l,
		vw("beschlossen", "decided", "die beschlossene Maßnahme", "the decided measure", finch),
		vw("die Maßnahme", "measure", "Maßnahmen ergreifen", "to take measures", "Cora"),
		vw("ausgearbeitet", "worked out", "ein ausgearbeiteter Plan", "a detailed plan", "Cora"),
		vw("erwartet", "awaited / expected", "lang erwartet", "long awaited", finch),
	)

	// 80 — Nominal style & noun-verb combinations
	s = addSkillL(db, de, u6, "Nomen-Verb-Verbindungen", "Funktionsverbgefüge.", "Hash", "#FF5C5C", 80, 2160)
	l = addLesson(db, s, "Nominalstil (C2)", 1, 30,
		char(finch, "Academic German favours noun-verb combinations: 'Kritik üben', 'in Anspruch nehmen', 'zum Ausdruck bringen', 'in Betracht ziehen'."),
		mc("to criticise", "Kritik ___ (üben)", "üben", "üben", "machen", "nehmen", "geben"),
		mc("to make use of", "etwas in Anspruch ___ (nehmen)", "nehmen", "nehmen", "üben", "bringen", "stellen"),
		fill("to express", "zum Ausdruck ___ (bringen)", "bringen"),
		tr("Translate", "to take into account", "in Betracht ziehen"),
		speak("Blaze", "Wir müssen alle Faktoren in Betracht ziehen."),
	)
	addVocab(db, l,
		vw("Kritik üben", "to criticise", "Kritik üben an …", "to criticise …", finch),
		vw("in Anspruch nehmen", "to make use of", "Hilfe in Anspruch nehmen", "to use help", "Cora"),
		vw("zum Ausdruck bringen", "to express", "Dank zum Ausdruck bringen", "to express thanks", "Cora"),
		vw("in Betracht ziehen", "to consider", "Optionen in Betracht ziehen", "to consider options", finch),
	)

	// 81 — Voice as stylistic choice
	s = addSkillL(db, de, u6, "Aktiv & Passiv (Stil)", "Voice for stylistic effect.", "Layers", "#17A3DD", 81, 2200)
	l = addLesson(db, s, "Stilwahl", 1, 30,
		char(finch, "Choosing voice is a stylistic decision: the passive (or 'man') depersonalises ('Es wurde beschlossen …'); the active is more direct. Passiversatz keeps prose elegant."),
		mc("man-construction", "Wie sagt ___ das auf Deutsch?", "man", "man", "es", "sich", "wird"),
		mc("sein + zu", "Die Frist ist unbedingt ___ . (to be met)", "einzuhalten", "einzuhalten", "eingehalten", "einhalten", "hält ein"),
		fill("Passiversatz (sein+zu)", "Fehler sind zu ___ . (avoid → vermeiden)", "vermeiden"),
		tr("Translate", "The goal is hard to achieve", "Das Ziel ist schwer zu erreichen"),
		speak("Blaze", "Diese Frist ist strikt einzuhalten."),
	)
	addVocab(db, l,
		vw("die Frist", "deadline", "die Frist einhalten", "to meet the deadline", finch),
		vw("einhalten", "to observe / meet", "Regeln einhalten", "to observe rules", "Cora"),
		vw("vermeiden", "to avoid", "Fehler vermeiden", "to avoid mistakes", "Cora"),
		vw("erreichen", "to achieve", "ein Ziel erreichen", "to achieve a goal", finch),
	)

	// 82 — Modal verbs: fine shades
	s = addSkillL(db, de, u6, "Modalverben (Nuance)", "Hearsay, claims, hedging.", "Sparkles", "#6C3FC5", 82, 2240)
	l = addLesson(db, s, "Feine Bedeutung", 1, 30,
		char(finch, "Modals hedge and nuance: 'Das mag sein' (may be), 'Er will alles gesehen haben' (he claims to have), 'Sie soll sehr klug sein' (she's said to be)."),
		mc("claim (subject's own)", "Er ___ alles gewusst haben. (claims to have)", "will", "will", "soll", "muss", "darf"),
		mc("hearsay", "Sie ___ sehr reich sein. (is said to be)", "soll", "soll", "will", "muss", "kann"),
		mc("concession", "Das ___ stimmen, aber … (may)", "mag", "mag", "muss", "will", "soll"),
		tr("Translate", "He claims to be ill", "Er will krank sein"),
		speak("Blaze", "Das mag sein, doch ich bin nicht überzeugt."),
	)
	addVocab(db, l,
		vw("angeblich", "allegedly", "Er ist angeblich reich.", "He's allegedly rich.", finch),
		vw("die Behauptung", "claim", "eine kühne Behauptung", "a bold claim", "Cora"),
		vw("das Gerücht", "rumour", "Gerüchten zufolge …", "according to rumours …", "Cora"),
		vw("einräumen", "to concede", "Ich räume ein, dass …", "I concede that …", finch),
	)

	// 83 — Nested clauses & rare connectors
	s = addSkillL(db, de, u6, "Verschachtelte Sätze", "geschweige denn, es sei denn …", "Link2", "#F5A623", 83, 2280)
	l = addLesson(db, s, "Komplexe Verbindungen", 1, 30,
		char(finch, "C2 handles nested clauses and rare connectors: 'geschweige denn' (let alone), 'es sei denn' (unless), 'sofern' (provided that), 'nicht zuletzt' (not least)."),
		mc("let alone", "Er kann kaum gehen, ___ laufen.", "geschweige denn", "geschweige denn", "sondern", "sowie", "weil"),
		mc("unless", "Ich komme, ___ es regnet.", "es sei denn", "es sei denn", "sofern", "obwohl", "damit"),
		mc("provided that", "___ du Zeit hast, treffen wir uns.", "Sofern", "Sofern", "Obwohl", "Indem", "Damit"),
		tr("Translate", "not least because of the costs", "nicht zuletzt wegen der Kosten"),
		speak("Blaze", "Ich helfe gern, es sei denn, ich bin verhindert."),
	)
	addVocab(db, l,
		vw("geschweige denn", "let alone", "kaum …, geschweige denn …", "hardly …, let alone …", finch),
		vw("es sei denn", "unless", "…, es sei denn, …", "…, unless …", finch),
		vw("sofern", "provided that", "sofern möglich", "provided possible", "Cora"),
		vw("nicht zuletzt", "not least", "nicht zuletzt deshalb", "not least for that reason", "Cora"),
	)

	// 84 — Word order & rhetoric
	s = addSkillL(db, de, u6, "Wortstellung & Rhetorik", "Rhythm, inversion, emphasis.", "PenLine", "#00C2A8", 84, 2320)
	l = addLesson(db, s, "Stilistische Wortstellung", 1, 30,
		char(finch, "Stylistic word order creates rhythm and force: fronting ('Kaum war er da, …'), inversion after a fronted element, and 'erst dann …' for emphasis."),
		mc("Inversion after 'Kaum'", "Kaum war er da, ___ das Telefon.", "klingelte", "klingelte", "es klingelte", "klingeln", "geklingelt"),
		mc("Front 'nie' (verb 2nd)", "___ habe ich das behauptet.", "Nie", "Nie", "Ich", "Habe", "Das"),
		fill("the … the", "Je länger ich nachdenke, ___ weniger verstehe ich.", "desto"),
		tr("Translate", "Only then did I understand", "Erst dann verstand ich"),
		speak("Blaze", "Kaum hatte ich es gesagt, bereute ich es schon."),
	)
	addVocab(db, l,
		vw("kaum", "hardly / scarcely", "kaum hatte er …", "scarcely had he …", finch),
		vw("hervorheben", "to emphasise", "etwas hervorheben", "to emphasise sth.", "Cora"),
		vw("der Rhythmus", "rhythm", "der Satzrhythmus", "the sentence rhythm", "Cora"),
		vw("bereuen", "to regret", "Ich bereue nichts.", "I regret nothing.", finch),
	)

	// 85 — Negation, emphasis & irony
	s = addSkillL(db, de, u6, "Negation & Ironie", "alles andere als, keineswegs …", "Quote", "#FF5C5C", 85, 2360)
	l = addLesson(db, s, "Feine Negation", 1, 30,
		char(finch, "Subtle negation and irony: 'alles andere als' (anything but), 'keineswegs' (by no means), 'nicht gerade' (not exactly) — often ironic understatement."),
		mc("anything but", "Das war ___ einfach.", "alles andere als", "alles andere als", "sehr", "ziemlich", "genau"),
		mc("by no means", "Das ist ___ sicher.", "keineswegs", "keineswegs", "durchaus", "sehr", "ziemlich"),
		mc("not exactly", "Er war ___ begeistert.", "nicht gerade", "nicht gerade", "sehr", "total", "echt"),
		tr("Translate", "That's by no means certain", "Das ist keineswegs sicher"),
		speak("Blaze", "Das war ja alles andere als einfach."),
	)
	addVocab(db, l,
		vw("keineswegs", "by no means", "keineswegs sicher", "by no means certain", finch),
		vw("alles andere als", "anything but", "alles andere als leicht", "anything but easy", finch),
		vw("nicht gerade", "not exactly", "nicht gerade billig", "not exactly cheap", "Cora"),
		vw("die Untertreibung", "understatement", "eine ironische Untertreibung", "an ironic understatement", "Cora"),
	)

	// 86 — Idioms
	s = addSkillL(db, de, u6, "Redewendungen", "Idioms & nuance.", "Quote", "#6C3FC5", 86, 2400)
	l = addLesson(db, s, "Redewendungen", 1, 30,
		char(finch, "Mastery shows in idioms: 'nur Bahnhof verstehen' = to understand nothing; 'die Daumen drücken' = to wish luck; 'ins Fettnäpfchen treten' = to put one's foot in it."),
		mc("Idiom meaning", "'Ich verstehe nur Bahnhof.'", "I don't understand a thing", "I don't understand a thing", "I love trains", "I'm leaving", "I'm late"),
		mc("Idiom meaning", "'ins Fettnäpfchen treten'", "To put one's foot in it", "To put one's foot in it", "To step in butter", "To cook", "To fall down"),
		mc("Idiom meaning", "'jemandem die Daumen drücken'", "To wish someone luck", "To wish someone luck", "To threaten", "To shake hands", "To leave"),
		tr("Translate (idiom)", "It's raining heavily", "Es regnet in Strömen"),
		speak("Blaze", "Ich drücke dir die Daumen!"),
	)
	addVocab(db, l,
		vw("die Redewendung", "idiom", "eine bekannte Redewendung", "a well-known idiom", finch),
		vw("ins Fettnäpfchen treten", "to put one's foot in it", "ins Fettnäpfchen treten", "to put one's foot in it", "Cora"),
		vw("der Ausdruck", "expression", "ein gehobener Ausdruck", "an elevated expression", "Cora"),
		vw("fließend", "fluent(ly)", "fließend sprechen", "to speak fluently", finch),
	)

	// 87 — Register & style
	s = addSkillL(db, de, u6, "Register & Stil", "Switching register at will.", "MessageCircle", "#17A3DD", 87, 2440)
	l = addLesson(db, s, "Register wechseln", 1, 30,
		char(finch, "C2 means switching register at will: formal ('Ich möchte Sie bitten …'), neutral ('Kannst du …'), colloquial ('Machst du mal …'), and academic. Match register to context."),
		mc("Most formal request", "Choose the most formal", "Könnten Sie mir bitte helfen?", "Könnten Sie mir bitte helfen?", "Hilf mir mal!", "Mach das!", "Na los!"),
		mc("Most casual greeting", "Choose the casual one", "Na, alles klar?", "Na, alles klar?", "Guten Tag", "Sehr geehrte Damen", "Grüß Gott"),
		mc("Academic register", "In an essay you'd write:", "Es lässt sich feststellen, dass …", "Es lässt sich feststellen, dass …", "Voll krass, dass …", "Ich find halt …", "Najaa …"),
		tr("Translate (formal)", "I would like to ask you a favour", "Ich möchte Sie um einen Gefallen bitten"),
		speak("Blaze", "Darf ich Sie um einen kleinen Gefallen bitten?"),
	)
	addVocab(db, l,
		vw("das Register", "register", "das Register wechseln", "to switch register", finch),
		vw("förmlich", "formal", "ein förmlicher Brief", "a formal letter", "Cora"),
		vw("umgangssprachlich", "colloquial", "umgangssprachlich gesagt", "colloquially speaking", "Cora"),
		vw("gehoben", "elevated", "gehobene Sprache", "elevated language", finch),
	)

	// 88 — Cohesion
	s = addSkillL(db, de, u6, "Kohäsion", "Linking long texts.", "Link2", "#F5A623", 88, 2480)
	l = addLesson(db, s, "Textverknüpfung", 1, 30,
		char(finch, "Long C2 texts cohere through reference and connectors: 'dies', 'im Folgenden', 'wie bereits erwähnt', 'einerseits/andererseits'."),
		mc("Reference word", "Das Problem ist komplex; ___ erschwert die Lösung.", "Dies", "Dies", "Welches", "Wer", "Wessen"),
		mc("as already …", "___ erwähnt, ist das Thema vielschichtig.", "Wie bereits", "Wie bereits", "Im Folgenden", "Dennoch", "Zudem"),
		fill("in the following", "Im ___ nenne ich drei Punkte. (Folgenden)", "Folgenden"),
		tr("Translate", "as mentioned above", "wie oben erwähnt"),
		speak("Blaze", "Wie bereits erwähnt, ist die Lage komplex."),
	)
	addVocab(db, l,
		vw("dies", "this", "dies bedeutet …", "this means …", finch),
		vw("im Folgenden", "in the following", "im Folgenden …", "in the following …", "Cora"),
		vw("wie bereits erwähnt", "as already mentioned", "wie bereits erwähnt", "as already mentioned", "Cora"),
		vw("der Bezug", "reference", "Bezug nehmen auf", "to refer to", finch),
	)

	// 89 — Writing (C2)
	s = addSkillL(db, de, u6, "Schreiben: Synthese (C2)", "Essays, reviews & syntheses.", "PenLine", "#00C2A8", 89, 2520)
	l = addLesson(db, s, "Anspruchsvolles Schreiben", 1, 30,
		char(finch, "C2 writing synthesises sources into a coherent, stylistically refined argument. Command register, vary sentence structure, and integrate viewpoints smoothly."),
		write("Write an essay (≈250 words): 'Inwiefern verändert die Digitalisierung die Arbeitswelt?' Analyse the issue, weigh different perspectives, and draw a nuanced conclusion.",
			"Die Digitalisierung hat die Arbeitswelt in den letzten Jahren tiefgreifend verändert — und dieser Wandel ist keineswegs abgeschlossen.\nEinerseits eröffnet sie neue Möglichkeiten: Tätigkeiten lassen sich ortsunabhängig erledigen, Prozesse werden effizienter, und ganz neue Berufsfelder entstehen. Andererseits gehen traditionelle Arbeitsplätze verloren, und der Druck, sich ständig weiterzubilden, wächst. Nicht zuletzt stellt sich die Frage, wie viel Kontrolle wir der Technik überlassen wollen.\nKritiker warnen zudem vor einer wachsenden Kluft zwischen gut und schlecht Qualifizierten. Befürworter hingegen betonen, dass jede technische Revolution langfristig mehr Beschäftigung geschaffen habe.\nMeines Erachtens kommt es entscheidend darauf an, den Wandel aktiv zu gestalten: durch Bildung, faire Regeln und den bewussten Einsatz der Technik. Die Digitalisierung ist somit weder Segen noch Fluch, sondern eine Gestaltungsaufgabe."),
		write("Write a review (≈180 words) of a book or film you know well, evaluating its strengths and weaknesses with a clear recommendation.",
			"Der Film überzeugt vor allem durch seine eindringliche Bildsprache und seine vielschichtigen Figuren. Die Handlung, die zunächst ruhig beginnt, gewinnt nach und nach an Spannung, ohne je ins Effekthascherische abzugleiten.\nBesonders gelungen ist das Zusammenspiel der Hauptdarsteller, das die emotionale Tiefe der Geschichte trägt. Schwächen zeigen sich allenfalls im Mittelteil, der etwas langatmig geraten ist.\nInsgesamt handelt es sich um ein anspruchsvolles Werk, das zum Nachdenken anregt. Wer intelligentes Kino schätzt, dem sei dieser Film uneingeschränkt empfohlen."),
		write("Read the claim 'Soziale Medien gefährden die Demokratie' and write a reasoned Stellungnahme (≈200 words) that weighs arguments and states your position.",
			"Die Behauptung, soziale Medien gefährdeten die Demokratie, ist ebenso verbreitet wie umstritten.\nFür diese These spricht, dass sich Falschinformationen über solche Plattformen rasend schnell verbreiten und dass Algorithmen die Nutzer in sogenannte Echokammern drängen. Folglich kann die öffentliche Debatte verzerrt werden.\nDagegen lässt sich einwenden, dass soziale Medien zugleich die Meinungsfreiheit stärken und auch jenen eine Stimme geben, die sonst ungehört blieben.\nMeines Erachtens liegt die Gefahr weniger in der Technik selbst als in ihrem unkritischen Gebrauch. Demokratien sollten daher nicht die Plattformen verbieten, sondern Medienkompetenz fördern und für mehr Transparenz sorgen."),
	)
	addVocab(db, l,
		vw("die Stellungnahme", "position statement", "eine Stellungnahme verfassen", "to write a statement", finch),
		vw("abwägen", "to weigh up", "Argumente abwägen", "to weigh arguments", "Cora"),
		vw("differenzieren", "to differentiate", "stärker differenzieren", "to differentiate more", "Cora"),
		vw("meines Erachtens", "in my view", "Meines Erachtens …", "In my view …", finch),
	)

	// 90 — Speaking (C2)
	s = addSkillL(db, de, u6, "Sprechen: Debatte (C2)", "Fluent, idiomatic, rhetorical.", "MessageCircle", "#6C3FC5", 90, 2560)
	addLesson(db, s, "Vortrag & Debatte", 1, 30,
		char(finch, "C2 speaking is spontaneous, precise and idiomatic. Structure a talk, argue rhetorically, concede gracefully ('Zwar …, doch …'), and respond to counterpoints with ease."),
		mc("Concede, then counter", "Choose the rhetorical move", "Zwar …, doch …", "Zwar …, doch …", "Ich heiße …", "Tschüss", "Wie bitte?"),
		speak("Blaze", "Im Folgenden möchte ich drei zentrale Aspekte beleuchten."),
		speak("Blaze", "Zwar mag das auf den ersten Blick überzeugen, doch bei näherer Betrachtung greift es zu kurz."),
		speak("Blaze", "Zusammenfassend lässt sich festhalten, dass die Frage vielschichtiger ist, als sie zunächst erscheint."),
	)
	addVocab(db, l,
		vw("beleuchten", "to shed light on", "einen Aspekt beleuchten", "to examine an aspect", finch),
		vw("vielschichtig", "multilayered", "ein vielschichtiges Thema", "a multilayered topic", "Cora"),
		vw("die Betrachtung", "consideration", "bei näherer Betrachtung", "on closer inspection", "Cora"),
		vw("überzeugend", "convincing", "ein überzeugendes Argument", "a convincing argument", finch),
	)

	addListeningL(db, de, u6, "Smalltalk auf Muttersprachler-Niveau",
		"A fast, idiomatic café chat. Listen carefully, then answer.", 6, 35,
		[]models.ListeningMatch{
			lm("Was geht ab?", "What's up? (casual)"),
			lm("Ich bin total im Stress", "I'm totally stressed"),
			lm("Mir reicht's", "I've had enough"),
			lm("Kopf hoch!", "Chin up!"),
		},
		[]models.ListeningLine{
			ln("Cora", "Na, Riko, was geht ab? Du siehst fix und fertig aus.", "Hey Riko, what's up? You look completely worn out."),
			ln("Riko", "Frag nicht. Ich bin total im Stress — die Prüfungen bringen mich noch um.", "Don't ask. I'm totally stressed — the exams will be the death of me."),
			ln("Cora", "Komm, halb so wild. Du hast doch wie ein Weltmeister gelernt.", "Come on, it's not that bad. You've studied like a champion."),
			ln("Riko", "Trotzdem, mir reicht's langsam. Ich brauche dringend Urlaub.", "Still, I've slowly had enough. I urgently need a holiday."),
			ln("Cora", "Kopf hoch! Nach der Prüfung lade ich dich auf einen Kaffee ein.", "Chin up! After the exam I'll treat you to a coffee."),
			ln("Riko", "Das lasse ich mir nicht zweimal sagen!", "I won't need to be told twice!"),
		},
		[]models.ListeningQuestion{
			lq("How does Riko look?", "Worn out", "Worn out", "Happy", "Bored", "Sleepy"),
			lq("Why is he stressed?", "Exams", "Exams", "Work", "Money", "Family"),
			lq("What does Cora offer?", "To treat him to a coffee", "To treat him to a coffee", "Money", "A book", "A ride"),
		},
	)
	addReadingL(db, de, u6, "Literarischer Auszug",
		"Read this literary passage, then answer.", 6, 35,
		[]models.ReadingLine{
			rl("Es war einer jener Abende, an denen die Stadt den Atem anzuhalten schien.", "It was one of those evenings on which the city seemed to hold its breath."),
			rl("Die Straßen lagen still, als hätte jemand die Welt für einen Augenblick beiseitegelegt.", "The streets lay silent, as if someone had set the world aside for a moment."),
			rl("Er ging, ohne ein Ziel zu haben, und gerade deshalb fühlte er sich seltsam frei.", "He walked without having a destination, and precisely for that reason he felt strangely free."),
			rl("Manchmal, dachte er, beginnt das Leben erst dort, wo die Pläne enden.", "Sometimes, he thought, life only begins where the plans end."),
		},
		[]models.ReadingQuestion{
			rq("How is the city described?", "As if holding its breath", "As if holding its breath", "Loud and busy", "On fire", "Crowded"),
			rq("Why did he feel free?", "He had no destination", "He had no destination", "He was rich", "He was young", "He was leaving"),
			rq("His final thought is that life begins…", "Where the plans end", "Where the plans end", "At dawn", "With money", "In the city"),
		},
	)
	addReadingL(db, de, u6, "Lerntipps für C2",
		"A smart study plan for C2. Read, then answer.", 7, 35,
		[]models.ReadingLine{
			rl("Auf dem C2-Niveau geht es nicht mehr um Korrektheit, sondern um Eleganz, Präzision und stilistische Wahl.", "At C2 level it's no longer about correctness, but about elegance, precision and stylistic choice."),
			rl("Lies anspruchsvolle Literatur und Qualitätsjournalismus und analysiere bewusst Ton, Ironie und Register.", "Read demanding literature and quality journalism and consciously analyse tone, irony and register."),
			rl("Übe das Zusammenfassen mehrerer Quellen zu einem kohärenten, eigenen Text.", "Practise synthesising several sources into one coherent text of your own."),
			rl("Sammle Redewendungen und Kollokationen und setze sie gezielt und situationsgerecht ein.", "Collect idioms and collocations and use them deliberately and appropriately."),
			rl("Wer dieses Niveau hält, braucht ständige Immersion und die Bereitschaft, auch feinste Nuancen zu hinterfragen.", "Whoever maintains this level needs constant immersion and the willingness to question even the finest nuances."),
		},
		[]models.ReadingQuestion{
			rq("What matters at C2?", "Elegance, precision, style", "Elegance, precision, style", "Only correctness", "Only vocabulary", "Nothing"),
			rq("What should you analyse when reading?", "Tone, irony, register", "Tone, irony, register", "Page count", "Font size", "Only grammar"),
			rq("What writing skill is highlighted?", "Synthesising several sources", "Synthesising several sources", "Copying texts", "Writing lists", "Translating word-by-word"),
		},
	)
}

// ===== French course =========================================================
func seedFrench(db *gorm.DB) {
	const fr = "fr"

	greet := addSkillL(db, fr, "Basics", "Salutations", "Greet people in French.", "Hand", "#6C3FC5", 1, 0)
	l := addLesson(db, greet, "Bonjour & Au revoir", 1, 15,
		char("Lumora", "Bonjour! First the words, then we practise. On y va!"),
		mc("Select the meaning", "Bonjour", "Hello", "Hello", "Goodbye", "Thanks", "Please"),
		tr("Translate this sentence", "Thank you", "Merci"),
		mc("What does 'Au revoir' mean?", "Au revoir", "Goodbye", "Hello", "Goodbye", "Sorry", "Welcome"),
		fill("Fill in the blank", "Comment ça ___ ? (How's it going?)", "va"),
		speak("Blaze", "Bonjour! Merci."),
	)
	addVocab(db, l,
		vw("Bonjour", "Hello / Good morning", "Bonjour ! Ça va ?", "Hello! How's it going?", "Lumora"),
		vw("Merci", "Thank you", "Merci beaucoup !", "Thank you very much!", "Cora"),
		vw("Au revoir", "Goodbye", "Au revoir, à demain.", "Goodbye, see you tomorrow.", "Lumora"),
		vw("S'il vous plaît", "Please", "Un café, s'il vous plaît.", "A coffee, please.", "Cora"),
	)

	cafe := addSkillL(db, fr, "Everyday Life", "Au Café", "Order food and drink.", "Coffee", "#F5A623", 2, 30)
	l = addLesson(db, cafe, "Au café", 1, 20,
		char("Blaze", "Un café ! I need un café. On commande !"),
		listen("Listen and choose the meaning", "Un café, s'il vous plaît.", "A coffee, please",
			"A coffee, please", "A tea, please", "A water, please", "The bill, please"),
		tr("Translate this sentence", "I would like water", "Je voudrais de l'eau"),
		mc("What does 'l'addition' mean?", "l'addition", "the bill", "the menu", "the bill", "the table", "the water"),
		fill("Fill in the blank", "Un café, s'il vous ___ . (please)", "plaît"),
		speak("Blaze", "Un café, s'il vous plaît."),
	)
	addVocab(db, l,
		vw("le café", "coffee", "Un café, s'il vous plaît.", "A coffee, please.", "Blaze"),
		vw("l'eau", "water", "De l'eau, s'il vous plaît.", "Water, please.", "Cora"),
		vw("l'addition", "the bill", "L'addition, s'il vous plaît.", "The bill, please.", "Lumora"),
		vw("s'il vous plaît", "please", "Merci, s'il vous plaît.", "Thanks, please.", "Cora"),
	)

	travel := addSkillL(db, fr, "Getting Around", "Voyage", "Find your way.", "Plane", "#06AECE", 3, 60)
	l = addLesson(db, travel, "À l'aéroport", 1, 20,
		char("Mira", "Du calme. Écoute. Le voyage commence par un mot."),
		mc("Select the meaning", "l'aéroport", "the airport", "the airport", "the station", "the hotel", "the street"),
		tr("Translate this sentence", "Where is the gate?", "Où est la porte ?"),
		mc("Select the meaning", "à gauche", "to the left", "to the left", "to the right", "up", "down"),
		speak("Blaze", "Où est la porte ? À gauche."),
	)
	addVocab(db, l,
		vw("l'aéroport", "the airport", "À l'aéroport, s'il vous plaît.", "To the airport, please.", "Mira"),
		vw("le vol", "the flight", "Mon vol est à trois heures.", "My flight is at three.", "Mira"),
		vw("à gauche", "to the left", "Tournez à gauche.", "Turn left.", "Lumora"),
		vw("à droite", "to the right", "Tournez à droite.", "Turn right.", "Lumora"),
	)

	convo := addSkillL(db, fr, "Connecting", "Faire connaissance", "Have your first chat.", "MessageCircle", "#17A3DD", 4, 90)
	l = addLesson(db, convo, "Conversation", 1, 20,
		char("Riko", "Tu veux me parler ? Bon — montre-moi ce que tu sais."),
		mc("Select the meaning", "Comment t'appelles-tu ?", "What's your name?", "What's your name?", "How are you?", "Where are you?", "How old are you?"),
		tr("Translate this sentence", "My name is Ana", "Je m'appelle Ana"),
		mc("Select the meaning", "D'où viens-tu ?", "Where are you from?", "Where are you from?", "What do you do?", "Where do you live?", "How are you?"),
		speak("Blaze", "Je m'appelle Ana. D'où viens-tu ?"),
	)
	addVocab(db, l,
		vw("Comment t'appelles-tu ?", "What's your name?", "Bonjour, comment t'appelles-tu ?", "Hi, what's your name?", "Riko"),
		vw("Je m'appelle ...", "My name is ...", "Je m'appelle Ana.", "My name is Ana.", "Lumora"),
		vw("D'où viens-tu ?", "Where are you from?", "D'où viens-tu ?", "Where are you from?", "Riko"),
	)

	addListeningL(db, fr, "Basics", "Un bonjour amical",
		"Listen to Lumora greet the professor.", 1, 15,
		[]models.ListeningMatch{
			lm("Bonjour", "Hello"),
			lm("Comment ça va ?", "How are you?"),
			lm("Merci", "Thank you"),
			lm("Au revoir", "Goodbye"),
		},
		[]models.ListeningLine{
			ln("Lumora", "Bonjour, professeur !", "Hello, professor!"),
			ln("Professor Finch", "Bonjour. Comment ça va ?", "Hello. How are you?"),
			ln("Lumora", "Très bien, merci !", "Very well, thank you!"),
			ln("Professor Finch", "Parfait. Au revoir !", "Perfect. Goodbye!"),
		},
		[]models.ListeningQuestion{
			lq("How does Lumora greet?", "Bonjour", "Bonjour", "Bonne nuit", "Au revoir", "Merci"),
			lq("How does she feel?", "Very well", "Very well", "Bad", "Tired", "Sad"),
		},
	)

	addListeningL(db, fr, "Everyday Life", "Au café",
		"Blaze orders at the cafe.", 2, 15,
		[]models.ListeningMatch{
			lm("le café", "coffee"),
			lm("l'eau", "water"),
			lm("l'addition", "the bill"),
			lm("s'il vous plaît", "please"),
		},
		[]models.ListeningLine{
			ln("Blaze", "Un café, s'il vous plaît.", "A coffee, please."),
			ln("Cora", "Et de l'eau, s'il vous plaît.", "And water, please."),
			ln("Blaze", "L'addition, s'il vous plaît.", "The bill, please."),
		},
		[]models.ListeningQuestion{
			lq("What does Blaze order?", "A coffee", "A coffee", "A tea", "A juice", "A beer"),
			lq("What do they ask for at the end?", "The bill", "The bill", "The menu", "A table", "Water"),
		},
	)

	addReadingL(db, fr, "Basics", "Une note",
		"Read Ana's note.", 1, 15,
		[]models.ReadingLine{
			rl("Bonjour ! Je m'appelle Ana.", "Hello! My name is Ana."),
			rl("Comment ça va ?", "How are you?"),
			rl("Merci et au revoir.", "Thank you and goodbye."),
		},
		[]models.ReadingQuestion{
			rq("What is the writer's name?", "Ana", "Ana", "Lumora", "Cora", "Mira"),
			rq("How does she end the note?", "Goodbye", "Goodbye", "Hello", "Please", "Sorry"),
		},
	)

	addReadingL(db, fr, "Everyday Life", "Au café",
		"Read the cafe order.", 2, 15,
		[]models.ReadingLine{
			rl("Je voudrais un café.", "I would like a coffee."),
			rl("Et de l'eau, s'il vous plaît.", "And water, please."),
			rl("L'addition, s'il vous plaît.", "The bill, please."),
		},
		[]models.ReadingQuestion{
			rq("What do they want first?", "A coffee", "A coffee", "A tea", "Juice", "Wine"),
			rq("What do they ask for last?", "The bill", "The bill", "The menu", "A table", "A spoon"),
		},
	)
}

func seedCharacters(db *gorm.DB) {
	characters := []models.Character{
		{Name: "Lumora", Species: "Fennec Fox", Role: "Your Guide", Personality: "Curious, playful, endlessly supportive.", Color: "#6C3FC5", Emoji: "🦊"},
		{Name: "Professor Finch", Species: "Eagle", Role: "Grammar Teacher", Personality: "Strict but secretly warm. Sighs dramatically at errors.", Color: "#8B6F47", Emoji: "🦅"},
		{Name: "Cora", Species: "Octopus", Role: "Vocabulary Friend", Personality: "Goofy, chaotic, can't stop making puns.", Color: "#00C2A8", Emoji: "🐙"},
		{Name: "Blaze", Species: "Fire Spirit", Role: "Speaking Coach", Personality: "Hypes you up. MORE FIRE!", Color: "#FF5C5C", Emoji: "🔥"},
		{Name: "Mira", Species: "Snow Leopard", Role: "Listening Guide", Personality: "Serene and wise. Loves music and poetry.", Color: "#9090A0", Emoji: "🐆"},
		{Name: "Riko", Species: "Red Panda", Role: "Your Rival", Personality: "Smug but lovable. Secretly rooting for you.", Color: "#F5A623", Emoji: "🐼"},
		{Name: "Zephyr", Species: "Wind Spirit", Role: "Writing Mentor", Personality: "Philosophical, poetic, occasionally pretentious.", Color: "#17A3DD", Emoji: "🌬️"},
		{Name: "Nana", Species: "Giant Tortoise", Role: "Wise Elder", Personality: "Slow-spoken and reassuring. Every journey takes time.", Color: "#06AECE", Emoji: "🐢"},
		{Name: "Pip", Species: "Hedgehog", Role: "Quest Giver", Personality: "Frantic, fast-talking, perpetually late.", Color: "#F5A623", Emoji: "🦔"},
	}
	db.Create(&characters)
}

func seedQuests(db *gorm.DB) {
	quests := []models.Quest{
		{Title: "Complete 2 lessons", Description: "Finish two lessons today", Icon: "📚", XPReward: 15, Target: 2},
		{Title: "Earn 30 XP", Description: "Rack up 30 XP before you rest", Icon: "⚡", XPReward: 10, Target: 30},
		{Title: "Practice speaking", Description: "Try a speaking exercise with Blaze", Icon: "🎤", XPReward: 20, Target: 1},
	}
	db.Create(&quests)
}

// --- small builders that keep the curriculum below compact & readable --------

func opts(items ...string) string {
	b, _ := json.Marshal(items)
	return string(b)
}

func char(name, q string) models.Exercise {
	return models.Exercise{Type: models.ExerciseCharacter, Character: name, Question: q}
}
func mc(prompt, q, correct string, options ...string) models.Exercise {
	return models.Exercise{Type: models.ExerciseMultipleChoice, Prompt: prompt, Question: q, CorrectAnswer: correct, OptionsJSON: opts(options...)}
}
func listen(prompt, q, correct string, options ...string) models.Exercise {
	return models.Exercise{Type: models.ExerciseListen, Prompt: prompt, Question: q, CorrectAnswer: correct, OptionsJSON: opts(options...)}
}
func match(prompt, q, correct string, options ...string) models.Exercise {
	return models.Exercise{Type: models.ExerciseMatch, Prompt: prompt, Question: q, CorrectAnswer: correct, OptionsJSON: opts(options...)}
}
func tr(prompt, q, a string) models.Exercise {
	return models.Exercise{Type: models.ExerciseTranslate, Prompt: prompt, Question: q, CorrectAnswer: a}
}
func fill(prompt, q, a string) models.Exercise {
	return models.Exercise{Type: models.ExerciseFill, Prompt: prompt, Question: q, CorrectAnswer: a}
}
func speak(name, q string) models.Exercise {
	return models.Exercise{Type: models.ExerciseSpeak, Character: name, Prompt: "Say it out loud", Question: q, CorrectAnswer: q}
}

// write builds a free-text writing task (e.g. a short email). CorrectAnswer
// holds an example answer the learner can reveal to compare.
func write(task, sample string) models.Exercise {
	return models.Exercise{Type: models.ExerciseWrite, Prompt: "Write your answer", Question: task, CorrectAnswer: sample}
}

// vw builds a vocabulary item (the "learn the word" phase).
func vw(word, translation, example, exampleTr, speaker string) models.VocabItem {
	return models.VocabItem{Word: word, Translation: translation, Example: example, ExampleTranslation: exampleTr, Speaker: speaker}
}

func addSkill(db *gorm.DB, unit, title, desc, icon, color string, order, reqXP int) uint {
	s := models.Skill{
		Language: "es", Unit: unit, Title: title, Description: desc,
		Icon: icon, Color: color, OrderIndex: order, RequiredXP: reqXP,
	}
	db.Create(&s)
	return s.ID
}

// addLesson creates a lesson with its ordered exercises and returns the new id.
func addLesson(db *gorm.DB, skillID uint, title string, order, xp int, exs ...models.Exercise) uint {
	l := models.Lesson{SkillID: skillID, Title: title, OrderIndex: order, XPReward: xp}
	db.Create(&l)
	for i := range exs {
		exs[i].LessonID = l.ID
		exs[i].OrderIndex = i + 1
	}
	db.Create(&exs)
	return l.ID
}

func addVocab(db *gorm.DB, lessonID uint, items ...models.VocabItem) {
	for i := range items {
		items[i].LessonID = lessonID
		items[i].OrderIndex = i + 1
	}
	db.Create(&items)
}

// seedSpanish builds a full Spanish course: four units, eight skills with a
// rising XP requirement. Every lesson teaches its vocabulary first, then quizzes
// it, and ends with a speaking exercise. Icons are lucide names.
func seedSpanish(db *gorm.DB) {
	// ===== Unit 1 — Basics =====
	greetings := addSkill(db, "Basics", "Greetings", "Say hello like a local.", "Hand", "#6C3FC5", 1, 0)

	l := addLesson(db, greetings, "Hello & Goodbye", 1, 15,
		char("Lumora", "¡Hola! First we'll learn the words, then we'll practise. Ready?"),
		mc("Select the meaning", "Hola", "Hello", "Hello", "Goodbye", "Thanks", "Please"),
		tr("Translate this sentence", "Good morning", "Buenos días"),
		mc("What does 'Adiós' mean?", "Adiós", "Goodbye", "Hello", "Goodbye", "Sorry", "Welcome"),
		fill("Fill in the blank", "¿Cómo ___ ? (How are you?)", "estás"),
		speak("Blaze", "¡Hola! Buenos días."),
	)
	addVocab(db, l,
		vw("Hola", "Hello", "¡Hola! ¿Qué tal?", "Hi! How's it going?", "Lumora"),
		vw("Buenos días", "Good morning", "Buenos días, señora.", "Good morning, ma'am.", "Cora"),
		vw("Adiós", "Goodbye", "Adiós, hasta mañana.", "Goodbye, see you tomorrow.", "Lumora"),
		vw("¿Cómo estás?", "How are you?", "Hola, ¿cómo estás?", "Hi, how are you?", "Cora"),
	)

	l = addLesson(db, greetings, "Polite Words", 2, 15,
		char("Cora", "Magic words time! These open every door (and I have eight hands for doors)."),
		mc("Select the meaning", "Gracias", "Thank you", "Thank you", "Please", "Sorry", "Hello"),
		mc("Select the meaning", "Por favor", "Please", "Please", "Thank you", "Goodbye", "Welcome"),
		tr("Translate this sentence", "You're welcome", "De nada"),
		fill("Fill in the blank", "Lo ___ . (I'm sorry)", "siento"),
		speak("Blaze", "Gracias, por favor."),
	)
	addVocab(db, l,
		vw("Gracias", "Thank you", "Muchas gracias.", "Thank you very much.", "Cora"),
		vw("Por favor", "Please", "Un café, por favor.", "A coffee, please.", "Cora"),
		vw("De nada", "You're welcome", "—Gracias. —De nada.", "—Thanks. —You're welcome.", "Lumora"),
		vw("Lo siento", "I'm sorry", "Lo siento mucho.", "I'm very sorry.", "Lumora"),
	)

	essentials := addSkill(db, "Basics", "Essentials", "The words you'll use every day.", "Sparkles", "#17A3DD", 2, 20)

	l = addLesson(db, essentials, "Yes, No & Maybe", 1, 15,
		char("Professor Finch", "Precision matters. 'Sí' and 'no' are small words with great power."),
		mc("Select the meaning", "Sí", "Yes", "Yes", "No", "Maybe", "Never"),
		mc("Select the meaning", "No", "No", "Yes", "No", "Maybe", "Always"),
		tr("Translate this sentence", "Maybe", "Quizás"),
		fill("Fill in the blank", "Tal ___ . (Maybe)", "vez"),
		speak("Blaze", "Sí, quizás."),
	)
	addVocab(db, l,
		vw("Sí", "Yes", "Sí, por favor.", "Yes, please.", "Cora"),
		vw("No", "No", "No, gracias.", "No, thank you.", "Cora"),
		vw("Quizás", "Maybe", "Quizás mañana.", "Maybe tomorrow.", "Lumora"),
	)

	l = addLesson(db, essentials, "Numbers 1–5", 2, 15,
		char("Cora", "Counting with an octopus is easy — I lose track after eight!"),
		mc("Select the meaning", "uno", "one", "one", "two", "three", "four"),
		mc("Select the meaning", "tres", "three", "one", "two", "three", "five"),
		tr("Translate this sentence", "five", "cinco"),
		fill("Fill in the blank", "uno, dos, ___ , cuatro", "tres"),
		speak("Blaze", "uno, dos, tres, cuatro, cinco."),
	)
	addVocab(db, l,
		vw("uno", "one", "Tengo uno.", "I have one.", "Cora"),
		vw("dos", "two", "Dos cafés, por favor.", "Two coffees, please.", "Cora"),
		vw("tres", "three", "Son las tres.", "It's three o'clock.", "Lumora"),
		vw("cuatro", "four", "Cuatro amigos.", "Four friends.", "Lumora"),
		vw("cinco", "five", "Cinco minutos.", "Five minutes.", "Lumora"),
	)

	// ===== Unit 2 — Everyday Life =====
	food := addSkill(db, "Everyday Life", "Food & Cafe", "Order anything, anywhere.", "Coffee", "#F5A623", 3, 45)

	l = addLesson(db, food, "Ordering at a Cafe", 1, 20,
		char("Blaze", "We're going to a cafe. I need a coffee or I might COMBUST. Words first!"),
		listen("Listen and choose the meaning", "Un café con leche, por favor.", "A coffee with milk, please",
			"A coffee with milk, please", "A tea with sugar, please", "A water, please", "The bill, please"),
		tr("Translate this sentence", "I would like a croissant, please.", "Quiero un croissant, por favor"),
		mc("What does 'la cuenta' mean?", "la cuenta", "The bill", "The menu", "The bill", "The waiter", "The tip"),
		match("What is 'agua'?", "agua", "water", "water", "sugar", "butter", "bread"),
		fill("Fill in the blank", "¿Me trae la ___ , por favor? (the bill)", "cuenta"),
		speak("Blaze", "Un café con leche, por favor."),
	)
	addVocab(db, l,
		vw("el café", "coffee", "Un café, por favor.", "A coffee, please.", "Blaze"),
		vw("con leche", "with milk", "Café con leche.", "Coffee with milk.", "Cora"),
		vw("el agua", "water", "Un agua, por favor.", "A water, please.", "Cora"),
		vw("la cuenta", "the bill", "La cuenta, por favor.", "The bill, please.", "Lumora"),
	)

	l = addLesson(db, food, "At the Restaurant", 2, 20,
		char("Cora", "Tonight, we dine! I've memorised the menu with all eight arms."),
		mc("Select the meaning", "el menú", "the menu", "the menu", "the bill", "the table", "the kitchen"),
		tr("Translate this sentence", "The table, please", "La mesa, por favor"),
		mc("Select the meaning", "delicioso", "delicious", "delicious", "expensive", "cold", "spicy"),
		fill("Fill in the blank", "La comida está ___ . (delicious)", "deliciosa"),
		speak("Blaze", "La mesa, por favor. El menú, por favor."),
	)
	addVocab(db, l,
		vw("el menú", "the menu", "El menú, por favor.", "The menu, please.", "Cora"),
		vw("la mesa", "the table", "Una mesa para dos.", "A table for two.", "Cora"),
		vw("delicioso", "delicious", "¡Qué delicioso!", "How delicious!", "Lumora"),
	)

	shopping := addSkill(db, "Everyday Life", "Shopping", "Markets, prices and bargains.", "ShoppingBag", "#00C2A8", 4, 80)

	l = addLesson(db, shopping, "At the Market", 1, 20,
		char("Cora", "Ooh, shopping! Eight arms means eight bags. Let's LEARN to bargain!"),
		mc("Select the meaning", "¿Cuánto cuesta?", "How much is it?", "How much is it?", "Where is it?", "What is it?", "Who is it?"),
		tr("Translate this sentence", "I want this", "Quiero esto"),
		mc("Select the meaning", "barato", "cheap", "cheap", "expensive", "free", "broken"),
		fill("Fill in the blank", "¿___ cuesta? (How much)", "Cuánto"),
		speak("Blaze", "¿Cuánto cuesta? Quiero esto."),
	)
	addVocab(db, l,
		vw("¿Cuánto cuesta?", "How much is it?", "¿Cuánto cuesta esto?", "How much is this?", "Cora"),
		vw("barato", "cheap", "Es muy barato.", "It's very cheap.", "Cora"),
		vw("caro", "expensive", "Es muy caro.", "It's very expensive.", "Lumora"),
	)

	// ===== Unit 3 — Getting Around =====
	travel := addSkill(db, "Getting Around", "Travel", "Find your way anywhere.", "Plane", "#06AECE", 5, 115)

	l = addLesson(db, travel, "At the Airport", 1, 20,
		char("Mira", "Breathe. Listen closely. The journey begins with a single word."),
		mc("Select the meaning", "el aeropuerto", "the airport", "the airport", "the station", "the hotel", "the street"),
		tr("Translate this sentence", "Where is the gate?", "¿Dónde está la puerta?"),
		mc("Select the meaning", "el vuelo", "the flight", "the flight", "the ticket", "the seat", "the bag"),
		speak("Blaze", "¿Dónde está la puerta? El vuelo."),
	)
	addVocab(db, l,
		vw("el aeropuerto", "the airport", "Voy al aeropuerto.", "I'm going to the airport.", "Mira"),
		vw("el vuelo", "the flight", "Mi vuelo es a las dos.", "My flight is at two.", "Mira"),
		vw("la puerta", "the gate", "¿Dónde está la puerta?", "Where is the gate?", "Lumora"),
	)

	l = addLesson(db, travel, "Directions", 2, 20,
		char("Mira", "Lost? Good. Being lost is how we learn the streets."),
		mc("Select the meaning", "a la derecha", "to the right", "to the right", "to the left", "straight ahead", "back"),
		mc("Select the meaning", "a la izquierda", "to the left", "to the right", "to the left", "near", "far"),
		tr("Translate this sentence", "Where is the hotel?", "¿Dónde está el hotel?"),
		fill("Fill in the blank", "Todo ___ . (Straight ahead)", "recto"),
		speak("Blaze", "A la derecha, luego todo recto."),
	)
	addVocab(db, l,
		vw("a la derecha", "to the right", "Gira a la derecha.", "Turn right.", "Mira"),
		vw("a la izquierda", "to the left", "Gira a la izquierda.", "Turn left.", "Mira"),
		vw("todo recto", "straight ahead", "Sigue todo recto.", "Keep going straight.", "Lumora"),
	)

	family := addSkill(db, "Getting Around", "Family & People", "Talk about the people you love.", "Users", "#6C3FC5", 6, 145)

	l = addLesson(db, family, "Family Members", 1, 20,
		char("Nana", "Family is the slow, strong root of every story. Let us name them."),
		mc("Select the meaning", "la madre", "the mother", "the mother", "the father", "the sister", "the friend"),
		mc("Select the meaning", "el hermano", "the brother", "the brother", "the uncle", "the son", "the cousin"),
		tr("Translate this sentence", "my family", "mi familia"),
		fill("Fill in the blank", "el ___ (the father)", "padre"),
		speak("Blaze", "Mi madre, mi padre, mi hermano."),
	)
	addVocab(db, l,
		vw("la madre", "the mother", "Mi madre es amable.", "My mother is kind.", "Nana"),
		vw("el padre", "the father", "Mi padre trabaja.", "My father works.", "Nana"),
		vw("el hermano", "the brother", "Tengo un hermano.", "I have a brother.", "Lumora"),
		vw("la familia", "the family", "Mi familia es grande.", "My family is big.", "Lumora"),
	)

	// ===== Unit 4 — Connecting =====
	convo := addSkill(db, "Connecting", "Conversations", "Hold your first real chat.", "MessageCircle", "#17A3DD", 7, 175)

	l = addLesson(db, convo, "Small Talk", 1, 20,
		char("Riko", "So you think you can talk to me? Hmph. Learn the words — then prove it."),
		mc("Select the meaning", "¿Cómo te llamas?", "What's your name?", "What's your name?", "How are you?", "Where are you?", "How old are you?"),
		tr("Translate this sentence", "My name is Ana", "Me llamo Ana"),
		mc("Select the meaning", "¿De dónde eres?", "Where are you from?", "Where are you from?", "What do you do?", "Where do you live?", "How are you?"),
		fill("Fill in the blank", "Soy ___ España. (I'm from Spain)", "de"),
		speak("Blaze", "Me llamo Ana. ¿De dónde eres?"),
	)
	addVocab(db, l,
		vw("¿Cómo te llamas?", "What's your name?", "Hola, ¿cómo te llamas?", "Hi, what's your name?", "Riko"),
		vw("Me llamo...", "My name is...", "Me llamo Ana.", "My name is Ana.", "Lumora"),
		vw("¿De dónde eres?", "Where are you from?", "¿De dónde eres?", "Where are you from?", "Riko"),
	)

	romance := addSkill(db, "Connecting", "Romance", "Words that make hearts glow.", "Heart", "#FF5C5C", 8, 205)

	l = addLesson(db, romance, "Sweet Nothings", 1, 25,
		char("Zephyr", "Ah, language as poetry. A well-formed sentence is a small eternity."),
		mc("Select the meaning", "Te quiero", "I love you", "I love you", "I miss you", "See you soon", "Good night"),
		tr("Translate this sentence", "You are beautiful", "Eres hermosa"),
		mc("Select the meaning", "mi amor", "my love", "my love", "my friend", "my dear", "my life"),
		speak("Blaze", "Te quiero, mi amor."),
	)
	addVocab(db, l,
		vw("Te quiero", "I love you", "Te quiero mucho.", "I love you a lot.", "Zephyr"),
		vw("mi amor", "my love", "Buenas noches, mi amor.", "Good night, my love.", "Zephyr"),
		vw("hermosa", "beautiful", "Eres muy hermosa.", "You are very beautiful.", "Lumora"),
	)
}

// --- listening sessions ------------------------------------------------------

func ln(character, text, translation string) models.ListeningLine {
	return models.ListeningLine{Character: character, Text: text, Translation: translation}
}
func lq(question, correct string, options ...string) models.ListeningQuestion {
	return models.ListeningQuestion{Prompt: "What did you hear?", Question: question, CorrectAnswer: correct, OptionsJSON: opts(options...)}
}
func lm(word, translation string) models.ListeningMatch {
	return models.ListeningMatch{Word: word, Translation: translation}
}

func addListening(db *gorm.DB, unit, title, desc string, order, xp int, matches []models.ListeningMatch, lines []models.ListeningLine, qs []models.ListeningQuestion) {
	s := models.ListeningSession{Language: "es", Unit: unit, Title: title, Description: desc, OrderIndex: order, XPReward: xp}
	db.Create(&s)
	for i := range matches {
		matches[i].SessionID = s.ID
		matches[i].OrderIndex = i + 1
	}
	db.Create(&matches)
	for i := range lines {
		lines[i].SessionID = s.ID
		lines[i].OrderIndex = i + 1
	}
	db.Create(&lines)
	for i := range qs {
		qs[i].SessionID = s.ID
		qs[i].OrderIndex = i + 1
	}
	db.Create(&qs)
}

// seedListening adds one verbal listening session per unit. Each is a short
// dialogue voiced by different characters, then comprehension questions.
func seedListening(db *gorm.DB) {
	addListening(db, "Basics", "A Friendly Hello",
		"Listen to Lumora greet Professor Finch, then answer.", 1, 15,
		[]models.ListeningMatch{
			lm("Buenos días", "Good morning"),
			lm("¿Cómo estás?", "How are you?"),
			lm("gracias", "thank you"),
			lm("Adiós", "Goodbye"),
		},
		[]models.ListeningLine{
			ln("Lumora", "¡Hola, profesor! Buenos días.", "Hello, professor! Good morning."),
			ln("Professor Finch", "Buenos días. ¿Cómo estás?", "Good morning. How are you?"),
			ln("Lumora", "Muy bien, gracias. ¿Y usted?", "Very well, thank you. And you?"),
			ln("Professor Finch", "Bien, bien. ¡Adiós!", "Well, well. Goodbye!"),
		},
		[]models.ListeningQuestion{
			lq("How does Lumora greet the professor?", "Buenos días", "Buenos días", "Buenas noches", "Adiós", "Gracias"),
			lq("What does Lumora say she is?", "Muy bien", "Muy bien", "Mal", "Cansada", "Triste"),
			lq("How does the professor say goodbye?", "Adiós", "Adiós", "Hola", "Por favor", "De nada"),
		},
	)

	addListening(db, "Everyday Life", "At the Cafe",
		"Blaze orders for the table. Listen, then answer.", 2, 15,
		[]models.ListeningMatch{
			lm("café con leche", "coffee with milk"),
			lm("agua", "water"),
			lm("¿Cuánto cuesta?", "How much is it?"),
			lm("la cuenta", "the bill"),
		},
		[]models.ListeningLine{
			ln("Blaze", "¡Hola! Un café con leche, por favor.", "Hello! A coffee with milk, please."),
			ln("Cora", "Para mí, un agua, por favor.", "For me, a water, please."),
			ln("Blaze", "¿Cuánto cuesta?", "How much is it?"),
			ln("Cora", "Y la cuenta, por favor.", "And the bill, please."),
		},
		[]models.ListeningQuestion{
			lq("What does Blaze order?", "A coffee with milk", "A coffee with milk", "A tea", "A water", "A croissant"),
			lq("What does Cora ask for to drink?", "Water", "Water", "Coffee", "Juice", "Milk"),
			lq("What do they ask for at the end?", "The bill", "The bill", "The menu", "A table", "The waiter"),
		},
	)

	addListening(db, "Getting Around", "Finding the Gate",
		"Mira helps a traveller at the airport. Listen, then answer.", 3, 15,
		[]models.ListeningMatch{
			lm("la puerta", "the gate"),
			lm("a la derecha", "to the right"),
			lm("todo recto", "straight ahead"),
			lm("el vuelo", "the flight"),
		},
		[]models.ListeningLine{
			ln("Mira", "Disculpe, ¿dónde está la puerta?", "Excuse me, where is the gate?"),
			ln("Riko", "A la derecha, luego todo recto.", "To the right, then straight ahead."),
			ln("Mira", "¿Y mi vuelo?", "And my flight?"),
			ln("Riko", "Tu vuelo es a las tres.", "Your flight is at three."),
		},
		[]models.ListeningQuestion{
			lq("What is Mira looking for?", "The gate", "The gate", "The hotel", "The bill", "The menu"),
			lq("Which way should she go first?", "To the right", "To the right", "To the left", "Back", "Straight only"),
			lq("When is the flight?", "At three", "At three", "At two", "At five", "At one"),
		},
	)

	addListening(db, "Connecting", "Nice to Meet You",
		"Riko and Zephyr make small talk. Listen, then answer.", 4, 20,
		[]models.ListeningMatch{
			lm("¿cómo te llamas?", "what's your name?"),
			lm("me llamo", "my name is"),
			lm("¿de dónde eres?", "where are you from?"),
			lm("mucho gusto", "nice to meet you"),
		},
		[]models.ListeningLine{
			ln("Riko", "Hola, ¿cómo te llamas?", "Hi, what's your name?"),
			ln("Zephyr", "Me llamo Zephyr. ¿Y tú?", "My name is Zephyr. And you?"),
			ln("Riko", "Soy Riko. ¿De dónde eres?", "I'm Riko. Where are you from?"),
			ln("Zephyr", "Soy de España. Mucho gusto.", "I'm from Spain. Nice to meet you."),
		},
		[]models.ListeningQuestion{
			lq("What is the second speaker's name?", "Zephyr", "Zephyr", "Riko", "Ana", "Mira"),
			lq("Where is Zephyr from?", "Spain", "Spain", "Mexico", "France", "Peru"),
			lq("What does the first speaker ask first?", "What's your name?", "What's your name?", "How are you?", "How old are you?", "Where do you live?"),
		},
	)
}

// --- reading sessions --------------------------------------------------------

func rl(text, translation string) models.ReadingLine {
	return models.ReadingLine{Text: text, Translation: translation}
}
func rq(question, correct string, options ...string) models.ReadingQuestion {
	return models.ReadingQuestion{Prompt: "Read, then answer", Question: question, CorrectAnswer: correct, OptionsJSON: opts(options...)}
}

func addReading(db *gorm.DB, unit, title, desc string, order, xp int, lines []models.ReadingLine, qs []models.ReadingQuestion) {
	s := models.ReadingSession{Language: "es", Unit: unit, Title: title, Description: desc, OrderIndex: order, XPReward: xp}
	db.Create(&s)
	for i := range lines {
		lines[i].SessionID = s.ID
		lines[i].OrderIndex = i + 1
	}
	db.Create(&lines)
	for i := range qs {
		qs[i].SessionID = s.ID
		qs[i].OrderIndex = i + 1
	}
	db.Create(&qs)
}

// seedReading adds one short reading passage per unit, each in Spanish with a
// sentence-by-sentence translation and comprehension questions.
func seedReading(db *gorm.DB) {
	addReading(db, "Basics", "Una nota de Ana",
		"Read Ana's little note, then answer.", 1, 15,
		[]models.ReadingLine{
			rl("¡Hola! Me llamo Ana.", "Hello! My name is Ana."),
			rl("Buenos días a todos.", "Good morning, everyone."),
			rl("Gracias y adiós.", "Thank you and goodbye."),
		},
		[]models.ReadingQuestion{
			rq("What is the writer's name?", "Ana", "Ana", "Lumora", "Cora", "Mira"),
			rq("What time of day does she greet?", "Morning", "Morning", "Night", "Afternoon", "Evening"),
		},
	)

	addReading(db, "Everyday Life", "En el café",
		"Read the cafe scene, then answer.", 2, 15,
		[]models.ReadingLine{
			rl("Quiero un café con leche.", "I want a coffee with milk."),
			rl("También quiero un agua, por favor.", "I also want a water, please."),
			rl("¿Cuánto cuesta? La cuenta, por favor.", "How much is it? The bill, please."),
		},
		[]models.ReadingQuestion{
			rq("What does the person want to drink first?", "Coffee with milk", "Coffee with milk", "Tea", "Juice", "Wine"),
			rq("What do they ask for at the end?", "The bill", "The bill", "The menu", "A table", "A spoon"),
		},
	)

	addReading(db, "Getting Around", "En el aeropuerto",
		"Read the airport notice, then answer.", 3, 15,
		[]models.ReadingLine{
			rl("El aeropuerto está a la derecha.", "The airport is to the right."),
			rl("La puerta nueve está todo recto.", "Gate nine is straight ahead."),
			rl("El vuelo es a las tres.", "The flight is at three."),
		},
		[]models.ReadingQuestion{
			rq("Where is gate nine?", "Straight ahead", "Straight ahead", "To the left", "Behind", "Upstairs"),
			rq("When is the flight?", "At three", "At three", "At two", "At nine", "At noon"),
		},
	)

	addReading(db, "Connecting", "Mucho gusto",
		"Read the introduction, then answer.", 4, 20,
		[]models.ReadingLine{
			rl("Hola, me llamo Zephyr.", "Hello, my name is Zephyr."),
			rl("Soy de España.", "I am from Spain."),
			rl("Mucho gusto. Te quiero, amigo.", "Nice to meet you. I love you, friend."),
		},
		[]models.ReadingQuestion{
			rq("Where is Zephyr from?", "Spain", "Spain", "Mexico", "Peru", "Chile"),
			rq("How does Zephyr end the note?", "Nice to meet you", "Nice to meet you", "Goodbye forever", "See you never", "Good night"),
		},
	)
}
