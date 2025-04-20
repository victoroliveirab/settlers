import { useCallback, useEffect, useRef, useState } from "react";

const parseTime = (isoDate: string) => new Date(isoDate).getTime();

export function useCountdownTimer(deadline: string | null, serverNow: string | null) {
  const [timeLeft, setTimeLeft] = useState(0);
  const deadlineRef = useRef<number | null>(null);
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const startTicking = useCallback(() => {
    const tick = () => {
      if (!deadlineRef.current) return;

      const now = Date.now();
      const remaining = Math.max(0, Math.floor((deadlineRef.current - now) / 1000));
      setTimeLeft(remaining);

      if (remaining > 0) {
        const nextTickIn = 1000 - (now % 1000);
        timeoutRef.current = setTimeout(tick, nextTickIn);
      }
    };

    tick();
  }, []);

  useEffect(() => {
    if (!deadline || !serverNow) return;

    const serverNowMs = parseTime(serverNow);
    const deadlineMs = parseTime(deadline);
    const now = Date.now();
    const offset = now - serverNowMs;

    const adjustedDeadline = deadlineMs + offset;
    deadlineRef.current = adjustedDeadline;
    setTimeLeft(Math.max(0, Math.floor((adjustedDeadline - now) / 1000)));

    startTicking();

    return () => {
      if (timeoutRef.current) clearTimeout(timeoutRef.current);
    };
  }, [deadline, serverNow, startTicking]);

  return timeLeft;
}
