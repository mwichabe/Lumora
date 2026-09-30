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

  /// Called after a 401 clears the session, so the router can bounce to auth.
  void Function()? onUnauthorized;

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
          handler.next(options);
        },
        onError: (error, handler) async {
          if (error.response?.statusCode == 401) {
            await SessionStorage.instance.clearSession();
            onUnauthorized?.call();
          }
          handler.next(error);
        },
      ),
    );

  /// Converts a DioException into the app-wide ApiException, extracting the
  /// backend's `{"error": "..."}` message when present (matches api.ts).
  ApiException toApiException(DioException e) {
    final status = e.response?.statusCode ?? 0;
    String message = 'request failed ($status)';
    final data = e.response?.data;
    if (data is Map && data['error'] != null) {
      message = data['error'].toString();
    } else if (status == 0) {
      message = 'network error';
    }
    return ApiException(message, status);
  }
}
