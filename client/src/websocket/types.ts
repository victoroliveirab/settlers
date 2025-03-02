export namespace SettlersCore {
  export type Resource = "Brick" | "Ore" | "Desert" | "Grain" | "Sheep" | "Lumber";
  export type Tile = {
    id: number;
    resource: Resource;
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
    "game.player-round": {
      currentRoundPlayer: SettlersCore.Player["name"];
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
