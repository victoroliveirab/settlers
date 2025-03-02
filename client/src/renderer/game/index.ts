import type { SettlersCore } from "../../websocket/types";
import BaseMapRenderer from "../maps/base";
import Base4MapRenderer from "../maps/base4";

export default class GameRenderer {
  private mapRenderer: BaseMapRenderer;
  constructor(
    private readonly root: HTMLElement,
    private readonly mapName: string,
  ) {
    this.root.style.display = "";
    const svg = this.root.querySelector<SVGElement>("#map")!;
    svg.setAttribute("width", "100%");
    svg.setAttribute("height", "100%");
    this.mapRenderer = new Base4MapRenderer(svg, 10, 10, 1);
  }

  drawMap(map: SettlersCore.Map) {
    this.mapRenderer.draw(map);
  }

  drawPlayers(players: SettlersCore.Player[], currentRoundPlayer: string) {
    const playersContainer = this.root.querySelector<HTMLDivElement>("#players")!;
    playersContainer.innerHTML = "";
    players.forEach((player) => {
      const div = document.createElement("div");
      div.style.background = player.color;
      const playerNameElement = document.createElement("h2");
      playerNameElement.textContent = player.name;
      div.appendChild(playerNameElement);

      const infoElement = document.createElement("ul");

      const numberOfCardsElement = document.createElement("li");
      numberOfCardsElement.textContent = `#R: ${0}`;
      const numberOfDevCardsElement = document.createElement("li");
      numberOfDevCardsElement.textContent = `#D: ${0}`;
      const longestRoadElement = document.createElement("li");
      longestRoadElement.textContent = `LG: ${0}`;
      const knightsElement = document.createElement("li");
      knightsElement.textContent = `#K: ${0}`;
      const points = document.createElement("li");
      points.textContent = `#P: ${0}`;

      infoElement.appendChild(numberOfCardsElement);
      infoElement.appendChild(numberOfDevCardsElement);
      infoElement.appendChild(longestRoadElement);
      infoElement.appendChild(knightsElement);
      infoElement.appendChild(points);
      div.appendChild(infoElement);

      if (player.name === currentRoundPlayer) {
        div.dataset.current = "true";
      }

      playersContainer.appendChild(div);
    });
  }

  drawDices(dices: [number, number], onClick?: () => void) {
    dices.forEach((dice, index) => {
      const selector = `#dice${index + 1}`;
      const element = this.root.querySelector<HTMLDivElement>(selector)!;
      element.textContent = String(dice);
    });
    if (onClick) {
      this.root.querySelector("#dice")?.addEventListener("click", onClick, { once: true });
    }
  }

  drawHud(hand: SettlersCore.Hand) {}

  drawSettlement(settlement: SettlersCore.Building, color: string) {
    const { id } = settlement;
    const spot = this.root.querySelector<SVGCircleElement>(`circle[data-id="${id}"]`);
    if (!spot) {
      console.warn("vertex not found:", id);
      return;
    }

    spot.style.opacity = "1";
    spot.style.fill = color;
    spot.dataset.disabled = "true";
  }

  drawRoad(road: SettlersCore.Building, color: string) {
    const { id } = road;
    const spot = this.root.querySelector<SVGCircleElement>(`rect[data-id="${id}"]`);
    if (!spot) {
      console.warn("edge not found:", id);
      return;
    }

    spot.style.opacity = "1";
    spot.style.fill = color;
    spot.dataset.disabled = "true";
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
    vertices.forEach((vertex) => {
      vertex.addEventListener("click", ref);
      vertex.dataset.disabled = "false";
    });
  }

  makeEdgesClickable(edgesIDs: number[], cb: (edgeID: number) => void) {
    const surface = this.root.querySelector("#map") as SVGElement;
    surface.classList.add("pulse-edges");
    const edges = Array.from(
      surface.querySelectorAll<Settlers.SVGEdge>("[data-type='edge']"),
    ).filter((edge) => edgesIDs.includes(+edge.dataset.id));
    let ref = function onVerticeClick(e: Event) {
      const vertice = e.target as Settlers.SVGVertice;
      const verticeID = Number(vertice.dataset.id);
      cb(verticeID);
      edges.forEach((edge) => {
        edge.removeEventListener("click", ref);
        edge.dataset.disabled = "true";
      });
      surface.classList.remove("pulse-edges");
    };
    edges.forEach((edge) => {
      edge.addEventListener("click", ref);
      edge.dataset.disabled = "false";
    });
  }

  renderNewLog(log: string) {
    const entry = document.createElement("p");
    entry.textContent = log;
    this.root.querySelector("#logs")?.append(entry);
  }
}
