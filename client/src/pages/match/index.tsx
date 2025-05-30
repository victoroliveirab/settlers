import { Card, CardContent } from "@/components/ui/card";

import { Actions } from "./components/actions";
import { Logs } from "./components/logs";
import { MainArea } from "./components/main-area";
import { Modals } from "./components/modals";

export const Match = () => {
  return (
    <main className="h-full p-6 flex">
      <Card className="w-6xl mx-auto bg-neutral-800 h-full">
        <CardContent className="flex gap-4 h-full relative">
          <div className="h-full flex flex-col gap-4 w-1/4 min-w-44">
            <Logs />
            <hr className="w-full h-px bg-border" />
            <Actions />
          </div>
          <hr className="w-px h-full bg-border" />
          <MainArea />
        </CardContent>
      </Card>
      <Modals />
    </main>
  );
};
