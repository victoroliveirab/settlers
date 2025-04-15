import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

import { useWebSocket } from "@/hooks/useWebSocket";
import { useRoomStore } from "@/state/room";

export const RoomCapacitySelect = () => {
  const [minPlayers, maxPlayers] = useRoomStore((state) => state.minMaxPlayers);
  const capacity = useRoomStore((state) => state.room.capacity);
  const { sendMessage } = useWebSocket();

  const onChange = (value: number) => {
    sendMessage({ type: "room.update-capacity", payload: { capacity: value } });
  };

  const options: number[] = [];
  for (let option = minPlayers; option <= maxPlayers; ++option) {
    options.push(option);
  }

  return (
    <Select defaultValue={String(capacity)} onValueChange={(value) => onChange(+value)}>
      <SelectTrigger>
        <SelectValue placeholder={capacity} />
      </SelectTrigger>
      <SelectContent>
        {options.map((option) => (
          <SelectItem key={option} value={String(option)}>
            {option}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
};
