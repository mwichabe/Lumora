package database

import "gorm.io/gorm"

// Japanese course, middle levels: N4 (A2) and N3 (B1). See seed_ja.go.

// ─────────────── N4 · A2 — Everyday Japanese ───────────────

func seedJapaneseN4(db *gorm.DB) {
	const ja = "ja"
	const u = "A2 · JLPT N4 — Everyday Japanese"
	finch := "Professor Finch"

	s := addSkillL(db, ja, u, "Plain Form & Casual Speech", "Dictionary, ない, た and なかった forms — talking with friends.", "MessageCircle", "#6C3FC5", 17, 180)
	l := addLesson(db, s, "Plain Forms", 1, 18,
		char(finch, "Friends and family speak in plain form, and grammar builds on it. Group 2: たべる → たべない → たべた → たべなかった. Group 1 negative changes う-sound to あ: かく → かかない, のむ → のまない, かう → かわない (う becomes わ). Past mirrors the て-form: かいた, のんだ. Irregular: する → しない/した, くる → こない/きた. And ある → ない."),
		mc("Plain negative of のむ", "のむ (drink)", "のまない", "のまない", "のみない", "のむない", "のめない"),
		mc("Plain negative of かう", "かう (buy)", "かわない", "かわない", "かあない", "かいない", "かうない"),
		mc("Plain past of よむ", "よむ (read)", "よんだ", "よんだ", "よみた", "よった", "よむた"),
		mc("Plain negative of くる", "くる (come)", "こない", "こない", "くない", "きない", "くらない"),
		mc("Plain negative of ある", "ある (there is)", "ない", "ない", "あらない", "ありない", "あない"),
		mc("Casual version", "きのう、なにをしましたか。", "きのう、なにした？", "きのう、なにした？", "きのう、なにするか？", "きのう、なにしない？", "きのう、なにしよう？"),
		speak("Blaze", "きょう、なにする？えいがみない？"),
	)
	addVocab(db, l,
		jv("〜ない", "〜nai", "plain negative", "きょうはいかない。", "kyō wa ikanai", "I'm not going today.", "Blaze"),
		jv("〜た", "〜ta", "plain past", "もうたべた？", "mō tabeta?", "Have you eaten already?", "Cora"),
		jv("友達", "ともだち tomodachi", "friend", "友達とあそぶ。", "ともだちとあそぶ。", "I hang out with friends.", "Pip"),
		jv("あそぶ", "asobu", "to play, hang out", "こうえんであそんだ。", "kōen de asonda", "We played in the park.", "Pip"),
	)

	s = addSkillL(db, ja, u, "The て-Form in Use", "Permission, prohibition and sequences.", "Shield", "#00C2A8", 18, 195)
	l = addLesson(db, s, "May I…? You mustn't…", 1, 18,
		char(finch, "〜てもいいですか asks permission: しゃしんをとってもいいですか (may I take a photo?). 〜てはいけません forbids: ここでたばこをすってはいけません. 〜てから means 'after doing': てをあらってからたべます. 〜なくてもいいです: you don't have to."),
		mc("May I sit here?", "ここにすわって＿ですか。", "もいい", "もいい", "はいけない", "から", "ください"),
		mc("You mustn't take photos here.", "ここでしゃしんをとって＿。", "はいけません", "はいけません", "もいいです", "ください", "からです"),
		mc("What does this mean?", "しゅくだいをしてから、テレビをみます。", "I'll watch TV after doing my homework.", "I'll watch TV after doing my homework.", "I'll do homework while watching TV.", "I watched TV instead of homework.", "Please do homework and watch TV."),
		mc("What does this mean?", "あしたはこなくてもいいです。", "You don't have to come tomorrow.", "You don't have to come tomorrow.", "You must come tomorrow.", "Don't come tomorrow.", "Can you come tomorrow?"),
		mc("Polite refusal of permission", "ここでたべてもいいですか。", "すみません、ちょっと…。", "すみません、ちょっと…。", "はい、たべてはいけません。", "いいえ、たべてもいいです。", "どういたしまして。"),
		speak("Riko", "すみません、しゃしんをとってもいいですか。"),
	)
	addVocab(db, l,
		jv("写真", "しゃしん shashin", "photo", "写真をとります。", "しゃしんをとります。", "I'll take a photo.", "Zephyr"),
		jv("すわる", "suwaru", "to sit", "ここにすわってもいいですか。", "koko ni suwatte mo ii desu ka", "May I sit here?", "Riko"),
		jv("〜てはいけない", "〜te wa ikenai", "must not", "はしってはいけない。", "hashitte wa ikenai", "You mustn't run.", finch),
		jv("〜なくてもいい", "〜nakute mo ii", "don't have to", "いそがなくてもいいです。", "isoganakute mo ii desu", "You don't have to hurry.", "Nana"),
	)

	s = addSkillL(db, ja, u, "Wants, Plans & Intentions", "たい, ほしい, つもり and the volitional form.", "Compass", "#F5A623", 19, 210)
	l = addLesson(db, s, "I Want To…, I Plan To…", 1, 18,
		char(finch, "〜たい (verb stem + たい) is 'want to do': にほんへいきたいです. ほしい is 'want a thing': あたらしいくるまがほしいです. Plans: dictionary form + つもりです (I intend to). The volitional form — たべよう, いこう — means 'let's' casually, and with とおもっています, 'I'm thinking of…'."),
		mc("I want to eat sushi.", "すしを＿です。", "たべたい", "たべたい", "たべほしい", "たべるたい", "たべよう"),
		mc("I want a new phone.", "あたらしいスマホが＿です。", "ほしい", "ほしい", "たい", "つもり", "よう"),
		mc("Volitional of いく", "いく (go)", "いこう", "いこう", "いきよう", "いくよう", "いけよう"),
		mc("What does this mean?", "らいねん、にほんでべんきょうするつもりです。", "I plan to study in Japan next year.", "I plan to study in Japan next year.", "I studied in Japan last year.", "I want to stop studying.", "I can study in Japan."),
		mc("What does this mean?", "なつやすみにほっかいどうへいこうとおもっています。", "I'm thinking of going to Hokkaido in the summer holidays.", "I'm thinking of going to Hokkaido in the summer holidays.", "I went to Hokkaido last summer.", "I don't want to go to Hokkaido.", "Let's not go to Hokkaido."),
		speak("Zephyr", "いつかふじさんにのぼりたいです。"),
	)
	addVocab(db, l,
		jv("〜たい", "〜tai", "want to (do)", "ねたい。", "netai", "I want to sleep.", "Pip"),
		jv("ほしい", "hoshii", "want (a thing)", "じかんがほしい。", "jikan ga hoshii", "I want more time.", "Mira"),
		jv("つもり", "tsumori", "intention, plan", "あしたいくつもりです。", "ashita iku tsumori desu", "I plan to go tomorrow.", finch),
		jv("来年", "らいねん rainen", "next year", "来年、日本へ行きます。", "らいねん、にほんへいきます。", "I'm going to Japan next year.", "Zephyr"),
	)

	s = addSkillL(db, ja, u, "Giving & Receiving", "あげる, くれる, もらう — and doing favours.", "Heart", "#FF5C5C", 20, 225)
	l = addLesson(db, s, "Who Gives to Whom", 1, 20,
		char(finch, "Giving depends on direction. あげる: I (or anyone) give outward — わたしはともだちにほんをあげました. くれる: someone gives to me or my group — ともだちがわたしにほんをくれました. もらう: I receive — ともだちにほんをもらいました. Add them to て-form for favours: てつだってくれました (they helped me)."),
		mc("My friend gave me a present.", "ともだちがプレゼントを＿。", "くれました", "くれました", "あげました", "もらいました", "やりました"),
		mc("I gave my mother flowers.", "ははにはなを＿。", "あげました", "あげました", "くれました", "もらいました", "ください"),
		mc("I received chocolate from Tanaka.", "たなかさんにチョコレートを＿。", "もらいました", "もらいました", "くれました", "あげました", "いただけます"),
		mc("What does this mean?", "せんせいがしゅくだいをてつだってくれました。", "The teacher helped me with my homework.", "The teacher helped me with my homework.", "I helped the teacher.", "The teacher gave me homework.", "I asked the teacher for help."),
		mc("What does this mean?", "あにになおしてもらいました。", "I had my older brother fix it.", "I had my older brother fix it.", "I fixed it for my brother.", "My brother broke it.", "I gave it to my brother."),
		speak("Mira", "たんじょうびに、ははがケーキをつくってくれました。"),
	)
	addVocab(db, l,
		jv("あげる", "ageru", "to give (outward)", "いもうとにペンをあげた。", "imōto ni pen o ageta", "I gave my little sister a pen.", "Mira"),
		jv("くれる", "kureru", "to give (to me)", "ちちがとけいをくれた。", "chichi ga tokei o kureta", "My father gave me a watch.", "Zephyr"),
		jv("もらう", "morau", "to receive", "てがみをもらった。", "tegami o moratta", "I got a letter.", "Nana"),
		jv("手伝う", "てつだう tetsudau", "to help", "手伝ってくれてありがとう。", "てつだってくれてありがとう。", "Thanks for helping me.", "Lumora"),
	)

	s = addSkillL(db, ja, u, "Comparisons", "より, のほうが, いちばん, ほど.", "Scale", "#17A3DD", 21, 240)
	l = addLesson(db, s, "Bigger, Better, Best", 1, 18,
		char(finch, "A は B より〜: A is more ~ than B — とうきょうはおおさかよりおおきいです. Asking: AとBとどちらが〜ですか → Aのほうが〜です. Superlative: 〜のなかで、〜がいちばん〜です. Not as ~ as: AはBほど〜ない."),
		mc("Tokyo is bigger than Osaka.", "とうきょうはおおさか＿おおきいです。", "より", "より", "ほど", "のほう", "いちばん"),
		mc("Which is more expensive, tea or coffee?", "おちゃとコーヒーと、＿がたかいですか。", "どちら", "どちら", "どこ", "なに", "だれ"),
		mc("Coffee is (the more expensive one).", "コーヒー＿たかいです。", "のほうが", "のほうが", "より", "ほど", "がいちばん"),
		mc("What does this mean?", "スポーツのなかで、サッカーがいちばんすきです。", "Of all sports, I like football best.", "Of all sports, I like football best.", "I like football more than other people.", "I don't like football much.", "Football is the oldest sport."),
		mc("What does this mean?", "きょうはきのうほどさむくないです。", "Today isn't as cold as yesterday.", "Today isn't as cold as yesterday.", "Today is colder than yesterday.", "Yesterday wasn't cold.", "It's cold today and yesterday."),
		speak("Cora", "わたしはいぬよりねこのほうがすきです。"),
	)
	addVocab(db, l,
		jv("より", "yori", "than", "バスよりでんしゃがはやい。", "basu yori densha ga hayai", "Trains are faster than buses.", "Riko"),
		jv("いちばん", "ichiban", "the most, number one", "これがいちばんおいしい。", "kore ga ichiban oishii", "This is the most delicious.", "Cora"),
		jv("どちら", "dochira", "which (of two)", "どちらがいいですか。", "dochira ga ii desu ka", "Which one is better?", "Cora"),
		jv("ほど", "hodo", "(not) as … as", "思ったほどむずかしくない。", "おもったほどむずかしくない。", "It's not as hard as I thought.", "Mira"),
	)

	s = addSkillL(db, ja, u, "The Potential Form", "Say what you can and can't do.", "CheckCircle", "#9B51E0", 22, 255)
	l = addLesson(db, s, "Can You…?", 1, 18,
		char(finch, "Potential = 'can'. Group 2: る → られる (たべられる). Group 1: う-sound → え-sound + る (かく → かける, よむ → よめる, はなす → はなせる). Irregular: する → できる, くる → こられる. The object usually takes が: にほんごがはなせます."),
		mc("Potential of はなす", "はなす (speak)", "はなせる", "はなせる", "はなされる", "はなさる", "はなしられる"),
		mc("Potential of たべる", "たべる (eat)", "たべられる", "たべられる", "たべれる", "たべける", "たべさせる"),
		mc("Potential of する", "する (do)", "できる", "できる", "しれる", "される", "せる"),
		mc("I can read kanji.", "かんじ＿よめます。", "が", "が", "で", "に", "へ"),
		mc("What does this mean?", "およげません。", "I can't swim.", "I can't swim.", "I don't swim.", "Don't swim.", "I didn't swim."),
		speak("Lumora", "にほんごがすこしはなせます。"),
	)
	addVocab(db, l,
		jv("できる", "dekiru", "can do; be done", "りょうりができます。", "ryōri ga dekimasu", "I can cook.", "Cora"),
		jv("およぐ", "oyogu", "to swim", "うみでおよげる？", "umi de oyogeru?", "Can you swim in the sea?", "Blaze"),
		jv("漢字", "かんじ kanji", "kanji", "漢字がよめます。", "かんじがよめます。", "I can read kanji.", finch),
		jv("少し", "すこし sukoshi", "a little", "少しわかります。", "すこしわかります。", "I understand a little.", "Lumora"),
	)

	s = addSkillL(db, ja, u, "Conditionals", "たら, ば, と and なら — if and when.", "Link", "#00C2A8", 23, 270)
	l = addLesson(db, s, "If…, When…", 1, 20,
		char(finch, "〜たら: if/when (one-off events, very common) — あめがふったら、いきません. 〜ば: if (hypothetical conditions) — やすければかいます. 〜と: whenever (natural consequences) — はるになると、さくらがさきます. 〜なら: if that's the case (responding to context) — きょうとへいくなら、きんかくじがいいですよ."),
		mc("If it rains, I won't go.", "あめがふっ＿、いきません。", "たら", "たら", "ても", "ので", "から"),
		mc("When spring comes, cherry blossoms bloom (always).", "はるになる＿、さくらがさきます。", "と", "と", "なら", "ても", "まで"),
		mc("ば-form of やすい", "やすい (cheap)", "やすければ", "やすければ", "やすいば", "やすかれば", "やすくば"),
		mc("Responding to 'I'm going to Kyoto': If you're going to Kyoto…", "きょうとへいく＿、きんかくじがいいですよ。", "なら", "なら", "と", "ば", "たら"),
		mc("What does this mean?", "このボタンをおすと、ドアがあきます。", "When you press this button, the door opens.", "When you press this button, the door opens.", "Don't press this button.", "The door opened, so I pressed it.", "Press the button after the door opens."),
		speak("Zephyr", "じかんがあったら、いっしょにいきましょう。"),
	)
	addVocab(db, l,
		jv("〜たら", "〜tara", "if / when", "ついたら、でんわしてね。", "tsuitara, denwa shite ne", "Call me when you arrive.", "Nana"),
		jv("〜ば", "〜ba", "if (condition)", "はやくおきれば、まにあう。", "hayaku okireba, ma ni au", "If you get up early, you'll make it.", finch),
		jv("〜なら", "〜nara", "if it's the case that", "すしなら、あのみせがいい。", "sushi nara, ano mise ga ii", "If it's sushi, that shop is good.", "Cora"),
		jv("桜", "さくら sakura", "cherry blossom", "桜がきれいですね。", "さくらがきれいですね。", "The cherry blossoms are beautiful.", "Mira"),
	)

	s = addSkillL(db, ja, u, "Explaining with んです", "Reasons with から and ので, and the explanatory の.", "Quote", "#F5A623", 24, 285)
	l = addLesson(db, s, "Because…", 1, 18,
		char(finch, "Reasons: 〜から (subjective, direct) and 〜ので (softer, more polite) both mean 'because'. 〜んです (〜のです) explains or seeks an explanation and sounds natural in conversation: どうしたんですか — what's the matter? あたまがいたいんです — (it's that) I have a headache."),
		mc("Because I'm busy, I can't go (polite, soft).", "いそがしい＿、いけません。", "ので", "ので", "のに", "ても", "まで"),
		mc("What does this ask?", "どうしたんですか。", "What's the matter?", "What's the matter?", "What will you do?", "How do you do?", "Where did you go?"),
		mc("Explain: it's that I have a headache.", "あたまがいたい＿。", "んです", "んです", "ですか", "でした", "ましょう"),
		mc("What does this mean?", "あめだから、タクシーでいこう。", "It's raining, so let's take a taxi.", "It's raining, so let's take a taxi.", "Let's walk even though it's raining.", "It stopped raining, so let's go.", "The taxi is late because of rain."),
		mc("Which sounds more polite as an excuse to a boss?", "Reason for being late", "でんしゃがおくれたので…", "でんしゃがおくれたので…", "でんしゃがおくれたから！", "でんしゃだし。", "でんしゃのせいだ。"),
		speak("Nana", "かぜをひいたので、きょうはやすみます。"),
	)
	addVocab(db, l,
		jv("〜ので", "〜node", "because (polite)", "ねつがあるので、やすみます。", "netsu ga aru node, yasumimasu", "I have a fever, so I'll rest.", "Nana"),
		jv("〜から", "〜kara", "because", "おなかがすいたから、たべよう。", "onaka ga suita kara, tabeyō", "I'm hungry, so let's eat.", "Blaze"),
		jv("どうしたんですか", "dō shita n desu ka", "what's wrong?", "げんきがないね。どうしたんですか。", "genki ga nai ne. dō shita n desu ka", "You look down. What's wrong?", "Lumora"),
		jv("遅れる", "おくれる okureru", "to be late", "でんしゃが遅れました。", "でんしゃがおくれました。", "The train was late.", "Riko"),
	)

	s = addSkillL(db, ja, u, "Kanji (N4)", "About 300 kanji in total: movement, time and daily life.", "PenLine", "#6C3FC5", 25, 300)
	l = addLesson(db, s, "Moving & Daily Life", 1, 18,
		char(finch, "Kanji combine into words. 行 (go) → 行く いく, 銀行 ぎんこう (bank). 来 (come) → 来る くる, 来週 らいしゅう (next week). 食 (eat) → 食べる たべる, 食堂 しょくどう (canteen). 飲 (drink) → 飲む のむ. 見 (see) → 見る みる. Kun'yomi is often used alone with kana; on'yomi in compounds."),
		mc("How is 来週 read?", "来週 (next week)", "らいしゅう", "らいしゅう", "くるしゅう", "きしゅう", "らいしゅ"),
		mc("What does 飲む mean?", "飲む (のむ)", "to drink", "to drink", "to eat", "to read", "to rest"),
		mc("How is 銀行 read?", "銀行 (bank)", "ぎんこう", "ぎんこう", "ぎんいく", "かねこう", "ぎんぎょう"),
		mc("Which kanji means 'to see'?", "みる", "見", "見", "貝", "目", "兄"),
		mc("What does 食堂 mean?", "食堂 (しょくどう)", "canteen, dining hall", "canteen, dining hall", "grocery store", "kitchen knife", "breakfast"),
		mc("How is 電車 read?", "電車", "でんしゃ", "でんしゃ", "でんくるま", "てんしゃ", "でんじゃ"),
		speak("Mira", "来週、友達と映画を見に行きます。"),
	)
	addVocab(db, l,
		jv("行く", "いく iku", "to go", "学校へ行く。", "がっこうへいく。", "I go to school.", finch),
		jv("来週", "らいしゅう raishū", "next week", "来週はいそがしい。", "らいしゅうはいそがしい。", "I'm busy next week.", "Riko"),
		jv("食べる", "たべる taberu", "to eat", "朝ご飯を食べる。", "あさごはんをたべる。", "I eat breakfast.", "Cora"),
		jv("飲む", "のむ nomu", "to drink", "お茶を飲む。", "おちゃをのむ。", "I drink tea.", "Nana"),
		jv("電車", "でんしゃ densha", "train", "電車で行きます。", "でんしゃでいきます。", "I'll go by train.", "Riko"),
	)

	s = addSkillL(db, ja, u, "Travel & Transport", "Tickets, directions, hotels and the shinkansen.", "Plane", "#17A3DD", 26, 315)
	l = addLesson(db, s, "Getting Around Japan", 1, 18,
		char("Riko", "Directions: まっすぐいってください (go straight), みぎにまがって (turn right), ひだりにまがって (turn left), 〜のかど (at the corner of). Buying tickets: きょうとまで、おとなにまいおねがいします. Ask: このでんしゃは〜にとまりますか (does this train stop at ~?)."),
		mc("Turn right at the corner.", "かどをみぎに＿ください。", "まがって", "まがって", "まがり", "まがる", "まがった"),
		mc("What does this ask?", "このでんしゃはしんじゅくにとまりますか。", "Does this train stop at Shinjuku?", "Does this train stop at Shinjuku?", "Is this the train to Shinjuku tomorrow?", "How far is Shinjuku?", "Is Shinjuku closed?"),
		mc("Two adult tickets to Kyoto, please.", "きょうとまで、おとな＿おねがいします。", "にまい", "にまい", "ふたり", "ふたつ", "にほん"),
		mc("What does this mean?", "まっすぐいくと、えきがみえます。", "If you go straight, you'll see the station.", "If you go straight, you'll see the station.", "The station is behind you.", "Turn back at the station.", "The station is far away."),
		mc("At the hotel: I have a reservation.", "Checking in", "よやくしているんですが…", "よやくしているんですが…", "よやくをしましょう。", "よやくがほしいです。", "よやくしませんでした。"),
		listen("Listen and choose", "つぎはきょうと、きょうとです", "Next stop: Kyoto.", "Next stop: Kyoto.", "This train goes to Tokyo.", "Change trains at Kyoto.", "Kyoto is the last stop."),
		speak("Riko", "すみません、しんかんせんののりばはどこですか。"),
	)
	addVocab(db, l,
		jv("切符", "きっぷ kippu", "ticket", "切符を買います。", "きっぷをかいます。", "I'll buy a ticket.", "Riko"),
		jv("新幹線", "しんかんせん shinkansen", "bullet train", "新幹線ははやい。", "しんかんせんははやい。", "The shinkansen is fast.", "Riko"),
		jv("予約", "よやく yoyaku", "reservation", "ホテルを予約しました。", "ホテルをよやくしました。", "I booked a hotel.", "Zephyr"),
		jv("まっすぐ", "massugu", "straight ahead", "まっすぐ行ってください。", "まっすぐいってください。", "Please go straight.", "Riko"),
		jv("右・左", "みぎ・ひだり migi / hidari", "right / left", "右にまがってください。", "みぎにまがってください。", "Please turn right.", "Riko"),
	)

	s = addSkillL(db, ja, u, "Health & the Body", "At the doctor's: symptoms and advice.", "HeartPulse", "#FF5C5C", 27, 330)
	l = addLesson(db, s, "At the Clinic", 1, 20,
		char("Nana", "Describing symptoms: あたまがいたいです (my head hurts), ねつがあります (I have a fever), せきがでます (I'm coughing). Advice uses 〜たほうがいいです (you'd better…): ゆっくりやすんだほうがいいですよ. And 〜ないほうがいい: you'd better not."),
		mc("I have a fever.", "ねつが＿。", "あります", "あります", "います", "します", "なります"),
		mc("You'd better rest.", "やすんだ＿がいいですよ。", "ほう", "ほう", "ほど", "より", "まで"),
		mc("What does this mean?", "おさけはのまないほうがいいです。", "You'd better not drink alcohol.", "You'd better not drink alcohol.", "You should drink alcohol.", "Alcohol is better than medicine.", "I can't drink alcohol."),
		mc("Which body part is おなか?", "おなかがいたい", "stomach", "stomach", "head", "throat", "back"),
		mc("What does the doctor ask?", "いつからですか。", "Since when?", "Since when?", "How long will it take?", "When is your next visit?", "What time is it?"),
		write("Write a short message (about 40 characters, in Japanese) to your teacher: say you have a fever and will rest today.",
			"せんせい、すみません。きのうからねつがあるので、きょうはがっこうをやすみます。あしたはいきます。"),
	)
	addVocab(db, l,
		jv("熱", "ねつ netsu", "fever", "熱があります。", "ねつがあります。", "I have a fever.", "Nana"),
		jv("痛い", "いたい itai", "painful, hurts", "のどが痛いです。", "のどがいたいです。", "My throat hurts.", "Nana"),
		jv("薬", "くすり kusuri", "medicine", "薬を飲んでください。", "くすりをのんでください。", "Please take the medicine.", "Nana"),
		jv("病院", "びょういん byōin", "hospital, clinic", "病院へ行ったほうがいい。", "びょういんへいったほうがいい。", "You'd better go to the clinic.", "Nana"),
	)
}

// ─────────────── N3 · B1 — Intermediate Japanese ───────────────

func seedJapaneseN3(db *gorm.DB) {
	const ja = "ja"
	const u = "B1 · JLPT N3 — Intermediate Japanese"
	finch := "Professor Finch"

	s := addSkillL(db, ja, u, "The Passive Form", "Things done to you — including the 'suffering' passive.", "Shield", "#6C3FC5", 28, 360)
	l := addLesson(db, s, "Done To…", 1, 20,
		char(finch, "Passive: Group 2 る → られる (たべられる), Group 1 う-sound → あ + れる (かく → かかれる, ふむ → ふまれる), する → される, くる → こられる. The doer takes に: せんせいにほめられました (I was praised by the teacher). Japanese also has an adversative passive for things that happen to you: あめにふられました — I got rained on."),
		mc("Passive of よむ", "よむ (read)", "よまれる", "よまれる", "よめる", "よませる", "よまさせる"),
		mc("I was scolded by my mother.", "ははに＿。", "しかられました", "しかられました", "しかりました", "しからせました", "しかれました"),
		mc("What does this mean?", "でんしゃであしをふまれました。", "Someone stepped on my foot on the train.", "Someone stepped on my foot on the train.", "I stepped on someone's foot.", "I hurt my foot running for the train.", "My foot was fine on the train."),
		mc("What nuance does this have?", "あめにふられました。", "It rained on me (and I suffered for it).", "It rained on me (and I suffered for it).", "I made it rain.", "I was happy it rained.", "It might rain."),
		mc("What does this mean?", "このおてらはせんねんまえにたてられました。", "This temple was built 1,000 years ago.", "This temple was built 1,000 years ago.", "This temple will be built in 1,000 years.", "Someone tried to build this temple.", "This temple has 1,000 rooms."),
		speak("Zephyr", "このほんはせかいじゅうでよまれています。"),
	)
	addVocab(db, l,
		jv("叱る", "しかる shikaru", "to scold", "先生に叱られた。", "せんせいにしかられた。", "I was scolded by the teacher.", finch),
		jv("褒める", "ほめる homeru", "to praise", "上司に褒められた。", "じょうしにほめられた。", "I was praised by my boss.", "Mira"),
		jv("踏む", "ふむ fumu", "to step on", "足を踏まれた。", "あしをふまれた。", "Someone stepped on my foot.", "Riko"),
		jv("建てる", "たてる tateru", "to build", "家を建てる。", "いえをたてる。", "To build a house.", "Zephyr"),
	)

	s = addSkillL(db, ja, u, "Causative & Causative-Passive", "Make or let someone do — and being made to.", "Users", "#00C2A8", 29, 380)
	l = addLesson(db, s, "Make, Let & Be Made To", 1, 20,
		char(finch, "Causative: Group 2 る → させる (たべさせる), Group 1 → あ + せる (いく → いかせる), する → させる, くる → こさせる. It means 'make' or 'let': こどもにやさいをたべさせます. 〜させてください asks to be allowed: かえらせてください. Causative-passive = 'be made to': ざんぎょうさせられました (I was made to work overtime)."),
		mc("Causative of いく", "いく (go)", "いかせる", "いかせる", "いかれる", "いける", "いこさせる"),
		mc("The mother made her child eat vegetables.", "はははこどもにやさいを＿。", "たべさせました", "たべさせました", "たべられました", "たべさせられました", "たべました"),
		mc("Please let me go home early today.", "きょうははやく＿ください。", "かえらせて", "かえらせて", "かえられて", "かえって", "かえさせられて"),
		mc("What does this mean?", "ぶちょうにおさけをのまされました。", "My manager made me drink (against my will).", "My manager made me drink (against my will).", "I made my manager drink.", "My manager let me drink if I wanted.", "My manager didn't drink."),
		mc("What does this mean?", "むすこにすきなことをさせています。", "I let my son do what he likes.", "I let my son do what he likes.", "My son makes me do things.", "I was made to do things by my son.", "My son doesn't like anything."),
		speak("Mira", "すみません、わたしにやらせてください。"),
	)
	addVocab(db, l,
		jv("〜させる", "〜saseru", "make / let (someone) do", "子どもに本を読ませる。", "こどもにほんをよませる。", "I have my child read books.", "Nana"),
		jv("〜させてください", "〜sasete kudasai", "please let me…", "考えさせてください。", "かんがえさせてください。", "Please let me think about it.", "Mira"),
		jv("残業", "ざんぎょう zangyō", "overtime work", "残業させられた。", "ざんぎょうさせられた。", "I was made to work overtime.", "Blaze"),
		jv("部長", "ぶちょう buchō", "department manager", "部長に報告します。", "ぶちょうにほうこくします。", "I'll report to the manager.", finch),
	)

	s = addSkillL(db, ja, u, "Hearsay & Conjecture", "そうだ, らしい, ようだ, みたい — hedging like a native.", "Brain", "#F5A623", 30, 400)
	l = addLesson(db, s, "It Seems…, I Heard…", 1, 20,
		char(finch, "Japanese hedges constantly. Plain form + そうです = I heard that (hearsay): あしたはあめだそうです. Stem + そうです = looks like (from appearance): おいしそうです. らしい: apparently (based on what one hears; also 'typical of'). ようです / みたいです: it seems (based on evidence) — みたい is the casual one."),
		mc("I heard it will rain tomorrow.", "あしたはあめが＿そうです。", "ふる", "ふる", "ふり", "ふりそう", "ふった"),
		mc("That cake looks delicious!", "そのケーキ、おいし＿！", "そう", "そう", "らしい", "みたい", "よう"),
		mc("What does this mean?", "だれもいないようです。", "It seems nobody is here.", "It seems nobody is here.", "Nobody wants to be here.", "Everybody is here.", "I heard someone is here."),
		mc("What does らしい mean here?", "かれはおとこらしい。", "He's manly (typical of a man).", "He's manly (typical of a man).", "He seems to be a man.", "He looks like a woman.", "He's not a man."),
		mc("Casual 'seems like'", "かぜをひいた＿。", "みたい", "みたい", "ようです", "そうだ", "らしいです"),
		speak("Cora", "このラーメン、おいしそう！"),
	)
	addVocab(db, l,
		jv("〜そうだ", "〜sō da", "I heard that…; looks…", "あの店は高いそうだ。", "あのみせはたかいそうだ。", "I hear that shop is expensive.", "Cora"),
		jv("〜らしい", "〜rashii", "apparently; typical of", "来月結婚するらしい。", "らいげつけっこんするらしい。", "Apparently they're getting married next month.", "Mira"),
		jv("〜ようだ", "〜yō da", "it seems", "かぎがかかっているようだ。", "かぎがかかっているようだ。", "It seems to be locked.", "Riko"),
		jv("〜みたい", "〜mitai", "seems (casual)", "雨みたいだね。", "あめみたいだね。", "Looks like rain, huh.", "Blaze"),
	)

	s = addSkillL(db, ja, u, "Change, Effort & Decisions", "ようになる, ようにする, ことになる, ことにする.", "Leaf", "#17A3DD", 31, 420)
	l = addLesson(db, s, "Becoming & Deciding", 1, 20,
		char(finch, "〜ようになる: come to be able / start to (a change): にほんごがはなせるようになりました. 〜ようにする: make an effort to: まいにちあるくようにしています. 〜ことにする: decide (yourself): たばこをやめることにしました. 〜ことになる: it's been decided (by circumstance): らいげつてんきんすることになりました."),
		mc("I became able to read kanji.", "かんじがよめる＿。", "ようになりました", "ようになりました", "ようにしました", "ことにしました", "ことになりました"),
		mc("I try to go to bed early.", "はやくねる＿います。", "ようにして", "ようにして", "ようになって", "ことにして", "ことになって"),
		mc("I decided to quit smoking.", "たばこをやめる＿。", "ことにしました", "ことにしました", "ことになりました", "ようになりました", "ようにしました"),
		mc("What does this mean?", "らいげつおおさかへてんきんすることになりました。", "It's been decided that I'll transfer to Osaka next month.", "It's been decided that I'll transfer to Osaka next month.", "I decided to move to Osaka by myself.", "I transferred to Osaka last month.", "I'm trying to transfer to Osaka."),
		mc("What does this mean?", "あかちゃんがあるけるようになりました。", "The baby has started walking.", "The baby has started walking.", "The baby tried to walk.", "The baby decided to walk.", "The baby can't walk yet."),
		speak("Lumora", "毎日日本語を話すようにしています。"),
	)
	addVocab(db, l,
		jv("〜ようになる", "〜yō ni naru", "come to (be able to)", "泳げるようになった。", "およげるようになった。", "I've learned to swim.", "Blaze"),
		jv("〜ようにする", "〜yō ni suru", "make an effort to", "野菜を食べるようにする。", "やさいをたべるようにする。", "I make a point of eating vegetables.", "Nana"),
		jv("〜ことにする", "〜koto ni suru", "decide to", "留学することにした。", "りゅうがくすることにした。", "I've decided to study abroad.", "Zephyr"),
		jv("転勤", "てんきん tenkin", "job transfer", "転勤が決まった。", "てんきんがきまった。", "My transfer has been decided.", finch),
	)

	s = addSkillL(db, ja, u, "Keigo Foundations", "Polite, honorific and humble speech — the basics.", "GraduationCap", "#9B51E0", 32, 440)
	l = addLesson(db, s, "Three Levels of Respect", 1, 20,
		char(finch, "Keigo has three layers. 丁寧語 (polite): です/ます. 尊敬語 (honorific) raises the other person's actions: いらっしゃる (go/come/be), めしあがる (eat), おっしゃる (say), ごらんになる (see), or お+stem+になる. 謙譲語 (humble) lowers your own: まいる (go/come), いただく (eat/receive), もうす (say), はいけんする (see), or お+stem+する."),
		mc("Honorific of いる/いく/くる", "the customer's action", "いらっしゃる", "いらっしゃる", "まいる", "おる", "うかがう"),
		mc("Humble 'say' (my name is…)", "わたしはスミスと＿。", "もうします", "もうします", "おっしゃいます", "いいます", "めしあがります"),
		mc("Honorific: Please eat (to a guest).", "どうぞ＿ください。", "めしあがって", "めしあがって", "いただいて", "たべさせて", "まいって"),
		mc("Humble: I'll carry your bag.", "おにもつを＿。", "おもちします", "おもちします", "おもちになります", "もってください", "もたれます"),
		mc("Which is humble (謙譲語)?", "Choose the humble verb", "はいけんする", "はいけんする", "ごらんになる", "おっしゃる", "なさる"),
		mc("What does this mean?", "しゃちょうはもうおかえりになりました。", "The company president has already gone home.", "The company president has already gone home.", "I sent the president home.", "The president hasn't arrived yet.", "Please go home, sir."),
		speak("Mira", "はじめまして。ミラともうします。どうぞよろしくおねがいいたします。"),
	)
	addVocab(db, l,
		jv("いらっしゃる", "irassharu", "go / come / be (honorific)", "先生はいらっしゃいますか。", "せんせいはいらっしゃいますか。", "Is the teacher in?", finch),
		jv("召し上がる", "めしあがる meshiagaru", "eat / drink (honorific)", "どうぞ召し上がってください。", "どうぞめしあがってください。", "Please help yourself.", "Cora"),
		jv("申す", "もうす mōsu", "say (humble)", "田中と申します。", "たなかともうします。", "My name is Tanaka.", "Mira"),
		jv("いただく", "itadaku", "receive / eat (humble)", "お茶をいただきます。", "おちゃをいただきます。", "I'll gladly have some tea.", "Nana"),
		jv("拝見する", "はいけんする haiken suru", "see (humble)", "資料を拝見しました。", "しりょうをはいけんしました。", "I've looked at the materials.", finch),
	)

	s = addSkillL(db, ja, u, "Nuance: ばかり, ところ, わけ, はず", "Just did, about to, no wonder, should be.", "Quote", "#FF5C5C", 33, 460)
	l = addLesson(db, s, "Shades of Meaning", 1, 20,
		char(finch, "た + ばかり: just did (feels recent): きたばかりです. Dictionary + ところ: about to; ている + ところ: in the middle of; た + ところ: just finished. はず: should be (logical expectation): もうついたはずです. わけ: no wonder / reason: つかれるわけだ. わけではない: it's not that…"),
		mc("I just arrived in Japan.", "にほんにきた＿です。", "ばかり", "ばかり", "はず", "わけ", "ところで"),
		mc("I'm just about to leave.", "いまでかける＿です。", "ところ", "ところ", "ばかり", "はず", "わけ"),
		mc("He should have arrived by now.", "かれはもうついた＿です。", "はず", "はず", "ばかり", "ところ", "わけ"),
		mc("What does this mean?", "にくがきらいなわけではありません。", "It's not that I dislike meat.", "It's not that I dislike meat.", "I definitely dislike meat.", "That's why I dislike meat.", "I should dislike meat."),
		mc("What does this mean?", "まいにちざんぎょうしているんだから、つかれるわけだ。", "No wonder you're tired — you do overtime every day.", "No wonder you're tired — you do overtime every day.", "You shouldn't be tired.", "You're about to be tired.", "It's not that you're tired."),
		speak("Riko", "いま、えきについたところです。"),
	)
	addVocab(db, l,
		jv("〜たばかり", "〜ta bakari", "just did", "買ったばかりのかさ。", "かったばかりのかさ。", "An umbrella I just bought.", "Riko"),
		jv("〜ところ", "〜tokoro", "about to / in the middle of / just did", "今食べているところ。", "いまたべているところ。", "I'm in the middle of eating.", "Cora"),
		jv("〜はず", "〜hazu", "should (be), expected", "会議は三時に終わるはずだ。", "かいぎはさんじにおわるはずだ。", "The meeting should end at three.", finch),
		jv("〜わけではない", "〜wake de wa nai", "it's not that…", "嫌いなわけではない。", "きらいなわけではない。", "It's not that I dislike it.", "Mira"),
	)

	s = addSkillL(db, ja, u, "Kanji (N3) & Headlines", "About 650 kanji: society, news and set compounds.", "PenLine", "#6C3FC5", 34, 480)
	l = addLesson(db, s, "Reading the News", 1, 20,
		char(finch, "Headlines drop particles and lean on two-kanji compounds: 首相 しゅしょう (prime minister), 経済 けいざい (economy), 政府 せいふ (government), 事故 じこ (accident), 発表 はっぴょう (announcement), 増加 ぞうか (increase), 減少 げんしょう (decrease). Reading them is the bridge to authentic news."),
		mc("How is 経済 read?", "経済 (economy)", "けいざい", "けいざい", "きょうさい", "けいさい", "けいざ"),
		mc("What does 事故 mean?", "事故 (じこ)", "accident", "accident", "event", "incident report", "festival"),
		mc("What does this headline say?", "観光客、過去最多に", "Tourist numbers hit a record high", "Tourist numbers hit a record high", "Tourists fall to a record low", "Tourism tax introduced", "Tourists warned about the weather"),
		mc("What does 減少 mean?", "人口が減少している", "decrease", "decrease", "increase", "stability", "movement"),
		mc("How is 政府 read?", "政府 (government)", "せいふ", "せいふ", "まさふ", "せいぶ", "しょうふ"),
		mc("What does this headline say?", "首相、来月訪米へ", "The prime minister will visit the US next month", "The prime minister will visit the US next month", "The US president visits Japan", "The prime minister resigns next month", "A US delegation arrives"),
		speak("Zephyr", "政府は新しい経済政策を発表しました。"),
	)
	addVocab(db, l,
		jv("経済", "けいざい keizai", "economy", "日本の経済。", "にほんのけいざい。", "The Japanese economy.", finch),
		jv("政府", "せいふ seifu", "government", "政府が発表した。", "せいふがはっぴょうした。", "The government announced it.", finch),
		jv("事故", "じこ jiko", "accident", "事故で電車が遅れた。", "じこででんしゃがおくれた。", "The train was delayed by an accident.", "Riko"),
		jv("増加", "ぞうか zōka", "increase", "観光客が増加した。", "かんこうきゃくがぞうかした。", "Tourist numbers increased.", "Zephyr"),
		jv("発表", "はっぴょう happyō", "announcement; presentation", "結果を発表します。", "けっかをはっぴょうします。", "I'll announce the results.", "Mira"),
	)

	s = addSkillL(db, ja, u, "Work & Society", "Meetings, requests at work and Japanese workplace norms.", "Briefcase", "#00C2A8", 35, 500)
	l = addLesson(db, s, "At the Office", 1, 20,
		char("Mira", "At work, polite requests matter: 〜ていただけませんか (could you kindly…?), 〜てもよろしいでしょうか (may I…?). ほうれんそう — 報告 (report), 連絡 (inform), 相談 (consult) — is the golden rule of Japanese teamwork. おつかれさまです greets colleagues all day."),
		mc("Could you kindly check this document?", "このしりょうをチェックして＿。", "いただけませんか", "いただけませんか", "くれませんか", "もらいますか", "あげませんか"),
		mc("What is おつかれさまです used for?", "said to colleagues", "Greeting / thanking colleagues for their work", "Greeting / thanking colleagues for their work", "Apologising for being late", "Saying goodbye forever", "Ordering food"),
		mc("What does ほうれんそう stand for at work?", "報・連・相", "Report, inform, consult", "Report, inform, consult", "Spinach, rice, soup", "Plan, do, check", "Arrive, work, leave"),
		mc("May I leave a little early today?", "きょうはすこしはやくかえっても＿。", "よろしいでしょうか", "よろしいでしょうか", "いけません", "ください", "ませんか"),
		mc("What does this mean?", "かいぎはごご２じからにへんこうになりました。", "The meeting has been changed to 2 p.m.", "The meeting has been changed to 2 p.m.", "The meeting was cancelled.", "The meeting lasted two hours.", "There are two meetings this afternoon."),
		write("Write a short email (about 80 characters, in Japanese) to your manager asking to move tomorrow's meeting to the afternoon, with a reason.",
			"ぶちょう、おつかれさまです。あしたのかいぎですが、ごぜんちゅうにきゃくさまとのやくそくがあるので、ごごにへんこうしていただけませんか。よろしくおねがいいたします。"),
	)
	addVocab(db, l,
		jv("会議", "かいぎ kaigi", "meeting", "三時から会議があります。", "さんじからかいぎがあります。", "There's a meeting from three.", "Mira"),
		jv("資料", "しりょう shiryō", "documents, materials", "資料を準備します。", "しりょうをじゅんびします。", "I'll prepare the materials.", "Mira"),
		jv("お疲れ様です", "おつかれさまです otsukaresama desu", "thanks for your hard work", "お疲れ様です！", "おつかれさまです！", "Good work today!", "Blaze"),
		jv("相談", "そうだん sōdan", "consultation", "ちょっと相談があるんですが。", "ちょっとそうだんがあるんですが。", "Could I ask your advice on something?", "Mira"),
		jv("変更", "へんこう henkō", "change", "予定を変更しました。", "よていをへんこうしました。", "I changed the schedule.", finch),
	)
}
