import { useEffect } from "react";

import type BaseMapRenderer from "@/core/maps/_base";
import { useMatchStore } from "@/state/match";

export function useUpdateSettlements(instance: BaseMapRenderer | null, tick: number) {
  const settlements = useMatchStore((state) => state.settlements);

  useEffect(() => {
    if (!instance || tick === 0) return;
    instance.updateSettlements(settlements);
  }, [instance, settlements, tick]);
}
