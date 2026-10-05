package controllers

// italianPapers is the hand-authored Italian proficiency-exam bank, keyed by
// CEFR level and modelled on the four-skill format of CILS / CELI. Questions
// are in English through B1 and in Italian from B2. Writing targets (in
// words) follow the German bank.
var italianPapers = map[string]paperContent{

	// ---------------- A1 ----------------
	"A1": {
		Listening: PaperListening{
			Title: "La Giornata di Marco — Marco's Day",
			Lines: []PaperLine{
				{Character: "Cora", Text: "Ciao! Mi chiamo Marco e sono uno studente. Abito a Verona con la mia famiglia.", Translation: "Hi! My name is Marco and I'm a student. I live in Verona with my family."},
				{Character: "Cora", Text: "La mattina mi sveglio alle sette e bevo un caffè. Vado all'università in bicicletta.", Translation: "In the morning I wake up at seven and drink a coffee. I go to university by bicycle."},
				{Character: "Cora", Text: "La sera gioco a calcio con gli amici e poi mangio con i miei genitori.", Translation: "In the evening I play football with friends and then eat with my parents."},
			},
			Questions: []PaperQuestion{
				{Question: "Where does Marco live?", Options: []string{"Verona", "Venice", "Rome", "Milan"}, CorrectAnswer: "Verona"},
				{Question: "What time does he wake up?", Options: []string{"7:00", "8:00", "6:00", "9:00"}, CorrectAnswer: "7:00"},
				{Question: "How does he go to university?", Options: []string{"By bicycle", "By bus", "On foot", "By car"}, CorrectAnswer: "By bicycle"},
				{Question: "What does he do in the evening?", Options: []string{"Plays football with friends", "Studies at the library", "Works in a bar", "Watches TV alone"}, CorrectAnswer: "Plays football with friends"},
			},
		},
		Reading: PaperReading{
			Title: "La Mia Amica Sofia — My Friend Sofia",
			Paragraphs: []string{
				"Ho un'amica italiana. Si chiama Sofia e ha venticinque anni.",
				"Sofia è infermiera e lavora in un ospedale a Bari. Ha un cane piccolo e bianco.",
				"Il sabato andiamo insieme al mare e mangiamo un gelato al pistacchio.",
			},
			Questions: []PaperQuestion{
				{Question: "How old is Sofia?", Options: []string{"25", "15", "35", "20"}, CorrectAnswer: "25"},
				{Question: "What is her job?", Options: []string{"Nurse", "Teacher", "Doctor", "Student"}, CorrectAnswer: "Nurse"},
				{Question: "What pet does she have?", Options: []string{"A small white dog", "A black cat", "A big dog", "A bird"}, CorrectAnswer: "A small white dog"},
				{Question: "What do they eat on Saturdays?", Options: []string{"Pistachio ice cream", "Pizza", "Pasta", "Fish"}, CorrectAnswer: "Pistachio ice cream"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write a short self-introduction in Italian: your name, age, where you're from, where you live, your family and what you like.",
			MinWords: 30,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Buongiorno! Mi chiamo Anna, sono keniota e studio l'italiano.",
			Speaker:     "Lumora",
			Translation: "Good morning! My name is Anna, I'm Kenyan and I'm studying Italian.",
		},
	},

	// ---------------- A2 ----------------
	"A2": {
		Listening: PaperListening{
			Title: "Una Vacanza in Sicilia — A Holiday in Sicily",
			Lines: []PaperLine{
				{Character: "Zephyr", Text: "L'estate scorsa sono andato in Sicilia con due amici. Siamo partiti da Roma in aereo.", Translation: "Last summer I went to Sicily with two friends. We left from Rome by plane."},
				{Character: "Zephyr", Text: "Abbiamo visitato Palermo e Siracusa e abbiamo mangiato gli arancini tutti i giorni!", Translation: "We visited Palermo and Syracuse and ate arancini every day!"},
				{Character: "Zephyr", Text: "L'ultimo giorno faceva troppo caldo, quindi non siamo andati sull'Etna. Ci torneremo l'anno prossimo.", Translation: "On the last day it was too hot, so we didn't go up Etna. We'll go back next year."},
			},
			Questions: []PaperQuestion{
				{Question: "How did they travel to Sicily?", Options: []string{"By plane", "By train", "By ferry", "By car"}, CorrectAnswer: "By plane"},
				{Question: "Which cities did they visit?", Options: []string{"Palermo and Syracuse", "Catania and Messina", "Rome and Naples", "Palermo and Catania"}, CorrectAnswer: "Palermo and Syracuse"},
				{Question: "Why didn't they go up Etna?", Options: []string{"It was too hot", "It was raining", "It was closed", "They had no time"}, CorrectAnswer: "It was too hot"},
				{Question: "What will they do next year?", Options: []string{"Go back to Sicily", "Go to Sardinia", "Stay at home", "Move to Sicily"}, CorrectAnswer: "Go back to Sicily"},
			},
		},
		Reading: PaperReading{
			Title: "Annuncio: Cercasi Cameriere — Notice: Waiter Wanted",
			Paragraphs: []string{
				"Il ristorante 'Da Gino' cerca un cameriere o una cameriera per la stagione estiva, da giugno a settembre.",
				"L'orario è dalle 18.00 alle 24.00, dal martedì alla domenica. È richiesta la conoscenza dell'inglese.",
				"Gli interessati possono inviare il curriculum entro il 30 aprile all'indirizzo email del ristorante.",
			},
			Questions: []PaperQuestion{
				{Question: "How long is the job?", Options: []string{"June to September", "All year", "One month", "April to June"}, CorrectAnswer: "June to September"},
				{Question: "What are the working hours?", Options: []string{"6 p.m. to midnight", "6 a.m. to noon", "Noon to 6 p.m.", "8 a.m. to 4 p.m."}, CorrectAnswer: "6 p.m. to midnight"},
				{Question: "Which day is the restaurant closed to this worker?", Options: []string{"Monday", "Sunday", "Tuesday", "Saturday"}, CorrectAnswer: "Monday"},
				{Question: "How should people apply?", Options: []string{"Send a CV by email by 30 April", "Phone the restaurant", "Visit on Monday", "Send a letter in June"}, CorrectAnswer: "Send a CV by email by 30 April"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write in Italian about your last weekend: where you went, what you did and how it was. Use the passato prossimo with both avere and essere, and the imperfetto at least once.",
			MinWords: 55,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Ieri sono andata al mercato e ho comprato frutta e verdura.",
			Speaker:     "Cora",
			Translation: "Yesterday I went to the market and bought fruit and vegetables.",
		},
	},

	// ---------------- B1 ----------------
	"B1": {
		Listening: PaperListening{
			Title: "Cambiare Casa — Moving House",
			Lines: []PaperLine{
				{Character: "Mira", Text: "Sto pensando di trasferirmi a Milano perché ho trovato un lavoro lì, ma gli affitti sono altissimi.", Translation: "I'm thinking of moving to Milan because I've found a job there, but rents are very high."},
				{Character: "Blaze", Text: "Potresti cercare casa in un paese vicino e prendere il treno ogni giorno.", Translation: "You could look for a place in a nearby town and take the train every day."},
				{Character: "Mira", Text: "Hai ragione. Se trovassi un appartamento vicino alla stazione, risparmierei molto.", Translation: "You're right. If I found a flat near the station, I'd save a lot."},
				{Character: "Blaze", Text: "Allora sabato andiamo insieme a vedere qualche appartamento!", Translation: "Then on Saturday let's go and look at some flats together!"},
			},
			Questions: []PaperQuestion{
				{Question: "Why does Mira want to move to Milan?", Options: []string{"She found a job there", "Her family lives there", "To study", "Rents are low"}, CorrectAnswer: "She found a job there"},
				{Question: "What is the problem?", Options: []string{"Rents are very high", "There's no train", "She doesn't like Milan", "The job is temporary"}, CorrectAnswer: "Rents are very high"},
				{Question: "What does Blaze suggest?", Options: []string{"Live in a nearby town and commute by train", "Stay where she is", "Buy a car", "Share his flat"}, CorrectAnswer: "Live in a nearby town and commute by train"},
				{Question: "What will they do on Saturday?", Options: []string{"Look at some flats together", "Start the new job", "Go on holiday", "Sign a contract"}, CorrectAnswer: "Look at some flats together"},
			},
		},
		Reading: PaperReading{
			Title: "Il Ritorno ai Borghi — Back to the Villages",
			Paragraphs: []string{
				"Negli ultimi anni molti piccoli borghi italiani si sono svuotati: i giovani si sono trasferiti nelle città in cerca di lavoro.",
				"Alcuni comuni hanno reagito vendendo case abbandonate a un euro, a condizione che gli acquirenti le ristrutturino entro tre anni.",
				"Grazie al lavoro da remoto, sempre più persone scelgono di vivere in questi luoghi, dove la vita è più tranquilla e costa meno.",
			},
			Questions: []PaperQuestion{
				{Question: "Why have many villages emptied?", Options: []string{"Young people moved to cities for work", "Natural disasters", "High taxes", "Tourism"}, CorrectAnswer: "Young people moved to cities for work"},
				{Question: "What have some towns done?", Options: []string{"Sold abandoned houses for one euro", "Built new airports", "Closed their schools", "Banned tourists"}, CorrectAnswer: "Sold abandoned houses for one euro"},
				{Question: "What is the condition for buyers?", Options: []string{"Renovate within three years", "Live there for life", "Open a shop", "Pay extra taxes"}, CorrectAnswer: "Renovate within three years"},
				{Question: "What is bringing people back?", Options: []string{"Remote work", "New factories", "Free cars", "Universities"}, CorrectAnswer: "Remote work"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write in Italian about a change you would like to see in your city. Explain why it matters, what would happen if it were made (se + congiuntivo imperfetto + condizionale), and give your opinion with 'penso che' + congiuntivo.",
			MinWords: 90,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Se avessi più tempo, farei volontariato in una biblioteca.",
			Speaker:     "Mira",
			Translation: "If I had more time, I'd volunteer in a library.",
		},
	},

	// ---------------- B2 ----------------
	"B2": {
		Listening: PaperListening{
			Title: "Riunione: Meno Carta in Ufficio — Meeting: Less Paper in the Office",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "Buongiorno a tutti. Oggi dobbiamo decidere come ridurre il consumo di carta in ufficio.", Translation: "Good morning, everyone. Today we must decide how to reduce paper use in the office."},
				{Character: "Mira", Text: "Propongo che tutti i documenti interni vengano condivisi solo in formato digitale.", Translation: "I propose that all internal documents be shared only in digital format."},
				{Character: "Professor Finch", Text: "Ottima idea. Tuttavia, alcuni nostri clienti anziani preferiscono ancora le lettere cartacee.", Translation: "Excellent idea. However, some of our older clients still prefer paper letters."},
				{Character: "Mira", Text: "In quel caso potremmo continuare a spedire la posta solo a loro, purché lo richiedano.", Translation: "In that case we could keep posting letters to them only, provided they request it."},
			},
			Questions: []PaperQuestion{
				{Question: "Qual è l'obiettivo della riunione?", Options: []string{"Ridurre il consumo di carta", "Assumere nuovo personale", "Aumentare gli stipendi", "Chiudere l'ufficio"}, CorrectAnswer: "Ridurre il consumo di carta"},
				{Question: "Che cosa propone Mira?", Options: []string{"Condividere i documenti interni solo in digitale", "Stampare più documenti", "Eliminare i clienti anziani", "Lavorare solo da casa"}, CorrectAnswer: "Condividere i documenti interni solo in digitale"},
				{Question: "Qual è l'obiezione?", Options: []string{"Alcuni clienti anziani preferiscono le lettere cartacee", "I computer sono pochi", "Il digitale costa troppo", "I dipendenti non sono d'accordo"}, CorrectAnswer: "Alcuni clienti anziani preferiscono le lettere cartacee"},
				{Question: "Quale soluzione viene proposta?", Options: []string{"Spedire la posta solo ai clienti che lo richiedono", "Smettere di usare il digitale", "Insegnare il computer a tutti i clienti", "Stampare tutto come prima"}, CorrectAnswer: "Spedire la posta solo ai clienti che lo richiedono"},
			},
		},
		Reading: PaperReading{
			Title: "Il Turismo di Massa — Mass Tourism",
			Paragraphs: []string{
				"Città come Venezia e Firenze accolgono ogni anno milioni di visitatori. Il turismo rappresenta una fonte di reddito fondamentale per l'economia locale.",
				"Tuttavia, l'eccessivo afflusso ha conseguenze pesanti: affitti alle stelle, residenti costretti a trasferirsi e centri storici trasformati in musei a cielo aperto.",
				"Venezia ha introdotto un contributo d'accesso per i visitatori giornalieri. Benché la misura sia stata criticata, molti ritengono che sia un primo passo verso un turismo più sostenibile.",
			},
			Questions: []PaperQuestion{
				{Question: "Che cosa rappresenta il turismo per queste città?", Options: []string{"Una fonte di reddito fondamentale", "Un problema senza vantaggi", "Un'attività marginale", "Una novità recente"}, CorrectAnswer: "Una fonte di reddito fondamentale"},
				{Question: "Quale conseguenza NON viene citata?", Options: []string{"La chiusura delle università", "Affitti alle stelle", "Residenti costretti a trasferirsi", "Centri storici come musei a cielo aperto"}, CorrectAnswer: "La chiusura delle università"},
				{Question: "Che cosa ha introdotto Venezia?", Options: []string{"Un contributo d'accesso per i visitatori giornalieri", "Il divieto di ingresso ai turisti", "Nuovi hotel", "Biglietti gratuiti"}, CorrectAnswer: "Un contributo d'accesso per i visitatori giornalieri"},
				{Question: "Che cosa significa 'affitti alle stelle'?", Options: []string{"Affitti altissimi", "Affitti con vista", "Affitti economici", "Affitti di lusso solo di notte"}, CorrectAnswer: "Affitti altissimi"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write a formal email in Italian to your local council (Gentile Sindaco / Spettabile Comune…) complaining about rubbish collection in your area, explaining the problem and proposing a solution. Close with Distinti saluti.",
			MinWords: 130,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Le sarei grata se potesse inviarmi maggiori informazioni sul corso.",
			Speaker:     "Mira",
			Translation: "I'd be grateful if you could send me more information about the course.",
		},
	},

	// ---------------- C1 ----------------
	"C1": {
		Listening: PaperListening{
			Title: "Lezione: Lingua e Identità — Lecture: Language and Identity",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "La lingua non è soltanto uno strumento di comunicazione: è il deposito della memoria collettiva di un popolo.", Translation: "Language is not just a tool of communication: it is the repository of a people's collective memory."},
				{Character: "Professor Finch", Text: "Basti pensare a come i dialetti italiani conservino parole e modi di dire che lo standard ha perduto.", Translation: "Just consider how Italian dialects preserve words and sayings that the standard language has lost."},
				{Character: "Professor Finch", Text: "Ciononostante, sarebbe un errore contrapporre dialetto e lingua: la ricchezza dell'italiano deriva proprio dalla loro convivenza.", Translation: "Nonetheless, it would be a mistake to set dialect against language: the richness of Italian derives precisely from their coexistence."},
			},
			Questions: []PaperQuestion{
				{Question: "Secondo il docente, che cos'è la lingua oltre a uno strumento di comunicazione?", Options: []string{"Il deposito della memoria collettiva", "Un insieme di regole", "Una materia scolastica", "Un mezzo commerciale"}, CorrectAnswer: "Il deposito della memoria collettiva"},
				{Question: "Che cosa conservano i dialetti?", Options: []string{"Parole e modi di dire perduti dallo standard", "Solo parole straniere", "Regole grammaticali nuove", "Nulla di importante"}, CorrectAnswer: "Parole e modi di dire perduti dallo standard"},
				{Question: "Che cosa sarebbe un errore, secondo il docente?", Options: []string{"Contrapporre dialetto e lingua", "Studiare i dialetti", "Parlare italiano standard", "Scrivere in dialetto"}, CorrectAnswer: "Contrapporre dialetto e lingua"},
				{Question: "Da che cosa deriva la ricchezza dell'italiano?", Options: []string{"Dalla convivenza di lingua e dialetti", "Dal latino soltanto", "Dalle parole inglesi", "Dalla burocrazia"}, CorrectAnswer: "Dalla convivenza di lingua e dialetti"},
			},
		},
		Reading: PaperReading{
			Title: "L'Elogio della Lentezza — In Praise of Slowness",
			Paragraphs: []string{
				"Viviamo in un'epoca che ha fatto della velocità un valore assoluto: consegne in giornata, risposte immediate, notizie che invecchiano nel giro di poche ore.",
				"Eppure, come ricorda il proverbio, chi va piano va sano e va lontano. Le decisioni prese in fretta finiscono spesso per generare rimpianti, nella vita privata come in politica.",
				"Non si tratta di rinunciare all'efficienza, bensì di imparare a distinguere ciò che richiede rapidità da ciò che esige pazienza e riflessione.",
			},
			Questions: []PaperQuestion{
				{Question: "Secondo l'autore, quale valore domina la nostra epoca?", Options: []string{"La velocità", "La pazienza", "La tradizione", "La solidarietà"}, CorrectAnswer: "La velocità"},
				{Question: "Perché viene citato il proverbio?", Options: []string{"Per sostenere che la fretta genera rimpianti", "Per lodare la velocità", "Per parlare di viaggi", "Per criticare la salute"}, CorrectAnswer: "Per sostenere che la fretta genera rimpianti"},
				{Question: "Che cosa significa 'bensì' nel testo?", Options: []string{"ma piuttosto", "quindi", "inoltre", "perché"}, CorrectAnswer: "ma piuttosto"},
				{Question: "Qual è la conclusione dell'autore?", Options: []string{"Distinguere ciò che richiede rapidità da ciò che esige riflessione", "Rinunciare all'efficienza", "Fare tutto in fretta", "Ignorare i proverbi"}, CorrectAnswer: "Distinguere ciò che richiede rapidità da ciò che esige riflessione"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write an argumentative essay in Italian: 'L'intelligenza artificiale minaccia il lavoro umano?' Present a thesis, arguments, the opposing view (d'altro canto) and a reasoned conclusion, using at least three congiuntivo forms.",
			MinWords: 180,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Alla luce di quanto detto, è innegabile che l'istruzione sia la chiave del progresso.",
			Speaker:     "Professor Finch",
			Translation: "In light of what has been said, it is undeniable that education is the key to progress.",
		},
	},

	// ---------------- C2 ----------------
	"C2": {
		Listening: PaperListening{
			Title: "Discorso di Commemorazione — A Commemoration Speech",
			Lines: []PaperLine{
				{Character: "Zephyr", Text: "Autorità, gentili ospiti, cari concittadini: siamo qui riuniti per ricordare chi, con il proprio coraggio, ha reso libera questa terra.", Translation: "Authorities, distinguished guests, dear fellow citizens: we are gathered here to remember those who, with their courage, made this land free."},
				{Character: "Zephyr", Text: "Furono uomini e donne comuni, che seppero scegliere da che parte stare quando scegliere significava rischiare tutto.", Translation: "They were ordinary men and women, who knew which side to take when choosing meant risking everything."},
				{Character: "Zephyr", Text: "A noi spetta il compito di custodire questa memoria, affinché le generazioni future non dimentichino. Vi ringrazio per l'attenzione.", Translation: "It falls to us to safeguard this memory, so that future generations do not forget. Thank you for your attention."},
			},
			Questions: []PaperQuestion{
				{Question: "Qual è l'occasione del discorso?", Options: []string{"Una commemorazione", "Un matrimonio", "Una premiazione sportiva", "Un'inaugurazione commerciale"}, CorrectAnswer: "Una commemorazione"},
				{Question: "Quale tempo verbale prevale nella seconda frase?", Options: []string{"Passato remoto", "Futuro", "Congiuntivo presente", "Imperfetto soltanto"}, CorrectAnswer: "Passato remoto"},
				{Question: "Chi viene ricordato?", Options: []string{"Uomini e donne comuni che rischiarono tutto", "Solo i politici", "Gli ospiti presenti", "Le generazioni future"}, CorrectAnswer: "Uomini e donne comuni che rischiarono tutto"},
				{Question: "Perché va custodita la memoria, secondo l'oratore?", Options: []string{"Affinché le generazioni future non dimentichino", "Per motivi economici", "Per obbligo di legge", "Per attirare turisti"}, CorrectAnswer: "Affinché le generazioni future non dimentichino"},
			},
		},
		Reading: PaperReading{
			Title: "Avviso dell'Ufficio Tributi — A Tax Office Notice",
			Paragraphs: []string{
				"Si rende noto che, a decorrere dal 1° gennaio, il versamento dell'imposta dovrà essere effettuato esclusivamente mediante modalità telematiche.",
				"I contribuenti impossibilitati all'utilizzo dei servizi online potranno avvalersi dell'assistenza presso gli sportelli comunali, previo appuntamento.",
				"Si invita pertanto l'utenza a prendere visione della documentazione pubblicata sul portale istituzionale, onde evitare l'applicazione di sanzioni.",
			},
			Questions: []PaperQuestion{
				{Question: "Come va pagata l'imposta dal 1° gennaio?", Options: []string{"Solo online", "Solo in contanti allo sportello", "Per posta", "Con assegno"}, CorrectAnswer: "Solo online"},
				{Question: "Che cosa possono fare i contribuenti che non usano internet?", Options: []string{"Chiedere assistenza allo sportello su appuntamento", "Non pagare", "Pagare l'anno dopo", "Inviare una lettera"}, CorrectAnswer: "Chiedere assistenza allo sportello su appuntamento"},
				{Question: "Che cosa significa 'onde evitare'?", Options: []string{"per evitare", "nonostante", "dopo aver evitato", "senza evitare"}, CorrectAnswer: "per evitare"},
				{Question: "Come riscriveresti 'effettuare il versamento' in modo semplice?", Options: []string{"pagare", "versare acqua", "evitare", "ricevere"}, CorrectAnswer: "pagare"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write a formal speech in Italian for a ceremony honouring teachers: open with the appropriate forms of address, reflect on the role of teaching using at least one proverb or literary reference, use congiuntivo and passato remoto where natural, and close formally.",
			MinWords: 230,
		},
		Speaking: PaperSpeaking{
			Phrase:      "È per me un onore prendere la parola in questa giornata così significativa.",
			Speaker:     "Zephyr",
			Translation: "It is an honour for me to take the floor on such a meaningful day.",
		},
	},

	// ---------------- FINAL · comprehensive ----------------
	"FINAL": {
		Listening: PaperListening{
			Title: "Tavola Rotonda: L'Italiano nel Mondo Digitale — Round Table: Italian in the Digital World",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "L'italiano non è mai stato tanto scritto quanto oggi: messaggi, social, chat. Paradossalmente, questo ha avvicinato la lingua scritta a quella parlata.", Translation: "Italian has never been written as much as today: messages, social media, chats. Paradoxically, this has brought the written language closer to the spoken."},
				{Character: "Mira", Text: "D'accordo, ma l'uso massiccio di anglicismi rischia di impoverire il lessico, soprattutto nel linguaggio aziendale.", Translation: "Agreed, but the massive use of anglicisms risks impoverishing the vocabulary, especially in business language."},
				{Character: "Professor Finch", Text: "È un rischio reale. Tuttavia, una lingua viva prende in prestito da sempre: il punto non è vietare, bensì offrire alternative italiane valide e farle conoscere.", Translation: "It's a real risk. However, a living language has always borrowed: the point is not to forbid but to offer valid Italian alternatives and make them known."},
			},
			Questions: []PaperQuestion{
				{Question: "Che effetto ha avuto la comunicazione digitale, secondo il primo relatore?", Options: []string{"Ha avvicinato scritto e parlato", "Ha eliminato lo scritto", "Ha diffuso i dialetti", "Ha ridotto l'uso dell'italiano"}, CorrectAnswer: "Ha avvicinato scritto e parlato"},
				{Question: "Quale rischio segnala Mira?", Options: []string{"L'impoverimento del lessico per gli anglicismi", "La scomparsa delle chat", "L'aumento degli errori di ortografia", "La perdita del congiuntivo"}, CorrectAnswer: "L'impoverimento del lessico per gli anglicismi"},
				{Question: "Qual è la posizione finale del primo relatore?", Options: []string{"Offrire alternative italiane valide invece di vietare", "Vietare tutti gli anglicismi", "Adottare l'inglese", "Ignorare il problema"}, CorrectAnswer: "Offrire alternative italiane valide invece di vietare"},
				{Question: "Che cosa significa 'bensì' in questo contesto?", Options: []string{"ma piuttosto", "inoltre", "quindi", "infatti"}, CorrectAnswer: "ma piuttosto"},
			},
		},
		Reading: PaperReading{
			Title: "Tradurre — On Translating",
			Paragraphs: []string{
				"Tradurre non significa sostituire parole, ma trasferire un intero mondo di significati da una cultura a un'altra.",
				"Si pensi all'espressione 'fare bella figura': nessuna traduzione letterale ne restituisce il peso sociale, legato all'importanza che in Italia si attribuisce all'apparenza e al decoro.",
				"Il buon traduttore, dunque, è colui che conosce non due lingue soltanto, ma due culture, e sa quando conservare l'originale e quando, invece, spiegarlo.",
			},
			Questions: []PaperQuestion{
				{Question: "Secondo il testo, che cosa significa tradurre?", Options: []string{"Trasferire un mondo di significati tra culture", "Sostituire parole", "Riassumere un testo", "Correggere errori"}, CorrectAnswer: "Trasferire un mondo di significati tra culture"},
				{Question: "Perché 'fare bella figura' è difficile da tradurre?", Options: []string{"Ha un peso sociale legato all'apparenza e al decoro", "È un'espressione dialettale", "È troppo lunga", "Non ha significato"}, CorrectAnswer: "Ha un peso sociale legato all'apparenza e al decoro"},
				{Question: "Chi è il buon traduttore?", Options: []string{"Chi conosce due culture, non solo due lingue", "Chi parla più lingue possibili", "Chi usa solo il dizionario", "Chi traduce velocemente"}, CorrectAnswer: "Chi conosce due culture, non solo due lingue"},
				{Question: "Quale forma è 'si pensi'?", Options: []string{"Congiuntivo esortativo", "Passato remoto", "Futuro", "Condizionale"}, CorrectAnswer: "Congiuntivo esortativo"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write an essay in Italian on 'La lingua italiana tra tradizione e globalizzazione'. Develop a nuanced argument across registers, engage a counter-position, cite at least one literary or cultural reference, and draw a reasoned conclusion.",
			MinWords: 300,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Una lingua è viva finché c'è qualcuno che la ama e la parla.",
			Speaker:     "Lumora",
			Translation: "A language is alive as long as someone loves it and speaks it.",
		},
	},
}
