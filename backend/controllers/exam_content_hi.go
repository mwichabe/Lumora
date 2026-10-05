package controllers

// hindiPapers is the hand-authored Hindi proficiency-exam bank, keyed by CEFR
// level. Texts are in Devanagari; A1–A2 texts carry a transliteration.
// Questions are in English through B1 and in Hindi from B2. Hindi separates
// words with spaces (postpositions are separate words), so writing targets
// are word counts, as for the European languages.
var hindiPapers = map[string]paperContent{

	// ---------------- A1 ----------------
	"A1": {
		Listening: PaperListening{
			Title: "राहुल का दिन — Rahul's Day",
			Lines: []PaperLine{
				{Character: "Cora", Text: "नमस्ते! मेरा नाम राहुल है। मैं छात्र हूँ और जयपुर में रहता हूँ।", Translation: "namaste! merā nām rāhul hai. mãi chātra hū̃ aur jaypur mẽ rahtā hū̃. — Hello! My name is Rahul. I'm a student and I live in Jaipur."},
				{Character: "Cora", Text: "मैं रोज़ सुबह छह बजे उठता हूँ और चाय पीता हूँ। मैं साइकिल से कॉलेज जाता हूँ।", Translation: "mãi roz subah chah baje uṭhtā hū̃ aur chāy pītā hū̃. mãi sāikil se kôlej jātā hū̃. — I get up at six every morning and drink tea. I go to college by bicycle."},
				{Character: "Cora", Text: "शाम को मैं दोस्तों के साथ क्रिकेट खेलता हूँ।", Translation: "shām ko mãi dostõ ke sāth krikeṭ kheltā hū̃. — In the evening I play cricket with friends."},
			},
			Questions: []PaperQuestion{
				{Question: "Where does Rahul live?", Options: []string{"Jaipur", "Delhi", "Mumbai", "Agra"}, CorrectAnswer: "Jaipur"},
				{Question: "What time does he get up?", Options: []string{"6:00", "7:00", "5:00", "8:00"}, CorrectAnswer: "6:00"},
				{Question: "How does he go to college?", Options: []string{"By bicycle", "By bus", "On foot", "By car"}, CorrectAnswer: "By bicycle"},
				{Question: "What does he do in the evening?", Options: []string{"Plays cricket with friends", "Reads books", "Watches TV", "Cooks"}, CorrectAnswer: "Plays cricket with friends"},
			},
		},
		Reading: PaperReading{
			Title: "मेरी दोस्त प्रिया — My Friend Priya",
			Paragraphs: []string{
				"मेरी एक दोस्त है। उसका नाम प्रिया है। (merī ek dost hai. uskā nām priyā hai.)",
				"प्रिया डॉक्टर है और वह मुंबई के एक अस्पताल में काम करती है। (priyā ḍākṭar hai aur vah mumbaī ke ek aspatāl mẽ kām kartī hai.)",
				"उसके पास एक छोटी सफ़ेद बिल्ली है। रविवार को हम साथ में खाना खाते हैं। (uske pās ek choṭī safed billī hai. ravivār ko ham sāth mẽ khānā khāte haĩ.)",
			},
			Questions: []PaperQuestion{
				{Question: "What is Priya's job?", Options: []string{"Doctor", "Teacher", "Student", "Engineer"}, CorrectAnswer: "Doctor"},
				{Question: "Where does she work?", Options: []string{"In a hospital in Mumbai", "In a school in Delhi", "At home", "In a shop"}, CorrectAnswer: "In a hospital in Mumbai"},
				{Question: "What pet does she have?", Options: []string{"A small white cat", "A big dog", "A parrot", "A black cat"}, CorrectAnswer: "A small white cat"},
				{Question: "When do they eat together?", Options: []string{"On Sunday", "On Monday", "Every day", "On Saturday"}, CorrectAnswer: "On Sunday"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write a short self-introduction in Hindi (in Devanagari): your name, where you're from, where you live, your family and what you like.",
			MinWords: 30,
		},
		Speaking: PaperSpeaking{
			Phrase:      "नमस्ते! मेरा नाम अन्ना है। मैं केन्या से हूँ और हिंदी सीख रही हूँ।",
			Speaker:     "Lumora",
			Translation: "namaste! merā nām annā hai. mãi kenyā se hū̃ aur hindī sīkh rahī hū̃. — Hello! My name is Anna. I'm from Kenya and I'm learning Hindi.",
		},
	},

	// ---------------- A2 ----------------
	"A2": {
		Listening: PaperListening{
			Title: "गोवा की छुट्टियाँ — A Holiday in Goa",
			Lines: []PaperLine{
				{Character: "Zephyr", Text: "पिछली सर्दियों में मैं अपने दोस्तों के साथ गोवा गया। हम हवाई जहाज़ से गए।", Translation: "pichlī sardiyõ mẽ mãi apne dostõ ke sāth govā gayā. ham havāī jahāz se gae. — Last winter I went to Goa with my friends. We went by plane."},
				{Character: "Zephyr", Text: "हमने समुद्र में तैराकी की और मछली-चावल खाया। हर शाम हम समुद्र किनारे घूमते थे।", Translation: "hamne samudra mẽ tairākī kī aur machlī-chāval khāyā. har shām ham samudra kināre ghūmte the. — We swam in the sea and ate fish-rice. Every evening we walked along the beach."},
				{Character: "Zephyr", Text: "आख़िरी दिन बारिश हुई, इसलिए हम क़िला देखने नहीं जा सके। अगले साल फिर जाएँगे!", Translation: "ākhirī din bārish huī, islie ham qilā dekhne nahī̃ jā sake. agle sāl phir jāẽge! — On the last day it rained, so we couldn't go to see the fort. We'll go again next year!"},
			},
			Questions: []PaperQuestion{
				{Question: "How did they travel to Goa?", Options: []string{"By plane", "By train", "By bus", "By car"}, CorrectAnswer: "By plane"},
				{Question: "What did they eat?", Options: []string{"Fish and rice", "Dal and roti", "Biryani", "Samosas"}, CorrectAnswer: "Fish and rice"},
				{Question: "What did they do every evening?", Options: []string{"Walked along the beach", "Watched films", "Went shopping", "Slept early"}, CorrectAnswer: "Walked along the beach"},
				{Question: "Why couldn't they see the fort?", Options: []string{"It rained", "It was closed", "They were tired", "They had no time"}, CorrectAnswer: "It rained"},
			},
		},
		Reading: PaperReading{
			Title: "विज्ञापन: नौकरी — Advertisement: A Job",
			Paragraphs: []string{
				"एक होटल को रिसेप्शनिस्ट की ज़रूरत है। काम का समय सुबह नौ बजे से शाम पाँच बजे तक है। (ek hoṭal ko risepshanisṭ kī zarūrat hai. kām kā samay subah nau baje se shām pā̃ch baje tak hai.)",
				"उम्मीदवार को हिंदी और अंग्रेज़ी दोनों आनी चाहिए। (ummīdvār ko hindī aur angrezī donõ ānī chāhie.)",
				"इच्छुक लोग पंद्रह तारीख़ से पहले फ़ोन करें। (icchuk log pandrah tārīkh se pahle fon karẽ.)",
			},
			Questions: []PaperQuestion{
				{Question: "What job is offered?", Options: []string{"Receptionist", "Cook", "Driver", "Manager"}, CorrectAnswer: "Receptionist"},
				{Question: "What are the working hours?", Options: []string{"9 a.m. to 5 p.m.", "5 a.m. to 9 a.m.", "9 p.m. to 5 a.m.", "Noon to 9 p.m."}, CorrectAnswer: "9 a.m. to 5 p.m."},
				{Question: "Which languages are required?", Options: []string{"Hindi and English", "Hindi only", "English only", "Urdu and English"}, CorrectAnswer: "Hindi and English"},
				{Question: "How should people apply?", Options: []string{"Phone before the 15th", "Email by the 5th", "Visit on Monday", "Send a letter"}, CorrectAnswer: "Phone before the 15th"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write in Hindi (in Devanagari) about your last weekend: where you went, what you did and how it was. Use the past tense with ने at least twice.",
			MinWords: 55,
		},
		Speaking: PaperSpeaking{
			Phrase:      "कल मैं बाज़ार गई और मैंने फल ख़रीदे।",
			Speaker:     "Cora",
			Translation: "kal mãi bāzār gaī aur mãine phal kharīde. — Yesterday I went to the market and bought fruit.",
		},
	},

	// ---------------- B1 ----------------
	"B1": {
		Listening: PaperListening{
			Title: "नए शहर में घर — A Home in a New City",
			Lines: []PaperLine{
				{Character: "Mira", Text: "मुझे बेंगलुरु में नई नौकरी मिली है, लेकिन वहाँ किराया बहुत ज़्यादा है।", Translation: "I've got a new job in Bengaluru, but rent there is very high."},
				{Character: "Blaze", Text: "तुम शहर के बाहर कमरा ले सकती हो और मेट्रो से दफ़्तर जा सकती हो।", Translation: "You could take a room outside the city and go to the office by metro."},
				{Character: "Mira", Text: "अच्छा सुझाव है। अगर मेट्रो स्टेशन के पास कमरा मिल जाए तो बहुत पैसे बचेंगे।", Translation: "Good suggestion. If I find a room near a metro station, I'll save a lot of money."},
				{Character: "Blaze", Text: "तो शनिवार को साथ चलकर कुछ कमरे देख लेते हैं।", Translation: "Then let's go together on Saturday and look at a few rooms."},
			},
			Questions: []PaperQuestion{
				{Question: "Why is Mira moving to Bengaluru?", Options: []string{"She got a new job", "To study", "Her family lives there", "Rent is cheap there"}, CorrectAnswer: "She got a new job"},
				{Question: "What is the problem?", Options: []string{"Rent is very high", "There's no metro", "She doesn't like the city", "The job is temporary"}, CorrectAnswer: "Rent is very high"},
				{Question: "What does Blaze suggest?", Options: []string{"A room outside the city and commuting by metro", "Staying with him", "Buying a car", "Refusing the job"}, CorrectAnswer: "A room outside the city and commuting by metro"},
				{Question: "What will they do on Saturday?", Options: []string{"Look at some rooms together", "Start the job", "Go on holiday", "Sign a contract"}, CorrectAnswer: "Look at some rooms together"},
			},
		},
		Reading: PaperReading{
			Title: "गाँव की लाइब्रेरी — The Village Library",
			Paragraphs: []string{
				"राजस्थान के एक छोटे गाँव में कुछ युवाओं ने मिलकर एक पुस्तकालय शुरू किया।",
				"पहले गाँव के बच्चों के पास स्कूल की किताबों के अलावा पढ़ने के लिए कुछ नहीं था।",
				"आज इस पुस्तकालय में दो हज़ार से ज़्यादा किताबें हैं, और शाम को बड़े लोग भी यहाँ अख़बार पढ़ने आते हैं।",
			},
			Questions: []PaperQuestion{
				{Question: "Who started the library?", Options: []string{"A group of young people", "The government", "A school", "A rich businessman"}, CorrectAnswer: "A group of young people"},
				{Question: "What was the situation before?", Options: []string{"Children had nothing to read except school books", "There were too many books", "There was no school", "Children didn't like reading"}, CorrectAnswer: "Children had nothing to read except school books"},
				{Question: "How many books does the library have now?", Options: []string{"More than 2,000", "About 200", "20,000", "Exactly 1,000"}, CorrectAnswer: "More than 2,000"},
				{Question: "Who comes in the evening?", Options: []string{"Adults, to read newspapers", "Only teachers", "Tourists", "Nobody"}, CorrectAnswer: "Adults, to read newspapers"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write in Hindi (in Devanagari) about a change you would like to see in your town. Explain why it matters and what would happen if it were made (use अगर… तो…).",
			MinWords: 90,
		},
		Speaking: PaperSpeaking{
			Phrase:      "अगर मेरे पास समय होता, तो मैं गाँव के बच्चों को पढ़ाता।",
			Speaker:     "Zephyr",
			Translation: "If I had time, I would teach the village children.",
		},
	},

	// ---------------- B2 ----------------
	"B2": {
		Listening: PaperListening{
			Title: "बैठक: काग़ज़ की बचत — Meeting: Saving Paper",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "आज की बैठक का उद्देश्य दफ़्तर में काग़ज़ का प्रयोग कम करना है।", Translation: "The purpose of today's meeting is to reduce paper use in the office."},
				{Character: "Mira", Text: "मेरा सुझाव है कि सभी आंतरिक दस्तावेज़ केवल ईमेल द्वारा भेजे जाएँ।", Translation: "I suggest that all internal documents be sent only by email."},
				{Character: "Professor Finch", Text: "अच्छा विचार है, परंतु हमारे कई बुज़ुर्ग ग्राहक अब भी डाक से पत्र पसंद करते हैं।", Translation: "A good idea, but many of our elderly customers still prefer letters by post."},
				{Character: "Mira", Text: "तो केवल उन्हीं ग्राहकों को डाक से पत्र भेजे जाएँ जो इसकी माँग करें।", Translation: "Then letters should be posted only to those customers who request it."},
			},
			Questions: []PaperQuestion{
				{Question: "बैठक का उद्देश्य क्या है?", Options: []string{"काग़ज़ का प्रयोग कम करना", "नए कर्मचारी रखना", "वेतन बढ़ाना", "दफ़्तर बंद करना"}, CorrectAnswer: "काग़ज़ का प्रयोग कम करना"},
				{Question: "मीरा क्या सुझाव देती है?", Options: []string{"आंतरिक दस्तावेज़ केवल ईमेल द्वारा भेजे जाएँ", "ज़्यादा काग़ज़ ख़रीदा जाए", "ग्राहकों को फ़ोन न किया जाए", "सब घर से काम करें"}, CorrectAnswer: "आंतरिक दस्तावेज़ केवल ईमेल द्वारा भेजे जाएँ"},
				{Question: "आपत्ति क्या है?", Options: []string{"बुज़ुर्ग ग्राहक डाक से पत्र पसंद करते हैं", "ईमेल महँगा है", "कंप्यूटर कम हैं", "कर्मचारी सहमत नहीं हैं"}, CorrectAnswer: "बुज़ुर्ग ग्राहक डाक से पत्र पसंद करते हैं"},
				{Question: "अंतिम समाधान क्या है?", Options: []string{"केवल माँग करने वाले ग्राहकों को डाक से पत्र भेजना", "ईमेल बंद करना", "सबको डाक से पत्र भेजना", "ग्राहकों को कंप्यूटर सिखाना"}, CorrectAnswer: "केवल माँग करने वाले ग्राहकों को डाक से पत्र भेजना"},
			},
		},
		Reading: PaperReading{
			Title: "पर्यटन और पर्यावरण — Tourism and the Environment",
			Paragraphs: []string{
				"हिमाचल प्रदेश जैसे पहाड़ी राज्यों में पर्यटन आय का मुख्य स्रोत है। हर वर्ष लाखों पर्यटक यहाँ आते हैं।",
				"परंतु पर्यटकों की बढ़ती संख्या से कई समस्याएँ पैदा हुई हैं: कचरा बढ़ रहा है, पानी की कमी हो रही है और जंगल कट रहे हैं।",
				"विशेषज्ञों का मानना है कि टिकाऊ पर्यटन ही इसका समाधान है, जिसमें स्थानीय लोगों को शामिल किया जाए और पर्यावरण की रक्षा की जाए।",
			},
			Questions: []PaperQuestion{
				{Question: "पहाड़ी राज्यों में पर्यटन क्या है?", Options: []string{"आय का मुख्य स्रोत", "एक छोटी समस्या", "नया शौक़", "सरकारी नियम"}, CorrectAnswer: "आय का मुख्य स्रोत"},
				{Question: "कौन-सी समस्या नहीं बताई गई है?", Options: []string{"स्कूलों का बंद होना", "कचरा बढ़ना", "पानी की कमी", "जंगल कटना"}, CorrectAnswer: "स्कूलों का बंद होना"},
				{Question: "विशेषज्ञ किस समाधान की बात करते हैं?", Options: []string{"टिकाऊ पर्यटन", "पर्यटन पर रोक", "ज़्यादा होटल बनाना", "जंगल काटना"}, CorrectAnswer: "टिकाऊ पर्यटन"},
				{Question: "'टिकाऊ' शब्द का अर्थ क्या है?", Options: []string{"लंबे समय तक चलने वाला", "बहुत महँगा", "बहुत तेज़", "अस्थायी"}, CorrectAnswer: "लंबे समय तक चलने वाला"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write a formal letter in Hindi (in Devanagari) to your municipal officer complaining about rubbish collection in your area and proposing a solution. Use सेवा में, विषय:, महोदय, सविनय निवेदन है कि… and भवदीय/भवदीया.",
			MinWords: 130,
		},
		Speaking: PaperSpeaking{
			Phrase:      "आपकी सहायता के लिए मैं हृदय से आभारी हूँ।",
			Speaker:     "Mira",
			Translation: "I am heartily grateful for your help.",
		},
	},

	// ---------------- C1 ----------------
	"C1": {
		Listening: PaperListening{
			Title: "व्याख्यान: भाषा और पहचान — Lecture: Language and Identity",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "भाषा केवल संवाद का साधन नहीं, बल्कि किसी समाज की सामूहिक स्मृति का भंडार है।", Translation: "Language is not merely a means of communication but a storehouse of a society's collective memory."},
				{Character: "Professor Finch", Text: "उदाहरण के लिए, हमारी लोकोक्तियों और मुहावरों में पीढ़ियों का अनुभव संचित है।", Translation: "For example, the experience of generations is stored in our proverbs and idioms."},
				{Character: "Professor Finch", Text: "फिर भी, भाषा को जड़ मानना भूल होगी; वह तभी जीवित रहती है जब नए विचारों के साथ बदलती रहे।", Translation: "Even so, it would be a mistake to regard language as fixed; it stays alive only if it keeps changing with new ideas."},
			},
			Questions: []PaperQuestion{
				{Question: "वक्ता के अनुसार भाषा संवाद के साधन के अलावा क्या है?", Options: []string{"सामूहिक स्मृति का भंडार", "केवल व्याकरण", "एक विषय", "व्यापार का साधन"}, CorrectAnswer: "सामूहिक स्मृति का भंडार"},
				{Question: "पीढ़ियों का अनुभव कहाँ संचित है?", Options: []string{"लोकोक्तियों और मुहावरों में", "केवल किताबों में", "फ़िल्मों में", "क़ानूनों में"}, CorrectAnswer: "लोकोक्तियों और मुहावरों में"},
				{Question: "वक्ता किस बात को भूल मानते हैं?", Options: []string{"भाषा को जड़ मानना", "भाषा को बदलना", "मुहावरे सीखना", "नए विचार अपनाना"}, CorrectAnswer: "भाषा को जड़ मानना"},
				{Question: "भाषा कब जीवित रहती है?", Options: []string{"जब नए विचारों के साथ बदलती रहे", "जब कभी न बदले", "जब केवल लिखी जाए", "जब कम लोग बोलें"}, CorrectAnswer: "जब नए विचारों के साथ बदलती रहे"},
			},
		},
		Reading: PaperReading{
			Title: "धीरे-धीरे रे मना — Slowly, O Mind",
			Paragraphs: []string{
				"आज का युग गति का युग है। हम तुरंत परिणाम चाहते हैं — तुरंत संदेश, तुरंत भुगतान, तुरंत सफलता।",
				"किंतु कबीर ने सदियों पहले कहा था: 'धीरे-धीरे रे मना, धीरे सब कुछ होय।' जल्दबाज़ी में लिए गए निर्णय प्रायः पछतावे का कारण बनते हैं।",
				"इसका अर्थ यह नहीं कि हम कार्यकुशलता छोड़ दें, अपितु यह समझें कि कौन-सा काम गति माँगता है और कौन-सा धैर्य।",
			},
			Questions: []PaperQuestion{
				{Question: "लेखक के अनुसार आज का युग किसका युग है?", Options: []string{"गति का", "धैर्य का", "परंपरा का", "कविता का"}, CorrectAnswer: "गति का"},
				{Question: "कबीर का उद्धरण किस बात का समर्थन करता है?", Options: []string{"हर काम अपने समय पर होता है", "जल्दी करना अच्छा है", "मन को रोकना असंभव है", "सफलता तुरंत मिलती है"}, CorrectAnswer: "हर काम अपने समय पर होता है"},
				{Question: "'अपितु' का अर्थ क्या है?", Options: []string{"बल्कि", "इसलिए", "क्योंकि", "और"}, CorrectAnswer: "बल्कि"},
				{Question: "लेखक का निष्कर्ष क्या है?", Options: []string{"समझें कि कौन-सा काम गति और कौन-सा धैर्य माँगता है", "कार्यकुशलता छोड़ दें", "सब कुछ जल्दी करें", "कबीर को न पढ़ें"}, CorrectAnswer: "समझें कि कौन-सा काम गति और कौन-सा धैर्य माँगता है"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write an argumentative essay in Hindi (in Devanagari): 'क्या कृत्रिम बुद्धिमत्ता मनुष्यों की नौकरियाँ छीन लेगी?' Present arguments for and against, use at least one मुहावरा or लोकोक्ति, and reach a reasoned conclusion with अतः.",
			MinWords: 180,
		},
		Speaking: PaperSpeaking{
			Phrase:      "निष्कर्ष के रूप में कहा जा सकता है कि शिक्षा ही प्रगति की कुंजी है।",
			Speaker:     "Professor Finch",
			Translation: "In conclusion, it can be said that education is the key to progress.",
		},
	},

	// ---------------- C2 ----------------
	"C2": {
		Listening: PaperListening{
			Title: "श्रद्धांजलि सभा — A Memorial Address",
			Lines: []PaperLine{
				{Character: "Zephyr", Text: "आदरणीय अतिथिगण, आज हम एक ऐसे शिक्षक को श्रद्धांजलि देने के लिए एकत्र हुए हैं जिन्होंने अपना पूरा जीवन ग्रामीण बच्चों की शिक्षा को समर्पित कर दिया।", Translation: "Respected guests, today we have gathered to pay tribute to a teacher who devoted his whole life to the education of rural children."},
				{Character: "Zephyr", Text: "उन्होंने हमें सिखाया कि ज्ञान बाँटने से घटता नहीं, बल्कि बढ़ता है।", Translation: "He taught us that knowledge does not diminish by sharing, but grows."},
				{Character: "Zephyr", Text: "उनके सपनों को साकार करना ही उन्हें हमारी सच्ची श्रद्धांजलि होगी। धन्यवाद।", Translation: "Making his dreams come true will be our true tribute to him. Thank you."},
			},
			Questions: []PaperQuestion{
				{Question: "यह भाषण किस अवसर पर दिया जा रहा है?", Options: []string{"श्रद्धांजलि सभा", "विवाह", "पुरस्कार वितरण", "उद्घाटन"}, CorrectAnswer: "श्रद्धांजलि सभा"},
				{Question: "शिक्षक ने अपना जीवन किसे समर्पित किया?", Options: []string{"ग्रामीण बच्चों की शिक्षा को", "राजनीति को", "व्यापार को", "संगीत को"}, CorrectAnswer: "ग्रामीण बच्चों की शिक्षा को"},
				{Question: "उन्होंने ज्ञान के बारे में क्या सिखाया?", Options: []string{"बाँटने से बढ़ता है", "बाँटने से घटता है", "छिपाकर रखना चाहिए", "केवल किताबों में होता है"}, CorrectAnswer: "बाँटने से बढ़ता है"},
				{Question: "वक्ता के अनुसार सच्ची श्रद्धांजलि क्या होगी?", Options: []string{"उनके सपनों को साकार करना", "उनकी मूर्ति बनाना", "छुट्टी मनाना", "भाषण देना"}, CorrectAnswer: "उनके सपनों को साकार करना"},
			},
		},
		Reading: PaperReading{
			Title: "राजभाषा विभाग की सूचना — An Official Language Department Notice",
			Paragraphs: []string{
				"सर्वसाधारण को सूचित किया जाता है कि दिनांक 1 अप्रैल से समस्त कार्यालयी पत्राचार हिंदी अथवा द्विभाषी रूप में किया जाएगा।",
				"जिन कर्मचारियों को हिंदी में टिप्पण एवं आलेखन में कठिनाई हो, वे विभाग द्वारा आयोजित प्रशिक्षण कार्यक्रम में भाग ले सकते हैं।",
				"इस संबंध में किसी भी जानकारी हेतु संबंधित अनुभाग अधिकारी से संपर्क किया जा सकता है।",
			},
			Questions: []PaperQuestion{
				{Question: "1 अप्रैल से क्या बदलेगा?", Options: []string{"पत्राचार हिंदी या द्विभाषी रूप में होगा", "कार्यालय बंद होंगे", "केवल अंग्रेज़ी में काम होगा", "वेतन बढ़ेगा"}, CorrectAnswer: "पत्राचार हिंदी या द्विभाषी रूप में होगा"},
				{Question: "कठिनाई वाले कर्मचारी क्या कर सकते हैं?", Options: []string{"प्रशिक्षण कार्यक्रम में भाग ले सकते हैं", "काम छोड़ सकते हैं", "अंग्रेज़ी में लिख सकते हैं", "छुट्टी ले सकते हैं"}, CorrectAnswer: "प्रशिक्षण कार्यक्रम में भाग ले सकते हैं"},
				{Question: "'द्विभाषी' का अर्थ क्या है?", Options: []string{"दो भाषाओं वाला", "दो पृष्ठों वाला", "दोहरा वेतन", "दो कार्यालय"}, CorrectAnswer: "दो भाषाओं वाला"},
				{Question: "'हेतु' का सरल अर्थ क्या है?", Options: []string{"के लिए", "के बाद", "के बिना", "के पास"}, CorrectAnswer: "के लिए"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write a formal speech in Hindi (in Devanagari) for a ceremony honouring teachers: address the guests in order of rank, reflect on the role of teachers using at least one doha, proverb or idiom, and close formally.",
			MinWords: 230,
		},
		Speaking: PaperSpeaking{
			Phrase:      "आदरणीय मुख्य अतिथि महोदय, आज इस मंच से बोलना मेरे लिए गर्व की बात है।",
			Speaker:     "Zephyr",
			Translation: "Respected chief guest, it is a matter of pride for me to speak from this stage today.",
		},
	},

	// ---------------- FINAL · comprehensive ----------------
	"FINAL": {
		Listening: PaperListening{
			Title: "परिचर्चा: डिजिटल युग में हिंदी — Panel: Hindi in the Digital Age",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "आज हिंदी जितनी लिखी जा रही है, उतनी पहले कभी नहीं लिखी गई — संदेशों, सोशल मीडिया और ब्लॉगों में।", Translation: "Hindi is being written today more than ever before — in messages, on social media and in blogs."},
				{Character: "Mira", Text: "यह सच है, किंतु रोमन लिपि में हिंदी लिखने की बढ़ती प्रवृत्ति देवनागरी के लिए चुनौती बन गई है।", Translation: "That's true, but the growing tendency to write Hindi in Roman script has become a challenge for Devanagari."},
				{Character: "Professor Finch", Text: "समाधान प्रतिबंध नहीं, सुविधा है: देवनागरी टाइपिंग को जितना सरल बनाएँगे, लोग उतना ही उसे अपनाएँगे।", Translation: "The solution is not prohibition but convenience: the simpler we make Devanagari typing, the more people will adopt it."},
			},
			Questions: []PaperQuestion{
				{Question: "पहले वक्ता के अनुसार आज हिंदी के बारे में क्या सच है?", Options: []string{"पहले से कहीं अधिक लिखी जा रही है", "कम बोली जा रही है", "केवल किताबों में है", "समाप्त हो रही है"}, CorrectAnswer: "पहले से कहीं अधिक लिखी जा रही है"},
				{Question: "मीरा किस चुनौती की बात करती हैं?", Options: []string{"रोमन लिपि में हिंदी लिखने की प्रवृत्ति", "हिंदी फ़िल्मों का कम होना", "व्याकरण की कठिनाई", "शब्दकोश की कमी"}, CorrectAnswer: "रोमन लिपि में हिंदी लिखने की प्रवृत्ति"},
				{Question: "प्रस्तावित समाधान क्या है?", Options: []string{"देवनागरी टाइपिंग को सरल बनाना", "रोमन लिपि पर प्रतिबंध", "अंग्रेज़ी अपनाना", "लिखना बंद करना"}, CorrectAnswer: "देवनागरी टाइपिंग को सरल बनाना"},
				{Question: "'प्रवृत्ति' का अर्थ क्या है?", Options: []string{"झुकाव / चलन", "प्रतिबंध", "समाधान", "पुस्तक"}, CorrectAnswer: "झुकाव / चलन"},
			},
		},
		Reading: PaperReading{
			Title: "अनुवाद की कला — The Art of Translation",
			Paragraphs: []string{
				"अनुवाद केवल शब्दों का स्थानांतरण नहीं, बल्कि एक संस्कृति के भावों को दूसरी संस्कृति तक पहुँचाने का प्रयास है।",
				"उदाहरण के लिए, 'जी' जैसे छोटे-से शब्द में जो आदर छिपा है, उसे अंग्रेज़ी में एक शब्द में व्यक्त करना लगभग असंभव है।",
				"अतः श्रेष्ठ अनुवादक वही है जो दो भाषाओं के साथ-साथ दो संस्कृतियों को भी समझता हो और जानता हो कि कब मूल शब्द रखना है और कब उसकी व्याख्या करनी है।",
			},
			Questions: []PaperQuestion{
				{Question: "लेखक के अनुसार अनुवाद क्या है?", Options: []string{"एक संस्कृति के भाव दूसरी तक पहुँचाने का प्रयास", "केवल शब्द बदलना", "सार लिखना", "व्याकरण सुधारना"}, CorrectAnswer: "एक संस्कृति के भाव दूसरी तक पहुँचाने का प्रयास"},
				{Question: "'जी' का उदाहरण क्यों दिया गया है?", Options: []string{"उसमें छिपे आदर को एक शब्द में अनुवाद करना कठिन है", "वह बहुत लंबा शब्द है", "वह उर्दू का शब्द है", "उसका कोई अर्थ नहीं"}, CorrectAnswer: "उसमें छिपे आदर को एक शब्द में अनुवाद करना कठिन है"},
				{Question: "श्रेष्ठ अनुवादक कौन है?", Options: []string{"जो दो संस्कृतियों को समझता हो", "जो सबसे तेज़ अनुवाद करे", "जो केवल शब्दकोश प्रयोग करे", "जो कई भाषाएँ जानता हो"}, CorrectAnswer: "जो दो संस्कृतियों को समझता हो"},
				{Question: "'अतः' का अर्थ क्या है?", Options: []string{"इसलिए", "लेकिन", "क्योंकि", "फिर भी"}, CorrectAnswer: "इसलिए"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write an essay in Hindi (in Devanagari) on 'वैश्वीकरण के युग में मातृभाषा का महत्त्व'. Develop a nuanced argument across registers, engage a counter-position, use idioms or a literary reference, and draw a reasoned conclusion.",
			MinWords: 300,
		},
		Speaking: PaperSpeaking{
			Phrase:      "भाषा संस्कृति का दर्पण है; उसे सँजोना अपनी पहचान को सँजोना है।",
			Speaker:     "Lumora",
			Translation: "Language is the mirror of culture; to cherish it is to cherish our own identity.",
		},
	},
}
