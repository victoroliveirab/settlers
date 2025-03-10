import type { SettlersCore } from "../../core/types";

type Param = {
  description: string;
  key: string;
  label: string;
  value: number;
  values: number[];
};

export default class PreMatchRenderer {
  constructor(readonly root: HTMLElement) {}

  renderParticipantList(
    participants: SettlersCore.Participant[],
    userName: SettlersCore.Player["name"],
    onReadyChange: (state: boolean) => void,
  ) {
    const container = this.root.querySelector<HTMLDivElement>("#player-list")!;
    container.innerHTML = "";
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
    console.log({ participants, userName, owner });
    if (userName !== owner) return;
    const isReady = participants.every((participant) => participant.ready);
    const startButton = this.root.querySelector<HTMLButtonElement>("#start-game")!;
    startButton.style.visibility = "visible";
    startButton.disabled = !isReady;
    if (isReady) {
      startButton.addEventListener("click", onClick, { once: true });
    }
  }

  renderParams(params: Param[], onChange: (key: string, value: number) => void) {
    const container = this.root.querySelector<HTMLDivElement>("#params")!;
    container.innerHTML = "";
    const selects: HTMLSelectElement[] = [];
    params.forEach((param) => {
      const element = document.createElement("div");
      element.classList.add("param");
      const label = document.createElement("h3");
      label.textContent = param.label;
      element.appendChild(label);

      const select = document.createElement("select");
      select.addEventListener(
        "change",
        () => {
          selects.forEach((el) => {
            el.disabled = true;
          });
          onChange(param.key, Number(select.value));
        },
        { once: true },
      );
      param.values.forEach((value) => {
        const option = document.createElement("option");
        option.textContent = String(value);
        option.value = String(value);
        select.appendChild(option);
      });
      select.value = String(param.value);

      selects.push(select);
      element.appendChild(select);
      container.appendChild(element);
    });
  }

  destroy() {
    this.root.remove();
  }
}
