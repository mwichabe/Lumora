import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../core/json_utils.dart';
import '../core/network/api_client.dart';
import '../core/storage/session_storage.dart';
import '../models/lesson.dart';
import '../models/listening_reading.dart';
import 'auth_provider.dart';

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

List<Skill> _parseSkills(Map<String, dynamic> j) => asList(j['skills'], (e) => Skill.fromJson(asMap(e)));

List<ListeningSession> _parseListening(Map<String, dynamic> j) =>
    asList(j['sessions'], (e) => ListeningSession.fromJson(asMap(e)));

List<ReadingSession> _parseReading(Map<String, dynamic> j) =>
    asList(j['sessions'], (e) => ReadingSession.fromJson(asMap(e)));

class LearnController extends Notifier<LearnState> {
  // Bumped on every load so a slow response from an earlier load (e.g. the
  // previous language) can't overwrite a newer one.
  int _generation = 0;
  bool _loading = false;

  @override
  LearnState build() {
    // The course belongs to one learner and one language: rebuild — and so
    // reload — whenever either changes (sign-in, sign-out, language switch).
    final (userId, _) = ref.watch(authProvider.select((a) => (a.user?.id, a.user?.targetLanguage)));
    // Finishing a lesson, listening or reading session earns XP: reload so the
    // path shows it done (and the next one unlocked) without a manual refresh.
    // Skipped while a load is running — that one will already include it.
    ref.listen(authProvider.select((a) => a.user?.xp), (prev, next) {
      if (prev != null && next != null && prev != next && !_loading) load();
    });
    // Deferred a tick: load() reads and writes `state`, which only exists once
    // build has returned.
    if (userId != null) Future.microtask(load);
    return LearnState.initial();
  }

  /// One saved copy per learner and language, so switching course (or account)
  /// never shows someone else's progress.
  String? _cacheKey() {
    final user = ref.read(authProvider).user;
    if (user == null) return null;
    return 'learn_${user.id}_${user.targetLanguage.isEmpty ? "es" : user.targetLanguage}';
  }

  /// Shows the course in two steps so the screen is never waiting on the
  /// network when it doesn't have to:
  ///
  ///  1. The copy saved from the last visit is put on screen straight away.
  ///  2. The three requests fire together and each replaces its part as it
  ///     lands — skills don't wait on listening and reading — then the fresh
  ///     responses are saved for next time.
  ///
  /// `loading` tracks skills only, and whatever is already on screen stays
  /// there while a reload is in flight.
  Future<void> load() async {
    final gen = ++_generation;
    _loading = true;
    final cacheKey = _cacheKey();
    state = LearnState(skills: state.skills, listening: state.listening, reading: state.reading, loading: true, error: false);

    // Set once the network has answered, so the saved copy (which may resolve
    // later than a fast response) never overwrites fresher data.
    var skillsLanded = false, listeningLanded = false, readingLanded = false;

    if (cacheKey != null && state.skills.isEmpty) {
      SessionStorage.instance.readCache(cacheKey).then((saved) {
        if (saved == null || gen != _generation) return;
        try {
          state = LearnState(
            skills: skillsLanded ? state.skills : _parseSkills(asMap(saved['skills'])),
            listening: listeningLanded ? state.listening : _parseListening(asMap(saved['listening'])),
            reading: readingLanded ? state.reading : _parseReading(asMap(saved['reading'])),
            // Still refreshing, but there is something to show: `loading` only
            // drives the empty-screen spinner, which checks for skills too.
            loading: state.loading,
            error: false,
          );
        } catch (_) {
          // An unreadable copy is simply ignored; the network fills the screen.
        }
      }, onError: (_) {});
    }

    final fresh = <String, dynamic>{};

    final skills = ApiClient.instance.getJson('/api/skills').then((json) {
      if (gen != _generation) return;
      skillsLanded = true;
      fresh['skills'] = json;
      state = LearnState(skills: _parseSkills(json), listening: state.listening, reading: state.reading, loading: false, error: false);
    }, onError: (_) {
      if (gen != _generation) return;
      state = LearnState(skills: state.skills, listening: state.listening, reading: state.reading, loading: false, error: true);
    });

    final listening = ApiClient.instance.getJson('/api/listening').then((json) {
      if (gen != _generation) return;
      listeningLanded = true;
      fresh['listening'] = json;
      state = LearnState(skills: state.skills, listening: _parseListening(json), reading: state.reading, loading: state.loading, error: state.error);
    }, onError: (_) {});

    final reading = ApiClient.instance.getJson('/api/reading').then((json) {
      if (gen != _generation) return;
      readingLanded = true;
      fresh['reading'] = json;
      state = LearnState(skills: state.skills, listening: state.listening, reading: _parseReading(json), loading: state.loading, error: state.error);
    }, onError: (_) {});

    await Future.wait([skills, listening, reading]);
    if (gen == _generation) _loading = false;

    // Only a complete set is saved: a partial one would show a course with its
    // listening or reading sessions missing on the next launch.
    if (gen == _generation && cacheKey != null && fresh.length == 3) {
      SessionStorage.instance.writeCache(cacheKey, fresh).catchError((_) {});
    }
  }
}

final learnProvider = NotifierProvider<LearnController, LearnState>(LearnController.new);

/// Which view the Learn tab shows. Shared state rather than local to the
/// screen so other screens can open it on a given view — Home's "Explore the
/// galaxy map" card lands on the roadmap.
enum LearnView { course, roadmap }

final learnViewProvider = StateProvider<LearnView>((ref) => LearnView.course);
