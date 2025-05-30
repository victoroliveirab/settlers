import { useEffect, useState } from "react";
import { ChartColumn } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

import { Statistics as StatisticsComponents } from "@/features/statistics";

import { useWebSocket } from "@/hooks/useWebSocket";
import { useRoomStore } from "@/state/room";
import { useMatchStore } from "@/state/match";
import { useMatchReportStore } from "@/state/match-report";

export const Statistics = () => {
  const [open, setOpen] = useState(false);
  const roomStatus = useRoomStore((state) => state.room.status);
  const players = useMatchStore((state) => state.players);
  const stats = useMatchReportStore((state) => state.statistics);
  const { sendMessage } = useWebSocket();

  const requestStatistics = () => {
    sendMessage({ type: "match.report", payload: {} });
  };

  useEffect(() => {
    if (!stats.diceStatsByPlayer) {
      setOpen(false);
    }
  }, [stats.diceStatsByPlayer]);

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button disabled={roomStatus !== "match"} onClick={requestStatistics}>
          <ChartColumn />
        </Button>
      </DialogTrigger>
      <DialogContent className="w-2xl" __TEMP_FREE_SIZE>
        <DialogHeader>
          <DialogTitle>Statistics</DialogTitle>
          <DialogDescription>
            <StatisticsComponents
              diceStats={stats.generalDiceStats}
              diceStatsByPlayer={stats.diceStatsByPlayer}
              players={players}
              pointsDistribution={null}
            />
          </DialogDescription>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
};
