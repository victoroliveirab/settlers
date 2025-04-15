import { Button } from "@/components/ui/button";
import { DialogClose } from "@/components/ui/dialog";

import { Trade } from "@/features/trade";

import { useWebSocket } from "@/hooks/useWebSocket";
import { useMatchStore } from "@/state/match";
import { useCallback } from "react";

export const PortTrade = ({ step }: { step: number }) => {
  const hand = useMatchStore((state) => state.hand);
  const ports = useMatchStore((state) => state.ownedPorts);
  const { sendMessage } = useWebSocket();

  const submitTrade = (
    given: SettlersCore.ResourceCollection,
    requested: SettlersCore.ResourceCollection,
  ) => {
    sendMessage({ type: "match.make-resource-port-trade", payload: { given, requested } });
  };

  const isResourceAvailableToBeGiven = useCallback(
    (resource: SettlersCore.Resource) => {
      return !(ports.includes(resource) && hand[resource] >= step);
    },
    [hand, step],
  );

  return (
    <Trade
      className="flex flex-col gap-4 py-2"
      givenResourcesAvailable={hand}
      givenResourceIsDisabledGetter={isResourceAvailableToBeGiven}
      givenStep={step}
    >
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
                    Math.trunc(totalGiven / step) !== totalRequested
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
