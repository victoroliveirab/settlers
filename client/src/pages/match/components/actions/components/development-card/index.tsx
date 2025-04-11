import { Button } from "@/components/ui/button";

import { useWebSocket } from "@/hooks/useWebSocket";
import { useMatchStore } from "@/state/match";

export const BuyDevelopmentCardButton = () => {
  const { sendMessage } = useWebSocket();
  const enabled = useMatchStore((state) => state.actions.buyDevCard);

  const onClick = () => {
    sendMessage({ type: "match.buy-dev-card", payload: {} });
  };

  return (
    <Button disabled={!enabled} onClick={onClick}>
      Buy Development Card
    </Button>
  );
};
