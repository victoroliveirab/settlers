import GameRenderer from "../renderer/game";
import PreGameRenderer from "../renderer/pre-game";
import WebSocketConnection from "../websocket";
import { SettlersCore } from "../websocket/types";

type UIPart =
  | "devHand"
  | "dice"
  | "hand"
  | "map"
  | "participantList"
  | "playerList"
  | "road"
  | "startButton"
  | "vertices";

const defaultUpdateUIState: Record<UIPart, boolean> = {
  devHand: false,
  dice: false,
  hand: false,
  map: false,
  participantList: false,
  playerList: false,
  road: false,
  startButton: false,
  vertices: false,
};

export default class GameState {
  private participants: SettlersCore.Participant[] = [];
  private players: SettlersCore.Player[] = [];
  private currentRoundPlayer: string = "";
  private map!: SettlersCore.Map;
  private owner: SettlersCore.Player | null = null;
  private cities: SettlersCore.Building[] = [];
  private roads: SettlersCore.Building[] = [];
  private settlements: SettlersCore.Building[] = [];
  private hand: SettlersCore.Hand = {
    Brick: 0,
    Grain: 0,
    Lumber: 0,
    Ore: 0,
    Sheep: 0,
  };
  private devHand: SettlersCore.DevHand = {
    Knight: 0,
    Monopoly: 0,
    "Road Building": 0,
    "Victory Point": 0,
    "Year of Plenty": 0,
  };

  private dices: [number, number] = [0, 0];
  private resourceCount!: Record<string, number>;

  private preGameRenderer: PreGameRenderer;
  private gameRenderer!: GameRenderer;
  private service!: WebSocketConnection;

  private shouldUpdateUIPart: Record<UIPart, boolean> = { ...defaultUpdateUIState };

  constructor(
    private readonly pregameRoot: HTMLElement,
    private readonly root: HTMLElement,
    readonly userName: string,
  ) {
    this.preGameRenderer = new PreGameRenderer(pregameRoot);
  }

  setService(service: WebSocketConnection) {
    this.service = service;
  }

  setMap(map: SettlersCore.Map) {
    if (this.map) {
      throw new Error("should not set map twice");
    }
    this.pregameRoot.remove();
    this.gameRenderer = new GameRenderer(this.root, "base4");
    this.map = map;
    this.shouldUpdateUIPart.map = true;
  }

  setParticipants(participants: SettlersCore.Participant[]) {
    this.participants = participants;
    this.shouldUpdateUIPart.participantList = true;
  }

  setOwner(player: SettlersCore.Player["name"]) {
    this.shouldUpdateUIPart.startButton = true;
    if (this.owner?.name === player) return;
    for (const participant of this.participants) {
      if (participant.player?.name === player) {
        this.owner = participant.player;
        return;
      }
    }
    console.warn(`trying to set player#${player} as owner, but they are not a participant`);
  }

  setPlayers(players: SettlersCore.Player[]) {
    this.players = players;
    this.shouldUpdateUIPart.playerList = true;
  }

  setCurrentRoundPlayer(player: string) {
    this.currentRoundPlayer = player;
    this.shouldUpdateUIPart.playerList = true;
    this.shouldUpdateUIPart.dice = true;
  }

  setHand(hand: SettlersCore.Hand) {
    this.hand = hand;
    this.shouldUpdateUIPart.hand = true;
  }

  setDevHand(devHand: SettlersCore.DevHand) {
    this.devHand = devHand;
    this.shouldUpdateUIPart.devHand = true;
  }

  setResourcesCounts(counts: Record<string, number>) {
    this.resourceCount = counts;
    this.shouldUpdateUIPart.playerList = true;
  }

  setDices(dice1: number, dice2: number) {
    this.dices[0] = dice1;
    this.dices[1] = dice2;
    this.shouldUpdateUIPart.dice = true;
  }

  addSettlement(settlement: SettlersCore.Building) {
    this.settlements.push(settlement);

    const color = this.players.find(({ name }) => name === settlement.owner)?.color;
    if (!color) {
      console.warn("addSettlement: color not found for owner:");
      return;
    }
    this.shouldUpdateUIPart.vertices = true;
  }

  addRoad(road: SettlersCore.Building) {
    this.roads.push(road);

    const color = this.players.find(({ name }) => name === road.owner)?.color;
    if (!color) {
      console.warn("addRoad: color not found for owner:");
      return;
    }
    this.shouldUpdateUIPart.road = true;
  }

  addLogs(logs: string[]) {
    logs.forEach((log) => {
      this.gameRenderer.renderNewLog(log);
    });
  }

  enableVerticesToBuildSettlement(vertices: number[], phase: "game" | "setup") {
    this.gameRenderer.makeVerticesClickable(vertices, (vertexID) => {
      this.service.onSettlementPositionChose(phase, vertexID);
    });
  }

  enableEdgesToBuildRoad(edges: number[], phase: "game" | "setup") {
    this.gameRenderer.makeEdgesClickable(edges, (edgeID) => {
      this.service.onRoadPositionChose(phase, edgeID);
    });
  }

  repaintScreen() {
    for (const [uiPart, shouldRerender] of Object.entries(this.shouldUpdateUIPart)) {
      if (shouldRerender) {
        switch (uiPart as UIPart) {
          case "map": {
            this.gameRenderer.drawMap(this.map);
            break;
          }
          case "participantList": {
            this.preGameRenderer.renderParticipantList(
              this.participants,
              this.userName,
              (state) => {
                this.service.onReadyChange(state);
              },
            );
            break;
          }
          case "playerList": {
            const players = this.players.map((player) => ({
              color: player.color,
              devHandCount: 0,
              isCurrentRound: player.name === this.currentRoundPlayer,
              knights: 0,
              longestRoad: 0,
              name: player.name,
              points: 0,
              resourceCount: this.resourceCount[player.name],
            }));
            this.gameRenderer.drawPlayers(players);
            break;
          }
          case "dice": {
            if (
              this.currentRoundPlayer === this.userName &&
              this.dices[0] === 0 &&
              this.dices[1] === 0
            ) {
              this.gameRenderer.drawDices(this.dices, () => {
                this.service.onDiceRollRequested();
              });
              this.gameRenderer.updatePassButton();
            } else {
              this.gameRenderer.drawDices(this.dices);
              this.gameRenderer.updatePassButton(() => {
                this.service.onEndRound();
              });
            }
            break;
          }
          case "hand": {
            this.gameRenderer.drawHand(this.hand);
            break;
          }
          case "devHand": {
            this.gameRenderer.drawDevHand(this.devHand);
            break;
          }
          case "startButton": {
            this.preGameRenderer.renderStartButton(
              this.participants,
              this.userName,
              this.owner?.name ?? null,
              () => {
                this.service.onClickStart();
              },
            );
            break;
          }
          case "vertices": {
            this.settlements.forEach((settlement) => {
              const color = this.players.find(({ name }) => name === settlement.owner)!.color;
              this.gameRenderer.drawSettlement(settlement, color);
            });
            break;
          }
          case "road": {
            this.roads.forEach((road) => {
              const color = this.players.find(({ name }) => name === road.owner)!.color;
              this.gameRenderer.drawRoad(road, color);
            });
            break;
          }
          default: {
            console.warn("unknown ui part:", uiPart);
          }
        }
      }
    }
    this.shouldUpdateUIPart = { ...defaultUpdateUIState };
  }
}
