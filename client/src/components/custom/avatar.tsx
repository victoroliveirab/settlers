import { cn } from "@/lib/utils";

export type AvatarProps = React.ComponentProps<"div">;

function Avatar({ className, ...props }: AvatarProps) {
  return (
    <div
      data-slot="avatar"
      className={cn("w-14 h-14 border border-solid border-transparent rounded-full", className)}
      {...props}
    />
  );
}

export { Avatar };
