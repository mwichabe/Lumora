import 'package:flutter/material.dart';

import '../core/theme/colors.dart';
import '../core/theme/radii.dart';
import '../models/user.dart';
import '../providers/hearts_provider.dart';
import 'lumora_button.dart';

Future<void> showOutOfHeartsModal(
  BuildContext context, {
  required HeartsStatus status,
  required int secondsToNext,
  String? note,
  String closeLabel = "I'll wait",
  required VoidCallback onClose,
  required Future<void> Function() onBuy,
}) {
  return showModalBottomSheet(
    context: context,
    isDismissible: true,
    backgroundColor: Colors.white,
    shape: const RoundedRectangleBorder(borderRadius: BorderRadius.vertical(top: Radius.circular(28))),
    builder: (context) => _OutOfHeartsSheet(
      status: status,
      secondsToNext: secondsToNext,
      note: note,
      closeLabel: closeLabel,
      onClose: onClose,
      onBuy: onBuy,
    ),
  );
}

class _OutOfHeartsSheet extends StatefulWidget {
  final HeartsStatus status;
  final int secondsToNext;
  final String? note;
  final String closeLabel;
  final VoidCallback onClose;
  final Future<void> Function() onBuy;

  const _OutOfHeartsSheet({
    required this.status,
    required this.secondsToNext,
    required this.note,
    required this.closeLabel,
    required this.onClose,
    required this.onBuy,
  });

  @override
  State<_OutOfHeartsSheet> createState() => _OutOfHeartsSheetState();
}

class _OutOfHeartsSheetState extends State<_OutOfHeartsSheet> {
  bool _buying = false;

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: double.infinity,
            padding: const EdgeInsets.fromLTRB(24, 32, 24, 20),
            decoration: BoxDecoration(
              color: LumoraColors.coral.withValues(alpha: 0.1),
              borderRadius: const BorderRadius.vertical(top: Radius.circular(28)),
            ),
            child: Column(
              children: [
                Container(
                  width: 64, height: 64,
                  decoration: BoxDecoration(color: LumoraColors.coral.withValues(alpha: 0.15), shape: BoxShape.circle),
                  child: const Icon(Icons.favorite, color: LumoraColors.coral, size: 32),
                ),
                const SizedBox(height: 12),
                const Text("You're out of hearts", style: TextStyle(fontSize: 20, fontWeight: FontWeight.w800)),
                const SizedBox(height: 4),
                Text(
                  widget.note ?? 'Hearts refill over time — or top up now to keep learning.',
                  textAlign: TextAlign.center,
                  style: const TextStyle(color: LumoraColors.slatey),
                ),
              ],
            ),
          ),
          Padding(
            padding: const EdgeInsets.fromLTRB(24, 20, 24, 24),
            child: Column(
              children: [
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.symmetric(vertical: 12),
                  decoration: BoxDecoration(color: LumoraColors.gray50, borderRadius: BorderRadius.circular(LumoraRadii.lg)),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      const Icon(Icons.access_time, size: 18, color: LumoraColors.slatey),
                      const SizedBox(width: 8),
                      Text(
                        widget.secondsToNext > 0 ? 'Next heart in ${fmtCountdown(widget.secondsToNext)}' : 'A heart is on its way…',
                        style: const TextStyle(fontWeight: FontWeight.w800),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 4),
                Text('One heart every ${widget.status.regenMinutes} minutes',
                    style: const TextStyle(fontSize: 12, color: LumoraColors.slatey)),
                const SizedBox(height: 16),
                if (widget.status.paymentsEnabled)
                  LumoraButton(
                    label: _buying
                        ? 'Opening checkout…'
                        : 'Refill ${widget.status.max} hearts — KES ${widget.status.refillPriceKes}'
                            '${widget.status.refillPriceUsd > 0 ? " (≈ \$${widget.status.refillPriceUsd.toStringAsFixed(2)})" : ""}',
                    full: true,
                    loading: _buying,
                    onPressed: () async {
                      setState(() => _buying = true);
                      await widget.onBuy();
                      if (mounted) setState(() => _buying = false);
                    },
                  ),
                const SizedBox(height: 8),
                LumoraButton(label: widget.closeLabel, full: true, variant: LumoraButtonVariant.outline, onPressed: widget.onClose),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
