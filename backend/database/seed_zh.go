package database

import (
	"gorm.io/gorm"

	"lumora/backend/models"
)

// ===== Mandarin Chinese course (HSK 1 → HSK 6, mapped onto CEFR A1 → C2) =====
//
// Structured on the HSK 3.0 framework (Chinese Proficiency Test). The app's
// level system is CEFR-based, so each HSK level is one CEFR-tagged unit, using
// the usual approximation:
//
//	HSK 1 = A1 · HSK 2 = A2 · HSK 3 = B1 · HSK 4 = B2 · HSK 5 = C1 · HSK 6 = C2
//
// (HSK 7–9, the shared advanced band, is covered by the Final exam.)
//
// HSK 1 opens with the foundations every learner needs first — Pinyin, the four
// tones plus the neutral tone, tone sandhi, strokes & radicals and measure
// words — before survival topics. Characters are Simplified (Mainland
// standard). Pinyin is shown with every word and example throughout the early
// levels and faded out of the questions from HSK 4.
//
// Every exercise carries its own answer options (mc / listen / match), so none
// rely on the server's generic distractor fallback, which is Spanish.

// zv builds a Chinese vocabulary item: the characters are what's spoken aloud;
// the Pinyin rides along in the translation so learners always see it.
func zv(hanzi, pinyin, english, example, examplePinyin, exampleEnglish, speaker string) models.VocabItem {
	return vw(hanzi, pinyin+" · "+english, example, examplePinyin+" — "+exampleEnglish, speaker)
}

func seedMandarin(db *gorm.DB) {
	seedMandarinHSK1(db)
	seedMandarinHSK2(db)
	seedMandarinHSK3(db)
	seedMandarinHSK4(db)
	seedMandarinHSK5(db)
	seedMandarinHSK6(db)
	seedMandarinListening(db)
	seedMandarinReading(db)
}

// ───────────────────────── HSK 1 · A1 — Survival Chinese ─────────────────────────

func seedMandarinHSK1(db *gorm.DB) {
	const zh = "zh"
	const u = "A1 · HSK 1 — Survival Chinese"
	finch := "Professor Finch"

	// ── Foundations: Pinyin & the four tones ──
	s := addSkillL(db, zh, u, "Pinyin & the Four Tones", "The sound system: initials, finals and tones.", "Music", "#6C3FC5", 1, 0)
	l := addLesson(db, s, "Four Tones + Neutral", 1, 15,
		char(finch, "Mandarin is tonal: the same syllable means different things with a different pitch. 1st tone is high and level (mā), 2nd rises (má), 3rd dips then rises (mǎ), 4th falls sharply (mà). A light, short syllable has the neutral tone (ma)."),
		listen("Listen and choose the syllable you hear", "妈", "mā — 1st tone (mother)", "mā — 1st tone (mother)", "má — 2nd tone (hemp)", "mǎ — 3rd tone (horse)", "mà — 4th tone (scold)"),
		listen("Listen and choose the syllable you hear", "马", "mǎ — 3rd tone (horse)", "mā — 1st tone (mother)", "má — 2nd tone (hemp)", "mǎ — 3rd tone (horse)", "mà — 4th tone (scold)"),
		mc("Which tone is high and level?", "Tone mark: ā", "1st tone", "1st tone", "2nd tone", "3rd tone", "4th tone"),
		mc("Which tone falls sharply?", "Tone mark: à", "4th tone", "1st tone", "2nd tone", "3rd tone", "4th tone"),
		mc("What does 吗 (ma) do at the end of a sentence?", "你好吗？", "Turns it into a yes/no question", "Turns it into a yes/no question", "Makes it past tense", "Makes it plural", "Makes it negative"),
		speak("Mira", "妈妈骂马吗？"),
	)
	addVocab(db, l,
		zv("妈", "mā", "mother (1st tone)", "妈妈好。", "Māma hǎo.", "Mum is well.", "Mira"),
		zv("麻", "má", "hemp (2nd tone)", "这是麻。", "Zhè shì má.", "This is hemp.", "Mira"),
		zv("马", "mǎ", "horse (3rd tone)", "这是马。", "Zhè shì mǎ.", "This is a horse.", "Mira"),
		zv("骂", "mà", "to scold (4th tone)", "别骂我。", "Bié mà wǒ.", "Don't scold me.", "Mira"),
		zv("吗", "ma", "question particle (neutral tone)", "你好吗？", "Nǐ hǎo ma?", "How are you?", "Lumora"),
	)
	l = addLesson(db, s, "Initials & Finals", 2, 15,
		char(finch, "A syllable is an initial (consonant) plus a final (vowel part). Watch out: q sounds like 'ch' in cheese, x like 'sh' in sheep, zh like 'j' in jam with the tongue curled, and c like 'ts' in cats."),
		mc("How is the Pinyin 'q' pronounced?", "qī (seven)", "Like 'ch' in cheese", "Like 'ch' in cheese", "Like 'k' in key", "Like 'q' in queen", "Like 'g' in go"),
		mc("How is the Pinyin 'x' pronounced?", "xiè (thanks)", "Like 'sh' in sheep", "Like 'sh' in sheep", "Like 'x' in box", "Like 'z' in zoo", "Like 'ks'"),
		mc("How is the Pinyin 'c' pronounced?", "cài (dish)", "Like 'ts' in cats", "Like 'ts' in cats", "Like 'k' in cat", "Like 's' in sun", "Like 'ch' in chair"),
		listen("Listen and choose what you hear", "七", "qī — seven", "qī — seven", "xī — west", "jī — chicken", "chī — eat"),
		listen("Listen and choose what you hear", "吃", "chī — eat", "qī — seven", "xī — west", "chī — eat", "shí — ten"),
		speak("Blaze", "七，吃，西，十"),
	)
	addVocab(db, l,
		zv("七", "qī", "seven", "我有七本书。", "Wǒ yǒu qī běn shū.", "I have seven books.", finch),
		zv("吃", "chī", "to eat", "我吃米饭。", "Wǒ chī mǐfàn.", "I eat rice.", "Cora"),
		zv("西", "xī", "west", "我在西边。", "Wǒ zài xībian.", "I'm on the west side.", finch),
		zv("十", "shí", "ten", "我十岁。", "Wǒ shí suì.", "I'm ten years old.", "Pip"),
	)

	// ── Foundations: tone changes ──
	s = addSkillL(db, zh, u, "Tone Changes", "Tone sandhi: 3rd tone, 不 bù and 一 yī.", "Waves", "#17A3DD", 2, 10)
	l = addLesson(db, s, "Tone Sandhi", 1, 15,
		char(finch, "Tones shift in speech. Two 3rd tones in a row: the first becomes 2nd — 你好 nǐ hǎo is said ní hǎo. 不 bù becomes bú before a 4th tone (不是 bú shì). 一 yī becomes yí before a 4th tone (一个 yí gè) and yì before 1st, 2nd and 3rd tones (一天 yì tiān)."),
		mc("How is 你好 (nǐ hǎo) actually pronounced?", "3rd + 3rd tone", "ní hǎo", "ní hǎo", "nǐ hǎo", "nì hào", "nī hāo"),
		mc("How is 不是 (bù shì) pronounced?", "不 before a 4th tone", "bú shì", "bú shì", "bù shì", "bǔ shì", "bū shì"),
		mc("How is 一个 (yī gè) pronounced?", "一 before a 4th tone", "yí gè", "yí gè", "yì gè", "yī gè", "yǐ gè"),
		mc("How is 一天 (yī tiān) pronounced?", "一 before a 1st tone", "yì tiān", "yì tiān", "yí tiān", "yī tiān", "yǐ tiān"),
		mc("How is 很好 (hěn hǎo) pronounced?", "3rd + 3rd tone", "hén hǎo", "hén hǎo", "hěn hǎo", "hèn hǎo", "hēn hǎo"),
		speak("Mira", "你好！不是。一个。一天。"),
	)
	addVocab(db, l,
		zv("你好", "nǐ hǎo (said ní hǎo)", "hello", "你好，老师！", "Nǐ hǎo, lǎoshī!", "Hello, teacher!", "Lumora"),
		zv("不是", "bú shì", "is not", "我不是老师。", "Wǒ bú shì lǎoshī.", "I'm not a teacher.", finch),
		zv("一个", "yí gè", "one (of something)", "一个人。", "Yí gè rén.", "One person.", finch),
		zv("很好", "hěn hǎo (said hén hǎo)", "very good", "我很好。", "Wǒ hěn hǎo.", "I'm very well.", "Cora"),
	)

	// ── Foundations: characters, strokes & radicals ──
	s = addSkillL(db, zh, u, "Strokes & Radicals", "How characters are built and written.", "PenLine", "#F5A623", 3, 20)
	l = addLesson(db, s, "Building Characters", 1, 15,
		char(finch, "Characters are built from strokes written in a fixed order: top before bottom, left before right, horizontal before vertical. Many contain a radical that hints at meaning — 氵 (water) in 河 river and 喝 is different: 口 (mouth) means drinking, speaking, eating."),
		mc("What does the radical 氵 usually indicate?", "河 hé (river), 海 hǎi (sea)", "Water", "Water", "Fire", "Person", "Tree"),
		mc("What does the radical 口 usually indicate?", "吃 chī (eat), 喝 hē (drink), 叫 jiào (call)", "Mouth", "Mouth", "Hand", "Heart", "Sun"),
		mc("What does the radical 亻 indicate?", "他 tā (he), 你 nǐ (you), 们 men", "Person", "Person", "Water", "Woman", "Speech"),
		mc("What does the radical 女 indicate?", "妈 mā (mother), 她 tā (she), 好 hǎo", "Woman", "Woman", "Child", "Mouth", "Earth"),
		mc("General stroke-order rule", "Which comes first?", "Top before bottom", "Top before bottom", "Bottom before top", "Right before left", "Inside before outside"),
		mc("What does 木 mean (and as a radical)?", "林 lín (grove), 树 shù (tree)", "Tree / wood", "Tree / wood", "Fire", "Water", "Metal"),
	)
	addVocab(db, l,
		zv("人", "rén", "person", "他是好人。", "Tā shì hǎo rén.", "He's a good person.", finch),
		zv("口", "kǒu", "mouth", "一口水。", "Yì kǒu shuǐ.", "A mouthful of water.", finch),
		zv("水", "shuǐ", "water", "我喝水。", "Wǒ hē shuǐ.", "I drink water.", "Cora"),
		zv("木", "mù", "wood, tree", "这是木头。", "Zhè shì mùtou.", "This is wood.", finch),
		zv("好", "hǎo", "good (女 woman + 子 child)", "很好！", "Hěn hǎo!", "Very good!", "Lumora"),
	)

	// ── Greetings & introductions ──
	s = addSkillL(db, zh, u, "Greetings & Introductions", "Say hello, give your name and nationality.", "Hand", "#00C2A8", 4, 30)
	l = addLesson(db, s, "Hello, I'm…", 1, 15,
		char("Lumora", "你好！Let's introduce ourselves. 我叫 (wǒ jiào) means 'my name is', and 我是…人 (wǒ shì … rén) says where you're from."),
		mc("What does this mean?", "你好 (nǐ hǎo)", "Hello", "Hello", "Goodbye", "Thank you", "Sorry"),
		mc("What does this mean?", "我叫小明。(Wǒ jiào Xiǎo Míng.)", "My name is Xiao Ming.", "My name is Xiao Ming.", "I am a student.", "I like Xiao Ming.", "Xiao Ming is here."),
		mc("Choose the Chinese", "I'm Kenyan.", "我是肯尼亚人。", "我是肯尼亚人。", "我叫肯尼亚。", "我去肯尼亚。", "我有肯尼亚人。"),
		mc("Complete the sentence", "你___什么名字？(What's your name?)", "叫", "叫", "是", "有", "在"),
		mc("What does this mean?", "很高兴认识你。(Hěn gāoxìng rènshi nǐ.)", "Nice to meet you.", "Nice to meet you.", "See you tomorrow.", "I'm very tired.", "I don't know you."),
		speak("Cora", "你好！我叫小明，我是肯尼亚人。"),
	)
	addVocab(db, l,
		zv("你好", "nǐ hǎo", "hello", "你好！", "Nǐ hǎo!", "Hello!", "Lumora"),
		zv("我", "wǒ", "I, me", "我是学生。", "Wǒ shì xuésheng.", "I am a student.", "Cora"),
		zv("叫", "jiào", "to be called", "我叫安娜。", "Wǒ jiào Ānnà.", "My name is Anna.", "Cora"),
		zv("是", "shì", "to be", "他是老师。", "Tā shì lǎoshī.", "He is a teacher.", finch),
		zv("中国人", "Zhōngguó rén", "Chinese person", "她是中国人。", "Tā shì Zhōngguó rén.", "She is Chinese.", "Mira"),
		zv("谢谢", "xièxie", "thank you", "谢谢你！", "Xièxie nǐ!", "Thank you!", "Lumora"),
	)
	l = addLesson(db, s, "Polite Phrases", 2, 15,
		char("Mira", "Politeness matters. 谢谢 thanks, 不客气 you're welcome, 对不起 sorry, 没关系 it doesn't matter, 再见 goodbye. 您 (nín) is the respectful 'you'."),
		mc("Reply to 谢谢！", "Thank you!", "不客气。(bú kèqi)", "不客气。(bú kèqi)", "对不起。(duìbuqǐ)", "再见。(zàijiàn)", "你好。(nǐ hǎo)"),
		mc("Reply to 对不起！", "Sorry!", "没关系。(méi guānxi)", "没关系。(méi guānxi)", "谢谢。(xièxie)", "不客气。(bú kèqi)", "你好。(nǐ hǎo)"),
		mc("Which 'you' is respectful?", "To an elder or a customer", "您 (nín)", "您 (nín)", "你 (nǐ)", "他 (tā)", "我 (wǒ)"),
		mc("What does this mean?", "再见！(Zàijiàn!)", "Goodbye!", "Goodbye!", "Hello!", "Thanks!", "Sorry!"),
		mc("What does this mean?", "老师好！(Lǎoshī hǎo!)", "Hello, teacher!", "Hello, teacher!", "The teacher is good.", "Goodbye, teacher!", "Thank you, teacher!"),
		speak("Mira", "谢谢！不客气。对不起！没关系。"),
	)
	addVocab(db, l,
		zv("不客气", "bú kèqi", "you're welcome", "谢谢！——不客气。", "Xièxie! — Bú kèqi.", "Thanks! — You're welcome.", "Mira"),
		zv("对不起", "duìbuqǐ", "sorry", "对不起，我来晚了。", "Duìbuqǐ, wǒ lái wǎn le.", "Sorry, I'm late.", "Mira"),
		zv("没关系", "méi guānxi", "it doesn't matter", "没关系！", "Méi guānxi!", "No problem!", "Mira"),
		zv("再见", "zàijiàn", "goodbye", "明天见！再见！", "Míngtiān jiàn! Zàijiàn!", "See you tomorrow! Bye!", "Lumora"),
		zv("您", "nín", "you (respectful)", "您好！", "Nín hǎo!", "Hello! (respectful)", finch),
	)

	// ── Numbers, age & measure words ──
	s = addSkillL(db, zh, u, "Numbers & Measure Words", "Count, tell your age, use 个 and 本.", "Hash", "#FF5C5C", 5, 42)
	l = addLesson(db, s, "Counting 1–100", 1, 15,
		char("Pip", "一 yī, 二 èr, 三 sān, 四 sì, 五 wǔ, 六 liù, 七 qī, 八 bā, 九 jiǔ, 十 shí. Then 十一 is 11, 二十 is 20, 九十九 is 99. Age uses 岁 suì: 我二十岁 — I'm 20. 8 is lucky (it sounds like 'prosper'); 4 is unlucky (sì sounds like 死, death)."),
		mc("What number is this?", "二十五 (èrshíwǔ)", "25", "25", "52", "205", "15"),
		mc("Choose the Chinese", "48", "四十八", "四十八", "八十四", "四八", "四百八"),
		mc("How old is she?", "她三十岁。(Tā sānshí suì.)", "30", "30", "3", "13", "33"),
		mc("Which number is considered lucky?", "It sounds like 发 (fā, 'to prosper')", "八 (8)", "八 (8)", "四 (4)", "七 (7)", "一 (1)"),
		mc("Complete the question", "你今年多大？(How old are you this year?) — 我二十___。", "岁", "岁", "个", "年", "号"),
		speak("Pip", "我今年二十八岁。"),
	)
	addVocab(db, l,
		zv("一", "yī", "one", "一个苹果。", "Yí gè píngguǒ.", "One apple.", "Pip"),
		zv("二", "èr", "two", "二十。", "Èrshí.", "Twenty.", "Pip"),
		zv("三", "sān", "three", "三个人。", "Sān gè rén.", "Three people.", "Pip"),
		zv("八", "bā", "eight", "八点。", "Bā diǎn.", "Eight o'clock.", "Pip"),
		zv("岁", "suì", "years old", "我二十岁。", "Wǒ èrshí suì.", "I'm twenty.", "Pip"),
	)
	l = addLesson(db, s, "Measure Words", 2, 15,
		char(finch, "Between a number and a noun, Chinese needs a measure word. 个 gè is the general one; 本 běn for books; 杯 bēi for cups; 口 kǒu for family members; 只 zhī for many animals. And 二 becomes 两 liǎng before a measure word: 两个人."),
		mc("Choose the measure word", "三___书 (three books)", "本", "本", "个", "杯", "只"),
		mc("Choose the measure word", "一___茶 (a cup of tea)", "杯", "杯", "本", "口", "只"),
		mc("Choose the measure word", "两___猫 (two cats)", "只", "只", "本", "杯", "口"),
		mc("Choose the correct phrase", "two people", "两个人", "两个人", "二个人", "两人个", "二人"),
		mc("What does this mean?", "我家有五口人。(Wǒ jiā yǒu wǔ kǒu rén.)", "There are five people in my family.", "There are five people in my family.", "I have five mouths.", "My home has five rooms.", "Five people came to my house."),
		speak("Blaze", "我要两杯茶和一本书。"),
	)
	addVocab(db, l,
		zv("个", "gè", "general measure word", "一个朋友。", "Yí gè péngyou.", "A friend.", finch),
		zv("本", "běn", "measure word for books", "两本书。", "Liǎng běn shū.", "Two books.", finch),
		zv("杯", "bēi", "cup (of)", "一杯水。", "Yì bēi shuǐ.", "A glass of water.", "Cora"),
		zv("两", "liǎng", "two (before measure words)", "两个人。", "Liǎng gè rén.", "Two people.", finch),
		zv("只", "zhī", "measure word for animals", "一只狗。", "Yì zhī gǒu.", "A dog.", "Pip"),
	)

	// ── Family & people ──
	s = addSkillL(db, zh, u, "Family & People", "Talk about your family with 有 and 的.", "Users", "#9B51E0", 6, 54)
	l = addLesson(db, s, "My Family", 1, 15,
		char("Nana", "有 yǒu means 'to have'; its negative is always 没有 méiyǒu (never 不有). 的 de shows possession: 我的妈妈 — my mother. With close family you can drop it: 我妈妈."),
		mc("What does this mean?", "我有一个哥哥。(Wǒ yǒu yí gè gēge.)", "I have an older brother.", "I have an older brother.", "I am an older brother.", "My brother has one.", "I have a younger sister."),
		mc("Make it negative", "我___妹妹。(I don't have a younger sister.)", "没有", "没有", "不有", "不是", "没是"),
		mc("What does this mean?", "这是我的爸爸。(Zhè shì wǒ de bàba.)", "This is my dad.", "This is my dad.", "My dad is here.", "I am a dad.", "Dad has this."),
		mc("Choose the Chinese", "older sister", "姐姐 (jiějie)", "姐姐 (jiějie)", "妹妹 (mèimei)", "哥哥 (gēge)", "弟弟 (dìdi)"),
		mc("Answer the question", "你家有几口人？(How many people are in your family?)", "我家有四口人。", "我家有四口人。", "我家是四个。", "我有四岁。", "我家在四号。"),
		speak("Nana", "我家有四口人：爸爸、妈妈、哥哥和我。"),
	)
	addVocab(db, l,
		zv("爸爸", "bàba", "dad", "我爸爸是医生。", "Wǒ bàba shì yīshēng.", "My dad is a doctor.", "Nana"),
		zv("妈妈", "māma", "mum", "妈妈在家。", "Māma zài jiā.", "Mum is at home.", "Nana"),
		zv("哥哥", "gēge", "older brother", "我有一个哥哥。", "Wǒ yǒu yí gè gēge.", "I have an older brother.", "Nana"),
		zv("姐姐", "jiějie", "older sister", "姐姐很漂亮。", "Jiějie hěn piàoliang.", "My sister is pretty.", "Nana"),
		zv("有", "yǒu", "to have; there is", "我有猫。", "Wǒ yǒu māo.", "I have a cat.", "Nana"),
		zv("没有", "méiyǒu", "not have", "我没有车。", "Wǒ méiyǒu chē.", "I don't have a car.", "Nana"),
	)

	// ── Food & ordering ──
	s = addSkillL(db, zh, u, "Food & Ordering", "Order food and drink, ask prices.", "Coffee", "#F5A623", 7, 68)
	l = addLesson(db, s, "At the Restaurant", 1, 15,
		char("Cora", "我要 (wǒ yào) — I want — is how you order. 多少钱？(duōshao qián?) asks the price. 吗 makes a yes/no question; 呢 asks 'and you?'. Try 饺子 jiǎozi (dumplings) and 茶 chá (tea)!"),
		mc("What does this mean?", "我要一碗面条。(Wǒ yào yì wǎn miàntiáo.)", "I'd like a bowl of noodles.", "I'd like a bowl of noodles.", "I eat noodles every day.", "The noodles are hot.", "I don't want noodles."),
		mc("How do you ask the price?", "How much is it?", "多少钱？(duōshao qián?)", "多少钱？(duōshao qián?)", "几点？(jǐ diǎn?)", "在哪儿？(zài nǎr?)", "是谁？(shì shéi?)"),
		mc("What does this mean?", "你喝茶吗？(Nǐ hē chá ma?)", "Do you drink tea?", "Do you drink tea?", "You drink tea.", "Where is the tea?", "What tea do you drink?"),
		mc("Choose the Chinese", "dumplings", "饺子 (jiǎozi)", "饺子 (jiǎozi)", "米饭 (mǐfàn)", "面条 (miàntiáo)", "包子 (bāozi)"),
		mc("Complete the dialogue", "我很好，你___？(I'm fine, and you?)", "呢", "呢", "吗", "的", "了"),
		mc("What does this mean?", "十五块。(Shíwǔ kuài.)", "15 yuan.", "15 yuan.", "50 yuan.", "15 bowls.", "5 yuan."),
		speak("Cora", "服务员，我要一盘饺子和一杯茶。多少钱？"),
	)
	addVocab(db, l,
		zv("要", "yào", "to want; would like", "我要米饭。", "Wǒ yào mǐfàn.", "I'd like rice.", "Cora"),
		zv("茶", "chá", "tea", "我喝绿茶。", "Wǒ hē lǜchá.", "I drink green tea.", "Cora"),
		zv("饺子", "jiǎozi", "dumplings", "我喜欢吃饺子。", "Wǒ xǐhuan chī jiǎozi.", "I like eating dumplings.", "Cora"),
		zv("多少钱", "duōshao qián", "how much (money)?", "这个多少钱？", "Zhège duōshao qián?", "How much is this?", "Cora"),
		zv("块", "kuài", "yuan (spoken)", "十块钱。", "Shí kuài qián.", "Ten yuan.", "Cora"),
		zv("喝", "hē", "to drink", "你喝什么？", "Nǐ hē shénme?", "What will you drink?", "Cora"),
	)

	// ── Time, days & where things are ──
	s = addSkillL(db, zh, u, "Time & Places", "Days, dates, clock time, 在 and question words.", "Clock", "#17A3DD", 8, 82)
	l = addLesson(db, s, "When and Where", 1, 15,
		char("Riko", "Time words come before the verb: 我明天去 — I'm going tomorrow. 在 zài says where something is: 书在桌子上 — the book is on the table. Question words stay where the answer goes: 你去哪儿？— where are you going?"),
		mc("Put the time in the right place", "I'll go to Beijing tomorrow.", "我明天去北京。", "我明天去北京。", "我去北京明天。", "明天北京去我。", "我去明天北京。"),
		mc("What does this mean?", "现在几点？(Xiànzài jǐ diǎn?)", "What time is it now?", "What time is it now?", "Where are you now?", "What day is it?", "How old are you?"),
		mc("What does this mean?", "书在桌子上。(Shū zài zhuōzi shang.)", "The book is on the table.", "The book is on the table.", "The book is under the table.", "The table is on the book.", "There's no book on the table."),
		mc("Choose the Chinese", "Wednesday", "星期三 (xīngqīsān)", "星期三 (xīngqīsān)", "星期一 (xīngqīyī)", "星期天 (xīngqītiān)", "三月 (sānyuè)"),
		mc("Choose the question word", "你去___？(Where are you going?)", "哪儿", "哪儿", "什么", "谁", "几"),
		speak("Riko", "现在八点。我九点去学校。"),
	)
	addVocab(db, l,
		zv("今天", "jīntiān", "today", "今天星期一。", "Jīntiān xīngqīyī.", "Today is Monday.", "Riko"),
		zv("明天", "míngtiān", "tomorrow", "明天见！", "Míngtiān jiàn!", "See you tomorrow!", "Riko"),
		zv("点", "diǎn", "o'clock", "三点。", "Sān diǎn.", "Three o'clock.", "Riko"),
		zv("在", "zài", "to be at; in", "我在家。", "Wǒ zài jiā.", "I'm at home.", "Riko"),
		zv("哪儿", "nǎr", "where", "你在哪儿？", "Nǐ zài nǎr?", "Where are you?", "Riko"),
		zv("星期", "xīngqī", "week", "这个星期很忙。", "Zhège xīngqī hěn máng.", "This week is busy.", "Riko"),
	)
}

// ───────────────────────── HSK 2 · A2 — Everyday Life ─────────────────────────

func seedMandarinHSK2(db *gorm.DB) {
	const zh = "zh"
	const u = "A2 · HSK 2 — Everyday Life"
	finch := "Professor Finch"

	s := addSkillL(db, zh, u, "Completed Actions with 了", "Talk about what happened and what changed.", "CheckCircle", "#6C3FC5", 9, 110)
	l := addLesson(db, s, "了 — Done & Changed", 1, 18,
		char(finch, "了 le after a verb marks a completed action: 我吃了饭 — I ate. 了 at the end of a sentence marks a change: 下雨了 — it's raining now. The negative of a completed action is 没 + verb, with no 了: 我没吃。"),
		mc("What does this mean?", "我买了三本书。(Wǒ mǎi le sān běn shū.)", "I bought three books.", "I bought three books.", "I will buy three books.", "I'm buying three books.", "I want three books."),
		mc("What does this mean?", "下雨了！(Xià yǔ le!)", "It's started raining!", "It's started raining!", "It rained yesterday.", "It will rain.", "It isn't raining."),
		mc("Make it negative", "I didn't go.", "我没去。", "我没去。", "我不去了。", "我没去了。", "我不去。"),
		mc("Choose the correct sentence", "He has arrived.", "他到了。", "他到了。", "他了到。", "他到过。", "他没到了。"),
		mc("What does this mean?", "我不想去了。(Wǒ bù xiǎng qù le.)", "I don't want to go any more.", "I don't want to go any more.", "I didn't want to go.", "I never go.", "I want to go now."),
		speak("Blaze", "昨天我去了商店，买了两件衣服。"),
	)
	addVocab(db, l,
		zv("了", "le", "completed action / change", "我吃了。", "Wǒ chī le.", "I've eaten.", finch),
		zv("昨天", "zuótiān", "yesterday", "昨天我很忙。", "Zuótiān wǒ hěn máng.", "I was busy yesterday.", "Cora"),
		zv("买", "mǎi", "to buy", "我买了苹果。", "Wǒ mǎi le píngguǒ.", "I bought apples.", "Cora"),
		zv("下雨", "xià yǔ", "to rain", "下雨了。", "Xià yǔ le.", "It's raining.", "Riko"),
		zv("没", "méi", "not (for past actions)", "我没看。", "Wǒ méi kàn.", "I didn't watch.", finch),
	)

	s = addSkillL(db, zh, u, "Experiences with 过", "Say what you have (never) done.", "Globe", "#00C2A8", 10, 124)
	l = addLesson(db, s, "Have You Ever…?", 1, 18,
		char("Zephyr", "过 guo after a verb means 'have done before' — experience. 我去过上海 — I've been to Shanghai. Negative: 没 + verb + 过: 我没吃过烤鸭 — I've never had roast duck. Ask with 吗 or 没有."),
		mc("What does this mean?", "我去过长城。(Wǒ qù guo Chángchéng.)", "I've been to the Great Wall.", "I've been to the Great Wall.", "I'm going to the Great Wall.", "I went past the Great Wall.", "I live by the Great Wall."),
		mc("Choose the negative", "I've never been to China.", "我没去过中国。", "我没去过中国。", "我不去过中国。", "我去过没中国。", "我没去中国了。"),
		mc("Choose the Chinese", "Have you ever eaten Peking duck?", "你吃过北京烤鸭吗？", "你吃过北京烤鸭吗？", "你吃了北京烤鸭吗？", "你过吃北京烤鸭吗？", "你吃北京烤鸭过吗？"),
		mc("Which particle marks a past experience?", "我去___北京。(I've been to Beijing.)", "过 (guo)", "过 (guo)", "了 (le)", "着 (zhe)", "在 (zài)"),
		mc("What does this mean?", "他学过汉语。(Tā xué guo Hànyǔ.)", "He has studied Chinese before.", "He has studied Chinese before.", "He is studying Chinese.", "He will study Chinese.", "He doesn't study Chinese."),
		speak("Zephyr", "我去过北京，可是没去过上海。"),
	)
	addVocab(db, l,
		zv("过", "guo", "(have) done before", "你来过吗？", "Nǐ lái guo ma?", "Have you been here before?", "Zephyr"),
		zv("长城", "Chángchéng", "the Great Wall", "长城很长。", "Chángchéng hěn cháng.", "The Great Wall is very long.", "Zephyr"),
		zv("北京", "Běijīng", "Beijing", "我住在北京。", "Wǒ zhù zài Běijīng.", "I live in Beijing.", "Zephyr"),
		zv("可是", "kěshì", "but", "很贵，可是很好。", "Hěn guì, kěshì hěn hǎo.", "Expensive, but very good.", "Zephyr"),
		zv("旅游", "lǚyóu", "to travel (for fun)", "我喜欢旅游。", "Wǒ xǐhuan lǚyóu.", "I like travelling.", "Zephyr"),
	)

	s = addSkillL(db, zh, u, "Comparisons with 比", "Compare people, prices and places.", "Scale", "#FF5C5C", 11, 138)
	l = addLesson(db, s, "Taller, Cheaper, Better", 1, 18,
		char(finch, "A 比 B + adjective: 他比我高 — he's taller than me. Never put 很 in a 比 sentence. Add 一点儿 (a bit) or 多了 (much) after the adjective. Negative: A 没有 B + adjective: 我没有他高."),
		mc("What does this mean?", "今天比昨天冷。(Jīntiān bǐ zuótiān lěng.)", "Today is colder than yesterday.", "Today is colder than yesterday.", "Yesterday was colder than today.", "Today is cold, yesterday too.", "Today isn't cold."),
		mc("Which sentence is correct?", "He is taller than me.", "他比我高。", "他比我高。", "他比我很高。", "他高比我。", "我比他高。"),
		mc("Choose the Chinese", "This one is a bit cheaper.", "这个便宜一点儿。", "这个便宜一点儿。", "这个很一点儿便宜。", "一点儿这个便宜。", "这个比便宜。"),
		mc("Make it negative", "I'm not as busy as you.", "我没有你忙。", "我没有你忙。", "我不比你忙。", "我比你不忙。", "我没比你忙。"),
		mc("What does this mean?", "北京比上海大多了。", "Beijing is much bigger than Shanghai.", "Beijing is much bigger than Shanghai.", "Shanghai is bigger than Beijing.", "Beijing and Shanghai are both big.", "Beijing is a bit bigger."),
		speak("Blaze", "坐地铁比坐出租车便宜多了。"),
	)
	addVocab(db, l,
		zv("比", "bǐ", "than (compared with)", "我比他大。", "Wǒ bǐ tā dà.", "I'm older than him.", finch),
		zv("便宜", "piányi", "cheap", "这个很便宜。", "Zhège hěn piányi.", "This is cheap.", "Cora"),
		zv("贵", "guì", "expensive", "太贵了！", "Tài guì le!", "Too expensive!", "Cora"),
		zv("一点儿", "yìdiǎnr", "a little", "快一点儿！", "Kuài yìdiǎnr!", "A bit faster!", "Blaze"),
		zv("冷", "lěng", "cold", "今天很冷。", "Jīntiān hěn lěng.", "It's cold today.", "Riko"),
	)

	s = addSkillL(db, zh, u, "Because…, so… & Although…, but…", "Link ideas with paired conjunctions.", "Link", "#17A3DD", 12, 152)
	l = addLesson(db, s, "Linking Ideas", 1, 18,
		char(finch, "Chinese often uses conjunctions in pairs: 因为…所以… (because… so…) and 虽然…但是… (although… but…). Unlike English, both halves are normally kept."),
		mc("Complete the sentence", "因为下雨，___我没去。", "所以", "所以", "但是", "虽然", "因为"),
		mc("Complete the sentence", "虽然很贵，___很好吃。", "但是", "但是", "所以", "因为", "还是"),
		mc("What does this mean?", "因为我病了，所以没上班。", "Because I was ill, I didn't go to work.", "Because I was ill, I didn't go to work.", "I was ill, but I went to work.", "Although I was ill, I went to work.", "I went to work so I got ill."),
		mc("What does this mean?", "虽然他很忙，但是每天都运动。", "Although he's busy, he exercises every day.", "Although he's busy, he exercises every day.", "Because he's busy, he doesn't exercise.", "He's busy so he exercises.", "He's not busy and exercises."),
		mc("Choose the Chinese", "Because it's far, I take a taxi.", "因为很远，所以我坐出租车。", "因为很远，所以我坐出租车。", "虽然很远，所以我坐出租车。", "因为很远，但是我坐出租车。", "所以很远，因为我坐出租车。"),
		speak("Mira", "虽然汉语很难，但是很有意思。"),
	)
	addVocab(db, l,
		zv("因为", "yīnwèi", "because", "因为我累了。", "Yīnwèi wǒ lèi le.", "Because I'm tired.", finch),
		zv("所以", "suǒyǐ", "so, therefore", "所以我回家了。", "Suǒyǐ wǒ huí jiā le.", "So I went home.", finch),
		zv("虽然", "suīrán", "although", "虽然很难…", "Suīrán hěn nán…", "Although it's hard…", finch),
		zv("但是", "dànshì", "but", "但是很有意思。", "Dànshì hěn yǒu yìsi.", "But it's interesting.", finch),
		zv("有意思", "yǒu yìsi", "interesting", "这本书很有意思。", "Zhè běn shū hěn yǒu yìsi.", "This book is interesting.", "Mira"),
	)

	s = addSkillL(db, zh, u, "Festivals & Habits", "Spring Festival, hobbies and daily routines.", "Sparkles", "#F5A623", 13, 166)
	l = addLesson(db, s, "Spring Festival", 1, 18,
		char("Nana", "春节 Chūnjié, the Spring Festival, is the most important holiday. Families gather for a reunion dinner, eat 饺子, give 红包 (red envelopes with money) and say 新年快乐 — Happy New Year!"),
		mc("What is 春节?", "Chūnjié", "Spring Festival (Chinese New Year)", "Spring Festival (Chinese New Year)", "Mid-Autumn Festival", "National Day", "Dragon Boat Festival"),
		mc("What is put inside a 红包?", "hóngbāo", "Money", "Money", "Tea", "Dumplings", "A letter"),
		mc("What do you say at New Year?", "Happy New Year!", "新年快乐！(Xīnnián kuàilè!)", "新年快乐！(Xīnnián kuàilè!)", "生日快乐！(Shēngrì kuàilè!)", "一路平安！(Yílù píng'ān!)", "祝你健康！(Zhù nǐ jiànkāng!)"),
		mc("What does this mean?", "我每天早上跑步。(Wǒ měitiān zǎoshang pǎobù.)", "I run every morning.", "I run every morning.", "I ran this morning.", "I'll run tomorrow morning.", "I don't like running."),
		mc("What does this mean?", "他一边吃饭一边看电视。", "He watches TV while eating.", "He watches TV while eating.", "He eats after watching TV.", "He doesn't watch TV.", "He eats in front of the TV shop."),
		speak("Nana", "春节的时候，我们一家人一起吃饺子。"),
	)
	addVocab(db, l,
		zv("春节", "Chūnjié", "Spring Festival", "春节快乐！", "Chūnjié kuàilè!", "Happy Spring Festival!", "Nana"),
		zv("红包", "hóngbāo", "red envelope (money gift)", "奶奶给我红包。", "Nǎinai gěi wǒ hóngbāo.", "Grandma gives me a red envelope.", "Nana"),
		zv("快乐", "kuàilè", "happy", "生日快乐！", "Shēngrì kuàilè!", "Happy birthday!", "Nana"),
		zv("每天", "měitiān", "every day", "我每天学汉语。", "Wǒ měitiān xué Hànyǔ.", "I study Chinese every day.", "Nana"),
		zv("一边…一边…", "yìbiān… yìbiān…", "while (doing two things)", "一边走一边说。", "Yìbiān zǒu yìbiān shuō.", "Talk while walking.", "Nana"),
	)
}

// ───────────────────────── HSK 3 · B1 — Social Chinese ─────────────────────────

func seedMandarinHSK3(db *gorm.DB) {
	const zh = "zh"
	const u = "B1 · HSK 3 — Social Chinese"
	finch := "Professor Finch"

	s := addSkillL(db, zh, u, "The 把 Construction", "Say what you do to something.", "Hand", "#6C3FC5", 14, 200)
	l := addLesson(db, s, "把 — Handling Things", 1, 20,
		char(finch, "把 bǎ moves the object before the verb to stress what happens to it: 我把作业做完了 — I finished the homework. The verb can't stand alone; it needs a result or complement after it (完, 好, 在…, 给…)."),
		mc("What does this mean?", "请把门关上。(Qǐng bǎ mén guān shang.)", "Please close the door.", "Please close the door.", "Please open the door.", "The door is closed.", "Please knock on the door."),
		mc("Choose the correct 把 sentence", "I put the book on the table.", "我把书放在桌子上。", "我把书放在桌子上。", "我把放书在桌子上。", "我放把书在桌子上。", "我书把放在桌子上。"),
		mc("Which is NOT a complete 把 sentence?", "A 把 verb needs something after it.", "我把饭吃。", "我把饭吃。", "我把饭吃完了。", "我把饭吃了。", "我把饭吃光了。"),
		mc("What does this mean?", "他把钱包丢了。(Tā bǎ qiánbāo diū le.)", "He lost his wallet.", "He lost his wallet.", "He found his wallet.", "He bought a wallet.", "He gave away his wallet."),
		mc("Choose the Chinese", "Please give this letter to the teacher.", "请把这封信交给老师。", "请把这封信交给老师。", "请交给老师把这封信。", "请这封信把交给老师。", "把请这封信交给老师。"),
		speak("Blaze", "我已经把报告写完了。"),
	)
	addVocab(db, l,
		zv("把", "bǎ", "(object marker)", "把书给我。", "Bǎ shū gěi wǒ.", "Give me the book.", finch),
		zv("放", "fàng", "to put", "放在这儿。", "Fàng zài zhèr.", "Put it here.", finch),
		zv("完", "wán", "to finish (result)", "我做完了。", "Wǒ zuò wán le.", "I've finished.", finch),
		zv("丢", "diū", "to lose", "我的手机丢了。", "Wǒ de shǒujī diū le.", "My phone is lost.", "Cora"),
		zv("已经", "yǐjīng", "already", "他已经走了。", "Tā yǐjīng zǒu le.", "He's already left.", "Blaze"),
	)

	s = addSkillL(db, zh, u, "The Passive with 被", "Say what happened to someone or something.", "Shield", "#00C2A8", 15, 230)
	l = addLesson(db, s, "被 — It Happened to Me", 1, 20,
		char(finch, "被 bèi marks the passive, often for something unwelcome: 我的自行车被偷了 — my bike was stolen. The doer can go after 被: 蛋糕被弟弟吃了 — the cake was eaten by my little brother."),
		mc("What does this mean?", "我的手机被偷了。", "My phone was stolen.", "My phone was stolen.", "I stole a phone.", "My phone is broken.", "I bought a phone."),
		mc("What does this mean?", "蛋糕被弟弟吃了。", "The cake was eaten by my little brother.", "The cake was eaten by my little brother.", "My little brother made a cake.", "My little brother wants cake.", "The cake was given to my brother."),
		mc("Choose the passive", "The window was broken by the wind.", "窗户被风吹坏了。", "窗户被风吹坏了。", "风被窗户吹坏了。", "窗户把风吹坏了。", "被窗户风吹坏了。"),
		mc("Which word marks the passive?", "我的钱包___偷了。", "被 (bèi)", "被 (bèi)", "把 (bǎ)", "比 (bǐ)", "跟 (gēn)"),
		mc("What does this mean?", "他被老师批评了。", "He was criticised by the teacher.", "He was criticised by the teacher.", "He criticised the teacher.", "The teacher praised him.", "He became a teacher."),
		speak("Riko", "我的自行车被人借走了。"),
	)
	addVocab(db, l,
		zv("被", "bèi", "(passive marker)", "我被雨淋湿了。", "Wǒ bèi yǔ lín shī le.", "I got soaked by the rain.", finch),
		zv("偷", "tōu", "to steal", "钱被偷了。", "Qián bèi tōu le.", "The money was stolen.", "Riko"),
		zv("坏", "huài", "broken; bad", "电脑坏了。", "Diànnǎo huài le.", "The computer's broken.", "Riko"),
		zv("批评", "pīpíng", "to criticise", "别批评他。", "Bié pīpíng tā.", "Don't criticise him.", finch),
		zv("借", "jiè", "to borrow; lend", "我能借你的笔吗？", "Wǒ néng jiè nǐ de bǐ ma?", "Can I borrow your pen?", "Cora"),
	)

	s = addSkillL(db, zh, u, "Health & the Doctor", "Describe symptoms and get advice.", "HeartPulse", "#FF5C5C", 16, 260)
	l = addLesson(db, s, "Seeing a Doctor", 1, 20,
		char("Mira", "At the doctor's: 我哪儿不舒服 — where I feel unwell; 头疼 headache; 发烧 fever; 咳嗽 cough. The doctor may say 多喝水，多休息 — drink lots of water, rest a lot."),
		mc("What does this mean?", "我头疼，还有点儿发烧。", "I have a headache and a slight fever.", "I have a headache and a slight fever.", "My head is hot, I'm fine.", "I'm going to the hospital.", "I had a fever yesterday."),
		mc("What is the doctor advising?", "你应该多休息，多喝水。", "Rest more and drink more water.", "Rest more and drink more water.", "Exercise more.", "Eat less.", "Go back to work."),
		mc("Choose the Chinese", "cough", "咳嗽 (késou)", "咳嗽 (késou)", "发烧 (fāshāo)", "感冒 (gǎnmào)", "头疼 (tóuténg)"),
		mc("What does this mean?", "你哪儿不舒服？", "Where does it hurt? / What's wrong?", "Where does it hurt? / What's wrong?", "Where are you comfortable?", "Where do you live?", "Are you comfortable?"),
		mc("Complete the sentence", "你___去医院看看。(You should go to the hospital.)", "应该", "应该", "已经", "比", "被"),
		speak("Mira", "医生，我感冒了，嗓子很疼。"),
	)
	addVocab(db, l,
		zv("舒服", "shūfu", "comfortable; well", "我不舒服。", "Wǒ bù shūfu.", "I don't feel well.", "Mira"),
		zv("发烧", "fāshāo", "to have a fever", "孩子发烧了。", "Háizi fāshāo le.", "The child has a fever.", "Mira"),
		zv("感冒", "gǎnmào", "a cold; to catch a cold", "我感冒了。", "Wǒ gǎnmào le.", "I've caught a cold.", "Mira"),
		zv("应该", "yīnggāi", "should", "你应该休息。", "Nǐ yīnggāi xiūxi.", "You should rest.", "Mira"),
		zv("医院", "yīyuàn", "hospital", "他在医院工作。", "Tā zài yīyuàn gōngzuò.", "He works at a hospital.", "Mira"),
	)

	s = addSkillL(db, zh, u, "Work & the Office", "Meetings, schedules and requests at work.", "Briefcase", "#17A3DD", 17, 290)
	l = addLesson(db, s, "At the Office", 1, 20,
		char("Zephyr", "Office Chinese: 开会 to hold a meeting; 加班 to work overtime; 请假 to ask for leave; 同事 colleague. Polite requests start with 请 or 麻烦你… (sorry to trouble you…)."),
		mc("What does this mean?", "下午三点开会。", "There's a meeting at 3 p.m.", "There's a meeting at 3 p.m.", "The meeting ended at 3.", "The office opens at 3.", "I'll be free at 3 p.m."),
		mc("What does this mean?", "我明天想请一天假。", "I'd like to take a day off tomorrow.", "I'd like to take a day off tomorrow.", "I'll work overtime tomorrow.", "I start tomorrow.", "I'm on holiday for a week."),
		mc("Choose the Chinese", "colleague", "同事 (tóngshì)", "同事 (tóngshì)", "经理 (jīnglǐ)", "客人 (kèrén)", "朋友 (péngyou)"),
		mc("What does this mean?", "麻烦你帮我打印一下。", "Sorry to trouble you — could you print this for me?", "Sorry to trouble you — could you print this for me?", "The printer is broken.", "I printed it for you.", "Don't print it."),
		mc("What does this mean?", "我们公司经常加班。", "Our company often works overtime.", "Our company often works overtime.", "Our company never works late.", "Our company is closing.", "We often have meetings."),
		speak("Zephyr", "经理，今天的会几点开始？"),
	)
	addVocab(db, l,
		zv("开会", "kāihuì", "to have a meeting", "我在开会。", "Wǒ zài kāihuì.", "I'm in a meeting.", "Zephyr"),
		zv("加班", "jiābān", "to work overtime", "今天要加班。", "Jīntiān yào jiābān.", "I have to work late today.", "Zephyr"),
		zv("请假", "qǐngjià", "to ask for leave", "我想请假。", "Wǒ xiǎng qǐngjià.", "I'd like to take leave.", "Zephyr"),
		zv("同事", "tóngshì", "colleague", "他是我的同事。", "Tā shì wǒ de tóngshì.", "He's my colleague.", "Zephyr"),
		zv("麻烦", "máfan", "trouble; to bother", "麻烦你了！", "Máfan nǐ le!", "Sorry for the trouble!", "Zephyr"),
	)
}

// ───────────────────────── HSK 4 · B2 — Adult Life ─────────────────────────

func seedMandarinHSK4(db *gorm.DB) {
	const zh = "zh"
	const u = "B2 · HSK 4 — Adult Life"
	finch := "Professor Finch"

	s := addSkillL(db, zh, u, "Renting & Housing", "Find a flat, negotiate rent, sign a contract.", "Home", "#6C3FC5", 18, 400)
	l := addLesson(db, s, "Renting a Flat", 1, 22,
		char("Riko", "Renting: 租房 to rent a place; 房租 the rent; 押金 deposit; 合同 contract; 房东 landlord; 中介 agent. 押一付三 means one month's deposit plus three months' rent upfront — common in Chinese cities."),
		mc("What does 押一付三 mean?", "yā yī fù sān", "One month's deposit, three months' rent upfront", "One month's deposit, three months' rent upfront", "Pay one month, get three free", "A three-month contract", "One room, three people"),
		mc("What does this mean?", "房租每个月四千块，水电费另算。", "Rent is 4,000 a month; utilities are extra.", "Rent is 4,000 a month; utilities are extra.", "Rent includes water and electricity.", "Rent is 4,000 a year.", "The deposit is 4,000."),
		mc("Choose the Chinese", "landlord", "房东 (fángdōng)", "房东 (fángdōng)", "中介 (zhōngjiè)", "邻居 (línjū)", "房客 (fángkè)"),
		mc("What is the tenant asking?", "如果提前退租，押金能退吗？", "If I move out early, can I get the deposit back?", "If I move out early, can I get the deposit back?", "Can I pay the rent late?", "Can I move in early?", "Is the deposit included in the rent?"),
		mc("Choose the best reply", "房东：合同签一年。 你：能不能先签半年？", "Asking for a six-month contract instead", "Asking for a six-month contract instead", "Agreeing to two years", "Refusing to sign", "Asking for a lower deposit"),
		speak("Riko", "这套房子离地铁站很近，房租也合适，我想签合同。"),
	)
	addVocab(db, l,
		zv("租", "zū", "to rent", "我在租房子。", "Wǒ zài zū fángzi.", "I'm renting a place.", "Riko"),
		zv("房租", "fángzū", "rent", "房租太贵了。", "Fángzū tài guì le.", "The rent is too high.", "Riko"),
		zv("押金", "yājīn", "deposit", "押金是一个月房租。", "Yājīn shì yí gè yuè fángzū.", "The deposit is one month's rent.", "Riko"),
		zv("合同", "hétong", "contract", "请先看合同。", "Qǐng xiān kàn hétong.", "Please read the contract first.", "Riko"),
		zv("房东", "fángdōng", "landlord", "房东人很好。", "Fángdōng rén hěn hǎo.", "The landlord is nice.", "Riko"),
	)

	s = addSkillL(db, zh, u, "Technology & Daily Life", "Mobile payment, apps and online shopping.", "Smartphone", "#00C2A8", 19, 430)
	l = addLesson(db, s, "A Cashless Society", 1, 22,
		char("Pip", "China is largely cashless: people 扫码支付 (scan a QR code to pay) with 微信 WeChat or 支付宝 Alipay. You can 网购 (shop online), 点外卖 (order delivery) and 叫车 (hail a car) — all from your phone."),
		mc("What does 扫码 mean?", "sǎo mǎ", "To scan a (QR) code", "To scan a (QR) code", "To sweep the floor", "To send a message", "To print a receipt"),
		mc("What does this mean?", "现在很多人出门不带现金，用手机支付。", "Many people now go out without cash and pay by phone.", "Many people now go out without cash and pay by phone.", "Many people lose their phones.", "Cash is more popular than phones.", "People don't go out any more."),
		mc("Choose the Chinese", "to order food delivery", "点外卖 (diǎn wàimài)", "点外卖 (diǎn wàimài)", "网购 (wǎnggòu)", "叫车 (jiào chē)", "做饭 (zuò fàn)"),
		mc("What is the concern?", "移动支付虽然方便，但是也有安全问题。", "Mobile payment is convenient but has security risks.", "Mobile payment is convenient but has security risks.", "Mobile payment is slow.", "Mobile payment isn't allowed.", "Mobile payment is too expensive."),
		mc("What does this mean?", "网上买的东西如果不满意，七天内可以退货。", "If you're not satisfied with online purchases, you can return them within 7 days.", "If you're not satisfied with online purchases, you can return them within 7 days.", "Online goods arrive in 7 days.", "You can't return online purchases.", "Online shopping has a 7-day sale."),
		speak("Pip", "我用手机扫码付款，又快又方便。"),
	)
	addVocab(db, l,
		zv("支付", "zhīfù", "to pay", "可以用手机支付吗？", "Kěyǐ yòng shǒujī zhīfù ma?", "Can I pay by phone?", "Pip"),
		zv("扫码", "sǎo mǎ", "to scan a QR code", "请扫码。", "Qǐng sǎo mǎ.", "Please scan the code.", "Pip"),
		zv("方便", "fāngbiàn", "convenient", "网购很方便。", "Wǎnggòu hěn fāngbiàn.", "Online shopping is convenient.", "Pip"),
		zv("安全", "ānquán", "safe; security", "密码要安全。", "Mìmǎ yào ānquán.", "Passwords must be secure.", "Pip"),
		zv("退货", "tuìhuò", "to return goods", "我想退货。", "Wǒ xiǎng tuìhuò.", "I'd like to return this.", "Pip"),
	)

	s = addSkillL(db, zh, u, "Feelings & Opinions", "Express emotions and argue a point.", "MessageCircle", "#FF5C5C", 20, 460)
	l = addLesson(db, s, "In My Opinion", 1, 22,
		char(finch, "Giving opinions: 我觉得… I think…; 在我看来… in my view…; 一方面…另一方面… on one hand… on the other. To concede: 不过 / 然而 — however. To conclude: 总的来说 — all in all."),
		mc("What does this mean?", "在我看来，学外语最重要的是坚持。", "In my view, the most important thing in learning a language is persistence.", "In my view, the most important thing in learning a language is persistence.", "I watch foreign films to learn.", "Learning languages is not important.", "Persistence is easy for me."),
		mc("Complete the structure", "一方面工资很高，___工作压力也很大。", "另一方面", "另一方面", "所以", "因为", "首先"),
		mc("Which phrase concludes an argument?", "all in all", "总的来说 (zǒng de lái shuō)", "总的来说 (zǒng de lái shuō)", "比如说 (bǐrú shuō)", "首先 (shǒuxiān)", "另外 (lìngwài)"),
		mc("What does this mean?", "他对这个结果非常失望。", "He's very disappointed with this result.", "He's very disappointed with this result.", "He's very pleased with the result.", "He expected this result.", "He doesn't care about the result."),
		mc("Choose the Chinese", "I'm worried about the exam.", "我很担心考试。", "我很担心考试。", "我很激动考试。", "我对考试很满意。", "考试让我很放心。"),
		speak("Mira", "我觉得在城市生活虽然方便，但是压力也很大。"),
	)
	addVocab(db, l,
		zv("觉得", "juéde", "to think, feel", "你觉得怎么样？", "Nǐ juéde zěnmeyàng?", "What do you think?", finch),
		zv("坚持", "jiānchí", "to persist", "坚持就是胜利。", "Jiānchí jiù shì shènglì.", "Persistence is victory.", finch),
		zv("压力", "yālì", "pressure, stress", "工作压力很大。", "Gōngzuò yālì hěn dà.", "Work is very stressful.", "Mira"),
		zv("失望", "shīwàng", "disappointed", "别让我失望。", "Bié ràng wǒ shīwàng.", "Don't let me down.", "Mira"),
		zv("担心", "dānxīn", "to worry", "别担心！", "Bié dānxīn!", "Don't worry!", "Mira"),
	)

	s = addSkillL(db, zh, u, "Jobs & Recruitment", "Interviews, CVs and workplace complaints.", "Briefcase", "#17A3DD", 21, 490)
	l = addLesson(db, s, "The Job Interview", 1, 22,
		char("Zephyr", "Recruitment: 招聘 recruiting; 应聘 applying; 简历 CV; 面试 interview; 经验 experience; 工资/薪水 salary; 福利 benefits. Interviewers often ask 你为什么想来我们公司？— why do you want to join us?"),
		mc("What does this mean?", "我们公司正在招聘销售经理。", "Our company is recruiting a sales manager.", "Our company is recruiting a sales manager.", "Our sales manager is leaving.", "We sell managers.", "Our company has no vacancies."),
		mc("Choose the Chinese", "job interview", "面试 (miànshì)", "面试 (miànshì)", "考试 (kǎoshì)", "简历 (jiǎnlì)", "合同 (hétong)"),
		mc("What does the interviewer want to know?", "请简单介绍一下你的工作经验。", "A brief summary of your work experience", "A brief summary of your work experience", "Your expected salary", "Your hobbies", "Your home address"),
		mc("What does this mean?", "这份工作工资不错，而且福利也很好。", "This job pays well, and the benefits are good too.", "This job pays well, and the benefits are good too.", "This job pays badly but has benefits.", "This job has no benefits.", "This job is temporary."),
		mc("What does this mean?", "我对你们的服务很不满意，我要投诉。", "I'm very unhappy with your service — I want to make a complaint.", "I'm very unhappy with your service — I want to make a complaint.", "I'm very happy with your service.", "I want to work for your service.", "I'd like to thank your staff."),
		speak("Zephyr", "我有三年的工作经验，而且会说英语和汉语。"),
	)
	addVocab(db, l,
		zv("招聘", "zhāopìn", "to recruit", "公司在招聘。", "Gōngsī zài zhāopìn.", "The company is hiring.", "Zephyr"),
		zv("面试", "miànshì", "interview", "明天有面试。", "Míngtiān yǒu miànshì.", "I have an interview tomorrow.", "Zephyr"),
		zv("经验", "jīngyàn", "experience", "他很有经验。", "Tā hěn yǒu jīngyàn.", "He's very experienced.", "Zephyr"),
		zv("简历", "jiǎnlì", "CV, résumé", "请发简历。", "Qǐng fā jiǎnlì.", "Please send your CV.", "Zephyr"),
		zv("投诉", "tóusù", "to complain (formally)", "我要投诉。", "Wǒ yào tóusù.", "I want to file a complaint.", "Zephyr"),
	)
}

// ───────────────────────── HSK 5 · C1 — Analytical Chinese ─────────────────────────

func seedMandarinHSK5(db *gorm.DB) {
	const zh = "zh"
	const u = "C1 · HSK 5 — Analytical Chinese"
	finch := "Professor Finch"

	s := addSkillL(db, zh, u, "Insurance & Finance", "Policies, claims and personal finance.", "Landmark", "#6C3FC5", 22, 700)
	l := addLesson(db, s, "Understanding a Policy", 1, 25,
		char(finch, "保险 insurance; 投保 to take out a policy; 理赔 a claim settlement; 保费 premium; 受益人 beneficiary. Formal texts prefer 若 (if), 须 (must) and 予以 (to give) over everyday words."),
		mc("What does 理赔 mean?", "lǐpéi", "Settling an insurance claim", "Settling an insurance claim", "Paying the premium", "Cancelling a policy", "Opening a bank account"),
		mc("What does this mean?", "若发生意外，保险公司将予以赔偿。", "If an accident occurs, the insurer will pay compensation.", "If an accident occurs, the insurer will pay compensation.", "Accidents are not covered.", "The insurer causes accidents.", "You must pay if there's an accident."),
		mc("Which word is formal for 'must'?", "合同期间，租客___按时交租。", "须 (xū)", "须 (xū)", "要 (yào)", "得 (děi)", "会 (huì)"),
		mc("What does this mean?", "理财不能只看收益，还要考虑风险。", "When managing money, consider risk, not just returns.", "When managing money, consider risk, not just returns.", "Only returns matter in finance.", "Avoid all risk.", "Profit is guaranteed."),
		mc("Choose the Chinese", "premium (insurance fee)", "保费 (bǎofèi)", "保费 (bǎofèi)", "利息 (lìxī)", "工资 (gōngzī)", "押金 (yājīn)"),
		speak("Zephyr", "投保之前，务必仔细阅读保险条款。"),
	)
	addVocab(db, l,
		zv("保险", "bǎoxiǎn", "insurance", "我买了医疗保险。", "Wǒ mǎi le yīliáo bǎoxiǎn.", "I bought health insurance.", finch),
		zv("理赔", "lǐpéi", "claim settlement", "理赔需要三天。", "Lǐpéi xūyào sān tiān.", "The claim takes three days.", finch),
		zv("风险", "fēngxiǎn", "risk", "投资有风险。", "Tóuzī yǒu fēngxiǎn.", "Investment carries risk.", finch),
		zv("收益", "shōuyì", "returns, earnings", "收益不错。", "Shōuyì búcuò.", "The returns are good.", finch),
		zv("务必", "wùbì", "must, be sure to", "务必准时到达。", "Wùbì zhǔnshí dàodá.", "Be sure to arrive on time.", finch),
	)

	s = addSkillL(db, zh, u, "Management & Teamwork", "Lead teams, resolve conflict, set goals.", "Users", "#00C2A8", 23, 730)
	l = addLesson(db, s, "Leading a Team", 1, 25,
		char("Zephyr", "Management talk: 制定目标 set goals; 分配任务 assign tasks; 沟通 communicate; 协调 coordinate; 提高效率 improve efficiency. 与其…不如… means 'rather than… it's better to…'."),
		mc("What does this mean?", "与其互相抱怨，不如一起想办法。", "Rather than complaining to each other, let's find a solution together.", "Rather than complaining to each other, let's find a solution together.", "Complaining is better than working.", "We have no solution.", "Let's complain together."),
		mc("Choose the Chinese", "to improve efficiency", "提高效率 (tígāo xiàolǜ)", "提高效率 (tígāo xiàolǜ)", "降低成本 (jiàngdī chéngběn)", "分配任务 (fēnpèi rènwu)", "制定目标 (zhìdìng mùbiāo)"),
		mc("What does this mean?", "良好的沟通是团队成功的关键。", "Good communication is the key to a team's success.", "Good communication is the key to a team's success.", "Teams should avoid communication.", "Success needs no team.", "The key to the office is lost."),
		mc("Complete the structure", "___他经验不足，___态度非常认真。", "尽管……但……", "尽管……但……", "因为……所以……", "只要……就……", "如果……就……"),
		mc("What does 协调 mean?", "xiétiáo", "To coordinate", "To coordinate", "To compete", "To criticise", "To resign"),
		speak("Zephyr", "我们应该先制定明确的目标，再合理分配任务。"),
	)
	addVocab(db, l,
		zv("沟通", "gōutōng", "to communicate", "我们需要多沟通。", "Wǒmen xūyào duō gōutōng.", "We need to communicate more.", "Zephyr"),
		zv("效率", "xiàolǜ", "efficiency", "效率很高。", "Xiàolǜ hěn gāo.", "Very efficient.", "Zephyr"),
		zv("团队", "tuánduì", "team", "我们的团队很强。", "Wǒmen de tuánduì hěn qiáng.", "Our team is strong.", "Zephyr"),
		zv("目标", "mùbiāo", "goal, target", "完成目标。", "Wánchéng mùbiāo.", "Achieve the goal.", "Zephyr"),
		zv("尽管", "jǐnguǎn", "although, despite", "尽管很累，他还在工作。", "Jǐnguǎn hěn lèi, tā hái zài gōngzuò.", "Despite being tired, he's still working.", "Zephyr"),
	)

	s = addSkillL(db, zh, u, "Research & Findings", "Read and report research results.", "FlaskConical", "#FF5C5C", 24, 760)
	l = addLesson(db, s, "What the Study Shows", 1, 25,
		char(finch, "Reporting research: 研究表明… studies show…; 数据显示… data show…; 据调查… according to a survey…; 由此可见… it can thus be seen…; 呈…趋势 shows a … trend."),
		mc("What does this mean?", "研究表明，适量运动有助于改善睡眠。", "Studies show moderate exercise helps improve sleep.", "Studies show moderate exercise helps improve sleep.", "Exercise makes sleep worse.", "Research on sleep is incomplete.", "Sleeping improves exercise."),
		mc("What does this mean?", "近年来，网购人数呈上升趋势。", "In recent years, the number of online shoppers has been rising.", "In recent years, the number of online shoppers has been rising.", "Online shopping is declining.", "Online shopping has stayed the same.", "Few people shop online."),
		mc("Which phrase draws a conclusion?", "it can thus be seen that…", "由此可见 (yóu cǐ kě jiàn)", "由此可见 (yóu cǐ kě jiàn)", "据调查 (jù diàochá)", "比如 (bǐrú)", "首先 (shǒuxiān)"),
		mc("Choose the Chinese", "according to a survey", "据调查 (jù diàochá)", "据调查 (jù diàochá)", "由此可见 (yóu cǐ kě jiàn)", "总而言之 (zǒng ér yán zhī)", "换句话说 (huàn jù huà shuō)"),
		mc("What does this mean?", "该结论仍需进一步验证。", "This conclusion still needs further verification.", "This conclusion still needs further verification.", "This conclusion is proven.", "This conclusion is wrong.", "This conclusion is old."),
		speak("Mira", "数据显示，年轻人的阅读时间正在逐年减少。"),
	)
	addVocab(db, l,
		zv("研究", "yánjiū", "research; to study", "研究表明…", "Yánjiū biǎomíng…", "Research shows…", finch),
		zv("表明", "biǎomíng", "to show, indicate", "结果表明…", "Jiéguǒ biǎomíng…", "The results show…", finch),
		zv("趋势", "qūshì", "trend", "上升趋势。", "Shàngshēng qūshì.", "An upward trend.", finch),
		zv("数据", "shùjù", "data", "数据很准确。", "Shùjù hěn zhǔnquè.", "The data are accurate.", finch),
		zv("验证", "yànzhèng", "to verify", "需要验证。", "Xūyào yànzhèng.", "It needs verifying.", finch),
	)

	s = addSkillL(db, zh, u, "Mental Health & Society", "Wellbeing, community and social issues.", "Brain", "#17A3DD", 25, 790)
	l = addLesson(db, s, "Wellbeing in Modern Life", 1, 25,
		char("Mira", "心理健康 mental health; 焦虑 anxiety; 缓解压力 relieve stress; 社区 community; 志愿者 volunteer. 不仅…而且… — not only… but also…"),
		mc("What does this mean?", "心理健康和身体健康同样重要。", "Mental health matters as much as physical health.", "Mental health matters as much as physical health.", "Physical health matters more.", "Mental health isn't important.", "Health is impossible."),
		mc("Choose the Chinese", "to relieve stress", "缓解压力 (huǎnjiě yālì)", "缓解压力 (huǎnjiě yālì)", "增加压力 (zēngjiā yālì)", "保持联系 (bǎochí liánxì)", "提高效率 (tígāo xiàolǜ)"),
		mc("Complete the structure", "志愿活动___帮助了别人，___让自己更快乐。", "不仅……而且……", "不仅……而且……", "虽然……但是……", "因为……所以……", "与其……不如……"),
		mc("What does 焦虑 mean?", "jiāolǜ", "Anxiety; anxious", "Anxiety; anxious", "Calm", "Excited", "Bored"),
		mc("What does this mean?", "社区为老人提供了免费的健康咨询。", "The community offers free health advice to elderly people.", "The community offers free health advice to elderly people.", "The elderly pay for health advice.", "The community has no services.", "Old people run the community."),
		speak("Mira", "适当的运动和充足的睡眠可以有效缓解焦虑。"),
	)
	addVocab(db, l,
		zv("心理", "xīnlǐ", "psychological, mental", "心理健康很重要。", "Xīnlǐ jiànkāng hěn zhòngyào.", "Mental health is important.", "Mira"),
		zv("焦虑", "jiāolǜ", "anxious; anxiety", "他很焦虑。", "Tā hěn jiāolǜ.", "He's very anxious.", "Mira"),
		zv("缓解", "huǎnjiě", "to relieve, ease", "缓解疼痛。", "Huǎnjiě téngtòng.", "Ease the pain.", "Mira"),
		zv("社区", "shèqū", "community", "我们社区很安静。", "Wǒmen shèqū hěn ānjìng.", "Our community is quiet.", "Mira"),
		zv("志愿者", "zhìyuànzhě", "volunteer", "我是志愿者。", "Wǒ shì zhìyuànzhě.", "I'm a volunteer.", "Mira"),
	)
}

// ───────────────────────── HSK 6 · C2 — Formal & Public Chinese ─────────────────────────

func seedMandarinHSK6(db *gorm.DB) {
	const zh = "zh"
	const u = "C2 · HSK 6 — Formal & Public Chinese"
	finch := "Professor Finch"

	s := addSkillL(db, zh, u, "AI & Technology Debate", "Discuss artificial intelligence and its impact.", "Cpu", "#6C3FC5", 26, 1000)
	l := addLesson(db, s, "人工智能的利与弊", 1, 28,
		char(finch, "Debate register: 毋庸置疑 beyond doubt; 不容忽视 cannot be ignored; 一味 blindly; 取而代之 to replace; 宁可…也不… would rather… than…. At this level questions are in Chinese."),
		mc("这句话是什么意思？", "人工智能的发展毋庸置疑地改变了我们的生活。", "AI has undoubtedly changed our lives.", "AI has undoubtedly changed our lives.", "AI has not changed our lives.", "It's doubtful AI matters.", "AI is still science fiction."),
		mc("“不容忽视”的意思是：", "隐私问题不容忽视。", "不能不重视", "不能不重视", "可以忽略", "已经解决", "不太重要"),
		mc("选择正确的词语", "我们不应该___地依赖技术。", "一味", "一味", "毋庸", "宁可", "何况"),
		mc("这句话是什么意思？", "有人担心机器人终将取而代之。", "Some worry robots will eventually replace (humans).", "Some worry robots will eventually replace (humans).", "Robots will never replace humans.", "People want to become robots.", "Robots are worried about people."),
		mc("选择正确的结构", "我___多花点时间，___不交一份粗糙的报告。", "宁可……也……", "宁可……也……", "不仅……而且……", "因为……所以……", "虽然……但是……"),
		speak("Zephyr", "技术本身无所谓好坏，关键在于人类如何使用它。"),
	)
	addVocab(db, l,
		zv("人工智能", "réngōng zhìnéng", "artificial intelligence", "人工智能发展很快。", "Réngōng zhìnéng fāzhǎn hěn kuài.", "AI is developing fast.", finch),
		zv("毋庸置疑", "wúyōng zhìyí", "beyond doubt", "这一点毋庸置疑。", "Zhè yì diǎn wúyōng zhìyí.", "This is beyond doubt.", finch),
		zv("不容忽视", "bù róng hūshì", "cannot be ignored", "风险不容忽视。", "Fēngxiǎn bù róng hūshì.", "The risks can't be ignored.", finch),
		zv("一味", "yíwèi", "blindly; persistently", "不要一味模仿。", "Búyào yíwèi mófǎng.", "Don't imitate blindly.", finch),
		zv("取而代之", "qǔ ér dài zhī", "to replace", "新技术取而代之。", "Xīn jìshù qǔ ér dài zhī.", "New technology takes its place.", finch),
	)

	s = addSkillL(db, zh, u, "Environment & Infrastructure", "Disasters, investment and sustainable cities.", "Leaf", "#00C2A8", 27, 1030)
	l = addLesson(db, s, "可持续发展", 1, 28,
		char("Riko", "Public-affairs vocabulary: 可持续发展 sustainable development; 基础设施 infrastructure; 自然灾害 natural disasters; 防范 to guard against; 兼顾 to give consideration to both."),
		mc("“基础设施”指的是：", "jīchǔ shèshī", "道路、桥梁、电网等公共设施", "道路、桥梁、电网等公共设施", "私人住宅", "网上商店", "学校课程"),
		mc("这句话是什么意思？", "经济发展必须兼顾环境保护。", "Economic growth must also take environmental protection into account.", "Economic growth must also take environmental protection into account.", "The economy matters more than the environment.", "Environmental protection stops growth.", "The economy and environment are unrelated."),
		mc("选择正确的词语", "政府加大投入，以___自然灾害。", "防范", "防范", "兼顾", "导致", "忽视"),
		mc("这句话是什么意思？", "过度开发导致生态环境日益恶化。", "Over-development is causing the ecosystem to deteriorate day by day.", "Over-development is causing the ecosystem to deteriorate day by day.", "Development is improving the environment.", "The environment is stable.", "Development has stopped."),
		mc("“可持续发展”的核心是：", "kě chíxù fāzhǎn", "满足当代需要，又不损害后代的利益", "满足当代需要，又不损害后代的利益", "最快地发展经济", "停止所有开发", "只保护动物"),
		speak("Riko", "城市化进程中，我们应当重视基础设施的长期规划。"),
	)
	addVocab(db, l,
		zv("可持续", "kě chíxù", "sustainable", "可持续发展。", "Kě chíxù fāzhǎn.", "Sustainable development.", "Riko"),
		zv("基础设施", "jīchǔ shèshī", "infrastructure", "改善基础设施。", "Gǎishàn jīchǔ shèshī.", "Improve infrastructure.", "Riko"),
		zv("自然灾害", "zìrán zāihài", "natural disaster", "预防自然灾害。", "Yùfáng zìrán zāihài.", "Prevent natural disasters.", "Riko"),
		zv("兼顾", "jiāngù", "to balance both", "兼顾工作和家庭。", "Jiāngù gōngzuò hé jiātíng.", "Balance work and family.", "Riko"),
		zv("恶化", "èhuà", "to deteriorate", "情况恶化了。", "Qíngkuàng èhuà le.", "The situation worsened.", "Riko"),
	)

	s = addSkillL(db, zh, u, "Classical Culture & Idioms", "成语 and classical echoes in modern Chinese.", "BookOpen", "#F5A623", 28, 1060)
	l = addLesson(db, s, "成语 — Four-Character Idioms", 1, 28,
		char("Nana", "成语 chéngyǔ are four-character idioms, often from classical stories. 画蛇添足 'drawing legs on a snake' — overdoing it; 守株待兔 'waiting by a stump for a rabbit' — waiting for luck instead of working; 半途而废 giving up halfway."),
		mc("“画蛇添足”的意思是：", "huà shé tiān zú", "做多余的事，反而坏事", "做多余的事，反而坏事", "画画很认真", "努力学习", "动作很快"),
		mc("“守株待兔”批评的是：", "shǒu zhū dài tù", "不努力，只想等好运", "不努力，只想等好运", "太努力工作", "喜欢小动物", "做事太急"),
		mc("选择合适的成语", "学习汉语不能___，要坚持下去。", "半途而废", "半途而废", "画蛇添足", "守株待兔", "一举两得"),
		mc("“一举两得”的意思是：", "yì jǔ liǎng dé", "做一件事得到两个好处", "做一件事得到两个好处", "得到两次奖励", "一个人做两件事", "失去两样东西"),
		mc("这句话是什么意思？", "这篇文章已经很好了，再改就是画蛇添足。", "The article is already good; further edits would spoil it.", "The article is already good; further edits would spoil it.", "The article needs a drawing.", "The article is about snakes.", "The article needs more work."),
		speak("Nana", "学习贵在坚持，千万不能半途而废。"),
	)
	addVocab(db, l,
		zv("成语", "chéngyǔ", "idiom (four-character)", "这个成语很有名。", "Zhège chéngyǔ hěn yǒumíng.", "This idiom is famous.", "Nana"),
		zv("画蛇添足", "huà shé tiān zú", "to overdo it", "别画蛇添足了。", "Bié huà shé tiān zú le.", "Don't overdo it.", "Nana"),
		zv("守株待兔", "shǒu zhū dài tù", "to wait for luck", "不能守株待兔。", "Bù néng shǒu zhū dài tù.", "You can't just wait for luck.", "Nana"),
		zv("半途而废", "bàntú ér fèi", "to give up halfway", "不要半途而废。", "Búyào bàntú ér fèi.", "Don't give up halfway.", "Nana"),
		zv("一举两得", "yì jǔ liǎng dé", "to kill two birds with one stone", "这真是一举两得。", "Zhè zhēn shì yì jǔ liǎng dé.", "That's two birds with one stone.", "Nana"),
	)

	s = addSkillL(db, zh, u, "Formal Speeches & Debate", "Structure a speech and rebut an argument.", "Mic", "#FF5C5C", 29, 1090)
	l = addLesson(db, s, "演讲与辩论", 1, 28,
		char(finch, "Formal speech: 尊敬的各位来宾 respected guests; 首先…其次…最后… firstly… secondly… finally; 诚然… admittedly…; 恰恰相反 quite the opposite; 综上所述 in summary."),
		mc("演讲开头常说：", "正式演讲的称呼", "尊敬的各位来宾", "尊敬的各位来宾", "喂，你好", "再见了", "吃了吗"),
		mc("“诚然”用来：", "chéngrán", "先承认对方观点的合理之处", "先承认对方观点的合理之处", "表示完全反对", "表示结束", "表示时间"),
		mc("选择合适的词语", "这项政策不但没有增加负担，___减轻了人们的压力。", "恰恰相反，它", "恰恰相反，它", "综上所述，它", "首先，它", "尊敬的，它"),
		mc("“综上所述”一般出现在：", "zōng shàng suǒ shù", "结论部分", "结论部分", "开头", "举例", "提问"),
		mc("这句话是什么意思？", "对方辩友的观点看似合理，实则经不起推敲。", "The opponent's view seems reasonable but doesn't stand up to scrutiny.", "The opponent's view seems reasonable but doesn't stand up to scrutiny.", "The opponent is completely right.", "The opponent didn't speak.", "Both sides agree."),
		speak("Lumora", "综上所述，我方认为教育公平是社会公平的基础。"),
	)
	addVocab(db, l,
		zv("尊敬", "zūnjìng", "respected; to respect", "尊敬的老师。", "Zūnjìng de lǎoshī.", "Respected teacher.", finch),
		zv("诚然", "chéngrán", "admittedly", "诚然，这有困难。", "Chéngrán, zhè yǒu kùnnan.", "Admittedly, it's difficult.", finch),
		zv("恰恰相反", "qiàqià xiāngfǎn", "quite the opposite", "事实恰恰相反。", "Shìshí qiàqià xiāngfǎn.", "The truth is quite the opposite.", finch),
		zv("综上所述", "zōng shàng suǒ shù", "in summary", "综上所述…", "Zōng shàng suǒ shù…", "In summary…", finch),
		zv("推敲", "tuīqiāo", "to scrutinise, weigh words", "经不起推敲。", "Jīng bu qǐ tuīqiāo.", "Doesn't stand up to scrutiny.", finch),
	)
}

// ───────────────────────── Listening & reading sessions ─────────────────────────
//
// Lines are pure characters (they're read aloud); the translation carries Pinyin
// at HSK 1–2 and English throughout. Questions are in English up to HSK 4 and in
// Chinese from HSK 5, mirroring the German course.

func seedMandarinListening(db *gorm.DB) {
	const zh = "zh"
	finch := "Professor Finch"

	addListeningL(db, zh, "A1 · HSK 1 — Survival Chinese", "第一天上课 — First Day of Class",
		"Professor Finch meets a new student. Listen, then answer.", 1, 20,
		[]models.ListeningMatch{
			lm("你好", "hello"), lm("你叫什么名字？", "What's your name?"),
			lm("你是哪国人？", "Where are you from?"), lm("认识你很高兴", "Nice to meet you"),
		},
		[]models.ListeningLine{
			ln(finch, "你好！欢迎来到汉语课。你叫什么名字？", "Nǐ hǎo! Huānyíng lái dào Hànyǔ kè. Nǐ jiào shénme míngzi? — Hello! Welcome to Chinese class. What's your name?"),
			ln("Lumora", "老师好！我叫露茉拉。", "Lǎoshī hǎo! Wǒ jiào Lùmòlā. — Hello, teacher! My name is Lumora."),
			ln(finch, "你是哪国人？", "Nǐ shì nǎ guó rén? — Which country are you from?"),
			ln("Lumora", "我是肯尼亚人。我今年二十岁。", "Wǒ shì Kěnníyà rén. Wǒ jīnnián èrshí suì. — I'm Kenyan. I'm twenty this year."),
			ln(finch, "认识你很高兴！", "Rènshi nǐ hěn gāoxìng! — Nice to meet you!"),
		},
		[]models.ListeningQuestion{
			lq("What is the student's name?", "Lumora", "Lumora", "Mira", "Cora", "Finch"),
			lq("Where is she from?", "Kenya", "Kenya", "China", "Japan", "Germany"),
			lq("How old is she?", "20", "20", "12", "22", "10"),
		},
	)
	addListeningL(db, zh, "A1 · HSK 1 — Survival Chinese", "在饭馆 — At the Restaurant",
		"Cora orders lunch. Listen, then answer.", 2, 20,
		[]models.ListeningMatch{
			lm("服务员", "waiter"), lm("我要…", "I'd like…"),
			lm("一共", "altogether"), lm("多少钱？", "How much?"),
		},
		[]models.ListeningLine{
			ln("Riko", "你好！你要吃什么？", "Nǐ hǎo! Nǐ yào chī shénme? — Hello! What would you like to eat?"),
			ln("Cora", "我要一盘饺子和一碗米饭。", "Wǒ yào yì pán jiǎozi hé yì wǎn mǐfàn. — I'd like a plate of dumplings and a bowl of rice."),
			ln("Riko", "你喝什么？茶还是水？", "Nǐ hē shénme? Chá háishi shuǐ? — What will you drink? Tea or water?"),
			ln("Cora", "一杯茶，谢谢。一共多少钱？", "Yì bēi chá, xièxie. Yígòng duōshao qián? — A cup of tea, thanks. How much altogether?"),
			ln("Riko", "一共三十五块。", "Yígòng sānshíwǔ kuài. — Thirty-five yuan altogether."),
		},
		[]models.ListeningQuestion{
			lq("What does Cora order to eat?", "Dumplings and rice", "Dumplings and rice", "Noodles", "Rice only", "Dumplings and noodles"),
			lq("What does she drink?", "Tea", "Tea", "Water", "Coffee", "Juice"),
			lq("How much is it altogether?", "35 yuan", "35 yuan", "53 yuan", "30 yuan", "25 yuan"),
		},
	)
	addListeningL(db, zh, "A2 · HSK 2 — Everyday Life", "周末去哪儿了？ — Weekend Plans",
		"Two friends talk about last weekend. Listen, then answer.", 3, 22,
		[]models.ListeningMatch{
			lm("上个周末", "last weekend"), lm("爬山", "to climb a mountain"),
			lm("累", "tired"), lm("下次", "next time"),
		},
		[]models.ListeningLine{
			ln("Blaze", "米拉，你上个周末去哪儿了？", "Mǐlā, nǐ shàng ge zhōumò qù nǎr le? — Mira, where did you go last weekend?"),
			ln("Mira", "我和朋友去爬山了。你去过香山吗？", "Wǒ hé péngyou qù pá shān le. Nǐ qù guo Xiāngshān ma? — I went hiking with friends. Have you been to Fragrant Hill?"),
			ln("Blaze", "没去过。山上冷吗？", "Méi qù guo. Shān shang lěng ma? — Never. Was it cold up there?"),
			ln("Mira", "比城里冷一点儿，但是风景特别漂亮。虽然很累，但是很开心。", "Bǐ chéng lǐ lěng yìdiǎnr, dànshì fēngjǐng tèbié piàoliang. Suīrán hěn lèi, dànshì hěn kāixīn. — A bit colder than the city, but the view was beautiful. Tiring, but fun."),
			ln("Blaze", "下次叫上我吧！", "Xià cì jiào shang wǒ ba! — Take me next time!"),
		},
		[]models.ListeningQuestion{
			lq("What did Mira do last weekend?", "Went hiking with friends", "Went hiking with friends", "Stayed at home", "Went shopping", "Went to work"),
			lq("Has Blaze been to Fragrant Hill?", "No, never", "No, never", "Yes, once", "Yes, often", "He lives there"),
			lq("How was the weather on the mountain?", "A bit colder than the city", "A bit colder than the city", "Much hotter", "Rainy", "The same as the city"),
			lq("How did Mira feel?", "Tired but happy", "Tired but happy", "Bored", "Ill", "Angry"),
		},
	)
	addListeningL(db, zh, "B1 · HSK 3 — Social Chinese", "看医生 — At the Doctor's",
		"Zephyr sees a doctor. Listen, then answer.", 4, 24,
		[]models.ListeningMatch{
			lm("不舒服", "unwell"), lm("发烧", "fever"),
			lm("开药", "to prescribe medicine"), lm("休息", "to rest"),
		},
		[]models.ListeningLine{
			ln("Mira", "你哪儿不舒服？", "Where don't you feel well?"),
			ln("Zephyr", "我从昨天晚上开始头疼，嗓子也疼，还有点儿发烧。", "Since last night I've had a headache and a sore throat, and a slight fever."),
			ln("Mira", "我看看……三十八度。你感冒了，不太严重。", "Let me see… 38 degrees. You have a cold; it's not too serious."),
			ln("Mira", "我给你开点儿药，一天吃三次，饭后吃。这两天别去上班了，多喝水，多休息。", "I'll prescribe some medicine — three times a day, after meals. Don't go to work for the next two days; drink plenty of water and rest."),
			ln("Zephyr", "好的，谢谢医生。我得把药方拿给药房吧？", "OK, thank you, doctor. I should take the prescription to the pharmacy, right?"),
		},
		[]models.ListeningQuestion{
			lq("Since when has Zephyr felt ill?", "Last night", "Last night", "This morning", "Last week", "Two days ago"),
			lq("What is his temperature?", "38 degrees", "38 degrees", "36 degrees", "39 degrees", "37 degrees"),
			lq("How often should he take the medicine?", "Three times a day, after meals", "Three times a day, after meals", "Once a day", "Twice a day before meals", "Only at night"),
			lq("What does the doctor advise about work?", "Stay home for two days", "Stay home for two days", "Go to work as usual", "Work half days", "Stay home for a week"),
		},
	)
	addListeningL(db, zh, "B2 · HSK 4 — Adult Life", "看房子 — Viewing a Flat",
		"Riko views a flat with an agent. Listen, then answer.", 5, 26,
		[]models.ListeningMatch{
			lm("中介", "letting agent"), lm("押一付三", "1 month deposit + 3 months' rent"),
			lm("家具", "furniture"), lm("物业费", "property management fee"),
		},
		[]models.ListeningLine{
			ln("Pip", "这套房子是两室一厅，离地铁站走路只要五分钟，家具家电都是齐全的。", "This flat has two bedrooms and a living room, five minutes' walk from the metro, fully furnished with appliances."),
			ln("Riko", "看起来不错。房租多少？", "Looks good. How much is the rent?"),
			ln("Pip", "每个月四千五，押一付三。物业费房东出，水电费你自己付。", "4,500 a month, one month's deposit and three months upfront. The landlord pays the management fee; you pay water and electricity."),
			ln("Riko", "有点儿贵。如果我签两年合同，房租能便宜一点儿吗？", "A bit expensive. If I sign a two-year contract, could the rent be a bit lower?"),
			ln("Pip", "我帮你问问房东，签两年的话应该可以降到四千二。", "I'll ask the landlord for you. For two years it could probably come down to 4,200."),
		},
		[]models.ListeningQuestion{
			lq("How far is the flat from the metro?", "5 minutes on foot", "5 minutes on foot", "15 minutes on foot", "5 minutes by bus", "Next door"),
			lq("What is the original rent?", "4,500 a month", "4,500 a month", "4,200 a month", "5,400 a month", "4,500 a year"),
			lq("Who pays the property management fee?", "The landlord", "The landlord", "The tenant", "The agent", "Nobody"),
			lq("How could the rent be reduced?", "By signing a two-year contract", "By signing a two-year contract", "By paying in cash", "By renting without furniture", "By moving in later"),
		},
	)
	addListeningL(db, zh, "C1 · HSK 5 — Analytical Chinese", "研究报告 — A Research Briefing",
		"A researcher presents survey results. Listen, then answer in Chinese.", 6, 28,
		[]models.ListeningMatch{
			lm("调查", "survey"), lm("受访者", "respondents"),
			lm("比例", "proportion"), lm("由此可见", "it can thus be seen"),
		},
		[]models.ListeningLine{
			ln("Mira", "各位好，今天我向大家介绍一项关于年轻人阅读习惯的调查。", "Hello everyone. Today I'll present a survey on young people's reading habits."),
			ln("Mira", "这次调查共有两千名受访者，年龄在十八到三十岁之间。", "There were 2,000 respondents, aged between 18 and 30."),
			ln("Mira", "数据显示，超过六成的人主要通过手机阅读，而每天阅读纸质书的人不到百分之十五。", "The data show over 60% read mainly on their phones, while under 15% read paper books daily."),
			ln("Mira", "值得注意的是，虽然阅读总时间并没有减少，但是深度阅读的比例明显下降了。", "Notably, total reading time hasn't fallen, but the share of in-depth reading has dropped markedly."),
			ln("Mira", "由此可见，问题不在于读得少，而在于读得浅。", "So the problem isn't reading less — it's reading more shallowly."),
		},
		[]models.ListeningQuestion{
			lq("这次调查有多少受访者？", "两千名", "两千名", "两百名", "一千八百名", "三千名"),
			lq("大多数年轻人主要用什么阅读？", "手机", "手机", "纸质书", "电脑", "报纸"),
			lq("关于阅读总时间，调查发现了什么？", "并没有减少", "并没有减少", "大大减少了", "增加了一倍", "无法统计"),
			lq("研究者的结论是什么？", "问题在于读得浅", "问题在于读得浅", "年轻人不读书", "纸质书最受欢迎", "手机阅读有益"),
		},
	)
	addListeningL(db, zh, "C2 · HSK 6 — Formal & Public Chinese", "人工智能与就业 — AI and Employment",
		"An excerpt from a public debate. Listen, then answer in Chinese.", 7, 30,
		[]models.ListeningMatch{
			lm("取而代之", "to replace"), lm("诚然", "admittedly"),
			lm("转型", "transformation"), lm("综上所述", "in summary"),
		},
		[]models.ListeningLine{
			ln(finch, "尊敬的各位来宾，今天我们讨论的话题是：人工智能是否会导致大规模失业。", "Respected guests, today's topic is whether AI will cause mass unemployment."),
			ln("Zephyr", "诚然，许多重复性的工作正在被机器取而代之，这一点毋庸置疑。", "Admittedly, many repetitive jobs are being replaced by machines; that is beyond doubt."),
			ln("Zephyr", "然而，历史告诉我们，每一次技术革命在淘汰旧岗位的同时，也创造了大量新岗位。", "Yet history shows every technological revolution, while eliminating old jobs, also creates many new ones."),
			ln("Zephyr", "关键不在于技术本身，而在于社会能否及时为劳动者提供培训，帮助他们实现转型。", "The key is not the technology itself but whether society can train workers in time to help them transition."),
			ln("Zephyr", "综上所述，与其恐惧变化，不如积极应对。", "In summary, rather than fear change, we should respond to it actively."),
		},
		[]models.ListeningQuestion{
			lq("发言人承认了什么？", "很多重复性工作正被机器取代", "很多重复性工作正被机器取代", "人工智能没有影响", "所有工作都会消失", "机器无法工作"),
			lq("他用什么来支持自己的观点？", "历史上技术革命的经验", "历史上技术革命的经验", "个人经历", "政府的规定", "公司的数据"),
			lq("他认为关键在于什么？", "能否为劳动者提供培训", "能否为劳动者提供培训", "停止发展技术", "减少工作时间", "提高工资"),
			lq("他的结论是什么？", "与其恐惧，不如积极应对", "与其恐惧，不如积极应对", "应该禁止人工智能", "失业无法避免", "变化不会发生"),
		},
	)
}

func seedMandarinReading(db *gorm.DB) {
	const zh = "zh"

	addReadingL(db, zh, "A1 · HSK 1 — Survival Chinese", "我的家 — My Family",
		"A short self-introduction. Read, then answer.", 1, 20,
		[]models.ReadingLine{
			rl("我叫王明。(Wǒ jiào Wáng Míng.)", "My name is Wang Ming."),
			rl("我是中国人，我家在北京。(Wǒ shì Zhōngguó rén, wǒ jiā zài Běijīng.)", "I'm Chinese; my home is in Beijing."),
			rl("我家有四口人：爸爸、妈妈、姐姐和我。(Wǒ jiā yǒu sì kǒu rén: bàba, māma, jiějie hé wǒ.)", "There are four people in my family: dad, mum, my older sister and me."),
			rl("爸爸是医生，妈妈是老师。我是学生。(Bàba shì yīshēng, māma shì lǎoshī. Wǒ shì xuésheng.)", "Dad is a doctor and Mum is a teacher. I'm a student."),
		},
		[]models.ReadingQuestion{
			rq("Where is Wang Ming's home?", "Beijing", "Beijing", "Shanghai", "Kenya", "Hong Kong"),
			rq("How many people are in his family?", "4", "4", "3", "5", "6"),
			rq("What is his mother's job?", "Teacher", "Teacher", "Doctor", "Student", "Waiter"),
		},
	)
	addReadingL(db, zh, "A1 · HSK 1 — Survival Chinese", "去商店 — Going Shopping",
		"Wang Ming buys fruit. Read, then answer.", 2, 20,
		[]models.ReadingLine{
			rl("今天是星期六。王明去商店买水果。(Jīntiān shì xīngqīliù. Wáng Míng qù shāngdiàn mǎi shuǐguǒ.)", "Today is Saturday. Wang Ming goes to the shop to buy fruit."),
			rl("他要三个苹果和一斤香蕉。(Tā yào sān gè píngguǒ hé yì jīn xiāngjiāo.)", "He wants three apples and half a kilo of bananas."),
			rl("苹果五块，香蕉三块，一共八块钱。(Píngguǒ wǔ kuài, xiāngjiāo sān kuài, yígòng bā kuài qián.)", "The apples are 5 yuan, the bananas 3: 8 yuan altogether."),
			rl("王明用手机付钱。(Wáng Míng yòng shǒujī fù qián.)", "Wang Ming pays with his phone."),
		},
		[]models.ReadingQuestion{
			rq("What day is it?", "Saturday", "Saturday", "Sunday", "Monday", "Friday"),
			rq("How many apples does he buy?", "3", "3", "5", "8", "1"),
			rq("How much does he pay altogether?", "8 yuan", "8 yuan", "5 yuan", "3 yuan", "11 yuan"),
			rq("How does he pay?", "With his phone", "With his phone", "In cash", "By card", "He doesn't pay"),
		},
	)
	addReadingL(db, zh, "A2 · HSK 2 — Everyday Life", "我的春节 — My Spring Festival",
		"A student describes the New Year holiday. Read, then answer.", 3, 22,
		[]models.ReadingLine{
			rl("春节是中国最重要的节日。(Chūnjié shì Zhōngguó zuì zhòngyào de jiérì.)", "The Spring Festival is China's most important holiday."),
			rl("今年春节，我坐了八个小时的火车回老家。(Jīnnián Chūnjié, wǒ zuò le bā gè xiǎoshí de huǒchē huí lǎojiā.)", "This Spring Festival I took an eight-hour train back to my hometown."),
			rl("除夕晚上，我们一家人一边包饺子，一边看电视。(Chúxī wǎnshang, wǒmen yì jiā rén yìbiān bāo jiǎozi, yìbiān kàn diànshì.)", "On New Year's Eve, our family made dumplings while watching TV."),
			rl("奶奶给了我一个红包。虽然回家很累，但是我觉得很幸福。(Nǎinai gěi le wǒ yí gè hóngbāo. Suīrán huí jiā hěn lèi, dànshì wǒ juéde hěn xìngfú.)", "Grandma gave me a red envelope. Although travelling home was tiring, I felt very happy."),
		},
		[]models.ReadingQuestion{
			rq("How did the writer travel home?", "By train", "By train", "By plane", "By car", "By bus"),
			rq("How long was the journey?", "8 hours", "8 hours", "18 hours", "2 hours", "4 hours"),
			rq("What did the family do on New Year's Eve?", "Made dumplings while watching TV", "Made dumplings while watching TV", "Went to a restaurant", "Went to bed early", "Watched fireworks only"),
			rq("Who gave the writer a red envelope?", "Grandma", "Grandma", "Mum", "Dad", "A friend"),
		},
	)
	addReadingL(db, zh, "B1 · HSK 3 — Social Chinese", "我的第一份工作 — My First Job",
		"A young professional writes about starting work. Read, then answer.", 4, 24,
		[]models.ReadingLine{
			rl("去年大学毕业以后，我在一家科技公司找到了第一份工作。", "After graduating last year, I found my first job at a tech company."),
			rl("刚开始的时候，我什么都不懂，经常把文件放错地方，还被经理批评过一次。", "At first I knew nothing — I often put files in the wrong place and was even criticised by my manager once."),
			rl("后来，同事们耐心地教我，我也每天下班后自己学习。", "Later, my colleagues patiently taught me, and I studied on my own after work every day."),
			rl("现在我已经能独立完成工作了。我明白了：只要努力，就一定会进步。", "Now I can work independently. I've learned that as long as you work hard, you will improve."),
		},
		[]models.ReadingQuestion{
			rq("Where did the writer find a job?", "At a tech company", "At a tech company", "At a school", "At a hospital", "At a bank"),
			rq("What happened at the beginning?", "They made mistakes and were criticised", "They made mistakes and were criticised", "They were promoted", "They quit", "They became manager"),
			rq("Who helped the writer?", "Colleagues", "Colleagues", "Parents", "Teachers", "Nobody"),
			rq("What lesson did the writer learn?", "Hard work leads to progress", "Hard work leads to progress", "Work is too hard", "Managers are always right", "Studying is useless"),
		},
	)
	addReadingL(db, zh, "B2 · HSK 4 — Adult Life", "无现金社会 — The Cashless Society",
		"A short article on mobile payment. Read, then answer.", 5, 26,
		[]models.ReadingLine{
			rl("在中国的大城市，很多人出门已经不带钱包了。从买早餐到交房租，几乎所有事情都可以用手机完成。", "In China's big cities many people no longer carry a wallet. From breakfast to rent, almost everything can be done by phone."),
			rl("移动支付给生活带来了极大的方便，也推动了外卖、网购等新行业的发展。", "Mobile payment has made life extremely convenient and driven the growth of new industries like food delivery and online shopping."),
			rl("然而，并不是所有人都能适应这种变化。一些老年人不会使用智能手机，在生活中遇到了不少困难。", "However, not everyone can adapt. Some elderly people can't use smartphones and face many difficulties."),
			rl("此外，个人信息安全也成为人们关心的问题。专家建议，在享受便利的同时，也要保护好自己的隐私。", "In addition, personal data security has become a concern. Experts advise protecting your privacy while enjoying the convenience."),
		},
		[]models.ReadingQuestion{
			rq("What do many people in big cities no longer carry?", "A wallet", "A wallet", "A phone", "Keys", "A bag"),
			rq("Which new industries has mobile payment helped?", "Food delivery and online shopping", "Food delivery and online shopping", "Farming and fishing", "Banking only", "Tourism only"),
			rq("Who has difficulty with the change?", "Some elderly people", "Some elderly people", "Young people", "Shop owners", "Banks"),
			rq("What do experts advise?", "Protect your privacy while enjoying convenience", "Protect your privacy while enjoying convenience", "Stop using phones", "Only use cash", "Share your data freely"),
		},
	)
	addReadingL(db, zh, "C1 · HSK 5 — Analytical Chinese", "城市里的孤独 — Loneliness in the City",
		"An opinion piece on modern urban life. Read, then answer in Chinese.", 6, 28,
		[]models.ReadingLine{
			rl("随着城市化的加快，越来越多的年轻人离开家乡，到大城市打拼。", "As urbanisation accelerates, more and more young people leave home to strive in big cities."),
			rl("他们虽然拥有更多的工作机会，却常常感到孤独。研究表明，独居人群出现焦虑情绪的比例明显高于其他人群。", "Though they have more job opportunities, they often feel lonely. Research shows people living alone are markedly more prone to anxiety."),
			rl("对此，一些社区开始组织读书会、运动小组等活动，帮助居民建立联系。", "In response, some communities have begun organising book clubs and sports groups to help residents connect."),
			rl("由此可见，解决孤独问题不仅需要个人的努力，而且需要整个社会的参与。", "It can thus be seen that solving loneliness needs not only individual effort but the participation of society as a whole."),
		},
		[]models.ReadingQuestion{
			rq("年轻人为什么去大城市？", "为了更多的工作机会", "为了更多的工作机会", "为了照顾父母", "为了上小学", "为了退休"),
			rq("研究发现独居的人更容易怎么样？", "出现焦虑情绪", "出现焦虑情绪", "身体更健康", "收入更高", "更爱运动"),
			rq("社区采取了什么措施？", "组织读书会和运动小组", "组织读书会和运动小组", "提高房租", "关闭公园", "减少活动"),
			rq("作者的结论是什么？", "需要个人和社会共同努力", "需要个人和社会共同努力", "只靠个人就够了", "孤独无法解决", "应该离开城市"),
		},
	)
	addReadingL(db, zh, "C2 · HSK 6 — Formal & Public Chinese", "守株待兔的现代启示 — A Classical Fable Today",
		"A reflection on an ancient fable. Read, then answer in Chinese.", 7, 30,
		[]models.ReadingLine{
			rl("《韩非子》中记载了一个故事：宋国有个农夫，偶然看见一只兔子撞死在树桩上，便从此放下农具，天天守在树旁，等待下一只兔子。", "The Han Feizi records a story: a farmer in the state of Song chanced to see a rabbit die by running into a tree stump, so he put down his tools and waited by the tree every day for the next rabbit."),
			rl("结果，兔子再也没有出现，他的田地却荒芜了。", "In the end, no rabbit ever came, and his fields went to waste."),
			rl("两千多年后的今天，这个寓言依然发人深省。有人把一次偶然的成功当作必然的规律，一味地等待机会，而不肯付出努力。", "More than two thousand years later the fable still gives food for thought. Some take a chance success for an inevitable rule, waiting blindly for opportunity instead of making an effort."),
			rl("诚然，机遇固然重要，但机遇只青睐有准备的人。与其守株待兔，不如脚踏实地。", "Admittedly, opportunity matters, but it favours only the prepared. Rather than wait by the stump, one should keep one's feet on the ground."),
		},
		[]models.ReadingQuestion{
			rq("农夫为什么放下了农具？", "他想等下一只撞树的兔子", "他想等下一只撞树的兔子", "他生病了", "他找到了新工作", "他的农具坏了"),
			rq("故事的结果是什么？", "兔子没再出现，田地荒芜了", "兔子没再出现，田地荒芜了", "他抓到了很多兔子", "他成了富人", "他搬到了城里"),
			rq("作者认为这个寓言批评的是什么？", "把偶然当必然，不肯努力", "把偶然当必然，不肯努力", "过于勤劳", "喜欢动物", "太相信别人"),
			rq("“脚踏实地”在文中的意思是：", "踏实做事、认真努力", "踏实做事、认真努力", "走路要小心", "待在原地", "去农村生活"),
		},
	)
}
