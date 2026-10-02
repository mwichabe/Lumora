import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../core/network/api_client.dart';
import '../core/storage/session_storage.dart';
import '../models/home_data.dart';
import 'auth_provider.dart';

/// The home screen's data: "continue where you left off", daily goal, quests.
///
/// Home lives in the bottom-nav shell and is never torn down, so it can't rely
/// on being rebuilt to pick up progress. Instead:
///  * it opens from the copy saved on the last visit, then refreshes quietly —
///    the screen never sits on a spinner (or an error) when it has data;
///  * any XP change (a lesson, listening, reading or practice session) refreshes
///    it in the background, so the next lesson and quest bars move on their own.
class HomeController extends AsyncNotifier<HomeData> {
  // Bumped per fetch so a slow, older response can't overwrite a newer one.
  int _generation = 0;
  // While a fetch is running, the setUser() it makes must not trigger another.
  int _inFlight = 0;

  @override
  Future<HomeData> build() async {
    // Home shows the active course, so refetch whenever the learner or their
    // language changes — switching language used to leave the previous
    // course's lesson on this card. Selecting just these two keeps the
    // setUser() in _fetch from re-triggering a rebuild.
    final (userId, lang) = ref.watch(authProvider.select((a) => (a.user?.id, a.user?.targetLanguage)));

    ref.listen(authProvider.select((a) => a.user?.xp), (prev, next) {
      if (_inFlight == 0 && next != null && next != state.valueOrNull?.user.xp) refresh();
    });

    final key = userId == null ? null : 'home_${userId}_$lang';
    final saved = key == null ? null : await _readSaved(key);
    if (saved != null) {
      // Show the saved copy now; the fresh one replaces it when it lands.
      // Deferred to the event queue so it runs after build's result is applied.
      Future(refresh);
      return saved;
    }
    _generation++;
    return _fetch(key);
  }

  String? _key() {
    final user = ref.read(authProvider).user;
    if (user == null) return null;
    return 'home_${user.id}_${user.targetLanguage}';
  }

  Future<HomeData?> _readSaved(String key) async {
    try {
      final json = await SessionStorage.instance.readCache(key);
      return json == null ? null : HomeData.fromJson(json);
    } catch (_) {
      return null; // an unreadable copy is ignored; the network fills the screen
    }
  }

  Future<HomeData> _fetch(String? key) async {
    _inFlight++;
    try {
      final json = await ApiClient.instance.getJson('/api/home');
      final data = HomeData.fromJson(json);
      if (key != null) SessionStorage.instance.writeCache(key, json).catchError((_) {});
      ref.read(authProvider.notifier).setUser(data.user);
      return data;
    } finally {
      _inFlight--;
    }
  }

  /// Fetches fresh data without blanking the screen: whatever is showing stays
  /// until the new response lands, and a failed refresh keeps it too.
  Future<void> refresh() async {
    final gen = ++_generation;
    if (!state.hasValue) state = const AsyncLoading();
    try {
      final data = await _fetch(_key());
      if (gen == _generation) state = AsyncData(data);
    } catch (e, st) {
      if (gen == _generation && !state.hasValue) state = AsyncError(e, st);
    }
  }
}

final homeProvider = AsyncNotifierProvider<HomeController, HomeData>(HomeController.new);
