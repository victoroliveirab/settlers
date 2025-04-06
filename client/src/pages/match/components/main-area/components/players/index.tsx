import { useMatchStore } from "@/state/match";

import { Player } from "./components/player";

export const Players = () => {
  const currentRoundPlayer = useMatchStore((state) => state.currentRoundPlayer);
  const players = useMatchStore((state) => state.players);
  const resourceCount = useMatchStore((state) => state.resourceCount);

  return (
    <ul className="h-full flex items-center gap-4">
      {players.map((player) => (
        <Player
          data={player}
          isPlayerRound={currentRoundPlayer === player.name}
          knightsUsed={0}
          longestRoad={0}
          numberOfCardsToDiscard={0}
          numberOfDevCards={0}
          numberOfResources={resourceCount[player.name]}
          points={0}
        />
      ))}
    </ul>
  );
};
