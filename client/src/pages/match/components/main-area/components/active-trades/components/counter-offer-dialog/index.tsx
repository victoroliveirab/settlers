import { Pencil } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

import { Trade } from "@/features/trade";
import { hasSameResourceInOfferAndRequest } from "@/features/trade/utils";

import { useWebSocket } from "@/hooks/useWebSocket";
import { useMatchStore } from "@/state/match";

interface ICounterOfferDialogProps {
  offeredResources: SettlersCore.ResourceCollection;
  requestedResources: SettlersCore.ResourceCollection;
  tradeID: number;
}

export const CounterOfferDialog = ({
  offeredResources,
  requestedResources,
  tradeID,
}: ICounterOfferDialogProps) => {
  const { sendMessage } = useWebSocket();
  const hand = useMatchStore((state) => state.hand);

  const onSubmit = (
    given: SettlersCore.ResourceCollection,
    requested: SettlersCore.ResourceCollection,
  ) => {
    sendMessage({
      type: "match.create-counter-trade-offer",
      payload: { given, requested, tradeID },
    });
  };

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button size="xs" variant="outline">
          <Pencil />
        </Button>
      </DialogTrigger>
      <DialogContent className="w-3xl h-fit">
        <DialogHeader className="h-fit">
          <DialogTitle>Create counter-offer</DialogTitle>
          <DialogDescription>
            <Trade
              givenResourcesAvailable={hand}
              initialStateGiven={requestedResources}
              initialStateRequested={offeredResources}
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
                        onClick={() => onSubmit(given, requested)}
                      >
                        Submit
                      </Button>
                    </li>
                  </ul>
                );
              }}
            </Trade>
          </DialogDescription>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
};
