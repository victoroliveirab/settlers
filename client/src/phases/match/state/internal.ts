import MatchStateManager from "../state";

export function setAvailableEdges(this: MatchStateManager, edges: number[]) {
  this.availableEdges = edges;
  this.shouldUpdateUIPart.edges = true;
}

export function setAvailableVertices(this: MatchStateManager, vertices: number[]) {
  this.availableVertices = vertices;
  this.shouldUpdateUIPart.vertices = true;
}
