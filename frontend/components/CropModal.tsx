"use client";

import { useEffect, useRef, useState } from "react";
import { Button } from "./Button";

/**
 * Lightweight circular image cropper — zoom + drag inside a round frame, then
 * export a square PNG via canvas. No external dependency.
 */
const VIEW = 280; // on-screen crop frame (px)
const OUT = 512; // exported image size (px)

export function CropModal({
  file,
  onCancel,
  onCropped,
}: {
  file: File;
  onCancel: () => void;
  onCropped: (f: File) => void;
}) {
  const imgRef = useRef<HTMLImageElement | null>(null);
  const dragRef = useRef<{ x: number; y: number; ox: number; oy: number } | null>(null);
  const [url, setUrl] = useState("");
  const [nat, setNat] = useState({ w: 0, h: 0 });
  const [zoom, setZoom] = useState(1);
  const [off, setOff] = useState({ x: 0, y: 0 });

  useEffect(() => {
    const u = URL.createObjectURL(file);
    setUrl(u);
    return () => URL.revokeObjectURL(u);
  }, [file]);

  const minDim = Math.min(nat.w || 1, nat.h || 1);
  const scale = (VIEW / minDim) * zoom;
  const dispW = nat.w * scale;
  const dispH = nat.h * scale;
  const left = VIEW / 2 - dispW / 2 + off.x;
  const top = VIEW / 2 - dispH / 2 + off.y;

  // Keep the frame fully covered by the image (no empty edges).
  function clamp(o: { x: number; y: number }) {
    const l = VIEW / 2 - dispW / 2 + o.x;
    const t = VIEW / 2 - dispH / 2 + o.y;
    let nx = o.x;
    let ny = o.y;
    if (l > 0) nx -= l;
    if (l < VIEW - dispW) nx += VIEW - dispW - l;
    if (t > 0) ny -= t;
    if (t < VIEW - dispH) ny += VIEW - dispH - t;
    return { x: nx, y: ny };
  }

  useEffect(() => {
    setOff((o) => clamp(o));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [zoom, nat.w, nat.h]);

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

  function save() {
    const img = imgRef.current;
    if (!img || !nat.w) return;
    const sx = (0 - left) / scale;
    const sy = (0 - top) / scale;
    const s = VIEW / scale;
    const canvas = document.createElement("canvas");
    canvas.width = OUT;
    canvas.height = OUT;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;
    ctx.drawImage(img, sx, sy, s, s, 0, 0, OUT, OUT);
    canvas.toBlob(
      (blob) => {
        if (blob) onCropped(new File([blob], "avatar.png", { type: "image/png" }));
      },
      "image/png",
      0.92
    );
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
      <div className="w-full max-w-sm rounded-2xl bg-white p-5 shadow-card-lg">
        <h3 className="mb-3 text-heading-sm font-extrabold text-ink">
          Crop your photo
        </h3>

        <div
          className="relative mx-auto touch-none select-none overflow-hidden rounded-lg bg-gray-100"
          style={{ width: VIEW, height: VIEW }}
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
              }}
            />
          )}
          {/* circular mask */}
          <div
            className="pointer-events-none absolute inset-0 rounded-full"
            style={{ boxShadow: "0 0 0 9999px rgba(0,0,0,0.45)" }}
          />
        </div>

        <div className="mt-4 flex items-center gap-3">
          <span className="text-label-md font-bold text-slatey">Zoom</span>
          <input
            type="range"
            min={1}
            max={3}
            step={0.01}
            value={zoom}
            onChange={(e) => setZoom(parseFloat(e.target.value))}
            className="flex-1 accent-purple"
          />
        </div>

        <div className="mt-4 flex gap-3">
          <Button variant="outline" className="flex-1" onClick={onCancel}>
            Cancel
          </Button>
          <Button className="flex-1" onClick={save}>
            Save photo
          </Button>
        </div>
      </div>
    </div>
  );
}
