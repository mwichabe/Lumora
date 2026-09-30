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
  // Bumped on every load so a slow response from an earlier load (e.g. the
  // previous language) can't overwrite a newer one.
  int _generation = 0;

  @override
  LearnState build() {
    load();
    return LearnState.initial();
  }

  /// Fires the three requests together and publishes each as it lands, like
  /// the web client: the course renders as soon as skills arrive instead of
  /// waiting on listening and reading too. `loading` tracks skills only, and
  /// whatever is already on screen stays there while a reload is in flight.
  Future<void> load() async {
    final gen = ++_generation;
    state = LearnState(skills: state.skills, listening: state.listening, reading: state.reading, loading: true, error: false);

    final skills = ApiClient.instance.skills().then((skills) {
      if (gen != _generation) return;
      state = LearnState(skills: skills, listening: state.listening, reading: state.reading, loading: false, error: false);
    }, onError: (_) {
      if (gen != _generation) return;
      state = LearnState(skills: state.skills, listening: state.listening, reading: state.reading, loading: false, error: true);
    });

    final listening = ApiClient.instance.listeningSessions().then((listening) {
      if (gen != _generation) return;
      state = LearnState(skills: state.skills, listening: listening, reading: state.reading, loading: state.loading, error: state.error);
    }, onError: (_) {});

    final reading = ApiClient.instance.readingSessions().then((reading) {
      if (gen != _generation) return;
      state = LearnState(skills: state.skills, listening: state.listening, reading: reading, loading: state.loading, error: state.error);
    }, onError: (_) {});

    await Future.wait([skills, listening, reading]);
  }
}

final learnProvider = NotifierProvider<LearnController, LearnState>(LearnController.new);
