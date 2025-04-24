import { useEffect } from "react";

import type BaseMapRenderer from "@/core/maps/_base";
import { useMatchStore } from "@/state/match";

export function useUpdateBlockedTiles(instance: BaseMapRenderer | null, tick: number) {
  const blockedTiles = useMatchStore((state) => state.blockedTiles);

  useEffect(() => {
    if (!instance || tick === 0) return;
    instance.updateRobbers(blockedTiles);
  }, [blockedTiles, instance, tick]);
}
