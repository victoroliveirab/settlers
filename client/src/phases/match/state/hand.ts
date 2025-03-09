import { SettlersCore } from "../../../core/types";
import MatchStateManager from "../state";

export function setHand(this: MatchStateManager, hand: SettlersCore.Hand) {
  this.hand = hand;
  this.shouldUpdateUIPart.hand = true;
}

export function setDevHand(this: MatchStateManager, devHand: SettlersCore.DevHand) {
  this.devHand = devHand;
  this.shouldUpdateUIPart.devHand = true;
}
