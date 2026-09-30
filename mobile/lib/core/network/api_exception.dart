class ApiException implements Exception {
  final String message;
  final int status;
  ApiException(this.message, this.status);

  @override
  String toString() => message;
}
