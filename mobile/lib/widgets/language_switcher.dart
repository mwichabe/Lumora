import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../core/languages.dart';
import '../core/network/api_client.dart';
import '../core/theme/colors.dart';
import '../core/theme/radii.dart';
import '../providers/auth_provider.dart';

class LanguageSwitcher extends ConsumerStatefulWidget {
  final VoidCallback? onChanged;
  const LanguageSwitcher({super.key, this.onChanged});

  @override
  ConsumerState<LanguageSwitcher> createState() => _LanguageSwitcherState();
}

class _LanguageSwitcherState extends ConsumerState<LanguageSwitcher> {
  List<String> _languages = [];
  bool _busy = false;

  @override
  void initState() {
    super.initState();
    ApiClient.instance.enrollments().then((r) {
      if (mounted) setState(() => _languages = r.$1);
    }).catchError((_) {});
  }

  Future<void> _pick(String code, String active) async {
    if (code == active || _busy) return;
    setState(() => _busy = true);
    try {
      final (languages, _, user) = await ApiClient.instance.switchLanguage(code);
      ref.read(authProvider.notifier).setUser(user);
      setState(() => _languages = languages);
      widget.onChanged?.call();
    } catch (_) {
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final active = ref.watch(authProvider).user?.targetLanguage ?? '';
    final meta = languageMeta(active);

    return PopupMenuButton<String>(
      color: Colors.white,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(LumoraRadii.lg)),
      onSelected: (v) => v == '__add' ? context.push('/onboarding/language?add=1') : _pick(v, active),
      itemBuilder: (context) => [
        const PopupMenuItem(enabled: false, child: Text('MY LANGUAGES', style: TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.gray500))),
        for (final code in _languages)
          PopupMenuItem(
            value: code,
            child: Row(children: [
              Text(languageMeta(code).flag, style: const TextStyle(fontSize: 18)),
              const SizedBox(width: 10),
              Expanded(child: Text(languageMeta(code).name)),
              if (code == active) const Icon(Icons.check, size: 16, color: LumoraColors.purple),
            ]),
          ),
        const PopupMenuDivider(),
        const PopupMenuItem(
          value: '__add',
          child: Row(children: [
            Icon(Icons.add, size: 18, color: LumoraColors.purple),
            SizedBox(width: 10),
            Text('Add a language', style: TextStyle(color: LumoraColors.purple, fontWeight: FontWeight.w700)),
          ]),
        ),
      ],
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
        decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.15), borderRadius: BorderRadius.circular(LumoraRadii.full)),
        child: Row(mainAxisSize: MainAxisSize.min, children: [
          Text(meta.flag, style: const TextStyle(fontSize: 16)),
          const SizedBox(width: 6),
          Text(meta.name, style: const TextStyle(color: Colors.white, fontWeight: FontWeight.w700, fontSize: 13)),
          const Icon(Icons.keyboard_arrow_down, color: Colors.white, size: 16),
        ]),
      ),
    );
  }
}
