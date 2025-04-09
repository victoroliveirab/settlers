import { GameCard } from "@/components/custom/game-card";

import { resourcesOrder } from "@/core/constants";

interface IMonopolyProps {
  onClick: (resource: SettlersCore.Resource) => void;
}

export const Monopoly = ({ onClick }: IMonopolyProps) => {
  return (
    <ul>
      {resourcesOrder.map((resource) => (
        <li key={resource}>
          <GameCard
            className="h-16 cursor-pointer"
            value={resource}
            onClick={() => onClick(resource)}
          />
        </li>
      ))}
    </ul>
  );
};
