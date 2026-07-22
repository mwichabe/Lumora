"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { motion } from "framer-motion";
import {
  Crop,
  RotateCcw,
  RotateCw,
  SlidersHorizontal,
  Sparkles,
  X,
} from "lucide-react";
import { Button } from "./Button";

/**
 * Image editor for shared photos: crop to an aspect ratio, rotate, zoom, and
 * apply a colour filter. Exports a JPEG File ready to upload.
 *
 * Built on canvas with no external dependency, and a sibling of CropModal
 * rather than a replacement — that one crops avatars to a circle, which is a
 * different job with a fixed output. This one keeps rectangular aspect ratios
 * and free-form framing.
 *
 * Everything happens client-side, so the server only ever receives the final
 * image. The original never leaves the device.
 */

const VIEW = 300; // on-screen crop frame, longest edge (px)
const OUT_MAX = 1600; // exported longest edge — matches the server's own cap

type Ratio = { label: string; value: number | null };

const RATIOS: Ratio[] = [
  { label: "Free", value: null },
  { label: "1:1", value: 1 },
  { label: "4:3", value: 4 / 3 },
  { label: "16:9", value: 16 / 9 },
  { label: "3:4", value: 3 / 4 },
];

// CSS filter strings, applied live to the preview and replayed onto the canvas
// at export so what you see is what gets sent.
const FILTERS: { label: string; css: string }[] = [
  { label: "None", css: "none" },
  { label: "Punch", css: "saturate(1.4) contrast(1.12)" },
  { label: "Warm", css: "sepia(0.25) saturate(1.25) brightness(1.05)" },
  { label: "Cool", css: "hue-rotate(-12deg) saturate(1.1) brightness(1.04)" },
  { label: "Mono", css: "grayscale(1) contrast(1.08)" },
  { label: "Faded", css: "saturate(0.75) brightness(1.1) contrast(0.92)" },
];

export function ImageEditorModal({
  file,
  title = "Edit photo",
  onCancel,
  onDone,
}: {
  file: File;
  title?: string;
  onCancel: () => void;
  onDone: (f: File) => void;
}) {
  const imgRef = useRef<HTMLImageElement | null>(null);
  const dragRef = useRef<{ x: number; y: number; ox: number; oy: number } | null>(null);

  const [url, setUrl] = useState("");
  const [nat, setNat] = useState({ w: 0, h: 0 });
  const [ratio, setRatio] = useState<Ratio>(RATIOS[0]);
  const [rotation, setRotation] = useState(0); // degrees, multiples of 90
  const [zoom, setZoom] = useState(1);
  const [off, setOff] = useState({ x: 0, y: 0 });
  const [filter, setFilter] = useState(FILTERS[0]);
  const [tab, setTab] = useState<"crop" | "filter">("crop");
  const [busy, setBusy] = useState(false);

  useEffect(() => {
    const u = URL.createObjectURL(file);
    setUrl(u);
    return () => URL.revokeObjectURL(u);
  }, [file]);

  // Rotating by 90° swaps the image's effective width and height.
  const swapped = rotation % 180 !== 0;
  const srcW = swapped ? nat.h : nat.w;
  const srcH = swapped ? nat.w : nat.h;

  // The crop frame: square-ish by default, or the chosen aspect ratio.
  const frame = useMemo(() => {
    const r = ratio.value ?? (srcW && srcH ? srcW / srcH : 1);
    return r >= 1
      ? { w: VIEW, h: Math.round(VIEW / r) }
      : { w: Math.round(VIEW * r), h: VIEW };
  }, [ratio, srcW, srcH]);

  // Scale so the image always covers the frame, then apply the zoom on top.
  const cover = srcW && srcH ? Math.max(frame.w / srcW, frame.h / srcH) : 1;
  const scale = cover * zoom;
  const dispW = srcW * scale;
  const dispH = srcH * scale;
  const left = frame.w / 2 - dispW / 2 + off.x;
  const top = frame.h / 2 - dispH / 2 + off.y;

  // Keep the frame fully covered — no empty edges, whatever the pan.
  const clamp = useCallback(
    (o: { x: number; y: number }) => {
      const l = frame.w / 2 - dispW / 2 + o.x;
      const t = frame.h / 2 - dispH / 2 + o.y;
      let nx = o.x;
      let ny = o.y;
      if (l > 0) nx -= l;
      if (l < frame.w - dispW) nx += frame.w - dispW - l;
      if (t > 0) ny -= t;
      if (t < frame.h - dispH) ny += frame.h - dispH - t;
      return { x: nx, y: ny };
    },
    [frame.w, frame.h, dispW, dispH]
  );

  useEffect(() => {
    setOff((o) => clamp(o));
  }, [clamp]);

  function onDown(e: React.PointerEvent) {
    dragRef.current = { x: e.clientX, y: e.clientY, ox: off.x, oy: off.y };
    (e.target as Element).setPointerCapture?.(e.pointerId);
  }
  function onMove(e: React.PointerEvent) {
    const d = dragRef.current;
    if (!d) return;
    setOff(clamp({ x: d.ox + (e.clientX - d.x), y: d.oy + (e.clientY - d.y) }));
  }
  function onUp() {
    dragRef.current = null;
  }

  function reset() {
    setZoom(1);
    setOff({ x: 0, y: 0 });
    setRotation(0);
    setRatio(RATIOS[0]);
    setFilter(FILTERS[0]);
  }

  /**
   * Export. Draws in two passes: first the rotated source onto a scratch canvas
   * so rotation is baked into pixel space, then the visible crop window from
   * that onto the output canvas with the filter applied.
   */
  function save() {
    const img = imgRef.current;
    if (!img || !nat.w || busy) return;
    setBusy(true);

    // Pass 1 — rotate.
    const rot = document.createElement("canvas");
    rot.width = srcW;
    rot.height = srcH;
    const rctx = rot.getContext("2d");
    if (!rctx) return setBusy(false);
    rctx.translate(srcW / 2, srcH / 2);
    rctx.rotate((rotation * Math.PI) / 180);
    rctx.drawImage(img, -nat.w / 2, -nat.h / 2, nat.w, nat.h);

    // Pass 2 — crop the frame out of it, at the display scale, with the filter.
    const sx = (0 - left) / scale;
    const sy = (0 - top) / scale;
    const sw = frame.w / scale;
    const sh = frame.h / scale;

    const outScale = Math.min(1, OUT_MAX / Math.max(sw, sh));
    const outW = Math.max(1, Math.round(sw * outScale));
    const outH = Math.max(1, Math.round(sh * outScale));

    const out = document.createElement("canvas");
    out.width = outW;
    out.height = outH;
    const ctx = out.getContext("2d");
    if (!ctx) return setBusy(false);
    // A white base so a transparent PNG doesn't come out with black corners
    // once it's encoded as JPEG.
    ctx.fillStyle = "#ffffff";
    ctx.fillRect(0, 0, outW, outH);
    ctx.filter = filter.css;
    ctx.drawImage(rot, sx, sy, sw, sh, 0, 0, outW, outH);

    out.toBlob(
      (blob) => {
        setBusy(false);
        if (!blob) return;
        const name = file.name.replace(/\.[^.]+$/, "") || "photo";
        onDone(new File([blob], `${name}.jpg`, { type: "image/jpeg" }));
      },
      "image/jpeg",
      0.9
    );
  }

  return (
    <div className="fixed inset-0 z-[70] flex items-center justify-center bg-black/70 p-4">
      <motion.div
        initial={{ opacity: 0, scale: 0.96, y: 12 }}
        animate={{ opacity: 1, scale: 1, y: 0 }}
        className="max-h-[92vh] w-full max-w-md overflow-y-auto rounded-2xl bg-white p-5 shadow-card-lg"
      >
        <div className="mb-3 flex items-center justify-between">
          <h3 className="text-heading-sm font-extrabold text-ink">{title}</h3>
          <button
            onClick={onCancel}
            aria-label="Cancel"
            className="rounded-full p-1.5 text-slatey transition hover:bg-gray-50 hover:text-ink"
          >
            <X size={18} />
          </button>
        </div>

        {/* Canvas */}
        <div className="flex justify-center">
          <div
            className="relative touch-none select-none overflow-hidden rounded-lg bg-gray-100"
            style={{ width: frame.w, height: frame.h }}
            onPointerDown={onDown}
            onPointerMove={onMove}
            onPointerUp={onUp}
            onPointerLeave={onUp}
          >
            {url && (
              // eslint-disable-next-line @next/next/no-img-element
              <img
                ref={imgRef}
                src={url}
                alt=""
                draggable={false}
                onLoad={(e) =>
                  setNat({
                    w: e.currentTarget.naturalWidth,
                    h: e.currentTarget.naturalHeight,
                  })
                }
                style={{
                  position: "absolute",
                  width: dispW,
                  height: dispH,
                  left,
                  top,
                  maxWidth: "none",
                  cursor: "grab",
                  filter: filter.css,
                  transform: `rotate(${rotation}deg)`,
                  transformOrigin: "center",
                }}
              />
            )}
            {/* Rule-of-thirds guides */}
            <div className="pointer-events-none absolute inset-0">
              <div className="absolute left-1/3 top-0 h-full w-px bg-white/40" />
              <div className="absolute left-2/3 top-0 h-full w-px bg-white/40" />
              <div className="absolute left-0 top-1/3 h-px w-full bg-white/40" />
              <div className="absolute left-0 top-2/3 h-px w-full bg-white/40" />
              <div className="absolute inset-0 ring-2 ring-inset ring-white/70" />
            </div>
          </div>
        </div>

        {/* Tabs */}
        <div className="mt-4 flex gap-2 rounded-full bg-gray-50 p-1">
          <TabButton active={tab === "crop"} onClick={() => setTab("crop")}>
            <Crop size={14} /> Crop
          </TabButton>
          <TabButton active={tab === "filter"} onClick={() => setTab("filter")}>
            <Sparkles size={14} /> Filter
          </TabButton>
        </div>

        {tab === "crop" ? (
          <div className="mt-4 space-y-4">
            <div className="flex flex-wrap gap-2">
              {RATIOS.map((r) => (
                <Chip
                  key={r.label}
                  active={r.label === ratio.label}
                  onClick={() => setRatio(r)}
                >
                  {r.label}
                </Chip>
              ))}
            </div>

            <div className="flex items-center gap-3">
              <SlidersHorizontal size={14} className="shrink-0 text-slatey" />
              <input
                type="range"
                min={1}
                max={3}
                step={0.01}
                value={zoom}
                onChange={(e) => setZoom(parseFloat(e.target.value))}
                aria-label="Zoom"
                className="flex-1 accent-purple"
              />
              <span className="w-10 text-right text-label-md font-bold text-slatey">
                {zoom.toFixed(1)}x
              </span>
            </div>

            <div className="flex gap-2">
              <IconAction
                onClick={() => setRotation((r) => (r + 270) % 360)}
                label="Rotate left"
              >
                <RotateCcw size={16} />
              </IconAction>
              <IconAction
                onClick={() => setRotation((r) => (r + 90) % 360)}
                label="Rotate right"
              >
                <RotateCw size={16} />
              </IconAction>
              <button
                onClick={reset}
                className="ml-auto text-label-lg font-bold text-slatey underline decoration-dotted underline-offset-4 hover:text-ink"
              >
                Reset
              </button>
            </div>
          </div>
        ) : (
          <div className="mt-4 grid grid-cols-3 gap-2">
            {FILTERS.map((f) => (
              <button
                key={f.label}
                onClick={() => setFilter(f)}
                className={`overflow-hidden rounded-lg border-2 transition ${
                  f.label === filter.label
                    ? "border-purple"
                    : "border-transparent hover:border-gray-100"
                }`}
              >
                <div className="relative h-14 w-full bg-gray-100">
                  {url && (
                    // eslint-disable-next-line @next/next/no-img-element
                    <img
                      src={url}
                      alt=""
                      className="h-full w-full object-cover"
                      style={{ filter: f.css }}
                    />
                  )}
                </div>
                <div className="py-1 text-label-md font-bold text-slatey">
                  {f.label}
                </div>
              </button>
            ))}
          </div>
        )}

        <div className="mt-5 flex gap-3">
          <Button variant="outline" className="flex-1" onClick={onCancel}>
            Cancel
          </Button>
          <Button className="flex-1" onClick={save} loading={busy}>
            Use photo
          </Button>
        </div>
      </motion.div>
    </div>
  );
}

function TabButton({
  active,
  onClick,
  children,
}: {
  active: boolean;
  onClick: () => void;
  children: React.ReactNode;
}) {
  return (
    <button
      onClick={onClick}
      className={`flex flex-1 items-center justify-center gap-1.5 rounded-full py-2 text-label-lg font-extrabold transition ${
        active ? "bg-white text-purple shadow-card" : "text-slatey"
      }`}
    >
      {children}
    </button>
  );
}

function Chip({
  active,
  onClick,
  children,
}: {
  active: boolean;
  onClick: () => void;
  children: React.ReactNode;
}) {
  return (
    <button
      onClick={onClick}
      className={`rounded-full px-3 py-1.5 text-label-lg font-bold transition ${
        active ? "bg-purple text-white" : "bg-gray-50 text-slatey hover:bg-gray-100"
      }`}
    >
      {children}
    </button>
  );
}

function IconAction({
  onClick,
  label,
  children,
}: {
  onClick: () => void;
  label: string;
  children: React.ReactNode;
}) {
  return (
    <button
      onClick={onClick}
      aria-label={label}
      title={label}
      className="flex h-9 w-9 items-center justify-center rounded-lg bg-gray-50 text-slatey transition hover:bg-gray-100 hover:text-ink"
    >
      {children}
    </button>
  );
}
