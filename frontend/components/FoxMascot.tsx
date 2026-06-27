"use client";

import { motion } from "framer-motion";

/**
 * Lumora the Fennec Fox — rendered as inline SVG so the app ships with no image
 * assets. Deep violet fur, amber-gold markings, teal eyes, oversized ears, and
 * a tail tip that glows when `glow` is set (per the brand brief).
 */
export function FoxMascot({
  size = 200,
  glow = false,
  bounce = false,
}: {
  size?: number;
  glow?: boolean;
  bounce?: boolean;
}) {
  return (
    <motion.div
      initial={bounce ? { y: 40, opacity: 0, scale: 0.8 } : false}
      animate={bounce ? { y: 0, opacity: 1, scale: 1 } : undefined}
      transition={{ type: "spring", stiffness: 300, damping: 18 }}
      style={{ width: size, height: size }}
    >
      <svg viewBox="0 0 200 200" width={size} height={size} aria-label="Lumora the fox">
        <defs>
          <radialGradient id="tailGlow" cx="50%" cy="50%" r="50%">
            <stop offset="0%" stopColor="#FFD27D" />
            <stop offset="100%" stopColor="#F5A623" stopOpacity="0" />
          </radialGradient>
          <linearGradient id="fur" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="#7B4AD6" />
            <stop offset="100%" stopColor="#5E33B0" />
          </linearGradient>
        </defs>

        {/* soft ground shadow */}
        <ellipse cx="100" cy="178" rx="46" ry="9" fill="#0F0F24" opacity="0.15" />

        {/* glowing tail */}
        {glow && <circle cx="150" cy="120" r="40" fill="url(#tailGlow)" />}
        <path
          d="M150 150 C176 140 178 104 156 96 C168 116 150 136 132 138 Z"
          fill="#6C3FC5"
          stroke="#0F0F24"
          strokeWidth="2.5"
        />
        <path d="M158 100 C168 112 158 128 146 132" fill="#F5A623" opacity="0.9" />

        {/* body */}
        <ellipse cx="100" cy="132" rx="40" ry="34" fill="url(#fur)" stroke="#0F0F24" strokeWidth="2.5" />
        <ellipse cx="100" cy="142" rx="22" ry="18" fill="#FFF1DD" />

        {/* ears */}
        <path d="M70 70 L58 28 L92 60 Z" fill="#6C3FC5" stroke="#0F0F24" strokeWidth="2.5" strokeLinejoin="round" />
        <path d="M130 70 L142 28 L108 60 Z" fill="#6C3FC5" stroke="#0F0F24" strokeWidth="2.5" strokeLinejoin="round" />
        <path d="M70 64 L63 38 L84 58 Z" fill="#F5A623" />
        <path d="M130 64 L137 38 L116 58 Z" fill="#F5A623" />

        {/* head */}
        <ellipse cx="100" cy="92" rx="38" ry="33" fill="url(#fur)" stroke="#0F0F24" strokeWidth="2.5" />
        {/* cheek / muzzle */}
        <path d="M100 78 C84 78 74 92 78 104 C84 116 116 116 122 104 C126 92 116 78 100 78 Z" fill="#FFF1DD" />

        {/* eyes */}
        <circle cx="86" cy="90" r="8" fill="#fff" />
        <circle cx="114" cy="90" r="8" fill="#fff" />
        <circle cx="87" cy="91" r="5" fill="#00C2A8" />
        <circle cx="115" cy="91" r="5" fill="#00C2A8" />
        <circle cx="89" cy="89" r="1.6" fill="#fff" />
        <circle cx="117" cy="89" r="1.6" fill="#fff" />

        {/* nose + smile */}
        <path d="M96 100 L104 100 L100 105 Z" fill="#0F0F24" />
        <path d="M100 105 C100 110 94 112 90 109" fill="none" stroke="#0F0F24" strokeWidth="2" strokeLinecap="round" />
        <path d="M100 105 C100 110 106 112 110 109" fill="none" stroke="#0F0F24" strokeWidth="2" strokeLinecap="round" />

        {/* amber cheek markings */}
        <circle cx="76" cy="100" r="4" fill="#F5A623" opacity="0.55" />
        <circle cx="124" cy="100" r="4" fill="#F5A623" opacity="0.55" />
      </svg>
    </motion.div>
  );
}
