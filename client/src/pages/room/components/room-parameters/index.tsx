import { useCallback, useEffect, useState } from "react";

import { useWebSocket } from "@/hooks/useWebSocket";
import { useRoomStore } from "@/state/room";
import { usePlayerStore } from "@/state/player";

import { CardParameter } from "./components/card-parameter";

export function RoomParameters() {
  const [isLoading, setIsLoading] = useState(false);
  const room = useRoomStore((state) => state.room);
  const params = useRoomStore((state) => state.params);
  const username = usePlayerStore((state) => state.username);
  const { sendMessage } = useWebSocket();

  const updateParam = useCallback((key: string, value: number) => {
    setIsLoading(true);
    sendMessage({ type: "room.update-param", payload: { key, value } });
  }, []);

  useEffect(() => {
    setIsLoading(false);
  }, [params]);

  return (
    <div className="grid grid-cols-3 gap-4 pr-4">
      {params.map((entry) => (
        <CardParameter
          key={entry.key}
          disabled={room.owner !== username}
          isLoading={isLoading}
          onChange={updateParam}
          param={entry}
        />
      ))}
    </div>
  );
}
