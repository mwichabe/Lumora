package database

import "gorm.io/gorm"

// Swahili course, B1 → C2. See seed_sw.go.

// ───────────────────── B1 — Mazungumzo: Conversation & Narrative ─────────────────────

func seedSwahiliB1(db *gorm.DB) {
	const sw = "sw"
	const u = "B1 · Mazungumzo — Conversation & Narrative"
	finch := "Professor Finch"

	s := addSkillL(db, sw, u, "Relative Clauses", "amba- and the relative infix: mtu aliyekuja, kitabu nilichonunua.", "Link", "#6C3FC5", 20, 300)
	l := addLesson(db, s, "The One Who…, The Thing That…", 1, 22,
		char(finch, "Two ways to say 'who/which'. 1) amba- + the class's relative particle: mtu ambaye (the person who), watu ambao, kitabu ambacho, vitabu ambavyo, gari ambalo, nyumba ambayo, mti ambao. 2) Infix the particle after the tense: mtu aliyekuja (the person who came), kitabu nilichonunua (the book I bought). The infix works with -li-, -na-, -ta- (as -taka-): atakayekuja (who will come)."),
		mc("the person who came", "mtu ___ alikuja", "ambaye", "ambaye", "ambacho", "ambalo", "ambayo"),
		mc("the book which I bought", "kitabu ___ nilinunua", "ambacho", "ambacho", "ambaye", "ambavyo", "ambalo"),
		mc("the car that broke down", "gari ___ liliharibika", "ambalo", "ambalo", "ambayo", "ambacho", "ambao"),
		mc("What does 'mtu aliyekuja' mean?", "mtu aliyekuja", "the person who came", "the person who came", "the person who will come", "the person is coming", "the person didn't come"),
		tr("Translate", "the food that I cooked", "chakula ambacho nilipika"),
		fill("Infix form: the children who are playing", "watoto wanaocheza → watoto wana___cheza", "o"),
		speak("Mira", "Huyu ndiye rafiki niliyekuambia habari zake."),
	)
	addVocab(db, l,
		vw("ambaye / ambao", "who (people, sg / pl)", "Mwalimu ambaye anafundisha hesabu.", "The teacher who teaches maths.", finch),
		vw("ambacho / ambavyo", "which (KI-/VI-)", "Kitabu ambacho nilisoma.", "The book that I read.", finch),
		vw("ambayo", "which (N- sg; MI- pl)", "Nyumba ambayo tunaishi.", "The house we live in.", "Nana"),
		vw("-ye-", "relative infix (person)", "Aliyeshinda ni nani?", "Who is the one who won?", "Blaze"),
		vw("ndiye", "it is he/she (who)", "Yeye ndiye mwalimu wetu.", "She's the one who is our teacher.", "Mira"),
	)

	s = addSkillL(db, sw, u, "Conditionals: -ki-, -nge-, -ngali-", "If, would and would have.", "Scale", "#00C2A8", 21, 320)
	l = addLesson(db, s, "If…, Would…", 1, 22,
		char(finch, "-ki- means 'if / when' (real conditions): Ukija, tutakwenda (if you come, we'll go); Mvua ikinyesha, sitakuja. -nge- is 'would' (hypothetical): Ningekuwa na pesa, ningenunua gari (if I had money, I'd buy a car). -ngali- is 'would have' (past unreal): Ungalikuja, ungalimwona (if you had come, you'd have seen her). Kama and ikiwa also mean 'if'."),
		mc("If you come, we'll go.", "Real condition", "Ukija, tutakwenda.", "Ukija, tutakwenda.", "Ungekuja, tungekwenda.", "Ulikuja, tulikwenda.", "Uje, twende."),
		mc("If I had money, I would buy a car.", "Hypothetical", "Ningekuwa na pesa, ningenunua gari.", "Ningekuwa na pesa, ningenunua gari.", "Nikiwa na pesa, nitanunua gari.", "Nilikuwa na pesa, nilinunua gari.", "Nina pesa, ninanunua gari."),
		mc("What does 'ungalikuja, ungalimwona' mean?", "-ngali-", "If you had come, you would have seen her", "If you had come, you would have seen her", "If you come, you'll see her", "When you came, you saw her", "Come and see her"),
		mc("If it rains, I won't come.", "Mvua ___, sitakuja.", "ikinyesha", "ikinyesha", "ingenyesha", "ilinyesha", "inanyesha"),
		tr("Translate", "If you study, you will pass", "Ukisoma, utafaulu"),
		fill("Hypothetical: I would help you.", "Ni___kusaidia.", "nge"),
		speak("Zephyr", "Ningekuwa na muda, ningejifunza kupiga gitaa."),
	)
	addVocab(db, l,
		vw("-ki-", "if / when (real)", "Ukifika, nipigie simu.", "When you arrive, call me.", "Riko"),
		vw("-nge-", "would (hypothetical)", "Ningependa kwenda.", "I would like to go.", "Mira"),
		vw("-ngali-", "would have", "Ningalijua, nisingekuja.", "Had I known, I wouldn't have come.", finch),
		vw("kufaulu", "to pass, succeed", "Alifaulu mtihani.", "She passed the exam.", "Pip"),
		vw("kama / ikiwa", "if", "Kama una muda, njoo.", "If you have time, come.", "Lumora"),
	)

	s = addSkillL(db, sw, u, "Habits & Storytelling", "hu- for habits, -po- for 'when', and narrative -ka-.", "BookOpen", "#F5A623", 22, 340)
	l = addLesson(db, s, "Hadithi — Telling a Story", 1, 22,
		char(finch, "hu- marks habits for any person: Mimi husoma kila siku (I read every day), Wao hucheza mpira Jumamosi. -po- means 'when': Nilipokuwa mtoto… (when I was a child…). Stories chain events with -ka- (and then): Alikwenda sokoni, akanunua samaki, akarudi nyumbani. Folk tales open with 'Paukwa?' — 'Pakawa!' and 'Hapo zamani za kale…' (once upon a time)."),
		mc("I usually wake up early.", "Habit", "Mimi huamka mapema.", "Mimi huamka mapema.", "Mimi niliamka mapema.", "Mimi nitaamka mapema.", "Mimi nimeamka mapema."),
		mc("when I was a child", "Nili___kuwa mtoto", "po", "po", "ye", "cho", "ki"),
		mc("She went to the market, bought fish and returned home.", "Narrative -ka-", "Alikwenda sokoni, akanunua samaki, akarudi nyumbani.", "Alikwenda sokoni, akanunua samaki, akarudi nyumbani.", "Anakwenda sokoni, ananunua samaki, anarudi nyumbani.", "Atakwenda sokoni, atanunua samaki, atarudi.", "Nenda sokoni, nunua samaki, rudi nyumbani."),
		mc("How do Swahili folk tales traditionally open?", "The storyteller says…", "Paukwa? — Pakawa!", "Paukwa? — Pakawa!", "Kitendawili? — Tega!", "Hodi? — Karibu!", "Mambo? — Poa!"),
		tr("Translate", "When I was a child, I lived in Mombasa", "Nilipokuwa mtoto, niliishi Mombasa"),
		write("Write a short story in Swahili (about 50 words) about something that happened to you on a journey. Use -li- and the narrative -ka- at least twice.",
			"Mwaka jana nilisafiri kwenda Arusha kwa basi. Njiani basi liliharibika, tukasimama kwa saa mbili. Abiria wakashuka, wakaanza kuongea na kucheka. Mzee mmoja akatuambia hadithi za zamani. Mwishowe fundi akafika, akatengeneza basi, tukaendelea na safari. Nilifika usiku, lakini nilifurahi sana."),
	)
	addVocab(db, l,
		vw("hu-", "habitual (usually)", "Yeye hunywa kahawa asubuhi.", "He usually drinks coffee in the morning.", "Nana"),
		vw("hadithi", "story", "Bibi alituambia hadithi.", "Grandma told us a story.", "Nana"),
		vw("Hapo zamani za kale", "once upon a time", "Hapo zamani za kale paliishi sungura.", "Once upon a time there lived a hare.", "Nana"),
		vw("kurudi", "to return", "Akarudi nyumbani.", "And he returned home.", "Zephyr"),
		vw("mwishowe", "in the end", "Mwishowe walifurahi.", "In the end they were happy.", "Lumora"),
	)

	s = addSkillL(db, sw, u, "Locatives: PA-, KU-, MU-", "pana, kuna, mna — specific, general and inside.", "Compass", "#FF5C5C", 23, 360)
	l = addLesson(db, s, "There Is, In There", 1, 22,
		char(finch, "Place has its own classes. PA- is a specific spot: Pana mtu mlangoni (there's someone at the door), mahali pazuri. KU- is general or remote: Kuna watu wengi sokoni (there are many people at the market). MU- is inside: Mna maji ndani ya chupa (there's water in the bottle). The -ni ending often triggers them: Nyumbani kuna wageni."),
		mc("There are many people at the market (general).", "___ watu wengi sokoni.", "Kuna", "Kuna", "Pana", "Mna", "Wana"),
		mc("There's water inside the bottle.", "___ maji ndani ya chupa.", "Mna", "Mna", "Kuna", "Pana", "Ina"),
		mc("a nice place (specific)", "mahali ___", "pazuri", "pazuri", "kuzuri", "mzuri", "nzuri"),
		mc("What does 'Hakuna matata' literally mean?", "Hakuna matata", "There are no problems", "There are no problems", "I have no money", "Nobody is here", "No more food"),
		tr("Translate", "There are guests at home", "Nyumbani kuna wageni"),
		fill("Negative: There's nobody here.", "Ha___ mtu hapa.", "pana"),
		speak("Riko", "Kuna duka karibu na hapa?"),
	)
	addVocab(db, l,
		vw("kuna / hakuna", "there is / there isn't", "Hakuna shida.", "No problem.", "Blaze"),
		vw("pana", "there is (at a specific spot)", "Pana nafasi hapa.", "There's room here.", "Riko"),
		vw("mna", "there is (inside)", "Mna nini ndani?", "What's inside?", "Pip"),
		vw("mahali", "place", "Mahali hapa ni pazuri.", "This place is lovely.", "Zephyr"),
		vw("mgeni / wageni", "guest, stranger / guests", "Karibuni wageni!", "Welcome, guests!", "Nana"),
	)

	s = addSkillL(db, sw, u, "Health & the Body", "At the clinic: symptoms, advice and sympathy.", "HeartPulse", "#17A3DD", 24, 380)
	l = addLesson(db, s, "Kwa Daktari — At the Doctor's", 1, 22,
		char("Nana", "Symptoms: Ninaumwa na kichwa (I have a headache), Nina homa (I have a fever), Tumbo linauma (my stomach hurts), Nakohoa (I'm coughing). The doctor may say: Kunywa dawa mara tatu kwa siku (take the medicine three times a day), Pumzika (rest). Friends say Pole, upone haraka! (sorry — get well soon!)."),
		mc("I have a headache.", "kichwa", "Ninaumwa na kichwa.", "Ninaumwa na kichwa.", "Kichwa changu ni kikubwa.", "Nina kichwa.", "Ninapenda kichwa."),
		mc("I have a fever.", "homa", "Nina homa.", "Nina homa.", "Ni homa.", "Niko homa.", "Homa yangu."),
		mc("Get well soon!", "To a sick friend", "Upone haraka!", "Upone haraka!", "Karibu tena!", "Safari njema!", "Hongera!"),
		mc("Take the medicine three times a day.", "Doctor's advice", "Kunywa dawa mara tatu kwa siku.", "Kunywa dawa mara tatu kwa siku.", "Kula dawa siku tatu.", "Dawa ni tatu.", "Usinywe dawa."),
		tr("Translate", "My stomach hurts", "Tumbo linauma"),
		tr("Translate", "You need to rest", "Unahitaji kupumzika"),
		speak("Nana", "Pole sana. Pumzika na unywe maji mengi."),
	)
	addVocab(db, l,
		vw("daktari", "doctor", "Nimekwenda kwa daktari.", "I've been to the doctor.", "Nana"),
		vw("homa", "fever", "Mtoto ana homa kali.", "The child has a high fever.", "Nana"),
		vw("dawa", "medicine", "Dawa hii ni chungu.", "This medicine is bitter.", "Pip"),
		vw("kuumwa", "to be ill; to hurt", "Anaumwa tangu jana.", "She's been ill since yesterday.", "Nana"),
		vw("kupumzika", "to rest", "Pumzika kidogo.", "Rest a little.", "Lumora"),
	)

	s = addSkillL(db, sw, u, "Work & Opinions", "Talk about jobs and give your view: nadhani, kwa maoni yangu.", "Briefcase", "#9B51E0", 25, 400)
	l = addLesson(db, s, "Kazini — At Work", 1, 22,
		char("Mira", "Opinions: Nadhani… / Nafikiri… (I think…), Kwa maoni yangu… (in my opinion…), Nakubali (I agree), Sikubali (I disagree), Ni kweli, lakini… (that's true, but…). At work: mkutano (meeting), ratiba (schedule), mshahara (salary), likizo (leave), meneja."),
		mc("In my opinion…", "Giving a view", "Kwa maoni yangu…", "Kwa maoni yangu…", "Kwa sababu yangu…", "Kwa kazi yangu…", "Kwa jina langu…"),
		mc("I disagree.", "kukubali", "Sikubali.", "Sikubali.", "Nakubali.", "Hakubali.", "Sijui."),
		mc("What is 'mshahara'?", "Mshahara wangu ni mdogo.", "salary", "salary", "meeting", "holiday", "manager"),
		mc("That's true, but…", "Polite disagreement", "Ni kweli, lakini…", "Ni kweli, lakini…", "Si kweli kabisa!", "Hapana, wewe!", "Sawa sawa."),
		tr("Translate", "The meeting will start at 3 p.m.", "Mkutano utaanza saa tisa mchana"),
		write("Write a short paragraph in Swahili (about 60 words) giving your opinion on working from home. Use 'kwa maoni yangu', 'kwa sababu' and 'lakini'.",
			"Kwa maoni yangu, kufanya kazi nyumbani kuna faida nyingi. Kwanza, hakuna haja ya kusafiri kila siku, kwa hiyo tunaokoa muda na pesa. Pili, mtu anaweza kupanga ratiba yake vizuri. Lakini kuna changamoto pia, kwa sababu wafanyakazi hawaonani mara kwa mara. Nadhani ni bora kuchanganya kazi ya nyumbani na ofisini."),
	)
	addVocab(db, l,
		vw("nadhani / nafikiri", "I think", "Nadhani uko sahihi.", "I think you're right.", "Mira"),
		vw("maoni", "opinion(s)", "Una maoni gani?", "What's your opinion?", finch),
		vw("mkutano", "meeting", "Tuna mkutano saa nne.", "We have a meeting at 10 a.m.", "Mira"),
		vw("mshahara", "salary", "Mshahara utalipwa kesho.", "Salaries will be paid tomorrow.", "Blaze"),
		vw("changamoto", "challenge", "Kazi hii ina changamoto nyingi.", "This job has many challenges.", "Zephyr"),
	)
}

// ───────────────────── B2 — Lugha Rasmi: Formal Swahili & Society ─────────────────────

func seedSwahiliB2(db *gorm.DB) {
	const sw = "sw"
	const u = "B2 · Lugha Rasmi — Formal Swahili & Society"
	finch := "Professor Finch"

	s := addSkillL(db, sw, u, "Verb Extensions I: Applicative & Causative", "-ia/-ea (for, to) and -isha/-esha (make, cause).", "Layers", "#6C3FC5", 26, 440)
	l := addLesson(db, s, "Doing For & Making Do", 1, 24,
		char(finch, "Extensions change a verb's meaning. Applicative -ia/-ea adds 'for / to / at': andika → andikia (write to), pika → pikia (cook for), soma → somea (read to/for). Vowel harmony: roots with e or o take -ea: soma → somea. Causative -isha/-esha (or -sha, -za) means 'make / cause': soma → somesha (teach, make study), ona → onyesha (show), lala → laza (put to sleep), kula → lisha (feed)."),
		mc("Write to me.", "andika + applicative", "Niandikie.", "Niandikie.", "Niandike.", "Niandikishe.", "Niandikwe."),
		mc("Mum cooked for us.", "pika + applicative", "Mama alitupikia.", "Mama alitupikia.", "Mama alitupika.", "Mama alitupikisha.", "Mama alipikwa."),
		mc("What does 'kuonyesha' mean?", "ona + causative", "to show", "to show", "to be seen", "to see each other", "to look for"),
		mc("The teacher teaches (makes study) the pupils.", "soma + causative", "Mwalimu anawasomesha wanafunzi.", "Mwalimu anawasomesha wanafunzi.", "Mwalimu anawasomea wanafunzi.", "Mwalimu anasomwa na wanafunzi.", "Mwalimu anasomana na wanafunzi."),
		tr("Translate", "Show me the way", "Nionyeshe njia"),
		fill("Applicative", "Ninakuletea zawadi. (I'm bringing you a gift — leta + ___)", "ea"),
		speak("Lumora", "Nitakuandikia barua pepe kesho asubuhi."),
	)
	addVocab(db, l,
		vw("kuandikia", "to write to", "Niandikie ujumbe.", "Write me a message.", finch),
		vw("kuletea", "to bring (to/for)", "Nimekuletea zawadi.", "I've brought you a present.", "Pip"),
		vw("kusomesha", "to teach, educate", "Alisomesha watoto wake wote.", "She put all her children through school.", "Nana"),
		vw("kuonyesha", "to show", "Nionyeshe picha.", "Show me the photo.", "Mira"),
		vw("kulisha", "to feed", "Analisha ng'ombe.", "He's feeding the cattle.", "Zephyr"),
	)

	s = addSkillL(db, sw, u, "Verb Extensions II: Passive, Stative & Reciprocal", "-wa, -ika/-eka, -ana.", "Shield", "#00C2A8", 27, 470)
	l = addLesson(db, s, "Done, Breakable, Each Other", 1, 24,
		char(finch, "Passive -wa: andika → andikwa (be written), piga → pigwa; the agent takes na: Barua iliandikwa na Juma. Verbs whose root ends in a vowel take -liwa/-lewa: fungua → funguliwa (be opened), and kula → kuliwa (be eaten). Stative -ika/-eka shows a state or possibility: vunja → vunjika (get broken), soma → someka (be legible). Reciprocal -ana: penda → pendana (love each other), ona → onana (see each other)."),
		mc("The letter was written by Juma.", "Passive", "Barua iliandikwa na Juma.", "Barua iliandikwa na Juma.", "Barua iliandika Juma.", "Juma aliandikiwa barua.", "Barua iliandikana na Juma."),
		mc("The cup got broken.", "vunja + stative", "Kikombe kimevunjika.", "Kikombe kimevunjika.", "Kikombe kimevunja.", "Kikombe kimevunjwa na.", "Kikombe kimevunjana."),
		mc("They love each other.", "penda + reciprocal", "Wanapendana.", "Wanapendana.", "Wanapendwa.", "Wanapendeka.", "Wanapendea."),
		mc("The door was opened.", "fungua + passive", "Mlango ulifunguliwa.", "Mlango ulifunguliwa.", "Mlango ulifungua.", "Mlango ulifungwa.", "Mlango ulifungukana."),
		tr("Translate", "We will see each other tomorrow", "Tutaonana kesho"),
		mc("What does 'Maandishi haya hayasomeki' mean?", "soma + stative", "This writing is illegible", "This writing is illegible", "Nobody has read this writing", "Don't read this writing", "This writing reads itself"),
		speak("Zephyr", "Daraja lilijengwa na serikali mwaka jana."),
	)
	addVocab(db, l,
		vw("-wa (passive)", "be done", "Gari liliibiwa usiku.", "The car was stolen at night.", finch),
		vw("kuvunjika", "to get broken", "Simu imevunjika.", "The phone is broken.", "Blaze"),
		vw("kuonana", "to see each other", "Tutaonana baadaye.", "See you later.", "Lumora"),
		vw("kujengwa", "to be built", "Shule mpya itajengwa.", "A new school will be built.", "Zephyr"),
		vw("kusikika", "to be audible", "Sauti yako haisikiki.", "Your voice can't be heard.", "Riko"),
	)

	s = addSkillL(db, sw, u, "Formal Letters & Register", "Barua rasmi: Kwa, Yah:, Ndugu, Wako mtiifu.", "PenLine", "#F5A623", 28, 500)
	l = addLesson(db, s, "Barua Rasmi — Formal Letters", 1, 24,
		char("Mira", "A formal Swahili letter has the sender's address top right, the date, the recipient (Kwa: Meneja Mkuu, …), then Ndugu/Mheshimiwa (Dear Sir/Madam), a subject line YAH: OMBI LA KAZI (Re: job application), the body, and Wako mtiifu (yours faithfully) with the name and signature. Formal register uses polite forms: Napenda kuomba… (I would like to request…), Naomba kuwasilisha… (I submit…)."),
		mc("What does 'YAH:' introduce?", "YAH: OMBI LA KAZI", "The subject line (Re:)", "The subject line (Re:)", "The date", "The signature", "The greeting"),
		mc("Formal closing (Yours faithfully)", "End of a formal letter", "Wako mtiifu,", "Wako mtiifu,", "Kwaheri rafiki,", "Tutaonana,", "Poa sana,"),
		mc("Formal opening to an honourable official", "Addressing a minister", "Mheshimiwa,", "Mheshimiwa,", "Mambo,", "Habari za jana,", "Rafiki yangu,"),
		mc("I would like to apply for the post of…", "Formal request", "Napenda kuomba nafasi ya kazi ya…", "Napenda kuomba nafasi ya kazi ya…", "Nipe kazi ya…", "Nataka kazi sasa…", "Kazi ni yangu…"),
		tr("Translate", "Thank you for your cooperation", "Asante kwa ushirikiano wako"),
		write("Write a short formal letter in Swahili (about 80 words) applying for a job as a teacher. Include Kwa, YAH:, a greeting, your qualifications, and Wako mtiifu.",
			"Kwa: Mkuu wa Shule,\nShule ya Sekondari Amani,\nS.L.P. 120, Nairobi.\n\nNdugu,\nYAH: OMBI LA KAZI YA UALIMU\n\nNapenda kuomba nafasi ya kazi ya mwalimu wa Kiswahili iliyotangazwa gazetini. Nina shahada ya elimu na uzoefu wa miaka mitatu katika kufundisha. Ninaamini kwamba ujuzi wangu utasaidia wanafunzi wenu. Naambatanisha vyeti vyangu. Natumaini kupata jibu lako.\n\nWako mtiifu,\nAmina Juma"),
	)
	addVocab(db, l,
		vw("barua rasmi", "formal letter", "Andika barua rasmi kwa meneja.", "Write a formal letter to the manager.", "Mira"),
		vw("ombi", "request, application", "Ombi lako limekubaliwa.", "Your application has been accepted.", finch),
		vw("Mheshimiwa", "Honourable (title)", "Mheshimiwa Waziri.", "The Honourable Minister.", finch),
		vw("ushirikiano", "cooperation", "Tunashukuru kwa ushirikiano.", "We are grateful for the cooperation.", "Mira"),
		vw("Wako mtiifu", "Yours faithfully", "Wako mtiifu, Juma.", "Yours faithfully, Juma.", "Mira"),
	)

	s = addSkillL(db, sw, u, "News & Society", "Habari: serikali, uchaguzi, uchumi — understand the news.", "Globe", "#FF5C5C", 29, 530)
	l = addLesson(db, s, "Taarifa ya Habari — The News", 1, 24,
		char(finch, "News Swahili (BBC Swahili, VOA, Radio Tanzania, KBC) has its own vocabulary: serikali (government), rais (president), bunge (parliament), uchaguzi (election), uchumi (economy), mfumuko wa bei (inflation), mafuriko (floods), ukame (drought), wananchi (citizens). The passive is everywhere: Imeripotiwa kwamba… (it has been reported that…)."),
		mc("What does 'uchaguzi' mean?", "uchaguzi mkuu", "election", "election", "economy", "parliament", "drought"),
		mc("What does 'mfumuko wa bei' mean?", "Mfumuko wa bei umeongezeka.", "inflation", "inflation", "price list", "a sale", "taxes"),
		mc("It has been reported that…", "News phrasing", "Imeripotiwa kwamba…", "Imeripotiwa kwamba…", "Nimeripoti kwamba…", "Ripoti yangu…", "Tunaripoti kesho…"),
		mc("What does this headline say?", "Mafuriko yaathiri maelfu ya wananchi", "Floods affect thousands of citizens", "Floods affect thousands of citizens", "Drought ends in the north", "Thousands vote in the election", "Citizens protest high prices"),
		tr("Translate", "The government has announced new measures", "Serikali imetangaza hatua mpya"),
		mc("What is 'bunge'?", "Bunge limepitisha sheria.", "parliament", "parliament", "court", "market", "army"),
		speak("Zephyr", "Habari kutoka Dodoma: Bunge limepitisha bajeti mpya leo."),
	)
	addVocab(db, l,
		vw("serikali", "government", "Serikali imetangaza likizo.", "The government has declared a holiday.", finch),
		vw("uchaguzi", "election", "Uchaguzi utafanyika mwezi ujao.", "The election will take place next month.", finch),
		vw("uchumi", "economy", "Uchumi unakua.", "The economy is growing.", "Zephyr"),
		vw("wananchi", "citizens, the public", "Wananchi walijitokeza kupiga kura.", "Citizens turned out to vote.", "Zephyr"),
		vw("ukame", "drought", "Ukame umeathiri mavuno.", "The drought has affected the harvest.", "Nana"),
	)

	s = addSkillL(db, sw, u, "Loanwords & Word Origins", "Arabic, Persian, Portuguese, German and English in Swahili.", "Languages", "#17A3DD", 30, 560)
	l = addLesson(db, s, "Where Swahili Words Come From", 1, 24,
		char(finch, "Swahili is a Bantu language shaped by centuries of Indian Ocean trade. Arabic gave many abstract and religious words: kitabu, dunia (world), elimu (education), habari, saa, rafiki, hesabu. From Persian and Hindi: chai (tea). Portuguese: meza (table), pesa, bendera (flag). German (Tanganyika): shule (school). English: kompyuta, benki, polisi, baiskeli. The Bantu core: watu, -kubwa, kula, maji."),
		mc("'shule' (school) comes from which language?", "shule", "German", "German", "Arabic", "English", "Portuguese"),
		mc("'meza' (table) comes from which language?", "meza", "Portuguese", "Portuguese", "Arabic", "German", "Hindi"),
		mc("Which of these is an Arabic loanword?", "Arabic origin", "elimu (education)", "elimu (education)", "watu (people)", "kula (eat)", "baiskeli (bicycle)"),
		mc("Which word is from English?", "English origin", "baiskeli", "baiskeli", "dunia", "meza", "maji"),
		mc("Why does Swahili have so many Arabic words?", "History", "Centuries of Indian Ocean trade and Islam on the coast", "Centuries of Indian Ocean trade and Islam on the coast", "It was invented in Arabia", "Recent television", "Colonial schools only"),
		tr("Translate", "The world is big", "Dunia ni kubwa"),
		speak("Mira", "Elimu ni ufunguo wa maisha."),
	)
	addVocab(db, l,
		vw("dunia", "world (Arabic)", "Dunia ni kijiji.", "The world is a village.", finch),
		vw("elimu", "education (Arabic)", "Elimu ni muhimu.", "Education is important.", finch),
		vw("meza", "table (Portuguese)", "Weka chakula mezani.", "Put the food on the table.", "Cora"),
		vw("shule", "school (German)", "Watoto wako shuleni.", "The children are at school.", "Pip"),
		vw("kompyuta", "computer (English)", "Kompyuta yangu ni mpya.", "My computer is new.", "Blaze"),
	)
}

// ───────────────────── C1 — Fasihi na Methali: Literature & Proverbs ─────────────────────

func seedSwahiliC1(db *gorm.DB) {
	const sw = "sw"
	const u = "C1 · Fasihi na Methali — Literature & Proverbs"
	finch := "Professor Finch"

	s := addSkillL(db, sw, u, "Proverbs (Methali)", "Wisdom in a sentence — used in speech, debate and teaching.", "Quote", "#6C3FC5", 31, 700)
	l := addLesson(db, s, "Methali za Kiswahili", 1, 26,
		char("Nana", "Methali carry East African wisdom and are used to clinch an argument. Haraka haraka haina baraka (haste has no blessing). Pole pole ndio mwendo (slowly is the way to go). Bandu bandu huisha gogo (chip by chip finishes the log — persistence). Usipoziba ufa utajenga ukuta (if you don't fill the crack you'll build a wall). Kidole kimoja hakivunji chawa (one finger can't kill a louse — unity). Asiyesikia la mkuu huvunjika guu (whoever ignores elders breaks a leg)."),
		mc("Which proverb warns against haste?", "Don't rush", "Haraka haraka haina baraka", "Haraka haraka haina baraka", "Bandu bandu huisha gogo", "Kidole kimoja hakivunji chawa", "Mgeni njoo mwenyeji apone"),
		mc("What does 'Usipoziba ufa utajenga ukuta' teach?", "ufa = crack", "Fix small problems before they grow", "Fix small problems before they grow", "Build strong houses", "Walls are better than doors", "Never repair anything"),
		mc("Which proverb is about unity?", "Working together", "Kidole kimoja hakivunji chawa", "Kidole kimoja hakivunji chawa", "Haraka haraka haina baraka", "Pole pole ndio mwendo", "Asiyesikia la mkuu huvunjika guu"),
		mc("What does 'Bandu bandu huisha gogo' mean?", "bandu = chip; gogo = log", "Little by little, a big task gets done", "Little by little, a big task gets done", "Wood is valuable", "Cut trees carefully", "Don't waste time"),
		mc("'Asiyesikia la mkuu huvunjika guu' advises…", "mkuu = elder", "Listen to the advice of elders", "Listen to the advice of elders", "Walk carefully", "Elders break legs", "Ignore old customs"),
		fill("Complete the proverb", "Pole pole ndio ___.", "mwendo"),
		speak("Nana", "Haraka haraka haina baraka, mjukuu wangu."),
	)
	addVocab(db, l,
		vw("methali", "proverb", "Wazee hutumia methali nyingi.", "Elders use many proverbs.", "Nana"),
		vw("Haraka haraka haina baraka", "haste has no blessing", "Usikimbie — haraka haraka haina baraka.", "Don't rush — haste brings no blessing.", "Nana"),
		vw("Pole pole ndio mwendo", "slowly is the way to go", "Jifunze polepole; pole pole ndio mwendo.", "Learn steadily; slow and steady.", "Lumora"),
		vw("Usipoziba ufa utajenga ukuta", "a stitch in time saves nine", "Tatua tatizo sasa — usipoziba ufa utajenga ukuta.", "Solve it now, before it grows.", finch),
		vw("Kidole kimoja hakivunji chawa", "unity is strength", "Tushirikiane; kidole kimoja hakivunji chawa.", "Let's cooperate; one finger can't do it alone.", "Zephyr"),
	)

	s = addSkillL(db, sw, u, "Idioms & Riddles (Nahau na Vitendawili)", "Kupiga domo, kuvunja moyo — and 'Kitendawili? Tega!'", "Brain", "#00C2A8", 32, 730)
	l = addLesson(db, s, "Figurative Swahili", 1, 26,
		char(finch, "Nahau (idioms) mean more than their words: kupiga domo (to chatter idly — 'beat the lip'), kuvunja moyo (to discourage — 'break the heart'), kupata jiko (to get married, of a man — 'get a stove'), kufa moyo (to lose heart). Riddles are a social game: the asker says Kitendawili! (riddle!), the listener Tega! (set it!). 'Kila nikienda ananifuata' (wherever I go it follows me) — jibu: kivuli (shadow)."),
		mc("What does 'kuvunja moyo' mean?", "Usinivunje moyo!", "to discourage", "to discourage", "to have a heart attack", "to fall in love", "to break a promise"),
		mc("What does 'kupiga domo' mean?", "Acheni kupiga domo!", "to chatter idly", "to chatter idly", "to kiss", "to shout an order", "to eat noisily"),
		mc("The riddle-giver says 'Kitendawili!' — you reply:", "Starting a riddle game", "Tega!", "Tega!", "Paukwa!", "Karibu!", "Pakawa!"),
		mc("'Kila nikienda ananifuata.' What is it?", "Riddle: wherever I go, it follows me", "Kivuli (a shadow)", "Kivuli (a shadow)", "Mbwa (a dog)", "Upepo (the wind)", "Jua (the sun)"),
		mc("'Nyumba yangu haina mlango.' What is it?", "Riddle: my house has no door", "Yai (an egg)", "Yai (an egg)", "Pango (a cave)", "Kaburi (a grave)", "Mfuko (a bag)"),
		mc("What does 'kufa moyo' mean?", "Usife moyo!", "to lose heart", "to lose heart", "to die suddenly", "to be very brave", "to love deeply"),
		speak("Pip", "Kitendawili! — Tega! — Kila nikienda ananifuata. — Kivuli!"),
	)
	addVocab(db, l,
		vw("nahau", "idiom", "Nahau hii ina maana gani?", "What does this idiom mean?", finch),
		vw("kuvunja moyo", "to discourage", "Matokeo yalinivunja moyo.", "The results discouraged me.", "Mira"),
		vw("kitendawili", "riddle", "Nipe kitendawili!", "Give me a riddle!", "Pip"),
		vw("Tega!", "Set it! (reply to Kitendawili)", "Kitendawili! — Tega!", "Riddle! — Go ahead!", "Pip"),
		vw("kivuli", "shadow", "Tukae kivulini.", "Let's sit in the shade.", "Zephyr"),
	)

	s = addSkillL(db, sw, u, "Literature & Poetry", "Shaaban Robert, Abunuwasi, and the rules of shairi.", "BookOpen", "#F5A623", 33, 760)
	l = addLesson(db, s, "Fasihi ya Kiswahili", 1, 26,
		char(finch, "Swahili literature is old and rich. Traditional poetry (shairi) has rules: beti (stanzas), mistari (lines), mizani (syllables per line, often 16 split 8+8) and vina (rhymes at the middle and end of each line). Shaaban Robert (1909–1962), 'the father of Swahili prose', wrote Kusadikika and Maisha Yangu. The trickster tales of Abunuwasi (Hekaya za Abunuwasi) are classroom favourites."),
		mc("What are 'mizani' in a shairi?", "Kanuni za shairi", "The syllable count per line", "The syllable count per line", "The rhymes", "The stanzas", "The poet's name"),
		mc("What are 'vina'?", "Kanuni za shairi", "Rhymes (middle and end of lines)", "Rhymes (middle and end of lines)", "Syllables", "Stanzas", "Titles"),
		mc("Who is called the father of Swahili prose?", "Mwandishi maarufu", "Shaaban Robert", "Shaaban Robert", "Abunuwasi", "Fumo Liyongo", "Julius Nyerere"),
		mc("What is a 'beti'?", "Shairi lina beti tano.", "A stanza", "A stanza", "A line", "A rhyme", "A proverb"),
		mc("Abunuwasi is best known as…", "Hekaya za Abunuwasi", "A clever trickster in folk tales", "A clever trickster in folk tales", "A colonial governor", "A modern rapper", "A famous runner"),
		mc("How many syllables are in a typical shairi line?", "8 + 8", "16", "16", "8", "12", "24"),
		speak("Zephyr", "Kiswahili ni lugha, ya pwani na ya bara."),
	)
	addVocab(db, l,
		vw("fasihi", "literature", "Ninasoma fasihi ya Kiswahili.", "I study Swahili literature.", finch),
		vw("shairi / mashairi", "poem / poems", "Aliandika shairi zuri.", "She wrote a beautiful poem.", "Mira"),
		vw("mwandishi", "writer, author", "Shaaban Robert ni mwandishi maarufu.", "Shaaban Robert is a famous writer.", finch),
		vw("riwaya", "novel", "Riwaya hii inasisimua.", "This novel is thrilling.", "Zephyr"),
		vw("mizani / vina", "syllable count / rhymes", "Shairi hili lina mizani kumi na sita.", "This poem has sixteen syllables per line.", finch),
	)

	s = addSkillL(db, sw, u, "Debate & Argument (Mjadala)", "Structure an argument and rebut politely.", "MessageCircle", "#FF5C5C", 34, 790)
	l = addLesson(db, s, "Hoja na Mjadala", 1, 26,
		char(finch, "School and radio debates (mijadala) follow a ritual: Mwenyekiti, waheshimiwa majaji, … (Chairperson, honourable judges…). Build hoja (arguments) with connectors: kwanza (first), pili (second), zaidi ya hayo (moreover), kwa upande mwingine (on the other hand), hata hivyo (nevertheless), kwa hiyo / hivyo basi (therefore), kwa kuhitimisha (in conclusion)."),
		mc("on the other hand", "Contrast", "kwa upande mwingine", "kwa upande mwingine", "zaidi ya hayo", "kwa hiyo", "kwanza"),
		mc("nevertheless", "Concession", "hata hivyo", "hata hivyo", "hivyo basi", "pili", "kwa mfano"),
		mc("in conclusion", "Ending", "kwa kuhitimisha", "kwa kuhitimisha", "kwa mfano", "kwanza kabisa", "zaidi ya hayo"),
		mc("How does a school debate speech traditionally open?", "Opening", "Mwenyekiti, waheshimiwa majaji…", "Mwenyekiti, waheshimiwa majaji…", "Mambo vipi wote!", "Kwaherini!", "Shikamoo watoto…"),
		mc("What is a 'hoja'?", "Hoja yako ina nguvu.", "an argument / point", "an argument / point", "a judge", "a vote", "a question"),
		write("Write an argumentative paragraph in Swahili (about 100 words) on 'Simu za mkononi zipigwe marufuku shuleni' (mobile phones should be banned in schools). Give two arguments, a counterargument with 'hata hivyo', and a conclusion.",
			"Mwenyekiti, waheshimiwa majaji, ninaunga mkono hoja kwamba simu za mkononi zipigwe marufuku shuleni. Kwanza, simu huwavuruga wanafunzi darasani na kupunguza umakini wao. Pili, baadhi ya wanafunzi huzitumia kuibia mitihani. Kwa upande mwingine, wapo wanaosema kwamba simu husaidia katika utafiti. Hata hivyo, shule zina maktaba na kompyuta kwa kazi hiyo. Kwa kuhitimisha, ni wazi kwamba madhara ya simu shuleni ni makubwa kuliko faida zake, hivyo basi zipigwe marufuku."),
	)
	addVocab(db, l,
		vw("mjadala", "debate", "Mjadala ulikuwa mkali.", "The debate was heated.", finch),
		vw("hoja", "argument, point", "Hoja yako ina mashiko.", "Your argument holds water.", "Mira"),
		vw("hata hivyo", "nevertheless", "Ni ghali; hata hivyo, nitanunua.", "It's expensive; nevertheless, I'll buy it.", "Cora"),
		vw("kwa upande mwingine", "on the other hand", "Kwa upande mwingine, kuna faida.", "On the other hand, there are benefits.", finch),
		vw("kuunga mkono", "to support (a motion)", "Ninaunga mkono hoja hii.", "I support this motion.", "Blaze"),
	)

	s = addSkillL(db, sw, u, "Business & Professional Swahili", "Contracts, invoices, board meetings and email.", "Briefcase", "#17A3DD", 35, 820)
	l = addLesson(db, s, "Biashara na Kazi", 1, 26,
		char("Mira", "Professional Swahili: mkataba (contract), ankara (invoice), risiti (receipt), bajeti, faida (profit), hasara (loss), mkurugenzi (director), bodi ya wakurugenzi (board of directors), barua pepe (email). Formal email: Habari ndugu…, Natumaini u mzima (I hope you're well), Nakuandikia kuhusu… (I'm writing about…), Kwa heshima (respectfully)."),
		mc("What is 'ankara'?", "Tafadhali tuma ankara.", "an invoice", "an invoice", "a contract", "a receipt", "a salary"),
		mc("What is 'hasara'?", "Kampuni ilipata hasara.", "a loss", "a loss", "a profit", "a tax", "a loan"),
		mc("I hope you are well (email opening)", "Polite email opening", "Natumaini u mzima.", "Natumaini u mzima.", "Mambo vipi?", "Uko wapi sasa?", "Nipe pesa."),
		mc("I'm writing to you about…", "Email purpose line", "Nakuandikia kuhusu…", "Nakuandikia kuhusu…", "Nimekuandika jana…", "Andika kuhusu…", "Niandikie kesho…"),
		tr("Translate", "Please sign the contract", "Tafadhali tia sahihi mkataba"),
		mc("What does 'mkurugenzi' mean?", "Mkurugenzi mtendaji", "director", "director", "accountant", "customer", "guard"),
		speak("Mira", "Natumaini u mzima. Nakuandikia kuhusu mkutano wa bodi wiki ijayo."),
	)
	addVocab(db, l,
		vw("mkataba", "contract", "Tumesaini mkataba.", "We've signed the contract.", "Mira"),
		vw("ankara", "invoice", "Ankara imelipwa.", "The invoice has been paid.", "Cora"),
		vw("faida / hasara", "profit / loss", "Biashara ilipata faida.", "The business made a profit.", "Blaze"),
		vw("barua pepe", "email", "Nitumie barua pepe.", "Send me an email.", "Mira"),
		vw("mkurugenzi", "director", "Mkurugenzi atahudhuria mkutano.", "The director will attend the meeting.", finch),
	)
}

// ───────────────────── C2 — Umahiri: Mastery ─────────────────────

func seedSwahiliC2(db *gorm.DB) {
	const sw = "sw"
	const u = "C2 · Umahiri — Mastery"
	finch := "Professor Finch"

	s := addSkillL(db, sw, u, "Sheng & Urban Swahili", "Nairobi street language — understand it, use it carefully.", "Smartphone", "#6C3FC5", 36, 1000)
	l := addLesson(db, s, "Sheng ya Mtaani", 1, 28,
		char("Blaze", "Sheng mixes Swahili, English and local languages and changes fast — it's youth and street identity in Nairobi. Niaje? — Poa! (What's up? — Cool!). msee (guy), maze / buda (bro), mbogi (crew), doh / mullah (money), keja (house), mathree (matatu), form ni gani? (what's the plan?), fiti (fine, good). Understand it; use standard Swahili in formal settings."),
		mc("What does 'keja' mean in Sheng?", "Niko kejani.", "house, home", "house, home", "money", "friend", "car"),
		mc("What does 'doh' mean?", "Sina doh.", "money", "money", "food", "time", "phone"),
		mc("'Form ni gani?' means…", "Sheng greeting", "What's the plan?", "What's the plan?", "What's your form at school?", "Fill in this form", "Where's the shop?"),
		mc("Where is Sheng mainly spoken?", "Origin", "Urban Kenya, especially Nairobi", "Urban Kenya, especially Nairobi", "Zanzibar", "Rural Tanzania", "Eastern Congo only"),
		mc("When should you avoid Sheng?", "Register", "In formal letters and job interviews", "In formal letters and job interviews", "With friends", "In music", "On social media"),
		mc("Standard Swahili for 'mathree'", "Tushike mathree.", "matatu", "matatu", "pikipiki", "treni", "ndege"),
		speak("Blaze", "Niaje msee? Poa sana! Form ni gani leo?"),
	)
	addVocab(db, l,
		vw("Niaje? — Poa!", "What's up? — Cool! (Sheng)", "Niaje buda? Poa tu.", "What's up, bro? All good.", "Blaze"),
		vw("msee", "guy, person (Sheng)", "Huyo msee ni fiti.", "That guy is great.", "Blaze"),
		vw("mbogi", "crew, gang of friends (Sheng)", "Mbogi yetu inakuja.", "Our crew is coming.", "Blaze"),
		vw("keja", "house, home (Sheng)", "Karibu kwa keja yangu.", "Welcome to my place.", "Blaze"),
		vw("doh", "money (Sheng)", "Doh imeisha.", "The money's finished.", "Blaze"),
	)

	s = addSkillL(db, sw, u, "Regional Varieties", "Kiunguja, Kimvita, Kiamu — and Kenya vs Tanzania vs Congo.", "Globe", "#00C2A8", 37, 1030)
	l = addLesson(db, s, "Lahaja za Kiswahili", 1, 28,
		char(finch, "Standard Swahili is based on Kiunguja (Zanzibar town). Other dialects (lahaja) include Kimvita (Mombasa) and Kiamu (Lamu), home of classical poetry. National usage differs: Tanzanians say daladala and hela (money), Kenyans matatu and pesa. In eastern DR Congo, Kingwana Swahili has French loanwords (e.g. bureau). Tanzania uses Swahili far more widely in schools and government."),
		mc("Standard Swahili is based on which dialect?", "Kiswahili sanifu", "Kiunguja (Zanzibar)", "Kiunguja (Zanzibar)", "Kimvita (Mombasa)", "Kiamu (Lamu)", "Kingwana (Congo)"),
		mc("Which dialect is from Mombasa?", "Lahaja", "Kimvita", "Kimvita", "Kiamu", "Kiunguja", "Kipemba"),
		mc("A Tanzanian word for money", "Tanzania usage", "hela", "hela", "doh", "keja", "mbogi"),
		mc("Congolese Swahili (Kingwana) borrows heavily from…", "Loanwords", "French", "French", "German", "Portuguese", "Japanese"),
		mc("What is a 'lahaja'?", "Lahaja ya Kiamu", "a dialect", "a dialect", "a proverb", "a poem", "a newspaper"),
		mc("Which island town is the heart of classical Swahili poetry?", "Utenzi", "Lamu", "Lamu", "Nairobi", "Dodoma", "Kampala"),
		speak("Zephyr", "Kiswahili sanifu kimetokana na lahaja ya Kiunguja."),
	)
	addVocab(db, l,
		vw("lahaja", "dialect", "Kuna lahaja nyingi za Kiswahili.", "There are many Swahili dialects.", finch),
		vw("Kiswahili sanifu", "Standard Swahili", "Tuandike kwa Kiswahili sanifu.", "Let's write in Standard Swahili.", finch),
		vw("Kiunguja", "Zanzibar dialect (basis of the standard)", "Kiunguja huzungumzwa Zanzibar.", "Kiunguja is spoken in Zanzibar.", "Zephyr"),
		vw("hela", "money (Tanzania)", "Sina hela leo.", "I have no money today.", "Cora"),
		vw("pwani", "coast", "Watu wa pwani.", "Coastal people.", "Zephyr"),
	)

	s = addSkillL(db, sw, u, "Classical Poetry (Utenzi)", "Utendi wa Tambuka, Fumo Liyongo and the Lamu tradition.", "Landmark", "#F5A623", 38, 1060)
	l = addLesson(db, s, "Utenzi na Ushairi wa Kale", 1, 28,
		char(finch, "Utenzi (or utendi) is a long narrative poem in four-line stanzas of eight-syllable lines — epics of faith, war and history. Utendi wa Tambuka (1728) is among the oldest dated Swahili manuscripts, originally in Arabic script. Fumo Liyongo, a legendary poet-hero of the northern coast, is celebrated in Utendi wa Fumo Liyongo. Modern writers like Shaaban Robert used the form for Utenzi wa Vita vya Uhuru."),
		mc("How many syllables are in each line of an utenzi?", "Utenzi", "8", "8", "16", "12", "4"),
		mc("Utendi wa Tambuka was originally written in…", "Hati ya zamani", "Arabic script", "Arabic script", "Latin script", "Cyrillic", "Ge'ez"),
		mc("Who is Fumo Liyongo?", "Shujaa wa pwani", "A legendary poet-hero of the northern coast", "A legendary poet-hero of the northern coast", "A Tanzanian president", "A modern novelist", "A Portuguese explorer"),
		mc("What is an utenzi?", "Fasihi", "A long narrative poem", "A long narrative poem", "A short proverb", "A riddle", "A newspaper column"),
		mc("When was Utendi wa Tambuka composed?", "Historia", "1728", "1728", "1928", "1828", "1628"),
		speak("Nana", "Bismillahi kwanza, ndiyo kauli ya mwanzo."),
	)
	addVocab(db, l,
		vw("utenzi / tenzi", "epic poem / poems", "Utenzi huu una beti mia.", "This epic has a hundred stanzas.", finch),
		vw("ushairi", "poetry", "Ushairi wa Kiswahili ni wa kale.", "Swahili poetry is ancient.", "Nana"),
		vw("shujaa", "hero", "Liyongo alikuwa shujaa.", "Liyongo was a hero.", "Zephyr"),
		vw("hati", "manuscript, document", "Hati ya zamani.", "An old manuscript.", finch),
	)

	s = addSkillL(db, sw, u, "Translation & Style", "Ufasaha: idiomatic, natural Swahili — and false friends.", "Languages", "#FF5C5C", 39, 1090)
	l = addLesson(db, s, "Ufasaha wa Lugha", 1, 28,
		char(finch, "Mastery means sounding natural, not translated. Watch false friends and calques: 'I am fine' is Sijambo / Nzuri, not *Mimi ni fine. Swahili prefers verbs to abstract nouns: 'the completion of the project' → mradi ulipokamilika. Agreement must stay consistent across long sentences: Vitabu vile vizuri ambavyo nilivinunua vimepotea. Choose register: kufariki (pass away, respectful) vs kufa (die)."),
		mc("Most natural Swahili for 'He passed away' (respectfully)", "Register", "Alifariki dunia.", "Alifariki dunia.", "Alikufa tu.", "Alipotea maisha.", "Alimaliza kifo."),
		mc("Which sentence has correct agreement throughout?", "Those good books I bought have got lost.", "Vitabu vile vizuri ambavyo nilivinunua vimepotea.", "Vitabu vile vizuri ambavyo nilivinunua vimepotea.", "Vitabu vile mzuri ambacho nilikinunua kimepotea.", "Vitabu ile nzuri ambayo nilinunua imepotea.", "Vitabu wale wazuri ambao niliwanunua wamepotea."),
		mc("Natural Swahili for 'I'm fine' (reply to Hujambo)", "Avoid calques", "Sijambo.", "Sijambo.", "Mimi ni fine.", "Mimi ni nzuri sana mtu.", "Niko poa ya kawaida."),
		mc("What does 'ufasaha' mean?", "Ufasaha wa lugha", "fluency, eloquence", "fluency, eloquence", "translation", "grammar mistake", "dialect"),
		tr("Translate naturally", "Congratulations on your new job", "Hongera kwa kazi mpya"),
		write("Translate this into natural Swahili (about 60 words), paying attention to agreement and register: 'Dear colleagues, I am pleased to inform you that our project has been completed successfully. I thank everyone who contributed. The final report will be sent to you next week. Let us continue working together with the same spirit.'",
			"Wapendwa wenzangu, nina furaha kuwajulisha kwamba mradi wetu umekamilika kwa mafanikio. Nawashukuru wote waliochangia. Ripoti ya mwisho itatumwa kwenu wiki ijayo. Tuendelee kufanya kazi pamoja kwa moyo uleule."),
	)
	addVocab(db, l,
		vw("ufasaha", "fluency, eloquence", "Anaongea kwa ufasaha.", "She speaks fluently.", finch),
		vw("kufariki", "to pass away (respectful)", "Mzee alifariki jana.", "The old man passed away yesterday.", "Nana"),
		vw("hongera", "congratulations", "Hongera kwa kufaulu!", "Congratulations on passing!", "Lumora"),
		vw("kukamilika", "to be completed", "Kazi imekamilika.", "The work is complete.", "Mira"),
		vw("kuchangia", "to contribute", "Asante kwa kuchangia.", "Thank you for contributing.", "Mira"),
	)

	s = addSkillL(db, sw, u, "Speeches & Ceremony (Hotuba)", "Weddings, funerals, national days — the most formal Swahili.", "Mic", "#17A3DD", 40, 1120)
	l = addLesson(db, s, "Hotuba Rasmi", 1, 30,
		char(finch, "Formal speeches open by honouring guests in order of rank: Mheshimiwa Mgeni Rasmi, waheshimiwa viongozi, mabibi na mabwana… (Honourable Guest of Honour, honourable leaders, ladies and gentlemen). Weddings: Tunawatakia maisha marefu yenye furaha. Funerals: Pole kwa msiba; Mungu aiweke roho yake mahali pema peponi. Closing: Asanteni kwa kunisikiliza."),
		mc("How do you address the guest of honour?", "Opening of a speech", "Mheshimiwa Mgeni Rasmi", "Mheshimiwa Mgeni Rasmi", "Mambo mgeni", "Rafiki mgeni", "Bwana mdogo"),
		mc("Condolence at a funeral", "Msiba", "Pole kwa msiba.", "Pole kwa msiba.", "Hongera sana!", "Safari njema!", "Karibu tena!"),
		mc("How do you close a speech?", "Final line", "Asanteni kwa kunisikiliza.", "Asanteni kwa kunisikiliza.", "Kwaheri, nimechoka.", "Poa, tuonane.", "Basi, ndiyo hiyo."),
		mc("A wedding blessing", "Harusi", "Tunawatakia maisha marefu yenye furaha.", "Tunawatakia maisha marefu yenye furaha.", "Pole kwa msiba.", "Upone haraka.", "Tutaonana kesho."),
		mc("What does 'mabibi na mabwana' mean?", "Hotuba", "ladies and gentlemen", "ladies and gentlemen", "grandparents", "brides and grooms", "teachers and pupils"),
		write("Write the opening of a formal speech in Swahili (about 100 words) for your school's prize-giving day: greet the guests in order of rank, thank the organisers, and congratulate the winners.",
			"Mheshimiwa Mgeni Rasmi, Mwalimu Mkuu, waheshimiwa walimu, wazazi, mabibi na mabwana, hamjambo! Ni heshima kubwa kwangu kusimama mbele yenu siku hii ya kutoa zawadi. Kwanza kabisa, napenda kuwashukuru waandaaji wote kwa maandalizi mazuri. Pili, nawapongeza wanafunzi wote walioshinda zawadi leo; juhudi zenu zimezaa matunda. Kumbukeni kwamba bandu bandu huisha gogo. Asanteni kwa kunisikiliza."),
	)
	addVocab(db, l,
		vw("hotuba", "speech", "Rais alitoa hotuba.", "The president gave a speech.", finch),
		vw("Mgeni Rasmi", "guest of honour", "Mgeni rasmi amewasili.", "The guest of honour has arrived.", finch),
		vw("mabibi na mabwana", "ladies and gentlemen", "Mabibi na mabwana, karibuni.", "Ladies and gentlemen, welcome.", "Mira"),
		vw("msiba", "bereavement, mourning", "Tuna msiba nyumbani.", "We're in mourning at home.", "Nana"),
		vw("kupongeza", "to congratulate", "Nawapongeza washindi.", "I congratulate the winners.", "Lumora"),
	)
}
