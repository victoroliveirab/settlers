import BaseMapRenderer, { type EventHandlers } from "./_base";

import Base4MapRenderer from "./base4";

export function mapRendererFactory(
  mapName: string,
  root: HTMLElement,
  colorByPlayer: Record<SettlersCore.Player["name"], SettlersCore.Player["color"]>,
  username: string,
  eventHandlers: EventHandlers,
): BaseMapRenderer | null {
  console.log({ mapName });
  switch (mapName) {
    case "base4": {
      return new Base4MapRenderer(root, colorByPlayer, username, eventHandlers);
    }
    default: {
      console.warn(`Unsupported map: ${mapName || "<empty>"}`);
      return null;
    }
  }
}
