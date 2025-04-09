import { useEffect, useRef } from "react";

import { useMatchStore } from "./match";

export function useStateDebug() {
  const matchState = useMatchStore((state) => state);
  const prevMatchState = useRef<typeof matchState | null>(null);
  useEffect(() => {
    console.debug("==============MATCHSTATE_CHANGE==============");
    console.debug("PREV:", prevMatchState.current);
    console.debug("CURR:", { ...matchState });
    console.debug("=============================================");
    prevMatchState.current = { ...matchState };
  }, [matchState]);
}
