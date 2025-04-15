import BaseMapRenderer, { EventHandlers } from "./_base";
import type { HexCoordinate, Point } from "./types";

// function getTrianglePolygonPointsFromMids(point1: Point, point2: Point, legWidth = 4) {
//   const { x: x1, y: y1 } = point1;
//   const { x: x2, y: y2 } = point2;
//   // 1. Get the apex point using the original function
//   const apex = getIsoscelesApexPoint(point1, point2);
//
//   // 2. Vector from A to B
//   const dx = x2 - x1;
//   const dy = y2 - y1;
//   const length = Math.hypot(dx, dy);
//
//   // 3. Perpendicular unit vector to AB
//   const perpX = -dy / length;
//   const perpY = dx / length;
//
//   // 4. Half-thickness for base offset
//   const offsetX = perpX * (legWidth / 2);
//   const offsetY = perpY * (legWidth / 2);
//
//   // 5. Base edge points from the midpoints
//   const base1 = { x: x1 + offsetX, y: y1 + offsetY };
//   const base2 = { x: x2 + offsetX, y: y2 + offsetY };
//   const base3 = { x: x2 - offsetX, y: y2 - offsetY };
//   const base4 = { x: x1 - offsetX, y: y1 - offsetY };
//
//   // 6. Return polygon points in a clockwise or counterclockwise order
//   return [base1, base2, apex, base4];
// }

// function getEdgePolygonPoints(midX, midY, apexX, apexY, width = 4) {
//   // Vector from midpoint to apex
//   const dx = apexX - midX;
//   const dy = apexY - midY;
//   const length = Math.hypot(dx, dy);
//
//   // Unit perpendicular vector
//   const perpX = -dy / length;
//   const perpY = dx / length;
//
//   // Offset for width
//   const offsetX = perpX * (width / 2);
//   const offsetY = perpY * (width / 2);
//
//   // Create the 4 corner points of the edge polygon
//   return [
//     { x: midX + offsetX, y: midY + offsetY },
//     { x: apexX + offsetX, y: apexY + offsetY },
//     { x: apexX - offsetX, y: apexY - offsetY },
//     { x: midX - offsetX, y: midY - offsetY },
//   ];
// }

function getIsoscelesApexPoint(point1: Point, point2: Point) {
  const { x: x1, y: y1 } = point1;
  const { x: x2, y: y2 } = point2;
  // Calculate the midpoint of AB
  const mx = (x1 + x2) / 2;
  const my = (y1 + y2) / 2;

  // Compute the vector from A to B
  const dx = x2 - x1;
  const dy = y2 - y1;

  // Length of AB
  const length = Math.hypot(dx, dy);

  // Height = 0.5 * length
  const h = length / 2;

  // Perpendicular direction vector (normalized)
  const perpDx = -dy / length;
  const perpDy = dx / length;

  // Coordinates of the apex point C
  const cx = mx - perpDx * h;
  const cy = my - perpDy * h;

  return { x: cx, y: cy };
}

// function generateArc(pointA: Point, pointB: Point): string {
//   const dx = pointB.x - pointA.x;
//   const dy = pointB.y - pointA.y;
//   const distance = Math.hypot(dx, dy);
//   const radius = distance; // Since arc diameter is 2 * distance between points
//
//   // Large arc flag = 0 (use smaller arc), sweep flag = 1 (clockwise)
//   const largeArcFlag = 0;
//   const sweepFlag = 1;
//
//   // SVG path string (from pointA to pointB, with given radius)
//   const arcPath = `M ${pointA.x} ${pointA.y} A ${radius} ${radius} 0 ${largeArcFlag} ${sweepFlag} ${pointB.x} ${pointB.y}`;
//
//   return arcPath;
// }

function generateArc(pointA: Point, pointB: Point, spacing: number = 0): string {
  const dx = pointB.x - pointA.x;
  const dy = pointB.y - pointA.y;
  const distance = Math.hypot(dx, dy);

  if (distance === 0 || spacing * 2 >= distance) {
    throw new Error("Spacing is too large or points are identical.");
  }

  // Unit direction vector from A to B
  const ux = dx / distance;
  const uy = dy / distance;

  // Shrink both points toward each other
  const paddedA: Point = {
    x: pointA.x + ux * spacing,
    y: pointA.y + uy * spacing,
  };

  const paddedB: Point = {
    x: pointB.x - ux * spacing,
    y: pointB.y - uy * spacing,
  };

  const paddedDistance = Math.hypot(paddedB.x - paddedA.x, paddedB.y - paddedA.y);
  const radius = paddedDistance; // Arc diameter is twice this, so radius = distance

  const largeArcFlag = 0;
  const sweepFlag = 1;

  return `M ${paddedA.x} ${paddedA.y} A ${radius} ${radius} 0 ${largeArcFlag} ${sweepFlag} ${paddedB.x} ${paddedB.y}`;
}

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
  constructor(
    root: HTMLElement,
    colorByPlayer: Record<SettlersCore.Player["name"], SettlersCore.Player["color"]>,
    username: string,
    eventHandlers: EventHandlers,
  ) {
    super(root, generateHexagonGrid(), colorByPlayer, username, eventHandlers);
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

    const portsLayer = document.createElementNS(this.ns, "g");

    svg.appendChild(tilesLayer);
    svg.appendChild(numberTokensLayer);
    svg.appendChild(edgesLayer);
    svg.appendChild(verticesLayer);
    svg.appendChild(portsLayer);

    const createdEdges = new Set<number>();
    const createdVertices = new Set<number>();
    const centerByVertex: Record<number, Point> = {};

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
        centerByVertex[vertexID] = circleCenter;
        createdVertices.add(vertexID);
        verticesLayer.appendChild(circle);
      });

      if (tile.blocked) {
        const robber = this.drawRobber(center);
        tilesLayer.appendChild(robber);
      }
    });

    ports.forEach((port) => {
      const [vertexID1, vertexID2] = port.vertices;
      const p1 = centerByVertex[vertexID1];
      const p2 = centerByVertex[vertexID2];
      let tileIndex = 0;
      let pointA: Point | undefined;
      let pointB: Point | undefined;
      let angle: number = 0;
      while (tileIndex < map.length) {
        const tile = map[tileIndex];
        if (tile.vertices[0] === vertexID1 && tile.vertices[5] === vertexID2) {
          pointA = p1;
          pointB = p2;
          angle = -30;
          break;
        }
        for (let i = 0; i < 5; ++i) {
          if (tile.vertices[i] === vertexID1 && tile.vertices[i + 1] === vertexID2) {
            pointA = p2;
            pointB = p1;

            if (i === 2) {
              angle = -30;
            } else if (i === 0 || i === 3) {
              angle = 30;
            }
            break;
          } else if (tile.vertices[i] === vertexID2 && tile.vertices[i + 1] === vertexID1) {
            pointA = p1;
            pointB = p2;

            if (i === 2) {
              angle = -30;
            } else if (i === 3) {
              angle = 30;
            }
            break;
          }
        }
        if (pointA && pointB) break;
        ++tileIndex;
      }
      if (!pointA || !pointB) {
        console.error("No points found", { port, map });
        return;
      }
      const portGroup = this.drawPort(pointA, pointB, spacing, port.type, angle);
      portsLayer.appendChild(portGroup);
    });
  }
}
