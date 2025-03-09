import BaseMapRenderer from "./_base";

import Base4MapRenderer from "./base4";

export function mapRendererFactory(
  mapName: string,
  root: SVGElement,
  width: number,
  height: number,
  spacing: number,
): BaseMapRenderer {
  switch (mapName) {
    case "base4": {
      return new Base4MapRenderer(root, width, height, spacing);
    }
    default: {
      throw new Error(`Unsupported map: ${mapName}`);
    }
  }
}
