import { useMatchStore } from "@/state/match";

import { Player } from "./components/player";

export const Players = () => {
  const currentRoundPlayer = useMatchStore((state) => state.currentRoundPlayer);
  const discardAmountByPlayer = useMatchStore((state) => state.discard.discardAmounts);
  const knightsUsedByPlayer = useMatchStore((state) => state.knightUsages);
  const longestRoadByPlayer = useMatchStore((state) => state.longestRoadSize);
  const players = useMatchStore((state) => state.players);
  const pointsByPlayer = useMatchStore((state) => state.points);
  const resourceCountByPlayer = useMatchStore((state) => state.resourceCount);

  return (
    <ul className="h-full flex items-center gap-4">
      {players.map((player) => (
        <Player
          data={player}
          isPlayerRound={currentRoundPlayer === player.name}
          knightsUsed={knightsUsedByPlayer[player.name]}
          longestRoad={longestRoadByPlayer[player.name]}
          numberOfCardsToDiscard={discardAmountByPlayer[player.name]}
          numberOfDevCards={0}
          numberOfResources={resourceCountByPlayer[player.name]}
          points={pointsByPlayer[player.name]}
        />
      ))}
    </ul>
  );
};
