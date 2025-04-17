import { createElement } from "react";

import { cn } from "@/lib/utils";

import styles from "./styles.module.css";

interface IDieProps {
  active?: boolean;
  as?: keyof React.JSX.IntrinsicElements;
  className?: string;
  value: number;
}

export const Die = ({ active = false, as = "div", className, value }: IDieProps) => {
  return createElement(as, {
    className: cn(
      "flex-1 aspect-square bg-neutral-200 max-w-24 relative rounded-md",
      styles.die,
      {
        [styles.animate]: active,
        [styles["die-0"]]: value === 0,
        [styles["die-1"]]: value === 1,
        [styles["die-2"]]: value === 2,
        [styles["die-3"]]: value === 3,
        [styles["die-4"]]: value === 4,
        [styles["die-5"]]: value === 5,
        [styles["die-6"]]: value === 6,
      },
      className,
    ),
    children: (
      <>
        <span className="absolute -translate-1/2 h-1/6 w-1/6 rounded-full bg-neutral-900" />
        <span className="absolute -translate-1/2 h-1/6 w-1/6 rounded-full bg-neutral-900" />
        <span className="absolute -translate-1/2 h-1/6 w-1/6 rounded-full bg-neutral-900" />
        <span className="absolute -translate-1/2 h-1/6 w-1/6 rounded-full bg-neutral-900" />
        <span className="absolute -translate-1/2 h-1/6 w-1/6 rounded-full bg-neutral-900" />
        <span className="absolute -translate-1/2 h-1/6 w-1/6 rounded-full bg-neutral-900" />
      </>
    ),
  });
};
