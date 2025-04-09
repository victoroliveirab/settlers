import { Button } from "@/components/ui/button";
import { DialogClose } from "@/components/ui/dialog";

import { Trade } from "@/features/trade";
import { hasSameResourceInOfferAndRequest } from "@/features/trade/utils";

import { useWebSocket } from "@/hooks/useWebSocket";
import { useMatchStore } from "@/state/match";

export const PlayerTrade = () => {
  const hand = useMatchStore((state) => state.hand);
  const { sendMessage } = useWebSocket();

  const submitTrade = (
    given: SettlersCore.ResourceCollection,
    requested: SettlersCore.ResourceCollection,
  ) => {
    sendMessage({ type: "match.create-trade-offer", payload: { given, requested } });
  };

  return (
    <Trade className="py-2" givenResourcesAvailable={hand}>
      {({ dirty, given, requested, totalGiven, totalRequested }) => {
        return (
          <ul className="flex items-center justify-end gap-1">
            <li>
              <DialogClose asChild>
                <Button>Close</Button>
              </DialogClose>
            </li>
            <li>
              <DialogClose asChild>
                <Button
                  disabled={
                    !dirty ||
                    totalGiven === 0 ||
                    totalRequested === 0 ||
                    hasSameResourceInOfferAndRequest({
                      given,
                      requested,
                    })
                  }
                  onClick={() => submitTrade(given, requested)}
                >
                  Submit
                </Button>
              </DialogClose>
            </li>
          </ul>
        );
      }}
    </Trade>
  );
};
