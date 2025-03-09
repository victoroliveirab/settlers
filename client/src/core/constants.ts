import { SettlersCore } from "./types";

export const resourcesOrder: SettlersCore.Resource[] = ["Lumber", "Brick", "Sheep", "Grain", "Ore"];

export const emojis = Object.freeze({
  devCards: {
    Knight: "⚔️",
    "Victory Point": "🎖️",
    "Road Building": "🛤️",
    "Year of Plenty": "🎁",
    Monopoly: "🎩",
  },
  misc: {
    discarding: "❌",
  },
  resources: {
    Lumber: "🌲",
    Brick: "🧱",
    Sheep: "🐑",
    Grain: "🌾",
    Ore: "⛰️",
  },
});

export const resourceEmojis = Object.freeze({});

export const developmentEmojis = Object.freeze({});

export const noop = () => {};
