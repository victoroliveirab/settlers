import { Check, Ellipsis, XIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

import { useWebSocket } from "@/hooks/useWebSocket";
import { useMatchStore } from "@/state/match";

interface IOpponentsStatusProps {
  data: Record<SettlersCore.Player["name"], { blocked: boolean; status: string }>;
  disabled: boolean;
  tradeID: number;
}

export const OpponentsStatus = ({ data, disabled, tradeID }: IOpponentsStatusProps) => {
  const { sendMessage } = useWebSocket();
  const players = useMatchStore((state) => state.players);

  const entries = Object.entries(data);

  const onButtonClick = (index: number) => {
    if (disabled) return;
    const [playerName, info] = entries[index];
    if (info.status === "Accepted") {
      sendMessage({ type: "match.finalize-trade-offer", payload: { player: playerName, tradeID } });
    }
  };

  return (
    <ul className="flex gap-1">
      {entries.map(([playerName, data], index) => {
        const playerColor = players.find((player) => player.name === playerName)?.color;
        return (
          <Button
            size="xs"
            className={cn({
              "cursor-default": disabled || data.status !== "Accepted",
            })}
            style={{
              background: playerColor?.background,
              color: playerColor?.foreground,
            }}
            onClick={() => onButtonClick(index)}
          >
            {data.status === "Accepted" ? (
              <Check />
            ) : data.status === "Declined" ? (
              <XIcon />
            ) : (
              <Ellipsis />
            )}
          </Button>
        );
      })}
    </ul>
  );
};
