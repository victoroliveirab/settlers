import MatchStateManager from "./state";
import type { SettlersWSServerMatch } from "./types";

function safeParse(text: string):
  | {
      [K in keyof SettlersWSServerMatch.IncomingMessages]: SettlersWSServerMatch.IncomingMessage<K>;
    }[keyof SettlersWSServerMatch.IncomingMessages]
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

export default class MatchWebSocketHandler {
  constructor(
    readonly ws: WebSocket,
    private readonly state: MatchStateManager,
  ) {
    ws.onmessage = this.onMessage.bind(this);
  }

  private onMessage(event: MessageEvent) {
    const message = safeParse(event.data);
    if (!message) return;

    switch (message.type) {
      case "setup.build-settlement": {
        const { vertices } = message.payload;
        this.state.setupSettlement(vertices);
        console.log("DONE");
        break;
      }
      case "setup.settlement-build.success": {
        const { logs, settlement } = message.payload;
        this.state.addSettlement(settlement);
        this.state.addLogs(logs);
        break;
      }
      case "setup.build-road": {
        const { edges } = message.payload;
        this.state.setupRoad(edges);
        break;
      }
      case "setup.road-build.success": {
        const { logs, road } = message.payload;
        this.state.addRoad(road);
        this.state.addLogs(logs);
        break;
      }
      case "setup.end": {
        const { logs, hands } = message.payload;
        // REFACTOR: don't like doing logic here, but for now it's fine
        const resourceCount = Object.entries(hands).reduce(
          (acc, [player, resources]) => ({
            ...acc,
            [player]: Object.values(resources).reduce((acc, quantity) => acc + quantity, 0),
          }),
          {} as Record<string, number>,
        );
        this.state.setResourcesCounts(resourceCount);
        this.state.addLogs(logs);
        break;
      }
      // Match, player round related
      case "game.your-round": {
        const { availableEdges, availableVertices, currentRoundPlayer, roundType } =
          message.payload;
        this.state.setAvailableEdges(availableEdges);
        this.state.setAvailableVertices(availableVertices);
        this.state.setRoundPlayer(currentRoundPlayer);
        this.state.setRoundType(roundType);
        break;
      }
      case "game.new-road.success": {
        const { availableEdges, hand, logs, road } = message.payload;
        this.state.setAvailableEdges(availableEdges);
        this.state.setHand(hand);
        this.state.addRoad(road);
        this.state.addLogs(logs);
        break;
      }
      // General
      case "game.dice-roll.success": {
        const { dices, hand, logs, resourceCount } = message.payload;
        this.state.setDice(dices[0], dices[1]);
        this.state.setHand(hand);
        this.state.setResourcesCounts(resourceCount);
        this.state.addLogs(logs);
        break;
      }
      // Match, opponent round related
      case "game.player-round": {
        const { currentRoundPlayer, roundType } = message.payload;
        this.state.setRoundPlayer(currentRoundPlayer);
        this.state.setRoundType(roundType);
        break;
      }
    }
    this.state.updateUI();
  }

  sendSetupNewSettlement(vertexID: number) {
    this.sendMessage({
      type: "setup.new-settlement",
      payload: { vertex: vertexID },
    });
  }

  sendSetupNewRoad(edgeID: number) {
    this.sendMessage({
      type: "setup.new-road",
      payload: { edge: edgeID },
    });
  }

  sendDiceRollRequest() {
    this.sendMessage({
      type: "game.dice-roll",
      payload: {},
    });
  }

  sendMatchNewSettlement(vertexID: number) {
    this.sendMessage({
      type: "game.new-settlement",
      payload: { vertex: vertexID },
    });
  }

  sendMatchNewRoad(edgeID: number) {
    this.sendMessage({
      type: "game.new-road",
      payload: { edge: edgeID },
    });
  }

  sendRobberNewPosition(tileID: number) {
    this.sendMessage({
      type: "game.move-robber",
      payload: { tile: tileID },
    });
  }

  sendEndRound() {
    this.sendMessage({
      type: "game.end-round",
      payload: {},
    });
  }

  private sendMessage<T extends keyof SettlersWSServerMatch.OutgoingMessages>(
    message: SettlersWSServerMatch.OutgoingMessage<T>,
  ) {
    this.ws.send(JSON.stringify(message));
  }
}
