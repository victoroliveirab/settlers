import { useEffect } from "react";

import type BaseMapRenderer from "@/core/maps/_base";
import { useMatchStore } from "@/state/match";

export function useUpdateRoads(instance: BaseMapRenderer | null) {
  const roads = useMatchStore((state) => state.roads);

  useEffect(() => {
    if (!instance) return;
    instance.updateRoads(roads);
  }, [instance, roads]);
}
