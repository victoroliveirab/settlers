import { SettlersCore } from "../../../core/types";
import MatchStateManager from "../state";

export function setPlayers(this: MatchStateManager, players: SettlersCore.Player[]) {
  this.players = players;
  this.shouldUpdateUIPart.playerList = true;
}

export function setRoundPlayer(this: MatchStateManager, player: SettlersCore.Player["name"]) {
  this.currentRoundPlayer = player;
  this.shouldUpdateUIPart.playerList = true;
  if (player === this.userName) this.shouldUpdateUIPart.diceAction = true;
}

export function setResourcesCounts(
  this: MatchStateManager,
  counts: Record<SettlersCore.Resource, number>,
) {
  this.resourceCount = counts;
  this.shouldUpdateUIPart.playerList = true;
}

export function setQuantitiesToDiscard(
  this: MatchStateManager,
  quantityByPlayers: Record<string, number>,
) {
  this.discardQuantityByPlayers = quantityByPlayers;
  this.shouldUpdateUIPart.playerList = true;
  this.shouldUpdateUIPart.discard = true;
}
