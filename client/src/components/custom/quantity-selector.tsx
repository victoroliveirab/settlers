import { Minus, Plus } from "lucide-react";

import { Button } from "../ui/button";

interface IQuantitySelectorProps {
  disabled?: boolean;
  onValueChange: (value: number) => void;
  max?: number;
  min?: number;
  step?: number;
  value: number;
}

export const QuantitySelector = ({
  disabled = false,
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
        disabled={disabled || value - step < min}
      >
        <Minus />
      </Button>
      <span>{value}</span>
      <Button
        size="xs"
        variant="ghost"
        onClick={() => onValueChange(Math.min(max, value + step))}
        disabled={disabled || value + step > max}
      >
        <Plus />
      </Button>
    </div>
  );
};
