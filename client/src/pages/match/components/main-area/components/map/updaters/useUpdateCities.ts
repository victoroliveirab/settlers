import { useEffect } from "react";

import type BaseMapRenderer from "@/core/maps/_base";
import { useMatchStore } from "@/state/match";

export function useUpdateCities(instance: BaseMapRenderer | null, tick: number) {
  const cities = useMatchStore((state) => state.cities);

  useEffect(() => {
    if (!instance || tick === 0) return;
    instance.updateCities(cities);
  }, [cities, instance, tick]);
}
