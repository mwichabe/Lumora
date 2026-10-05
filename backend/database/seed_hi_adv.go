package database

import "gorm.io/gorm"

// Hindi course, B1 → C2. See seed_hi.go. From here on, exercises are in
// Devanagari only; vocabulary keeps its transliteration as a pronunciation aid.

// ───────────────────── B1 — बातचीत: Conversation & Narrative ─────────────────────

func seedHindiB1(db *gorm.DB) {
	const hi = "hi"
	const u = "B1 · बातचीत — Conversation & Narrative"
	finch := "Professor Finch"

	s := addSkillL(db, hi, u, "Compound Verbs", "खा लेना, दे देना, चला जाना — the colour of natural Hindi.", "Layers", "#6C3FC5", 23, 300)
	l := addLesson(db, s, "लेना, देना, जाना…", 1, 22,
		char(finch, "Natives pair a main verb's stem with a 'light' verb that adds nuance. लेना: for oneself — खाना खा लो (go ahead, eat). देना: for others — मैंने उसे बता दिया (I told him). जाना: completion or change — वह सो गया (he fell asleep), चले जाओ (go away). पड़ना: suddenness — वह रो पड़ी (she burst into tears). बैठना: rash action — क्या कर बैठे! (what have you gone and done!)."),
		mc("He fell asleep.", "सोना + जाना", "वह सो गया।", "वह सो गया।", "वह सोया दिया।", "वह सो लिया।", "वह सोता गया।"),
		mc("I told him (for his benefit).", "बताना + देना", "मैंने उसे बता दिया।", "मैंने उसे बता दिया।", "मैंने उसे बता लिया।", "मैं उसे बता गया।", "मैंने उसे बताया गया।"),
		mc("What nuance does पड़ना add in 'वह रो पड़ी'?", "रोना + पड़ना", "Suddenness — she burst into tears", "Suddenness — she burst into tears", "Completion — she finished crying", "For herself", "Ability — she could cry"),
		mc("Go ahead and eat! (to a friend)", "खाना + लेना", "खाना खा लो!", "खाना खा लो!", "खाना खा दो!", "खाना खा जाओ बैठो!", "खाना खाते लो!"),
		mc("Which light verb usually shows the action benefits someone else?", "Compound verbs", "देना", "देना", "लेना", "जाना", "पड़ना"),
		speak("Mira", "चाय ठंडी हो जाएगी, जल्दी पी लो!"),
	)
	addVocab(db, l,
		hv("खा लेना", "khā lenā", "to eat up (for oneself)", "पहले खा लो।", "pahle khā lo.", "Eat first.", "Cora"),
		hv("बता देना", "batā denā", "to tell (someone)", "मुझे बता देना।", "mujhe batā denā.", "Do let me know.", "Mira"),
		hv("सो जाना", "so jānā", "to fall asleep", "बच्चा सो गया।", "bacchā so gayā.", "The child fell asleep.", "Nana"),
		hv("रो पड़ना", "ro paṛnā", "to burst into tears", "वह अचानक रो पड़ी।", "vah acānak ro paṛī.", "She suddenly burst into tears.", "Nana"),
	)

	s = addSkillL(db, hi, u, "Can, Already, Must", "सकना, चुकना, पड़ना / चाहिए / है.", "Shield", "#00C2A8", 24, 320)
	l = addLesson(db, s, "Ability & Obligation", 1, 22,
		char(finch, "सकना (can) follows the stem: मैं तैर सकता हूँ (I can swim), वह नहीं आ सकती (she can't come). चुकना (already): मैं खा चुका हूँ (I've already eaten). Obligation uses मुझे + infinitive: मुझे जाना है (I have to go — planned), मुझे जाना पड़ेगा (I'll have to — forced), आपको आराम करना चाहिए (you should rest)."),
		mc("I can swim. (man)", "सकना", "मैं तैर सकता हूँ।", "मैं तैर सकता हूँ।", "मैं तैरना सकता हूँ।", "मुझे तैर सकता है।", "मैं तैरता सकता।"),
		mc("I have already eaten. (woman)", "चुकना", "मैं खा चुकी हूँ।", "मैं खा चुकी हूँ।", "मैं खा चुका हूँ।", "मैंने खा चुकी।", "मुझे खा चुकी है।"),
		mc("You should rest. (polite)", "चाहिए", "आपको आराम करना चाहिए।", "आपको आराम करना चाहिए।", "आप आराम करना चाहिए।", "आपको आराम चाहिए करना।", "आपने आराम करना चाहिए।"),
		mc("I'll have to go to the office tomorrow.", "पड़ना (forced)", "मुझे कल दफ़्तर जाना पड़ेगा।", "मुझे कल दफ़्तर जाना पड़ेगा।", "मैं कल दफ़्तर जाना पड़ेगा।", "मुझे कल दफ़्तर जा पड़ेगा।", "मेरा कल दफ़्तर जाना पड़ा।"),
		mc("She can't come.", "सकना, negative", "वह नहीं आ सकती।", "वह नहीं आ सकती।", "वह नहीं आ सकता।", "वह आना नहीं सकती।", "उसे नहीं आ सकती।"),
		speak("Blaze", "माफ़ कीजिए, मुझे अभी जाना है, मैं कल फिर आ सकता हूँ।"),
	)
	addVocab(db, l,
		hv("सकना", "saknā", "can, be able", "क्या आप मेरी मदद कर सकते हैं?", "kyā āp merī madad kar sakte haĩ?", "Can you help me?", "Lumora"),
		hv("चुकना", "cuknā", "to have already (done)", "फ़िल्म शुरू हो चुकी है।", "film shurū ho cukī hai.", "The film has already started.", "Blaze"),
		hv("चाहिए", "chāhie", "should, ought", "तुम्हें सच बोलना चाहिए।", "tumhẽ sac bolnā chāhie.", "You should tell the truth.", finch),
		hv("आराम", "ārām", "rest (m)", "थोड़ा आराम कीजिए।", "thoṛā ārām kījie.", "Rest a little.", "Nana"),
	)

	s = addSkillL(db, hi, u, "Subjunctive & Conditionals", "शायद, अगर… तो… — possibility and 'if'.", "Scale", "#F5A623", 25, 340)
	l = addLesson(db, s, "अगर… तो…", 1, 22,
		char(finch, "The subjunctive (future without गा/गी/गे) expresses possibility, wishes and suggestions: शायद वह आए (maybe he'll come), चलें? (shall we go?), मैं चाहता हूँ कि तुम आओ. Conditions: अगर… तो… — अगर बारिश होगी तो हम नहीं जाएँगे (if it rains, we won't go). Unreal past: अगर मैं जानता तो बताता (if I had known, I'd have told)."),
		mc("Maybe she'll come.", "शायद + subjunctive", "शायद वह आए।", "शायद वह आए।", "शायद वह आती है।", "शायद वह आई थी।", "शायद वह आएगा।"),
		mc("If it rains, we won't go.", "अगर… तो", "अगर बारिश होगी तो हम नहीं जाएँगे।", "अगर बारिश होगी तो हम नहीं जाएँगे।", "अगर बारिश तो होगी हम नहीं जाएँगे।", "बारिश अगर हम नहीं जाएँगे।", "अगर बारिश होगी हम जाएँगे तो।"),
		mc("If I had known, I would have told you. (man)", "Unreal past", "अगर मैं जानता तो बताता।", "अगर मैं जानता तो बताता।", "अगर मैं जानूँगा तो बताऊँगा।", "अगर मैं जानता हूँ तो बताता हूँ।", "मैं जानता तो अगर।"),
		mc("I want you to come.", "चाहना कि + subjunctive", "मैं चाहता हूँ कि तुम आओ।", "मैं चाहता हूँ कि तुम आओ।", "मैं चाहता हूँ कि तुम आते हो।", "मैं चाहता तुम आना।", "मैं चाहता हूँ तुम्हें आना।"),
		mc("Shall we go? (suggestion)", "Subjunctive", "चलें?", "चलें?", "चलेंगे था?", "चलता?", "चलो गया?"),
		speak("Zephyr", "अगर मेरे पास समय होता, तो मैं पूरे भारत की यात्रा करता।"),
	)
	addVocab(db, l,
		hv("शायद", "shāyad", "perhaps, maybe", "शायद कल बारिश हो।", "shāyad kal bārish ho.", "Maybe it'll rain tomorrow.", "Zephyr"),
		hv("अगर… तो…", "agar… to…", "if… then…", "अगर तुम आओ तो अच्छा हो।", "agar tum āo to acchā ho.", "It'd be nice if you came.", finch),
		hv("चाहना", "chāhnā", "to want, wish", "मैं चाहती हूँ कि सब ख़ुश रहें।", "mãi chāhtī hū̃ ki sab khush rahẽ.", "I wish everyone stays happy.", "Mira"),
		hv("यात्रा", "yātrā", "journey (f)", "आपकी यात्रा कैसी रही?", "āpkī yātrā kaisī rahī?", "How was your journey?", "Riko"),
	)

	s = addSkillL(db, hi, u, "Relative Clauses", "जो… वह, जहाँ… वहाँ, जब… तब.", "Link", "#FF5C5C", 26, 360)
	l = addLesson(db, s, "जो… वह…", 1, 22,
		char(finch, "Hindi relative clauses come in pairs: जो… वह (the one who… he/that), जहाँ… वहाँ (where… there), जब… तब (when… then), जैसा… वैसा (as… so). जो लड़का वहाँ खड़ा है, वह मेरा भाई है — the boy who is standing there is my brother. जो oblique is जिस: जिस घर में मैं रहता हूँ…"),
		mc("The girl who is singing is my sister.", "जो… वह", "जो लड़की गा रही है, वह मेरी बहन है।", "जो लड़की गा रही है, वह मेरी बहन है।", "लड़की जो गा रही, मेरी बहन वह।", "वह लड़की जो, गा रही मेरी बहन है।", "जो लड़की गा रही है, जो मेरी बहन है।"),
		mc("Wherever you go, I'll go too.", "जहाँ… वहाँ", "जहाँ तुम जाओगे, वहाँ मैं भी जाऊँगा।", "जहाँ तुम जाओगे, वहाँ मैं भी जाऊँगा।", "जहाँ तुम जाओगे, जहाँ मैं जाऊँगा।", "वहाँ तुम जाओगे, जहाँ मैं जाऊँगा।", "जब तुम जाओगे, वहाँ मैं।"),
		mc("When I was small, then I lived in a village.", "जब… तब", "जब मैं छोटा था, तब मैं गाँव में रहता था।", "जब मैं छोटा था, तब मैं गाँव में रहता था।", "जब मैं छोटा था, जब मैं गाँव में रहता था।", "तब मैं छोटा था, जब गाँव।", "जो मैं छोटा था, वह गाँव।"),
		mc("the house in which I live", "जो → जिस (oblique)", "जिस घर में मैं रहता हूँ", "जिस घर में मैं रहता हूँ", "जो घर में मैं रहता हूँ", "घर जिसमें मैं रहता", "वह घर में जो रहता हूँ"),
		mc("As you sow, so shall you reap.", "जैसा… वैसा", "जैसा बोओगे, वैसा काटोगे।", "जैसा बोओगे, वैसा काटोगे।", "जो बोओगे, जहाँ काटोगे।", "जब बोओगे, वैसा काटा।", "जैसा बोया, जो काटेगा।"),
		speak("Mira", "जिस किताब की मैंने बात की थी, वह यही है।"),
	)
	addVocab(db, l,
		hv("जो… वह", "jo… vah", "the one who… (he/she/that)", "जो मेहनत करता है, वह सफल होता है।", "jo mehnat kartā hai, vah saphal hotā hai.", "Whoever works hard succeeds.", finch),
		hv("जहाँ… वहाँ", "jahā̃… vahā̃", "where… there", "जहाँ चाह, वहाँ राह।", "jahā̃ cāh, vahā̃ rāh.", "Where there's a will, there's a way.", "Nana"),
		hv("जब… तब", "jab… tab", "when… then", "जब मैं आऊँगा, तब बात करेंगे।", "jab mãi āū̃gā, tab bāt karẽge.", "We'll talk when I come.", "Riko"),
		hv("जिस", "jis", "which (oblique of जो)", "जिस दिन तुम आए…", "jis din tum āe…", "The day you came…", "Mira"),
	)

	s = addSkillL(db, hi, u, "Honorifics in Practice", "जी, आप + -इए, and respectful plurals.", "GraduationCap", "#17A3DD", 27, 380)
	l = addLesson(db, s, "आदर — Respect", 1, 22,
		char(finch, "Respect shows in grammar. A single respected person takes plural verbs: पिता जी आ रहे हैं (Father is coming), दादी जी सो रही हैं. Add जी to names, titles and even answers (जी हाँ, जी नहीं). Use उनका / उन्होंने for respected third persons, not उसका / उसने. With आप, the most polite requests use -इएगा: बैठिएगा (do please sit)."),
		mc("My grandmother is sleeping. (respectful)", "Respectful plural", "दादी जी सो रही हैं।", "दादी जी सो रही हैं।", "दादी जी सो रही है।", "दादी सो रहा है।", "दादी जी सोती है।"),
		mc("The teacher said… (respectful)", "उन्होंने vs उसने", "अध्यापक जी ने कहा…", "अध्यापक जी ने कहा…", "अध्यापक ने कहती…", "अध्यापक जी कहा था ने…", "अध्यापक को कहा…"),
		mc("Respectful 'his/her' for an elder", "उनका vs उसका", "उनका", "उनका", "उसका", "तेरा", "इसका"),
		mc("Most polite: 'please do come again'", "-इएगा", "फिर आइएगा।", "फिर आइएगा।", "फिर आना।", "फिर आ।", "फिर आओ।"),
		mc("A polite 'yes'", "जी", "जी हाँ", "जी हाँ", "हाँ रे", "अरे हाँ", "हाँ तू"),
		speak("Lumora", "जी हाँ, गुप्ता जी अभी दफ़्तर में हैं। आप बैठिएगा।"),
	)
	addVocab(db, l,
		hv("जी", "jī", "honorific particle", "शर्मा जी, नमस्ते।", "sharmā jī, namaste.", "Hello, Mr Sharma.", "Lumora"),
		hv("उन्होंने", "unhõne", "he/she/they (respectful, with ने)", "उन्होंने मेरी मदद की।", "unhõne merī madad kī.", "He kindly helped me.", finch),
		hv("आदर", "ādar", "respect (m)", "बड़ों का आदर करो।", "baṛõ kā ādar karo.", "Respect your elders.", "Nana"),
		hv("आइएगा", "āiegā", "do please come (very polite)", "शाम को ज़रूर आइएगा।", "shām ko zarūr āiegā.", "Do come in the evening.", "Lumora"),
	)

	s = addSkillL(db, hi, u, "Health & the Body", "मेरे सिर में दर्द है, बुख़ार, दवाई.", "HeartPulse", "#9B51E0", 28, 400)
	l = addLesson(db, s, "डॉक्टर के पास", 1, 22,
		char("Nana", "Pain is 'in' the body part: मेरे सिर में दर्द है (I have a headache), पेट में दर्द (stomach ache). Illness takes मुझे: मुझे बुख़ार है, मुझे ज़ुकाम है (I have a cold), मुझे खाँसी है (cough). The doctor: यह दवाई दिन में दो बार लीजिए (take this medicine twice a day). Wish well: जल्दी ठीक हो जाइए!"),
		mc("I have a headache.", "सिर", "मेरे सिर में दर्द है।", "मेरे सिर में दर्द है।", "मेरा सिर दर्द हूँ।", "मुझे सिर है दर्द।", "मैं सिर में दर्द हूँ।"),
		mc("I have a fever.", "बुख़ार", "मुझे बुख़ार है।", "मुझे बुख़ार है।", "मैं बुख़ार हूँ।", "मेरा बुख़ार हूँ।", "बुख़ार मैं है।"),
		mc("Take this medicine twice a day. (polite)", "Doctor", "यह दवाई दिन में दो बार लीजिए।", "यह दवाई दिन में दो बार लीजिए।", "यह दवाई दो दिन में बार लो।", "दवाई दिन दो बार दीजिए।", "यह दवाई दिन में दो बार दो मत।"),
		mc("Get well soon! (polite)", "Wishing well", "जल्दी ठीक हो जाइए!", "जल्दी ठीक हो जाइए!", "जल्दी बीमार हो जाइए!", "धीरे ठीक मत होइए!", "जल्दी आइए!"),
		mc("What is ज़ुकाम?", "मुझे ज़ुकाम है।", "A cold", "A cold", "A fever", "A broken bone", "A stomach ache"),
		speak("Nana", "आपको आराम करना चाहिए और ख़ूब पानी पीना चाहिए।"),
	)
	addVocab(db, l,
		hv("दर्द", "dard", "pain (m)", "मेरे पेट में दर्द है।", "mere peṭ mẽ dard hai.", "I have a stomach ache.", "Nana"),
		hv("बुख़ार", "bukhār", "fever (m)", "बच्चे को बुख़ार है।", "bacce ko bukhār hai.", "The child has a fever.", "Nana"),
		hv("दवाई", "davāī", "medicine (f)", "दवाई खाने के बाद लीजिए।", "davāī khāne ke bād lījie.", "Take the medicine after food.", "Nana"),
		hv("डॉक्टर", "ḍākṭar", "doctor", "डॉक्टर साहब आ गए।", "ḍākṭar sāhab ā gae.", "The doctor has come.", "Riko"),
		hv("बीमार", "bīmār", "ill, sick", "वह दो दिन से बीमार है।", "vah do din se bīmār hai.", "She's been ill for two days.", "Mira"),
	)

	s = addSkillL(db, hi, u, "Opinions & Work", "मेरे ख़याल से, मेरी राय में — दफ़्तर और नौकरी.", "Briefcase", "#00C2A8", 29, 420)
	l = addLesson(db, s, "दफ़्तर में", 1, 22,
		char("Mira", "Opinions: मेरे ख़याल से / मेरे विचार से (I think), मेरी राय में (in my opinion), मैं सहमत हूँ / असहमत हूँ (I agree / disagree), आप ठीक कह रहे हैं, लेकिन… (you're right, but…). Work: दफ़्तर (office), नौकरी (job), तनख़्वाह / वेतन (salary), छुट्टी (leave), बैठक / मीटिंग (meeting)."),
		mc("In my opinion…", "Opinion", "मेरी राय में…", "मेरी राय में…", "मेरा राय में…", "मुझे राय है…", "राय मेरे से…"),
		mc("I disagree.", "Disagreeing", "मैं असहमत हूँ।", "मैं असहमत हूँ।", "मैं सहमत हूँ।", "मुझे असहमत है।", "मैं असहमत नहीं हूँ।"),
		mc("What is तनख़्वाह?", "तनख़्वाह कब मिलेगी?", "Salary", "Salary", "Holiday", "Meeting", "Office"),
		mc("You're right, but…", "Polite disagreement", "आप ठीक कह रहे हैं, लेकिन…", "आप ठीक कह रहे हैं, लेकिन…", "आप ग़लत हैं, बस!", "चुप रहिए…", "ठीक है, ठीक है…"),
		mc("I need leave tomorrow.", "छुट्टी", "मुझे कल छुट्टी चाहिए।", "मुझे कल छुट्टी चाहिए।", "मैं कल छुट्टी चाहिए।", "कल छुट्टी मेरी है।", "मुझे कल छुट्टी हूँ।"),
		write("Write a short paragraph in Hindi (about 60 words, in Devanagari) giving your opinion on working from home. Use 'मेरी राय में', 'लेकिन' and 'क्योंकि'.",
			"मेरी राय में घर से काम करने के कई फ़ायदे हैं। रोज़ ट्रैफ़िक में समय बर्बाद नहीं होता और हम अपने परिवार के साथ ज़्यादा समय बिता सकते हैं। लेकिन कुछ मुश्किलें भी हैं, क्योंकि साथियों से मिलना कम हो जाता है और कभी-कभी अकेलापन महसूस होता है। इसलिए मुझे लगता है कि हफ़्ते में दो-तीन दिन घर से और बाक़ी दिन दफ़्तर से काम करना सबसे अच्छा है।"),
	)
	addVocab(db, l,
		hv("राय", "rāy", "opinion (f)", "आपकी क्या राय है?", "āpkī kyā rāy hai?", "What's your opinion?", "Mira"),
		hv("सहमत", "sahmat", "in agreement", "मैं आपसे सहमत हूँ।", "mãi āpse sahmat hū̃.", "I agree with you.", finch),
		hv("नौकरी", "naukrī", "job (f)", "उसे नई नौकरी मिली।", "use naī naukrī milī.", "She got a new job.", "Blaze"),
		hv("दफ़्तर", "daftar", "office (m)", "मैं नौ बजे दफ़्तर जाता हूँ।", "mãi nau baje daftar jātā hū̃.", "I go to the office at nine.", "Mira"),
		hv("छुट्टी", "chuṭṭī", "holiday, leave (f)", "कल छुट्टी है।", "kal chuṭṭī hai.", "Tomorrow is a holiday.", "Pip"),
	)
}

// ───────────────────── B2 — औपचारिक हिंदी: Formal Hindi & Society ─────────────────────

func seedHindiB2(db *gorm.DB) {
	const hi = "hi"
	const u = "B2 · औपचारिक हिंदी — Formal Hindi & Society"
	finch := "Professor Finch"

	s := addSkillL(db, hi, u, "The Passive", "किया जाता है, बनाया गया — and द्वारा.", "Shield", "#6C3FC5", 30, 460)
	l := addLesson(db, s, "कर्मवाच्य", 1, 24,
		char(finch, "Passive = perfective participle + जाना: यहाँ हिंदी बोली जाती है (Hindi is spoken here), ताजमहल सत्रहवीं सदी में बनाया गया (the Taj Mahal was built in the 17th century). The agent takes द्वारा (formal) or से: यह उपन्यास प्रेमचंद द्वारा लिखा गया. Negative passive often means 'can't': मुझसे खाया नहीं जाता (I can't bring myself to eat)."),
		mc("Hindi is spoken here.", "Passive", "यहाँ हिंदी बोली जाती है।", "यहाँ हिंदी बोली जाती है।", "यहाँ हिंदी बोलता है।", "यहाँ हिंदी बोला जाता है।", "यहाँ हिंदी बोल जाती है।"),
		mc("This novel was written by Premchand.", "द्वारा", "यह उपन्यास प्रेमचंद द्वारा लिखा गया।", "यह उपन्यास प्रेमचंद द्वारा लिखा गया।", "यह उपन्यास प्रेमचंद ने लिखी गई।", "प्रेमचंद द्वारा उपन्यास लिखता है।", "यह उपन्यास प्रेमचंद को लिखा।"),
		mc("The letter was sent yesterday.", "चिट्ठी (f)", "चिट्ठी कल भेजी गई।", "चिट्ठी कल भेजी गई।", "चिट्ठी कल भेजा गया।", "चिट्ठी कल भेजती गई।", "चिट्ठी कल भेजी जाती है।"),
		mc("What does 'मुझसे चला नहीं जाता' mean?", "Negative passive", "I can't (bring myself to) walk", "I can't (bring myself to) walk", "I don't go anywhere", "Nobody walked with me", "I was not allowed to go"),
		mc("When was the Taj Mahal built?", "ताजमहल ___ बनाया गया।", "सत्रहवीं सदी में", "सत्रहवीं सदी में", "पिछले साल", "दसवीं सदी में", "उन्नीसवीं सदी में"),
		speak("Zephyr", "इस मंदिर का निर्माण सन् 1950 में किया गया था।"),
	)
	addVocab(db, l,
		hv("द्वारा", "dvārā", "by (agent, formal)", "सरकार द्वारा घोषित।", "sarkār dvārā ghoṣit.", "Announced by the government.", finch),
		hv("बनाया गया", "banāyā gayā", "was made / built", "यह पुल 2010 में बनाया गया।", "yah pul 2010 mẽ banāyā gayā.", "This bridge was built in 2010.", "Zephyr"),
		hv("उपन्यास", "upanyās", "novel (m)", "मैंने एक उपन्यास पढ़ा।", "mãine ek upanyās paṛhā.", "I read a novel.", finch),
		hv("सदी", "sadī", "century (f)", "इक्कीसवीं सदी।", "ikkīsvī̃ sadī.", "The twenty-first century.", "Zephyr"),
	)

	s = addSkillL(db, hi, u, "Causatives", "करना → कराना / करवाना, पढ़ना → पढ़ाना.", "Users", "#00C2A8", 31, 490)
	l = addLesson(db, s, "Making Things Happen", 1, 24,
		char(finch, "Causatives are formed from the root. First causative (-आ): पढ़ना → पढ़ाना (to teach — make read), सीखना → सिखाना (to teach — make learn), हँसना → हँसाना (to make laugh), देखना → दिखाना (to show). Second causative (-वा, through someone else): बनाना → बनवाना (have something built), करना → करवाना: मैंने घर रंगवाया (I had the house painted)."),
		mc("The teacher teaches the children.", "पढ़ना → पढ़ाना", "अध्यापक बच्चों को पढ़ाते हैं।", "अध्यापक बच्चों को पढ़ाते हैं।", "अध्यापक बच्चों को पढ़ते हैं।", "अध्यापक बच्चों से पढ़वाता हैं।", "अध्यापक बच्चे पढ़ते हैं।"),
		mc("Show me the photo.", "देखना → दिखाना", "मुझे फ़ोटो दिखाइए।", "मुझे फ़ोटो दिखाइए।", "मुझे फ़ोटो देखिए।", "मुझे फ़ोटो दिखवाइए मत।", "फ़ोटो मुझको देखो।"),
		mc("I had the house painted.", "Second causative", "मैंने घर रंगवाया।", "मैंने घर रंगवाया।", "मैंने घर रंगा।", "मैं घर रंगवाया।", "मैंने घर रंगाया गया।"),
		mc("What does हँसाना mean?", "हँसना → हँसाना", "to make (someone) laugh", "to make (someone) laugh", "to laugh", "to be laughed at", "to stop laughing"),
		mc("Shah Jahan had the Taj Mahal built.", "बनवाना", "शाहजहाँ ने ताजमहल बनवाया।", "शाहजहाँ ने ताजमहल बनवाया।", "शाहजहाँ ने ताजमहल बनाया गया।", "शाहजहाँ ताजमहल बना।", "शाहजहाँ ने ताजमहल बना।"),
		speak("Mira", "मेरी दादी ने मुझे खाना बनाना सिखाया।"),
	)
	addVocab(db, l,
		hv("पढ़ाना", "paṛhānā", "to teach", "वह गणित पढ़ाती हैं।", "vah gaṇit paṛhātī haĩ.", "She teaches maths.", finch),
		hv("सिखाना", "sikhānā", "to teach (a skill)", "मुझे गाड़ी चलाना सिखाओ।", "mujhe gāṛī calānā sikhāo.", "Teach me to drive.", "Blaze"),
		hv("दिखाना", "dikhānā", "to show", "मुझे रास्ता दिखाइए।", "mujhe rāstā dikhāie.", "Please show me the way.", "Riko"),
		hv("बनवाना", "banvānā", "to have (something) made", "मैंने नया सूट बनवाया।", "mãine nayā sūṭ banvāyā.", "I had a new suit made.", "Zephyr"),
	)

	s = addSkillL(db, hi, u, "Participles & Linking Actions", "खाकर, पहुँचते ही, गाते हुए, बैठा हुआ.", "Link", "#F5A623", 32, 520)
	l = addLesson(db, s, "Having Done, While Doing", 1, 24,
		char(finch, "Participles make Hindi flow. Stem + कर = having done: खाना खाकर मैं सो गया (after eating, I slept). -ते ही = as soon as: घर पहुँचते ही बारिश शुरू हो गई. -ते हुए = while: वह गाते हुए खाना बना रही थी. -आ हुआ = in a state: बैठा हुआ आदमी (the seated man), खुली हुई खिड़की (the open window)."),
		mc("After eating, I slept.", "-कर", "खाना खाकर मैं सो गया।", "खाना खाकर मैं सो गया।", "खाना खाते मैं सो गया।", "खाना खाया कर सो गया।", "खाना खाकर सोता गया था।"),
		mc("As soon as I reached home, it started raining.", "-ते ही", "घर पहुँचते ही बारिश शुरू हो गई।", "घर पहुँचते ही बारिश शुरू हो गई।", "घर पहुँचकर ही बारिश।", "घर पहुँचा हुआ बारिश।", "घर पहुँचते हुए बारिश शुरू।"),
		mc("She was cooking while singing.", "-ते हुए", "वह गाते हुए खाना बना रही थी।", "वह गाते हुए खाना बना रही थी।", "वह गाकर हुए खाना बनाती।", "वह गाती ही खाना बना।", "वह गाया हुआ खाना बना रही।"),
		mc("the open window (खिड़की, f)", "-आ हुआ", "खुली हुई खिड़की", "खुली हुई खिड़की", "खुला हुआ खिड़की", "खुलते हुए खिड़की", "खोलकर खिड़की"),
		mc("What does 'सोच-समझकर' mean?", "सोच-समझकर बोलिए।", "Having thought it through, carefully", "Having thought it through, carefully", "Without thinking", "While sleeping", "As soon as you think"),
		speak("Riko", "स्टेशन पहुँचते ही मुझे फ़ोन कीजिए।"),
	)
	addVocab(db, l,
		hv("-कर", "-kar", "having done", "हाथ धोकर खाना खाओ।", "hāth dhokar khānā khāo.", "Wash your hands and then eat.", "Nana"),
		hv("-ते ही", "-te hī", "as soon as", "सुबह उठते ही चाय पीता हूँ।", "subah uṭhte hī chāy pītā hū̃.", "I drink tea as soon as I get up.", "Cora"),
		hv("-ते हुए", "-te hue", "while doing", "वह मुस्कुराते हुए बोली।", "vah muskurāte hue bolī.", "She said, smiling.", "Mira"),
		hv("खिड़की", "khiṛkī", "window (f)", "खिड़की खोल दीजिए।", "khiṛkī khol dījie.", "Please open the window.", "Riko"),
	)

	s = addSkillL(db, hi, u, "Formal vs Everyday Vocabulary", "शुद्ध हिंदी and everyday Hindustani: प्रश्न / सवाल, सहायता / मदद.", "Languages", "#FF5C5C", 33, 550)
	l = addLesson(db, s, "दो रजिस्टर — Two Registers", 1, 24,
		char(finch, "Hindi has two layers. Formal Hindi (news, government, writing) favours Sanskrit-derived words; everyday speech mixes in Persian- and Arabic-derived words shared with Urdu. धन्यवाद / शुक्रिया (thanks), प्रश्न / सवाल (question), उत्तर / जवाब (answer), सहायता / मदद (help), समय / वक़्त (time), जीवन / ज़िंदगी (life), प्रयास / कोशिश (effort). Knowing both lets you switch register."),
		mc("Formal word for 'question'", "Sanskrit-derived", "प्रश्न", "प्रश्न", "सवाल", "जवाब", "वक़्त"),
		mc("Everyday word for 'help'", "Persian/Arabic-derived", "मदद", "मदद", "सहायता", "प्रयास", "उत्तर"),
		mc("Formal equivalent of ज़िंदगी", "life", "जीवन", "जीवन", "समय", "वक़्त", "कोशिश"),
		mc("Which word would a news anchor more likely use for 'effort'?", "Formal register", "प्रयास", "प्रयास", "कोशिश", "जवाब", "मदद"),
		mc("Everyday equivalent of समय", "time", "वक़्त", "वक़्त", "जीवन", "उत्तर", "सहायता"),
		match("Match the pair", "उत्तर", "जवाब (answer)", "जवाब (answer)", "सवाल (question)", "मदद (help)", "शुक्रिया (thanks)"),
		speak("Zephyr", "आपकी सहायता के लिए हार्दिक धन्यवाद।"),
	)
	addVocab(db, l,
		hv("प्रश्न / सवाल", "prashn / savāl", "question (formal / everyday)", "क्या कोई प्रश्न है?", "kyā koī prashn hai?", "Are there any questions?", finch),
		hv("उत्तर / जवाब", "uttar / javāb", "answer (formal / everyday)", "उत्तर सही है।", "uttar sahī hai.", "The answer is correct.", finch),
		hv("सहायता / मदद", "sahāytā / madad", "help (formal / everyday)", "आपकी सहायता चाहिए।", "āpkī sahāytā chāhie.", "We need your help.", "Mira"),
		hv("जीवन / ज़िंदगी", "jīvan / zindagī", "life (formal / everyday)", "ज़िंदगी ख़ूबसूरत है।", "zindagī khūbsūrat hai.", "Life is beautiful.", "Lumora"),
		hv("प्रयास / कोशिश", "prayās / koshish", "effort (formal / everyday)", "कोशिश करते रहो।", "koshish karte raho.", "Keep trying.", "Blaze"),
	)

	s = addSkillL(db, hi, u, "News & Society", "सरकार, चुनाव, अर्थव्यवस्था, महँगाई.", "Globe", "#17A3DD", 34, 580)
	l = addLesson(db, s, "समाचार", 1, 24,
		char(finch, "News Hindi (आज तक, NDTV इंडिया, दैनिक जागरण) uses formal vocabulary: सरकार (government), प्रधानमंत्री (prime minister), संसद (parliament), चुनाव (election), अर्थव्यवस्था (economy), महँगाई (inflation, rising prices), बाढ़ (flood), सूखा (drought), नागरिक (citizen). Headlines use the passive and drop verbs: 'राजधानी में भारी बारिश, यातायात ठप'."),
		mc("What is चुनाव?", "आम चुनाव", "election", "election", "economy", "parliament", "flood"),
		mc("What is महँगाई?", "महँगाई बढ़ी", "inflation, rising prices", "inflation, rising prices", "salary", "tax", "unemployment"),
		mc("What is संसद?", "संसद का सत्र", "parliament", "parliament", "court", "army", "city council"),
		mc("What does the headline say?", "राजधानी में भारी बारिश, यातायात ठप", "Heavy rain in the capital, traffic at a standstill", "Heavy rain in the capital, traffic at a standstill", "Drought in the capital", "The capital celebrates a festival", "New roads open in the capital"),
		mc("What is बाढ़?", "असम में बाढ़", "flood", "flood", "drought", "storm", "earthquake"),
		speak("Zephyr", "सरकार ने किसानों के लिए नई योजना की घोषणा की है।"),
	)
	addVocab(db, l,
		hv("सरकार", "sarkār", "government (f)", "सरकार ने नियम बदले।", "sarkār ne niyam badle.", "The government changed the rules.", finch),
		hv("चुनाव", "cunāv", "election (m)", "चुनाव अगले महीने है।", "cunāv agle mahīne hai.", "The election is next month.", finch),
		hv("अर्थव्यवस्था", "arthvyavasthā", "economy (f)", "अर्थव्यवस्था बढ़ रही है।", "arthvyavasthā baṛh rahī hai.", "The economy is growing.", "Zephyr"),
		hv("महँगाई", "mahãgāī", "inflation, high prices (f)", "महँगाई से सब परेशान हैं।", "mahãgāī se sab pareshān haĩ.", "Everyone is troubled by rising prices.", "Cora"),
		hv("घोषणा", "ghoṣṇā", "announcement (f)", "नई घोषणा हुई।", "naī ghoṣṇā huī.", "A new announcement was made.", "Zephyr"),
	)

	s = addSkillL(db, hi, u, "Formal Letters", "सेवा में, महोदय, विषय:, सविनय निवेदन, भवदीय.", "PenLine", "#9B51E0", 35, 610)
	l = addLesson(db, s, "औपचारिक पत्र", 1, 24,
		char("Mira", "A formal Hindi letter opens with सेवा में (to), the recipient's title and address, then विषय: (subject). Greeting: महोदय / महोदया (Sir / Madam). The body often begins सविनय निवेदन है कि… (I respectfully submit that…). Close with भवदीय / भवदीया (yours faithfully) and your name. Students write प्रार्थी (applicant) or आपका आज्ञाकारी शिष्य."),
		mc("How does a formal letter's address line begin?", "Opening", "सेवा में,", "सेवा में,", "प्रिय दोस्त,", "नमस्ते यार,", "अरे सुनो,"),
		mc("What comes after 'विषय:'?", "विषय: छुट्टी हेतु प्रार्थना-पत्र", "The subject of the letter", "The subject of the letter", "The date", "The signature", "A poem"),
		mc("I respectfully submit that…", "Body opening", "सविनय निवेदन है कि…", "सविनय निवेदन है कि…", "मुझे बताओ कि…", "सुनो, बात यह है…", "चलो, शुरू करते हैं…"),
		mc("Yours faithfully (man)", "Closing", "भवदीय", "भवदीय", "तुम्हारा दोस्त", "प्यार सहित", "फिर मिलेंगे"),
		mc("Formal greeting to a woman", "Salutation", "महोदया", "महोदया", "प्रिय बहना", "दीदी", "मैडम जी यार"),
		write("Write a short formal letter in Hindi (about 80 words, in Devanagari) to your school principal asking for two days' leave because of illness. Use सेवा में, विषय:, महोदय, सविनय निवेदन है कि… and भवदीय/भवदीया.",
			"सेवा में,\nप्रधानाचार्य जी,\nसरस्वती विद्यालय, दिल्ली।\n\nविषय: दो दिन के अवकाश हेतु प्रार्थना-पत्र\n\nमहोदय,\nसविनय निवेदन है कि मुझे कल रात से तेज़ बुख़ार है। डॉक्टर ने मुझे दो दिन आराम करने की सलाह दी है। अतः आपसे प्रार्थना है कि मुझे दिनांक 10 और 11 मार्च का अवकाश प्रदान करने की कृपा करें।\n\nसधन्यवाद।\nभवदीया,\nआशा कुमारी\nकक्षा दस"),
	)
	addVocab(db, l,
		hv("सेवा में", "sevā mẽ", "to (formal letter address)", "सेवा में, प्रबंधक महोदय।", "sevā mẽ, prabandhak mahoday.", "To the Manager.", "Mira"),
		hv("महोदय / महोदया", "mahoday / mahodayā", "Sir / Madam", "आदरणीय महोदय,", "ādaraṇīy mahoday,", "Respected Sir,", finch),
		hv("विषय", "viṣay", "subject (m)", "विषय: नौकरी के लिए आवेदन", "viṣay: naukrī ke lie āvedan", "Subject: job application", "Mira"),
		hv("निवेदन", "nivedan", "humble request (m)", "मेरा निवेदन है कि…", "merā nivedan hai ki…", "My humble request is that…", finch),
		hv("भवदीय", "bhavdīy", "yours faithfully", "भवदीय, राकेश शर्मा", "bhavdīy, rākesh sharmā", "Yours faithfully, Rakesh Sharma", "Mira"),
	)
}

// ───────────────────── C1 — साहित्य और मुहावरे: Literature & Idioms ─────────────────────

func seedHindiC1(db *gorm.DB) {
	const hi = "hi"
	const u = "C1 · साहित्य और मुहावरे — Literature & Idioms"
	finch := "Professor Finch"

	s := addSkillL(db, hi, u, "Idioms (मुहावरे)", "नौ दो ग्यारह होना, दाल में कुछ काला होना…", "Quote", "#6C3FC5", 36, 700)
	l := addLesson(db, s, "मुहावरे", 1, 26,
		char("Nana", "मुहावरे are phrases whose meaning isn't literal, and Hindi is full of them. नौ दो ग्यारह होना (nine-two-eleven — to run away), दाल में कुछ काला होना (something black in the lentils — something fishy), आँखों का तारा (star of the eyes — the apple of one's eye), ऊँट के मुँह में जीरा (cumin in a camel's mouth — a drop in the ocean), आग बबूला होना (to be furious), नाक कटना (to lose face)."),
		mc("What does 'नौ दो ग्यारह होना' mean?", "चोर पुलिस को देखकर नौ दो ग्यारह हो गया।", "to run away", "to run away", "to do maths", "to be eleven years old", "to fall asleep"),
		mc("What does 'दाल में कुछ काला है' mean?", "Idiom", "Something's fishy", "Something's fishy", "The food is burnt", "The lentils are spoiled", "It's dinner time"),
		mc("Who is someone's 'आँखों का तारा'?", "वह अपनी माँ की आँखों का तारा है।", "A beloved, darling person", "A beloved, darling person", "An astronomer", "Someone with eye trouble", "A film star"),
		mc("'ऊँट के मुँह में जीरा' describes…", "Idiom", "Something far too little", "Something far too little", "A tasty dish", "A desert journey", "A great gift"),
		mc("What does 'आग बबूला होना' mean?", "पिता जी आग बबूला हो गए।", "to become furious", "to become furious", "to catch fire", "to cook quickly", "to feel cold"),
		mc("What does 'नाक कटना' mean?", "उसकी हरकत से परिवार की नाक कट गई।", "to lose face / be disgraced", "to lose face / be disgraced", "to have an accident", "to sneeze", "to smell something bad"),
		speak("Nana", "पुलिस को देखते ही चोर नौ दो ग्यारह हो गया।"),
	)
	addVocab(db, l,
		hv("मुहावरा", "muhāvrā", "idiom (m)", "इस मुहावरे का अर्थ बताइए।", "is muhāvre kā arth batāie.", "Tell me the meaning of this idiom.", finch),
		hv("नौ दो ग्यारह होना", "nau do gyārah honā", "to run away", "बच्चे नौ दो ग्यारह हो गए।", "bacce nau do gyārah ho gae.", "The kids ran off.", "Pip"),
		hv("दाल में कुछ काला", "dāl mẽ kuch kālā", "something fishy", "मुझे दाल में कुछ काला लगता है।", "mujhe dāl mẽ kuch kālā lagtā hai.", "I smell something fishy.", "Blaze"),
		hv("आँखों का तारा", "ā̃khõ kā tārā", "apple of one's eye", "बेटी पिता की आँखों का तारा है।", "beṭī pitā kī ā̃khõ kā tārā hai.", "The daughter is her father's darling.", "Nana"),
	)

	s = addSkillL(db, hi, u, "Proverbs (लोकोक्तियाँ)", "जैसी करनी वैसी भरनी, बूँद-बूँद से घड़ा भरता है…", "Quote", "#00C2A8", 37, 730)
	l = addLesson(db, s, "लोकोक्तियाँ", 1, 26,
		char("Nana", "Proverbs are complete sayings: जैसी करनी वैसी भरनी (as you sow, so you reap), बूँद-बूँद से घड़ा भरता है (drop by drop the pot fills), अब पछताए होत क्या जब चिड़िया चुग गई खेत (what's the use of regret once the birds have eaten the field — no use crying over spilt milk), नाच न जाने आँगन टेढ़ा (can't dance, blames the crooked courtyard), अधजल गगरी छलकत जाए (a half-full pot spills most — empty vessels make the most noise)."),
		mc("'जैसी करनी वैसी भरनी' means…", "Proverb", "As you sow, so shall you reap", "As you sow, so shall you reap", "Work hard to get rich", "Fill your pot every day", "Don't do housework"),
		mc("Which proverb praises small, steady savings?", "Saving little by little", "बूँद-बूँद से घड़ा भरता है", "बूँद-बूँद से घड़ा भरता है", "नाच न जाने आँगन टेढ़ा", "अधजल गगरी छलकत जाए", "जैसी करनी वैसी भरनी"),
		mc("'नाच न जाने आँगन टेढ़ा' mocks someone who…", "Proverb", "Blames circumstances for their own lack of skill", "Blames circumstances for their own lack of skill", "Dances very well", "Builds crooked houses", "Never goes outside"),
		mc("'अब पछताए होत क्या जब चिड़िया चुग गई खेत' advises…", "Proverb", "Act in time — regret afterwards is useless", "Act in time — regret afterwards is useless", "Feed the birds", "Plant more crops", "Never trust birds"),
		mc("'अधजल गगरी छलकत जाए' describes people who…", "Proverb", "Know little but show off the most", "Know little but show off the most", "Drink a lot of water", "Are always thirsty", "Work quietly"),
		speak("Nana", "बेटा, बूँद-बूँद से घड़ा भरता है, रोज़ थोड़ा पढ़ो।"),
	)
	addVocab(db, l,
		hv("लोकोक्ति", "lokokti", "proverb (f)", "दादी लोकोक्तियों में बात करती हैं।", "dādī lokoktiyõ mẽ bāt kartī haĩ.", "Grandma speaks in proverbs.", "Nana"),
		hv("जैसी करनी वैसी भरनी", "jaisī karnī vaisī bharnī", "as you sow, so you reap", "उसे सज़ा मिली — जैसी करनी वैसी भरनी।", "use sazā milī — jaisī karnī vaisī bharnī.", "He was punished — as you sow…", finch),
		hv("बूँद-बूँद से घड़ा भरता है", "bū̃d-bū̃d se ghaṛā bhartā hai", "every little helps", "रोज़ बचत करो — बूँद-बूँद से घड़ा भरता है।", "roz bacat karo — bū̃d-bū̃d se ghaṛā bhartā hai.", "Save daily — every little helps.", "Cora"),
		hv("पछताना", "pachtānā", "to regret", "बाद में पछताओगे।", "bād mẽ pachtāoge.", "You'll regret it later.", "Mira"),
	)

	s = addSkillL(db, hi, u, "Literature: Kabir to Premchand", "दोहे, कहानियाँ, उपन्यास — the classics of Hindi.", "BookOpen", "#F5A623", 38, 760)
	l = addLesson(db, s, "हिंदी साहित्य", 1, 26,
		char(finch, "Hindi literature spans centuries. Kabir (15th c.) wrote pithy दोहे (couplets) mocking hypocrisy: बुरा जो देखन मैं चला, बुरा न मिलिया कोय। जो दिल खोजा आपना, मुझसे बुरा न कोय॥ (I went looking for the wicked and found none; when I searched my own heart, none was worse than me). Tulsidas wrote the रामचरितमानस (1574). Premchand (1880–1936) founded modern Hindi fiction: the short story ईदगाह and the novel गोदान."),
		mc("What is Kabir's doha about?", "बुरा जो देखन मैं चला…", "Looking at one's own faults before others'", "Looking at one's own faults before others'", "Searching for a lost friend", "A journey to the city", "Love for one's mother"),
		mc("Who wrote the रामचरितमानस?", "1574", "Tulsidas", "Tulsidas", "Kabir", "Premchand", "Mirabai"),
		mc("Who is known as the father of modern Hindi fiction?", "गोदान, ईदगाह", "Premchand", "Premchand", "Kabir", "Tulsidas", "Surdas"),
		mc("In Premchand's 'ईदगाह', what does Hamid buy at the fair?", "ईदगाह", "Tongs (चिमटा) for his grandmother", "Tongs (चिमटा) for his grandmother", "Sweets for himself", "A toy horse", "A new shirt"),
		mc("What is a दोहा?", "Kabir's form", "A rhyming couplet", "A rhyming couplet", "A novel", "A play", "A long epic in prose"),
		mc("Why does Hamid buy that gift?", "ईदगाह", "His grandmother burns her hands making rotis", "His grandmother burns her hands making rotis", "It was the cheapest toy", "His friends told him to", "He wanted to sell it later"),
		speak("Zephyr", "बुरा जो देखन मैं चला, बुरा न मिलिया कोय।"),
	)
	addVocab(db, l,
		hv("साहित्य", "sāhitya", "literature (m)", "हिंदी साहित्य बहुत समृद्ध है।", "hindī sāhitya bahut samṛddh hai.", "Hindi literature is very rich.", finch),
		hv("दोहा", "dohā", "couplet (m)", "कबीर के दोहे प्रसिद्ध हैं।", "kabīr ke dohe prasiddh haĩ.", "Kabir's couplets are famous.", "Nana"),
		hv("कहानी", "kahānī", "story (f)", "प्रेमचंद की कहानियाँ यथार्थवादी हैं।", "premcand kī kahāniyā̃ yathārthvādī haĩ.", "Premchand's stories are realist.", finch),
		hv("लेखक / कवि", "lekhak / kavi", "writer / poet", "कबीर एक संत कवि थे।", "kabīr ek sant kavi the.", "Kabir was a saint-poet.", "Zephyr"),
	)

	s = addSkillL(db, hi, u, "Debate & Argument", "वाद-विवाद: पक्ष, विपक्ष, इसके अलावा, दूसरी ओर, अतः.", "MessageCircle", "#FF5C5C", 39, 790)
	l = addLesson(db, s, "वाद-विवाद", 1, 26,
		char(finch, "School and TV debates (वाद-विवाद) argue पक्ष (for) and विपक्ष (against). Openers: आदरणीय निर्णायक मंडल… (respected judges). Connectors: सबसे पहले (first of all), इसके अलावा (besides), दूसरी ओर (on the other hand), फिर भी (even so), इसलिए / अतः (therefore), निष्कर्ष के रूप में (in conclusion)."),
		mc("on the other hand", "Contrast", "दूसरी ओर", "दूसरी ओर", "इसके अलावा", "अतः", "सबसे पहले"),
		mc("therefore (formal)", "Consequence", "अतः", "अतः", "फिर भी", "दूसरी ओर", "इसके अलावा"),
		mc("What is the 'विपक्ष' in a debate?", "पक्ष और विपक्ष", "The side against the motion", "The side against the motion", "The judges", "The audience", "The motion itself"),
		mc("in conclusion", "Ending", "निष्कर्ष के रूप में", "निष्कर्ष के रूप में", "सबसे पहले", "उदाहरण के लिए", "इसके अलावा"),
		mc("even so / nevertheless", "Concession", "फिर भी", "फिर भी", "इसलिए", "क्योंकि", "अतः"),
		write("Write an argumentative paragraph in Hindi (about 100 words, in Devanagari) on 'सोशल मीडिया युवाओं के लिए हानिकारक है'. Give two arguments, a counter-argument with 'दूसरी ओर', and a conclusion with 'अतः'.",
			"आदरणीय निर्णायक मंडल, मैं इस विषय के पक्ष में बोल रहा हूँ कि सोशल मीडिया युवाओं के लिए हानिकारक है। सबसे पहले, घंटों फ़ोन चलाने से पढ़ाई और नींद दोनों पर बुरा असर पड़ता है। इसके अलावा, दूसरों की चमकदार ज़िंदगी देखकर युवा अपने आप को कम समझने लगते हैं। दूसरी ओर, यह सच है कि सोशल मीडिया से जानकारी और नए अवसर भी मिलते हैं। फिर भी, बिना संयम के इसका उपयोग नुक़सानदेह है। अतः हमें इसका सीमित और समझदारी से प्रयोग करना चाहिए।"),
	)
	addVocab(db, l,
		hv("वाद-विवाद", "vād-vivād", "debate (m)", "कल वाद-विवाद प्रतियोगिता है।", "kal vād-vivād pratiyogitā hai.", "There's a debate competition tomorrow.", finch),
		hv("पक्ष / विपक्ष", "pakṣ / vipakṣ", "for / against (side)", "मैं विपक्ष में बोलूँगी।", "mãi vipakṣ mẽ bolū̃gī.", "I'll speak against the motion.", "Mira"),
		hv("दूसरी ओर", "dūsrī or", "on the other hand", "दूसरी ओर, इसके फ़ायदे भी हैं।", "dūsrī or, iske fāyde bhī haĩ.", "On the other hand, it has benefits too.", finch),
		hv("अतः", "ataḥ", "therefore (formal)", "अतः यह सिद्ध होता है…", "ataḥ yah siddh hotā hai…", "Therefore it is proven…", finch),
		hv("निष्कर्ष", "niṣkarṣ", "conclusion (m)", "निष्कर्ष स्पष्ट है।", "niṣkarṣ spaṣṭ hai.", "The conclusion is clear.", "Zephyr"),
	)

	s = addSkillL(db, hi, u, "Bollywood & Popular Culture", "Film Hindi, Urdu poetry words, festivals.", "Music", "#17A3DD", 40, 820)
	l = addLesson(db, s, "फ़िल्में और त्योहार", 1, 26,
		char("Blaze", "Hindi cinema, based in Mumbai (hence 'Bollywood'), is a huge source of everyday Hindi. Film songs are sung by playback singers while actors lip-sync. Film lyrics lean on Urdu poetic words: दिल (heart), इश्क़ / मोहब्बत (love), ख़्वाब (dream), सनम (beloved). Festivals frame the year: दिवाली (lights), होली (colours), ईद, रक्षाबंधन, दुर्गा पूजा."),
		mc("Where is the Hindi film industry based?", "Bollywood", "Mumbai", "Mumbai", "Delhi", "Kolkata", "Chennai"),
		mc("What is a 'playback singer'?", "फ़िल्मी गाने", "A singer who records songs that actors lip-sync", "A singer who records songs that actors lip-sync", "An actor who sings live", "A DJ", "A music critic"),
		mc("Which festival is the festival of colours?", "रंगों का त्योहार", "होली", "होली", "दिवाली", "ईद", "रक्षाबंधन"),
		mc("What does ख़्वाब mean in film songs?", "ख़्वाब", "dream", "dream", "heart", "moon", "rain"),
		mc("What is celebrated at रक्षाबंधन?", "राखी", "The bond between brothers and sisters", "The bond between brothers and sisters", "The harvest", "The new year", "A wedding"),
		speak("Mira", "दिवाली की हार्दिक शुभकामनाएँ!"),
	)
	addVocab(db, l,
		hv("दिल", "dil", "heart (m)", "मेरा दिल ख़ुश है।", "merā dil khush hai.", "My heart is happy.", "Mira"),
		hv("इश्क़ / मोहब्बत", "ishq / mohabbat", "love (m / f)", "यह इश्क़ की कहानी है।", "yah ishq kī kahānī hai.", "It's a love story.", "Blaze"),
		hv("त्योहार", "tyohār", "festival (m)", "दिवाली मेरा पसंदीदा त्योहार है।", "divālī merā pasandīdā tyohār hai.", "Diwali is my favourite festival.", "Pip"),
		hv("शुभकामनाएँ", "shubhkāmnāẽ", "best wishes (f pl)", "नए साल की शुभकामनाएँ!", "nae sāl kī shubhkāmnāẽ!", "Happy New Year!", "Lumora"),
	)
}

// ───────────────────── C2 — प्रवीणता: Mastery ─────────────────────

func seedHindiC2(db *gorm.DB) {
	const hi = "hi"
	const u = "C2 · प्रवीणता — Mastery"
	finch := "Professor Finch"

	s := addSkillL(db, hi, u, "Hindi, Urdu & Hindustani", "One spoken language, two scripts, two literary registers.", "Languages", "#6C3FC5", 41, 1000)
	l := addLesson(db, s, "हिंदी–उर्दू", 1, 28,
		char(finch, "Everyday Hindi and Urdu share their grammar and core vocabulary — often called Hindustani. They differ in script (Devanagari vs Perso-Arabic Nastaliq) and in formal vocabulary (Sanskrit-derived vs Persian/Arabic-derived). The nuqta letters — क़ ख़ ग़ ज़ फ़ — mark Persian/Arabic sounds; many Hindi speakers merge them with क ख ग ज फ. Urdu poetry forms like the ग़ज़ल are part of shared culture."),
		mc("What do everyday Hindi and Urdu share?", "Hindustani", "Grammar and core vocabulary", "Grammar and core vocabulary", "Only their script", "Nothing", "Only religious words"),
		mc("In which script is Urdu written?", "उर्दू लिपि", "Perso-Arabic (Nastaliq)", "Perso-Arabic (Nastaliq)", "Devanagari", "Latin", "Gurmukhi"),
		mc("What does the dot in क़ (qa) indicate?", "नुक़्ता", "A sound from Persian/Arabic", "A sound from Persian/Arabic", "Aspiration", "A long vowel", "Nasalization"),
		mc("Which word comes from Persian/Arabic?", "Origin", "क़िस्मत (fate)", "क़िस्मत (fate)", "जीवन (life)", "विद्यालय (school)", "प्रश्न (question)"),
		mc("What is a ग़ज़ल?", "ग़ज़ल", "A poem of rhyming couplets (from Persian/Urdu)", "A poem of rhyming couplets (from Persian/Urdu)", "A folk dance", "A type of bread", "A wedding ritual"),
		speak("Zephyr", "क़िस्मत अपनी जगह है, लेकिन मेहनत ज़रूरी है।"),
	)
	addVocab(db, l,
		hv("उर्दू", "urdū", "Urdu (f)", "उर्दू शायरी बहुत सुंदर है।", "urdū shāyrī bahut sundar hai.", "Urdu poetry is very beautiful.", finch),
		hv("क़िस्मत", "qismat", "fate, luck (f)", "यह उसकी क़िस्मत थी।", "yah uskī qismat thī.", "It was his fate.", "Zephyr"),
		hv("ग़ज़ल", "ghazal", "ghazal (f)", "उसने एक ग़ज़ल गाई।", "usne ek ghazal gāī.", "She sang a ghazal.", "Mira"),
		hv("शायरी", "shāyrī", "(Urdu) poetry (f)", "उसे शायरी का शौक़ है।", "use shāyrī kā shauq hai.", "He loves poetry.", "Blaze"),
	)

	s = addSkillL(db, hi, u, "Regional Varieties & Hinglish", "Bhojpuri, Bambaiya Hindi and code-mixing.", "Globe", "#00C2A8", 42, 1030)
	l = addLesson(db, s, "बोलियाँ", 1, 28,
		char(finch, "Hindi is a family of varieties. Bhojpuri (eastern UP and Bihar): का हाल बा? (how are you?). Bambaiya Hindi (Mumbai street speech): क्या बोलता है? (what d'you say?), बिंदास (carefree, cool), झकास (awesome). Haryanvi, Awadhi, Braj and Rajasthani have rich traditions — Tulsidas wrote in Awadhi. Urban speakers mix in English freely (Hinglish): मैं थोड़ा busy हूँ, बाद में call करता हूँ."),
		mc("'का हाल बा?' is from which variety?", "Greeting", "Bhojpuri", "Bhojpuri", "Bambaiya", "Standard Hindi", "Urdu"),
		mc("What does 'बिंदास' mean?", "Bambaiya", "carefree, cool", "carefree, cool", "angry", "tired", "expensive"),
		mc("In which variety did Tulsidas write the रामचरितमानस?", "Language history", "Awadhi", "Awadhi", "Bhojpuri", "Braj", "Standard Hindi"),
		mc("What is Hinglish?", "मैं busy हूँ", "Mixing Hindi and English", "Mixing Hindi and English", "Hindi written in English letters only", "An English dialect of India", "Formal Hindi"),
		mc("Where is Bambaiya Hindi spoken?", "Region", "Mumbai", "Mumbai", "Lucknow", "Patna", "Jaipur"),
		speak("Blaze", "क्या बोलता है भिड़ू? सब बिंदास!"),
	)
	addVocab(db, l,
		hv("बोली", "bolī", "dialect, speech variety (f)", "हर इलाक़े की अपनी बोली है।", "har ilāqe kī apnī bolī hai.", "Every region has its own dialect.", finch),
		hv("बिंदास", "bindās", "carefree, cool (Mumbai)", "वह बिंदास लड़की है।", "vah bindās laṛkī hai.", "She's a carefree girl.", "Blaze"),
		hv("का हाल बा? (भोजपुरी)", "kā hāl bā?", "how are you? (Bhojpuri)", "का हाल बा, भइया?", "kā hāl bā, bhaiyā?", "How are you, brother?", "Zephyr"),
		hv("हिंग्लिश", "hinglish", "Hindi–English mixing", "शहरों में लोग हिंग्लिश बोलते हैं।", "shaharõ mẽ log hinglish bolte haĩ.", "People in cities speak Hinglish.", "Blaze"),
	)

	s = addSkillL(db, hi, u, "Sanskritized & Technical Hindi", "संधि, समास and पारिभाषिक शब्दावली.", "FlaskConical", "#F5A623", 43, 1060)
	l = addLesson(db, s, "शब्द-निर्माण", 1, 28,
		char(finch, "Formal Hindi builds words like Sanskrit. संधि (sound joining): विद्या + आलय = विद्यालय (school), सूर्य + उदय = सूर्योदय (sunrise). समास (compounding): राजपुत्र = राजा का पुत्र (prince), माता-पिता (parents). Official technical vocabulary (पारिभाषिक शब्दावली) offers Sanskrit coinages — दूरभाष (telephone), संगणक (computer) — though everyday speech often prefers the English loans फ़ोन and कंप्यूटर."),
		mc("विद्या + आलय =", "संधि", "विद्यालय", "विद्यालय", "विद्याआलय", "विद्यालाय", "विद्योलय"),
		mc("What does सूर्योदय mean?", "सूर्य + उदय", "sunrise", "sunrise", "sunset", "solar system", "sunlight"),
		mc("What is राजपुत्र?", "समास", "a prince (king's son)", "a prince (king's son)", "a king's palace", "a royal city", "a soldier"),
		mc("Formal Sanskrit-based word for 'telephone'", "पारिभाषिक", "दूरभाष", "दूरभाष", "संगणक", "दूरदर्शन", "विद्यालय"),
		mc("What is दूरदर्शन?", "दूर + दर्शन", "television (also India's public broadcaster)", "television (also India's public broadcaster)", "telescope", "telephone", "computer"),
		speak("Zephyr", "सूर्योदय से पहले हम विद्यालय पहुँच गए।"),
	)
	addVocab(db, l,
		hv("संधि", "sandhi", "sound joining (f)", "विद्यालय संधि से बना है।", "vidyālay sandhi se banā hai.", "विद्यालय is formed by sandhi.", finch),
		hv("समास", "samās", "compound word (m)", "माता-पिता द्वंद्व समास है।", "mātā-pitā dvandva samās hai.", "माता-पिता is a dvandva compound.", finch),
		hv("विद्यालय", "vidyālay", "school (formal) (m)", "हमारा विद्यालय बड़ा है।", "hamārā vidyālay baṛā hai.", "Our school is big.", "Pip"),
		hv("दूरभाष / संगणक", "dūrbhāṣ / sanganak", "telephone / computer (formal)", "दूरभाष संख्या लिखिए।", "dūrbhāṣ sankhyā likhie.", "Write the telephone number.", "Mira"),
	)

	s = addSkillL(db, hi, u, "Poetry: Doha, Chaupai & Ghazal", "मात्राएँ, छंद and how Hindi verse is built.", "Landmark", "#FF5C5C", 44, 1090)
	l = addLesson(db, s, "छंद", 1, 28,
		char(finch, "Hindi verse counts मात्राएँ (morae: short syllables 1, long 2). A दोहा has two lines, each split 13 + 11 मात्राएँ — Kabir's and Rahim's couplets. A चौपाई has four quarters of 16 मात्राएँ each — the metre of Tulsidas's रामचरितमानस. The ग़ज़ल, from Persian/Urdu, is built of शेर (couplets) sharing a रदीफ़ (refrain) and क़ाफ़िया (rhyme)."),
		mc("How are a doha's half-lines split?", "दोहा", "13 + 11 मात्राएँ", "13 + 11 मात्राएँ", "16 + 16", "8 + 8", "5 + 7"),
		mc("Which metre is the रामचरितमानस mostly written in?", "Tulsidas", "चौपाई", "चौपाई", "ग़ज़ल", "सॉनेट", "हाइकु"),
		mc("What is a मात्रा in prosody?", "छंद", "A unit of syllable length", "A unit of syllable length", "A stanza", "A rhyme", "A poet"),
		mc("What is a शेर in a ghazal?", "ग़ज़ल", "A couplet", "A couplet", "A lion", "A refrain only", "A title"),
		mc("What is the रदीफ़?", "ग़ज़ल", "The repeated refrain at the end of lines", "The repeated refrain at the end of lines", "The poet's pen-name", "The first word", "The rhythm of drums"),
		write("Write two lines of your own in Hindi (in Devanagari) in the spirit of a doha about patience, then explain in one Hindi sentence what they mean.",
			"धीरे-धीरे चल सखी, मंज़िल होगी पास।\nजल्दी में जो दौड़ता, खोता है विश्वास॥\nइसका अर्थ है कि धैर्य से चलने वाला ही अपनी मंज़िल तक पहुँचता है।"),
	)
	addVocab(db, l,
		hv("छंद", "chand", "metre, verse form (m)", "दोहा एक छंद है।", "dohā ek chand hai.", "The doha is a verse form.", finch),
		hv("मात्रा", "mātrā", "mora; vowel sign (f)", "इस पंक्ति में तेरह मात्राएँ हैं।", "is pankti mẽ terah mātrāẽ haĩ.", "This line has thirteen morae.", finch),
		hv("चौपाई", "caupāī", "chaupai (16-mora quatrain)", "मानस की चौपाइयाँ प्रसिद्ध हैं।", "mānas kī caupāiyā̃ prasiddh haĩ.", "The chaupais of the Manas are famous.", "Nana"),
		hv("शेर", "sher", "couplet (in a ghazal); lion", "उसने एक शेर सुनाया।", "usne ek sher sunāyā.", "He recited a couplet.", "Blaze"),
	)

	s = addSkillL(db, hi, u, "Speeches & Ceremony", "मंच का संचालन, बधाई, श्रद्धांजलि, धन्यवाद ज्ञापन.", "Mic", "#17A3DD", 45, 1120)
	l = addLesson(db, s, "मंच से", 1, 30,
		char(finch, "Formal speeches open by honouring guests in order: आदरणीय मुख्य अतिथि महोदय, सम्माननीय अध्यक्ष जी, उपस्थित सज्जनो और देवियो… Congratulate: हार्दिक बधाई. At a memorial: विनम्र श्रद्धांजलि. The vote of thanks (धन्यवाद ज्ञापन) closes the event; national occasions often end with जय हिंद!"),
		mc("How does a formal speech open?", "Opening", "आदरणीय मुख्य अतिथि महोदय…", "आदरणीय मुख्य अतिथि महोदय…", "हाय दोस्तो!", "सुनो सब लोग…", "क्या हाल है भिड़ू…"),
		mc("Heartfelt congratulations", "बधाई", "हार्दिक बधाई!", "हार्दिक बधाई!", "विनम्र श्रद्धांजलि!", "फिर मिलेंगे!", "शुभ रात्रि!"),
		mc("At a memorial, you offer…", "Tribute", "विनम्र श्रद्धांजलि", "विनम्र श्रद्धांजलि", "हार्दिक बधाई", "जन्मदिन मुबारक", "शुभ यात्रा"),
		mc("What is धन्यवाद ज्ञापन?", "Programme item", "The vote of thanks", "The vote of thanks", "The welcome song", "The main lecture", "The prize list"),
		mc("'Ladies and gentlemen' (formal)", "Address", "देवियो और सज्जनो", "देवियो और सज्जनो", "भाइयो-बहनो यार", "अरे सब", "दोस्त लोग"),
		write("Write the opening of a formal speech in Hindi (about 100 words, in Devanagari) for your school's annual prize-giving: honour the chief guest and others in order, thank the organisers and congratulate the winners.",
			"आदरणीय मुख्य अतिथि महोदय, सम्माननीय प्रधानाचार्य जी, आदरणीय शिक्षकगण, प्रिय अभिभावको और मेरे साथियो, आप सभी का हार्दिक स्वागत है। आज के इस वार्षिक पुरस्कार वितरण समारोह में बोलना मेरे लिए गर्व की बात है। सबसे पहले मैं उन सभी का धन्यवाद करना चाहता हूँ जिन्होंने इस कार्यक्रम का इतना सुंदर आयोजन किया। साथ ही, आज पुरस्कार पाने वाले सभी विद्यार्थियों को हार्दिक बधाई। आपकी मेहनत हम सबके लिए प्रेरणा है। याद रखिए, बूँद-बूँद से घड़ा भरता है। धन्यवाद।"),
	)
	addVocab(db, l,
		hv("मुख्य अतिथि", "mukhya atithi", "chief guest", "मुख्य अतिथि पधार चुके हैं।", "mukhya atithi padhār cuke haĩ.", "The chief guest has arrived.", finch),
		hv("हार्दिक बधाई", "hārdik badhāī", "heartfelt congratulations", "आपको हार्दिक बधाई!", "āpko hārdik badhāī!", "Heartfelt congratulations to you!", "Lumora"),
		hv("श्रद्धांजलि", "shraddhānjali", "tribute (f)", "हम उन्हें श्रद्धांजलि देते हैं।", "ham unhẽ shraddhānjali dete haĩ.", "We pay tribute to them.", "Nana"),
		hv("समारोह", "samāroh", "ceremony (m)", "समारोह शाम को है।", "samāroh shām ko hai.", "The ceremony is in the evening.", "Mira"),
		hv("सज्जनो और देवियो", "sajjano aur deviyo", "ladies and gentlemen", "देवियो और सज्जनो, स्वागत है।", "deviyo aur sajjano, svāgat hai.", "Ladies and gentlemen, welcome.", finch),
	)
}
