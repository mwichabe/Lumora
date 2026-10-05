package controllers

// japanesePapers is the hand-authored Japanese proficiency-exam bank, modelled
// on the JLPT and keyed by the app's CEFR levels:
//
//	A1 = N5 · A2 = N4 · B1 = N3 · B2 = N2 · C1 = N1 · C2 = beyond N1
//	FINAL = comprehensive mastery (N1+ across every register)
//
// Like the JLPT, listening and reading lead; the app adds writing and
// speaking, which the JLPT doesn't test. Writing is scored in characters, so
// MinWords here is a character count — the clients count each kana and kanji
// as one unit. Questions are in English through N3 and in Japanese from N2.
// Reading texts carry kana readings at N5–N4.
var japanesePapers = map[string]paperContent{

	// ---------------- A1 · N5 ----------------
	"A1": {
		Listening: PaperListening{
			Title: "けんさんのいちにち — Ken's Day",
			Lines: []PaperLine{
				{Character: "Cora", Text: "わたしはけんです。がくせいです。まいあさしちじにおきます。", Translation: "I'm Ken. I'm a student. I get up at seven every morning."},
				{Character: "Cora", Text: "はちじにでんしゃでがっこうへいきます。ひるごはんはがっこうでたべます。すしがすきです。", Translation: "At eight I go to school by train. I eat lunch at school. I like sushi."},
				{Character: "Cora", Text: "ごごごじにうちへかえります。よるはともだちとでんわではなします。", Translation: "I go home at five in the afternoon. In the evening I talk with friends on the phone."},
			},
			Questions: []PaperQuestion{
				{Question: "What time does Ken get up?", Options: []string{"7:00", "8:00", "5:00", "9:00"}, CorrectAnswer: "7:00"},
				{Question: "How does he go to school?", Options: []string{"By train", "By bus", "On foot", "By bicycle"}, CorrectAnswer: "By train"},
				{Question: "What food does he like?", Options: []string{"Sushi", "Ramen", "Bread", "Curry"}, CorrectAnswer: "Sushi"},
				{Question: "What does he do in the evening?", Options: []string{"Talks with friends on the phone", "Watches TV", "Studies at the library", "Plays football"}, CorrectAnswer: "Talks with friends on the phone"},
			},
		},
		Reading: PaperReading{
			Title: "わたしのともだち — My Friend",
			Paragraphs: []string{
				"わたしのともだちはミラさんです。ミラさんはケニアじんです。(Watashi no tomodachi wa Mira-san desu. Mira-san wa Keniajin desu.)",
				"ミラさんはおおさかのだいがくで日本語をべんきょうしています。(Mira-san wa Ōsaka no daigaku de nihongo o benkyō shite imasu.)",
				"どようびにふたりでこうえんへいきました。こうえんにはおおきい木がありました。(Doyōbi ni futari de kōen e ikimashita. Kōen ni wa ōkii ki ga arimashita.)",
			},
			Questions: []PaperQuestion{
				{Question: "Where is Mira from?", Options: []string{"Kenya", "Japan", "China", "France"}, CorrectAnswer: "Kenya"},
				{Question: "Where does she study?", Options: []string{"At a university in Osaka", "At a school in Tokyo", "At home", "In Kenya"}, CorrectAnswer: "At a university in Osaka"},
				{Question: "When did they go to the park?", Options: []string{"Saturday", "Sunday", "Friday", "Monday"}, CorrectAnswer: "Saturday"},
				{Question: "What was in the park?", Options: []string{"A big tree", "A small cat", "A river", "A shop"}, CorrectAnswer: "A big tree"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write a short self-introduction in Japanese (hiragana, katakana or kanji): your name, where you're from, what you do, and one thing you like.",
			MinWords: 40,
		},
		Speaking: PaperSpeaking{
			Phrase:      "はじめまして。わたしはアンナです。ケニアからきました。よろしくおねがいします。",
			Speaker:     "Lumora",
			Translation: "Hajimemashite. Watashi wa Anna desu. Kenia kara kimashita. Yoroshiku onegaishimasu. — Nice to meet you. I'm Anna. I'm from Kenya. Pleased to meet you.",
		},
	},

	// ---------------- A2 · N4 ----------------
	"A2": {
		Listening: PaperListening{
			Title: "京都旅行 — A Trip to Kyoto",
			Lines: []PaperLine{
				{Character: "Zephyr", Text: "先月、友達と京都へ旅行に行きました。東京から新幹線で二時間ぐらいかかりました。", Translation: "Last month I went on a trip to Kyoto with a friend. It took about two hours from Tokyo by shinkansen."},
				{Character: "Zephyr", Text: "金閣寺を見たり、抹茶を飲んだりしました。京都は東京より静かで、とてもよかったです。", Translation: "We saw Kinkaku-ji, drank matcha, and so on. Kyoto was quieter than Tokyo, and very nice."},
				{Character: "Zephyr", Text: "でも、雨が降ったので、嵐山へは行けませんでした。来年の秋、もう一度行きたいと思っています。", Translation: "But because it rained, we couldn't go to Arashiyama. I'm thinking of going again next autumn."},
			},
			Questions: []PaperQuestion{
				{Question: "How long did the trip from Tokyo take?", Options: []string{"About two hours", "About one hour", "About four hours", "About thirty minutes"}, CorrectAnswer: "About two hours"},
				{Question: "How did he find Kyoto compared with Tokyo?", Options: []string{"Quieter", "Busier", "More expensive", "Colder"}, CorrectAnswer: "Quieter"},
				{Question: "Why couldn't they go to Arashiyama?", Options: []string{"It rained", "It was closed", "They had no time", "They got lost"}, CorrectAnswer: "It rained"},
				{Question: "When does he want to go again?", Options: []string{"Next autumn", "Next month", "Next spring", "Never"}, CorrectAnswer: "Next autumn"},
			},
		},
		Reading: PaperReading{
			Title: "アルバイトのお知らせ — Part-Time Job Notice",
			Paragraphs: []string{
				"喫茶店でアルバイトをしてくれる人をさがしています。(きっさてんでアルバイトをしてくれるひとをさがしています。)",
				"時間は土曜日と日曜日の午前十時から午後四時までです。日本語が少し話せれば大丈夫です。(じかんはどようびとにちようびのごぜんじゅうじからごごよじまでです。にほんごがすこしはなせればだいじょうぶです。)",
				"興味がある人は、今月の二十日までにお店に電話してください。(きょうみがあるひとは、こんげつのはつかまでにおみせにでんわしてください。)",
			},
			Questions: []PaperQuestion{
				{Question: "What kind of place is hiring?", Options: []string{"A café", "A hotel", "A bookshop", "A school"}, CorrectAnswer: "A café"},
				{Question: "Which days is the job?", Options: []string{"Saturday and Sunday", "Monday to Friday", "Only Sunday", "Every day"}, CorrectAnswer: "Saturday and Sunday"},
				{Question: "How much Japanese do you need?", Options: []string{"A little is fine", "Fluent", "None", "Business level"}, CorrectAnswer: "A little is fine"},
				{Question: "How should you apply?", Options: []string{"Phone the shop by the 20th", "Email by the 20th", "Visit on the 2nd", "Send a letter"}, CorrectAnswer: "Phone the shop by the 20th"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write in Japanese about your last weekend: where you went, what you did, and how it was. Use at least one 〜たり〜たりしました or 〜から/〜ので.",
			MinWords: 70,
		},
		Speaking: PaperSpeaking{
			Phrase:      "週末は友達と映画を見たり、買い物をしたりしました。",
			Speaker:     "Blaze",
			Translation: "しゅうまつはともだちとえいがをみたり、かいものをしたりしました。 — At the weekend I watched films and went shopping with friends.",
		},
	},

	// ---------------- B1 · N3 ----------------
	"B1": {
		Listening: PaperListening{
			Title: "引っ越しの相談 — Talking About Moving",
			Lines: []PaperLine{
				{Character: "Mira", Text: "最近、会社の近くに引っ越そうかと思っているんだ。今の家からだと、通勤に一時間半もかかるから。", Translation: "Lately I've been thinking of moving near the office. From my current place the commute takes an hour and a half."},
				{Character: "Blaze", Text: "でも、会社の近くは家賃が高いんじゃない？", Translation: "But isn't rent expensive near the office?"},
				{Character: "Mira", Text: "うん、二万円ぐらい高くなるらしい。でも、通勤時間が減れば、朝ゆっくりできるし、運動する時間も作れると思うんだ。", Translation: "Yeah, apparently it'd be about 20,000 yen more. But if the commute gets shorter, I can take my time in the morning and make time to exercise."},
				{Character: "Blaze", Text: "なるほどね。じゃあ、まず週末にいくつか部屋を見に行ってみたら？", Translation: "I see. Then why not go and look at a few flats this weekend first?"},
			},
			Questions: []PaperQuestion{
				{Question: "Why is Mira thinking of moving?", Options: []string{"Her commute is too long", "Her flat is too small", "Her rent is too high", "She changed jobs"}, CorrectAnswer: "Her commute is too long"},
				{Question: "How long is her commute now?", Options: []string{"An hour and a half", "Thirty minutes", "Two hours", "One hour"}, CorrectAnswer: "An hour and a half"},
				{Question: "How much more would the rent be?", Options: []string{"About ¥20,000", "About ¥2,000", "About ¥200,000", "The same"}, CorrectAnswer: "About ¥20,000"},
				{Question: "What does Blaze suggest?", Options: []string{"Look at some flats this weekend", "Stay where she is", "Change jobs", "Buy a car"}, CorrectAnswer: "Look at some flats this weekend"},
			},
		},
		Reading: PaperReading{
			Title: "スマホと睡眠 — Smartphones and Sleep",
			Paragraphs: []string{
				"寝る前にスマートフォンを使う人が増えている。ある調査によると、二十代の約七割が、ベッドの中でスマホを見ているそうだ。",
				"画面から出る光は、眠りを助けるホルモンの働きを弱めるため、寝つきが悪くなると言われている。",
				"専門家は、寝る一時間前にはスマホを使わないようにすることや、部屋を暗くすることをすすめている。",
			},
			Questions: []PaperQuestion{
				{Question: "According to the survey, how many people in their twenties use phones in bed?", Options: []string{"About 70%", "About 7%", "About 17%", "Almost 100%"}, CorrectAnswer: "About 70%"},
				{Question: "Why does screen light affect sleep?", Options: []string{"It weakens a hormone that helps sleep", "It makes the room too hot", "It hurts the eyes permanently", "It causes headaches"}, CorrectAnswer: "It weakens a hormone that helps sleep"},
				{Question: "What do experts recommend?", Options: []string{"Avoid phones for an hour before bed", "Use phones only in bed", "Sleep with the lights on", "Buy a new phone"}, CorrectAnswer: "Avoid phones for an hour before bed"},
				{Question: "What is the article's main topic?", Options: []string{"Smartphones and sleep", "Phone prices", "Hormones in food", "Bedroom design"}, CorrectAnswer: "Smartphones and sleep"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write in Japanese about a habit you started (or want to start) and how it has changed you. Use 〜ようになった or 〜ようにしている, and give a reason.",
			MinWords: 130,
		},
		Speaking: PaperSpeaking{
			Phrase:      "最近、毎朝三十分歩くようにしています。おかげで、よく眠れるようになりました。",
			Speaker:     "Mira",
			Translation: "Lately I make a point of walking for thirty minutes every morning. Thanks to that, I've started sleeping well.",
		},
	},

	// ---------------- B2 · N2 ----------------
	"B2": {
		Listening: PaperListening{
			Title: "会議：新商品の販売計画 — Meeting: Launch Plan for a New Product",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "それでは、来月発売予定の新商品について、販売計画を確認させていただきます。", Translation: "Now, allow me to go over the sales plan for the new product due out next month."},
				{Character: "Mira", Text: "当初はテレビ広告を中心に考えておりましたが、若い世代へのアピールを重視し、SNSでの宣伝に予算の半分を充てることにいたしました。", Translation: "We had initially planned to focus on TV advertising, but to prioritise appeal to younger people, we've decided to devote half the budget to social-media promotion."},
				{Character: "Professor Finch", Text: "なるほど。ただ、高齢のお客様も多い商品ですので、店頭での説明は引き続き必要ではないでしょうか。", Translation: "I see. However, as many of our customers for this product are elderly, isn't in-store explanation still necessary?"},
				{Character: "Mira", Text: "おっしゃるとおりです。主要店舗には説明スタッフを配置する予定です。", Translation: "Exactly as you say. We plan to station explanation staff at the main stores."},
			},
			Questions: []PaperQuestion{
				{Question: "当初はどんな宣伝を中心に考えていましたか。", Options: []string{"テレビ広告", "SNSでの宣伝", "店頭での説明", "新聞広告"}, CorrectAnswer: "テレビ広告"},
				{Question: "予算の半分は何に使われますか。", Options: []string{"SNSでの宣伝", "テレビ広告", "店員の給料", "商品開発"}, CorrectAnswer: "SNSでの宣伝"},
				{Question: "なぜ店頭での説明が必要だと言っていますか。", Options: []string{"高齢の客が多いから", "若い客が多いから", "商品が高いから", "店が少ないから"}, CorrectAnswer: "高齢の客が多いから"},
				{Question: "最終的にどうすることになりましたか。", Options: []string{"主要店舗に説明スタッフを置く", "SNSの宣伝をやめる", "発売を延期する", "テレビ広告だけにする"}, CorrectAnswer: "主要店舗に説明スタッフを置く"},
			},
		},
		Reading: PaperReading{
			Title: "外国人観光客の増加と課題 — Inbound Tourism and Its Challenges",
			Paragraphs: []string{
				"近年、日本を訪れる外国人観光客は増加の一途をたどっている。観光業は地域経済にとって重要な収入源となっており、各地で受け入れ体制の整備が進められている。",
				"その一方で、一部の人気観光地では、混雑やごみ問題、住民の生活への影響といった「観光公害」が深刻化している。",
				"観光客を制限せざるを得ないという意見もあるが、それは経済的な損失につながりかねない。重要なのは、観光客を地方へ分散させるなど、地域と観光が共存できる仕組みを作ることだろう。",
			},
			Questions: []PaperQuestion{
				{Question: "観光業は地域にとってどのようなものだと述べられていますか。", Options: []string{"重要な収入源", "大きな負担", "不必要なもの", "新しい公害"}, CorrectAnswer: "重要な収入源"},
				{Question: "「観光公害」の例として挙げられていないものはどれですか。", Options: []string{"物価の下落", "混雑", "ごみ問題", "住民生活への影響"}, CorrectAnswer: "物価の下落"},
				{Question: "観光客を制限することについて、筆者はどう考えていますか。", Options: []string{"経済的な損失につながるおそれがある", "すぐに実施すべきだ", "まったく問題ない", "観光公害を完全に解決する"}, CorrectAnswer: "経済的な損失につながるおそれがある"},
				{Question: "筆者が最も重要だと考えていることは何ですか。", Options: []string{"地域と観光が共存できる仕組みを作ること", "観光客をすべて受け入れること", "観光業をやめること", "都市部に観光客を集めること"}, CorrectAnswer: "地域と観光が共存できる仕組みを作ること"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write a formal email in Japanese to a client apologising for a delayed delivery, explaining the reason, and proposing a new delivery date. Use appropriate business keigo.",
			MinWords: 220,
		},
		Speaking: PaperSpeaking{
			Phrase:      "お忙しいところ恐れ入りますが、資料をご確認いただけますでしょうか。",
			Speaker:     "Mira",
			Translation: "I'm sorry to trouble you when you're busy, but could you kindly check the materials?",
		},
	},

	// ---------------- C1 · N1 ----------------
	"C1": {
		Listening: PaperListening{
			Title: "講演：人工知能と翻訳の未来 — Talk: AI and the Future of Translation",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "機械翻訳の精度は、ここ数年で飛躍的に向上したと言わざるを得ません。日常的な文書であれば、もはや人間の翻訳と見分けがつかないことも少なくありません。", Translation: "One must admit that machine translation accuracy has improved dramatically in recent years. For everyday documents, it's often indistinguishable from human translation."},
				{Character: "Professor Finch", Text: "しかしながら、文学作品や外交文書のように、文脈や文化的背景に即した微妙なニュアンスが求められる分野では、依然として人間の判断が不可欠です。", Translation: "However, in fields such as literature or diplomatic documents, where subtle nuance in line with context and cultural background is required, human judgment remains indispensable."},
				{Character: "Professor Finch", Text: "今後、翻訳者の役割は、ゼロから訳すことから、機械の訳を検証し、磨き上げることへと移っていくに違いありません。", Translation: "From now on, the translator's role will surely shift from translating from scratch to verifying and polishing machine output."},
			},
			Questions: []PaperQuestion{
				{Question: "話し手は機械翻訳の現状をどう評価していますか。", Options: []string{"大きく進歩したと認めている", "まったく役に立たないと考えている", "人間より常に優れていると考えている", "評価を避けている"}, CorrectAnswer: "大きく進歩したと認めている"},
				{Question: "人間の判断が依然として不可欠な分野はどれですか。", Options: []string{"文学作品や外交文書", "日常的な文書", "天気予報", "料理のレシピ"}, CorrectAnswer: "文学作品や外交文書"},
				{Question: "その理由は何ですか。", Options: []string{"文脈や文化に即した微妙なニュアンスが必要だから", "機械が高価だから", "文書が短いから", "法律で禁止されているから"}, CorrectAnswer: "文脈や文化に即した微妙なニュアンスが必要だから"},
				{Question: "今後、翻訳者の役割はどう変わると述べていますか。", Options: []string{"機械の訳を検証し磨き上げる方向へ移る", "完全になくなる", "今と全く変わらない", "機械を作る仕事になる"}, CorrectAnswer: "機械の訳を検証し磨き上げる方向へ移る"},
			},
		},
		Reading: PaperReading{
			Title: "「便利さ」の代償 — The Price of Convenience",
			Paragraphs: []string{
				"現代社会は、あらゆるものを「便利」にすべく進歩を続けてきた。二十四時間営業の店、即日配達、ワンクリックの決済——私たちは待つことを忘れつつある。",
				"しかし、便利さを追求するあまり、私たちは何かを失ってはいないだろうか。待つ時間とは、本来、期待を育み、物事の価値を噛みしめる時間でもあったはずだ。",
				"利便性を否定するつもりは毛頭ない。ただ、その恩恵に浴する一方で、不便さの中にこそ宿る豊かさにも目を向けるべきではないだろうか。",
			},
			Questions: []PaperQuestion{
				{Question: "筆者によると、現代人は何を忘れつつありますか。", Options: []string{"待つこと", "働くこと", "買い物をすること", "便利さ"}, CorrectAnswer: "待つこと"},
				{Question: "筆者は「待つ時間」をどのようなものだと考えていますか。", Options: []string{"期待を育み、価値を噛みしめる時間", "無駄な時間", "退屈で避けるべき時間", "仕事のための時間"}, CorrectAnswer: "期待を育み、価値を噛みしめる時間"},
				{Question: "筆者は利便性についてどう考えていますか。", Options: []string{"否定するつもりはまったくない", "すべて否定すべきだ", "もっと追求すべきだ", "関心がない"}, CorrectAnswer: "否定するつもりはまったくない"},
				{Question: "筆者の主張として最も適切なものはどれですか。", Options: []string{"不便さの中にある豊かさにも目を向けるべきだ", "便利さはすべて悪である", "即日配達を禁止すべきだ", "店は二十四時間営業すべきだ"}, CorrectAnswer: "不便さの中にある豊かさにも目を向けるべきだ"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write an argumentative essay in Japanese (である style): 「AIの発達は人間の仕事を奪うのか」. Present arguments for and against, address a counterargument, and reach a reasoned conclusion.",
			MinWords: 320,
		},
		Speaking: PaperSpeaking{
			Phrase:      "以上のことから、技術の進歩は脅威であると同時に、新たな可能性をもたらすものだと考えられます。",
			Speaker:     "Professor Finch",
			Translation: "From the above, technological progress can be considered both a threat and something that brings new possibilities.",
		},
	},

	// ---------------- C2 · Beyond N1 ----------------
	"C2": {
		Listening: PaperListening{
			Title: "祝辞 — A Wedding Speech",
			Lines: []PaperLine{
				{Character: "Zephyr", Text: "ただいまご紹介にあずかりました、新婦の上司の中村と申します。僭越ではございますが、一言ご挨拶を申し上げます。", Translation: "I am Nakamura, the bride's manager, just introduced. If I may be so bold, I'd like to say a few words."},
				{Character: "Zephyr", Text: "美咲さんは、入社以来、どんな困難にも笑顔で向き合い、周囲を明るく照らしてくれる存在でございました。", Translation: "Ever since joining the company, Misaki has faced every difficulty with a smile and has been a presence who brightens everyone around her."},
				{Character: "Zephyr", Text: "お二人には、互いを思いやる心を大切に、温かい家庭を築いていかれますよう、心よりお祈り申し上げます。本日は誠におめでとうございます。", Translation: "I sincerely hope that the two of you will cherish your consideration for each other and build a warm home. Heartfelt congratulations today."},
			},
			Questions: []PaperQuestion{
				{Question: "話し手と新婦の関係は何ですか。", Options: []string{"会社の上司", "大学の友人", "新郎の父", "幼なじみ"}, CorrectAnswer: "会社の上司"},
				{Question: "「僭越ではございますが」はどんな役割を果たしていますか。", Options: []string{"へりくだって話し始める", "相手を批判する", "話を終える", "冗談を言う"}, CorrectAnswer: "へりくだって話し始める"},
				{Question: "新婦はどのような人物として描かれていますか。", Options: []string{"困難にも笑顔で向き合う明るい人", "厳しく無口な人", "仕事を休みがちな人", "目立たない人"}, CorrectAnswer: "困難にも笑顔で向き合う明るい人"},
				{Question: "この祝辞で避けられている「忌み言葉」の例はどれですか。", Options: []string{"別れる", "築く", "祈る", "照らす"}, CorrectAnswer: "別れる"},
			},
		},
		Reading: PaperReading{
			Title: "徒然草 第七段より — From Tsurezuregusa, Section 7",
			Paragraphs: []string{
				"あだし野の露消ゆる時なく、鳥部山の煙立ち去らでのみ住み果つる習ひならば、いかにもののあはれもなからん。世は定めなきこそいみじけれ。",
				"（現代語訳）もしあだし野の露が消えることもなく、鳥部山の煙が立ち去ることもなく、人がいつまでも生き続けるのが世の習いであったなら、どんなにもののあわれもないことだろう。この世は無常であるからこそ素晴らしいのだ。",
				"兼好法師はこのように、命のはかなさを嘆くのではなく、むしろ無常であるがゆえの美しさを肯定している。",
			},
			Questions: []PaperQuestion{
				{Question: "「世は定めなきこそいみじけれ」の意味として最も適切なものはどれですか。", Options: []string{"この世は無常だからこそ素晴らしい", "この世は定めがあるから安心だ", "この世はいつまでも変わらない", "この世はつまらない"}, CorrectAnswer: "この世は無常だからこそ素晴らしい"},
				{Question: "露や煙は何の象徴として使われていますか。", Options: []string{"命のはかなさ", "豊かな自然", "都の繁栄", "季節の楽しみ"}, CorrectAnswer: "命のはかなさ"},
				{Question: "兼好法師の態度として正しいものはどれですか。", Options: []string{"無常を肯定的にとらえている", "無常を激しく嘆いている", "永遠の命を求めている", "自然を批判している"}, CorrectAnswer: "無常を肯定的にとらえている"},
				{Question: "「なからん」の「ん（む）」はどのような意味を表しますか。", Options: []string{"推量（〜だろう）", "過去", "打消", "断定"}, CorrectAnswer: "推量（〜だろう）"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write a formal speech in Japanese for a company's 50th-anniversary ceremony: humble opening, appropriate seasonal or ceremonial expressions, thanks to stakeholders, and a forward-looking close (ご清聴ありがとうございました).",
			MinWords: 420,
		},
		Speaking: PaperSpeaking{
			Phrase:      "本日はご多忙のところ、かくも多数ご臨席を賜り、誠にありがとうございます。",
			Speaker:     "Zephyr",
			Translation: "Thank you most sincerely for honouring us with your presence in such numbers today, busy as you are.",
		},
	},

	// ---------------- FINAL · comprehensive ----------------
	"FINAL": {
		Listening: PaperListening{
			Title: "シンポジウム：言語と文化の継承 — Symposium: Passing On Language and Culture",
			Lines: []PaperLine{
				{Character: "Professor Finch", Text: "方言の衰退は、単なる言葉の消失にとどまらず、その土地で培われてきた世界の見方そのものの喪失を意味します。", Translation: "The decline of dialects means not merely the loss of words, but the loss of the very way of seeing the world cultivated in that region."},
				{Character: "Mira", Text: "おっしゃることはもっともですが、標準語の普及がもたらした教育や経済の機会を軽視するわけにもいかないのではないでしょうか。", Translation: "What you say is entirely reasonable, but surely we can't make light of the educational and economic opportunities that the spread of standard Japanese brought?"},
				{Character: "Professor Finch", Text: "ええ、二者択一ではないのです。標準語を共通の基盤としつつ、方言を記録し、次の世代が誇りを持って使える環境を整えることこそが肝要でしょう。", Translation: "Yes — it isn't either/or. What matters is keeping standard Japanese as a common foundation while recording dialects and creating an environment where the next generation can use them with pride."},
			},
			Questions: []PaperQuestion{
				{Question: "最初の話し手によれば、方言の衰退は何を意味しますか。", Options: []string{"土地で培われた世界の見方の喪失", "言葉が少し変わるだけのこと", "経済の発展", "教育の向上"}, CorrectAnswer: "土地で培われた世界の見方の喪失"},
				{Question: "二人目の話し手はどのような点を指摘していますか。", Options: []string{"標準語がもたらした機会も軽視できない", "方言はすぐに禁止すべきだ", "標準語は不要だ", "方言の記録は無駄だ"}, CorrectAnswer: "標準語がもたらした機会も軽視できない"},
				{Question: "「二者択一ではない」とはどういう意味ですか。", Options: []string{"どちらか一方を選ぶ問題ではない", "二つとも捨てるべきだ", "二つのうち方言を選ぶべきだ", "議論する価値がない"}, CorrectAnswer: "どちらか一方を選ぶ問題ではない"},
				{Question: "最終的に「肝要」だとされているのは何ですか。", Options: []string{"標準語を基盤としつつ方言を継承できる環境を整えること", "方言だけで教育を行うこと", "標準語の使用をやめること", "方言の研究を中止すること"}, CorrectAnswer: "標準語を基盤としつつ方言を継承できる環境を整えること"},
			},
		},
		Reading: PaperReading{
			Title: "翻訳という営み — On the Act of Translation",
			Paragraphs: []string{
				"翻訳とは、ある言語の意味を別の言語に移し替える作業であると一般には理解されている。しかし、言葉が文化と不可分である以上、完全な等価など望むべくもない。",
				"たとえば「よろしくお願いします」という表現は、文脈いかんによって依頼にも挨拶にもなり、英語に一語で置き換えることはできない。訳者は常に、何を残し何を捨てるかという選択を迫られているのである。",
				"してみれば、翻訳とは喪失の技術であると同時に、二つの言語のあいだに新たな意味を生み出す創造的な営みにほかならない。",
			},
			Questions: []PaperQuestion{
				{Question: "筆者によれば、なぜ完全な等価は望めないのですか。", Options: []string{"言葉が文化と切り離せないから", "辞書が不十分だから", "訳者の能力が足りないから", "言語の数が多すぎるから"}, CorrectAnswer: "言葉が文化と切り離せないから"},
				{Question: "「よろしくお願いします」の例は何を示すために使われていますか。", Options: []string{"一語で訳せない、文脈に依存する表現があること", "日本語が簡単であること", "英語の方が正確であること", "挨拶は不要であること"}, CorrectAnswer: "一語で訳せない、文脈に依存する表現があること"},
				{Question: "「望むべくもない」の意味として正しいものはどれですか。", Options: []string{"望むことなどできない", "強く望むべきだ", "すでに実現している", "望む人が多い"}, CorrectAnswer: "望むことなどできない"},
				{Question: "筆者の結論として最も適切なものはどれですか。", Options: []string{"翻訳は喪失であると同時に創造的な営みである", "翻訳は機械に任せるべきだ", "翻訳は不可能なので無意味だ", "翻訳は単なる置き換え作業だ"}, CorrectAnswer: "翻訳は喪失であると同時に創造的な営みである"},
			},
		},
		Writing: PaperWriting{
			Prompt:   "write an essay in Japanese (である style) on 「グローバル化の時代に、母語と方言をどう守るか」. Develop a nuanced argument across registers, engage a counter-position, cite at least one concrete example, and draw a reasoned conclusion.",
			MinWords: 500,
		},
		Speaking: PaperSpeaking{
			Phrase:      "言葉を学ぶということは、その言葉を育んできた人々の心に触れることにほかなりません。",
			Speaker:     "Lumora",
			Translation: "To learn a language is nothing other than to touch the hearts of the people who nurtured it.",
		},
	},
}
