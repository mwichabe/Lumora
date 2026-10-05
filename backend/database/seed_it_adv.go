package database

import "gorm.io/gorm"

// Italian course, B1 → C2. See seed_it.go.

// ───────────────────── B1 — Racconti e Opinioni: Stories & Opinions ─────────────────────

func seedItalianB1(db *gorm.DB) {
	const it = "it"
	const u = "B1 · Racconti e Opinioni — Stories & Opinions"
	finch := "Professor Finch"

	s := addSkillL(db, it, u, "The Future Tense", "parlerò, sarò, avrò, andrò — and the future of probability.", "Compass", "#6C3FC5", 20, 300)
	l := addLesson(db, s, "Domani…", 1, 22,
		char(finch, "Futuro semplice: -are and -ere → -erò (parlerò, scriverò), -ire → -irò (dormirò). Irregular stems: essere → sarò, avere → avrò, andare → andrò, fare → farò, venire → verrò, potere → potrò, volere → vorrò. Italians also use it to guess: Sarà a casa (he's probably at home), Avrà trent'anni (she must be about thirty)."),
		mc("Next year I will go to Sicily.", "andare", "L'anno prossimo andrò in Sicilia.", "L'anno prossimo andrò in Sicilia.", "L'anno prossimo anderò in Sicilia.", "L'anno prossimo vado in Sicilia ieri.", "L'anno prossimo sono andato in Sicilia."),
		mc("Future of essere (noi)", "essere", "saremo", "saremo", "esseremo", "siamo", "saremmo"),
		mc("What does 'Sarà stanco' suggest?", "Future of probability", "He's probably tired", "He's probably tired", "He will be tired next year", "He was tired", "He would be tired if…"),
		mc("They will come tomorrow.", "venire", "Verranno domani.", "Verranno domani.", "Veniranno domani.", "Vengono ieri.", "Verrebbero domani."),
		tr("Translate", "I will call you tomorrow (tu)", "Ti chiamerò domani"),
		fill("Future", "Quando ___ (avere, io) tempo, ti scriverò.", "avrò"),
		speak("Zephyr", "Quest'estate farò un viaggio in Puglia e starò al mare."),
	)
	addVocab(db, l,
		vw("sarò / avrò", "I will be / I will have", "Sarò a casa alle otto.", "I'll be home at eight.", finch),
		vw("andrò / farò", "I will go / I will do", "Domani farò la spesa.", "Tomorrow I'll do the shopping.", "Cora"),
		vw("l'anno prossimo", "next year", "L'anno prossimo studierò a Bologna.", "Next year I'll study in Bologna.", "Zephyr"),
		vw("fra / tra", "in (time from now)", "Arrivo fra dieci minuti.", "I'll arrive in ten minutes.", "Riko"),
	)

	s = addSkillL(db, it, u, "The Conditional", "vorrei, potrei, dovresti — wishes, politeness and advice.", "Scale", "#00C2A8", 21, 320)
	l = addLesson(db, s, "I Would…", 1, 22,
		char(finch, "Condizionale presente: the future stem + -ei, -esti, -ebbe, -emmo, -este, -ebbero: parlerei, sarei, avrei, andrei. It softens requests (Potrebbe aiutarmi? — could you help me?), gives advice (Dovresti riposare — you should rest) and expresses wishes (Mi piacerebbe vivere a Roma). Careful: saremo (we will be) vs saremmo (we would be) — one m vs two!"),
		mc("You should study more. (tu)", "Advice", "Dovresti studiare di più.", "Dovresti studiare di più.", "Devi studiare di più ieri.", "Dovrai studiare di più.", "Dovevi studiare di più."),
		mc("Could you help me? (formal)", "Polite request", "Potrebbe aiutarmi?", "Potrebbe aiutarmi?", "Può aiutare me subito?", "Potrà aiutarmi?", "Aiutami, Lei!"),
		mc("We would be happy. (one m or two?)", "essere", "Saremmo felici.", "Saremmo felici.", "Saremo felici.", "Siamo felici.", "Sareste felici."),
		mc("I'd love to live in Rome.", "piacere", "Mi piacerebbe vivere a Roma.", "Mi piacerebbe vivere a Roma.", "Mi piace vivevo a Roma.", "Mi piacerà vivere a Roma ieri.", "Piacerei vivere a Roma."),
		tr("Translate", "I would like a room with a view", "Vorrei una camera con vista"),
		fill("Conditional", "Al tuo posto, io non lo ___. (fare)", "farei"),
		speak("Mira", "Sarebbe bello andare in Toscana insieme, non pensi?"),
	)
	addVocab(db, l,
		vw("vorrei / potrei", "I would like / I could", "Potrei avere il menù?", "Could I have the menu?", "Cora"),
		vw("dovresti", "you should", "Dovresti chiamarla.", "You should call her.", "Mira"),
		vw("mi piacerebbe", "I would like / I'd love", "Mi piacerebbe imparare a cucinare.", "I'd love to learn to cook.", "Cora"),
		vw("al tuo posto", "if I were you", "Al tuo posto, accetterei.", "If I were you, I'd accept.", finch),
	)

	s = addSkillL(db, it, u, "The Present Subjunctive", "penso che sia, voglio che tu venga — the congiuntivo begins.", "Brain", "#F5A623", 22, 340)
	l = addLesson(db, s, "Forming the Congiuntivo", 1, 24,
		char(finch, "The congiuntivo expresses opinion, doubt, wishes and emotion after che. Regular endings: -are → -i (che io parli, che tu parli, che lui parli, che parliamo, che parliate, che parlino); -ere/-ire → -a (che scriva, che dorma). Key irregulars: essere → sia, avere → abbia, andare → vada, fare → faccia, venire → venga, potere → possa, volere → voglia, dovere → debba."),
		mc("I think he's right.", "pensare che + congiuntivo", "Penso che abbia ragione.", "Penso che abbia ragione.", "Penso che ha ragione.", "Penso che avrà ragione ieri.", "Penso che avere ragione."),
		mc("Congiuntivo of essere (lui)", "che lui ___", "sia", "sia", "è", "fosse (only)", "essa"),
		mc("I want you to come. (tu)", "volere che", "Voglio che tu venga.", "Voglio che tu venga.", "Voglio che tu vieni.", "Voglio tu venire.", "Voglio che tu verrai."),
		mc("Congiuntivo of andare (loro)", "che loro ___", "vadano", "vadano", "vanno", "andino", "vadino"),
		tr("Translate", "I hope that you are well (tu)", "Spero che tu stia bene"),
		fill("Congiuntivo", "Credo che lei ___ (essere) italiana.", "sia"),
		speak("Lumora", "Spero che domani faccia bel tempo."),
	)
	addVocab(db, l,
		vw("che io sia / abbia", "that I be / have (subjunctive)", "Non credo che sia vero.", "I don't think it's true.", finch),
		vw("pensare che", "to think that (+ subjunctive)", "Penso che sia tardi.", "I think it's late.", "Mira"),
		vw("sperare che", "to hope that", "Spero che tu venga.", "I hope you'll come.", "Lumora"),
		vw("credere che", "to believe that", "Credo che abbiano ragione.", "I believe they're right.", "Zephyr"),
	)
	l = addLesson(db, s, "When to Use It", 2, 24,
		char(finch, "Use the congiuntivo after verbs of opinion (penso, credo, mi sembra che), will and wishes (voglio, preferisco che), emotion (sono contento che, ho paura che), doubt (dubito che) and impersonal expressions (è importante / bisogna / è possibile che). NOT after certainty: So che è vero; È chiaro che ha ragione. If the subject is the same, use di + infinitive: Penso di partire (I think I'll leave)."),
		mc("I know that he's right. (certainty)", "sapere che", "So che ha ragione.", "So che ha ragione.", "So che abbia ragione.", "So che avere ragione.", "So che sia ragione."),
		mc("It's important that you study.", "Impersonal", "È importante che tu studi.", "È importante che tu studi.", "È importante che tu studia.", "È importante tu studiare.", "È importante che tu studierai."),
		mc("I'm happy you're here.", "Emotion", "Sono contento che tu sia qui.", "Sono contento che tu sia qui.", "Sono contento che tu sei qui.", "Sono contento di tu essere qui.", "Sono contento che tu fossi qui ora."),
		mc("I think I'll leave tomorrow. (same subject)", "penso di + infinitive", "Penso di partire domani.", "Penso di partire domani.", "Penso che io parta domani.", "Penso partire domani.", "Penso che parto domani."),
		tr("Translate", "I doubt that it is true", "Dubito che sia vero"),
		fill("Same subject", "Spero ___ vederti presto.", "di"),
		speak("Mira", "Bisogna che tutti facciano la loro parte."),
	)
	addVocab(db, l,
		vw("bisogna che", "it's necessary that", "Bisogna che tu parta subito.", "You need to leave right away.", finch),
		vw("è importante che", "it's important that", "È importante che tu lo sappia.", "It's important that you know.", "Mira"),
		vw("dubitare", "to doubt", "Dubito che venga.", "I doubt he'll come.", "Zephyr"),
		vw("avere paura che", "to be afraid that", "Ho paura che piova.", "I'm afraid it'll rain.", "Nana"),
	)

	s = addSkillL(db, it, u, "Relative Pronouns", "che, cui, il quale, chi.", "Link", "#FF5C5C", 23, 360)
	l = addLesson(db, s, "Che & Cui", 1, 22,
		char(finch, "che is the all-purpose 'who / which / that' for subjects and direct objects: il ragazzo che parla, il libro che leggo. After a preposition use cui: la città in cui vivo, l'amico con cui esco, il motivo per cui. il quale / la quale replaces che or cui in formal writing. chi means 'he/she who': Chi dorme non piglia pesci."),
		mc("The girl who is singing is my sister.", "Subject", "La ragazza che canta è mia sorella.", "La ragazza che canta è mia sorella.", "La ragazza cui canta è mia sorella.", "La ragazza chi canta è mia sorella.", "La ragazza quale canta è mia sorella."),
		mc("The city I live in", "After a preposition", "La città in cui vivo", "La città in cui vivo", "La città in che vivo", "La città che vivo in", "La città chi vivo"),
		mc("The friend I go out with", "con + relative", "L'amico con cui esco", "L'amico con cui esco", "L'amico con che esco", "L'amico che esco con", "L'amico con chi esco"),
		mc("That's the reason why I left.", "per cui", "È il motivo per cui sono partito.", "È il motivo per cui sono partito.", "È il motivo per che sono partito.", "È il motivo chi sono partito.", "È il motivo cui sono partito per."),
		tr("Translate", "the book that I am reading", "il libro che sto leggendo"),
		fill("Relative", "La persona a ___ ho scritto non ha risposto.", "cui"),
		speak("Zephyr", "Questo è il paese in cui sono nato."),
	)
	addVocab(db, l,
		vw("che", "who, which, that", "Il film che ho visto era bello.", "The film I saw was good.", finch),
		vw("cui", "which / whom (after prepositions)", "La ragione per cui studio.", "The reason why I study.", finch),
		vw("il quale / la quale", "who, which (formal)", "Il direttore, il quale ha parlato…", "The director, who spoke…", "Mira"),
		vw("chi", "he/she who, whoever", "Chi cerca trova.", "Whoever seeks finds.", "Nana"),
	)

	s = addSkillL(db, it, u, "Combined Pronouns, Ci & Ne", "glielo, me lo, ci vado, ne ho due.", "Layers", "#17A3DD", 24, 380)
	l = addLesson(db, s, "Glielo Do", 1, 22,
		char(finch, "When an indirect and a direct pronoun combine, mi/ti/ci/vi become me/te/ce/ve: Me lo dai? (will you give it to me?), Te la presto. gli and le both become glie- joined to the next pronoun: Glielo do (I give it to him/her). ci replaces a place (Ci vado domani — I'm going there tomorrow) and appears in ci vuole (it takes). ne replaces 'of it/them' or 'from there': Ne parliamo dopo."),
		mc("Will you give it to me? (the book)", "mi + lo", "Me lo dai?", "Me lo dai?", "Mi lo dai?", "Lo mi dai?", "Me la dai? (book)"),
		mc("I give it to her. (the letter, la lettera)", "le + la", "Gliela do.", "Gliela do.", "Le la do.", "Gliela da.", "La le do."),
		mc("Are you going to Rome? — Yes, I'm going there tomorrow.", "ci = there", "Sì, ci vado domani.", "Sì, ci vado domani.", "Sì, ne vado domani.", "Sì, lo vado domani.", "Sì, vado ci domani."),
		mc("It takes two hours.", "volerci", "Ci vogliono due ore.", "Ci vogliono due ore.", "Ci vuole due ore.", "Vogliono due ore ci.", "Ne vogliono due ore."),
		tr("Translate", "Let's talk about it later", "Ne parliamo dopo"),
		fill("Combined pronoun", "Il libro? ___ lo presto volentieri. (to you, tu)", "Te"),
		speak("Mira", "Le chiavi? Te le porto stasera."),
	)
	addVocab(db, l,
		vw("me lo / te lo", "it to me / it to you", "Me lo spieghi?", "Will you explain it to me?", finch),
		vw("glielo / gliela", "it to him/her", "Glielo dico io.", "I'll tell him/her.", "Mira"),
		vw("ci vuole / ci vogliono", "it takes (time, effort)", "Ci vuole pazienza.", "It takes patience.", "Nana"),
		vw("prestare", "to lend", "Mi presti la penna?", "Can you lend me the pen?", "Pip"),
	)

	s = addSkillL(db, it, u, "Health & the Body", "Dal medico: mi fa male, la ricetta, la farmacia.", "HeartPulse", "#9B51E0", 25, 400)
	l = addLesson(db, s, "Dal Medico", 1, 22,
		char("Nana", "Symptoms: Mi fa male la testa (my head hurts — literally 'the head hurts me'), Mi fanno male i piedi (plural!), Ho la febbre (I have a temperature), Ho il raffreddore (a cold), Ho la tosse (a cough). The doctor writes la ricetta (the prescription); you buy medicine in farmacia. Wish someone well: Guarisci presto! / Buona guarigione!"),
		mc("My head hurts.", "fare male", "Mi fa male la testa.", "Mi fa male la testa.", "Mi fanno male la testa.", "Ho male testa io.", "La testa fa male a me ieri."),
		mc("My feet hurt.", "Plural subject", "Mi fanno male i piedi.", "Mi fanno male i piedi.", "Mi fa male i piedi.", "Ho male piedi.", "I piedi mi fa male."),
		mc("What is 'la ricetta' at the doctor's?", "Il medico scrive la ricetta.", "The prescription", "The prescription", "The recipe only", "The bill", "The appointment"),
		mc("I have a temperature.", "Symptoms", "Ho la febbre.", "Ho la febbre.", "Sono febbre.", "Ho febbre la.", "Mi fa febbre."),
		tr("Translate", "I have a cold", "Ho il raffreddore"),
		tr("Translate", "Get well soon (tu)", "Guarisci presto"),
		speak("Nana", "Prenda questa medicina due volte al giorno e si riposi."),
	)
	addVocab(db, l,
		vw("mi fa male…", "… hurts (me)", "Mi fa male la schiena.", "My back hurts.", "Nana"),
		vw("la febbre", "fever, temperature", "Il bambino ha la febbre.", "The child has a temperature.", "Nana"),
		vw("la ricetta", "prescription; recipe", "Ha la ricetta del medico?", "Do you have the doctor's prescription?", "Nana"),
		vw("la farmacia", "pharmacy", "La farmacia è aperta fino alle otto.", "The pharmacy is open until eight.", "Riko"),
		vw("guarire", "to recover, heal", "Guarisci presto!", "Get well soon!", "Lumora"),
	)

	s = addSkillL(db, it, u, "Opinions & Work", "Secondo me, sono d'accordo — and the job interview.", "Briefcase", "#00C2A8", 26, 420)
	l = addLesson(db, s, "Al Lavoro", 1, 22,
		char("Mira", "Opinions: Secondo me… / A mio parere… (in my opinion), Sono d'accordo / Non sono d'accordo, Hai ragione, però… At work and in interviews use Lei: Mi parli della sua esperienza (tell me about your experience). Useful words: il colloquio (interview), il curriculum, lo stipendio (salary), le ferie (holidays), il contratto."),
		mc("In my opinion…", "Opinion", "Secondo me…", "Secondo me…", "Secondo io…", "Per mio…", "A me parere…"),
		mc("I don't agree.", "Disagreeing", "Non sono d'accordo.", "Non sono d'accordo.", "Non ho d'accordo.", "Non accordo.", "Sono non d'accordo."),
		mc("What is 'lo stipendio'?", "Lo stipendio è basso.", "Salary", "Salary", "Interview", "Holiday", "Office"),
		mc("Interviewer (formal): Tell me about yourself.", "Colloquio", "Mi parli di Lei.", "Mi parli di Lei.", "Parlami di te, dai!", "Parla di Lei, tu.", "Mi parla di te."),
		tr("Translate", "I have a job interview tomorrow", "Domani ho un colloquio di lavoro"),
		write("Write a short paragraph in Italian (about 60 words) giving your opinion on working from home. Use 'secondo me', 'però' and at least one 'penso che' + congiuntivo.",
			"Secondo me, lavorare da casa ha molti vantaggi. Non bisogna perdere tempo nel traffico e si può organizzare meglio la giornata. Però penso che sia difficile separare il lavoro dalla vita privata, e a volte ci si sente soli. Per questo credo che la soluzione migliore sia alternare: due o tre giorni a casa e il resto in ufficio."),
	)
	addVocab(db, l,
		vw("secondo me", "in my opinion", "Secondo me hai ragione.", "In my opinion you're right.", "Mira"),
		vw("essere d'accordo", "to agree", "Siamo tutti d'accordo.", "We all agree.", finch),
		vw("il colloquio", "interview", "Com'è andato il colloquio?", "How did the interview go?", "Mira"),
		vw("lo stipendio", "salary", "Lo stipendio arriva a fine mese.", "Salary comes at the end of the month.", "Blaze"),
		vw("le ferie", "holidays (from work)", "Ad agosto vado in ferie.", "In August I go on holiday.", "Zephyr"),
	)
}

// ───────────────────── B2 — L'Italiano in Contesto: Italian in Context ─────────────────────

func seedItalianB2(db *gorm.DB) {
	const it = "it"
	const u = "B2 · L'Italiano in Contesto — Italian in Context"
	finch := "Professor Finch"

	s := addSkillL(db, it, u, "Imperfect & Pluperfect Subjunctive", "se fossi, vorrei che tu venissi, se avessi saputo.", "Brain", "#6C3FC5", 27, 460)
	l := addLesson(db, s, "Fossi, Avessi, Venissi", 1, 24,
		char(finch, "Congiuntivo imperfetto: -are → -assi (parlassi), -ere → -essi (scrivessi), -ire → -issi (dormissi): che io parlassi, tu parlassi, lui parlasse, noi parlassimo, voi parlaste, loro parlassero. essere → fossi, fare → facessi, dare → dessi, stare → stessi. Use it after a past or conditional main verb: Volevo che tu venissi; Vorrei che fosse vero. Trapassato: avessi/fossi + participle — Pensavo che fosse già partito."),
		mc("I wish it were true.", "vorrei che", "Vorrei che fosse vero.", "Vorrei che fosse vero.", "Vorrei che sia vero.", "Vorrei che è vero.", "Vorrei che sarebbe vero."),
		mc("I wanted you to come. (tu)", "Past main verb", "Volevo che tu venissi.", "Volevo che tu venissi.", "Volevo che tu venga.", "Volevo che tu vieni.", "Volevo che tu verresti."),
		mc("Congiuntivo imperfetto of essere (io)", "che io ___", "fossi", "fossi", "sia", "ero", "sarei"),
		mc("I thought he had already left.", "Trapassato", "Pensavo che fosse già partito.", "Pensavo che fosse già partito.", "Pensavo che sia già partito.", "Pensavo che era già partito ieri.", "Pensavo che sarebbe già partito ora."),
		tr("Translate", "It seemed that they were tired", "Sembrava che fossero stanchi"),
		fill("Congiuntivo imperfetto", "Se io ___ (avere) più tempo…", "avessi"),
		speak("Mira", "Magari fosse sempre estate!"),
	)
	addVocab(db, l,
		vw("fossi / avessi", "(if) I were / had", "Se fossi ricco, viaggerei.", "If I were rich, I'd travel.", finch),
		vw("magari", "if only; maybe", "Magari potessi venire!", "If only I could come!", "Mira"),
		vw("sembrare che", "to seem that", "Sembra che piova.", "It seems to be raining.", "Zephyr"),
		vw("come se", "as if (+ congiuntivo)", "Parla come se sapesse tutto.", "He talks as if he knew everything.", finch),
	)

	s = addSkillL(db, it, u, "Hypothetical Sentences", "Il periodo ipotetico: real, possible and impossible.", "Scale", "#00C2A8", 28, 490)
	l = addLesson(db, s, "Se…", 1, 24,
		char(finch, "Three types. Real: se + present → present/future (Se piove, resto a casa). Possible: se + congiuntivo imperfetto → condizionale presente (Se avessi tempo, verrei). Impossible (past): se + congiuntivo trapassato → condizionale passato (Se avessi studiato, avrei superato l'esame). Never put the conditional after se!"),
		mc("If I had time, I would come.", "Type 2", "Se avessi tempo, verrei.", "Se avessi tempo, verrei.", "Se avrei tempo, verrei.", "Se ho tempo, verrei.", "Se avessi tempo, vengo ieri."),
		mc("If I had studied, I would have passed.", "Type 3", "Se avessi studiato, avrei superato l'esame.", "Se avessi studiato, avrei superato l'esame.", "Se avrei studiato, avessi superato l'esame.", "Se studiavo, superavo l'esame (only).", "Se ho studiato, supererei l'esame."),
		mc("If it rains, I'll stay at home.", "Type 1", "Se piove, resto a casa.", "Se piove, resto a casa.", "Se piovesse, resterò a casa.", "Se pioverebbe, resto a casa.", "Se avesse piovuto, resto a casa."),
		mc("Which is wrong?", "Spot the error", "Se sarei ricco, comprerei una villa.", "Se sarei ricco, comprerei una villa.", "Se fossi ricco, comprerei una villa.", "Se sono ricco, compro una villa.", "Se fossi stato ricco, avrei comprato una villa."),
		tr("Translate", "If you had told me, I would have helped you (tu)", "Se me l'avessi detto, ti avrei aiutato"),
		fill("Type 2", "Se ___ (essere, io) in te, accetterei.", "fossi"),
		speak("Zephyr", "Se avessi saputo che venivi, avrei preparato la cena."),
	)
	addVocab(db, l,
		vw("se fossi in te", "if I were you", "Se fossi in te, ci penserei.", "If I were you, I'd think about it.", finch),
		vw("avrei fatto", "I would have done", "Avrei fatto lo stesso.", "I'd have done the same.", "Mira"),
		vw("superare un esame", "to pass an exam", "Ha superato l'esame.", "She passed the exam.", "Pip"),
		vw("altrimenti", "otherwise", "Sbrigati, altrimenti perdiamo il treno.", "Hurry, otherwise we'll miss the train.", "Riko"),
	)

	s = addSkillL(db, it, u, "Passive & Impersonal si", "è stato costruito, viene scritto, va fatto, si mangia bene.", "Shield", "#F5A623", 29, 520)
	l = addLesson(db, s, "Being Done", 1, 24,
		char(finch, "Passive: essere + participle (agreeing) + da: Il Colosseo è stato costruito dai Romani. venire replaces essere in simple tenses for an ongoing action: La pasta viene servita calda. andare + participle means 'must be': Il modulo va compilato (the form must be filled in). The impersonal si: In Italia si mangia bene (one eats well); with a plural noun the verb is plural: Si vendono case."),
		mc("The Colosseum was built by the Romans.", "Passive", "Il Colosseo è stato costruito dai Romani.", "Il Colosseo è stato costruito dai Romani.", "Il Colosseo ha costruito i Romani.", "Il Colosseo è costruita dai Romani.", "Il Colosseo si è costruito i Romani."),
		mc("The form must be filled in.", "andare + participle", "Il modulo va compilato.", "Il modulo va compilato.", "Il modulo va a compilare.", "Il modulo compila.", "Il modulo è andato compilato."),
		mc("Houses for sale (impersonal si)", "Cartello", "Si vendono case.", "Si vendono case.", "Si vende case.", "Vendono si case.", "Case si vende."),
		mc("In Italy people eat well.", "Impersonal si", "In Italia si mangia bene.", "In Italia si mangia bene.", "In Italia ci mangiamo bene.", "In Italia si mangiano bene.", "In Italia mangia bene."),
		tr("Translate", "The letter was written by Maria", "La lettera è stata scritta da Maria"),
		fill("Passive with venire", "Il pane ___ fatto ogni mattina. (is made)", "viene"),
		speak("Zephyr", "Questo palazzo è stato progettato da un architetto famoso."),
	)
	addVocab(db, l,
		vw("è stato costruito", "was built", "Il ponte è stato costruito nel 1500.", "The bridge was built in 1500.", "Zephyr"),
		vw("va fatto", "must be done", "Il lavoro va finito oggi.", "The work must be finished today.", finch),
		vw("si dice che", "they say that", "Si dice che sia un genio.", "They say he's a genius.", "Mira"),
		vw("il modulo", "form (document)", "Compili il modulo, per favore.", "Fill in the form, please.", "Riko"),
	)

	s = addSkillL(db, it, u, "Formal Register & Letters", "Gentile, Egregio, Le scrivo per…, Distinti saluti.", "PenLine", "#FF5C5C", 30, 550)
	l = addLesson(db, s, "La Lettera Formale", 1, 24,
		char("Mira", "Formal letters and emails: Gentile Signora Rossi / Gentile Dottore (Dear…), Egregio Direttore (very formal, for men in authority). Le scrivo per… (I'm writing to you to…), Le sarei grato/a se potesse… (I'd be grateful if you could…), In attesa di un Suo cortese riscontro (awaiting your kind reply). Close: Cordiali saluti or Distinti saluti. Lei, Suo and La are capitalised in letters as a mark of respect."),
		mc("Formal opening to a woman", "Email", "Gentile Signora Rossi,", "Gentile Signora Rossi,", "Ciao Rossi,", "Cara Rossi bella,", "Ehi Signora,"),
		mc("I'm writing to you to…", "Purpose", "Le scrivo per…", "Le scrivo per…", "Ti scrivo per…", "Scrivo a Lei che…", "Mi scrivo per…"),
		mc("Formal closing", "End of a formal letter", "Distinti saluti", "Distinti saluti", "Baci e abbracci", "Ciao ciao", "A presto, tesoro"),
		mc("I'd be grateful if you could send me…", "Formal request", "Le sarei grato se potesse inviarmi…", "Le sarei grato se potesse inviarmi…", "Mandami subito…", "Voglio che mi manda…", "Ti sarei grato se puoi…"),
		tr("Translate", "Thank you for your attention", "La ringrazio per l'attenzione"),
		write("Write a formal email in Italian (about 80 words) to a language school asking about an intensive B2 course: dates, price and accommodation. Use Gentile…, Le scrivo per…, and Distinti saluti.",
			"Gentile Segreteria,\nLe scrivo per avere informazioni sul corso intensivo di livello B2. Vorrei sapere quando inizia il prossimo corso, quanto costa e se la scuola offre la possibilità di alloggiare presso una famiglia. Inoltre, Le sarei grata se potesse indicarmi il numero di ore settimanali. In attesa di un Suo cortese riscontro, porgo distinti saluti.\nAnna Kamau"),
	)
	addVocab(db, l,
		vw("Gentile…", "Dear… (formal)", "Gentile Professore,", "Dear Professor,", "Mira"),
		vw("Le scrivo per…", "I'm writing to you to…", "Le scrivo per un'informazione.", "I'm writing to ask for some information.", "Mira"),
		vw("Distinti saluti", "Yours faithfully", "Distinti saluti, Marco Bianchi.", "Yours faithfully, Marco Bianchi.", finch),
		vw("il riscontro", "reply, feedback", "Attendo un Suo riscontro.", "I await your reply.", finch),
	)

	s = addSkillL(db, it, u, "News & Media", "Il governo, le elezioni, il PIL — read the headlines.", "Globe", "#17A3DD", 31, 580)
	l = addLesson(db, s, "Al Telegiornale", 1, 24,
		char(finch, "News Italian (Corriere della Sera, la Repubblica, RAI Tg1) has its own vocabulary: il governo, il Parlamento, le elezioni, il presidente del Consiglio (prime minister), il PIL (GDP), la disoccupazione (unemployment), l'inflazione, lo sciopero (strike). Headlines drop articles and verbs: 'Sciopero dei treni, disagi in tutta Italia'."),
		mc("What does 'sciopero' mean?", "Sciopero dei treni", "strike", "strike", "storm", "train timetable", "delay"),
		mc("What is 'il PIL'?", "Il PIL cresce dello 0,5%", "GDP", "GDP", "a political party", "the police", "a tax"),
		mc("Who is 'il presidente del Consiglio'?", "Italian politics", "The prime minister", "The prime minister", "The head of a city council", "The President of the Republic", "The mayor"),
		mc("What does this headline say?", "Disoccupazione in calo al Sud", "Unemployment is falling in the South", "Unemployment is falling in the South", "Unemployment is rising in the North", "Jobs are moving south", "The South has no jobs"),
		tr("Translate", "The government has approved the budget", "Il governo ha approvato il bilancio"),
		mc("What does 'disagi' mean in 'disagi in tutta Italia'?", "Sciopero dei treni, disagi in tutta Italia", "disruption, inconvenience", "disruption, inconvenience", "celebrations", "accidents", "discounts"),
		speak("Zephyr", "Secondo l'ISTAT, la disoccupazione giovanile è scesa per il terzo mese consecutivo."),
	)
	addVocab(db, l,
		vw("il governo", "the government", "Il governo ha presentato la riforma.", "The government presented the reform.", finch),
		vw("le elezioni", "elections", "Le elezioni si terranno a giugno.", "The elections will be held in June.", finch),
		vw("lo sciopero", "strike", "Domani c'è sciopero dei mezzi.", "Tomorrow there's a public transport strike.", "Riko"),
		vw("la disoccupazione", "unemployment", "La disoccupazione è in calo.", "Unemployment is falling.", "Zephyr"),
	)

	s = addSkillL(db, it, u, "Connectors & Discourse", "tuttavia, pertanto, infatti, benché, affinché, nonostante.", "Quote", "#9B51E0", 32, 610)
	l = addLesson(db, s, "Linking Ideas", 1, 24,
		char(finch, "Connectors raise your level instantly. Contrast: tuttavia, però, invece. Cause/consequence: infatti (in fact — confirming), quindi, pertanto (therefore), dato che, siccome. Concession + congiuntivo: benché / sebbene / nonostante (although) — Benché sia tardi, continuo. Purpose + congiuntivo: affinché / perché (so that) — Te lo dico affinché tu lo sappia."),
		mc("Although it's late, I'll keep working.", "Concession", "Benché sia tardi, continuo a lavorare.", "Benché sia tardi, continuo a lavorare.", "Benché è tardi, continuo a lavorare.", "Perché sia tardi, continuo.", "Benché sarà tardi, continuo."),
		mc("I'm telling you so that you know.", "Purpose", "Te lo dico affinché tu lo sappia.", "Te lo dico affinché tu lo sappia.", "Te lo dico affinché tu lo sai.", "Te lo dico perché lo sapessi ieri.", "Te lo dico per tu sapere."),
		mc("Which connector means 'therefore' (formal)?", "Consequence", "pertanto", "pertanto", "tuttavia", "benché", "infatti"),
		mc("Which connector introduces a contrast?", "Contrast", "tuttavia", "tuttavia", "quindi", "infatti", "pertanto"),
		tr("Translate", "Despite the rain, we went out", "Nonostante la pioggia, siamo usciti"),
		fill("Concession", "Sebbene lui ___ (essere) stanco, è venuto.", "fosse"),
		speak("Mira", "Il progetto è ambizioso; tuttavia, crediamo che sia realizzabile."),
	)
	addVocab(db, l,
		vw("tuttavia", "however, nevertheless", "È caro; tuttavia, lo compro.", "It's expensive; nevertheless, I'll buy it.", finch),
		vw("pertanto", "therefore", "Il negozio è chiuso, pertanto torneremo domani.", "The shop is closed, so we'll come back tomorrow.", finch),
		vw("benché / sebbene", "although (+ congiuntivo)", "Benché piova, esco.", "Although it's raining, I'm going out.", "Mira"),
		vw("affinché", "so that (+ congiuntivo)", "Parlo piano affinché tutti capiscano.", "I speak slowly so everyone understands.", "Lumora"),
		vw("infatti", "indeed, in fact", "Era stanco; infatti è andato a letto presto.", "He was tired; indeed, he went to bed early.", "Zephyr"),
	)
}

// ───────────────────── C1 — Sfumature: Nuance & Style ─────────────────────

func seedItalianC1(db *gorm.DB) {
	const it = "it"
	const u = "C1 · Sfumature — Nuance & Style"
	finch := "Professor Finch"

	s := addSkillL(db, it, u, "Passato Remoto & Literary Italian", "fu, ebbe, nacque, disse — the tense of history and novels.", "Landmark", "#6C3FC5", 33, 700)
	l := addLesson(db, s, "C'era una volta…", 1, 26,
		char(finch, "The passato remoto narrates completed events in history and literature (and is still spoken in parts of the South and Tuscany). Regular: parlai, parlasti, parlò, parlammo, parlaste, parlarono. Key irregulars: essere → fui, fosti, fu, fummo, foste, furono; avere → ebbi, ebbe, ebbero; nascere → nacque; dire → disse; fare → fece; vedere → vide; venire → venne."),
		mc("Dante was born in Florence.", "nascere", "Dante nacque a Firenze.", "Dante nacque a Firenze.", "Dante nascette a Firenze.", "Dante nascé a Firenze.", "Dante nascerà a Firenze."),
		mc("Passato remoto of essere (lui)", "Garibaldi ___ un eroe.", "fu", "fu", "fosse", "era stato", "fù"),
		mc("He said nothing.", "dire", "Non disse nulla.", "Non disse nulla.", "Non dicette nulla.", "Non dirò nulla.", "Non dise nulla."),
		mc("Where is the passato remoto still used in speech?", "Uso regionale", "In parts of the South and Tuscany", "In parts of the South and Tuscany", "Only in Milan", "Nowhere at all", "Only in Switzerland"),
		mc("They saw the sea for the first time.", "vedere", "Videro il mare per la prima volta.", "Videro il mare per la prima volta.", "Vedettero il mare.", "Vedrono il mare.", "Videvano il mare."),
		tr("Translate", "The war ended in 1945", "La guerra finì nel 1945"),
		speak("Nana", "C'era una volta un re che ebbe tre figlie."),
	)
	addVocab(db, l,
		vw("fu", "was (passato remoto)", "Fu un giorno indimenticabile.", "It was an unforgettable day.", finch),
		vw("nacque", "was born", "Leonardo nacque nel 1452.", "Leonardo was born in 1452.", finch),
		vw("disse / fece", "said / did, made", "Disse di sì e lo fece.", "He said yes and did it.", "Nana"),
		vw("c'era una volta", "once upon a time", "C'era una volta un pezzo di legno.", "Once upon a time there was a piece of wood.", "Nana"),
	)

	s = addSkillL(db, it, u, "Idioms & Proverbs", "In bocca al lupo, non vedo l'ora, costa un occhio della testa.", "Quote", "#00C2A8", 34, 730)
	l = addLesson(db, s, "Modi di Dire", 1, 26,
		char("Nana", "Idioms are everywhere. In bocca al lupo! (good luck — 'into the wolf's mouth') — the reply is Crepi! (may it die). Non vedo l'ora (I can't wait), costa un occhio della testa (it costs an arm and a leg), essere al verde (to be broke), avere le mani bucate (to be a spendthrift), prendere in giro (to tease). Proverbs: Chi dorme non piglia pesci; Tra il dire e il fare c'è di mezzo il mare; Chi va piano va sano e va lontano."),
		mc("Someone says 'In bocca al lupo!' You reply:", "Before an exam", "Crepi!", "Crepi!", "Grazie mille, lupo!", "Anche a te!", "Salute!"),
		mc("What does 'essere al verde' mean?", "Sono al verde.", "to be broke", "to be broke", "to be inexperienced", "to be envious", "to be in the countryside"),
		mc("What does 'costa un occhio della testa' mean?", "Quella borsa costa un occhio della testa.", "It's extremely expensive", "It's extremely expensive", "It's ugly", "It's eye-catching", "It's on sale"),
		mc("'Tra il dire e il fare c'è di mezzo il mare' means…", "Proverbio", "Easier said than done", "Easier said than done", "The sea is far away", "Talk while you work", "Travel broadens the mind"),
		mc("'Non vedo l'ora!' means…", "Non vedo l'ora di vederti!", "I can't wait!", "I can't wait!", "I can't see the time", "I don't have time", "I'm late"),
		mc("Who 'ha le mani bucate'?", "Ha le mani bucate.", "Someone who spends money freely", "Someone who spends money freely", "Someone clumsy", "Someone injured", "Someone generous with time"),
		speak("Nana", "Chi va piano va sano e va lontano."),
	)
	addVocab(db, l,
		vw("In bocca al lupo! — Crepi!", "Good luck! — Thanks!", "Domani hai l'esame? In bocca al lupo!", "Exam tomorrow? Good luck!", "Nana"),
		vw("non vedo l'ora", "I can't wait", "Non vedo l'ora delle vacanze.", "I can't wait for the holidays.", "Pip"),
		vw("essere al verde", "to be broke", "A fine mese sono sempre al verde.", "At the end of the month I'm always broke.", "Blaze"),
		vw("prendere in giro", "to tease, make fun of", "Non prendermi in giro!", "Don't make fun of me!", "Pip"),
		vw("Chi dorme non piglia pesci", "the early bird catches the worm", "Svegliati! Chi dorme non piglia pesci.", "Wake up! You snooze, you lose.", "Nana"),
	)

	s = addSkillL(db, it, u, "Italian Literature", "Dante, Petrarca, Boccaccio, Manzoni, Calvino.", "BookOpen", "#F5A623", 35, 760)
	l = addLesson(db, s, "Le Tre Corone e Oltre", 1, 26,
		char(finch, "Modern Italian grew out of 14th-century Florentine, thanks to the 'tre corone': Dante (la Divina Commedia — 'Nel mezzo del cammin di nostra vita'), Petrarca (the sonnets of the Canzoniere) and Boccaccio (the hundred tales of the Decameron). Manzoni's I promessi sposi (1840s) fixed the modern language — he 'rinsed his clothes in the Arno'. The 20th century brought Pirandello, Calvino, Morante, Eco and Ferrante."),
		mc("Who wrote the Divina Commedia?", "'Nel mezzo del cammin di nostra vita…'", "Dante Alighieri", "Dante Alighieri", "Petrarca", "Boccaccio", "Manzoni"),
		mc("Which work is a collection of a hundred tales?", "Dieci giovani, dieci giorni", "Il Decameron", "Il Decameron", "Il Canzoniere", "I promessi sposi", "Le città invisibili"),
		mc("Why is Manzoni important for the language?", "I promessi sposi", "He shaped modern Italian on the Florentine model", "He shaped modern Italian on the Florentine model", "He invented the alphabet", "He wrote only in Latin", "He translated Shakespeare"),
		mc("What does 'risciacquare i panni in Arno' refer to?", "Manzoni", "Revising his novel into Florentine Italian", "Revising his novel into Florentine Italian", "Doing laundry in Florence", "A famous painting", "A flood"),
		mc("Who wrote 'Le città invisibili'?", "1972", "Italo Calvino", "Italo Calvino", "Umberto Eco", "Pirandello", "Dante"),
		mc("Which dialect became the basis of standard Italian?", "Storia della lingua", "Florentine (Tuscan)", "Florentine (Tuscan)", "Roman", "Neapolitan", "Venetian"),
		speak("Zephyr", "Nel mezzo del cammin di nostra vita mi ritrovai per una selva oscura."),
	)
	addVocab(db, l,
		vw("la Divina Commedia", "the Divine Comedy", "Ho letto l'Inferno di Dante.", "I've read Dante's Inferno.", finch),
		vw("il romanzo", "novel", "Un romanzo storico.", "A historical novel.", "Zephyr"),
		vw("il sonetto", "sonnet", "Petrarca scrisse molti sonetti.", "Petrarch wrote many sonnets.", finch),
		vw("lo scrittore / la scrittrice", "writer (m / f)", "Elena Ferrante è una scrittrice famosa.", "Elena Ferrante is a famous writer.", "Mira"),
	)

	s = addSkillL(db, it, u, "Argumentative Writing", "Il tema argomentativo: è innegabile che, d'altro canto, in conclusione.", "PenLine", "#FF5C5C", 36, 790)
	l = addLesson(db, s, "Il Tema Argomentativo", 1, 26,
		char(finch, "An Italian argumentative essay states a tesi, supports it with argomenti, addresses the antitesi and concludes. Formal moves: È innegabile che… (it's undeniable that), Va sottolineato che… (it should be stressed that), D'altro canto… (on the other hand), Ciononostante… (nonetheless), Alla luce di quanto detto… / In conclusione… — many trigger the congiuntivo."),
		mc("On the other hand…", "Antitesi", "D'altro canto…", "D'altro canto…", "Per esempio…", "In primo luogo…", "Insomma…"),
		mc("It's undeniable that technology has changed our lives.", "+ indicativo (certainty)", "È innegabile che la tecnologia ha cambiato la nostra vita.", "È innegabile che la tecnologia ha cambiato la nostra vita.", "È innegabile la tecnologia cambiare.", "Innegabile tecnologia cambiò.", "È innegabile per la tecnologia."),
		mc("In light of what has been said…", "Conclusion", "Alla luce di quanto detto…", "Alla luce di quanto detto…", "Alla luce del sole…", "Per quanto detto prima di ieri…", "Detto fatto…"),
		mc("What is the 'antitesi' in a tema?", "Struttura", "The opposing view", "The opposing view", "The title", "The thesis", "The bibliography"),
		mc("Nonetheless…", "Concession", "Ciononostante…", "Ciononostante…", "Quindi…", "Infatti…", "Cioè…"),
		write("Write an argumentative paragraph in Italian (about 120 words) on 'I social network fanno più male che bene ai giovani'. State your thesis, give two arguments, address the opposing view with 'd'altro canto', and conclude.",
			"Ritengo che i social network, se usati senza consapevolezza, facciano più male che bene ai giovani. In primo luogo, favoriscono il confronto continuo con modelli irraggiungibili, con conseguenze sull'autostima. In secondo luogo, riducono il tempo dedicato allo studio e alle relazioni reali. D'altro canto, è innegabile che offrano opportunità di informazione e di espressione creativa. Ciononostante, alla luce di quanto detto, credo che la scuola e la famiglia debbano educare a un uso critico di questi strumenti, affinché diventino una risorsa e non una dipendenza."),
	)
	addVocab(db, l,
		vw("la tesi / l'antitesi", "thesis / opposing view", "Esponi la tua tesi.", "State your thesis.", finch),
		vw("è innegabile che", "it's undeniable that", "È innegabile che il clima stia cambiando.", "It's undeniable that the climate is changing.", finch),
		vw("d'altro canto", "on the other hand", "D'altro canto, ci sono dei rischi.", "On the other hand, there are risks.", "Mira"),
		vw("ciononostante", "nonetheless", "Ciononostante, ha accettato.", "Nonetheless, she accepted.", "Zephyr"),
		vw("alla luce di", "in light of", "Alla luce dei fatti, cambio idea.", "In light of the facts, I'm changing my mind.", finch),
	)

	s = addSkillL(db, it, u, "Gestures & Pragmatics", "Gesti, la bella figura, dare del tu.", "Hand", "#17A3DD", 37, 820)
	l = addLesson(db, s, "Parlare con le Mani", 1, 26,
		char("Blaze", "Italians speak with their hands. Fingertips pinched together and shaken: Ma che vuoi? (what do you want?!). A finger twisted into the cheek: buono! (delicious). The chin flick outward: non me ne importa (I don't care). Pragmatics: fare bella figura (make a good impression) matters socially; switching from Lei to tu is offered by the older or senior person: Possiamo darci del tu?"),
		mc("Fingertips pinched together, hand shaken, means…", "Gesto", "Ma che vuoi? (What do you want?!)", "Ma che vuoi? (What do you want?!)", "Delicious", "Goodbye", "Come here"),
		mc("A finger twisted into the cheek means…", "Gesto", "Delicious!", "Delicious!", "You're crazy", "I'm thinking", "Quiet!"),
		mc("What does 'fare bella figura' mean?", "Pragmatica", "To make a good impression", "To make a good impression", "To draw a nice picture", "To look in the mirror", "To be very thin"),
		mc("Who usually proposes switching from Lei to tu?", "Possiamo darci del tu?", "The older or more senior person", "The older or more senior person", "The younger person, always", "Nobody — it's forbidden", "Only strangers"),
		mc("What does 'non me ne importa' mean?", "Gesto del mento", "I don't care", "I don't care", "I'm very important", "I'll import it", "I don't understand"),
		speak("Blaze", "Ma dai, possiamo darci del tu!"),
	)
	addVocab(db, l,
		vw("fare bella / brutta figura", "make a good / bad impression", "Non voglio fare brutta figura.", "I don't want to make a bad impression.", "Mira"),
		vw("darsi del tu", "to switch to 'tu'", "Ci diamo del tu?", "Shall we use 'tu'?", "Blaze"),
		vw("il gesto", "gesture", "Un gesto molto italiano.", "A very Italian gesture.", "Blaze"),
		vw("non me ne importa", "I don't care", "Non me ne importa niente.", "I couldn't care less.", "Pip"),
	)
}

// ───────────────────── C2 — Padronanza: Mastery ─────────────────────

func seedItalianC2(db *gorm.DB) {
	const it = "it"
	const u = "C2 · Padronanza — Mastery"
	finch := "Professor Finch"

	s := addSkillL(db, it, u, "Regional Varieties & Dialects", "Toscano, romanesco, napoletano, milanese — and regional Italian.", "Globe", "#6C3FC5", 38, 1000)
	l := addLesson(db, s, "L'Italia dei Dialetti", 1, 28,
		char(finch, "Italy's dialects are distinct languages descended from Latin, alongside regional varieties of Italian. Tuscan 'gorgia' softens c between vowels (la hasa for la casa). Romanesco: 'daje!' (come on!), 'mo'' (now). Neapolitan: 'guagliò' (kid), 'jamme' (let's go). Milanese/northern: 'il' before names (la Giulia), 'ciapa' (take). In the North the passato prossimo dominates; in the South, the passato remoto."),
		mc("What is the Tuscan 'gorgia'?", "la hasa", "Softening c between vowels to an h-like sound", "Softening c between vowels to an h-like sound", "A Tuscan dish", "A rolled r", "Dropping final vowels"),
		mc("'Daje!' is typical of…", "Dialetto", "Rome", "Rome", "Milan", "Venice", "Turin"),
		mc("What does Neapolitan 'guagliò' mean?", "Uè, guagliò!", "kid, young man", "kid, young man", "money", "food", "friend's mother"),
		mc("Using the article before first names (la Giulia) is typical of…", "Uso regionale", "Northern Italy", "Northern Italy", "Sicily", "Rome", "Sardinia only"),
		mc("Where is the passato remoto more common in speech?", "Uso regionale", "The South", "The South", "The North", "Rome only", "Nowhere"),
		mc("What are Italian dialects, linguistically?", "Storia", "Separate languages descended from Latin", "Separate languages descended from Latin", "Corrupted forms of standard Italian", "Slang", "Foreign languages"),
		speak("Blaze", "Daje, mo' annamo!"),
	)
	addVocab(db, l,
		vw("il dialetto", "dialect", "Mio nonno parla in dialetto.", "My grandfather speaks dialect.", "Nana"),
		vw("daje (romanesco)", "come on!", "Daje, ce la fai!", "Come on, you can do it!", "Blaze"),
		vw("mo' (centro-sud)", "now", "Arrivo mo'.", "I'm coming now.", "Blaze"),
		vw("l'italiano regionale", "regional Italian", "Ogni regione ha il suo italiano.", "Every region has its own Italian.", finch),
	)

	s = addSkillL(db, it, u, "Bureaucratic & Formal Style", "Il burocratese: nominalisation, impersonal forms and how to simplify.", "Briefcase", "#00C2A8", 39, 1030)
	l = addLesson(db, s, "Il Burocratese", 1, 28,
		char(finch, "Official Italian (burocratese) loves nouns and impersonal forms: 'Si fa presente che…' (please note that), 'Si prega di…' (you are requested to), 'In data odierna' (today), 'A seguito della Sua richiesta' (following your request), 'Procedere alla verifica' instead of 'verificare'. Mastery means reading it — and rewriting it plainly when you can."),
		mc("Plain Italian for 'in data odierna'", "Burocratese", "oggi", "oggi", "domani", "ieri", "in futuro"),
		mc("Plain Italian for 'procedere alla verifica'", "Burocratese", "verificare", "verificare", "procedere", "andare via", "verificato"),
		mc("What does 'Si prega di non fumare' mean?", "Cartello", "Please do not smoke", "Please do not smoke", "Smoking is encouraged", "Pray before smoking", "Smoke only here"),
		mc("What does 'A seguito della Sua richiesta' introduce?", "Lettera ufficiale", "A reply following your request", "A reply following your request", "A complaint", "An invitation", "A bill"),
		tr("Rewrite plainly", "Si fa presente che l'ufficio resterà chiuso → (plain Italian: Please note the office will stay closed)", "L'ufficio resterà chiuso"),
		mc("Which is the plain version of 'effettuare il pagamento'?", "Semplificare", "pagare", "pagare", "effettuare", "pagamento", "fare pagato"),
		speak("Zephyr", "Si comunica che, a decorrere da lunedì, gli uffici osserveranno il nuovo orario."),
	)
	addVocab(db, l,
		vw("il burocratese", "bureaucratic jargon", "Questo testo è scritto in burocratese.", "This text is written in bureaucratese.", finch),
		vw("si fa presente che", "please note that", "Si fa presente che il termine è scaduto.", "Please note that the deadline has passed.", finch),
		vw("a decorrere da", "with effect from", "A decorrere dal primo marzo.", "With effect from 1 March.", "Zephyr"),
		vw("in data odierna", "today (formal)", "Le scrivo in data odierna.", "I am writing to you today.", "Mira"),
	)

	s = addSkillL(db, it, u, "Translation & False Friends", "camera, libreria, morbido, eventualmente, attualmente.", "Languages", "#F5A623", 40, 1060)
	l = addLesson(db, s, "Falsi Amici", 1, 28,
		char(finch, "False friends trap even advanced learners: la camera = room (camera = macchina fotografica), la libreria = bookshop (library = biblioteca), morbido = soft, eventualmente = possibly / if necessary (not 'eventually' = alla fine), attualmente = currently (not 'actually' = in realtà), la fattoria = farm, il magazzino = warehouse, sensibile = sensitive (sensible = ragionevole)."),
		mc("What does 'attualmente' mean?", "Attualmente lavoro a Roma.", "currently", "currently", "actually", "eventually", "fortunately"),
		mc("Where do you borrow books?", "Not the libreria!", "in biblioteca", "in biblioteca", "in libreria", "in fattoria", "in magazzino"),
		mc("What does 'eventualmente' mean?", "Eventualmente ti chiamo.", "possibly / if necessary", "possibly / if necessary", "eventually, in the end", "immediately", "never"),
		mc("'He's a very sensible person' in Italian", "sensible ≠ sensibile", "È una persona molto ragionevole.", "È una persona molto ragionevole.", "È una persona molto sensibile.", "È una persona molto sensata in camera.", "È una persona eventuale."),
		mc("What is 'la camera'?", "Prenoto una camera.", "a room", "a room", "a camera", "a chamber orchestra", "a cupboard"),
		write("Translate into natural Italian (about 60 words): 'Currently I work in a bookshop in Bologna, but eventually I would like to open my own. My colleagues are very sensible and the atmosphere is relaxed. If necessary, I could also work at weekends.'",
			"Attualmente lavoro in una libreria a Bologna, ma alla fine vorrei aprirne una mia. I miei colleghi sono molto ragionevoli e l'atmosfera è rilassata. Eventualmente, potrei lavorare anche nei fine settimana."),
	)
	addVocab(db, l,
		vw("attualmente", "currently (not 'actually')", "Attualmente vivo a Torino.", "I currently live in Turin.", finch),
		vw("la libreria / la biblioteca", "bookshop / library", "Compro libri in libreria.", "I buy books at the bookshop.", "Mira"),
		vw("eventualmente", "possibly, if necessary", "Eventualmente ti richiamo.", "I'll call you back if necessary.", "Riko"),
		vw("ragionevole / sensibile", "sensible / sensitive", "Una scelta ragionevole.", "A sensible choice.", finch),
	)

	s = addSkillL(db, it, u, "Speeches & Ceremony", "Discorsi: matrimoni, premiazioni, commemorazioni.", "Mic", "#FF5C5C", 41, 1090)
	l = addLesson(db, s, "Il Discorso", 1, 30,
		char(finch, "Formal speeches open with Autorità, gentili ospiti, cari amici… (Authorities, distinguished guests, dear friends). Weddings: Un brindisi agli sposi! — Viva gli sposi! Thanks: Desidero ringraziare di cuore… Commemorations use the passato prossimo and remoto with care and a respectful register. Close: Vi ringrazio per l'attenzione."),
		mc("How does a formal speech often open?", "Apertura", "Autorità, gentili ospiti, cari amici…", "Autorità, gentili ospiti, cari amici…", "Ciao a tutti, ragazzi!", "Allora, insomma…", "Ehi, ascoltate!"),
		mc("A wedding toast", "Matrimonio", "Viva gli sposi!", "Viva gli sposi!", "In bocca al lupo!", "Buon viaggio!", "Condoglianze!"),
		mc("How do you close a speech?", "Chiusura", "Vi ringrazio per l'attenzione.", "Vi ringrazio per l'attenzione.", "Basta così, ciao.", "Fine.", "Andiamo a mangiare."),
		mc("I wish to thank wholeheartedly…", "Ringraziamenti", "Desidero ringraziare di cuore…", "Desidero ringraziare di cuore…", "Voglio grazie a cuore…", "Grazie da cuore voglio…", "Ringrazio cuore…"),
		mc("Condolences at a funeral", "Lutto", "Sentite condoglianze.", "Sentite condoglianze.", "Auguri!", "Complimenti!", "Viva!"),
		write("Write the opening of a formal speech in Italian (about 100 words) for a school prize-giving: greet the authorities and guests, thank the organisers, congratulate the winners and use at least one congiuntivo.",
			"Autorità, gentili ospiti, cari studenti, buonasera. È per me un onore prendere la parola in questa giornata di festa. Desidero innanzitutto ringraziare di cuore gli organizzatori, senza i quali questa cerimonia non sarebbe stata possibile. Rivolgo poi le mie più sincere congratulazioni ai vincitori: il vostro impegno dimostra che lo studio, quando è accompagnato dalla passione, porta lontano. Mi auguro che questo premio sia per voi non un traguardo, ma un punto di partenza. Vi ringrazio per l'attenzione."),
	)
	addVocab(db, l,
		vw("il discorso", "speech", "Ha fatto un bel discorso.", "She gave a fine speech.", finch),
		vw("gli sposi", "the bride and groom", "Un brindisi agli sposi!", "A toast to the newlyweds!", "Cora"),
		vw("ringraziare di cuore", "to thank wholeheartedly", "Vi ringrazio di cuore.", "I thank you wholeheartedly.", "Lumora"),
		vw("le condoglianze", "condolences", "Sentite condoglianze.", "My deepest condolences.", "Nana"),
		vw("prendere la parola", "to take the floor", "Ora prende la parola il sindaco.", "The mayor will now speak.", finch),
	)

	s = addSkillL(db, it, u, "Cinema, Song & Cultural References", "Fellini, Sanremo and the cantautori — the references natives expect you to catch.", "Music", "#17A3DD", 42, 1120)
	l = addLesson(db, s, "Riferimenti Culturali", 1, 30,
		char(finch, "Mastery includes the shared references Italians quote without explaining. 'La dolce vita' (Fellini, 1960) gave English the word paparazzo, from a photographer character. The Sanremo festival each February is a national event. Cantautori (singer-songwriters) like Fabrizio De André, Lucio Dalla and Franco Battiato are studied as poets. Neorealism (Rossellini, De Sica's Ladri di biciclette) shaped world cinema."),
		mc("Which English word comes from a character in Fellini's 'La dolce vita'?", "Fellini, 1960", "paparazzo", "paparazzo", "fiasco", "graffiti", "fresco"),
		mc("What is the Festival di Sanremo?", "Ogni febbraio", "A national song festival", "A national song festival", "A film festival in Venice", "A food fair", "A football tournament"),
		mc("What is a 'cantautore'?", "De André, Dalla, Battiato", "A singer-songwriter", "A singer-songwriter", "An opera singer", "A conductor", "A radio host"),
		mc("Which film is a classic of Neorealism?", "Il Neorealismo", "Ladri di biciclette", "Ladri di biciclette", "La vita è bella", "Il Gattopardo", "Nuovo Cinema Paradiso"),
		mc("Where is Italy's oldest international film festival held?", "Mostra del Cinema", "Venice", "Venice", "Rome", "Cannes", "Turin"),
		speak("Zephyr", "Il Festival di Sanremo è uno degli appuntamenti più seguiti dell'anno."),
	)
	addVocab(db, l,
		vw("il cantautore / la cantautrice", "singer-songwriter", "De André è un grande cantautore.", "De André is a great singer-songwriter.", "Mira"),
		vw("il Neorealismo", "Neorealism (film movement)", "Il Neorealismo nacque nel dopoguerra.", "Neorealism was born after the war.", finch),
		vw("la colonna sonora", "soundtrack", "La colonna sonora è di Morricone.", "The soundtrack is by Morricone.", "Zephyr"),
		vw("il paparazzo", "paparazzo (press photographer)", "I paparazzi aspettavano l'attrice.", "The paparazzi were waiting for the actress.", "Blaze"),
	)
}
