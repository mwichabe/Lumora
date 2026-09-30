import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../core/network/api_client.dart';
import '../../core/theme/colors.dart';
import '../../widgets/auth_field.dart';
import '../../widgets/fox_mascot.dart';
import '../../widgets/lumora_button.dart';

/// Always shows the same "check your email" confirmation, whether or not the
/// account exists — enumeration-safe, matching frontend/app/forgot-password.
class ForgotPasswordScreen extends StatefulWidget {
  const ForgotPasswordScreen({super.key});

  @override
  State<ForgotPasswordScreen> createState() => _ForgotPasswordScreenState();
}

class _ForgotPasswordScreenState extends State<ForgotPasswordScreen> {
  final _email = TextEditingController();
  bool _sent = false;
  bool _loading = false;

  Future<void> _submit() async {
    if (_email.text.trim().isEmpty) return;
    setState(() => _loading = true);
    try {
      await ApiClient.instance.forgotPassword(_email.text.trim());
    } catch (_) {
      // still show the confirmation — enumeration-safe by design
    } finally {
      if (mounted) {
        setState(() {
          _loading = false;
          _sent = true;
        });
      }
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
                  Row(children: [AuthBackButton(onTap: () => context.canPop() ? context.pop() : context.go('/welcome'))]),
                  const SizedBox(height: 24),
                  const Center(child: FoxMascot(size: 72)),
                  const SizedBox(height: 24),
                  AuthCard(child: _sent ? _sentContent() : _formContent()),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _sentContent() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        const AuthIconBadge(icon: Icons.mark_email_read_rounded, color: LumoraColors.teal),
        const SizedBox(height: 16),
        const Text('Check your email', style: TextStyle(fontSize: 20, fontWeight: FontWeight.w800, color: LumoraColors.ink)),
        const SizedBox(height: 8),
        Text(
          'If an account exists for ${_email.text.trim()}, we\'ve sent a link to reset your password.',
          textAlign: TextAlign.center,
          style: const TextStyle(color: LumoraColors.slatey, height: 1.4),
        ),
        const SizedBox(height: 24),
        LumoraButton(label: 'Back to sign in', full: true, onPressed: () => context.go('/welcome')),
      ],
    );
  }

  Widget _formContent() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        const AuthIconBadge(icon: Icons.lock_reset_rounded),
        const SizedBox(height: 16),
        const Text('Forgot password?', style: TextStyle(fontSize: 20, fontWeight: FontWeight.w800, color: LumoraColors.ink)),
        const SizedBox(height: 6),
        const Text("Enter your email and we'll send you a reset link.", style: TextStyle(color: LumoraColors.slatey)),
        const SizedBox(height: 24),
        const AuthFieldLabel('Email'),
        TextField(
          controller: _email,
          keyboardType: TextInputType.emailAddress,
          decoration: authFieldDecoration(icon: Icons.mail_outline_rounded, hint: 'you@example.com'),
        ),
        const SizedBox(height: 20),
        LumoraButton(label: 'Send reset link', full: true, loading: _loading, onPressed: _submit),
      ],
    );
  }
}
