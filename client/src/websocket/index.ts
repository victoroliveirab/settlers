import GameRenderer from "../renderer/game";
import PreGameRenderer from "../renderer/pre-game";
import GameState from "../state";
import { SettlersCore, SettlersWSServer } from "./types";

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
        const { currentRoundPlayer, logs, map, players, resourceCount } = message.payload;
        this.stateManager.setMap(map);
        this.stateManager.setPlayers(players);
        this.stateManager.setCurrentRoundPlayer(currentRoundPlayer);
        this.stateManager.setResourcesCounts(resourceCount);
        this.stateManager.addLogs(logs);
        break;
      }
      case "setup.build-settlement": {
        const { vertices } = message.payload;
        this.stateManager.enableVerticesToBuildSettlement(vertices, "setup");
        break;
      }
      case "setup.settlement-build.success": {
        const { logs, settlement } = message.payload;
        this.stateManager.addSettlement(settlement);
        this.stateManager.addLogs(logs);
        break;
      }
      case "setup.build-road": {
        const { edges } = message.payload;
        this.stateManager.enableEdgesToBuildRoad(edges, "setup");
        break;
      }
      case "setup.road-build.success": {
        const { logs, road } = message.payload;
        this.stateManager.addRoad(road);
        this.stateManager.addLogs(logs);
        break;
      }
      case "setup.player-round-changed": {
        const { currentRoundPlayer } = message.payload;
        this.stateManager.setCurrentRoundPlayer(currentRoundPlayer);
        break;
      }
      case "setup.end": {
        const { hands, logs } = message.payload;
        // REFACTOR: don't like doing logic here, but for now it's fine
        const resourceCount = Object.entries(hands).reduce(
          (acc, [player, resources]) => ({
            ...acc,
            [player]: Object.values(resources).reduce((acc, quantity) => acc + quantity, 0),
          }),
          {} as Record<string, number>,
        );
        this.stateManager.setPhase("game");
        this.stateManager.setResourcesCounts(resourceCount);
        this.stateManager.setHand(hands[this.stateManager.userName]);
        this.stateManager.addLogs(logs);
        break;
      }
      case "game.player-round": {
        const { currentRoundPlayer } = message.payload;
        this.stateManager.setCurrentRoundPlayer(currentRoundPlayer);
        break;
      }
      case "game.dice-roll.success": {
        const { dices, hand, logs, resourceCount } = message.payload;
        this.stateManager.setDices(dices[0], dices[1]);
        this.stateManager.setHand(hand);
        this.stateManager.setResourcesCounts(resourceCount);
        this.stateManager.addLogs(logs);
        break;
      }
      case "game.discard-cards-request": {
        const { quantityByPlayers } = message.payload;
        this.stateManager.setQuantitiesToDiscard(quantityByPlayers);
        break;
      }
      case "game.discard-cards.success": {
        const { hand, resourceCount } = message.payload;
        this.stateManager.setHand(hand);
        this.stateManager.setResourcesCounts(resourceCount);
        break;
      }
      case "game.discarded-cards": {
        const { resourceCount, quantityByPlayers, logs } = message.payload;
        this.stateManager.setQuantitiesToDiscard(quantityByPlayers);
        this.stateManager.setResourcesCounts(resourceCount);
        this.stateManager.addLogs(logs);
        break;
      }
      case "game.move-robber-request": {
        const { availableTiles } = message.payload;
        this.stateManager.enableRobberMovement(availableTiles);
        break;
      }
      case "hydrate": {
        const { currentRoundPlayer, dice, hand, map, players, resourceCount, roads, settlements } =
          message.payload.state;
        this.stateManager.setPhase("game");
        this.stateManager.setMap(map);
        this.stateManager.setPlayers(players);
        this.stateManager.setCurrentRoundPlayer(currentRoundPlayer);
        this.stateManager.setHand(hand);
        this.stateManager.setResourcesCounts(resourceCount);
        this.stateManager.setDices(dice[0], dice[1]);
        Object.values(settlements).forEach((settlement) => {
          this.stateManager.addSettlement(settlement);
        });
        Object.values(roads).forEach((road) => {
          this.stateManager.addRoad(road);
        });
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

  onRoadPositionChose(phase: "game" | "setup", edgeID: number) {
    this.sendMessage({
      type: phase === "game" ? "game.new-road" : "setup.new-road",
      payload: {
        edge: edgeID,
      },
    });
  }

  onDiceRollRequested() {
    this.sendMessage({
      type: "game.dice-roll",
      payload: {},
    });
  }

  onDiscardCardsSelected(cards: Record<SettlersCore.Resource, number>) {
    this.sendMessage({
      type: "game.discard-cards",
      payload: {
        resources: cards,
      },
    });
  }

  onRobberNewPositionSelected(tileID: number) {
    this.sendMessage({
      type: "game.move-robber",
      payload: {
        tile: tileID,
      },
    });
  }

  onEndRound() {
    this.sendMessage({
      type: "game.end-round",
      payload: {},
    });
  }

  private sendMessage<T extends keyof SettlersWSServer.OutgoingMessages>(
    message: SettlersWSServer.OutgoingMessage<T>,
  ) {
    this.ws.send(JSON.stringify(message));
  }
}
