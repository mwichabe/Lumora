import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/languages.dart';
import '../../core/network/api_client.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../providers/auth_provider.dart';

/// Two modes: first-time onboarding (stashes the choice and continues to the
/// daily-goal picker) or "add a language" from an existing account (enrolls
/// immediately and returns to Learn). Mirrors frontend/app/onboarding/language.
class OnboardingLanguageScreen extends ConsumerStatefulWidget {
  final bool adding;
  const OnboardingLanguageScreen({super.key, this.adding = false});

  /// Holds the chosen language between this screen and the goal picker — the
  /// account isn't created with a targetLanguage until /api/auth/setup runs.
  static String? pendingChoice;

  @override
  ConsumerState<OnboardingLanguageScreen> createState() => _OnboardingLanguageScreenState();
}

class _OnboardingLanguageScreenState extends ConsumerState<OnboardingLanguageScreen> {
  final _search = TextEditingController();
  String? _picked;
  bool _loading = false;

  Future<void> _choose(LanguageMeta lang) async {
    if (!lang.available) return;
    setState(() => _picked = lang.code);

    if (widget.adding) {
      setState(() => _loading = true);
      try {
        final (_, _, user) = await ApiClient.instance.enrollLanguage(lang.code);
        ref.read(authProvider.notifier).setUser(user);
        if (mounted) context.go('/learn');
      } finally {
        if (mounted) setState(() => _loading = false);
      }
    } else {
      OnboardingLanguageScreen.pendingChoice = lang.code;
      if (mounted) context.push('/onboarding/goal');
    }
  }

  @override
  Widget build(BuildContext context) {
    final query = _search.text.trim().toLowerCase();
    final list = kLanguages.where((l) =>
        query.isEmpty || l.name.toLowerCase().contains(query) || l.nativeName.toLowerCase().contains(query));

    return Scaffold(
      backgroundColor: LumoraColors.cream,
      appBar: AppBar(
        backgroundColor: LumoraColors.cream,
        foregroundColor: LumoraColors.ink,
        elevation: 0,
        leading: widget.adding ? const BackButton() : null,
        automaticallyImplyLeading: widget.adding,
        title: const Text('Choose a language', style: TextStyle(fontWeight: FontWeight.w800)),
      ),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(16),
            child: TextField(
              controller: _search,
              onChanged: (_) => setState(() {}),
              decoration: const InputDecoration(
                hintText: 'Search languages',
                prefixIcon: Icon(Icons.search),
              ),
            ),
          ),
          if (_loading) const LinearProgressIndicator(minHeight: 2),
          Expanded(
            child: ListView.separated(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              itemCount: list.length,
              separatorBuilder: (_, _) => const SizedBox(height: 8),
              itemBuilder: (context, i) {
                final lang = list.elementAt(i);
                final selected = _picked == lang.code;
                return Opacity(
                  opacity: lang.available ? 1 : 0.5,
                  child: Material(
                    color: selected ? LumoraColors.purpleLight : Colors.white,
                    borderRadius: BorderRadius.circular(LumoraRadii.lg),
                    child: InkWell(
                      borderRadius: BorderRadius.circular(LumoraRadii.lg),
                      onTap: lang.available ? () => _choose(lang) : null,
                      child: Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                        child: Row(
                          children: [
                            Text(lang.flag, style: const TextStyle(fontSize: 28)),
                            const SizedBox(width: 14),
                            Expanded(
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Text(lang.name, style: const TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.ink)),
                                  Text(lang.nativeName, style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
                                ],
                              ),
                            ),
                            if (!lang.available)
                              Container(
                                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                                decoration: BoxDecoration(
                                  color: LumoraColors.gray100,
                                  borderRadius: BorderRadius.circular(LumoraRadii.full),
                                ),
                                child: const Text('Soon', style: TextStyle(fontSize: 11, fontWeight: FontWeight.w700, color: LumoraColors.slatey)),
                              )
                            else
                              const Icon(Icons.chevron_right, color: LumoraColors.gray300),
                          ],
                        ),
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
