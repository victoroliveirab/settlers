import { useEffect } from "react";

import type BaseMapRenderer from "@/core/maps/_base";
import { useMatchStore } from "@/state/match";

export function useUpdateCities(instance: BaseMapRenderer | null) {
  const cities = useMatchStore((state) => state.cities);

  useEffect(() => {
    if (!instance) return;
    instance.updateCities(cities);
  }, [cities, instance]);
}
