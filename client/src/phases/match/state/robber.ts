import MatchStateManager from "../state";

export function enableRobber(this: MatchStateManager, availableTiles: number[]) {
  if (this.currentRoundPlayer === this.userName) {
    this.renderer.makeTilesClickable(availableTiles, (tileID) => {
      this.handler.sendRobberNewPosition(tileID);
    });
  }
}
