import { GameCard } from "@/components/custom/game-card";

import { resourcesOrder } from "@/core/constants";

interface IResourceListProps {
  resources: Record<SettlersCore.Resource, number>;
}

export const ResourceList = ({ resources }: IResourceListProps) => {
  return (
    <ul className="flex gap-1 justify-end">
      {resourcesOrder.map((resource) => {
        if (resources[resource] === 0) return null;
        return Array.from({ length: resources[resource] }).map((_, index) => (
          <GameCard key={`${resource}-${index}`} as="li" size="xs" value={resource} />
        ));
      })}
    </ul>
  );
};
