import { type ReactNode, useState } from "react";

import { GameCard } from "@/components/custom/game-card";
import { QuantitySelector } from "@/components/custom/quantity-selector";

import { resourcesOrder } from "@/core/constants";
import { sumOfResources } from "@/core/utils";

import { isDirty } from "./utils";

const defaultInitialState: SettlersCore.ResourceCollection = {
  Lumber: 0,
  Brick: 0,
  Sheep: 0,
  Grain: 0,
  Ore: 0,
};

interface ITradeProps {
  className?: string;
  children: (props: {
    dirty: boolean;
    given: SettlersCore.ResourceCollection;
    requested: SettlersCore.ResourceCollection;
    totalGiven: number;
    totalRequested: number;
  }) => ReactNode;
  givenResourcesAvailable: SettlersCore.ResourceCollection;
  givenResourceIsDisabledGetter?: (resource: SettlersCore.Resource) => boolean;
  givenStep?: number;
  givenStepGetter?: (resource: SettlersCore.Resource) => number;
  initialStateGiven?: SettlersCore.ResourceCollection;
  initialStateRequested?: SettlersCore.ResourceCollection;
}

export const Trade = ({
  className,
  children,
  givenResourcesAvailable,
  givenResourceIsDisabledGetter,
  givenStep,
  givenStepGetter,
  initialStateGiven = defaultInitialState,
  initialStateRequested = defaultInitialState,
}: ITradeProps) => {
  const [given, setGiven] = useState(initialStateGiven);
  const [requested, setRequested] = useState(initialStateRequested);

  const totalGiven = sumOfResources(given);
  const totalRequested = sumOfResources(requested);

  const isGivenResourceDisabled = (resource: SettlersCore.Resource) => {
    const alreadyAddedAtRequested = requested[resource] > 0;
    if (alreadyAddedAtRequested) return true;

    if (givenResourceIsDisabledGetter) {
      return givenResourceIsDisabledGetter(resource);
    }
    return false;
  };

  const isRequestedResourceDisabled = (resource: SettlersCore.Resource) => {
    return given[resource] > 0;
  };

  return (
    <div className={className}>
      <div className="flex flex-col gap-2">
        <h3>You give:</h3>
        <ul className="flex justify-center gap-4">
          {resourcesOrder.map((resource) => (
            <li key={resource} className="flex flex-col gap-2 text-center">
              <div className="flex justify-center">
                <GameCard size="md" value={resource} />
              </div>
              <QuantitySelector
                disabled={isGivenResourceDisabled(resource)}
                min={0}
                max={givenResourcesAvailable[resource]}
                onValueChange={(value) => {
                  setGiven({
                    ...given,
                    [resource]: value,
                  });
                }}
                step={givenStepGetter ? givenStepGetter(resource) ?? 1 : givenStep ?? 1}
                value={given[resource]}
              />
            </li>
          ))}
        </ul>
      </div>
      <div className="flex flex-col gap-2">
        <h3>You receive:</h3>
        <ul className="flex justify-center gap-4">
          {resourcesOrder.map((resource) => (
            <li key={resource} className="flex flex-col gap-2 text-center">
              <div className="flex justify-center">
                <GameCard size="md" value={resource} />
              </div>
              <QuantitySelector
                disabled={isRequestedResourceDisabled(resource)}
                onValueChange={(value) => {
                  setRequested({
                    ...requested,
                    [resource]: value,
                  });
                }}
                value={requested[resource]}
              />
            </li>
          ))}
        </ul>
      </div>
      {children({
        dirty: isDirty({
          currentState: {
            given,
            requested,
          },
          initialState: {
            given: initialStateGiven,
            requested: initialStateRequested,
          },
        }),
        given,
        requested,
        totalGiven,
        totalRequested,
      })}
    </div>
  );
};
