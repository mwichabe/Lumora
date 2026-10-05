package database

import (
	"gorm.io/gorm"

	"lumora/backend/models"
)

// Italian listening and reading sessions — at least one of each per level.

const (
	itA1 = "A1 · Primi Passi — First Steps"
	itA2 = "A2 · La Vita Quotidiana — Everyday Life"
	itB1 = "B1 · Racconti e Opinioni — Stories & Opinions"
	itB2 = "B2 · L'Italiano in Contesto — Italian in Context"
	itC1 = "C1 · Sfumature — Nuance & Style"
	itC2 = "C2 · Padronanza — Mastery"
)

func seedItalianListening(db *gorm.DB) {
	const it = "it"
	finch := "Professor Finch"

	addListeningL(db, it, itA1, "Piacere! — Nice to Meet You",
		"Professor Finch meets a new student. Listen, then answer.", 1, 20,
		[]models.ListeningMatch{
			lm("Come ti chiami?", "What's your name?"), lm("Di dove sei?", "Where are you from?"),
			lm("Piacere!", "Nice to meet you!"), lm("Quanti anni hai?", "How old are you?"),
		},
		[]models.ListeningLine{
			ln(finch, "Ciao! Benvenuta al corso. Come ti chiami?", "Hi! Welcome to the course. What's your name?"),
			ln("Lumora", "Mi chiamo Lumora. Piacere!", "My name is Lumora. Nice to meet you!"),
			ln(finch, "Piacere mio. Di dove sei?", "The pleasure is mine. Where are you from?"),
			ln("Lumora", "Sono keniota, di Nairobi, ma adesso abito a Bologna. Ho ventidue anni.", "I'm Kenyan, from Nairobi, but now I live in Bologna. I'm twenty-two."),
		},
		[]models.ListeningQuestion{
			lq("Where is Lumora from?", "Nairobi", "Nairobi", "Bologna", "Rome", "Milan"),
			lq("Where does she live now?", "Bologna", "Bologna", "Nairobi", "Florence", "Naples"),
			lq("How old is she?", "22", "22", "12", "20", "32"),
		},
	)
	addListeningL(db, it, itA1, "Al Bar — At the Café",
		"Cora orders breakfast. Listen, then answer.", 2, 20,
		[]models.ListeningMatch{
			lm("Vorrei…", "I'd like…"), lm("cornetto", "croissant"),
			lm("Quant'è?", "How much is it?"), lm("lo scontrino", "the receipt"),
		},
		[]models.ListeningLine{
			ln("Riko", "Buongiorno! Cosa prende?", "Good morning! What will you have?"),
			ln("Cora", "Vorrei un cappuccino e un cornetto alla marmellata, per favore.", "I'd like a cappuccino and a jam croissant, please."),
			ln("Riko", "Subito. Altro?", "Right away. Anything else?"),
			ln("Cora", "No, grazie. Quant'è?", "No, thanks. How much is it?"),
			ln("Riko", "Due euro e cinquanta. Ecco lo scontrino.", "Two euros fifty. Here's the receipt."),
		},
		[]models.ListeningQuestion{
			lq("What does Cora drink?", "A cappuccino", "A cappuccino", "An espresso", "A tea", "An orange juice"),
			lq("What kind of croissant does she order?", "Jam", "Jam", "Chocolate", "Plain", "Custard"),
			lq("How much does she pay?", "€2.50", "€2.50", "€5.20", "€2.15", "€1.50"),
		},
	)
	addListeningL(db, it, itA2, "Il Weekend a Firenze — A Weekend in Florence",
		"Zephyr talks about a trip. Listen, then answer.", 3, 22,
		[]models.ListeningMatch{
			lm("sono andato", "I went"), lm("abbiamo visitato", "we visited"),
			lm("pioveva", "it was raining"), lm("tornerò", "I'll come back"),
		},
		[]models.ListeningLine{
			ln("Zephyr", "Il fine settimana scorso sono andato a Firenze con mia sorella. Siamo partiti sabato mattina in treno.", "Last weekend I went to Florence with my sister. We left on Saturday morning by train."),
			ln("Zephyr", "Abbiamo visitato gli Uffizi e abbiamo mangiato la bistecca alla fiorentina.", "We visited the Uffizi and ate Florentine steak."),
			ln("Zephyr", "Domenica pioveva, quindi non siamo saliti sul campanile di Giotto. Ma tornerò in primavera!", "On Sunday it was raining, so we didn't climb Giotto's bell tower. But I'll come back in spring!"),
		},
		[]models.ListeningQuestion{
			lq("Who did he go with?", "His sister", "His sister", "His brother", "His friends", "His parents"),
			lq("How did they travel?", "By train", "By train", "By car", "By plane", "By bus"),
			lq("What did they eat?", "Florentine steak", "Florentine steak", "Pizza", "Seafood pasta", "Gelato only"),
			lq("Why didn't they climb the bell tower?", "It was raining", "It was raining", "It was closed", "They were tired", "It was too expensive"),
		},
	)
	addListeningL(db, it, itB1, "Dal Medico — At the Doctor's",
		"A patient visits the doctor. Listen, then answer.", 4, 24,
		[]models.ListeningMatch{
			lm("mi fa male", "it hurts"), lm("la febbre", "temperature"),
			lm("la ricetta", "prescription"), lm("si riposi", "rest (formal)"),
		},
		[]models.ListeningLine{
			ln("Nana", "Buongiorno, mi dica. Che cosa si sente?", "Good morning, tell me. How are you feeling?"),
			ln("Blaze", "Da tre giorni mi fa male la gola e stanotte ho avuto la febbre.", "My throat has hurt for three days and last night I had a temperature."),
			ln("Nana", "Vediamo… È un'infiammazione. Le scrivo la ricetta per un antibiotico: lo prenda due volte al giorno per una settimana.", "Let's see… It's an inflammation. I'll write you a prescription for an antibiotic: take it twice a day for a week."),
			ln("Blaze", "Posso andare al lavoro?", "Can I go to work?"),
			ln("Nana", "Meglio di no. Si riposi almeno fino a venerdì e beva molta acqua.", "Better not. Rest at least until Friday and drink plenty of water."),
		},
		[]models.ListeningQuestion{
			lq("How long has his throat hurt?", "Three days", "Three days", "One day", "A week", "Since this morning"),
			lq("What does the doctor prescribe?", "An antibiotic", "An antibiotic", "Nothing", "Vitamins", "Cough syrup only"),
			lq("How often should he take it?", "Twice a day for a week", "Twice a day for a week", "Once a day for three days", "Three times a day", "Only when it hurts"),
			lq("What does the doctor advise about work?", "Rest at least until Friday", "Rest at least until Friday", "Go back tomorrow", "Work from home", "Change jobs"),
		},
	)
	addListeningL(db, it, itB2, "Il Telegiornale — The News",
		"A news bulletin. Listen, then answer.", 5, 26,
		[]models.ListeningMatch{
			lm("lo sciopero", "strike"), lm("i disagi", "disruption"),
			lm("i sindacati", "trade unions"), lm("il governo", "government"),
		},
		[]models.ListeningLine{
			ln("Mira", "Buonasera. Apriamo con lo sciopero nazionale dei trasporti previsto per venerdì prossimo.", "Good evening. We open with the national transport strike planned for next Friday."),
			ln("Mira", "I sindacati chiedono aumenti salariali e maggiori investimenti nella sicurezza. Sono attesi forti disagi soprattutto nelle grandi città.", "The unions are demanding pay rises and more investment in safety. Major disruption is expected, especially in the big cities."),
			ln("Mira", "Il governo ha convocato le parti sociali per domani, nella speranza di evitare la protesta. Saranno comunque garantiti i servizi minimi nelle fasce orarie di punta.", "The government has called the social partners to a meeting tomorrow, hoping to avoid the protest. Minimum services will in any case be guaranteed at peak times."),
		},
		[]models.ListeningQuestion{
			lq("When is the strike planned?", "Next Friday", "Next Friday", "Tomorrow", "Next Monday", "This weekend"),
			lq("What are the unions asking for?", "Pay rises and more investment in safety", "Pay rises and more investment in safety", "Fewer working hours only", "New trains", "Lower ticket prices"),
			lq("What is the government doing?", "Meeting the social partners tomorrow", "Meeting the social partners tomorrow", "Cancelling the strike by law", "Ignoring the unions", "Raising taxes"),
			lq("What will be guaranteed?", "Minimum services at peak times", "Minimum services at peak times", "Free travel", "Full service all day", "Nothing at all"),
		},
	)
	addListeningL(db, it, itC1, "Lezione: Dante e la Lingua Italiana — Lecture: Dante and Italian",
		"An excerpt from a university lecture. Listen, then answer.", 6, 28,
		[]models.ListeningMatch{
			lm("il volgare", "the vernacular"), lm("il latino", "Latin"),
			lm("il fiorentino", "Florentine"), lm("l'unità", "unification"),
		},
		[]models.ListeningLine{
			ln(finch, "Quando Dante scrisse la Commedia, la lingua della cultura era il latino. La sua scelta di scrivere in volgare fiorentino fu rivoluzionaria.", "When Dante wrote the Comedy, the language of culture was Latin. His choice to write in the Florentine vernacular was revolutionary."),
			ln(finch, "Nei secoli successivi, grazie anche a Petrarca e Boccaccio, il fiorentino divenne il modello della lingua letteraria.", "In the following centuries, thanks also to Petrarch and Boccaccio, Florentine became the model for the literary language."),
			ln(finch, "Tuttavia, al momento dell'Unità d'Italia, nel 1861, si stima che solo una piccola percentuale della popolazione parlasse italiano: la maggioranza usava i dialetti.", "However, at the time of Italian unification in 1861, it is estimated that only a small percentage of the population spoke Italian: the majority used dialects."),
		},
		[]models.ListeningQuestion{
			lq("What was the language of culture in Dante's time?", "Latin", "Latin", "Florentine", "French", "Greek"),
			lq("Why was Dante's choice revolutionary?", "He wrote in the Florentine vernacular", "He wrote in the Florentine vernacular", "He wrote in Latin", "He wrote in prose", "He wrote in French"),
			lq("Which writers also helped make Florentine the model?", "Petrarch and Boccaccio", "Petrarch and Boccaccio", "Manzoni and Calvino", "Virgil and Ovid", "Pirandello and Eco"),
			lq("What was the situation in 1861?", "Only a small percentage spoke Italian", "Only a small percentage spoke Italian", "Everyone spoke Italian", "Latin was still spoken", "Dialects had disappeared"),
		},
	)
	addListeningL(db, it, itC2, "A Roma — In Rome",
		"Two Romans chat informally. Listen, then answer.", 7, 30,
		[]models.ListeningMatch{
			lm("daje", "come on"), lm("mo'", "now"),
			lm("annamo", "let's go (romanesco)"), lm("'na cifra", "loads (slang)"),
		},
		[]models.ListeningLine{
			ln("Blaze", "Aò, che fai stasera? Annamo a magnà da Checco?", "Hey, what are you doing tonight? Shall we go and eat at Checco's?"),
			ln("Pip", "Daje! Però mo' sto a lavorà, finisco alle otto.", "Come on! But I'm working right now, I finish at eight."),
			ln("Blaze", "Tranquillo, t'aspetto. Ce sta 'na cifra de gente il venerdì, quindi prenoto io.", "No worries, I'll wait for you. There are loads of people on Fridays, so I'll book."),
			ln("Pip", "Perfetto, ce vediamo là!", "Perfect, see you there!"),
		},
		[]models.ListeningQuestion{
			lq("What do they plan to do?", "Eat at Checco's", "Eat at Checco's", "Watch a match", "Go to the cinema", "Go dancing"),
			lq("Why can't Pip go right away?", "He's working until eight", "He's working until eight", "He's ill", "He has no money", "He's on a train"),
			lq("Why will Blaze book?", "It's crowded on Fridays", "It's crowded on Fridays", "It's Pip's birthday", "It's a new restaurant", "It's cheaper online"),
			lq("Which variety is this?", "Romanesco (Roman dialect)", "Romanesco (Roman dialect)", "Neapolitan", "Standard formal Italian", "Milanese"),
		},
	)
}

func seedItalianReading(db *gorm.DB) {
	const it = "it"

	addReadingL(db, it, itA1, "La Mia Famiglia — My Family",
		"Giulia introduces herself. Read, then answer.", 1, 20,
		[]models.ReadingLine{
			rl("Mi chiamo Giulia e ho diciannove anni.", "My name is Giulia and I'm nineteen."),
			rl("Abito a Torino con i miei genitori e mio fratello Paolo.", "I live in Turin with my parents and my brother Paolo."),
			rl("Mio padre è ingegnere e mia madre è insegnante.", "My father is an engineer and my mother is a teacher."),
			rl("Studio lingue all'università. Mi piace molto viaggiare.", "I study languages at university. I really like travelling."),
		},
		[]models.ReadingQuestion{
			rq("How old is Giulia?", "19", "19", "9", "17", "29"),
			rq("Where does she live?", "Turin", "Turin", "Rome", "Milan", "Naples"),
			rq("What is her mother's job?", "Teacher", "Teacher", "Engineer", "Doctor", "Student"),
			rq("What does Giulia study?", "Languages", "Languages", "Engineering", "Medicine", "History"),
		},
	)
	addReadingL(db, it, itA1, "Al Mercato — At the Market",
		"A shopping list. Read, then answer.", 2, 20,
		[]models.ReadingLine{
			rl("Oggi vado al mercato. Compro un chilo di pomodori, due peperoni e il basilico.", "Today I'm going to the market. I'm buying a kilo of tomatoes, two peppers and basil."),
			rl("Il fruttivendolo è molto gentile e le verdure sono fresche.", "The greengrocer is very kind and the vegetables are fresh."),
			rl("Spendo dieci euro. Stasera preparo la pasta al pomodoro per i miei amici.", "I spend ten euros. Tonight I'm making tomato pasta for my friends."),
		},
		[]models.ReadingQuestion{
			rq("How many tomatoes does the writer buy?", "A kilo", "A kilo", "Two", "Half a kilo", "Ten"),
			rq("How much does the writer spend?", "€10", "€10", "€2", "€12", "€20"),
			rq("What will the writer cook tonight?", "Tomato pasta", "Tomato pasta", "Pizza", "Soup", "Risotto"),
		},
	)
	addReadingL(db, it, itA2, "Una Mail da Napoli — An Email from Naples",
		"Anna writes to a friend. Read, then answer.", 3, 22,
		[]models.ReadingLine{
			rl("Cara Sara, come stai? Io sono a Napoli da una settimana per un corso di italiano.", "Dear Sara, how are you? I've been in Naples for a week for an Italian course."),
			rl("Ogni mattina vado a lezione e il pomeriggio visito la città. Ieri sono andata a Pompei: era incredibile!", "Every morning I go to class and in the afternoon I visit the city. Yesterday I went to Pompeii: it was incredible!"),
			rl("La famiglia che mi ospita è simpaticissima e la signora cucina benissimo.", "The family I'm staying with is very nice and the lady cooks wonderfully."),
			rl("Torno a casa sabato prossimo. Un abbraccio, Anna", "I'm going home next Saturday. A hug, Anna"),
		},
		[]models.ReadingQuestion{
			rq("Why is Anna in Naples?", "For an Italian course", "For an Italian course", "For work", "On holiday with family", "To visit Sara"),
			rq("Where did she go yesterday?", "Pompeii", "Pompeii", "Capri", "Rome", "The beach"),
			rq("Who is she staying with?", "A host family", "A host family", "Sara", "In a hotel", "Friends from school"),
			rq("When is she going home?", "Next Saturday", "Next Saturday", "Tomorrow", "In a month", "Last Saturday"),
		},
	)
	addReadingL(db, it, itB1, "Il Rito del Caffè — The Coffee Ritual",
		"A short article. Read, then answer.", 4, 24,
		[]models.ReadingLine{
			rl("Per gli italiani, il caffè non è solo una bevanda, ma un momento sociale.", "For Italians, coffee isn't just a drink but a social moment."),
			rl("Si beve in piedi al bancone, spesso in pochi secondi, scambiando due parole con il barista.", "It's drunk standing at the counter, often in a few seconds, exchanging a few words with the barista."),
			rl("A Napoli esiste la tradizione del 'caffè sospeso': chi può paga un caffè in più per una persona che non può permetterselo.", "In Naples there's the tradition of the 'suspended coffee': whoever can pays for an extra coffee for someone who can't afford one."),
			rl("Questa usanza, nata nel dopoguerra, si è diffusa in molti paesi del mondo.", "This custom, born after the war, has spread to many countries around the world."),
		},
		[]models.ReadingQuestion{
			rq("How is coffee usually drunk at the bar?", "Standing at the counter", "Standing at the counter", "Sitting outside for hours", "At home only", "Taken away in a paper cup"),
			rq("What is a 'caffè sospeso'?", "A coffee paid in advance for someone who can't afford it", "A coffee paid in advance for someone who can't afford it", "A coffee you drink later", "A cold coffee", "A coffee with liquor"),
			rq("Where does the tradition come from?", "Naples", "Naples", "Milan", "Venice", "Turin"),
			rq("When did the custom begin?", "After the war", "After the war", "In the Middle Ages", "Last year", "In the 1990s"),
		},
	)
	addReadingL(db, it, itB2, "La Fuga dei Cervelli — The Brain Drain",
		"An opinion column. Read, then answer.", 5, 26,
		[]models.ReadingLine{
			rl("Ogni anno migliaia di giovani laureati lasciano l'Italia per lavorare all'estero.", "Every year thousands of young graduates leave Italy to work abroad."),
			rl("Le cause sono note: stipendi bassi, precarietà e poche opportunità di carriera, soprattutto nel Mezzogiorno.", "The causes are well known: low salaries, insecure contracts and few career opportunities, especially in the South."),
			rl("Benché alcuni vedano in questo fenomeno un'occasione di crescita personale, il Paese perde competenze in cui ha investito per anni.", "Although some see this phenomenon as a chance for personal growth, the country loses skills it has invested in for years."),
			rl("Se non si interverrà con politiche serie, il divario con gli altri Paesi europei è destinato ad aumentare.", "Unless serious policies are adopted, the gap with other European countries is set to grow."),
		},
		[]models.ReadingQuestion{
			rq("Who is leaving Italy?", "Young graduates", "Young graduates", "Retirees", "Tourists", "Farmers"),
			rq("Which is NOT given as a cause?", "High taxes on travel", "High taxes on travel", "Low salaries", "Insecure contracts", "Few career opportunities"),
			rq("What does the country lose, according to the writer?", "Skills it has invested in", "Skills it has invested in", "Tourism", "Its language", "Its pensions"),
			rq("What does the writer warn about?", "The gap with other European countries will grow without serious policies", "The gap with other European countries will grow without serious policies", "Italy will close its universities", "Graduates will all return soon", "Salaries will rise automatically"),
		},
	)
	addReadingL(db, it, itC1, "Da 'Il Barone Rampante' — From 'The Baron in the Trees'",
		"A retelling of the premise of Italo Calvino's novel. Read, then answer.", 6, 28,
		[]models.ReadingLine{
			rl("Era il giugno del 1767 quando mio fratello Cosimo, che aveva dodici anni, pranzò con noi per l'ultima volta.", "It was June 1767 when my brother Cosimo, who was twelve, had lunch with us for the last time."),
			rl("Rifiutò il piatto di lumache, si alzò da tavola e salì sull'elce del giardino.", "He refused the plate of snails, got up from the table and climbed the holm oak in the garden."),
			rl("Da quel giorno non toccò più terra: visse tutta la vita sugli alberi, senza però mai smettere di occuparsi del mondo di sotto.", "From that day he never touched the ground again: he lived his whole life in the trees, yet never stopped caring about the world below."),
		},
		[]models.ReadingQuestion{
			rq("What tense dominates the passage?", "Passato remoto", "Passato remoto", "Passato prossimo", "Future", "Present subjunctive"),
			rq("What did Cosimo refuse?", "A plate of snails", "A plate of snails", "To go to school", "To marry", "To speak Italian"),
			rq("Where did he live from that day?", "In the trees", "In the trees", "In another city", "In a monastery", "On a boat"),
			rq("What is suggested about his relationship with the world?", "He still cared about the world below", "He still cared about the world below", "He forgot everyone", "He hated his family", "He never spoke again"),
		},
	)
	addReadingL(db, it, itC2, "Comunicazione dell'Ufficio Anagrafe — A Registry Office Notice",
		"A bureaucratic notice. Read, then answer.", 7, 30,
		[]models.ReadingLine{
			rl("Si comunica alla cittadinanza che, a decorrere dal giorno 1° marzo, l'Ufficio Anagrafe osserverà il seguente orario di apertura al pubblico: lunedì-venerdì, ore 9.00-12.30.", "Citizens are informed that, with effect from 1 March, the Registry Office will be open to the public as follows: Monday to Friday, 9.00 to 12.30."),
			rl("Si fa presente che il rilascio della carta d'identità elettronica avverrà esclusivamente previo appuntamento, da prenotarsi tramite il portale istituzionale.", "Please note that electronic identity cards will be issued exclusively by appointment, to be booked via the official portal."),
			rl("Si prega l'utenza di presentarsi muniti di documento di riconoscimento in corso di validità e della ricevuta di avvenuto pagamento.", "Users are requested to bring a valid identity document and the receipt confirming payment."),
		},
		[]models.ReadingQuestion{
			rq("When does the new timetable start?", "1 March", "1 March", "1 May", "Next Monday", "Immediately"),
			rq("How can you get an electronic identity card?", "Only by booking an appointment online", "Only by booking an appointment online", "By walking in any afternoon", "By post", "At the police station"),
			rq("What does 'previo appuntamento' mean?", "By prior appointment", "By prior appointment", "Without appointment", "After the appointment", "Previous appointment cancelled"),
			rq("What must you bring?", "A valid ID document and the payment receipt", "A valid ID document and the payment receipt", "Two photos only", "A birth certificate", "Nothing"),
		},
	)
}
