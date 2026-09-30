/// Backend base URL. Override at build/run time with:
///   flutter run --dart-define=API_URL=https://your-api.onrender.com
///
/// Defaults to the deployed Render API so the app works out of the box; point
/// it at http://10.0.2.2:8080 (Android emulator) or http://localhost:8080
/// (iOS simulator / desktop / web) for local backend development.
class Env {
  Env._();

  static const apiUrl = String.fromEnvironment(
    'API_URL',
    defaultValue: 'https://lumora-api-kqsn.onrender.com',
  );
}

/// Resolves a backend-relative media path (avatar/chat/idea attachment URLs)
/// to an absolute URL, mirroring frontend/lib/api.ts `mediaUrl`.
String mediaUrl(String? path) {
  if (path == null || path.isEmpty) return '';
  if (path.startsWith('http')) return path;
  return '${Env.apiUrl}$path';
}
