import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../core/network/api_client.dart';
import '../../core/network/api_exception.dart';
import '../../core/theme/colors.dart';
import '../../widgets/auth_field.dart';
import '../../widgets/fox_mascot.dart';
import '../../widgets/lumora_button.dart';

class ResetPasswordScreen extends StatefulWidget {
  final String? token;
  const ResetPasswordScreen({super.key, this.token});

  @override
  State<ResetPasswordScreen> createState() => _ResetPasswordScreenState();
}

class _ResetPasswordScreenState extends State<ResetPasswordScreen> {
  final _password = TextEditingController();
  final _confirm = TextEditingController();
  bool _obscure = true;
  bool _loading = false;
  bool _done = false;
  String? _error;

  Future<void> _submit() async {
    if (widget.token == null || widget.token!.isEmpty) {
      setState(() => _error = 'This reset link is invalid or has expired.');
      return;
    }
    if (_password.text.length < 6) {
      setState(() => _error = 'Password must be at least 6 characters.');
      return;
    }
    if (_password.text != _confirm.text) {
      setState(() => _error = 'Passwords do not match.');
      return;
    }
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      await ApiClient.instance.resetPassword(widget.token!, _password.text);
      setState(() => _done = true);
    } on ApiException catch (e) {
      setState(() => _error = e.message);
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: LumoraColors.cream,
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.fromLTRB(24, 20, 24, 24),
          child: Center(
            child: ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 420),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  if (!_done) Row(children: [AuthBackButton(onTap: () => context.canPop() ? context.pop() : context.go('/welcome'))]),
                  const SizedBox(height: 24),
                  const Center(child: FoxMascot(size: 72)),
                  const SizedBox(height: 24),
                  AuthCard(child: _done ? _doneContent() : _formContent()),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _doneContent() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        const AuthIconBadge(icon: Icons.check_circle_rounded, color: LumoraColors.teal),
        const SizedBox(height: 16),
        const Text('Password updated', style: TextStyle(fontSize: 20, fontWeight: FontWeight.w800, color: LumoraColors.ink)),
        const SizedBox(height: 8),
        const Text('You can now log in with your new password.', textAlign: TextAlign.center, style: TextStyle(color: LumoraColors.slatey)),
        const SizedBox(height: 24),
        LumoraButton(label: 'Back to sign in', full: true, onPressed: () => context.go('/welcome')),
      ],
    );
  }

  Widget _formContent() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        const AuthIconBadge(icon: Icons.lock_outline_rounded),
        const SizedBox(height: 16),
        const Text('Set a new password', style: TextStyle(fontSize: 20, fontWeight: FontWeight.w800, color: LumoraColors.ink)),
        const SizedBox(height: 6),
        const Text('Choose a strong password you\'ll remember.', style: TextStyle(color: LumoraColors.slatey)),
        const SizedBox(height: 24),
        const AuthFieldLabel('New password'),
        TextField(
          controller: _password,
          obscureText: _obscure,
          decoration: authFieldDecoration(
            icon: Icons.lock_outline_rounded,
            hint: '••••••••',
            suffixIcon: IconButton(
              icon: Icon(_obscure ? Icons.visibility_off_outlined : Icons.visibility_outlined, size: 20, color: LumoraColors.gray500),
              onPressed: () => setState(() => _obscure = !_obscure),
            ),
          ),
        ),
        const SizedBox(height: 16),
        const AuthFieldLabel('Confirm password'),
        TextField(
          controller: _confirm,
          obscureText: _obscure,
          decoration: authFieldDecoration(icon: Icons.lock_outline_rounded, hint: '••••••••'),
        ),
        if (_error != null) ...[
          const SizedBox(height: 12),
          Container(
            width: double.infinity,
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
            decoration: BoxDecoration(color: LumoraColors.coralLight, borderRadius: BorderRadius.circular(12)),
            child: Row(
              children: [
                const Icon(Icons.error_outline_rounded, size: 16, color: LumoraColors.coral),
                const SizedBox(width: 8),
                Expanded(child: Text(_error!, style: const TextStyle(color: LumoraColors.coral, fontSize: 13, fontWeight: FontWeight.w600))),
              ],
            ),
          ),
        ],
        const SizedBox(height: 20),
        LumoraButton(label: 'Reset password', full: true, loading: _loading, onPressed: _submit),
      ],
    );
  }
}
