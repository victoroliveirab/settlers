import { SettlersCore } from "../../../core/types";
import MatchStateManager from "../state";

export function addRoad(this: MatchStateManager, road: SettlersCore.Building) {
  this.roads.push(road);
  this.shouldUpdateUIPart.edges = true;
}

export function setupRoad(this: MatchStateManager, edges: number[]) {
  this.renderer.enableEdges(edges, true, (edgeID) => {
    this.handler.sendSetupNewRoad(edgeID);
  });
}

export function matchRoad(this: MatchStateManager, edges: number[]) {
  this.renderer.enableEdges(edges, false, (edgeID) => {
    this.handler.sendMatchNewRoad(edgeID);
  });
}
