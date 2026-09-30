import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:image_picker/image_picker.dart';

import '../../core/languages.dart';
import '../../core/network/api_client.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../core/theme/shadows.dart';
import '../../models/character.dart';
import '../../providers/auth_provider.dart';
import '../../widgets/avatar.dart';
import '../../widgets/character_bio_modal.dart';
import '../../widgets/confirm_dialog.dart';

class ProfileScreen extends ConsumerStatefulWidget {
  const ProfileScreen({super.key});

  @override
  ConsumerState<ProfileScreen> createState() => _ProfileScreenState();
}

class _ProfileScreenState extends ConsumerState<ProfileScreen> {
  List<CharacterWithFriendship> _characters = [];
  List<String> _languages = [];
  int _certCount = 0;
  bool _switching = false;
  bool _uploadingAvatar = false;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      _characters = await ApiClient.instance.characters();
    } catch (_) {}
    try {
      final (langs, _) = await ApiClient.instance.enrollments();
      _languages = langs;
    } catch (_) {}
    try {
      _certCount = (await ApiClient.instance.certificates()).length;
    } catch (_) {}
    if (mounted) setState(() {});
  }

  Future<void> _switchLanguage(String code) async {
    final user = ref.read(authProvider).user;
    if (code == user?.targetLanguage || _switching) return;
    setState(() => _switching = true);
    try {
      final (langs, _, updatedUser) = await ApiClient.instance.switchLanguage(code);
      ref.read(authProvider.notifier).setUser(updatedUser);
      setState(() => _languages = langs);
    } catch (_) {
    } finally {
      if (mounted) setState(() => _switching = false);
    }
  }

  Future<void> _pickAvatar() async {
    final picker = ImagePicker();
    final file = await picker.pickImage(source: ImageSource.gallery, imageQuality: 90);
    if (file == null) return;
    final bytes = await file.readAsBytes();
    setState(() => _uploadingAvatar = true);
    try {
      final user = await ApiClient.instance.uploadAvatar(UploadFile(bytes, file.name));
      ref.read(authProvider.notifier).setUser(user);
    } catch (_) {
    } finally {
      if (mounted) setState(() => _uploadingAvatar = false);
    }
  }

  Future<void> _confirmSignOut() async {
    final confirm = await showLumoraConfirmDialog(
      context,
      title: 'Sign out?',
      message: "You'll need to sign in again to continue learning.",
      confirmLabel: 'Sign out',
      danger: true,
    );
    if (confirm) {
      await ref.read(authProvider.notifier).logout();
      if (mounted) context.go('/welcome');
    }
  }

  @override
  Widget build(BuildContext context) {
    final user = ref.watch(authProvider).user;
    if (user == null) return const SizedBox.shrink();
    final fluencyPct = (user.fluencyScore / 1000).clamp(0, 1).toDouble();
    final lessonsDone = (user.xp / 10).floor().clamp(0, 1 << 31);

    return Scaffold(
      backgroundColor: LumoraColors.cream,
      body: RefreshIndicator(
        onRefresh: _load,
        child: ListView(
          padding: EdgeInsets.zero,
          children: [
            Container(
              padding: const EdgeInsets.fromLTRB(24, 56, 24, 40),
              decoration: const BoxDecoration(
                gradient: LinearGradient(begin: Alignment.topCenter, end: Alignment.bottomCenter, colors: [LumoraColors.purple, LumoraColors.purpleDark]),
                borderRadius: BorderRadius.vertical(bottom: Radius.circular(32)),
              ),
              child: Row(children: [
                GestureDetector(
                  onTap: _pickAvatar,
                  child: Stack(children: [
                    Container(
                      decoration: BoxDecoration(shape: BoxShape.circle, border: Border.all(color: Colors.white30, width: 4)),
                      child: LumoraAvatar(name: user.name, avatarColor: user.avatarColor, avatarUrl: user.avatarUrl, size: 80),
                    ),
                    Positioned(
                      right: 0, bottom: 0,
                      child: Container(
                        width: 28, height: 28,
                        decoration: BoxDecoration(color: Colors.white, shape: BoxShape.circle, border: Border.all(color: LumoraColors.purple, width: 2)),
                        child: _uploadingAvatar
                            ? const Padding(padding: EdgeInsets.all(6), child: CircularProgressIndicator(strokeWidth: 2))
                            : const Icon(Icons.camera_alt, size: 14, color: LumoraColors.purple),
                      ),
                    ),
                  ]),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                    Text(user.name, maxLines: 1, overflow: TextOverflow.ellipsis, style: const TextStyle(color: Colors.white, fontSize: 22, fontWeight: FontWeight.w800)),
                    Text(user.email, maxLines: 1, overflow: TextOverflow.ellipsis, style: const TextStyle(color: LumoraColors.purpleLight, fontSize: 12)),
                    const SizedBox(height: 6),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 2),
                      decoration: BoxDecoration(color: LumoraColors.amber, borderRadius: BorderRadius.circular(LumoraRadii.full)),
                      child: Text('${user.cefrLevel.isEmpty ? "A1" : user.cefrLevel} · ${user.levelName.isEmpty ? "Spark" : user.levelName}',
                          style: const TextStyle(color: LumoraColors.ink, fontSize: 11, fontWeight: FontWeight.w800)),
                    ),
                  ]),
                ),
              ]),
            ),
            Transform.translate(
              offset: const Offset(0, -24),
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20),
                child: Container(
                  padding: const EdgeInsets.all(20),
                  decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
                  child: Row(children: [
                    _FluencyRing(pct: fluencyPct, score: user.fluencyScore),
                    const SizedBox(width: 16),
                    Expanded(
                      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                        const Text('Fluency Score', style: TextStyle(fontWeight: FontWeight.w800, fontSize: 16)),
                        Text(user.targetLanguage.isNotEmpty ? 'Your ${languageName(user.targetLanguage)} mastery, out of 1000.' : 'Your mastery, out of 1000.',
                            style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
                      ]),
                    ),
                  ]),
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20),
              child: GridView.count(
                shrinkWrap: true,
                physics: const NeverScrollableScrollPhysics(),
                crossAxisCount: 2,
                mainAxisSpacing: 12,
                crossAxisSpacing: 12,
                childAspectRatio: 2.6,
                children: [
                  _StatTile(icon: Icons.local_fire_department, tint: LumoraColors.coral, label: 'Day Streak', value: user.streak),
                  _StatTile(icon: Icons.bolt, tint: LumoraColors.amber, label: 'Total XP', value: user.xp),
                  _StatTile(icon: Icons.menu_book, tint: LumoraColors.purple, label: 'Lessons', value: lessonsDone),
                  _StatTile(icon: Icons.public, tint: LumoraColors.teal, label: 'Languages', value: _languages.length.clamp(user.targetLanguage.isNotEmpty ? 1 : 0, 999)),
                ],
              ),
            ),
            const SizedBox(height: 24),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20),
              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Row(children: [
                  const Expanded(child: Text('My Languages', style: TextStyle(fontSize: 18, fontWeight: FontWeight.w800))),
                  TextButton.icon(onPressed: () => context.push('/onboarding/language?add=1'), icon: const Icon(Icons.add, size: 16), label: const Text('Add')),
                ]),
                Container(
                  decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
                  child: Column(children: [
                    if (_languages.isEmpty) const Padding(padding: EdgeInsets.all(16), child: Text('No languages yet. Tap "Add" to start a course.', style: TextStyle(color: LumoraColors.slatey))),
                    for (var i = 0; i < _languages.length; i++)
                      InkWell(
                        onTap: () => _switchLanguage(_languages[i]),
                        child: Container(
                          decoration: BoxDecoration(border: i > 0 ? const Border(top: BorderSide(color: LumoraColors.gray100)) : null),
                          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                          child: Row(children: [
                            Text(languageMeta(_languages[i]).flag, style: const TextStyle(fontSize: 24)),
                            const SizedBox(width: 12),
                            Expanded(
                              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                                Text(languageMeta(_languages[i]).name, style: const TextStyle(fontWeight: FontWeight.w800)),
                                Text(_languages[i] == user.targetLanguage ? 'Active course' : 'Tap to switch', style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
                              ]),
                            ),
                            if (_languages[i] == user.targetLanguage)
                              Container(padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4), decoration: BoxDecoration(color: LumoraColors.purple, borderRadius: BorderRadius.circular(LumoraRadii.full)),
                                  child: const Text('Active', style: TextStyle(color: Colors.white, fontSize: 11, fontWeight: FontWeight.w800)))
                            else
                              const Icon(Icons.chevron_right, color: LumoraColors.gray300),
                          ]),
                        ),
                      ),
                  ]),
                ),
              ]),
            ),
            const SizedBox(height: 24),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20),
              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Row(children: [
                  const Expanded(child: Text('Certificates', style: TextStyle(fontSize: 18, fontWeight: FontWeight.w800))),
                  TextButton.icon(onPressed: () => context.push('/exam'), icon: const Icon(Icons.school_outlined, size: 16), label: const Text('Get certified')),
                ]),
                Material(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(LumoraRadii.xl),
                  child: InkWell(
                    borderRadius: BorderRadius.circular(LumoraRadii.xl),
                    onTap: () => context.push('/certificates'),
                    child: Container(
                      padding: const EdgeInsets.all(16),
                      decoration: BoxDecoration(borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
                      child: Row(children: [
                        Container(width: 44, height: 44, decoration: BoxDecoration(color: LumoraColors.purpleLight, borderRadius: BorderRadius.circular(LumoraRadii.md)),
                            child: const Icon(Icons.workspace_premium, color: LumoraColors.purple)),
                        const SizedBox(width: 12),
                        Expanded(
                          child: Text(_certCount == 0 ? 'Take the proficiency exam to earn your first certificate.' : 'View your $_certCount certificate${_certCount == 1 ? "" : "s"}',
                              style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 13)),
                        ),
                        const Icon(Icons.chevron_right, color: LumoraColors.gray300),
                      ]),
                    ),
                  ),
                ),
              ]),
            ),
            const SizedBox(height: 24),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20),
              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                const Text('Your Companions', style: TextStyle(fontSize: 18, fontWeight: FontWeight.w800)),
                const SizedBox(height: 8),
                GridView.builder(
                  shrinkWrap: true,
                  physics: const NeverScrollableScrollPhysics(),
                  gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(crossAxisCount: 3, mainAxisSpacing: 10, crossAxisSpacing: 10, childAspectRatio: 0.8),
                  itemCount: _characters.length,
                  itemBuilder: (context, i) {
                    final c = _characters[i];
                    return Material(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(LumoraRadii.lg),
                      child: InkWell(
                        borderRadius: BorderRadius.circular(LumoraRadii.lg),
                        onTap: () => showCharacterBioModal(context, c),
                        child: Container(
                          padding: const EdgeInsets.all(8),
                          decoration: BoxDecoration(borderRadius: BorderRadius.circular(LumoraRadii.lg), boxShadow: LumoraShadows.card),
                          child: Column(mainAxisAlignment: MainAxisAlignment.center, children: [
                            CircleAvatar(radius: 24, backgroundColor: LumoraColors.purpleLight, child: Text(c.emoji, style: const TextStyle(fontSize: 22))),
                            const SizedBox(height: 6),
                            Text(c.name, style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 11)),
                            Text(c.role, style: const TextStyle(color: LumoraColors.slatey, fontSize: 9)),
                            const SizedBox(height: 4),
                            Row(mainAxisAlignment: MainAxisAlignment.center, children: [
                              for (var d = 0; d < 3; d++)
                                Container(margin: const EdgeInsets.symmetric(horizontal: 1), width: 6, height: 6,
                                    decoration: BoxDecoration(color: d < c.friendshipLevel ? LumoraColors.amber : LumoraColors.gray300, shape: BoxShape.circle)),
                            ]),
                          ]),
                        ),
                      ),
                    );
                  },
                ),
              ]),
            ),
            const SizedBox(height: 24),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20),
              child: Container(
                decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
                child: Column(children: [
                  _SettingsRow(label: 'Account Settings', onTap: () => context.push('/profile/settings')),
                  _SettingsRow(label: 'Notifications', onTap: () => context.push('/notifications')),
                  _SettingsRow(label: 'Help & Support', onTap: () => context.push('/profile/help')),
                  InkWell(
                    onTap: _confirmSignOut,
                    child: Container(
                      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
                      child: const Row(children: [
                        Icon(Icons.logout, color: LumoraColors.coral, size: 18),
                        SizedBox(width: 10),
                        Expanded(child: Text('Sign Out', style: TextStyle(color: LumoraColors.coral, fontWeight: FontWeight.w800))),
                        Icon(Icons.chevron_right, color: LumoraColors.coral),
                      ]),
                    ),
                  ),
                ]),
              ),
            ),
            const SizedBox(height: 32),
          ],
        ),
      ),
    );
  }
}

class _FluencyRing extends StatelessWidget {
  final double pct;
  final int score;
  const _FluencyRing({required this.pct, required this.score});

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: 84, height: 84,
      child: Stack(alignment: Alignment.center, children: [
        SizedBox(
          width: 84, height: 84,
          child: CircularProgressIndicator(value: pct, strokeWidth: 8, backgroundColor: LumoraColors.purpleLight, valueColor: const AlwaysStoppedAnimation(LumoraColors.purple)),
        ),
        Column(mainAxisSize: MainAxisSize.min, children: [
          Text('$score', style: const TextStyle(color: LumoraColors.purple, fontWeight: FontWeight.w800, fontSize: 18)),
          const Text('/ 1000', style: TextStyle(color: LumoraColors.slatey, fontSize: 9)),
        ]),
      ]),
    );
  }
}

class _StatTile extends StatelessWidget {
  final IconData icon;
  final Color tint;
  final String label;
  final int value;
  const _StatTile({required this.icon, required this.tint, required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
      child: Row(children: [
        Container(width: 36, height: 36, decoration: BoxDecoration(color: tint.withValues(alpha: 0.1), borderRadius: BorderRadius.circular(LumoraRadii.md)), child: Icon(icon, size: 18, color: tint)),
        const SizedBox(width: 10),
        Expanded(
          child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Text('$value', style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 16)),
            Text(label, style: const TextStyle(color: LumoraColors.slatey, fontSize: 10)),
          ]),
        ),
      ]),
    );
  }
}

class _SettingsRow extends StatelessWidget {
  final String label;
  final VoidCallback onTap;
  const _SettingsRow({required this.label, required this.onTap});
  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
        decoration: const BoxDecoration(border: Border(bottom: BorderSide(color: LumoraColors.gray100))),
        child: Row(children: [
          Expanded(child: Text(label, style: const TextStyle(fontWeight: FontWeight.w700))),
          const Icon(Icons.chevron_right, color: LumoraColors.gray300),
        ]),
      ),
    );
  }
}
