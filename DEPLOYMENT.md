# Lumora — Deployment Guide (Render + Neon + Netlify), free tier

| Piece | Host | Cost |
|---|---|---|
| Go API | Render web service, free plan | $0 |
| Postgres | **Neon** free tier | $0 |
| Next.js client | Netlify | $0 |

**Why the database isn't on Render:** Render's own free Postgres is **deleted 30
days after creation** (+14 day grace), with no backups. Neon's free tier has no
expiry, so your users' data doesn't sit on a countdown. Your API still runs on
Render — only the database lives elsewhere, reached over `DATABASE_URL`.

**What the free tier costs you** (not money):
- The API **spins down after ~15 min idle**; the next request takes ~50s to wake.
- Neon free is 0.5 GB. Avatars are downscaled to ~15-25 KB each, so that's
  thousands of users before it matters.

---

## Step 0 — Push to GitHub ✅ done

The repo is at `github.com/mwichabe/Lumora` and `render.yaml` is committed on
`main`, so Render can already see the blueprint. Re-run this after any change:

```bash
cd /home/mwichabe/projects/lumora
git add . && git commit -m "..." && git push origin main
```

Confirm no secrets ever go with it — this must print nothing:

```bash
git ls-files | grep -E "\.env$|\.env\.local$|\.db$"
```

**Verified against a fresh clone of your GitHub repo:** `render.yaml` is present,
no `.env` leaked, the blueprint's `buildCommand` compiles, and `startCommand`
boots on Render's `PORT=10000` and answers `healthCheckPath` with
`{"status":"ok"}`. The blueprint works as written.

---

## PART 1 — Database on Neon

### Step 1.1 — Create the project

1. Sign up at <https://neon.com> (GitHub login works).
2. **Create project** → name it `lumora` → region **Europe (Frankfurt)** to match
   the Render region and keep latency low.
3. Neon shows a **connection string**. Copy it. It looks like:

```
postgresql://lumora_owner:npg_XXXXXXXX@ep-cool-name-12345.eu-central-1.aws.neon.tech/lumora?sslmode=require
```

Keep `?sslmode=require` on the end — Neon rejects unencrypted connections.
This whole string is the `DATABASE_URL` you paste into Render in Step 2.3.

### Step 1.2 — There are no tables to create

Don't run any SQL. On first boot the API auto-migrates all 24 tables and seeds
the starter content (182 skills, 187 lessons, 1186 exercises, 9 characters).
Verified locally against Postgres 16: migrate + seed takes ~14s, and restarts
skip the seed (guarded by row counts, so content is never duplicated).

---

## PART 2 — Backend on Render

### Step 2.1 — Create the Blueprint

1. <https://dashboard.render.com> → **New** → **Blueprint**.
2. Connect GitHub → pick the `lumora` repo.
3. Render reads `render.yaml` and shows `lumora-api` on the **free** plan → **Apply**.

### Step 2.2 — Where to paste each variable

Render → `lumora-api` → **Environment**. These are the vars marked `sync: false`
in `render.yaml` — Render can't know them, so you enter them by hand.

Values come from your local `backend/.env`.

| Variable | What to paste | Where it comes from |
|---|---|---|
| `DATABASE_URL` | `postgresql://...neon.tech/lumora?sslmode=require` | Neon, Step 1.1 |
| `CORS_ORIGINS` | `https://<your-site>.netlify.app` | Placeholder now, fixed in Step 4.1. **No trailing slash** |
| `APP_URL` | `https://<your-site>.netlify.app` | Same |
| `SMTP_HOST` | `smtp.gmail.com` | your `.env` |
| `SMTP_USER` | your Gmail address | your `.env` |
| `SMTP_PASS` | your 16-char App Password | your `.env` |
| `SMTP_FROM` | same as `SMTP_USER` | Gmail requires the match |
| `PAYSTACK_SECRET_KEY` | `sk_test_…` → **`sk_live_…`** | your `.env` has **test** keys |
| `PAYSTACK_PUBLIC_KEY` | `pk_test_…` → **`pk_live_…`** | your `.env` has **test** keys |

Set automatically — **don't touch**: `JWT_SECRET` (generated; changing it logs
out every user), `SMTP_PORT`, `SMTP_FROM_NAME`, `EXAM_PRICE_KES`, `KES_PER_USD`,
`HEART_REGEN_MINUTES`. **Never set `PORT`** — Render injects it.

There is no Go version setting. Render's native Go runtime always builds with
the latest stable Go and can't be pinned (only a Docker deploy can); `go.mod`
requires >= 1.25.0, which latest stable satisfies. `GO_VERSION` is not a real
Render setting — it would be silently ignored.

Not used any more: `DB_PATH` and `UPLOADS_DIR` are local-dev only. In production
`DATABASE_URL` takes over and nothing is written to the filesystem.

### Step 2.3 — Verify

```bash
curl https://lumora-api.onrender.com/api/health
# {"service":"lumora","status":"ok"}
```

First deploy takes a few minutes (Go build + migrate + seed). A ~50s delay on
later requests is the free-tier spindown, not a bug.

---

## PART 3 — Frontend on Netlify

### Step 3.1 — Create the site

<https://app.netlify.com> → **Add new site** → **Import an existing project** →
GitHub → `lumora`. `netlify.toml` supplies base/command/publish and the Next.js
plugin, so **leave the build settings blank**.

### Step 3.2 — Paste these BEFORE the first build

Site configuration → **Environment variables**:

| Variable | Value |
|---|---|
| `NEXT_PUBLIC_API_URL` | `https://lumora-api.onrender.com` |
| `NEXT_PUBLIC_SITE_URL` | `https://<your-site>.netlify.app` |

This ordering is not optional. `NEXT_PUBLIC_*` values are **inlined into the JS
bundle at build time**, not read at runtime. Build without them and the app ships
`undefined` as its API URL; editing the var afterwards changes nothing until you
**redeploy**.

Then **Deploy site**, and rename it under Site details → **Change site name** if
you want something friendlier than `random-name-123.netlify.app`.

---

## PART 4 — Close the loop

### Step 4.1 — Point Render at the real Netlify URL

Render → `lumora-api` → Environment:

- `CORS_ORIGINS` = `https://<your-real-site>.netlify.app`
- `APP_URL` = `https://<your-real-site>.netlify.app`

Saving redeploys automatically. Exact origin, **no trailing slash** — `APP_URL`
is concatenated with paths (`APP_URL + "/payment/callback"`), so a stray slash
produces broken `//payment/callback` links inside real emails.

### Step 4.2 — Paystack webhook

Paystack → **Settings → API Keys & Webhooks** → Webhook URL:

```
https://lumora-api.onrender.com/api/paystack/webhook
```

Without it, a user who closes the tab mid-payment is charged but never unlocked.

### Step 4.3 — Redeploy Netlify if you renamed the site

Deploys → **Trigger deploy** → **Clear cache and deploy site**, so the corrected
`NEXT_PUBLIC_SITE_URL` gets inlined.

---

## Step 5 — Smoke test

- [ ] `curl https://lumora-api.onrender.com/api/health` → `{"status":"ok"}`
- [ ] Netlify site loads, no CORS errors in the console
- [ ] Sign up → welcome email → account persists
- [ ] Upload a profile photo → **still there after a Render redeploy** (this is
      the thing the free tier used to break)
- [ ] Complete a lesson → progress persists
- [ ] Buy an exam attempt → receipt email + unlock
- [ ] Proctored exam prompts for camera + screen share (needs HTTPS — both hosts give it)

---

## Troubleshooting

| Symptom | Cause |
|---|---|
| CORS errors on every call | `CORS_ORIGINS` ≠ Netlify origin exactly (trailing slash, http vs https) |
| Calls go to `undefined/api/...` | `NEXT_PUBLIC_API_URL` missing at build time → set it and **redeploy** |
| `SSL is not enabled on the server` | `?sslmode=require` missing from `DATABASE_URL` |
| First request takes ~50s | Free-tier spindown after 15 min idle |
| All users logged out after deploy | `JWT_SECRET` changed |
| Payment succeeds, exam locked | Paystack webhook URL not set (Step 4.2) |
| Netlify build fails on `output: 'standalone'` | It was removed from `next.config.js` — incompatible with Netlify; don't re-add |

---

## Local development is unchanged

Leave `DATABASE_URL` empty in `backend/.env` and the app still uses the SQLite
file at `DB_PATH`. Same code, same behaviour — only the driver differs.

One local note: avatars now live in the database, so any profile photo uploaded
before this change (its `avatarUrl` still points at the removed `/uploads/…`
path) shows as broken until re-uploaded. Production is unaffected — its database
starts fresh.
