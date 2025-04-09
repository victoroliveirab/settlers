import { GameCard } from "@/components/custom/game-card";

import { cn } from "@/lib/utils";

import { developmentOrder } from "@/core/constants";
import { useWebSocket } from "@/hooks/useWebSocket";
import { useMatchStore } from "@/state/match";

export const DevHand = () => {
  const { sendMessage } = useWebSocket();
  const devHand = useMatchStore((state) => state.devHand);
  const devHandPermissions = useMatchStore((state) => state.devHandPermissions);

  const onDevCardClick = (kind: SettlersCore.DevelopmentCard) => {
    sendMessage({ type: "match.dev-card-click", payload: { kind } });
  };

  return (
    <ul className="flex items-center h-full gap-1">
      {developmentOrder.map((card) => {
        if (devHand[card] === 0) return null;
        return Array.from({ length: devHand[card] }).map((_, index) => (
          <GameCard
            key={`${card}-${index}`}
            className={cn({
              "cursor-not-allowed": !devHandPermissions[card],
              "cursor-pointer": devHandPermissions[card],
            })}
            as="li"
            value={card}
            onClick={() => onDevCardClick(card)}
          />
        ));
      })}
    </ul>
  );
};
