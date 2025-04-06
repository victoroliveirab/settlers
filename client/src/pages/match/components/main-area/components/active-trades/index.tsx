import { Check, XIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";

import { useWebSocket } from "@/hooks/useWebSocket";
import { useMatchStore } from "@/state/match";
import { usePlayerStore } from "@/state/player";

import { ResourceList } from "./components/resource-list";
import { OpponentsStatus } from "./components/opponents-status";
import { CounterOfferDialog } from "./components/counter-offer-dialog";

const TradeCard = ({
  creator,
  isCounterOffer,
  offeredResources,
  opponents,
  requestedResources,
  tradeID,
}: {
  creator: string;
  isCounterOffer: boolean;
  offeredResources: SettlersCore.ResourceCollection;
  opponents: Record<
    SettlersCore.Player["name"],
    { status: "Open" | "Accepted" | "Declined"; blocked: boolean }
  >;
  requestedResources: SettlersCore.ResourceCollection;
  tradeID: number;
}) => {
  const { sendMessage } = useWebSocket();
  const username = usePlayerStore((state) => state.username);
  console.log({ creator, isCounterOffer, tradeID });

  const submitOfferAccept = () => {
    console.log("submitting offer accept");
    sendMessage({ type: "match.accept-trade-offer", payload: { tradeID } });
  };

  const submitOfferReject = () => {
    console.log("submitting offer reject");
    sendMessage({ type: "match.reject-trade-offer", payload: { tradeID } });
  };

  const submitOfferCancel = () => {
    console.log("submitting offer cancel");
    sendMessage({ type: "match.cancel-trade-offer", payload: { tradeID } });
  };

  return (
    <>
      <Card className="w-full">
        <CardContent className="w-full">
          <div className="grid grid-cols-2 gap-2 grid-flow-row items-center w-full">
            <h3 className="text-xs">Offered</h3>
            <ResourceList resources={offeredResources} />
            <h3 className="text-xs">Requested</h3>
            <ResourceList resources={requestedResources} />
            <OpponentsStatus data={opponents} disabled={creator !== username} tradeID={tradeID} />
            <div className="flex gap-1 justify-end col-span-2">
              {creator !== username && !isCounterOffer && (
                <CounterOfferDialog
                  offeredResources={offeredResources}
                  requestedResources={requestedResources}
                  tradeID={tradeID}
                />
              )}
              <Button
                size="xs"
                variant="destructive"
                onClick={creator === username ? submitOfferCancel : submitOfferReject}
              >
                <XIcon />
              </Button>
              {creator !== username && (
                <Button size="xs" variant="success" onClick={submitOfferAccept}>
                  <Check />
                </Button>
              )}
            </div>
          </div>
        </CardContent>
      </Card>
    </>
  );
};

export const ActiveTrades = () => {
  const activeTrades = useMatchStore((state) => state.activeTradeOffers);

  if (activeTrades.length === 0) return null;

  return (
    <>
      {activeTrades.map((tradeOffer) => (
        <TradeCard
          creator={tradeOffer.player}
          isCounterOffer={tradeOffer.parent >= 0}
          offeredResources={tradeOffer.offer}
          opponents={tradeOffer.opponents}
          requestedResources={tradeOffer.request}
          tradeID={tradeOffer.id}
        />
      ))}
    </>
  );
};
