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
  icons: {
    icon: [{ url: "/icon.svg", type: "image/svg+xml" }],
    shortcut: "/icon.svg",
    apple: "/icon.svg",
  },
  appleWebApp: { capable: true, title: "Lumora", statusBarStyle: "default" },
  openGraph: {
    type: "website",
    siteName: "Lumora",
    title: "Lumora — Learn a language. Fall in love with it.",
    description:
      "A next-generation language learning app where every lesson is an adventure.",
    images: ["/logo.svg"],
  },
  twitter: {
    card: "summary",
    title: "Lumora",
    description:
      "A next-generation language learning app where every lesson is an adventure.",
    images: ["/logo.svg"],
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
