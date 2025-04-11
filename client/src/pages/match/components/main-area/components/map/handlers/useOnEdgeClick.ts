import { useRef } from "react";

import { useWebSocket } from "@/hooks/useWebSocket";

export function useOnEdgeClick() {
  const { sendMessage } = useWebSocket();
  const fn = useRef((edgeID: number) => {
    sendMessage({ type: "match.edge-click", payload: { edge: edgeID } });
  });
  return fn.current;
}
