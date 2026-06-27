# Lumora 🦊

A gamified, premium language-learning web app — built as a working MVP with a
**Next.js 14** front end and a **Go (Fiber) MVC** back end.

Lumora the Fennec Fox guides learners through a galaxy-map skill tree, XP, streaks,
hearts, gems, leagues, daily quests, and a Fluency Score (0–1000). The fully seeded
sample course is **Spanish** (starting with *Greetings* and *Ordering at a Café*).

---

## Architecture (MVC)

```
lumora/
├── backend/                  Go + Fiber API (the Model & Controller layers)
│   ├── main.go               App entry, middleware, server start
│   ├── config/               Env-driven configuration
│   ├── models/               GORM models — User, Skill, Lesson, Exercise, Quest…
│   ├── controllers/          Request handlers (auth, lessons, progress, quests…)
│   ├── routes/               Route table → controllers
│   ├── middleware/           JWT auth guard
│   ├── database/             SQLite connect, AutoMigrate, seed data
│   └── utils/                JWT helpers
│
└── frontend/                 Next.js 14 App Router (the View layer)
    ├── app/                  Screens: splash, onboarding, home, learn (galaxy map),
    │                         lesson player, lesson complete, practice, leaderboard, profile
    ├── components/           Button, FoxMascot, widgets, AppShell, BottomTabBar
    └── lib/                  Typed API client, auth context, shared types
```

The Go service owns the **Models** (data) and **Controllers** (business logic);
the Next.js app is the **View**. They communicate over a small JSON REST API.

---

## Prerequisites

- **Go 1.22+** (the backend uses a pure-Go SQLite driver, so **no CGO / no C compiler needed**)
- **Node.js 18+** and npm

> This project was packaged without running `go build` or `npm install` (the build
> environment had no network access). Run the install steps below once on your
> machine and both apps will start normally.

---

## 1. Run the backend (port 8080)

```bash
cd backend
cp .env.example .env          # optional: adjust PORT / JWT_SECRET / CORS_ORIGINS
go mod tidy                   # downloads dependencies
go run .                      # starts the API on http://localhost:8080
```

On first run it creates `lumora.db` (SQLite) and seeds the characters, daily quests,
and the Spanish skill tree automatically.

Health check: `GET http://localhost:8080/api/health` → `{"status":"ok"}`

## 2. Run the frontend (port 3000)

```bash
cd frontend
cp .env.local.example .env.local   # sets NEXT_PUBLIC_API_URL=http://localhost:8080
npm install
npm run dev                        # starts the web app on http://localhost:3000
```

Open **http://localhost:3000**, create an account, choose Spanish, set a daily goal,
and start the first lesson.

---

## API endpoints

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| GET  | `/api/health` | — | Service health check |
| POST | `/api/auth/register` | — | Create account → returns JWT |
| POST | `/api/auth/login` | — | Log in → returns JWT |
| GET  | `/api/auth/me` | ✓ | Current user |
| POST | `/api/auth/setup` | ✓ | Save target language + daily goal |
| GET  | `/api/skills` | ✓ | Galaxy map (skills with unlock/complete state) |
| GET  | `/api/lessons/:id` | ✓ | A lesson with its exercises |
| POST | `/api/lessons/:id/complete` | ✓ | Award XP/gems, update streak & quests |
| GET  | `/api/home` | ✓ | Home dashboard (user, next lesson, quests) |
| GET  | `/api/quests/daily` | ✓ | Today's daily quests |
| GET  | `/api/characters` | ✓ | Companion characters + friendship levels |
| GET  | `/api/leaderboard` | ✓ | Current league standings |

Protected routes expect an `Authorization: Bearer <token>` header. The frontend
stores the token in `localStorage` under `lumora_token` and attaches it automatically.

---

## Tech stack

**Backend:** Go 1.22, Fiber v2, GORM, pure-Go SQLite (`glebarez/sqlite`),
JWT (`golang-jwt/jwt/v5`), bcrypt.

**Frontend:** Next.js 14 (App Router), TypeScript, Tailwind CSS, Framer Motion,
lucide-react, Nunito font. Design tokens (colours, radii, type scale) are encoded
in `tailwind.config.ts` from the Lumora design spec.

---

## ⚠️ Security note

The original Figma design spec document contained a live **Figma personal access
token**. It was **not** included anywhere in this codebase. If you haven't already,
**revoke/rotate that token** in your Figma account settings — once a token is shared
in a document it should be considered compromised.

---

## Notes & scope

- Spanish is the fully-seeded course; other languages appear in onboarding but route
  into the same starter content for this MVP.
- Audio uses the browser's built-in Speech Synthesis (no external TTS service).
- The mascot and character art are rendered as inline SVG / emoji — no binary image
  assets are required to run the app.
# Lumora
