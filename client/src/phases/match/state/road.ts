import { SettlersCore } from "../../../core/types";
import MatchStateManager from "../state";

export function addRoad(this: MatchStateManager, road: SettlersCore.Building) {
  this.roads.push(road);
  this.shouldUpdateUIPart.edges = true;
}

export function setupRoad(this: MatchStateManager, edges: number[]) {
  this.renderer.enableEdges(edges, (edgeID) => {
    this.handler.sendSetupNewRoad(edgeID);
  });
}
