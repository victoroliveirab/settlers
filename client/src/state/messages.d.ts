export namespace SettlersRoom {
  export type IncomingMessages = {
    "room.connect.success": {
      minMaxPlayers: [number, number];
      room: SettlersServer.Room;
      params: SettlersServer.RoomParam[];
    };
    "room.new-update": {
      room: SettlersServer.Room;
      params: SettlersServer.RoomParam[];
    };
    "room.update-capacity.success": {
      minMaxPlayers: [number, number];
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
      ports: SettlersCore.Ports;
      resourceCount: Record<SettlersCore.Player["name"], number>;
      roomStatus: string;
    };
  };

  export type OutgoingMessages = {
    "room.player-change-color": {
      color: string;
    };
    "room.update-capacity": {
      capacity: number;
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
      deadline: string;
      player: string;
      serverNow: string;
      subDeadline: string | null;
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
    "setup.update-ports": {
      ports: SettlersCore.PortType[];
    };
    "setup.update-longest-road-size": {
      longestRoadSizeByPlayer: Record<SettlersCore.Player["name"], number>;
    };
    "setup.update-points": {
      points: Record<SettlersCore.Player["name"], number>;
    };
    "setup.update-logs": string[];

    "match.update-round-player": {
      deadline: string;
      player: string;
      serverNow: string;
      subDeadline: string | null;
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
    "match.update-ports": {
      ports: SettlersCore.PortType[];
    };
    "match.update-resource-count": {
      resourceCount: Record<SettlersCore.Player["name"], number>;
    };
    "match.update-dev-hand-count": {
      devHandCount: Record<SettlersCore.Player["name"], number>;
    };
    "match.update-logs": string[];
    "match.update-dice": {
      dice: [number, number];
      enabled: boolean;
    };
    "match.update-hand": {
      hand: SettlersCore.Hand;
    };
    "match.update-dev-hand": {
      devHand: SettlersCore.DevHand;
    };
    "match.update-dev-hand-permissions": {
      devHandPermissions: Record<SettlersCore.DevelopmentCard, boolean>;
    };
    "match.update-pass": {
      enabled: boolean;
    };
    "match.update-trade": {
      enabled: boolean;
    };
    "match.update-buy-dev-card": {
      enabled: boolean;
    };
    "match.update-robber-movement": {
      availableTiles: number[];
      enabled: boolean;
      highlight: boolean;
    };
    "match.update-pick-robbed": {
      enabled: boolean;
      options: SettlersCore.Player["name"][] | null;
    };
    "match.update-discard-phase": {
      discardAmounts: Record<SettlersCore.Player["name"], number>;
      enabled: boolean;
    };
    "match.update-trade-offers": {
      offers: {
        creator: SettlersCore.Player["name"];
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
        requester: SettlersCore.Player["name"];
        responses: Record<
          SettlersCore.Player["name"],
          { status: "Open" | "Accepted" | "Declined"; blocked: boolean }
        >;
        status: "Open" | "Closed";
        timestamp: number;
      }[];
    };
    "match.update-points": {
      points: Record<SettlersCore.Player["name"], number>;
    };
    "match.update-longest-road-size": {
      longestRoadSizeByPlayer: Record<SettlersCore.Player["name"], number>;
    };
    "match.update-knight-usage": {
      knightUsesByPlayer: Record<SettlersCore.Player["name"], number>;
    };
    "match.update-year-of-plenty": {
      enabled: boolean;
    };
    "match.update-monopoly": {
      enabled: boolean;
    };
    "match.statistics.success": {
      statistics: {
        diceStatsByPlayer: Record<SettlersCore.Player["name"], Record<number, number>>;
        generalDiceStats: Record<number, number>;
        longestRoadEvolution: Record<SettlersCore.Player["name"], number[]>;
        numberOfRobberiesByPlayer: Record<SettlersCore.Player["name"], number>;
        pointsEvolution: Record<SettlersCore.Player["name"], number[]> | null;
      };
    };

    "over.data": {
      points: Record<SettlersCore.Player["name"], number>;
      roomStatus: string;
      roundsPlayed: number;
      statistics: {
        diceStatsByPlayer: Record<SettlersCore.Player["name"], Record<number, number>>;
        generalDiceStats: Record<number, number>;
        longestRoadEvolution: Record<SettlersCore.Player["name"], number[]>;
        numberOfRobberiesByPlayer: Record<SettlersCore.Player["name"], number>;
        pointsEvolution: Record<SettlersCore.Player["name"], number[]> | null;
      };
      startDatetime: string;
      endDatetime: string;
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
    "match.dev-card-click": {
      kind: SettlersCore.DevelopmentCard;
    };
    "match.discard-cards": {
      resources: SettlersCore.ResourceCollection;
    };
    "match.make-bank-trade": {
      given: SettlersCore.ResourceCollection;
      requested: SettlersCore.ResourceCollection;
    };
    "match.make-general-port-trade": {
      given: SettlersCore.ResourceCollection;
      requested: SettlersCore.ResourceCollection;
    };
    "match.make-resource-port-trade": {
      given: SettlersCore.ResourceCollection;
      requested: SettlersCore.ResourceCollection;
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
    "match.rob-player": {
      player: string;
    };
    "match.buy-dev-card": {};
    "match.year-of-plenty": {
      resource1: SettlersCore.Resource;
      resource2: SettlersCore.Resource;
    };
    "match.monopoly": {
      resource: SettlersCore.Resource;
    };
    "match.end-round": {};
    "match.statistics": {};
  };

  type SingleIncomingMessage<T extends keyof IncomingMessages = keyof IncomingMessages> = {
    type: T;
    payload: IncomingMessages[T];
  };

  export type MatchSetupHydrateMessage = {
    type: "setup.hydrate";
    payload: {
      devHandCount: Record<SettlersCore.Player["name"], number>;
      edgeUpdate: SingleIncomingMessage<"setup.update-edges">;
      map: SettlersCore.Map;
      mapName: string;
      mapUpdate: SingleIncomingMessage<"setup.update-map">;
      players: SettlersCore.Player[];
      ports: SettlersCore.Ports;
      resourceCount: Record<SettlersCore.Player["name"], number>;
      roundPlayerUpdate: SingleIncomingMessage<"setup.update-round-player">;
      roomStatus: string;
      vertexUpdate: SingleIncomingMessage<"setup.update-vertices">;
    };
  };

  export type MatchHydrateMessage = {
    type: "match.hydrate";
    payload: {
      buyDevCardUpdate: SingleIncomingMessage<"match.update-buy-dev-card">;
      devHandCount: Record<SettlersCore.Player["name"], number>;
      devHandUpdate: SingleIncomingMessage<"match.update-dev-hand">;
      devHandPermissionsUpdate: SingleIncomingMessage<"match.update-dev-hand-permissions">;
      diceUpdate: SingleIncomingMessage<"match.update-dice">;
      discardUpdate: SingleIncomingMessage<"match.update-discard-phase">;
      edgeUpdate: SingleIncomingMessage<"match.update-edges">;
      handUpdate: SingleIncomingMessage<"match.update-hand">;
      knightsUsageUpdate: SingleIncomingMessage<"match.update-knight-usage">;
      longestRoadUpdate: SingleIncomingMessage<"match.update-longest-road-size">;
      map: SettlersCore.Map;
      mapName: string;
      mapUpdate: SingleIncomingMessage<"match.update-map">;
      passActionState: SingleIncomingMessage<"match.update-pass">;
      players: SettlersCore.Player[];
      pointsUpdate: SingleIncomingMessage<"match.update-points">;
      ports: SettlersCore.Ports;
      portsUpdate: SingleIncomingMessage<"match.update-ports">;
      resourceCount: Record<SettlersCore.Player["name"], number>;
      robbablePlayersUpdate: SingleIncomingMessage<"match.update-pick-robbed">;
      robberMovementUpdate: SingleIncomingMessage<"match.update-robber-movement">;
      roomStatus: string;
      roundPlayerUpdate: SingleIncomingMessage<"match.update-round-player">;
      tradeActionState: SingleIncomingMessage<"match.update-trade">;
      tradeOffersUpdate: SingleIncomingMessage<"match.update-trade-offers">;
      vertexUpdate: SingleIncomingMessage<"match.update-vertices">;
      yearOfPlentyUpdate: SingleIncomingMessage<"match.update-year-of-plenty">;
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
