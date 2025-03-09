import { roundTypesByName } from "../../../core/constants";
import { SettlersCore } from "../../../core/types";
import MatchStateManager from "../state";

export function setRoundPlayer(this: MatchStateManager, player: SettlersCore.Player["name"]) {
  this.currentRoundPlayer = player;
  this.shouldUpdateUIPart.playerList = true;
  if (this.currentRoundPlayer !== this.userName) return;
  if (
    this.roundType === roundTypesByName.FirstRound ||
    this.roundType === roundTypesByName.BetweenRounds
  ) {
    this.shouldUpdateUIPart.diceAction = true;
  } else {
    this.shouldUpdateUIPart.diceAction = false;
  }
}

export function setRoundType(this: MatchStateManager, roundType: number) {
  this.roundType = roundType;
  if (this.currentRoundPlayer !== this.userName) return;
  if (
    this.roundType === roundTypesByName.FirstRound ||
    this.roundType === roundTypesByName.BetweenRounds
  ) {
    this.shouldUpdateUIPart.diceAction = true;
  } else {
    this.shouldUpdateUIPart.diceAction = false;
  }
}
