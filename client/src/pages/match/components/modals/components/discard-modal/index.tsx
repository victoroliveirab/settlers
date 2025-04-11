import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { cn } from "@/lib/utils";

import { Discard } from "@/features/discard";

import { useWebSocket } from "@/hooks/useWebSocket";
import { useMatchStore } from "@/state/match";
import { usePlayerStore } from "@/state/player";

export const DiscardModal = () => {
  const { sendMessage } = useWebSocket();
  const hand = useMatchStore((state) => state.hand);
  const discardState = useMatchStore((state) => state.discard);
  const username = usePlayerStore((state) => state.username);

  const targetSelected = discardState.discardAmounts[username ?? ""];

  const submitDiscard = (resources: SettlersCore.ResourceCollection) => {
    sendMessage({ type: "match.discard-cards", payload: { resources } });
  };

  return (
    <Dialog open={discardState.enabled}>
      <DialogContent className="w-fit" hideCloseButton>
        <DialogHeader>
          <DialogTitle>Discard {targetSelected} cards</DialogTitle>
          <DialogDescription className="flex flex-col gap-2">
            <Discard resources={hand}>
              {({ selected, totalSelected }) => {
                return (
                  <div className="flex items-center justify-between">
                    <p
                      className={cn("text-xs text-gray-500 select-none", {
                        "text-red-600": totalSelected > targetSelected,
                        "text-green-600": totalSelected === targetSelected,
                      })}
                    >{`${totalSelected}/${targetSelected}`}</p>
                    <Button
                      disabled={totalSelected !== targetSelected}
                      onClick={() => submitDiscard(selected)}
                    >
                      Discard
                    </Button>
                  </div>
                );
              }}
            </Discard>
          </DialogDescription>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
};
