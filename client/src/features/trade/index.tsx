import { type ReactNode, useState } from "react";

import { GameCard } from "@/components/custom/game-card";
import { QuantitySelector } from "@/components/custom/quantity-selector";
import { cn } from "@/lib/utils";

import { resourcesOrder } from "@/core/constants";

import { isDirty } from "./utils";

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
  initialStateGiven: SettlersCore.ResourceCollection;
  initialStateRequested: SettlersCore.ResourceCollection;
}

export const Trade = ({
  className,
  children,
  givenResourcesAvailable,
  initialStateGiven,
  initialStateRequested,
}: ITradeProps) => {
  const [given, setGiven] = useState(initialStateGiven);
  const [requested, setRequested] = useState(initialStateRequested);

  const totalGiven = Object.values(given).reduce((acc, value) => acc + value, 0);
  const totalRequested = Object.values(requested).reduce((acc, value) => acc + value, 0);

  return (
    <div className={cn("flex flex-col gap-4", className)}>
      <div className="flex flex-col gap-2">
        <h3>You give:</h3>
        <ul className="flex justify-center gap-6">
          {resourcesOrder.map((resource) => (
            <li className="flex flex-col gap-2 text-center">
              <GameCard className="h-16" value={resource} />
              <QuantitySelector
                min={0}
                max={givenResourcesAvailable[resource]}
                onValueChange={(value) => {
                  setGiven({
                    ...given,
                    [resource]: value,
                  });
                }}
                value={given[resource]}
              />
            </li>
          ))}
        </ul>
      </div>
      <div className="flex flex-col gap-2">
        <h3>You receive:</h3>
        <ul className="flex justify-center gap-6">
          {resourcesOrder.map((resource) => (
            <li className="flex flex-col gap-3 text-center">
              <GameCard className="h-16" value={resource} />
              <QuantitySelector
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
