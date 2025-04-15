import { useRoomStore } from "@/state/room";

import { ParticipantAvatar } from "./components/participant-avatar";
import { Check } from "lucide-react";

export function ParticipantList() {
  const { capacity, participants, owner } = useRoomStore((state) => state.room);
  return (
    <ul className="h-fit">
      {Array.from({ length: capacity }).map((_, index) => {
        const participant = participants[index] ?? {
          bot: false,
          player: null,
          ready: false,
        };
        return (
          <li key={index} className="flex flex-col items-center not-last:mb-4">
            <ParticipantAvatar
              color={participant.player?.color.background}
              empty={!participant.player}
              owner={participant.player?.name === owner}
            >
              {participant.ready && (
                <Check
                  className="absolute top-0 right-0"
                  color={participant.player?.color.foreground}
                />
              )}
            </ParticipantAvatar>
            <p className="flex items-center truncate min-w-0 text-center">
              {participant.player?.name ?? "empty"}
            </p>
          </li>
        );
      })}
    </ul>
  );
}
