import { Button } from "@/components/ui/button";
import { ChevronDown, ChevronUp } from "lucide-react";

interface ICollapsibleTogglerProps {
  numberOfTrades: number;
  open: boolean;
}

export const CollapsibleToggler = ({ open, numberOfTrades }: ICollapsibleTogglerProps) => {
  const text = numberOfTrades === 1 ? "trade" : "trades";
  const icon = open ? <ChevronDown /> : <ChevronUp />;

  return (
    <Button className="w-1/2" variant="secondary">
      {`${numberOfTrades} ${text} opened`} <span>{icon}</span>
    </Button>
  );
};
