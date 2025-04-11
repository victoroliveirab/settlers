import { Button } from "@/components/ui/button";
import { DialogClose } from "@/components/ui/dialog";

import { Trade } from "@/features/trade";
import { hasSameResourceInOfferAndRequest } from "@/features/trade/utils";

import { useWebSocket } from "@/hooks/useWebSocket";
import { useMatchStore } from "@/state/match";

// TODO: get from game params
const STEP = 4;

export const BankTrade = () => {
  const { sendMessage } = useWebSocket();
  const hand = useMatchStore((state) => state.hand);

  const submitTrade = (given: SettlersCore.Resource, requested: SettlersCore.Resource) => {
    sendMessage({ type: "match.make-bank-trade", payload: { given, requested } });
  };

  return (
    <Trade className="py-2" givenResourcesAvailable={hand} givenStep={STEP}>
      {({ dirty, given, requested, totalGiven, totalRequested }) => {
        const givenResource = Object.keys(given).find(
          (resource) => given[resource as SettlersCore.Resource] > 0,
        ) as SettlersCore.Resource;
        const requestedResource = Object.keys(requested).find(
          (resource) => requested[resource as SettlersCore.Resource] > 0,
        ) as SettlersCore.Resource;

        return (
          <div className="flex items-center w-full">
            <p className="text-xs text-gray-500 flex-1">Submit each trade individually</p>
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
                      totalGiven !== STEP ||
                      totalRequested !== 1 ||
                      hasSameResourceInOfferAndRequest({
                        given,
                        requested,
                      })
                    }
                    onClick={() => submitTrade(givenResource, requestedResource)}
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
