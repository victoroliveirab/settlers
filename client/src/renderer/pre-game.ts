import type { SettlersCore } from "../websocket/types";

export default class PreGameRenderer {
  constructor(
    private readonly root: HTMLElement,
    private readonly userID: number,
  ) {}

  renderPlayerList(
    participants: SettlersCore.Participant[],
    onReadyChange: (state: boolean) => void,
  ) {
    this.root.innerHTML = "";
    const container = document.createElement("div");
    container.classList.add("pre-game-container");
    this.root.appendChild(container);
    const playersContainer = document.createElement("div");
    playersContainer.classList.add("players");
    container.appendChild(playersContainer);
    for (const participant of participants) {
      this.renderParticipant(participant, playersContainer, onReadyChange);
    }
  }

  private renderParticipant(
    participant: SettlersCore.Participant,
    container: HTMLElement,
    onReadyChange: (state: boolean) => void,
  ) {
    const element = document.createElement("div");
    element.classList.add("pre-game-spot");
    if (participant.player) {
      element.classList.add("pre-game-player");
      element.style.background = participant.player.color;
      element.textContent = participant.player.username;

      const readyCheckbox = document.createElement("input");
      readyCheckbox.type = "checkbox";
      const isCheckboxActive = participant.player.id === this.userID;
      console.log({ participant, userID: this.userID });
      readyCheckbox.disabled = !isCheckboxActive;
      readyCheckbox.checked = participant.ready;

      if (isCheckboxActive) {
        readyCheckbox.addEventListener(
          "change",
          () => {
            onReadyChange(readyCheckbox.checked);
          },
          { once: true },
        );
      }

      element.appendChild(readyCheckbox);
    }
    container.appendChild(element);
  }

  renderStartButton(participants: SettlersCore.Participant[]) {
    const startButton = document.createElement("button");
    startButton.textContent = "Start";
    startButton.disabled = true;
    const isReady = participants.every((participant) => participant.ready);
    this.root.querySelector(".pre-game-container")?.appendChild(startButton);
    if (isReady) {
      // TODO: have on BE room owner
      const isRoomOwner = participants[0].player?.id === this.userID;
      startButton.disabled = !isRoomOwner;
    }
  }
}
