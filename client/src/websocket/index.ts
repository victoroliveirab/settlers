import GameRenderer from "../renderer/game";
import PreGameRenderer from "../renderer/pre-game";
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
  private gameRenderer: GameRenderer | null = null;

  constructor(
    url: string,
    private readonly preGameRenderer: PreGameRenderer,
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
        this.preGameRenderer.renderPlayerList(
          message.payload.participants,
          this.onReadyChange.bind(this),
        );
        this.preGameRenderer.renderStartButton(
          message.payload.participants,
          message.payload.owner,
          this.onClickReady.bind(this),
        );
        break;
      }
      case "room.new-update": {
        this.preGameRenderer.renderPlayerList(
          message.payload.participants,
          this.onReadyChange.bind(this),
        );
        this.preGameRenderer.renderStartButton(
          message.payload.participants,
          message.payload.owner,
          this.onClickReady.bind(this),
        );
        break;
      }
      case "game.start": {
        // TODO: get map name from payload
        const { map } = message.payload;
        this.gameRenderer = new GameRenderer(this.preGameRenderer.root, "base4");
        this.gameRenderer.drawMap(map);
        break;
      }
      case "hydrate": {
        const { map } = message.payload.state;
        this.gameRenderer = new GameRenderer(this.preGameRenderer.root, "base4");
        this.gameRenderer.drawMap(map);
        break;
      }
    }
  }

  private onReadyChange(state: boolean) {
    this.sendMessage({
      type: "room.toggle-ready",
      payload: {
        ready: state,
        roomID: this.roomID,
      },
    });
  }

  private onClickReady() {
    this.sendMessage({
      type: "room.start-game",
      payload: {},
    });
  }

  private sendMessage<T extends keyof SettlersWSServer.OutgoingMessages>(
    message: SettlersWSServer.OutgoingMessage<T>,
  ) {
    this.ws.send(JSON.stringify(message));
  }
}
