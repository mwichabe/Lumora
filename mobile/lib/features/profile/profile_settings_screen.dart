import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:image_picker/image_picker.dart';

import '../../core/network/api_client.dart';
import '../../core/network/api_exception.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../core/theme/shadows.dart';
import '../../providers/auth_provider.dart';
import '../../widgets/avatar.dart';
import '../../widgets/confirm_dialog.dart';
import '../../widgets/lumora_button.dart';

const _kGoals = [('Casual', 10, '5 min/day'), ('Regular', 20, '10 min/day'), ('Serious', 30, '15 min/day'), ('Intense', 50, '20+ min/day')];

class ProfileSettingsScreen extends ConsumerStatefulWidget {
  const ProfileSettingsScreen({super.key});

  @override
  ConsumerState<ProfileSettingsScreen> createState() => _ProfileSettingsScreenState();
}

class _ProfileSettingsScreenState extends ConsumerState<ProfileSettingsScreen> {
  late final TextEditingController _name;
  int _goal = 20;
  bool _busy = false;
  bool _photoBusy = false;
  String? _msg;
  bool _msgOk = false;

  final _current = TextEditingController();
  final _next = TextEditingController();
  final _confirm = TextEditingController();
  bool _pwBusy = false;
  String? _pwMsg;
  bool _pwOk = false;

  bool _dangerOpen = false;
  final _dangerPassword = TextEditingController();
  bool _dangerBusy = false;
  String? _dangerErr;

  @override
  void initState() {
    super.initState();
    final user = ref.read(authProvider).user!;
    _name = TextEditingController(text: user.name);
    _goal = user.dailyGoalXp;
  }

  bool get _changed {
    final user = ref.read(authProvider).user!;
    return _name.text.trim() != user.name || _goal != user.dailyGoalXp;
  }

  Future<void> _saveProfile() async {
    setState(() { _busy = true; _msg = null; });
    try {
      final user = await ApiClient.instance.updateProfile(name: _name.text.trim(), dailyGoalXp: _goal);
      ref.read(authProvider.notifier).setUser(user);
      setState(() { _msg = 'Profile updated.'; _msgOk = true; });
    } on ApiException catch (e) {
      setState(() { _msg = e.message; _msgOk = false; });
    } finally {
      if (mounted) setState(() => _busy = false);
    }
  }

  Future<void> _pickPhoto() async {
    final picker = ImagePicker();
    final file = await picker.pickImage(source: ImageSource.gallery, imageQuality: 90);
    if (file == null) return;
    final bytes = await file.readAsBytes();
    setState(() => _photoBusy = true);
    try {
      final user = await ApiClient.instance.uploadAvatar(UploadFile(bytes, file.name));
      ref.read(authProvider.notifier).setUser(user);
    } catch (_) {
    } finally {
      if (mounted) setState(() => _photoBusy = false);
    }
  }

  Future<void> _removePhoto() async {
    setState(() => _photoBusy = true);
    try {
      final user = await ApiClient.instance.removeAvatar();
      ref.read(authProvider.notifier).setUser(user);
    } catch (_) {
    } finally {
      if (mounted) setState(() => _photoBusy = false);
    }
  }

  Future<void> _savePassword() async {
    if (_current.text.isEmpty || _next.text.length < 6 || _next.text != _confirm.text) {
      setState(() {
        _pwOk = false;
        _pwMsg = _next.text != _confirm.text ? "Passwords don't match." : 'New password must be 6+ characters.';
      });
      return;
    }
    setState(() { _pwBusy = true; _pwMsg = null; });
    try {
      await ApiClient.instance.changePassword(_current.text, _next.text);
      setState(() { _pwMsg = 'Password updated.'; _pwOk = true; });
      _current.clear(); _next.clear(); _confirm.clear();
    } on ApiException catch (e) {
      setState(() { _pwMsg = e.message; _pwOk = false; });
    } finally {
      if (mounted) setState(() => _pwBusy = false);
    }
  }

  Future<void> _deleteAccount() async {
    if (_dangerPassword.text.isEmpty) {
      setState(() => _dangerErr = 'Enter your password to confirm.');
      return;
    }
    setState(() { _dangerBusy = true; _dangerErr = null; });
    try {
      await ApiClient.instance.deleteAccount(_dangerPassword.text);
      await ref.read(authProvider.notifier).logout();
      if (mounted) context.go('/welcome');
    } on ApiException catch (e) {
      setState(() { _dangerErr = e.message; _dangerBusy = false; });
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
    final user = ref.watch(authProvider).user!;
    return Scaffold(
      backgroundColor: LumoraColors.cream,
      appBar: AppBar(
        backgroundColor: Colors.white,
        foregroundColor: LumoraColors.ink,
        elevation: 0,
        leading: BackButton(onPressed: () => context.canPop() ? context.pop() : context.go('/profile')),
        title: const Text('Account Settings', style: TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.ink)),
      ),
      body: ListView(
        padding: const EdgeInsets.all(20),
        children: [
          _Card(title: 'Profile', child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            const Text('Profile photo', style: TextStyle(color: LumoraColors.slatey, fontWeight: FontWeight.w700, fontSize: 13)),
            const SizedBox(height: 8),
            Row(children: [
              LumoraAvatar(name: user.name, avatarColor: user.avatarColor, avatarUrl: user.avatarUrl, size: 64),
              const SizedBox(width: 16),
              Wrap(spacing: 8, children: [
                LumoraButton(label: _photoBusy ? 'Working…' : (user.avatarUrl.isNotEmpty ? 'Change photo' : 'Upload photo'), loading: _photoBusy, onPressed: _pickPhoto,
                    icon: const Icon(Icons.camera_alt, size: 16, color: Colors.white)),
                if (user.avatarUrl.isNotEmpty)
                  LumoraButton(label: 'Remove', variant: LumoraButtonVariant.outline, onPressed: _photoBusy ? null : _removePhoto),
              ]),
            ]),
            const SizedBox(height: 20),
            const Text('Display name', style: TextStyle(color: LumoraColors.slatey, fontWeight: FontWeight.w700, fontSize: 13)),
            const SizedBox(height: 6),
            TextField(controller: _name, onChanged: (_) => setState(() {})),
            const SizedBox(height: 16),
            const Text('Daily goal', style: TextStyle(color: LumoraColors.slatey, fontWeight: FontWeight.w700, fontSize: 13)),
            const SizedBox(height: 8),
            GridView.count(
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              crossAxisCount: 2,
              mainAxisSpacing: 8, crossAxisSpacing: 8, childAspectRatio: 2.2,
              children: [
                for (final g in _kGoals)
                  InkWell(
                    onTap: () => setState(() => _goal = g.$2),
                    child: Container(
                      decoration: BoxDecoration(
                        color: _goal == g.$2 ? LumoraColors.purpleLight : Colors.white,
                        border: Border.all(color: _goal == g.$2 ? LumoraColors.purple : LumoraColors.gray100, width: 2),
                        borderRadius: BorderRadius.circular(LumoraRadii.lg),
                      ),
                      child: Column(mainAxisAlignment: MainAxisAlignment.center, children: [
                        Text(g.$1, style: const TextStyle(fontWeight: FontWeight.w800)),
                        Text(g.$3, style: const TextStyle(color: LumoraColors.slatey, fontSize: 11)),
                      ]),
                    ),
                  ),
              ],
            ),
            const SizedBox(height: 16),
            LumoraButton(label: 'Save changes', full: true, loading: _busy, onPressed: _changed ? _saveProfile : null),
            if (_msg != null) Padding(padding: const EdgeInsets.only(top: 8), child: _Notice(text: _msg!, ok: _msgOk)),
          ])),
          const SizedBox(height: 16),
          _Card(title: 'Password', child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            TextField(controller: _current, obscureText: true, decoration: const InputDecoration(labelText: 'Current password')),
            const SizedBox(height: 12),
            TextField(controller: _next, obscureText: true, decoration: const InputDecoration(labelText: 'New password')),
            const SizedBox(height: 12),
            TextField(controller: _confirm, obscureText: true, decoration: const InputDecoration(labelText: 'Confirm new password')),
            const SizedBox(height: 16),
            LumoraButton(label: 'Update password', full: true, loading: _pwBusy, onPressed: _savePassword),
            if (_pwMsg != null) Padding(padding: const EdgeInsets.only(top: 8), child: _Notice(text: _pwMsg!, ok: _pwOk)),
          ])),
          const SizedBox(height: 16),
          Container(
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(color: LumoraColors.coralLight.withValues(alpha: 0.6), border: Border.all(color: LumoraColors.coral.withValues(alpha: 0.3), width: 2), borderRadius: BorderRadius.circular(LumoraRadii.xl)),
            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              const Row(children: [Icon(Icons.warning_amber_rounded, color: LumoraColors.coral, size: 18), SizedBox(width: 8), Text('Danger zone', style: TextStyle(color: LumoraColors.coral, fontWeight: FontWeight.w800, fontSize: 16))]),
              const SizedBox(height: 4),
              const Text('Deleting your account is permanent and removes all your progress, streaks and data. This cannot be undone.', style: TextStyle(color: LumoraColors.slatey, fontSize: 13)),
              const SizedBox(height: 12),
              if (!_dangerOpen)
                LumoraButton(label: 'Delete account', variant: LumoraButtonVariant.danger, onPressed: () => setState(() => _dangerOpen = true))
              else
                Column(children: [
                  TextField(controller: _dangerPassword, obscureText: true, decoration: const InputDecoration(hintText: 'Enter your password to confirm')),
                  if (_dangerErr != null) Padding(padding: const EdgeInsets.only(top: 6), child: Text(_dangerErr!, style: const TextStyle(color: LumoraColors.coral, fontSize: 12))),
                  const SizedBox(height: 12),
                  Row(children: [
                    Expanded(child: LumoraButton(label: 'Cancel', variant: LumoraButtonVariant.outline, full: true, onPressed: () => setState(() { _dangerOpen = false; _dangerPassword.clear(); _dangerErr = null; }))),
                    const SizedBox(width: 8),
                    Expanded(child: LumoraButton(label: 'Confirm delete', variant: LumoraButtonVariant.danger, full: true, loading: _dangerBusy, onPressed: _deleteAccount)),
                  ]),
                ]),
            ]),
          ),
          const SizedBox(height: 16),
          Material(
            color: Colors.white,
            borderRadius: BorderRadius.circular(LumoraRadii.xl),
            child: InkWell(
              borderRadius: BorderRadius.circular(LumoraRadii.xl),
              onTap: _confirmSignOut,
              child: Container(
                padding: const EdgeInsets.all(16),
                decoration: BoxDecoration(borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
                child: const Row(mainAxisAlignment: MainAxisAlignment.center, children: [
                  Icon(Icons.logout, color: LumoraColors.slatey, size: 18),
                  SizedBox(width: 8),
                  Text('Sign out', style: TextStyle(color: LumoraColors.slatey, fontWeight: FontWeight.w800)),
                ]),
              ),
            ),
          ),
          const SizedBox(height: 32),
        ],
      ),
    );
  }
}

class _Card extends StatelessWidget {
  final String title;
  final Widget child;
  const _Card({required this.title, required this.child});
  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Text(title, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w800)),
        const SizedBox(height: 12),
        child,
      ]),
    );
  }
}

class _Notice extends StatelessWidget {
  final String text;
  final bool ok;
  const _Notice({required this.text, required this.ok});
  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      decoration: BoxDecoration(color: ok ? LumoraColors.tealLight : LumoraColors.coralLight, borderRadius: BorderRadius.circular(LumoraRadii.sm)),
      child: Text(text, style: TextStyle(color: ok ? LumoraColors.teal : LumoraColors.coral, fontWeight: FontWeight.w700, fontSize: 12)),
    );
  }
}
