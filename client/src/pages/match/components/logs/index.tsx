import { useEffect, useRef } from "react";

import { ScrollArea } from "@/components/ui/scroll-area";

import { useMatchStore } from "@/state/match";

import { parseLog } from "./parser";

export const Logs = () => {
  const ref = useRef<HTMLDivElement>(null);
  const logsRef = useRef<HTMLUListElement>(null);
  const logs = useMatchStore((state) => state.logs);

  useEffect(() => {
    logs.forEach((log) => {
      const root = document.createElement("li");
      root.classList.add("flex", "items-center", "flex-wrap", "gap-1", "text-neutral-100");
      const success = parseLog(root, log);
      if (success) {
        logsRef.current?.appendChild(root);
      }
    });
    ref.current?.scrollTo(0, ref.current.scrollHeight);
  }, [logs]);

  return (
    <ScrollArea ref={ref}>
      <ul className="h-full w-full flex flex-col gap-2" ref={logsRef} />
    </ScrollArea>
  );
};
