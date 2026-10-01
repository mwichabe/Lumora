package controllers

// mandarinPapers is the hand-authored Mandarin proficiency-exam bank, modelled
// on HSK 3.0 and keyed by the app's CEFR levels:
//
//	A1 = HSK 1 · A2 = HSK 2 · B1 = HSK 3 · B2 = HSK 4 · C1 = HSK 5 · C2 = HSK 6
//	FINAL = HSK 7–9 (the shared advanced band)
//
// Listening and reading come first (as in the HSK); writing is scored in
// Chinese characters, so MinWords here is a character count — the clients
// count each Han character as one unit. Questions are in English up to HSK 4
// and in Chinese from HSK 5. Reading texts carry Pinyin at HSK 1–2.
var mandarinPapers = map[string]paperContent{

	// ---------------- A1 · HSK 1 ----------------
	"A1": {
		Listening: PaperListening{
			Title: "小明的一天 — Xiao Ming's Day",
			Lines: []PaperLine{
				{Character: "Cora", Text: "你好！我叫小明，我是学生。我今年十九岁，我家在上海。", Translation: "Hello! My name is Xiao Ming and I'm a student. I'm nineteen this year; my home is in Shanghai."},
				{Character: "Cora", Text: "我每天七点起床，八点去学校。中午我在学校吃饭，我喜欢吃米饭和鱼。", Translation: "I get up at seven every day and go to school at eight. At noon I eat at school; I like rice and fish."},
				{Character: "Cora", Text: "下午五点我回家。晚上我和妈妈一起看电视。", Translation: "At five in the afternoon I go home. In the evening I watch TV with my mum."},
			},
			Questions: []PaperQuestion{
				{Question: "Where is Xiao Ming's home?", Options: []string{"Shanghai", "Beijing", "Hong Kong", "Nanjing"}, CorrectAnswer: "Shanghai"},
				{Question: "How old is he?", Options: []string{"19", "9", "17", "29"}, CorrectAnswer: "19"},
				{Question: "When does he go to school?", Options: []string{"8 o'clock", "7 o'clock", "9 o'clock", "5 o'clock"}, CorrectAnswer: "8 o'clock"},
				{Question: "What does he like to eat?", Options: []string{"Rice and fish", "Noodles", "Dumplings", "Bread"}, CorrectAnswer: "Rice and fish"},
			},
		},
		Reading: PaperReading{
			Title: "我的朋友 — My Friend",
			Paragraphs: []string{
				"我有一个中国朋友，她叫李月。(Wǒ yǒu yí gè Zhōngguó péngyou, tā jiào Lǐ Yuè.)",
				"李月是老师，她在北京工作。她有一只猫，很可爱。(Lǐ Yuè shì lǎoshī, tā zài Běijīng gōngzuò. Tā yǒu yì zhī māo, hěn kě'ài.)",
				"星期六我们一起去商店买东西，然后去饭馆吃饺子。(Xīngqīliù wǒmen yìqǐ qù shāngdiàn mǎi dōngxi, ránhòu qù fànguǎn chī jiǎozi.)",
			},
			Questions: []PaperQuestion{
				{Question: "What is the friend's job?", Options: []string{"Teacher", "Doctor", "Student", "Waiter"}, CorrectAnswer: "Teacher"},
				{Question: "Where does she work?", Options: []string{"Beijing", "Shanghai", "Kenya", "Hong Kong"}, CorrectAnswer: "Beijing"},
				{Question: "What pet does she have?", Options: []string{"A cat", "A dog", "A fish", "A bird"}, CorrectAnswer: "A cat"},
				{Question: "What do they eat on Saturday?", Options: []string{"Dumplings", "Noodles", "Rice", "Fish"}, CorrectAnswer: "Dumplings"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write a short self-introduction in Chinese characters: your name, nationality, age and one thing you like.",
			MinWords: 25,
		},
		Speaking: PaperSpeaking{
			Phrase:      "你好！我叫安娜，我是肯尼亚人。",
			Speaker:     "Lumora",
			Translation: "Nǐ hǎo! Wǒ jiào Ānnà, wǒ shì Kěnníyà rén. — Hello! My name is Anna; I'm Kenyan.",
		},
	},

	// ---------------- A2 · HSK 2 ----------------
	"A2": {
		Listening: PaperListening{
			Title: "去旅游 — A Trip to Xi'an",
			Lines: []PaperLine{
				{Character: "Zephyr", Text: "上个月我和朋友去西安旅游了。我们是坐高铁去的，比坐飞机便宜多了。", Translation: "Last month I travelled to Xi'an with friends. We went by high-speed train — much cheaper than flying."},
				{Character: "Zephyr", Text: "我们看了兵马俑，还吃了很多好吃的东西。我以前没去过西安，觉得那儿特别有意思。", Translation: "We saw the Terracotta Warriors and ate lots of good food. I'd never been to Xi'an before and found it fascinating."},
				{Character: "Zephyr", Text: "虽然天气很热，但是大家都很开心。我们打算明年去成都看熊猫。", Translation: "Although it was hot, everyone was happy. We plan to go to Chengdu to see the pandas next year."},
			},
			Questions: []PaperQuestion{
				{Question: "When did he go to Xi'an?", Options: []string{"Last month", "Last year", "Next month", "Last week"}, CorrectAnswer: "Last month"},
				{Question: "How did they travel?", Options: []string{"By high-speed train", "By plane", "By car", "By bus"}, CorrectAnswer: "By high-speed train"},
				{Question: "Had he been to Xi'an before?", Options: []string{"No, never", "Yes, once", "Yes, many times", "He lives there"}, CorrectAnswer: "No, never"},
				{Question: "What do they plan for next year?", Options: []string{"See pandas in Chengdu", "Go back to Xi'an", "Visit Beijing", "Stay at home"}, CorrectAnswer: "See pandas in Chengdu"},
			},
		},
		Reading: PaperReading{
			Title: "新同事 — A New Colleague",
			Paragraphs: []string{
				"我们公司来了一个新同事，他叫张华。(Wǒmen gōngsī lái le yí gè xīn tóngshì, tā jiào Zhāng Huá.)",
				"张华比我小三岁，可是他的工作经验比我多。他以前在上海工作过两年。(Zhāng Huá bǐ wǒ xiǎo sān suì, kěshì tā de gōngzuò jīngyàn bǐ wǒ duō. Tā yǐqián zài Shànghǎi gōngzuò guo liǎng nián.)",
				"因为他住得离公司很远，所以每天早上六点就起床了。他一边坐地铁，一边学英语。(Yīnwèi tā zhù de lí gōngsī hěn yuǎn, suǒyǐ měitiān zǎoshang liù diǎn jiù qǐchuáng le. Tā yìbiān zuò dìtiě, yìbiān xué Yīngyǔ.)",
			},
			Questions: []PaperQuestion{
				{Question: "How does Zhang Hua's age compare with the writer's?", Options: []string{"He is 3 years younger", "He is 3 years older", "They are the same age", "He is 13 years younger"}, CorrectAnswer: "He is 3 years younger"},
				{Question: "Where did he work before?", Options: []string{"Shanghai", "Beijing", "Xi'an", "Abroad"}, CorrectAnswer: "Shanghai"},
				{Question: "Why does he get up at six?", Options: []string{"He lives far from the office", "He likes running", "He starts work at six", "He has a baby"}, CorrectAnswer: "He lives far from the office"},
				{Question: "What does he do on the metro?", Options: []string{"Studies English", "Sleeps", "Reads the news", "Plays games"}, CorrectAnswer: "Studies English"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write a short message in Chinese characters to a friend about your last weekend: where you went, what you did and how you felt.",
			MinWords: 40,
		},
		Speaking: PaperSpeaking{
			Phrase:      "虽然汉语很难，但是我每天都学习。",
			Speaker:     "Mira",
			Translation: "Suīrán Hànyǔ hěn nán, dànshì wǒ měitiān dōu xuéxí. — Although Chinese is hard, I study every day.",
		},
	},

	// ---------------- B1 · HSK 3 ----------------
	"B1": {
		Listening: PaperListening{
			Title: "电话预约 — Booking by Phone",
			Lines: []PaperLine{
				{Character: "Mira", Text: "您好，这里是阳光牙科诊所，请问有什么可以帮您？", Translation: "Hello, Sunshine Dental Clinic, how can I help you?"},
				{Character: "Blaze", Text: "您好，我想预约看牙。我的牙已经疼了好几天了。", Translation: "Hello, I'd like to book a dental appointment. My tooth has hurt for several days."},
				{Character: "Mira", Text: "明天上午的号已经约满了，下午三点可以吗？请您提前十分钟到，把身份证带来。", Translation: "Tomorrow morning is fully booked; is 3 p.m. OK? Please arrive ten minutes early and bring your ID card."},
				{Character: "Blaze", Text: "好的，下午三点没问题。如果我有事来不了，需要打电话取消吗？", Translation: "Fine, 3 p.m. is no problem. If I can't come, do I need to call to cancel?"},
				{Character: "Mira", Text: "是的，请至少提前一天通知我们。", Translation: "Yes, please tell us at least a day in advance."},
			},
			Questions: []PaperQuestion{
				{Question: "Why is Blaze calling?", Options: []string{"To book a dental appointment", "To cancel an appointment", "To ask for a job", "To complain"}, CorrectAnswer: "To book a dental appointment"},
				{Question: "Why can't he go tomorrow morning?", Options: []string{"It is fully booked", "He has to work", "The clinic is closed", "The dentist is ill"}, CorrectAnswer: "It is fully booked"},
				{Question: "What must he bring?", Options: []string{"His ID card", "His insurance card", "Money", "A photo"}, CorrectAnswer: "His ID card"},
				{Question: "How early must he cancel?", Options: []string{"At least one day before", "One hour before", "One week before", "He can't cancel"}, CorrectAnswer: "At least one day before"},
			},
		},
		Reading: PaperReading{
			Title: "垃圾分类 — Sorting Rubbish",
			Paragraphs: []string{
				"从去年开始，我们小区要求每家每户把垃圾分成四类：可回收物、厨余垃圾、有害垃圾和其他垃圾。",
				"刚开始的时候，很多居民觉得麻烦，经常把垃圾放错。后来，社区派了志愿者在垃圾桶旁边帮忙，还给大家发了说明手册。",
				"现在，大部分人已经养成了分类的习惯。小区不仅变得更干净了，而且回收的东西也被重新利用了。",
			},
			Questions: []PaperQuestion{
				{Question: "Into how many categories is rubbish sorted?", Options: []string{"Four", "Two", "Three", "Five"}, CorrectAnswer: "Four"},
				{Question: "How did residents feel at first?", Options: []string{"It was troublesome", "It was easy", "It was fun", "They didn't notice"}, CorrectAnswer: "It was troublesome"},
				{Question: "How did the community help?", Options: []string{"Volunteers and instruction booklets", "Fines", "New bins only", "TV adverts"}, CorrectAnswer: "Volunteers and instruction booklets"},
				{Question: "What is the result now?", Options: []string{"Cleaner area and recycled materials reused", "More rubbish", "Residents moved away", "Nothing changed"}, CorrectAnswer: "Cleaner area and recycled materials reused"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write a paragraph in Chinese characters about a healthy habit you have (or want to build): what it is, why it matters, and how you keep it up.",
			MinWords: 80,
		},
		Speaking: PaperSpeaking{
			Phrase:      "我已经把今天的作业写完了，现在想出去散散步。",
			Speaker:     "Blaze",
			Translation: "Wǒ yǐjīng bǎ jīntiān de zuòyè xiě wán le, xiànzài xiǎng chūqu sànsan bù. — I've finished today's homework; now I'd like to go for a walk.",
		},
	},

	// ---------------- B2 · HSK 4 ----------------
	"B2": {
		Listening: PaperListening{
			Title: "求职面试 — A Job Interview",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "请你先简单介绍一下自己，为什么想来我们公司工作？", Translation: "First, briefly introduce yourself. Why do you want to work for us?"},
				{Character: "Riko", Text: "我叫刘洋，大学学的是市场营销，毕业后在一家电商公司做了三年运营。", Translation: "I'm Liu Yang. I studied marketing and, after graduating, worked in operations at an e-commerce company for three years."},
				{Character: "Riko", Text: "贵公司在短视频带货方面做得很成功，我希望能把自己的经验用在这里，同时学到更多东西。", Translation: "Your company has been very successful in short-video sales; I'd like to apply my experience here and learn more."},
				{Character: "Professor Finch", Text: "这个职位经常需要出差，有时候周末也要加班，你能接受吗？", Translation: "This role involves frequent travel and sometimes weekend work. Can you accept that?"},
				{Character: "Riko", Text: "可以接受。不过我想了解一下，加班的话有没有调休或者补贴？", Translation: "Yes. But I'd like to know: is there time off in lieu or an allowance for overtime?"},
			},
			Questions: []PaperQuestion{
				{Question: "What did Liu Yang study?", Options: []string{"Marketing", "Computer science", "Law", "Medicine"}, CorrectAnswer: "Marketing"},
				{Question: "How long has he worked in operations?", Options: []string{"Three years", "One year", "Five years", "Six months"}, CorrectAnswer: "Three years"},
				{Question: "What is the company good at?", Options: []string{"Selling through short videos", "Making phones", "Building houses", "Running hotels"}, CorrectAnswer: "Selling through short videos"},
				{Question: "What does he ask about at the end?", Options: []string{"Time off in lieu or allowances for overtime", "The holiday dates", "The office location", "The dress code"}, CorrectAnswer: "Time off in lieu or allowances for overtime"},
			},
		},
		Reading: PaperReading{
			Title: "共享单车的变化 — Bike-Sharing Then and Now",
			Paragraphs: []string{
				"几年前，共享单车在中国的城市里迅速流行起来。人们只要用手机扫一下码，就能骑走一辆自行车，非常方便。",
				"然而，由于投放数量过多，很多单车被随意停放，甚至堆满了人行道，给城市管理带来了不少麻烦。",
				"后来，政府出台了相关规定，限制投放数量，并划定了专门的停车区域。企业也开始利用定位技术管理车辆。",
				"如今，共享单车已经成为城市交通的一部分。它的经历说明，新技术的发展既需要创新，也离不开合理的管理。",
			},
			Questions: []PaperQuestion{
				{Question: "Why was bike-sharing popular?", Options: []string{"Scan a code and ride — very convenient", "It was free forever", "Cars were banned", "Buses stopped running"}, CorrectAnswer: "Scan a code and ride — very convenient"},
				{Question: "What problem appeared?", Options: []string{"Too many bikes parked everywhere", "Too few bikes", "Bikes were too expensive", "Nobody used them"}, CorrectAnswer: "Too many bikes parked everywhere"},
				{Question: "What did the government do?", Options: []string{"Limited numbers and set parking areas", "Banned all bikes", "Bought the companies", "Nothing"}, CorrectAnswer: "Limited numbers and set parking areas"},
				{Question: "What is the main lesson?", Options: []string{"Innovation needs sensible management", "Technology always fails", "Bikes are better than cars", "Phones are dangerous"}, CorrectAnswer: "Innovation needs sensible management"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write a short essay in Chinese characters: 'Is working from home a good idea?' Give your view, two reasons and an example.",
			MinWords: 150,
		},
		Speaking: PaperSpeaking{
			Phrase:      "在我看来，移动支付虽然方便，但是我们也应该注意个人信息的安全。",
			Speaker:     "Pip",
			Translation: "In my view, although mobile payment is convenient, we should also pay attention to the security of personal information.",
		},
	},

	// ---------------- C1 · HSK 5 ----------------
	"C1": {
		Listening: PaperListening{
			Title: "老龄化社会 — An Ageing Society",
			Lines: []PaperLine{
				{Character: "Mira", Text: "据统计，到二〇三五年，中国六十岁以上的人口将超过四亿，老龄化问题日益突出。", Translation: "Statistics show that by 2035 China's over-60 population will exceed 400 million; ageing is an increasingly prominent issue."},
				{Character: "Mira", Text: "一方面，养老金和医疗支出不断增加，给社会保障体系带来了压力；另一方面，劳动力的减少也可能影响经济增长。", Translation: "On one hand, pensions and medical spending keep rising, straining social security; on the other, a shrinking workforce may affect growth."},
				{Character: "Mira", Text: "对此，一些城市开始发展社区养老，让老人在熟悉的环境中得到照顾，同时鼓励身体健康的老人参与志愿服务。", Translation: "In response, some cities are developing community-based elder care, so older people are looked after in familiar surroundings, while healthy elders are encouraged to volunteer."},
				{Character: "Mira", Text: "专家指出，老龄化既是挑战，也蕴含着机遇，比如所谓的“银发经济”正在快速发展。", Translation: "Experts note ageing is both a challenge and an opportunity — the so-called 'silver economy' is growing fast."},
			},
			Questions: []PaperQuestion{
				{Question: "到二〇三五年，中国六十岁以上的人口将超过多少？", Options: []string{"四亿", "四千万", "十亿", "一亿"}, CorrectAnswer: "四亿"},
				{Question: "老龄化给社会保障体系带来了什么？", Options: []string{"压力", "更多收入", "更少支出", "没有影响"}, CorrectAnswer: "压力"},
				{Question: "“社区养老”的好处是什么？", Options: []string{"老人在熟悉的环境中得到照顾", "老人搬到大城市", "老人重新上班", "减少医院数量"}, CorrectAnswer: "老人在熟悉的环境中得到照顾"},
				{Question: "专家怎么看待老龄化？", Options: []string{"既是挑战也是机遇", "只是一个问题", "完全没有问题", "无法解决"}, CorrectAnswer: "既是挑战也是机遇"},
				{Question: "“银发经济”指的是：", Options: []string{"面向老年人的经济", "白银市场", "高收入人群", "年轻人消费"}, CorrectAnswer: "面向老年人的经济"},
			},
		},
		Reading: PaperReading{
			Title: "慢下来的勇气 — The Courage to Slow Down",
			Paragraphs: []string{
				"在快节奏的现代社会，“效率”几乎成了衡量一切的标准。人们习惯一边吃饭一边回复消息，一边走路一边听课，仿佛停下来就意味着落后。",
				"然而，心理学研究表明，长期处于多任务状态不仅会降低工作质量，而且容易导致焦虑和疲劳。大脑需要专注，也需要休息。",
				"近年来，越来越多的人开始尝试“慢生活”：每天留出一段不看手机的时间，认真读一本书，或者只是安静地散步。",
				"由此可见，慢下来并不等于懒惰，而是一种对生活的重新选择。有时候，敢于慢下来，反而需要更大的勇气。",
			},
			Questions: []PaperQuestion{
				{Question: "在现代社会，什么几乎成了衡量一切的标准？", Options: []string{"效率", "金钱", "健康", "年龄"}, CorrectAnswer: "效率"},
				{Question: "长期多任务会带来什么问题？", Options: []string{"工作质量下降，容易焦虑", "收入增加", "更有创造力", "睡得更好"}, CorrectAnswer: "工作质量下降，容易焦虑"},
				{Question: "“慢生活”的例子是：", Options: []string{"每天留出不看手机的时间", "加班到深夜", "同时做三件事", "每天换工作"}, CorrectAnswer: "每天留出不看手机的时间"},
				{Question: "作者认为慢下来是：", Options: []string{"对生活的重新选择", "一种懒惰", "浪费时间", "不负责任"}, CorrectAnswer: "对生活的重新选择"},
				{Question: "最后一句的意思是：", Options: []string{"敢于慢下来需要更大的勇气", "慢的人没有勇气", "勇气不重要", "快一点更勇敢"}, CorrectAnswer: "敢于慢下来需要更大的勇气"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write an essay in Chinese characters (用中文写作): 网络对人际关系的影响 — the internet's effect on relationships. State your view, analyse two aspects with examples, and conclude.",
			MinWords: 250,
		},
		Speaking: PaperSpeaking{
			Phrase:      "研究表明，适当的休息不仅能提高效率，而且有助于缓解焦虑。",
			Speaker:     "Mira",
			Translation: "Research shows that proper rest not only improves efficiency but also helps relieve anxiety.",
		},
	},

	// ---------------- C2 · HSK 6 ----------------
	"C2": {
		Listening: PaperListening{
			Title: "文化遗产的保护与开发 — Heritage: Protect or Develop?",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "尊敬的各位来宾，近年来，各地纷纷将古镇、古村落开发为旅游景点，由此引发了广泛争议。", Translation: "Respected guests, in recent years towns everywhere have turned ancient towns and villages into tourist sites, sparking wide debate."},
				{Character: "Professor Finch", Text: "诚然，旅游开发为当地带来了可观的经济收益，也让更多人了解了传统文化。", Translation: "Admittedly, tourism has brought considerable income and introduced more people to traditional culture."},
				{Character: "Professor Finch", Text: "然而，过度商业化使不少古镇千篇一律：同样的小吃街，同样的纪念品，原有的生活气息荡然无存。", Translation: "Yet over-commercialisation has made many ancient towns identical: the same snack streets, the same souvenirs, the original way of life completely gone."},
				{Character: "Professor Finch", Text: "保护与开发并非水火不容。关键在于因地制宜，在尊重历史原貌的前提下适度开发，并让当地居民真正参与其中、从中受益。", Translation: "Protection and development are not irreconcilable. The key is to adapt to local conditions, develop moderately while respecting the historical character, and let residents genuinely take part and benefit."},
				{Character: "Professor Finch", Text: "唯有如此，文化遗产才能既“活”下来，又“传”下去。", Translation: "Only thus can cultural heritage both 'live on' and 'be passed down'."},
			},
			Questions: []PaperQuestion{
				{Question: "引发争议的现象是什么？", Options: []string{"古镇被开发为旅游景点", "古镇被拆除", "游客减少", "居民搬进古镇"}, CorrectAnswer: "古镇被开发为旅游景点"},
				{Question: "发言人承认旅游开发的好处是：", Options: []string{"带来经济收益并传播传统文化", "保护了自然环境", "减少了游客", "降低了物价"}, CorrectAnswer: "带来经济收益并传播传统文化"},
				{Question: "“千篇一律”在这里的意思是：", Options: []string{"都一样，没有特色", "有很多文章", "非常有特色", "历史悠久"}, CorrectAnswer: "都一样，没有特色"},
				{Question: "发言人认为保护与开发的关系是：", Options: []string{"并非水火不容", "完全对立", "只能选一个", "没有关系"}, CorrectAnswer: "并非水火不容"},
				{Question: "他提出的关键做法包括：", Options: []string{"因地制宜、适度开发、让居民受益", "全面停止开发", "只靠政府投资", "只发展小吃街"}, CorrectAnswer: "因地制宜、适度开发、让居民受益"},
			},
		},
		Reading: PaperReading{
			Title: "算法时代的选择 — Choice in the Age of Algorithms",
			Paragraphs: []string{
				"打开手机，推荐系统总能精准地送上我们“感兴趣”的内容。表面上看，算法让信息获取变得前所未有地便捷；实际上，它也在悄然塑造着我们的视野。",
				"长期沉浸在个性化推荐之中，人们接触到的观点日益同质化，逐渐形成所谓的“信息茧房”。久而久之，不同群体之间的理解与沟通愈发困难。",
				"对此，有人主张加强对平台的监管，要求算法更加透明；也有人认为，与其一味指责技术，不如提高公众的媒介素养，主动走出舒适区。",
				"笔者以为，技术本身无所谓善恶，问题在于使用者能否保持清醒的判断。宁可多花一点时间去求证，也不轻易被算法牵着鼻子走——这或许是数字时代最可贵的自觉。",
			},
			Questions: []PaperQuestion{
				{Question: "第一段指出算法的影响是：", Options: []string{"既便捷，又在悄然塑造视野", "只有好处", "只有坏处", "与我们无关"}, CorrectAnswer: "既便捷，又在悄然塑造视野"},
				{Question: "“信息茧房”指的是：", Options: []string{"只接触相似观点而封闭起来", "储存信息的房间", "一种新手机", "信息太多"}, CorrectAnswer: "只接触相似观点而封闭起来"},
				{Question: "第三段提到了哪两种主张？", Options: []string{"加强监管；提高公众媒介素养", "禁止手机；关闭平台", "降低价格；增加广告", "多看推荐；少读书"}, CorrectAnswer: "加强监管；提高公众媒介素养"},
				{Question: "作者对技术的看法是：", Options: []string{"技术本身无所谓善恶", "技术是邪恶的", "技术完全有益", "应该放弃技术"}, CorrectAnswer: "技术本身无所谓善恶"},
				{Question: "“被算法牵着鼻子走”的意思是：", Options: []string{"被算法控制而失去判断", "用鼻子操作手机", "主动使用算法", "拒绝使用算法"}, CorrectAnswer: "被算法控制而失去判断"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write an argumentative essay in Chinese characters (用中文写作): 经济发展与环境保护能否兼顾？ Present arguments and counter-arguments, use formal connectors, and reach a reasoned conclusion.",
			MinWords: 350,
		},
		Speaking: PaperSpeaking{
			Phrase:      "诚然，发展经济至关重要，但我们宁可放慢脚步，也不能以牺牲环境为代价。",
			Speaker:     "Zephyr",
			Translation: "Admittedly, economic development is vital, but we would rather slow down than sacrifice the environment.",
		},
	},

	// ---------------- FINAL · HSK 7–9 (advanced band) ----------------
	"FINAL": {
		Listening: PaperListening{
			Title: "学术讲座：语言与思维 — Lecture: Language and Thought",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "各位同学，今天的讲座围绕一个经久不衰的问题展开：语言究竟在多大程度上影响我们的思维？", Translation: "Students, today's lecture centres on an enduring question: to what extent does language shape our thinking?"},
				{Character: "Professor Finch", Text: "早在上世纪，“语言相对论”便提出，不同语言的使用者会以不同的方式认知世界。这一观点一度引发学界激烈争论。", Translation: "Last century, 'linguistic relativity' proposed that speakers of different languages perceive the world differently — a view that once sparked fierce academic debate."},
				{Character: "Professor Finch", Text: "近年来的实证研究表明，语言固然不能决定思维，却能在颜色辨认、空间方位乃至时间表述等方面施加微妙而可测的影响。", Translation: "Recent empirical studies show language cannot determine thought, but it exerts subtle, measurable effects on colour discrimination, spatial orientation and even how time is expressed."},
				{Character: "Professor Finch", Text: "以汉语为例，人们常用“上”“下”来描述时间的先后，如“上个月”“下周”，这或许在潜移默化中塑造了一种纵向的时间观。", Translation: "In Chinese, for instance, 'up' and 'down' describe earlier and later — 'last month', 'next week' — which may subtly foster a vertical sense of time."},
				{Character: "Professor Finch", Text: "由此可见，掌握一门外语，不仅是多了一种交流工具，更是获得了一种重新审视世界的视角。", Translation: "Thus, mastering a foreign language means not just another tool for communication but a new perspective from which to re-examine the world."},
			},
			Questions: []PaperQuestion{
				{Question: "讲座讨论的核心问题是什么？", Options: []string{"语言在多大程度上影响思维", "如何快速学习外语", "汉语的历史", "颜色的种类"}, CorrectAnswer: "语言在多大程度上影响思维"},
				{Question: "“语言相对论”认为：", Options: []string{"不同语言使用者以不同方式认知世界", "所有语言都一样", "语言决定一切", "思维与语言无关"}, CorrectAnswer: "不同语言使用者以不同方式认知世界"},
				{Question: "近年来的实证研究发现：", Options: []string{"语言不能决定思维，但有微妙影响", "语言完全决定思维", "语言对思维没有任何影响", "研究无法进行"}, CorrectAnswer: "语言不能决定思维，但有微妙影响"},
				{Question: "汉语用什么来描述时间先后？", Options: []string{"“上”和“下”", "“左”和“右”", "“前”和“后”只用于空间", "颜色词"}, CorrectAnswer: "“上”和“下”"},
				{Question: "“潜移默化”的意思是：", Options: []string{"在不知不觉中受到影响", "突然改变", "明显的变化", "拒绝改变"}, CorrectAnswer: "在不知不觉中受到影响"},
				{Question: "讲座的结论是：", Options: []string{"学外语能获得审视世界的新视角", "学外语没有必要", "母语最重要", "语言只是工具"}, CorrectAnswer: "学外语能获得审视世界的新视角"},
			},
		},
		Reading: PaperReading{
			Title: "翻译之难 — The Difficulty of Translation",
			Paragraphs: []string{
				"严复曾提出翻译的三大标准：“信、达、雅”，即忠实于原文、表达通顺、文辞优美。百余年来，这一主张始终被译界奉为圭臬，却也备受质疑。",
				"质疑者认为，三者之间往往难以兼顾：过分追求“雅”，难免偏离原意；一味强调“信”，译文又可能佶屈聱牙。尤其在文学翻译中，双关、典故与韵律几乎无法完整移植。",
				"然而，正是这种“不可译”之处，彰显了翻译的创造性。优秀的译者并非被动地搬运文字，而是在两种文化之间斡旋、取舍，力求在新的语境中重现原作的神韵。",
				"诚如一位学者所言，翻译是一门“带着镣铐跳舞”的艺术。镣铐无法挣脱，舞姿却可以因人而异——这恰恰是翻译经久不衰的魅力所在。",
			},
			Questions: []PaperQuestion{
				{Question: "“信、达、雅”分别指：", Options: []string{"忠实原文、表达通顺、文辞优美", "速度、价格、质量", "词汇、语法、发音", "长度、格式、字体"}, CorrectAnswer: "忠实原文、表达通顺、文辞优美"},
				{Question: "质疑者的主要观点是：", Options: []string{"三条标准往往难以兼顾", "三条标准都没有用", "翻译非常容易", "只需要追求“雅”"}, CorrectAnswer: "三条标准往往难以兼顾"},
				{Question: "文学翻译中特别难以移植的是：", Options: []string{"双关、典故与韵律", "人名和地名", "数字", "标点符号"}, CorrectAnswer: "双关、典故与韵律"},
				{Question: "作者如何看待“不可译”之处？", Options: []string{"它彰显了翻译的创造性", "它说明翻译毫无意义", "应该删掉这些内容", "只能直译"}, CorrectAnswer: "它彰显了翻译的创造性"},
				{Question: "“带着镣铐跳舞”比喻的是：", Options: []string{"在原文限制下进行创造", "翻译工作很危险", "译者喜欢跳舞", "完全自由地改写"}, CorrectAnswer: "在原文限制下进行创造"},
				{Question: "“奉为圭臬”的意思是：", Options: []string{"当作准则来遵守", "当作笑话", "完全否定", "暂时忘记"}, CorrectAnswer: "当作准则来遵守"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write a comprehensive, well-structured essay in Chinese characters (用中文写作): 学习外语的意义 — the value of learning a foreign language in the age of machine translation. Argue, rebut a counter-view, and conclude with a reasoned position.",
			MinWords: 500,
		},
		Speaking: PaperSpeaking{
			Phrase:      "掌握一门外语，不仅是多了一种交流工具，更是获得了一种重新审视世界的视角。",
			Speaker:     "Lumora",
			Translation: "Mastering a foreign language means not just another tool for communication, but a new perspective from which to re-examine the world.",
		},
	},
}
