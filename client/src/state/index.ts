import GameRenderer from "../renderer/game";
import PreGameRenderer from "../renderer/pre-game";
import WebSocketConnection from "../websocket";
import { SettlersCore } from "../websocket/types";

type UIPart =
  | "dice"
  | "hud"
  | "map"
  | "participantList"
  | "playerList"
  | "road"
  | "startButton"
  | "vertices";

const defaultUpdateUIState: Record<UIPart, boolean> = {
  dice: false,
  hud: false,
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
  private hand!: SettlersCore.Hand;
  private dices: [number, number] = [0, 0];

  private preGameRenderer: PreGameRenderer;
  private gameRenderer!: GameRenderer;
  private service!: WebSocketConnection;

  private shouldUpdateUIPart: Record<UIPart, boolean> = { ...defaultUpdateUIState };
  private hasSetInitialState: boolean = false;

  constructor(
    private readonly pregameRoot: HTMLElement,
    private readonly root: HTMLElement,
    private readonly userName: string,
  ) {
    this.preGameRenderer = new PreGameRenderer(pregameRoot);
  }

  setService(service: WebSocketConnection) {
    this.service = service;
  }

  setInitialState(
    map: SettlersCore.Map,
    players: SettlersCore.Player[],
    hand: SettlersCore.Hand | null,
    currentRoundPlayer: string,
  ) {
    if (this.hasSetInitialState) {
      throw new Error("should not set initial state twice");
    }
    this.pregameRoot.remove();
    this.gameRenderer = new GameRenderer(this.root, "base4");
    this.map = map;
    this.players = players;
    if (hand) this.hand = hand;
    this.setCurrentRoundPlayer(currentRoundPlayer);

    this.hasSetInitialState = true;
    this.shouldUpdateUIPart.map = true;
    this.shouldUpdateUIPart.playerList = true;
    this.shouldUpdateUIPart.hud = true;
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

  setCurrentRoundPlayer(player: string) {
    this.currentRoundPlayer = player;
    this.shouldUpdateUIPart.playerList = true;
    if (this.userName === player) {
      console.log("HERE?");
      this.shouldUpdateUIPart.dice = true;
    }
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
            this.gameRenderer.drawPlayers(this.players, this.currentRoundPlayer);
            break;
          }
          case "dice": {
            if (this.currentRoundPlayer === this.userName) {
              console.log("IF");
              this.gameRenderer.drawDices(this.dices, () => {
                this.service.onDiceRollRequested();
              });
            } else {
              console.log("ELSE");
              this.gameRenderer.drawDices(this.dices);
            }
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
