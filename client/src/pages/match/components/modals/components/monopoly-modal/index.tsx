import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

import { Monopoly } from "@/features/monopoly";

import { useWebSocket } from "@/hooks/useWebSocket";
import { useMatchStore } from "@/state/match";

export const PickRobbedModal = () => {
  const { sendMessage } = useWebSocket();
  const robbablePlayersState = useMatchStore((state) => state.robbablePlayers);
  const players = useMatchStore((state) => state.players);

  const onResourceClick = (resource: SettlersCore.Resource) => {
    sendMessage({ type: "match.monopoly", payload: { resource } });
  };

  return (
    <Dialog open={robbablePlayersState.enabled}>
      <DialogContent className="w-fit min-w-60">
        <DialogHeader>
          <DialogTitle>Pick Resource</DialogTitle>
          <DialogDescription>
            <Monopoly onClick={onResourceClick} />
          </DialogDescription>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
};
