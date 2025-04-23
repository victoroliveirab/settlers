import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

import { GeneralDiceStatistics } from "./components/general-dice";
import { DiceByPlayer } from "./components/dice-by-player";

export interface IStatisticsProps {
  diceStats: Record<number, number>;
  diceStatsByPlayer: Record<string, Record<number, number>>;
  players: SettlersCore.Player[];
}

export const Statistics = ({ diceStats, diceStatsByPlayer, players }: IStatisticsProps) => {
  return (
    <Tabs defaultValue="general-dice">
      <TabsList className="w-full">
        <TabsTrigger className="cursor-pointer" value="general-dice">
          Dice
        </TabsTrigger>
        <TabsTrigger className="cursor-pointer" value="dice-player">
          Dice (player)
        </TabsTrigger>
      </TabsList>
      <TabsContent value="general-dice">
        <GeneralDiceStatistics data={diceStats} />
      </TabsContent>
      <TabsContent value="dice-player">
        <DiceByPlayer data={diceStatsByPlayer} players={players} />
      </TabsContent>
    </Tabs>
  );
};
