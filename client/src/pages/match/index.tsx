import { Card, CardContent } from "@/components/ui/card";

import { Actions } from "./components/actions";
import { Logs } from "./components/logs";
import { MainArea } from "./components/main-area";
import { CardCounts } from "./components/card-counts";
import { ScrollArea } from "@/components/ui/scroll-area";

export const Match = () => {
  return (
    <main className="h-full p-6 flex">
      <Card className="w-6xl mx-auto bg-neutral-800 h-full">
        <CardContent className="flex gap-4 h-full">
          <div className="h-full flex flex-col gap-4 w-1/5 min-w-44">
            <ScrollArea>
              <Logs />
            </ScrollArea>
            <hr className="w-full h-px bg-border" />
            <CardCounts />
            <Actions />
          </div>
          <hr className="w-px h-full bg-border" />
          <MainArea />
        </CardContent>
      </Card>
    </main>
  );
};
