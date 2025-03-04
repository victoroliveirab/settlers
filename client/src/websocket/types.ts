export namespace SettlersCore {
  export type Resource = "Brick" | "Ore" | "Grain" | "Sheep" | "Lumber";
  export type DevelopmentCard =
    | "Knight"
    | "Victory Point"
    | "Year of Plenty"
    | "Road Building"
    | "Monopoly";
  export type TileType = Resource | "Desert";
  export type Tile = {
    id: number;
    resource: TileType;
    token: number;
    edges: number[];
    vertices: number[];
    coordinates: { q: number; r: number; s: number };
  };
  export type Map = Tile[];
  export type Player = {
    color: string;
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
  export type Settlements = Record<Building["id"], Building>;
  export type Cities = Record<Building["id"], Building>;
  export type Roads = Record<Building["id"], Building>;
  export type Hand = Record<Resource, number>;
  export type DevHand = Record<DevelopmentCard, number>;
}

export namespace SettlersWSServer {
  export type IncomingMessages = {
    "room.join.success": {
      id: string;
      capacity: number;
      map: string;
      participants: SettlersCore.Participant[];
      owner: SettlersCore.Player["name"];
    };
    "room.new-update": {
      id: string;
      capacity: number;
      map: string;
      participants: SettlersCore.Participant[];
      owner: SettlersCore.Player["name"];
    };
    "game.start": {
      currentRoundPlayer: SettlersCore.Player["name"];
      logs: string[];
      map: SettlersCore.Map;
      players: SettlersCore.Player[];
      resourceCount: Record<SettlersCore.Player["name"], number>;
    };
    "setup.build-settlement": {
      vertices: number[];
    };
    "setup.settlement-build.success": {
      logs: string[];
      settlement: SettlersCore.Building;
    };
    "setup.build-road": {
      edges: number[];
    };
    "setup.road-build.success": {
      logs: string[];
      road: SettlersCore.Building;
    };
    "setup.player-round-changed": {
      currentRoundPlayer: SettlersCore.Player["name"];
    };
    "setup.end": {
      hands: Record<SettlersCore.Player["name"], SettlersCore.Hand>;
      logs: string[];
    };
    "game.player-round": {
      currentRoundPlayer: SettlersCore.Player["name"];
    };
    "game.dice-roll.success": {
      dices: [number, number];
      hand: SettlersCore.Hand;
      logs: string[];
      resourceCount: Record<SettlersCore.Player["name"], number>;
    };
    "game.discard-cards-request": {
      quantityByPlayers: Record<SettlersCore.Player["name"], number>;
    };
    "game.discard-cards.success": {
      hand: SettlersCore.Hand;
      resourceCount: Record<SettlersCore.Player["name"], number>;
    };
    "game.discarded-cards": {
      logs: string[];
      quantityByPlayers: Record<SettlersCore.Player["name"], number>;
    };
    hydrate: {
      state: {
        map: SettlersCore.Map;
        settlements: SettlersCore.Settlements;
        cities: SettlersCore.Cities;
        roads: SettlersCore.Roads;
        round: number;
        players: SettlersCore.Player[];
        currentRoundPlayer: SettlersCore.Player["name"];
        hand: SettlersCore.Hand;
        resourceCount: Record<SettlersCore.Player["name"], number>;
        dice: [number, number];
      };
    };
  };

  export type OutgoingMessages = {
    "room.join": {
      roomID: string;
    };
    "room.toggle-ready": {
      roomID: string;
      ready: boolean;
    };
    "room.start-game": {};
    "game.new-settlement": {
      vertex: number;
    };
    "setup.new-settlement": {
      vertex: number;
    };
    "game.new-road": {
      edge: number;
    };
    "setup.new-road": {
      edge: number;
    };
    "game.dice-roll": {};
    "game.end-round": {};
    "game.discard-cards": {
      resources: Record<SettlersCore.Resource, number>;
    };
  };

  export type IncomingMessage<T extends keyof IncomingMessages> = {
    type: T;
    payload: IncomingMessages[T];
  };
  export type OutgoingMessage<T extends keyof OutgoingMessages> = {
    type: T;
    payload: OutgoingMessages[T];
  };
}
