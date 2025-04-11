import { emojis } from "@/core/constants";

interface ICardProps {
  resource: SettlersCore.Resource;
}

export const Card = ({ resource }: ICardProps) => {
  const emoji = emojis.resources[resource];
  return (
    <li className="h-full aspect-[3/4] rounded-md bg-neutral-300 flex items-center justify-center">
      <span>{emoji}</span>
    </li>
  );
};
