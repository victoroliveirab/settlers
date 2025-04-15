import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

import { useMatchStore } from "@/state/match";

import { GeneralTrade } from "./components/general-trade";
import { PortTrade } from "./components/port-trade";
import { PlayerTrade } from "./components/player-trade";

export const Trade = () => {
  const tradeState = useMatchStore((state) => state.actions.trade);
  const hand = useMatchStore((state) => state.hand);
  const ports = useMatchStore((state) => state.ownedPorts);

  const hasGeneralPort = ports.includes("General");
  const bankStep = hasGeneralPort ? 3 : 4; // TODO: get from game params

  const isBankTradeDisabled = Object.values(hand).every((quantity) => quantity < bankStep);
  const isPortTradeDisabled = ports.filter((type) => type !== "General").length === 0;
  const isPlayerTradeDisabled = Object.values(hand).every((quantity) => quantity === 0);

  return (
    <Dialog>
      <DialogTrigger className="w-full" asChild disabled={!tradeState}>
        <Button className="w-full" disabled={!tradeState}>
          Trade
        </Button>
      </DialogTrigger>
      <DialogContent className="w-fit h-fit">
        <DialogHeader>
          <DialogTitle>Trade menu</DialogTitle>
          <DialogDescription className="flex flex-col gap-4">
            <Tabs defaultValue="player" className="w-full">
              <TabsList className="w-full">
                <TabsTrigger className="cursor-pointer" value="bank" disabled={isBankTradeDisabled}>
                  {hasGeneralPort ? "General" : "Bank"}
                </TabsTrigger>
                <TabsTrigger className="cursor-pointer" value="port" disabled={isPortTradeDisabled}>
                  Port
                </TabsTrigger>
                <TabsTrigger
                  className="cursor-pointer"
                  value="player"
                  disabled={isPlayerTradeDisabled}
                >
                  Player
                </TabsTrigger>
              </TabsList>
              <TabsContent value="bank">
                <GeneralTrade isBankTrade={!hasGeneralPort} step={bankStep} />
              </TabsContent>
              <TabsContent value="port">
                <PortTrade step={2} />
              </TabsContent>
              <TabsContent value="player">
                <PlayerTrade />
              </TabsContent>
            </Tabs>
          </DialogDescription>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
};
