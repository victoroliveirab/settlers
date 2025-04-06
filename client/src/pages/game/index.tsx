import { useRoomStore } from "@/state/room";

import { Room } from "../room";
import { Match } from "../match";

export const Game = () => {
  const room = useRoomStore((state) => state.room);
  if (room.status === "prematch") {
    return <Room />;
  }
  return <Match />;
};
