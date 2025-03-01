export namespace SettlersCore {
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
}

export namespace SettlersWSServer {
  export type IncomingMessages = {
    "room.join.success": {
      id: string;
      capacity: number;
      map: string;
      participants: SettlersCore.Participant[];
      owner: SettlersCore.Player["id"];
    };
    "room.new-update": {
      id: string;
      capacity: number;
      map: string;
      participants: SettlersCore.Participant[];
      owner: SettlersCore.Player["id"];
    };
  };

  export type OutgoingMessages = {
    "room.join": {
      roomID: string;
    };
    "room.toggle-ready": {
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
