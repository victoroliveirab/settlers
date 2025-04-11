import { useContext } from "react";

import { WebSocketContext } from "@/context";

export function useWebSocket() {
  return useContext(WebSocketContext);
}
