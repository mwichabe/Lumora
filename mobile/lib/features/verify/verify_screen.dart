import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

import '../../core/languages.dart';
import '../../core/network/api_client.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../core/theme/shadows.dart';
import '../../models/exam.dart';
import '../../widgets/fox_mascot.dart';

/// Public, unauthenticated certificate lookup (frontend/app/verify/[serial]).
class VerifyScreen extends StatefulWidget {
  final String serial;
  const VerifyScreen({super.key, required this.serial});

  @override
  State<VerifyScreen> createState() => _VerifyScreenState();
}

class _VerifyScreenState extends State<VerifyScreen> {
  CertVerification? _result;

  @override
  void initState() {
    super.initState();
    ApiClient.instance.verifyCertificate(widget.serial).then((r) {
      if (mounted) setState(() => _result = r);
    });
  }

  @override
  Widget build(BuildContext context) {
    final r = _result;
    return Scaffold(
      backgroundColor: LumoraColors.cream,
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: r == null
              ? const FoxMascot(size: 110, glow: true)
              : Container(
                  padding: const EdgeInsets.all(24),
                  decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl2), boxShadow: LumoraShadows.cardLg),
                  child: Column(mainAxisSize: MainAxisSize.min, children: [
                    Icon(r.valid ? Icons.verified_rounded : Icons.error_outline, size: 56, color: r.valid ? LumoraColors.teal : LumoraColors.coral),
                    const SizedBox(height: 12),
                    Text(r.valid ? 'Certificate verified' : 'Certificate not found', style: const TextStyle(fontSize: 20, fontWeight: FontWeight.w800)),
                    if (r.valid) ...[
                      const SizedBox(height: 16),
                      Text(r.userName ?? '', style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w700)),
                      Text('${r.level} — ${languageName(r.language ?? "es")}', style: const TextStyle(color: LumoraColors.purple, fontWeight: FontWeight.w800)),
                      const SizedBox(height: 8),
                      Text('Score: ${r.score}%'),
                      Text(
                        r.issuedAt != null ? 'Issued ${DateFormat.yMMMd().format(DateTime.tryParse(r.issuedAt!) ?? DateTime.now())}' : '',
                        style: const TextStyle(color: LumoraColors.slatey, fontSize: 12),
                      ),
                      const SizedBox(height: 8),
                      Text('Serial: ${r.serial}', style: const TextStyle(fontFamily: 'monospace', fontSize: 12)),
                    ] else ...[
                      const SizedBox(height: 8),
                      const Text('This serial does not match any issued Lumora certificate.', textAlign: TextAlign.center, style: TextStyle(color: LumoraColors.slatey)),
                    ],
                  ]),
                ),
        ),
      ),
    );
  }
}
