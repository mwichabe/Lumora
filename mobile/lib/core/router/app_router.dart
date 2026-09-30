import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../features/auth/welcome_screen.dart';
import '../../features/auth/forgot_password_screen.dart';
import '../../features/auth/reset_password_screen.dart';
import '../../features/onboarding/splash_screen.dart';
import '../../features/onboarding/language_screen.dart';
import '../../features/onboarding/goal_screen.dart';
import '../../features/home/home_screen.dart';
import '../../features/learn/learn_screen.dart';
import '../../features/lesson/lesson_screen.dart';
import '../../features/lesson/lesson_complete_screen.dart';
import '../../features/listening/listening_screen.dart';
import '../../features/reading/reading_screen.dart';
import '../../features/practice/practice_screen.dart';
import '../../features/practice/practice_run_screen.dart';
import '../../features/leaderboard/leaderboard_screen.dart';
import '../../features/notifications/notifications_screen.dart';
import '../../features/chat/chat_list_screen.dart';
import '../../features/chat/chat_thread_screen.dart';
import '../../features/ideas/ideas_screen.dart';
import '../../features/ideas/idea_detail_screen.dart';
import '../../features/ideas/idea_thread_screen.dart';
import '../../features/exam/exam_screen.dart';
import '../../features/payments/payment_callback_screen.dart';
import '../../features/certificates/certificates_screen.dart';
import '../../features/certificates/certificate_detail_screen.dart';
import '../../features/verify/verify_screen.dart';
import '../../features/profile/profile_screen.dart';
import '../../features/profile/profile_settings_screen.dart';
import '../../features/profile/profile_help_screen.dart';
import '../storage/session_storage.dart';
import '../../providers/auth_provider.dart';
import 'app_shell.dart';

/// Bridges Riverpod's auth state into a Listenable go_router can watch, so
/// `redirect` re-runs whenever sign-in state changes (matches the web splash
/// screen's routing rules in frontend/app/page.tsx).
class _AuthRefresh extends ChangeNotifier {
  _AuthRefresh(Ref ref) {
    ref.listen(authProvider, (_, _) => notifyListeners());
  }
}

final routerProvider = Provider<GoRouter>((ref) {
  final refresh = _AuthRefresh(ref);

  return GoRouter(
    initialLocation: '/',
    refreshListenable: refresh,
    redirect: (context, state) async {
      final loc = state.matchedLocation;
      final auth = ref.read(authProvider);

      const publicPaths = ['/welcome', '/forgot-password', '/reset-password'];
      final isPublic = publicPaths.any((p) => loc.startsWith(p)) || loc.startsWith('/verify/');

      if (loc == '/') return null; // splash decides for itself
      if (auth.loading) return null;

      if (!auth.isAuthenticated) {
        return isPublic ? null : '/welcome';
      }

      // Authenticated but hasn't finished onboarding yet.
      if (auth.user!.targetLanguage.isEmpty &&
          !loc.startsWith('/onboarding') &&
          !isPublic) {
        return '/onboarding/language';
      }

      if (loc == '/welcome') {
        await SessionStorage.instance.setLastRoute('');
        return '/home';
      }

      return null;
    },
    routes: [
      GoRoute(path: '/', builder: (context, state) => const SplashScreen()),
      GoRoute(path: '/welcome', builder: (context, state) => const WelcomeScreen()),
      GoRoute(path: '/forgot-password', builder: (context, state) => const ForgotPasswordScreen()),
      GoRoute(
        path: '/reset-password',
        builder: (context, state) => ResetPasswordScreen(token: state.uri.queryParameters['token']),
      ),
      GoRoute(
        path: '/onboarding/language',
        builder: (context, state) => OnboardingLanguageScreen(
          adding: state.uri.queryParameters['add'] == '1',
        ),
      ),
      GoRoute(path: '/onboarding/goal', builder: (context, state) => const OnboardingGoalScreen()),

      GoRoute(
        path: '/lesson/:id',
        builder: (context, state) => LessonScreen(id: int.parse(state.pathParameters['id']!)),
      ),
      GoRoute(
        path: '/lesson/:id/complete',
        builder: (context, state) =>
            LessonCompleteScreen(id: int.parse(state.pathParameters['id']!)),
      ),
      GoRoute(
        path: '/listening/:id',
        builder: (context, state) =>
            ListeningScreen(id: int.parse(state.pathParameters['id']!)),
      ),
      GoRoute(
        path: '/reading/:id',
        builder: (context, state) => ReadingScreen(id: int.parse(state.pathParameters['id']!)),
      ),
      GoRoute(
        path: '/practice/run',
        builder: (context, state) => PracticeRunScreen(mode: state.uri.queryParameters['mode'] ?? ''),
      ),

      GoRoute(path: '/notifications', builder: (context, state) => const NotificationsScreen()),
      GoRoute(path: '/chat', builder: (context, state) => const ChatListScreen()),
      GoRoute(
        path: '/chat/:id',
        builder: (context, state) => ChatThreadScreen(userId: int.parse(state.pathParameters['id']!)),
      ),
      GoRoute(
        path: '/ideas',
        builder: (context, state) =>
            IdeasScreen(ideaId: int.tryParse(state.uri.queryParameters['idea'] ?? '')),
      ),
      GoRoute(
        path: '/ideas/:id',
        builder: (context, state) => IdeaDetailScreen(id: int.parse(state.pathParameters['id']!)),
      ),
      GoRoute(
        path: '/ideas/:id/thread',
        builder: (context, state) => IdeaThreadScreen(id: int.parse(state.pathParameters['id']!)),
      ),
      GoRoute(
        path: '/exam',
        builder: (context, state) => const ExamScreen(),
      ),
      GoRoute(path: '/payment/callback', builder: (context, state) => PaymentCallbackScreen(
            reference: state.uri.queryParameters['reference'] ?? state.uri.queryParameters['trxref'],
          )),
      GoRoute(path: '/certificates', builder: (context, state) => const CertificatesScreen()),
      GoRoute(
        path: '/certificates/:id',
        builder: (context, state) =>
            CertificateDetailScreen(id: int.parse(state.pathParameters['id']!)),
      ),
      GoRoute(
        path: '/verify/:serial',
        builder: (context, state) => VerifyScreen(serial: state.pathParameters['serial']!),
      ),
      GoRoute(path: '/profile/settings', builder: (context, state) => const ProfileSettingsScreen()),
      GoRoute(path: '/profile/help', builder: (context, state) => const ProfileHelpScreen()),

      StatefulShellRoute.indexedStack(
        builder: (context, state, shell) => AppShell(shell: shell),
        branches: [
          StatefulShellBranch(routes: [
            GoRoute(path: '/home', builder: (context, state) => const HomeScreen()),
          ]),
          StatefulShellBranch(routes: [
            GoRoute(path: '/learn', builder: (context, state) => const LearnScreen()),
          ]),
          StatefulShellBranch(routes: [
            GoRoute(path: '/practice', builder: (context, state) => const PracticeScreen()),
          ]),
          StatefulShellBranch(routes: [
            GoRoute(path: '/leaderboard', builder: (context, state) => const LeaderboardScreen()),
          ]),
        ],
      ),
      GoRoute(path: '/profile', builder: (context, state) => const ProfileScreen()),
    ],
  );
});
