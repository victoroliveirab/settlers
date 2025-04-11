import { Button } from "@/components/ui/button";

import { useWebSocket } from "@/hooks/useWebSocket";
import { usePlayerStore } from "@/state/player";
import { useRoomStore } from "@/state/room";

export const StartButton = () => {
  const { sendMessage } = useWebSocket();
  const room = useRoomStore((state) => state.room);
  const currentUsername = usePlayerStore((state) => state.username);

  const isEnabled =
    currentUsername === room.owner &&
    room.participants.every((entry) => !!entry.player && entry.ready);

  const onClick = () => {
    sendMessage({ type: "room.start-game", payload: {} });
  };

  return (
    <Button className="w-20" disabled={!isEnabled} onClick={onClick}>
      Start
    </Button>
  );
};
