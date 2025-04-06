import { Minus, Plus } from "lucide-react";

import { Button } from "../ui/button";

interface IQuantitySelectorProps {
  onValueChange: (value: number) => void;
  max?: number;
  min?: number;
  step?: number;
  value: number;
}

export const QuantitySelector = ({
  max = Infinity,
  min = 0,
  onValueChange,
  step = 1,
  value,
}: IQuantitySelectorProps) => {
  return (
    <div className="flex items-center gap-1">
      <Button
        size="xs"
        variant="ghost"
        onClick={() => onValueChange(Math.max(min, value - step))}
        // disabled={value <= min}
      >
        <Minus />
      </Button>
      <span>{value}</span>
      <Button
        size="xs"
        variant="ghost"
        onClick={() => onValueChange(Math.min(max, value + step))}
        disabled={value >= max}
      >
        <Plus />
      </Button>
    </div>
  );
};
