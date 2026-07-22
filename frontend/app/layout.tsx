import type { Metadata, Viewport } from "next";
import { Nunito } from "next/font/google";
import "./globals.css";
import { AuthProvider } from "@/lib/auth";

const nunito = Nunito({
  subsets: ["latin"],
  weight: ["400", "600", "700", "800"],
  variable: "--font-nunito",
});

const siteUrl = process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000";

export const metadata: Metadata = {
  metadataBase: new URL(siteUrl),
  title: {
    default: "Lumora — Learn a language. Fall in love with it.",
    template: "%s · Lumora",
  },
  description:
    "A next-generation language learning app where every lesson is an adventure.",
  applicationName: "Lumora",
  manifest: "/manifest.webmanifest",
  // No `icons` field: app/icon.svg is picked up by Next's file convention,
  // which then suppresses this metadata field entirely (verified — an `apple`
  // entry here emits no link tag at all). The Apple touch icon is served from
  // public/apple-touch-icon.png instead, which iOS requests from the site root
  // by convention whether or not a <link> advertises it. Android and the
  // install prompt get their icons from app/manifest.ts.
  appleWebApp: { capable: true, title: "Lumora", statusBarStyle: "default" },

  // The share card. This was pointing at /logo.svg, which renders nowhere:
  // Facebook, X, LinkedIn, Slack, WhatsApp and iMessage all ignore SVG, so
  // every shared link appeared with no image at all. It's a 1200x630 PNG now —
  // the size every scraper expects — showing Lumora the fennec fox.
  openGraph: {
    type: "website",
    siteName: "Lumora",
    url: siteUrl,
    title: "Lumora — Learn a language. Fall in love with it.",
    description:
      "A next-generation language learning app where every lesson is an adventure.",
    images: [
      {
        url: "/og.png",
        width: 1200,
        height: 630,
        type: "image/png",
        alt: "Lumora the fennec fox beside the words: Learn a language. Fall in love with it.",
      },
    ],
  },
  twitter: {
    // summary_large_image, not summary: the latter crops to a small square
    // thumbnail and throws away the artwork.
    card: "summary_large_image",
    title: "Lumora — Learn a language. Fall in love with it.",
    description:
      "A next-generation language learning app where every lesson is an adventure.",
    images: [
      {
        url: "/og.png",
        alt: "Lumora the fennec fox beside the words: Learn a language. Fall in love with it.",
      },
    ],
  },
  robots: { index: true, follow: true },
};

export const viewport: Viewport = {
  width: "device-width",
  initialScale: 1,
  maximumScale: 1,
  themeColor: "#6C3FC5",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className={nunito.variable}>
      <body className="font-sans text-ink">
        <AuthProvider>{children}</AuthProvider>
      </body>
    </html>
  );
}
