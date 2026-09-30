# Lumora Mobile

A Flutter (Android + iOS) port of the Lumora web app, talking to the same Go/Fiber
backend (`../backend`). Built with Riverpod, go_router and Dio, and matching the
web app's design spec 1:1 (`lib/core/theme`).

## Features

Everything in the web app has a mobile counterpart:

- Auth (register/login), forgot/reset password, onboarding (language + daily goal)
- Home dashboard, galaxy-map/course view, lesson player (all 8 exercise types)
- Listening & reading comprehension sessions
- Practice hub (vocab quiz, listening/speaking drills, mistake review, daily mix)
- Hearts economy (regen countdown, refill via Paystack)
- Daily quests, companion characters
- All-time leaderboard is superseded in the product by the weekly **League**
  (10 tiers, pods, promotion/demotion zones, group goals, tournament stage,
  casual mode, peer reporting, end-of-week ceremony)
- Notifications (polling)
- Chat (1:1 DMs, images, edit/delete, auto-translation, polling — the backend
  has no WebSocket, so the web app polls too)
- Ideas workspace (board, voting, tasks, threaded discussion, reactions,
  silent brainstorm sessions)
- Exam hub (A1–C2 + Final Mastery): pay for an attempt in-app, then sit the
  exam on the web; certificates (list, detail, share), public certificate
  verification
- Profile (avatar upload, fluency ring, stats, companions), account settings
  (password change, delete account), help/FAQ

## Intentional platform adaptations

- **Exams are web-only**: sitting an exam requires screen-share + camera
  proctoring through browser APIs that don't exist on mobile, so the app
  doesn't offer it. The exam screen lists the levels, takes payment, and links
  out to the web app (`WEB_URL`), where the paid attempt is waiting on the same
  account. The API enforces this too: `/api/exam/paper`, `/start` and `/submit`
  reject requests that don't come from the web app's origin.
- **Payments**: Paystack checkout runs inside the app in a WebView
  (`features/payments/checkout_screen.dart`). The app intercepts Paystack's
  redirect to `/payment/callback`, verifies the reference with the API and
  closes the checkout itself — no browser hop, no manual "I've paid" step. On
  platforms without a WebView (desktop, web) it falls back to the system
  browser with a manual confirm.
- **Chat/notifications**: REST polling, matching the backend (no WebSocket
  exists server-side).
- **Speech**: `flutter_tts` for character voices, `speech_to_text` for
  pronunciation scoring — native equivalents of the web app's Web Speech API.

## Setup

```bash
flutter pub get
```

The default API base URL points at the deployed Render backend. To point at a
local backend instead:

```bash
# Android emulator (10.0.2.2 maps to the host machine's localhost)
flutter run --dart-define=API_URL=http://10.0.2.2:8080

# iOS simulator / desktop / web
flutter run --dart-define=API_URL=http://localhost:8080
```

## Running

```bash
flutter devices        # list available targets
flutter run             # Android/iOS device or emulator
flutter run -d chrome    # web
flutter run -d linux     # Linux desktop
```

## Building

```bash
flutter build apk --release      # Android
flutter build ios --release      # iOS (requires Xcode + signing)
```
