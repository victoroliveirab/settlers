import { cn } from "@/lib/utils";

export type AvatarProps = React.ComponentProps<"div"> & {
  background?: string;
  borderStyle?: React.CSSProperties["borderStyle"];
  foreground?: string;
  withBorder?: boolean;
};

function Avatar({
  background,
  borderStyle,
  className,
  foreground,
  withBorder = false,
  ...props
}: AvatarProps) {
  return (
    <div
      data-slot="avatar"
      className={cn("w-14 h-14 border border-solid rounded-full", className)}
      style={{
        background,
        borderColor: withBorder ? foreground : "transparent",
        borderStyle: borderStyle || "unset",
        ...props.style,
      }}
      {...props}
    />
  );
}

export { Avatar };
