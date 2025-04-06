import { GameCard } from "@/components/custom/game-card";
import { resourcesOrder } from "@/core/constants";

import { useMatchStore } from "@/state/match";

export const Hand = () => {
  const hand = useMatchStore((state) => state.hand);
  return (
    <ul className="flex items-center h-full gap-1">
      {resourcesOrder.map((resource) => {
        if (hand[resource] === 0) return null;
        return Array.from({ length: hand[resource] }).map((_, index) => (
          <GameCard key={`${resource}-${index}`} as="li" value={resource} />
        ));
      })}
    </ul>
  );
};
