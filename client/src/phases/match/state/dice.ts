import MatchStateManager from "../state";

export function setDice(this: MatchStateManager, dice1: number, dice2: number) {
  this.dice[0] = dice1;
  this.dice[1] = dice2;
  this.shouldUpdateUIPart.dice = true;
  this.shouldUpdateUIPart.passAction = true;
}
