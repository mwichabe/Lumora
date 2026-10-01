import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';
import 'package:share_plus/share_plus.dart';

import '../../core/env.dart';
import '../../core/languages.dart';
import '../../core/network/api_client.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../core/theme/shadows.dart';
import '../../models/exam.dart';
import '../../widgets/fox_mascot.dart';

class CertificateDetailScreen extends StatefulWidget {
  final int id;
  const CertificateDetailScreen({super.key, required this.id});

  @override
  State<CertificateDetailScreen> createState() => _CertificateDetailScreenState();
}

class _CertificateDetailScreenState extends State<CertificateDetailScreen> {
  Certificate? _cert;

  @override
  void initState() {
    super.initState();
    ApiClient.instance.certificate(widget.id).then((c) {
      if (mounted) setState(() => _cert = c);
    }).catchError((_) {});
  }

  @override
  Widget build(BuildContext context) {
    final c = _cert;
    return Scaffold(
      backgroundColor: LumoraColors.cream,
      appBar: AppBar(
        backgroundColor: LumoraColors.cream,
        foregroundColor: LumoraColors.ink,
        elevation: 0,
        leading: BackButton(onPressed: () => context.canPop() ? context.pop() : context.go('/certificates')),
        actions: c == null
            ? null
            : [
                IconButton(
                  icon: const Icon(Icons.share_outlined),
                  onPressed: () => Share.share('Verify my Lumora ${levelDisplay(c.level, c.language)} certificate: ${Env.apiUrl}/api/verify/${c.serial}'),
                ),
              ],
      ),
      body: c == null
          ? const Center(child: FoxMascot(size: 110, glow: true))
          : SingleChildScrollView(
              padding: const EdgeInsets.all(20),
              child: Column(
                children: [
                  Container(
                    padding: const EdgeInsets.all(24),
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(LumoraRadii.xl2),
                      border: Border.all(color: LumoraColors.purple, width: 3),
                      boxShadow: LumoraShadows.cardLg,
                    ),
                    child: Column(
                      children: [
                        const Icon(Icons.workspace_premium_rounded, color: LumoraColors.amber, size: 56),
                        const SizedBox(height: 8),
                        const Text('CERTIFICATE OF PROFICIENCY', style: TextStyle(fontSize: 12, fontWeight: FontWeight.w800, color: LumoraColors.purple, letterSpacing: 1)),
                        const SizedBox(height: 16),
                        Text(c.userName, style: const TextStyle(fontSize: 24, fontWeight: FontWeight.w800)),
                        const SizedBox(height: 4),
                        Text('has achieved', style: TextStyle(color: LumoraColors.slatey.withValues(alpha: 0.8), fontStyle: FontStyle.italic)),
                        const SizedBox(height: 8),
                        Text('${levelDisplay(c.level, c.language)} — ${languageName(c.language)}', style: const TextStyle(fontSize: 20, fontWeight: FontWeight.w800, color: LumoraColors.purple)),
                        const SizedBox(height: 16),
                        Text('Overall score: ${c.score}%', style: const TextStyle(fontWeight: FontWeight.w700)),
                        const SizedBox(height: 16),
                        Row(mainAxisAlignment: MainAxisAlignment.spaceAround, children: [
                          _score('Listening', c.listening),
                          _score('Reading', c.reading),
                          _score('Writing', c.writing),
                          _score('Speaking', c.speaking),
                        ]),
                        const SizedBox(height: 16),
                        Text(DateFormat.yMMMd().format(DateTime.tryParse(c.issuedAt) ?? DateTime.now()), style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
                      ],
                    ),
                  ),
                  const SizedBox(height: 16),
                  Container(
                    width: double.infinity,
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Row(children: [Icon(Icons.verified_rounded, size: 16, color: LumoraColors.teal), SizedBox(width: 6), Text('VERIFICATION', style: TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.gray500))]),
                        const SizedBox(height: 8),
                        Row(children: [
                          Expanded(child: SelectableText(c.serial, style: const TextStyle(fontFamily: 'monospace', fontWeight: FontWeight.w700))),
                          IconButton(
                            icon: const Icon(Icons.copy, size: 16),
                            onPressed: () {
                              Clipboard.setData(ClipboardData(text: c.serial));
                              ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Serial copied')));
                            },
                          ),
                        ]),
                        Text('Anyone can verify this certificate at $verifyPath/${c.serial}', style: const TextStyle(color: LumoraColors.slatey, fontSize: 11)),
                      ],
                    ),
                  ),
                ],
              ),
            ),
    );
  }

  String get verifyPath => '${Env.apiUrl}/verify';

  Widget _score(String label, int value) {
    return Column(children: [
      Text('$value%', style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 16)),
      Text(label, style: const TextStyle(color: LumoraColors.slatey, fontSize: 10)),
    ]);
  }
}
