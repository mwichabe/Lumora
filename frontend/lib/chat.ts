"use client";

import { useEffect, useState } from "react";
import { api } from "./api";

/** Polls the unread direct-message count for badges. */
export function useChatUnread(): number {
  const [count, setCount] = useState(0);
  useEffect(() => {
    const load = () =>
      api.chatUnread().then((d) => setCount(d.count)).catch(() => {});
    load();
    const id = setInterval(load, 15000);
    return () => clearInterval(id);
  }, []);
  return count;
}
