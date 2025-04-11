import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

import { PickRobbed } from "@/features/pick-robbed";

import { useWebSocket } from "@/hooks/useWebSocket";
import { useMatchStore } from "@/state/match";

export const PickRobbedModal = () => {
  const { sendMessage } = useWebSocket();
  const robbablePlayersState = useMatchStore((state) => state.robbablePlayers);
  const players = useMatchStore((state) => state.players);

  const robbablePlayers = players.filter((player) =>
    robbablePlayersState.options.includes(player.name),
  );

  const onPlayerChosen = (playerName: string) => {
    sendMessage({ type: "match.rob-player", payload: { player: playerName } });
  };

  return (
    <Dialog open={robbablePlayersState.enabled}>
      <DialogContent className="w-fit min-w-60">
        <DialogHeader>
          <DialogTitle>Pick Robbed Player</DialogTitle>
          <DialogDescription>
            <PickRobbed onClick={onPlayerChosen} players={robbablePlayers} />
          </DialogDescription>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
};
