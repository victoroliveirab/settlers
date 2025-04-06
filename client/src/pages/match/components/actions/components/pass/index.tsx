import { Button } from "@/components/ui/button";

import { useWebSocket } from "@/hooks/useWebSocket";
import { useMatchStore } from "@/state/match";

export const Pass = () => {
  const { sendMessage } = useWebSocket();
  const enabled = useMatchStore((state) => state.actions.pass);

  const onPassClick = () => {
    if (enabled) {
      sendMessage({ type: "match.end-round", payload: {} });
    }
  };

  return (
    <Button className="w-full" disabled={!enabled} onClick={onPassClick}>
      Pass
    </Button>
  );
};
