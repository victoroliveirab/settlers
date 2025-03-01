import type { SettlersCore } from "../../websocket/types";
import BaseMapRenderer from "../maps/base";
import Base4MapRenderer from "../maps/base4";

export default class GameRenderer {
  private mapRenderer: BaseMapRenderer;
  constructor(
    private readonly root: HTMLElement,
    private readonly mapName: string,
  ) {
    this.root.innerHTML = "";
    const svg = document.createElementNS("http://www.w3.org/2000/svg", "svg");
    svg.id = "map";
    svg.setAttribute("width", "720px");
    svg.setAttribute("height", "800px");
    svg.setAttribute("viewBox", "-60 -60 120 120");
    svg.setAttribute("version", "1.1");
    svg.setAttribute("xmlns", "http://www.w3.org/2000/svg");
    this.root.appendChild(svg);
    this.mapRenderer = new Base4MapRenderer(svg, 10, 10, 1);
  }

  drawMap(map: SettlersCore.Map) {
    this.mapRenderer.draw(map);
  }
}
