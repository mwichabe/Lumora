package database

import (
	"gorm.io/gorm"

	"lumora/backend/models"
)

// Japanese listening and reading sessions — at least one of each per level.
// Transcripts carry readings in the early levels and fade them out from N3.

const (
	jaN5 = "A1 · JLPT N5 — Foundations & Survival Japanese"
	jaN4 = "A2 · JLPT N4 — Everyday Japanese"
	jaN3 = "B1 · JLPT N3 — Intermediate Japanese"
	jaN2 = "B2 · JLPT N2 — Upper-Intermediate Japanese"
	jaN1 = "C1 · JLPT N1 — Advanced Japanese"
	jaC2 = "C2 · Beyond N1 — Native-Level Japanese"
)

func seedJapaneseListening(db *gorm.DB) {
	const ja = "ja"
	finch := "Professor Finch"

	addListeningL(db, ja, jaN5, "はじめまして — First Day of Class",
		"Professor Finch meets a new student. Listen, then answer.", 1, 20,
		[]models.ListeningMatch{
			lm("はじめまして", "nice to meet you"), lm("おなまえは？", "your name?"),
			lm("〜からきました", "I'm from ~"), lm("よろしくおねがいします", "please treat me well"),
		},
		[]models.ListeningLine{
			ln(finch, "はじめまして。フィンチです。おなまえは？", "Hajimemashite. Finchi desu. Onamae wa? — Nice to meet you. I'm Finch. Your name?"),
			ln("Lumora", "ルモラです。ケニアからきました。", "Rumora desu. Kenia kara kimashita. — I'm Lumora. I'm from Kenya."),
			ln(finch, "がくせいですか。", "Gakusei desu ka. — Are you a student?"),
			ln("Lumora", "はい、だいがくせいです。にじゅうさいです。よろしくおねがいします。", "Hai, daigakusei desu. Nijussai desu. Yoroshiku onegaishimasu. — Yes, I'm a university student. I'm twenty. Nice to meet you."),
		},
		[]models.ListeningQuestion{
			lq("Where is Lumora from?", "Kenya", "Kenya", "Japan", "China", "Germany"),
			lq("What is she?", "A university student", "A university student", "A teacher", "A doctor", "A chef"),
			lq("How old is she?", "20", "20", "12", "22", "2"),
		},
	)
	addListeningL(db, ja, jaN5, "レストランで — At the Restaurant",
		"Cora orders lunch. Listen, then answer.", 2, 20,
		[]models.ListeningMatch{
			lm("いらっしゃいませ", "welcome"), lm("ラーメン", "ramen"),
			lm("みず", "water"), lm("ぜんぶで", "in total"),
		},
		[]models.ListeningLine{
			ln("Riko", "いらっしゃいませ。なんめいさまですか。", "Irasshaimase. Nanmei-sama desu ka. — Welcome. How many people?"),
			ln("Cora", "ひとりです。ラーメンとぎょうざをおねがいします。", "Hitori desu. Rāmen to gyōza o onegaishimasu. — Just one. Ramen and gyoza, please."),
			ln("Riko", "おのみものは？", "Onomimono wa? — Anything to drink?"),
			ln("Cora", "みずをください。", "Mizu o kudasai. — Water, please."),
			ln("Riko", "ぜんぶで、きゅうひゃくごじゅうえんです。", "Zenbu de, kyūhyaku gojū en desu. — That's 950 yen in total."),
		},
		[]models.ListeningQuestion{
			lq("How many people are eating?", "One", "One", "Two", "Three", "Four"),
			lq("What does Cora order to eat?", "Ramen and gyoza", "Ramen and gyoza", "Sushi", "Rice and fish", "Udon"),
			lq("How much is it in total?", "¥950", "¥950", "¥590", "¥1,950", "¥905"),
		},
	)
	addListeningL(db, ja, jaN4, "週末の予定 — Weekend Plans",
		"Two friends make plans. Listen, then answer.", 3, 22,
		[]models.ListeningMatch{
			lm("ひま", "free (not busy)"), lm("映画", "film"),
			lm("〜ませんか", "won't you…?"), lm("駅の前", "in front of the station"),
		},
		[]models.ListeningLine{
			ln("Blaze", "ねえ、土曜日ひま？一緒に映画を見に行かない？", "ねえ、どようびひま？いっしょにえいがをみにいかない？ — Hey, are you free Saturday? Want to go see a film?"),
			ln("Mira", "ごめん、土曜日はバイトがあるんだ。日曜日なら大丈夫だよ。", "ごめん、どようびはバイトがあるんだ。にちようびならだいじょうぶだよ。 — Sorry, I've got my part-time job Saturday. Sunday's fine, though."),
			ln("Blaze", "じゃあ、日曜日の二時に駅の前で会おう。", "じゃあ、にちようびのにじにえきのまえであおう。 — Then let's meet in front of the station at two on Sunday."),
			ln("Mira", "うん、いいよ。映画の後でご飯も食べたいな。", "うん、いいよ。えいがのあとでごはんもたべたいな。 — OK. I'd like to eat after the film too."),
		},
		[]models.ListeningQuestion{
			lq("Why can't Mira go on Saturday?", "She has her part-time job", "She has her part-time job", "She's sick", "She's travelling", "She has an exam"),
			lq("When will they meet?", "Sunday at 2:00", "Sunday at 2:00", "Saturday at 2:00", "Sunday at 12:00", "Monday at 2:00"),
			lq("Where will they meet?", "In front of the station", "In front of the station", "At the cinema", "At a restaurant", "At Mira's house"),
			lq("What does Mira want to do after the film?", "Eat a meal", "Eat a meal", "Go shopping", "Go home", "Watch another film"),
		},
	)
	addListeningL(db, ja, jaN3, "病院で — At the Clinic",
		"A patient sees the doctor. Listen, then answer.", 4, 24,
		[]models.ListeningMatch{
			lm("せき", "cough"), lm("熱", "fever"),
			lm("薬", "medicine"), lm("無理をしない", "don't overdo it"),
		},
		[]models.ListeningLine{
			ln("Nana", "今日はどうしましたか。", "What brings you in today?"),
			ln("Zephyr", "三日前からせきが止まらなくて、昨日の夜から熱も出てきたんです。", "I haven't been able to stop coughing for three days, and since last night I've had a fever too."),
			ln("Nana", "風邪のようですね。薬を出しますので、一日三回、食後に飲んでください。", "It looks like a cold. I'll prescribe medicine — take it three times a day after meals."),
			ln("Zephyr", "仕事は休んだほうがいいでしょうか。", "Should I take time off work?"),
			ln("Nana", "熱が下がるまでは、無理をしないで休んでください。", "Until the fever goes down, don't push yourself — please rest."),
		},
		[]models.ListeningQuestion{
			lq("How long has he been coughing?", "Three days", "Three days", "One day", "A week", "Since this morning"),
			lq("When did the fever start?", "Last night", "Last night", "Three days ago", "This morning", "A week ago"),
			lq("How often should he take the medicine?", "Three times a day after meals", "Three times a day after meals", "Once a day before bed", "Twice a day before meals", "Only when coughing"),
			lq("What does the doctor advise about work?", "Rest until the fever goes down", "Rest until the fever goes down", "Go to work as usual", "Work from home", "Quit his job"),
		},
	)
	addListeningL(db, ja, jaN2, "取引先への電話 — A Business Call",
		"A call to a client's office. Listen, then answer.", 5, 26,
		[]models.ListeningMatch{
			lm("お世話になっております", "thank you for your support"), lm("席を外しております", "is away from their desk"),
			lm("折り返し", "call back"), lm("伝言", "message"),
		},
		[]models.ListeningLine{
			ln("Mira", "お電話ありがとうございます。山田商事でございます。", "Thank you for calling. This is Yamada Trading."),
			ln("Blaze", "いつもお世話になっております。ABC社の佐藤と申しますが、営業部の田中様はいらっしゃいますか。", "Thank you for your continued support. My name is Sato from ABC; is Mr Tanaka in Sales available?"),
			ln("Mira", "申し訳ございません。田中はただいま会議中で、席を外しております。四時ごろには戻る予定ですが。", "I'm very sorry. Tanaka is in a meeting and away from his desk right now. He's expected back around four."),
			ln("Blaze", "では、恐れ入りますが、戻られましたら折り返しお電話いただけますでしょうか。", "In that case, sorry to trouble you, but could you ask him to call me back when he returns?"),
			ln("Mira", "かしこまりました。佐藤様からお電話があったと、田中に申し伝えます。", "Certainly. I'll let Tanaka know that you called, Mr Sato."),
		},
		[]models.ListeningQuestion{
			lq("Who is the caller?", "Sato from ABC", "Sato from ABC", "Tanaka from Yamada Trading", "Yamada from ABC", "Tanaka from ABC"),
			lq("Why can't Tanaka take the call?", "He's in a meeting", "He's in a meeting", "He's on holiday", "He's left the company", "He's on another call"),
			lq("When is Tanaka expected back?", "Around 4:00", "Around 4:00", "Around 2:00", "Tomorrow", "In an hour"),
			lq("What does the caller ask for?", "A call back", "A call back", "An email", "A meeting tomorrow", "To wait on the line"),
		},
	)
	addListeningL(db, ja, jaN1, "講義：都市と過疎化 — Lecture: Cities and Depopulation",
		"An excerpt from a university lecture. Listen, then answer.", 6, 28,
		[]models.ListeningMatch{
			lm("過疎化", "rural depopulation"), lm("一極集中", "over-concentration in one place"),
			lm("地方創生", "regional revitalisation"), lm("示唆", "suggestion, implication"),
		},
		[]models.ListeningLine{
			ln(finch, "本日は、地方の過疎化と東京一極集中について考えてみたいと思います。", "Today I'd like to consider rural depopulation and the over-concentration of people in Tokyo."),
			ln(finch, "戦後の高度経済成長期以降、若年層は仕事を求めて都市部へ移動し続けてきました。", "Since the high-growth era after the war, young people have kept moving to cities in search of work."),
			ln(finch, "その結果、地方では高齢化が進み、学校や病院の維持すら困難になっている地域も少なくありません。", "As a result, regions have aged, and quite a few areas struggle even to keep schools and hospitals running."),
			ln(finch, "一方で、近年はリモートワークの普及により、地方へ移住する人も増えつつあります。これは、従来の一極集中の構造が変わりうることを示唆していると言えるでしょう。", "On the other hand, with remote work spreading, more people have recently been moving to the regions. This suggests the traditional over-concentration may be able to change."),
		},
		[]models.ListeningQuestion{
			lq("What is the topic of the lecture?", "Rural depopulation and concentration in Tokyo", "Rural depopulation and concentration in Tokyo", "The history of remote work", "Hospital funding in Tokyo", "The post-war economy only"),
			lq("Why did young people move to cities?", "To find work", "To find work", "For education only", "Because of disasters", "For cheaper housing"),
			lq("What problem do some regions face?", "Even maintaining schools and hospitals is hard", "Even maintaining schools and hospitals is hard", "Too many young people", "Rising house prices", "Overcrowded trains"),
			lq("What recent trend does the lecturer mention?", "More people moving to the regions thanks to remote work", "More people moving to the regions thanks to remote work", "Tokyo's population exploding further", "Schools reopening in the countryside", "Companies banning remote work"),
		},
	)
	addListeningL(db, ja, jaC2, "大阪の商店街で — An Osaka Shopping Street",
		"A shopkeeper chats in Kansai dialect. Listen, then answer.", 7, 30,
		[]models.ListeningMatch{
			lm("まいど", "hello / thanks (Kansai shop greeting)"), lm("ほんま", "really"),
			lm("あかん", "no good"), lm("おまけ", "a free extra"),
		},
		[]models.ListeningLine{
			ln("Blaze", "まいど！お姉ちゃん、今日はええトマト入ってるで。", "Hello there! Miss, we've got some great tomatoes in today."),
			ln("Cora", "ほんまに？ほな、三つもらうわ。なんぼ？", "Really? Then I'll take three. How much?"),
			ln("Blaze", "三つで三百円や。せやけど、一個おまけしとくわ。", "Three for 300 yen. But I'll throw in one extra for free."),
			ln("Cora", "えっ、ほんま？おおきに！また来るわ。", "Really? Thanks so much! I'll come again."),
			ln("Blaze", "あかんあかん、お礼なんかええねん。また寄ってや！", "No, no, no need to thank me. Drop by again!"),
		},
		[]models.ListeningQuestion{
			lq("Which dialect is spoken?", "Kansai dialect", "Kansai dialect", "Standard Tokyo Japanese", "Hakata dialect", "Okinawan"),
			lq("What does なんぼ mean here?", "How much?", "How much?", "How many?", "Where?", "Why?"),
			lq("How many tomatoes does Cora get in the end?", "Four", "Four", "Three", "Two", "Five"),
			lq("What does おおきに mean?", "Thank you", "Thank you", "Goodbye", "Delicious", "Expensive"),
		},
	)
}

func seedJapaneseReading(db *gorm.DB) {
	const ja = "ja"

	addReadingL(db, ja, jaN5, "わたしのかぞく — My Family",
		"A short self-introduction. Read, then answer.", 1, 20,
		[]models.ReadingLine{
			rl("わたしはさとうけんです。(Watashi wa Satō Ken desu.)", "I'm Ken Sato."),
			rl("とうきょうにすんでいます。(Tōkyō ni sunde imasu.)", "I live in Tokyo."),
			rl("かぞくはよにんです。ちちとははといもうとがいます。(Kazoku wa yonin desu. Chichi to haha to imōto ga imasu.)", "There are four in my family: my father, my mother and my little sister."),
			rl("ちちはいしゃです。ははは日本語のせんせいです。(Chichi wa isha desu. Haha wa nihongo no sensei desu.)", "My father is a doctor. My mother is a Japanese teacher."),
			rl("わたしはねこがすきです。(Watashi wa neko ga suki desu.)", "I like cats."),
		},
		[]models.ReadingQuestion{
			rq("Where does Ken live?", "Tokyo", "Tokyo", "Osaka", "Kyoto", "Kenya"),
			rq("How many people are in his family?", "4", "4", "3", "5", "6"),
			rq("What is his mother's job?", "Japanese teacher", "Japanese teacher", "Doctor", "Student", "Chef"),
			rq("What does Ken like?", "Cats", "Cats", "Dogs", "Sushi", "Trains"),
		},
	)
	addReadingL(db, ja, jaN5, "かいもの — Going Shopping",
		"Ken goes to the shop. Read, then answer.", 2, 20,
		[]models.ReadingLine{
			rl("きのう、コンビニへいきました。(Kinō, konbini e ikimashita.)", "Yesterday I went to the convenience store."),
			rl("パンとぎゅうにゅうをかいました。(Pan to gyūnyū o kaimashita.)", "I bought bread and milk."),
			rl("パンはひゃくごじゅうえんでした。ぎゅうにゅうはにひゃくえんでした。(Pan wa hyaku gojū en deshita. Gyūnyū wa nihyaku en deshita.)", "The bread was 150 yen. The milk was 200 yen."),
			rl("コーヒーもほしかったですが、たかかったです。(Kōhī mo hoshikatta desu ga, takakatta desu.)", "I wanted coffee too, but it was expensive."),
		},
		[]models.ReadingQuestion{
			rq("Where did Ken go?", "A convenience store", "A convenience store", "A supermarket", "A café", "A restaurant"),
			rq("How much did he spend in total?", "¥350", "¥350", "¥150", "¥200", "¥500"),
			rq("Why didn't he buy coffee?", "It was expensive", "It was expensive", "They had none", "He doesn't like it", "He forgot"),
		},
	)
	addReadingL(db, ja, jaN4, "日本での生活 — Life in Japan",
		"A student writes about her first months. Read, then answer.", 3, 22,
		[]models.ReadingLine{
			rl("わたしは三か月前に日本へ来ました。(わたしはさんかげつまえににほんへきました。)", "I came to Japan three months ago."),
			rl("はじめは日本語があまりわからなくて、たいへんでした。(はじめはにほんごがあまりわからなくて、たいへんでした。)", "At first I didn't understand much Japanese, and it was hard."),
			rl("でも、毎日友達と話したので、少しずつ話せるようになりました。(でも、まいにちともだちとはなしたので、すこしずつはなせるようになりました。)", "But because I talked with friends every day, little by little I became able to speak."),
			rl("来年は、京都で日本の文化を勉強するつもりです。(らいねんは、きょうとでにほんのぶんかをべんきょうするつもりです。)", "Next year I plan to study Japanese culture in Kyoto."),
		},
		[]models.ReadingQuestion{
			rq("When did she come to Japan?", "Three months ago", "Three months ago", "Three years ago", "Last week", "Three days ago"),
			rq("Why was it hard at first?", "She didn't understand much Japanese", "She didn't understand much Japanese", "She had no friends", "The food was strange", "She was ill"),
			rq("How did her Japanese improve?", "By talking with friends every day", "By talking with friends every day", "By watching anime", "By taking classes only", "By reading newspapers"),
			rq("What does she plan to do next year?", "Study Japanese culture in Kyoto", "Study Japanese culture in Kyoto", "Go back home", "Work in Tokyo", "Travel to Osaka"),
		},
	)
	addReadingL(db, ja, jaN3, "自動販売機の国 — The Land of Vending Machines",
		"A short article. Read, then answer.", 4, 24,
		[]models.ReadingLine{
			rl("日本には約四百万台の自動販売機があると言われている。", "It's said there are about four million vending machines in Japan."),
			rl("飲み物だけでなく、食べ物や本など、さまざまな物が売られている。", "Not only drinks — all sorts of things, like food and books, are sold."),
			rl("これほど多い理由の一つは、治安がよく、機械が壊されにくいことだそうだ。", "One reason there are so many is apparently that public safety is good, so machines rarely get damaged."),
			rl("しかし最近は、コンビニの増加や電気代の上昇により、台数は少しずつ減っている。", "Recently, however, with more convenience stores and rising electricity costs, the number is gradually falling."),
		},
		[]models.ReadingQuestion{
			rq("About how many vending machines are in Japan?", "About 4 million", "About 4 million", "About 400,000", "About 40 million", "About 4,000"),
			rq("What is one reason there are so many?", "Good public safety, so they're rarely damaged", "Good public safety, so they're rarely damaged", "Electricity is free", "There are no shops", "The government requires them"),
			rq("What is happening to their number recently?", "It's slowly decreasing", "It's slowly decreasing", "It's rising fast", "It's unchanged", "They've been banned"),
			rq("Which is NOT given as a cause of the decrease?", "Damage by criminals", "Damage by criminals", "More convenience stores", "Rising electricity costs"),
		},
	)
	addReadingL(db, ja, jaN2, "ワークライフバランス — Work-Life Balance",
		"An opinion column. Read, then answer.", 5, 26,
		[]models.ReadingLine{
			rl("かつて日本では、長時間労働こそが会社への忠誠の証とされてきた。", "In the past, long working hours were seen in Japan as proof of loyalty to one's company."),
			rl("しかし、過労による健康被害が社会問題となったことをきっかけに、働き方改革が進められている。", "However, after health damage from overwork became a social issue, work-style reform has been pushed forward."),
			rl("残業時間の上限規制が導入された一方で、仕事量が減らないまま時間だけが制限され、かえって負担が増えたという声も少なくない。", "While a cap on overtime was introduced, quite a few people say the burden actually grew, since only hours were limited while the workload stayed the same."),
			rl("制度を変えるだけでなく、業務そのものを見直さない限り、真の改革は実現しないだろう。", "Unless the work itself is re-examined, not just the rules, true reform probably won't happen."),
		},
		[]models.ReadingQuestion{
			rq("How were long hours viewed in the past?", "As proof of loyalty to the company", "As proof of loyalty to the company", "As a sign of poor efficiency", "As illegal", "As unusual"),
			rq("What triggered work-style reform?", "Health damage from overwork", "Health damage from overwork", "A fall in exports", "Pressure from other countries", "A shortage of offices"),
			rq("Why do some say the burden increased?", "Hours were capped but the workload stayed the same", "Hours were capped but the workload stayed the same", "Salaries were cut", "Holidays were cancelled", "Offices were closed"),
			rq("What does the writer conclude?", "The work itself must be re-examined", "The work itself must be re-examined", "The reform has already succeeded", "Overtime should be unlimited again", "Nothing can be done"),
		},
	)
	addReadingL(db, ja, jaN1, "「間」の美学 — The Aesthetics of Ma",
		"An essay on Japanese culture. Read, then answer.", 6, 28,
		[]models.ReadingLine{
			rl("日本の文化を語る上で欠かせないのが、「間」という概念である。", "An indispensable concept when discussing Japanese culture is ma (the space between)."),
			rl("能楽における沈黙、日本庭園の余白、会話の中の一瞬の沈黙——いずれも、何もないことにこそ意味を見いだす感性に支えられている。", "Silence in Noh theatre, empty space in Japanese gardens, a momentary pause in conversation — all rest on a sensibility that finds meaning precisely in nothingness."),
			rl("西洋の芸術が空間を埋めることで表現を豊かにしてきたとすれば、日本の芸術は空間を残すことで想像の余地を与えてきたと言えよう。", "If Western art enriched expression by filling space, Japanese art, one might say, has left space to give room for the imagination."),
			rl("情報があふれる現代においてこそ、「間」の価値は改めて見直されるべきではないだろうか。", "Is it not precisely in our information-saturated age that the value of ma should be reconsidered?"),
		},
		[]models.ReadingQuestion{
			rq("What is 'ma' as described here?", "Meaningful space or pause", "Meaningful space or pause", "A kind of garden plant", "A Noh mask", "A type of poem"),
			rq("Which is given as an example of ma?", "A momentary silence in conversation", "A momentary silence in conversation", "A crowded festival", "A loud drum performance", "A detailed painting"),
			rq("How does the writer contrast Japanese and Western art?", "Japanese art leaves space; Western art fills it", "Japanese art leaves space; Western art fills it", "Both always fill space", "Western art has no meaning", "Japanese art copies Western art"),
			rq("What is the writer's final suggestion?", "Ma deserves to be re-valued today", "Ma deserves to be re-valued today", "Ma is outdated", "Gardens should be filled", "Silence should be avoided"),
		},
	)
	addReadingL(db, ja, jaC2, "方丈記 冒頭 — The Opening of Hōjōki",
		"A classical text (1212) with a modern translation. Read, then answer.", 7, 30,
		[]models.ReadingLine{
			rl("ゆく河の流れは絶えずして、しかももとの水にあらず。", "The flow of the river never ceases, and yet the water is never the same."),
			rl("よどみに浮かぶうたかたは、かつ消えかつ結びて、久しくとどまりたるためしなし。", "The bubbles floating on its pools vanish and form again, never lingering long."),
			rl("世の中にある人とすみかと、またかくのごとし。", "So it is with people in this world and their dwellings."),
		},
		[]models.ReadingQuestion{
			rq("What is the main theme of this passage?", "Impermanence", "Impermanence", "Love", "War", "Wealth"),
			rq("What does あらず mean here?", "is not", "is not", "is", "was", "will be"),
			rq("What are compared to the bubbles?", "People and their dwellings", "People and their dwellings", "Fish in the river", "The seasons", "Mountains"),
			rq("This text belongs to which kind of Japanese?", "Classical Japanese (文語)", "Classical Japanese (文語)", "Modern business Japanese", "Kansai dialect", "Okinawan"),
		},
	)
}
