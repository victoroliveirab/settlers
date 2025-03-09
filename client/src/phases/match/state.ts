import type { SettlersCore } from "../../core/types";
import MatchRenderer from "./renderer";
import { setDice } from "./state/dice";
import { setDevHand, setHand } from "./state/hand";
import {
  setPlayers,
  setQuantitiesToDiscard,
  setResourcesCounts,
  setRoundPlayer,
} from "./state/players";
import { addRoad, setupRoad } from "./state/road";
import { enableRobber } from "./state/robber";
import { addSettlement, setupSettlement } from "./state/settlement";
import MatchWebSocketHandler from "./websocket";

type InitialData = {
  firstPlayer: SettlersCore.Player["name"];
  logs: string[];
  map: SettlersCore.Map;
  mapName: string;
  players: SettlersCore.Player[];
  resourceCount: Record<SettlersCore.Player["name"], number>;
};

type UIPart =
  | "devHand"
  | "dice"
  | "diceAction"
  | "discard"
  | "hand"
  | "passAction"
  | "playerList"
  | "edges"
  | "vertices";

const defaultUpdateUIState: Record<UIPart, boolean> = {
  devHand: false,
  dice: false,
  diceAction: false,
  discard: false,
  edges: false,
  hand: false,
  passAction: false,
  playerList: false,
  vertices: false,
};

export default class MatchStateManager {
  protected handler: MatchWebSocketHandler;
  protected renderer: MatchRenderer;

  // Dice
  protected dice: [number, number] = [0, 0];
  setDice: (dice1: number, dice2: number) => void;

  // Hands
  protected devHand: SettlersCore.DevHand = {
    Knight: 0,
    Monopoly: 0,
    "Road Building": 0,
    "Victory Point": 0,
    "Year of Plenty": 0,
  };
  protected hand: SettlersCore.Hand = {
    Brick: 0,
    Grain: 0,
    Lumber: 0,
    Ore: 0,
    Sheep: 0,
  };
  setDevHand: (devHand: SettlersCore.DevHand) => void;
  setHand: (hand: SettlersCore.Hand) => void;

  // Map
  private map: SettlersCore.Map;

  // Players
  protected currentRoundPlayer!: string;
  protected discardQuantityByPlayers: Record<SettlersCore.Player["name"], number>;
  protected resourceCount!: Record<SettlersCore.Player["name"], number>;
  protected players!: SettlersCore.Player[];
  setPlayers: (players: SettlersCore.Player[]) => void;
  setQuantitiesToDiscard: (quantityByPlayers: Record<SettlersCore.Player["name"], number>) => void;
  setResourcesCounts: (counts: Record<SettlersCore.Resource, number>) => void;
  setRoundPlayer: (player: SettlersCore.Player["name"]) => void;

  // Roads
  protected roads: SettlersCore.Building[] = [];
  addRoad: (road: SettlersCore.Building) => void;
  setupRoad: (vertices: number[]) => void;

  // Robber
  protected robbers: number[];
  enableRobber: (availableTiles: number[]) => void;

  // Settlements
  protected settlements: SettlersCore.Building[];
  addSettlement: (settlement: SettlersCore.Building) => void;
  setupSettlement: (vertices: number[]) => void;

  // Buildings
  protected cities: SettlersCore.Building[] = [];

  // Internal
  protected shouldUpdateUIPart: Record<UIPart, boolean> = { ...defaultUpdateUIState };

  constructor(
    ws: WebSocket,
    private readonly root: HTMLElement,
    readonly userName: string,
    data: InitialData,
  ) {
    this.renderer = new MatchRenderer(root, data.mapName);
    this.handler = new MatchWebSocketHandler(ws, this);

    // Dice
    this.setDice = setDice.bind(this);

    // Hands
    this.setDevHand = setDevHand.bind(this);
    this.setHand = setHand.bind(this);

    // Map
    this.map = data.map;

    // Player
    this.setResourcesCounts = setResourcesCounts.bind(this);
    this.setRoundPlayer = setRoundPlayer.bind(this);
    this.setPlayers = setPlayers.bind(this);
    this.setPlayers(data.players);
    this.discardQuantityByPlayers = this.players.reduce(
      (acc, player) => ({
        ...acc,
        [player.name]: 0,
      }),
      {},
    );
    this.setQuantitiesToDiscard = setQuantitiesToDiscard.bind(this);
    this.setResourcesCounts(data.resourceCount);
    this.setRoundPlayer(data.firstPlayer);

    // Road
    this.roads = [];
    this.addRoad = addRoad.bind(this);
    this.setupRoad = setupRoad.bind(this);

    // Robbers
    this.robbers = this.map.filter((tile) => tile.blocked).map((tile) => tile.id);
    this.enableRobber = enableRobber.bind(this);

    // Settlement
    this.settlements = [];
    this.addSettlement = addSettlement.bind(this);
    this.setupSettlement = setupSettlement.bind(this);

    // Has to be drawn before everything else
    this.renderer.drawMap(this.map);
    this.renderer.drawRobbers(this.robbers);
  }

  // Logs

  addLogs(logs: string[]) {
    logs.forEach((log) => {
      this.renderer.renderNewLog(log);
    });
  }

  updateUI() {
    for (const [uiPart, shouldRerender] of Object.entries(this.shouldUpdateUIPart)) {
      if (shouldRerender) {
        switch (uiPart as UIPart) {
          case "devHand": {
            this.renderer.drawDevHand(this.devHand);
            break;
          }
          case "dice": {
            this.renderer.drawDices(this.dice);
            break;
          }
          case "diceAction": {
            this.renderer.attachClickHandlerToDice(() => {
              this.handler.sendDiceRollRequest();
            });
            break;
          }
          case "edges": {
            this.roads.forEach((road) => {
              const color = this.players.find(({ name }) => name === road.owner)!.color;
              this.renderer.drawRoad(road, color);
            });
            break;
          }
          case "hand": {
            this.renderer.drawHand(this.hand);
            break;
          }
          case "passAction": {
            if (this.dice[0] > 0 && this.dice[1] > 0 && this.currentRoundPlayer === this.userName) {
              this.renderer.updatePassButton(() => {
                this.handler.sendEndRound();
              });
            } else {
              this.renderer.updatePassButton();
            }
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
              quantityToDiscard: this.discardQuantityByPlayers[player.name],
              resourceCount: this.resourceCount[player.name],
            }));
            this.renderer.drawPlayers(players);
            break;
          }
          case "vertices": {
            console.log("VERTICES");
            this.settlements.forEach((settlement) => {
              const color = this.players.find(({ name }) => name === settlement.owner)!.color;
              this.renderer.drawSettlement(settlement, color);
            });
            this.cities.forEach((city) => {
              // TODO: add draw city handler
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
