import { useEffect } from "react";

import type BaseMapRenderer from "@/core/maps/_base";
import { useMatchStore } from "@/state/match";

export function useUpdateEdges(instance: BaseMapRenderer | null) {
  const edgesState = useMatchStore((state) => state.edges);

  useEffect(() => {
    if (!instance) return;
    instance.updateEdges(edgesState.available, edgesState.enabled, edgesState.highlight);
  }, [edgesState, instance]);
}
