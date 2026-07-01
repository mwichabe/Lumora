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

	// The essential A1 grammar, writing & speaking — the foundation that turns
	// memorised phrases into sentences you can build yourself.
	seedSpanishA1Grammar(db)
}

// seedSpanishA1Grammar adds the grammar-focused A1 unit: ser/estar, articles &
// gender, the present tense (-ar/-er/-ir + key irregulars), gustar, numbers &
// time, daily routine (reflexives), directions, questions & negation, and the
// productive skills (writing an email, speaking an introduction).
func seedSpanishA1Grammar(db *gorm.DB) {
	const u = "A1 · Gramática Esencial"

	// ── Ser vs Estar ──
	s := addSkill(db, u, "Ser y Estar", "Two ways to say 'to be'.", "Sparkles", "#00C2A8", 9, 210)
	l := addLesson(db, s, "Ser o Estar", 1, 18,
		char("Professor Finch", "Both mean 'to be'. SER is who you are; ESTAR is how or where you are."),
		mc("Profession uses ser", "Yo ___ profesor.", "soy", "soy", "estoy", "eres", "está"),
		mc("Location uses estar", "Yo ___ en Madrid.", "estoy", "estoy", "soy", "es", "están"),
		tr("Translate this sentence", "She is tired", "Está cansada"),
		mc("Which verb for location?", "¿Dónde estás? →", "estar", "estar", "ser", "tener", "ir"),
		fill("Fill in the blank", "Nosotros ___ estudiantes. (ser)", "somos"),
		speak("Blaze", "Soy de España. Estoy en Madrid."),
	)
	addVocab(db, l,
		vw("ser", "to be (essence)", "Soy profesor.", "I am a teacher.", "Professor Finch"),
		vw("estar", "to be (state/place)", "Estoy bien.", "I am fine.", "Professor Finch"),
		vw("soy", "I am (ser)", "Soy de México.", "I'm from Mexico.", "Cora"),
		vw("estoy", "I am (estar)", "Estoy cansado.", "I'm tired.", "Cora"),
		vw("es", "he/she/it is", "Ella es alta.", "She is tall.", "Lumora"),
	)

	// ── Articles & gender ──
	s = addSkill(db, u, "Artículos y Género", "el / la / un / una.", "Hash", "#F5A623", 10, 226)
	l = addLesson(db, s, "El, La, Un, Una", 1, 18,
		char("Cora", "Every noun has a gender! 'el' for masculine, 'la' for feminine — and plurals too."),
		mc("the book (masculine)", "___ libro", "el", "el", "la", "los", "las"),
		mc("the house (feminine)", "___ casa", "la", "la", "el", "un", "unos"),
		tr("Translate: a friend (female)", "a female friend", "una amiga"),
		fill("Fill in the blank", "___ niños juegan. (the, masc. plural)", "los"),
		mc("indefinite plural feminine", "___ casas", "unas", "unas", "unos", "una", "un"),
		speak("Blaze", "el libro, la casa, un amigo, una amiga."),
	)
	addVocab(db, l,
		vw("el / la", "the (m / f)", "el libro, la mesa", "the book, the table", "Cora"),
		vw("un / una", "a (m / f)", "un café, una silla", "a coffee, a chair", "Cora"),
		vw("los / las", "the (plural)", "los niños, las niñas", "the boys, the girls", "Lumora"),
		vw("unos / unas", "some", "unos libros, unas casas", "some books, some houses", "Lumora"),
	)

	// ── Present tense -AR ──
	s = addSkill(db, u, "Presente: verbos -AR", "hablo, hablas, habla…", "PenLine", "#17A3DD", 11, 242)
	l = addLesson(db, s, "Verbos en -AR", 1, 18,
		char("Professor Finch", "Regular -ar verbs: drop -ar, add -o, -as, -a, -amos, -áis, -an."),
		mc("I speak", "yo ___", "hablo", "hablo", "hablas", "habla", "hablan"),
		mc("you speak (tú)", "tú ___", "hablas", "hablas", "hablo", "habla", "habláis"),
		fill("Fill in the blank", "Nosotros ___ español. (hablar)", "hablamos"),
		tr("Translate this sentence", "They work", "Ellos trabajan"),
		mc("she studies", "ella ___", "estudia", "estudia", "estudio", "estudias", "estudian"),
		speak("Blaze", "Yo hablo, tú hablas, él habla."),
	)
	addVocab(db, l,
		vw("hablar", "to speak", "Hablo español.", "I speak Spanish.", "Professor Finch"),
		vw("trabajar", "to work", "Trabajo en casa.", "I work at home.", "Cora"),
		vw("estudiar", "to study", "Estudio inglés.", "I study English.", "Cora"),
		vw("hablamos", "we speak", "Hablamos mucho.", "We talk a lot.", "Lumora"),
	)

	// ── Present tense -ER / -IR ──
	s = addSkill(db, u, "Presente: -ER / -IR", "como, vivo, escribo…", "PenLine", "#6C3FC5", 12, 258)
	l = addLesson(db, s, "Verbos en -ER / -IR", 1, 18,
		char("Professor Finch", "-er → como, comes, come. -ir → vivo, vives, vive."),
		mc("I eat", "yo ___", "como", "como", "comes", "come", "comen"),
		mc("you live (tú)", "tú ___", "vives", "vives", "vivo", "vive", "vivimos"),
		fill("Fill in the blank", "Nosotros ___ en Madrid. (vivir)", "vivimos"),
		tr("Translate this sentence", "He drinks water", "Él bebe agua"),
		mc("they write", "ellos ___", "escriben", "escriben", "escribo", "escribe", "escribimos"),
		speak("Blaze", "Como pan y bebo agua. Vivo en España."),
	)
	addVocab(db, l,
		vw("comer", "to eat", "Como una manzana.", "I eat an apple.", "Professor Finch"),
		vw("beber", "to drink", "Bebo agua.", "I drink water.", "Cora"),
		vw("vivir", "to live", "Vivo en Madrid.", "I live in Madrid.", "Cora"),
		vw("escribir", "to write", "Escribo una carta.", "I write a letter.", "Lumora"),
	)

	// ── Key irregular verbs ──
	s = addSkill(db, u, "Verbos Irregulares", "tener, ir, hacer, querer, poder.", "Sparkles", "#FF5C5C", 13, 276)
	l = addLesson(db, s, "Irregulares Clave", 1, 18,
		char("Professor Finch", "Learn these by heart: tengo, voy, hago, quiero, puedo."),
		mc("I have", "yo ___", "tengo", "tengo", "tienes", "tiene", "tienen"),
		mc("I go", "yo ___", "voy", "voy", "vas", "va", "van"),
		fill("Fill in the blank", "Yo ___ la tarea. (hacer)", "hago"),
		tr("Translate this sentence", "I want a coffee", "Quiero un café"),
		mc("I can", "yo ___", "puedo", "puedo", "puede", "podemos", "pueden"),
		speak("Blaze", "Tengo hambre. Voy a casa."),
	)
	addVocab(db, l,
		vw("tener", "to have", "Tengo un perro.", "I have a dog.", "Professor Finch"),
		vw("ir", "to go", "Voy al cine.", "I go to the cinema.", "Cora"),
		vw("hacer", "to do / make", "Hago deporte.", "I do sport.", "Cora"),
		vw("querer", "to want", "Quiero agua.", "I want water.", "Lumora"),
		vw("poder", "to be able to", "Puedo ayudar.", "I can help.", "Lumora"),
	)

	// ── Gustar ──
	s = addSkill(db, u, "Me gusta", "Saying what you like.", "Heart", "#FF5C5C", 14, 294)
	l = addLesson(db, s, "Gustar", 1, 18,
		char("Cora", "'Me gusta' + one thing, 'me gustan' + many things. Add no to say you don't!"),
		mc("I like coffee", "___ el café", "Me gusta", "Me gusta", "Me gustan", "Te gusta", "Le gusta"),
		mc("I like the books", "Me ___ los libros", "gustan", "gustan", "gusta", "gusto", "gustas"),
		tr("Translate this sentence", "Do you like music?", "¿Te gusta la música?"),
		fill("Fill in the blank", "A él ___ gusta el fútbol. (to him)", "le"),
		mc("I don't like tea", "No me ___ el té", "gusta", "gusta", "gustan", "gusto", "gustas"),
		speak("Blaze", "Me gusta el café. Me gustan los libros."),
	)
	addVocab(db, l,
		vw("me gusta", "I like (it)", "Me gusta bailar.", "I like to dance.", "Cora"),
		vw("me gustan", "I like (them)", "Me gustan los gatos.", "I like cats.", "Cora"),
		vw("te gusta", "you like", "¿Te gusta el cine?", "Do you like cinema?", "Lumora"),
		vw("le gusta", "he/she likes", "Le gusta leer.", "He likes to read.", "Lumora"),
	)

	// ── Numbers, time & date ──
	s = addSkill(db, u, "Números, Hora y Fecha", "0–100, telling time, days.", "Clock", "#F5A623", 15, 312)
	l = addLesson(db, s, "¿Qué hora es?", 1, 18,
		char("Lumora", "Numbers, the clock and the calendar — the everyday essentials."),
		mc("ten", "diez =", "ten", "ten", "eight", "twenty", "a hundred"),
		mc("it's two o'clock", "Son las ___", "dos", "dos", "tres", "doce", "diez"),
		tr("Translate this sentence", "It's half past three", "Son las tres y media"),
		fill("Fill in the blank", "Hoy es ___. (Monday)", "lunes"),
		mc("twenty", "veinte =", "twenty", "twenty", "twelve", "two", "a hundred"),
		speak("Blaze", "Son las dos y media. Hoy es lunes."),
	)
	addVocab(db, l,
		vw("diez", "ten", "Tengo diez euros.", "I have ten euros.", "Lumora"),
		vw("veinte", "twenty", "Veinte minutos.", "Twenty minutes.", "Cora"),
		vw("cien", "a hundred", "Cien personas.", "A hundred people.", "Cora"),
		vw("¿qué hora es?", "what time is it?", "¿Qué hora es?", "What time is it?", "Lumora"),
		vw("y media", "half past", "Son las tres y media.", "It's half past three.", "Lumora"),
		vw("lunes", "Monday", "Hoy es lunes.", "Today is Monday.", "Cora"),
	)

	// ── Daily routine (reflexives) ──
	s = addSkill(db, u, "Rutina Diaria", "Reflexive verbs: me levanto…", "Coffee", "#6C3FC5", 16, 330)
	l = addLesson(db, s, "Mi Día", 1, 18,
		char("Mira", "Routines use reflexive verbs: me levanto, me ducho, me acuesto."),
		mc("I get up", "___ a las siete.", "Me levanto", "Me levanto", "Me acuesto", "Me ducho", "Me llamo"),
		tr("Translate this sentence", "I go to bed late", "Me acuesto tarde"),
		fill("Fill in the blank", "Por la mañana ___ ducho. (myself)", "me"),
		match("What is 'dormir'?", "dormir", "to sleep", "to sleep", "to eat", "to wake", "to wash"),
		listen("Listen and choose the meaning", "Me levanto temprano.", "I get up early", "I get up early", "I go to bed late", "I eat breakfast", "I work a lot"),
		speak("Blaze", "Me levanto a las siete y me acuesto a las once."),
	)
	addVocab(db, l,
		vw("levantarse", "to get up", "Me levanto temprano.", "I get up early.", "Mira"),
		vw("ducharse", "to shower", "Me ducho por la mañana.", "I shower in the morning.", "Cora"),
		vw("acostarse", "to go to bed", "Me acuesto tarde.", "I go to bed late.", "Cora"),
		vw("desayunar", "to have breakfast", "Desayuno café y pan.", "I have coffee and bread.", "Lumora"),
		vw("dormir", "to sleep", "Duermo ocho horas.", "I sleep eight hours.", "Lumora"),
	)

	// ── City & directions ──
	s = addSkill(db, u, "Ciudad y Direcciones", "¿Dónde está…?, hay, a la derecha.", "Compass", "#17A3DD", 17, 350)
	l = addLesson(db, s, "¿Dónde está?", 1, 18,
		char("Riko", "Ask the way: ¿dónde está…? Then follow: a la derecha, a la izquierda, todo recto."),
		mc("Where is the bathroom?", "¿Dónde está el ___?", "baño", "baño", "menú", "vuelo", "amor"),
		mc("to the right", "a la ___", "derecha", "derecha", "izquierda", "recto", "mesa"),
		tr("Translate this sentence", "Turn left", "Gira a la izquierda"),
		fill("Fill in the blank", "___ un banco aquí. (there is)", "Hay"),
		match("What is 'todo recto'?", "todo recto", "straight ahead", "straight ahead", "to the right", "behind", "upstairs"),
		speak("Blaze", "¿Dónde está el hotel? Todo recto y a la derecha."),
	)
	addVocab(db, l,
		vw("¿dónde está?", "where is?", "¿Dónde está el baño?", "Where is the bathroom?", "Riko"),
		vw("a la derecha", "to the right", "Gira a la derecha.", "Turn right.", "Cora"),
		vw("a la izquierda", "to the left", "Está a la izquierda.", "It's on the left.", "Cora"),
		vw("todo recto", "straight ahead", "Sigue todo recto.", "Go straight ahead.", "Lumora"),
		vw("hay", "there is / are", "Hay un parque.", "There is a park.", "Lumora"),
	)

	// ── Questions & negation ──
	s = addSkill(db, u, "Preguntas y Negación", "qué, dónde, cómo… and saying no.", "MessageCircle", "#00C2A8", 18, 370)
	l = addLesson(db, s, "Preguntar y Negar", 1, 18,
		char("Cora", "Question words: qué, dónde, cómo, cuándo, quién, cuánto. To negate, put 'no' before the verb."),
		mc("what?", "¿___?", "Qué", "Qué", "Dónde", "Cómo", "Quién"),
		mc("where is it?", "¿___ está?", "Dónde", "Dónde", "Qué", "Cuándo", "Quién"),
		tr("Translate this sentence", "I don't speak English", "No hablo inglés"),
		fill("Fill in the blank", "¿___ años tienes? (how many)", "Cuántos"),
		mc("who is it?", "¿___ es?", "Quién", "Quién", "Qué", "Cómo", "Dónde"),
		speak("Blaze", "¿Dónde está? No, no hablo francés."),
	)
	addVocab(db, l,
		vw("qué", "what", "¿Qué es esto?", "What is this?", "Cora"),
		vw("dónde", "where", "¿Dónde vives?", "Where do you live?", "Cora"),
		vw("cómo", "how", "¿Cómo estás?", "How are you?", "Lumora"),
		vw("cuándo", "when", "¿Cuándo vienes?", "When do you come?", "Lumora"),
		vw("quién", "who", "¿Quién es?", "Who is it?", "Riko"),
		vw("cuánto", "how much", "¿Cuánto cuesta?", "How much is it?", "Riko"),
	)

	// ── Writing: email / postcard ──
	s = addSkill(db, u, "Escribir: Correo", "A1 writing — short messages & emails.", "PenLine", "#00C2A8", 19, 390)
	l = addLesson(db, s, "Escribe un Mensaje", 1, 20,
		char("Lumora", "A1 writing is short messages: greet, say a little, sign off. Let's practise."),
		mc("Start an informal message with…", "Greeting", "Hola", "Hola", "Estimado señor", "Atentamente", "Adiós"),
		write("Write a short postcard from Madrid (30–40 words): greet, say where you are and what you do.",
			"¡Hola! Estoy en Madrid. La ciudad es muy bonita. Por la mañana estudio español y por la tarde visito museos. Hace sol. Un abrazo, Ana."),
		write("Write a short email introducing yourself: your name, nationality, job and what you like.",
			"Hola, me llamo Pablo. Soy de México y soy profesor. Me gusta el fútbol y la música. ¿Y tú? Saludos, Pablo."),
		mc("Informal sign-off?", "Closing", "Un abrazo", "Un abrazo", "Atentamente", "Por favor", "Buenos días"),
		fill("Fill in the blank", "Un ___ (a hug — informal close)", "abrazo"),
		speak("Blaze", "Hola, me llamo Ana. Un abrazo."),
	)
	addVocab(db, l,
		vw("Hola / Querido", "Hi / Dear", "Querido amigo,", "Dear friend,", "Lumora"),
		vw("Un abrazo", "Hugs (informal close)", "Un abrazo, Ana.", "Hugs, Ana.", "Cora"),
		vw("Saludos", "Regards", "Saludos, Pablo.", "Regards, Pablo.", "Cora"),
		vw("Atentamente", "Sincerely (formal)", "Atentamente, Sr. Ruiz.", "Sincerely, Mr. Ruiz.", "Professor Finch"),
		vw("Estimado", "Dear (formal)", "Estimada señora:", "Dear madam:", "Professor Finch"),
	)

	// ── Speaking: introduce yourself ──
	s = addSkill(db, u, "Hablar: Preséntate", "Introduce yourself out loud.", "MessageCircle", "#6C3FC5", 20, 410)
	l = addLesson(db, s, "Preséntate", 1, 20,
		char("Blaze", "Time to SPEAK! Say each line out loud — I believe in you!"),
		speak("Lumora", "Hola, me llamo Ana."),
		speak("Blaze", "Soy de España y tengo veinte años."),
		speak("Blaze", "Soy estudiante. Me gusta la música y el café."),
		mc("To say your age", "Tengo ___ años.", "veinte", "veinte", "soy", "me llamo", "vivo"),
		speak("Blaze", "Vivo en Madrid y hablo español e inglés."),
		speak("Blaze", "Mucho gusto. ¿Y tú, cómo te llamas?"),
	)
	addVocab(db, l,
		vw("me llamo…", "my name is…", "Me llamo Ana.", "My name is Ana.", "Lumora"),
		vw("soy de…", "I'm from…", "Soy de España.", "I'm from Spain.", "Cora"),
		vw("tengo … años", "I'm … years old", "Tengo veinte años.", "I'm twenty.", "Cora"),
		vw("soy estudiante", "I'm a student", "Soy estudiante.", "I'm a student.", "Lumora"),
		vw("mucho gusto", "nice to meet you", "Mucho gusto.", "Nice to meet you.", "Riko"),
	)

	// A2 builds on A1: the past tenses, near future, object pronouns,
	// comparisons and the everyday topics that let you narrate and handle
	// routine situations.
	seedSpanishA2(db)
}

// seedSpanishA2 adds the elementary (A2) unit: preterite (regular + irregular),
// present perfect and when to use each, near future, present continuous, object
// pronouns, comparatives, gustar-type verbs, por/para, the affirmative
// imperative, health, clothing/travel, and the productive A2 skills.
func seedSpanishA2(db *gorm.DB) {
	const u = "A2 · Elemental"
	finch := "Professor Finch"

	// ── Preterite: regular ──
	s := addSkill(db, u, "Indefinido: Regulares", "The simple past: hablé, comí, viví.", "Clock", "#F5A623", 21, 430)
	l := addLesson(db, s, "El Pasado Regular", 1, 20,
		char(finch, "The preterite says what happened. Endings: -é/-aste/-ó and -í/-iste/-ió."),
		mc("I spoke", "yo ___", "hablé", "hablé", "hablo", "hablaré", "hablaba"),
		mc("you ate (tú)", "tú ___", "comiste", "comiste", "comes", "comió", "comías"),
		fill("Fill in the blank", "Ayer ___ en casa. (comer, yo)", "comí"),
		tr("Translate this sentence", "She lived in Madrid", "Vivió en Madrid"),
		mc("they worked", "ellos ___", "trabajaron", "trabajaron", "trabajan", "trabajaban", "trabajarán"),
		speak("Blaze", "Ayer hablé con mi amigo y comí en casa."),
	)
	addVocab(db, l,
		vw("hablé", "I spoke", "Ayer hablé con Ana.", "Yesterday I spoke with Ana.", finch),
		vw("comí", "I ate", "Comí una pizza.", "I ate a pizza.", "Cora"),
		vw("viví", "I lived", "Viví en México.", "I lived in Mexico.", "Cora"),
		vw("ayer", "yesterday", "Ayer trabajé.", "Yesterday I worked.", "Lumora"),
		vw("trabajó", "he/she worked", "Ella trabajó mucho.", "She worked a lot.", "Lumora"),
	)

	// ── Preterite: irregular ──
	s = addSkill(db, u, "Indefinido: Irregulares", "fui, tuve, hice, dije, estuve.", "Clock", "#6C3FC5", 22, 455)
	l = addLesson(db, s, "Pasados Irregulares", 1, 20,
		char(finch, "Learn these by heart: fui (ser/ir), tuve (tener), hice (hacer), dije (decir), estuve (estar)."),
		mc("I went / I was", "yo ___", "fui", "fui", "voy", "iba", "fue"),
		mc("I had", "yo ___", "tuve", "tuve", "tengo", "tenía", "tuvo"),
		fill("Fill in the blank", "Ayer ___ mucho trabajo. (hacer, yo)", "hice"),
		tr("Translate this sentence", "He said the truth", "Dijo la verdad"),
		mc("we were (estar)", "nosotros ___", "estuvimos", "estuvimos", "estamos", "éramos", "fuimos"),
		speak("Blaze", "Ayer fui al cine y tuve un buen día."),
	)
	addVocab(db, l,
		vw("fui", "I went / I was", "Fui al cine.", "I went to the cinema.", finch),
		vw("tuve", "I had", "Tuve una idea.", "I had an idea.", "Cora"),
		vw("hice", "I did / made", "Hice la tarea.", "I did the homework.", "Cora"),
		vw("dije", "I said", "Dije la verdad.", "I told the truth.", "Lumora"),
		vw("estuve", "I was (estar)", "Estuve en casa.", "I was at home.", "Lumora"),
	)

	// ── Present perfect ──
	s = addSkill(db, u, "Pretérito Perfecto", "he hablado — the recent past.", "Clock", "#17A3DD", 23, 480)
	l = addLesson(db, s, "He hablado", 1, 20,
		char(finch, "Recent past: haber (he/has/ha/hemos/habéis/han) + participle (hablado, comido, vivido)."),
		mc("I have spoken", "___ hablado", "He", "He", "Has", "Ha", "Hemos"),
		mc("have you eaten? (tú)", "¿___ comido?", "Has", "Has", "He", "Ha", "Han"),
		fill("Fill in the blank", "Hoy ___ trabajado mucho. (yo, haber)", "he"),
		tr("Translate this sentence", "We have lived here", "Hemos vivido aquí"),
		mc("irregular participle: written", "escribir →", "escrito", "escrito", "escribido", "escribo", "escrita"),
		speak("Blaze", "Hoy he comido bien y he estudiado español."),
	)
	addVocab(db, l,
		vw("he comido", "I have eaten", "Hoy he comido tarde.", "Today I've eaten late.", finch),
		vw("hoy", "today", "Hoy he trabajado.", "Today I have worked.", "Cora"),
		vw("ya", "already", "Ya he terminado.", "I've already finished.", "Cora"),
		vw("todavía no", "not yet", "Todavía no he comido.", "I haven't eaten yet.", "Lumora"),
		vw("hecho", "done/made (part.)", "He hecho la cama.", "I've made the bed.", "Lumora"),
	)

	// ── Indefinido vs Perfecto ──
	s = addSkill(db, u, "Indefinido vs Perfecto", "ayer vs hoy — which past?", "Layers", "#FF5C5C", 24, 508)
	l = addLesson(db, s, "¿Cuál Pasado?", 1, 20,
		char(finch, "Perfecto for today/this week (recent, relevant). Indefinido for finished time (ayer, en 2020)."),
		mc("Ayer ___ al cine. (ir)", "ayer = finished", "fui", "fui", "he ido", "voy", "iba"),
		mc("Hoy ___ al cine. (ir)", "hoy = recent", "he ido", "he ido", "fui", "iba", "voy"),
		fill("Fill in the blank", "El año pasado ___ a México. (viajar, yo)", "viajé"),
		tr("Translate this sentence", "This week I have worked a lot", "Esta semana he trabajado mucho"),
		mc("marker for indefinido", "Which signals finished past?", "ayer", "ayer", "hoy", "esta semana", "ya"),
		speak("Blaze", "Ayer comí pizza. Hoy he comido ensalada."),
	)
	addVocab(db, l,
		vw("ayer", "yesterday", "Ayer llovió.", "Yesterday it rained.", finch),
		vw("el año pasado", "last year", "El año pasado viajé.", "Last year I travelled.", "Cora"),
		vw("esta semana", "this week", "Esta semana he leído.", "This week I've read.", "Cora"),
		vw("hoy", "today", "Hoy he salido.", "Today I went out.", "Lumora"),
		vw("ya", "already", "Ya lo he visto.", "I've already seen it.", "Lumora"),
	)

	// ── Near future ──
	s = addSkill(db, u, "Futuro: ir a + infinitivo", "Plans: voy a estudiar.", "Compass", "#00C2A8", 25, 536)
	l = addLesson(db, s, "Voy a…", 1, 20,
		char("Lumora", "Talk about plans with ir a + infinitive: Voy a viajar."),
		mc("I'm going to study", "Voy ___ estudiar", "a", "a", "de", "en", "que"),
		mc("we are going to travel", "Vamos a ___", "viajar", "viajar", "viajo", "viajamos", "viajado"),
		fill("Fill in the blank", "Mañana ___ a comer fuera. (ir, nosotros)", "vamos"),
		tr("Translate this sentence", "She is going to work tomorrow", "Va a trabajar mañana"),
		mc("tomorrow", "mañana =", "tomorrow", "tomorrow", "yesterday", "today", "now"),
		speak("Blaze", "Mañana voy a estudiar y voy a hacer deporte."),
	)
	addVocab(db, l,
		vw("voy a", "I'm going to", "Voy a leer.", "I'm going to read.", "Lumora"),
		vw("vas a", "you're going to", "¿Vas a venir?", "Are you going to come?", "Cora"),
		vw("vamos a", "we're going to", "Vamos a cenar.", "We're going to have dinner.", "Cora"),
		vw("mañana", "tomorrow", "Mañana descanso.", "Tomorrow I rest.", "Lumora"),
		vw("el plan", "the plan", "¿Cuál es el plan?", "What's the plan?", "Riko"),
	)

	// ── Present continuous ──
	s = addSkill(db, u, "Estar + Gerundio", "Happening now: estoy comiendo.", "Sparkles", "#F5A623", 26, 564)
	l = addLesson(db, s, "Ahora Mismo", 1, 20,
		char("Cora", "Right now: estar + gerund (-ando / -iendo). Estoy comiendo."),
		mc("I am eating", "Estoy ___", "comiendo", "comiendo", "comer", "comido", "como"),
		mc("she is speaking", "Está ___", "hablando", "hablando", "habla", "hablar", "hablado"),
		fill("Fill in the blank", "Ahora ___ estudiando. (yo, estar)", "estoy"),
		tr("Translate this sentence", "They are working", "Están trabajando"),
		mc("gerund of vivir", "vivir →", "viviendo", "viviendo", "vivir", "vivido", "vivo"),
		speak("Blaze", "Estoy aprendiendo español ahora mismo."),
	)
	addVocab(db, l,
		vw("estoy + -ando", "I am …-ing", "Estoy cocinando.", "I'm cooking.", "Cora"),
		vw("hablando", "speaking", "Está hablando.", "He's speaking.", "Cora"),
		vw("comiendo", "eating", "Estoy comiendo.", "I'm eating.", "Lumora"),
		vw("viviendo", "living", "Estamos viviendo aquí.", "We're living here.", "Lumora"),
		vw("ahora", "now", "Ahora trabajo.", "Now I'm working.", "Riko"),
	)

	// ── Object pronouns ──
	s = addSkill(db, u, "Pronombres de Objeto", "lo / la / le / les.", "Hash", "#6C3FC5", 27, 594)
	l = addLesson(db, s, "Lo, La, Le, Les", 1, 20,
		char(finch, "Direct: lo/la/los/las (it/them). Indirect: le/les (to him/her/them). Placed before the verb."),
		mc("I see it (the book)", "___ veo", "Lo", "Lo", "La", "Le", "Les"),
		mc("I give him the book", "___ doy el libro", "Le", "Le", "Lo", "La", "Les"),
		fill("Fill in the blank", "¿La carta? ___ escribo. (it, fem.)", "La"),
		tr("Translate this sentence", "I buy them (the shoes)", "Los compro"),
		mc("to them", "___ hablo a ellos", "Les", "Les", "Lo", "La", "Le"),
		speak("Blaze", "¿El café? Lo quiero. ¿A María? Le hablo."),
	)
	addVocab(db, l,
		vw("lo / la", "it (m / f)", "Lo veo. La leo.", "I see it. I read it.", finch),
		vw("los / las", "them", "Los compro.", "I buy them.", "Cora"),
		vw("le", "(to) him/her", "Le doy el libro.", "I give him the book.", "Cora"),
		vw("les", "(to) them", "Les hablo.", "I talk to them.", "Lumora"),
		vw("me / te", "me / you", "Me llamas.", "You call me.", "Lumora"),
	)

	// ── Comparatives & superlatives ──
	s = addSkill(db, u, "Comparativos", "más/menos que, el más, mejor.", "Hash", "#17A3DD", 28, 624)
	l = addLesson(db, s, "Más y Menos", 1, 20,
		char(finch, "Compare with más/menos … que. The most: el/la más …. Irregulars: mejor, peor."),
		mc("taller than", "más alto ___ tú", "que", "que", "de", "como", "y"),
		mc("the best", "el ___", "mejor", "mejor", "más bueno", "bien", "bueno"),
		fill("Fill in the blank", "Madrid es ___ grande que mi pueblo. (more)", "más"),
		tr("Translate this sentence", "She is less tall than me", "Es menos alta que yo"),
		mc("as … as", "tan alto ___ tú", "como", "como", "que", "de", "más"),
		speak("Blaze", "Soy más alto que mi hermano, pero él juega mejor."),
	)
	addVocab(db, l,
		vw("más … que", "more … than", "Más alto que tú.", "Taller than you.", finch),
		vw("menos … que", "less … than", "Menos caro que eso.", "Less expensive than that.", "Cora"),
		vw("tan … como", "as … as", "Tan rápido como tú.", "As fast as you.", "Cora"),
		vw("el más", "the most", "El más grande.", "The biggest.", "Lumora"),
		vw("mejor / peor", "better / worse", "Es mejor.", "It's better.", "Lumora"),
	)

	// ── Gustar-type verbs ──
	s = addSkill(db, u, "Verbos como Gustar", "encantar, interesar, doler.", "Heart", "#FF5C5C", 29, 654)
	l = addLesson(db, s, "Me encanta", 1, 20,
		char("Cora", "Built like gustar: me encanta, me interesa, me duele. Use me/te/le + verb."),
		mc("I love dancing", "Me ___ bailar", "encanta", "encanta", "encantan", "encanto", "encantas"),
		mc("my head hurts", "Me ___ la cabeza", "duele", "duele", "duelen", "dolor", "duelo"),
		fill("Fill in the blank", "Me ___ los museos. (interesar, plural)", "interesan"),
		tr("Translate this sentence", "Her feet hurt", "Le duelen los pies"),
		mc("I'm interested in art", "Me interesa el ___", "arte", "arte", "agua", "amor", "aire"),
		speak("Blaze", "Me encanta el café, pero me duele la cabeza."),
	)
	addVocab(db, l,
		vw("encantar", "to love", "Me encanta leer.", "I love reading.", "Cora"),
		vw("interesar", "to interest", "Me interesa el arte.", "Art interests me.", "Cora"),
		vw("doler", "to hurt", "Me duele la espalda.", "My back hurts.", "Lumora"),
		vw("me encanta", "I love (it)", "Me encanta el mar.", "I love the sea.", "Lumora"),
		vw("me duele", "it hurts (me)", "Me duele aquí.", "It hurts here.", "Mira"),
	)

	// ── Por vs Para ──
	s = addSkill(db, u, "Por y Para", "reason vs purpose.", "Link2", "#00C2A8", 30, 686)
	l = addLesson(db, s, "Por o Para", 1, 20,
		char(finch, "POR = reason, exchange, through. PARA = purpose, destination, deadline."),
		mc("thanks for the gift", "Gracias ___ el regalo", "por", "por", "para", "de", "a"),
		mc("this is for you", "Esto es ___ ti", "para", "para", "por", "de", "a"),
		fill("Fill in the blank", "Estudio ___ aprender. (in order to)", "para"),
		tr("Translate this sentence", "I travel through Spain", "Viajo por España"),
		mc("for tomorrow (deadline)", "La tarea es ___ mañana", "para", "para", "por", "de", "en"),
		speak("Blaze", "Gracias por todo. Esto es para ti."),
	)
	addVocab(db, l,
		vw("por", "for / through / because of", "Gracias por venir.", "Thanks for coming.", finch),
		vw("para", "for / in order to", "Es para ti.", "It's for you.", "Cora"),
		vw("por favor", "please", "Un café, por favor.", "A coffee, please.", "Cora"),
		vw("para mí", "for me", "Para mí, un té.", "For me, a tea.", "Lumora"),
		vw("por eso", "that's why", "Por eso estudio.", "That's why I study.", "Lumora"),
	)

	// ── Affirmative imperative ──
	s = addSkill(db, u, "Imperativo Afirmativo", "¡Habla! ¡Ven! ¡Haz!", "MessageCircle", "#F5A623", 31, 718)
	l = addLesson(db, s, "¡Órdenes!", 1, 20,
		char("Blaze", "Give commands! tú: habla, come, vive. Irregulars: ven, haz, di, pon, ve."),
		mc("speak! (tú)", "¡___!", "Habla", "Habla", "Hablas", "Hablar", "Hablo"),
		mc("come! (tú)", "¡___ aquí!", "Ven", "Ven", "Viene", "Venir", "Vienes"),
		fill("Fill in the blank", "¡___ la tarea! (to do: haz)", "Haz"),
		tr("Translate this sentence", "Eat the fruit!", "¡Come la fruta!"),
		mc("tell me! (decir)", "¡___me la verdad!", "Di", "Di", "Dice", "Decir", "Digo"),
		speak("Blaze", "¡Habla más alto! ¡Ven aquí!"),
	)
	addVocab(db, l,
		vw("¡habla!", "speak!", "¡Habla despacio!", "Speak slowly!", "Blaze"),
		vw("¡come!", "eat!", "¡Come la sopa!", "Eat the soup!", "Cora"),
		vw("¡ven!", "come!", "¡Ven conmigo!", "Come with me!", "Cora"),
		vw("¡haz!", "do/make!", "¡Haz la cama!", "Make the bed!", "Lumora"),
		vw("¡di!", "say/tell!", "¡Dime!", "Tell me!", "Lumora"),
	)

	// ── Health & body ──
	s = addSkill(db, u, "Salud y el Cuerpo", "me duele…, body, symptoms.", "Heart", "#06AECE", 32, 750)
	l = addLesson(db, s, "En el Médico", 1, 20,
		char("Mira", "At the doctor: parts of the body and how to say what hurts."),
		mc("my head hurts", "Me duele la ___", "cabeza", "cabeza", "mano", "pie", "espalda"),
		mc("I'm ill", "Estoy ___", "enfermo", "enfermo", "bien", "contento", "alto"),
		fill("Fill in the blank", "Me duele el ___. (stomach)", "estómago"),
		tr("Translate this sentence", "I have a fever", "Tengo fiebre"),
		match("What is 'la mano'?", "la mano", "the hand", "the hand", "the foot", "the head", "the arm"),
		speak("Blaze", "Me duele la cabeza y tengo fiebre."),
	)
	addVocab(db, l,
		vw("la cabeza", "the head", "Me duele la cabeza.", "My head hurts.", "Mira"),
		vw("el estómago", "the stomach", "Me duele el estómago.", "My stomach hurts.", "Cora"),
		vw("la mano", "the hand", "Me lavo las manos.", "I wash my hands.", "Cora"),
		vw("la fiebre", "the fever", "Tengo fiebre.", "I have a fever.", "Lumora"),
		vw("enfermo", "ill", "Estoy enfermo.", "I'm ill.", "Lumora"),
	)

	// ── Clothes, shopping & travel ──
	s = addSkill(db, u, "Ropa, Compras y Viajes", "tallas, colores, transporte.", "ShoppingBag", "#6C3FC5", 33, 784)
	l = addLesson(db, s, "De Compras y de Viaje", 1, 20,
		char("Cora", "Clothes & sizes, plus getting around: la talla, el color, el tren, el billete."),
		mc("what size?", "¿Qué ___ tiene?", "talla", "talla", "color", "precio", "tienda"),
		mc("the train", "el ___", "tren", "tren", "avión", "coche", "barco"),
		fill("Fill in the blank", "Quiero un ___ de tren. (ticket)", "billete"),
		tr("Translate this sentence", "These shoes are too big", "Estos zapatos son muy grandes"),
		mc("can I try it on?", "¿Puedo ___?", "probármelo", "probármelo", "comerlo", "verlo", "pagarlo"),
		speak("Blaze", "¿Qué talla tiene? Quiero un billete de tren."),
	)
	addVocab(db, l,
		vw("la talla", "the size", "¿Qué talla usas?", "What size do you wear?", "Cora"),
		vw("el color", "the colour", "Me gusta el color azul.", "I like blue.", "Cora"),
		vw("el billete", "the ticket", "Un billete, por favor.", "A ticket, please.", "Lumora"),
		vw("el tren", "the train", "El tren llega tarde.", "The train is late.", "Lumora"),
		vw("los zapatos", "the shoes", "Estos zapatos son caros.", "These shoes are expensive.", "Riko"),
	)

	// ── Writing (A2) ──
	s = addSkill(db, u, "Escribir: Correos y Reseñas", "50–80 words with connectors.", "PenLine", "#00C2A8", 34, 818)
	l = addLesson(db, s, "Escribe con Conectores", 1, 22,
		char("Lumora", "A2 writing links ideas: primero, luego, después, porque. Tell a little story."),
		mc("first", "___, desayuno.", "Primero", "Primero", "Luego", "Pero", "Porque"),
		write("Write an email to a friend about your last weekend (50–80 words): what you did, using connectors.",
			"¡Hola! El fin de semana pasado fui a la playa con mis amigos. Primero nadamos y luego comimos en un restaurante. La comida estuvo deliciosa. Después fuimos al cine. ¡Lo pasé muy bien! Un abrazo, Ana."),
		write("Write a short restaurant review (50–80 words): what you ate and your opinion.",
			"Ayer cené en el restaurante La Plaza. Pedí paella y estaba muy rica. El camarero fue amable y el precio fue barato. Volveré pronto. ¡Lo recomiendo!"),
		mc("then / next", "Primero como y ___ estudio.", "luego", "luego", "pero", "porque", "primero"),
		fill("Fill in the blank", "El fin de semana ___ a la playa. (ir, yo)", "fui"),
		speak("Blaze", "El fin de semana pasado fui a la playa con mis amigos."),
	)
	addVocab(db, l,
		vw("primero", "first", "Primero estudio.", "First I study.", "Lumora"),
		vw("luego", "then", "Luego como.", "Then I eat.", "Cora"),
		vw("después", "afterwards", "Después salgo.", "Afterwards I go out.", "Cora"),
		vw("porque", "because", "No salgo porque llueve.", "I don't go out because it's raining.", "Lumora"),
		vw("el fin de semana pasado", "last weekend", "El fin de semana pasado descansé.", "Last weekend I rested.", "Lumora"),
	)

	// ── Speaking (A2) ──
	s = addSkill(db, u, "Hablar: Narra y Planea", "Tell a story, make plans aloud.", "MessageCircle", "#6C3FC5", 35, 852)
	l = addLesson(db, s, "Narra tu Fin de Semana", 1, 22,
		char("Blaze", "Tell what you did and what you'll do — out loud!"),
		speak("Lumora", "El sábado fui al parque con mi familia."),
		speak("Blaze", "El fin de semana pasado visité a mis abuelos y comimos juntos."),
		speak("Blaze", "Ayer estudié español y por la tarde vi una película."),
		mc("next weekend (future)", "El próximo fin de semana ___ a viajar.", "voy", "voy", "fui", "iba", "he ido"),
		speak("Blaze", "El próximo fin de semana voy a ir a la playa."),
		speak("Blaze", "Me encantó la película. ¿Y tú, qué hiciste?"),
	)
	addVocab(db, l,
		vw("fui", "I went", "Fui al parque.", "I went to the park.", "Lumora"),
		vw("visité", "I visited", "Visité a mi familia.", "I visited my family.", "Cora"),
		vw("el próximo fin de semana", "next weekend", "El próximo fin de semana viajo.", "Next weekend I travel.", "Cora"),
		vw("me encantó", "I loved it", "Me encantó la película.", "I loved the film.", "Lumora"),
		vw("¿qué hiciste?", "what did you do?", "¿Qué hiciste ayer?", "What did you do yesterday?", "Riko"),
	)

	// B1 makes you an independent user: the past-tense contrasts, future &
	// conditional, the present subjunctive, relative pronouns and connected,
	// opinionated discourse.
	seedSpanishB1(db)
}

// seedSpanishB1 adds the intermediate (B1) unit: imperfect & narrative past,
// pluperfect, simple future, conditional, the present subjunctive (forms +
// desire/emotion + doubt/opinion), relative pronouns, impersonal/passive se,
// discourse connectors, reported speech, prepositional verbs, and the
// productive B1 skills (opinion writing & sustained speaking).
func seedSpanishB1(db *gorm.DB) {
	const u = "B1 · Intermedio"
	finch := "Professor Finch"

	// ── Imperfect ──
	s := addSkill(db, u, "El Imperfecto", "Habits & background: hablaba, comía.", "Clock", "#F5A623", 36, 880)
	l := addLesson(db, s, "El Imperfecto", 1, 24,
		char(finch, "The imperfect paints the past: habits, descriptions, background. -aba (hablaba) and -ía (comía, vivía)."),
		mc("I used to speak", "yo ___", "hablaba", "hablaba", "hablé", "hablo", "hablaría"),
		mc("we used to eat", "nosotros ___", "comíamos", "comíamos", "comimos", "comemos", "comeríamos"),
		fill("Fill in the blank", "Cuando era niño, ___ en Madrid. (vivir, yo)", "vivía"),
		tr("Translate this sentence", "She was tall (description)", "Era alta"),
		mc("imperfect of ir (yo)", "yo ___ a la escuela", "iba", "iba", "fui", "voy", "iría"),
		speak("Blaze", "Cuando era niño, jugaba en el parque todos los días."),
	)
	addVocab(db, l,
		vw("hablaba", "I used to speak", "Hablaba mucho.", "I used to talk a lot.", finch),
		vw("comía", "I used to eat", "Comía en casa.", "I used to eat at home.", "Cora"),
		vw("vivía", "I used to live", "Vivía en un pueblo.", "I used to live in a village.", "Cora"),
		vw("era", "was (ser)", "Era muy alto.", "He was very tall.", "Lumora"),
		vw("iba", "used to go", "Iba al colegio a pie.", "I used to walk to school.", "Lumora"),
	)

	// ── Indefinido vs Imperfecto ──
	s = addSkill(db, u, "Indefinido vs Imperfecto", "Narrating: background vs action.", "Layers", "#6C3FC5", 37, 914)
	l = addLesson(db, s, "Narrar el Pasado", 1, 24,
		char(finch, "Imperfecto = the scene (ongoing). Indefinido = what happened (the event)."),
		mc("background: it was raining", "___ cuando salí", "Llovía", "Llovía", "Llovió", "Llueve", "Lloverá"),
		mc("event: suddenly it rained", "De repente ___", "llovió", "llovió", "llovía", "llueve", "lloverá"),
		fill("Fill in the blank", "Mientras ___, sonó el teléfono. (comer, yo)", "comía"),
		tr("Translate this sentence", "I was reading when he arrived", "Leía cuando llegó"),
		mc("the interrupting action", "Dormía cuando ___ el teléfono", "sonó", "sonó", "sonaba", "suena", "sonará"),
		speak("Blaze", "Hacía sol cuando llegué a la playa."),
	)
	addVocab(db, l,
		vw("mientras", "while", "Mientras comía, leía.", "While eating, I read.", finch),
		vw("de repente", "suddenly", "De repente, llegó.", "Suddenly, he arrived.", "Cora"),
		vw("cuando", "when", "Cuando llegué, dormías.", "When I arrived, you were sleeping.", "Cora"),
		vw("llovía", "it was raining", "Llovía mucho.", "It was raining a lot.", "Lumora"),
		vw("hacía sol", "it was sunny", "Hacía sol ese día.", "It was sunny that day.", "Lumora"),
	)

	// ── Pluperfect ──
	s = addSkill(db, u, "Pluscuamperfecto", "había + participle.", "Clock", "#17A3DD", 38, 948)
	l = addLesson(db, s, "Había hablado", 1, 24,
		char(finch, "The past before the past: había + participle. Cuando llegué, ya habían comido."),
		mc("I had spoken", "___ hablado", "Había", "Había", "He", "Habría", "Habrá"),
		mc("they had eaten", "___ comido", "habían", "habían", "han", "habrían", "habrán"),
		fill("Fill in the blank", "Cuando llegué, la película ya ___ empezado. (haber)", "había"),
		tr("Translate this sentence", "We had already left", "Ya habíamos salido"),
		mc("participle of ver", "ver →", "visto", "visto", "veído", "vido", "veo"),
		speak("Blaze", "Cuando llegué a casa, mis amigos ya se habían ido."),
	)
	addVocab(db, l,
		vw("había hablado", "I had spoken", "Ya había hablado con ella.", "I had already spoken to her.", finch),
		vw("ya", "already", "Ya había comido.", "I had already eaten.", "Cora"),
		vw("todavía no", "not yet", "Todavía no había salido.", "He hadn't left yet.", "Cora"),
		vw("habíamos", "we had", "Habíamos terminado.", "We had finished.", "Lumora"),
		vw("antes", "before", "Antes había vivido allí.", "Before, I had lived there.", "Lumora"),
	)

	// ── Simple future ──
	s = addSkill(db, u, "Futuro Simple", "hablaré, vendré, haré.", "Compass", "#00C2A8", 39, 982)
	l = addLesson(db, s, "El Futuro", 1, 24,
		char("Lumora", "Simple future: infinitive + é/ás/á/emos/éis/án. Irregular stems: vendr-, har-, tendr-."),
		mc("I will speak", "yo ___", "hablaré", "hablaré", "hablo", "hablaba", "hablaría"),
		mc("they will come", "ellos ___", "vendrán", "vendrán", "vienen", "venían", "vendrían"),
		fill("Fill in the blank", "Mañana ___ el trabajo. (terminar, yo)", "terminaré"),
		tr("Translate this sentence", "We will travel next year", "Viajaremos el próximo año"),
		mc("future of hacer (yo)", "yo ___", "haré", "haré", "hago", "hacía", "haría"),
		speak("Blaze", "El año que viene viajaré por Sudamérica."),
	)
	addVocab(db, l,
		vw("hablaré", "I will speak", "Hablaré con él.", "I'll speak with him.", "Lumora"),
		vw("vendré", "I will come", "Vendré pronto.", "I'll come soon.", "Cora"),
		vw("haré", "I will do/make", "Haré la cena.", "I'll make dinner.", "Cora"),
		vw("tendré", "I will have", "Tendré tiempo.", "I'll have time.", "Lumora"),
		vw("el año que viene", "next year", "El año que viene estudio.", "Next year I study.", "Lumora"),
	)

	// ── Conditional ──
	s = addSkill(db, u, "El Condicional", "Politeness & hypotheticals: hablaría.", "Sparkles", "#F5A623", 40, 1018)
	l = addLesson(db, s, "El Condicional", 1, 24,
		char(finch, "The conditional (would): infinitive + ía. hablaría. Great for politeness and hypotheticals."),
		mc("I would speak", "yo ___", "hablaría", "hablaría", "hablaré", "hablaba", "hablo"),
		mc("could you help me?", "¿___ ayudarme?", "Podrías", "Podrías", "Puedes", "Podías", "Pudiste"),
		fill("Fill in the blank", "Yo ___ un café, por favor. (querer, polite)", "querría"),
		tr("Translate this sentence", "I would like to travel", "Me gustaría viajar"),
		mc("conditional of hacer (yo)", "yo ___", "haría", "haría", "haré", "hago", "hice"),
		speak("Blaze", "Me gustaría visitar Argentina algún día."),
	)
	addVocab(db, l,
		vw("hablaría", "I would speak", "Hablaría con ella.", "I would speak to her.", finch),
		vw("me gustaría", "I would like", "Me gustaría ir.", "I would like to go.", "Cora"),
		vw("podría", "could", "¿Podría ayudar?", "Could you help?", "Cora"),
		vw("querría", "I would like (want)", "Querría un té.", "I'd like a tea.", "Lumora"),
		vw("tendría", "would have", "Tendría que estudiar.", "I would have to study.", "Lumora"),
	)

	// ── Present subjunctive: forms ──
	s = addSkill(db, u, "Subjuntivo: Formas", "hable, coma, viva, sea, vaya.", "Quote", "#FF5C5C", 41, 1054)
	l = addLesson(db, s, "Formas del Subjuntivo", 1, 24,
		char(finch, "From the yo-present, swap the ending: hablo→hable, como→coma, vivo→viva. Irregulars: sea, vaya, tenga, haga."),
		mc("that I speak (subj.)", "que yo ___", "hable", "hable", "hablo", "hablé", "hablaré"),
		mc("that you eat (tú, subj.)", "que tú ___", "comas", "comas", "comes", "comiste", "comías"),
		fill("Fill in the blank", "Espero que ___ bien. (estar, tú - subj.)", "estés"),
		tr("Translate this sentence", "that we be (ser)", "que seamos"),
		mc("subjunctive of tener (yo)", "que yo ___", "tenga", "tenga", "tengo", "tuve", "tendré"),
		speak("Blaze", "Quiero que hables conmigo."),
	)
	addVocab(db, l,
		vw("hable", "(that) I/he speak", "que yo hable", "that I speak", finch),
		vw("coma", "(that) I/he eat", "que él coma", "that he eats", "Cora"),
		vw("viva", "(that) I/he live", "que viva feliz", "that he live happily", "Cora"),
		vw("sea", "(that) be (ser)", "que sea verdad", "that it be true", "Lumora"),
		vw("vaya", "(that) go (ir)", "que vaya a casa", "that he go home", "Lumora"),
	)

	// ── Subjunctive: desire & emotion ──
	s = addSkill(db, u, "Subjuntivo: Deseo y Emoción", "Quiero que vengas.", "Heart", "#FF5C5C", 42, 1090)
	l = addLesson(db, s, "Quiero que…", 1, 24,
		char("Cora", "With wishes & emotions (two subjects), use the subjunctive: Quiero que vengas. Me alegra que estés aquí."),
		mc("I want you to come", "Quiero que ___", "vengas", "vengas", "vienes", "venir", "vendrás"),
		mc("I hope you're well", "Espero que ___ bien", "estés", "estés", "estás", "eres", "estarás"),
		fill("Fill in the blank", "Me alegra que ___ aquí. (estar, tú - subj.)", "estés"),
		tr("Translate this sentence", "I want him to study", "Quiero que estudie"),
		mc("I hope it's sunny", "Ojalá ___ sol mañana", "haga", "haga", "hace", "hizo", "hará"),
		speak("Blaze", "Espero que tengas un buen día."),
	)
	addVocab(db, l,
		vw("quiero que", "I want (that)", "Quiero que vengas.", "I want you to come.", "Cora"),
		vw("espero que", "I hope (that)", "Espero que estés bien.", "I hope you're well.", "Cora"),
		vw("me alegra que", "I'm glad (that)", "Me alegra que vengas.", "I'm glad you're coming.", "Lumora"),
		vw("ojalá", "hopefully", "Ojalá llueva.", "Hopefully it rains.", "Lumora"),
		vw("que vengas", "that you come", "Quiero que vengas ya.", "I want you to come now.", "Riko"),
	)

	// ── Subjunctive: doubt & opinion ──
	s = addSkill(db, u, "Subjuntivo: Duda y Opinión", "No creo que sea verdad.", "Quote", "#6C3FC5", 43, 1128)
	l = addLesson(db, s, "No creo que…", 1, 24,
		char(finch, "Doubt and impersonal opinion trigger the subjunctive: No creo que sea verdad. Es importante que estudies."),
		mc("I don't think it's true", "No creo que ___ verdad", "sea", "sea", "es", "era", "será"),
		mc("it's important that you study", "Es importante que ___", "estudies", "estudies", "estudias", "estudiar", "estudiaste"),
		fill("Fill in the blank", "Quizás ___ mañana. (venir, él - subj.)", "venga"),
		tr("Translate this sentence", "It's possible that it rains", "Es posible que llueva"),
		mc("certainty → indicative", "Creo que ___ verdad", "es", "es", "sea", "fuera", "haya sido"),
		speak("Blaze", "No creo que sea difícil, pero es importante que practiques."),
	)
	addVocab(db, l,
		vw("no creo que", "I don't think (that)", "No creo que venga.", "I don't think he'll come.", finch),
		vw("es importante que", "it's important that", "Es importante que estudies.", "It's important you study.", "Cora"),
		vw("es posible que", "it's possible that", "Es posible que llueva.", "It may rain.", "Cora"),
		vw("quizás", "maybe", "Quizás venga.", "Maybe he'll come.", "Lumora"),
		vw("dudo que", "I doubt that", "Dudo que sea cierto.", "I doubt it's true.", "Lumora"),
	)

	// ── Relative pronouns ──
	s = addSkill(db, u, "Pronombres Relativos", "que, quien, lo que, donde.", "Link2", "#00C2A8", 44, 1166)
	l = addLesson(db, s, "Que, Quien, Lo que", 1, 24,
		char(finch, "Relatives join clauses: que (that/which), quien (who, after prepositions), lo que (what), donde (where)."),
		mc("the book that I read", "el libro ___ leí", "que", "que", "quien", "lo que", "cuyo"),
		mc("I don't understand what you say", "No entiendo ___ dices", "lo que", "lo que", "que", "quien", "cual"),
		fill("Fill in the blank", "La ciudad ___ vivo es bonita. (where)", "donde"),
		tr("Translate this sentence", "the woman with whom I work", "la mujer con quien trabajo"),
		mc("the man who came", "el hombre ___ vino", "que", "que", "lo que", "cuyo", "donde"),
		speak("Blaze", "Esta es la casa donde vivo y el coche que compré."),
	)
	addVocab(db, l,
		vw("que", "that / which", "el libro que leí", "the book I read", finch),
		vw("quien", "who", "la persona con quien hablo", "the person I talk to", "Cora"),
		vw("lo que", "what", "lo que dices", "what you say", "Cora"),
		vw("donde", "where", "la casa donde vivo", "the house where I live", "Lumora"),
		vw("cuyo", "whose", "el autor cuyo libro leí", "the author whose book I read", "Lumora"),
	)

	// ── Impersonal / passive se ──
	s = addSkill(db, u, "Se Impersonal y Pasiva", "Se habla español. Se lo doy.", "Hash", "#17A3DD", 45, 1206)
	l = addLesson(db, s, "El Se", 1, 24,
		char(finch, "se makes things impersonal/passive: Se habla español. Se venden casas. And le+lo → se lo."),
		mc("Spanish is spoken here", "Se ___ español", "habla", "habla", "hablan", "hablado", "hablar"),
		mc("houses are sold", "Se ___ casas", "venden", "venden", "vende", "vendido", "vender"),
		fill("Fill in the blank", "Le doy el libro → ___ lo doy. (le+lo)", "Se"),
		tr("Translate this sentence", "How do you say this?", "¿Cómo se dice esto?"),
		mc("one lives well here", "Aquí se ___ bien", "vive", "vive", "viven", "vivir", "vivido"),
		speak("Blaze", "Aquí se habla español. Se lo doy a María."),
	)
	addVocab(db, l,
		vw("se habla", "is spoken", "Aquí se habla inglés.", "English is spoken here.", finch),
		vw("se vende", "is sold", "Se vende pan.", "Bread is sold.", "Cora"),
		vw("se dice", "is said / one says", "¿Cómo se dice?", "How do you say it?", "Cora"),
		vw("se lo doy", "I give it to him", "Se lo doy mañana.", "I'll give it to him tomorrow.", "Lumora"),
		vw("se puede", "one can", "Aquí se puede fumar.", "You can smoke here.", "Lumora"),
	)

	// ── Discourse connectors ──
	s = addSkill(db, u, "Conectores del Discurso", "sin embargo, por lo tanto, aunque.", "MessageCircle", "#F5A623", 46, 1246)
	l = addLesson(db, s, "Conectores", 1, 24,
		char("Cora", "Connect ideas: sin embargo (however), por lo tanto (therefore), aunque (although), además (besides)."),
		mc("however", "Quiero ir. ___, no puedo.", "Sin embargo", "Sin embargo", "Por lo tanto", "Porque", "Además"),
		mc("therefore", "Llueve, ___ no salgo.", "por lo tanto", "por lo tanto", "sin embargo", "aunque", "pero"),
		fill("Fill in the blank", "___ es caro, lo compro. (although)", "Aunque"),
		tr("Translate this sentence", "Besides, it's late", "Además, es tarde"),
		mc("in spite of", "___ la lluvia, salí.", "A pesar de", "A pesar de", "Gracias a", "Sin", "Para"),
		speak("Blaze", "Quiero ir; sin embargo, no tengo tiempo."),
	)
	addVocab(db, l,
		vw("sin embargo", "however", "Es caro; sin embargo, es bueno.", "It's expensive; however, it's good.", "Cora"),
		vw("por lo tanto", "therefore", "Llueve, por lo tanto me quedo.", "It's raining, therefore I stay.", "Cora"),
		vw("aunque", "although", "Aunque llueve, salgo.", "Although it rains, I go out.", "Lumora"),
		vw("además", "besides", "Además, es tarde.", "Besides, it's late.", "Lumora"),
		vw("a pesar de", "in spite of", "A pesar de todo, vine.", "In spite of everything, I came.", "Riko"),
	)

	// ── Reported speech ──
	s = addSkill(db, u, "Estilo Indirecto", "Dice que… Dijo que…", "Languages", "#6C3FC5", 47, 1288)
	l = addLesson(db, s, "Estilo Indirecto", 1, 24,
		char(finch, "Reporting speech: 'Estoy cansado' → Dice que está cansado. In the past: Dijo que estaba cansado."),
		mc("He says he's tired", "Dice que ___ cansado", "está", "está", "estás", "estoy", "estaba"),
		mc("She said she was coming", "Dijo que ___", "venía", "venía", "viene", "vendrá", "vino"),
		fill("Fill in the blank", "'Tengo hambre' → Dice que ___ hambre.", "tiene"),
		tr("Translate this sentence", "He says he can't", "Dice que no puede"),
		mc("'Vendré' → Dijo que ___", "future → conditional", "vendría", "vendría", "vendrá", "viene", "venía"),
		speak("Blaze", "María dice que está bien y que vendrá mañana."),
	)
	addVocab(db, l,
		vw("dice que", "he/she says (that)", "Dice que está bien.", "He says he's fine.", finch),
		vw("dijo que", "he/she said (that)", "Dijo que vendría.", "He said he'd come.", "Cora"),
		vw("preguntó si", "asked whether", "Preguntó si venías.", "He asked if you were coming.", "Cora"),
		vw("que estaba", "that he was", "Dijo que estaba cansado.", "He said he was tired.", "Lumora"),
		vw("que vendría", "that he'd come", "Dijo que vendría.", "He said he'd come.", "Lumora"),
	)

	// ── Prepositional verbs ──
	s = addSkill(db, u, "Verbos con Preposición", "pensar en, depender de, soñar con.", "Link2", "#06AECE", 48, 1330)
	l = addLesson(db, s, "Verbos + Preposición", 1, 24,
		char(finch, "Many verbs take a fixed preposition: pensar EN, depender DE, soñar CON, ayudar A."),
		mc("I think about you", "Pienso ___ ti", "en", "en", "de", "con", "a"),
		mc("it depends on the weather", "Depende ___ tiempo", "del", "del", "en el", "con el", "al"),
		fill("Fill in the blank", "Sueño ___ viajar. (with)", "con"),
		tr("Translate this sentence", "I help you to study", "Te ayudo a estudiar"),
		mc("I remember (acordarse de)", "Me acuerdo ___ ti", "de", "de", "en", "con", "a"),
		speak("Blaze", "Pienso en mi familia y sueño con viajar."),
	)
	addVocab(db, l,
		vw("pensar en", "to think about", "Pienso en ti.", "I think about you.", finch),
		vw("depender de", "to depend on", "Depende de ti.", "It depends on you.", "Cora"),
		vw("soñar con", "to dream of", "Sueño con viajar.", "I dream of travelling.", "Cora"),
		vw("ayudar a", "to help to", "Te ayudo a estudiar.", "I help you study.", "Lumora"),
		vw("acordarse de", "to remember", "Me acuerdo de ti.", "I remember you.", "Lumora"),
	)

	// ── Writing (B1) ──
	s = addSkill(db, u, "Escribir: Opinión y Carta", "150–200 words, argued.", "PenLine", "#00C2A8", 49, 1374)
	l = addLesson(db, s, "Escribe tu Opinión", 1, 26,
		char("Lumora", "B1 writing has structure: introduction, body, conclusion — and connectors to argue your point."),
		mc("In my opinion…", "___, la tecnología ayuda.", "En mi opinión", "En mi opinión", "Sin embargo", "Por lo tanto", "Además"),
		write("Write an opinion text (150–200 words): Are social media good or bad? Give reasons and examples.",
			"En mi opinión, las redes sociales tienen ventajas y desventajas. Por un lado, nos permiten comunicarnos con amigos lejanos y compartir información rápidamente. Sin embargo, también pueden crear adicción y problemas de privacidad. Por ejemplo, muchos jóvenes pasan demasiado tiempo frente a la pantalla. Creo que es importante usarlas con moderación. En conclusión, las redes sociales son útiles si las usamos de forma responsable."),
		write("Write a formal email of complaint about a hotel (around 150 words).",
			"Estimado señor: Le escribo para quejarme de mi estancia en su hotel la semana pasada. La habitación estaba sucia y el aire acondicionado no funcionaba. Además, el servicio fue muy lento. Por estos motivos, solicito una compensación. Espero su pronta respuesta. Atentamente, Pablo Ruiz."),
		mc("to conclude", "___, las redes son útiles.", "En conclusión", "En conclusión", "Por un lado", "Sin embargo", "Primero"),
		fill("Fill in the blank", "Por un ___, es útil; por otro, es peligroso. (side)", "lado"),
		speak("Blaze", "En mi opinión, es importante leer cada día."),
	)
	addVocab(db, l,
		vw("en mi opinión", "in my opinion", "En mi opinión, sí.", "In my opinion, yes.", "Lumora"),
		vw("por un lado", "on one hand", "Por un lado, es útil.", "On one hand, it's useful.", "Cora"),
		vw("por otro lado", "on the other hand", "Por otro lado, es caro.", "On the other hand, it's costly.", "Cora"),
		vw("en conclusión", "in conclusion", "En conclusión, sí.", "In conclusion, yes.", "Lumora"),
		vw("creo que", "I think that", "Creo que es verdad.", "I think it's true.", "Riko"),
	)

	// ── Speaking (B1) ──
	s = addSkill(db, u, "Hablar: Opina y Debate", "Narrate, give & defend opinions.", "MessageCircle", "#6C3FC5", 50, 1420)
	l = addLesson(db, s, "Opina en Voz Alta", 1, 26,
		char("Blaze", "Now express yourself: narrate, give opinions, agree and disagree — out loud!"),
		speak("Lumora", "Cuando era niño, vivía en un pueblo pequeño y muy tranquilo."),
		speak("Blaze", "En mi opinión, viajar es la mejor forma de aprender."),
		speak("Blaze", "No creo que sea fácil, pero quiero que lo intentes."),
		mc("I agree", "___ contigo.", "Estoy de acuerdo", "Estoy de acuerdo", "No estoy de acuerdo", "Depende", "Quizás"),
		speak("Blaze", "Por un lado estoy de acuerdo, pero por otro lado tengo dudas."),
		speak("Blaze", "Me gustaría viajar más, aunque ahora no tengo tiempo."),
	)
	addVocab(db, l,
		vw("estoy de acuerdo", "I agree", "Estoy de acuerdo contigo.", "I agree with you.", "Blaze"),
		vw("no estoy de acuerdo", "I disagree", "No estoy de acuerdo.", "I disagree.", "Cora"),
		vw("en mi opinión", "in my opinion", "En mi opinión, no.", "In my opinion, no.", "Cora"),
		vw("cuando era niño", "when I was a child", "Cuando era niño, jugaba.", "When I was a child, I played.", "Lumora"),
		vw("me gustaría", "I would like", "Me gustaría ir.", "I'd like to go.", "Lumora"),
	)

	// B2: fluency, nuance and range — the imperfect subjunctive, every kind of
	// conditional, adverbial subjunctive, advanced periphrasis, the passive,
	// idioms and argued, register-aware discourse.
	seedSpanishB2(db)
}

// seedSpanishB2 adds the upper-intermediate (B2) unit: imperfect subjunctive,
// the three si-clause types, adverbial subjunctive, subjunctive vs indicative,
// future/conditional perfect, verbal periphrasis, passive & causatives,
// advanced connectors, full reported speech, idioms, and the productive B2
// skills (essay writing & debate).
func seedSpanishB2(db *gorm.DB) {
	const u = "B2 · Avanzado"
	finch := "Professor Finch"

	// ── Imperfect subjunctive: forms ──
	s := addSkill(db, u, "Imperfecto de Subjuntivo", "hablara / hablase, tuviera, fuera.", "Quote", "#FF5C5C", 51, 1460)
	l := addLesson(db, s, "Formas del Pasado Subjuntivo", 1, 28,
		char(finch, "Take the 3rd-person plural preterite, drop -ron, add -ra: hablaron→hablara, tuvieron→tuviera, fueron→fuera."),
		mc("that I spoke (subj.)", "que yo ___", "hablara", "hablara", "hablaba", "hablé", "hable"),
		mc("that you had (subj.)", "que tú ___", "tuvieras", "tuvieras", "tienes", "tuviste", "tengas"),
		fill("Fill in the blank", "Quería que ___ a casa. (venir, tú - subj.)", "vinieras"),
		tr("Translate this sentence", "as if it were true", "como si fuera verdad"),
		mc("imperfect subj. of hacer (yo)", "que yo ___", "hiciera", "hiciera", "hago", "hice", "haga"),
		speak("Blaze", "Ojalá tuviera más tiempo para viajar."),
	)
	addVocab(db, l,
		vw("hablara", "(that) I/he spoke", "que yo hablara", "that I spoke", finch),
		vw("tuviera", "(that) I/he had", "si tuviera dinero", "if I had money", "Cora"),
		vw("fuera", "(that) were (ser/ir)", "como si fuera fácil", "as if it were easy", "Cora"),
		vw("viniera", "(that) came", "quería que viniera", "I wanted him to come", "Lumora"),
		vw("hiciera", "(that) did/made", "si hiciera sol", "if it were sunny", "Lumora"),
	)

	// ── Si clauses: unreal present ──
	s = addSkill(db, u, "Condicional: Si + Imperfecto Subj.", "Si tuviera…, iría.", "Layers", "#6C3FC5", 52, 1510)
	l = addLesson(db, s, "Hipótesis Improbable", 1, 28,
		char(finch, "Unreal present: Si + imperfect subjunctive + conditional. Si tuviera tiempo, viajaría."),
		mc("If I had money…", "Si ___ dinero, viajaría.", "tuviera", "tuviera", "tengo", "tendría", "tenía"),
		mc("…I would travel", "Si tuviera dinero, ___.", "viajaría", "viajaría", "viajo", "viajara", "viajaré"),
		fill("Fill in the blank", "Si ___ tú, estudiaría más. (ser)", "fueras"),
		tr("Translate this sentence", "If I could, I would help you", "Si pudiera, te ayudaría"),
		mc("If it rained, we'd stay", "Si ___, nos quedaríamos.", "lloviera", "lloviera", "llueve", "llovería", "llovía"),
		speak("Blaze", "Si tuviera más tiempo, aprendería a tocar la guitarra."),
	)
	addVocab(db, l,
		vw("si tuviera", "if I had", "Si tuviera tiempo…", "If I had time…", finch),
		vw("viajaría", "I would travel", "Viajaría a Japón.", "I'd travel to Japan.", "Cora"),
		vw("si pudiera", "if I could", "Si pudiera, iría.", "If I could, I'd go.", "Cora"),
		vw("si fuera", "if I were", "Si fuera rico…", "If I were rich…", "Lumora"),
		vw("ayudaría", "I would help", "Te ayudaría.", "I would help you.", "Lumora"),
	)

	// ── Si clauses: unreal past ──
	s = addSkill(db, u, "Condicional Perfecto: Si hubiera…", "Si hubiera sabido, habría ido.", "Layers", "#FF5C5C", 53, 1560)
	l = addLesson(db, s, "Hipótesis Imposible (pasado)", 1, 28,
		char(finch, "Unreal past: Si + pluperfect subjunctive (hubiera + participle) + conditional perfect (habría + participle)."),
		mc("If I had known…", "Si ___ sabido, habría ido.", "hubiera", "hubiera", "había", "habré", "habría"),
		mc("…I would have gone", "Si hubiera sabido, ___ ido.", "habría", "habría", "hubiera", "había", "he"),
		fill("Fill in the blank", "Si me lo ___ dicho, te habría ayudado. (haber, tú - subj.)", "hubieras"),
		tr("Translate this sentence", "If you had studied, you would have passed", "Si hubieras estudiado, habrías aprobado"),
		mc("conditional perfect of ir (yo)", "yo ___ ido", "habría", "habría", "hubiera", "había", "habré"),
		speak("Blaze", "Si hubiera salido antes, no habría perdido el tren."),
	)
	addVocab(db, l,
		vw("si hubiera sabido", "if I had known", "Si hubiera sabido, vengo.", "Had I known, I'd come.", finch),
		vw("habría ido", "I would have gone", "Habría ido contigo.", "I'd have gone with you.", "Cora"),
		vw("hubieras", "you had (subj.)", "Si hubieras venido…", "If you had come…", "Cora"),
		vw("habrías", "you would have", "Habrías ganado.", "You would have won.", "Lumora"),
		vw("habríamos", "we would have", "Habríamos llegado.", "We'd have arrived.", "Lumora"),
	)

	// ── Adverbial subjunctive ──
	s = addSkill(db, u, "Subjuntivo en Adverbiales", "para que, sin que, antes de que.", "Quote", "#00C2A8", 54, 1610)
	l = addLesson(db, s, "Conjunciones + Subjuntivo", 1, 28,
		char(finch, "Some conjunctions always need the subjunctive: para que, sin que, antes de que, a menos que."),
		mc("so that you learn", "Te lo explico para que ___", "aprendas", "aprendas", "aprendes", "aprender", "aprendiste"),
		mc("without him noticing", "Salí sin que él ___ cuenta", "se diera", "se diera", "se da", "se dio", "se daba"),
		fill("Fill in the blank", "Llámame antes de que ___. (salir, tú - subj.)", "salgas"),
		tr("Translate this sentence", "unless it rains", "a menos que llueva"),
		mc("future time → subjunctive", "Cuando ___, te llamo.", "llegue", "llegue", "llego", "llegué", "llegaré"),
		speak("Blaze", "Trabajo mucho para que mi familia esté bien."),
	)
	addVocab(db, l,
		vw("para que", "so that", "Estudio para que apruebes.", "I study so you pass.", finch),
		vw("sin que", "without (subject)", "Entró sin que lo viéramos.", "He came in without us seeing.", "Cora"),
		vw("antes de que", "before", "Antes de que llegues…", "Before you arrive…", "Cora"),
		vw("a menos que", "unless", "No iré a menos que vengas.", "I won't go unless you come.", "Lumora"),
		vw("cuando + subj.", "when (future)", "Cuando llegues, llámame.", "When you arrive, call me.", "Lumora"),
	)

	// ── Subjunctive vs indicative ──
	s = addSkill(db, u, "Subjuntivo vs Indicativo", "aunque, el hecho de que, cuando.", "Quote", "#6C3FC5", 55, 1660)
	l = addLesson(db, s, "¿Subjuntivo o Indicativo?", 1, 28,
		char(finch, "Same word, two moods: 'aunque llueve' (fact) vs 'aunque llueva' (possibility). Indicative = real; subjunctive = unknown/valued."),
		mc("although it's raining (fact)", "Aunque ___, salgo.", "llueve", "llueve", "llueva", "lloviera", "lloverá"),
		mc("even if it rains (hypothetical)", "Aunque ___, saldré.", "llueva", "llueva", "llueve", "llovía", "llovió"),
		fill("Fill in the blank", "Cuando ___ pequeño, jugaba. (ser - habitual past = indicative)", "era"),
		tr("Translate this sentence", "I don't deny that it's difficult", "No niego que sea difícil"),
		mc("the fact that you came (valued)", "El hecho de que ___ me alegra.", "vinieras", "vinieras", "viniste", "venías", "vendrás"),
		speak("Blaze", "Aunque sea difícil, seguiré intentándolo."),
	)
	addVocab(db, l,
		vw("aunque (+ind.)", "although (fact)", "Aunque llueve, salgo.", "Although it's raining, I go.", finch),
		vw("aunque (+subj.)", "even if (maybe)", "Aunque llueva, iré.", "Even if it rains, I'll go.", "Cora"),
		vw("el hecho de que", "the fact that", "El hecho de que vengas…", "The fact that you come…", "Cora"),
		vw("no niego que", "I don't deny that", "No niego que sea duro.", "I don't deny it's hard.", "Lumora"),
		vw("siempre que", "as long as", "Iré siempre que vengas.", "I'll go as long as you come.", "Lumora"),
	)

	// ── Future / conditional perfect ──
	s = addSkill(db, u, "Futuro y Condicional Perfecto", "habré / habría + participio.", "Clock", "#17A3DD", 56, 1710)
	l = addLesson(db, s, "Acciones Acabadas", 1, 28,
		char(finch, "Future perfect (habré hablado) = will have done. Conditional perfect (habría hablado) = would have done. Also for guesses about the past."),
		mc("I will have finished", "___ terminado", "Habré", "Habré", "Había", "Habría", "He"),
		mc("they would have left", "___ salido", "Habrían", "Habrían", "Habrán", "Habían", "Han"),
		fill("Fill in the blank", "Para mañana ya ___ llegado. (haber, ellos - future perf.)", "habrán"),
		tr("Translate this sentence", "He must have arrived already (guess)", "Habrá llegado ya"),
		mc("guess about the past", "¿Quién llamó? — ___ sido Ana.", "Habrá", "Habrá", "Había", "Hubo", "Ha"),
		speak("Blaze", "Para el verano habré terminado el curso."),
	)
	addVocab(db, l,
		vw("habré + part.", "I will have …", "Habré comido.", "I'll have eaten.", finch),
		vw("habría + part.", "I would have …", "Habría ganado.", "I'd have won.", "Cora"),
		vw("habrá", "must have (guess)", "Habrá salido.", "He must have left.", "Cora"),
		vw("para entonces", "by then", "Para entonces habré vuelto.", "By then I'll have returned.", "Lumora"),
		vw("ya", "already", "Ya habrán llegado.", "They'll have arrived already.", "Lumora"),
	)

	// ── Verbal periphrasis ──
	s = addSkill(db, u, "Perífrasis Verbales", "llevar + gerundio, acabar de, volver a.", "Link2", "#F5A623", 57, 1760)
	l = addLesson(db, s, "Perífrasis", 1, 28,
		char(finch, "Useful constructions: llevar + gerund (duration), acabar de (just did), volver a (again), estar a punto de (about to), dejar de (stop)."),
		mc("I've been studying for an hour", "___ una hora estudiando", "Llevo", "Llevo", "Hago", "Tengo", "Estoy"),
		mc("I have just eaten", "___ de comer", "Acabo", "Acabo", "Vuelvo", "Dejo", "Llevo"),
		fill("Fill in the blank", "___ a empezar de nuevo. (start again: volver)", "Vuelvo"),
		tr("Translate this sentence", "I'm about to leave", "Estoy a punto de salir"),
		mc("I stopped smoking", "___ de fumar", "Dejé", "Dejé", "Volví", "Acabé", "Llevé"),
		speak("Blaze", "Llevo dos años estudiando español y acabo de aprobar un examen."),
	)
	addVocab(db, l,
		vw("llevar + gerundio", "to have been …-ing", "Llevo un año aquí.", "I've been here a year.", finch),
		vw("acabar de", "to have just", "Acabo de llegar.", "I've just arrived.", "Cora"),
		vw("volver a", "to do again", "Vuelvo a intentarlo.", "I try again.", "Cora"),
		vw("estar a punto de", "to be about to", "Está a punto de llover.", "It's about to rain.", "Lumora"),
		vw("dejar de", "to stop …-ing", "Dejé de fumar.", "I stopped smoking.", "Lumora"),
	)

	// ── Passive & causatives ──
	s = addSkill(db, u, "Voz Pasiva y Causativas", "ser + participio, hacer que.", "Hash", "#00C2A8", 58, 1810)
	l = addLesson(db, s, "La Pasiva", 1, 28,
		char(finch, "Passive with ser + participle (agreeing): La carta fue escrita por Ana. Causative: hacer que + subjunctive."),
		mc("the house was built", "La casa ___ construida", "fue", "fue", "fui", "es", "está"),
		mc("the letters were written", "Las cartas fueron ___", "escritas", "escritas", "escrito", "escritos", "escribir"),
		fill("Fill in the blank", "El libro fue escrito ___ Cervantes. (by)", "por"),
		tr("Translate this sentence", "He makes me study", "Hace que estudie"),
		mc("passive 'se': Spanish is spoken", "___ habla español", "Se", "Se", "Es", "Está", "Le"),
		speak("Blaze", "El cuadro fue pintado por un artista famoso."),
	)
	addVocab(db, l,
		vw("ser + participio", "to be (passive)", "Fue construido en 1900.", "It was built in 1900.", finch),
		vw("fue escrito", "was written", "Fue escrito por ella.", "It was written by her.", "Cora"),
		vw("por", "by (agent)", "Hecho por expertos.", "Made by experts.", "Cora"),
		vw("hacer que", "to make (someone)", "Hace que trabajemos.", "He makes us work.", "Lumora"),
		vw("la pasiva refleja", "passive 'se'", "Se venden pisos.", "Flats are sold.", "Lumora"),
	)

	// ── Advanced connectors ──
	s = addSkill(db, u, "Conectores Avanzados", "de hecho, en cambio, por consiguiente.", "MessageCircle", "#6C3FC5", 59, 1860)
	l = addLesson(db, s, "Marcadores del Discurso", 1, 28,
		char("Cora", "Sound natural and precise: de hecho (in fact), en cambio (on the other hand), por consiguiente (consequently), cabe destacar (it's worth noting)."),
		mc("in fact", "Es caro; ___, es el mejor.", "de hecho", "de hecho", "en cambio", "sin embargo", "además"),
		mc("on the other hand (contrast)", "Él es alto; ella, ___, es baja.", "en cambio", "en cambio", "de hecho", "por tanto", "además"),
		fill("Fill in the blank", "Llueve; por ___, no salimos. (consequently: consiguiente)", "consiguiente"),
		tr("Translate this sentence", "It's worth noting that it's free", "Cabe destacar que es gratis"),
		mc("to sum up", "___, fue un éxito.", "En resumen", "En resumen", "De hecho", "En cambio", "Es decir"),
		speak("Blaze", "El proyecto es difícil; no obstante, cabe destacar su importancia."),
	)
	addVocab(db, l,
		vw("de hecho", "in fact", "De hecho, ya lo sabía.", "In fact, I already knew.", "Cora"),
		vw("en cambio", "on the other hand", "Yo sí; él, en cambio, no.", "I do; he, however, doesn't.", "Cora"),
		vw("por consiguiente", "consequently", "Por consiguiente, fallé.", "Consequently, I failed.", "Lumora"),
		vw("cabe destacar", "it's worth noting", "Cabe destacar su esfuerzo.", "Worth noting his effort.", "Lumora"),
		vw("no obstante", "nevertheless", "No obstante, lo intenté.", "Nevertheless, I tried.", "Riko"),
	)

	// ── Full reported speech ──
	s = addSkill(db, u, "Estilo Indirecto Avanzado", "Tense back-shift in full.", "Languages", "#17A3DD", 60, 1910)
	l = addLesson(db, s, "Transformar el Discurso", 1, 28,
		char(finch, "When reporting in the past, tenses shift back: presente→imperfecto, indefinido→pluscuamperfecto, futuro→condicional, imperative→subjunctive."),
		mc("'Vivo aquí' → Dijo que ___ allí", "present → imperfect", "vivía", "vivía", "vive", "vivió", "viviría"),
		mc("'Llegué tarde' → Dijo que ___ tarde", "indefinido → pluperfect", "había llegado", "había llegado", "llegó", "llega", "llegaría"),
		fill("Fill in the blank", "'Ven' → Me pidió que ___. (venir - subj.)", "viniera"),
		tr("Translate this sentence", "She said she would come", "Dijo que vendría"),
		mc("'¿Estás bien?' → Preguntó si ___ bien", "question back-shift", "estaba", "estaba", "estás", "estés", "estarás"),
		speak("Blaze", "Me dijo que había estado enfermo y que vendría al día siguiente."),
	)
	addVocab(db, l,
		vw("dijo que vivía", "said he lived", "Dijo que vivía solo.", "He said he lived alone.", finch),
		vw("había llegado", "had arrived", "Dijo que había llegado.", "He said he had arrived.", "Cora"),
		vw("me pidió que", "asked me to", "Me pidió que esperara.", "He asked me to wait.", "Cora"),
		vw("al día siguiente", "the next day", "Vendría al día siguiente.", "He'd come the next day.", "Lumora"),
		vw("aquel día", "that day", "Dijo que aquel día llovía.", "He said that day it rained.", "Lumora"),
	)

	// ── Idioms ──
	s = addSkill(db, u, "Expresiones Idiomáticas", "echar de menos, dar igual, tener ganas.", "Sparkles", "#FF5C5C", 61, 1960)
	l = addLesson(db, s, "Modismos Útiles", 1, 28,
		char("Cora", "Sound like a native: echar de menos (to miss), dar igual (not to mind), tener ganas de (to look forward to), valer la pena (to be worth it)."),
		mc("I miss my family", "___ de menos a mi familia", "Echo", "Echo", "Tengo", "Doy", "Hago"),
		mc("I don't mind", "Me ___ igual", "da", "da", "tiene", "echa", "hace"),
		fill("Fill in the blank", "Tengo ___ de verte. (look forward to: ganas)", "ganas"),
		tr("Translate this sentence", "It's worth it", "Vale la pena"),
		mc("to pull someone's leg", "tomar el ___", "pelo", "pelo", "café", "sol", "tiempo"),
		speak("Blaze", "Echo de menos a mis amigos, pero valió la pena venir."),
	)
	addVocab(db, l,
		vw("echar de menos", "to miss", "Te echo de menos.", "I miss you.", "Cora"),
		vw("dar igual", "not to mind", "Me da igual.", "I don't mind.", "Cora"),
		vw("tener ganas de", "to look forward to", "Tengo ganas de ir.", "I look forward to going.", "Lumora"),
		vw("valer la pena", "to be worth it", "Vale la pena.", "It's worth it.", "Lumora"),
		vw("tomar el pelo", "to pull one's leg", "Me tomas el pelo.", "You're teasing me.", "Riko"),
	)

	// ── Society & current affairs vocabulary ──
	s = addSkill(db, u, "Sociedad y Actualidad", "Environment, politics, technology.", "BookOpen", "#06AECE", 62, 2010)
	l = addLesson(db, s, "Temas de Actualidad", 1, 28,
		char(finch, "B2 deals with abstract themes: el medio ambiente, la desigualdad, la globalización, el desarrollo."),
		mc("climate change", "el cambio ___", "climático", "climático", "ambiente", "global", "natural"),
		mc("unemployment", "el ___", "desempleo", "desempleo", "empleo", "trabajo", "sueldo"),
		fill("Fill in the blank", "La ___ entre ricos y pobres crece. (inequality)", "desigualdad"),
		tr("Translate this sentence", "We must protect the environment", "Hay que proteger el medio ambiente"),
		mc("development", "el ___ sostenible", "desarrollo", "desarrollo", "desempleo", "destino", "deseo"),
		speak("Blaze", "El cambio climático es uno de los mayores retos de hoy."),
	)
	addVocab(db, l,
		vw("el medio ambiente", "the environment", "Cuidamos el medio ambiente.", "We care for the environment.", finch),
		vw("el cambio climático", "climate change", "El cambio climático avanza.", "Climate change advances.", "Cora"),
		vw("la desigualdad", "inequality", "Hay mucha desigualdad.", "There's much inequality.", "Cora"),
		vw("el desempleo", "unemployment", "El desempleo bajó.", "Unemployment fell.", "Lumora"),
		vw("el desarrollo", "development", "el desarrollo sostenible", "sustainable development", "Lumora"),
	)

	// ── Writing (B2) ──
	s = addSkill(db, u, "Escribir: Ensayo y Carta Formal", "250+ words, register-aware.", "PenLine", "#00C2A8", 63, 2060)
	l = addLesson(db, s, "Escribe un Ensayo", 1, 30,
		char("Lumora", "B2 writing is structured and varied: a clear thesis, balanced arguments, formal register, rich connectors. 250+ words."),
		mc("formal opening", "___ señor/a:", "Estimado/a", "Estimado/a", "Hola", "Querido", "Oye"),
		write("Write an argumentative essay (250+ words): 'Is remote work better than office work?' Present both sides and your conclusion.",
			"El teletrabajo se ha convertido en un tema de debate. Por un lado, ofrece flexibilidad y ahorra tiempo de transporte; de hecho, muchos trabajadores afirman ser más productivos en casa. Por otro lado, puede provocar aislamiento y dificultar la separación entre la vida laboral y personal. Cabe destacar que no todos los empleos permiten esta modalidad. En mi opinión, lo ideal sería un modelo híbrido que combine las ventajas de ambos. En conclusión, el teletrabajo no es mejor ni peor, sino una opción que depende de cada persona y profesión."),
		write("Write a formal letter to a newspaper (around 200 words) giving your opinion on a social issue.",
			"Estimado director: Le escribo para expresar mi preocupación por la contaminación en nuestra ciudad. A pesar de las campañas, el tráfico sigue aumentando y, por consiguiente, la calidad del aire empeora. Considero que las autoridades deberían fomentar el transporte público y crear más zonas verdes. No obstante, la responsabilidad también es ciudadana. Atentamente, Pablo Ruiz."),
		mc("formal connector", "___, cabe destacar su importancia.", "Asimismo", "Asimismo", "Oye", "Vale", "Pues"),
		fill("Fill in the blank", "Por un lado…; por otro ___ . (side)", "lado"),
		speak("Blaze", "En conclusión, lo ideal sería un modelo híbrido."),
	)
	addVocab(db, l,
		vw("Estimado/a", "Dear (formal)", "Estimado director:", "Dear editor:", "Lumora"),
		vw("asimismo", "likewise / also", "Asimismo, propongo…", "Likewise, I propose…", "Cora"),
		vw("considero que", "I consider that", "Considero que es clave.", "I consider it key.", "Cora"),
		vw("lo ideal sería", "the ideal would be", "Lo ideal sería un acuerdo.", "The ideal would be a deal.", "Lumora"),
		vw("en definitiva", "ultimately", "En definitiva, depende.", "Ultimately, it depends.", "Lumora"),
	)

	// ── Speaking (B2) ──
	s = addSkill(db, u, "Hablar: Debate y Argumentación", "Defend a position fluently.", "MessageCircle", "#6C3FC5", 64, 2110)
	l = addLesson(db, s, "Defiende tu Postura", 1, 30,
		char("Blaze", "Debate time! State your view, support it, concede a point, and rebut — all out loud."),
		speak("Lumora", "Desde mi punto de vista, la educación es la base del progreso."),
		speak("Blaze", "Es cierto que tiene desventajas; sin embargo, las ventajas son mayores."),
		speak("Blaze", "Si tuviéramos más recursos, podríamos resolver el problema."),
		mc("to concede a point", "___ razón, pero…", "Tienes", "Tienes", "Haces", "Das", "Eres"),
		speak("Blaze", "Por un lado entiendo tu postura; por otro, no la comparto del todo."),
		speak("Blaze", "En definitiva, creo que deberíamos buscar un punto intermedio."),
	)
	addVocab(db, l,
		vw("desde mi punto de vista", "from my point of view", "Desde mi punto de vista, sí.", "From my view, yes.", "Blaze"),
		vw("es cierto que", "it's true that", "Es cierto que es caro.", "It's true it's expensive.", "Cora"),
		vw("tienes razón", "you're right", "Tienes razón en parte.", "You're partly right.", "Cora"),
		vw("no la comparto", "I don't share it", "No comparto esa idea.", "I don't share that idea.", "Lumora"),
		vw("un punto intermedio", "a middle ground", "Busquemos un punto intermedio.", "Let's find middle ground.", "Lumora"),
	)

	// C1: sophistication, stylistic range and nuance — total subjunctive
	// command, mixed conditionals, advanced relatives, the neuter 'lo',
	// register shifts, word formation, idioms/proverbs and academic discourse.
	seedSpanishC1(db)
}

// seedSpanishC1 adds the advanced (C1) unit: full subjunctive across tenses,
// mixed conditionals, indicative/subjunctive nuance, advanced relatives, the
// neuter 'lo' and nominalisation, advanced periphrasis, reformulation
// connectors, word formation, register, idioms & proverbs, and the productive
// C1 skills (argumentative essay & formal exposition).
func seedSpanishC1(db *gorm.DB) {
	const u = "C1 · Superior"
	finch := "Professor Finch"

	// ── Subjunctive: full command (perfect subjunctive) ──
	s := addSkill(db, u, "Subjuntivo: Dominio Total", "haya hablado, hubiera hablado.", "Quote", "#FF5C5C", 65, 2160)
	l := addLesson(db, s, "Todos los Subjuntivos", 1, 30,
		char(finch, "Master all four: presente (hable), perfecto (haya hablado), imperfecto (hablara), pluscuamperfecto (hubiera hablado)."),
		mc("I hope he has arrived", "Espero que ___ llegado", "haya", "haya", "ha", "había", "habrá"),
		mc("I doubted he had come", "Dudaba que ___ venido", "hubiera", "hubiera", "había", "ha", "habría"),
		fill("Fill in the blank", "Me alegro de que ___ aprobado. (haber, tú - perfect subj.)", "hayas"),
		tr("Translate this sentence", "as if nothing had happened", "como si nada hubiera pasado"),
		mc("present subj. after 'cuando' (future)", "Cuando ___ tiempo, te aviso.", "tenga", "tenga", "tengo", "tendré", "tenía"),
		speak("Blaze", "Me alegra que hayas venido y ojalá hubiéramos hablado antes."),
	)
	addVocab(db, l,
		vw("haya hablado", "(that) I have spoken", "que yo haya hablado", "that I have spoken", finch),
		vw("hubiera hablado", "(that) I had spoken", "como si hubiera hablado", "as if I had spoken", "Cora"),
		vw("como si", "as if (+ past subj.)", "Habla como si supiera.", "He talks as if he knew.", "Cora"),
		vw("me alegro de que", "I'm glad that", "Me alegro de que vengas.", "I'm glad you're coming.", "Lumora"),
		vw("ojalá hubiera", "I wish I had", "Ojalá hubiera ido.", "I wish I had gone.", "Lumora"),
	)

	// ── Mixed conditionals ──
	s = addSkill(db, u, "Condicionales Mixtas", "Si hubiera…, ahora …ría.", "Layers", "#6C3FC5", 66, 2220)
	l = addLesson(db, s, "Hipótesis Mixtas", 1, 30,
		char(finch, "Past condition with present result: Si hubiera estudiado medicina, ahora sería médico."),
		mc("If I had saved (then)…", "Si ___ ahorrado, ahora tendría casa.", "hubiera", "hubiera", "había", "habría", "tuviera"),
		mc("…now I would have a house", "Si hubiera ahorrado, ahora ___ casa.", "tendría", "tendría", "tengo", "tuviera", "habría tenido"),
		fill("Fill in the blank", "Si me ___ hecho caso, no estarías así. (haber, tú - subj.)", "hubieras"),
		tr("Translate this sentence", "If I were taller, I would have been a model", "Si fuera más alto, habría sido modelo"),
		mc("present condition, past result", "Si ___ responsable, no habrías fallado.", "fueras", "fueras", "eres", "serías", "fuiste"),
		speak("Blaze", "Si hubiera nacido en España, ahora hablaría español perfecto."),
	)
	addVocab(db, l,
		vw("si hubiera …, ahora …ría", "had I …, now I'd …", "Si hubiera estudiado, ahora trabajaría.", "Had I studied, I'd work now.", finch),
		vw("ahora sería", "now I would be", "Ahora sería médico.", "I'd be a doctor now.", "Cora"),
		vw("de haberlo sabido", "had I known", "De haberlo sabido, vengo.", "Had I known, I'd come.", "Cora"),
		vw("en tu lugar", "in your place", "En tu lugar, lo haría.", "In your place, I'd do it.", "Lumora"),
		vw("a no ser que", "unless", "Iré, a no ser que llueva.", "I'll go, unless it rains.", "Lumora"),
	)

	// ── Indicative vs subjunctive nuance ──
	s = addSkill(db, u, "Indicativo/Subjuntivo: Matices", "probability, concession, attitude.", "Quote", "#00C2A8", 67, 2280)
	l = addLesson(db, s, "Matices del Modo", 1, 30,
		char(finch, "The mood changes the meaning: 'quizás viene' (likely) vs 'quizás venga' (less sure); 'el hecho de que' values, not states."),
		mc("probably (more certain)", "Quizás ___ hoy. (likely → indicative)", "viene", "viene", "venga", "viniera", "vendría"),
		mc("possibly (less certain)", "Tal vez ___ mañana. (doubt → subjunctive)", "venga", "venga", "viene", "vino", "vendrá"),
		fill("Fill in the blank", "No es que no ___, es que no puedo. (querer - subj.)", "quiera"),
		tr("Translate this sentence", "The fact that he lied bothers me", "El hecho de que mintiera me molesta"),
		mc("certainty → indicative", "Es evidente que ___ verdad.", "es", "es", "sea", "fuera", "haya sido"),
		speak("Blaze", "No es que sea difícil, es que requiere práctica."),
	)
	addVocab(db, l,
		vw("quizás / tal vez", "maybe", "Tal vez venga.", "Maybe he'll come.", finch),
		vw("no es que", "it's not that", "No es que no quiera.", "It's not that I don't want to.", "Cora"),
		vw("el hecho de que", "the fact that", "El hecho de que mienta…", "The fact that he lies…", "Cora"),
		vw("es evidente que", "it's clear that", "Es evidente que sí.", "It's clearly so.", "Lumora"),
		vw("siempre y cuando", "provided that", "Iré siempre y cuando vengas.", "I'll go provided you come.", "Lumora"),
	)

	// ── Advanced relatives ──
	s = addSkill(db, u, "Relativos Avanzados", "el cual, lo cual, cuyo.", "Link2", "#17A3DD", 68, 2340)
	l = addLesson(db, s, "El cual, Lo cual", 1, 30,
		char(finch, "After prepositions or for clarity, use el/la/los/las cual(es). 'lo cual' refers to a whole idea. 'cuyo' = whose."),
		mc("the reason for which I came", "la razón por la ___ vine", "cual", "cual", "que", "quien", "cuyo"),
		mc("…, which annoyed me", "No vino, ___ me molestó.", "lo cual", "lo cual", "el cual", "que", "cuyo"),
		fill("Fill in the blank", "El autor ___ novela leí es famoso. (whose)", "cuya"),
		tr("Translate this sentence", "the house in which I grew up", "la casa en la cual crecí"),
		mc("everything (that), which → neuter", "Hizo lo que pudo, ___ fue suficiente.", "lo cual", "lo cual", "el cual", "la cual", "que"),
		speak("Blaze", "Llegó tarde, lo cual no me sorprendió en absoluto."),
	)
	addVocab(db, l,
		vw("el cual / la cual", "which (after prep.)", "el motivo por el cual…", "the reason for which…", finch),
		vw("lo cual", "which (whole idea)", "Se fue, lo cual me dolió.", "He left, which hurt me.", "Cora"),
		vw("cuyo / cuya", "whose", "el libro cuyo autor…", "the book whose author…", "Cora"),
		vw("quienes", "those who", "Quienes estudian, aprueban.", "Those who study, pass.", "Lumora"),
		vw("en el cual", "in which", "el año en el cual nací", "the year in which I was born", "Lumora"),
	)

	// ── Neuter 'lo' & nominalisation ──
	s = addSkill(db, u, "Lo Neutro y Nominalización", "lo bueno, lo que, lo de.", "Hash", "#F5A623", 69, 2400)
	l = addLesson(db, s, "El Uso de 'Lo'", 1, 30,
		char(finch, "'lo' + adjective abstracts a quality (lo importante = the important thing). 'lo que' = what. 'lo de' = the matter of."),
		mc("the good thing", "___ bueno es que vino.", "Lo", "Lo", "El", "La", "Los"),
		mc("what matters", "___ que importa es la salud.", "Lo", "Lo", "El", "Que", "La"),
		fill("Fill in the blank", "___ de ayer fue increíble. (the matter of)", "Lo"),
		tr("Translate this sentence", "the best thing about the trip", "lo mejor del viaje"),
		mc("how + adjective (exclamation)", "No sabes ___ difícil que es.", "lo", "lo", "el", "que", "cuán"),
		speak("Blaze", "Lo importante es disfrutar; lo demás no importa tanto."),
	)
	addVocab(db, l,
		vw("lo bueno", "the good thing", "Lo bueno es que aprendí.", "The good thing is I learned.", finch),
		vw("lo que", "what", "Lo que dices es cierto.", "What you say is true.", "Cora"),
		vw("lo de", "the matter of", "Lo de ayer fue raro.", "Yesterday's thing was odd.", "Cora"),
		vw("lo mejor", "the best (thing)", "Lo mejor del día.", "The best of the day.", "Lumora"),
		vw("lo demás", "the rest", "Lo demás no importa.", "The rest doesn't matter.", "Lumora"),
	)

	// ── Advanced periphrasis ──
	s = addSkill(db, u, "Perífrasis Avanzadas", "ir/seguir/llevar + gerundio.", "Link2", "#FF5C5C", 70, 2460)
	l = addLesson(db, s, "Matices con Perífrasis", 1, 30,
		char(finch, "Nuance through periphrasis: ir + gerund (gradual), seguir + gerund (still), llevar + gerund (duration), tener + participle (result)."),
		mc("I'm gradually understanding", "Voy ___ poco a poco.", "entendiendo", "entendiendo", "entender", "entendido", "entiendo"),
		mc("I'm still working", "Sigo ___ aquí.", "trabajando", "trabajando", "trabajar", "trabajado", "trabajo"),
		fill("Fill in the blank", "___ tres horas esperando. (llevar, yo - duration)", "Llevo"),
		tr("Translate this sentence", "I have three chapters written", "Tengo tres capítulos escritos"),
		mc("he ended up accepting", "Acabó ___ la oferta.", "aceptando", "aceptando", "aceptar", "aceptado", "acepta"),
		speak("Blaze", "Llevo años estudiando y voy mejorando poco a poco."),
	)
	addVocab(db, l,
		vw("ir + gerundio", "to gradually …", "Voy aprendiendo.", "I'm gradually learning.", finch),
		vw("seguir + gerundio", "to still …", "Sigo viviendo aquí.", "I still live here.", "Cora"),
		vw("llevar + gerundio", "to have been …-ing", "Llevo años aquí.", "I've been here for years.", "Cora"),
		vw("tener + participio", "to have (done)", "Tengo hecho el trabajo.", "I have the work done.", "Lumora"),
		vw("acabar + gerundio", "to end up …-ing", "Acabó aceptando.", "He ended up accepting.", "Lumora"),
	)

	// ── Reformulation connectors ──
	s = addSkill(db, u, "Conectores y Reformulación", "en resumidas cuentas, por ende.", "MessageCircle", "#06AECE", 71, 2520)
	l = addLesson(db, s, "Reformular y Concluir", 1, 30,
		char("Cora", "Refine your discourse: en resumidas cuentas (in short), dicho de otro modo (in other words), por ende (hence), cabe señalar (it should be noted)."),
		mc("in short", "___, fue un éxito rotundo.", "En resumidas cuentas", "En resumidas cuentas", "Por ejemplo", "En cambio", "De hecho"),
		mc("in other words", "Es complejo; ___, no es fácil.", "dicho de otro modo", "dicho de otro modo", "por ende", "es decir que no", "sin embargo"),
		fill("Fill in the blank", "Llueve; por ___, suspenden el evento. (hence)", "ende"),
		tr("Translate this sentence", "It should be noted that it's optional", "Cabe señalar que es opcional"),
		mc("that is to say", "Es bilingüe, ___, habla dos lenguas.", "es decir", "es decir", "no obstante", "en cambio", "a pesar de"),
		speak("Blaze", "En resumidas cuentas, dicho de otro modo, mereció la pena."),
	)
	addVocab(db, l,
		vw("en resumidas cuentas", "in short", "En resumidas cuentas, sí.", "In short, yes.", "Cora"),
		vw("dicho de otro modo", "in other words", "Dicho de otro modo, no.", "In other words, no.", "Cora"),
		vw("por ende", "hence", "Falló; por ende, repite.", "He failed; hence, he repeats.", "Lumora"),
		vw("cabe señalar", "it should be noted", "Cabe señalar el riesgo.", "The risk should be noted.", "Lumora"),
		vw("es decir", "that is to say", "Es tarde, es decir, vamos.", "It's late, that is, let's go.", "Riko"),
	)

	// ── Word formation ──
	s = addSkill(db, u, "Formación de Palabras", "prefixes, suffixes, derivation.", "Languages", "#6C3FC5", 72, 2580)
	l = addLesson(db, s, "Derivación y Afijos", 1, 30,
		char(finch, "Build words: des-/in- (negation), re- (again); -ción/-dad/-eza (nouns); -oso/-able (adjectives); -ón (augmentative), -ito (diminutive)."),
		mc("opposite of 'ordenado'", "des___", "ordenado", "ordenado", "ordenar", "orden", "ordenando"),
		mc("noun from 'feliz'", "la ___", "felicidad", "felicidad", "felizmente", "felizar", "feliza"),
		fill("Fill in the blank", "Casa grande → un cas___ (augmentative)", "ón"),
		tr("Translate this sentence", "It's unforgettable", "Es inolvidable"),
		mc("adjective from 'cuidado'", "una persona ___", "cuidadosa", "cuidadosa", "cuidar", "cuidado", "cuidadamente"),
		speak("Blaze", "La rapidez y la amabilidad del servicio fueron increíbles."),
	)
	addVocab(db, l,
		vw("des-", "un-/dis- (prefix)", "deshacer, desordenado", "to undo, untidy", finch),
		vw("-ción", "-tion (noun suffix)", "la educación", "education", "Cora"),
		vw("-dad", "-ity (noun suffix)", "la felicidad", "happiness", "Cora"),
		vw("-mente", "-ly (adverb)", "rápidamente", "quickly", "Lumora"),
		vw("-ito / -ón", "diminutive / augmentative", "perrito, casón", "little dog, big house", "Lumora"),
	)

	// ── Register ──
	s = addSkill(db, u, "Registro: Formal y Coloquial", "Match the situation.", "Quote", "#17A3DD", 73, 2640)
	l = addLesson(db, s, "Cambiar de Registro", 1, 30,
		char("Cora", "C1 means shifting register: formal (le agradezco, ¿sería tan amable?) vs colloquial (mola, qué guay, vale, tío)."),
		mc("colloquial 'cool'", "¡Qué ___!", "guay", "guay", "amable", "estimado", "cordial"),
		mc("formal request", "¿___ tan amable de ayudarme?", "Sería", "Sería", "Eres", "Estás", "Vas"),
		fill("Fill in the blank", "Formal thanks: Le ___ su atención. (agradecer)", "agradezco"),
		tr("Translate this sentence", "It's great! (colloquial)", "¡Mola mucho!"),
		mc("informal 'okay'", "— ¿Vamos? — ___.", "Vale", "Vale", "Estimado", "Atentamente", "Le ruego"),
		speak("Blaze", "Con amigos digo '¡qué guay!'; en el trabajo, 'me parece excelente'."),
	)
	addVocab(db, l,
		vw("le agradezco", "I thank you (formal)", "Le agradezco su ayuda.", "I thank you for your help.", finch),
		vw("¿sería tan amable?", "would you be so kind?", "¿Sería tan amable?", "Would you be so kind?", "Cora"),
		vw("mola / qué guay", "it's cool (colloq.)", "¡Mola mucho!", "It's really cool!", "Cora"),
		vw("vale", "okay (colloq.)", "Vale, de acuerdo.", "Okay, agreed.", "Lumora"),
		vw("le ruego", "I beg/request (formal)", "Le ruego disculpas.", "I beg your pardon.", "Lumora"),
	)

	// ── Idioms & proverbs ──
	s = addSkill(db, u, "Modismos y Refranes", "no hay mal que por bien no venga.", "Sparkles", "#FF5C5C", 74, 2700)
	l = addLesson(db, s, "Refranes y Frases Hechas", 1, 30,
		char("Cora", "Proverbs carry culture: 'No hay mal que por bien no venga', 'Más vale tarde que nunca', 'A quien madruga, Dios le ayuda'."),
		mc("every cloud has a silver lining", "No hay mal que por bien no ___", "venga", "venga", "viene", "vino", "vendrá"),
		mc("better late than never", "Más vale tarde que ___", "nunca", "nunca", "siempre", "ahora", "pronto"),
		fill("Fill in the blank", "A quien ___ , Dios le ayuda. (rises early)", "madruga"),
		tr("Translate this sentence", "to cost an arm and a leg", "costar un ojo de la cara"),
		mc("to be a piece of cake", "Es pan ___", "comido", "comido", "tostado", "duro", "rico"),
		speak("Blaze", "Más vale tarde que nunca: por fin aprobé el examen."),
	)
	addVocab(db, l,
		vw("no hay mal que…", "every cloud has a silver lining", "No hay mal que por bien no venga.", "Every cloud has a silver lining.", "Cora"),
		vw("más vale tarde que nunca", "better late than never", "Más vale tarde que nunca.", "Better late than never.", "Cora"),
		vw("costar un ojo de la cara", "to cost a fortune", "Cuesta un ojo de la cara.", "It costs a fortune.", "Lumora"),
		vw("ser pan comido", "to be a piece of cake", "Es pan comido.", "It's a piece of cake.", "Lumora"),
		vw("estar en las nubes", "to have one's head in the clouds", "Estás en las nubes.", "You're daydreaming.", "Riko"),
	)

	// ── Impersonal / advanced passive ──
	s = addSkill(db, u, "Impersonalidad y Pasiva", "se dice que, uno, voz pasiva.", "Hash", "#00C2A8", 75, 2760)
	l = addLesson(db, s, "Construcciones Impersonales", 1, 30,
		char(finch, "Express generality: se dice que (it's said), uno nunca sabe (one never knows), and the full passive for formal style."),
		mc("it's said that…", "___ dice que es bueno.", "Se", "Se", "Le", "Lo", "Uno"),
		mc("one never knows", "___ nunca sabe.", "Uno", "Uno", "Se", "Le", "Tú"),
		fill("Fill in the blank", "La ley fue ___ por el parlamento. (approve, fem.)", "aprobada"),
		tr("Translate this sentence", "Mistakes were made", "Se cometieron errores"),
		mc("people say (impersonal 3rd pl.)", "___ que va a llover.", "Dicen", "Dicen", "Dice", "Se dicen", "Digo"),
		speak("Blaze", "Se dice que la práctica hace al maestro."),
	)
	addVocab(db, l,
		vw("se dice que", "it's said that", "Se dice que es difícil.", "It's said it's hard.", finch),
		vw("uno", "one (impersonal)", "Uno nunca sabe.", "One never knows.", "Cora"),
		vw("fue aprobada", "was approved", "La ley fue aprobada.", "The law was approved.", "Cora"),
		vw("se cometieron", "were committed", "Se cometieron errores.", "Mistakes were made.", "Lumora"),
		vw("dicen que", "they say that", "Dicen que es genial.", "They say it's great.", "Lumora"),
	)

	// ── Culture, art & current affairs vocabulary ──
	s = addSkill(db, u, "Léxico: Cultura y Actualidad", "arts, media, abstract terms.", "BookOpen", "#F5A623", 76, 2820)
	l = addLesson(db, s, "Cultura y Sociedad", 1, 30,
		char(finch, "C1 vocabulary is precise and abstract: el patrimonio, la trama, la crítica, el prejuicio, la libertad de expresión."),
		mc("heritage", "el ___ cultural", "patrimonio", "patrimonio", "prejuicio", "estreno", "guion"),
		mc("the plot (of a film/book)", "la ___ de la novela", "trama", "trama", "crítica", "escena", "obra"),
		fill("Fill in the blank", "La libertad de ___ es un derecho. (expression)", "expresión"),
		tr("Translate this sentence", "The film received good reviews", "La película recibió buenas críticas"),
		mc("prejudice", "un ___ social", "prejuicio", "prejuicio", "patrimonio", "guion", "ensayo"),
		speak("Blaze", "La trama de la novela critica los prejuicios de la sociedad."),
	)
	addVocab(db, l,
		vw("el patrimonio", "heritage", "patrimonio de la humanidad", "world heritage", finch),
		vw("la trama", "the plot", "La trama es compleja.", "The plot is complex.", "Cora"),
		vw("la crítica", "review / criticism", "Buenas críticas.", "Good reviews.", "Cora"),
		vw("el prejuicio", "prejudice", "Hay muchos prejuicios.", "There are many prejudices.", "Lumora"),
		vw("la libertad de expresión", "freedom of speech", "Defiendo la libertad de expresión.", "I defend free speech.", "Lumora"),
	)

	// ── Writing (C1) ──
	s = addSkill(db, u, "Escribir: Ensayo Argumentativo", "300+ words, cohesive & precise.", "PenLine", "#00C2A8", 77, 2880)
	l = addLesson(db, s, "Ensayo Avanzado", 1, 32,
		char("Lumora", "C1 writing is cohesive and nuanced: a thesis, developed arguments, counter-arguments, refined connectors and varied syntax. 300+ words."),
		mc("to introduce a thesis", "___, conviene definir el problema.", "En primer lugar", "En primer lugar", "Vale", "Oye", "Total"),
		write("Write an argumentative essay (300+ words): 'Should artificial intelligence be regulated?' Develop arguments, counter-arguments and a reasoned conclusion.",
			"La irrupción de la inteligencia artificial ha transformado nuestra sociedad a un ritmo vertiginoso. En primer lugar, conviene señalar sus innegables beneficios: optimiza procesos, impulsa la investigación médica y facilita tareas cotidianas. No obstante, este avance plantea serios dilemas éticos, como la pérdida de empleos o el uso indebido de datos personales. Hay quienes sostienen que regularla frenaría la innovación; sin embargo, la ausencia de límites podría tener consecuencias imprevisibles. Dicho de otro modo, no se trata de prohibir, sino de establecer un marco que garantice un desarrollo responsable. En definitiva, considero que una regulación equilibrada resulta imprescindible para que la tecnología esté al servicio de las personas y no al revés."),
		write("Write a formal report (around 250 words) summarising a problem in your city and proposing solutions.",
			"El presente informe analiza el problema del tráfico en el centro de la ciudad. En los últimos años, la congestión ha aumentado considerablemente, lo cual repercute en la calidad del aire y en la salud de los ciudadanos. Por un lado, cabe destacar la falta de transporte público eficiente; por otro, el uso excesivo del vehículo privado. Como solución, se propone ampliar la red de metro y crear carriles bici. Asimismo, sería conveniente fomentar el teletrabajo. En conclusión, solo una estrategia integral permitirá revertir esta tendencia."),
		mc("nuanced concession", "Hay quienes opinan lo contrario; ___, no les falta razón.", "de hecho", "de hecho", "vale", "total", "oye"),
		fill("Fill in the blank", "En ___ lugar, conviene matizar. (first)", "primer"),
		speak("Blaze", "En definitiva, una regulación equilibrada resulta imprescindible."),
	)
	addVocab(db, l,
		vw("en primer lugar", "firstly", "En primer lugar, definamos.", "Firstly, let's define.", "Lumora"),
		vw("conviene señalar", "it's worth pointing out", "Conviene señalar el riesgo.", "Worth pointing out the risk.", "Cora"),
		vw("hay quienes sostienen", "some maintain", "Hay quienes sostienen eso.", "Some maintain that.", "Cora"),
		vw("no se trata de", "it's not about", "No se trata de prohibir.", "It's not about banning.", "Lumora"),
		vw("resulta imprescindible", "it's essential", "Resulta imprescindible actuar.", "It's essential to act.", "Lumora"),
	)

	// ── Speaking (C1) ──
	s = addSkill(db, u, "Hablar: Exposición y Debate", "Present and argue with nuance.", "MessageCircle", "#6C3FC5", 78, 2940)
	l = addLesson(db, s, "Expón y Argumenta", 1, 32,
		char("Blaze", "Deliver a clear, nuanced argument: frame it, develop it, concede, rebut and conclude — fluently."),
		speak("Lumora", "El tema que voy a tratar suscita un intenso debate en la actualidad."),
		speak("Blaze", "Si bien es cierto que existen riesgos, los beneficios los superan con creces."),
		speak("Blaze", "Cabe matizar que no todos los casos son iguales."),
		mc("to nuance a claim", "___ que matizar esa afirmación.", "Habría", "Habría", "Hace", "Tiene", "Va"),
		speak("Blaze", "Dicho esto, no comparto del todo esa visión tan pesimista."),
		speak("Blaze", "En resumidas cuentas, la solución pasa por el equilibrio."),
	)
	addVocab(db, l,
		vw("suscitar debate", "to spark debate", "El tema suscita debate.", "The topic sparks debate.", "Blaze"),
		vw("si bien es cierto que", "while it's true that", "Si bien es cierto que…", "While it's true that…", "Cora"),
		vw("cabe matizar", "it should be nuanced", "Cabe matizar esto.", "This should be nuanced.", "Cora"),
		vw("dicho esto", "that said", "Dicho esto, discrepo.", "That said, I disagree.", "Lumora"),
		vw("superar con creces", "to far exceed", "Los superan con creces.", "They far exceed them.", "Lumora"),
	)

	// C2: mastery — stylistic precision, rhetoric, register virtuosity,
	// regional awareness and near-native idiomaticity. The pinnacle of the
	// course.
	seedSpanishC2(db)
}

// seedSpanishC2 adds the mastery (C2) unit: subtle mood/aspect, literary tenses,
// cleft sentences & emphasis, complex subordination, hedging, reported speech
// with attitude, fine ser/estar–por/para–article use, rhetorical devices,
// regional variation, advanced idioms, register virtuosity, specialised lexis,
// and the productive C2 skills (critical essay & high-level oratory).
func seedSpanishC2(db *gorm.DB) {
	const u = "C2 · Maestría"
	finch := "Professor Finch"

	// ── Subtle mood ──
	s := addSkill(db, u, "Matices del Modo", "Subjective vs factual shades.", "Quote", "#FF5C5C", 79, 3000)
	l := addLesson(db, s, "Indicativo o Subjuntivo (sutil)", 1, 32,
		char(finch, "At C2 the mood conveys attitude: 'el que diga eso miente' (whoever may say) vs 'el que dice eso miente' (the one who says)."),
		mc("whoever says that (unknown)", "El que ___ eso, miente.", "diga", "diga", "dice", "decía", "dijo"),
		mc("not that I doubt it, but…", "No es que lo ___, pero…", "dude", "dude", "dudo", "dudaba", "dudaría"),
		fill("Fill in the blank", "Por mucho que ___, no lo lograrás. (insistir, tú - subj.)", "insistas"),
		tr("Translate this sentence", "Say what you may, I won't change my mind", "Digas lo que digas, no cambiaré de idea"),
		mc("hope tinged with doubt", "Que yo ___, nunca falló.", "sepa", "sepa", "sé", "sabía", "sabré"),
		speak("Blaze", "Digan lo que digan, mantendré mi postura."),
	)
	addVocab(db, l,
		vw("digas lo que digas", "say what you may", "Digas lo que digas, no.", "Whatever you say, no.", finch),
		vw("por mucho que", "however much", "Por mucho que insistas…", "However much you insist…", "Cora"),
		vw("que yo sepa", "as far as I know", "Que yo sepa, no vino.", "As far as I know, he didn't come.", "Cora"),
		vw("no sea que", "lest", "Apúrate, no sea que llueva.", "Hurry, lest it rain.", "Lumora"),
		vw("el que / quien", "whoever", "Quien lo diga, miente.", "Whoever says it lies.", "Lumora"),
	)

	// ── Literary tenses & aspect ──
	s = addSkill(db, u, "Tiempos y Aspecto Literarios", "Stylistic & narrative uses.", "Clock", "#6C3FC5", 80, 3070)
	l = addLesson(db, s, "Usos Estilísticos del Tiempo", 1, 32,
		char(finch, "Literature exploits tense: the historical present for vividness, the imperfect for atmosphere, and the rare pretérito anterior (apenas hubo llegado)."),
		mc("historical present (vivid past)", "En 1492, Colón ___ a América.", "llega", "llega", "llegará", "llegaría", "llegaba"),
		mc("as soon as he had finished", "Apenas ___ terminado, salió.", "hubo", "hubo", "había", "ha", "habrá"),
		fill("Fill in the blank", "El sol ___ mientras ella lloraba. (set - imperfect, atmosphere)", "caía"),
		tr("Translate this sentence", "He would die three years later (narrative future-in-past)", "Moriría tres años después"),
		mc("narrative conditional", "Aquel error le ___ caro con el tiempo.", "costaría", "costaría", "cuesta", "costó", "cueste"),
		speak("Blaze", "Caía la tarde cuando, de pronto, todo cambió."),
	)
	addVocab(db, l,
		vw("presente histórico", "historical present", "En 1936 estalla la guerra.", "In 1936 war breaks out.", finch),
		vw("apenas hubo …", "as soon as he had …", "Apenas hubo llegado…", "As soon as he had arrived…", "Cora"),
		vw("caía la tarde", "evening was falling", "Caía la tarde.", "Evening was falling.", "Cora"),
		vw("moriría … después", "would die … later", "Moriría años después.", "He would die years later.", "Lumora"),
		vw("de pronto", "suddenly", "De pronto, calló.", "Suddenly, he fell silent.", "Lumora"),
	)

	// ── Cleft sentences & emphasis ──
	s = addSkill(db, u, "Oraciones Escindidas", "Fue … quien; lo que … es.", "Hash", "#00C2A8", 81, 3140)
	l = addLesson(db, s, "Énfasis y Foco", 1, 32,
		char(finch, "Highlight with cleft structures: 'Fue Ana quien llamó', 'Lo que necesito es tiempo', 'Por eso es por lo que vine'."),
		mc("It was Ana who called", "Fue Ana ___ llamó.", "quien", "quien", "que", "la cual", "cuya"),
		mc("What I need is time", "___ necesito es tiempo.", "Lo que", "Lo que", "El que", "Que", "Cual"),
		fill("Fill in the blank", "Fue allí ___ lo conocí. (where, cleft)", "donde"),
		tr("Translate this sentence", "It's you who decides", "Eres tú quien decide"),
		mc("emphatic reason", "Por eso es ___ lo vine.", "por lo que", "por lo que", "que", "para que", "lo cual"),
		speak("Blaze", "Lo que de verdad importa es cómo reaccionamos."),
	)
	addVocab(db, l,
		vw("fue … quien", "it was … who", "Fue ella quien ganó.", "It was she who won.", finch),
		vw("lo que … es", "what … is", "Lo que quiero es paz.", "What I want is peace.", "Cora"),
		vw("es … donde", "it's … where", "Es aquí donde vivo.", "It's here that I live.", "Cora"),
		vw("eres tú quien", "it's you who", "Eres tú quien manda.", "It's you who's in charge.", "Lumora"),
		vw("por eso es por lo que", "that's why", "Por eso es por lo que vine.", "That's why I came.", "Lumora"),
	)

	// ── Complex subordination & ellipsis ──
	s = addSkill(db, u, "Subordinación y Elipsis", "Dense, economical syntax.", "Link2", "#17A3DD", 82, 3210)
	l = addLesson(db, s, "Sintaxis Compleja", 1, 32,
		char(finch, "Pack meaning elegantly: gerund/participle clauses (terminada la reunión, …) and ellipsis (unos prefieren té; otros, café)."),
		mc("the meeting over, we left", "___ la reunión, nos fuimos.", "Terminada", "Terminada", "Terminado", "Terminando", "Terminar"),
		mc("once said this…", "___ esto, continuó.", "Dicho", "Dicho", "Diciendo", "Decir", "Dije"),
		fill("Fill in the blank", "Unos querían ir; otros, ___ . (ellipsis: no)", "no"),
		tr("Translate this sentence", "Knowing the truth, he stayed silent", "Conociendo la verdad, calló"),
		mc("absolute participle", "Una vez ___ el problema, brindamos.", "resuelto", "resuelto", "resolviendo", "resolver", "resuelve"),
		speak("Blaze", "Terminado el discurso, y conmovidos todos, aplaudieron."),
	)
	addVocab(db, l,
		vw("terminada la reunión", "the meeting over", "Terminada la cena, salimos.", "Dinner over, we left.", finch),
		vw("dicho esto", "this said", "Dicho esto, me marcho.", "This said, I leave.", "Cora"),
		vw("una vez + part.", "once …", "Una vez resuelto, descansé.", "Once solved, I rested.", "Cora"),
		vw("conociendo", "knowing", "Conociéndolo, no me extraña.", "Knowing him, no surprise.", "Lumora"),
		vw("unos…; otros…", "some…; others…", "Unos sí; otros, no.", "Some do; others don't.", "Lumora"),
	)

	// ── Hedging & nuance connectors ──
	s = addSkill(db, u, "Atenuación y Matiz", "Hedging: por así decirlo.", "MessageCircle", "#F5A623", 83, 3280)
	l = addLesson(db, s, "Suavizar y Matizar", 1, 32,
		char("Cora", "Sound diplomatic and precise: por así decirlo (so to speak), en cierto modo (in a way), si no me equivoco, hasta cierto punto."),
		mc("so to speak", "Es, ___, un genio.", "por así decirlo", "por así decirlo", "sin duda", "por supuesto", "en absoluto"),
		mc("in a way", "___, tienes razón.", "En cierto modo", "En cierto modo", "Jamás", "Sin falta", "Por completo"),
		fill("Fill in the blank", "Si no me ___, fue en 1999. (be mistaken)", "equivoco"),
		tr("Translate this sentence", "up to a point, I agree", "hasta cierto punto, estoy de acuerdo"),
		mc("if I may say so", "Es, ___, discutible.", "todo sea dicho", "todo sea dicho", "sin más", "ni hablar", "por cierto"),
		speak("Blaze", "En cierto modo, y por así decirlo, ambos llevamos razón."),
	)
	addVocab(db, l,
		vw("por así decirlo", "so to speak", "Es, por así decirlo, único.", "It's, so to speak, unique.", "Cora"),
		vw("en cierto modo", "in a way", "En cierto modo, sí.", "In a way, yes.", "Cora"),
		vw("si no me equivoco", "if I'm not mistaken", "Si no me equivoco, hoy.", "If I'm not mistaken, today.", "Lumora"),
		vw("hasta cierto punto", "up to a point", "Hasta cierto punto, vale.", "Up to a point, okay.", "Lumora"),
		vw("todo sea dicho", "it must be said", "Todo sea dicho, ayudó.", "It must be said, he helped.", "Riko"),
	)

	// ── Reported speech with attitude ──
	s = addSkill(db, u, "Discurso Referido con Actitud", "matizar, recalcar, insinuar.", "Languages", "#06AECE", 84, 3350)
	l = addLesson(db, s, "Verbos de Habla", 1, 32,
		char(finch, "Beyond 'decir': choose verbs that convey attitude — matizó (nuanced), recalcó (stressed), insinuó (implied), reprochó (reproached)."),
		mc("he stressed that…", "___ que era urgente.", "Recalcó", "Recalcó", "Preguntó", "Negó", "Dudó"),
		mc("she hinted that…", "___ que algo iba mal.", "Insinuó", "Insinuó", "Gritó", "Repitió", "Aclaró"),
		fill("Fill in the blank", "___ que, en realidad, no era tan simple. (nuanced: matizar)", "Matizó"),
		tr("Translate this sentence", "He denied having said it", "Negó haberlo dicho"),
		mc("he reproached me for…", "Me ___ no haber llamado.", "reprochó", "reprochó", "felicitó", "preguntó", "aclaró"),
		speak("Blaze", "Recalcó la urgencia, aunque insinuó que aún había tiempo."),
	)
	addVocab(db, l,
		vw("matizar", "to nuance/qualify", "Matizó su respuesta.", "He qualified his answer.", finch),
		vw("recalcar", "to stress", "Recalcó la idea.", "He stressed the idea.", "Cora"),
		vw("insinuar", "to imply", "Insinuó una crítica.", "He implied a criticism.", "Cora"),
		vw("reprochar", "to reproach", "Me reprochó mi tardanza.", "He reproached my lateness.", "Lumora"),
		vw("negar haber + part.", "to deny having", "Negó haber mentido.", "He denied having lied.", "Lumora"),
	)

	// ── Fine ser/estar, por/para, articles ──
	s = addSkill(db, u, "Usos Sutiles", "ser/estar, por/para, artículos.", "Sparkles", "#6C3FC5", 85, 3420)
	l = addLesson(db, s, "Detalles que Marcan la Diferencia", 1, 32,
		char(finch, "Subtlety: 'es claro' (it's logical) vs 'está claro' (it's evident); 'no para' vs 'no por'; the zero article in abstractions (tener paciencia)."),
		mc("evident (result/state)", "Está ___ que no vendrá.", "claro", "claro", "clara", "claros", "claramente"),
		mc("considering he's a child (para)", "___ ser un niño, sabe mucho.", "Para", "Para", "Por", "De", "A"),
		fill("Fill in the blank", "Hay que tener ___ . (patience — zero article)", "paciencia"),
		tr("Translate this sentence", "He is being unbearable (temporary)", "Está insoportable"),
		mc("ripe (estar) vs green (ser)", "El plátano está ___ . (ripe)", "maduro", "maduro", "verde", "joven", "viejo"),
		speak("Blaze", "Para ser tan joven, es de una madurez sorprendente."),
	)
	addVocab(db, l,
		vw("está claro", "it's evident", "Está claro que sí.", "It's clearly so.", finch),
		vw("para ser", "considering (he's)", "Para ser novato, va bien.", "For a novice, he's doing well.", "Cora"),
		vw("tener paciencia", "to be patient", "Ten paciencia.", "Be patient.", "Cora"),
		vw("estar insoportable", "to be (acting) unbearable", "Hoy estás insoportable.", "You're unbearable today.", "Lumora"),
		vw("ser vs estar (matiz)", "subtle ser/estar", "Es aburrido / está aburrido.", "He's boring / he's bored.", "Lumora"),
	)

	// ── Rhetorical devices ──
	s = addSkill(db, u, "Recursos Retóricos", "metáfora, ironía, hipérbole.", "Quote", "#FF5C5C", 86, 3490)
	l = addLesson(db, s, "El Arte de Persuadir", 1, 32,
		char("Cora", "Style persuades: metáfora (her eyes, two oceans), ironía (saying the opposite), hipérbole (exaggeration), paralelismo (repetition of structure)."),
		mc("'Te lo he dicho mil veces' is…", "exaggeration =", "hipérbole", "hipérbole", "ironía", "metáfora", "símil"),
		mc("'¡Qué bien, otra avería!' is…", "saying the opposite =", "ironía", "ironía", "hipérbole", "metáfora", "elipsis"),
		fill("Fill in the blank", "Comparación con 'como': es un ___ . (simile)", "símil"),
		tr("Translate this sentence", "Her smile was sunshine (metaphor)", "Su sonrisa era el sol"),
		mc("repetition of structure", "'Vine, vi, vencí' is…", "paralelismo", "paralelismo", "ironía", "hipérbole", "metáfora"),
		speak("Blaze", "Sus palabras, afiladas como cuchillos, cortaron el silencio."),
	)
	addVocab(db, l,
		vw("la metáfora", "metaphor", "una metáfora bella", "a beautiful metaphor", "Cora"),
		vw("la ironía", "irony", "Lo dijo con ironía.", "He said it ironically.", "Cora"),
		vw("la hipérbole", "hyperbole", "Es pura hipérbole.", "It's pure hyperbole.", "Lumora"),
		vw("el paralelismo", "parallelism", "Usó el paralelismo.", "He used parallelism.", "Lumora"),
		vw("el símil", "simile", "un símil acertado", "an apt simile", "Riko"),
	)

	// ── Regional variation ──
	s = addSkill(db, u, "Variación Regional", "España vs Latinoamérica.", "Languages", "#17A3DD", 87, 3560)
	l = addLesson(db, s, "Un Idioma, Mil Voces", 1, 32,
		char(finch, "Spanish varies by region: coche/carro, ordenador/computadora, móvil/celular, zumo/jugo; vosotros (Spain) vs ustedes (LatAm)."),
		mc("'car' in Latin America", "el ___ (LatAm)", "carro", "carro", "coche", "auto-stop", "camión"),
		mc("'computer' in Spain", "el ___ (Spain)", "ordenador", "ordenador", "computadora", "móvil", "celular"),
		fill("Fill in the blank", "'Juice': in Spain 'zumo'; in LatAm '___'.", "jugo"),
		tr("Translate this sentence", "you all speak (Spain, informal)", "vosotros habláis"),
		mc("'mobile phone' in Latin America", "el ___ (LatAm)", "celular", "celular", "móvil", "ordenador", "coche"),
		speak("Blaze", "En España cojo el coche; en México, manejo el carro."),
	)
	addVocab(db, l,
		vw("coche / carro", "car (Spain / LatAm)", "Cojo el coche.", "I take the car.", finch),
		vw("ordenador / computadora", "computer", "Uso el ordenador.", "I use the computer.", "Cora"),
		vw("móvil / celular", "mobile phone", "¿Y tu móvil?", "And your phone?", "Cora"),
		vw("zumo / jugo", "juice", "Un zumo de naranja.", "An orange juice.", "Lumora"),
		vw("vosotros / ustedes", "you (pl.)", "¿Vosotros venís?", "Are you all coming?", "Lumora"),
	)

	// ── Advanced idioms & colloquialisms ──
	s = addSkill(db, u, "Modismos Avanzados", "Native-like colloquial range.", "Sparkles", "#F5A623", 88, 3630)
	l = addLesson(db, s, "Como un Nativo", 1, 32,
		char("Cora", "Real fluency: 'estar hasta las narices' (fed up), 'no tener pelos en la lengua' (to be blunt), 'ponerse las pilas' (to get going)."),
		mc("to be fed up", "Estoy hasta las ___ .", "narices", "narices", "manos", "nubes", "uñas"),
		mc("to be blunt/outspoken", "No tiene pelos en la ___ .", "lengua", "lengua", "mano", "cabeza", "boca"),
		fill("Fill in the blank", "¡Ponte las ___ ! (get going)", "pilas"),
		tr("Translate this sentence", "It rang a bell", "Me sonó de algo"),
		mc("to throw in the towel", "Tiró la ___ .", "toalla", "toalla", "casa", "mesa", "red"),
		speak("Blaze", "Estoy hasta las narices, pero voy a ponerme las pilas."),
	)
	addVocab(db, l,
		vw("estar hasta las narices", "to be fed up", "Estoy hasta las narices.", "I'm fed up.", "Cora"),
		vw("no tener pelos en la lengua", "to be blunt", "No tiene pelos en la lengua.", "He's very blunt.", "Cora"),
		vw("ponerse las pilas", "to get one's act together", "Ponte las pilas.", "Get your act together.", "Lumora"),
		vw("tirar la toalla", "to throw in the towel", "No tires la toalla.", "Don't give up.", "Lumora"),
		vw("sonarle a alguien", "to ring a bell", "Me suena ese nombre.", "That name rings a bell.", "Riko"),
	)

	// ── Register virtuosity ──
	s = addSkill(db, u, "Registro: del Aula a la Calle", "Academic vs creative voice.", "Quote", "#00C2A8", 89, 3700)
	l = addLesson(db, s, "Dominar el Tono", 1, 32,
		char(finch, "Move fluidly between registers: academic ('cabe colegir que…'), journalistic, and creative/colloquial — choosing each consciously."),
		mc("academic 'one may infer'", "Cabe ___ que el dato es clave.", "colegir", "colegir", "molar", "pillar", "currar"),
		mc("colloquial 'to work' (Spain)", "Tengo que ___ mañana.", "currar", "currar", "deliberar", "constatar", "esgrimir"),
		fill("Fill in the blank", "Formal: 'a tenor de lo ___' (expounded). (exponer)", "expuesto"),
		tr("Translate this sentence", "The data corroborates the hypothesis (academic)", "Los datos corroboran la hipótesis"),
		mc("colloquial 'cool/great' (Spain)", "¡Esto ___ un montón!", "mola", "mola", "constata", "infiere", "deduce"),
		speak("Blaze", "En el ensayo escribo 'cabe colegir'; con amigos, 'mola un montón'."),
	)
	addVocab(db, l,
		vw("cabe colegir que", "one may infer that", "Cabe colegir que sí.", "One may infer so.", finch),
		vw("a tenor de", "in accordance with", "A tenor de lo dicho…", "In line with what was said…", "Cora"),
		vw("corroborar", "to corroborate", "Los datos lo corroboran.", "The data corroborates it.", "Cora"),
		vw("currar", "to work (colloq.)", "Curro mucho.", "I work a lot.", "Lumora"),
		vw("molar (un montón)", "to be cool (colloq.)", "Mola un montón.", "It's really cool.", "Lumora"),
	)

	// ── Specialised & abstract lexis ──
	s = addSkill(db, u, "Léxico Especializado", "philosophy, politics, science.", "BookOpen", "#6C3FC5", 90, 3770)
	l = addLesson(db, s, "Términos Abstractos", 1, 32,
		char(finch, "C2 commands precise abstraction: la cosmovisión, la sostenibilidad, el sesgo, la coyuntura, el paradigma."),
		mc("worldview", "su ___ del mundo", "cosmovisión", "cosmovisión", "coyuntura", "sesgo", "paradigma"),
		mc("(cognitive) bias", "un ___ cognitivo", "sesgo", "sesgo", "paradigma", "auge", "ocaso"),
		fill("Fill in the blank", "La ___ política actual es delicada. (situation/juncture)", "coyuntura"),
		tr("Translate this sentence", "a paradigm shift", "un cambio de paradigma"),
		mc("sustainability", "la ___ ambiental", "sostenibilidad", "sostenibilidad", "cosmovisión", "coyuntura", "sesgo"),
		speak("Blaze", "El nuevo paradigma exige replantear nuestra cosmovisión."),
	)
	addVocab(db, l,
		vw("la cosmovisión", "worldview", "una cosmovisión holística", "a holistic worldview", finch),
		vw("el sesgo", "bias", "un sesgo evidente", "an evident bias", "Cora"),
		vw("la coyuntura", "juncture/situation", "la coyuntura económica", "the economic situation", "Cora"),
		vw("el paradigma", "paradigm", "un cambio de paradigma", "a paradigm shift", "Lumora"),
		vw("la sostenibilidad", "sustainability", "la sostenibilidad del planeta", "the planet's sustainability", "Lumora"),
	)

	// ── Writing (C2) ──
	s = addSkill(db, u, "Escribir: Ensayo Crítico", "400+ words, publication-level.", "PenLine", "#00C2A8", 91, 3840)
	l = addLesson(db, s, "Ensayo de Maestría", 1, 34,
		char("Lumora", "C2 writing is publication-level: a compelling thesis, sophisticated argument, varied syntax, rhetorical control and a distinctive voice. 400+ words."),
		mc("elegant opening of an essay", "___ que el lenguaje moldea el pensamiento.", "Pocos discutirían", "Pocos discutirían", "Mola que", "Oye que", "Total que"),
		write("Write a critical essay (400+ words): 'Does language shape the way we think?' Build a nuanced argument, engage objections, and close memorably.",
			"Pocos discutirían que el lenguaje es mucho más que un mero vehículo de comunicación. La célebre hipótesis de Sapir-Whorf sostiene que la lengua que hablamos moldea, en cierto modo, nuestra percepción de la realidad. Si bien esta tesis, en su versión más radical, ha sido matizada por la lingüística contemporánea, no puede negarse que las categorías de cada idioma orientan sutilmente la atención y la memoria. Quienes defienden lo contrario suelen pasar por alto que traducir nunca es trasvasar sin pérdida: cada lengua ilumina matices que otra apenas insinúa. Dicho esto, conviene evitar el determinismo: el pensamiento también desborda las palabras y las reinventa. En definitiva, lengua y pensamiento se entretejen en un diálogo incesante en el que ninguno tiene la última palabra. Acaso ahí, en esa tensión fértil, resida la verdadera riqueza de hablar más de un idioma."),
		write("Write a critical book or film review (around 300 words) with a clear judgement and stylistic flair.",
			"La última novela del autor confirma su madurez narrativa. A través de una prosa medida y de una estructura fragmentaria, el relato indaga en la memoria y el desarraigo con una hondura poco común. Cabe destacar el pulso con que dosifica la información: nada sobra, nada falta. No obstante, el tramo final, quizá demasiado simbólico, exige al lector un esfuerzo que no siempre se ve recompensado. Con todo, estamos ante una obra valiente que se aparta de lo previsible. En suma, una lectura imprescindible para quienes busquen literatura que interpele."),
		mc("to engage an objection", "___ sostienen lo contrario olvidan un matiz.", "Quienes", "Quienes", "Que", "Cuyos", "Lo que"),
		fill("Fill in the blank", "En ___ , lengua y pensamiento se entretejen. (ultimately)", "definitiva"),
		speak("Blaze", "Acaso ahí, en esa tensión fértil, resida la verdadera riqueza."),
	)
	addVocab(db, l,
		vw("pocos discutirían que", "few would dispute that", "Pocos discutirían que sí.", "Few would dispute it.", "Lumora"),
		vw("pasar por alto", "to overlook", "Pasan por alto un matiz.", "They overlook a nuance.", "Cora"),
		vw("dicho esto", "that said", "Dicho esto, matizo.", "That said, I qualify.", "Cora"),
		vw("en suma", "in sum", "En suma, recomendable.", "In sum, recommendable.", "Lumora"),
		vw("acaso", "perhaps (literary)", "Acaso tenga razón.", "Perhaps he's right.", "Lumora"),
	)

	// ── Speaking (C2) ──
	s = addSkill(db, u, "Hablar: Oratoria de Alto Nivel", "Persuade with native flair.", "MessageCircle", "#6C3FC5", 92, 3910)
	l = addLesson(db, s, "El Arte de la Oratoria", 1, 34,
		char("Blaze", "Speak to move: open with a hook, build with rhetoric, anticipate objections and land a memorable close — effortlessly."),
		speak("Lumora", "Imaginemos por un momento un mundo en el que nadie tuviera miedo a equivocarse."),
		speak("Blaze", "No se trata, como algunos pretenden, de elegir entre libertad y seguridad."),
		speak("Blaze", "Y es que, a fin de cuentas, las grandes decisiones se toman en los pequeños gestos."),
		mc("rhetorical concession", "Se dirá que es utópico; ___, ¿quién no soñó alguna vez?", "ahora bien", "ahora bien", "o sea", "vale", "total"),
		speak("Blaze", "Permítanme, para terminar, dejarles con una pregunta incómoda."),
		speak("Blaze", "Porque, en el fondo, no hay mayor riesgo que no arriesgar nada."),
	)
	addVocab(db, l,
		vw("imaginemos por un momento", "let's imagine for a moment", "Imaginemos un mundo mejor.", "Let's imagine a better world.", "Blaze"),
		vw("no se trata de", "it's not about", "No se trata de ganar.", "It's not about winning.", "Cora"),
		vw("a fin de cuentas", "ultimately", "A fin de cuentas, importa.", "Ultimately, it matters.", "Cora"),
		vw("ahora bien", "that said / however", "Ahora bien, hay un riesgo.", "That said, there's a risk.", "Lumora"),
		vw("en el fondo", "deep down", "En el fondo, lo sabes.", "Deep down, you know it.", "Lumora"),
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

	addListening(db, "A1 · Gramática Esencial", "La rutina de Ana",
		"Ana describes her daily routine to Mira. Listen, then answer.", 5, 20,
		[]models.ListeningMatch{
			lm("me levanto", "I get up"),
			lm("desayuno", "I have breakfast"),
			lm("trabajo", "I work"),
			lm("me acuesto", "I go to bed"),
		},
		[]models.ListeningLine{
			ln("Mira", "Ana, ¿cómo es tu día?", "Ana, what is your day like?"),
			ln("Lumora", "Me levanto a las siete y desayuno café con pan.", "I get up at seven and have coffee with bread."),
			ln("Mira", "¿Y por la tarde?", "And in the afternoon?"),
			ln("Lumora", "Trabajo en una oficina y estudio español. Me acuesto a las once.", "I work in an office and study Spanish. I go to bed at eleven."),
		},
		[]models.ListeningQuestion{
			lq("When does Ana get up?", "At seven", "At seven", "At eleven", "At nine", "At six"),
			lq("What does she have for breakfast?", "Coffee with bread", "Coffee with bread", "Tea", "Eggs", "Fruit"),
			lq("What does she study?", "Spanish", "Spanish", "English", "French", "Music"),
			lq("When does she go to bed?", "At eleven", "At eleven", "At seven", "At ten", "At midnight"),
		},
	)

	addListening(db, "A1 · Gramática Esencial", "En la ciudad",
		"Riko asks for directions in town. Listen, then answer.", 6, 20,
		[]models.ListeningMatch{
			lm("¿dónde está?", "where is?"),
			lm("todo recto", "straight ahead"),
			lm("a la derecha", "to the right"),
			lm("hay", "there is"),
		},
		[]models.ListeningLine{
			ln("Riko", "Perdone, ¿dónde está el museo?", "Excuse me, where is the museum?"),
			ln("Cora", "Sigue todo recto y luego a la derecha.", "Go straight ahead and then to the right."),
			ln("Riko", "¿Hay un café cerca?", "Is there a café nearby?"),
			ln("Cora", "Sí, hay uno a la izquierda.", "Yes, there is one on the left."),
		},
		[]models.ListeningQuestion{
			lq("What is Riko looking for?", "The museum", "The museum", "The hotel", "The station", "The bank"),
			lq("Which way after going straight?", "To the right", "To the right", "To the left", "Back", "Up"),
			lq("Where is the café?", "On the left", "On the left", "On the right", "Straight ahead", "Behind"),
		},
	)

	addListening(db, "A2 · Elemental", "El fin de semana de Pablo",
		"Pablo tells Cora about his weekend. Listen, then answer.", 7, 22,
		[]models.ListeningMatch{
			lm("fui", "I went"),
			lm("comimos", "we ate"),
			lm("fue", "it was"),
			lm("volví", "I came back"),
		},
		[]models.ListeningLine{
			ln("Cora", "Pablo, ¿qué hiciste el fin de semana?", "Pablo, what did you do at the weekend?"),
			ln("Riko", "El sábado fui a la montaña con unos amigos.", "On Saturday I went to the mountains with some friends."),
			ln("Cora", "¿Y qué tal?", "And how was it?"),
			ln("Riko", "Fue genial. Comimos en un pueblo y volví el domingo por la noche.", "It was great. We ate in a village and I came back on Sunday night."),
		},
		[]models.ListeningQuestion{
			lq("Where did Pablo go on Saturday?", "To the mountains", "To the mountains", "To the beach", "To the city", "To work"),
			lq("Who did he go with?", "Some friends", "Some friends", "His family", "Alone", "His boss"),
			lq("How was the weekend?", "Great", "Great", "Boring", "Bad", "Tiring"),
			lq("When did he come back?", "Sunday night", "Sunday night", "Saturday", "Monday", "Friday"),
		},
	)

	addListening(db, "A2 · Elemental", "En la consulta",
		"Ana visits the doctor. Listen, then answer.", 8, 22,
		[]models.ListeningMatch{
			lm("me duele", "it hurts"),
			lm("la cabeza", "the head"),
			lm("fiebre", "fever"),
			lm("descansar", "to rest"),
		},
		[]models.ListeningLine{
			ln("Mira", "Buenos días. ¿Qué le pasa?", "Good morning. What's wrong?"),
			ln("Lumora", "Me duele la cabeza y tengo fiebre.", "My head hurts and I have a fever."),
			ln("Mira", "¿Desde cuándo?", "Since when?"),
			ln("Lumora", "Desde ayer. Estoy muy cansada.", "Since yesterday. I'm very tired."),
			ln("Mira", "Tiene que descansar y beber agua.", "You have to rest and drink water."),
		},
		[]models.ListeningQuestion{
			lq("What hurts?", "Her head", "Her head", "Her stomach", "Her hand", "Her back"),
			lq("What else does she have?", "A fever", "A fever", "A cough", "A cold", "Nothing"),
			lq("Since when?", "Yesterday", "Yesterday", "Last week", "Today", "An hour ago"),
			lq("What must she do?", "Rest and drink water", "Rest and drink water", "Work", "Run", "Eat more"),
		},
	)

	addListening(db, "A2 · Elemental", "En la tienda de ropa",
		"Riko shops for a jacket. Listen, then answer.", 9, 22,
		[]models.ListeningMatch{
			lm("la talla", "the size"),
			lm("probarme", "to try on"),
			lm("¿cuánto cuesta?", "how much is it?"),
			lm("la tarjeta", "the card"),
		},
		[]models.ListeningLine{
			ln("Riko", "Hola, busco una chaqueta azul.", "Hi, I'm looking for a blue jacket."),
			ln("Cora", "¿Qué talla usa?", "What size do you wear?"),
			ln("Riko", "La mediana. ¿Puedo probármela?", "Medium. Can I try it on?"),
			ln("Cora", "Claro. Cuesta cuarenta euros.", "Of course. It costs forty euros."),
			ln("Riko", "Perfecto. Pago con tarjeta.", "Perfect. I'll pay by card."),
		},
		[]models.ListeningQuestion{
			lq("What is Riko looking for?", "A blue jacket", "A blue jacket", "Red shoes", "A green hat", "Black trousers"),
			lq("What size does he want?", "Medium", "Medium", "Small", "Large", "Extra large"),
			lq("How much does it cost?", "40 euros", "40 euros", "14 euros", "50 euros", "44 euros"),
			lq("How does he pay?", "By card", "By card", "In cash", "By phone", "He doesn't"),
		},
	)

	addListening(db, "B1 · Intermedio", "La entrevista de trabajo",
		"Marta has a job interview with Mr. Finch. Listen, then answer.", 10, 26,
		[]models.ListeningMatch{
			lm("la experiencia", "experience"),
			lm("trabajé", "I worked"),
			lm("me gustaría", "I would like"),
			lm("el sueldo", "the salary"),
		},
		[]models.ListeningLine{
			ln("Professor Finch", "Buenos días. Hábleme de su experiencia.", "Good morning. Tell me about your experience."),
			ln("Mira", "Antes trabajé tres años en una empresa de marketing.", "Before, I worked three years at a marketing company."),
			ln("Professor Finch", "¿Y por qué quiere este puesto?", "And why do you want this job?"),
			ln("Mira", "Me gustaría aprender más y crecer profesionalmente.", "I would like to learn more and grow professionally."),
			ln("Professor Finch", "Perfecto. El sueldo se hablará más adelante.", "Perfect. We'll discuss the salary later."),
		},
		[]models.ListeningQuestion{
			lq("How long did Marta work before?", "Three years", "Three years", "One year", "Five years", "Two years"),
			lq("In what sector?", "Marketing", "Marketing", "Health", "Education", "Tourism"),
			lq("Why does she want the job?", "To learn and grow", "To learn and grow", "For the money", "It's near home", "Her friend works there"),
			lq("What will be discussed later?", "The salary", "The salary", "The hours", "The holidays", "The office"),
		},
	)

	addListening(db, "B1 · Intermedio", "Recuerdos de la infancia",
		"Nana remembers her childhood. Listen, then answer.", 11, 26,
		[]models.ListeningMatch{
			lm("cuando era niña", "when I was a girl"),
			lm("jugábamos", "we used to play"),
			lm("el campo", "the countryside"),
			lm("éramos felices", "we were happy"),
		},
		[]models.ListeningLine{
			ln("Pip", "Nana, ¿cómo era tu vida de pequeña?", "Nana, what was your life like as a child?"),
			ln("Nana", "Cuando era niña, vivía en el campo con mis abuelos.", "When I was a girl, I lived in the countryside with my grandparents."),
			ln("Pip", "¿Y qué hacíais?", "And what did you all do?"),
			ln("Nana", "Jugábamos fuera todo el día. No teníamos móviles, pero éramos felices.", "We played outside all day. We had no phones, but we were happy."),
		},
		[]models.ListeningQuestion{
			lq("Where did Nana live as a child?", "In the countryside", "In the countryside", "In a city", "By the sea", "Abroad"),
			lq("Who did she live with?", "Her grandparents", "Her grandparents", "Her parents", "Her aunt", "Alone"),
			lq("What did they do all day?", "Played outside", "Played outside", "Watched TV", "Studied", "Worked"),
			lq("How were they?", "Happy", "Happy", "Bored", "Sad", "Tired"),
		},
	)

	addListening(db, "B1 · Intermedio", "El medio ambiente",
		"Zephyr and Cora discuss the environment. Listen, then answer.", 12, 26,
		[]models.ListeningMatch{
			lm("la contaminación", "pollution"),
			lm("reciclar", "to recycle"),
			lm("creo que", "I think that"),
			lm("el transporte público", "public transport"),
		},
		[]models.ListeningLine{
			ln("Zephyr", "Creo que la contaminación es el mayor problema de las ciudades.", "I think pollution is the biggest problem in cities."),
			ln("Cora", "Estoy de acuerdo. Por eso es importante que reciclemos.", "I agree. That's why it's important that we recycle."),
			ln("Zephyr", "Sí, y deberíamos usar más el transporte público.", "Yes, and we should use public transport more."),
			ln("Cora", "Sin embargo, mucha gente prefiere el coche.", "However, many people prefer the car."),
		},
		[]models.ListeningQuestion{
			lq("What is the biggest problem, says Zephyr?", "Pollution", "Pollution", "Noise", "Traffic jams", "Rubbish"),
			lq("What does Cora say is important?", "To recycle", "To recycle", "To drive", "To save money", "To plant trees"),
			lq("What should people use more?", "Public transport", "Public transport", "Bikes only", "Their cars", "Taxis"),
			lq("What do many people prefer?", "The car", "The car", "The bus", "Walking", "The train"),
		},
	)

	addListening(db, "B2 · Avanzado", "Debate: ¿teletrabajo?",
		"Zephyr and Mira debate remote work. Listen, then answer.", 13, 30,
		[]models.ListeningMatch{
			lm("desde mi punto de vista", "from my point of view"),
			lm("la productividad", "productivity"),
			lm("el aislamiento", "isolation"),
			lm("un modelo híbrido", "a hybrid model"),
		},
		[]models.ListeningLine{
			ln("Zephyr", "Desde mi punto de vista, el teletrabajo aumenta la productividad.", "From my point of view, remote work increases productivity."),
			ln("Mira", "Es cierto, pero también puede provocar aislamiento.", "True, but it can also cause isolation."),
			ln("Zephyr", "Si las empresas ofrecieran apoyo, no habría ese problema.", "If companies offered support, there wouldn't be that problem."),
			ln("Mira", "Por eso creo que lo ideal sería un modelo híbrido.", "That's why I think the ideal would be a hybrid model."),
		},
		[]models.ListeningQuestion{
			lq("What does Zephyr say remote work increases?", "Productivity", "Productivity", "Costs", "Traffic", "Stress"),
			lq("What problem does Mira mention?", "Isolation", "Isolation", "Noise", "Low pay", "Long hours"),
			lq("Under what condition would the problem disappear?", "If companies offered support", "If companies offered support", "If salaries rose", "If offices closed", "If hours were shorter"),
			lq("What does Mira propose?", "A hybrid model", "A hybrid model", "Office only", "Remote only", "A four-day week"),
		},
	)

	addListening(db, "B2 · Avanzado", "Noticias de la mañana",
		"A short news bulletin. Listen, then answer.", 14, 30,
		[]models.ListeningMatch{
			lm("según fuentes", "according to sources"),
			lm("el gobierno", "the government"),
			lm("ha anunciado", "has announced"),
			lm("las medidas", "the measures"),
		},
		[]models.ListeningLine{
			ln("Professor Finch", "Buenos días. El gobierno ha anunciado nuevas medidas contra la contaminación.", "Good morning. The government has announced new measures against pollution."),
			ln("Professor Finch", "Según fuentes oficiales, se reducirá el tráfico en el centro.", "According to official sources, traffic in the centre will be reduced."),
			ln("Professor Finch", "Además, cabe destacar que se invertirá en transporte público.", "Also, it's worth noting that there will be investment in public transport."),
			ln("Professor Finch", "Los expertos consideran que estas medidas son necesarias.", "Experts consider these measures necessary."),
		},
		[]models.ListeningQuestion{
			lq("What has the government announced?", "Measures against pollution", "Measures against pollution", "Tax cuts", "New schools", "A holiday"),
			lq("What will be reduced in the centre?", "Traffic", "Traffic", "Prices", "Crime", "Noise"),
			lq("What will they invest in?", "Public transport", "Public transport", "Roads", "Hospitals", "Parks"),
			lq("What do experts think?", "The measures are necessary", "The measures are necessary", "They are useless", "They are too costly", "They are too late"),
		},
	)

	addListening(db, "B2 · Avanzado", "Si pudiera volver atrás",
		"Nana reflects on her life choices. Listen, then answer.", 15, 30,
		[]models.ListeningMatch{
			lm("si pudiera", "if I could"),
			lm("habría estudiado", "I would have studied"),
			lm("me arrepiento", "I regret"),
			lm("valió la pena", "it was worth it"),
		},
		[]models.ListeningLine{
			ln("Pip", "Nana, ¿cambiarías algo de tu vida?", "Nana, would you change anything in your life?"),
			ln("Nana", "Si pudiera volver atrás, habría estudiado música.", "If I could go back, I would have studied music."),
			ln("Pip", "¿Te arrepientes entonces?", "Do you regret it then?"),
			ln("Nana", "No mucho. Aunque cometí errores, todo valió la pena.", "Not much. Although I made mistakes, it was all worth it."),
		},
		[]models.ListeningQuestion{
			lq("What would Nana have studied?", "Music", "Music", "Medicine", "Law", "Art"),
			lq("Does she regret her life?", "Not much", "Not much", "Very much", "Completely", "She doesn't say"),
			lq("What does she admit she made?", "Mistakes", "Mistakes", "Money", "Friends", "Plans"),
			lq("What does she conclude?", "It was worth it", "It was worth it", "It was a waste", "She'd change everything", "She's unsure"),
		},
	)

	addListening(db, "C1 · Superior", "Tertulia: la inteligencia artificial",
		"Two experts discuss AI on a talk show. Listen, then answer.", 16, 32,
		[]models.ListeningMatch{
			lm("si bien es cierto", "while it's true"),
			lm("plantea dilemas", "raises dilemmas"),
			lm("cabe matizar", "it should be nuanced"),
			lm("a largo plazo", "in the long term"),
		},
		[]models.ListeningLine{
			ln("Professor Finch", "La inteligencia artificial plantea dilemas éticos que no podemos ignorar.", "AI raises ethical dilemmas we can't ignore."),
			ln("Zephyr", "Si bien es cierto que entraña riesgos, sus beneficios son enormes.", "While it's true it entails risks, its benefits are huge."),
			ln("Professor Finch", "Cabe matizar que todo depende de cómo se regule.", "It should be nuanced that it all depends on how it's regulated."),
			ln("Zephyr", "Exacto. A largo plazo, lo ideal sería un marco común.", "Exactly. In the long term, the ideal would be a common framework."),
		},
		[]models.ListeningQuestion{
			lq("What does Finch say AI raises?", "Ethical dilemmas", "Ethical dilemmas", "Profits", "Jobs", "Taxes"),
			lq("What does Zephyr emphasise?", "Its benefits are huge", "Its benefits are huge", "It's useless", "It's cheap", "It's simple"),
			lq("On what does it all depend, says Finch?", "How it's regulated", "How it's regulated", "Who pays", "The country", "The year"),
			lq("What would be ideal long term?", "A common framework", "A common framework", "A total ban", "No rules", "Higher prices"),
		},
	)

	addListening(db, "C1 · Superior", "Conferencia sobre la lectura",
		"A short lecture on the value of reading. Listen, then answer.", 17, 32,
		[]models.ListeningMatch{
			lm("fomentar", "to foster"),
			lm("el pensamiento crítico", "critical thinking"),
			lm("en resumidas cuentas", "in short"),
			lm("imprescindible", "essential"),
		},
		[]models.ListeningLine{
			ln("Professor Finch", "Hoy quisiera reflexionar sobre el papel de la lectura en la sociedad.", "Today I'd like to reflect on the role of reading in society."),
			ln("Professor Finch", "Leer no solo amplía el vocabulario, sino que fomenta el pensamiento crítico.", "Reading not only widens vocabulary, but fosters critical thinking."),
			ln("Professor Finch", "Quienes leen con regularidad desarrollan mayor empatía.", "Those who read regularly develop greater empathy."),
			ln("Professor Finch", "En resumidas cuentas, la lectura es imprescindible para una mente libre.", "In short, reading is essential for a free mind."),
		},
		[]models.ListeningQuestion{
			lq("What is the lecture about?", "The role of reading", "The role of reading", "Writing essays", "Learning maths", "Public speaking"),
			lq("What does reading foster, besides vocabulary?", "Critical thinking", "Critical thinking", "Wealth", "Speed", "Memory only"),
			lq("What do regular readers develop?", "Greater empathy", "Greater empathy", "Worse eyesight", "More stress", "Less time"),
			lq("How does he sum up reading?", "Essential for a free mind", "Essential for a free mind", "A waste of time", "Only for students", "Outdated"),
		},
	)

	addListening(db, "C1 · Superior", "Entrevista a una escritora",
		"An author is interviewed about her new novel. Listen, then answer.", 18, 32,
		[]models.ListeningMatch{
			lm("la trama", "the plot"),
			lm("me inspiré en", "I was inspired by"),
			lm("los prejuicios", "prejudices"),
			lm("dicho esto", "that said"),
		},
		[]models.ListeningLine{
			ln("Cora", "Su nueva novela aborda temas muy actuales. ¿De dónde surge la trama?", "Your new novel tackles very current themes. Where does the plot come from?"),
			ln("Mira", "Me inspiré en historias reales sobre la migración.", "I was inspired by real stories about migration."),
			ln("Cora", "¿Pretende denunciar algo?", "Do you intend to denounce something?"),
			ln("Mira", "Quiero retratar los prejuicios sin juzgar. Dicho esto, cada lector sacará sus conclusiones.", "I want to portray prejudices without judging. That said, each reader will draw their own conclusions."),
		},
		[]models.ListeningQuestion{
			lq("What inspired the plot?", "Real migration stories", "Real migration stories", "A dream", "Her childhood", "A film"),
			lq("What does the author want to portray?", "Prejudices without judging", "Prejudices without judging", "Wealth", "War heroes", "Romance only"),
			lq("Who will draw their own conclusions?", "Each reader", "Each reader", "The critics", "The author", "No one"),
			lq("How would you describe the novel's themes?", "Very current", "Very current", "Historical only", "Light", "Unclear"),
		},
	)

	addListening(db, "C2 · Maestría", "Debate: ¿libertad o seguridad?",
		"A sharp philosophical debate. Listen for nuance, then answer.", 19, 34,
		[]models.ListeningMatch{
			lm("a fin de cuentas", "ultimately"),
			lm("falso dilema", "false dilemma"),
			lm("ahora bien", "that said"),
			lm("no se trata de", "it's not about"),
		},
		[]models.ListeningLine{
			ln("Zephyr", "No se trata, como algunos pretenden, de elegir entre libertad y seguridad.", "It's not about, as some claim, choosing between freedom and security."),
			ln("Professor Finch", "Ahora bien, en situaciones extremas, ¿no priorizaríamos la seguridad?", "That said, in extreme situations, wouldn't we prioritise security?"),
			ln("Zephyr", "Eso es un falso dilema. A fin de cuentas, sin libertad no hay seguridad que valga.", "That's a false dilemma. Ultimately, without freedom no security is worth anything."),
			ln("Professor Finch", "Permítame discrepar, aunque reconozco la fuerza de su argumento.", "Allow me to disagree, though I acknowledge the strength of your argument."),
		},
		[]models.ListeningQuestion{
			lq("What does Zephyr say it is NOT about?", "Choosing freedom or security", "Choosing freedom or security", "Money", "Politics", "History"),
			lq("What does Zephyr call the choice?", "A false dilemma", "A false dilemma", "A fair point", "A law", "A tradition"),
			lq("What does Finch do at the end?", "Politely disagrees", "Politely disagrees", "Fully agrees", "Stays silent", "Changes topic"),
			lq("What is the tone of the exchange?", "Sharp but respectful", "Sharp but respectful", "Hostile", "Indifferent", "Comic"),
		},
	)

	addListening(db, "C2 · Maestría", "Tertulia sobre el humor",
		"Two critics dissect irony and satire. Listen, then answer.", 20, 34,
		[]models.ListeningMatch{
			lm("el sarcasmo", "sarcasm"),
			lm("dar en el clavo", "to hit the nail on the head"),
			lm("a costa de", "at the expense of"),
			lm("doble sentido", "double meaning"),
		},
		[]models.ListeningLine{
			ln("Cora", "La buena sátira da en el clavo sin caer en el insulto fácil.", "Good satire hits the nail on the head without resorting to cheap insults."),
			ln("Mira", "Cierto, aunque el sarcasmo, mal usado, hiere a costa de los débiles.", "True, although sarcasm, misused, wounds at the expense of the weak."),
			ln("Cora", "De ahí que el humor inteligente juegue con el doble sentido, no con la crueldad.", "Hence intelligent humour plays with double meaning, not cruelty."),
			ln("Mira", "En el fondo, reírse de uno mismo es la forma más alta de lucidez.", "Deep down, laughing at oneself is the highest form of clarity."),
		},
		[]models.ListeningQuestion{
			lq("What does good satire do, per Cora?", "Hits the point without cheap insults", "Hits the point without cheap insults", "Always offends", "Avoids politics", "Stays silent"),
			lq("What does misused sarcasm do?", "Wounds the weak", "Wounds the weak", "Educates", "Heals", "Bores"),
			lq("What does intelligent humour play with?", "Double meaning", "Double meaning", "Cruelty", "Volume", "Speed"),
			lq("What does Mira call the highest clarity?", "Laughing at oneself", "Laughing at oneself", "Winning debates", "Staying serious", "Mocking others"),
		},
	)

	addListening(db, "C2 · Maestría", "Acentos del español",
		"A linguist discusses regional variety. Listen, then answer.", 21, 34,
		[]models.ListeningMatch{
			lm("la riqueza", "the richness"),
			lm("no hay uno mejor", "none is better"),
			lm("el seseo", "the 'seseo'"),
			lm("el prestigio", "prestige"),
		},
		[]models.ListeningLine{
			ln("Professor Finch", "El español no es un bloque uniforme, sino un mosaico de variedades.", "Spanish is not a uniform block, but a mosaic of varieties."),
			ln("Lumora", "¿Y hay alguna más correcta que otra?", "And is any more correct than another?"),
			ln("Professor Finch", "En absoluto. No hay uno mejor; el prestigio es una cuestión social, no lingüística.", "Not at all. None is better; prestige is a social matter, not a linguistic one."),
			ln("Lumora", "Entonces, la diversidad es, en realidad, su mayor riqueza.", "So diversity is, in fact, its greatest richness."),
		},
		[]models.ListeningQuestion{
			lq("How does Finch describe Spanish?", "A mosaic of varieties", "A mosaic of varieties", "A uniform block", "A dying language", "A simple code"),
			lq("Is one variety more correct?", "No, none is better", "No, none is better", "Yes, Spain's", "Yes, Mexico's", "He won't say"),
			lq("What kind of matter is prestige?", "Social, not linguistic", "Social, not linguistic", "Purely grammatical", "Legal", "Economic"),
			lq("What is the language's greatest richness?", "Its diversity", "Its diversity", "Its rules", "Its age", "Its difficulty"),
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

	addReading(db, "A1 · Gramática Esencial", "Un día de Pablo",
		"Read about Pablo's day, then answer.", 5, 20,
		[]models.ReadingLine{
			rl("Me llamo Pablo y soy de México.", "My name is Pablo and I'm from Mexico."),
			rl("Por la mañana me levanto a las siete y desayuno.", "In the morning I get up at seven and have breakfast."),
			rl("Trabajo en una oficina y como a la una.", "I work in an office and eat at one."),
			rl("Por la tarde estudio español. Me gusta mucho.", "In the afternoon I study Spanish. I like it a lot."),
		},
		[]models.ReadingQuestion{
			rq("Where is Pablo from?", "Mexico", "Mexico", "Spain", "Peru", "Chile"),
			rq("What time does he get up?", "Seven", "Seven", "One", "Eight", "Nine"),
			rq("Where does he work?", "In an office", "In an office", "At home", "In a shop", "At school"),
			rq("What does he study?", "Spanish", "Spanish", "English", "Maths", "Music"),
		},
	)

	addReading(db, "A1 · Gramática Esencial", "La familia de Marta",
		"Read about Marta's family, then answer.", 6, 20,
		[]models.ReadingLine{
			rl("Mi familia es pequeña.", "My family is small."),
			rl("Tengo un hermano y una hermana.", "I have a brother and a sister."),
			rl("Mi padre es alto y mi madre es simpática.", "My father is tall and my mother is nice."),
			rl("Vivimos en una casa con un perro.", "We live in a house with a dog."),
		},
		[]models.ReadingQuestion{
			rq("How many siblings does Marta have?", "Two", "Two", "One", "Three", "None"),
			rq("How is her father described?", "Tall", "Tall", "Short", "Funny", "Old"),
			rq("What pet do they have?", "A dog", "A dog", "A cat", "A bird", "A fish"),
			rq("How is the family described?", "Small", "Small", "Big", "Loud", "New"),
		},
	)

	addReading(db, "A2 · Elemental", "Mis vacaciones",
		"Read about Lucía's holiday, then answer.", 7, 22,
		[]models.ReadingLine{
			rl("El verano pasado fui a Sevilla con mi familia.", "Last summer I went to Seville with my family."),
			rl("Visitamos la catedral y comimos tapas todos los días.", "We visited the cathedral and ate tapas every day."),
			rl("Hizo mucho calor, pero lo pasamos genial.", "It was very hot, but we had a great time."),
			rl("Volvimos a casa muy contentos.", "We came back home very happy."),
		},
		[]models.ReadingQuestion{
			rq("Where did Lucía go?", "Seville", "Seville", "Madrid", "Barcelona", "Valencia"),
			rq("Who did she travel with?", "Her family", "Her family", "Her friends", "Alone", "Her class"),
			rq("What did they eat?", "Tapas", "Tapas", "Pizza", "Sushi", "Burgers"),
			rq("What was the weather like?", "Very hot", "Very hot", "Cold", "Rainy", "Snowy"),
		},
	)

	addReading(db, "A2 · Elemental", "Una reseña",
		"Read the restaurant review, then answer.", 8, 22,
		[]models.ReadingLine{
			rl("Ayer cené en el restaurante El Sol.", "Yesterday I had dinner at El Sol restaurant."),
			rl("Pedí pescado y estaba muy rico.", "I ordered fish and it was very tasty."),
			rl("El camarero fue muy amable, pero la música estaba demasiado alta.", "The waiter was very friendly, but the music was too loud."),
			rl("En general, lo recomiendo. Volveré pronto.", "Overall, I recommend it. I'll come back soon."),
		},
		[]models.ReadingQuestion{
			rq("When did the writer go?", "Yesterday", "Yesterday", "Last week", "Today", "Last year"),
			rq("What did they order?", "Fish", "Fish", "Meat", "Pasta", "Soup"),
			rq("What was the problem?", "The music was too loud", "The music was too loud", "The food was cold", "The waiter was rude", "It was expensive"),
			rq("Will they return?", "Yes, soon", "Yes, soon", "Never", "Not sure", "Only once"),
		},
	)

	addReading(db, "A2 · Elemental", "Planes para el verano",
		"Read about Marco's summer plans, then answer.", 9, 22,
		[]models.ReadingLine{
			rl("Este verano voy a viajar a México.", "This summer I'm going to travel to Mexico."),
			rl("Voy a visitar las pirámides y a practicar mi español.", "I'm going to visit the pyramids and practise my Spanish."),
			rl("Mi hermana va a venir conmigo.", "My sister is going to come with me."),
			rl("Vamos a quedarnos dos semanas.", "We're going to stay two weeks."),
		},
		[]models.ReadingQuestion{
			rq("Where is Marco going?", "Mexico", "Mexico", "Spain", "Peru", "Cuba"),
			rq("What will he visit?", "The pyramids", "The pyramids", "The beach", "A museum", "The mountains"),
			rq("Who is coming with him?", "His sister", "His sister", "His brother", "His friend", "No one"),
			rq("How long will they stay?", "Two weeks", "Two weeks", "One week", "A month", "Three days"),
		},
	)

	addReading(db, "B1 · Intermedio", "La tecnología en nuestra vida",
		"Read the opinion article, then answer.", 10, 26,
		[]models.ReadingLine{
			rl("Hoy en día, la tecnología está presente en todo.", "Nowadays, technology is present in everything."),
			rl("Por un lado, nos permite comunicarnos y trabajar desde casa.", "On one hand, it lets us communicate and work from home."),
			rl("Sin embargo, muchos creen que pasamos demasiado tiempo con el móvil.", "However, many think we spend too much time on our phones."),
			rl("En mi opinión, la tecnología es útil si la usamos con moderación.", "In my opinion, technology is useful if we use it in moderation."),
		},
		[]models.ReadingQuestion{
			rq("What does technology let us do, on one hand?", "Work from home", "Work from home", "Sleep more", "Travel free", "Earn more"),
			rq("What is the concern?", "Too much phone time", "Too much phone time", "It's expensive", "It's slow", "It's boring"),
			rq("What is the writer's opinion?", "Useful in moderation", "Useful in moderation", "Always bad", "Always good", "Useless"),
			rq("Which connector introduces the contrast?", "Sin embargo", "Sin embargo", "Por un lado", "En mi opinión", "Hoy en día"),
		},
	)

	addReading(db, "B1 · Intermedio", "Un viaje inolvidable",
		"Read the travel story, then answer.", 11, 26,
		[]models.ReadingLine{
			rl("Hacía mucho calor cuando llegamos a Granada.", "It was very hot when we arrived in Granada."),
			rl("Mientras paseábamos, descubrimos la Alhambra.", "While we were walking, we discovered the Alhambra."),
			rl("Nunca había visto un lugar tan bonito.", "I had never seen such a beautiful place."),
			rl("Aunque estábamos cansados, decidimos quedarnos un día más.", "Although we were tired, we decided to stay one more day."),
		},
		[]models.ReadingQuestion{
			rq("What was the weather like on arrival?", "Very hot", "Very hot", "Cold", "Rainy", "Windy"),
			rq("What did they discover while walking?", "The Alhambra", "The Alhambra", "A market", "A beach", "A museum"),
			rq("Had the writer seen such a place before?", "Never", "Never", "Many times", "Once", "Twice"),
			rq("What did they decide despite being tired?", "To stay one more day", "To stay one more day", "To go home", "To sleep", "To complain"),
		},
	)

	addReading(db, "B1 · Intermedio", "Buscar trabajo",
		"Read the job advice, then answer.", 12, 26,
		[]models.ReadingLine{
			rl("Cuando buscas trabajo, es importante que prepares bien tu currículum.", "When you look for work, it's important to prepare your CV well."),
			rl("Antes de la entrevista, investiga la empresa.", "Before the interview, research the company."),
			rl("Te recomiendo que llegues pronto y que hagas preguntas.", "I recommend that you arrive early and ask questions."),
			rl("No creo que sea fácil, pero con esfuerzo lo conseguirás.", "I don't think it's easy, but with effort you'll get it."),
		},
		[]models.ReadingQuestion{
			rq("What is important to prepare?", "Your CV", "Your CV", "Your car", "Your lunch", "Your clothes only"),
			rq("What should you do before the interview?", "Research the company", "Research the company", "Sleep", "Call a friend", "Buy a gift"),
			rq("What does the writer recommend?", "Arrive early and ask questions", "Arrive early and ask questions", "Arrive late", "Stay silent", "Talk about money"),
			rq("What is the writer's view on difficulty?", "Not easy but achievable", "Not easy but achievable", "Impossible", "Very easy", "Pointless"),
		},
	)

	addReading(db, "B2 · Avanzado", "El impacto del turismo",
		"Read the opinion essay, then answer.", 13, 30,
		[]models.ReadingLine{
			rl("El turismo masivo se ha convertido en un arma de doble filo.", "Mass tourism has become a double-edged sword."),
			rl("Por un lado, genera empleo y riqueza para muchas regiones.", "On one hand, it generates jobs and wealth for many regions."),
			rl("Sin embargo, también provoca la subida de los alquileres y daña el medio ambiente.", "However, it also drives up rents and damages the environment."),
			rl("Cabe destacar que, si no se regulara, algunas ciudades se volverían invivibles.", "It's worth noting that, if it weren't regulated, some cities would become unlivable."),
		},
		[]models.ReadingQuestion{
			rq("How is mass tourism described?", "A double-edged sword", "A double-edged sword", "A disaster", "A blessing", "A mystery"),
			rq("What benefit is mentioned?", "Jobs and wealth", "Jobs and wealth", "Cheaper rents", "Cleaner air", "Less traffic"),
			rq("What problem does it cause?", "Higher rents", "Higher rents", "Lower wages", "Fewer tourists", "More rain"),
			rq("What would happen without regulation?", "Cities would become unlivable", "Cities would become unlivable", "Tourism would stop", "Rents would fall", "Nothing"),
		},
	)

	addReading(db, "B2 · Avanzado", "La salud mental",
		"Read the article on mental health, then answer.", 14, 30,
		[]models.ReadingLine{
			rl("Durante mucho tiempo, la salud mental fue un tema tabú.", "For a long time, mental health was a taboo subject."),
			rl("Hoy, aunque todavía existe estigma, se habla de ello más abiertamente.", "Today, although stigma still exists, it is talked about more openly."),
			rl("Los expertos recomiendan que pidamos ayuda sin vergüenza.", "Experts recommend that we ask for help without shame."),
			rl("De hecho, cuidar la mente es tan importante como cuidar el cuerpo.", "In fact, caring for the mind is as important as caring for the body."),
		},
		[]models.ReadingQuestion{
			rq("How was mental health seen for a long time?", "A taboo", "A taboo", "A joke", "A luxury", "A myth"),
			rq("How is it discussed today?", "More openly", "More openly", "Never", "Less than before", "Only by doctors"),
			rq("What do experts recommend?", "Asking for help without shame", "Asking for help without shame", "Ignoring it", "Hiding it", "Waiting"),
			rq("Caring for the mind is as important as…?", "Caring for the body", "Caring for the body", "Earning money", "Studying", "Sleeping"),
		},
	)

	addReading(db, "B2 · Avanzado", "La globalización",
		"Read the essay on globalization, then answer.", 15, 30,
		[]models.ReadingLine{
			rl("La globalización ha conectado el mundo como nunca antes.", "Globalization has connected the world like never before."),
			rl("Gracias a ella, podemos acceder a productos y culturas de todo el planeta.", "Thanks to it, we can access products and cultures from all over the planet."),
			rl("No obstante, algunos critican que favorece a las grandes empresas.", "Nevertheless, some criticize that it favours large companies."),
			rl("Si fuéramos capaces de repartir mejor la riqueza, sus beneficios serían mayores.", "If we were able to distribute wealth better, its benefits would be greater."),
		},
		[]models.ReadingQuestion{
			rq("What has globalization done?", "Connected the world", "Connected the world", "Divided the world", "Slowed the world", "Ended trade"),
			rq("What can we access thanks to it?", "Products and cultures", "Products and cultures", "Only money", "Free travel", "Cheap homes"),
			rq("What do some critics say?", "It favours large companies", "It favours large companies", "It helps the poor", "It's too slow", "It's harmless"),
			rq("What would make its benefits greater?", "Better wealth distribution", "Better wealth distribution", "More companies", "Less trade", "Higher prices"),
		},
	)

	addReading(db, "C1 · Superior", "Editorial: la era digital",
		"Read the editorial, noting its tone, then answer.", 16, 32,
		[]models.ReadingLine{
			rl("Vivimos inmersos en una vorágine tecnológica de la que resulta difícil sustraerse.", "We live immersed in a technological whirlwind from which it's hard to escape."),
			rl("Si bien las redes nos acercan, también han erosionado, paradójicamente, la conversación pausada.", "While networks bring us closer, they have, paradoxically, eroded unhurried conversation."),
			rl("No se trata de demonizar el progreso, sino de aprender a convivir con él.", "It's not about demonizing progress, but learning to live with it."),
			rl("En definitiva, la tecnología debería estar al servicio del ser humano, y no al revés.", "Ultimately, technology should serve humans, not the other way around."),
		},
		[]models.ReadingQuestion{
			rq("How does the author describe our era?", "A technological whirlwind", "A technological whirlwind", "A golden age", "A quiet time", "A disaster"),
			rq("What have networks paradoxically eroded?", "Unhurried conversation", "Unhurried conversation", "Our memory", "The economy", "Friendship entirely"),
			rq("What is the author's stance on progress?", "Not to demonize but coexist with it", "Not to demonize but coexist with it", "To reject it", "To worship it", "To ignore it"),
			rq("What is the editorial's tone?", "Reflective and balanced", "Reflective and balanced", "Furious", "Indifferent", "Comic"),
		},
	)

	addReading(db, "C1 · Superior", "Fragmento literario",
		"Read the literary excerpt, then answer.", 17, 32,
		[]models.ReadingLine{
			rl("Cuando regresó al pueblo, todo le pareció más pequeño, como si los años lo hubieran encogido.", "When he returned to the village, everything seemed smaller, as if the years had shrunk it."),
			rl("Las calles que antaño recorría con entusiasmo guardaban ahora un silencio extraño.", "The streets he once walked with enthusiasm now held a strange silence."),
			rl("Comprendió, no sin cierta melancolía, que el lugar no había cambiado: era él quien ya no era el mismo.", "He understood, not without some melancholy, that the place hadn't changed: it was he who was no longer the same."),
		},
		[]models.ReadingQuestion{
			rq("How did the village seem on his return?", "Smaller", "Smaller", "Bigger", "Noisier", "Brand new"),
			rq("What did the streets hold now?", "A strange silence", "A strange silence", "Loud music", "Crowds", "Markets"),
			rq("What did he finally understand?", "He had changed, not the place", "He had changed, not the place", "The place had changed", "Nothing had meaning", "He should leave"),
			rq("What feeling pervades the passage?", "Melancholy", "Melancholy", "Joy", "Anger", "Fear"),
		},
	)

	addReading(db, "C1 · Superior", "La identidad cultural",
		"Read the essay on cultural identity, then answer.", 18, 32,
		[]models.ReadingLine{
			rl("La identidad cultural no es un bloque inmutable, sino un proceso en constante construcción.", "Cultural identity is not an immutable block, but a process in constant construction."),
			rl("Cada generación reinterpreta sus tradiciones a la luz de nuevos contextos.", "Each generation reinterprets its traditions in the light of new contexts."),
			rl("Quienes defienden una identidad 'pura' suelen olvidar que toda cultura es, en esencia, mestiza.", "Those who defend a 'pure' identity often forget that every culture is, in essence, mixed."),
			rl("Por consiguiente, abrazar la diversidad no debilita lo propio: lo enriquece.", "Consequently, embracing diversity does not weaken what is one's own: it enriches it."),
		},
		[]models.ReadingQuestion{
			rq("How is cultural identity described?", "A process in constant construction", "A process in constant construction", "An immutable block", "A modern invention", "A lost cause"),
			rq("What does each generation do with traditions?", "Reinterprets them", "Reinterprets them", "Abandons them", "Freezes them", "Forgets them"),
			rq("What do defenders of a 'pure' identity forget?", "Every culture is mixed", "Every culture is mixed", "Culture is useless", "Traditions matter", "Languages die"),
			rq("What does embracing diversity do, per the author?", "Enriches one's own culture", "Enriches one's own culture", "Weakens it", "Erases it", "Has no effect"),
		},
	)

	addReading(db, "C2 · Maestría", "Elogio de la duda",
		"Read the philosophical essay, attending to its rhetoric, then answer.", 19, 34,
		[]models.ReadingLine{
			rl("Vivimos en una época que premia la certeza y desconfía de quien titubea.", "We live in an age that rewards certainty and distrusts those who hesitate."),
			rl("Y sin embargo, acaso sea la duda, y no la convicción ciega, el verdadero motor del pensamiento.", "And yet, perhaps it is doubt, not blind conviction, that is the true engine of thought."),
			rl("Quien nunca duda no piensa: se limita a repetir certezas heredadas.", "Whoever never doubts does not think: they merely repeat inherited certainties."),
			rl("Dudar, lejos de ser debilidad, es el primer acto de una mente libre.", "To doubt, far from being weakness, is the first act of a free mind."),
		},
		[]models.ReadingQuestion{
			rq("What does the age reward, per the author?", "Certainty", "Certainty", "Doubt", "Silence", "Wealth"),
			rq("What does the author call the true engine of thought?", "Doubt", "Doubt", "Blind conviction", "Memory", "Fear"),
			rq("What does someone who never doubts do?", "Repeats inherited certainties", "Repeats inherited certainties", "Thinks deeply", "Stays silent", "Learns fast"),
			rq("How does the author reframe doubting?", "The first act of a free mind", "The first act of a free mind", "A weakness", "A waste", "A sin"),
		},
	)

	addReading(db, "C2 · Maestría", "Columna: la prisa",
		"Read this ironic opinion column, attending to tone, then answer.", 20, 34,
		[]models.ReadingLine{
			rl("Enhorabuena: hemos inventado mil formas de ahorrar tiempo y ninguna de disfrutarlo.", "Congratulations: we've invented a thousand ways to save time and none to enjoy it."),
			rl("Corremos para llegar antes, sin recordar ya adónde íbamos ni por qué.", "We rush to arrive sooner, no longer remembering where we were going or why."),
			rl("Al parecer, el ocio se ha vuelto sospechoso, casi un delito que conviene justificar.", "Apparently, leisure has become suspect, almost a crime one had better justify."),
			rl("Quizá detenerse, hoy, sea el más subversivo de los actos.", "Perhaps stopping, today, is the most subversive of acts."),
		},
		[]models.ReadingQuestion{
			rq("What is the column's tone?", "Ironic", "Ironic", "Sincere praise", "Neutral", "Academic"),
			rq("What have we invented, per the writer?", "Ways to save time, not enjoy it", "Ways to save time, not enjoy it", "Ways to travel", "Ways to earn money", "Ways to rest"),
			rq("How is leisure now seen?", "As suspect, almost a crime", "As suspect, almost a crime", "As sacred", "As normal", "As free"),
			rq("What does the writer call subversive today?", "Stopping", "Stopping", "Working", "Spending", "Competing"),
		},
	)

	addReading(db, "C2 · Maestría", "Crítica literaria",
		"Read the literary review, then answer.", 21, 34,
		[]models.ReadingLine{
			rl("La novela, de prosa contenida y estructura fragmentaria, indaga en la memoria y el desarraigo.", "The novel, of restrained prose and fragmentary structure, probes memory and rootlessness."),
			rl("Cabe destacar el pulso con que el autor dosifica la información: nada sobra, nada falta.", "It's worth noting the steady hand with which the author measures out information: nothing is excessive, nothing lacking."),
			rl("No obstante, el desenlace, quizá demasiado simbólico, exige un esfuerzo no siempre recompensado.", "However, the ending, perhaps too symbolic, demands an effort not always rewarded."),
			rl("Con todo, estamos ante una obra valiente que se aparta de lo previsible.", "All in all, this is a bold work that departs from the predictable."),
		},
		[]models.ReadingQuestion{
			rq("How is the novel's prose described?", "Restrained", "Restrained", "Excessive", "Careless", "Comic"),
			rq("What does the reviewer praise?", "How information is measured out", "How information is measured out", "The cover", "The length", "The price"),
			rq("What is the reviewer's reservation?", "The ending is too symbolic", "The ending is too symbolic", "The start is dull", "It's too short", "It's too cheap"),
			rq("What is the overall judgement?", "A bold, unpredictable work", "A bold, unpredictable work", "A failure", "Mediocre", "Forgettable"),
		},
	)
}
