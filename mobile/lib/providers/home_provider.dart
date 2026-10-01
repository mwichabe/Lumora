import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../core/network/api_client.dart';
import '../models/home_data.dart';
import 'auth_provider.dart';

class HomeController extends AsyncNotifier<HomeData> {
  @override
  Future<HomeData> build() {
    // Home shows the active course ("continue where you left off", quests),
    // so refetch whenever the learner or their language changes — switching
    // language used to leave the previous course's lesson on this card.
    // Selecting just these two keeps the setUser() in _fetch from re-triggering.
    ref.watch(authProvider.select((a) => (a.user?.id, a.user?.targetLanguage)));
    return _fetch();
  }

  Future<HomeData> _fetch() async {
    final data = await ApiClient.instance.home();
    ref.read(authProvider.notifier).setUser(data.user);
    return data;
  }

  Future<void> refresh() async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(_fetch);
  }
}

final homeProvider = AsyncNotifierProvider<HomeController, HomeData>(HomeController.new);
