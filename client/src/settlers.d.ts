type SVGDataset = {
  disabled: "false" | "true";
  id: string;
};

declare namespace Settlers {
  interface SVGEdge extends Omit<SVGRectElement, "dataset"> {
    dataset: SVGDataset;
  }
  interface SVGVertice extends Omit<SVGCircleElement, "dataset"> {
    dataset: SVGDataset;
  }
  interface SVGHexagon extends Omit<SVGPolygonElement, "dataset"> {
    dataset: SVGDataset;
  }
}
