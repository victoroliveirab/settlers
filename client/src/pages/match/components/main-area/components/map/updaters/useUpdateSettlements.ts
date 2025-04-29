import { useEffect } from "react";

import type BaseMapRenderer from "@/core/maps/_base";
import { audio } from "@/lib/audio";
import { useMatchStore } from "@/state/match";
import { useStableObject } from "@/hooks/useStableObject";

export function useUpdateSettlements(instance: BaseMapRenderer | null, tick: number) {
  const settlements = useStableObject(
    useMatchStore((state) => state.settlements),
    (prev, curr, key) => prev[+key]?.owner === curr[+key]?.owner,
  );

  useEffect(() => {
    if (!instance || tick === 0) return;
    instance.updateSettlements(settlements);
    audio.playAudio("settlement-building");
  }, [instance, settlements, tick]);
}
