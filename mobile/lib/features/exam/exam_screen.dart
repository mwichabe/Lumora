import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:url_launcher/url_launcher.dart';

import '../../core/env.dart';
import '../../core/languages.dart';
import '../../core/network/api_client.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../core/theme/shadows.dart';
import '../../models/user.dart';
import '../../providers/auth_provider.dart';
import '../../widgets/fox_mascot.dart';
import '../../widgets/lumora_button.dart';
import '../payments/checkout_screen.dart';

const _kLevels = [
  ('A1', 'Beginner'), ('A2', 'Elementary'), ('B1', 'Intermediate'),
  ('B2', 'Upper-Int.'), ('C1', 'Advanced'), ('C2', 'Mastery'),
];
const _kFinalCode = 'FINAL';

String get _examUrl => '${Env.webUrl}/exam';

/// The exam hub. Exams themselves are only taken on the web: they're proctored
/// with the camera *and* a screen share, which a phone can't provide. So this
/// screen covers everything around the exam — see the levels and what you've
/// earned, pay for an attempt in-app — and then hands off to the web app to
/// sit it. The attempt is tied to the account, not the device, so an attempt
/// bought here is waiting when the learner signs in on the web.
class ExamScreen extends ConsumerStatefulWidget {
  const ExamScreen({super.key});

  @override
  ConsumerState<ExamScreen> createState() => _ExamScreenState();
}

class _ExamScreenState extends ConsumerState<ExamScreen> {
  bool _loading = true;
  List<String> _completed = [];
  PaymentStatus? _payStatus;

  /// The level whose detail page is open, or null for the level list.
  String? _level;

  @override
  void initState() {
    super.initState();
    _bootstrap();
  }

  Future<void> _bootstrap() async {
    final lang = ref.read(authProvider).user?.targetLanguage ?? 'es';
    await Future.wait([
      ApiClient.instance.certificates().then((certs) {
        _completed = certs.where((c) => c.language == lang).map((c) => c.level).toList();
      }, onError: (_) {}),
      _refreshPayments(),
    ]);
    if (mounted) setState(() => _loading = false);
  }

  Future<void> _refreshPayments() async {
    try {
      final status = await ApiClient.instance.paymentStatus();
      if (mounted) setState(() => _payStatus = status);
    } catch (_) {}
  }

  String get _lang => ref.read(authProvider).user?.targetLanguage ?? 'es';

  bool _paid(String level) => _payStatus?.paid[level] ?? false;

  /// True when this level still needs a paid attempt before it can be sat.
  bool _mustPay(String level) => (_payStatus?.paymentsEnabled ?? false) && !_paid(level);

  Future<void> _pay(String level) async {
    final messenger = ScaffoldMessenger.of(context);
    final paid = await payWithPaystack(context, level: level);
    // Re-read either way: a payment that settles late is applied by the server.
    await _refreshPayments();
    if (paid) messenger.showSnackBar(SnackBar(content: Text('Payment received — your ${levelDisplay(level, _lang)} attempt is ready 🎉')));
  }

  Future<void> _openOnWeb() async {
    final messenger = ScaffoldMessenger.of(context);
    final opened = await launchUrl(Uri.parse(_examUrl), mode: LaunchMode.externalApplication).catchError((_) => false);
    if (!opened) messenger.showSnackBar(const SnackBar(content: Text("Couldn't open the browser — copy the link instead.")));
  }

  Future<void> _copyLink() async {
    final messenger = ScaffoldMessenger.of(context);
    await Clipboard.setData(ClipboardData(text: _examUrl));
    messenger.showSnackBar(const SnackBar(content: Text('Link copied')));
  }

  @override
  Widget build(BuildContext context) {
    final level = _level;
    return PopScope(
      // Back from a level's detail returns to the list, not out of the hub.
      canPop: level == null,
      onPopInvokedWithResult: (didPop, _) {
        if (!didPop) setState(() => _level = null);
      },
      child: Scaffold(
        backgroundColor: LumoraColors.cream,
        body: SafeArea(
          child: _loading
              ? const Center(child: FoxMascot(size: 110, glow: true))
              : Column(
                  children: [
                    Padding(
                      padding: const EdgeInsets.fromLTRB(16, 8, 16, 8),
                      child: Row(children: [
                        IconButton(
                          icon: Icon(level == null ? Icons.close : Icons.arrow_back_rounded, color: LumoraColors.gray500),
                          onPressed: () {
                            if (level != null) {
                              setState(() => _level = null);
                            } else if (context.canPop()) {
                              context.pop();
                            } else {
                              context.go('/profile');
                            }
                          },
                        ),
                      ]),
                    ),
                    Expanded(
                      child: Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20),
                        child: level == null ? _buildLevels() : _buildLevel(level),
                      ),
                    ),
                  ],
                ),
        ),
      ),
    );
  }

  Widget _buildLevels() {
    final firstOpen = _kLevels.map((l) => l.$1).firstWhere((c) => !_completed.contains(c), orElse: () => 'A1');
    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const SizedBox(height: 12),
          const Text('Proficiency Exam', style: TextStyle(fontSize: 24, fontWeight: FontWeight.w800)),
          const SizedBox(height: 4),
          const Text('Pick a level to test your listening, reading, writing and speaking.', style: TextStyle(color: LumoraColors.slatey)),
          const SizedBox(height: 16),
          const _WebOnlyBanner(),
          const SizedBox(height: 16),
          for (final l in _kLevels)
            Padding(
              padding: const EdgeInsets.only(bottom: 8),
              child: _LevelCard(
                code: l.$1, name: _lang == 'zh' ? '${levelDisplay(l.$1, _lang).split(' · ').first} · ${l.$2}' : l.$2, done: _completed.contains(l.$1),
                price: _payStatus?.prices[l.$1], usd: _payStatus?.pricesUsd[l.$1], paid: _paid(l.$1),
                onTap: () => setState(() => _level = l.$1),
              ),
            ),
          const Divider(height: 32),
          _LevelCard(
            code: _kFinalCode, name: _lang == 'zh' ? 'Final — HSK 7–9 advanced band' : 'Final Mastery — comprehensive A1→C2', done: _completed.contains(_kFinalCode),
            price: _payStatus?.prices[_kFinalCode], usd: _payStatus?.pricesUsd[_kFinalCode], paid: _paid(_kFinalCode),
            onTap: () => setState(() => _level = _kFinalCode),
          ),
          const SizedBox(height: 8),
          Center(child: Text('Suggested next: $firstOpen', style: const TextStyle(color: LumoraColors.gray500, fontSize: 12))),
          const SizedBox(height: 24),
        ],
      ),
    );
  }

  Widget _buildLevel(String level) {
    return SingleChildScrollView(
      child: ConstrainedBox(
        constraints: const BoxConstraints(maxWidth: 480),
        child: _mustPay(level) ? _buildPay(level) : _buildTakeOnWeb(level),
      ),
    );
  }

  Widget _buildPay(String level) {
    final price = _payStatus?.prices[level] ?? 0;
    final usd = _payStatus?.pricesUsd[level] ?? 0;
    return Column(
      children: [
        const SizedBox(height: 24),
        const Icon(Icons.lock_rounded, size: 48, color: LumoraColors.amber),
        const SizedBox(height: 12),
        Text('Unlock the $level exam', style: const TextStyle(fontSize: 22, fontWeight: FontWeight.w800)),
        const SizedBox(height: 8),
        Text('KES $price${usd > 0 ? " (≈ \$${usd.toStringAsFixed(2)})" : ""}', style: const TextStyle(fontSize: 16, color: LumoraColors.slatey)),
        const SizedBox(height: 6),
        const Text('One attempt, with a certificate when you pass.', textAlign: TextAlign.center, style: TextStyle(color: LumoraColors.slatey)),
        const SizedBox(height: 24),
        LumoraButton(label: 'Pay KES $price', full: true, onPressed: () => _pay(level)),
        const SizedBox(height: 16),
        const _WebOnlyBanner(
          text: 'Pay here, then sit the exam on the web. Your attempt is saved to your account, so it will be waiting when you sign in there.',
        ),
        const SizedBox(height: 24),
      ],
    );
  }

  Widget _buildTakeOnWeb(String level) {
    final paid = _paid(level);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        const SizedBox(height: 24),
        Center(
          child: Container(
            width: 80, height: 80,
            decoration: BoxDecoration(color: LumoraColors.teal.withValues(alpha: 0.1), shape: BoxShape.circle),
            child: Icon(paid ? Icons.check_circle : Icons.laptop_mac_rounded, color: LumoraColors.teal, size: 44),
          ),
        ),
        const SizedBox(height: 16),
        Text(paid ? 'Your ${levelDisplay(level, _lang)} attempt is ready' : 'Take the ${levelDisplay(level, _lang)} exam on the web',
            textAlign: TextAlign.center, style: const TextStyle(fontSize: 22, fontWeight: FontWeight.w800)),
        const SizedBox(height: 8),
        const Text(
          'Exams are proctored with your camera and a screen share, so they run in a web browser on a computer — not in the app.',
          textAlign: TextAlign.center,
          style: TextStyle(color: LumoraColors.slatey, height: 1.4),
        ),
        const SizedBox(height: 20),
        Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
          child: Column(children: [
            const _Step(n: '1', text: 'On a computer, open the link below in Chrome, Edge or Firefox.'),
            _Step(n: '2', text: 'Sign in as ${ref.watch(authProvider).user?.email ?? "the same account"}.'),
            _Step(n: '3', text: 'Open Exam, choose $level and begin.', last: true),
          ]),
        ),
        const SizedBox(height: 12),
        Material(
          color: LumoraColors.purpleLight,
          borderRadius: BorderRadius.circular(LumoraRadii.lg),
          child: InkWell(
            borderRadius: BorderRadius.circular(LumoraRadii.lg),
            onTap: _copyLink,
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
              child: Row(children: [
                Expanded(
                  child: Text(_examUrl.replaceFirst(RegExp(r'^https?://'), ''),
                      overflow: TextOverflow.ellipsis, style: const TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.purple)),
                ),
                const SizedBox(width: 8),
                const Icon(Icons.copy_rounded, size: 18, color: LumoraColors.purple),
              ]),
            ),
          ),
        ),
        const SizedBox(height: 16),
        LumoraButton(label: 'Open in browser', full: true, variant: LumoraButtonVariant.outline, onPressed: _openOnWeb),
        const SizedBox(height: 24),
      ],
    );
  }
}

class _WebOnlyBanner extends StatelessWidget {
  final String text;
  const _WebOnlyBanner({this.text = 'Exams are taken on the web, on a computer. You can pay for an attempt here.'});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(color: LumoraColors.amberLight, borderRadius: BorderRadius.circular(LumoraRadii.lg)),
      child: Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
        const Icon(Icons.laptop_mac_rounded, size: 18, color: LumoraColors.amber),
        const SizedBox(width: 10),
        Expanded(child: Text(text, style: const TextStyle(color: LumoraColors.slatey, fontSize: 13, height: 1.35))),
      ]),
    );
  }
}

class _Step extends StatelessWidget {
  final String n;
  final String text;
  final bool last;
  const _Step({required this.n, required this.text, this.last = false});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.only(bottom: last ? 0 : 12),
      child: Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Container(
          width: 24, height: 24,
          decoration: const BoxDecoration(color: LumoraColors.purple, shape: BoxShape.circle),
          child: Center(child: Text(n, style: const TextStyle(color: Colors.white, fontWeight: FontWeight.w800, fontSize: 12))),
        ),
        const SizedBox(width: 10),
        Expanded(child: Padding(padding: const EdgeInsets.only(top: 2), child: Text(text, style: const TextStyle(color: LumoraColors.ink, height: 1.35)))),
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
          decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
          child: Row(children: [
            Container(width: 44, height: 44, decoration: BoxDecoration(color: LumoraColors.purpleLight, borderRadius: BorderRadius.circular(LumoraRadii.md)),
                child: Center(child: Text(code, style: const TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.purple, fontSize: 12)))),
            const SizedBox(width: 12),
            Expanded(
              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Text(name, style: const TextStyle(fontWeight: FontWeight.w800)),
                if (paid) const Text('Paid — take it on the web', style: TextStyle(color: LumoraColors.purple, fontSize: 11, fontWeight: FontWeight.w700))
                else if (done) const Text('Certificate earned', style: TextStyle(color: LumoraColors.teal, fontSize: 11, fontWeight: FontWeight.w700))
                else if (price != null && price! > 0)
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
