import { useRef } from "react";

import { useWebSocket } from "./useWebSocket";

export function useRequestDiceRoll() {
  const { sendMessage } = useWebSocket();
  const fn = useRef(() => {
    sendMessage({ type: "match.dice-roll", payload: {} });
  });
  return fn.current;
}
