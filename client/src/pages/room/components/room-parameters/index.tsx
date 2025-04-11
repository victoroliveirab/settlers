import { useRoomStore } from "@/state/room";

import { CardParameter } from "./components/card-parameter";

export function RoomParameters() {
  const params = useRoomStore((state) => state.params);
  return (
    <div className="grid grid-cols-3 gap-4 pr-4">
      {params.map((entry) => (
        <CardParameter key={entry.key} param={entry} />
      ))}
    </div>
  );
}
