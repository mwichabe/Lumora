package controllers

// germanPapers is the hand-authored German proficiency-exam bank (A1–C2). Each
// level scales up: longer listening passages, denser reading texts, more and
// harder questions (German-language questions from B2 upward), longer required
// writing, and more complex speaking prompts.
var germanPapers = map[string]paperContent{

	// ---------------- FINAL (comprehensive A1→C2 mastery) ----------------
	"FINAL": {
		Listening: PaperListening{
			Title: "Abschlussvortrag: Der Wert der Mehrsprachigkeit",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "Herzlich willkommen zu unserem Abschlussvortrag. Heute möchte ich über den Wert der Mehrsprachigkeit sprechen. In einer zunehmend vernetzten Welt ist die Fähigkeit, mehrere Sprachen zu beherrschen, längst kein Luxus mehr, sondern eine Schlüsselkompetenz.", Translation: ""},
				{Character: "Professor Finch", Text: "Studien zeigen, dass mehrsprachige Menschen nicht nur beruflich im Vorteil sind, sondern auch kognitiv flexibler denken. Wer zwischen Sprachen wechselt, trainiert das Gehirn, Perspektiven zu wechseln und Probleme kreativer zu lösen.", Translation: ""},
				{Character: "Professor Finch", Text: "Gleichwohl wäre es naiv zu behaupten, das Erlernen einer Sprache sei mühelos. Es erfordert Geduld, Disziplin und die Bereitschaft, Fehler zu machen. Doch gerade darin liegt der eigentliche Gewinn: Wer eine Sprache lernt, lernt zugleich, bescheiden und beharrlich zu sein.", Translation: ""},
				{Character: "Professor Finch", Text: "Abschließend lässt sich sagen: Jede neue Sprache öffnet nicht nur Türen, sondern auch Horizonte. Sie ist, um es mit einem Philosophen zu sagen, die Grenze unserer Welt — und wer sie erweitert, erweitert sich selbst.", Translation: ""},
			},
			Questions: []PaperQuestion{
				{Question: "Wie bezeichnet der Redner die Mehrsprachigkeit heute?", Options: []string{"Als Schlüsselkompetenz", "Als Luxus", "Als Modeerscheinung", "Als Pflicht"}, CorrectAnswer: "Als Schlüsselkompetenz"},
				{Question: "Welchen kognitiven Vorteil nennt er?", Options: []string{"Flexibleres Denken", "Ein besseres Gedächtnis allein", "Schnelleres Rechnen", "Weniger Schlafbedarf"}, CorrectAnswer: "Flexibleres Denken"},
				{Question: "Was trainiert, wer zwischen Sprachen wechselt?", Options: []string{"Perspektiven zu wechseln", "Schneller zu sprechen", "Auswendig zu lernen", "Lauter zu reden"}, CorrectAnswer: "Perspektiven zu wechseln"},
				{Question: "Was erfordert das Sprachenlernen laut Redner?", Options: []string{"Geduld und Disziplin", "Nur Talent", "Viel Geld", "Perfekte Aussprache"}, CorrectAnswer: "Geduld und Disziplin"},
				{Question: "Worin liegt der eigentliche Gewinn?", Options: []string{"Bescheiden und beharrlich zu werden", "Schnell reich zu werden", "Nie mehr Fehler zu machen", "Berühmt zu werden"}, CorrectAnswer: "Bescheiden und beharrlich zu werden"},
				{Question: "Wie beschreibt er eine neue Sprache am Ende?", Options: []string{"Als Grenze und Erweiterung unserer Welt", "Als reines Werkzeug", "Als Zeitverschwendung", "Als Modeerscheinung"}, CorrectAnswer: "Als Grenze und Erweiterung unserer Welt"},
			},
		},
		Reading: PaperReading{
			Title: "Lebenslanges Lernen",
			Paragraphs: []string{
				"Die Vorstellung, dass Bildung mit dem Schulabschluss endet, gehört längst der Vergangenheit an. In einer Welt, die sich mit atemberaubender Geschwindigkeit wandelt, wird das lebenslange Lernen zur Notwendigkeit.",
				"Wer heute aufhört zu lernen, riskiert, den Anschluss zu verlieren. Neue Technologien entstehen, ganze Berufsfelder verschwinden, und Wissen, das gestern noch aktuell war, ist morgen womöglich überholt.",
				"Kritiker wenden ein, der ständige Druck, sich weiterzubilden, erzeuge Stress und Erschöpfung. Dieser Einwand ist nicht von der Hand zu weisen. Dennoch geht es beim lebenslangen Lernen nicht um pausenlose Selbstoptimierung, sondern um Neugier.",
				"Denn wer neugierig bleibt, altert nicht im Geiste. Bildung ist, recht verstanden, kein Wettlauf, sondern eine Haltung — die Bereitschaft, sich immer wieder aufs Neue verwundern zu lassen.",
			},
			Questions: []PaperQuestion{
				{Question: "Was gehört laut Text der Vergangenheit an?", Options: []string{"Dass Bildung mit dem Schulabschluss endet", "Dass man überhaupt lernt", "Dass Schulen existieren", "Dass Technik sich ändert"}, CorrectAnswer: "Dass Bildung mit dem Schulabschluss endet"},
				{Question: "Was riskiert, wer aufhört zu lernen?", Options: []string{"Den Anschluss zu verlieren", "Reich zu werden", "Mehr Freizeit", "Bessere Gesundheit"}, CorrectAnswer: "Den Anschluss zu verlieren"},
				{Question: "Was kann mit dem Wissen von gestern geschehen?", Options: []string{"Es ist morgen überholt", "Es bleibt ewig gültig", "Es wird wertvoller", "Es verschwindet nie"}, CorrectAnswer: "Es ist morgen überholt"},
				{Question: "Welchen Einwand bringen Kritiker vor?", Options: []string{"Der ständige Druck erzeuge Stress", "Lernen sei zu billig", "Es gebe zu wenige Schulen", "Niemand wolle lernen"}, CorrectAnswer: "Der ständige Druck erzeuge Stress"},
				{Question: "Worum geht es beim lebenslangen Lernen laut Autor?", Options: []string{"Um Neugier", "Um Selbstoptimierung ohne Pause", "Um Prüfungen", "Um Wettbewerb"}, CorrectAnswer: "Um Neugier"},
				{Question: "Wie versteht der Autor Bildung am Ende?", Options: []string{"Als eine Haltung", "Als einen Wettlauf", "Als eine Pflicht", "Als ein Produkt"}, CorrectAnswer: "Als eine Haltung"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "a comprehensive, well-structured essay in German (300+ words): 'Is it worth learning several languages?' Present clear arguments and counter-arguments, use varied grammar and connectors, and reach a reasoned conclusion.",
			MinWords: 300,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Jede neue Sprache öffnet nicht nur Türen, sondern auch Horizonte; und wer ihre Grenzen erweitert, erweitert zugleich sich selbst.",
			Speaker:     "Lumora",
			Translation: "Every new language opens not only doors but horizons; and whoever widens its limits widens themselves.",
		},
	},

	// ---------------- A1 ----------------
	"A1": {
		Listening: PaperListening{
			Title: "Lenas Tag",
			Lines: []PaperLine{
				{Character: "Cora", Text: "Hallo! Ich heiße Lena und ich wohne in Hamburg. Ich bin dreiundzwanzig Jahre alt und ich bin Studentin. Am Morgen stehe ich um sieben Uhr auf. Ich trinke einen Kaffee und ich esse ein Brötchen mit Käse.", Translation: "Hello! My name is Lena and I live in Hamburg. I am twenty-three years old and I am a student. In the morning I get up at seven. I drink a coffee and eat a roll with cheese."},
				{Character: "Cora", Text: "Um neun Uhr fahre ich mit dem Fahrrad zur Universität. Der Weg dauert zwanzig Minuten. Am Nachmittag lerne ich in der Bibliothek. Am Abend koche ich mit meiner Freundin. Wir essen gern Gemüse und Nudeln.", Translation: "At nine o'clock I ride my bike to the university. The way takes twenty minutes. In the afternoon I study in the library. In the evening I cook with my friend. We like to eat vegetables and pasta."},
			},
			Questions: []PaperQuestion{
				{Question: "Where does Lena live?", Options: []string{"Hamburg", "Berlin", "München", "Köln"}, CorrectAnswer: "Hamburg"},
				{Question: "How old is Lena?", Options: []string{"23", "32", "13", "20"}, CorrectAnswer: "23"},
				{Question: "How does Lena get to the university?", Options: []string{"By bicycle", "By bus", "By car", "On foot"}, CorrectAnswer: "By bicycle"},
				{Question: "What does Lena eat in the morning?", Options: []string{"A roll with cheese", "Soup", "Pizza", "Rice"}, CorrectAnswer: "A roll with cheese"},
			},
		},
		Reading: PaperReading{
			Title: "Eine E-Mail aus Berlin",
			Paragraphs: []string{
				"Hallo Tom, wie geht es dir? Ich bin jetzt in Berlin. Die Stadt ist sehr groß und schön. Ich wohne in einer kleinen Wohnung im Zentrum.",
				"Morgens gehe ich in einen Deutschkurs und nachmittags besuche ich Museen. Das Wetter ist kalt, aber es regnet nicht. Am Wochenende möchte ich einen Park besuchen. Viele Grüße, Anna.",
			},
			Questions: []PaperQuestion{
				{Question: "Where is Anna now?", Options: []string{"In Berlin", "In Hamburg", "In Wien", "In Zürich"}, CorrectAnswer: "In Berlin"},
				{Question: "What does Anna do in the morning?", Options: []string{"She goes to a German course", "She works in an office", "She visits museums", "She sleeps"}, CorrectAnswer: "She goes to a German course"},
				{Question: "How is the weather?", Options: []string{"Cold but no rain", "Warm and sunny", "Rainy", "Snowy"}, CorrectAnswer: "Cold but no rain"},
				{Question: "What does Anna want to do at the weekend?", Options: []string{"Visit a park", "Go to a museum", "Cook dinner", "Take the train"}, CorrectAnswer: "Visit a park"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "a short email to a friend: introduce yourself, say where you live, and describe what you do every day.",
			MinWords: 30,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Guten Morgen! Ich heiße Anna und ich wohne in Berlin.",
			Speaker:     "Lumora",
			Translation: "Good morning! My name is Anna and I live in Berlin.",
		},
	},

	// ---------------- A2 ----------------
	"A2": {
		Listening: PaperListening{
			Title: "Nachricht vom Reisebüro",
			Lines: []PaperLine{
				{Character: "Ricco", Text: "Guten Tag, hier ist das Reisebüro Sonnenschein. Sie haben eine Reise nach Österreich gebucht. Ihr Zug fährt am Samstag um acht Uhr fünfzehn vom Hauptbahnhof ab. Bitte seien Sie zwanzig Minuten früher da.", Translation: "Hello, this is the Sonnenschein travel agency. You have booked a trip to Austria. Your train leaves on Saturday at 8:15 from the main station. Please be there twenty minutes early."},
				{Character: "Ricco", Text: "Das Hotel liegt direkt am See. Das Frühstück ist im Preis inbegriffen, aber das Abendessen müssen Sie extra bezahlen. Wenn Sie Fragen haben, rufen Sie uns bitte an. Wir wünschen Ihnen eine gute Reise. Auf Wiederhören!", Translation: "The hotel is right by the lake. Breakfast is included in the price, but you have to pay extra for dinner. If you have questions, please call us. We wish you a good trip. Goodbye!"},
			},
			Questions: []PaperQuestion{
				{Question: "When does the train leave?", Options: []string{"At 8:15", "At 8:50", "At 9:15", "At 7:20"}, CorrectAnswer: "At 8:15"},
				{Question: "Where is the hotel located?", Options: []string{"By the lake", "In the mountains", "Near the station", "In the city centre"}, CorrectAnswer: "By the lake"},
				{Question: "What is included in the price?", Options: []string{"Breakfast", "Dinner", "Both meals", "Nothing"}, CorrectAnswer: "Breakfast"},
				{Question: "What must the traveller pay extra for?", Options: []string{"Dinner", "Breakfast", "The train ticket", "The phone call"}, CorrectAnswer: "Dinner"},
			},
		},
		Reading: PaperReading{
			Title: "Mein neues Hobby",
			Paragraphs: []string{
				"Seit drei Monaten gehe ich jeden Dienstag in einen Sportverein. Früher hatte ich nach der Arbeit keine Energie mehr, aber eine Kollegin hat mich überredet, es einmal zu versuchen.",
				"Jetzt spiele ich Volleyball mit acht anderen Leuten. Wir trainieren zwei Stunden und danach gehen wir oft zusammen etwas essen. Ich habe schon viele neue Freunde gefunden.",
				"Am Anfang war es schwer, weil ich sehr unsportlich war. Heute fühle ich mich gesünder und ich schlafe auch besser. Nächstes Jahr möchte ich an einem kleinen Turnier teilnehmen.",
			},
			Questions: []PaperQuestion{
				{Question: "How often does the writer go to the sports club?", Options: []string{"Once a week", "Every day", "Twice a week", "Once a month"}, CorrectAnswer: "Once a week"},
				{Question: "Why did the writer start this hobby?", Options: []string{"A colleague convinced them", "A doctor recommended it", "They were bored", "They wanted to win a prize"}, CorrectAnswer: "A colleague convinced them"},
				{Question: "What sport does the writer play?", Options: []string{"Volleyball", "Football", "Tennis", "Basketball"}, CorrectAnswer: "Volleyball"},
				{Question: "What does the writer want to do next year?", Options: []string{"Take part in a small tournament", "Change clubs", "Stop the hobby", "Become a coach"}, CorrectAnswer: "Take part in a small tournament"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "an email to a colleague describing what you did last weekend and your plans for the next one.",
			MinWords: 55,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Am Wochenende fahre ich mit dem Zug nach München, um meine Familie zu besuchen.",
			Speaker:     "Lumora",
			Translation: "At the weekend I take the train to Munich to visit my family.",
		},
	},

	// ---------------- B1 ----------------
	"B1": {
		Listening: PaperListening{
			Title: "Radiobeitrag: Leben in der Stadt",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "Willkommen zu unserem Beitrag über das Leben in deutschen Großstädten. Immer mehr junge Menschen ziehen vom Land in die Stadt, weil sie dort bessere Chancen auf einen Arbeitsplatz und ein Studium haben. Gleichzeitig steigen aber die Mieten, sodass viele sich kaum noch eine Wohnung leisten können.", Translation: ""},
				{Character: "Professor Finch", Text: "Eine mögliche Lösung ist der Ausbau des öffentlichen Nahverkehrs. Wenn Busse und Bahnen günstig und zuverlässig sind, müssen die Menschen nicht unbedingt im teuren Zentrum wohnen. Sie können in einem Vorort leben und trotzdem schnell zur Arbeit kommen.", Translation: ""},
				{Character: "Professor Finch", Text: "Außerdem fordern viele Bürger mehr Grünflächen und Fahrradwege. Studien zeigen, dass Parks nicht nur die Luft verbessern, sondern auch das Wohlbefinden der Menschen erhöhen. Die Stadt der Zukunft muss also lebenswert und bezahlbar zugleich sein.", Translation: ""},
			},
			Questions: []PaperQuestion{
				{Question: "Why do many young people move to the city?", Options: []string{"Better job and study opportunities", "Cheaper housing", "Cleaner air", "Quieter life"}, CorrectAnswer: "Better job and study opportunities"},
				{Question: "What problem does the report mention?", Options: []string{"Rising rents", "Falling salaries", "Too few universities", "Bad weather"}, CorrectAnswer: "Rising rents"},
				{Question: "What solution is suggested?", Options: []string{"Expanding public transport", "Building more offices", "Banning cars completely", "Lowering taxes"}, CorrectAnswer: "Expanding public transport"},
				{Question: "According to studies, what do parks do?", Options: []string{"Improve air and well-being", "Increase rents", "Reduce jobs", "Make cities louder"}, CorrectAnswer: "Improve air and well-being"},
				{Question: "How is the 'city of the future' described?", Options: []string{"Liveable and affordable", "Large and expensive", "Empty and quiet", "Modern but polluted"}, CorrectAnswer: "Liveable and affordable"},
			},
		},
		Reading: PaperReading{
			Title: "Ehrenamt: Zeit, die sich lohnt",
			Paragraphs: []string{
				"In Deutschland engagieren sich Millionen Menschen ehrenamtlich, das heißt, sie arbeiten freiwillig und ohne Bezahlung für andere. Manche helfen in der Feuerwehr, andere betreuen Kinder oder begleiten ältere Menschen beim Einkaufen.",
				"Lukas, 29, erzählt: „Nach meinem Umzug in eine neue Stadt kannte ich niemanden. Über das Ehrenamt habe ich nicht nur sinnvolle Aufgaben gefunden, sondern auch viele nette Leute kennengelernt. Heute leite ich eine Gruppe, die geflüchteten Familien beim Deutschlernen hilft.“",
				"Natürlich kostet das Engagement Zeit, und nicht jeder kann sich neben Beruf und Familie noch zusätzlich verpflichten. Trotzdem sagen die meisten Freiwilligen, dass sie mehr zurückbekommen, als sie geben. Sie fühlen sich gebraucht und sind stolz, etwas Positives für die Gesellschaft zu tun.",
			},
			Questions: []PaperQuestion{
				{Question: "What does 'ehrenamtlich' mean here?", Options: []string{"Working voluntarily without pay", "Working part-time", "Working from home", "Working abroad"}, CorrectAnswer: "Working voluntarily without pay"},
				{Question: "Why did Lukas start volunteering?", Options: []string{"He knew nobody in his new city", "He needed money", "His boss told him to", "He was retired"}, CorrectAnswer: "He knew nobody in his new city"},
				{Question: "What does Lukas's group do now?", Options: []string{"Helps refugee families learn German", "Cleans parks", "Trains firefighters", "Organises concerts"}, CorrectAnswer: "Helps refugee families learn German"},
				{Question: "What drawback of volunteering is mentioned?", Options: []string{"It takes time", "It is dangerous", "It is boring", "It is illegal"}, CorrectAnswer: "It takes time"},
				{Question: "How do most volunteers feel?", Options: []string{"They get back more than they give", "They feel exploited", "They want to stop", "They earn a good salary"}, CorrectAnswer: "They get back more than they give"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "a forum post giving your opinion on whether young people should do voluntary work, with reasons and an example.",
			MinWords: 90,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Meiner Meinung nach ist es wichtig, dass wir die Umwelt schützen und weniger mit dem Auto fahren.",
			Speaker:     "Lumora",
			Translation: "In my opinion it is important that we protect the environment and drive less.",
		},
	},

	// ---------------- B2 ----------------
	"B2": {
		Listening: PaperListening{
			Title: "Diskussion: Homeoffice",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "Seit der Pandemie hat sich die Arbeitswelt grundlegend verändert. Viele Unternehmen haben festgestellt, dass ihre Mitarbeiterinnen und Mitarbeiter auch von zu Hause aus produktiv sein können. Das spart Bürokosten und gibt den Beschäftigten mehr Flexibilität bei der Gestaltung ihres Tages.", Translation: ""},
				{Character: "Professor Finch", Text: "Allerdings warnen Fachleute vor den Schattenseiten. Wer ständig zu Hause arbeitet, verliert leicht den Kontakt zu Kolleginnen und Kollegen. Die Grenze zwischen Beruf und Privatleben verschwimmt, und manche Menschen fühlen sich isoliert oder arbeiten sogar mehr als zuvor.", Translation: ""},
				{Character: "Professor Finch", Text: "Die meisten Experten plädieren deshalb für ein hybrides Modell. An einigen Tagen kommen die Mitarbeiter ins Büro, um sich auszutauschen und gemeinsam an Projekten zu arbeiten, an anderen Tagen bleiben sie zu Hause, um konzentriert und ungestört Aufgaben zu erledigen.", Translation: ""},
			},
			Questions: []PaperQuestion{
				{Question: "Was haben viele Unternehmen seit der Pandemie festgestellt?", Options: []string{"Mitarbeiter können auch zu Hause produktiv sein", "Homeoffice ist unmöglich", "Büros sind billiger als früher", "Niemand will mehr arbeiten"}, CorrectAnswer: "Mitarbeiter können auch zu Hause produktiv sein"},
				{Question: "Welcher Vorteil des Homeoffice wird genannt?", Options: []string{"Mehr Flexibilität", "Höhere Gehälter", "Längere Ferien", "Weniger Aufgaben"}, CorrectAnswer: "Mehr Flexibilität"},
				{Question: "Vor welcher Schattenseite warnen die Fachleute?", Options: []string{"Verlust des Kontakts zu Kollegen", "Zu viele Meetings", "Zu hohe Bürokosten", "Schlechtes Internet"}, CorrectAnswer: "Verlust des Kontakts zu Kollegen"},
				{Question: "Was passiert mit der Grenze zwischen Beruf und Privatleben?", Options: []string{"Sie verschwimmt", "Sie wird klarer", "Sie verschwindet ganz", "Sie wird gesetzlich geregelt"}, CorrectAnswer: "Sie verschwimmt"},
				{Question: "Für welches Modell plädieren die meisten Experten?", Options: []string{"Ein hybrides Modell", "Nur Büroarbeit", "Nur Homeoffice", "Eine Vier-Tage-Woche"}, CorrectAnswer: "Ein hybrides Modell"},
			},
		},
		Reading: PaperReading{
			Title: "Soziale Medien und Jugendliche",
			Paragraphs: []string{
				"Kaum eine Erfindung hat den Alltag junger Menschen so stark geprägt wie soziale Medien. Plattformen wie Instagram oder TikTok ermöglichen es, in Sekundenschnelle mit Freunden in Kontakt zu bleiben, kreative Inhalte zu teilen und sich über aktuelle Themen zu informieren.",
				"Kritiker weisen jedoch darauf hin, dass der ständige Vergleich mit anderen das Selbstwertgefühl belasten kann. Wer täglich perfekt inszenierte Bilder sieht, vergisst leicht, dass diese sorgfältig ausgewählt und bearbeitet wurden. Studien deuten darauf hin, dass eine intensive Nutzung mit Schlafproblemen und Unzufriedenheit zusammenhängen kann.",
				"Fachleute fordern daher nicht ein Verbot, sondern einen bewussteren Umgang. Schulen sollten Medienkompetenz vermitteln, damit Jugendliche Inhalte kritisch hinterfragen und ihre Bildschirmzeit selbst sinnvoll begrenzen können. Letztlich liegt die Verantwortung bei den Plattformen, den Eltern und den Nutzern gemeinsam.",
			},
			Questions: []PaperQuestion{
				{Question: "Welchen Vorteil sozialer Medien nennt der Text?", Options: []string{"Schnell mit Freunden in Kontakt bleiben", "Bessere Schulnoten", "Mehr Schlaf", "Weniger Werbung"}, CorrectAnswer: "Schnell mit Freunden in Kontakt bleiben"},
				{Question: "Was kritisieren die Kritiker?", Options: []string{"Den ständigen Vergleich mit anderen", "Die hohen Kosten", "Die langsame Technik", "Den Mangel an Inhalten"}, CorrectAnswer: "Den ständigen Vergleich mit anderen"},
				{Question: "Was vergessen viele Nutzer laut Text?", Options: []string{"Dass die Bilder bearbeitet sind", "Dass die Plattformen kostenlos sind", "Dass sie Freunde haben", "Dass sie lernen müssen"}, CorrectAnswer: "Dass die Bilder bearbeitet sind"},
				{Question: "Womit hängt eine intensive Nutzung laut Studien zusammen?", Options: []string{"Mit Schlafproblemen", "Mit besseren Noten", "Mit mehr Sport", "Mit mehr Freunden"}, CorrectAnswer: "Mit Schlafproblemen"},
				{Question: "Was fordern die Fachleute?", Options: []string{"Einen bewussteren Umgang", "Ein komplettes Verbot", "Mehr Werbung", "Höhere Preise"}, CorrectAnswer: "Einen bewussteren Umgang"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "a formal letter of complaint to a hotel about several problems during your stay, describing what happened and asking for an appropriate solution.",
			MinWords: 130,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Obwohl soziale Medien viele Vorteile bieten, sollten wir die Risiken für die Privatsphäre nicht unterschätzen.",
			Speaker:     "Lumora",
			Translation: "Although social media offer many advantages, we should not underestimate the risks to privacy.",
		},
	},

	// ---------------- C1 ----------------
	"C1": {
		Listening: PaperListening{
			Title: "Vortrag: Urbanisierung und Nachhaltigkeit",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "Bis zur Mitte dieses Jahrhunderts werden voraussichtlich zwei Drittel der Weltbevölkerung in Städten leben. Diese rasante Urbanisierung bringt enorme Herausforderungen mit sich: Der Energieverbrauch steigt, der Verkehr nimmt zu, und der bezahlbare Wohnraum wird knapp. Wer glaubt, man könne diese Probleme allein durch technische Innovationen lösen, greift jedoch zu kurz.", Translation: ""},
				{Character: "Professor Finch", Text: "Entscheidend ist vielmehr ein integriertes Konzept, das ökologische, soziale und ökonomische Aspekte miteinander verbindet. Begrünte Dächer und Solaranlagen sind sinnvoll, aber sie entfalten ihre Wirkung erst, wenn sie in eine kluge Stadtplanung eingebettet sind, die kurze Wege ermöglicht und verschiedene Funktionen – Wohnen, Arbeiten, Erholung – miteinander verzahnt.", Translation: ""},
				{Character: "Professor Finch", Text: "Nicht zuletzt hängt der Erfolg davon ab, ob es gelingt, die Bürgerinnen und Bürger einzubeziehen. Wo Menschen das Gefühl haben, mitgestalten zu können, identifizieren sie sich mit ihrer Umgebung und gehen sorgsamer mit gemeinsamen Ressourcen um. Nachhaltigkeit ist somit nicht nur eine technische, sondern vor allem eine gesellschaftliche Aufgabe.", Translation: ""},
			},
			Questions: []PaperQuestion{
				{Question: "Wie viele Menschen werden bis Mitte des Jahrhunderts in Städten leben?", Options: []string{"Etwa zwei Drittel der Weltbevölkerung", "Die Hälfte aller Menschen", "Ein Viertel der Weltbevölkerung", "Fast alle Menschen"}, CorrectAnswer: "Etwa zwei Drittel der Weltbevölkerung"},
				{Question: "Welche Position vertritt der Redner zu technischen Innovationen?", Options: []string{"Sie allein reichen nicht aus", "Sie lösen alle Probleme", "Sie sind überflüssig", "Sie sind zu teuer"}, CorrectAnswer: "Sie allein reichen nicht aus"},
				{Question: "Was versteht der Redner unter einem 'integrierten Konzept'?", Options: []string{"Die Verbindung ökologischer, sozialer und ökonomischer Aspekte", "Den Bau möglichst hoher Gebäude", "Den Verzicht auf jede Technik", "Die Trennung von Wohnen und Arbeiten"}, CorrectAnswer: "Die Verbindung ökologischer, sozialer und ökonomischer Aspekte"},
				{Question: "Wann entfalten begrünte Dächer laut Redner ihre Wirkung?", Options: []string{"Wenn sie in eine kluge Stadtplanung eingebettet sind", "Wenn sie besonders groß sind", "Wenn der Staat sie bezahlt", "Wenn es viel regnet"}, CorrectAnswer: "Wenn sie in eine kluge Stadtplanung eingebettet sind"},
				{Question: "Warum ist die Einbeziehung der Bürger wichtig?", Options: []string{"Sie gehen dann sorgsamer mit Ressourcen um", "Sie zahlen dann höhere Steuern", "Sie ziehen dann aufs Land", "Sie bauen dann selbst Häuser"}, CorrectAnswer: "Sie gehen dann sorgsamer mit Ressourcen um"},
				{Question: "Wie beschreibt der Redner Nachhaltigkeit abschließend?", Options: []string{"Als vor allem gesellschaftliche Aufgabe", "Als rein technisches Problem", "Als unlösbare Aufgabe", "Als Sache der Politik allein"}, CorrectAnswer: "Als vor allem gesellschaftliche Aufgabe"},
			},
		},
		Reading: PaperReading{
			Title: "Die Kunst des Scheiterns",
			Paragraphs: []string{
				"In einer Leistungsgesellschaft, die Erfolg verehrt, gilt das Scheitern als Makel, den man möglichst verbirgt. Dabei zeigt ein Blick in die Geschichte der Wissenschaft und der Wirtschaft, dass Fehlschläge oft die Voraussetzung für Durchbrüche sind. Wer nie scheitert, hat vermutlich nie etwas wirklich Neues gewagt.",
				"In den letzten Jahren hat sich daher in manchen Unternehmen eine sogenannte Fehlerkultur etabliert. Statt Schuldige zu suchen, fragt man, was aus einem Misserfolg zu lernen ist. Diese Haltung verlangt Mut, denn sie setzt voraus, dass Vorgesetzte Fehler nicht bestrafen, sondern als Quelle wertvoller Erkenntnisse betrachten.",
				"Kritiker wenden ein, dass eine allzu nachsichtige Haltung zur Gleichgültigkeit führen könne. Entscheidend ist deshalb die Unterscheidung zwischen vermeidbaren Nachlässigkeiten und produktiven Fehlern, die im Streben nach Innovation unvermeidlich entstehen. Nur wer diese Differenzierung beherrscht, kann aus dem Scheitern tatsächlich Kapital schlagen.",
			},
			Questions: []PaperQuestion{
				{Question: "Wie wird Scheitern in einer Leistungsgesellschaft meist betrachtet?", Options: []string{"Als Makel, den man verbirgt", "Als Zeichen von Mut", "Als völlig normal", "Als gesetzlich verboten"}, CorrectAnswer: "Als Makel, den man verbirgt"},
				{Question: "Was zeigt der Blick in die Geschichte laut Text?", Options: []string{"Fehlschläge sind oft Voraussetzung für Durchbrüche", "Erfolg kommt ohne Anstrengung", "Niemand scheitert wirklich", "Wissenschaft braucht keine Fehler"}, CorrectAnswer: "Fehlschläge sind oft Voraussetzung für Durchbrüche"},
				{Question: "Was kennzeichnet eine 'Fehlerkultur'?", Options: []string{"Man fragt, was man aus dem Misserfolg lernt", "Man bestraft die Schuldigen", "Man verschweigt Fehler", "Man vermeidet jedes Risiko"}, CorrectAnswer: "Man fragt, was man aus dem Misserfolg lernt"},
				{Question: "Was setzt diese Haltung laut Text voraus?", Options: []string{"Dass Vorgesetzte Fehler nicht bestrafen", "Dass Mitarbeiter nie Fehler machen", "Dass der Staat eingreift", "Dass alle gleich viel verdienen"}, CorrectAnswer: "Dass Vorgesetzte Fehler nicht bestrafen"},
				{Question: "Welchen Einwand bringen Kritiker vor?", Options: []string{"Zu viel Nachsicht könne zu Gleichgültigkeit führen", "Fehler seien immer gut", "Erfolg sei unwichtig", "Innovation sei überschätzt"}, CorrectAnswer: "Zu viel Nachsicht könne zu Gleichgültigkeit führen"},
				{Question: "Worauf kommt es laut Text entscheidend an?", Options: []string{"Die Unterscheidung zwischen Nachlässigkeit und produktiven Fehlern", "Die Vermeidung jeglicher Fehler", "Die Bestrafung von Misserfolg", "Den Verzicht auf Innovation"}, CorrectAnswer: "Die Unterscheidung zwischen Nachlässigkeit und produktiven Fehlern"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "an argumentative essay weighing the advantages and disadvantages of remote work, taking a clear position and supporting it with examples.",
			MinWords: 180,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Die zunehmende Urbanisierung stellt unsere Städte vor enorme Herausforderungen, die sich nur durch nachhaltige Konzepte bewältigen lassen.",
			Speaker:     "Lumora",
			Translation: "Increasing urbanisation poses enormous challenges for our cities that can only be met through sustainable concepts.",
		},
	},

	// ---------------- C2 ----------------
	"C2": {
		Listening: PaperListening{
			Title: "Essay-Vortrag: Sprache und Denken",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "Die Frage, inwiefern die Sprache, die wir sprechen, unser Denken formt, beschäftigt Philosophen und Linguisten seit Jahrhunderten. Die sogenannte Sapir-Whorf-Hypothese behauptet in ihrer starken Form, dass die Struktur einer Sprache die Wahrnehmung der Wirklichkeit determiniert. In dieser radikalen Fassung gilt sie heute allerdings als widerlegt.", Translation: ""},
				{Character: "Professor Finch", Text: "Differenzierter ist die schwache Variante, die lediglich von einer Beeinflussung ausgeht. Tatsächlich legen empirische Studien nahe, dass grammatische Kategorien, etwa die Art, wie eine Sprache Zeit oder Raum gliedert, subtile Auswirkungen auf Gedächtnis und Aufmerksamkeit haben können. Von einer Determination kann jedoch keine Rede sein.", Translation: ""},
				{Character: "Professor Finch", Text: "Bemerkenswert ist, dass diese Debatte längst nicht nur akademisch ist. Im Zeitalter automatischer Übersetzung und künstlicher Intelligenz stellt sich erneut die Frage, ob Bedeutung restlos von einer Sprache in eine andere übertragbar ist – oder ob mit jeder Übersetzung unweigerlich Nuancen verloren gehen, die sich der vollständigen Rekonstruktion entziehen.", Translation: ""},
			},
			Questions: []PaperQuestion{
				{Question: "Was behauptet die starke Form der Sapir-Whorf-Hypothese?", Options: []string{"Die Sprachstruktur determiniert die Wahrnehmung der Wirklichkeit", "Sprache hat keinen Einfluss auf das Denken", "Alle Sprachen sind identisch", "Denken entsteht ohne Sprache"}, CorrectAnswer: "Die Sprachstruktur determiniert die Wahrnehmung der Wirklichkeit"},
				{Question: "Wie gilt diese starke Form heute?", Options: []string{"Als widerlegt", "Als bewiesen", "Als unerforscht", "Als verboten"}, CorrectAnswer: "Als widerlegt"},
				{Question: "Was besagt die schwache Variante?", Options: []string{"Sprache beeinflusst das Denken, determiniert es aber nicht", "Sprache hat keinerlei Wirkung", "Denken bestimmt die Sprache vollständig", "Grammatik ist bedeutungslos"}, CorrectAnswer: "Sprache beeinflusst das Denken, determiniert es aber nicht"},
				{Question: "Worauf können grammatische Kategorien laut Studien wirken?", Options: []string{"Auf Gedächtnis und Aufmerksamkeit", "Auf die Körpergröße", "Auf das Wetter", "Auf den Wortschatz allein"}, CorrectAnswer: "Auf Gedächtnis und Aufmerksamkeit"},
				{Question: "Warum ist die Debatte laut Redner aktuell?", Options: []string{"Wegen automatischer Übersetzung und künstlicher Intelligenz", "Wegen steigender Mieten", "Wegen des Klimawandels", "Wegen neuer Schulgesetze"}, CorrectAnswer: "Wegen automatischer Übersetzung und künstlicher Intelligenz"},
				{Question: "Welche Frage stellt sich beim Übersetzen?", Options: []string{"Ob Bedeutung restlos übertragbar ist", "Ob Maschinen schneller sind als Menschen", "Ob Sprachen aussterben", "Ob Grammatik nötig ist"}, CorrectAnswer: "Ob Bedeutung restlos übertragbar ist"},
			},
		},
		Reading: PaperReading{
			Title: "Über den Wert der Muße",
			Paragraphs: []string{
				"In einer Kultur, die Geschäftigkeit mit Bedeutung verwechselt, ist die Muße in Verruf geraten. Wer innehält, gilt schnell als unproduktiv, ja beinahe als verdächtig. Dabei war die Muße in der antiken Philosophie kein Gegenteil der Arbeit, sondern deren eigentliches Ziel: ein Zustand freier, zweckloser Betrachtung, in dem der Mensch zu sich selbst findet.",
				"Die moderne Ökonomie hat dieses Verständnis auf den Kopf gestellt. Freizeit wird heute weniger als Raum der Besinnung denn als Konsumgelegenheit begriffen, die ihrerseits durchgeplant und optimiert sein will. Selbst die Erholung gerät so unter den Imperativ der Effizienz, und das, was eigentlich entlasten sollte, wird zur neuen Belastung.",
				"Ob sich diese Entwicklung umkehren lässt, ist fraglich. Gleichwohl mehren sich die Stimmen, die in der bewussten Verlangsamung keine nostalgische Flucht, sondern eine notwendige Korrektur sehen. Vielleicht besteht die eigentliche Kunst des Lebens nicht darin, möglichst viel zu erledigen, sondern darin, dem scheinbar Nutzlosen wieder einen Eigenwert zuzugestehen.",
			},
			Questions: []PaperQuestion{
				{Question: "Womit wird Geschäftigkeit laut Text verwechselt?", Options: []string{"Mit Bedeutung", "Mit Reichtum", "Mit Glück", "Mit Gesundheit"}, CorrectAnswer: "Mit Bedeutung"},
				{Question: "Wie verstand die antike Philosophie die Muße?", Options: []string{"Als eigentliches Ziel der Arbeit", "Als reine Faulheit", "Als Strafe", "Als Zeitverschwendung"}, CorrectAnswer: "Als eigentliches Ziel der Arbeit"},
				{Question: "Wie wird Freizeit laut Text heute meist begriffen?", Options: []string{"Als Konsumgelegenheit", "Als Raum der Besinnung", "Als verbotene Zeit", "Als reine Arbeit"}, CorrectAnswer: "Als Konsumgelegenheit"},
				{Question: "Was geschieht laut Text mit der Erholung?", Options: []string{"Sie gerät unter den Imperativ der Effizienz", "Sie wird abgeschafft", "Sie wird kostenlos", "Sie verliert jede Bedeutung"}, CorrectAnswer: "Sie gerät unter den Imperativ der Effizienz"},
				{Question: "Wie sehen manche Stimmen die bewusste Verlangsamung?", Options: []string{"Als notwendige Korrektur", "Als nostalgische Flucht", "Als wirtschaftlichen Schaden", "Als vorübergehende Mode"}, CorrectAnswer: "Als notwendige Korrektur"},
				{Question: "Worin könnte laut Schluss die 'Kunst des Lebens' bestehen?", Options: []string{"Dem scheinbar Nutzlosen einen Eigenwert zuzugestehen", "Möglichst viel zu erledigen", "Stets effizient zu sein", "Die Freizeit zu optimieren"}, CorrectAnswer: "Dem scheinbar Nutzlosen einen Eigenwert zuzugestehen"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "a nuanced argumentative essay on whether meaning can be fully translated between languages, engaging with counter-arguments and reaching a reasoned conclusion.",
			MinWords: 230,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Die Frage, inwieweit unsere Sprache das Denken prägt, beschäftigt die Philosophie seit Jahrhunderten und lässt sich keineswegs eindeutig beantworten.",
			Speaker:     "Lumora",
			Translation: "The question of how far our language shapes thought has occupied philosophy for centuries and cannot be answered unambiguously.",
		},
	},
}
