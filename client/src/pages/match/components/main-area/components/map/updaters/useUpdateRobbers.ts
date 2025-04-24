import { useEffect } from "react";

import type BaseMapRenderer from "@/core/maps/_base";
import { useMatchStore } from "@/state/match";

export function useUpdateRobbers(instance: BaseMapRenderer | null, tick: number) {
  const robbers = useMatchStore((state) => state.robber);

  useEffect(() => {
    if (!instance || tick === 0) return;
    instance.updateTiles(robbers.availableTiles, robbers.enabled, robbers.highlight);
  }, [instance, robbers, tick]);
}
