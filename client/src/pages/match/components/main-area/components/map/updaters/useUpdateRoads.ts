import { useEffect } from "react";

import type BaseMapRenderer from "@/core/maps/_base";
import { useMatchStore } from "@/state/match";

export function useUpdateRoads(instance: BaseMapRenderer | null, tick: number) {
  const roads = useMatchStore((state) => state.roads);

  useEffect(() => {
    if (!instance || tick === 0) return;
    instance.updateRoads(roads);
  }, [instance, roads, tick]);
}
