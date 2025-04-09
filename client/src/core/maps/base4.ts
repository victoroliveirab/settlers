import BaseMapRenderer, { EventHandlers } from "./_base";
import type { HexCoordinate } from "./types";

function generateHexagonGrid() {
  const grid: HexCoordinate[] = [];

  const rows: [number, number[]][] = [
    // [q-value, array of r-values]
    [-2, [0, 1, 2]],
    [-1, [-1, 0, 1, 2]],
    [0, [-2, -1, 0, 1, 2]],
    [1, [-2, -1, 0, 1]],
    [2, [-2, -1, 0]],
  ];

  for (const [q, qValues] of rows) {
    for (const r of qValues) {
      const s = -q - r;
      grid.push({ q, r, s });
    }
  }

  return grid;
}

export default class Base4MapRenderer extends BaseMapRenderer {
  private readonly grid: HexCoordinate[];
  constructor(
    root: HTMLElement,
    colorByPlayer: Record<SettlersCore.Player["name"], SettlersCore.Player["color"]>,
    username: string,
    eventHandlers: EventHandlers,
  ) {
    super(root, colorByPlayer, username, eventHandlers);
    this.grid = generateHexagonGrid();
  }

  render(map: SettlersCore.Map, ports: { vertices: [number, number]; type: string }[]): void {
    const hexSize = this.hexSize;
    const spacing = hexSize * this.spacingProportion;

    const [min, max] = this.getRectBounds(this.grid, hexSize, spacing, this.outerPadding);
    const svgWidth = max.x - min.x;
    const svgHeight = max.y - min.y;

    this.root.innerHTML = "";
    const svg = document.createElementNS(this.ns, "svg");
    svg.id = "map";
    svg.setAttribute("width", String(svgWidth));
    svg.setAttribute("height", String(svgHeight));
    svg.setAttribute("viewBox", `0 0 ${svgWidth} ${svgHeight}`);
    this.root.appendChild(svg);

    const tilesLayer = document.createElementNS(this.ns, "g");
    tilesLayer.id = this.tilesGroupID;
    tilesLayer.addEventListener("click", (e) => {
      const element = e.target as SVGPolygonElement;
      if (element.dataset.disabled !== "false" || !element.dataset.id) return;
      this.eventHandlers.onTileClick(Number(element.dataset.id));
    });

    const numberTokensLayer = document.createElementNS(this.ns, "g");

    const edgesLayer = document.createElementNS(this.ns, "g");
    edgesLayer.id = this.edgesGroupID;
    edgesLayer.addEventListener("click", (e) => {
      const element = e.target as SVGPolygonElement;
      if (element.dataset.disabled !== "false" || !element.dataset.id) return;
      this.eventHandlers.onEdgeClick(Number(element.dataset.id));
    });

    const verticesLayer = document.createElementNS(this.ns, "g");
    verticesLayer.id = this.verticesGroupID;
    verticesLayer.addEventListener("click", (e) => {
      const element = e.target as SVGCircleElement;
      if (element.dataset.disabled !== "false" || !element.dataset.id) return;
      this.eventHandlers.onVertexClick(Number(element.dataset.id));
    });

    svg.appendChild(tilesLayer);
    svg.appendChild(numberTokensLayer);
    svg.appendChild(edgesLayer);
    svg.appendChild(verticesLayer);

    const createdEdges = new Set<number>();
    const createdVertices = new Set<number>();

    this.grid.forEach((hex, index) => {
      const tile = map[index];
      const point = this.hexCoordinateToPoint(hex, hexSize, spacing);
      const center = {
        x: point.x - min.x,
        y: point.y - min.y,
      };
      const hexagon = this.drawHexagon(center, hexSize, tile.id, tile.resource);
      if (tile.resource !== "Desert") {
        const token = this.drawNumberToken(center, tile.token);
        numberTokensLayer.append(token);
      }
      tilesLayer.appendChild(hexagon);

      const edgesCoordinates = this.getEdgesCoordinatesAroundHexagon(center, hexSize, spacing, 0);
      edgesCoordinates.forEach((edge, i) => {
        const edgeID = tile.edges[i];
        if (createdEdges.has(edgeID)) return;
        const element = this.drawEdge(
          edge,
          edgeID,
          false,
          this.colorByPlayer[this.username].background,
        );
        edgesLayer.appendChild(element);
        createdEdges.add(edgeID);
      });

      const verticesCoordinates = this.getHexagonPoints(center, hexSize);
      const centers = this.getVirtualMiddlePoints(verticesCoordinates, spacing);
      centers.forEach((circleCenter, i) => {
        const vertexID = tile.vertices[i];
        if (createdVertices.has(vertexID)) return;
        const circle = this.drawVertex(
          circleCenter,
          spacing * 1.3,
          vertexID,
          false,
          this.colorByPlayer[this.username].background,
        );
        createdVertices.add(vertexID);
        verticesLayer.appendChild(circle);
      });

      if (tile.blocked) {
        const robber = this.drawRobber(center);
        tilesLayer.appendChild(robber);
      }
    });
  }
}
