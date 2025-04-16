import { useEffect, useState } from "react";

import { cn } from "@/lib/utils";

import { useMatchStore } from "@/state/match";

// Parses to ms
const parseTime = (isoDate: string) => new Date(isoDate).getTime();

export const Timer = () => {
  const currentRoundPlayer = useMatchStore((state) => state.currentRoundPlayer);
  const [timeLeft, setTimeLeft] = useState(0);

  useEffect(() => {
    if (!currentRoundPlayer) return;
    const { deadline, serverNow, subDeadline } = currentRoundPlayer;
    const diff = parseTime(subDeadline || deadline) - parseTime(serverNow);
    setTimeLeft(Math.max(0, Math.floor(diff / 1000)));
  }, [currentRoundPlayer]);

  useEffect(() => {
    if (timeLeft <= 0) return;

    const interval = setInterval(() => {
      setTimeLeft((prev) => {
        const next = prev - 1;
        return next > 0 ? next : 0;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [timeLeft]);

  if (!currentRoundPlayer) return null;

  const minutes = String(Math.floor(timeLeft / 60)).padStart(2, "0");
  const seconds = String(timeLeft % 60).padStart(2, "0");

  return (
    <div className="rounded-lg p-1 bg-neutral-100 select-none w-20 text-center">
      <span
        className={cn("text-xl text-neutral-800", {
          "text-red-800 animate-blink": timeLeft <= 5,
        })}
      >
        {`${minutes}:${seconds}`}
      </span>
    </div>
  );
};
