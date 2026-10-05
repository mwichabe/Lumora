package database

import (
	"gorm.io/gorm"

	"lumora/backend/models"
)

// Swahili listening and reading sessions — at least one of each per level.

const (
	swA1 = "A1 · Mwanzo — Survival Swahili"
	swA2 = "A2 · Maisha ya Kila Siku — Everyday Life"
	swB1 = "B1 · Mazungumzo — Conversation & Narrative"
	swB2 = "B2 · Lugha Rasmi — Formal Swahili & Society"
	swC1 = "C1 · Fasihi na Methali — Literature & Proverbs"
	swC2 = "C2 · Umahiri — Mastery"
)

func seedSwahiliListening(db *gorm.DB) {
	const sw = "sw"
	finch := "Professor Finch"

	addListeningL(db, sw, swA1, "Darasa la Kwanza — First Day of Class",
		"Professor Finch meets a new student. Listen, then answer.", 1, 20,
		[]models.ListeningMatch{
			lm("Hujambo?", "How are you?"), lm("Jina lako nani?", "What's your name?"),
			lm("Unatoka wapi?", "Where are you from?"), lm("Nimefurahi kukujua", "Pleased to meet you"),
		},
		[]models.ListeningLine{
			ln(finch, "Hujambo? Karibu darasani. Jina lako nani?", "How are you? Welcome to class. What's your name?"),
			ln("Lumora", "Sijambo, mwalimu. Jina langu ni Lumora.", "I'm fine, teacher. My name is Lumora."),
			ln(finch, "Unatoka wapi, Lumora?", "Where are you from, Lumora?"),
			ln("Lumora", "Ninatoka Kenya, lakini ninaishi Dar es Salaam. Nina miaka ishirini.", "I'm from Kenya, but I live in Dar es Salaam. I'm twenty years old."),
			ln(finch, "Vizuri sana. Nimefurahi kukujua!", "Very good. Pleased to meet you!"),
		},
		[]models.ListeningQuestion{
			lq("Where is Lumora from?", "Kenya", "Kenya", "Tanzania", "Uganda", "Rwanda"),
			lq("Where does she live?", "Dar es Salaam", "Dar es Salaam", "Nairobi", "Mombasa", "Arusha"),
			lq("How old is she?", "20", "20", "12", "30", "22"),
		},
	)
	addListeningL(db, sw, swA1, "Sokoni — At the Market",
		"Cora buys fruit and bargains. Listen, then answer.", 2, 20,
		[]models.ListeningMatch{
			lm("Bei gani?", "How much?"), lm("ghali", "expensive"),
			lm("punguza", "reduce (the price)"), lm("embe", "mango"),
		},
		[]models.ListeningLine{
			ln("Riko", "Karibu mama! Embe tamu sana leo.", "Welcome, madam! Very sweet mangoes today."),
			ln("Cora", "Asante. Embe moja bei gani?", "Thanks. How much is one mango?"),
			ln("Riko", "Shilingi hamsini tu.", "Just fifty shillings."),
			ln("Cora", "Ghali sana! Nitanunua matano kwa shilingi mia mbili.", "Very expensive! I'll buy five for two hundred shillings."),
			ln("Riko", "Haya, sawa. Chukua!", "Alright, OK. Take them!"),
		},
		[]models.ListeningQuestion{
			lq("What is Cora buying?", "Mangoes", "Mangoes", "Bananas", "Oranges", "Tomatoes"),
			lq("What was the first price for one?", "50 shillings", "50 shillings", "15 shillings", "500 shillings", "200 shillings"),
			lq("How many does she buy in the end?", "Five", "Five", "One", "Two", "Ten"),
			lq("What does she pay in total?", "200 shillings", "200 shillings", "250 shillings", "50 shillings", "100 shillings"),
		},
	)
	addListeningL(db, sw, swA2, "Safari ya Mombasa — A Trip to Mombasa",
		"Zephyr talks about a holiday. Listen, then answer.", 3, 22,
		[]models.ListeningMatch{
			lm("treni", "train"), lm("pwani", "coast"),
			lm("tulikula", "we ate"), lm("tutarudi", "we will return"),
		},
		[]models.ListeningLine{
			ln("Zephyr", "Mwezi uliopita nilisafiri kwenda Mombasa na familia yangu.", "Last month I travelled to Mombasa with my family."),
			ln("Zephyr", "Tulipanda treni kutoka Nairobi asubuhi, tukafika saa sita mchana.", "We took the train from Nairobi in the morning and arrived at noon."),
			ln("Zephyr", "Tulikaa siku nne. Tuliogelea baharini na tulikula samaki na wali wa nazi.", "We stayed four days. We swam in the sea and ate fish and coconut rice."),
			ln("Zephyr", "Watoto walifurahi sana. Mwaka ujao tutarudi tena!", "The children were very happy. Next year we'll go back again!"),
		},
		[]models.ListeningQuestion{
			lq("How did they travel to Mombasa?", "By train", "By train", "By bus", "By plane", "By car"),
			lq("When did they arrive?", "At noon", "At noon", "In the evening", "At night", "In the morning"),
			lq("How long did they stay?", "Four days", "Four days", "One week", "Two days", "A month"),
			lq("What did they eat?", "Fish and coconut rice", "Fish and coconut rice", "Ugali and meat", "Chips and chicken", "Beans and bread"),
		},
	)
	addListeningL(db, sw, swB1, "Kwa Daktari — At the Clinic",
		"A patient visits the doctor. Listen, then answer.", 4, 24,
		[]models.ListeningMatch{
			lm("homa", "fever"), lm("kikohozi", "cough"),
			lm("dawa", "medicine"), lm("pumzika", "rest"),
		},
		[]models.ListeningLine{
			ln("Nana", "Karibu. Una tatizo gani leo?", "Welcome. What's the problem today?"),
			ln("Blaze", "Daktari, nimekuwa na homa na kikohozi tangu juzi, na kichwa kinaniuma sana.", "Doctor, I've had a fever and a cough since the day before yesterday, and my head hurts a lot."),
			ln("Nana", "Pole sana. Inaonekana ni mafua. Nitakupa dawa; kunywa vidonge viwili mara tatu kwa siku baada ya chakula.", "I'm sorry. It looks like flu. I'll give you medicine; take two tablets three times a day after meals."),
			ln("Blaze", "Je, naweza kwenda kazini kesho?", "Can I go to work tomorrow?"),
			ln("Nana", "Hapana. Pumzika nyumbani siku mbili na unywe maji mengi.", "No. Rest at home for two days and drink lots of water."),
		},
		[]models.ListeningQuestion{
			lq("Since when has he been ill?", "The day before yesterday", "The day before yesterday", "Last week", "This morning", "Yesterday evening"),
			lq("What does the doctor think it is?", "Flu", "Flu", "Malaria", "A broken arm", "Food poisoning"),
			lq("How should he take the medicine?", "Two tablets three times a day after meals", "Two tablets three times a day after meals", "One tablet at night", "Three tablets before breakfast", "Only when in pain"),
			lq("What is the advice about work?", "Rest at home for two days", "Rest at home for two days", "Go to work as usual", "Work half days", "Rest for a week"),
		},
	)
	addListeningL(db, sw, swB2, "Taarifa ya Habari — The News Bulletin",
		"A radio news bulletin. Listen, then answer.", 5, 26,
		[]models.ListeningMatch{
			lm("serikali", "government"), lm("mafuriko", "floods"),
			lm("wananchi", "citizens"), lm("msaada", "aid, help"),
		},
		[]models.ListeningLine{
			ln("Mira", "Hii ni taarifa ya habari kutoka Dodoma. Serikali imetangaza hatua mpya za kusaidia wakulima walioathiriwa na mafuriko.", "This is the news from Dodoma. The government has announced new measures to help farmers affected by floods."),
			ln("Mira", "Waziri wa Kilimo amesema kwamba wakulima zaidi ya elfu kumi watapewa mbegu na mbolea bure.", "The Minister of Agriculture said more than ten thousand farmers will be given free seeds and fertiliser."),
			ln("Mira", "Aidha, imeripotiwa kwamba barabara kadhaa zimeharibika, na wananchi wametakiwa kuwa waangalifu wanaposafiri.", "Furthermore, it has been reported that several roads have been damaged, and citizens have been asked to be careful when travelling."),
		},
		[]models.ListeningQuestion{
			lq("Where is the bulletin from?", "Dodoma", "Dodoma", "Nairobi", "Kampala", "Zanzibar"),
			lq("Who is the government helping?", "Farmers affected by floods", "Farmers affected by floods", "Teachers on strike", "Tourists", "Fishermen affected by drought"),
			lq("What will the farmers receive?", "Free seeds and fertiliser", "Free seeds and fertiliser", "Cash payments", "New tractors", "Land"),
			lq("What are citizens asked to do?", "Be careful when travelling", "Be careful when travelling", "Stay at home", "Vote", "Plant trees"),
		},
	)
	addListeningL(db, sw, swC1, "Mjadala: Elimu na Teknolojia — Debate: Education and Technology",
		"A school debate. Listen, then answer.", 6, 28,
		[]models.ListeningMatch{
			lm("hoja", "argument"), lm("kwa upande mwingine", "on the other hand"),
			lm("hata hivyo", "nevertheless"), lm("kwa kuhitimisha", "in conclusion"),
		},
		[]models.ListeningLine{
			ln(finch, "Mwenyekiti, waheshimiwa majaji, ninapinga hoja kwamba teknolojia imeharibu elimu.", "Chairperson, honourable judges, I oppose the motion that technology has ruined education."),
			ln(finch, "Kwanza, mtandao umewawezesha wanafunzi wa vijijini kupata vitabu na mihadhara ambayo hawangeipata kamwe.", "First, the internet has enabled rural students to access books and lectures they would never have obtained."),
			ln(finch, "Kwa upande mwingine, ni kweli kwamba baadhi ya wanafunzi hutumia simu vibaya. Hata hivyo, kosa si la teknolojia bali ni la matumizi.", "On the other hand, it is true that some students misuse phones. Nevertheless, the fault is not technology's but its use."),
			ln(finch, "Kwa kuhitimisha, teknolojia ni chombo; sisi ndio tunaoamua kama itajenga au itabomoa.", "In conclusion, technology is a tool; it is we who decide whether it builds or destroys."),
		},
		[]models.ListeningQuestion{
			lq("What is the speaker's position?", "Against the motion that technology ruined education", "Against the motion that technology ruined education", "For the motion", "Undecided", "Against using the internet"),
			lq("What is the first argument?", "The internet gives rural students access to books and lectures", "The internet gives rural students access to books and lectures", "Phones are cheap", "Teachers prefer computers", "Exams are online"),
			lq("What concession does the speaker make?", "Some students misuse phones", "Some students misuse phones", "Technology is expensive", "Books are better", "Teachers are lazy"),
			lq("What is the conclusion?", "Technology is a tool; how we use it matters", "Technology is a tool; how we use it matters", "Technology should be banned", "Schools should close", "Phones are dangerous"),
		},
	)
	addListeningL(db, sw, swC2, "Mtaani Nairobi — On a Nairobi Street",
		"Two young friends chat in Sheng. Listen, then answer.", 7, 30,
		[]models.ListeningMatch{
			lm("Niaje?", "What's up?"), lm("keja", "home"),
			lm("doh", "money"), lm("mathree", "matatu"),
		},
		[]models.ListeningLine{
			ln("Blaze", "Niaje msee! Form ni gani leo jioni?", "What's up, man! What's the plan this evening?"),
			ln("Pip", "Poa sana! Mbogi inakuja kejani kucheki mechi.", "All good! The crew is coming over to my place to watch the match."),
			ln("Blaze", "Fiti! Lakini doh imenishinda, sina hata ya mathree.", "Great! But I'm broke — I don't even have matatu fare."),
			ln("Pip", "Usijali buda, nitakutumia kidogo kwa M-Pesa.", "Don't worry, bro, I'll send you a bit by M-Pesa."),
		},
		[]models.ListeningQuestion{
			lq("What is the plan for the evening?", "Watching a match at Pip's place", "Watching a match at Pip's place", "Going to a wedding", "Studying for exams", "Visiting a grandmother"),
			lq("What is Blaze's problem?", "He has no money", "He has no money", "He is sick", "He is lost", "He has to work"),
			lq("How will Pip help?", "Send him money by M-Pesa", "Send him money by M-Pesa", "Pick him up by car", "Lend him a phone", "Cook for him"),
			lq("What kind of language is this?", "Sheng (Nairobi urban slang)", "Sheng (Nairobi urban slang)", "Classical poetry", "Formal Standard Swahili", "Kiamu dialect"),
		},
	)
}

func seedSwahiliReading(db *gorm.DB) {
	const sw = "sw"

	addReadingL(db, sw, swA1, "Familia Yangu — My Family",
		"Amina introduces her family. Read, then answer.", 1, 20,
		[]models.ReadingLine{
			rl("Jina langu ni Amina. Ninaishi Arusha, Tanzania.", "My name is Amina. I live in Arusha, Tanzania."),
			rl("Familia yangu ina watu watano: baba, mama, kaka yangu, dada yangu na mimi.", "My family has five people: dad, mum, my brother, my sister and me."),
			rl("Baba yangu ni mkulima na mama yangu ni mwalimu.", "My father is a farmer and my mother is a teacher."),
			rl("Mimi ni mwanafunzi. Ninapenda kusoma vitabu na kucheza mpira.", "I'm a student. I like reading books and playing football."),
		},
		[]models.ReadingQuestion{
			rq("Where does Amina live?", "Arusha", "Arusha", "Nairobi", "Mombasa", "Dodoma"),
			rq("How many people are in her family?", "5", "5", "4", "6", "3"),
			rq("What is her father's job?", "Farmer", "Farmer", "Teacher", "Doctor", "Driver"),
			rq("What does Amina like?", "Reading books and playing football", "Reading books and playing football", "Cooking", "Singing", "Swimming"),
		},
	)
	addReadingL(db, sw, swA1, "Siku ya Juma — Juma's Day",
		"A day in Juma's life. Read, then answer.", 2, 20,
		[]models.ReadingLine{
			rl("Juma anaamka saa kumi na moja alfajiri.", "Juma wakes up at 5 a.m. (dawn)."),
			rl("Saa moja asubuhi anakunywa chai na anakula mandazi.", "At 7 a.m. he drinks tea and eats mandazi."),
			rl("Saa mbili anakwenda kazini kwa daladala.", "At 8 a.m. he goes to work by daladala."),
			rl("Jioni anarudi nyumbani na anapika wali na maharagwe.", "In the evening he returns home and cooks rice and beans."),
		},
		[]models.ReadingQuestion{
			rq("What time does Juma wake up (English time)?", "5 a.m.", "5 a.m.", "11 a.m.", "7 a.m.", "1 a.m."),
			rq("What does he eat for breakfast?", "Mandazi", "Mandazi", "Ugali", "Rice", "Bread and eggs"),
			rq("How does he get to work?", "By daladala", "By daladala", "On foot", "By bicycle", "By train"),
			rq("What does he cook in the evening?", "Rice and beans", "Rice and beans", "Fish", "Ugali and meat", "Chapati"),
		},
	)
	addReadingL(db, sw, swA2, "Barua kwa Rafiki — A Letter to a Friend",
		"Neema writes to her friend. Read, then answer.", 3, 22,
		[]models.ReadingLine{
			rl("Mpendwa Rehema, habari za siku nyingi?", "Dear Rehema, how have you been? (news of many days?)"),
			rl("Mimi ni mzima. Mwezi uliopita nilihamia Mwanza kwa sababu nimepata kazi mpya benki.", "I'm well. Last month I moved to Mwanza because I got a new job at a bank."),
			rl("Mji huu ni mzuri sana; uko karibu na Ziwa Victoria. Kila Jumamosi ninakwenda ziwani kupumzika.", "This town is very nice; it's near Lake Victoria. Every Saturday I go to the lake to relax."),
			rl("Desemba nitakuja Dar es Salaam. Tutaonana! Wako akupendaye, Neema.", "In December I'll come to Dar es Salaam. See you then! Your loving friend, Neema."),
		},
		[]models.ReadingQuestion{
			rq("Why did Neema move to Mwanza?", "She got a new job at a bank", "She got a new job at a bank", "To study", "To get married", "To visit family"),
			rq("What is Mwanza near?", "Lake Victoria", "Lake Victoria", "The Indian Ocean", "Mount Kilimanjaro", "Lake Tanganyika"),
			rq("What does she do every Saturday?", "Goes to the lake to relax", "Goes to the lake to relax", "Works at the bank", "Visits Rehema", "Goes to the market"),
			rq("When will she come to Dar es Salaam?", "In December", "In December", "Next week", "In January", "Next year"),
		},
	)
	addReadingL(db, sw, swB1, "Sungura na Fisi — The Hare and the Hyena",
		"A folk tale. Read, then answer.", 4, 24,
		[]models.ReadingLine{
			rl("Hapo zamani za kale, Sungura na Fisi walikuwa marafiki.", "Once upon a time, Hare and Hyena were friends."),
			rl("Siku moja walikubaliana kulima shamba pamoja na kugawana mavuno.", "One day they agreed to farm a field together and share the harvest."),
			rl("Fisi alikuwa mvivu; alilala chini ya mti wakati Sungura akilima peke yake.", "Hyena was lazy; he slept under a tree while Hare farmed alone."),
			rl("Mavuno yalipofika, Fisi alitaka sehemu kubwa. Sungura akamwambia: 'Asiyefanya kazi asile.' Wanyama wote wakakubaliana na Sungura.", "When the harvest came, Hyena wanted the bigger share. Hare told him: 'Whoever doesn't work shouldn't eat.' All the animals agreed with Hare."),
		},
		[]models.ReadingQuestion{
			rq("What did Hare and Hyena agree to do?", "Farm a field together and share the harvest", "Farm a field together and share the harvest", "Build a house", "Go hunting", "Race each other"),
			rq("What did Hyena do while Hare worked?", "Slept under a tree", "Slept under a tree", "Helped him", "Went to the market", "Cooked food"),
			rq("What did Hyena want at harvest time?", "The bigger share", "The bigger share", "Nothing", "An equal share", "To leave"),
			rq("What is the moral?", "Whoever doesn't work shouldn't eat", "Whoever doesn't work shouldn't eat", "Friends should share everything equally", "Sleep is important", "Hares are lazy"),
		},
	)
	addReadingL(db, sw, swB2, "Simu za Mkononi na Maisha ya Kisasa — Mobile Phones and Modern Life",
		"A newspaper column. Read, then answer.", 5, 26,
		[]models.ReadingLine{
			rl("Katika miaka ishirini iliyopita, simu za mkononi zimebadilisha maisha ya Waafrika Mashariki kwa kiasi kikubwa.", "Over the past twenty years, mobile phones have changed East Africans' lives enormously."),
			rl("Huduma za kutuma pesa kwa simu, kama M-Pesa, zimewawezesha wakulima na wafanyabiashara wadogo kupata huduma za kifedha bila kuwa na akaunti ya benki.", "Mobile money services such as M-Pesa have enabled farmers and small traders to access financial services without a bank account."),
			rl("Hata hivyo, wataalamu wanaonya kwamba matumizi yasiyo na kiasi yanaweza kuathiri afya ya akili na mahusiano ya kifamilia.", "However, experts warn that excessive use can affect mental health and family relationships."),
			rl("Kwa hiyo, changamoto iliyopo si kuacha teknolojia, bali kuitumia kwa busara.", "So the challenge is not to abandon technology but to use it wisely."),
		},
		[]models.ReadingQuestion{
			rq("What have mobile money services made possible?", "Financial services without a bank account", "Financial services without a bank account", "Free phone calls", "Higher crop prices", "Cheaper school fees"),
			rq("Who has benefited, according to the column?", "Farmers and small traders", "Farmers and small traders", "Only large banks", "Tourists", "Politicians"),
			rq("What do experts warn about?", "Excessive use can harm mental health and family relationships", "Excessive use can harm mental health and family relationships", "Phones are too expensive", "M-Pesa will close", "Phones cause floods"),
			rq("What is the writer's conclusion?", "Use technology wisely rather than abandon it", "Use technology wisely rather than abandon it", "Ban mobile phones", "Return to bank accounts only", "Technology has no benefits"),
		},
	)
	addReadingL(db, sw, swC1, "Shairi: Kiswahili — A Poem: Swahili",
		"A traditional-style shairi in praise of the language. Read, then answer.", 6, 28,
		[]models.ReadingLine{
			rl("Kiswahili ni lugha, ya Afrika Mashariki,", "Swahili is the language of East Africa,"),
			rl("Kimeenea kwa wingi, mipakani hakisiki,", "it has spread far and wide; borders do not stop it,"),
			rl("Wazee kwa vijana, kinawaweka rafiki,", "old and young alike, it keeps them as friends,"),
			rl("Tukienzi kwa bidii, lugha yetu ya haki.", "let us honour it with diligence, our rightful language."),
		},
		[]models.ReadingQuestion{
			rq("What is the poem's main message?", "Praise of Swahili and a call to honour it", "Praise of Swahili and a call to honour it", "A love story", "A warning about borders", "A description of a journey"),
			rq("Which rhyme (kina) ends most lines?", "-ki", "-ki", "-gha", "-na", "-li"),
			rq("What does the poem say borders do?", "They don't stop the language", "They don't stop the language", "They divide speakers", "They create new dialects", "They protect the language"),
			rq("What does 'tukienzi' ask us to do?", "Honour it", "Honour it", "Forget it", "Change it", "Translate it"),
		},
	)
	addReadingL(db, sw, swC2, "Kiswahili kama Lugha ya Afrika — Swahili as a Language of Africa",
		"An essay on Swahili's continental role. Read, then answer.", 7, 30,
		[]models.ReadingLine{
			rl("Kiswahili kimevuka mipaka ya pwani ya Afrika Mashariki na sasa ni miongoni mwa lugha zinazozungumzwa na watu wengi zaidi barani Afrika.", "Swahili has crossed the borders of the East African coast and is now among the most widely spoken languages on the African continent."),
			rl("Umoja wa Afrika umekitambua kuwa mojawapo ya lugha zake rasmi za kazi, na UNESCO imetenga tarehe saba Julai kuwa Siku ya Kiswahili Duniani.", "The African Union has recognised it as one of its official working languages, and UNESCO has designated 7 July as World Kiswahili Language Day."),
			rl("Hata hivyo, ukuaji huu unaleta swali muhimu: je, Kiswahili kitaweza kuhifadhi ufasaha na utajiri wake wa kifasihi huku kikienea kwa kasi kubwa?", "Yet this growth raises an important question: will Swahili be able to preserve its eloquence and literary richness while spreading so fast?"),
			rl("Jibu linategemea jinsi tutakavyowekeza katika elimu, uchapishaji na tafsiri, ili lugha hii ibaki chombo cha fikra pevu na si cha mawasiliano ya juu juu tu.", "The answer depends on how we invest in education, publishing and translation, so that the language remains a vehicle of mature thought and not merely of superficial communication."),
		},
		[]models.ReadingQuestion{
			rq("Which date is World Kiswahili Language Day?", "7 July", "7 July", "7 June", "1 May", "12 December"),
			rq("Which body uses Swahili as an official working language?", "The African Union", "The African Union", "The European Union", "NATO", "The Arab League"),
			rq("What concern does the writer raise?", "Whether Swahili can keep its eloquence and literary richness as it spreads", "Whether Swahili can keep its eloquence and literary richness as it spreads", "Whether people will stop speaking it", "Whether it should replace English", "Whether it has enough speakers"),
			rq("What does the answer depend on, according to the writer?", "Investment in education, publishing and translation", "Investment in education, publishing and translation", "Government bans on other languages", "Tourism", "Social media alone"),
		},
	)
}
