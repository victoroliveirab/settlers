import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

import { Statistics } from "@/features/statistics";

import { useMatchReportStore } from "@/state/match-report";
import { useMatchStore } from "@/state/match";

export const PostMatch = () => {
  const state = useMatchReportStore((state) => state);
  const players = useMatchStore((state) => state.players);

  return (
    <main className="h-full flex items-center justify-center max-w-4xl mx-auto">
      <Card className="w-full overflow-hidden">
        <CardHeader>
          <CardTitle>Room #{state.roomName}</CardTitle>
          <CardDescription>4-player game, base map</CardDescription>
        </CardHeader>
        <CardContent className="flex gap-4 justify-center max-h-[50vh]">
          <Statistics
            diceStats={state.statistics.generalDiceStats}
            diceStatsByPlayer={state.statistics.diceStatsByPlayer}
            players={players}
            pointsDistribution={state.pointsDistribution}
          />
        </CardContent>
        <CardFooter className="mt-8 flex items-center justify-end gap-2">
          <Button>Share</Button>
          <Button>Back to Lobby</Button>
        </CardFooter>
      </Card>
    </main>
  );
};
