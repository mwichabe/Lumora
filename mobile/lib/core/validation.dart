/// Email format rules, shared with the backend (backend/utils/email_validate.go)
/// and web (frontend/lib/validation.ts) — keep the three in step.
///
///  * local part: letters, digits and !#$%&'*+/=?^_`{|}~- in dot-separated
///    runs (no leading, trailing or doubled dots), at most 64 characters
///  * domain: labels of letters, digits and inner hyphens, then a 2–63 letter TLD
///  * whole address at most 254 characters
final _emailPattern = RegExp(
  r"^[a-z0-9!#$%&'*+/=?^_`{|}~-]+(\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*"
  r'@([a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,63}$',
);

String normaliseEmail(String email) => email.trim().toLowerCase();

/// What's wrong with an address, or null when it's valid. Usable directly as a
/// TextFormField validator.
String? emailError(String? raw) {
  final email = normaliseEmail(raw ?? '');
  if (email.isEmpty) return 'Enter your email address.';
  if (email.length > 254) return 'That email address is too long.';
  final parts = email.split('@');
  if (parts.length != 2) return 'Enter a valid email address, like name@example.com.';
  if (parts[0].length > 64) return 'The part before the @ is too long.';
  if (!_emailPattern.hasMatch(email)) return 'Enter a valid email address, like name@example.com.';
  return null;
}
