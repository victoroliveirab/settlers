import type { SettlersCore } from "../../core/types";

type Param = {
  description: string;
  key: string;
  label: string;
  value: number;
  values: number[];
};

type Color = {
  disabled: boolean;
  label: string;
  value: string;
};

export default class PreMatchRenderer {
  constructor(readonly root: HTMLElement) {}

  renderParticipantList(
    participants: SettlersCore.Participant[],
    userName: SettlersCore.Player["name"],
    colors: Color[],
    onReadyChange: (state: boolean) => void,
    onColorChange: (color: string) => void,
  ) {
    const container = this.root.querySelector<HTMLDivElement>("#player-list")!;
    container.innerHTML = "";
    const playersContainer = document.createElement("div");
    playersContainer.classList.add("players");
    container.appendChild(playersContainer);
    for (const participant of participants) {
      this.renderParticipant(
        participant,
        userName,
        colors,
        playersContainer,
        onReadyChange,
        onColorChange,
      );
    }
  }

  private renderParticipant(
    participant: SettlersCore.Participant,
    userName: SettlersCore.Player["name"],
    colors: Color[],
    container: HTMLElement,
    onReadyChange: (state: boolean) => void,
    onColorChange: (color: string) => void,
  ) {
    const element = document.createElement("div");
    element.classList.add("pre-game-spot");
    if (participant.player) {
      element.classList.add("pre-game-player");
      element.style.background = participant.player.color;

      const playerName = document.createElement("h2");
      playerName.textContent = participant.player.name;
      element.appendChild(playerName);

      const colorPicker = document.createElement("select");
      colors.forEach((color) => {
        const option = document.createElement("option");
        option.value = color.value;
        option.textContent = color.label;
        option.disabled = color.disabled && participant.player?.color !== color.value;
        colorPicker.appendChild(option);
      });
      colorPicker.value = participant.player.color;

      if (participant.player.name !== userName) {
        colorPicker.disabled = true;
      } else {
        colorPicker.addEventListener(
          "change",
          () => {
            colorPicker.disabled = true;
            onColorChange(colorPicker.value);
          },
          { once: true },
        );
      }

      element.appendChild(colorPicker);

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
    if (userName !== owner) return;
    const isReady = participants.every((participant) => participant.ready);
    const startButton = this.root.querySelector<HTMLButtonElement>("#start-game")!;
    startButton.style.visibility = "visible";
    startButton.disabled = !isReady;
    if (isReady) {
      startButton.addEventListener("click", onClick, { once: true });
    }
  }

  renderParams(params: Param[], onChange?: (key: string, value: number) => void) {
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
      if (onChange) {
        select.disabled = false;
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
      } else {
        select.disabled = true;
      }
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
