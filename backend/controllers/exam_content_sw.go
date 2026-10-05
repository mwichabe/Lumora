package controllers

// swahiliPapers is the hand-authored Swahili (Kiswahili) proficiency-exam bank,
// keyed by CEFR level like the German and Spanish banks. Questions are in
// English through B1 and in Swahili from B2. Swahili is agglutinative — one
// word carries subject, tense, object and verb — so writing targets (in words)
// sit a little below the equivalent German ones.
var swahiliPapers = map[string]paperContent{

	// ---------------- A1 ----------------
	"A1": {
		Listening: PaperListening{
			Title: "Siku ya Baraka — Baraka's Day",
			Lines: []PaperLine{
				{Character: "Cora", Text: "Habari! Jina langu ni Baraka. Mimi ni mwanafunzi na ninaishi Mombasa.", Translation: "Hello! My name is Baraka. I'm a student and I live in Mombasa."},
				{Character: "Cora", Text: "Kila siku ninaamka saa kumi na mbili asubuhi. Ninakunywa chai na ninakula mkate.", Translation: "Every day I wake up at 6 a.m. I drink tea and eat bread."},
				{Character: "Cora", Text: "Ninakwenda shule kwa baiskeli. Jioni ninacheza mpira na rafiki zangu.", Translation: "I go to school by bicycle. In the evening I play football with my friends."},
			},
			Questions: []PaperQuestion{
				{Question: "Where does Baraka live?", Options: []string{"Mombasa", "Nairobi", "Arusha", "Kampala"}, CorrectAnswer: "Mombasa"},
				{Question: "What does he have for breakfast?", Options: []string{"Tea and bread", "Coffee and eggs", "Ugali", "Rice and beans"}, CorrectAnswer: "Tea and bread"},
				{Question: "How does he go to school?", Options: []string{"By bicycle", "By bus", "On foot", "By car"}, CorrectAnswer: "By bicycle"},
				{Question: "What does he do in the evening?", Options: []string{"Plays football with friends", "Reads books", "Cooks dinner", "Watches TV"}, CorrectAnswer: "Plays football with friends"},
			},
		},
		Reading: PaperReading{
			Title: "Rafiki Yangu Neema — My Friend Neema",
			Paragraphs: []string{
				"Nina rafiki anaitwa Neema. Yeye ni Mtanzania na anaishi Dar es Salaam.",
				"Neema ni daktari. Anafanya kazi hospitalini. Ana watoto wawili na paka mmoja.",
				"Jumamosi tulikwenda sokoni pamoja. Tulinunua matunda na mboga, kisha tukanywa chai.",
			},
			Questions: []PaperQuestion{
				{Question: "Where is Neema from?", Options: []string{"Tanzania", "Kenya", "Uganda", "Rwanda"}, CorrectAnswer: "Tanzania"},
				{Question: "What is her job?", Options: []string{"Doctor", "Teacher", "Farmer", "Driver"}, CorrectAnswer: "Doctor"},
				{Question: "How many children does she have?", Options: []string{"Two", "One", "Three", "None"}, CorrectAnswer: "Two"},
				{Question: "What did they buy at the market?", Options: []string{"Fruit and vegetables", "Clothes", "Fish", "Shoes"}, CorrectAnswer: "Fruit and vegetables"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write a short self-introduction in Swahili: your name, where you're from, where you live, your family, and what you like.",
			MinWords: 30,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Habari! Jina langu ni Anna. Ninatoka Kenya na ninajifunza Kiswahili.",
			Speaker:     "Lumora",
			Translation: "Hello! My name is Anna. I'm from Kenya and I'm learning Swahili.",
		},
	},

	// ---------------- A2 ----------------
	"A2": {
		Listening: PaperListening{
			Title: "Likizo Zanzibar — A Holiday in Zanzibar",
			Lines: []PaperLine{
				{Character: "Zephyr", Text: "Mwezi uliopita nilikwenda Zanzibar na rafiki yangu. Tulisafiri kwa boti kutoka Dar es Salaam.", Translation: "Last month I went to Zanzibar with my friend. We travelled by boat from Dar es Salaam."},
				{Character: "Zephyr", Text: "Tulitembelea Mji Mkongwe na tulikula samaki wa kuchoma usiku kwenye bustani ya Forodhani.", Translation: "We visited Stone Town and ate grilled fish at night in the Forodhani gardens."},
				{Character: "Zephyr", Text: "Siku ya tatu mvua kubwa ilinyesha, kwa hiyo hatukwenda ufukweni. Mwakani nitarudi tena!", Translation: "On the third day heavy rain fell, so we didn't go to the beach. Next year I'll go back again!"},
			},
			Questions: []PaperQuestion{
				{Question: "How did they travel to Zanzibar?", Options: []string{"By boat", "By plane", "By bus", "By train"}, CorrectAnswer: "By boat"},
				{Question: "What did they eat at night?", Options: []string{"Grilled fish", "Ugali and meat", "Rice and beans", "Chips"}, CorrectAnswer: "Grilled fish"},
				{Question: "Why didn't they go to the beach on the third day?", Options: []string{"It rained heavily", "They were tired", "The beach was closed", "They went home"}, CorrectAnswer: "It rained heavily"},
				{Question: "What will he do next year?", Options: []string{"Go back to Zanzibar", "Move to Zanzibar", "Visit Mombasa", "Stay at home"}, CorrectAnswer: "Go back to Zanzibar"},
			},
		},
		Reading: PaperReading{
			Title: "Tangazo: Nafasi ya Kazi — Notice: Job Vacancy",
			Paragraphs: []string{
				"Hoteli ya Bahari inatafuta mhudumu wa mgahawa. Kazi ni kuanzia saa mbili asubuhi hadi saa kumi jioni.",
				"Mwombaji anahitaji kujua Kiswahili na Kiingereza. Uzoefu wa mwaka mmoja ni faida.",
				"Watu wanaopenda nafasi hii wapige simu kabla ya tarehe kumi na tano mwezi huu.",
			},
			Questions: []PaperQuestion{
				{Question: "What job is offered?", Options: []string{"Restaurant waiter", "Hotel manager", "Driver", "Cook"}, CorrectAnswer: "Restaurant waiter"},
				{Question: "What are the working hours (English time)?", Options: []string{"8 a.m. to 4 p.m.", "2 a.m. to 10 a.m.", "8 a.m. to 10 p.m.", "6 a.m. to 4 p.m."}, CorrectAnswer: "8 a.m. to 4 p.m."},
				{Question: "Which languages are required?", Options: []string{"Swahili and English", "Swahili only", "English and French", "Arabic"}, CorrectAnswer: "Swahili and English"},
				{Question: "How should people apply?", Options: []string{"Phone before the 15th", "Email before the 5th", "Visit on Monday", "Send a letter"}, CorrectAnswer: "Phone before the 15th"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write in Swahili about your last weekend: where you went, what you did, and how it was. Use the past tense (-li-) and at least one negative (-ku-).",
			MinWords: 50,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Jana nilikwenda sokoni, nikanunua matunda, kisha nikarudi nyumbani.",
			Speaker:     "Blaze",
			Translation: "Yesterday I went to the market, bought fruit, then went back home.",
		},
	},

	// ---------------- B1 ----------------
	"B1": {
		Listening: PaperListening{
			Title: "Kuhamia Mjini — Moving to the City",
			Lines: []PaperLine{
				{Character: "Mira", Text: "Ninafikiria kuhamia Nairobi kwa sababu nimepata kazi mpya huko, lakini kodi ya nyumba ni ghali sana.", Translation: "I'm thinking of moving to Nairobi because I've found a new job there, but rent is very expensive."},
				{Character: "Blaze", Text: "Ungeweza kuishi nje kidogo ya mji na kusafiri kwa treni kila siku.", Translation: "You could live a little outside the city and commute by train every day."},
				{Character: "Mira", Text: "Ni kweli. Nikipata nyumba karibu na kituo cha treni, nitaokoa pesa nyingi.", Translation: "That's true. If I find a house near a train station, I'll save a lot of money."},
				{Character: "Blaze", Text: "Basi Jumamosi tuende pamoja tukaangalie nyumba kadhaa.", Translation: "Then on Saturday let's go together and look at a few houses."},
			},
			Questions: []PaperQuestion{
				{Question: "Why does Mira want to move to Nairobi?", Options: []string{"She found a new job there", "Her family lives there", "She wants to study", "Rent is cheap there"}, CorrectAnswer: "She found a new job there"},
				{Question: "What is her worry?", Options: []string{"Rent is very expensive", "The city is dangerous", "She has no friends there", "The job is boring"}, CorrectAnswer: "Rent is very expensive"},
				{Question: "What does Blaze suggest?", Options: []string{"Live outside the city and commute by train", "Stay where she is", "Buy a car", "Share a flat with him"}, CorrectAnswer: "Live outside the city and commute by train"},
				{Question: "What will they do on Saturday?", Options: []string{"Look at some houses together", "Start the new job", "Travel to the coast", "Sign a contract"}, CorrectAnswer: "Look at some houses together"},
			},
		},
		Reading: PaperReading{
			Title: "Maji Safi Kijijini — Clean Water in the Village",
			Paragraphs: []string{
				"Kwa miaka mingi, wanawake wa kijiji cha Mlimani walitembea kilomita tano kila siku kuchota maji mtoni.",
				"Mwaka jana, kikundi cha vijana kilianzisha mradi wa kuchimba kisima kwa msaada wa shirika moja la maendeleo.",
				"Sasa kisima kipo katikati ya kijiji. Wasichana wengi ambao zamani walikosa masomo kwa sababu ya kuchota maji wanakwenda shule kila siku.",
			},
			Questions: []PaperQuestion{
				{Question: "How far did the women walk for water before?", Options: []string{"Five kilometres", "Fifteen kilometres", "Half a kilometre", "Fifty kilometres"}, CorrectAnswer: "Five kilometres"},
				{Question: "Who started the well project?", Options: []string{"A group of young people", "The government", "The village elders", "A school"}, CorrectAnswer: "A group of young people"},
				{Question: "Where is the well now?", Options: []string{"In the middle of the village", "By the river", "At the school", "In the next village"}, CorrectAnswer: "In the middle of the village"},
				{Question: "What is one result of the project?", Options: []string{"More girls go to school every day", "The river dried up", "People moved to the city", "Water became expensive"}, CorrectAnswer: "More girls go to school every day"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write in Swahili about a change you would like to see in your community. Say why it matters and what would happen if it were made (use -ki- or -nge-).",
			MinWords: 80,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Ningekuwa na muda zaidi, ningejitolea kufundisha watoto wa kijijini.",
			Speaker:     "Mira",
			Translation: "If I had more time, I would volunteer to teach children in the village.",
		},
	},

	// ---------------- B2 ----------------
	"B2": {
		Listening: PaperListening{
			Title: "Mkutano wa Wafanyakazi — A Staff Meeting",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "Ndugu wafanyakazi, lengo la mkutano huu ni kujadili mpango wa kupunguza matumizi ya karatasi ofisini.", Translation: "Dear staff, the purpose of this meeting is to discuss the plan to reduce paper use in the office."},
				{Character: "Mira", Text: "Napendekeza kwamba ripoti zote zitumwe kwa barua pepe badala ya kuchapishwa.", Translation: "I propose that all reports be sent by email instead of being printed."},
				{Character: "Professor Finch", Text: "Pendekezo zuri. Hata hivyo, baadhi ya wateja wetu wazee hawatumii barua pepe; tutawahudumia vipi?", Translation: "A good proposal. However, some of our older clients don't use email; how will we serve them?"},
				{Character: "Mira", Text: "Tunaweza kuendelea kuwatumia nakala za karatasi wao peke yao, na wengine wote watumie mfumo wa kidijitali.", Translation: "We can continue sending paper copies to them alone, and everyone else can use the digital system."},
			},
			Questions: []PaperQuestion{
				{Question: "Lengo la mkutano ni nini?", Options: []string{"Kupunguza matumizi ya karatasi", "Kuongeza mishahara", "Kuajiri wafanyakazi wapya", "Kufunga ofisi"}, CorrectAnswer: "Kupunguza matumizi ya karatasi"},
				{Question: "Mira anapendekeza nini?", Options: []string{"Ripoti zitumwe kwa barua pepe", "Ripoti zichapishwe zaidi", "Wateja waache kuja", "Ofisi ifungwe"}, CorrectAnswer: "Ripoti zitumwe kwa barua pepe"},
				{Question: "Wasiwasi wa mwenyekiti ni upi?", Options: []string{"Wateja wazee hawatumii barua pepe", "Barua pepe ni ghali", "Kompyuta hazitoshi", "Wafanyakazi hawapendi mabadiliko"}, CorrectAnswer: "Wateja wazee hawatumii barua pepe"},
				{Question: "Suluhisho ni lipi?", Options: []string{"Kuwatumia wateja wazee nakala za karatasi peke yao", "Kuacha kutumia barua pepe", "Kuwafundisha wateja wote kompyuta", "Kuchapisha kila kitu"}, CorrectAnswer: "Kuwatumia wateja wazee nakala za karatasi peke yao"},
			},
		},
		Reading: PaperReading{
			Title: "Utalii na Mazingira — Tourism and the Environment",
			Paragraphs: []string{
				"Utalii ni mojawapo ya vyanzo vikuu vya mapato kwa nchi za Afrika Mashariki. Mbuga za wanyama kama Serengeti na Maasai Mara huvutia wageni kutoka pande zote za dunia.",
				"Hata hivyo, ongezeko la watalii limeleta changamoto: msongamano wa magari mbugani unavuruga wanyama, na taka zinaongezeka.",
				"Wataalamu wanapendekeza utalii endelevu, ambapo jamii za wenyeji hushirikishwa na sehemu ya mapato hutumika kuhifadhi mazingira.",
			},
			Questions: []PaperQuestion{
				{Question: "Kwa mujibu wa makala, utalii ni nini kwa Afrika Mashariki?", Options: []string{"Chanzo kikuu cha mapato", "Tatizo kubwa tu", "Shughuli isiyo na faida", "Jambo jipya"}, CorrectAnswer: "Chanzo kikuu cha mapato"},
				{Question: "Ni changamoto ipi imetajwa?", Options: []string{"Msongamano wa magari unavuruga wanyama", "Ukosefu wa hoteli", "Bei ya chakula", "Uhaba wa mvua"}, CorrectAnswer: "Msongamano wa magari unavuruga wanyama"},
				{Question: "Utalii endelevu unahusisha nini?", Options: []string{"Kushirikisha jamii za wenyeji na kuhifadhi mazingira", "Kufunga mbuga zote", "Kuongeza idadi ya magari", "Kupiga marufuku watalii"}, CorrectAnswer: "Kushirikisha jamii za wenyeji na kuhifadhi mazingira"},
				{Question: "Neno 'wenyeji' lina maana gani hapa?", Options: []string{"Wakazi wa eneo husika", "Wageni kutoka nje", "Wanyama wa porini", "Viongozi wa serikali"}, CorrectAnswer: "Wakazi wa eneo husika"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write a formal letter in Swahili to your local council (Kwa: Mkurugenzi wa Halmashauri…) complaining about rubbish collection in your area and proposing a solution. Use YAH: and Wako mtiifu.",
			MinWords: 120,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Napenda kuchukua nafasi hii kuwashukuru wote kwa ushirikiano wenu.",
			Speaker:     "Mira",
			Translation: "I'd like to take this opportunity to thank everyone for your cooperation.",
		},
	},

	// ---------------- C1 ----------------
	"C1": {
		Listening: PaperListening{
			Title: "Mhadhara: Lugha na Utambulisho — Lecture: Language and Identity",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "Lugha si chombo cha mawasiliano peke yake; ni hazina ya historia, mila na mtazamo wa jamii kuhusu dunia.", Translation: "Language is not merely a tool of communication; it is a treasury of a community's history, customs and worldview."},
				{Character: "Professor Finch", Text: "Kiswahili, kwa mfano, kimebeba athari za biashara ya karne nyingi katika Bahari ya Hindi, jambo linalodhihirika katika msamiati wake wa Kiarabu, Kiajemi na Kireno.", Translation: "Swahili, for example, carries the influence of centuries of Indian Ocean trade, evident in its Arabic, Persian and Portuguese vocabulary."},
				{Character: "Professor Finch", Text: "Hata hivyo, lugha inapoenea, hukabiliwa na hatari ya kupoteza lahaja zake. Kwa hiyo, ni wajibu wetu kuzihifadhi bila kudhoofisha Kiswahili sanifu.", Translation: "However, as a language spreads, it risks losing its dialects. It is therefore our duty to preserve them without weakening Standard Swahili."},
			},
			Questions: []PaperQuestion{
				{Question: "Mhadhiri anasema lugha ni nini zaidi ya chombo cha mawasiliano?", Options: []string{"Hazina ya historia, mila na mtazamo", "Somo la shule tu", "Chombo cha biashara pekee", "Mfumo wa sarufi tu"}, CorrectAnswer: "Hazina ya historia, mila na mtazamo"},
				{Question: "Athari ya biashara ya Bahari ya Hindi inaonekana wapi?", Options: []string{"Katika msamiati wa Kiswahili", "Katika matamshi tu", "Katika ngeli", "Katika muziki"}, CorrectAnswer: "Katika msamiati wa Kiswahili"},
				{Question: "Ni hatari gani inayotajwa?", Options: []string{"Kupoteza lahaja", "Kupoteza wazungumzaji wote", "Kuchanganya na Kiingereza", "Kukosa walimu"}, CorrectAnswer: "Kupoteza lahaja"},
				{Question: "Mhadhiri anapendekeza nini?", Options: []string{"Kuhifadhi lahaja bila kudhoofisha Kiswahili sanifu", "Kuacha lahaja zote", "Kuacha Kiswahili sanifu", "Kufundisha Kiarabu"}, CorrectAnswer: "Kuhifadhi lahaja bila kudhoofisha Kiswahili sanifu"},
			},
		},
		Reading: PaperReading{
			Title: "Haraka Haraka Haina Baraka — Haste Has No Blessing",
			Paragraphs: []string{
				"Katika ulimwengu wa leo, kasi imekuwa kipimo cha mafanikio. Tunataka huduma za papo kwa papo, majibu ya haraka na faida za muda mfupi.",
				"Lakini wahenga walituonya: haraka haraka haina baraka. Maamuzi yanayofanywa kwa pupa mara nyingi huzaa majuto, iwe katika biashara, siasa au mahusiano.",
				"Hii haimaanishi kwamba tuache ufanisi. Badala yake, tujifunze kutofautisha kati ya mambo yanayohitaji kasi na yale yanayohitaji subira na tafakuri.",
			},
			Questions: []PaperQuestion{
				{Question: "Kwa mujibu wa mwandishi, kipimo cha mafanikio leo ni kipi?", Options: []string{"Kasi", "Subira", "Elimu", "Utajiri wa kale"}, CorrectAnswer: "Kasi"},
				{Question: "Methali 'haraka haraka haina baraka' inatumika kuonyesha nini?", Options: []string{"Maamuzi ya pupa huzaa majuto", "Kasi huleta baraka", "Biashara ni mbaya", "Siasa haina maana"}, CorrectAnswer: "Maamuzi ya pupa huzaa majuto"},
				{Question: "Neno 'pupa' lina maana gani?", Options: []string{"Haraka isiyo na busara", "Upole", "Busara", "Furaha"}, CorrectAnswer: "Haraka isiyo na busara"},
				{Question: "Hitimisho la mwandishi ni lipi?", Options: []string{"Kutofautisha mambo yanayohitaji kasi na yanayohitaji subira", "Kuacha ufanisi kabisa", "Kufanya kila kitu haraka", "Kuwasikiliza wahenga peke yao"}, CorrectAnswer: "Kutofautisha mambo yanayohitaji kasi na yanayohitaji subira"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write an argumentative essay in Swahili: 'Mitandao ya kijamii imeleta madhara kuliko faida kwa vijana.' Present arguments for and against, use at least one methali, and reach a reasoned conclusion.",
			MinWords: 170,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Kwa kuhitimisha, ni wazi kwamba elimu ndiyo ufunguo wa maendeleo ya jamii yoyote.",
			Speaker:     "Professor Finch",
			Translation: "In conclusion, it is clear that education is the key to the development of any society.",
		},
	},

	// ---------------- C2 ----------------
	"C2": {
		Listening: PaperListening{
			Title: "Hotuba ya Siku ya Kiswahili — A World Kiswahili Day Speech",
			Lines: []PaperLine{
				{Character: "Zephyr", Text: "Mheshimiwa Mgeni Rasmi, waheshimiwa viongozi, mabibi na mabwana: leo tunaadhimisha Siku ya Kiswahili Duniani.", Translation: "Honourable Guest of Honour, honourable leaders, ladies and gentlemen: today we celebrate World Kiswahili Day."},
				{Character: "Zephyr", Text: "Lugha hii, iliyozaliwa katika miji ya pwani, sasa imevuka mipaka ya bara letu na kufundishwa katika vyuo vikuu duniani kote.", Translation: "This language, born in the coastal towns, has now crossed our continent's borders and is taught in universities worldwide."},
				{Character: "Zephyr", Text: "Ni jukumu letu kuienzi, kuiandikia na kuitumia kwa ufasaha, ili vizazi vijavyo virithi hazina hii ikiwa imekamilika. Asanteni kwa kunisikiliza.", Translation: "It is our duty to honour it, write in it and use it eloquently, so that future generations inherit this treasure whole. Thank you for listening."},
			},
			Questions: []PaperQuestion{
				{Question: "Hotuba inatolewa katika tukio gani?", Options: []string{"Siku ya Kiswahili Duniani", "Siku ya Uhuru", "Mahafali ya chuo", "Harusi"}, CorrectAnswer: "Siku ya Kiswahili Duniani"},
				{Question: "Kwa mujibu wa msemaji, Kiswahili kilizaliwa wapi?", Options: []string{"Katika miji ya pwani", "Katika milima ya bara", "Ulaya", "Arabuni"}, CorrectAnswer: "Katika miji ya pwani"},
				{Question: "Msemaji anaona jukumu letu ni lipi?", Options: []string{"Kukienzi na kukitumia kwa ufasaha", "Kukibadilisha kabisa", "Kukiacha", "Kukifundisha nje tu"}, CorrectAnswer: "Kukienzi na kukitumia kwa ufasaha"},
				{Question: "Hotuba inafungwa kwa maneno gani?", Options: []string{"Asanteni kwa kunisikiliza", "Kwaherini wote", "Mungu ibariki Afrika", "Hongera sana"}, CorrectAnswer: "Asanteni kwa kunisikiliza"},
			},
		},
		Reading: PaperReading{
			Title: "Ubeti wa Utenzi — A Stanza of Epic Poetry",
			Paragraphs: []string{
				"Bismillahi kwanza, / ndiyo kauli ya mwanzo, / na Mungu ndiye mwenza, / wa kuandika utenzi.",
				"(Maelezo) Utenzi huanza kwa kumtaja Mungu, kwa kufuata desturi ya washairi wa pwani walioandika kwa hati za Kiarabu. Kila mstari una mizani minane.",
				"Mtindo huu, unaopatikana katika Utendi wa Tambuka na tenzi nyingine za kale, unaonyesha jinsi dini, historia na sanaa vilivyoungana katika fasihi ya Kiswahili.",
			},
			Questions: []PaperQuestion{
				{Question: "Kila mstari wa utenzi huu una mizani mingapi?", Options: []string{"Minane", "Kumi na sita", "Kumi na mbili", "Minne"}, CorrectAnswer: "Minane"},
				{Question: "Kwa nini utenzi huanza kwa kumtaja Mungu?", Options: []string{"Ni desturi ya washairi wa pwani", "Ni sheria ya serikali", "Ni kwa ajili ya mizani tu", "Ni kosa la mwandishi"}, CorrectAnswer: "Ni desturi ya washairi wa pwani"},
				{Question: "Tenzi za kale ziliandikwa kwa hati gani?", Options: []string{"Za Kiarabu", "Za Kilatini", "Za Kigiriki", "Za Kichina"}, CorrectAnswer: "Za Kiarabu"},
				{Question: "Mtindo huu unaonyesha nini kuhusu fasihi ya Kiswahili?", Options: []string{"Dini, historia na sanaa viliungana ndani yake", "Ni fasihi mpya kabisa", "Haikuhusiana na dini", "Iliandikwa na wageni tu"}, CorrectAnswer: "Dini, historia na sanaa viliungana ndani yake"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write a formal speech in Swahili for a national ceremony honouring teachers: greet the guests in order of rank, use at least two methali, reflect on the role of teachers in society, and close formally.",
			MinWords: 220,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Mheshimiwa Mgeni Rasmi, ni heshima kubwa kusimama mbele yenu siku hii adhimu.",
			Speaker:     "Zephyr",
			Translation: "Honourable Guest of Honour, it is a great honour to stand before you on this distinguished day.",
		},
	},

	// ---------------- FINAL · comprehensive ----------------
	"FINAL": {
		Listening: PaperListening{
			Title: "Kongamano: Kiswahili katika Karne ya 21 — Symposium: Swahili in the 21st Century",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "Teknolojia imekipa Kiswahili jukwaa jipya: mitandao ya kijamii, programu za simu na hata mifumo ya akili bandia sasa zinatumia lugha hii.", Translation: "Technology has given Swahili a new platform: social media, phone apps and even artificial-intelligence systems now use the language."},
				{Character: "Mira", Text: "Ni kweli, lakini istilahi nyingi za kisayansi bado zinakopwa moja kwa moja kutoka Kiingereza bila kufanyiwa uchunguzi wa kina.", Translation: "That's true, but many scientific terms are still borrowed directly from English without careful study."},
				{Character: "Professor Finch", Text: "Ndiyo maana mabaraza ya Kiswahili yana jukumu kubwa la kuunda istilahi zinazoeleweka na kukubalika, bila kuzuia ubunifu wa wazungumzaji wenyewe.", Translation: "That is why the Swahili councils have a great responsibility to create terms that are understood and accepted, without stifling speakers' own creativity."},
			},
			Questions: []PaperQuestion{
				{Question: "Teknolojia imekipa Kiswahili nini?", Options: []string{"Jukwaa jipya", "Matatizo pekee", "Wazungumzaji wachache", "Sarufi mpya"}, CorrectAnswer: "Jukwaa jipya"},
				{Question: "Mira anaibua tatizo gani?", Options: []string{"Istilahi za kisayansi zinakopwa bila uchunguzi wa kina", "Hakuna simu za kutosha", "Watu wameacha kuzungumza Kiswahili", "Mitandao ni ghali"}, CorrectAnswer: "Istilahi za kisayansi zinakopwa bila uchunguzi wa kina"},
				{Question: "Mabaraza ya Kiswahili yana jukumu gani?", Options: []string{"Kuunda istilahi zinazoeleweka na kukubalika", "Kupiga marufuku maneno ya kigeni yote", "Kuchapisha magazeti", "Kufundisha Kiingereza"}, CorrectAnswer: "Kuunda istilahi zinazoeleweka na kukubalika"},
				{Question: "Msemaji anaonya dhidi ya nini?", Options: []string{"Kuzuia ubunifu wa wazungumzaji", "Kutumia teknolojia", "Kujifunza lugha nyingine", "Kuandika mashairi"}, CorrectAnswer: "Kuzuia ubunifu wa wazungumzaji"},
			},
		},
		Reading: PaperReading{
			Title: "Tafsiri na Utamaduni — Translation and Culture",
			Paragraphs: []string{
				"Tafsiri si kubadilisha maneno ya lugha moja kwa ya lugha nyingine tu; ni kuhamisha dhana, hisia na muktadha wa kitamaduni.",
				"Kwa mfano, salamu 'Shikamoo' haiwezi kutafsiriwa kwa neno moja la Kiingereza, kwa sababu imebeba heshima kwa wazee ambayo ni msingi wa mahusiano katika jamii za Afrika Mashariki.",
				"Kwa hiyo mfasiri bora ni yule anayeelewa sio lugha mbili tu, bali tamaduni mbili, na anayejua lini kuhifadhi neno asilia na lini kulieleza.",
			},
			Questions: []PaperQuestion{
				{Question: "Kwa mujibu wa makala, tafsiri inahusu nini?", Options: []string{"Kuhamisha dhana, hisia na muktadha wa kitamaduni", "Kubadilisha maneno tu", "Kufupisha maandishi", "Kuandika upya hadithi"}, CorrectAnswer: "Kuhamisha dhana, hisia na muktadha wa kitamaduni"},
				{Question: "Kwa nini 'Shikamoo' ni vigumu kutafsiri?", Options: []string{"Imebeba heshima kwa wazee", "Ni neno refu mno", "Ni la Kiarabu", "Halina maana"}, CorrectAnswer: "Imebeba heshima kwa wazee"},
				{Question: "Mfasiri bora ni yupi?", Options: []string{"Anayeelewa tamaduni mbili", "Anayejua lugha nyingi zaidi", "Anayetumia kamusi tu", "Anayetafsiri haraka"}, CorrectAnswer: "Anayeelewa tamaduni mbili"},
				{Question: "Neno 'muktadha' lina maana gani?", Options: []string{"Mazingira yanayozunguka maana", "Kamusi", "Sarufi", "Matamshi"}, CorrectAnswer: "Mazingira yanayozunguka maana"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write an essay in Swahili on 'Nafasi ya Kiswahili katika kuunganisha Afrika'. Develop a nuanced argument, engage a counter-position, use appropriate methali and formal register, and draw a reasoned conclusion.",
			MinWords: 280,
		},
		Speaking: PaperSpeaking{
			Phrase:      "Lugha ni kioo cha utamaduni; tukiitunza, tunajitunza wenyewe.",
			Speaker:     "Lumora",
			Translation: "Language is the mirror of culture; if we care for it, we care for ourselves.",
		},
	},
}
