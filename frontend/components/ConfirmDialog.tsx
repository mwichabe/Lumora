"use client";

import { Button } from "./Button";

/** A small centered confirmation modal. */
export function ConfirmDialog({
  open,
  title,
  message,
  confirmLabel = "Confirm",
  cancelLabel = "Cancel",
  danger = false,
  onConfirm,
  onCancel,
}: {
  open: boolean;
  title: string;
  message?: string;
  confirmLabel?: string;
  cancelLabel?: string;
  danger?: boolean;
  onConfirm: () => void;
  onCancel: () => void;
}) {
  if (!open) return null;
  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
      onClick={onCancel}
    >
      <div
        className="w-full max-w-xs rounded-2xl bg-white p-5 text-center shadow-card-lg"
        onClick={(e) => e.stopPropagation()}
      >
        <h3 className="text-heading-sm font-extrabold text-ink">{title}</h3>
        {message && <p className="mt-1.5 text-body-md text-slatey">{message}</p>}
        <div className="mt-5 flex gap-3">
          <Button variant="outline" className="flex-1" onClick={onCancel}>
            {cancelLabel}
          </Button>
          <Button
            variant={danger ? "danger" : "primary"}
            className="flex-1"
            onClick={onConfirm}
          >
            {confirmLabel}
          </Button>
        </div>
      </div>
    </div>
  );
}
