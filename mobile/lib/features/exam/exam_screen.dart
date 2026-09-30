import 'dart:async';

import 'package:camera/camera.dart' as cam;
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:permission_handler/permission_handler.dart';
import 'package:url_launcher/url_launcher.dart';

import '../../core/languages.dart';
import '../../core/network/api_client.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../core/theme/shadows.dart';
import '../../core/voices.dart';
import '../../models/exam.dart';
import '../../models/user.dart';
import '../../providers/auth_provider.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../widgets/confirm_dialog.dart';
import '../../widgets/fox_mascot.dart';
import '../../widgets/lumora_button.dart';

enum _Phase { intro, pay, rules, listening, reading, writing, speaking, result }

const _kLevels = [
  ('A1', 'Beginner'), ('A2', 'Elementary'), ('B1', 'Intermediate'),
  ('B2', 'Upper-Int.'), ('C1', 'Advanced'), ('C2', 'Mastery'),
];
const _kFinalCode = 'FINAL';
const _kFallbackDuration = [600, 780, 960, 1140, 1320, 1500];

int _levelIdx(String code) {
  final i = _kLevels.indexWhere((l) => l.$1 == code);
  return i < 0 ? 0 : i;
}

class ExamScreen extends ConsumerStatefulWidget {
  const ExamScreen({super.key});

  @override
  ConsumerState<ExamScreen> createState() => _ExamScreenState();
}

class _ExamScreenState extends ConsumerState<ExamScreen> with WidgetsBindingObserver {
  bool _loading = true;
  List<String> _completed = [];
  PaymentStatus? _payStatus;

  String _level = 'A1';
  _Phase _phase = _Phase.intro;
  bool _paperLoading = false;
  ExamPaper? _paper;
  bool _unlocking = false;

  int _listening = 0, _reading = 0, _writing = 0, _speaking = 0;
  ExamResult? _result;
  bool _submitting = false;
  String? _setupError;

  int _secondsLeft = 0;
  Timer? _timer;
  bool _ended = false;
  bool _active = false;

  cam.CameraController? _camController;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
    _bootstrap();
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    _timer?.cancel();
    _camController?.dispose();
    Voices.instance.stopSpeaking();
    super.dispose();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    if (state == AppLifecycleState.paused && _active && !_ended) {
      _endExam();
    }
  }

  Future<void> _bootstrap() async {
    final lang = ref.read(authProvider).user?.targetLanguage ?? 'es';
    try {
      final certs = await ApiClient.instance.certificates();
      _completed = certs.where((c) => c.language == lang).map((c) => c.level).toList();
    } catch (_) {}
    try {
      _payStatus = await ApiClient.instance.paymentStatus();
    } catch (_) {}
    if (mounted) setState(() => _loading = false);
  }

  Future<void> _selectLevel(String lvl) async {
    setState(() {
      _level = lvl;
      _listening = _reading = _writing = _speaking = 0;
      _setupError = null;
      _paper = null;
      _paperLoading = true;
      final mustPay = (_payStatus?.paymentsEnabled ?? false) && !(_payStatus?.paid[lvl] ?? false);
      _phase = mustPay ? _Phase.pay : _Phase.rules;
    });
    try {
      _paper = await ApiClient.instance.examPaper(lvl);
    } catch (_) {
      _paper = null;
    } finally {
      if (mounted) setState(() => _paperLoading = false);
    }
  }

  Future<void> _payForLevel() async {
    setState(() => _unlocking = true);
    try {
      final (url, _) = await ApiClient.instance.initializePayment(level: _level);
      if (url != null) await launchUrl(Uri.parse(url), mode: LaunchMode.externalApplication);
    } catch (_) {
    } finally {
      if (mounted) setState(() => _unlocking = false);
    }
  }

  Future<void> _refreshPaymentAndContinue() async {
    try {
      _payStatus = await ApiClient.instance.paymentStatus();
      if (mounted && (_payStatus?.paid[_level] ?? false)) {
        setState(() => _phase = _Phase.rules);
      } else if (mounted) {
        setState(() {});
      }
    } catch (_) {}
  }

  Future<void> _beginExam() async {
    setState(() => _setupError = null);
    final status = await Permission.camera.request();
    if (!status.isGranted) {
      setState(() => _setupError = 'Camera access is required to take this exam. Allow your camera in Settings, then try again.');
      return;
    }
    try {
      final cameras = await cam.availableCameras();
      final front = cameras.firstWhere((c) => c.lensDirection == cam.CameraLensDirection.front, orElse: () => cameras.first);
      _camController = cam.CameraController(front, cam.ResolutionPreset.low, enableAudio: false);
      await _camController!.initialize();
    } catch (_) {
      setState(() => _setupError = "Couldn't start the camera. Check no other app is using it, then try again.");
      return;
    }

    ApiClient.instance.startExam(_level, ref.read(authProvider).user?.targetLanguage ?? 'es').catchError((_) {});

    final fbIdx = _levelIdx(_level).clamp(0, 5);
    final duration = _paper?.durationSeconds ?? _kFallbackDuration[fbIdx];
    setState(() {
      _secondsLeft = duration;
      _active = true;
      _ended = false;
      _phase = _Phase.listening;
    });
    _timer = Timer.periodic(const Duration(seconds: 1), (_) {
      if (!mounted) return;
      setState(() => _secondsLeft = (_secondsLeft - 1).clamp(0, 1 << 31));
      if (_secondsLeft <= 0 && _active && !_ended) _endExam();
    });
  }

  Future<void> _endExam() async {
    if (_ended) return;
    _ended = true;
    _active = false;
    _timer?.cancel();
    await _camController?.dispose();
    _camController = null;
    Voices.instance.stopSpeaking();
    await _submit();
  }

  Future<void> _submit() async {
    setState(() {
      _submitting = true;
      _phase = _Phase.result;
    });
    try {
      final lang = ref.read(authProvider).user?.targetLanguage ?? 'es';
      final r = await ApiClient.instance.submitExam(
        language: lang, level: _level, listening: _listening, reading: _reading, writing: _writing, speaking: _speaking,
      );
      setState(() => _result = r);
    } catch (_) {
      setState(() => _result = null);
    } finally {
      if (mounted) setState(() => _submitting = false);
    }
  }

  void _advance(String section, int score, _Phase next) {
    if (_ended) return;
    setState(() {
      switch (section) {
        case 'listening': _listening = score; break;
        case 'reading': _reading = score; break;
        case 'writing': _writing = score; break;
        case 'speaking': _speaking = score; break;
      }
    });
    if (next == _Phase.result) {
      _ended = true;
      _active = false;
      _timer?.cancel();
      _camController?.dispose();
      _camController = null;
      _submit();
    } else {
      setState(() => _phase = next);
    }
  }

  void _retake() {
    _timer?.cancel();
    _camController?.dispose();
    _camController = null;
    Voices.instance.stopSpeaking();
    setState(() {
      _ended = false;
      _active = false;
      _setupError = null;
      _result = null;
      _secondsLeft = 0;
      _listening = _reading = _writing = _speaking = 0;
      _phase = _Phase.intro;
    });
    ApiClient.instance.paymentStatus().then((p) { if (mounted) setState(() => _payStatus = p); }).catchError((_) {});
  }

  String _fmtTime(int s) {
    final m = s ~/ 60, r = s % 60;
    return '$m:${r.toString().padLeft(2, '0')}';
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: LumoraColors.cream,
      body: SafeArea(
        child: _loading
            ? const Center(child: FoxMascot(size: 110, glow: true))
            : Column(
                children: [
                  if (_phase != _Phase.result)
                    Padding(
                      padding: const EdgeInsets.fromLTRB(16, 8, 16, 8),
                      child: Row(children: [
                        IconButton(icon: const Icon(Icons.close, color: LumoraColors.gray500), onPressed: () async {
                          if (_active && !_ended) {
                            final leave = await showLumoraConfirmDialog(
                              context,
                              title: 'Leave exam?',
                              message: "Leaving now will end your exam and score only what you've completed.",
                              cancelLabel: 'Stay',
                              confirmLabel: 'Leave',
                              danger: true,
                            );
                            if (!leave) return;
                          }
                          Voices.instance.stopSpeaking();
                          if (mounted) context.go('/profile');
                        }),
                        if (_active)
                          Expanded(
                            child: Align(
                              alignment: Alignment.centerRight,
                              child: Container(
                                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
                                decoration: BoxDecoration(color: _secondsLeft < 60 ? LumoraColors.coralLight : LumoraColors.gray50, borderRadius: BorderRadius.circular(LumoraRadii.full)),
                                child: Text(_fmtTime(_secondsLeft), style: TextStyle(fontWeight: FontWeight.w800, color: _secondsLeft < 60 ? LumoraColors.coral : LumoraColors.ink)),
                              ),
                            ),
                          ),
                      ]),
                    ),
                  Expanded(
                    child: Stack(
                      children: [
                        Padding(padding: const EdgeInsets.symmetric(horizontal: 20), child: _buildPhase()),
                        if (_active && _camController != null && _camController!.value.isInitialized)
                          Positioned(
                            top: 8, right: 8,
                            child: Container(
                              width: 72, height: 96,
                              decoration: BoxDecoration(borderRadius: BorderRadius.circular(LumoraRadii.md), border: Border.all(color: LumoraColors.coral, width: 2), boxShadow: LumoraShadows.card),
                              child: ClipRRect(borderRadius: BorderRadius.circular(LumoraRadii.sm), child: cam.CameraPreview(_camController!)),
                            ),
                          ),
                        if (_active && _camController != null)
                          const Positioned(
                            top: 12, right: 84,
                            child: Icon(Icons.fiber_manual_record, color: LumoraColors.coral, size: 14),
                          ),
                      ],
                    ),
                  ),
                ],
              ),
      ),
    );
  }

  Widget _buildPhase() {
    switch (_phase) {
      case _Phase.intro:
        return _buildIntro();
      case _Phase.pay:
        return _buildPay();
      case _Phase.rules:
        return _buildRules();
      case _Phase.listening:
        return _buildListening();
      case _Phase.reading:
        return _buildReading();
      case _Phase.writing:
        return _buildWriting();
      case _Phase.speaking:
        return _buildSpeaking();
      case _Phase.result:
        return _buildResult();
    }
  }

  Widget _buildIntro() {
    final firstOpen = _kLevels.map((l) => l.$1).firstWhere((c) => !_completed.contains(c), orElse: () => 'A1');
    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const SizedBox(height: 12),
          const Text('Proficiency Exam', style: TextStyle(fontSize: 24, fontWeight: FontWeight.w800)),
          const SizedBox(height: 4),
          const Text('Pick a level to test your listening, reading, writing and speaking.', style: TextStyle(color: LumoraColors.slatey)),
          const SizedBox(height: 20),
          for (final l in _kLevels)
            Padding(
              padding: const EdgeInsets.only(bottom: 8),
              child: _LevelCard(
                code: l.$1, name: l.$2, done: _completed.contains(l.$1),
                price: _payStatus?.prices[l.$1], usd: _payStatus?.pricesUsd[l.$1], paid: _payStatus?.paid[l.$1] ?? false,
                onTap: () => _selectLevel(l.$1),
              ),
            ),
          const Divider(height: 32),
          _LevelCard(
            code: _kFinalCode, name: 'Final Mastery — comprehensive A1→C2', done: _completed.contains(_kFinalCode),
            price: _payStatus?.prices[_kFinalCode], usd: _payStatus?.pricesUsd[_kFinalCode], paid: _payStatus?.paid[_kFinalCode] ?? false,
            onTap: () => _selectLevel(_kFinalCode),
          ),
          const SizedBox(height: 8),
          Center(child: Text('Suggested next: $firstOpen', style: const TextStyle(color: LumoraColors.gray500, fontSize: 12))),
          const SizedBox(height: 24),
        ],
      ),
    );
  }

  Widget _buildPay() {
    final price = _payStatus?.prices[_level] ?? 0;
    final usd = _payStatus?.pricesUsd[_level] ?? 0;
    return Center(
      child: Column(mainAxisSize: MainAxisSize.min, children: [
        const Icon(Icons.lock_rounded, size: 48, color: LumoraColors.amber),
        const SizedBox(height: 12),
        Text('Unlock the $_level exam', style: const TextStyle(fontSize: 20, fontWeight: FontWeight.w800)),
        const SizedBox(height: 8),
        Text('KES $price${usd > 0 ? " (≈ \$${usd.toStringAsFixed(2)})" : ""}', style: const TextStyle(fontSize: 16, color: LumoraColors.slatey)),
        const SizedBox(height: 20),
        LumoraButton(label: _unlocking ? 'Redirecting…' : 'Pay with Paystack', full: true, loading: _unlocking, onPressed: _payForLevel),
        const SizedBox(height: 8),
        LumoraButton(label: "I've paid — continue", full: true, variant: LumoraButtonVariant.outline, onPressed: _refreshPaymentAndContinue),
      ]),
    );
  }

  Widget _buildRules() {
    return _RulesView(
      level: _level,
      paperLoading: _paperLoading,
      setupError: _setupError,
      onBegin: _beginExam,
    );
  }

  Widget _buildListening() {
    final listening = _paper?.listening;
    if (listening == null) return _sectionEmpty('listening', () => _advance('listening', 0, _Phase.reading));
    return _McqSectionView(
      icon: Icons.headphones_rounded, title: 'Listening', tint: LumoraColors.teal,
      intro: _ListeningIntro(listening: listening),
      questions: listening.questions,
      onDone: (score) => _advance('listening', score, _Phase.reading),
    );
  }

  Widget _buildReading() {
    final reading = _paper?.reading;
    if (reading == null) return _sectionEmpty('reading', () => _advance('reading', 0, _Phase.writing));
    return _McqSectionView(
      icon: Icons.menu_book_rounded, title: 'Reading', tint: const Color(0xFF17A3DD),
      intro: _ReadingIntro(reading: reading),
      questions: reading.questions,
      onDone: (score) => _advance('reading', score, _Phase.writing),
    );
  }

  Widget _buildWriting() {
    final writing = _paper?.writing ?? const ExamWritingPaper(prompt: 'about your day.', minWords: 40);
    final lang = ref.read(authProvider).user?.targetLanguage ?? 'es';
    return _WritingView(prompt: writing.prompt, minWords: writing.minWords, lang: lang, onDone: (score) => _advance('writing', score, _Phase.speaking));
  }

  Widget _buildSpeaking() {
    final speaking = _paper?.speaking ?? const ExamSpeakingPaper(phrase: 'Hola', speaker: 'Lumora', translation: 'Hello');
    return _SpeakingView(phrase: speaking.phrase, speaker: speaking.speaker, onDone: (score) => _advance('speaking', score, _Phase.result));
  }

  Widget _sectionEmpty(String section, VoidCallback onSkip) {
    WidgetsBinding.instance.addPostFrameCallback((_) => onSkip());
    return const SizedBox.shrink();
  }

  Widget _buildResult() {
    if (_submitting) {
      return const Center(child: FoxMascot(size: 120, glow: true, bounce: true));
    }
    final r = _result;
    if (r == null) {
      return Center(
        child: Column(mainAxisSize: MainAxisSize.min, children: [
          const Icon(Icons.error_outline, size: 48, color: LumoraColors.coral),
          const SizedBox(height: 12),
          const Text("Couldn't submit your exam", style: TextStyle(fontWeight: FontWeight.w800)),
          const SizedBox(height: 16),
          LumoraButton(label: 'Try again', onPressed: _retake),
        ]),
      );
    }
    return SingleChildScrollView(
      child: Column(
        children: [
          const SizedBox(height: 20),
          Icon(r.passed ? Icons.emoji_events_rounded : Icons.sentiment_neutral_rounded, size: 64, color: r.passed ? LumoraColors.amber : LumoraColors.gray500),
          const SizedBox(height: 12),
          Text(r.passed ? 'You passed!' : 'Not quite there', style: const TextStyle(fontSize: 22, fontWeight: FontWeight.w800)),
          const SizedBox(height: 4),
          Text('${r.overall}% overall${r.passMark != null ? " · pass mark ${r.passMark}%" : ""}', style: const TextStyle(color: LumoraColors.slatey)),
          if (r.sections != null) ...[
            const SizedBox(height: 20),
            _sectionScore('Listening', r.sections!.listening),
            _sectionScore('Reading', r.sections!.reading),
            _sectionScore('Writing', r.sections!.writing),
            _sectionScore('Speaking', r.sections!.speaking),
          ],
          const SizedBox(height: 24),
          if (r.passed && r.certificate != null)
            LumoraButton(label: 'View certificate', full: true, onPressed: () => context.push('/certificates/${r.certificate!.id}')),
          const SizedBox(height: 8),
          LumoraButton(label: 'Take another exam', full: true, variant: LumoraButtonVariant.outline, onPressed: _retake),
          const SizedBox(height: 8),
          LumoraButton(label: 'Back to profile', full: true, variant: LumoraButtonVariant.ghost, onPressed: () => context.go('/profile')),
          const SizedBox(height: 24),
        ],
      ),
    );
  }

  Widget _sectionScore(String label, int score) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(children: [
        SizedBox(width: 90, child: Text(label, style: const TextStyle(fontWeight: FontWeight.w700))),
        Expanded(
          child: ClipRRect(
            borderRadius: BorderRadius.circular(LumoraRadii.full),
            child: LinearProgressIndicator(value: score / 100, minHeight: 8, backgroundColor: LumoraColors.gray100, valueColor: const AlwaysStoppedAnimation(LumoraColors.purple)),
          ),
        ),
        const SizedBox(width: 8),
        Text('$score%', style: const TextStyle(fontWeight: FontWeight.w800)),
      ]),
    );
  }
}

class _LevelCard extends StatelessWidget {
  final String code;
  final String name;
  final bool done;
  final int? price;
  final double? usd;
  final bool paid;
  final VoidCallback onTap;
  const _LevelCard({required this.code, required this.name, required this.done, this.price, this.usd, required this.paid, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return Material(
      color: Colors.white,
      borderRadius: BorderRadius.circular(LumoraRadii.xl),
      child: InkWell(
        borderRadius: BorderRadius.circular(LumoraRadii.xl),
        onTap: onTap,
        child: Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
          child: Row(children: [
            Container(width: 44, height: 44, decoration: BoxDecoration(color: LumoraColors.purpleLight, borderRadius: BorderRadius.circular(LumoraRadii.md)),
                child: Center(child: Text(code, style: const TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.purple, fontSize: 12)))),
            const SizedBox(width: 12),
            Expanded(
              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Text(name, style: const TextStyle(fontWeight: FontWeight.w800)),
                if (done) const Text('Certificate earned', style: TextStyle(color: LumoraColors.teal, fontSize: 11, fontWeight: FontWeight.w700))
                else if (price != null && price! > 0 && !paid)
                  Text('KES $price${usd != null && usd! > 0 ? " (≈ \$${usd!.toStringAsFixed(2)})" : ""}', style: const TextStyle(color: LumoraColors.slatey, fontSize: 11)),
              ]),
            ),
            const Icon(Icons.chevron_right, color: LumoraColors.gray300),
          ]),
        ),
      ),
    );
  }
}

class _RulesView extends StatefulWidget {
  final String level;
  final bool paperLoading;
  final String? setupError;
  final Future<void> Function() onBegin;
  const _RulesView({required this.level, required this.paperLoading, required this.setupError, required this.onBegin});

  @override
  State<_RulesView> createState() => _RulesViewState();
}

class _RulesViewState extends State<_RulesView> {
  bool _agree = false;
  bool _starting = false;

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const SizedBox(height: 12),
          Text('${widget.level} Exam — before you begin', style: const TextStyle(fontSize: 20, fontWeight: FontWeight.w800)),
          const SizedBox(height: 12),
          const _RuleRow(icon: Icons.timer_outlined, text: 'The exam is timed. Once started, the clock cannot be paused.'),
          const _RuleRow(icon: Icons.camera_alt_outlined, text: 'Your camera stays on for proctoring for the duration of the exam.'),
          const _RuleRow(icon: Icons.phonelink_off_outlined, text: 'Leaving the app or backgrounding it ends the exam immediately and scores only what you completed.'),
          const _RuleRow(icon: Icons.replay_outlined, text: "You can't pause or restart once you begin."),
          const SizedBox(height: 16),
          CheckboxListTile(
            value: _agree,
            onChanged: (v) => setState(() => _agree = v ?? false),
            title: const Text('I understand and agree to be proctored for this exam.'),
            controlAffinity: ListTileControlAffinity.leading,
            contentPadding: EdgeInsets.zero,
          ),
          if (widget.setupError != null)
            Padding(
              padding: const EdgeInsets.only(top: 8),
              child: Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(color: LumoraColors.coralLight, borderRadius: BorderRadius.circular(LumoraRadii.md)),
                child: Text(widget.setupError!, style: const TextStyle(color: LumoraColors.coral, fontSize: 13)),
              ),
            ),
          const SizedBox(height: 16),
          LumoraButton(
            label: widget.paperLoading ? 'Preparing exam…' : 'Begin exam',
            full: true,
            loading: _starting,
            onPressed: (!_agree || widget.paperLoading)
                ? null
                : () async {
                    setState(() => _starting = true);
                    await widget.onBegin();
                    if (mounted) setState(() => _starting = false);
                  },
          ),
          const SizedBox(height: 24),
        ],
      ),
    );
  }
}

class _RuleRow extends StatelessWidget {
  final IconData icon;
  final String text;
  const _RuleRow({required this.icon, required this.text});
  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 6),
      child: Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Icon(icon, size: 18, color: LumoraColors.purple),
        const SizedBox(width: 10),
        Expanded(child: Text(text, style: const TextStyle(fontSize: 13))),
      ]),
    );
  }
}

class _ListeningIntro extends StatefulWidget {
  final ExamListeningPaper listening;
  const _ListeningIntro({required this.listening});
  @override
  State<_ListeningIntro> createState() => _ListeningIntroState();
}

class _ListeningIntroState extends State<_ListeningIntro> {
  bool _played = false;
  @override
  Widget build(BuildContext context) {
    return Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      Text(widget.listening.title, style: const TextStyle(fontWeight: FontWeight.w800)),
      const SizedBox(height: 8),
      LumoraButton(
        label: _played ? 'Play again' : 'Play the recording',
        full: true,
        icon: const Icon(Icons.volume_up_rounded, color: Colors.white, size: 18),
        onPressed: () {
          setState(() => _played = true);
          Voices.instance.speakSequence([for (final l in widget.listening.lines) (character: l.character, text: l.text)]);
        },
      ),
      const SizedBox(height: 4),
      const Text('You can replay the recording while you answer.', textAlign: TextAlign.center, style: TextStyle(color: LumoraColors.slatey, fontSize: 12)),
    ]);
  }
}

class _ReadingIntro extends StatelessWidget {
  final ExamReadingPaper reading;
  const _ReadingIntro({required this.reading});
  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(color: Colors.white, border: Border.all(color: LumoraColors.gray100), borderRadius: BorderRadius.circular(LumoraRadii.xl)),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Text(reading.title, style: const TextStyle(fontWeight: FontWeight.w800)),
        for (final p in reading.paragraphs) Padding(padding: const EdgeInsets.only(top: 6), child: Text(p, style: const TextStyle(height: 1.4))),
      ]),
    );
  }
}

class _McqSectionView extends StatefulWidget {
  final IconData icon;
  final String title;
  final Color tint;
  final Widget intro;
  final List<PaperQuestion> questions;
  final ValueChanged<int> onDone;
  const _McqSectionView({required this.icon, required this.title, required this.tint, required this.intro, required this.questions, required this.onDone});

  @override
  State<_McqSectionView> createState() => _McqSectionViewState();
}

class _McqSectionViewState extends State<_McqSectionView> {
  final Map<int, String> _answers = {};

  void _finish() {
    var correct = 0;
    for (var i = 0; i < widget.questions.length; i++) {
      if (_answers[i] == widget.questions[i].correctAnswer) correct++;
    }
    widget.onDone(widget.questions.isEmpty ? 0 : ((correct / widget.questions.length) * 100).round());
  }

  @override
  Widget build(BuildContext context) {
    final allAnswered = List.generate(widget.questions.length, (i) => _answers.containsKey(i)).every((x) => x);
    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const SizedBox(height: 8),
          Row(children: [Icon(widget.icon, size: 18, color: widget.tint), const SizedBox(width: 8), Text(widget.title.toUpperCase(), style: TextStyle(color: widget.tint, fontWeight: FontWeight.w800, fontSize: 12))]),
          const SizedBox(height: 10),
          widget.intro,
          const SizedBox(height: 16),
          for (var qi = 0; qi < widget.questions.length; qi++)
            Padding(
              padding: const EdgeInsets.only(bottom: 16),
              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Text('${qi + 1}. ${widget.questions[qi].question}', style: const TextStyle(fontWeight: FontWeight.w800)),
                const SizedBox(height: 8),
                for (final opt in widget.questions[qi].options)
                  Padding(
                    padding: const EdgeInsets.only(bottom: 6),
                    child: Material(
                      color: _answers[qi] == opt ? LumoraColors.purpleLight : Colors.white,
                      borderRadius: BorderRadius.circular(LumoraRadii.md),
                      child: InkWell(
                        borderRadius: BorderRadius.circular(LumoraRadii.md),
                        onTap: () => setState(() => _answers[qi] = opt),
                        child: Container(
                          height: 48, alignment: Alignment.centerLeft, padding: const EdgeInsets.symmetric(horizontal: 14),
                          decoration: BoxDecoration(border: Border.all(color: _answers[qi] == opt ? LumoraColors.purple : LumoraColors.gray100, width: 2), borderRadius: BorderRadius.circular(LumoraRadii.md)),
                          child: Text(opt, style: const TextStyle(fontWeight: FontWeight.w600)),
                        ),
                      ),
                    ),
                  ),
              ]),
            ),
          LumoraButton(label: 'Submit section', full: true, onPressed: allAnswered ? _finish : null),
          const SizedBox(height: 24),
        ],
      ),
    );
  }
}

class _WritingView extends StatefulWidget {
  final String prompt;
  final int minWords;
  final String lang;
  final ValueChanged<int> onDone;
  const _WritingView({required this.prompt, required this.minWords, required this.lang, required this.onDone});
  @override
  State<_WritingView> createState() => _WritingViewState();
}

class _WritingViewState extends State<_WritingView> {
  final _controller = TextEditingController();
  int _words = 0;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SizedBox(height: 8),
        const Row(children: [Icon(Icons.edit_rounded, size: 18, color: LumoraColors.amber), SizedBox(width: 8), Text('WRITING', style: TextStyle(color: LumoraColors.amber, fontWeight: FontWeight.w800, fontSize: 12))]),
        const SizedBox(height: 10),
        Text('Write in ${languageName(widget.lang)}, ${widget.prompt}'),
        const SizedBox(height: 12),
        TextField(
          controller: _controller,
          maxLines: 8,
          onChanged: (v) => setState(() => _words = v.trim().isEmpty ? 0 : v.trim().split(RegExp(r'\s+')).length),
          decoration: const InputDecoration(hintText: 'Type your answer here…'),
        ),
        const SizedBox(height: 8),
        Text('$_words words ${_words < widget.minWords ? "· aim for at least ${widget.minWords}" : "· nice!"}', style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
        const Spacer(),
        Padding(
          padding: const EdgeInsets.only(bottom: 16, top: 12),
          child: LumoraButton(label: 'Submit section', full: true, onPressed: _words < 3 ? null : () => widget.onDone(((_words / widget.minWords) * 100).clamp(0, 100).round())),
        ),
      ],
    );
  }
}

class _SpeakingView extends StatefulWidget {
  final String phrase;
  final String speaker;
  final ValueChanged<int> onDone;
  const _SpeakingView({required this.phrase, required this.speaker, required this.onDone});
  @override
  State<_SpeakingView> createState() => _SpeakingViewState();
}

class _SpeakingViewState extends State<_SpeakingView> {
  bool _listening = false;
  int? _score;

  Future<void> _record() async {
    if (_listening) return;
    setState(() { _listening = true; _score = null; });
    try {
      final said = await Voices.instance.recognizeSpeech();
      setState(() => _score = scorePronunciation(widget.phrase, said));
    } catch (_) {
      setState(() => _score = null);
    } finally {
      if (mounted) setState(() => _listening = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        const SizedBox(height: 8),
        const Row(children: [Icon(Icons.mic_rounded, size: 18, color: LumoraColors.coral), SizedBox(width: 8), Text('SPEAKING', style: TextStyle(color: LumoraColors.coral, fontWeight: FontWeight.w800, fontSize: 12))]),
        const SizedBox(height: 20),
        Text('"${widget.phrase}"', style: const TextStyle(fontSize: 22, fontWeight: FontWeight.w800, color: LumoraColors.purple)),
        const SizedBox(height: 24),
        Row(mainAxisAlignment: MainAxisAlignment.center, children: [
          InkWell(onTap: () => Voices.instance.speakAs(widget.speaker, widget.phrase),
              child: Container(width: 52, height: 52, decoration: const BoxDecoration(color: LumoraColors.purple, shape: BoxShape.circle), child: const Icon(Icons.volume_up_rounded, color: Colors.white))),
          const SizedBox(width: 20),
          InkWell(onTap: _record,
              child: Container(width: 72, height: 72, decoration: const BoxDecoration(color: LumoraColors.coral, shape: BoxShape.circle), child: Icon(_listening ? Icons.hearing : Icons.mic, color: Colors.white, size: 32))),
        ]),
        if (_score != null) ...[
          const SizedBox(height: 20),
          Text('$_score%', style: const TextStyle(fontSize: 32, fontWeight: FontWeight.w800)),
        ],
        const Spacer(),
        Padding(padding: const EdgeInsets.only(bottom: 16), child: LumoraButton(label: 'Submit section', full: true, onPressed: () => widget.onDone(_score ?? 0))),
      ],
    );
  }
}
