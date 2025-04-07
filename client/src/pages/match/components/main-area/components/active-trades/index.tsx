import { useEffect, useState } from "react";
import { Check, XIcon } from "lucide-react";

import { Avatar } from "@/components/custom/avatar";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { ScrollArea } from "@/components/ui/scroll-area";

import { useWebSocket } from "@/hooks/useWebSocket";
import { useMatchStore } from "@/state/match";
import { usePlayerStore } from "@/state/player";

import { ResourceList } from "./components/resource-list";
import { OpponentsStatus } from "./components/opponents-status";
import { CounterOfferDialog } from "./components/counter-offer-dialog";
import { CollapsibleToggler } from "./components/collapsible-toggler";

const TradeCard = ({
  creator,
  offeredResources,
  opponents,
  requestedResources,
  requester,
  tradeID,
}: {
  creator: string;
  offeredResources: SettlersCore.ResourceCollection;
  opponents: Record<
    SettlersCore.Player["name"],
    { status: "Open" | "Accepted" | "Declined"; blocked: boolean }
  >;
  requestedResources: SettlersCore.ResourceCollection;
  requester: string;
  tradeID: number;
}) => {
  const { sendMessage } = useWebSocket();
  const username = usePlayerStore((state) => state.username);
  const players = useMatchStore((state) => state.players);

  const isTradeRequester = requester === username;

  const submitOfferAccept = () => {
    sendMessage({ type: "match.accept-trade-offer", payload: { tradeID } });
  };

  const submitOfferReject = () => {
    sendMessage({ type: "match.reject-trade-offer", payload: { tradeID } });
  };

  const submitOfferCancel = () => {
    sendMessage({ type: "match.cancel-trade-offer", payload: { tradeID } });
  };

  return (
    <>
      <Card className="w-full bg-gray-50">
        <CardContent className="w-full">
          <div className="grid grid-cols-2 gap-2 grid-flow-row items-center w-full">
            <h3 className="text-xs inline-flex items-start gap-1">
              <Avatar
                background={players.find((player) => player.name === requester)?.color.background}
                className="h-4 w-4"
              />
              <span>gives</span>
            </h3>
            <ResourceList resources={offeredResources} />
            <h3 className="text-xs inline-flex items-start gap-1">
              <Avatar
                background={
                  isTradeRequester
                    ? undefined
                    : players.find((player) => player.name === username)?.color.background
                }
                borderStyle={isTradeRequester ? "dashed" : undefined}
                className="h-4 w-4"
                foreground={isTradeRequester ? "black" : undefined}
                withBorder={isTradeRequester}
              />
              <span>gives</span>
            </h3>
            <ResourceList resources={requestedResources} />
            <OpponentsStatus data={opponents} disabled={!isTradeRequester} tradeID={tradeID} />
            <div className="flex gap-1 justify-end col-span-2">
              {creator !== username && (
                <CounterOfferDialog
                  isTradeRequester={isTradeRequester}
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
              {!isTradeRequester && (
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
  const [open, setOpen] = useState(false);
  const activeTrades = useMatchStore((state) => state.activeTradeOffers);

  useEffect(() => {
    // TODO: use last activeTrades.length with a ref to set this
    if (activeTrades.length === 0) {
      setOpen(false);
    } else {
      setOpen(true);
    }
  }, [activeTrades]);

  if (activeTrades.length === 0) return null;

  return (
    <Collapsible open={open} onOpenChange={setOpen}>
      <CollapsibleTrigger asChild>
        <div className="flex items-end justify-end">
          <CollapsibleToggler open={open} numberOfTrades={activeTrades.length} />
        </div>
      </CollapsibleTrigger>
      <CollapsibleContent>
        <ScrollArea>
          <div className="flex flex-col gap-1 mt-1">
            {activeTrades.map((tradeOffer) => (
              <TradeCard
                creator={tradeOffer.creator}
                offeredResources={tradeOffer.offer}
                opponents={tradeOffer.responses}
                requestedResources={tradeOffer.request}
                requester={tradeOffer.requester}
                tradeID={tradeOffer.id}
              />
            ))}
          </div>
        </ScrollArea>
      </CollapsibleContent>
    </Collapsible>
  );
};
