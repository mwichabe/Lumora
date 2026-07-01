"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { CheckCircle2, XCircle } from "lucide-react";
import { FoxMascot } from "@/components/FoxMascot";
import { Button } from "@/components/Button";
import { api } from "@/lib/api";
import { useAuth } from "@/lib/auth";

type State = "checking" | "success" | "failed";

export default function PaymentCallbackPage() {
  const router = useRouter();
  const { refresh } = useAuth();
  const [state, setState] = useState<State>("checking");
  const [product, setProduct] = useState("");

  useEffect(() => {
    // Paystack appends ?reference=... (and ?trxref=...) to the callback URL.
    const params = new URLSearchParams(window.location.search);
    const reference = params.get("reference") || params.get("trxref");
    if (!reference) {
      setState("failed");
      return;
    }
    api
      .verifyPayment(reference)
      .then((r) => {
        setProduct(r.product);
        setState(r.success ? "success" : "failed");
        refresh(); // pull any updated state into the app
        // Surface the notification + refresh hearts/badges right away.
        window.dispatchEvent(new Event("lumora:notifications"));
        window.dispatchEvent(new Event("lumora:hearts"));
      })
      .catch(() => setState("failed"));
  }, [refresh]);

  const isHearts = product === "hearts_refill";

  return (
    <div className="flex min-h-[100dvh] w-full items-center justify-center bg-cream px-6">
      <div className="w-full max-w-sm text-center">
        {state === "checking" && (
          <>
            <FoxMascot size={120} glow />
            <p className="mt-5 text-heading-sm font-extrabold text-ink">
              Confirming your payment…
            </p>
            <p className="mt-1 text-body-md text-slatey">
              Hang tight, this only takes a moment.
            </p>
          </>
        )}

        {state === "success" && (
          <>
            <span className="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-teal/10 text-teal">
              <CheckCircle2 size={44} />
            </span>
            <h1 className="mt-4 text-heading-xl font-extrabold text-ink">
              {isHearts ? "Hearts refilled! ❤️" : "You're all set! 🎉"}
            </h1>
            <p className="mt-1 text-body-md text-slatey">
              {isHearts
                ? "Your hearts are full again — jump back into your lesson."
                : "Your exam attempt is ready. Good luck!"}
            </p>
            <div className="mt-6 space-y-2">
              {isHearts ? (
                <Button full onClick={() => router.push("/learn")}>
                  Continue learning
                </Button>
              ) : (
                <Button full onClick={() => router.push("/exam")}>
                  Start the exam
                </Button>
              )}
              <Button full variant="outline" onClick={() => router.push("/home")}>
                Back to home
              </Button>
            </div>
          </>
        )}

        {state === "failed" && (
          <>
            <span className="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-coral/10 text-coral">
              <XCircle size={44} />
            </span>
            <h1 className="mt-4 text-heading-xl font-extrabold text-ink">
              Payment not completed
            </h1>
            <p className="mt-1 text-body-md text-slatey">
              We couldn&apos;t confirm your payment. If you were charged, it will
              be applied automatically — otherwise you can try again.
            </p>
            <div className="mt-6 space-y-2">
              <Button full onClick={() => router.push("/exam")}>
                Try again
              </Button>
              <Button full variant="outline" onClick={() => router.push("/profile")}>
                Back to profile
              </Button>
            </div>
          </>
        )}
      </div>
    </div>
  );
}
