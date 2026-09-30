import 'package:flutter/material.dart';

import '../core/theme/colors.dart';
import '../core/theme/radii.dart';
import '../core/theme/shadows.dart';
import '../core/voices.dart';
import 'lumora_button.dart';

class ReviewItem {
  final String? prompt;
  final String question;
  final String correctAnswer;
  final String? playText;
  final String? speaker;
  const ReviewItem({this.prompt, required this.question, required this.correctAnswer, this.playText, this.speaker});
}

/// End-of-session recap of missed items, ported from
/// frontend/components/MistakesReview.tsx.
class MistakesReview extends StatelessWidget {
  final List<ReviewItem> items;
  final VoidCallback onDone;
  final String finishLabel;

  const MistakesReview({super.key, required this.items, required this.onDone, this.finishLabel = 'Finish'});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(children: [
          Container(
            width: 36, height: 36,
            decoration: BoxDecoration(color: LumoraColors.amberLight, borderRadius: BorderRadius.circular(LumoraRadii.md)),
            child: const Icon(Icons.replay_rounded, color: LumoraColors.amber, size: 18),
          ),
          const SizedBox(width: 10),
          Expanded(
            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              const Text('Review your mistakes', style: TextStyle(fontSize: 18, fontWeight: FontWeight.w800)),
              Text('Go over these ${items.length == 1 ? "one" : items.length} before you finish.',
                  style: const TextStyle(color: LumoraColors.slatey, fontSize: 12)),
            ]),
          ),
        ]),
        const SizedBox(height: 16),
        Expanded(
          child: ListView.separated(
            itemCount: items.length,
            separatorBuilder: (_, _) => const SizedBox(height: 10),
            itemBuilder: (context, i) {
              final it = items[i];
              return Container(
                padding: const EdgeInsets.all(16),
                decoration: BoxDecoration(
                  color: Colors.white,
                  border: Border.all(color: LumoraColors.gray100),
                  borderRadius: BorderRadius.circular(LumoraRadii.xl),
                  boxShadow: LumoraShadows.card,
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    if (it.prompt != null && it.prompt!.isNotEmpty)
                      Text(it.prompt!.toUpperCase(),
                          style: const TextStyle(fontSize: 10, fontWeight: FontWeight.w800, color: LumoraColors.gray500)),
                    const SizedBox(height: 2),
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Expanded(child: Text(it.question, style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 16))),
                        if (it.playText != null)
                          InkWell(
                            onTap: () => Voices.instance.speakAs(it.speaker, it.playText!),
                            child: Container(
                              width: 36, height: 36,
                              decoration: const BoxDecoration(color: LumoraColors.purple, shape: BoxShape.circle),
                              child: const Icon(Icons.volume_up_rounded, size: 16, color: Colors.white),
                            ),
                          ),
                      ],
                    ),
                    const SizedBox(height: 8),
                    Container(
                      width: double.infinity,
                      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                      decoration: BoxDecoration(color: LumoraColors.tealLight, borderRadius: BorderRadius.circular(LumoraRadii.sm)),
                      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                        const Text('CORRECT ANSWER', style: TextStyle(fontSize: 10, fontWeight: FontWeight.w800, color: LumoraColors.teal)),
                        Text(it.correctAnswer, style: const TextStyle(fontWeight: FontWeight.w800)),
                      ]),
                    ),
                  ],
                ),
              );
            },
          ),
        ),
        const SizedBox(height: 12),
        LumoraButton(label: finishLabel, full: true, onPressed: onDone),
      ],
    );
  }
}
