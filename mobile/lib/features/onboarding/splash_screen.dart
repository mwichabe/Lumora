import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../core/storage/session_storage.dart';
import '../../providers/auth_provider.dart';
import '../../widgets/fox_mascot.dart';

/// Matches frontend/app/page.tsx: an animated brand splash that decides where
/// to route the user once auth state has settled — onboarding, the last
/// screen they were on, or the welcome/auth screen.
class SplashScreen extends ConsumerStatefulWidget {
  const SplashScreen({super.key});

  @override
  ConsumerState<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends ConsumerState<SplashScreen> {
  bool _navigated = false;

  static const _appRoutes = ['/home', '/learn', '/practice', '/leaderboard', '/profile', '/lesson'];

  Future<void> _maybeNavigate(AuthState auth) async {
    if (_navigated || auth.loading) return;
    _navigated = true;

    String target;
    if (auth.user != null) {
      if (auth.user!.targetLanguage.isEmpty) {
        target = '/onboarding/language';
      } else {
        final last = await SessionStorage.instance.getLastRoute();
        target = (last != null && _appRoutes.any((r) => last.startsWith(r))) ? last : '/home';
      }
    } else {
      target = '/welcome';
    }

    final delay = auth.user != null ? 400 : 1600;
    await Future.delayed(Duration(milliseconds: delay));
    if (mounted) context.go(target);
  }

  @override
  Widget build(BuildContext context) {
    final auth = ref.watch(authProvider);
    _maybeNavigate(auth);

    return Scaffold(
      body: Container(
        width: double.infinity,
        height: double.infinity,
        decoration: const BoxDecoration(
          gradient: RadialGradient(
            center: Alignment.topCenter,
            radius: 1.3,
            colors: [Color(0xFF7B4AD6), Color(0xFF4A2A9E), Color(0xFF1F1640), Color(0xFF0F0F24)],
            stops: [0, 0.38, 0.72, 1],
          ),
        ),
        child: SafeArea(
          child: Column(
            children: [
              const Spacer(flex: 3),
              const FoxMascot(size: 180, glow: true, bounce: true),
              const SizedBox(height: 16),
              ShaderMask(
                shaderCallback: (bounds) => const LinearGradient(
                  begin: Alignment.topCenter,
                  end: Alignment.bottomCenter,
                  colors: [Colors.white, Color(0xFFF5A623)],
                ).createShader(bounds),
                child: const Text(
                  'LUMORA',
                  style: TextStyle(
                    fontSize: 56,
                    fontWeight: FontWeight.w800,
                    color: Colors.white,
                    letterSpacing: -1.5,
                  ),
                ),
              ),
              const SizedBox(height: 8),
              const Text(
                'Learn a language. Fall in love with it.',
                style: TextStyle(color: Colors.white70, fontStyle: FontStyle.italic),
              ),
              const Spacer(flex: 4),
              SizedBox(
                width: 200,
                child: ClipRRect(
                  borderRadius: BorderRadius.circular(4),
                  child: const LinearProgressIndicator(
                    minHeight: 4,
                    backgroundColor: Colors.white24,
                    valueColor: AlwaysStoppedAnimation(Color(0xFFF5A623)),
                  ),
                ),
              ),
              const SizedBox(height: 12),
              const Text('LOADING',
                  style: TextStyle(
                      color: Colors.white38, fontSize: 12, fontWeight: FontWeight.w700, letterSpacing: 2)),
              const SizedBox(height: 32),
            ],
          ),
        ),
      ),
    );
  }
}
