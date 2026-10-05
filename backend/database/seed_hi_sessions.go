package database

import (
	"gorm.io/gorm"

	"lumora/backend/models"
)

// Hindi listening and reading sessions — at least one of each per level.
// Transcripts carry a transliteration at A1–A2 and drop it from B1.

const (
	hiA1 = "A1 · शुरुआत — Foundations & Survival Hindi"
	hiA2 = "A2 · रोज़मर्रा की ज़िंदगी — Everyday Life"
	hiB1 = "B1 · बातचीत — Conversation & Narrative"
	hiB2 = "B2 · औपचारिक हिंदी — Formal Hindi & Society"
	hiC1 = "C1 · साहित्य और मुहावरे — Literature & Idioms"
	hiC2 = "C2 · प्रवीणता — Mastery"
)

func seedHindiListening(db *gorm.DB) {
	const hi = "hi"
	finch := "Professor Finch"

	addListeningL(db, hi, hiA1, "पहली मुलाक़ात — First Meeting",
		"Professor Finch meets a new student. Listen, then answer.", 1, 20,
		[]models.ListeningMatch{
			lm("नमस्ते", "hello"), lm("आपका नाम क्या है?", "What's your name?"),
			lm("आप कहाँ से हैं?", "Where are you from?"), lm("आपसे मिलकर ख़ुशी हुई", "Pleased to meet you"),
		},
		[]models.ListeningLine{
			ln(finch, "नमस्ते! आपका नाम क्या है?", "namaste! āpkā nām kyā hai? — Hello! What's your name?"),
			ln("Lumora", "नमस्ते जी। मेरा नाम लुमोरा है।", "namaste jī. merā nām lumorā hai. — Hello. My name is Lumora."),
			ln(finch, "आप कहाँ से हैं?", "āp kahā̃ se haĩ? — Where are you from?"),
			ln("Lumora", "मैं केन्या से हूँ, लेकिन अभी दिल्ली में रहती हूँ। मैं बीस साल की हूँ।", "mãi kenyā se hū̃, lekin abhī dillī mẽ rahtī hū̃. mãi bīs sāl kī hū̃. — I'm from Kenya, but I live in Delhi now. I'm twenty."),
			ln(finch, "आपसे मिलकर ख़ुशी हुई!", "āpse milkar khushī huī! — Pleased to meet you!"),
		},
		[]models.ListeningQuestion{
			lq("Where is Lumora from?", "Kenya", "Kenya", "India", "Nepal", "Delhi"),
			lq("Where does she live now?", "Delhi", "Delhi", "Mumbai", "Kenya", "Jaipur"),
			lq("How old is she?", "20", "20", "12", "25", "2"),
		},
	)
	addListeningL(db, hi, hiA1, "चाय की दुकान पर — At the Tea Stall",
		"Cora orders at a tea stall. Listen, then answer.", 2, 20,
		[]models.ListeningMatch{
			lm("भैया", "brother (friendly address)"), lm("समोसा", "samosa"),
			lm("कितने हुए?", "how much is it?"), lm("बीस रुपये", "twenty rupees"),
		},
		[]models.ListeningLine{
			ln("Cora", "भैया, दो चाय और एक समोसा दीजिए।", "bhaiyā, do chāy aur ek samosā dījie. — Brother, two teas and one samosa, please."),
			ln("Riko", "जी, अभी लाया। चीनी कम या ज़्यादा?", "jī, abhī lāyā. cīnī kam yā zyādā? — Coming right up. Less sugar or more?"),
			ln("Cora", "चीनी कम, कृपया। कितने हुए?", "cīnī kam, kṛpayā. kitne hue? — Less sugar, please. How much is it?"),
			ln("Riko", "चालीस रुपये।", "cālīs rupaye. — Forty rupees."),
		},
		[]models.ListeningQuestion{
			lq("How many teas does Cora order?", "Two", "Two", "One", "Three", "Four"),
			lq("What does she order to eat?", "A samosa", "A samosa", "Two samosas", "A biscuit", "Nothing"),
			lq("How does she want her tea?", "With less sugar", "With less sugar", "With more sugar", "With no milk", "Very hot"),
			lq("How much does she pay?", "₹40", "₹40", "₹14", "₹400", "₹20"),
		},
	)
	addListeningL(db, hi, hiA2, "जयपुर की यात्रा — A Trip to Jaipur",
		"Zephyr talks about a holiday. Listen, then answer.", 3, 22,
		[]models.ListeningMatch{
			lm("पिछले महीने", "last month"), lm("किला", "fort"),
			lm("गर्मी", "heat"), lm("फिर से", "again"),
		},
		[]models.ListeningLine{
			ln("Zephyr", "पिछले महीने मैं अपने भाई के साथ जयपुर गया। हम ट्रेन से गए।", "pichle mahīne mãi apne bhāī ke sāth jaypur gayā. ham ṭren se gae. — Last month I went to Jaipur with my brother. We went by train."),
			ln("Zephyr", "हमने आमेर का किला देखा और दाल-बाटी खाई। बहुत स्वादिष्ट थी!", "hamne āmer kā kilā dekhā aur dāl-bāṭī khāī. bahut svādiṣṭ thī! — We saw Amer Fort and ate dal-baati. It was delicious!"),
			ln("Zephyr", "लेकिन बहुत गर्मी थी, इसलिए हम दोपहर में होटल में आराम करते थे। सर्दियों में मैं फिर से जाऊँगा।", "lekin bahut garmī thī, islie ham dopahar mẽ hoṭal mẽ ārām karte the. sardiyõ mẽ mãi phir se jāū̃gā. — But it was very hot, so we rested in the hotel in the afternoons. I'll go again in winter."),
		},
		[]models.ListeningQuestion{
			lq("Who did he go with?", "His brother", "His brother", "His sister", "His friends", "His parents"),
			lq("How did they travel?", "By train", "By train", "By bus", "By plane", "By car"),
			lq("What did they eat?", "Dal-baati", "Dal-baati", "Biryani", "Samosas", "Dosa"),
			lq("When will he go again?", "In winter", "In winter", "Next month", "In summer", "Never"),
		},
	)
	addListeningL(db, hi, hiB1, "डॉक्टर के पास — At the Doctor's",
		"A patient visits the doctor. Listen, then answer.", 4, 24,
		[]models.ListeningMatch{
			lm("बुख़ार", "fever"), lm("खाँसी", "cough"),
			lm("दवाई", "medicine"), lm("आराम", "rest"),
		},
		[]models.ListeningLine{
			ln("Nana", "बैठिए। बताइए, क्या तकलीफ़ है?", "Please sit. Tell me, what's the trouble?"),
			ln("Blaze", "डॉक्टर साहब, तीन दिन से मुझे खाँसी है और कल रात से बुख़ार भी है।", "Doctor, I've had a cough for three days and since last night a fever too."),
			ln("Nana", "लगता है वायरल बुख़ार है। यह दवाई दिन में तीन बार खाने के बाद लीजिए।", "It seems to be a viral fever. Take this medicine three times a day after meals."),
			ln("Blaze", "क्या मैं कल दफ़्तर जा सकता हूँ?", "Can I go to the office tomorrow?"),
			ln("Nana", "नहीं, कम से कम दो दिन घर पर आराम कीजिए और ख़ूब पानी पीजिए।", "No, rest at home for at least two days and drink plenty of water."),
		},
		[]models.ListeningQuestion{
			lq("How long has he had a cough?", "Three days", "Three days", "One day", "A week", "Since this morning"),
			lq("What does the doctor think it is?", "A viral fever", "A viral fever", "Malaria", "Food poisoning", "Nothing serious at all"),
			lq("How should he take the medicine?", "Three times a day after meals", "Three times a day after meals", "Once at night", "Twice before meals", "Only when coughing"),
			lq("What is the advice about work?", "Rest at home for at least two days", "Rest at home for at least two days", "Go to work tomorrow", "Work half days", "Take a week off"),
		},
	)
	addListeningL(db, hi, hiB2, "समाचार बुलेटिन — The News Bulletin",
		"A radio news bulletin. Listen, then answer.", 5, 26,
		[]models.ListeningMatch{
			lm("सरकार", "government"), lm("किसान", "farmers"),
			lm("योजना", "scheme"), lm("बाढ़", "flood"),
		},
		[]models.ListeningLine{
			ln("Mira", "नमस्कार, ये हैं आज के मुख्य समाचार। केंद्र सरकार ने बाढ़ से प्रभावित किसानों के लिए एक नई योजना की घोषणा की है।", "Hello, here are today's main headlines. The central government has announced a new scheme for flood-affected farmers."),
			ln("Mira", "इस योजना के अंतर्गत लगभग दस लाख किसानों को मुफ़्त बीज और खाद दी जाएगी।", "Under this scheme, about one million farmers will be given free seeds and fertiliser."),
			ln("Mira", "इस बीच, मौसम विभाग ने अगले दो दिनों में भारी बारिश की चेतावनी जारी की है और लोगों से सावधान रहने की अपील की है।", "Meanwhile, the weather department has issued a warning of heavy rain over the next two days and urged people to stay alert."),
		},
		[]models.ListeningQuestion{
			lq("Who is the new scheme for?", "Flood-affected farmers", "Flood-affected farmers", "Students", "Teachers", "Factory workers"),
			lq("About how many farmers will benefit?", "One million (ten lakh)", "One million (ten lakh)", "Ten thousand", "One hundred thousand", "Ten million"),
			lq("What will they receive?", "Free seeds and fertiliser", "Free seeds and fertiliser", "Cash", "New tractors", "Land"),
			lq("What warning has been issued?", "Heavy rain in the next two days", "Heavy rain in the next two days", "A heatwave", "An earthquake", "A cyclone tomorrow"),
		},
	)
	addListeningL(db, hi, hiC1, "व्याख्यान: प्रेमचंद का यथार्थ — Lecture: Premchand's Realism",
		"An excerpt from a literature lecture. Listen, then answer.", 6, 28,
		[]models.ListeningMatch{
			lm("यथार्थ", "reality, realism"), lm("किसान", "farmer"),
			lm("समाज", "society"), lm("उपन्यास", "novel"),
		},
		[]models.ListeningLine{
			ln(finch, "प्रेमचंद से पहले हिंदी कथा-साहित्य में राजा-रानियों और जादू की कहानियाँ अधिक थीं।", "Before Premchand, Hindi fiction was dominated by tales of kings, queens and magic."),
			ln(finch, "प्रेमचंद ने पहली बार किसानों, मज़दूरों और स्त्रियों के जीवन को साहित्य के केंद्र में रखा।", "Premchand was the first to place the lives of farmers, labourers and women at the centre of literature."),
			ln(finch, "गोदान का होरी केवल एक पात्र नहीं, बल्कि पूरे भारतीय किसान वर्ग की पीड़ा का प्रतीक है। यही कारण है कि उनकी रचनाएँ आज भी प्रासंगिक हैं।", "Hori in Godaan is not just a character but a symbol of the suffering of the whole Indian peasantry. That is why his works remain relevant today."),
		},
		[]models.ListeningQuestion{
			lq("What dominated Hindi fiction before Premchand?", "Tales of kings, queens and magic", "Tales of kings, queens and magic", "Realist novels", "Poetry only", "Detective stories"),
			lq("What did Premchand put at the centre of literature?", "The lives of farmers, labourers and women", "The lives of farmers, labourers and women", "Royal courts", "City businessmen", "Foreign travel"),
			lq("Who is Hori?", "A character in Godaan who symbolises the peasantry", "A character in Godaan who symbolises the peasantry", "Premchand's son", "A king", "A poet"),
			lq("Why are his works still relevant, according to the lecturer?", "They express the suffering of a whole class", "They express the suffering of a whole class", "They are short", "They are funny", "They are about magic"),
		},
	)
	addListeningL(db, hi, hiC2, "मुंबई की गलियों में — On a Mumbai Street",
		"Two friends chat in Bambaiya Hindi. Listen, then answer.", 7, 30,
		[]models.ListeningMatch{
			lm("भिड़ू", "buddy (Bambaiya)"), lm("बिंदास", "carefree, cool"),
			lm("झकास", "awesome"), lm("लोकल", "local train"),
		},
		[]models.ListeningLine{
			ln("Blaze", "क्या भिड़ू, क्या बोलता है? आज शाम का क्या प्लान है?", "Hey buddy, what's up? What's the plan this evening?"),
			ln("Pip", "बिंदास! आज मरीन ड्राइव पे भुट्टा खाने चलते हैं।", "All cool! Let's go eat roasted corn at Marine Drive."),
			ln("Blaze", "झकास आइडिया! पर लोकल में बहुत भीड़ होगी, छह बजे के बाद निकलते हैं।", "Awesome idea! But the local train will be packed — let's leave after six."),
			ln("Pip", "ठीक है, स्टेशन पे मिलते हैं।", "OK, see you at the station."),
		},
		[]models.ListeningQuestion{
			lq("What do they plan to do?", "Eat roasted corn at Marine Drive", "Eat roasted corn at Marine Drive", "Watch a film", "Play cricket", "Go shopping"),
			lq("Why will they leave after six?", "The local train will be crowded", "The local train will be crowded", "Pip is working", "It's raining", "The shop opens then"),
			lq("Where will they meet?", "At the station", "At the station", "At Marine Drive", "At Pip's house", "At a café"),
			lq("Which variety is this?", "Bambaiya Hindi (Mumbai)", "Bambaiya Hindi (Mumbai)", "Bhojpuri", "Formal Hindi", "Urdu poetry"),
		},
	)
}

func seedHindiReading(db *gorm.DB) {
	const hi = "hi"

	addReadingL(db, hi, hiA1, "मेरा परिवार — My Family",
		"Riya introduces her family. Read, then answer.", 1, 20,
		[]models.ReadingLine{
			rl("मेरा नाम रिया है। मैं लखनऊ में रहती हूँ। (merā nām riyā hai. mãi lakhnaū mẽ rahtī hū̃.)", "My name is Riya. I live in Lucknow."),
			rl("मेरे परिवार में चार लोग हैं: पिता जी, माँ, मेरा छोटा भाई और मैं। (mere parivār mẽ cār log haĩ…)", "There are four people in my family: father, mother, my younger brother and me."),
			rl("मेरे पिता जी डॉक्टर हैं और मेरी माँ अध्यापिका हैं। (mere pitā jī ḍākṭar haĩ aur merī mā̃ adhyāpikā haĩ.)", "My father is a doctor and my mother is a teacher."),
			rl("मुझे गाना गाना और किताबें पढ़ना पसंद है। (mujhe gānā gānā aur kitābẽ paṛhnā pasand hai.)", "I like singing and reading books."),
		},
		[]models.ReadingQuestion{
			rq("Where does Riya live?", "Lucknow", "Lucknow", "Delhi", "Mumbai", "Jaipur"),
			rq("How many people are in her family?", "4", "4", "3", "5", "6"),
			rq("What is her mother's job?", "Teacher", "Teacher", "Doctor", "Engineer", "Shopkeeper"),
			rq("What does Riya like?", "Singing and reading books", "Singing and reading books", "Cooking", "Cricket", "Dancing"),
		},
	)
	addReadingL(db, hi, hiA1, "बाज़ार — The Market",
		"A shopping trip. Read, then answer.", 2, 20,
		[]models.ReadingLine{
			rl("आज मैं बाज़ार गई। (āj mãi bāzār gaī.)", "Today I went to the market."),
			rl("मैंने एक किलो टमाटर, आधा किलो प्याज़ और छह केले ख़रीदे। (mãine ek kilo ṭamāṭar, ādhā kilo pyāz aur chah kele kharīde.)", "I bought a kilo of tomatoes, half a kilo of onions and six bananas."),
			rl("सब मिलाकर सौ रुपये हुए। (sab milākar sau rupaye hue.)", "It came to a hundred rupees altogether."),
			rl("शाम को मैंने परिवार के लिए सब्ज़ी बनाई। (shām ko mãine parivār ke lie sabzī banāī.)", "In the evening I cooked vegetables for the family."),
		},
		[]models.ReadingQuestion{
			rq("How many bananas did she buy?", "Six", "Six", "One", "Two", "Ten"),
			rq("How much did she spend in total?", "₹100", "₹100", "₹10", "₹60", "₹1,000"),
			rq("What did she cook in the evening?", "Vegetables (sabzi)", "Vegetables (sabzi)", "Rice", "Sweets", "Tea"),
		},
	)
	addReadingL(db, hi, hiA2, "दोस्त को चिट्ठी — A Letter to a Friend",
		"Ankit writes to his friend. Read, then answer.", 3, 22,
		[]models.ReadingLine{
			rl("प्रिय राहुल, तुम कैसे हो? (priy rāhul, tum kaise ho?)", "Dear Rahul, how are you?"),
			rl("मैं पिछले हफ़्ते पुणे आया, क्योंकि मुझे यहाँ नई नौकरी मिली है। (mãi pichle hafte puṇe āyā, kyõki mujhe yahā̃ naī naukrī milī hai.)", "I came to Pune last week because I got a new job here."),
			rl("यहाँ का मौसम बहुत अच्छा है और लोग भी मिलनसार हैं। हर रविवार मैं पहाड़ों पर घूमने जाता हूँ। (…har ravivār mãi pahāṛõ par ghūmne jātā hū̃.)", "The weather here is very nice and the people are friendly. Every Sunday I go walking in the hills."),
			rl("दिवाली पर मैं घर आऊँगा। तब मिलेंगे! तुम्हारा दोस्त, अंकित", "I'll come home for Diwali. See you then! Your friend, Ankit"),
		},
		[]models.ReadingQuestion{
			rq("Why did Ankit move to Pune?", "He got a new job", "He got a new job", "To study", "To get married", "To visit Rahul"),
			rq("What does he do every Sunday?", "Goes walking in the hills", "Goes walking in the hills", "Works", "Visits family", "Watches films"),
			rq("When will he come home?", "For Diwali", "For Diwali", "Next week", "For Holi", "Next year"),
			rq("How does he describe the people in Pune?", "Friendly", "Friendly", "Rude", "Busy", "Quiet"),
		},
	)
	addReadingL(db, hi, hiB1, "प्यासा कौआ — The Thirsty Crow",
		"A folk tale. Read, then answer.", 4, 24,
		[]models.ReadingLine{
			rl("गर्मी का दिन था। एक कौआ बहुत प्यासा था।", "It was a summer day. A crow was very thirsty."),
			rl("उसे एक घड़ा दिखा, लेकिन उसमें पानी बहुत नीचे था।", "He saw a pot, but the water in it was very low."),
			rl("कौए ने आसपास से कंकड़ उठाकर एक-एक करके घड़े में डाले।", "The crow picked up pebbles from nearby and dropped them into the pot one by one."),
			rl("धीरे-धीरे पानी ऊपर आ गया और कौए ने पानी पीकर अपनी प्यास बुझाई। सीख: जहाँ चाह, वहाँ राह।", "Slowly the water rose, and the crow drank and quenched his thirst. Moral: where there's a will, there's a way."),
		},
		[]models.ReadingQuestion{
			rq("What was the problem?", "The water in the pot was too low", "The water in the pot was too low", "There was no pot", "The water was dirty", "The pot was broken"),
			rq("What did the crow do?", "Dropped pebbles into the pot", "Dropped pebbles into the pot", "Broke the pot", "Called other crows", "Flew away"),
			rq("What does 'उठाकर' mean here?", "having picked up", "having picked up", "while flying", "as soon as", "without lifting"),
			rq("What is the moral?", "Where there's a will, there's a way", "Where there's a will, there's a way", "Never trust strangers", "Slow and steady wins the race", "Share with others"),
		},
	)
	addReadingL(db, hi, hiB2, "डिजिटल भुगतान और बदलता भारत — Digital Payments and a Changing India",
		"A newspaper column. Read, then answer.", 5, 26,
		[]models.ReadingLine{
			rl("पिछले कुछ वर्षों में भारत में डिजिटल भुगतान का प्रयोग तेज़ी से बढ़ा है।", "In the last few years, the use of digital payments in India has grown rapidly."),
			rl("आज सब्ज़ी वाले से लेकर बड़े दुकानदार तक मोबाइल फ़ोन से भुगतान स्वीकार करते हैं।", "Today everyone from vegetable sellers to big shopkeepers accepts payment by mobile phone."),
			rl("इससे लेन-देन आसान और पारदर्शी हुआ है, किंतु बुज़ुर्गों और ग्रामीण क्षेत्रों के लोगों के लिए यह अब भी एक चुनौती है।", "This has made transactions easier and more transparent, but it is still a challenge for the elderly and people in rural areas."),
			rl("अतः आवश्यक है कि डिजिटल साक्षरता पर विशेष ध्यान दिया जाए।", "It is therefore necessary that special attention be paid to digital literacy."),
		},
		[]models.ReadingQuestion{
			rq("What has grown rapidly in India?", "Digital payments", "Digital payments", "Cash use", "Bank branches", "Credit cards only"),
			rq("Who accepts mobile payments today?", "Everyone from vegetable sellers to big shopkeepers", "Everyone from vegetable sellers to big shopkeepers", "Only big shops", "Only banks", "Only in cities"),
			rq("For whom is it still a challenge?", "The elderly and rural people", "The elderly and rural people", "Students", "Shopkeepers", "Bankers"),
			rq("What does the writer recommend?", "Special attention to digital literacy", "Special attention to digital literacy", "Banning cash", "Closing banks", "Using only cash"),
		},
	)
	addReadingL(db, hi, hiC1, "कबीर के दोहे — Couplets of Kabir",
		"Two classic couplets with explanations. Read, then answer.", 6, 28,
		[]models.ReadingLine{
			rl("बुरा जो देखन मैं चला, बुरा न मिलिया कोय। जो दिल खोजा आपना, मुझसे बुरा न कोय॥", "I set out to find the wicked and found no one; when I searched my own heart, there was no one worse than me."),
			rl("(भावार्थ) दूसरों में दोष ढूँढने से पहले हमें अपने भीतर झाँकना चाहिए।", "(Meaning) Before looking for faults in others, we should look within ourselves."),
			rl("धीरे-धीरे रे मना, धीरे सब कुछ होय। माली सींचे सौ घड़ा, ऋतु आए फल होय॥", "Slowly, O mind, everything happens slowly; the gardener may pour a hundred pots of water, but fruit comes only in its season."),
			rl("(भावार्थ) धैर्य रखना चाहिए; हर काम अपने समय पर ही पूरा होता है।", "(Meaning) One must be patient; every task is completed in its own time."),
		},
		[]models.ReadingQuestion{
			rq("What does the first couplet teach?", "Look at your own faults before others'", "Look at your own faults before others'", "Avoid wicked people", "Search for friends", "Trust your heart always"),
			rq("What does the gardener image in the second couplet show?", "Things happen in their own time", "Things happen in their own time", "Hard work is useless", "Gardens need lots of water", "Fruit is expensive"),
			rq("What does 'भावार्थ' mean?", "The meaning / gist", "The meaning / gist", "The poet's name", "The rhyme", "The title"),
			rq("What poetic form are these?", "Dohas (couplets)", "Dohas (couplets)", "Ghazals", "Chaupais", "Sonnets"),
		},
	)
	addReadingL(db, hi, hiC2, "हिंदी: एक जीवंत भाषा — Hindi: A Living Language",
		"An essay on Hindi today. Read, then answer.", 7, 30,
		[]models.ReadingLine{
			rl("हिंदी विश्व की सबसे अधिक बोली जाने वाली भाषाओं में से एक है और प्रवासी भारतीयों के माध्यम से इसका प्रसार अनेक देशों में हुआ है।", "Hindi is among the most widely spoken languages in the world and has spread to many countries through the Indian diaspora."),
			rl("हर वर्ष 14 सितंबर को हिंदी दिवस मनाया जाता है, क्योंकि इसी दिन 1949 में संविधान सभा ने हिंदी को राजभाषा के रूप में अपनाया था।", "Every year Hindi Day is celebrated on 14 September, because on that day in 1949 the Constituent Assembly adopted Hindi as the official language."),
			rl("आज हिंदी के सामने चुनौती यह है कि वह अंग्रेज़ी के बढ़ते प्रभाव के बीच अपनी सहजता और समृद्धि दोनों को कैसे बनाए रखे।", "Today the challenge for Hindi is how to keep both its naturalness and its richness amid the growing influence of English."),
			rl("इसका उत्तर भाषा को जड़ बनाने में नहीं, बल्कि उसे नए विचारों और विज्ञान की भाषा बनाने में है।", "The answer lies not in making the language rigid, but in making it a language of new ideas and of science."),
		},
		[]models.ReadingQuestion{
			rq("When is Hindi Day (हिंदी दिवस) celebrated?", "14 September", "14 September", "26 January", "15 August", "2 October"),
			rq("Why that date?", "Hindi was adopted as the official language on that day in 1949", "Hindi was adopted as the official language on that day in 1949", "It is a poet's birthday", "It is Independence Day", "The first Hindi film was released"),
			rq("What challenge does Hindi face, according to the essay?", "Keeping its naturalness and richness amid English influence", "Keeping its naturalness and richness amid English influence", "Losing all its speakers", "Changing its script", "Being banned in schools"),
			rq("What solution does the writer propose?", "Making Hindi a language of new ideas and science", "Making Hindi a language of new ideas and science", "Freezing the language as it is", "Replacing it with English", "Using only Sanskrit words"),
		},
	)
}
