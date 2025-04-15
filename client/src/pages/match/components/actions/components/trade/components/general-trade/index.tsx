import { Button } from "@/components/ui/button";
import { DialogClose } from "@/components/ui/dialog";

import { Trade } from "@/features/trade";

import { useWebSocket } from "@/hooks/useWebSocket";
import { useMatchStore } from "@/state/match";

export const GeneralTrade = ({ isBankTrade, step }: { isBankTrade: boolean; step: number }) => {
  const { sendMessage } = useWebSocket();
  const hand = useMatchStore((state) => state.hand);

  const submitTrade = (
    given: SettlersCore.ResourceCollection,
    requested: SettlersCore.ResourceCollection,
  ) => {
    if (isBankTrade) {
      sendMessage({ type: "match.make-bank-trade", payload: { given, requested } });
    } else {
      sendMessage({ type: "match.make-general-port-trade", payload: { given, requested } });
    }
  };

  return (
    <Trade className="flex flex-col gap-4 py-2" givenResourcesAvailable={hand} givenStep={step}>
      {({ dirty, given, requested, totalGiven, totalRequested }) => {
        return (
          <div className="flex items-center w-full">
            <ul className="flex items-center justify-end gap-1">
              <li>
                <DialogClose asChild>
                  <Button>Close</Button>
                </DialogClose>
              </li>
              <li>
                <DialogClose asChild>
                  <Button
                    disabled={!dirty || totalGiven !== step || totalRequested !== 1}
                    onClick={() => submitTrade(given, requested)}
                  >
                    Submit
                  </Button>
                </DialogClose>
              </li>
            </ul>
          </div>
        );
      }}
    </Trade>
  );
};
