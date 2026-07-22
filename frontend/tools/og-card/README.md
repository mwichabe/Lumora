# Share card and app icons

Source for the images in `public/`:

| File | Size | Used by |
| --- | --- | --- |
| `og.png` | 1200×630 | Open Graph + Twitter card (`app/layout.tsx`) |
| `icon-192.png`, `icon-512.png` | square | PWA manifest, favicon fallback |
| `icon-maskable-512.png` | 512×512 | Android adaptive icon (`purpose: maskable`) |
| `apple-touch-icon.png` | 180×180 | iOS home screen |

They're committed as PNGs rather than generated at request time, so a social
scraper never waits on a serverless function and Netlify serves them straight
from the CDN.

## Why PNG and not the SVG logo

The card used to point at `/logo.svg`. The file exists and looks fine in a
browser — but **no major scraper renders SVG**: Facebook, X, LinkedIn, Slack,
WhatsApp and iMessage all skip it, so shared links appeared with no image.
Open Graph in practice means PNG or JPEG.

The same applies to the PWA icons: Android's install prompt and maskable
adaptive icons need a raster.

## Regenerating

Needs Google Chrome and network access (the script fetches Nunito, the brand
font, from Google Fonts):

```bash
cd frontend/tools/og-card
npm install puppeteer-core        # not a project dependency; this is a one-off tool
node render.mjs
cp og.png icon-*.png apple-touch-icon.png ../../public/
```

Edit `card.html` to change the artwork — it's a plain 1200×630 HTML page using
the same colours, type and mascot SVG as the app. The starfield uses a seeded
PRNG so re-running produces the same sky rather than a different one each time.

## Checking it after deploy

Scrapers cache aggressively; force a re-fetch rather than trusting what you see:

- Facebook — <https://developers.facebook.com/tools/debug/>
- X — <https://cards-dev.twitter.com/validator>
- LinkedIn — <https://www.linkedin.com/post-inspector/>

`NEXT_PUBLIC_SITE_URL` must be set at build time on Netlify. It feeds
`metadataBase`, which is what turns `/og.png` into the absolute URL scrapers
require — a relative path is silently ignored by most of them.
