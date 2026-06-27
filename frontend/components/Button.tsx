"use client";

import { ButtonHTMLAttributes, forwardRef } from "react";
import { motion } from "framer-motion";

type Variant = "primary" | "secondary" | "outline" | "ghost" | "danger";

interface Props extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
  full?: boolean;
  loading?: boolean;
}

// Heights, fills and radii come straight from Design Spec section 5.1.
const base =
  "inline-flex items-center justify-center font-extrabold tracking-wide transition-colors select-none disabled:opacity-40 disabled:cursor-not-allowed";

const variants: Record<Variant, string> = {
  primary: "bg-purple text-white rounded-full h-[52px] px-6 hover:bg-purple-dark",
  secondary: "bg-amber text-ink rounded-full h-[52px] px-6 hover:brightness-95",
  outline:
    "bg-transparent text-purple border-2 border-purple rounded-full h-[52px] px-6 hover:bg-purple-light",
  ghost: "bg-transparent text-purple rounded-md h-11 px-4 hover:bg-purple-light",
  danger: "bg-coral text-white rounded-full h-[52px] px-6 hover:brightness-95",
};

export const Button = forwardRef<HTMLButtonElement, Props>(function Button(
  { variant = "primary", full, loading, className = "", children, ...rest },
  ref
) {
  return (
    <motion.button
      ref={ref}
      whileTap={{ scale: 0.96 }}
      transition={{ type: "spring", stiffness: 400, damping: 30 }}
      className={`${base} ${variants[variant]} ${full ? "w-full" : ""} ${className}`}
      {...(rest as any)}
    >
      {loading ? (
        <span className="h-5 w-5 animate-spin rounded-full border-2 border-white/40 border-t-white" />
      ) : (
        children
      )}
    </motion.button>
  );
});
