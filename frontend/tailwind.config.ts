import type { Config } from "tailwindcss";

// Every value here is lifted directly from the Lumora Figma Design Spec
// (Sections 2-4): colour tokens, radius scale, spacing, and type scale.
const config: Config = {
  content: ["./app/**/*.{ts,tsx}", "./components/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        // Brand
        purple: { DEFAULT: "#6C3FC5", dark: "#3A1F8A", light: "#EDE7F6" },
        amber: { DEFAULT: "#F5A623", light: "#FFF8E7" },
        teal: { DEFAULT: "#00C2A8", light: "#E0FAF7" },
        coral: { DEFAULT: "#FF5C5C", light: "#FFE8E8" },
        // Neutrals
        cream: "#FAFAF7",
        gray: {
          50: "#F5F5F5",
          100: "#EBEBEB",
          300: "#CCCCCC",
          500: "#9090A0",
        },
        ink: "#1A1A2E", // "Dark" — primary text
        slatey: "#4A4A6A", // "Mid" — secondary text
        space: "#0F0F24", // deep space for the galaxy map
      },
      borderRadius: {
        xs: "4px",
        sm: "8px",
        md: "12px",
        lg: "16px",
        xl: "24px",
        "2xl": "32px",
        full: "9999px",
      },
      fontFamily: {
        sans: ["var(--font-nunito)", "-apple-system", "BlinkMacSystemFont", "Segoe UI", "sans-serif"],
      },
      fontSize: {
        // name: [size, lineHeight]
        "display-xl": ["32px", "40px"],
        "display-lg": ["28px", "36px"],
        "heading-xl": ["24px", "32px"],
        "heading-lg": ["20px", "28px"],
        "heading-md": ["18px", "26px"],
        "heading-sm": ["16px", "24px"],
        "body-lg": ["16px", "24px"],
        "body-md": ["14px", "22px"],
        "body-sm": ["12px", "18px"],
        "label-lg": ["14px", "20px"],
        "label-md": ["12px", "16px"],
        "label-sm": ["10px", "14px"],
      },
      boxShadow: {
        card: "0 2px 8px rgba(0,0,0,0.08)",
        "card-lg": "0 4px 16px rgba(0,0,0,0.08)",
        skill: "0 4px 16px rgba(108,63,197,0.15)",
        quest: "0 2px 8px rgba(245,166,35,0.15)",
        float: "0 8px 24px rgba(108,63,197,0.25)",
      },
      keyframes: {
        wiggle: {
          "0%,100%": { transform: "rotate(-4deg)" },
          "50%": { transform: "rotate(4deg)" },
        },
        "float-up": {
          "0%": { transform: "translateY(0)", opacity: "0" },
          "30%": { opacity: "1" },
          "100%": { transform: "translateY(-32px)", opacity: "0" },
        },
        pulseRing: {
          "0%,100%": { transform: "scale(1)", opacity: "1" },
          "50%": { transform: "scale(1.08)", opacity: "0.85" },
        },
      },
      animation: {
        wiggle: "wiggle 0.6s ease-in-out",
        "float-up": "float-up 0.8s ease-out forwards",
        pulseRing: "pulseRing 2s ease-in-out infinite",
      },
    },
  },
  plugins: [],
};

export default config;
