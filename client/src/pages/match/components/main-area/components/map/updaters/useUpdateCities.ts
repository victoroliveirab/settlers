import { useEffect } from "react";

import type BaseMapRenderer from "@/core/maps/_base";
import { audio } from "@/lib/audio";
import { useStableObject } from "@/hooks/useStableObject";
import { useMatchStore } from "@/state/match";

export function useUpdateCities(instance: BaseMapRenderer | null, tick: number) {
  const cities = useStableObject(
    useMatchStore((state) => state.cities),
    (prev, curr, key) => prev[+key]?.owner === curr[+key]?.owner,
  );

  useEffect(() => {
    if (!instance || tick === 0) return;
    instance.updateCities(cities);
    audio.playAudio("city-building");
  }, [cities, instance, tick]);
}
