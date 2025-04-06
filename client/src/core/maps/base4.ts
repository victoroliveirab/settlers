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

const OUTER_PADDING = 20;
const HEXAGON_SIZE = 64;

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
    const hexSize = HEXAGON_SIZE;
    const spacing = hexSize / 8;

    const [min, max] = this.getRectBounds(this.grid, hexSize, spacing, OUTER_PADDING);
    const svgWidth = max.x - min.x;
    const svgHeight = max.y - min.y;

    this.root.innerHTML = "";
    const svg = document.createElementNS(this.ns, "svg");
    svg.id = "map";
    svg.setAttribute("width", String(svgWidth));
    svg.setAttribute("height", String(svgHeight));
    svg.setAttribute("viewBox", `0 0 ${svgWidth} ${svgHeight}`);

    const hexagons = document.createElementNS(this.ns, "g");
    const tokens = document.createElementNS(this.ns, "g");

    const edges = document.createElementNS(this.ns, "g");
    edges.id = this.edgesGroupID;
    edges.addEventListener("click", (e) => {
      const element = e.target as SVGPolygonElement;
      if (element.dataset.disabled !== "false" || !element.dataset.id) return;
      this.eventHandlers.onEdgeClick(Number(element.dataset.id));
    });

    const vertices = document.createElementNS(this.ns, "g");
    vertices.id = this.verticesGroupID;
    vertices.addEventListener("click", (e) => {
      const element = e.target as SVGCircleElement;
      if (element.dataset.disabled !== "false" || !element.dataset.id) return;
      this.eventHandlers.onVertexClick(Number(element.dataset.id));
    });

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
      hexagon.id = String(tile.id);
      if (tile.resource !== "Desert") {
        const token = this.drawNumberToken(center, tile.token);
        tokens.append(token);
      }
      hexagons.appendChild(hexagon);

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
        edges.appendChild(element);
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
        vertices.appendChild(circle);
      });
    });

    svg.appendChild(hexagons);
    svg.appendChild(tokens);
    svg.appendChild(edges);
    svg.appendChild(vertices);

    this.root.appendChild(svg);
  }

  drawRobbers(tilesIDs: number[]) {
    const robbers = this.root.querySelectorAll("rect[data-robber='true']");
    robbers.forEach((robber) => {
      robber.remove();
    });

    tilesIDs.forEach((tileID) => {
      const robber = this.generateRobber();
      const hexagon = this.root.querySelector(`polygon[data-id="${tileID}"]`)!;
      const group = hexagon.parentElement!;
      group.appendChild(robber);
    });
  }
}
