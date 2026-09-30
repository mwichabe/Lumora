import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../core/network/api_client.dart';
import '../models/lesson.dart';
import '../models/listening_reading.dart';

class LearnState {
  final List<Skill> skills;
  final List<ListeningSession> listening;
  final List<ReadingSession> reading;
  final bool loading;
  final bool error;

  const LearnState({
    required this.skills,
    required this.listening,
    required this.reading,
    required this.loading,
    required this.error,
  });

  factory LearnState.initial() =>
      const LearnState(skills: [], listening: [], reading: [], loading: true, error: false);
}

class LearnController extends Notifier<LearnState> {
  @override
  LearnState build() {
    load();
    return LearnState.initial();
  }

  Future<void> load() async {
    state = LearnState(skills: state.skills, listening: state.listening, reading: state.reading, loading: true, error: false);
    try {
      final skills = await ApiClient.instance.skills();
      List<ListeningSession> listening = [];
      List<ReadingSession> reading = [];
      try {
        listening = await ApiClient.instance.listeningSessions();
      } catch (_) {}
      try {
        reading = await ApiClient.instance.readingSessions();
      } catch (_) {}
      state = LearnState(skills: skills, listening: listening, reading: reading, loading: false, error: false);
    } catch (_) {
      state = LearnState(skills: state.skills, listening: state.listening, reading: state.reading, loading: false, error: true);
    }
  }
}

final learnProvider = NotifierProvider<LearnController, LearnState>(LearnController.new);
