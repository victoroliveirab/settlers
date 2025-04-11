import { useState } from "react";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

import { YearOfPlenty } from "@/features/year-of-plenty";

import { useWebSocket } from "@/hooks/useWebSocket";
import { useMatchStore } from "@/state/match";

export const YearOfPlentyModal = () => {
  const [resource1, setResource1] = useState<SettlersCore.Resource | null>(null);
  const [resource2, setResource2] = useState<SettlersCore.Resource | null>(null);
  const { sendMessage } = useWebSocket();

  const yearOfPlentyState = useMatchStore((state) => state.yearOfPlenty);

  const submitRequest = () => {
    if (!resource1 || !resource2) return;
    sendMessage({ type: "match.year-of-plenty", payload: { resource1, resource2 } });
  };

  return (
    <Dialog open={yearOfPlentyState.enabled}>
      <DialogContent className="w-fit min-w-60">
        <DialogHeader>
          <DialogTitle>Pick Resources</DialogTitle>
          <DialogDescription className="flex flex-col gap-4">
            <YearOfPlenty
              onClickResource1={setResource1}
              onClickResource2={setResource2}
              resource1={resource1}
              resource2={resource2}
            />
            <div className="flex justify-end">
              <Button disabled={!resource1 || !resource2} onClick={submitRequest}>
                Submit
              </Button>
            </div>
          </DialogDescription>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
};
