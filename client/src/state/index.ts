import GameRenderer from "../renderer/game";
import PreGameRenderer from "../renderer/pre-game";
import WebSocketConnection from "../websocket";
import { SettlersCore } from "../websocket/types";

type UIPart = "map" | "participantList" | "startButton";

const defaultUpdateUIState: Record<UIPart, boolean> = {
  map: false,
  participantList: false,
  startButton: false,
};

export default class GameState {
  private participants: SettlersCore.Participant[] = [];
  private players: SettlersCore.Player[] = [];
  private map!: SettlersCore.Map;
  private owner: SettlersCore.Player | null = null;

  private preGameRenderer: PreGameRenderer;
  private gameRenderer!: GameRenderer;
  private service!: WebSocketConnection;

  private shouldUpdateUIPart: Record<UIPart, boolean> = { ...defaultUpdateUIState };
  private hasSetInitialState: boolean = false;

  constructor(
    private readonly root: HTMLElement,
    private readonly userName: string,
  ) {
    this.preGameRenderer = new PreGameRenderer(root);
  }

  setService(service: WebSocketConnection) {
    this.service = service;
  }

  setInitialState(map: SettlersCore.Map, players: SettlersCore.Player[]) {
    if (this.hasSetInitialState) {
      throw new Error("should not set initial state twice");
    }
    this.gameRenderer = new GameRenderer(this.root, "base4");
    this.map = map;
    this.players = players;
    this.hasSetInitialState = true;
    this.shouldUpdateUIPart.map = true;
  }

  setParticipants(participants: SettlersCore.Participant[]) {
    this.participants = participants;
    this.shouldUpdateUIPart.participantList = true;
  }

  setOwner(playerName: SettlersCore.Player["name"]) {
    this.shouldUpdateUIPart.startButton = true;
    if (this.owner?.name === playerName) return;
    for (const participant of this.participants) {
      if (participant.player?.name === playerName) {
        this.owner = participant.player;
        return;
      }
    }
    console.warn(`trying to set player#${playerName} as owner, but they are not a participant`);
  }

  enableVerticesToBuildSettlement(vertices: number[], phase: "game" | "setup") {
    this.gameRenderer.makeVerticesClickable(vertices, (vertexID) => {
      this.service.onSettlementPositionChose(phase, vertexID);
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
          default: {
            console.warn("unknown ui part:", uiPart);
          }
        }
      }
    }
    this.shouldUpdateUIPart = { ...defaultUpdateUIState };
  }
}
