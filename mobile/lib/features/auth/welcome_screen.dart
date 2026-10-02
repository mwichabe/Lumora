import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/validation.dart';
import '../../core/languages.dart';
import '../../core/network/api_exception.dart';
import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../providers/auth_provider.dart';
import '../../widgets/auth_field.dart';
import '../../widgets/fox_mascot.dart';
import '../../widgets/lumora_button.dart';

enum _Mode { intro, signup, signin }

/// The web app has no separate /login or /signup route — both live here as
/// three states of one screen (frontend/app/onboarding/welcome/page.tsx).
///
/// Light, airy layout on the cream background: purple is used only as an
/// accent (button, focus rings, the mascot's halo), never as a full-bleed
/// fill or gradient, so the screen reads as a modern product form rather
/// than a marketing splash.
class WelcomeScreen extends ConsumerStatefulWidget {
  const WelcomeScreen({super.key});

  @override
  ConsumerState<WelcomeScreen> createState() => _WelcomeScreenState();
}

class _WelcomeScreenState extends ConsumerState<WelcomeScreen> {
  _Mode _mode = _Mode.intro;
  final _formKey = GlobalKey<FormState>();
  final _name = TextEditingController();
  final _email = TextEditingController();
  final _password = TextEditingController();
  bool _obscure = true;
  bool _loading = false;
  String? _error;

  @override
  void dispose() {
    _name.dispose();
    _email.dispose();
    _password.dispose();
    super.dispose();
  }

  void _goTo(_Mode mode) {
    setState(() {
      _mode = mode;
      _error = null;
    });
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      if (_mode == _Mode.signup) {
        await ref.read(authProvider.notifier).register(normaliseEmail(_email.text), _password.text, _name.text.trim());
        if (mounted) context.go('/onboarding/language');
      } else {
        await ref.read(authProvider.notifier).login(normaliseEmail(_email.text), _password.text);
        if (mounted) context.go('/home');
      }
    } on ApiException catch (e) {
      setState(() => _error = e.message);
    } catch (_) {
      setState(() => _error = 'Something went wrong. Please try again.');
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: LumoraColors.cream,
      body: SafeArea(
        child: LayoutBuilder(
          builder: (context, constraints) {
            return SingleChildScrollView(
              padding: const EdgeInsets.fromLTRB(24, 32, 24, 24),
              child: ConstrainedBox(
                constraints: BoxConstraints(minHeight: constraints.maxHeight - 56),
                child: Center(
                  child: ConstrainedBox(
                    constraints: const BoxConstraints(maxWidth: 420),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      children: [
                        _buildBrandHeader(),
                        const SizedBox(height: 28),
                        AnimatedSwitcher(
                          duration: const Duration(milliseconds: 220),
                          switchInCurve: Curves.easeOut,
                          switchOutCurve: Curves.easeIn,
                          transitionBuilder: (child, animation) => FadeTransition(
                            opacity: animation,
                            child: SlideTransition(
                              position: Tween(begin: const Offset(0, 0.03), end: Offset.zero).animate(animation),
                              child: child,
                            ),
                          ),
                          child: AuthCard(
                            key: ValueKey(_mode),
                            child: _mode == _Mode.intro ? _introContent() : _formContent(),
                          ),
                        ),
                        if (_mode == _Mode.intro) ...[
                          const SizedBox(height: 24),
                          _trustStrip(),
                        ],
                      ],
                    ),
                  ),
                ),
              ),
            );
          },
        ),
      ),
    );
  }

  Widget _buildBrandHeader() {
    return Column(
      children: [
        Stack(
          alignment: Alignment.center,
          children: [
            Container(
              width: 108,
              height: 108,
              decoration: const BoxDecoration(color: LumoraColors.purpleLight, shape: BoxShape.circle),
            ),
            const FoxMascot(size: 84, bounce: true),
          ],
        ),
        const SizedBox(height: 16),
        const Text(
          'LUMORA',
          textAlign: TextAlign.center,
          style: TextStyle(color: LumoraColors.ink, fontSize: 30, fontWeight: FontWeight.w800, letterSpacing: -0.5),
        ),
        const SizedBox(height: 4),
        const Text(
          'Language Learning, Reimagined',
          textAlign: TextAlign.center,
          style: TextStyle(color: LumoraColors.slatey, fontSize: 14, fontWeight: FontWeight.w600),
        ),
      ],
    );
  }

  Widget _trustStrip() {
    final items = [('${kAvailableLanguages.length}', 'languages'), ('5 min', 'a day'), ('Free', 'to start')];
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        for (var i = 0; i < items.length; i++) ...[
          if (i > 0)
            Container(width: 1, height: 14, margin: const EdgeInsets.symmetric(horizontal: 14), color: LumoraColors.gray300),
          RichText(
            text: TextSpan(
              children: [
                TextSpan(text: items[i].$1, style: const TextStyle(color: LumoraColors.ink, fontWeight: FontWeight.w800, fontSize: 13)),
                TextSpan(text: ' ${items[i].$2}', style: const TextStyle(color: LumoraColors.slatey, fontSize: 13)),
              ],
            ),
          ),
        ],
      ],
    );
  }

  Widget _introContent() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        const Text('Get started', style: TextStyle(fontSize: 22, fontWeight: FontWeight.w800, color: LumoraColors.ink)),
        const SizedBox(height: 6),
        const Text(
          'Learn Spanish, German, French and more — completely free to start.',
          style: TextStyle(color: LumoraColors.slatey, fontSize: 14, height: 1.4),
        ),
        const SizedBox(height: 24),
        LumoraButton(label: 'Create account', full: true, onPressed: () => _goTo(_Mode.signup)),
        const SizedBox(height: 12),
        LumoraButton(
          label: 'I already have an account',
          full: true,
          variant: LumoraButtonVariant.outline,
          onPressed: () => _goTo(_Mode.signin),
        ),
      ],
    );
  }

  Widget _formContent() {
    final isSignup = _mode == _Mode.signup;
    return Form(
      key: _formKey,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Row(
            children: [
              AuthBackButton(onTap: () => _goTo(_Mode.intro)),
              const SizedBox(width: 8),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      isSignup ? 'Create your account' : 'Welcome back',
                      style: const TextStyle(fontSize: 19, fontWeight: FontWeight.w800, color: LumoraColors.ink),
                    ),
                    Text(
                      isSignup ? 'Start learning in minutes.' : 'Good to see you again.',
                      style: const TextStyle(fontSize: 12, color: LumoraColors.slatey),
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 20),
          if (isSignup) ...[
            const AuthFieldLabel('Name'),
            TextFormField(
              controller: _name,
              textCapitalization: TextCapitalization.words,
              decoration: authFieldDecoration(icon: Icons.person_outline_rounded, hint: 'Your name'),
              validator: (v) => (v == null || v.trim().isEmpty) ? 'Enter your name' : null,
            ),
            const SizedBox(height: 16),
          ],
          const AuthFieldLabel('Email'),
          TextFormField(
            controller: _email,
            keyboardType: TextInputType.emailAddress,
            decoration: authFieldDecoration(icon: Icons.mail_outline_rounded, hint: 'you@example.com'),
            autovalidateMode: AutovalidateMode.onUserInteraction,
            autocorrect: false,
            validator: emailError,
          ),
          const SizedBox(height: 16),
          const AuthFieldLabel('Password'),
          TextFormField(
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
            validator: (v) => (v == null || v.length < 6) ? 'At least 6 characters' : null,
          ),
          if (!isSignup)
            Align(
              alignment: Alignment.centerRight,
              child: TextButton(
                style: TextButton.styleFrom(padding: const EdgeInsets.symmetric(vertical: 8), minimumSize: Size.zero, tapTargetSize: MaterialTapTargetSize.shrinkWrap),
                onPressed: () => context.push('/forgot-password'),
                child: const Text('Forgot password?', style: TextStyle(fontSize: 13, fontWeight: FontWeight.w700)),
              ),
            ),
          if (_error != null) ...[
            const SizedBox(height: 12),
            Container(
              width: double.infinity,
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
              decoration: BoxDecoration(color: LumoraColors.coralLight, borderRadius: BorderRadius.circular(LumoraRadii.md)),
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
          LumoraButton(
            label: isSignup ? 'Create account' : 'Log in',
            full: true,
            loading: _loading,
            onPressed: _loading ? null : _submit,
          ),
          const SizedBox(height: 16),
          Center(
            child: GestureDetector(
              onTap: () => _goTo(isSignup ? _Mode.signin : _Mode.signup),
              child: RichText(
                text: TextSpan(
                  style: const TextStyle(fontSize: 13, color: LumoraColors.slatey),
                  children: [
                    TextSpan(text: isSignup ? 'Already have an account? ' : 'New to Lumora? '),
                    TextSpan(
                      text: isSignup ? 'Log in' : 'Create account',
                      style: const TextStyle(color: LumoraColors.purple, fontWeight: FontWeight.w800),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

}
