import { useState } from "react";
import { Info } from "lucide-react";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Slider } from "@/components/ui/slider";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";

import type { RoomParam } from "@/state/room";

export function CardParameter({ param }: { param: RoomParam }) {
  const [value, setValue] = useState(param.value);
  const minValue = param.values[0];
  const maxValue = param.values[param.values.length - 1];
  const step = Math.round((maxValue - minValue) / param.values.length);

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-sm">{param.label}</CardTitle>
        <CardDescription className="truncate">
          <span className="flex items-center gap-1">
            <span className="truncate">{param.description}</span>
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger>
                  <Info className="shrink-0" size={16} />
                </TooltipTrigger>
                <TooltipContent>
                  <p>{param.description}</p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          </span>
        </CardDescription>
      </CardHeader>
      <CardContent className="flex-1 flex items-end">
        <div className="w-full flex flex-col gap-1">
          <div className="w-full flex items-center gap-2">
            <span className="text-muted-foreground text-xs">{minValue}</span>
            <Slider
              className="flex-1"
              defaultValue={[value]}
              step={step}
              max={maxValue}
              min={minValue}
              value={[value]}
              onValueChange={([newValue]) => setValue(newValue)}
            />
            <span className="text-muted-foreground text-xs">{maxValue}</span>
          </div>
          <p className="text-xs text-center">{value}</p>
        </div>
      </CardContent>
    </Card>
  );
}
