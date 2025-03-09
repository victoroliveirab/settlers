import { SettlersCore } from "../../websocket/types";
import BaseMapRenderer from "./_base";

export default class Base4MapRenderer extends BaseMapRenderer {
  constructor(
    protected readonly root: SVGElement,
    protected readonly width: number,
    protected readonly height: number,
    protected readonly spacing: number,
  ) {
    super(root, width, height, spacing);
  }

  draw(map: SettlersCore.Map): void {
    this.root.innerHTML = "";
    const polygonPoints = this.hexagonVerticesCoordinates.map(({ x, y }) => `${x},${y}`).join(" ");
    const vertices = new Set<number>();
    const edges = new Set<number>();
    map.forEach((element) => {
      const g = document.createElementNS(this.ns, "g");

      const { q, r, s } = element.coordinates;
      const x = (this.pixelMatrix[0][0] * q + this.pixelMatrix[0][1] * r) * this.width;
      const y = (this.pixelMatrix[1][0] * q + this.pixelMatrix[1][1] * r) * this.height;
      const transform = [`${x * this.spacing}px`, `${y * this.spacing}px`];
      g.style.transform = `translate(${transform.join(",")})`;
      g.dataset.q = String(q);
      g.dataset.r = String(r);
      g.dataset.s = String(s);

      const hexagon = document.createElementNS(this.ns, "polygon");
      hexagon.setAttribute("points", polygonPoints);
      hexagon.setAttribute("fill", this.colorByResource[element.resource]);
      hexagon.setAttribute("stroke", "#deb887");
      hexagon.setAttribute("stroke-width", "2");
      hexagon.dataset.type = "tile";
      hexagon.dataset.id = String(element.id);
      hexagon.dataset.disabled = "true";
      g.append(hexagon);

      if (element.resource !== "Desert") {
        const { circle, text, frequency } = this.generateNumberToken(element.token);
        g.append(circle);
        g.append(text);
        g.append(frequency);
      }

      const edgesCandidates = this.generateEdges();
      edgesCandidates.forEach((edge, edgeIndex) => {
        const edgeId = element.edges[edgeIndex];
        if (!edges.has(edgeId)) {
          edges.add(edgeId);
          edge.dataset.id = String(edgeId);
          g.append(edge);
        }
      });

      const verticesCandidates = this.generateVertices();
      verticesCandidates.forEach((vertice, verticeIndex) => {
        const verticeId = element.vertices[verticeIndex];
        if (!vertices.has(verticeId)) {
          vertices.add(verticeId);
          vertice.dataset.id = String(verticeId);
          g.append(vertice);
        }
      });

      this.root.prepend(g);
    });
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

  private generateEdges() {
    const edges: SVGRectElement[] = [];
    this.pathCoordinates.forEach(([x, y, w, h, r]) => {
      const rect = document.createElementNS(this.ns, "rect");
      rect.setAttribute("x", String(x));
      rect.setAttribute("y", String(y));
      rect.setAttribute("width", String(w));
      rect.setAttribute("height", String(h));
      rect.setAttribute("fill", "aqua");
      rect.style.transform = `rotate(${r}deg)`;
      rect.classList.add("edge-spot");
      rect.dataset.type = "edge";
      rect.dataset.disabled = "true";
      edges.push(rect);
    });
    return edges;
  }

  private generateVertices() {
    const circles: SVGCircleElement[] = [];
    this.verticesPoints.forEach((vertice) => {
      const circle = document.createElementNS(this.ns, "circle");
      // Perhaps these could css set with selectors
      circle.setAttribute("cx", String(vertice[0]));
      circle.setAttribute("cy", String(vertice[1]));
      circle.setAttribute("r", "2");
      circle.setAttribute("fill", "red");
      circle.setAttribute("stroke", "black");
      circle.setAttribute("stroke-width", "0.1px");
      circle.classList.add("settlement-spot");
      circle.dataset.type = "vertice";
      circle.dataset.disabled = "true";
      circles.push(circle);
    });
    return circles;
  }
}
