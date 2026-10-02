/**
 * Email format rules, shared with the backend (backend/utils/email_validate.go)
 * and mobile (mobile/lib/core/validation.dart) — keep the three in step.
 *
 * - local part: letters, digits and !#$%&'*+/=?^_`{|}~- in dot-separated runs
 *   (no leading, trailing or doubled dots), at most 64 characters
 * - domain: labels of letters, digits and inner hyphens, then a 2–63 letter TLD
 * - whole address at most 254 characters
 */
const EMAIL_PATTERN =
  /^[a-z0-9!#$%&'*+/=?^_`{|}~-]+(\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@([a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,63}$/;

export function normaliseEmail(email: string): string {
  return email.trim().toLowerCase();
}

/** What's wrong with an address, or "" when it's valid. */
export function emailError(raw: string): string {
  const email = normaliseEmail(raw);
  if (!email) return "Enter your email address.";
  if (email.length > 254) return "That email address is too long.";
  const at = email.split("@");
  if (at.length !== 2) return "Enter a valid email address, like name@example.com.";
  if (at[0].length > 64) return "The part before the @ is too long.";
  if (!EMAIL_PATTERN.test(email)) {
    return "Enter a valid email address, like name@example.com.";
  }
  return "";
}

export const isValidEmail = (email: string) => emailError(email) === "";
