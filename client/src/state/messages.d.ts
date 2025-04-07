export namespace SettlersRoom {
  export type IncomingMessages = {
    "room.connect.success": {
      room: SettlersServer.Room;
      params: SettlersServer.RoomParam[];
    };
    "room.new-update": {
      room: SettlersServer.Room;
      params: SettlersServer.RoomParam[];
    };
    "room.toggle-ready.success": {
      room: SettlersServer.Room;
      params: SettlersServer.RoomParam[];
    };
    "room.update-param.success": {
      room: SettlersServer.Room;
      params: SettlersServer.RoomParam[];
    };
    "room.start-game.success": {
      logs: string[];
      map: SettlersCore.Map;
      mapName: string;
      players: SettlersCore.Player[];
      resourceCount: Record<SettlersCore.Player["name"], number>;
    };
  };

  export type OutgoingMessages = {
    "room.player-change-color": {
      color: string;
    };
    "room.update-param": {
      key: string;
      value: number;
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

export namespace SettlersMatch {
  export type IncomingMessages = {
    "setup.update-round-player": {
      player: string;
    };
    "setup.update-vertices": {
      availableSettlementVertices: number[];
      availableCityVertices: number[];
      enabled: boolean;
      highlight: boolean;
    };
    "setup.update-edges": {
      availableEdges: number[];
      enabled: boolean;
      highlight: boolean;
    };
    "setup.update-map": {
      blockedTiles: number[];
      cities: SettlersCore.Cities;
      roads: SettlersCore.Roads;
      settlements: SettlersCore.Settlements;
    };
    "setup.update-logs": string[];

    "match.update-round-player": {
      player: string;
    };
    "match.update-vertices": {
      availableSettlementVertices: number[];
      availableCityVertices: number[];
      enabled: boolean;
      highlight: boolean;
    };
    "match.update-edges": {
      availableEdges: number[];
      enabled: boolean;
      highlight: boolean;
    };
    "match.update-map": {
      blockedTiles: number[];
      cities: SettlersCore.Cities;
      roads: SettlersCore.Roads;
      settlements: SettlersCore.Settlements;
    };
    "match.update-resource-count": {
      resourceCount: Record<SettlersCore.Player["name"], number>;
    };
    "match.update-logs": string[];
    "match.update-dice": {
      dice: [number, number];
      enabled: boolean;
    };
    "match.update-hand": {
      hand: SettlersCore.Hand;
    };
    "match.update-pass": {
      enabled: boolean;
    };
    "match.update-trade": {
      enabled: boolean;
    };
    "match.update-robber-movement": {
      availableTiles: number[];
      enabled: boolean;
    };
    "match.update-discard-phase": {
      discardAmounts: Record<SettlersCore.Player["name"], number>;
      enabled: boolean;
    };
    "match.update-trade-offers": {
      offers: {
        counters: number[];
        finalized: boolean;
        id: number;
        offer: Record<SettlersCore.Resource, number>;
        opponents: Record<
          SettlersCore.Player["name"],
          { status: "Open" | "Accepted" | "Declined"; blocked: boolean }
        >;
        parent: number;
        player: SettlersCore.Player["name"];
        request: Record<SettlersCore.Resource, number>;
        status: "Open" | "Closed";
        timestamp: number;
      }[];
    };
  };

  export type OutgoingMessages = {
    "match.dice-roll": {};
    "match.vertex-click": {
      vertex: number;
    };
    "match.edge-click": {
      edge: number;
    };
    "match.tile-click": {
      tile: number;
    };
    "match.pass-click": {};
    "match.discard-cards": {
      resources: SettlersCore.ResourceCollection;
    };
    "match.create-trade-offer": {
      given: SettlersCore.ResourceCollection;
      requested: SettlersCore.ResourceCollection;
    };
    "match.create-counter-trade-offer": {
      given: SettlersCore.ResourceCollection;
      requested: SettlersCore.ResourceCollection;
      tradeID: number;
    };
    "match.accept-trade-offer": {
      tradeID: number;
    };
    "match.reject-trade-offer": {
      tradeID: number;
    };
    "match.cancel-trade-offer": {
      tradeID: number;
    };
    "match.finalize-trade-offer": {
      accepter: string;
      tradeID: number;
    };
    "match.end-round": {};
  };

  type SingleIncomingMessage<T extends keyof IncomingMessages = keyof IncomingMessages> = {
    type: T;
    payload: IncomingMessages[T];
  };

  export type MatchSetupHydrateMessage = {
    type: "setup.hydrate";
    payload: {
      edgeUpdate: SingleIncomingMessage<"setup.update-edges">;
      map: SettlersCore.Map;
      mapName: string;
      mapUpdate: SingleIncomingMessage<"setup.update-map">;
      players: SettlersCore.Player[];
      ports: SettlersCore.Ports;
      resourceCount: Record<SettlersCore.Player["name"], number>;
      roundPlayerUpdate: SingleIncomingMessage<"setup.update-round-player">;
      vertexUpdate: SingleIncomingMessage<"setup.update-vertices">;
    };
  };

  export type MatchHydrateMessage = {
    type: "match.hydrate";
    payload: {
      diceUpdate: SingleIncomingMessage<"match.update-dice">;
      edgeUpdate: SingleIncomingMessage<"match.update-edges">;
      handUpdate: SingleIncomingMessage<"match.update-hand">;
      map: SettlersCore.Map;
      mapName: string;
      mapUpdate: SingleIncomingMessage<"match.update-map">;
      passActionState: SingleIncomingMessage<"match.update-pass">;
      players: SettlersCore.Player[];
      ports: SettlersCore.Ports;
      resourceCount: Record<SettlersCore.Player["name"], number>;
      robberMovementUpdate: SingleIncomingMessage<"match.update-robber-movement">;
      roundPlayerUpdate: SingleIncomingMessage<"match.update-round-player">;
      tradeActionState: SingleIncomingMessage<"match.update-trade">;
      tradeOffersUpdate: SingleIncomingMessage<"match.update-trade-offers">;
      vertexUpdate: SingleIncomingMessage<"match.update-vertices">;
    };
  };

  type BulkUpdateMessage<T extends keyof IncomingMessages = keyof IncomingMessages> = {
    type: T;
    payload: IncomingMessages[T];
  };
  export type IncomingMessage<T extends keyof IncomingMessages = keyof IncomingMessages> =
    | SingleIncomingMessage<T>
    | { type: "setup.bulk-update"; payload: BulkUpdateMessage[] }
    | { type: "match.bulk-update"; payload: BulkUpdateMessage[] };
  export type OutgoingMessage<T extends keyof OutgoingMessages> = {
    type: T;
    payload: OutgoingMessages[T];
  };
}

export type SettlersIncomingMessage =
  | {
      [K in keyof SettlersRoom.IncomingMessages]: SettlersRoom.IncomingMessage<K>;
    }[keyof SettlersRoom.IncomingMessages]
  | {
      [K in keyof SettlersMatch.IncomingMessages]: SettlersMatch.IncomingMessage<K>;
    }[keyof SettlersMatch.IncomingMessages]
  | SettlersMatch.MatchSetupHydrateMessage
  | SettlersMatch.MatchHydrateMessage;

export type SettlersOutgoingMessages =
  | {
      [K in keyof SettlersRoom.OutgoingMessages]: SettlersRoom.OutgoingMessage<K>;
    }[keyof SettlersRoom.OutgoingMessages]
  | {
      [K in keyof SettlersMatch.OutgoingMessages]: SettlersMatch.OutgoingMessage<K>;
    }[keyof SettlersMatch.OutgoingMessages];
