import { SettlersCore } from "../../websocket/types";
import { Point } from "./types";

const tenSq3_2 = (10 * Math.sqrt(3)) / 2;

const pathCoordinates = Object.freeze([
  [0 - tenSq3_2 / 2, -10 + 0.25, 10 - 2.25 / 2, 2, 30],
  [10 - 2.25 - 0.125, -5 + 0.5, 2, 10 - 2.25 / 2, 0],
  [0 - tenSq3_2 / 2, 10 - 2.25, 10 - 2.25 / 2, 2, -30],
  [0 - tenSq3_2 / 2, 10 - 2.25, 10 - 2.25 / 2, 2, 30],
  [-10 + 0.5 - 0.125, -5 + 0.5, 2, 10 - 2.25 / 2, 0],
  [0 - tenSq3_2 / 2, -10 + 0.25, 10 - 2.25 / 2, 2, -30],
]);

const pixelMatrix = Object.freeze([
  [Math.sqrt(3), Math.sqrt(3) / 2],
  [0, 3 / 2],
]);

const verticesPoints = Object.freeze([
  [0, -10],
  [tenSq3_2, -5],
  [tenSq3_2, 5],
  [0, 10],
  [-tenSq3_2, 5],
  [-tenSq3_2, -5],
]);

const colorByResource = Object.freeze({
  Brick: "#D2691E",
  Lumber: "#228B22",
  Grain: "#FFD700",
  Ore: "#A9A9A9",
  Sheep: "#98FB98",
  Desert: "#878878",
});

export default abstract class BaseMapRenderer {
  protected readonly ns = "http://www.w3.org/2000/svg";
  protected readonly hexagonVerticesCoordinates: Point[] = [];
  constructor(
    protected readonly root: SVGElement,
    protected readonly width: number,
    protected readonly height: number,
    protected readonly spacing: number,
  ) {
    for (let i = 0; i < 6; ++i) {
      this.hexagonVerticesCoordinates.push({
        x: width * Math.cos((2 * Math.PI * i) / 6 + Math.PI / 6),
        y: width * Math.sin((2 * Math.PI * i) / 6 + Math.PI / 6),
      });
    }
  }

  protected get pathCoordinates() {
    return pathCoordinates;
  }

  protected get pixelMatrix() {
    return pixelMatrix;
  }

  protected get verticesPoints() {
    return verticesPoints;
  }

  protected get colorByResource() {
    return colorByResource;
  }

  abstract draw(map: SettlersCore.Map): void;

  protected generateNumberToken(number: number) {
    const circle = document.createElementNS(this.ns, "circle");
    circle.setAttribute("cx", "0");
    circle.setAttribute("cy", "0");
    circle.setAttribute("r", "3");
    circle.setAttribute("fill", "white");
    circle.setAttribute("stroke", "black");
    circle.setAttribute("stroke-width", "0.1px");

    const text = document.createElementNS(this.ns, "text");
    text.setAttribute("x", "0");
    text.setAttribute("y", "0.5");
    text.setAttribute("text-anchor", "middle");
    text.setAttribute("font-size", "3.5px");
    text.setAttribute("fill", number === 6 || number === 8 ? "red" : "black");
    text.textContent = String(number);

    const frequency = document.createElementNS(this.ns, "text");
    frequency.setAttribute("x", "0");
    frequency.setAttribute("y", "1.5");
    frequency.setAttribute("text-anchor", "middle");
    frequency.setAttribute("font-size", "3.5px");
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

    return {
      circle,
      frequency,
      text,
    };
  }
}
