import type { HexagonDef, HexCoordinate, Point } from "./types";
import { resourceColors } from "../constants";

export type EventHandlers = {
  onEdgeClick: (edgeID: number) => void;
  onTileClick: (tileID: number) => void;
  onVertexClick: (vertexID: number) => void;
};

export default abstract class BaseMapRenderer {
  protected readonly ns = "http://www.w3.org/2000/svg";

  protected readonly edgesGroupID = "map-edges";
  protected readonly tilesGroupID = "map-tiles";
  protected readonly verticesGroupID = "map-vertices";

  protected hexSize!: number;
  protected outerPadding!: number;
  protected spacingProportion!: number;

  private darkenTilesTimeoutRef: NodeJS.Timeout | null = null;

  constructor(
    protected readonly root: HTMLElement,
    protected readonly grid: HexCoordinate[],
    protected readonly colorByPlayer: Record<
      SettlersCore.Player["name"],
      SettlersCore.Player["color"]
    >,
    protected readonly username: string,
    protected readonly eventHandlers: EventHandlers,
  ) {}

  protected hexCoordinateToPoint(coordinate: HexCoordinate, hexSize: number, spacing: number) {
    const spacingFactor = 1 + spacing / (Math.sqrt(3) * hexSize);

    const x =
      spacingFactor * hexSize * (Math.sqrt(3) * coordinate.q + (Math.sqrt(3) / 2) * coordinate.r);
    const y = spacingFactor * hexSize * ((3 / 2) * coordinate.r);

    return { x, y };
  }

  protected getHexagonPoints(center: Point, size: number) {
    const points: Point[] = [];
    for (let i = 0; i < 6; i++) {
      const angle = ((2 * Math.PI) / 6) * (i + 0.5); // +0.5 to make it pointy-top
      const x = center.x + size * Math.cos(angle);
      const y = center.y + size * Math.sin(angle);
      points.push({ x, y });
    }
    // FIXME: fix the generation to generate in the correct order instead of this
    return [points[5], ...points.slice(0, 5)] as HexagonDef;
  }

  protected getRectBounds(
    grid: HexCoordinate[],
    hexSize: number,
    spacing: number,
    outerPadding: number,
  ): [Point, Point] {
    let minX = Infinity,
      maxX = -Infinity,
      minY = Infinity,
      maxY = -Infinity;
    grid.forEach((coordinates) => {
      const pixel = this.hexCoordinateToPoint(coordinates, hexSize, spacing);
      minX = Math.min(minX, pixel.x - hexSize);
      maxX = Math.max(maxX, pixel.x + hexSize);
      minY = Math.min(minY, pixel.y - hexSize);
      maxY = Math.max(maxY, pixel.y + hexSize);
    });
    return [
      {
        x: minX - outerPadding,
        y: minY - outerPadding,
      },
      {
        x: maxX + outerPadding,
        y: maxY + outerPadding,
      },
    ];
  }

  protected drawHexagon(
    center: Point,
    size: number,
    id: number,
    tileType?: SettlersCore.TileType,
  ): SVGPolygonElement {
    const points = this.getHexagonPoints(center, size);
    const polygon = document.createElementNS(this.ns, "polygon");
    polygon.dataset.id = String(id);
    polygon.dataset.type = "tile";
    polygon.setAttribute("points", points.map(({ x, y }) => `${x},${y}`).join(" "));
    polygon.setAttribute("fill", tileType ? resourceColors[tileType] : "transparent");
    polygon.setAttribute("stroke", tileType ? "#333" : "transparent");
    polygon.setAttribute("stroke-width", "1");
    return polygon;
  }

  protected drawNumberToken(center: Point, number: number) {
    const numberTokenCircleRadius = Math.max(10, this.hexSize / 3);
    const textFontSize = 0.8 * numberTokenCircleRadius;
    const frequencyFontSize = 0.9 * textFontSize;

    const g = document.createElementNS(this.ns, "g");
    const circle = document.createElementNS(this.ns, "circle");
    circle.setAttribute("cx", String(center.x));
    circle.setAttribute("cy", String(center.y));
    circle.setAttribute("r", String(numberTokenCircleRadius));
    circle.setAttribute("fill", "white");
    circle.setAttribute("stroke", "black");
    circle.setAttribute("stroke-width", "0.1px");

    const text = document.createElementNS(this.ns, "text");
    text.setAttribute("x", String(center.x));
    text.setAttribute("y", String(center.y + 2));
    text.setAttribute("text-anchor", "middle");
    text.setAttribute("font-size", String(textFontSize));
    text.setAttribute("fill", number === 6 || number === 8 ? "red" : "black");
    text.textContent = String(number);

    const frequency = document.createElementNS(this.ns, "text");
    frequency.setAttribute("x", String(center.x));
    frequency.setAttribute("y", String(center.y + 8));
    frequency.setAttribute("text-anchor", "middle");
    frequency.setAttribute("font-size", String(frequencyFontSize));
    frequency.setAttribute("fill", number === 6 || number === 8 ? "red" : "black");
    let dots = "";
    if (number === 2 || number === 12) {
      dots = ".";
    } else if (number === 3 || number === 11) {
      dots = "..";
    } else if (number === 4 || number === 10) {
      dots = "...";
    } else if (number === 5 || number === 9) {
      dots = "....";
    } else if (number === 6 || number === 8) {
      dots = ".....";
    }
    frequency.textContent = dots;

    g.append(circle);
    g.append(text);
    g.append(frequency);
    return g;
  }

  private getPortArc(pointA: Point, pointB: Point, spacing: number) {
    const dx = pointB.x - pointA.x;
    const dy = pointB.y - pointA.y;
    const distance = Math.hypot(dx, dy);

    if (distance === 0 || spacing * 2 >= distance) {
      throw new Error("Spacing is too large or points are identical.");
    }

    const ux = dx / distance;
    const uy = dy / distance;

    const paddedA: Point = {
      x: pointA.x + ux * spacing,
      y: pointA.y + uy * spacing,
    };

    const paddedB: Point = {
      x: pointB.x - ux * spacing,
      y: pointB.y - uy * spacing,
    };

    const paddedDistance = Math.hypot(paddedB.x - paddedA.x, paddedB.y - paddedA.y);
    const radius = paddedDistance;

    const largeArcFlag = 0;
    const sweepFlag = 0;

    return `M ${paddedA.x} ${paddedA.y} A ${radius} ${radius} 0 ${largeArcFlag} ${sweepFlag} ${paddedB.x} ${paddedB.y}`;
  }

  private getIsoscelesApexPoint(pointA: Point, pointB: Point) {
    const { x: x1, y: y1 } = pointA;
    const { x: x2, y: y2 } = pointB;

    const mx = (x1 + x2) / 2;
    const my = (y1 + y2) / 2;

    const dx = x2 - x1;
    const dy = y2 - y1;

    const length = Math.hypot(dx, dy);

    const h = length / 2;

    const perpDx = -dy / length;
    const perpDy = dx / length;

    const cx = mx + perpDx * h;
    const cy = my + perpDy * h;

    return { x: cx, y: cy };
  }

  protected drawPort(
    pointA: Point,
    pointB: Point,
    spacing: number,
    portType: string,
    angle: number,
  ) {
    const arcPath = this.getPortArc(pointA, pointB, spacing);
    const middlePoint = this.getIsoscelesApexPoint(pointA, pointB);

    const g = document.createElementNS(this.ns, "g");

    const path = document.createElementNS(this.ns, "path");
    path.setAttribute("d", arcPath);
    path.setAttribute("stroke", "gray");
    path.setAttribute("fill", "none");
    path.setAttribute("stroke-width", "3");
    g.appendChild(path);

    const text = document.createElementNS(this.ns, "text");
    text.setAttribute("x", String(middlePoint.x));
    text.setAttribute("y", String(middlePoint.y));
    text.setAttribute("text-anchor", "middle");
    text.setAttribute("dominant-baseline", "middle");
    text.setAttribute("font-size", "10");
    text.textContent = portType;
    text.setAttribute("fill", "white");
    text.setAttribute("transform", `rotate(${angle}, ${middlePoint.x}, ${middlePoint.y})`);
    g.appendChild(text);

    return g;
  }

  protected getVirtualMiddlePoints(hexPoints: HexagonDef, spacing: number): HexagonDef {
    const points: Point[] = [];
    const spacing1_2 = spacing / 2;
    const r = (spacing * Math.sqrt(3)) / 6;
    // 0 degrees (North)
    const points5 = hexPoints[5];
    points.push({
      x: points5.x,
      y: points5.y - spacing1_2,
    });
    // 60 degrees (Northeast)
    const point0 = hexPoints[0];
    points.push({
      x: point0.x + spacing1_2,
      y: point0.y - r,
    });
    // 120 degrees (Southeast)
    const points1 = hexPoints[1];
    points.push({
      x: points1.x + spacing1_2,
      y: points1.y + r,
    });
    // 180 degrees (South)
    const points2 = hexPoints[2];
    points.push({
      x: points2.x,
      y: points2.y + spacing1_2,
    });
    // 240 degrees (Southwest)
    const points3 = hexPoints[3];
    points.push({
      x: points3.x - spacing1_2,
      y: points3.y + r,
    });
    // 300 degrees (Northwest)
    const points4 = hexPoints[4];
    points.push({
      x: points4.x - spacing1_2,
      y: points4.y - r,
    });
    return points as HexagonDef;
  }

  protected getEdgesCoordinatesAroundHexagon(
    center: Point,
    hexSize: number,
    spacing: number,
    edgePadding: number,
  ) {
    const edgesPoints: [Point, Point, Point, Point][] = [];

    const vertices: Point[] = [];
    for (let i = 0; i < 6; i++) {
      const angle = (Math.PI / 3) * i + Math.PI / 2;
      const x = center.x + hexSize * Math.cos(angle);
      const y = center.y + hexSize * Math.sin(angle);
      vertices.push({ x, y });
    }

    for (let i = 0; i < 6; i++) {
      const currVertex = vertices[i];
      const nextVertex = vertices[(i + 1) % 6];

      const edgeAngle = (Math.PI / 3) * i + Math.PI / 2 + Math.PI / 6;
      const dirX = Math.cos(edgeAngle);
      const dirY = Math.sin(edgeAngle);

      const perpAngle = edgeAngle - Math.PI / 2;
      const perpX = Math.cos(perpAngle);
      const perpY = Math.sin(perpAngle);

      const width = Math.sqrt(
        Math.pow(nextVertex.x - currVertex.x, 2) + Math.pow(nextVertex.y - currVertex.y, 2),
      );
      const halfWidth = width / 2;

      const midX = (currVertex.x + nextVertex.x) / 2;
      const midY = (currVertex.y + nextVertex.y) / 2;

      const rectPoints: [Point, Point, Point, Point] = [
        {
          x: midX - perpX * halfWidth + dirX * edgePadding,
          y: midY - perpY * halfWidth + dirY * edgePadding,
        },
        {
          x: midX + perpX * halfWidth + dirX * edgePadding,
          y: midY + perpY * halfWidth + dirY * edgePadding,
        },
        {
          x: midX + perpX * halfWidth + dirX * (spacing - edgePadding),
          y: midY + perpY * halfWidth + dirY * (spacing - edgePadding),
        },
        {
          x: midX - perpX * halfWidth + dirX * (spacing - edgePadding),
          y: midY - perpY * halfWidth + dirY * (spacing - edgePadding),
        },
      ];
      edgesPoints.push(rectPoints);
    }
    // FIXME: fix the generation to generate in the correct order instead of this
    return [...edgesPoints.slice(3), ...edgesPoints.slice(0, 3)];
  }

  // protected getPortBridgeCoordinatesByVertexIndex(vertexCenter: Point, vertexIndex: number) {
  //   const bridgeLength = this.hexSize / 2;
  // }

  protected drawVertex(
    center: Point,
    radius: number,
    id: number,
    enabled: boolean,
    color: string,
    owned: boolean = false,
  ) {
    const vertex = document.createElementNS(this.ns, "circle");
    vertex.dataset.type = "vertex";
    vertex.dataset.id = String(id);
    vertex.dataset.disabled = String(!enabled);
    vertex.dataset.owned = String(owned);
    vertex.setAttribute("cx", String(center.x));
    vertex.setAttribute("cy", String(center.y));
    vertex.setAttribute("r", String(radius));
    vertex.setAttribute("fill", color);
    return vertex;
  }

  protected drawEdge(
    points: [Point, Point, Point, Point],
    id: number,
    enabled: boolean,
    color: string,
    owned: boolean = false,
  ) {
    const coordinates = points.map(({ x, y }) => `${x},${y}`).join(" ");
    const edge = document.createElementNS(this.ns, "polygon");
    edge.dataset.type = "edge";
    edge.dataset.id = String(id);
    edge.dataset.disabled = String(!enabled);
    edge.dataset.owned = String(owned);
    edge.setAttribute("points", coordinates);
    edge.setAttribute("fill", color);
    return edge;
  }

  protected drawRobber(center: Point) {
    const side = this.hexSize / 3;
    const strokeWidth = side / 5;
    const x = center.x;
    const y = center.y - this.hexSize / 2 - 2;
    const height = (Math.sqrt(3) / 2) * side;
    const robber = document.createElementNS(this.ns, "polygon");
    robber.dataset.type = "robber";
    const points = [
      { x: x, y: y - (2 * height) / 3 },
      { x: x - side / 2, y: y + height / 3 },
      { x: x + side / 2, y: y + height / 3 },
    ];
    robber.setAttribute("points", points.map(({ x, y }) => `${x},${y}`).join(" "));
    robber.setAttribute("fill", "white");
    robber.setAttribute("stroke", "black");
    robber.setAttribute("stroke-width", String(strokeWidth));
    return robber;
  }

  updateEdges(availableEdges: number[], enabled: boolean, highlight: boolean) {
    const edgeGroup = this.root.querySelector<SVGGElement>(`#${this.edgesGroupID}`);
    if (!edgeGroup) {
      console.error("SVG Edge Group not found");
      return;
    }
    const edges = Array.from<SVGCircleElement>(this.root.querySelectorAll("[data-type='edge']"));
    if (!enabled) {
      edges.forEach((edge) => {
        edge.dataset.disabled = "true";
      });
      edgeGroup.classList.remove("pulse");
      return;
    }
    edges.forEach((edge) => {
      const edgeID = Number(edge.dataset.id);
      const enabled = availableEdges.includes(edgeID);
      const disabled = !enabled;
      edge.dataset.disabled = String(disabled);
    });
    if (highlight) {
      edgeGroup.classList.add("pulse");
    } else {
      edgeGroup.classList.remove("pulse");
    }
  }

  updateRoads(roads: SettlersCore.Roads) {
    for (const road of Object.values(roads)) {
      const edge = this.root.querySelector<SVGPolygonElement>(
        `[data-type='edge'][data-id='${road.id}']`,
      );
      if (!edge) {
        console.error(`edge#${road.id} not found`);
        continue;
      }
      edge.dataset.disabled = "true";
      edge.dataset.owned = "true";
      edge.setAttribute("fill", this.colorByPlayer[road.owner].background);
    }
  }

  updateVertices(
    availableSettlementVertices: number[],
    availableCityVertices: number[],
    enabled: boolean,
    highlight: boolean,
  ) {
    const vertexGroup = this.root.querySelector<SVGGElement>(`#${this.verticesGroupID}`);
    if (!vertexGroup) {
      console.error("SVG Vertex Group not found");
      return;
    }
    const vertices = Array.from<SVGCircleElement>(
      this.root.querySelectorAll("[data-type='vertex']"),
    );
    if (!enabled) {
      vertices.forEach((vertex) => {
        vertex.dataset.disabled = "true";
      });
      vertexGroup.classList.remove("pulse");
      return;
    }
    vertices.forEach((vertex) => {
      const vertexID = Number(vertex.dataset.id);
      const enabled =
        availableSettlementVertices.includes(vertexID) || availableCityVertices.includes(vertexID);
      const disabled = !enabled;
      vertex.dataset.disabled = String(disabled);
    });
    if (highlight) {
      vertexGroup.classList.add("pulse");
    } else {
      vertexGroup.classList.remove("pulse");
    }
  }

  updateSettlements(settlements: SettlersCore.Settlements) {
    for (const settlement of Object.values(settlements)) {
      const vertex = this.root.querySelector<SVGCircleElement>(
        `[data-type='vertex'][data-id='${settlement.id}']`,
      );
      if (!vertex) {
        console.error(`vertex#${settlement.id} not found`);
        continue;
      }
      vertex.dataset.disabled = "false";
      vertex.dataset.owned = "true";
      vertex.setAttribute("fill", this.colorByPlayer[settlement.owner].background);
    }
  }

  updateCities(cities: SettlersCore.Cities) {
    const spacing = this.hexSize * this.spacingProportion;
    for (const city of Object.values(cities)) {
      const vertex = this.root.querySelector<SVGCircleElement>(
        `[data-type='vertex'][data-id='${city.id}']`,
      );
      if (!vertex) {
        console.error(`vertex#${city.id} not found`);
        continue;
      }
      vertex.dataset.disabled = "false";
      vertex.dataset.owned = "true";
      vertex.setAttribute("fill", this.colorByPlayer[city.owner].background);
      vertex.setAttribute("r", String(spacing * 1.6));
      vertex.setAttribute("stroke-width", "2");
      vertex.setAttribute("stroke", this.colorByPlayer[city.owner].foreground);
    }
  }

  updateTiles(availableTiles: number[], enabled: boolean, highlight: boolean) {
    const tileGroup = this.root.querySelector<SVGPolygonElement>(`#${this.tilesGroupID}`);
    if (!tileGroup) {
      console.error("SVG Tile Group not found");
      return;
    }
    const tiles = Array.from<SVGPolygonElement>(this.root.querySelectorAll("[data-type='tile']"));
    if (!enabled) {
      tiles.forEach((tile) => {
        tile.dataset.disabled = "true";
      });
      tileGroup.classList.remove("pulse");
      return;
    }
    tiles.forEach((tile) => {
      const tileID = Number(tile.dataset.id);
      const enabled = availableTiles.includes(tileID);
      const disabled = !enabled;
      tile.dataset.disabled = String(disabled);
    });
    if (highlight) {
      tileGroup.classList.add("pulse");
    } else {
      tileGroup.classList.remove("pulse");
    }
  }

  updateRobbers(blockedTiles: number[]) {
    const robbers = Array.from(
      this.root.querySelectorAll<SVGPolygonElement>("[data-type='robber']"),
    );
    robbers.forEach((robber) => {
      robber.remove();
    });
    for (let i = 0; i < blockedTiles.length; ++i) {
      const tileID = blockedTiles[i];
      const tile = this.root.querySelector<SVGPolygonElement>(
        `[data-type='tile'][data-id='${tileID}']`,
      );
      if (!tile) {
        console.error(`tile${tileID} not found`);
        continue;
      }
      const points = Array.from(tile.points);
      let xSum = 0,
        ySum = 0;
      for (const point of points) {
        xSum += point.x;
        ySum += point.y;
      }
      const center = {
        x: xSum / points.length,
        y: ySum / points.length,
      };
      const robber = this.drawRobber(center);
      tile.after(robber);
    }
  }

  darkenTiles(tilesIDs: number[]) {
    if (tilesIDs.length === 0) return;
    const tileGroup = this.root.querySelector<SVGPolygonElement>(`#${this.tilesGroupID}`);
    if (!tileGroup) {
      console.error("SVG Tile Group not found");
      return;
    }
    if (tileGroup.classList.contains("pulse")) {
      console.warn("Already blinking tiles");
      return;
    }
    if (this.darkenTilesTimeoutRef) {
      clearTimeout(this.darkenTilesTimeoutRef);
    }
    const tiles = Array.from<SVGPolygonElement>(this.root.querySelectorAll("[data-type='tile']"));
    for (const tile of tiles) {
      if (!tilesIDs.includes(Number(tile.dataset.id))) continue;
      tile.classList.add("darken");
    }
    this.darkenTilesTimeoutRef = setTimeout(() => {
      for (const tile of tiles) {
        tile.classList.remove("darken");
      }
    }, 1000);
  }

  abstract render(map: SettlersCore.Map, ports: SettlersCore.Ports): void;
}
