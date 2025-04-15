import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { emojis } from "@/core/constants";
import { CirclePlay } from "lucide-react";

interface IPlayerProps {
  data: SettlersCore.Player;
  isPlayerRound: boolean;
  knightsUsed: number;
  longestRoad: number;
  numberOfCardsToDiscard: number;
  numberOfDevCards: number;
  numberOfResources: number;
  points: number;
}

export const Player = ({
  data,
  isPlayerRound,
  knightsUsed,
  longestRoad,
  numberOfCardsToDiscard,
  numberOfDevCards,
  numberOfResources,
  points,
}: IPlayerProps) => {
  const playerName = data.name;
  const playerColor = data.color;

  return (
    <li className="h-full">
      <Card
        className="relative h-full aspect-[3/2] flex flex-col justify-between border-2 border-dotted"
        variant="dense"
        style={{
          background: playerColor.background,
          color: playerColor.foreground,
          borderColor: isPlayerRound ? playerColor.foreground : "transparent",
        }}
      >
        <CardHeader>
          <CardTitle>{playerName}</CardTitle>
        </CardHeader>
        <CardContent>
          <ul className="grid grid-cols-2 grid-flow-row gap-1 text-xs">
            <li>#R: {numberOfResources}</li>
            <li>#D: {numberOfDevCards}</li>
            <li>LG: {longestRoad}</li>
            <li>#K: {knightsUsed}</li>
            <li>#P: {points}</li>
            {numberOfCardsToDiscard > 0 && <li>{emojis.misc.discarding}</li>}
          </ul>
        </CardContent>
        {isPlayerRound && (
          <div className="absolute top-0 right-0">
            <CirclePlay />
          </div>
        )}
      </Card>
    </li>
  );
};
