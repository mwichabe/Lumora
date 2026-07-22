"use client";

import {
  useCallback,
  useEffect,
  useLayoutEffect,
  useRef,
  useState,
} from "react";
import { createPortal } from "react-dom";
import { AnimatePresence, motion } from "framer-motion";
import { MoreHorizontal } from "lucide-react";

/**
 * The "…" menu on a message.
 *
 * Two things here are deliberate, because the obvious implementation gets both
 * wrong:
 *
 * 1. The panel is rendered into a portal at document.body and positioned with
 *    `fixed`. Message lists are `overflow-y-auto`, and an absolutely positioned
 *    child of a scrolling container is clipped at its edges — so a menu on the
 *    top message gets sliced in half by the header. A portal escapes the
 *    clipping and every ancestor stacking context with it.
 *
 * 2. The trigger is always visible, not hover-revealed. Hover doesn't exist on
 *    a touchscreen, so `opacity-0 group-hover:opacity-100` means the button is
 *    permanently invisible on a phone and the actions are unreachable. It sits
 *    at low opacity and lifts on hover or focus, which reads as quiet on
 *    desktop and is still tappable everywhere.
 *
 * It flips above the trigger when there isn't room below, and closes on outside
 * click, Escape, scroll or resize.
 */

const GAP = 6; // px between trigger and panel
const EDGE = 8; // keep-off-the-viewport-edge margin

export function ActionMenu({
  label = "Message actions",
  align = "right",
  width = 176,
  children,
}: {
  label?: string;
  align?: "left" | "right";
  width?: number;
  /** Receives a `close` callback so items can dismiss the menu when chosen. */
  children: (close: () => void) => React.ReactNode;
}) {
  const triggerRef = useRef<HTMLButtonElement | null>(null);
  const panelRef = useRef<HTMLDivElement | null>(null);
  const [open, setOpen] = useState(false);
  const [mounted, setMounted] = useState(false);
  const [pos, setPos] = useState<{ top: number; left: number; above: boolean } | null>(
    null
  );

  useEffect(() => setMounted(true), []);

  const place = useCallback(() => {
    const trigger = triggerRef.current;
    if (!trigger) return;
    const r = trigger.getBoundingClientRect();
    // Measured once rendered; the estimate only governs the first frame.
    const h = panelRef.current?.offsetHeight ?? 96;

    const below = r.bottom + GAP;
    const above = below + h > window.innerHeight - EDGE;

    let left = align === "right" ? r.right - width : r.left;
    left = Math.min(Math.max(left, EDGE), window.innerWidth - width - EDGE);

    setPos({ top: above ? r.top - h - GAP : below, left, above });
  }, [align, width]);

  useLayoutEffect(() => {
    if (open) place();
  }, [open, place]);

  useEffect(() => {
    if (!open) return;

    const onPointerDown = (e: PointerEvent) => {
      const t = e.target as Node;
      if (panelRef.current?.contains(t) || triggerRef.current?.contains(t)) return;
      setOpen(false);
    };
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") setOpen(false);
    };
    // Scrolling the thread should dismiss rather than leave the panel stranded
    // where the message used to be. Captured, so it catches the inner
    // scroll container too.
    const onScroll = () => setOpen(false);

    document.addEventListener("pointerdown", onPointerDown);
    document.addEventListener("keydown", onKey);
    window.addEventListener("scroll", onScroll, true);
    window.addEventListener("resize", onScroll);
    return () => {
      document.removeEventListener("pointerdown", onPointerDown);
      document.removeEventListener("keydown", onKey);
      window.removeEventListener("scroll", onScroll, true);
      window.removeEventListener("resize", onScroll);
    };
  }, [open]);

  return (
    <>
      <button
        ref={triggerRef}
        onClick={() => setOpen((o) => !o)}
        aria-label={label}
        aria-haspopup="menu"
        aria-expanded={open}
        className={`rounded-full p-1 transition hover:bg-gray-100 hover:text-slatey focus-visible:opacity-100 ${
          open ? "bg-gray-100 text-slatey opacity-100" : "text-gray-300 opacity-70"
        } hover:opacity-100`}
      >
        <MoreHorizontal size={16} />
      </button>

      {mounted &&
        createPortal(
          <AnimatePresence>
            {open && (
              <motion.div
                ref={panelRef}
                role="menu"
                initial={{ opacity: 0, scale: 0.96, y: pos?.above ? 4 : -4 }}
                animate={{ opacity: 1, scale: 1, y: 0 }}
                exit={{ opacity: 0, scale: 0.96 }}
                transition={{ duration: 0.12 }}
                style={{
                  position: "fixed",
                  top: pos?.top ?? -9999,
                  left: pos?.left ?? -9999,
                  width,
                }}
                className="z-[90] rounded-2xl bg-white p-1.5 shadow-card-lg"
              >
                {children(() => setOpen(false))}
              </motion.div>
            )}
          </AnimatePresence>,
          document.body
        )}
    </>
  );
}

/** A row inside an ActionMenu. */
export function ActionMenuItem({
  icon,
  tone,
  disabled,
  onClick,
  children,
}: {
  icon?: React.ReactNode;
  tone?: "danger";
  disabled?: boolean;
  onClick: () => void;
  children: React.ReactNode;
}) {
  return (
    <button
      role="menuitem"
      onClick={onClick}
      disabled={disabled}
      className={`flex w-full items-center gap-2 rounded-lg px-3 py-2 text-left text-body-sm font-bold transition disabled:opacity-40 ${
        tone === "danger"
          ? "text-coral hover:bg-coral-light"
          : "text-slatey hover:bg-gray-50 hover:text-ink"
      }`}
    >
      {icon}
      {children}
    </button>
  );
}

/** Explanatory text inside an ActionMenu — why an action isn't offered. */
export function ActionMenuNote({ children }: { children: React.ReactNode }) {
  return <p className="px-3 py-1.5 text-label-sm text-gray-500">{children}</p>;
}
