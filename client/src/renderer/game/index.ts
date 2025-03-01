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

  makeVerticesClickable(verticesIDs: number[], cb: (vertexID: number) => void) {
    const surface = this.root.querySelector("#map") as SVGElement;
    surface.classList.add("pulse-settlements");
    const vertices = Array.from(
      surface.querySelectorAll<Settlers.SVGVertice>("[data-type='vertice']"),
    ).filter((vertex) => verticesIDs.includes(+vertex.dataset.id));
    let ref = function onVerticeClick(e: Event) {
      const vertex = e.target as Settlers.SVGVertice;
      const vertexID = Number(vertex.dataset.id);
      cb(vertexID);
      vertices.forEach((vertex) => {
        vertex.removeEventListener("click", ref);
        vertex.dataset.disabled = "true";
      });
      surface.classList.remove("pulse-settlements");
    };
    console.log("ENABLING vertices:", vertices);
    vertices.forEach((vertex) => {
      vertex.addEventListener("click", ref);
      vertex.dataset.disabled = "false";
    });
  }
}
