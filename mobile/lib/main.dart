import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'core/network/dio_client.dart';
import 'core/router/app_router.dart';
import 'core/theme/app_theme.dart';

void main() {
  runApp(const ProviderScope(child: LumoraApp()));
}

class LumoraApp extends ConsumerStatefulWidget {
  const LumoraApp({super.key});

  @override
  ConsumerState<LumoraApp> createState() => _LumoraAppState();
}

class _LumoraAppState extends ConsumerState<LumoraApp> {
  @override
  void initState() {
    super.initState();
    DioClient.instance.onUnauthorized = () {
      final router = ref.read(routerProvider);
      router.go('/welcome');
    };
  }

  @override
  Widget build(BuildContext context) {
    final router = ref.watch(routerProvider);
    return MaterialApp.router(
      title: 'Lumora',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.light(),
      routerConfig: router,
    );
  }
}
