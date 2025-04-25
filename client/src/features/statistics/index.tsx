import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

import { GeneralDiceStatistics } from "./components/general-dice";
import { DiceByPlayer } from "./components/dice-by-player";
import { Points } from "./components/points";

export interface IStatisticsProps {
  diceStats: Record<number, number>;
  diceStatsByPlayer: Record<string, Record<number, number>>;
  players: SettlersCore.Player[];
  pointsDistribution: Record<SettlersCore.Player["name"], Record<string, number>> | null;
}

export const Statistics = ({
  diceStats,
  diceStatsByPlayer,
  players,
  pointsDistribution,
}: IStatisticsProps) => {
  const hasPointsDistribution = pointsDistribution !== null;
  const defaultView = hasPointsDistribution ? "points" : "general-dice";
  return (
    <Tabs className="w-full h-full" defaultValue={defaultView}>
      <TabsList className="w-full">
        {hasPointsDistribution && (
          <TabsTrigger className="cursor-pointer" value="points">
            Points
          </TabsTrigger>
        )}
        <TabsTrigger className="cursor-pointer" value="general-dice">
          Dice
        </TabsTrigger>
        <TabsTrigger className="cursor-pointer" value="dice-player">
          Dice (player)
        </TabsTrigger>
      </TabsList>
      <TabsContent value="points">
        <Points data={pointsDistribution!} players={players} />
      </TabsContent>
      <TabsContent value="general-dice">
        <GeneralDiceStatistics data={diceStats} />
      </TabsContent>
      <TabsContent value="dice-player">
        <DiceByPlayer data={diceStatsByPlayer} players={players} />
      </TabsContent>
    </Tabs>
  );
};
