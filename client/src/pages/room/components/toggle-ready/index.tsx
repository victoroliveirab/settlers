import { Toggle } from "@/components/ui/toggle";
import { useWebSocket } from "@/hooks/useWebSocket";
import { useRoomStore } from "@/state/room";

export function ToggleReady() {
  const { sendMessage } = useWebSocket();
  const roomID = useRoomStore((state) => state.room.id);

  const onReadyChange = (ready: boolean) => {
    sendMessage({ type: "room.toggle-ready", payload: { ready, roomID } });
  };

  return (
    <Toggle
      className="bg-green-700 text-primary-foreground hover:bg-green-500"
      onPressedChange={onReadyChange}
    >
      Ready
    </Toggle>
  );
}
