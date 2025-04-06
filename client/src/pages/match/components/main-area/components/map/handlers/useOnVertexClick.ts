import { useRef } from "react";

import { useWebSocket } from "@/hooks/useWebSocket";

export function useOnVertexClick() {
  const { sendMessage } = useWebSocket();
  const fn = useRef((vertexID: number) => {
    sendMessage({ type: "match.vertex-click", payload: { vertex: vertexID } });
  });
  return fn.current;
}
