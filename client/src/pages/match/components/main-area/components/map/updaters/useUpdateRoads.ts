import { useEffect } from "react";

import type BaseMapRenderer from "@/core/maps/_base";
import { audio } from "@/lib/audio";
import { useStableObject } from "@/hooks/useStableObject";
import { useMatchStore } from "@/state/match";

export function useUpdateRoads(instance: BaseMapRenderer | null, tick: number) {
  const roads = useStableObject(
    useMatchStore((state) => state.roads),
    (prev, curr, key) => prev[+key]?.owner === curr[+key]?.owner,
  );

  useEffect(() => {
    if (!instance || tick === 0) return;
    instance.updateRoads(roads);
    audio.playAudio("road-building");
  }, [instance, roads, tick]);
}
