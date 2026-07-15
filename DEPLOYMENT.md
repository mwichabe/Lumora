# Lumora — Deployment Guide (Render + Netlify)

Backend (Go + SQLite) → **Render**. Frontend (Next.js) → **Netlify**.

The two hosts reference each other's URLs, so deployment is a loop: Render first
with placeholder URLs, then Netlify, then back to Render with the real URL.
Budget ~30 minutes.

---

## Step 0 — Push to GitHub first

Both platforms deploy from a Git repo, so nothing works until this is pushed.

```bash
cd /home/mwichabe/projects/lumora
git add .
git commit -m "Add Render blueprint and Netlify config"
git push origin main
```

Confirm your secrets are NOT in the push — `backend/.env` is gitignored and must
stay that way. Verify with:

```bash
git ls-files | grep -E "\.env$|\.env\.local$|\.db$"   # must print nothing
```

---

## PART 1 — Backend + Database on Render

### Step 1.1 — Create the Blueprint

1. Go to <https://dashboard.render.com> → **New** → **Blueprint**.
2. Connect your GitHub account and pick the `lumora` repo.
3. Render detects `render.yaml` and shows the `lumora-api` service. Click **Apply**.

### Step 1.2 — The database (there's nothing to create)

There is no separate database to provision. Lumora uses **SQLite**, a single file
on the `lumora-data` disk that `render.yaml` mounts at `/var/data`. On first boot
the app auto-migrates the schema and seeds starter content (skills, lessons,
characters) into `/var/data/lumora.db`.

**This disk is why the service must be on the paid `starter` plan (~$7/mo).**
Render only allows disks on paid instances. On a free instance the filesystem is
ephemeral — every deploy and idle-spindown silently wipes all accounts, progress
and certificates, while the app keeps *appearing* to work because it re-seeds
itself. If that cost is a blocker, tell me and I'll migrate you to Render's free
Postgres tier instead (needs a GORM driver swap + object storage for avatars).

### Step 1.3 — Set the environment variables

`render.yaml` already handles some vars; the rest you enter by hand.

**Already set for you — do not touch:**

| Variable | Value | Why |
|---|---|---|
| `DB_PATH` | `/var/data/lumora.db` | Must sit on the disk |
| `UPLOADS_DIR` | `/var/data/uploads` | Must sit on the disk |
| `JWT_SECRET` | auto-generated | Changing it logs out every user |
| `GO_VERSION` | `1.22` | Matches `go.mod` |
| `SMTP_PORT`, `SMTP_FROM_NAME`, `EXAM_PRICE_KES`, `KES_PER_USD`, `HEART_REGEN_MINUTES` | defaults | Tune if you like |
| `PORT` | *(injected by Render)* | Never set this yourself |

**You must add these** in the Render dashboard → `lumora-api` → **Environment**.
Copy the values from your local `backend/.env`:

| Variable | Value to enter | Notes |
|---|---|---|
| `CORS_ORIGINS` | `https://<your-site>.netlify.app` | Placeholder for now — fixed in Step 3.1. **No trailing slash.** |
| `APP_URL` | `https://<your-site>.netlify.app` | Same. Used in emails + Paystack callback |
| `SMTP_HOST` | `smtp.gmail.com` | From your `.env` |
| `SMTP_USER` | your Gmail address | From your `.env` |
| `SMTP_PASS` | your 16-char App Password | From your `.env` |
| `SMTP_FROM` | same as `SMTP_USER` | Gmail requires these to match |
| `PAYSTACK_SECRET_KEY` | `sk_test_…` → `sk_live_…` | Your `.env` has **test** keys. Test keys = fake payments. Swap to live when ready to charge |
| `PAYSTACK_PUBLIC_KEY` | `pk_test_…` → `pk_live_…` | Same |

Leave `SMTP_HOST` empty if you want to disable email entirely; the app skips it.

### Step 1.4 — Verify the backend

Wait for the deploy to go green, then:

```bash
curl https://lumora-api.onrender.com/api/health
# {"service":"lumora","status":"ok"}
```

If it hangs ~30s the first time, that's a Render cold start, not a bug.

---

## PART 2 — Frontend on Netlify

### Step 2.1 — Create the site

1. <https://app.netlify.com> → **Add new site** → **Import an existing project**.
2. Pick GitHub → the `lumora` repo.
3. `netlify.toml` supplies base/command/publish and the Next.js plugin.
   **Leave the build settings blank** — don't override them in the UI.

### Step 2.2 — Set env vars BEFORE the first build

This ordering matters. `NEXT_PUBLIC_*` values are **inlined into the JS bundle at
build time**, not read at runtime. If you build without them, the app ships with
`undefined` as its API URL and every request fails — and editing the var later
does nothing until you trigger a **redeploy**.

Site configuration → **Environment variables** → add:

| Variable | Value |
|---|---|
| `NEXT_PUBLIC_API_URL` | `https://lumora-api.onrender.com` (your real Render URL, no trailing slash) |
| `NEXT_PUBLIC_SITE_URL` | `https://<your-site>.netlify.app` (your real Netlify URL) |

Then **Deploy site**.

### Step 2.3 — Note your real URL

Netlify assigns something like `random-name-123.netlify.app`. Rename it under
Site configuration → **Site details** → **Change site name**, or attach a custom
domain. Whatever you settle on is the URL used in Step 3.1 — changing it later
means redoing that step.

---

## PART 3 — Close the loop

### Step 3.1 — Point Render at the real Netlify URL

Back in Render → `lumora-api` → Environment, replace the placeholders:

- `CORS_ORIGINS` = `https://<your-real-site>.netlify.app`
- `APP_URL` = `https://<your-real-site>.netlify.app`

Saving triggers an automatic redeploy. **Exact origin, no trailing slash, no
path** — the browser matches this string literally, and a mismatch shows up as
CORS errors on every API call.

### Step 3.2 — Paystack webhook

Paystack dashboard → **Settings → API Keys & Webhooks** → Webhook URL:

```
https://lumora-api.onrender.com/api/paystack/webhook
```

Without this, a user who closes the tab mid-payment gets charged but never
unlocked — the webhook is what fulfils the purchase server-side.

### Step 3.3 — If `NEXT_PUBLIC_SITE_URL` was a guess

If you renamed the site after building, redeploy Netlify so the corrected value
gets inlined: **Deploys** → **Trigger deploy** → **Clear cache and deploy site**.

---

## Step 4 — Smoke test

- [ ] `curl https://lumora-api.onrender.com/api/health` → `{"status":"ok"}`
- [ ] Netlify site loads; no CORS errors in the browser console
- [ ] Sign up → welcome email arrives → account persists
- [ ] Upload a profile photo → still there after a Render redeploy (proves the disk)
- [ ] Complete a lesson → progress persists across reload
- [ ] Buy an exam attempt → receipt email + exam unlocks
- [ ] Camera/screen-share prompts appear in the proctored exam (needs HTTPS — both hosts give you this)

---

## Troubleshooting

| Symptom | Cause |
|---|---|
| CORS errors on every API call | `CORS_ORIGINS` doesn't exactly match the Netlify origin (trailing slash, `http` vs `https`) |
| API calls go to `undefined/api/...` | `NEXT_PUBLIC_API_URL` was missing at build time — set it and **redeploy** |
| Netlify build fails on `output: 'standalone'` | It was removed from `next.config.js`; don't add it back — it's incompatible with Netlify |
| All users logged out after a deploy | `JWT_SECRET` changed |
| Avatars/accounts vanish after deploy | Service isn't on a paid plan, so the disk isn't attached |
| First request takes ~30s | Render cold start on idle |
| Payment succeeds but exam stays locked | Webhook URL not set in Paystack (Step 3.2) |
