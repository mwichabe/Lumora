import 'package:dio/dio.dart';

import '../env.dart';
import '../storage/session_storage.dart';
import 'api_exception.dart';

/// One shared Dio instance for the whole app: attaches the bearer token to
/// every request and clears the session on a genuine 401, mirroring the
/// `request`/`upload` helpers in the web client's lib/api.ts.
class DioClient {
  DioClient._();
  static final DioClient instance = DioClient._();

  /// Called after a 401 clears the session, so the app can sign out.
  void Function()? onUnauthorized;

  static const _sentWithToken = 'lumora.sentWithToken';

  late final Dio dio = Dio(
    BaseOptions(
      baseUrl: Env.apiUrl,
      connectTimeout: const Duration(seconds: 20),
      receiveTimeout: const Duration(seconds: 30),
      sendTimeout: const Duration(seconds: 30),
    ),
  )..interceptors.add(
      InterceptorsWrapper(
        onRequest: (options, handler) async {
          final token = await SessionStorage.instance.getToken();
          if (token != null) {
            options.headers['Authorization'] = 'Bearer $token';
          }
          // Remember which session this request belongs to (see onError).
          options.extra[_sentWithToken] = token;
          handler.next(options);
        },
        onError: (error, handler) async {
          // A 401 only ends the session if it's about the session we have NOW.
          // Requests sent with no token (background polls after signing out,
          // a wrong password on the login form) or with a previous token (in
          // flight across sign-out → sign-in) get 401s too, and acting on those
          // used to delete the freshly issued token — after which every request
          // failed and the app kept bouncing back to home.
          if (error.response?.statusCode == 401) {
            final sentWith = error.requestOptions.extra[_sentWithToken] as String?;
            final current = await SessionStorage.instance.getToken();
            if (sentWith != null && sentWith == current) {
              await SessionStorage.instance.clearSession();
              onUnauthorized?.call();
            }
          }
          handler.next(error);
        },
      ),
    )
    ..interceptors.add(
      InterceptorsWrapper(
        onError: (error, handler) async {
          // The API sleeps when idle, and the first requests after it wakes
          // time out or get a 502/503 from the proxy. Retrying reads here means
          // a cold start costs a few seconds instead of a "Try again" button
          // on every screen. Only GETs: repeating a write could apply it twice.
          final options = error.requestOptions;
          final attempt = (options.extra[_retryAttempt] as int?) ?? 0;
          if (options.method.toUpperCase() != 'GET' || attempt >= _maxRetries || !_isTransient(error)) {
            return handler.next(error);
          }
          options.extra[_retryAttempt] = attempt + 1;
          await Future.delayed(Duration(milliseconds: 1500 * (attempt + 1)));
          try {
            handler.resolve(await dio.fetch(options));
          } on DioException catch (e) {
            handler.next(e);
          }
        },
      ),
    );

  static const _retryAttempt = 'lumora.retryAttempt';
  static const _maxRetries = 2;

  static bool _isTransient(DioException e) {
    switch (e.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.receiveTimeout:
      case DioExceptionType.sendTimeout:
      case DioExceptionType.connectionError:
        return true;
      default:
        final status = e.response?.statusCode ?? 0;
        return status == 502 || status == 503 || status == 504;
    }
  }

  /// Converts a DioException into the app-wide ApiException with a message fit
  /// to show the user: the backend's own `{"error": "..."}` when it sent one,
  /// otherwise a plain-language fallback (matches api.ts).
  ApiException toApiException(DioException e) {
    final status = e.response?.statusCode ?? 0;
    final data = e.response?.data;
    if (data is Map && data['error'] != null && data['error'].toString().isNotEmpty) {
      return ApiException(data['error'].toString(), status);
    }
    return ApiException(_fallbackMessage(e, status), status);
  }

  static String _fallbackMessage(DioException e, int status) {
    switch (e.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.receiveTimeout:
      case DioExceptionType.sendTimeout:
        return 'Lumora is taking longer than usual to respond — it may be waking up. Please try again in a moment.';
      case DioExceptionType.connectionError:
        return "We can't reach Lumora right now. Check your internet connection and try again.";
      default:
        break;
    }
    if (status == 401) return 'Please sign in to continue.';
    if (status == 403) return "You don't have access to that.";
    if (status == 404) return "We couldn't find what you were looking for.";
    if (status == 429) return 'Too many attempts. Please wait a moment and try again.';
    if (status >= 500) return 'Lumora is having trouble right now — it may be waking up. Please try again in a moment.';
    if (status == 0) return "We can't reach Lumora right now. Check your internet connection and try again.";
    return 'Something went wrong. Please try again.';
  }
}
