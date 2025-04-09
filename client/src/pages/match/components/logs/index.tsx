import { useEffect, useRef } from "react";

import { ScrollArea } from "@/components/ui/scroll-area";

import { useMatchStore } from "@/state/match";

export const Logs = () => {
  const ref = useRef<HTMLDivElement>(null);
  const logs = useMatchStore((state) => state.logs);

  useEffect(() => {
    logs.forEach((log) => {
      const element = document.createElement("p");
      element.textContent = log;
      ref.current?.appendChild(element);
    });
    ref.current?.scrollTo(0, ref.current.scrollHeight);
  }, [logs]);

  return (
    <ScrollArea ref={ref}>
      <section className="bg-red-50" />
    </ScrollArea>
  );
};
