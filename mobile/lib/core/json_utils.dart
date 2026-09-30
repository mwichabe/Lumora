/// Defensive JSON coercion helpers.
///
/// The Go backend marshals a nil slice/pointer as `null` rather than `[]`/omitted,
/// so every list/primitive read from the API is coerced here instead of trusting
/// the JSON shape at each call site (mirrors the intent of `normaliseMessage` in
/// the web client's lib/api.ts).
int asInt(dynamic v, [int def = 0]) {
  if (v == null) return def;
  if (v is int) return v;
  if (v is double) return v.toInt();
  if (v is String) return int.tryParse(v) ?? def;
  return def;
}

double asDouble(dynamic v, [double def = 0]) {
  if (v == null) return def;
  if (v is double) return v;
  if (v is int) return v.toDouble();
  if (v is String) return double.tryParse(v) ?? def;
  return def;
}

String asString(dynamic v, [String def = '']) => v == null ? def : v.toString();

bool asBool(dynamic v, [bool def = false]) => v is bool ? v : def;

List<T> asList<T>(dynamic v, T Function(dynamic) mapper) {
  if (v is! List) return <T>[];
  return v.map(mapper).toList();
}

List<String> asStringList(dynamic v) => asList<String>(v, (e) => asString(e));

Map<String, dynamic> asMap(dynamic v) =>
    v is Map ? Map<String, dynamic>.from(v) : <String, dynamic>{};

Map<String, int> asIntMap(dynamic v) {
  final m = asMap(v);
  return m.map((k, val) => MapEntry(k, asInt(val)));
}

Map<String, double> asDoubleMap(dynamic v) {
  final m = asMap(v);
  return m.map((k, val) => MapEntry(k, asDouble(val)));
}
