import { useRoomStore } from "@/state/room";

import { Room } from "../room";
import { Match } from "../match";
import { PostMatch } from "../post-match";

export const Game = () => {
  const room = useRoomStore((state) => state.room);
  if (room.status === "prematch") {
    return <Room />;
  } else if (room.status === "over") {
    return <PostMatch />;
  }
  return <Match />;
};
