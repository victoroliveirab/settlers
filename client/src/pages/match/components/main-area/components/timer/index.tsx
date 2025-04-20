import { cn } from "@/lib/utils";

import { useCountdownTimer } from "@/hooks/useCountdownTimer";
import { useMatchStore } from "@/state/match";

export const Timer = () => {
  const currentRoundPlayer = useMatchStore((state) => state.currentRoundPlayer);
  const deadline = currentRoundPlayer?.subDeadline || currentRoundPlayer?.deadline || null;
  const serverNow = currentRoundPlayer?.serverNow || null;
  const timeLeft = useCountdownTimer(deadline, serverNow);

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
