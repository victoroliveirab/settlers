import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";

import { useRoomStore } from "@/state/room";

import { ParticipantList } from "./components/participant-list";
import { RoomParameters } from "./components/room-parameters";
import { ToggleReady } from "./components/toggle-ready";
import { StartButton } from "./components/button-start";

export const Room = () => {
  const room = useRoomStore((state) => state.room);
  return (
    <main className="h-full flex items-center justify-center w-5xl mx-auto">
      <Card className="w-full max-h-[80vh] overflow-hidden">
        <CardHeader>
          <CardTitle>Room #{room.id}</CardTitle>
          <CardDescription>Card Description</CardDescription>
        </CardHeader>
        <CardContent className="flex gap-4 justify-center h-[50vh]">
          <div className="w-48">
            <ScrollArea className="h-full max-h-full">
              <ParticipantList />
            </ScrollArea>
          </div>
          <div className="flex-auto basis-0">
            <ScrollArea className="h-full max-h-full">
              <RoomParameters />
            </ScrollArea>
          </div>
        </CardContent>
        <CardFooter>
          <div className="w-full flex items-center justify-end gap-4">
            <ToggleReady />
            <StartButton />
          </div>
        </CardFooter>
      </Card>
    </main>
  );
};
