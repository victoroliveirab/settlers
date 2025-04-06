export type Point = {
  x: number;
  y: number;
};

export type HexCoordinate = {
  q: number;
  r: number;
  s: number;
};

export type HexagonDef = [Point, Point, Point, Point, Point, Point];
export type RectangleDef = [Point, Point, Point, Point];
