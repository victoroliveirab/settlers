import { SettlersCore } from "../../core/types";

export namespace SettlersWSServerMatch {
  export type IncomingMessages = {
    // Setup Phase
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
    "setup.hydrate": {
      state: {
        map: SettlersCore.Map;
        settlements: SettlersCore.Settlements;
        cities: SettlersCore.Cities;
        roads: SettlersCore.Roads;
        players: SettlersCore.Player[];
        currentRoundPlayer: SettlersCore.Player["name"];
      };
    };

    // Actual Match
    "game.player-round": {
      currentRoundPlayer: SettlersCore.Player["name"];
      roundType: number;
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
      resourceCount: Record<SettlersCore.Player["name"], number>;
      quantityByPlayers: Record<SettlersCore.Player["name"], number>;
    };
    "game.move-robber-request": {
      availableTiles: SettlersCore.Tile["id"][];
    };
  };

  export type OutgoingMessages = {
    // Setup Phase
    "setup.new-settlement": {
      vertex: number;
    };
    "setup.new-road": {
      edge: number;
    };

    // Match Phase
    "game.dice-roll": {};
    "game.move-robber": {
      tile: number;
    };
    "game.end-round": {};
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
