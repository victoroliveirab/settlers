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
    id: number;
    roomID: string;
    username: string;
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
      owner: SettlersCore.Player["username"];
    };
    "room.new-update": {
      id: string;
      capacity: number;
      map: string;
      participants: SettlersCore.Participant[];
      owner: SettlersCore.Player["username"];
    };
    "game.start": {
      logs: string[];
      map: SettlersCore.Map;
      players: SettlersCore.Player[];
    };
    hydrate: {
      state: {
        map: SettlersCore.Map;
        settlements: SettlersCore.Settlements;
        cities: SettlersCore.Cities;
        roads: SettlersCore.Roads;
        round: number;
        players: SettlersCore.Player[];
        currentRoundPlayer: SettlersCore.Player["id"];
        hand: SettlersCore.Hand;
        resourceCount: Record<SettlersCore.Player["id"], number>;
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
