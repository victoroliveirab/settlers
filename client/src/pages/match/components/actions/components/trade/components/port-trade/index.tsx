import { useMemo, useState } from "react";

import { Button } from "@/components/ui/button";
import { DialogClose } from "@/components/ui/dialog";
import { GameCard } from "@/components/custom/game-card";
import { QuantitySelector } from "@/components/custom/quantity-selector";

import { resourcesOrder } from "@/core/constants";
import { useMatchStore } from "@/state/match";

// This is V1, we can do so much better with this
// But let's go the simplest path
export const PortTrade = () => {
  const [state, setState] = useState({
    given: {
      Lumber: 0,
      Brick: 0,
      Sheep: 0,
      Grain: 0,
      Ore: 0,
    },
    taken: {
      Lumber: 0,
      Brick: 0,
      Sheep: 0,
      Grain: 0,
      Ore: 0,
    },
  });
  const hand = useMatchStore((state) => state.hand);

  const isTradeDisabled = useMemo(() => {
    const totalGiven = Object.values(state.given).reduce((acc, value) => acc + value, 0);
    if (totalGiven === 0) return true;
    const totalTaken = Object.values(state.taken).reduce((acc, value) => acc + value, 0);
    if (totalTaken === 0) return true;

    const targetNumberOfTaken = totalGiven / 4;
    if (totalTaken !== targetNumberOfTaken) return true;

    for (const [resource, quantity] of Object.entries(state.given)) {
      if (quantity > hand[resource as SettlersCore.Resource]) return true;
    }

    return false;
  }, [state]);

  return (
    <div className="flex flex-col gap-4 py-2">
      <ul className="flex justify-center gap-6">
        {resourcesOrder.map((resource) => (
          <li className="flex flex-col gap-2 text-center">
            <GameCard className="h-16" value={resource} />
            <QuantitySelector
              min={0}
              max={16}
              // max={Math.floor(hand[resource] / 4)}
              onValueChange={(value) => {
                setState({
                  ...state,
                  given: {
                    ...state.given,
                    [resource]: value,
                  },
                });
              }}
              value={state.given[resource]}
            />
          </li>
        ))}
      </ul>
      <ul className="flex justify-center gap-6">
        {resourcesOrder.map((resource) => (
          <li className="flex flex-col gap-3 text-center">
            <GameCard className="h-16" value={resource} />
            <QuantitySelector
              onValueChange={(value) => {
                setState({
                  ...state,
                  taken: {
                    ...state.taken,
                    [resource]: value,
                  },
                });
              }}
              value={state.taken[resource]}
            />
          </li>
        ))}
      </ul>
      <ul className="flex items-center justify-end gap-1">
        <li>
          <DialogClose asChild>
            <Button>Cancel</Button>
          </DialogClose>
        </li>
        <li>
          <Button disabled={isTradeDisabled}>Trade</Button>
        </li>
      </ul>
    </div>
  );
};
