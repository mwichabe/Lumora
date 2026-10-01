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
  String? _photoMsg;
  bool _photoOk = false;
  bool _showPasswords = false;
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
    setState(() { _photoBusy = true; _photoMsg = null; });
    try {
      final user = await ApiClient.instance.uploadAvatar(UploadFile(bytes, file.name));
      ref.read(authProvider.notifier).setUser(user);
      setState(() { _photoMsg = 'Photo updated.'; _photoOk = true; });
    } on ApiException catch (e) {
      setState(() { _photoMsg = e.message; _photoOk = false; });
    } finally {
      if (mounted) setState(() => _photoBusy = false);
    }
  }

  Future<void> _removePhoto() async {
    final confirm = await showLumoraConfirmDialog(
      context,
      title: 'Remove photo?',
      message: 'Your profile will show your initial instead.',
      confirmLabel: 'Remove',
      danger: true,
    );
    if (!confirm) return;
    setState(() { _photoBusy = true; _photoMsg = null; });
    try {
      final user = await ApiClient.instance.removeAvatar();
      ref.read(authProvider.notifier).setUser(user);
      setState(() { _photoMsg = 'Photo removed.'; _photoOk = true; });
    } on ApiException catch (e) {
      setState(() { _photoMsg = e.message; _photoOk = false; });
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
    final hasPhoto = user.avatarUrl.isNotEmpty;
    return Scaffold(
      backgroundColor: LumoraColors.cream,
      appBar: AppBar(
        backgroundColor: Colors.white,
        foregroundColor: LumoraColors.ink,
        surfaceTintColor: Colors.white,
        elevation: 0,
        leading: BackButton(onPressed: () => context.canPop() ? context.pop() : context.go('/profile')),
        title: const Text('Account Settings', style: TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.ink)),
      ),
      body: ListView(
        padding: const EdgeInsets.fromLTRB(16, 16, 16, 32),
        children: [
          // ---- Profile -------------------------------------------------------
          _Section(
            icon: Icons.person_rounded,
            tint: LumoraColors.purple,
            title: 'Profile',
            subtitle: 'How other learners see you',
            child: Column(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
              Row(children: [
                GestureDetector(
                  onTap: _photoBusy ? null : _pickPhoto,
                  child: Stack(clipBehavior: Clip.none, children: [
                    Container(
                      padding: const EdgeInsets.all(3),
                      decoration: const BoxDecoration(shape: BoxShape.circle, color: LumoraColors.purpleLight),
                      child: LumoraAvatar(name: user.name, avatarColor: user.avatarColor, avatarUrl: user.avatarUrl, size: 72),
                    ),
                    Positioned(
                      right: -2,
                      bottom: -2,
                      child: Container(
                        width: 28,
                        height: 28,
                        decoration: BoxDecoration(color: LumoraColors.purple, shape: BoxShape.circle, border: Border.all(color: Colors.white, width: 2)),
                        child: _photoBusy
                            ? const Padding(padding: EdgeInsets.all(6), child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white))
                            : const Icon(Icons.camera_alt_rounded, size: 14, color: Colors.white),
                      ),
                    ),
                  ]),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                    Text(user.name, maxLines: 1, overflow: TextOverflow.ellipsis,
                        style: const TextStyle(fontSize: 17, fontWeight: FontWeight.w800, color: LumoraColors.ink)),
                    const SizedBox(height: 2),
                    Text(user.email, maxLines: 1, overflow: TextOverflow.ellipsis, style: const TextStyle(fontSize: 13, color: LumoraColors.slatey)),
                  ]),
                ),
              ]),
              const SizedBox(height: 16),
              // Both photo actions share the row, so neither can be pushed off
              // the card on a narrow phone.
              Row(children: [
                Expanded(
                  child: _PillButton(
                    icon: Icons.photo_library_rounded,
                    label: hasPhoto ? 'Change photo' : 'Upload photo',
                    filled: true,
                    color: LumoraColors.purple,
                    onTap: _photoBusy ? null : _pickPhoto,
                  ),
                ),
                if (hasPhoto) ...[
                  const SizedBox(width: 10),
                  Expanded(
                    child: _PillButton(
                      icon: Icons.delete_outline_rounded,
                      label: 'Remove',
                      color: LumoraColors.coral,
                      onTap: _photoBusy ? null : _removePhoto,
                    ),
                  ),
                ],
              ]),
              if (_photoMsg != null) Padding(padding: const EdgeInsets.only(top: 10), child: _Notice(text: _photoMsg!, ok: _photoOk)),
              const SizedBox(height: 20),
              const _FieldLabel('Display name'),
              TextField(
                controller: _name,
                textCapitalization: TextCapitalization.words,
                onChanged: (_) => setState(() {}),
                decoration: const InputDecoration(prefixIcon: Icon(Icons.badge_outlined, color: LumoraColors.gray500)),
              ),
              const SizedBox(height: 18),
              const _FieldLabel('Daily goal'),
              GridView.count(
                shrinkWrap: true,
                physics: const NeverScrollableScrollPhysics(),
                crossAxisCount: 2,
                mainAxisSpacing: 10,
                crossAxisSpacing: 10,
                childAspectRatio: 2.4,
                children: [
                  for (final g in _kGoals)
                    _GoalTile(
                      label: g.$1,
                      detail: g.$3,
                      selected: _goal == g.$2,
                      onTap: () => setState(() => _goal = g.$2),
                    ),
                ],
              ),
              const SizedBox(height: 18),
              LumoraButton(label: 'Save changes', full: true, loading: _busy, onPressed: _changed ? _saveProfile : null),
              if (_msg != null) Padding(padding: const EdgeInsets.only(top: 10), child: _Notice(text: _msg!, ok: _msgOk)),
            ]),
          ),
          const SizedBox(height: 16),

          // ---- Password ------------------------------------------------------
          _Section(
            icon: Icons.lock_rounded,
            tint: LumoraColors.teal,
            title: 'Password',
            subtitle: 'At least 6 characters',
            trailing: TextButton.icon(
              onPressed: () => setState(() => _showPasswords = !_showPasswords),
              icon: Icon(_showPasswords ? Icons.visibility_off_rounded : Icons.visibility_rounded, size: 18),
              label: Text(_showPasswords ? 'Hide' : 'Show'),
            ),
            child: Column(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
              const _FieldLabel('Current password'),
              TextField(controller: _current, obscureText: !_showPasswords, decoration: _passwordDecoration()),
              const SizedBox(height: 14),
              const _FieldLabel('New password'),
              TextField(controller: _next, obscureText: !_showPasswords, decoration: _passwordDecoration()),
              const SizedBox(height: 14),
              const _FieldLabel('Confirm new password'),
              TextField(controller: _confirm, obscureText: !_showPasswords, decoration: _passwordDecoration()),
              const SizedBox(height: 18),
              LumoraButton(label: 'Update password', full: true, loading: _pwBusy, onPressed: _savePassword),
              if (_pwMsg != null) Padding(padding: const EdgeInsets.only(top: 10), child: _Notice(text: _pwMsg!, ok: _pwOk)),
            ]),
          ),
          const SizedBox(height: 16),

          // ---- Sign out ------------------------------------------------------
          Material(
            color: Colors.white,
            borderRadius: BorderRadius.circular(LumoraRadii.xl),
            child: InkWell(
              borderRadius: BorderRadius.circular(LumoraRadii.xl),
              onTap: _confirmSignOut,
              child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
                child: const Row(children: [
                  _SectionIcon(icon: Icons.logout_rounded, tint: LumoraColors.slatey),
                  SizedBox(width: 12),
                  Expanded(child: Text('Sign out', style: TextStyle(color: LumoraColors.ink, fontWeight: FontWeight.w800, fontSize: 15))),
                  Icon(Icons.chevron_right_rounded, color: LumoraColors.gray300),
                ]),
              ),
            ),
          ),
          const SizedBox(height: 16),

          // ---- Danger zone ---------------------------------------------------
          Container(
            padding: const EdgeInsets.all(18),
            decoration: BoxDecoration(
              color: Colors.white,
              border: Border.all(color: LumoraColors.coral.withValues(alpha: 0.35), width: 1.5),
              borderRadius: BorderRadius.circular(LumoraRadii.xl),
            ),
            child: Column(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
              const Row(children: [
                _SectionIcon(icon: Icons.warning_amber_rounded, tint: LumoraColors.coral),
                SizedBox(width: 12),
                Expanded(
                  child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                    Text('Delete account', style: TextStyle(color: LumoraColors.coral, fontWeight: FontWeight.w800, fontSize: 16)),
                    Text('Permanent — this cannot be undone', style: TextStyle(color: LumoraColors.slatey, fontSize: 12)),
                  ]),
                ),
              ]),
              const SizedBox(height: 12),
              const Text('Deleting your account removes all your progress, streaks, certificates and data.',
                  style: TextStyle(color: LumoraColors.slatey, fontSize: 13, height: 1.4)),
              const SizedBox(height: 14),
              if (!_dangerOpen)
                LumoraButton(label: 'Delete my account', variant: LumoraButtonVariant.danger, full: true, onPressed: () => setState(() => _dangerOpen = true))
              else ...[
                const _FieldLabel('Enter your password to confirm'),
                TextField(controller: _dangerPassword, obscureText: true, decoration: _passwordDecoration()),
                if (_dangerErr != null) Padding(padding: const EdgeInsets.only(top: 10), child: _Notice(text: _dangerErr!, ok: false)),
                const SizedBox(height: 14),
                Row(children: [
                  Expanded(
                    child: LumoraButton(
                      label: 'Cancel',
                      variant: LumoraButtonVariant.outline,
                      full: true,
                      onPressed: () => setState(() { _dangerOpen = false; _dangerPassword.clear(); _dangerErr = null; }),
                    ),
                  ),
                  const SizedBox(width: 10),
                  Expanded(child: LumoraButton(label: 'Delete', variant: LumoraButtonVariant.danger, full: true, loading: _dangerBusy, onPressed: _deleteAccount)),
                ]),
              ],
            ]),
          ),
        ],
      ),
    );
  }

  InputDecoration _passwordDecoration() =>
      const InputDecoration(prefixIcon: Icon(Icons.lock_outline_rounded, color: LumoraColors.gray500));
}

/// A white settings card with an icon, title and subtitle above its content.
class _Section extends StatelessWidget {
  final IconData icon;
  final Color tint;
  final String title;
  final String subtitle;
  final Widget? trailing;
  final Widget child;
  const _Section({required this.icon, required this.tint, required this.title, required this.subtitle, required this.child, this.trailing});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(18),
      decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
      child: Column(crossAxisAlignment: CrossAxisAlignment.stretch, children: [
        Row(children: [
          _SectionIcon(icon: icon, tint: tint),
          const SizedBox(width: 12),
          Expanded(
            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Text(title, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w800, color: LumoraColors.ink)),
              Text(subtitle, style: const TextStyle(fontSize: 12, color: LumoraColors.slatey)),
            ]),
          ),
          ?trailing,
        ]),
        const Padding(padding: EdgeInsets.symmetric(vertical: 14), child: Divider(height: 1)),
        child,
      ]),
    );
  }
}

class _SectionIcon extends StatelessWidget {
  final IconData icon;
  final Color tint;
  const _SectionIcon({required this.icon, required this.tint});

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 38,
      height: 38,
      decoration: BoxDecoration(color: tint.withValues(alpha: 0.12), borderRadius: BorderRadius.circular(LumoraRadii.md)),
      child: Icon(icon, color: tint, size: 20),
    );
  }
}

class _FieldLabel extends StatelessWidget {
  final String text;
  const _FieldLabel(this.text);
  @override
  Widget build(BuildContext context) => Padding(
        padding: const EdgeInsets.only(bottom: 6, left: 2),
        child: Text(text, style: const TextStyle(color: LumoraColors.slatey, fontWeight: FontWeight.w700, fontSize: 13)),
      );
}

/// A compact pill button that shrinks its label instead of overflowing.
class _PillButton extends StatelessWidget {
  final IconData icon;
  final String label;
  final Color color;
  final bool filled;
  final VoidCallback? onTap;
  const _PillButton({required this.icon, required this.label, required this.color, this.filled = false, this.onTap});

  @override
  Widget build(BuildContext context) {
    final fg = filled ? Colors.white : color;
    return Opacity(
      opacity: onTap == null ? 0.5 : 1,
      child: Material(
        color: filled ? color : color.withValues(alpha: 0.08),
        shape: StadiumBorder(side: filled ? BorderSide.none : BorderSide(color: color, width: 1.5)),
        child: InkWell(
          customBorder: const StadiumBorder(),
          onTap: onTap,
          child: SizedBox(
            height: 46,
            child: Row(mainAxisAlignment: MainAxisAlignment.center, children: [
              Icon(icon, size: 18, color: fg),
              const SizedBox(width: 6),
              Flexible(
                child: Text(label, maxLines: 1, overflow: TextOverflow.ellipsis,
                    style: TextStyle(color: fg, fontWeight: FontWeight.w800, fontSize: 14)),
              ),
            ]),
          ),
        ),
      ),
    );
  }
}

class _GoalTile extends StatelessWidget {
  final String label;
  final String detail;
  final bool selected;
  final VoidCallback onTap;
  const _GoalTile({required this.label, required this.detail, required this.selected, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return Material(
      color: selected ? LumoraColors.purpleLight : Colors.white,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(LumoraRadii.lg),
        side: BorderSide(color: selected ? LumoraColors.purple : LumoraColors.gray100, width: 2),
      ),
      child: InkWell(
        borderRadius: BorderRadius.circular(LumoraRadii.lg),
        onTap: onTap,
        child: Stack(children: [
          Center(
            child: Column(mainAxisSize: MainAxisSize.min, children: [
              Text(label, style: TextStyle(fontWeight: FontWeight.w800, color: selected ? LumoraColors.purple : LumoraColors.ink)),
              Text(detail, style: const TextStyle(color: LumoraColors.slatey, fontSize: 11)),
            ]),
          ),
          if (selected)
            const Positioned(top: 6, right: 6, child: Icon(Icons.check_circle_rounded, size: 16, color: LumoraColors.purple)),
        ]),
      ),
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
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
      decoration: BoxDecoration(color: ok ? LumoraColors.tealLight : LumoraColors.coralLight, borderRadius: BorderRadius.circular(LumoraRadii.md)),
      child: Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Icon(ok ? Icons.check_circle_rounded : Icons.error_outline_rounded, size: 16, color: ok ? LumoraColors.teal : LumoraColors.coral),
        const SizedBox(width: 8),
        Expanded(child: Text(text, style: TextStyle(color: ok ? const Color(0xFF0B6B5E) : LumoraColors.coral, fontWeight: FontWeight.w700, fontSize: 12.5))),
      ]),
    );
  }
}
