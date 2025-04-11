import { useEffect } from "react";

import type BaseMapRenderer from "@/core/maps/_base";
import { useMatchStore } from "@/state/match";

export function useDarkenTilesAfterDiceRoll(instance: BaseMapRenderer | null) {
  const map = useMatchStore((state) => state.map);
  const dice = useMatchStore((state) => state.dice.value);

  useEffect(() => {
    if (!instance) return;
    if (dice[0] === 0 || dice[1] === 0) return;
    const sum = dice[0] + dice[1];
    const tilesIDs = map.filter(({ token }) => token === sum).map(({ id }) => id);
    instance.darkenTiles(tilesIDs);
  }, [dice, instance, map]);
}
