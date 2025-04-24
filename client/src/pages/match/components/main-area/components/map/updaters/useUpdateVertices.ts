import { useEffect } from "react";

import type BaseMapRenderer from "@/core/maps/_base";
import { useMatchStore } from "@/state/match";

export function useUpdateVertices(instance: BaseMapRenderer | null, tick: number) {
  const verticesState = useMatchStore((state) => state.vertices);

  useEffect(() => {
    if (!instance || tick === 0) return;
    instance.updateVertices(
      verticesState.availableForSettlement,
      verticesState.availableForCity,
      verticesState.enabled,
      verticesState.highlight,
    );
  }, [instance, tick, verticesState]);
}
