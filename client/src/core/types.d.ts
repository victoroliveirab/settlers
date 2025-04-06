declare namespace SettlersCore {
  export type Resource = "Brick" | "Ore" | "Grain" | "Sheep" | "Lumber";
  export type DevelopmentCard =
    | "Knight"
    | "Victory Point"
    | "Year of Plenty"
    | "Road Building"
    | "Monopoly";
  export type TileType = Resource | "Desert";
  export type Tile = {
    blocked: boolean;
    id: number;
    resource: TileType;
    token: number;
    edges: number[];
    vertices: number[];
    coordinates: { q: number; r: number; s: number };
  };
  export type Port = {
    type: string;
    vertices: [number, number];
  };
  export type Map = Tile[];
  export type Ports = Port[];
  export type Player = {
    color: {
      background: string;
      foreground: string;
    };
    name: string;
  };
  export type Participant = {
    bot: boolean;
    player: Player | null;
    ready: boolean;
  };
  export type Building = {
    id: number;
    owner: string;
  };
  export type ResourceCollection = Record<Resource, number>;
  export type Settlements = Record<Building["id"], Building>;
  export type Cities = Record<Building["id"], Building>;
  export type Roads = Record<Building["id"], Building>;
  export type Hand = ResourceCollection;
  export type DevHand = Record<DevelopmentCard, number>;
}
