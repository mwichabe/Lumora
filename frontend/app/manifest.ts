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
    icons: [
      { src: "/logo.svg", type: "image/svg+xml", sizes: "any", purpose: "any" },
      {
        src: "/logo.svg",
        type: "image/svg+xml",
        sizes: "any",
        purpose: "maskable",
      },
    ],
  };
}
