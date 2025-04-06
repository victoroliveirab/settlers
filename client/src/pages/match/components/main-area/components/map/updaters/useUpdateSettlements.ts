import { useEffect } from "react";

import type BaseMapRenderer from "@/core/maps/_base";
import { useMatchStore } from "@/state/match";

export function useUpdateSettlements(instance: BaseMapRenderer | null) {
  const settlements = useMatchStore((state) => state.settlements);

  useEffect(() => {
    if (!instance) return;
    instance.updateSettlements(settlements);
  }, [instance, settlements]);
}
