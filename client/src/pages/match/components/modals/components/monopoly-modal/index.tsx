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

export const MonopolyModal = () => {
  const { sendMessage } = useWebSocket();
  const monopolyState = useMatchStore((state) => state.monopoly);

  const onResourceClick = (resource: SettlersCore.Resource) => {
    sendMessage({ type: "match.monopoly", payload: { resource } });
  };

  return (
    <Dialog open={monopolyState.enabled}>
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
