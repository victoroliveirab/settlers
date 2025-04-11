import { useEffect } from "react";

import type BaseMapRenderer from "@/core/maps/_base";
import { useMatchStore } from "@/state/match";

export function useUpdateVertices(instance: BaseMapRenderer | null) {
  const verticesState = useMatchStore((state) => state.vertices);

  useEffect(() => {
    if (!instance) return;
    instance.updateVertices(
      verticesState.availableForSettlement,
      verticesState.availableForCity,
      verticesState.enabled,
      verticesState.highlight,
    );
  }, [instance, verticesState]);
}
