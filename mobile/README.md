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
- Proctored proficiency exam (A1–C2 + Final Mastery), Paystack payment gate,
  certificates (list, detail, share), public certificate verification
- Profile (avatar upload, fluency ring, stats, companions), account settings
  (password change, delete account), help/FAQ

## Intentional platform adaptations

- **Exam proctoring**: the web app requires both screen-share and camera via
  browser APIs that don't exist on mobile. The app requires the **camera**
  (shown as a PiP preview during the exam) and uses app-lifecycle detection
  (backgrounding the app) as the mobile equivalent of the web's "tab switch
  ends the exam" rule.
- **Payments**: Paystack checkout opens in the system browser
  (`url_launcher`), same as the web app's redirect. Returning to the app
  automatically isn't guaranteed without platform-specific deep-link setup, so
  the exam/hearts screens include an explicit "I've paid — continue" action
  that re-checks payment status.
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
