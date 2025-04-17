import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";

import { usePlayerStore } from "@/state/player";
import { useMatchStore } from "@/state/match";

import { Player } from "./components/player";

export const Players = () => {
  const currentRoundPlayer = useMatchStore((state) => state.currentRoundPlayer);
  const discardAmountByPlayer = useMatchStore((state) => state.discard.discardAmounts);
  const knightsUsedByPlayer = useMatchStore((state) => state.knightUsages);
  const longestRoadByPlayer = useMatchStore((state) => state.longestRoadSize);
  const players = useMatchStore((state) => state.players);
  const pointsByPlayer = useMatchStore((state) => state.points);
  const devHand = useMatchStore((state) => state.devHand);
  const resourceCountByPlayer = useMatchStore((state) => state.resourceCount);
  const username = usePlayerStore((state) => state.username);

  return (
    <ScrollArea>
      <ul className="h-30 flex items-center gap-4">
        {players.map((player) => (
          <Player
            data={player}
            isPlayerRound={currentRoundPlayer?.player === player.name}
            knightsUsed={knightsUsedByPlayer[player.name]}
            longestRoad={longestRoadByPlayer[player.name]}
            numberOfCardsToDiscard={discardAmountByPlayer[player.name]}
            numberOfDevCards={0}
            numberOfResources={resourceCountByPlayer[player.name]}
            points={pointsByPlayer[player.name]}
            extraPoints={player.name === username ? devHand["Victory Point"] : 0}
          />
        ))}
      </ul>
      <ScrollBar orientation="horizontal" />
    </ScrollArea>
  );
};
