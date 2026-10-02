import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../providers/league_provider.dart';
import '../../providers/learn_provider.dart';
import '../../widgets/more_sheet.dart';
import '../theme/colors.dart';

/// The 4-tab bottom nav (Home/Learn/Practice/Leagues) + "More" sheet, matching
/// frontend/components/BottomTabBar.tsx. Wraps the app's four primary
/// destinations via StatefulShellRoute so each tab keeps its own navigation
/// stack and scroll position.
class AppShell extends ConsumerWidget {
  final StatefulNavigationShell shell;
  const AppShell({super.key, required this.shell});

  static const _tabs = [
    (label: 'Home', icon: Icons.home_rounded),
    (label: 'Learn', icon: Icons.menu_book_rounded),
    (label: 'Practice', icon: Icons.mic_rounded),
    (label: 'Leagues', icon: Icons.emoji_events_rounded),
  ];

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // Start fetching the course as soon as the signed-in app is on screen, not
    // when the Learn tab is first opened — by then it's already there.
    ref.listen(learnProvider, (_, _) {});
    // Likewise the league, so the Leagues tab is filled in on first tap.
    ref.listen(leagueProvider, (_, _) {});

    return Scaffold(
      body: shell,
      bottomNavigationBar: SafeArea(
        top: false,
        child: SizedBox(
          height: 64,
          child: Row(
            children: [
              for (var i = 0; i < _tabs.length; i++)
                Expanded(
                  child: _TabButton(
                    label: _tabs[i].label,
                    icon: _tabs[i].icon,
                    active: shell.currentIndex == i,
                    onTap: () => shell.goBranch(i, initialLocation: i == shell.currentIndex),
                  ),
                ),
              Expanded(
                child: _TabButton(
                  label: 'More',
                  icon: Icons.more_horiz_rounded,
                  active: false,
                  onTap: () => showMoreSheet(context),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _TabButton extends StatelessWidget {
  final String label;
  final IconData icon;
  final bool active;
  final VoidCallback onTap;
  const _TabButton({required this.label, required this.icon, required this.active, required this.onTap});

  @override
  Widget build(BuildContext context) {
    final color = active ? LumoraColors.purple : LumoraColors.gray300;
    return InkWell(
      onTap: onTap,
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          if (active)
            Container(margin: const EdgeInsets.only(bottom: 4), width: 6, height: 6,
                decoration: const BoxDecoration(color: LumoraColors.purple, shape: BoxShape.circle))
          else
            const SizedBox(height: 10),
          Icon(icon, size: 24, color: color),
          const SizedBox(height: 2),
          Text(label, style: TextStyle(fontSize: 10, fontWeight: FontWeight.w700, color: color)),
        ],
      ),
    );
  }
}
