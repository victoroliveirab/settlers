import { useEffect } from "react";

import type BaseMapRenderer from "@/core/maps/_base";
import { useMatchStore } from "@/state/match";

export function useUpdateRobbers(instance: BaseMapRenderer | null) {
  const robbers = useMatchStore((state) => state.robber);

  useEffect(() => {
    if (!instance) return;
    instance.updateTiles(robbers.availableTiles, robbers.enabled, robbers.highlight);
  }, [instance, robbers]);
}
