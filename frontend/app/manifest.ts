import type { MetadataRoute } from "next";

export default function manifest(): MetadataRoute.Manifest {
  return {
    name: "Lumora — Learn a language",
    short_name: "Lumora",
    description:
      "A next-generation language learning app where every lesson is an adventure.",
    start_url: "/",
    display: "standalone",
    background_color: "#eceaf3",
    theme_color: "#6C3FC5",
    // PNG, not SVG: Android's install prompt and the home-screen icon pipeline
    // are unreliable with SVG, and a maskable icon has to be a raster the
    // launcher can crop to whatever shape the device uses. The maskable variant
    // keeps the fox inside the 80% safe zone so a circular mask doesn't clip
    // his ears off.
    icons: [
      { src: "/icon-192.png", type: "image/png", sizes: "192x192", purpose: "any" },
      { src: "/icon-512.png", type: "image/png", sizes: "512x512", purpose: "any" },
      {
        src: "/icon-maskable-512.png",
        type: "image/png",
        sizes: "512x512",
        purpose: "maskable",
      },
    ],
  };
}
