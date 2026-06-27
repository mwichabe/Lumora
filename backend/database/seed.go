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
}

func seedCharacters(db *gorm.DB) {
	characters := []models.Character{
		{Name: "Lumora", Species: "Fennec Fox", Role: "Your Guide", Personality: "Curious, playful, endlessly supportive.", Color: "#6C3FC5", Emoji: "🦊"},
		{Name: "Professor Finch", Species: "Tawny Owl", Role: "Grammar Teacher", Personality: "Strict but secretly warm. Sighs dramatically at errors.", Color: "#8B6F47", Emoji: "🦉"},
		{Name: "Cora", Species: "Octopus", Role: "Vocabulary Friend", Personality: "Goofy, chaotic, can't stop making puns.", Color: "#00C2A8", Emoji: "🐙"},
		{Name: "Blaze", Species: "Fire Salamander", Role: "Speaking Coach", Personality: "Hypes you up. MORE FIRE!", Color: "#FF5C5C", Emoji: "🦎"},
		{Name: "Mira", Species: "Snow Leopard", Role: "Listening Guide", Personality: "Serene and wise. Loves music and poetry.", Color: "#9090A0", Emoji: "🐆"},
		{Name: "Riko", Species: "Red Panda", Role: "Your Rival", Personality: "Smug but lovable. Secretly rooting for you.", Color: "#F5A623", Emoji: "🐼"},
		{Name: "Zephyr", Species: "Wind Spirit", Role: "Writing Mentor", Personality: "Philosophical, poetic, occasionally pretentious.", Color: "#17A3DD", Emoji: "🦅"},
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
