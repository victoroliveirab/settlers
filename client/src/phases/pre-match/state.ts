import { SettlersCore } from "../../core/types";
import MatchStateManager from "../match/state";
import PreMatchRenderer from "./renderer";
import PreMatchWebSocketHandler from "./websocket";

type UIPart = "participantList" | "startButton";

const defaultUpdateUIState: Record<UIPart, boolean> = {
  participantList: false,
  startButton: false,
};

export default class PreMatchStateManager {
  private handler: PreMatchWebSocketHandler;
  private renderer: PreMatchRenderer;
  private participants: SettlersCore.Participant[] = [];
  private owner: SettlersCore.Participant["player"] | null = null;

  private shouldUpdateUIPart: Record<UIPart, boolean> = { ...defaultUpdateUIState };

  constructor(
    ws: WebSocket,
    firstMessage: MessageEvent,
    root: HTMLElement,
    readonly matchRoot: HTMLElement,
    readonly userName: string,
    readonly roomID: string,
  ) {
    this.renderer = new PreMatchRenderer(root);
    this.handler = new PreMatchWebSocketHandler(ws, firstMessage, this);
  }

  handleStartSetup(
    players: SettlersCore.Player[],
    initialResourceCount: Record<SettlersCore.Player["name"], number>,
    firstPlayer: SettlersCore.Player["name"],
    map: SettlersCore.Map,
    mapName: string,
    logs: string[],
  ) {
    this.renderer.destroy();
    const state = new MatchStateManager(
      this.handler.ws,
      this.matchRoot,
      this.userName,
      mapName,
      map,
      players,
    );
    state.setResourcesCounts(initialResourceCount);
    state.setRoundPlayer(firstPlayer);
    state.addLogs(logs);
    state.updateUI();
  }

  setParticipants(participants: SettlersCore.Participant[]) {
    this.shouldUpdateUIPart.participantList = true;
    this.participants = participants;
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

  updateUI() {
    for (const [uiPart, shouldRerender] of Object.entries(this.shouldUpdateUIPart)) {
      if (!shouldRerender) return;
      switch (uiPart as UIPart) {
        case "participantList": {
          this.renderer.renderParticipantList(this.participants, this.userName, (state) => {
            this.handler.sendReadyState(this.roomID, state);
          });
          break;
        }
        case "startButton": {
          this.renderer.renderStartButton(
            this.participants,
            this.userName,
            this.owner?.name ?? null,
            () => {
              this.handler.sendStartGame();
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
}
