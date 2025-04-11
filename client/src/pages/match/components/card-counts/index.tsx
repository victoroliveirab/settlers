import { emojis, resourcesOrder } from "@/core/constants";

export const CardCounts = () => {
  return (
    <section className="h-fit">
      <ul className="flex flex-col gap-1 bg-white">
        {resourcesOrder.map((resource) => (
          <li key={resource} className="inline-flex gap-1">
            <span>{emojis.resources[resource]}:</span>
            <span>12</span>
          </li>
        ))}
        <li className="inline-flex gap-1">
          <span>{emojis.devCards.Knight}:</span>
          <span>12</span>
        </li>
      </ul>
    </section>
  );
};
