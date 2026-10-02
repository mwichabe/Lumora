import 'package:flutter_test/flutter_test.dart';
import 'package:lumora_mobile/core/validation.dart';

void main() {
  test('accepts well-formed addresses', () {
    for (final e in [
      'name@example.com', 'first.last@example.co.uk', 'user+tag@sub.domain.io',
      "o'brien@example.ie", 'a@b.co', 'x_y-z@my-domain.com', ' Name@Example.COM ',
    ]) {
      expect(emailError(e), isNull, reason: e);
    }
  });

  test('rejects malformed addresses', () {
    for (final e in [
      '', 'plainaddress', '@example.com', 'name@', 'name@example', 'name@example.c',
      'name@@example.com', 'name@ex@ample.com', '.name@example.com', 'name.@example.com',
      'na..me@example.com', 'name@-example.com', 'name@example-.com', 'name@example..com',
      'name@.example.com', 'name @example.com', 'name@example.123', 'name@example.com.',
    ]) {
      expect(emailError(e), isNotNull, reason: e);
    }
  });
}
