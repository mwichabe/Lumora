package controllers

// spanishPapers is the hand-authored Spanish proficiency-exam bank (A1–C2),
// mirroring the German bank. Each level scales up: longer listening passages,
// denser reading texts, more and harder questions (Spanish-language questions
// from B2 upward), longer required writing, and more complex speaking prompts.
var spanishPapers = map[string]paperContent{

	// ---------------- FINAL (comprehensive A1→C2 mastery) ----------------
	"FINAL": {
		Listening: PaperListening{
			Title: "Conferencia final: el valor del multilingüismo",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "Bienvenidos a nuestra conferencia final. Hoy quisiera reflexionar sobre el valor del multilingüismo. En un mundo cada vez más interconectado, dominar varias lenguas ha dejado de ser un lujo para convertirse en una competencia clave.", Translation: ""},
				{Character: "Professor Finch", Text: "Diversos estudios demuestran que las personas multilingües no solo tienen ventajas profesionales, sino que también piensan con mayor flexibilidad. Quien alterna entre idiomas entrena su mente para cambiar de perspectiva y resolver problemas de forma más creativa.", Translation: ""},
				{Character: "Professor Finch", Text: "No obstante, sería ingenuo afirmar que aprender una lengua es tarea fácil. Exige paciencia, disciplina y la disposición a equivocarse. Y precisamente ahí reside la verdadera recompensa: quien aprende un idioma aprende también a ser humilde y perseverante.", Translation: ""},
				{Character: "Professor Finch", Text: "En definitiva, cada nueva lengua no solo abre puertas, sino también horizontes. Es, por así decirlo, el límite de nuestro mundo; y quien lo amplía, se amplía a sí mismo.", Translation: ""},
			},
			Questions: []PaperQuestion{
				{Question: "¿Cómo describe el ponente el multilingüismo hoy?", Options: []string{"Como una competencia clave", "Como un lujo", "Como una moda", "Como una obligación"}, CorrectAnswer: "Como una competencia clave"},
				{Question: "¿Qué ventaja cognitiva menciona?", Options: []string{"Pensar con mayor flexibilidad", "Solo una mejor memoria", "Calcular más rápido", "Dormir menos"}, CorrectAnswer: "Pensar con mayor flexibilidad"},
				{Question: "¿Qué entrena quien alterna entre idiomas?", Options: []string{"Cambiar de perspectiva", "Hablar más alto", "Memorizar listas", "Escribir más rápido"}, CorrectAnswer: "Cambiar de perspectiva"},
				{Question: "¿Qué exige aprender una lengua, según el ponente?", Options: []string{"Paciencia y disciplina", "Solo talento", "Mucho dinero", "Una pronunciación perfecta"}, CorrectAnswer: "Paciencia y disciplina"},
				{Question: "¿Dónde reside la verdadera recompensa?", Options: []string{"En volverse humilde y perseverante", "En hacerse rico", "En no equivocarse nunca", "En ser famoso"}, CorrectAnswer: "En volverse humilde y perseverante"},
				{Question: "¿Cómo describe una nueva lengua al final?", Options: []string{"Como el límite y la ampliación de nuestro mundo", "Como una simple herramienta", "Como una pérdida de tiempo", "Como una moda pasajera"}, CorrectAnswer: "Como el límite y la ampliación de nuestro mundo"},
			},
		},
		Reading: PaperReading{
			Title: "El aprendizaje permanente",
			Paragraphs: []string{
				"La idea de que la educación termina con el fin de la escuela pertenece hace tiempo al pasado. En un mundo que cambia a una velocidad vertiginosa, el aprendizaje permanente se ha convertido en una necesidad.",
				"Quien hoy deja de aprender corre el riesgo de quedarse atrás. Surgen nuevas tecnologías, desaparecen profesiones enteras y el conocimiento que ayer era actual puede quedar obsoleto mañana.",
				"Los críticos objetan que la presión constante por formarse genera estrés y agotamiento. Este reparo no es desdeñable. Sin embargo, el aprendizaje permanente no consiste en una optimización sin tregua, sino en la curiosidad.",
				"Porque quien conserva la curiosidad no envejece de espíritu. La educación, bien entendida, no es una carrera, sino una actitud: la disposición a dejarse asombrar una y otra vez.",
			},
			Questions: []PaperQuestion{
				{Question: "¿Qué pertenece al pasado, según el texto?", Options: []string{"Que la educación termine con la escuela", "Que existan las escuelas", "Que la gente aprenda", "Que cambie la tecnología"}, CorrectAnswer: "Que la educación termine con la escuela"},
				{Question: "¿Qué arriesga quien deja de aprender?", Options: []string{"Quedarse atrás", "Hacerse rico", "Tener más tiempo libre", "Mejorar su salud"}, CorrectAnswer: "Quedarse atrás"},
				{Question: "¿Qué puede ocurrir con el conocimiento de ayer?", Options: []string{"Quedar obsoleto mañana", "Volverse eterno", "Ganar valor siempre", "No cambiar nunca"}, CorrectAnswer: "Quedar obsoleto mañana"},
				{Question: "¿Qué objetan los críticos?", Options: []string{"Que la presión constante genera estrés", "Que aprender es demasiado barato", "Que hay pocas escuelas", "Que nadie quiere aprender"}, CorrectAnswer: "Que la presión constante genera estrés"},
				{Question: "¿En qué consiste el aprendizaje permanente, según el autor?", Options: []string{"En la curiosidad", "En optimizarse sin tregua", "En aprobar exámenes", "En competir"}, CorrectAnswer: "En la curiosidad"},
				{Question: "¿Cómo entiende el autor la educación al final?", Options: []string{"Como una actitud", "Como una carrera", "Como una obligación", "Como un producto"}, CorrectAnswer: "Como una actitud"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "a comprehensive, well-structured essay in Spanish (300+ words): 'Is it worth learning several languages?' Present clear arguments and counter-arguments, use varied grammar and connectors, and reach a reasoned conclusion.",
			MinWords: 300,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Cada nueva lengua no solo abre puertas, sino también horizontes; y quien amplía sus límites se amplía, al mismo tiempo, a sí mismo.",
			Speaker:     "Lumora",
			Translation: "Every new language opens not only doors but horizons; and whoever widens its limits widens, at the same time, themselves.",
		},
	},

	// ---------------- A1 ----------------
	"A1": {
		Listening: PaperListening{
			Title: "El día de Marta",
			Lines: []PaperLine{
				{Character: "Cora", Text: "¡Hola! Me llamo Marta y vivo en Madrid. Tengo veinticinco años y soy enfermera. Por la mañana me levanto a las siete, desayuno café con tostadas y voy al hospital en metro.", Translation: "Hi! My name is Marta and I live in Madrid. I'm twenty-five and I'm a nurse. In the morning I get up at seven, have coffee with toast and go to the hospital by metro."},
				{Character: "Cora", Text: "El viaje dura veinte minutos. Por la tarde estudio inglés y por la noche ceno con mi familia. Los fines de semana me gusta pasear por el parque.", Translation: "The journey takes twenty minutes. In the afternoon I study English and in the evening I have dinner with my family. At weekends I like to walk in the park."},
			},
			Questions: []PaperQuestion{
				{Question: "Where does Marta live?", Options: []string{"Madrid", "Barcelona", "Sevilla", "Valencia"}, CorrectAnswer: "Madrid"},
				{Question: "What is her job?", Options: []string{"Nurse", "Teacher", "Doctor", "Waiter"}, CorrectAnswer: "Nurse"},
				{Question: "How does she go to the hospital?", Options: []string{"By metro", "By car", "On foot", "By bike"}, CorrectAnswer: "By metro"},
				{Question: "What does she like to do at weekends?", Options: []string{"Walk in the park", "Cook", "Study", "Sleep"}, CorrectAnswer: "Walk in the park"},
			},
		},
		Reading: PaperReading{
			Title: "Una postal desde Valencia",
			Paragraphs: []string{
				"¡Hola, Ana! Estoy de vacaciones en Valencia. La ciudad es muy bonita y hace mucho calor.",
				"Por la mañana voy a la playa y por la tarde como paella en un restaurante. Mañana quiero visitar el centro. Un abrazo, Pablo.",
			},
			Questions: []PaperQuestion{
				{Question: "Where is Pablo?", Options: []string{"In Valencia", "In Madrid", "In Sevilla", "In Bilbao"}, CorrectAnswer: "In Valencia"},
				{Question: "What is the weather like?", Options: []string{"Very hot", "Cold", "Rainy", "Windy"}, CorrectAnswer: "Very hot"},
				{Question: "What does Pablo eat?", Options: []string{"Paella", "Pizza", "Soup", "Fish"}, CorrectAnswer: "Paella"},
				{Question: "What does he want to do tomorrow?", Options: []string{"Visit the centre", "Go home", "Swim", "Cook"}, CorrectAnswer: "Visit the centre"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "a short email to a friend: introduce yourself, say where you live, and describe what you do every day.",
			MinWords: 30,
		},
		Speaking: PaperSpeaking{
			Phrase:      "¡Hola! Me llamo Ana, soy de España y vivo en Madrid.",
			Speaker:     "Lumora",
			Translation: "Hi! My name is Ana, I'm from Spain and I live in Madrid.",
		},
	},

	// ---------------- A2 ----------------
	"A2": {
		Listening: PaperListening{
			Title: "El fin de semana de Luis",
			Lines: []PaperLine{
				{Character: "Riko", Text: "El fin de semana pasado fui a la montaña con mis amigos. El sábado salimos temprano y llegamos a mediodía. Comimos en un pueblo pequeño y por la tarde dimos un paseo largo.", Translation: "Last weekend I went to the mountains with my friends. On Saturday we left early and arrived at midday. We ate in a small village and in the afternoon took a long walk."},
				{Character: "Riko", Text: "El domingo llovió un poco, así que visitamos un museo. Volvimos a casa muy cansados pero contentos. Fue un fin de semana estupendo.", Translation: "On Sunday it rained a little, so we visited a museum. We went home very tired but happy. It was a great weekend."},
			},
			Questions: []PaperQuestion{
				{Question: "Where did Luis go?", Options: []string{"To the mountains", "To the beach", "To the city", "To work"}, CorrectAnswer: "To the mountains"},
				{Question: "When did they arrive on Saturday?", Options: []string{"At midday", "Early morning", "At night", "In the evening"}, CorrectAnswer: "At midday"},
				{Question: "What did they do on Sunday (it rained)?", Options: []string{"Visited a museum", "Went swimming", "Stayed in bed", "Climbed a mountain"}, CorrectAnswer: "Visited a museum"},
				{Question: "How was the weekend?", Options: []string{"Great", "Boring", "Terrible", "Stressful"}, CorrectAnswer: "Great"},
			},
		},
		Reading: PaperReading{
			Title: "Una reseña de restaurante",
			Paragraphs: []string{
				"Ayer cené en el restaurante La Cocina. Pedí pescado con verduras y de postre un flan. La comida estaba muy rica y el camarero fue muy amable.",
				"Sin embargo, el restaurante estaba un poco lleno y tuvimos que esperar veinte minutos. En general, lo recomiendo y volveré pronto.",
			},
			Questions: []PaperQuestion{
				{Question: "When did the writer go?", Options: []string{"Yesterday", "Last week", "Today", "Last month"}, CorrectAnswer: "Yesterday"},
				{Question: "What did they order?", Options: []string{"Fish with vegetables", "Pizza", "Steak", "Pasta"}, CorrectAnswer: "Fish with vegetables"},
				{Question: "What was the problem?", Options: []string{"It was full and they waited", "The food was cold", "It was expensive", "The waiter was rude"}, CorrectAnswer: "It was full and they waited"},
				{Question: "Will they return?", Options: []string{"Yes, soon", "Never", "Not sure", "Only once"}, CorrectAnswer: "Yes, soon"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "an email to a friend about your last weekend: what you did and your plans for the next one.",
			MinWords: 55,
		},
		Speaking: PaperSpeaking{
			Phrase:      "El fin de semana pasado fui a la playa con mi familia y comimos paella.",
			Speaker:     "Lumora",
			Translation: "Last weekend I went to the beach with my family and we ate paella.",
		},
	},

	// ---------------- B1 ----------------
	"B1": {
		Listening: PaperListening{
			Title: "Una entrevista de trabajo",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "Buenos días y bienvenida. Cuénteme, ¿qué experiencia tiene?", Translation: ""},
				{Character: "Mira", Text: "Antes trabajé tres años en una agencia de viajes, donde organizaba reservas y atendía a los clientes.", Translation: ""},
				{Character: "Professor Finch", Text: "Muy bien. ¿Y por qué le interesa este puesto?", Translation: ""},
				{Character: "Mira", Text: "Me gustaría asumir nuevos retos y aprender en una empresa internacional. Además, hablo tres idiomas.", Translation: ""},
				{Character: "Professor Finch", Text: "Perfecto. Le informaremos de nuestra decisión la próxima semana.", Translation: ""},
			},
			Questions: []PaperQuestion{
				{Question: "Where did Mira work before?", Options: []string{"A travel agency", "A hospital", "A school", "A bank"}, CorrectAnswer: "A travel agency"},
				{Question: "For how long?", Options: []string{"Three years", "One year", "Five years", "Six months"}, CorrectAnswer: "Three years"},
				{Question: "Why does she want this job?", Options: []string{"New challenges and an international company", "More money", "It's near home", "Shorter hours"}, CorrectAnswer: "New challenges and an international company"},
				{Question: "How many languages does she speak?", Options: []string{"Three", "Two", "One", "Four"}, CorrectAnswer: "Three"},
				{Question: "When will they inform her?", Options: []string{"Next week", "Tomorrow", "Next month", "The same day"}, CorrectAnswer: "Next week"},
			},
		},
		Reading: PaperReading{
			Title: "El teletrabajo",
			Paragraphs: []string{
				"En los últimos años, el teletrabajo se ha vuelto muy común. Por un lado, permite ahorrar tiempo de transporte y ofrece más flexibilidad; muchos trabajadores afirman que son más productivos en casa.",
				"Sin embargo, también tiene desventajas: algunas personas se sienten aisladas y les cuesta separar la vida laboral de la personal.",
				"En mi opinión, lo ideal sería un modelo mixto que combine el trabajo en casa y en la oficina, según las necesidades de cada persona.",
			},
			Questions: []PaperQuestion{
				{Question: "What has become common?", Options: []string{"Remote work", "Night shifts", "Long holidays", "Office parties"}, CorrectAnswer: "Remote work"},
				{Question: "Which is an advantage mentioned?", Options: []string{"More flexibility", "Higher salary", "Free meals", "More meetings"}, CorrectAnswer: "More flexibility"},
				{Question: "Which is a disadvantage?", Options: []string{"Feeling isolated", "Too much travel", "Noisy offices", "Less pay"}, CorrectAnswer: "Feeling isolated"},
				{Question: "What does the writer prefer?", Options: []string{"A mixed model", "Only office work", "Only remote work", "No work"}, CorrectAnswer: "A mixed model"},
				{Question: "What is the tone of the text?", Options: []string{"Balanced", "Angry", "Sad", "Sarcastic"}, CorrectAnswer: "Balanced"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "a forum post giving your opinion on whether technology makes life better, with reasons and an example.",
			MinWords: 90,
		},
		Speaking: PaperSpeaking{
			Phrase:      "En mi opinión, viajar es la mejor forma de aprender sobre otras culturas.",
			Speaker:     "Lumora",
			Translation: "In my opinion, travelling is the best way to learn about other cultures.",
		},
	},

	// ---------------- B2 ----------------
	"B2": {
		Listening: PaperListening{
			Title: "Debate: el medio ambiente",
			Lines: []PaperLine{
				{Character: "Zephyr", Text: "Desde mi punto de vista, la contaminación es el mayor problema de nuestras ciudades.", Translation: ""},
				{Character: "Cora", Text: "Estoy de acuerdo, pero también deberíamos hablar del consumo excesivo de energía.", Translation: ""},
				{Character: "Zephyr", Text: "Cierto. Si usáramos más el transporte público, reduciríamos notablemente las emisiones.", Translation: ""},
				{Character: "Cora", Text: "Sin embargo, mucha gente no renunciará a su coche a menos que existan alternativas eficaces.", Translation: ""},
				{Character: "Professor Finch", Text: "En definitiva, hace falta una combinación de medidas políticas y responsabilidad individual.", Translation: ""},
			},
			Questions: []PaperQuestion{
				{Question: "¿Cuál es, según Zephyr, el mayor problema de las ciudades?", Options: []string{"La contaminación", "El ruido", "El tráfico", "La basura"}, CorrectAnswer: "La contaminación"},
				{Question: "¿Qué tema añade Cora?", Options: []string{"El consumo excesivo de energía", "El precio de la vivienda", "La falta de parques", "El turismo"}, CorrectAnswer: "El consumo excesivo de energía"},
				{Question: "¿Qué reduciría las emisiones?", Options: []string{"Usar más el transporte público", "Construir más carreteras", "Bajar los impuestos", "Cerrar las fábricas"}, CorrectAnswer: "Usar más el transporte público"},
				{Question: "¿Bajo qué condición dejaría la gente el coche?", Options: []string{"Si existen alternativas eficaces", "Si sube la gasolina", "Si llueve", "Nunca"}, CorrectAnswer: "Si existen alternativas eficaces"},
				{Question: "¿Qué hace falta, en conclusión?", Options: []string{"Medidas políticas y responsabilidad individual", "Solo nuevas leyes", "Solo esfuerzo personal", "No hacer nada"}, CorrectAnswer: "Medidas políticas y responsabilidad individual"},
			},
		},
		Reading: PaperReading{
			Title: "Las redes sociales",
			Paragraphs: []string{
				"Pocas invenciones han transformado tanto nuestra forma de relacionarnos como las redes sociales. Gracias a ellas, podemos mantener el contacto con personas de todo el mundo y acceder a información al instante.",
				"No obstante, su uso excesivo plantea riesgos. La comparación constante con vidas aparentemente perfectas puede afectar a la autoestima, sobre todo entre los más jóvenes. Además, la difusión de noticias falsas se ha convertido en un problema serio.",
				"Por lo tanto, no se trata de prohibir estas herramientas, sino de educar en un uso crítico y responsable.",
			},
			Questions: []PaperQuestion{
				{Question: "¿Qué permiten las redes sociales, según el texto?", Options: []string{"Mantener el contacto y acceder a información", "Ganar dinero fácil", "Viajar gratis", "Aprender idiomas solos"}, CorrectAnswer: "Mantener el contacto y acceder a información"},
				{Question: "¿Qué puede afectar a la autoestima?", Options: []string{"La comparación constante", "El precio del móvil", "La falta de wifi", "Los anuncios"}, CorrectAnswer: "La comparación constante"},
				{Question: "¿A quiénes afecta sobre todo?", Options: []string{"A los más jóvenes", "A los mayores", "A los políticos", "A nadie"}, CorrectAnswer: "A los más jóvenes"},
				{Question: "¿Qué problema serio se menciona?", Options: []string{"La difusión de noticias falsas", "La lentitud de internet", "El exceso de fotos", "La publicidad"}, CorrectAnswer: "La difusión de noticias falsas"},
				{Question: "¿Cuál es la conclusión del autor?", Options: []string{"Educar en un uso crítico y responsable", "Prohibir las redes", "Usarlas sin límite", "Cerrar internet"}, CorrectAnswer: "Educar en un uso crítico y responsable"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "an argumentative essay weighing the advantages and disadvantages of remote work, taking a clear position supported with examples.",
			MinWords: 130,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Aunque las redes sociales tienen muchas ventajas, no deberíamos subestimar sus riesgos para la privacidad.",
			Speaker:     "Lumora",
			Translation: "Although social media have many advantages, we shouldn't underestimate their risks to privacy.",
		},
	},

	// ---------------- C1 ----------------
	"C1": {
		Listening: PaperListening{
			Title: "Tertulia: la inteligencia artificial",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "La inteligencia artificial plantea dilemas éticos que no podemos ignorar.", Translation: ""},
				{Character: "Zephyr", Text: "Si bien es cierto que entraña riesgos, sus beneficios en la medicina son innegables.", Translation: ""},
				{Character: "Professor Finch", Text: "Cabe matizar, no obstante, que todo depende de cómo se regule su uso.", Translation: ""},
				{Character: "Zephyr", Text: "Exacto. El problema no es la tecnología en sí, sino la falta de un marco legal común.", Translation: ""},
				{Character: "Professor Finch", Text: "En resumidas cuentas, necesitamos una regulación que proteja sin frenar la innovación.", Translation: ""},
			},
			Questions: []PaperQuestion{
				{Question: "¿Qué plantea la inteligencia artificial, según Finch?", Options: []string{"Dilemas éticos", "Problemas económicos", "Cuestiones de moda", "Nada relevante"}, CorrectAnswer: "Dilemas éticos"},
				{Question: "¿En qué campo destaca Zephyr sus beneficios?", Options: []string{"La medicina", "El deporte", "La cocina", "El turismo"}, CorrectAnswer: "La medicina"},
				{Question: "¿De qué depende todo, según Finch?", Options: []string{"De cómo se regule su uso", "De su precio", "De la moda", "Del idioma"}, CorrectAnswer: "De cómo se regule su uso"},
				{Question: "¿Cuál es el verdadero problema para Zephyr?", Options: []string{"La falta de un marco legal común", "La tecnología en sí", "El exceso de inversión", "La falta de interés"}, CorrectAnswer: "La falta de un marco legal común"},
				{Question: "¿Qué se necesita, en resumen?", Options: []string{"Regulación sin frenar la innovación", "Prohibir la IA", "Dejarla sin control", "Subir los impuestos"}, CorrectAnswer: "Regulación sin frenar la innovación"},
				{Question: "¿Qué tono tiene la conversación?", Options: []string{"Reflexivo y matizado", "Agresivo", "Indiferente", "Humorístico"}, CorrectAnswer: "Reflexivo y matizado"},
			},
		},
		Reading: PaperReading{
			Title: "Editorial: la era digital",
			Paragraphs: []string{
				"Vivimos inmersos en una vorágine tecnológica de la que resulta difícil sustraerse. Cada nueva aplicación promete hacernos la vida más fácil y, sin embargo, rara vez nos preguntamos qué perdemos a cambio.",
				"Si bien las redes nos acercan a quienes están lejos, han erosionado, paradójicamente, la conversación pausada y la atención sostenida.",
				"No se trata de demonizar el progreso —sería tan ingenuo como inútil—, sino de aprender a convivir con él de manera consciente.",
				"En definitiva, la tecnología debería estar al servicio del ser humano, y no al revés.",
			},
			Questions: []PaperQuestion{
				{Question: "¿Cómo describe el autor nuestra época?", Options: []string{"Una vorágine tecnológica", "Una edad dorada", "Una época tranquila", "Un desierto cultural"}, CorrectAnswer: "Una vorágine tecnológica"},
				{Question: "¿Qué rara vez nos preguntamos?", Options: []string{"Qué perdemos a cambio", "Cuánto cuesta", "Quién lo inventó", "Dónde se fabrica"}, CorrectAnswer: "Qué perdemos a cambio"},
				{Question: "¿Qué han erosionado las redes, paradójicamente?", Options: []string{"La conversación pausada", "La economía", "La memoria", "El idioma"}, CorrectAnswer: "La conversación pausada"},
				{Question: "¿Cuál es la postura del autor sobre el progreso?", Options: []string{"Convivir con él de forma consciente", "Rechazarlo del todo", "Adorarlo sin crítica", "Ignorarlo"}, CorrectAnswer: "Convivir con él de forma consciente"},
				{Question: "Según el autor, ¿al servicio de quién debe estar la tecnología?", Options: []string{"Del ser humano", "De las empresas", "De los gobiernos", "De sí misma"}, CorrectAnswer: "Del ser humano"},
				{Question: "¿Qué tono predomina en el editorial?", Options: []string{"Reflexivo y equilibrado", "Furioso", "Indiferente", "Cómico"}, CorrectAnswer: "Reflexivo y equilibrado"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "an argumentative essay on whether artificial intelligence should be regulated: develop arguments, counter-arguments and a reasoned conclusion.",
			MinWords: 180,
		},
		Speaking: PaperSpeaking{
			Phrase:      "La creciente digitalización plantea retos que solo podremos afrontar con una regulación inteligente y consensuada.",
			Speaker:     "Lumora",
			Translation: "Growing digitalisation poses challenges we can only meet with smart, agreed-upon regulation.",
		},
	},

	// ---------------- C2 ----------------
	"C2": {
		Listening: PaperListening{
			Title: "Debate: ¿libertad o seguridad?",
			Lines: []PaperLine{
				{Character: "Zephyr", Text: "No se trata, como algunos pretenden, de elegir entre libertad y seguridad; es un falso dilema.", Translation: ""},
				{Character: "Professor Finch", Text: "Ahora bien, en situaciones excepcionales, ¿no priorizaríamos la seguridad por encima de todo?", Translation: ""},
				{Character: "Zephyr", Text: "Hacerlo sería renunciar a aquello que decimos proteger. A fin de cuentas, sin libertad, la seguridad carece de sentido.", Translation: ""},
				{Character: "Professor Finch", Text: "Permítame discrepar, aunque reconozco la solidez de su argumento.", Translation: ""},
				{Character: "Zephyr", Text: "Quizá la verdadera cuestión no sea cuál elegir, sino cómo equilibrarlas sin sacrificar ninguna.", Translation: ""},
			},
			Questions: []PaperQuestion{
				{Question: "¿Cómo califica Zephyr la elección entre libertad y seguridad?", Options: []string{"Un falso dilema", "Una decisión sencilla", "Un asunto menor", "Una ley justa"}, CorrectAnswer: "Un falso dilema"},
				{Question: "¿Qué plantea Finch sobre las situaciones excepcionales?", Options: []string{"Priorizar la seguridad", "Eliminar la seguridad", "Ignorar el problema", "Pedir más libertad"}, CorrectAnswer: "Priorizar la seguridad"},
				{Question: "Según Zephyr, ¿qué carece de sentido sin libertad?", Options: []string{"La seguridad", "La economía", "La justicia", "La educación"}, CorrectAnswer: "La seguridad"},
				{Question: "¿Cómo reacciona Finch ante el argumento de Zephyr?", Options: []string{"Discrepa con respeto", "Lo acepta sin más", "Se enfada", "Lo ignora"}, CorrectAnswer: "Discrepa con respeto"},
				{Question: "¿Cuál es, para Zephyr, la verdadera cuestión?", Options: []string{"Cómo equilibrarlas sin sacrificar ninguna", "Cuál eliminar", "Quién tiene razón", "Cuándo decidir"}, CorrectAnswer: "Cómo equilibrarlas sin sacrificar ninguna"},
				{Question: "¿Qué registro emplean los hablantes?", Options: []string{"Formal y matizado", "Vulgar", "Infantil", "Técnico y frío"}, CorrectAnswer: "Formal y matizado"},
			},
		},
		Reading: PaperReading{
			Title: "Elogio de la duda",
			Paragraphs: []string{
				"Vivimos en una época que premia la certeza y desconfía de quien titubea. El que duda parece débil; el que afirma con rotundidad, en cambio, se gana el aplauso, aunque se equivoque.",
				"Y sin embargo, acaso sea la duda, y no la convicción ciega, el verdadero motor del pensamiento. Quien nunca duda no piensa: se limita a repetir certezas heredadas que jamás somete a examen.",
				"Dudar no es renunciar a toda creencia, sino sostenerla con la humildad de quien sabe que podría estar equivocado.",
				"Lejos de ser una debilidad, la duda es el primer acto de una mente verdaderamente libre.",
			},
			Questions: []PaperQuestion{
				{Question: "¿Qué premia la época actual, según el autor?", Options: []string{"La certeza", "La duda", "El silencio", "La riqueza"}, CorrectAnswer: "La certeza"},
				{Question: "¿Quién se gana el aplauso?", Options: []string{"El que afirma con rotundidad", "El que titubea", "El que calla", "El que pregunta"}, CorrectAnswer: "El que afirma con rotundidad"},
				{Question: "¿Qué es, para el autor, el verdadero motor del pensamiento?", Options: []string{"La duda", "La convicción ciega", "La memoria", "El miedo"}, CorrectAnswer: "La duda"},
				{Question: "¿Qué hace quien nunca duda?", Options: []string{"Repetir certezas heredadas", "Pensar con libertad", "Aprender deprisa", "Crear ideas nuevas"}, CorrectAnswer: "Repetir certezas heredadas"},
				{Question: "¿Qué significa dudar, según el texto?", Options: []string{"Sostener creencias con humildad", "Renunciar a todo", "No tener ideas", "Cambiar de opinión siempre"}, CorrectAnswer: "Sostener creencias con humildad"},
				{Question: "¿Cómo reformula el autor la duda al final?", Options: []string{"El primer acto de una mente libre", "Una debilidad", "Una pérdida de tiempo", "Un error"}, CorrectAnswer: "El primer acto de una mente libre"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "a nuanced critical essay on whether language shapes thought: build an argument, engage with objections, and reach a memorable conclusion.",
			MinWords: 230,
		},
		Speaking: PaperSpeaking{
			Phrase:      "La cuestión de hasta qué punto el lenguaje moldea nuestro pensamiento ha ocupado a la filosofía durante siglos y no admite respuestas simples.",
			Speaker:     "Lumora",
			Translation: "The question of how far language shapes our thought has occupied philosophy for centuries and admits no simple answers.",
		},
	},
}
