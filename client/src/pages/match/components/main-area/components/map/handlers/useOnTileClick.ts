import { useRef } from "react";

import { useWebSocket } from "@/hooks/useWebSocket";

export function useOnTileClick() {
  const { sendMessage } = useWebSocket();
  const fn = useRef((tileID: number) => {
    sendMessage({ type: "match.tile-click", payload: { tile: tileID } });
  });
  return fn.current;
}
