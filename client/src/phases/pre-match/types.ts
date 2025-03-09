import type { SettlersCore } from "../../core/types";

export namespace SettlersWSServerPreMatch {
  export type IncomingMessages = {
    "room.connect.success": {
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
  };

  export type OutgoingMessages = {
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
