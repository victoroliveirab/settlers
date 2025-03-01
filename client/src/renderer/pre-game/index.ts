import type { SettlersCore } from "../../websocket/types";

export default class PreGameRenderer {
  constructor(readonly root: HTMLElement) {}

  renderParticipantList(
    participants: SettlersCore.Participant[],
    userName: SettlersCore.Player["name"],
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
      this.renderParticipant(participant, userName, playersContainer, onReadyChange);
    }
  }

  private renderParticipant(
    participant: SettlersCore.Participant,
    userName: SettlersCore.Player["name"],
    container: HTMLElement,
    onReadyChange: (state: boolean) => void,
  ) {
    const element = document.createElement("div");
    element.classList.add("pre-game-spot");
    if (participant.player) {
      element.classList.add("pre-game-player");
      element.style.background = participant.player.color;
      element.textContent = participant.player.name;

      const readyCheckbox = document.createElement("input");
      readyCheckbox.type = "checkbox";
      const isCheckboxActive = participant.player.name === userName;
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

  renderStartButton(
    participants: SettlersCore.Participant[],
    userName: SettlersCore.Player["name"],
    owner: SettlersCore.Player["name"] | null,
    onClick: () => void,
  ) {
    const startButton = document.createElement("button");
    startButton.textContent = "Start";
    startButton.disabled = true;
    const isReady = participants.every((participant) => participant.ready);
    this.root.querySelector(".pre-game-container")?.appendChild(startButton);
    if (isReady) {
      const isRoomOwner = owner === userName;
      startButton.disabled = !isRoomOwner;
      if (isRoomOwner) {
        startButton.addEventListener("click", onClick, { once: true });
      }
    }
  }
}
