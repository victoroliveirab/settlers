import { SettlersCore } from "../../../core/types";
import MatchStateManager from "../state";

export function addSettlement(this: MatchStateManager, settlement: SettlersCore.Building) {
  this.settlements.push(settlement);
  this.shouldUpdateUIPart.vertices = true;
}

export function setupSettlement(this: MatchStateManager, vertices: number[]) {
  console.log("SETUP SETTLEMENT");
  this.renderer.enableVertices(vertices, (vertexID) => {
    this.handler.sendSetupNewSettlement(vertexID);
  });
}
