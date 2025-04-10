import { type ReactNode, useState } from "react";

import { GameCard } from "@/components/custom/game-card";
import { QuantitySelector } from "@/components/custom/quantity-selector";

import { resourcesOrder } from "@/core/constants";
import { sumOfResources } from "@/core/utils";

interface IDiscardProps {
  children: (props: {
    dirty: boolean;
    selected: SettlersCore.ResourceCollection;
    totalSelected: number;
  }) => ReactNode;
  resources: SettlersCore.ResourceCollection;
}

export const Discard = ({ children, resources }: IDiscardProps) => {
  const [state, setState] = useState({
    Lumber: 0,
    Brick: 0,
    Sheep: 0,
    Grain: 0,
    Ore: 0,
  });

  const totalSelected = sumOfResources(state);

  return (
    <>
      <ul className="flex justify-center gap-4">
        {resourcesOrder.map((resource) => (
          <li key={resource} className="flex flex-col gap-2 text-center">
            <div className="flex justify-center">
              <GameCard size="md" value={resource} />
            </div>
            <QuantitySelector
              min={0}
              max={resources[resource]}
              onValueChange={(value) => {
                setState({
                  ...state,
                  [resource]: value,
                });
              }}
              value={state[resource]}
            />
          </li>
        ))}
      </ul>
      {children({
        dirty: totalSelected > 0,
        selected: state,
        totalSelected,
      })}
    </>
  );
};
