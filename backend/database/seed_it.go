package database

import "gorm.io/gorm"

// ===== Italian (Italiano) course, A1 → C2 =====
//
// Standard Italian, CEFR-aligned (the scale used by CILS, CELI, PLIDA and
// ROMA TRE). The order follows the guide's principles: sounds and politeness
// first; essere/avere early because they unlock the past tenses; gender and
// the article always learned together with the noun; spoken past tenses
// (passato prossimo + imperfetto) at A2 and the literary passato remoto only
// at C1; the subjunctive introduced at B1 rather than left to become a block.
//
// Italian is written in the Latin alphabet, so lessons mix choosing with
// typing from A1 (translate / fill exercises are typed or picked per the
// level rules in controllers/lesson_writing.go). Multiple-choice questions
// always carry their own options.

func seedItalian(db *gorm.DB) {
	seedItalianA1(db)
	seedItalianA2(db)
	seedItalianB1(db)
	seedItalianB2(db)
	seedItalianC1(db)
	seedItalianC2(db)
	seedItalianListening(db)
	seedItalianReading(db)
}

// ───────────────────── A1 — Primi Passi: First Steps ─────────────────────

func seedItalianA1(db *gorm.DB) {
	const it = "it"
	const u = "A1 · Primi Passi — First Steps"
	finch := "Professor Finch"

	s := addSkillL(db, it, u, "Sounds & Spelling", "Pure vowels, soft and hard c/g, gl, gn, sc — and double consonants.", "Music", "#6C3FC5", 1, 0)
	l := addLesson(db, s, "C, G & Friends", 1, 15,
		char(finch, "Benvenuti! Italian spelling is regular once you know a few rules. C and G are hard before a, o, u (casa, gatto) and soft before e, i (cena = 'chay-na', gelato = 'jeh-lah-to'). To keep them hard before e/i, add h: che (keh), chi (kee), ghiaccio. gl before i sounds like 'lli' in 'million' (figlio), gn like 'ny' in 'canyon' (gnocchi), sc before e/i like 'sh' (pesce)."),
		mc("How is 'cena' (dinner) pronounced?", "c before e", "chay-na (soft c)", "chay-na (soft c)", "kay-na (hard c)", "say-na", "tsay-na"),
		mc("How is 'chi' (who) pronounced?", "ch before i", "kee", "kee", "chee", "shee", "jee"),
		mc("'gn' in 'gnocchi' sounds like…", "gnocchi", "'ny' in 'canyon'", "'ny' in 'canyon'", "'g' + 'n' separately", "'n' only", "'gn' in 'signal'"),
		mc("'sc' in 'pesce' (fish) sounds like…", "pesce", "'sh' in 'shoe'", "'sh' in 'shoe'", "'sk' in 'skate'", "'s' + 'ch'", "'z'"),
		listen("Listen and choose the word", "figlio", "figlio (son)", "figlio (son)", "filo (thread)", "figlia (daughter)", "foglio (sheet)"),
		mc("How is 'gelato' pronounced?", "g before e", "jeh-LAH-to", "jeh-LAH-to", "geh-LAH-to (hard g)", "heh-LAH-to", "zheh-LAH-to"),
		speak("Mira", "Che bello! Un gelato, per favore."),
	)
	addVocab(db, l,
		vw("ciao", "hi; bye (informal)", "Ciao, come stai?", "Hi, how are you?", "Lumora"),
		vw("la cena", "dinner", "La cena è pronta.", "Dinner is ready.", "Cora"),
		vw("il gelato", "ice cream", "Vorrei un gelato.", "I'd like an ice cream.", "Pip"),
		vw("il figlio", "son", "Mio figlio ha cinque anni.", "My son is five.", "Nana"),
		vw("il pesce", "fish (as food)", "Il pesce è fresco.", "The fish is fresh.", "Cora"),
	)
	l = addLesson(db, s, "Double Consonants & Stress", 2, 15,
		char(finch, "Double consonants are held longer — and they change meaning: casa (house) vs cassa (till), pala (shovel) vs palla (ball), nono (ninth) vs nonno (grandfather), caro (dear) vs carro (cart). Stress usually falls on the second-to-last syllable (a-MI-co); a written accent marks a stressed final vowel: città, perché, caffè. The r is lightly rolled."),
		mc("Which word means 'grandfather'?", "nonno vs nono", "nonno", "nonno", "nono", "noni", "nonna"),
		mc("Which word means 'ball'?", "palla vs pala", "palla", "palla", "pala", "pala-la", "pallo"),
		listen("Listen and choose", "cassa", "cassa (till, checkout)", "cassa (till, checkout)", "casa (house)", "cosa (thing)", "caso (case)"),
		mc("Where is the stress in 'città' (city)?", "città", "On the last syllable: cit-TÀ", "On the last syllable: cit-TÀ", "On the first: CIT-tà", "No stress", "On the double t"),
		mc("Where is the stress in 'amico' (friend)?", "a-mi-co", "a-MI-co", "a-MI-co", "A-mi-co", "a-mi-CO", "Evenly spread"),
		speak("Blaze", "Il nonno ha una palla rossa."),
	)
	addVocab(db, l,
		vw("la casa", "house, home", "Torno a casa.", "I'm going home.", "Nana"),
		vw("il nonno / la nonna", "grandfather / grandmother", "Mia nonna cucina bene.", "My grandmother cooks well.", "Nana"),
		vw("la città", "city (stress on the end)", "Roma è una città antica.", "Rome is an ancient city.", "Zephyr"),
		vw("il caffè", "coffee", "Un caffè, per favore.", "A coffee, please.", "Cora"),
		vw("l'amico / l'amica", "friend (m / f)", "Marco è un mio amico.", "Marco is a friend of mine.", "Blaze"),
	)

	s = addSkillL(db, it, u, "Greetings & Politeness", "Buongiorno, per favore, grazie — and tu vs Lei.", "Hand", "#00C2A8", 2, 10)
	l = addLesson(db, s, "Buongiorno!", 1, 15,
		char("Lumora", "Buongiorno in the morning and early afternoon, Buonasera from mid-afternoon, Buonanotte at bedtime. Ciao is informal (hello and bye); Arrivederci is the polite goodbye. Per favore (please), Grazie (thank you) — Prego (you're welcome). Scusa (sorry, to a friend) / Scusi (excuse me, polite). Italians use Lei (formal 'you') with strangers and elders, tu with friends."),
		mc("Polite goodbye to a shopkeeper", "Leaving a shop", "Arrivederci!", "Arrivederci!", "Ciao bella!", "Buonanotte!", "Prego!"),
		mc("Someone says 'Grazie!' You reply:", "Grazie!", "Prego!", "Prego!", "Scusi!", "Ciao!", "Per favore!"),
		mc("Excuse me (to a stranger, polite)", "Getting attention", "Scusi!", "Scusi!", "Scusa!", "Ciao!", "Grazie!"),
		mc("Which 'you' do you use with an older stranger?", "Register", "Lei", "Lei", "tu", "voi only", "loro"),
		tr("Translate", "Thank you very much", "Grazie mille"),
		fill("Complete the evening greeting", "Buona___! (Good evening)", "sera"),
		speak("Lumora", "Buongiorno! Come sta? Bene, grazie, e Lei?"),
	)
	addVocab(db, l,
		vw("buongiorno", "good morning / good day", "Buongiorno, signora!", "Good morning, madam!", "Lumora"),
		vw("buonasera", "good evening", "Buonasera a tutti.", "Good evening, everyone.", "Zephyr"),
		vw("per favore", "please", "Un'acqua, per favore.", "A water, please.", "Cora"),
		vw("grazie / prego", "thank you / you're welcome", "Grazie! — Prego!", "Thanks! — You're welcome!", "Mira"),
		vw("arrivederci", "goodbye (polite)", "Arrivederci, a domani.", "Goodbye, see you tomorrow.", "Nana"),
		vw("Come sta? / Come stai?", "How are you? (formal / informal)", "Come stai, Luca?", "How are you, Luca?", "Blaze"),
	)

	s = addSkillL(db, it, u, "Essere & Avere", "To be and to have — the two verbs that unlock everything.", "CheckCircle", "#F5A623", 3, 20)
	l = addLesson(db, s, "Essere — To Be", 1, 15,
		char(finch, "essere: io sono, tu sei, lui/lei è, noi siamo, voi siete, loro sono. Subject pronouns are usually dropped because the ending says who: Sono italiano (I'm Italian), Sei stanca? (are you tired?). Note è (is) has an accent; e without one means 'and'. Negative: non before the verb — Non sono di Roma."),
		mc("We are friends.", "essere", "Siamo amici.", "Siamo amici.", "Sono amici io.", "Siete amici.", "Abbiamo amici."),
		mc("Which form goes with 'voi'?", "voi ___ italiani", "siete", "siete", "siamo", "sono", "sei"),
		mc("What is the difference between è and e?", "è / e", "è = is, e = and", "è = is, e = and", "No difference", "e = is, è = and", "è is plural"),
		tr("Translate", "I am not tired", "Non sono stanco"),
		tr("Translate", "Where are you from? (informal)", "Di dove sei"),
		fill("Complete", "Lei ___ medico. (She is a doctor.)", "è"),
		speak("Lumora", "Sono Lumora, sono di Nairobi. E tu, di dove sei?"),
	)
	addVocab(db, l,
		vw("essere", "to be", "Sono felice.", "I'm happy.", finch),
		vw("io sono / tu sei / lui è", "I am / you are / he is", "Lui è mio fratello.", "He is my brother.", finch),
		vw("stanco / stanca", "tired (m / f)", "Sei stanca?", "Are you tired?", "Nana"),
		vw("Di dove sei?", "Where are you from?", "Di dove sei? Sono di Milano.", "Where are you from? I'm from Milan.", "Mira"),
	)
	l = addLesson(db, s, "Avere — To Have", 2, 15,
		char(finch, "avere: ho, hai, ha, abbiamo, avete, hanno — the h is silent. Italian uses avere where English uses 'be': ho vent'anni (I'm twenty — 'I have twenty years'), ho fame (I'm hungry), ho sete (thirsty), ho freddo (cold), ho ragione (I'm right), ho fretta (in a hurry)."),
		mc("I'm hungry.", "avere fame", "Ho fame.", "Ho fame.", "Sono fame.", "Ho affamato.", "Ha fame."),
		mc("How old are you? (informal)", "Asking age", "Quanti anni hai?", "Quanti anni hai?", "Quanti anni sei?", "Come anni hai?", "Che età sei?"),
		mc("They have a dog.", "avere", "Hanno un cane.", "Hanno un cane.", "Sono un cane.", "Ha un cane.", "Abbiamo un cane."),
		tr("Translate", "I am twenty years old", "Ho vent'anni"),
		tr("Translate", "We are thirsty", "Abbiamo sete"),
		fill("Complete", "Voi ___ ragione. (You all are right.)", "avete"),
		speak("Pip", "Ho dieci anni e ho sempre fame!"),
	)
	addVocab(db, l,
		vw("avere", "to have", "Hai una penna?", "Do you have a pen?", finch),
		vw("avere fame / sete", "to be hungry / thirsty", "Abbiamo fame.", "We're hungry.", "Pip"),
		vw("l'anno / gli anni", "year / years", "Ho trent'anni.", "I'm thirty.", "Zephyr"),
		vw("il cane / il gatto", "dog / cat", "Ho un gatto nero.", "I have a black cat.", "Pip"),
		vw("avere ragione", "to be right", "Hai ragione tu.", "You're right.", "Mira"),
	)

	s = addSkillL(db, it, u, "Gender & Articles", "il, lo, la, l', i, gli, le — always learn the article with the noun.", "Layers", "#FF5C5C", 4, 32)
	l = addLesson(db, s, "The Definite Article", 1, 18,
		char(finch, "Every noun is masculine or feminine; -o is usually masculine (il libro), -a feminine (la casa), -e can be either (il pane, la notte). Masculine: il before most consonants, lo before s + consonant, z, gn, ps (lo studente, lo zaino), l' before a vowel (l'amico). Plural: i / gli (gli studenti, gli amici). Feminine: la / l' → le."),
		mc("Choose the article", "___ studente", "lo", "lo", "il", "la", "l'"),
		mc("Choose the article", "___ zaino (backpack)", "lo", "lo", "il", "la", "gli"),
		mc("Plural of 'l'amico'", "l'amico", "gli amici", "gli amici", "i amici", "le amiche", "l'amici"),
		mc("Plural of 'la casa'", "la casa", "le case", "le case", "la case", "i casi", "gli case"),
		mc("Choose the article", "___ notte (night, feminine)", "la", "la", "il", "lo", "l'"),
		tr("Translate", "the books", "i libri"),
		fill("Article", "___ zio di Marco (Marco's uncle)", "lo"),
		speak("Cora", "Il libro, lo zaino, l'amica, gli studenti, le case."),
	)
	addVocab(db, l,
		vw("il libro / i libri", "book / books", "I libri sono sul tavolo.", "The books are on the table.", finch),
		vw("lo studente / gli studenti", "student / students", "Gli studenti studiano.", "The students study.", finch),
		vw("lo zaino", "backpack", "Lo zaino è pesante.", "The backpack is heavy.", "Pip"),
		vw("la notte", "night", "Buona notte!", "Good night!", "Nana"),
		vw("il pane", "bread", "Il pane è caldo.", "The bread is warm.", "Cora"),
	)
	l = addLesson(db, s, "A, An & Plurals", 2, 18,
		char(finch, "Indefinite articles: un (most masculine nouns: un libro, un amico), uno (before s+consonant, z: uno zaino), una (feminine: una casa), un' (feminine before a vowel: un'amica). Plurals change the ending: -o → -i (libro → libri), -a → -e (casa → case), -e → -i (pane → pani; notte → notti). Words with a final accent or from other languages don't change: il caffè → i caffè, il bar → i bar."),
		mc("Choose the article", "___ amica (a female friend)", "un'", "un'", "un", "una", "uno"),
		mc("Choose the article", "___ amico (a male friend)", "un", "un", "un'", "uno", "una"),
		mc("Choose the article", "___ specchio (mirror)", "uno", "uno", "un", "una", "un'"),
		mc("Plural of 'il caffè'", "il caffè", "i caffè", "i caffè", "i caffi", "gli caffè", "le caffè"),
		tr("Translate", "a house and two cars", "una casa e due macchine"),
		fill("Plural", "la stazione → le ___", "stazioni"),
		speak("Blaze", "Vorrei un caffè e una pasta, per favore."),
	)
	addVocab(db, l,
		vw("un / una", "a, an (m / f)", "Ho un fratello e una sorella.", "I have a brother and a sister.", finch),
		vw("la macchina", "car", "La macchina è nuova.", "The car is new.", "Blaze"),
		vw("la stazione", "station", "Dov'è la stazione?", "Where is the station?", "Riko"),
		vw("lo specchio", "mirror", "Guardo nello specchio.", "I look in the mirror.", "Mira"),
	)

	s = addSkillL(db, it, u, "Adjective Agreement", "rosso, rossa, rossi, rosse — and bello, grande, buono.", "Sparkles", "#9B51E0", 5, 44)
	l = addLesson(db, s, "Describing Things", 1, 18,
		char(finch, "Adjectives agree in gender and number and usually follow the noun. -o adjectives have four forms: un vestito rosso, una gonna rossa, vestiti rossi, gonne rosse. -e adjectives have two: un libro interessante, una storia interessante, libri/storie interessanti. A few common ones go before the noun: una bella giornata, un buon amico, una grande città."),
		mc("Agreement", "le scarpe ___ (the red shoes)", "rosse", "rosse", "rossi", "rossa", "rosso"),
		mc("Agreement", "una ragazza ___ (an intelligent girl)", "intelligente", "intelligente", "intelligenta", "intelligenti", "intelligento"),
		mc("Agreement", "i ragazzi ___ (the tall boys)", "alti", "alti", "alte", "alto", "alta"),
		mc("What a beautiful day!", "Exclamation", "Che bella giornata!", "Che bella giornata!", "Che bello giornata!", "Che giornata bello!", "Che belli giornata!"),
		tr("Translate", "The houses are big", "Le case sono grandi"),
		fill("Agreement", "Le mele sono ___ (green, verde).", "verdi"),
		speak("Mira", "Ho una macchina rossa e due gatti neri."),
	)
	addVocab(db, l,
		vw("rosso / rossa", "red", "Una mela rossa.", "A red apple.", "Pip"),
		vw("grande", "big, great", "Una grande città.", "A big city.", "Zephyr"),
		vw("piccolo / piccola", "small", "Un paese piccolo.", "A small village.", "Nana"),
		vw("bello / bella", "beautiful, nice", "Che bella giornata!", "What a beautiful day!", "Mira"),
		vw("interessante", "interesting", "Un libro interessante.", "An interesting book.", finch),
	)

	s = addSkillL(db, it, u, "Present Tense: -are, -ere, -ire", "parlo, scrivo, dormo, capisco.", "MessageCircle", "#17A3DD", 6, 56)
	l = addLesson(db, s, "Regular Verbs", 1, 18,
		char(finch, "Three conjugations. parlare: parlo, parli, parla, parliamo, parlate, parlano. scrivere: scrivo, scrivi, scrive, scriviamo, scrivete, scrivono. dormire: dormo, dormi, dorme, dormiamo, dormite, dormono. The stress in the 'loro' form stays early: PAR-la-no, not par-LA-no."),
		mc("We speak Italian.", "parlare", "Parliamo italiano.", "Parliamo italiano.", "Parlamo italiano.", "Parlano italiano.", "Parlate italiano."),
		mc("They write.", "scrivere", "Scrivono.", "Scrivono.", "Scrivano.", "Scriviamo.", "Scrive."),
		mc("You (tu) sleep a lot.", "dormire", "Dormi molto.", "Dormi molto.", "Dorme molto.", "Dormo molto.", "Dormite molto."),
		tr("Translate", "I live in Rome", "Abito a Roma"),
		tr("Translate", "Do you (voi) speak English?", "Parlate inglese"),
		fill("Complete", "Lei ___ una lettera. (scrivere)", "scrive"),
		speak("Cora", "Parlo italiano, scrivo un messaggio e poi dormo."),
	)
	addVocab(db, l,
		vw("parlare", "to speak", "Parli inglese?", "Do you speak English?", "Lumora"),
		vw("abitare", "to live (reside)", "Abito a Firenze.", "I live in Florence.", "Zephyr"),
		vw("scrivere", "to write", "Scrivo una mail.", "I'm writing an email.", finch),
		vw("leggere", "to read", "Leggo il giornale.", "I read the newspaper.", finch),
		vw("dormire", "to sleep", "Dormo otto ore.", "I sleep eight hours.", "Nana"),
	)
	l = addLesson(db, s, "-isc- Verbs & Negatives", 2, 18,
		char(finch, "Many -ire verbs add -isc- in all forms except noi and voi: capire → capisco, capisci, capisce, capiamo, capite, capiscono. Others like finire and preferire do the same. Negative: put non before the verb — Non capisco (I don't understand). Questions just use rising intonation: Capisci?"),
		mc("I understand.", "capire", "Capisco.", "Capisco.", "Capo.", "Capiscio.", "Capiamo."),
		mc("We finish at five.", "finire", "Finiamo alle cinque.", "Finiamo alle cinque.", "Finisciamo alle cinque.", "Finiscono alle cinque.", "Finisco alle cinque."),
		mc("They prefer tea.", "preferire", "Preferiscono il tè.", "Preferiscono il tè.", "Preferono il tè.", "Preferiamo il tè.", "Preferisce il tè."),
		tr("Translate", "I don't understand", "Non capisco"),
		tr("Translate", "She prefers coffee", "Preferisce il caffè"),
		fill("Complete", "Voi ___ il lavoro? (finire)", "finite"),
		speak("Lumora", "Scusi, non capisco. Può ripetere, per favore?"),
	)
	addVocab(db, l,
		vw("capire", "to understand", "Non capisco la domanda.", "I don't understand the question.", "Lumora"),
		vw("finire", "to finish", "Il film finisce tardi.", "The film ends late.", "Zephyr"),
		vw("preferire", "to prefer", "Preferisco il mare.", "I prefer the sea.", "Mira"),
		vw("ripetere", "to repeat", "Può ripetere, per favore?", "Can you repeat, please?", "Lumora"),
	)

	s = addSkillL(db, it, u, "Irregular Verbs", "andare, fare, venire, stare, dire, uscire.", "Compass", "#00C2A8", 7, 68)
	l = addLesson(db, s, "The Everyday Irregulars", 1, 18,
		char(finch, "The most-used verbs are irregular. andare: vado, vai, va, andiamo, andate, vanno. fare: faccio, fai, fa, facciamo, fate, fanno. venire: vengo, vieni, viene, veniamo, venite, vengono. stare: sto, stai, sta, stiamo, state, stanno (Come stai?). dire: dico, dici, dice, diciamo, dite, dicono. uscire: esco, esci, esce, usciamo, uscite, escono."),
		mc("I'm going to the cinema.", "andare", "Vado al cinema.", "Vado al cinema.", "Ando al cinema.", "Va al cinema.", "Vanno al cinema."),
		mc("What are you doing? (tu)", "fare", "Cosa fai?", "Cosa fai?", "Cosa fa?", "Cosa faci?", "Cosa facci?"),
		mc("Are you coming too? (tu)", "venire", "Vieni anche tu?", "Vieni anche tu?", "Veni anche tu?", "Viene anche tu?", "Vengo anche tu?"),
		mc("We're going out tonight.", "uscire", "Usciamo stasera.", "Usciamo stasera.", "Esciamo stasera.", "Escono stasera.", "Usciamo ieri."),
		tr("Translate", "What do you say? (tu)", "Cosa dici"),
		fill("Complete", "Loro ___ a casa. (stare)", "stanno"),
		speak("Blaze", "Stasera esco con gli amici: andiamo a mangiare una pizza!"),
	)
	addVocab(db, l,
		vw("andare", "to go", "Andiamo al mare!", "Let's go to the seaside!", "Blaze"),
		vw("fare", "to do, make", "Faccio una passeggiata.", "I'm taking a walk.", "Zephyr"),
		vw("venire", "to come", "Vieni alla festa?", "Are you coming to the party?", "Mira"),
		vw("stare", "to stay; to be (well)", "Sto bene, grazie.", "I'm well, thanks.", "Lumora"),
		vw("uscire", "to go out", "Esco alle otto.", "I go out at eight.", "Blaze"),
	)

	s = addSkillL(db, it, u, "Numbers, Time & Days", "Che ore sono? Sono le tre. Months and days.", "Clock", "#F5A623", 8, 80)
	l = addLesson(db, s, "Che Ore Sono?", 1, 18,
		char(finch, "Numbers: uno, due, tre, quattro, cinque, sei, sette, otto, nove, dieci, undici… venti, trenta, cento, mille (plural mila: duemila). Time: Che ore sono? — Sono le tre (it's three); È l'una (it's one); È mezzogiorno / mezzanotte. Add e mezza (half past), e un quarto (quarter past), meno un quarto (quarter to). Days are lowercase: lunedì, martedì, mercoledì, giovedì, venerdì, sabato, domenica."),
		mc("It's 3:30.", "Time", "Sono le tre e mezza.", "Sono le tre e mezza.", "È le tre e mezza.", "Sono tre e mezzo le.", "Sono le tre meno mezza."),
		mc("It's one o'clock.", "Singular hour", "È l'una.", "È l'una.", "Sono le una.", "È uno.", "Sono l'una."),
		mc("Which day is giovedì?", "giovedì", "Thursday", "Thursday", "Tuesday", "Friday", "Wednesday"),
		mc("What number is 'diciassette'?", "diciassette", "17", "17", "7", "70", "27"),
		tr("Translate", "What time is it?", "Che ore sono"),
		fill("Numbers", "mille, due___ (two thousand)", "mila"),
		speak("Riko", "Il treno parte alle otto e un quarto di lunedì."),
	)
	addVocab(db, l,
		vw("Che ore sono?", "What time is it?", "Scusi, che ore sono?", "Excuse me, what time is it?", "Riko"),
		vw("e mezza / e un quarto", "half past / quarter past", "Sono le due e un quarto.", "It's quarter past two.", "Riko"),
		vw("lunedì", "Monday", "Lunedì lavoro.", "On Monday I work.", "Zephyr"),
		vw("oggi / domani / ieri", "today / tomorrow / yesterday", "Domani è sabato.", "Tomorrow is Saturday.", "Pip"),
		vw("cento / mille", "hundred / thousand", "Costa mille euro.", "It costs a thousand euros.", "Cora"),
	)

	s = addSkillL(db, it, u, "Family & Possessives", "il mio libro, mia madre — when the article drops.", "Heart", "#FF5C5C", 9, 92)
	l = addLesson(db, s, "La Mia Famiglia", 1, 18,
		char("Nana", "Possessives agree with the thing owned and normally take the article: il mio libro, la mia casa, i miei amici, le mie scarpe. Exception: singular, unmodified family members drop it — mia madre, mio padre, tuo fratello — but plurals keep it: i miei genitori, le mie sorelle. Loro always keeps it: la loro madre."),
		mc("my mother", "Family, singular", "mia madre", "mia madre", "la mia madre", "il mio madre", "mio madre"),
		mc("my parents", "Family, plural", "i miei genitori", "i miei genitori", "miei genitori", "le mie genitori", "i mio genitori"),
		mc("our house", "Not family", "la nostra casa", "la nostra casa", "nostra casa", "il nostro casa", "la nostro casa"),
		mc("their father", "loro", "il loro padre", "il loro padre", "loro padre", "il suo padre", "suo loro padre"),
		tr("Translate", "my brother is tall", "mio fratello è alto"),
		fill("Possessive", "___ mia sorellina (my little sister — a suffix brings the article back)", "la"),
		speak("Nana", "Questa è mia madre, e questi sono i miei nonni."),
	)
	addVocab(db, l,
		vw("la madre / il padre", "mother / father", "Mia madre è insegnante.", "My mother is a teacher.", "Nana"),
		vw("i genitori", "parents", "I miei genitori abitano a Napoli.", "My parents live in Naples.", "Nana"),
		vw("il fratello / la sorella", "brother / sister", "Ho due sorelle.", "I have two sisters.", "Pip"),
		vw("la famiglia", "family", "La mia famiglia è grande.", "My family is big.", "Lumora"),
		vw("il marito / la moglie", "husband / wife", "Mia moglie è di Torino.", "My wife is from Turin.", "Zephyr"),
	)

	s = addSkillL(db, it, u, "Al Bar & Food", "Vorrei…, il conto, and how Italians do coffee.", "Coffee", "#6C3FC5", 10, 104)
	l = addLesson(db, s, "At the Café", 1, 18,
		char("Cora", "Order with Vorrei… (I'd like) or Prendo… (I'll have). Un caffè is an espresso; cappuccino is a morning drink — ordering one after lunch marks you as a tourist! At many bars you pay first at the cassa, then show the receipt at the counter. At a restaurant: Il conto, per favore. The coperto is a small cover charge."),
		mc("Polite way to order", "Ordering", "Vorrei un caffè, per favore.", "Vorrei un caffè, per favore.", "Voglio caffè subito.", "Dammi un caffè.", "Caffè!"),
		mc("What is 'un caffè' in Italy?", "Al bar", "An espresso", "An espresso", "A large filter coffee", "A latte", "An iced coffee"),
		mc("The bill, please.", "At the end of a meal", "Il conto, per favore.", "Il conto, per favore.", "La cassa, per favore.", "Il prezzo, grazie.", "Pago io!"),
		mc("What is the 'coperto'?", "Sul conto: coperto €2", "A small cover charge", "A small cover charge", "A tip", "A free starter", "The menu"),
		tr("Translate", "I'll have a pizza margherita", "Prendo una pizza margherita"),
		tr("Translate", "The pasta is delicious", "La pasta è buonissima"),
		speak("Cora", "Buongiorno! Vorrei un cappuccino e un cornetto, per favore."),
	)
	addVocab(db, l,
		vw("vorrei", "I would like", "Vorrei un bicchiere d'acqua.", "I'd like a glass of water.", "Cora"),
		vw("il conto", "the bill", "Ci porta il conto?", "Could you bring us the bill?", "Cora"),
		vw("il cornetto", "croissant (Italian style)", "Un cornetto alla crema.", "A custard croissant.", "Pip"),
		vw("buono / buonissimo", "good / delicious", "Il gelato è buonissimo!", "The ice cream is delicious!", "Pip"),
		vw("l'acqua (frizzante / naturale)", "water (sparkling / still)", "Acqua frizzante, per favore.", "Sparkling water, please.", "Cora"),
	)

	s = addSkillL(db, it, u, "Directions & Prepositions", "a, in, da, di, su — and al, nel, dal, sul, del.", "Plane", "#17A3DD", 11, 116)
	l = addLesson(db, s, "Dov'è…?", 1, 18,
		char("Riko", "Directions: vada dritto (go straight — formal), giri a destra / a sinistra (turn right/left), accanto a (next to), di fronte a (opposite). Prepositions merge with the article: a + il = al, in + il = nel, da + il = dal, su + il = sul, di + il = del (vado al bar, è nel cassetto, vengo dal medico). Cities take a (a Roma), countries in (in Italia)."),
		mc("I'm going to Rome.", "a / in", "Vado a Roma.", "Vado a Roma.", "Vado in Roma.", "Vado da Roma.", "Vado al Roma."),
		mc("I live in Italy.", "Countries", "Abito in Italia.", "Abito in Italia.", "Abito a Italia.", "Abito nell'Italia.", "Abito da Italia."),
		mc("It's in the drawer.", "in + il", "È nel cassetto.", "È nel cassetto.", "È in il cassetto.", "È al cassetto.", "È del cassetto."),
		mc("Turn left. (formal)", "Directions", "Giri a sinistra.", "Giri a sinistra.", "Gira destra.", "Vada sinistra.", "Giri di sinistra."),
		tr("Translate", "Where is the station?", "Dov'è la stazione"),
		fill("a + il", "Andiamo ___ ristorante. (to the restaurant)", "al"),
		speak("Riko", "Scusi, per il Colosseo? Vada dritto e poi giri a destra."),
	)
	addVocab(db, l,
		vw("Dov'è…?", "Where is…?", "Dov'è il bagno?", "Where's the bathroom?", "Riko"),
		vw("a destra / a sinistra", "to the right / left", "Il museo è a destra.", "The museum is on the right.", "Riko"),
		vw("dritto", "straight on", "Sempre dritto.", "Keep straight on.", "Riko"),
		vw("vicino / lontano", "near / far", "È vicino al centro.", "It's near the centre.", "Zephyr"),
		vw("accanto a", "next to", "La farmacia è accanto alla banca.", "The pharmacy is next to the bank.", "Riko"),
	)
}

// ───────────────────── A2 — La Vita Quotidiana: Everyday Life ─────────────────────

func seedItalianA2(db *gorm.DB) {
	const it = "it"
	const u = "A2 · La Vita Quotidiana — Everyday Life"
	finch := "Professor Finch"

	s := addSkillL(db, it, u, "Passato Prossimo with Avere", "ho mangiato, ho fatto — and the irregular participles.", "Clock", "#6C3FC5", 12, 140)
	l := addLesson(db, s, "What I Did", 1, 20,
		char(finch, "The spoken past: avere (present) + past participle. -are → -ato (mangiato), -ere → -uto (venduto), -ire → -ito (dormito). Many common participles are irregular: fare → fatto, dire → detto, vedere → visto, scrivere → scritto, prendere → preso, leggere → letto, aprire → aperto, mettere → messo."),
		mc("I ate pizza.", "mangiare", "Ho mangiato la pizza.", "Ho mangiato la pizza.", "Sono mangiato la pizza.", "Ho mangiare la pizza.", "Mangiavo la pizza ieri una volta."),
		mc("Past participle of 'fare'", "fare", "fatto", "fatto", "fato", "faciuto", "fatto stato"),
		mc("Past participle of 'vedere'", "vedere", "visto", "visto", "vedato", "vidato", "vesto"),
		mc("We read the book.", "leggere", "Abbiamo letto il libro.", "Abbiamo letto il libro.", "Abbiamo leggiuto il libro.", "Siamo letti il libro.", "Abbiamo legato il libro."),
		tr("Translate", "What did you do yesterday? (tu)", "Cosa hai fatto ieri"),
		fill("Participle", "Ho ___ una lettera. (scrivere)", "scritto"),
		speak("Blaze", "Ieri ho visto un film e ho mangiato una pizza."),
	)
	addVocab(db, l,
		vw("ho mangiato", "I ate / have eaten", "Hai già mangiato?", "Have you eaten already?", "Cora"),
		vw("fatto", "done, made (fare)", "Cosa hai fatto?", "What did you do?", finch),
		vw("visto", "seen (vedere)", "Ho visto Marco.", "I saw Marco.", "Mira"),
		vw("preso", "taken (prendere)", "Ho preso il treno.", "I took the train.", "Riko"),
		vw("già", "already", "L'ho già fatto.", "I've already done it.", finch),
	)

	s = addSkillL(db, it, u, "Passato Prossimo with Essere", "sono andato/andata — movement, change and agreement.", "CheckCircle", "#00C2A8", 13, 155)
	l = addLesson(db, s, "Where I Went", 1, 20,
		char(finch, "Verbs of movement and change (andare, venire, partire, arrivare, uscire, tornare, nascere, morire, diventare), plus essere itself and reflexives, use essere — and then the participle agrees with the subject: Marco è andato, Giulia è andata, siamo partiti, le ragazze sono uscite. essere's own participle is stato: sono stato a Roma."),
		mc("Giulia went to Milan.", "andare", "Giulia è andata a Milano.", "Giulia è andata a Milano.", "Giulia ha andato a Milano.", "Giulia è andato a Milano.", "Giulia ha andata a Milano."),
		mc("We left at eight.", "partire (noi, mixed group)", "Siamo partiti alle otto.", "Siamo partiti alle otto.", "Abbiamo partito alle otto.", "Siamo partito alle otto.", "Sono partiti alle otto."),
		mc("Which auxiliary for 'arrivare'?", "arrivare", "essere", "essere", "avere", "both equally", "stare"),
		mc("I've been to Rome (Marco speaking).", "essere → stato", "Sono stato a Roma.", "Sono stato a Roma.", "Ho stato a Roma.", "Sono essere a Roma.", "Sono stata a Roma."),
		tr("Translate", "She was born in Naples", "È nata a Napoli"),
		fill("Agreement", "Le ragazze sono usci___ tardi.", "te"),
		speak("Mira", "Sono arrivata a Venezia ieri sera e sono uscita subito!"),
	)
	addVocab(db, l,
		vw("sono andato / andata", "I went (m / f)", "Sono andata al mare.", "I went to the seaside.", "Mira"),
		vw("partire", "to leave, depart", "Il treno è partito.", "The train has left.", "Riko"),
		vw("arrivare", "to arrive", "Siamo arrivati tardi.", "We arrived late.", "Riko"),
		vw("nascere", "to be born", "Sono nato nel 2000.", "I was born in 2000.", "Pip"),
		vw("tornare", "to return", "Quando sei tornato?", "When did you get back?", "Zephyr"),
	)

	s = addSkillL(db, it, u, "The Imperfetto", "era, avevo, facevo — background, habits and descriptions.", "BookOpen", "#F5A623", 14, 170)
	l = addLesson(db, s, "When I Was Little…", 1, 20,
		char(finch, "The imperfetto describes how things were, habits in the past and ongoing background: -avo, -evo, -ivo + endings (parlavo, avevo, dormivo). essere is irregular: ero, eri, era, eravamo, eravate, erano; fare → facevo. Contrast: the passato prossimo is the event, the imperfetto the scene — Mentre leggevo, è arrivato Luca (while I was reading, Luca arrived)."),
		mc("When I was little, I lived in Rome.", "Background / habit", "Quando ero piccolo, abitavo a Roma.", "Quando ero piccolo, abitavo a Roma.", "Quando sono stato piccolo, ho abitato a Roma.", "Quando ero piccolo, ho abitato a Roma ogni giorno.", "Quando sarò piccolo, abiterò a Roma."),
		mc("While I was cooking, the phone rang.", "Background + event", "Mentre cucinavo, ha squillato il telefono.", "Mentre cucinavo, ha squillato il telefono.", "Mentre ho cucinato, squillava il telefono.", "Mentre cucino, ha squillato il telefono.", "Mentre cucinavo, squillava il telefono una volta."),
		mc("Imperfetto of essere (noi)", "essere", "eravamo", "eravamo", "erano", "siamo stati", "essevamo"),
		mc("Which tense for 'It was sunny and warm'?", "Description", "Imperfetto: C'era il sole e faceva caldo.", "Imperfetto: C'era il sole e faceva caldo.", "Passato prossimo: C'è stato il sole e ha fatto caldo.", "Futuro: Ci sarà il sole.", "Presente: C'è il sole."),
		tr("Translate", "Every summer we went to the sea", "Ogni estate andavamo al mare"),
		fill("Imperfetto", "Da bambino, ___ molto. (leggere, io)", "leggevo"),
		speak("Nana", "Da giovane, andavo a ballare ogni sabato."),
	)
	addVocab(db, l,
		vw("ero / era", "I was / he, she was", "Era una bella giornata.", "It was a beautiful day.", "Nana"),
		vw("mentre", "while", "Mentre studiavo, pioveva.", "While I was studying, it was raining.", finch),
		vw("da bambino / da bambina", "as a child", "Da bambina giocavo in giardino.", "As a child I played in the garden.", "Nana"),
		vw("ogni estate", "every summer", "Ogni estate andavamo in montagna.", "Every summer we went to the mountains.", "Zephyr"),
	)

	s = addSkillL(db, it, u, "Reflexive Verbs & Daily Routine", "mi sveglio, mi alzo, mi vesto.", "Clock", "#FF5C5C", 15, 185)
	l = addLesson(db, s, "La Mia Giornata", 1, 20,
		char(finch, "Reflexive verbs take mi, ti, si, ci, vi, si: svegliarsi → mi sveglio (I wake up), alzarsi → mi alzo, lavarsi, vestirsi, chiamarsi (Mi chiamo Anna), divertirsi (to have fun). In the passato prossimo they use essere and agree: Giulia si è svegliata tardi; ci siamo divertiti."),
		mc("I wake up at seven.", "svegliarsi", "Mi sveglio alle sette.", "Mi sveglio alle sette.", "Sveglio alle sette.", "Si sveglio alle sette.", "Mi svegli alle sette."),
		mc("We had fun.", "divertirsi, passato prossimo", "Ci siamo divertiti.", "Ci siamo divertiti.", "Abbiamo divertiti.", "Ci abbiamo divertito.", "Si siamo divertiti."),
		mc("She got dressed quickly.", "vestirsi", "Si è vestita in fretta.", "Si è vestita in fretta.", "Ha vestita in fretta.", "Si ha vestito in fretta.", "Si è vestito in fretta (she)."),
		mc("What's your name? (tu)", "chiamarsi", "Come ti chiami?", "Come ti chiami?", "Come si chiami?", "Come chiami?", "Come ti chiama?"),
		tr("Translate", "I get up early", "Mi alzo presto"),
		fill("Reflexive", "Loro ___ lavano le mani. (They wash their hands.)", "si"),
		speak("Pip", "Mi sveglio, mi lavo, mi vesto e vado a scuola."),
	)
	addVocab(db, l,
		vw("svegliarsi", "to wake up", "A che ora ti svegli?", "What time do you wake up?", "Pip"),
		vw("alzarsi", "to get up", "Mi alzo alle sei.", "I get up at six.", "Blaze"),
		vw("vestirsi", "to get dressed", "Mi vesto in fretta.", "I get dressed quickly.", "Mira"),
		vw("divertirsi", "to have fun", "Divertiti!", "Have fun!", "Blaze"),
		vw("chiamarsi", "to be called", "Mi chiamo Luca.", "My name is Luca.", "Lumora"),
	)

	s = addSkillL(db, it, u, "Modal Verbs", "potere, volere, dovere — can, want, must.", "Shield", "#17A3DD", 16, 200)
	l = addLesson(db, s, "Can, Want, Must", 1, 20,
		char(finch, "potere: posso, puoi, può, possiamo, potete, possono. volere: voglio, vuoi, vuole, vogliamo, volete, vogliono. dovere: devo, devi, deve, dobbiamo, dovete, devono. They take an infinitive: Posso entrare? (may I come in?), Devo partire (I have to leave). For politeness use vorrei / potrei instead of voglio / posso."),
		mc("May I come in?", "potere", "Posso entrare?", "Posso entrare?", "Puoi entrare io?", "Devo entrare?", "Voglio entrare?"),
		mc("We have to work tomorrow.", "dovere", "Domani dobbiamo lavorare.", "Domani dobbiamo lavorare.", "Domani dovemo lavorare.", "Domani devono lavorare.", "Domani possiamo lavorare."),
		mc("They want to go to the beach.", "volere", "Vogliono andare in spiaggia.", "Vogliono andare in spiaggia.", "Volono andare in spiaggia.", "Vuole andare in spiaggia.", "Vogliamo andare in spiaggia."),
		mc("Most polite way to ask for something", "Politeness", "Vorrei…", "Vorrei…", "Voglio…", "Devo avere…", "Dammi…"),
		tr("Translate", "Can you help me? (tu)", "Puoi aiutarmi"),
		fill("Modal", "Lui non ___ venire stasera. (potere)", "può"),
		speak("Lumora", "Scusi, posso pagare con la carta?"),
	)
	addVocab(db, l,
		vw("potere", "can, may", "Posso aiutarti?", "Can I help you?", "Lumora"),
		vw("volere", "to want", "Vuoi un caffè?", "Do you want a coffee?", "Cora"),
		vw("dovere", "to have to, must", "Devo andare.", "I have to go.", "Zephyr"),
		vw("aiutare", "to help", "Mi puoi aiutare?", "Can you help me?", "Lumora"),
	)

	s = addSkillL(db, it, u, "Object Pronouns", "lo, la, li, le; gli, le; and ne.", "Link", "#9B51E0", 17, 215)
	l = addLesson(db, s, "Lo, La, Gli & Ne", 1, 20,
		char(finch, "Direct object pronouns: mi, ti, lo, la, ci, vi, li, le — they go before the conjugated verb: Lo vedo (I see him/it), Le compro (I buy them, f). Indirect: mi, ti, gli (to him; colloquially also to them), le (to her), ci, vi. ne replaces 'of it / some': Quante mele vuoi? Ne voglio due. With an infinitive they attach to the end: Voglio vederlo."),
		mc("Do you know Maria? — Yes, I know her.", "Conosci Maria?", "Sì, la conosco.", "Sì, la conosco.", "Sì, lo conosco.", "Sì, le conosco.", "Sì, conosco la."),
		mc("I'm writing to him.", "Indirect, masculine", "Gli scrivo.", "Gli scrivo.", "Lo scrivo.", "Le scrivo a lui.", "Scrivo gli."),
		mc("How many do you want? — I want two (of them).", "ne", "Ne voglio due.", "Ne voglio due.", "Li voglio due.", "Voglio due ne.", "Lo voglio due."),
		mc("I want to see it.", "Pronoun with an infinitive", "Voglio vederlo.", "Voglio vederlo.", "Voglio lo vedere.", "Lo voglio vedo.", "Vedere lo voglio."),
		tr("Translate", "I buy them (the books)", "Li compro"),
		fill("Indirect pronoun", "Telefono a Giulia → ___ telefono.", "Le"),
		speak("Mira", "Il libro? L'ho già letto e te lo presto!"),
	)
	addVocab(db, l,
		vw("lo / la", "him, it / her, it", "Lo vedo domani.", "I'll see him tomorrow.", finch),
		vw("gli / le", "to him / to her", "Le ho detto la verità.", "I told her the truth.", "Mira"),
		vw("ne", "of it, some", "Ne prendo un po'.", "I'll take a bit (of it).", "Cora"),
		vw("conoscere", "to know (a person, place)", "Conosci Roma?", "Do you know Rome?", "Zephyr"),
	)

	s = addSkillL(db, it, u, "Shopping & Clothes", "Quanto costa? Che taglia porta? — and the sales.", "ShoppingBag", "#00C2A8", 18, 230)
	l = addLesson(db, s, "In Negozio", 1, 20,
		char("Cora", "Quanto costa? / Quanto costano? (how much is it / are they?). Clothes: la maglietta, i pantaloni, la giacca, le scarpe. The assistant asks Che taglia porta? (what size do you wear — formal) or Che numero? for shoes. Posso provarlo? (can I try it on?). È troppo caro (it's too expensive); i saldi (the sales), lo sconto (discount)."),
		mc("How much do these shoes cost?", "Plural", "Quanto costano queste scarpe?", "Quanto costano queste scarpe?", "Quanto costa queste scarpe?", "Quanti costano queste scarpe?", "Quanto prezzo scarpe?"),
		mc("Can I try it on? (the jacket)", "la giacca", "Posso provarla?", "Posso provarla?", "Posso provarlo?", "Posso la provare?", "Provo lei?"),
		mc("What are 'i saldi'?", "In vetrina: SALDI -50%", "The sales", "The sales", "The changing rooms", "The receipts", "The shop assistants"),
		mc("It's too expensive.", "Price", "È troppo caro.", "È troppo caro.", "È molto economico.", "È troppo carro.", "Costa poco."),
		tr("Translate", "I'm looking for a black jacket", "Cerco una giacca nera"),
		fill("Agreement", "Le scarpe sono ___. (expensive — caro)", "care"),
		speak("Cora", "Buongiorno, cerco una maglietta blu, taglia media."),
	)
	addVocab(db, l,
		vw("Quanto costa?", "How much is it?", "Quanto costa questa borsa?", "How much is this bag?", "Cora"),
		vw("la taglia", "size (clothes)", "Che taglia porta?", "What size do you take?", "Cora"),
		vw("provare", "to try (on)", "Posso provarlo?", "Can I try it on?", "Mira"),
		vw("caro / economico", "expensive / cheap", "Questo negozio è caro.", "This shop is expensive.", "Cora"),
		vw("i saldi", "the sales", "Compro tutto ai saldi.", "I buy everything in the sales.", "Blaze"),
	)

	s = addSkillL(db, it, u, "Travel & Weather", "Treni, biglietti, prenotazioni — and che tempo fa?", "Plane", "#F5A623", 19, 245)
	l = addLesson(db, s, "In Viaggio", 1, 20,
		char("Riko", "At the station: un biglietto di andata e ritorno (return ticket) per Firenze; Da che binario parte? (which platform?). Validate regional tickets before boarding! Booking: Vorrei prenotare una camera doppia. Weather uses fare: fa caldo / freddo / bel tempo; piove (it's raining), nevica, c'è il sole, c'è vento."),
		mc("A return ticket to Florence, please.", "At the ticket office", "Un biglietto di andata e ritorno per Firenze, per favore.", "Un biglietto di andata e ritorno per Firenze, per favore.", "Un biglietto solo andata da Firenze.", "Una camera per Firenze.", "Un binario per Firenze."),
		mc("Which platform does it leave from?", "binario", "Da che binario parte?", "Da che binario parte?", "A che ora parte?", "Quanto costa il binario?", "Dov'è il treno ieri?"),
		mc("It's cold.", "Weather", "Fa freddo.", "Fa freddo.", "È freddo il tempo io.", "Ho freddo il tempo.", "C'è freddo fa."),
		mc("It's raining.", "Weather", "Piove.", "Piove.", "Fa pioggia.", "È piovere.", "Ha piovuto domani."),
		tr("Translate", "I would like to book a double room", "Vorrei prenotare una camera doppia"),
		write("Write a short postcard in Italian (about 40 words) from a holiday: where you are, what the weather is like, what you did yesterday (passato prossimo) and a greeting.",
			"Cara Giulia, sono a Napoli con la mia famiglia. Fa molto caldo e c'è sempre il sole. Ieri abbiamo visitato Pompei e abbiamo mangiato una pizza buonissima. Domani andiamo a Capri in barca. Un abbraccio, Anna"),
	)
	addVocab(db, l,
		vw("il biglietto", "ticket", "Ho comprato il biglietto online.", "I bought the ticket online.", "Riko"),
		vw("il binario", "platform", "Il treno parte dal binario tre.", "The train leaves from platform three.", "Riko"),
		vw("prenotare", "to book", "Ho prenotato un tavolo.", "I've booked a table.", "Cora"),
		vw("Che tempo fa?", "What's the weather like?", "Che tempo fa a Milano?", "What's the weather like in Milan?", "Zephyr"),
		vw("piove / nevica", "it's raining / snowing", "Oggi piove.", "It's raining today.", "Zephyr"),
	)
}
