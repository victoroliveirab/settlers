import { createElement } from "react";
import { cva, VariantProps } from "class-variance-authority";
import { type ClassValue } from "clsx";

import { cn } from "@/lib/utils";

import { emojis, resourcesOrder } from "@/core/constants";

import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "../ui/tooltip";

export type GameCardProps<T extends keyof React.JSX.IntrinsicElements = "div"> =
  React.ComponentProps<T> &
    VariantProps<typeof cardsVariants> & {
      as?: keyof React.JSX.IntrinsicElements;
      innerClassName?: ClassValue[];
      value: string;
    };

const cardsVariants = cva("h-full aspect-[3/4] rounded-md flex items-center justify-center", {
  variants: {
    size: {
      xs: "h-6 p-1",
      sm: "h-10",
      md: "h-14",
      lg: "h-20",
    },
    value: {
      default: "bg-neutral-300",
      Lumber: "bg-[green]",
      Brick: "bg-[maroon]",
      Sheep: "bg-[limegreen]",
      Grain: "bg-[orange]",
      Ore: "bg-[gray]",
    },
  },
  defaultVariants: {
    size: "md",
    value: "default",
  },
});

function isResource(type: string): type is SettlersCore.Resource {
  return (resourcesOrder as string[]).includes(type);
}

function GameCard<T extends keyof React.JSX.IntrinsicElements>({
  as = "div",
  className,
  innerClassName,
  size,
  value,
  ...props
}: GameCardProps<T>) {
  const emoji = isResource(value) ? emojis.resources[value] : emojis.devCards[value];
  return createElement(as, {
    className: cn(cardsVariants({ className, size, value })),
    children: (
      <div
        className={cn(
          "flex items-center justify-center",
          {
            "rounded-md bg-neutral-200 h-3/4 aspect-[3/4]": size !== "xs",
          },
          innerClassName,
        )}
      >
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger>
              <span
                className={cn("select-none", {
                  "text-xs": size === "xs",
                  "text-sm": size === "sm",
                  "text-md": size === "md",
                  "text-lg": size === "lg",
                })}
              >
                {emoji}
              </span>
            </TooltipTrigger>
            <TooltipContent>{value}</TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </div>
    ),
    ...props,
  });
}

export { GameCard };
