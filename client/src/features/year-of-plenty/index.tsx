import { GameCard } from "@/components/custom/game-card";
import { cn } from "@/lib/utils";

import { resourcesOrder } from "@/core/constants";

interface IYearOfPlentyProps {
  onClickResource1: (resource: SettlersCore.Resource) => void;
  onClickResource2: (resource: SettlersCore.Resource) => void;
  resource1: SettlersCore.Resource | null;
  resource2: SettlersCore.Resource | null;
}

export const YearOfPlenty = ({
  onClickResource1,
  onClickResource2,
  resource1,
  resource2,
}: IYearOfPlentyProps) => {
  return (
    <>
      <ul className="flex justify-center gap-6">
        {resourcesOrder.map((resource) => (
          <li key={resource}>
            <GameCard
              className={cn("h-16 cursor-pointer", {
                "border-2 border-solid border-black": resource === resource1,
              })}
              value={resource}
              onClick={() => onClickResource1(resource)}
            />
          </li>
        ))}
      </ul>
      <ul className="flex justify-center gap-6">
        {resourcesOrder.map((resource) => (
          <li key={resource}>
            <GameCard
              className={cn("h-16 cursor-pointer", {
                "border-2 border-solid border-black": resource === resource2,
              })}
              value={resource}
              onClick={() => onClickResource2(resource)}
            />
          </li>
        ))}
      </ul>
    </>
  );
};
