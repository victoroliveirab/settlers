import type { SettlersCore } from "../../websocket/types";
import BaseMapRenderer from "../maps/base";
import Base4MapRenderer from "../maps/base4";

type Player = {
  color: string;
  devHandCount: number;
  isCurrentRound: boolean;
  knights: number;
  longestRoad: number;
  name: string;
  points: number;
  quantityToDiscard: number;
  resourceCount: number;
};

const resourcesOrder: SettlersCore.Resource[] = ["Lumber", "Brick", "Sheep", "Grain", "Ore"];

const resourceEmojis = Object.freeze({
  Lumber: "🌲",
  Brick: "🧱",
  Sheep: "🐑",
  Grain: "🌾",
  Ore: "⛰️",
});

const developmentEmojis = Object.freeze({
  Knight: "⚔️",
  "Victory Point": "🎖️",
  "Road Building": "🛤️",
  "Year of Plenty": "🎁",
  Monopoly: "🎩",
});

const noop = () => {};

export default class GameRenderer {
  private mapRenderer: BaseMapRenderer;
  private diceEventHandler: () => void = noop;
  private passEventHandler: () => void = noop;
  private tradeEventHandler: () => void = noop;

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

  drawRobbers(tilesIDs: number[]) {
    this.mapRenderer.drawRobbers(tilesIDs);
  }

  drawPlayers(players: Player[]) {
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
      numberOfCardsElement.textContent = `#R: ${player.resourceCount}`;
      const numberOfDevCardsElement = document.createElement("li");
      numberOfDevCardsElement.textContent = `#D: ${player.devHandCount}`;
      const longestRoadElement = document.createElement("li");
      longestRoadElement.textContent = `LG: ${player.longestRoad}`;
      const knightsElement = document.createElement("li");
      knightsElement.textContent = `#K: ${player.knights}`;
      const points = document.createElement("li");
      points.textContent = `#P: ${player.points}`;

      infoElement.appendChild(numberOfCardsElement);
      infoElement.appendChild(numberOfDevCardsElement);
      infoElement.appendChild(longestRoadElement);
      infoElement.appendChild(knightsElement);
      infoElement.appendChild(points);

      if (player.quantityToDiscard > 0) {
        const discarding = document.createElement("li");
        discarding.textContent = "❌";
        infoElement.appendChild(discarding);
      }

      div.appendChild(infoElement);

      if (player.isCurrentRound) {
        div.dataset.current = "true";
      }

      playersContainer.appendChild(div);
    });
  }

  drawDices(dices: [number, number]) {
    const element = this.root.querySelector("#dice");
    if (!element) return;
    element.removeEventListener("click", this.diceEventHandler);
    dices.forEach((dice, index) => {
      const selector = `#dice${index + 1}`;
      const element = this.root.querySelector<HTMLDivElement>(selector)!;
      element.textContent = String(dice);
    });
  }

  attachClickHandlerToDice(onClick: () => void) {
    const element = this.root.querySelector("#dice");
    if (!element) return;
    element.removeEventListener("click", this.diceEventHandler);
    this.diceEventHandler = () => {
      element.classList.remove("pulse");
      onClick();
    };
    element.addEventListener("click", this.diceEventHandler, { once: true });
    element.classList.add("pulse");
  }

  drawHand(hand: SettlersCore.Hand) {
    const resourcesElement = this.root.querySelector("#resources")!;
    resourcesElement.innerHTML = "";
    resourcesOrder.forEach((resource) => {
      if (hand[resource] > 0) {
        for (let i = 0; i < hand[resource]; ++i) {
          const element = document.createElement("li");
          element.dataset.type = resource;
          const text = document.createElement("span");
          text.textContent = resourceEmojis[resource];
          element.appendChild(text);
          resourcesElement.appendChild(element);
        }
      }
    });
  }

  drawDevHand(devHand: SettlersCore.DevHand) {
    const devElement = this.root.querySelector("#dev")!;
    devElement.innerHTML = "";
    const devTypes: SettlersCore.DevelopmentCard[] = [
      "Knight",
      "Year of Plenty",
      "Road Building",
      "Monopoly",
      "Victory Point",
    ];
    devTypes.forEach((type) => {
      if (devHand[type] > 0) {
        for (let i = 0; i < devHand[type]; ++i) {
          const element = document.createElement("li");
          element.dataset.type = type;
          const text = document.createElement("span");
          text.textContent = developmentEmojis[type];
          element.appendChild(text);
          devElement.appendChild(element);
        }
      }
    });
  }

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

  updatePassButton(onClick?: () => void) {
    const button = this.root.querySelector<HTMLButtonElement>("#action-pass");
    if (!button) {
      console.warn("pass button not found");
      return;
    }
    if (onClick) {
      this.passEventHandler = () => {
        button.disabled = true;
        onClick();
      };
      button.disabled = false;
      button.addEventListener("click", this.passEventHandler, { once: true });
    } else {
      this.passEventHandler = noop;
      button.disabled = true;
      button.removeEventListener("click", this.passEventHandler);
    }
  }

  renderDiscardModal(
    hand: SettlersCore.Hand,
    quantityToDiscard: number,
    onSubmit: (selectedCards: SettlersCore.Resource[]) => void,
  ) {
    const container = this.root.querySelector<HTMLDivElement>("#discard");
    if (!container) {
      console.warn("discard dom node not found");
      return;
    }
    if (container.style.display === "flex") return; // Already rendered
    container.style.display = "flex";
    const subtitle = container.querySelector<HTMLHeadingElement>("h4")!;
    subtitle.textContent = `Discard ${quantityToDiscard} cards`;

    let numberOfSelectedCards = 0;
    const counter = container.querySelector<HTMLHeadingElement>("h5")!;
    counter.textContent = `0/${quantityToDiscard}`;

    const discardButton = container.querySelector<HTMLButtonElement>("#action-discard")!;
    discardButton.disabled = true;

    const list = container.querySelector<HTMLUListElement>("#discard-card-list")!;
    list.innerHTML = "";

    discardButton.addEventListener(
      "click",
      () => {
        const selectedCards = list.querySelectorAll<HTMLLIElement>('li[data-selected="true"]');
        if (selectedCards.length !== quantityToDiscard) return;

        onSubmit(Array.from(selectedCards).map((el) => el.dataset.type as SettlersCore.Resource));
      },
      { once: true },
    );

    resourcesOrder.forEach((resource) => {
      if (hand[resource] === 0) return;
      for (let i = 0; i < hand[resource]; ++i) {
        const card = document.createElement("li");
        card.dataset.type = resource;
        card.dataset.selected = "false";
        const text = document.createElement("span");
        text.textContent = resourceEmojis[resource];
        card.appendChild(text);
        card.addEventListener("click", () => {
          const selected = card.dataset.selected;
          if (selected === "false") {
            card.dataset.selected = "true";
          } else {
            card.dataset.selected = "false";
          }

          numberOfSelectedCards = list.querySelectorAll('li[data-selected="true"]').length;
          counter.textContent = `${numberOfSelectedCards}/${quantityToDiscard}`;

          if (numberOfSelectedCards === quantityToDiscard) {
            discardButton.disabled = false;
          } else {
            discardButton.disabled = true;
          }
        });
        list.appendChild(card);
      }
    });
  }

  hideDiscardModal() {
    const container = this.root.querySelector<HTMLDivElement>("#discard");
    if (!container) {
      console.warn("discard dom node not found");
      return;
    }

    container.style.display = "";

    const list = container.querySelector<HTMLUListElement>("#discard-card-list")!;
    list.innerHTML = "";
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

  makeTilesClickable(tilesIDs: number[], cb: (tileID: number) => void) {
    const surface = this.root.querySelector("#map") as SVGElement;
    surface.classList.add("pulse-tiles");
    const tiles = Array.from(
      surface.querySelectorAll<Settlers.SVGHexagon>("[data-type='tile']"),
    ).filter((tile) => tilesIDs.includes(+tile.dataset.id));
    let ref = function onTileClick(e: Event) {
      const tile = e.target as Settlers.SVGHexagon;
      const tileID = Number(tile.dataset.id);
      cb(tileID);
      tiles.forEach((tile) => {
        tile.removeEventListener("click", ref);
        tile.dataset.disabled = "true";
      });
      surface.classList.remove("pulse-tiles");
    };
    tiles.forEach((tile) => {
      tile.addEventListener("click", ref);
      tile.dataset.disabled = "false";
    });
  }

  renderNewLog(log: string) {
    const entry = document.createElement("p");
    entry.textContent = log;
    this.root.querySelector("#logs")?.append(entry);
  }
}
