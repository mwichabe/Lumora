package database

import "gorm.io/gorm"

// Japanese course, advanced levels: N2 (B2), N1 (C1) and beyond N1 (C2).
// See seed_ja.go.

// ─────────────── N2 · B2 — Upper-Intermediate Japanese ───────────────

func seedJapaneseN2(db *gorm.DB) {
	const ja = "ja"
	const u = "B2 · JLPT N2 — Upper-Intermediate Japanese"
	finch := "Professor Finch"

	s := addSkillL(db, ja, u, "Business Keigo in Full", "させていただく, ご〜いたす, うかがう — keigo that works.", "Briefcase", "#6C3FC5", 36, 560)
	l := addLesson(db, s, "Keigo at Work", 1, 22,
		char(finch, "Business keigo combines forms: humble お/ご + stem + する/いたす (ご案内いたします, I'll show you the way), honorific お + stem + ください (お待ちください). 〜させていただく humbly asks leave to do something — use it, but don't overuse it. うかがう humbly means 'visit / ask', おる replaces いる, ございます replaces あります. 承知しました (understood) is the right word for a superior; 了解しました can sound too casual."),
		mc("I'll show you to the meeting room (humble).", "かいぎしつへ＿。", "ごあんないいたします", "ごあんないいたします", "ごあんないになります", "あんないしてください", "ごあんないなさいます"),
		mc("Reply to your manager's instruction", "Understood. (to a superior)", "承知しました", "承知しました", "了解！", "わかったよ", "オッケーです"),
		mc("Humble: I'll visit you tomorrow at ten.", "あした十時に＿。", "うかがいます", "うかがいます", "いらっしゃいます", "おいでになります", "まいられます"),
		mc("Humble 'be': Tanaka is currently out.", "田中はただいま外出して＿。", "おります", "おります", "いらっしゃいます", "おいでです", "ございます"),
		mc("Which is correct?", "Speaking of your own company's boss to a client", "部長の山田は席を外しております。", "部長の山田は席を外しております。", "山田部長は席を外していらっしゃいます。", "山田部長様はいらっしゃいません。", "山田さんはお休みになっています。"),
		mc("Please take a seat (honorific request).", "どうぞ＿ください。", "おかけ", "おかけ", "すわって", "おすわりして", "かけられて"),
		speak("Mira", "本日はお忙しいところ、お越しいただきありがとうございます。"),
	)
	addVocab(db, l,
		jv("承知しました", "しょうちしました shōchi shimashita", "understood, certainly (formal)", "かしこまりました。承知しました。", "かしこまりました。しょうちしました。", "Certainly. Understood.", "Mira"),
		jv("伺う", "うかがう ukagau", "visit; ask (humble)", "明日伺います。", "あしたうかがいます。", "I'll visit tomorrow.", finch),
		jv("おる", "oru", "be (humble)", "ただいま席を外しております。", "ただいませきをはずしております。", "They're away from their desk right now.", "Mira"),
		jv("ございます", "gozaimasu", "there is / be (very polite)", "こちらにございます。", "こちらにございます。", "It's right here.", "Cora"),
		jv("〜させていただく", "〜sasete itadaku", "(humbly) do", "説明させていただきます。", "せつめいさせていただきます。", "Allow me to explain.", finch),
	)

	s = addSkillL(db, ja, u, "Formal Grammar I", "にとって, に対して, によって, について, として.", "Link", "#00C2A8", 37, 590)
	l = addLesson(db, s, "Relating Ideas", 1, 22,
		char(finch, "〜にとって: for / from the viewpoint of (わたしにとって、かぞくがいちばんだいじだ). 〜に対して: towards, against, or in contrast to. 〜によって: by (an agent), by means of, or depending on (ひとによってちがう). 〜について: about. 〜として: as (in the role of)."),
		mc("For me, family matters most.", "わたし＿、家族が一番大切だ。", "にとって", "にとって", "に対して", "によって", "について"),
		mc("Customs differ depending on the country.", "国＿習慣が違う。", "によって", "によって", "にとって", "に対して", "として"),
		mc("He works in Japan as an engineer.", "彼はエンジニア＿日本で働いている。", "として", "として", "について", "にとって", "によって"),
		mc("What does this mean?", "兄が活発なのに対して、弟はおとなしい。", "Whereas the older brother is lively, the younger one is quiet.", "Whereas the older brother is lively, the younger one is quiet.", "The older brother is opposed to the younger one.", "Both brothers are lively.", "The younger brother is livelier."),
		mc("Let's talk about the environment.", "環境問題＿話し合いましょう。", "について", "について", "にとって", "として", "に対して"),
		speak("Zephyr", "この問題について、皆さんの意見を聞かせてください。"),
	)
	addVocab(db, l,
		jv("〜にとって", "〜ni totte", "for, from the viewpoint of", "学生にとって大切なこと。", "がくせいにとってたいせつなこと。", "What matters to students.", finch),
		jv("〜に対して", "〜ni taishite", "towards; in contrast to", "質問に対して答える。", "しつもんにたいしてこたえる。", "To answer a question.", finch),
		jv("〜によって", "〜ni yotte", "by; depending on", "人によって考え方が違う。", "ひとによってかんがえかたがちがう。", "Ways of thinking differ from person to person.", "Mira"),
		jv("習慣", "しゅうかん shūkan", "custom, habit", "日本の習慣を学ぶ。", "にほんのしゅうかんをまなぶ。", "To learn Japanese customs.", "Zephyr"),
	)

	s = addSkillL(db, ja, u, "Formal Grammar II", "ざるを得ない, わけにはいかない, に違いない, ものの, 一方だ.", "Scale", "#F5A623", 38, 620)
	l = addLesson(db, s, "Obligation, Certainty & Concession", 1, 22,
		char(finch, "〜ざるを得ない: have no choice but to (行かざるを得ない; する → せざるを得ない). 〜わけにはいかない: can't (for social or moral reasons). 〜に違いない: must be, surely. 〜ものの: although (it's true that…). 〜一方だ: keeps on (getting worse/more)."),
		mc("I have no choice but to cancel.", "中止せ＿。", "ざるを得ない", "ざるを得ない", "わけにはいかない", "に違いない", "一方だ"),
		mc("I can't leave early — it's an important meeting.", "大事な会議なので、早く帰る＿。", "わけにはいかない", "わけにはいかない", "ざるを得ない", "に違いない", "ものの"),
		mc("The lights are off — they must be out.", "電気が消えている。留守＿。", "に違いない", "に違いない", "わけにはいかない", "ものの", "ざるを得ない"),
		mc("Prices just keep rising.", "物価は上がる＿。", "一方だ", "一方だ", "ものの", "に違いない", "わけだ"),
		mc("What does this mean?", "仕事は見つかったものの、給料が安い。", "I found a job, but the pay is low.", "I found a job, but the pay is low.", "Because I found a job, I'm rich.", "I can't find a job with good pay.", "The pay was low, so I quit."),
		speak("Blaze", "雨がひどいので、試合を中止せざるを得ません。"),
	)
	addVocab(db, l,
		jv("〜ざるを得ない", "〜zaru o enai", "have no choice but to", "認めざるを得ない。", "みとめざるをえない。", "I have to admit it.", finch),
		jv("〜わけにはいかない", "〜wake ni wa ikanai", "can't (socially)", "約束を破るわけにはいかない。", "やくそくをやぶるわけにはいかない。", "I can't break a promise.", "Mira"),
		jv("〜に違いない", "〜ni chigainai", "must be, surely", "彼が犯人に違いない。", "かれがはんにんにちがいない。", "He must be the culprit.", "Zephyr"),
		jv("物価", "ぶっか bukka", "prices (cost of living)", "物価が上がった。", "ぶっかがあがった。", "Prices have gone up.", "Cora"),
	)

	s = addSkillL(db, ja, u, "Kanji Compounds & the Newspaper", "About 1,000 kanji: readings, compounds and editorials.", "PenLine", "#FF5C5C", 39, 650)
	l = addLesson(db, s, "One Kanji, Many Readings", 1, 22,
		char(finch, "Common kanji have many readings. 生: せい in 生活 (life), しょう in 一生 (a lifetime), なま in 生ビール (draught beer), う(まれる) in 生まれる (be born), い(きる) in 生きる (live). Compounds use on'yomi; words with okurigana use kun'yomi. Newspaper staples: 影響 えいきょう (influence), 対策 たいさく (countermeasure), 重要 じゅうよう (important)."),
		mc("How is 一生 read?", "一生 (a lifetime)", "いっしょう", "いっしょう", "いちせい", "ひとなま", "いっせい"),
		mc("How is 生ビール read?", "生ビール", "なまビール", "なまビール", "せいビール", "いきビール", "しょうビール"),
		mc("What does 影響 mean?", "影響 (えいきょう)", "influence, effect", "influence, effect", "shadow", "echo", "image"),
		mc("What does 対策 mean?", "感染対策", "countermeasure", "countermeasure", "opposition party", "policy debate", "survey"),
		mc("How is 重要 read?", "重要 (important)", "じゅうよう", "じゅうよう", "おもよう", "ちょうよう", "じゅよう"),
		mc("What does this mean?", "少子高齢化が経済に与える影響", "The effect of the falling birthrate and ageing population on the economy", "The effect of the falling birthrate and ageing population on the economy", "Children's influence on elderly people", "An economic plan for young families", "The rising birthrate in cities"),
		speak("Zephyr", "政府は少子化対策を重要な課題としている。"),
	)
	addVocab(db, l,
		jv("影響", "えいきょう eikyō", "influence", "天気が売上に影響する。", "てんきがうりあげにえいきょうする。", "The weather affects sales.", finch),
		jv("対策", "たいさく taisaku", "countermeasure", "対策を考える。", "たいさくをかんがえる。", "To think up countermeasures.", finch),
		jv("重要", "じゅうよう jūyō", "important", "重要な会議。", "じゅうようなかいぎ。", "An important meeting.", "Mira"),
		jv("少子高齢化", "しょうしこうれいか shōshi kōreika", "falling birthrate & ageing society", "少子高齢化が進む。", "しょうしこうれいかがすすむ。", "The population keeps ageing.", "Zephyr"),
		jv("生活", "せいかつ seikatsu", "daily life, living", "東京の生活に慣れた。", "とうきょうのせいかつになれた。", "I've got used to life in Tokyo.", "Riko"),
	)

	s = addSkillL(db, ja, u, "Onomatopoeia", "擬音語・擬態語: the sound-words natives use constantly.", "Music", "#17A3DD", 40, 680)
	l = addLesson(db, s, "Doki-doki & Pera-pera", 1, 22,
		char(finch, "Japanese has thousands of sound-symbolic words. 擬音語 imitate sounds: ザーザー (pouring rain), ゴロゴロ (thunder rumbling). 擬態語 describe states and feelings: ドキドキ (heart pounding), ワクワク (excited), ペラペラ (fluent), ニコニコ (smiling), イライラ (irritated), キラキラ (sparkling), ゴロゴロ (lazing about)."),
		mc("Her Japanese is fluent.", "彼女は日本語が＿だ。", "ペラペラ", "ペラペラ", "ドキドキ", "ゴロゴロ", "ザーザー"),
		mc("My heart is pounding before the exam.", "試験の前で、＿している。", "ドキドキ", "ドキドキ", "ニコニコ", "キラキラ", "ペラペラ"),
		mc("What does イライラ describe?", "電車が来なくてイライラする。", "Irritation", "Irritation", "Excitement", "Sleepiness", "Happiness"),
		mc("What does this mean?", "雨がザーザー降っている。", "It's pouring with rain.", "It's pouring with rain.", "It's drizzling.", "The rain has stopped.", "It's snowing heavily."),
		mc("Lazing around at home on Sunday", "日曜日は家で＿していた。", "ゴロゴロ", "ゴロゴロ", "ワクワク", "ペラペラ", "ドキドキ"),
		mc("What does ワクワク describe?", "旅行が楽しみでワクワクする。", "Excited anticipation", "Excited anticipation", "Nervous fear", "Boredom", "Anger"),
		speak("Pip", "明日は遠足だから、ワクワクして眠れない！"),
	)
	addVocab(db, l,
		jv("ドキドキ", "dokidoki", "heart pounding (nerves, love)", "告白してドキドキした。", "こくはくしてドキドキした。", "My heart raced when I confessed.", "Mira"),
		jv("ワクワク", "wakuwaku", "excited", "ワクワクする話。", "ワクワクするはなし。", "An exciting story.", "Pip"),
		jv("ペラペラ", "perapera", "fluent; flimsy", "英語がペラペラだ。", "えいごがペラペラだ。", "They're fluent in English.", "Lumora"),
		jv("イライラ", "iraira", "irritated", "待たされてイライラした。", "またされてイライラした。", "I got irritated waiting.", "Blaze"),
		jv("ニコニコ", "nikoniko", "smiling", "いつもニコニコしている。", "いつもニコニコしている。", "They're always smiling.", "Nana"),
	)

	s = addSkillL(db, ja, u, "Business Email & Phone", "Set phrases for email, calls and apologies.", "Smartphone", "#00C2A8", 41, 710)
	l = addLesson(db, s, "Professional Communication", 1, 24,
		char("Mira", "Business emails open with お世話になっております and close with よろしくお願いいたします. Softeners: 恐れ入りますが (I'm sorry to trouble you, but), お手数ですが (sorry for the trouble). On the phone: 少々お待ちください (one moment please), 折り返しお電話いたします (I'll call you back), 申し訳ございません (I'm very sorry)."),
		mc("Standard opening line of a business email", "To a client you already work with", "いつもお世話になっております。", "いつもお世話になっております。", "はじめまして、元気？", "お疲れ様でした。", "おはようございます、ご苦労様です。"),
		mc("I'll call you back.", "後ほど＿。", "折り返しお電話いたします", "折り返しお電話いたします", "電話してください", "電話をもらいます", "電話してあげます"),
		mc("Polite softener before a request", "＿、ご確認いただけますでしょうか。", "お手数ですが", "お手数ですが", "どうでもいいですが", "ちなみに", "ところで"),
		mc("The most formal apology", "Apologising to a customer", "大変申し訳ございません。", "大変申し訳ございません。", "ごめんね。", "すまない。", "悪かったです。"),
		mc("Why is ご苦労様です risky to say to your boss?", "ご苦労様です", "It's traditionally used toward subordinates", "It's traditionally used toward subordinates", "It's too formal", "It means goodbye", "It's only used at night"),
		write("Write a short business email (about 120 characters, in Japanese) to a client apologising that the delivery will be two days late, and asking for their understanding.",
			"株式会社山田商事 山田様\nいつもお世話になっております。\nご注文いただいた商品ですが、配送の遅れにより、お届けが二日ほど遅れる見込みです。ご迷惑をおかけし、大変申し訳ございません。何卒ご理解のほど、よろしくお願いいたします。"),
	)
	addVocab(db, l,
		jv("お世話になっております", "おせわになっております osewa ni natte orimasu", "thank you for your continued support", "いつもお世話になっております。", "いつもおせわになっております。", "Thank you, as always, for your support.", "Mira"),
		jv("恐れ入りますが", "おそれいりますが osoreirimasu ga", "I'm sorry to trouble you, but", "恐れ入りますが、お名前をお願いします。", "おそれいりますが、おなまえをおねがいします。", "Sorry to trouble you — may I have your name?", "Cora"),
		jv("申し訳ございません", "もうしわけございません mōshiwake gozaimasen", "I'm deeply sorry", "ご迷惑をおかけして申し訳ございません。", "ごめいわくをおかけしてもうしわけございません。", "I'm very sorry for the trouble.", "Mira"),
		jv("折り返し", "おりかえし orikaeshi", "(calling) back", "折り返しご連絡します。", "おりかえしごれんらくします。", "I'll get back to you.", "Riko"),
	)
}

// ─────────────── N1 · C1 — Advanced Japanese ───────────────

func seedJapaneseN1(db *gorm.DB) {
	const ja = "ja"
	const u = "C1 · JLPT N1 — Advanced Japanese"
	finch := "Professor Finch"

	s := addSkillL(db, ja, u, "Advanced Grammar I", "をもって, ならでは, べく, をよそに, を皮切りに.", "Layers", "#6C3FC5", 42, 820)
	l := addLesson(db, s, "Formal Written Patterns", 1, 25,
		char(finch, "N1 grammar lives in formal writing. 〜をもって: with / as of (本日をもって閉店します). 〜ならでは: unique to (京都ならではの景色). 〜べく: in order to (試験に合格すべく). 〜をよそに: in disregard of (反対をよそに). 〜を皮切りに: starting with (東京を皮切りに全国で)."),
		mc("The shop closes as of today.", "本日＿閉店いたします。", "をもって", "をもって", "ならでは", "をよそに", "べく"),
		mc("Scenery you can only see in Kyoto", "京都＿の景色", "ならでは", "ならでは", "をもって", "べく", "を皮切りに"),
		mc("He studied every day in order to pass.", "合格す＿、毎日勉強した。", "べく", "べく", "ならでは", "をよそに", "をもって"),
		mc("Despite residents' opposition, construction began.", "住民の反対＿、工事が始まった。", "をよそに", "をよそに", "を皮切りに", "べく", "ならでは"),
		mc("What does this mean?", "東京公演を皮切りに、全国ツアーが始まる。", "Starting with the Tokyo show, a nationwide tour begins.", "Starting with the Tokyo show, a nationwide tour begins.", "The tour ends in Tokyo.", "The Tokyo show was cancelled.", "Only Tokyo will have a show."),
		speak("Zephyr", "本日をもちまして、当店は閉店いたします。"),
	)
	addVocab(db, l,
		jv("〜をもって", "〜o motte", "with; as of", "以上をもって終わります。", "いじょうをもっておわります。", "With that, we conclude.", finch),
		jv("〜ならでは", "〜naradewa", "unique to", "日本ならではの文化。", "にほんならではのぶんか。", "Culture unique to Japan.", "Zephyr"),
		jv("〜べく", "〜beku", "in order to", "夢を叶えるべく努力する。", "ゆめをかなえるべくどりょくする。", "Strive to make a dream come true.", "Blaze"),
		jv("〜をよそに", "〜o yoso ni", "in disregard of", "心配をよそに出かけた。", "しんぱいをよそにでかけた。", "They went out despite everyone's worry.", "Nana"),
	)

	s = addSkillL(db, ja, u, "Advanced Grammar II", "いかんによらず, に即して, ずにはおかない, とあって, といえども.", "Scale", "#00C2A8", 43, 850)
	l = addLesson(db, s, "Nuanced Conditions", 1, 25,
		char(finch, "〜いかんによらず: regardless of (理由のいかんによらず). 〜に即して: in accordance with (事実に即して). 〜ずにはおかない: will inevitably cause (感動させずにはおかない). 〜とあって: because it's (a special case) — 連休とあって混雑している. 〜といえども: even though / even (専門家といえども)."),
		mc("Regardless of the reason, late submissions are not accepted.", "理由の＿、遅れた提出は受け付けない。", "いかんによらず", "いかんによらず", "に即して", "とあって", "をもって"),
		mc("Judge in accordance with the facts.", "事実＿判断する。", "に即して", "に即して", "をよそに", "ならでは", "とあって"),
		mc("Because it's a long weekend, it's very crowded.", "連休＿、どこも混雑している。", "とあって", "とあって", "といえども", "に即して", "べく"),
		mc("Even experts make mistakes.", "専門家＿、間違えることはある。", "といえども", "といえども", "とあって", "ならでは", "をもって"),
		mc("What does this mean?", "この映画は観る者を感動させずにはおかない。", "This film cannot fail to move its viewers.", "This film cannot fail to move its viewers.", "This film doesn't move anyone.", "Viewers must not be moved.", "This film was made to annoy viewers."),
		speak("Mira", "理由のいかんによらず、規則は守らなければならない。"),
	)
	addVocab(db, l,
		jv("〜いかんによらず", "〜ikan ni yorazu", "regardless of", "結果のいかんによらず報告する。", "けっかのいかんによらずほうこくする。", "Report regardless of the outcome.", finch),
		jv("〜に即して", "〜ni sokushite", "in accordance with", "現実に即して考える。", "げんじつにそくしてかんがえる。", "Think in line with reality.", finch),
		jv("〜とあって", "〜to atte", "because it's (special)", "初日とあって大勢の客が来た。", "しょにちとあっておおぜいのきゃくがきた。", "Being opening day, crowds came.", "Riko"),
		jv("〜といえども", "〜to iedomo", "even though, even", "子供といえども容赦しない。", "こどもといえどもようしゃしない。", "No mercy, even for a child.", "Zephyr"),
	)

	s = addSkillL(db, ja, u, "Four-Character Idioms & Proverbs", "四字熟語 and ことわざ.", "BookOpen", "#F5A623", 44, 880)
	l = addLesson(db, s, "Wisdom in Four Characters", 1, 25,
		char(finch, "四字熟語 pack ideas into four kanji: 一石二鳥 (two birds with one stone), 一期一会 (treasure each meeting — it may never recur), 十人十色 (to each their own), 自業自得 (you reap what you sow). Proverbs: 猿も木から落ちる (even monkeys fall from trees), 七転び八起き (fall seven times, rise eight), 石の上にも三年 (perseverance prevails)."),
		mc("What does 一石二鳥 mean?", "一石二鳥", "Killing two birds with one stone", "Killing two birds with one stone", "A bird in the hand", "One stone for every bird", "A once-in-a-lifetime chance"),
		mc("What does 一期一会 express?", "一期一会", "Treasure every encounter — it may never happen again", "Treasure every encounter — it may never happen again", "Meet once a year", "First come, first served", "Every meeting is boring"),
		mc("Even experts make mistakes.", "Which proverb?", "猿も木から落ちる", "猿も木から落ちる", "石の上にも三年", "七転び八起き", "十人十色"),
		mc("What does 十人十色 mean?", "十人十色", "Everyone has different tastes", "Everyone has different tastes", "Ten people make a team", "Colourful clothes", "A crowd of ten"),
		mc("He failed after never studying — it's his own fault.", "勉強しなかったのだから、＿だ。", "自業自得", "自業自得", "一期一会", "一石二鳥", "十人十色"),
		mc("Never give up after failure", "Which proverb?", "七転び八起き", "七転び八起き", "猿も木から落ちる", "自業自得", "一期一会"),
		speak("Nana", "七転び八起き。あきらめないで頑張りましょう。"),
	)
	addVocab(db, l,
		jv("一石二鳥", "いっせきにちょう isseki nichō", "two birds with one stone", "運動しながら通勤できて一石二鳥だ。", "うんどうしながらつうきんできていっせきにちょうだ。", "Exercise while commuting — two birds with one stone.", finch),
		jv("一期一会", "いちごいちえ ichigo ichie", "once-in-a-lifetime encounter", "一期一会を大切に。", "いちごいちえをたいせつに。", "Cherish every encounter.", "Nana"),
		jv("十人十色", "じゅうにんといろ jūnin toiro", "to each their own", "好みは十人十色だ。", "このみはじゅうにんといろだ。", "Tastes differ from person to person.", "Mira"),
		jv("自業自得", "じごうじとく jigō jitoku", "you reap what you sow", "それは自業自得だよ。", "それはじごうじとくだよ。", "That's on you.", "Blaze"),
		jv("七転び八起き", "ななころびやおき nanakorobi yaoki", "fall seven times, rise eight", "人生は七転び八起きだ。", "じんせいはななころびやおきだ。", "Life is about getting back up.", "Nana"),
	)

	s = addSkillL(db, ja, u, "Editorials & Abstract Argument", "Read and write opinion pieces in である style.", "Quote", "#FF5C5C", 45, 910)
	l = addLesson(db, s, "Arguing in Writing", 1, 25,
		char(finch, "Editorials use plain である style and hedged argument: 〜と言わざるを得ない (one must say), 〜のではないだろうか (isn't it the case that…?), 〜にほかならない (is nothing other than), 〜とは言い難い (can hardly be called), 〜べきである (ought to)."),
		mc("Which style do editorials use?", "Formal written Japanese", "だ・である体", "だ・である体", "です・ます体", "Casual speech", "Kansai dialect"),
		mc("What does this hedge mean?", "見直す必要があるのではないだろうか。", "Isn't there a need to reconsider? (suggesting there is)", "Isn't there a need to reconsider? (suggesting there is)", "There's definitely no need to reconsider.", "We already reconsidered.", "I don't know whether to reconsider."),
		mc("Success is nothing other than the result of effort.", "成功は努力の結果＿。", "にほかならない", "にほかならない", "とは言い難い", "べきである", "に違いない"),
		mc("This policy can hardly be called effective.", "この政策は効果的＿。", "とは言い難い", "とは言い難い", "にほかならない", "べきである", "ざるを得ない"),
		mc("The government ought to act immediately.", "政府は直ちに対応す＿。", "べきである", "べきである", "とは言い難い", "にほかならない", "わけだ"),
		write("Write an opinion paragraph (about 200 characters, in Japanese, である style) on whether remote work should become the norm. Give one argument for, one against, and your conclusion.",
			"在宅勤務を標準とすべきかについては意見が分かれる。通勤時間がなくなり、生産性や生活の質が向上するという利点は大きい。一方で、対面での交流が減り、新人の育成が難しくなるという問題も無視できない。したがって、全面的な在宅勤務ではなく、職種に応じて出社と組み合わせる柔軟な働き方こそが望ましいのではないだろうか。"),
	)
	addVocab(db, l,
		jv("〜べきである", "〜beki de aru", "ought to", "早急に対策を取るべきである。", "さっきゅうにたいさくをとるべきである。", "Measures should be taken urgently.", finch),
		jv("〜にほかならない", "〜ni hoka naranai", "is nothing but", "これは偶然にほかならない。", "これはぐうぜんにほかならない。", "This is nothing but chance.", finch),
		jv("〜とは言い難い", "〜to wa iigatai", "can hardly be said", "十分とは言い難い。", "じゅうぶんとはいいがたい。", "It can hardly be called enough.", "Mira"),
		jv("一方で", "いっぽうで ippō de", "on the other hand", "一方で、問題もある。", "いっぽうで、もんだいもある。", "On the other hand, there are problems too.", "Zephyr"),
	)

	s = addSkillL(db, ja, u, "Lectures & Academic Japanese", "Follow lectures and read academic prose.", "GraduationCap", "#17A3DD", 46, 940)
	l = addLesson(db, s, "In the Lecture Hall", 1, 25,
		char(finch, "Academic Japanese favours set expressions: 〜とされる (it is held that), 〜と考えられる (it is thought that), 〜に起因する (be caused by), 〜を示唆する (suggest), 本稿では (in this paper), 以上のことから (from the above). Lecturers signpost: まず、次に、最後に、つまり."),
		mc("What does 本稿では mean?", "本稿では〜について論じる。", "In this paper", "In this paper", "In this book's sequel", "In the last chapter", "In the newspaper"),
		mc("This accident is attributed to human error.", "この事故は人為的ミス＿。", "に起因する", "に起因する", "を示唆する", "とされる", "にほかならない"),
		mc("What does 示唆する mean?", "結果は新たな可能性を示唆している。", "to suggest, imply", "to suggest, imply", "to deny", "to prove conclusively", "to ignore"),
		mc("Which phrase signals a conclusion?", "Wrapping up an argument", "以上のことから", "以上のことから", "まず", "例えば", "ところで"),
		mc("What does this mean?", "この現象は気候変動が原因だと考えられる。", "This phenomenon is thought to be caused by climate change.", "This phenomenon is thought to be caused by climate change.", "Climate change is caused by this phenomenon.", "This phenomenon has nothing to do with climate.", "Climate change was disproved."),
		listen("Listen and choose", "まず、研究の背景について説明します", "First, I'll explain the background of the research.", "First, I'll explain the background of the research.", "Finally, here are the conclusions.", "Next, let's look at the data.", "Are there any questions?"),
		speak("Professor Finch", "以上のことから、本研究の仮説は支持されたと考えられる。"),
	)
	addVocab(db, l,
		jv("〜とされる", "〜to sareru", "is considered / held to be", "最古の記録とされる。", "さいこのきろくとされる。", "It is considered the oldest record.", finch),
		jv("示唆する", "しさする shisa suru", "to suggest, imply", "データが変化を示唆する。", "データがへんかをしさする。", "The data suggest a change.", finch),
		jv("〜に起因する", "〜にきいんする ni kiin suru", "be caused by", "ストレスに起因する病気。", "ストレスにきいんするびょうき。", "Illness caused by stress.", "Nana"),
		jv("仮説", "かせつ kasetsu", "hypothesis", "仮説を検証する。", "かせつをけんしょうする。", "To test a hypothesis.", finch),
	)
}

// ─────────────── Beyond N1 · C2 — Native-Level Japanese ───────────────

func seedJapaneseBeyondN1(db *gorm.DB) {
	const ja = "ja"
	const u = "C2 · Beyond N1 — Native-Level Japanese"
	finch := "Professor Finch"

	s := addSkillL(db, ja, u, "Classical Japanese (文語)", "Read the classics: けり, なり, ず, べし and む.", "Landmark", "#6C3FC5", 47, 1000)
	l := addLesson(db, s, "Reading the Classics", 1, 28,
		char(finch, "Classical Japanese (文語) still appears in literature, proverbs, legal phrases and poetry. Key auxiliaries: 〜けり (past, often with a sense of realisation), 〜なり (copula 'is'), 〜ず (negative), 〜べし (should / surely), 〜む (will / probably). Old kana spellings differ too: いふ is read いう, をかし is read おかし (charming)."),
		mc("What does 〜ず mean in classical Japanese?", "知らず (しらず)", "negative (not)", "negative (not)", "past tense", "question", "plural"),
		mc("How is the classical spelling いふ read today?", "いふ (to say)", "いう", "いう", "いふ", "いぶ", "ゆう only in Kansai"),
		mc("What does いとをかし mean?", "Makura no Sōshi: いとをかし", "Very charming / delightful", "Very charming / delightful", "Very funny / ridiculous", "Very strange", "Very sad"),
		mc("What does 〜べし express?", "学ぶべし", "should / must", "should / must", "did not", "might have", "is"),
		mc("The opening of The Tale of the Heike", "祇園精舎の鐘の声、諸行無常の響きあり", "The impermanence of all things", "The impermanence of all things", "A love story in the capital", "A comic tale of a monk", "Instructions for a festival"),
		mc("What does 月日は百代の過客にして describe? (Bashō)", "月日は百代の過客にして、行きかふ年もまた旅人なり", "Time as an eternal traveller", "Time as an eternal traveller", "The moon over a hundred guests", "A hotel full of travellers", "A calendar for the year"),
		speak("Nana", "祇園精舎の鐘の声、諸行無常の響きあり。"),
	)
	addVocab(db, l,
		jv("〜けり", "〜keri", "classical past (with realisation)", "昔、男ありけり。", "むかし、おとこありけり。", "Once, there was a man.", finch),
		jv("〜なり", "〜nari", "classical copula (is)", "旅人なり。", "たびびとなり。", "(It) is a traveller.", finch),
		jv("をかし", "おかし okashi", "charming, delightful (classical)", "いとをかし。", "いとおかし。", "Most delightful.", "Nana"),
		jv("諸行無常", "しょぎょうむじょう shogyō mujō", "impermanence of all things", "諸行無常の世の中。", "しょぎょうむじょうのよのなか。", "A world where nothing lasts.", "Zephyr"),
	)

	s = addSkillL(db, ja, u, "Literature & Style", "Haiku, kigo and literary technique.", "BookOpen", "#00C2A8", 48, 1030)
	l = addLesson(db, s, "Haiku & Literary Japanese", 1, 28,
		char(finch, "Haiku follow 5-7-5 morae and contain a seasonal word (季語) and often a cutting word (切れ字) like や or かな. 古池や 蛙飛びこむ 水の音 (Bashō): an old pond — a frog jumps in — the sound of water. 蛙 (frog) is a spring kigo. Literary techniques include 体言止め (ending a sentence on a noun) and 擬人法 (personification)."),
		mc("What is the mora pattern of a haiku?", "Haiku structure", "5-7-5", "5-7-5", "7-7-7", "5-5-5", "5-7-5-7-7"),
		mc("What is 季語?", "季語 (きご)", "A seasonal word", "A seasonal word", "A cutting word", "The poet's name", "A rhyme"),
		mc("In 古池や, what is や?", "古池や", "A cutting word (切れ字)", "A cutting word (切れ字)", "A question particle", "The word 'and'", "A season word"),
		mc("Which season does 蛙 (frog) signal?", "古池や 蛙飛びこむ 水の音", "Spring", "Spring", "Summer", "Autumn", "Winter"),
		mc("Which technique ends a sentence on a noun?", "静かな夜、遠くに響く鐘の音。", "体言止め", "体言止め", "擬人法", "倒置法", "反復法"),
		mc("5-7-5-7-7 poetry is called…", "百人一首 contains these", "短歌 (tanka)", "短歌 (tanka)", "俳句 (haiku)", "川柳 (senryū)", "詩 (free verse)"),
		write("Write your own haiku in Japanese (5-7-5 morae) that includes a seasonal word, then explain in one Japanese sentence which season it shows.",
			"夕立や 傘をたたんで 空を見る\nこの俳句は「夕立」という夏の季語を使っています。"),
	)
	addVocab(db, l,
		jv("俳句", "はいく haiku", "haiku", "俳句を詠む。", "はいくをよむ。", "To compose a haiku.", "Nana"),
		jv("季語", "きご kigo", "seasonal word", "季語は春です。", "きごははるです。", "The season word is spring.", finch),
		jv("短歌", "たんか tanka", "tanka (5-7-5-7-7)", "短歌を作る。", "たんかをつくる。", "To write a tanka.", "Nana"),
		jv("体言止め", "たいげんどめ taigendome", "ending on a noun", "体言止めで余韻を残す。", "たいげんどめでよいんをのこす。", "Ending on a noun leaves a lingering echo.", finch),
	)

	s = addSkillL(db, ja, u, "Dialects: Kansai & Beyond", "関西弁 and regional Japanese.", "Globe", "#F5A623", 49, 1060)
	l = addLesson(db, s, "Speaking Like Osaka", 1, 28,
		char("Blaze", "Kansai-ben (Osaka, Kyoto, Kobe) is everywhere on TV. ほんま = really (本当), あかん = no good / don't, おおきに = thanks, 〜へん = negative (わからへん = わからない), 〜や replaces だ (そうや = そうだ), なんでやねん = what on earth! (the classic comedy retort). Other regions: Hakata 〜ばい/〜けん (because), Okinawa めんそーれ (welcome)."),
		mc("What does あかん mean?", "それはあかん！", "No good / you mustn't", "No good / you mustn't", "Delicious", "Thank you", "Really?"),
		mc("Standard Japanese for わからへん", "わからへん", "わからない", "わからない", "わかる", "わかった", "わかります"),
		mc("What does ほんま mean?", "ほんまに？", "really", "really", "maybe", "never", "quickly"),
		mc("Kansai thanks", "Said in Kyoto and Osaka shops", "おおきに", "おおきに", "めんそーれ", "なんでやねん", "あかん"),
		mc("Standard Japanese for そうや", "そうや！", "そうだ", "そうだ", "そうですか", "そうじゃない", "そうしよう"),
		mc("What does めんそーれ mean, and where?", "めんそーれ", "Welcome — Okinawa", "Welcome — Okinawa", "Goodbye — Osaka", "Thank you — Kyoto", "Delicious — Hokkaido"),
		speak("Blaze", "ほんまに？そんなん、あかんやん！"),
	)
	addVocab(db, l,
		jv("ほんま", "honma", "really (Kansai)", "ほんまにおいしいわ。", "honma ni oishii wa", "It's really delicious.", "Blaze"),
		jv("あかん", "akan", "no good, mustn't (Kansai)", "遅刻したらあかんで。", "ちこくしたらあかんで。", "You mustn't be late.", "Blaze"),
		jv("おおきに", "ōkini", "thank you (Kansai)", "まいど、おおきに！", "maido, ōkini!", "Thanks, as always!", "Cora"),
		jv("〜へん", "〜hen", "negative ending (Kansai)", "行かへん。", "いかへん。", "I'm not going.", "Blaze"),
		jv("めんそーれ", "mensōre", "welcome (Okinawa)", "沖縄へめんそーれ！", "おきなわへめんそーれ！", "Welcome to Okinawa!", "Zephyr"),
	)

	s = addSkillL(db, ja, u, "Pragmatics: Reading the Air", "本音と建前, 空気を読む, 内と外 and indirect refusals.", "Brain", "#FF5C5C", 50, 1090)
	l = addLesson(db, s, "What Isn't Said", 1, 28,
		char(finch, "Near-native Japanese is as much about what goes unsaid. 本音 (true feelings) vs 建前 (public face). 空気を読む: read the air — sense what's expected. Refusals are indirect: ちょっと… or 検討します often mean no. 内と外 (in-group / out-group) changes keigo: to a client you call your own boss by name with no さん and humble verbs."),
		mc("A client asks for your boss, Suzuki. You say:", "To an outsider", "鈴木はただいま外出しております。", "鈴木はただいま外出しております。", "鈴木部長はいらっしゃいません。", "鈴木さんは外出されています。", "鈴木様は出かけました。"),
		mc("What does 前向きに検討します often mean in business?", "前向きに検討します", "A polite way of saying probably not", "A polite way of saying probably not", "Definitely yes", "We'll decide today", "We don't understand"),
		mc("What does 空気を読む mean?", "空気を読む", "Sense the mood and act appropriately", "Sense the mood and act appropriately", "Check the weather", "Read aloud", "Breathe deeply"),
		mc("What is 建前?", "本音と建前", "The polite public face", "The polite public face", "True feelings", "A building's foundation", "A formal letter"),
		mc("Someone invites you and you reply ちょっと…. What do they understand?", "今晩、飲みに行きませんか。— 今晩はちょっと…", "You're declining", "You're declining", "You'll come a bit later", "You'd like a small drink", "You didn't hear"),
		mc("What does 遠慮 refer to?", "遠慮しないでください", "Holding back out of consideration", "Holding back out of consideration", "Travelling far away", "Being rude on purpose", "Speaking loudly"),
		speak("Mira", "せっかくですが、今回は遠慮させていただきます。"),
	)
	addVocab(db, l,
		jv("本音", "ほんね honne", "true feelings", "本音を言うと、行きたくない。", "ほんねをいうと、いきたくない。", "Honestly, I don't want to go.", "Mira"),
		jv("建前", "たてまえ tatemae", "public face, official stance", "それは建前だ。", "それはたてまえだ。", "That's just the official line.", finch),
		jv("空気を読む", "くうきをよむ kūki o yomu", "read the room", "空気を読んで黙った。", "くうきをよんでだまった。", "I read the room and stayed quiet.", "Zephyr"),
		jv("遠慮", "えんりょ enryo", "restraint, holding back", "ご遠慮ください。", "ごえんりょください。", "Please refrain.", "Nana"),
		jv("検討", "けんとう kentō", "consideration, review", "前向きに検討します。", "まえむきにけんとうします。", "We'll consider it positively.", finch),
	)

	s = addSkillL(db, ja, u, "Ceremonial Speech & Formal Letters", "Weddings, speeches, 拝啓/敬具 and seasonal greetings.", "Mic", "#17A3DD", 51, 1120)
	l = addLesson(db, s, "The Most Formal Japanese", 1, 30,
		char(finch, "Formal letters open with 拝啓 and close with 敬具, with a seasonal greeting (時候の挨拶) such as 陽春の候 (in this bright spring season). Speeches start with 僭越ではございますが (if I may be so bold) and end with ご清聴ありがとうございました. At weddings, avoid 忌み言葉 — words like 切れる, 別れる, 終わる — and repeated words like たびたび."),
		mc("How does a formal letter that starts with 拝啓 end?", "拝啓 … ", "敬具", "敬具", "草々", "以上", "かしこ only"),
		mc("Which word must you avoid in a wedding speech?", "忌み言葉", "別れる", "別れる", "幸せ", "始まる", "結ぶ"),
		mc("How do you end a speech?", "Final line", "ご清聴ありがとうございました。", "ご清聴ありがとうございました。", "以上、さようなら。", "お疲れ様でした。", "よろしくね。"),
		mc("What is 陽春の候?", "Seasonal greeting", "A spring greeting for letters", "A spring greeting for letters", "A winter farewell", "A summer festival", "A New Year's dish"),
		mc("What does 僭越ではございますが signal?", "Opening a speech", "Humbly taking the floor", "Humbly taking the floor", "Apologising for lateness", "Disagreeing strongly", "Ending the speech"),
		write("Write the opening of a formal wedding speech in Japanese (about 150 characters): introduce yourself humbly, congratulate the couple, and avoid 忌み言葉.",
			"ただいまご紹介にあずかりました、新郎の友人の山田と申します。僭越ではございますが、一言お祝いの言葉を述べさせていただきます。健一さん、美咲さん、ご結婚誠におめでとうございます。また、ご両家の皆様にも心よりお祝い申し上げます。お二人の末永いお幸せをお祈りしております。"),
	)
	addVocab(db, l,
		jv("拝啓", "はいけい haikei", "Dear Sir/Madam (formal letter opening)", "拝啓 陽春の候、", "はいけい ようしゅんのこう、", "Dear Sir, in this bright spring season,", finch),
		jv("敬具", "けいぐ keigu", "Yours faithfully (closing)", "まずはお礼まで。敬具", "まずはおれいまで。けいぐ", "With thanks. Yours faithfully.", finch),
		jv("僭越ながら", "せんえつながら sen'etsu nagara", "if I may be so bold", "僭越ながら、乾杯の音頭を取らせていただきます。", "せんえつながら、かんぱいのおんどをとらせていただきます。", "If I may, I'll lead the toast.", "Mira"),
		jv("ご清聴", "ごせいちょう goseichō", "(your) kind attention", "ご清聴ありがとうございました。", "ごせいちょうありがとうございました。", "Thank you for your kind attention.", "Lumora"),
		jv("忌み言葉", "いみことば imikotoba", "taboo words (at ceremonies)", "結婚式では忌み言葉を避ける。", "けっこんしきではいみことばをさける。", "Avoid taboo words at weddings.", "Nana"),
	)
}
