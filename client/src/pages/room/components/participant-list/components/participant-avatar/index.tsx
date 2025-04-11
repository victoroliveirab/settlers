import clsx from "clsx";

import { Avatar, type AvatarProps } from "@/components/custom/avatar";

type ParticipantAvatarProps = React.PropsWithChildren<
  {
    bot?: boolean;
    color?: string;
    empty?: boolean;
    owner?: boolean;
  } & Omit<AvatarProps, "children">
>;

export function ParticipantAvatar({
  bot = false,
  children,
  color = "transparent",
  empty = false,
  owner = false,
  ...props
}: ParticipantAvatarProps) {
  return (
    <Avatar
      className={clsx("relative", {
        "border-dashed border-black": empty,
      })}
      style={{ background: empty ? "transparent" : color }}
      {...props}
    >
      {!empty && (
        <span className="absolute top-1/2 left-1/2 -translate-1/2">
          {owner ? "👑" : bot ? "🤖" : "👤"}
        </span>
      )}
      {children}
    </Avatar>
  );
}
