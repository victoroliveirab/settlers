import { Avatar } from "@/components/custom/avatar";

interface IPickRobbedProps {
  onClick: (player: string) => void;
  players: SettlersCore.Player[];
}

export const PickRobbed = ({ onClick, players }: IPickRobbedProps) => {
  return (
    <ul className="flex gap-4">
      {players.map((player) => (
        <li
          key={player.name}
          className="flex flex-col items-center gap-1 w-24 cursor-pointer"
          onClick={() => onClick(player.name)}
        >
          <Avatar key={player.name} background={player.color.background} />
          <h3 className="text-center">{player.name}</h3>
        </li>
      ))}
    </ul>
  );
};
