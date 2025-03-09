import PreMatchStateManager from "./state";
import type { SettlersWSServerPreMatch } from "./types";

function safeParse(text: string):
  | {
      [K in keyof SettlersWSServerPreMatch.IncomingMessages]: SettlersWSServerPreMatch.IncomingMessage<K>;
    }[keyof SettlersWSServerPreMatch.IncomingMessages]
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

export default class PreMatchWebSocketHandler {
  constructor(
    readonly ws: WebSocket,
    firstEvent: MessageEvent,
    private readonly state: PreMatchStateManager,
  ) {
    ws.onmessage = this.onMessage.bind(this);
    this.onMessage(firstEvent);
  }

  private onMessage(event: MessageEvent) {
    const message = safeParse(event.data);
    if (!message) return;

    switch (message.type) {
      case "room.connect.success": {
        const { owner, participants } = message.payload;
        this.state.setParticipants(participants);
        this.state.setOwner(owner);
        break;
      }
      case "room.new-update": {
        const { owner, participants } = message.payload;
        this.state.setParticipants(participants);
        this.state.setOwner(owner);
        break;
      }
      case "game.start": {
        const { currentRoundPlayer, logs, map, players, resourceCount } = message.payload;
        this.state.handleStartSetup(players, resourceCount, currentRoundPlayer, map, "base4", logs);
        break;
      }
      default: {
        console.warn(`Unknown message type: ${(message as any).type}`);
        return;
      }
    }
    this.state.updateUI();
  }

  sendReadyState(roomID: string, state: boolean) {
    this.sendMessage({
      type: "room.toggle-ready",
      payload: {
        ready: state,
        roomID,
      },
    });
  }

  sendStartGame() {
    this.sendMessage({
      type: "room.start-game",
      payload: {},
    });
  }

  private sendMessage<T extends keyof SettlersWSServerPreMatch.OutgoingMessages>(
    message: SettlersWSServerPreMatch.OutgoingMessage<T>,
  ) {
    this.ws.send(JSON.stringify(message));
  }
}
