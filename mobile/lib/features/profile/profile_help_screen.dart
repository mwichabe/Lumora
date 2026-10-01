import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:url_launcher/url_launcher.dart';

import '../../core/theme/colors.dart';
import '../../core/theme/radii.dart';
import '../../core/theme/shadows.dart';

const _kFaq = [
  ('How do I add another language?', 'Open the Learn tab and tap the language switcher (top-right), or go to Profile → My Languages → Add. Your progress is saved per language.'),
  ('How are streaks calculated?', 'Complete at least one lesson each day to keep your streak alive. Miss a day and it resets — so even a quick 5-minute lesson counts!'),
  ("Why can't I open a listening or reading session?", "Those unlock once you've completed all the lessons in that unit. Finish the unit's lessons first, then they open up."),
  ('How does the fluency score work?', 'In speaking practice we listen to what you say and compare it to the target phrase, giving a 0–100% score.'),
  ('How do I change my daily goal or password?', 'Go to Profile → Account Settings. You can edit your name, avatar, daily goal, change your password, or delete your account there.'),
];

class ProfileHelpScreen extends StatefulWidget {
  const ProfileHelpScreen({super.key});

  @override
  State<ProfileHelpScreen> createState() => _ProfileHelpScreenState();
}

class _ProfileHelpScreenState extends State<ProfileHelpScreen> {
  int? _open = 0;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: LumoraColors.cream,
      appBar: AppBar(
        backgroundColor: Colors.white,
        foregroundColor: LumoraColors.ink,
        elevation: 0,
        leading: BackButton(onPressed: () => context.canPop() ? context.pop() : context.go('/profile')),
        title: const Text('Help & Support', style: TextStyle(fontWeight: FontWeight.w800, color: LumoraColors.ink)),
      ),
      body: ListView(
        padding: const EdgeInsets.all(20),
        children: [
          const Text('FREQUENTLY ASKED', style: TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.gray500)),
          const SizedBox(height: 10),
          Container(
            decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
            child: Column(
              children: [
                for (var i = 0; i < _kFaq.length; i++)
                  Column(children: [
                    if (i > 0) const Divider(height: 1),
                    InkWell(
                      onTap: () => setState(() => _open = _open == i ? null : i),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                        child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                          Row(children: [
                            Expanded(child: Text(_kFaq[i].$1, style: const TextStyle(fontWeight: FontWeight.w800))),
                            AnimatedRotation(turns: _open == i ? 0.5 : 0, duration: const Duration(milliseconds: 150), child: const Icon(Icons.keyboard_arrow_down, color: LumoraColors.gray500)),
                          ]),
                          if (_open == i) Padding(padding: const EdgeInsets.only(top: 8), child: Text(_kFaq[i].$2, style: const TextStyle(color: LumoraColors.slatey))),
                        ]),
                      ),
                    ),
                  ]),
              ],
            ),
          ),
          const SizedBox(height: 24),
          const Text('STILL NEED HELP?', style: TextStyle(fontSize: 11, fontWeight: FontWeight.w800, color: LumoraColors.gray500)),
          const SizedBox(height: 10),
          Material(
            color: Colors.white,
            borderRadius: BorderRadius.circular(LumoraRadii.xl),
            child: InkWell(
              borderRadius: BorderRadius.circular(LumoraRadii.xl),
              onTap: () => launchUrl(Uri.parse('mailto:support@lumora.app?subject=Lumora%20Support')),
              child: Container(
                padding: const EdgeInsets.all(16),
                decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(LumoraRadii.xl), boxShadow: LumoraShadows.card),
                child: Row(children: [
                  Container(width: 48, height: 48, decoration: BoxDecoration(color: LumoraColors.purpleLight, borderRadius: BorderRadius.circular(LumoraRadii.md)), child: const Icon(Icons.mail_outline, color: LumoraColors.purple)),
                  const SizedBox(width: 12),
                  const Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                    Text('Contact support', style: TextStyle(fontWeight: FontWeight.w800)),
                    Text('support@lumora.app', style: TextStyle(color: LumoraColors.slatey, fontSize: 12)),
                  ]),
                ]),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
