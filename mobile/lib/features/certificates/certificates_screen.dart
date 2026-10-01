import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../core/languages.dart';
import '../../core/network/api_client.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../core/theme/shadows.dart';
import '../../models/exam.dart';
import '../../widgets/confirm_dialog.dart';
import '../../widgets/fox_mascot.dart';

const _kLevelOrder = ['A1', 'A2', 'B1', 'B2', 'C1', 'C2', 'FINAL'];

class CertificatesScreen extends StatefulWidget {
  const CertificatesScreen({super.key});

  @override
  State<CertificatesScreen> createState() => _CertificatesScreenState();
}

class _CertificatesScreenState extends State<CertificatesScreen> {
  List<Certificate>? _certs;
  String _filter = 'all';

  @override
  void initState() {
    super.initState();
    ApiClient.instance.certificates().then((c) {
      if (mounted) setState(() => _certs = c);
    }).catchError((_) {
      if (mounted) setState(() => _certs = []);
    });
  }

  Future<void> _delete(Certificate c) async {
    final confirm = await showLumoraConfirmDialog(
      context,
      title: 'Delete certificate?',
      message: 'This removes your ${levelDisplay(c.level, c.language)} ${languageName(c.language)} certificate. This cannot be undone.',
      confirmLabel: 'Delete',
      danger: true,
    );
    if (confirm) {
      try {
        await ApiClient.instance.deleteCertificate(c.id);
        setState(() => _certs = _certs!.where((x) => x.id != c.id).toList());
      } catch (_) {}
    }
  }

  @override
  Widget build(BuildContext context) {
    final certs = _certs;
    if (certs == null) {
      return const Scaffold(body: Center(child: FoxMascot(size: 110, glow: true)));
    }
    final languages = certs.map((c) => c.language).toSet().toList();
    final visible = (_filter == 'all' ? certs : certs.where((c) => c.language == _filter).toList())
      ..sort((a, b) {
        final l = a.language.compareTo(b.language);
        if (l != 0) return l;
        return _kLevelOrder.indexOf(a.level).compareTo(_kLevelOrder.indexOf(b.level));
      });

    return Scaffold(
      backgroundColor: LumoraColors.cream,
      appBar: AppBar(
        backgroundColor: LumoraColors.cream,
        foregroundColor: LumoraColors.ink,
        elevation: 0,
        leading: BackButton(onPressed: () => context.canPop() ? context.pop() : context.go('/profile')),
        title: const Text('Certificates', style: TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.ink)),
      ),
      body: certs.isEmpty
          ? const Center(child: Padding(padding: EdgeInsets.all(32), child: Text('Your achievements will appear here.', style: TextStyle(color: LumoraColors.slatey))))
          : Column(
              children: [
                if (languages.length > 1)
                  SizedBox(
                    height: 40,
                    child: ListView(
                      scrollDirection: Axis.horizontal,
                      padding: const EdgeInsets.symmetric(horizontal: 16),
                      children: [
                        Padding(padding: const EdgeInsets.only(right: 8), child: ChoiceChip(label: const Text('All'), selected: _filter == 'all', onSelected: (_) => setState(() => _filter = 'all'))),
                        for (final l in languages)
                          Padding(padding: const EdgeInsets.only(right: 8), child: ChoiceChip(label: Text(languageName(l)), selected: _filter == l, onSelected: (_) => setState(() => _filter = l))),
                      ],
                    ),
                  ),
                Expanded(
                  child: GridView.builder(
                    padding: const EdgeInsets.all(16),
                    gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(crossAxisCount: 2, mainAxisSpacing: 12, crossAxisSpacing: 12, childAspectRatio: 0.85),
                    itemCount: visible.length,
                    itemBuilder: (context, i) {
                      final c = visible[i];
                      return Material(
                        color: Colors.white,
                        borderRadius: BorderRadius.circular(LumoraRadii.xl),
                        child: InkWell(
                          borderRadius: BorderRadius.circular(LumoraRadii.xl),
                          onTap: () => context.push('/certificates/${c.id}'),
                          onLongPress: () => _delete(c),
                          child: Container(
                            padding: const EdgeInsets.all(16),
                            decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Container(
                                  width: 44, height: 44,
                                  decoration: BoxDecoration(color: LumoraColors.purpleLight, borderRadius: BorderRadius.circular(LumoraRadii.md)),
                                  child: const Icon(Icons.workspace_premium_rounded, color: LumoraColors.purple),
                                ),
                                const Spacer(),
                                Text(levelDisplay(c.level, c.language), maxLines: 1, overflow: TextOverflow.ellipsis,
                                    style: TextStyle(fontSize: c.language == 'zh' ? 16 : 20, fontWeight: FontWeight.w800)),
                                Text(languageName(c.language), style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
                                const SizedBox(height: 4),
                                Text('${c.score}% score', style: const TextStyle(color: LumoraColors.teal, fontSize: 11, fontWeight: FontWeight.w700)),
                              ],
                            ),
                          ),
                        ),
                      );
                    },
                  ),
                ),
              ],
            ),
    );
  }
}
