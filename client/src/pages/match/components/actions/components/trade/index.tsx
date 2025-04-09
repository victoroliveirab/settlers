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

import { BankTrade } from "./components/bank-trade";
import { PortTrade } from "./components/port-trade";
import { PlayerTrade } from "./components/player-trade";

export const Trade = () => {
  const tradeState = useMatchStore((state) => state.actions.trade);

  return (
    <Dialog>
      <DialogTrigger className="w-full" asChild disabled={!tradeState}>
        <Button className="w-full" disabled={!tradeState}>
          Trade
        </Button>
      </DialogTrigger>
      <DialogContent className="w-3xl h-fit">
        <DialogHeader className="h-fit">
          <DialogTitle>Trade menu</DialogTitle>
          <DialogDescription>
            <Tabs defaultValue="player" className="w-full">
              <TabsList className="w-full">
                <TabsTrigger className="cursor-pointer" value="bank">
                  Bank
                </TabsTrigger>
                <TabsTrigger className="cursor-pointer" value="port">
                  Port
                </TabsTrigger>
                <TabsTrigger className="cursor-pointer" value="player">
                  Player
                </TabsTrigger>
              </TabsList>
              <TabsContent value="bank">
                <BankTrade />
              </TabsContent>
              <TabsContent value="port">
                <PortTrade />
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
