import { useEffect } from "react";

import type BaseMapRenderer from "@/core/maps/_base";
import { useMatchStore } from "@/state/match";

export function useUpdateBlockedTiles(instance: BaseMapRenderer | null) {
  const blockedTiles = useMatchStore((state) => state.blockedTiles);

  useEffect(() => {
    if (!instance) return;
    instance.updateRobbers(blockedTiles);
  }, [blockedTiles, instance]);
}
