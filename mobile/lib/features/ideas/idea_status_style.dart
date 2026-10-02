import 'package:flutter/material.dart';

import '../../core/theme/colors.dart';
import '../../models/idea.dart';

/// Per-status colour and icon, shared by the board and the detail screen.
Color ideaStatusColor(IdeaStatus status) => switch (status) {
      IdeaStatus.draft => LumoraColors.gray500,
      IdeaStatus.underReview => const Color(0xFFD48806),
      IdeaStatus.approved => LumoraColors.teal,
      IdeaStatus.inProgress => LumoraColors.purple,
      IdeaStatus.completed => const Color(0xFF0B9E6E),
      IdeaStatus.archived => LumoraColors.gray500,
    };

IconData ideaStatusIcon(IdeaStatus status) => switch (status) {
      IdeaStatus.draft => Icons.edit_note_rounded,
      IdeaStatus.underReview => Icons.hourglass_top_rounded,
      IdeaStatus.approved => Icons.thumb_up_alt_rounded,
      IdeaStatus.inProgress => Icons.construction_rounded,
      IdeaStatus.completed => Icons.check_circle_rounded,
      IdeaStatus.archived => Icons.inventory_2_rounded,
    };

/// A compact pill showing the idea's status.
class IdeaStatusBadge extends StatelessWidget {
  final IdeaStatus status;
  const IdeaStatusBadge({super.key, required this.status});

  @override
  Widget build(BuildContext context) {
    final c = ideaStatusColor(status);
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
      decoration: BoxDecoration(color: c.withValues(alpha: 0.12), borderRadius: BorderRadius.circular(999)),
      child: Row(mainAxisSize: MainAxisSize.min, children: [
        Icon(ideaStatusIcon(status), size: 12, color: c),
        const SizedBox(width: 4),
        Text(ideaStatusLabel(status), style: TextStyle(color: c, fontSize: 11, fontWeight: FontWeight.w800)),
      ]),
    );
  }
}
