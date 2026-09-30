import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/languages.dart';
import '../../core/network/api_client.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../core/theme/shadows.dart';
import '../../providers/auth_provider.dart';
import '../../widgets/stat_widgets.dart';

class _Mode {
  final String key;
  final String title;
  final String desc;
  final IconData icon;
  final Color tint;
  final bool disabled;
  final int badge;
  const _Mode({
    required this.key, required this.title, required this.desc, required this.icon,
    required this.tint, required this.disabled, this.badge = 0,
  });
}

class PracticeScreen extends ConsumerStatefulWidget {
  const PracticeScreen({super.key});

  @override
  ConsumerState<PracticeScreen> createState() => _PracticeScreenState();
}

class _PracticeScreenState extends ConsumerState<PracticeScreen> {
  int? _vocabCount;
  int _mistakeCount = 0;
  int _listeningCount = 0;
  int _readingCount = 0;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final r = await ApiClient.instance.practice();
      if (!mounted) return;
      setState(() {
        _vocabCount = r.vocab.length;
        _mistakeCount = r.mistakes.length;
        _listeningCount = r.listeningCount;
        _readingCount = r.readingCount;
      });
    } catch (_) {
      if (mounted) setState(() => _vocabCount = 0);
    }
  }

  @override
  Widget build(BuildContext context) {
    final lang = ref.watch(authProvider).user?.targetLanguage ?? '';
    final hasVocab = (_vocabCount ?? 0) > 0;

    final modes = [
      _Mode(key: 'quiz', title: 'Vocabulary Quiz', desc: 'Fresh words every time — choose the meaning.', icon: Icons.checklist_rounded, tint: LumoraColors.purple, disabled: !hasVocab),
      _Mode(key: 'listening', title: 'Listening Comprehension', desc: _listeningCount > 0 ? '$_listeningCount unlocked conversations · listen & answer' : 'Unlock units to practise listening.', icon: Icons.forum_rounded, tint: LumoraColors.teal, disabled: _listeningCount == 0),
      _Mode(key: 'reading', title: 'Reading Comprehension', desc: _readingCount > 0 ? '$_readingCount unlocked passages · read & answer' : 'Unlock units to practise reading.', icon: Icons.menu_book_rounded, tint: const Color(0xFF17A3DD), disabled: _readingCount == 0),
      _Mode(key: 'listen', title: 'Listening Drill', desc: 'Hear a word, choose what it means.', icon: Icons.headphones_rounded, tint: const Color(0xFF0E9F8A), disabled: !hasVocab),
      _Mode(key: 'speak', title: 'Speaking Practice', desc: 'Say phrases aloud and get a fluency score.', icon: Icons.mic_rounded, tint: LumoraColors.coral, disabled: !hasVocab),
      _Mode(key: 'mistakes', title: 'Review Mistakes', desc: _mistakeCount > 0 ? '$_mistakeCount to review' : 'Anything you miss appears here.', icon: Icons.refresh_rounded, tint: LumoraColors.amber, disabled: _mistakeCount == 0, badge: _mistakeCount),
    ];

    return Scaffold(
      backgroundColor: LumoraColors.cream,
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.fromLTRB(20, 12, 20, 32),
          children: [
            const Text('Practice', style: TextStyle(fontSize: 28, fontWeight: FontWeight.w800)),
            const SizedBox(height: 4),
            Text('Strengthen your ${languageName(lang.isEmpty ? "es" : lang)} between lessons.', style: const TextStyle(color: LumoraColors.slatey)),
            const SizedBox(height: 16),
            const SpeechBubble(
              child: Text.rich(TextSpan(children: [
                TextSpan(text: 'Blaze: ', style: TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.coral)),
                TextSpan(text: "Ready to warm up? Pick a drill and let's get fluent! 🔥"),
              ])),
            ),
            const SizedBox(height: 16),
            Row(children: [
              Expanded(child: _StatTile(icon: Icons.bookmark_rounded, color: LumoraColors.purple, label: 'WORDS READY', value: _vocabCount == null ? '—' : '$_vocabCount', hint: 'grows as you unlock units')),
              const SizedBox(width: 12),
              Expanded(child: _StatTile(icon: Icons.refresh_rounded, color: LumoraColors.amber, label: 'TO REVIEW', value: '$_mistakeCount', hint: _mistakeCount > 0 ? 'mistakes to fix' : 'all caught up')),
            ]),
            if (hasVocab) ...[
              const SizedBox(height: 14),
              Material(
                color: LumoraColors.purple,
                borderRadius: BorderRadius.circular(LumoraRadii.xl),
                child: InkWell(
                  borderRadius: BorderRadius.circular(LumoraRadii.xl),
                  onTap: () => context.push('/practice/run?mode=mix'),
                  child: Container(
                    padding: const EdgeInsets.all(18),
                    decoration: BoxDecoration(borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.cardLg),
                    child: Row(children: [
                      Container(width: 52, height: 52, decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.2), borderRadius: BorderRadius.circular(LumoraRadii.lg)),
                          child: const Icon(Icons.auto_awesome_rounded, color: Colors.white, size: 26)),
                      const SizedBox(width: 14),
                      Expanded(
                        child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                          Row(children: [
                            const Text('Daily Mix', style: TextStyle(color: Colors.white, fontWeight: FontWeight.w800, fontSize: 16)),
                            const SizedBox(width: 8),
                            Container(padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2), decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.25), borderRadius: BorderRadius.circular(LumoraRadii.full)),
                                child: const Text('Recommended', style: TextStyle(color: Colors.white, fontSize: 10, fontWeight: FontWeight.w800))),
                          ]),
                          const Text('A fresh blend of new words — quiz, listening & speaking.', style: TextStyle(color: Colors.white70, fontSize: 12)),
                        ]),
                      ),
                      const Icon(Icons.chevron_right, color: Colors.white70),
                    ]),
                  ),
                ),
              ),
            ],
            const SizedBox(height: 14),
            for (final m in modes)
              Padding(
                padding: const EdgeInsets.only(bottom: 10),
                child: Opacity(
                  opacity: m.disabled ? 0.6 : 1,
                  child: Material(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(LumoraRadii.xl),
                    child: InkWell(
                      borderRadius: BorderRadius.circular(LumoraRadii.xl),
                      onTap: m.disabled ? null : () => context.push('/practice/run?mode=${m.key}'),
                      child: Container(
                        padding: const EdgeInsets.all(16),
                        decoration: BoxDecoration(borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
                        child: Row(children: [
                          Container(width: 48, height: 48, decoration: BoxDecoration(color: m.tint.withValues(alpha: 0.1), borderRadius: BorderRadius.circular(LumoraRadii.lg)),
                              child: Icon(m.icon, color: m.tint, size: 24)),
                          const SizedBox(width: 14),
                          Expanded(
                            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                              Row(children: [
                                Flexible(child: Text(m.title, overflow: TextOverflow.ellipsis, style: const TextStyle(fontWeight: FontWeight.w800))),
                                if (m.badge > 0) ...[
                                  const SizedBox(width: 8),
                                  Container(padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2), decoration: BoxDecoration(color: LumoraColors.amber, borderRadius: BorderRadius.circular(LumoraRadii.full)),
                                      child: Text('${m.badge}', style: const TextStyle(fontSize: 11, fontWeight: FontWeight.w800))),
                                ],
                              ]),
                              Text(m.desc, style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
                            ]),
                          ),
                          if (!m.disabled) const Icon(Icons.chevron_right, color: LumoraColors.gray300),
                        ]),
                      ),
                    ),
                  ),
                ),
              ),
            if (_vocabCount == 0)
              const Padding(
                padding: EdgeInsets.only(top: 8),
                child: Text('Complete a lesson first to unlock practice drills.', textAlign: TextAlign.center, style: TextStyle(color: LumoraColors.slatey, fontSize: 12)),
              ),
          ],
        ),
      ),
    );
  }
}

class _StatTile extends StatelessWidget {
  final IconData icon;
  final Color color;
  final String label;
  final String value;
  final String hint;
  const _StatTile({required this.icon, required this.color, required this.label, required this.value, required this.hint});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            Icon(icon, size: 18, color: color),
            const SizedBox(width: 6),
            Text(label, style: TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: color)),
          ]),
          const SizedBox(height: 4),
          Text(value, style: const TextStyle(fontSize: 24, fontWeight: FontWeight.w800)),
          Text(hint, style: const TextStyle(color: LumoraColors.slatey, fontSize: 11)),
        ],
      ),
    );
  }
}
