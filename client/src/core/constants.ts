import { SettlersCore } from "./types";

export const roundTypes = [
  "SettlementSetup#1",
  "RoadSetup#1",
  "SettlementSetup#2",
  "RoadSetup#2",
  "FirstRound",
  "Regular",
  "MoveRobber(7)",
  "MoveRobber(Knight)",
  "ChooseRobbedPlayer",
  "BetweenRounds",
  "BuildRoadDevelopment(1)",
  "BuildRoadDevelopment(2)",
  "MonopolyPickResource",
  "YearOfPlentyPickResources",
  "DiscardPhase",
  "GameOver",
];

export const roundTypesByName = {
  "SettlementSetup#1": 0,
  "RoadSetup#1": 1,
  "SettlementSetup#2": 2,
  "RoadSetup#2": 3,
  FirstRound: 4,
  Regular: 5,
  "MoveRobber(7)": 6,
  "MoveRobber(Knight)": 7,
  ChooseRobbedPlayer: 8,
  BetweenRounds: 9,
  "BuildRoadDevelopment(1)": 10,
  "BuildRoadDevelopment(2)": 11,
  MonopolyPickResource: 12,
  YearOfPlentyPickResources: 13,
  DiscardPhase: 14,
  GameOver: 15,
};

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

export const noop = () => {};
