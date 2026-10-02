import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:google_fonts/google_fonts.dart';

import 'core/network/dio_client.dart';
import 'core/router/app_router.dart';
import 'core/theme/app_theme.dart';
import 'core/voices.dart';
import 'providers/auth_provider.dart';

void main() {
  // Nunito ships in assets/google_fonts/ — never fetch it over the network.
  GoogleFonts.config.allowRuntimeFetching = false;
  runApp(const ProviderScope(child: LumoraApp()));
}

class LumoraApp extends ConsumerStatefulWidget {
  const LumoraApp({super.key});

  @override
  ConsumerState<LumoraApp> createState() => _LumoraAppState();
}

class _LumoraAppState extends ConsumerState<LumoraApp> {
  final _messenger = GlobalKey<ScaffoldMessengerState>();

  @override
  void initState() {
    super.initState();
    // Speech follows the course: German lessons are read by a German voice,
    // French by a French one. Without this every course used the default
    // (Spanish) locale and the wrong accent.
    ref.listenManual(
      authProvider.select((a) => a.user?.targetLanguage),
      (_, lang) => Voices.instance.setSpeechLanguage(lang),
      fireImmediately: true,
    );
    Voices.instance.onMissingVoice = (message) => _messenger.currentState?.showSnackBar(
          SnackBar(content: Text(message), duration: const Duration(seconds: 8)),
        );
    // The session expired server-side: sign out properly. Navigating alone
    // isn't enough — while the auth state still says "signed in", the router
    // sends /welcome straight back to /home. Updating the state lets the
    // router's redirect take the user to sign in from wherever they are.
    DioClient.instance.onUnauthorized = () {
      ref.read(authProvider.notifier).logout();
    };
  }

  @override
  Widget build(BuildContext context) {
    final router = ref.watch(routerProvider);
    return MaterialApp.router(
      title: 'Lumora',
      debugShowCheckedModeBanner: false,
      scaffoldMessengerKey: _messenger,
      theme: AppTheme.light(),
      routerConfig: router,
    );
  }
}
