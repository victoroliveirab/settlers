import { useEffect, useState } from "react";

import { Die } from "@/components/custom/die";
import { cn } from "@/lib/utils";

import { useWebSocket } from "@/hooks/useWebSocket";
import { useMatchStore } from "@/state/match";

import { Pass } from "./components/pass";
import { Trade } from "./components/trade";
import { BuyDevelopmentCardButton } from "./components/development-card";

export const Actions = () => {
  const { sendMessage } = useWebSocket();
  const diceState = useMatchStore((state) => state.dice);
  const [lastDice, setLastDice] = useState(diceState.value);

  const onDiceClick = () => {
    if (diceState.enabled) {
      sendMessage({ type: "match.dice-roll", payload: {} });
    }
  };

  // Keep last dice value rendered
  useEffect(() => {
    if (diceState.value[0] === 0 && diceState.value[1] === 0) return;
    setLastDice(diceState.value);
  }, [diceState.value]);

  return (
    <section className="h-fit flex flex-col justify-between gap-2">
      <BuyDevelopmentCardButton />
      <ul className="flex gap-1">
        <li className="flex-1">
          <Pass />
        </li>
        <li className="flex-1">
          <Trade />
        </li>
      </ul>
      <ul
        className={cn("flex justify-center gap-2", {
          "cursor-pointer": diceState.enabled,
        })}
        onClick={onDiceClick}
      >
        <Die active={diceState.enabled} value={lastDice[0]} />
        <Die active={diceState.enabled} value={lastDice[1]} />
      </ul>
    </section>
  );
};
