import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../core/network/api_client.dart';
import '../models/user.dart';

class HeartsState {
  final HeartsStatus? status;
  final int secondsToNext;
  const HeartsState({required this.status, required this.secondsToNext});

  int get hearts => status?.hearts ?? 0;
  int get max => status?.max ?? 5;
  bool get full => status?.full ?? true;
}

/// Server-authoritative hearts, mirrored from frontend/lib/hearts.ts: a 60s
/// background poll plus a local 1s countdown to the next regenerated heart.
class HeartsController extends Notifier<HeartsState> {
  Timer? _poll;
  Timer? _tick;

  @override
  HeartsState build() {
    reload();
    _poll = Timer.periodic(const Duration(seconds: 60), (_) => reload());
    ref.onDispose(() {
      _poll?.cancel();
      _tick?.cancel();
    });
    return const HeartsState(status: null, secondsToNext: 0);
  }

  void _restartTick() {
    _tick?.cancel();
    if (state.status == null || state.full || state.secondsToNext <= 0) return;
    _tick = Timer.periodic(const Duration(seconds: 1), (_) {
      if (state.secondsToNext <= 1) {
        _tick?.cancel();
        reload();
      } else {
        state = HeartsState(status: state.status, secondsToNext: state.secondsToNext - 1);
      }
    });
  }

  Future<void> reload() async {
    try {
      final s = await ApiClient.instance.heartsStatus();
      state = HeartsState(status: s, secondsToNext: s.secondsToNext);
      _restartTick();
    } catch (_) {
      // offline / unauthenticated — keep whatever we had
    }
  }

  Future<HeartsStatus?> lose() async {
    try {
      final s = await ApiClient.instance.loseHeart();
      state = HeartsState(status: s, secondsToNext: s.secondsToNext);
      _restartTick();
      return s;
    } catch (_) {
      return null;
    }
  }
}

final heartsProvider = NotifierProvider<HeartsController, HeartsState>(HeartsController.new);

String fmtCountdown(int s) {
  final m = s ~/ 60;
  final r = s % 60;
  return '$m:${r.toString().padLeft(2, '0')}';
}
