import { SettlersCore } from "../../../core/types";
import MatchStateManager from "../state";

export function addSettlement(this: MatchStateManager, settlement: SettlersCore.Building) {
  this.settlements.push(settlement);
  this.shouldUpdateUIPart.vertices = true;
}

export function setupSettlement(this: MatchStateManager, vertices: number[]) {
  this.renderer.enableVertices(vertices, true, (vertexID) => {
    this.handler.sendSetupNewSettlement(vertexID);
  });
}

export function matchSettlement(this: MatchStateManager, vertices: number[]) {
  this.renderer.enableVertices(vertices, false, (vertexID) => {
    this.handler.sendMatchNewSettlement(vertexID);
  });
}
