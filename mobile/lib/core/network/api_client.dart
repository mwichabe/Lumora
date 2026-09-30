import 'package:dio/dio.dart';

import '../json_utils.dart';
import '../../models/user.dart';
import '../../models/lesson.dart';
import '../../models/listening_reading.dart';
import '../../models/exam.dart';
import '../../models/notification.dart';
import '../../models/chat.dart';
import '../../models/idea.dart';
import '../../models/quest.dart';
import '../../models/character.dart';
import '../../models/league.dart';
import '../../models/home_data.dart';
import 'dio_client.dart';

/// Bytes + a filename, used for any multipart upload (avatar, chat image,
/// idea attachment) so the client stays platform-agnostic (works from
/// image_picker's XFile on mobile, web and desktop alike).
class UploadFile {
  final List<int> bytes;
  final String filename;
  const UploadFile(this.bytes, this.filename);
}

/// Full REST surface of the Lumora API. One-to-one port of
/// frontend/lib/api.ts — every method name and shape mirrors it so the
/// backend contract stays a single source of truth.
class ApiClient {
  ApiClient._();
  static final ApiClient instance = ApiClient._();

  Dio get _dio => DioClient.instance.dio;

  Future<T> _req<T>(
    String method,
    String path, {
    Map<String, dynamic>? data,
    Map<String, dynamic>? query,
    required T Function(Map<String, dynamic>) map,
  }) async {
    try {
      final res = await _dio.request(
        path,
        data: data,
        queryParameters: query,
        options: Options(method: method),
      );
      return map(asMap(res.data));
    } on DioException catch (e) {
      throw DioClient.instance.toApiException(e);
    }
  }

  Future<void> _reqVoid(String method, String path, {Map<String, dynamic>? data, Map<String, dynamic>? query}) async {
    try {
      await _dio.request(path, data: data, queryParameters: query, options: Options(method: method));
    } on DioException catch (e) {
      throw DioClient.instance.toApiException(e);
    }
  }

  Future<Map<String, dynamic>> _upload(String path, UploadFile file, Map<String, String> fields) async {
    try {
      final form = FormData.fromMap({
        ...fields,
        'file': MultipartFile.fromBytes(file.bytes, filename: file.filename),
      });
      final res = await _dio.post(path, data: form);
      return asMap(res.data);
    } on DioException catch (e) {
      throw DioClient.instance.toApiException(e);
    }
  }

  // --- Auth ---------------------------------------------------------------

  Future<(String token, User user)> register(String email, String password, String name) => _req(
        'POST',
        '/api/auth/register',
        data: {'email': email, 'password': password, 'name': name},
        map: (j) => (asString(j['token']), User.fromJson(asMap(j['user']))),
      );

  Future<(String token, User user)> login(String email, String password) => _req(
        'POST',
        '/api/auth/login',
        data: {'email': email, 'password': password},
        map: (j) => (asString(j['token']), User.fromJson(asMap(j['user']))),
      );

  Future<User> me() => _req('GET', '/api/auth/me', map: (j) => User.fromJson(asMap(j['user'])));

  Future<void> forgotPassword(String email) =>
      _reqVoid('POST', '/api/auth/forgot-password', data: {'email': email});

  Future<void> resetPassword(String token, String password) =>
      _reqVoid('POST', '/api/auth/reset-password', data: {'token': token, 'password': password});

  Future<User> setup(String targetLanguage, int dailyGoalXp, String reason) => _req(
        'POST',
        '/api/auth/setup',
        data: {'targetLanguage': targetLanguage, 'dailyGoalXp': dailyGoalXp, 'reason': reason},
        map: (j) => User.fromJson(asMap(j['user'])),
      );

  Future<User> updateProfile({String? name, String? avatarColor, int? dailyGoalXp}) => _req(
        'PATCH',
        '/api/auth/profile',
        data: {
          if (name != null) 'name': name,
          if (avatarColor != null) 'avatarColor': avatarColor,
          if (dailyGoalXp != null) 'dailyGoalXp': dailyGoalXp,
        },
        map: (j) => User.fromJson(asMap(j['user'])),
      );

  Future<User> uploadAvatar(UploadFile file) async {
    final j = await _upload('/api/auth/avatar', file, const {});
    return User.fromJson(asMap(j['user']));
  }

  Future<User> removeAvatar() =>
      _req('DELETE', '/api/auth/avatar', map: (j) => User.fromJson(asMap(j['user'])));

  Future<void> changePassword(String currentPassword, String newPassword) => _reqVoid(
        'POST',
        '/api/auth/password',
        data: {'currentPassword': currentPassword, 'newPassword': newPassword},
      );

  Future<void> deleteAccount(String password) =>
      _reqVoid('DELETE', '/api/auth/account', data: {'password': password});

  // --- Content & progress ---------------------------------------------------

  Future<HomeData> home() => _req('GET', '/api/home', map: HomeData.fromJson);

  Future<List<Skill>> skills() =>
      _req('GET', '/api/skills', map: (j) => asList(j['skills'], (e) => Skill.fromJson(asMap(e))));

  Future<Lesson> lesson(int id) =>
      _req('GET', '/api/lessons/$id', map: (j) => Lesson.fromJson(asMap(j['lesson'])));

  Future<(int xpEarned, int accuracy, User user, bool firstClear)> completeLesson(
    int id,
    int accuracy,
  ) =>
      _req(
        'POST',
        '/api/lessons/$id/complete',
        data: {'accuracy': accuracy},
        map: (j) => (
          asInt(j['xpEarned']),
          asInt(j['accuracy']),
          User.fromJson(asMap(j['user'])),
          asBool(j['firstClear']),
        ),
      );

  Future<List<ListeningSession>> listeningSessions() => _req(
        'GET',
        '/api/listening',
        map: (j) => asList(j['sessions'], (e) => ListeningSession.fromJson(asMap(e))),
      );

  Future<ListeningSession> listeningSession(int id) => _req(
        'GET',
        '/api/listening/$id',
        map: (j) => ListeningSession.fromJson(asMap(j['session'])),
      );

  Future<(int xpEarned, User user)> completeListening(int id) => _req(
        'POST',
        '/api/listening/$id/complete',
        map: (j) => (asInt(j['xpEarned']), User.fromJson(asMap(j['user']))),
      );

  Future<List<ReadingSession>> readingSessions() => _req(
        'GET',
        '/api/reading',
        map: (j) => asList(j['sessions'], (e) => ReadingSession.fromJson(asMap(e))),
      );

  Future<ReadingSession> readingSession(int id) =>
      _req('GET', '/api/reading/$id', map: (j) => ReadingSession.fromJson(asMap(j['session'])));

  Future<(int xpEarned, User user)> completeReading(int id) => _req(
        'POST',
        '/api/reading/$id/complete',
        map: (j) => (asInt(j['xpEarned']), User.fromJson(asMap(j['user']))),
      );

  Future<(List<String> languages, String active)> enrollments() => _req(
        'GET',
        '/api/enrollments',
        map: (j) => (asStringList(j['languages']), asString(j['active'])),
      );

  Future<(List<String> languages, String active, User user)> enrollLanguage(String language) =>
      _req(
        'POST',
        '/api/enrollments',
        data: {'language': language},
        map: (j) =>
            (asStringList(j['languages']), asString(j['active']), User.fromJson(asMap(j['user']))),
      );

  Future<(List<String> languages, String active, User user)> switchLanguage(String language) =>
      _req(
        'POST',
        '/api/enrollments/active',
        data: {'language': language},
        map: (j) =>
            (asStringList(j['languages']), asString(j['active']), User.fromJson(asMap(j['user']))),
      );

  Future<
      ({
        List<VocabItem> vocab,
        List<Mistake> mistakes,
        int listeningCount,
        int readingCount,
      })> practice() => _req(
        'GET',
        '/api/practice',
        map: (j) => (
          vocab: asList(j['vocab'], (e) => VocabItem.fromJson(asMap(e))),
          mistakes: asList(j['mistakes'], (e) => Mistake.fromJson(asMap(e))),
          listeningCount: asInt(j['listeningCount']),
          readingCount: asInt(j['readingCount']),
        ),
      );

  Future<List<ListeningSession>> practiceListening() => _req(
        'GET',
        '/api/practice/listening',
        map: (j) => asList(j['sessions'], (e) => ListeningSession.fromJson(asMap(e))),
      );

  Future<List<ReadingSession>> practiceReading() => _req(
        'GET',
        '/api/practice/reading',
        map: (j) => asList(j['sessions'], (e) => ReadingSession.fromJson(asMap(e))),
      );

  Future<void> recordMistake({
    required String prompt,
    required String question,
    required String correctAnswer,
  }) =>
      _reqVoid('POST', '/api/mistakes',
          data: {'prompt': prompt, 'question': question, 'correctAnswer': correctAnswer});

  Future<void> resolveMistakes(List<int> ids) =>
      _reqVoid('POST', '/api/mistakes/resolve', data: {'ids': ids});

  Future<(int xpEarned, User user)> completePractice(int xp) => _req(
        'POST',
        '/api/practice/complete',
        data: {'xp': xp},
        map: (j) => (asInt(j['xpEarned']), User.fromJson(asMap(j['user']))),
      );

  // --- Notifications ----------------------------------------------------

  Future<(List<AppNotification> notifications, int unread)> notifications() => _req(
        'GET',
        '/api/notifications',
        map: (j) => (
          asList(j['notifications'], (e) => AppNotification.fromJson(asMap(e))),
          asInt(j['unread']),
        ),
      );

  Future<void> markNotificationsRead() => _reqVoid('POST', '/api/notifications/read');

  Future<int> markNotificationRead(int id) =>
      _req('POST', '/api/notifications/$id/read', map: (j) => asInt(j['unread']));

  Future<int> deleteNotification(int id) =>
      _req('DELETE', '/api/notifications/$id', map: (j) => asInt(j['unread']));

  // --- Exam / certificates ------------------------------------------------

  Future<ExamResult> submitExam({
    required String language,
    required String level,
    required int listening,
    required int reading,
    required int writing,
    required int speaking,
  }) =>
      _req(
        'POST',
        '/api/exam/submit',
        data: {
          'language': language,
          'level': level,
          'listening': listening,
          'reading': reading,
          'writing': writing,
          'speaking': speaking,
        },
        map: ExamResult.fromJson,
      );

  Future<PaymentStatus> paymentStatus() =>
      _req('GET', '/api/payments/status', map: PaymentStatus.fromJson);

  Future<(String? authorizationUrl, String? reference)> initializePayment({
    String? level,
    String? product,
  }) =>
      _req(
        'POST',
        '/api/payments/initialize',
        data: {if (level != null) 'level': level, if (product != null) 'product': product},
        map: (j) => (
          j['authorizationUrl'] == null ? null : asString(j['authorizationUrl']),
          j['reference'] == null ? null : asString(j['reference']),
        ),
      );

  Future<({String status, bool success, String level, String product})> verifyPayment(
    String reference,
  ) =>
      _req(
        'GET',
        '/api/payments/verify',
        query: {'reference': reference},
        map: (j) => (
          status: asString(j['status']),
          success: asBool(j['success']),
          level: asString(j['level']),
          product: asString(j['product']),
        ),
      );

  Future<ExamMeta> examMeta() => _req('GET', '/api/exam/meta', map: ExamMeta.fromJson);

  Future<void> startExam(String level, String language) =>
      _reqVoid('POST', '/api/exam/start', data: {'level': level, 'language': language});

  Future<HeartsStatus> heartsStatus() => _req('GET', '/api/hearts', map: HeartsStatus.fromJson);

  Future<HeartsStatus> loseHeart() => _req('POST', '/api/hearts/lose', map: HeartsStatus.fromJson);

  Future<(String? authorizationUrl, String? reference)> buyHearts() => initializePayment(product: 'hearts');

  Future<ExamPaper> examPaper(String level) =>
      _req('GET', '/api/exam/paper', query: {'level': level}, map: ExamPaper.fromJson);

  Future<List<Certificate>> certificates() => _req(
        'GET',
        '/api/certificates',
        map: (j) => asList(j['certificates'], (e) => Certificate.fromJson(asMap(e))),
      );

  Future<Certificate> certificate(int id) =>
      _req('GET', '/api/certificates/$id', map: (j) => Certificate.fromJson(asMap(j['certificate'])));

  Future<void> deleteCertificate(int id) => _reqVoid('DELETE', '/api/certificates/$id');

  /// Public — no auth header needed.
  Future<CertVerification> verifyCertificate(String serial) async {
    try {
      final res = await Dio(BaseOptions(baseUrl: _dio.options.baseUrl))
          .get('/api/verify/${Uri.encodeComponent(serial)}');
      return CertVerification.fromJson(asMap(res.data));
    } catch (_) {
      return const CertVerification(valid: false);
    }
  }

  // --- Chat -----------------------------------------------------------------

  Future<List<ChatUser>> chatContacts() => _req(
        'GET',
        '/api/chat/contacts',
        map: (j) => asList(j['contacts'], (e) => ChatUser.fromJson(asMap(e))),
      );

  Future<List<ChatThread>> chatThreads() => _req(
        'GET',
        '/api/chat/threads',
        map: (j) => asList(j['threads'], (e) => ChatThread.fromJson(asMap(e))),
      );

  Future<int> chatUnread() => _req('GET', '/api/chat/unread', map: (j) => asInt(j['count']));

  Future<(List<ChatMessage> messages, ChatUser user)> chatMessages(int id) => _req(
        'GET',
        '/api/chat/with/$id',
        map: (j) => (
          asList(j['messages'], (e) => ChatMessage.fromJson(asMap(e))),
          ChatUser.fromJson(asMap(j['user'])),
        ),
      );

  Future<ChatMessage> sendChatMessage(int id, String body) => _req(
        'POST',
        '/api/chat/with/$id',
        data: {'body': body},
        map: (j) => ChatMessage.fromJson(asMap(j['message'])),
      );

  Future<ChatMessage> sendChatImage(int id, UploadFile file, String caption) async {
    final j = await _upload('/api/chat/with/$id/image', file, {'body': caption});
    return ChatMessage.fromJson(asMap(j['message']));
  }

  Future<ChatMessage> editChatMessage(int messageId, String body) => _req(
        'PATCH',
        '/api/chat/messages/$messageId',
        data: {'body': body},
        map: (j) => ChatMessage.fromJson(asMap(j['message'])),
      );

  Future<MessageTranslation?> translateChatMessage(int messageId) => _req(
        'POST',
        '/api/chat/messages/$messageId/translate',
        map: (j) => j['translation'] == null ? null : MessageTranslation.fromJson(asMap(j['translation'])),
      );

  Future<ChatMessage> deleteChatMessage(int messageId) => _req(
        'DELETE',
        '/api/chat/messages/$messageId',
        map: (j) => ChatMessage.fromJson(asMap(j['message'])),
      );

  // --- Quests / characters / leaderboard -------------------------------

  Future<List<UserQuest>> quests() => _req(
        'GET',
        '/api/quests/daily',
        map: (j) => asList(j['quests'], (e) => UserQuest.fromJson(asMap(e))),
      );

  Future<List<CharacterWithFriendship>> characters() => _req(
        'GET',
        '/api/characters',
        map: (j) => asList(j['characters'], (e) => CharacterWithFriendship.fromJson(asMap(e))),
      );

  Future<(String league, List<LeaderRow> rows, int userRank)> leaderboard() => _req(
        'GET',
        '/api/leaderboard',
        map: (j) => (
          asString(j['league']),
          asList(j['rows'], (e) => LeaderRow.fromJson(asMap(e))),
          asInt(j['userRank']),
        ),
      );

  // --- Weekly league ----------------------------------------------------

  Future<LeagueStandings> league() => _req('GET', '/api/league', map: LeagueStandings.fromJson);

  Future<LeagueResult?> leagueResult() => _req(
        'GET',
        '/api/league/result',
        map: (j) => j['result'] == null ? null : LeagueResult.fromJson(asMap(j['result'])),
      );

  Future<void> markLeagueResultSeen(String seasonId) =>
      _reqVoid('POST', '/api/league/result/seen', query: {'season': seasonId});

  Future<List<LeagueHistoryEntry>> leagueHistory() => _req(
        'GET',
        '/api/league/history',
        map: (j) => asList(j['history'], (e) => LeagueHistoryEntry.fromJson(asMap(e))),
      );

  Future<bool> setLeagueCasual(bool enabled) => _req(
        'POST',
        '/api/league/casual',
        data: {'enabled': enabled},
        map: (j) => asBool(j['casual']),
      );

  Future<bool> reportLeagueMember(int id, String reason) => _req(
        'POST',
        '/api/league/report/$id',
        data: {'reason': reason},
        map: (j) => asBool(j['alreadyReported']),
      );

  // --- Ideas workspace ----------------------------------------------------

  Future<IdeaBoard> ideas({String? status, String? tag, String? sort, String? q}) => _req(
        'GET',
        '/api/ideas',
        query: {
          if (status != null && status.isNotEmpty) 'status': status,
          if (tag != null && tag.isNotEmpty) 'tag': tag,
          if (sort != null && sort.isNotEmpty) 'sort': sort,
          if (q != null && q.isNotEmpty) 'q': q,
        },
        map: IdeaBoard.fromJson,
      );

  Future<Idea> createIdea({required String title, String? description, List<String>? tags}) => _req(
        'POST',
        '/api/ideas',
        data: {
          'title': title,
          if (description != null) 'description': description,
          if (tags != null) 'tags': tags,
        },
        map: (j) => Idea.fromJson(asMap(j['idea'])),
      );

  Future<IdeaDetail> idea(int id) => _req('GET', '/api/ideas/$id', map: IdeaDetail.fromJson);

  Future<Idea> updateIdea(
    int id, {
    String? title,
    String? description,
    IdeaStatus? status,
    List<String>? tags,
  }) =>
      _req(
        'PATCH',
        '/api/ideas/$id',
        data: {
          if (title != null) 'title': title,
          if (description != null) 'description': description,
          if (status != null) 'status': ideaStatusToString(status),
          if (tags != null) 'tags': tags,
        },
        map: (j) => Idea.fromJson(asMap(j['idea'])),
      );

  Future<void> deleteIdea(int id) => _reqVoid('DELETE', '/api/ideas/$id');

  Future<Idea> voteIdea(int id, int value) => _req(
        'POST',
        '/api/ideas/$id/vote',
        data: {'value': value},
        map: (j) => Idea.fromJson(asMap(j['idea'])),
      );

  Future<bool> starIdea(int id) =>
      _req('POST', '/api/ideas/$id/star', map: (j) => asBool(j['starred']));

  Future<Idea> archiveIdea(int id, String reason) => _req(
        'POST',
        '/api/ideas/$id/archive',
        data: {'reason': reason},
        map: (j) => Idea.fromJson(asMap(j['idea'])),
      );

  Future<Idea> restoreIdea(int id) =>
      _req('POST', '/api/ideas/$id/restore', map: (j) => Idea.fromJson(asMap(j['idea'])));

  Future<(Idea idea, Idea target)> mergeIdea(int id, int targetId) => _req(
        'POST',
        '/api/ideas/$id/merge',
        data: {'targetId': targetId},
        map: (j) => (Idea.fromJson(asMap(j['idea'])), Idea.fromJson(asMap(j['target']))),
      );

  Future<List<SimilarIdea>> similarIdeas(String q) => _req(
        'GET',
        '/api/ideas/similar',
        query: {'q': q},
        map: (j) => asList(j['similar'], (e) => SimilarIdea.fromJson(asMap(e))),
      );

  Future<ThreadSummary> ideaSummary(int id) =>
      _req('GET', '/api/ideas/$id/summary', map: ThreadSummary.fromJson);

  Future<(IdeaTask task, Idea idea)> createIdeaTask(int id, {String? title, String? sprint}) => _req(
        'POST',
        '/api/ideas/$id/tasks',
        data: {if (title != null) 'title': title, if (sprint != null) 'sprint': sprint},
        map: (j) => (IdeaTask.fromJson(asMap(j['task'])), Idea.fromJson(asMap(j['idea']))),
      );

  Future<IdeaTask> updateIdeaTask(int taskId, {String? status, String? sprint}) => _req(
        'PATCH',
        '/api/ideas/tasks/$taskId',
        data: {if (status != null) 'status': status, if (sprint != null) 'sprint': sprint},
        map: (j) => IdeaTask.fromJson(asMap(j['task'])),
      );

  IdeaMessage _normaliseMessage(IdeaMessage m) => m; // reactions/replies default to [] already

  Future<IdeaThread> ideaMessages(int id) => _req('GET', '/api/ideas/$id/messages', map: IdeaThread.fromJson);

  Future<IdeaMessage> postIdeaMessage(int id, {required String body, int? parentId, String? kind}) => _req(
        'POST',
        '/api/ideas/$id/messages',
        data: {'body': body, if (parentId != null) 'parentId': parentId, if (kind != null) 'kind': kind},
        map: (j) => _normaliseMessage(IdeaMessage.fromJson(asMap(j['message']))),
      );

  Future<IdeaMessage> postIdeaAttachment(
    int id,
    UploadFile file, {
    String? body,
    int? parentId,
    required String kind, // image | voice
    int? duration,
  }) async {
    final j = await _upload('/api/ideas/$id/messages', file, {
      'body': body ?? '',
      'kind': kind,
      'parentId': parentId?.toString() ?? '',
      'duration': duration?.toString() ?? '',
    });
    return _normaliseMessage(IdeaMessage.fromJson(asMap(j['message'])));
  }

  Future<IdeaMessage> editIdeaMessage(int messageId, String body) => _req(
        'PATCH',
        '/api/ideas/messages/$messageId',
        data: {'body': body},
        map: (j) => _normaliseMessage(IdeaMessage.fromJson(asMap(j['message']))),
      );

  Future<IdeaMessage> deleteIdeaMessage(int messageId) => _req(
        'DELETE',
        '/api/ideas/messages/$messageId',
        map: (j) => _normaliseMessage(IdeaMessage.fromJson(asMap(j['message']))),
      );

  Future<IdeaMessage> reactToIdeaMessage(int messageId, String emoji) => _req(
        'POST',
        '/api/ideas/messages/$messageId/react',
        data: {'emoji': emoji},
        map: (j) => _normaliseMessage(IdeaMessage.fromJson(asMap(j['message']))),
      );

  Future<MessageTranslation?> translateIdeaMessage(int messageId) => _req(
        'POST',
        '/api/ideas/messages/$messageId/translate',
        map: (j) => j['translation'] == null ? null : MessageTranslation.fromJson(asMap(j['translation'])),
      );

  Future<BrainstormSession> startBrainstorm(int id, int minutes, String topic) => _req(
        'POST',
        '/api/ideas/$id/brainstorm',
        data: {'minutes': minutes, 'topic': topic},
        map: (j) => BrainstormSession.fromJson(asMap(j['brainstorm'])),
      );

  Future<void> stopBrainstorm(int id) => _reqVoid('DELETE', '/api/ideas/$id/brainstorm');
}
