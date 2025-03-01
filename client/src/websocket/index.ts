import GameRenderer from "../renderer/game";
import PreGameRenderer from "../renderer/pre-game";
import GameState from "../state";
import { SettlersWSServer } from "./types";

function safeParse(text: string):
  | {
      [K in keyof SettlersWSServer.IncomingMessages]: SettlersWSServer.IncomingMessage<K>;
    }[keyof SettlersWSServer.IncomingMessages]
  | null {
  try {
    const parsed = JSON.parse(text);
    if (parsed && typeof parsed.type === "string") {
      return parsed;
    }
  } catch (err) {
    console.error("Error while safe parsing the following text blob:", text);
    return null;
  }
  return null;
}

export default class WebSocketConnection {
  private ws!: WebSocket;
  private roomID: string;

  constructor(
    url: string,
    private readonly stateManager: GameState,
  ) {
    this.roomID = window.location.pathname.split("/").at(-1)!;
    const ws = new WebSocket(url);
    ws.onopen = (e) => {
      this.ws = ws;
      console.log("websocket connection opened", e);
      ws.send(
        JSON.stringify({
          type: "room.join",
          payload: {
            roomID: this.roomID,
          },
        }),
      );
      this.stateManager.setService(this);
    };
    ws.onclose = (e) => {
      console.log("websocket connection closed", e);
    };
    ws.onerror = (e) => {
      console.error("websocket error", e);
    };
    ws.onmessage = this.onMessage.bind(this);
  }

  private onMessage(event: MessageEvent) {
    const message = safeParse(event.data);
    console.log(message);

    if (!message) return;

    switch (message.type) {
      case "room.join.success": {
        const { owner, participants } = message.payload;
        this.stateManager.setParticipants(participants);
        this.stateManager.setOwner(owner);
        break;
      }
      case "room.new-update": {
        const { owner, participants } = message.payload;
        this.stateManager.setParticipants(participants);
        this.stateManager.setOwner(owner);
        break;
      }
      case "game.start": {
        // TODO: get map name from payload
        const { map, players } = message.payload;
        console.log("SETTING INITIAL STATE");
        this.stateManager.setInitialState(map, players);
        break;
      }
      case "setup.build-settlement": {
        const { vertices } = message.payload;
        this.stateManager.enableVerticesToBuildSettlement(vertices, "setup");
        break;
      }
      case "hydrate": {
        const { map, players } = message.payload.state;
        this.stateManager.setInitialState(map, players);
        break;
      }
    }
    this.stateManager.repaintScreen();
  }

  onReadyChange(state: boolean) {
    this.sendMessage({
      type: "room.toggle-ready",
      payload: {
        ready: state,
        roomID: this.roomID,
      },
    });
  }

  onClickStart() {
    this.sendMessage({
      type: "room.start-game",
      payload: {},
    });
  }

  onSettlementPositionChose(phase: "game" | "setup", vertexID: number) {
    this.sendMessage({
      type: phase === "game" ? "game.new-settlement" : "setup.new-settlement",
      payload: {
        vertex: vertexID,
      },
    });
  }

  private sendMessage<T extends keyof SettlersWSServer.OutgoingMessages>(
    message: SettlersWSServer.OutgoingMessage<T>,
  ) {
    this.ws.send(JSON.stringify(message));
  }
}
