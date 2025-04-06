import {
  setActiveTradeOffers,
  setCities,
  setCurrentRoundPlayer,
  setDice,
  setEdges,
  setHand,
  setLogs,
  setMap,
  setMapName,
  setPassAction,
  setPlayers,
  setPorts,
  setResourceCount,
  setRoads,
  setSettlements,
  setTradeAction,
  setVertices,
} from "./match";
import type { SettlersIncomingMessage } from "./messages";
import { setRoom, setRoomParams, setRoomStatus } from "./room";

export function reducer(message: SettlersIncomingMessage) {
  if (message.type === "match.bulk-update" || message.type === "setup.bulk-update") {
    for (const msg of message.payload) {
      // TODO: fix type
      reducer(msg as any);
    }
    return;
  }
  switch (message.type) {
    case "setup.update-logs":
    case "match.update-logs": {
      setLogs(message.payload);
      break;
    }
    case "room.connect.success": {
      setRoom(message.payload.room);
      setRoomParams(message.payload.params);
      break;
    }
    case "room.new-update": {
      setRoom(message.payload.room);
      setRoomParams(message.payload.params);
      break;
    }
    case "room.toggle-ready.success": {
      setRoom(message.payload.room);
      setRoomParams(message.payload.params);
      break;
    }
    case "room.update-param.success": {
      setRoom(message.payload.room);
      setRoomParams(message.payload.params);
      break;
    }
    case "room.start-game.success": {
      setRoomStatus("setup"); // This should come from the API
      setMapName(message.payload.mapName);
      setMap(message.payload.map);
      setPlayers(message.payload.players);
      break;
    }
    case "setup.update-map":
    case "match.update-map": {
      setCities(message.payload.cities);
      setRoads(message.payload.roads);
      setSettlements(message.payload.settlements);
      break;
    }
    case "setup.update-edges":
    case "match.update-edges": {
      setEdges(message.payload.availableEdges, message.payload.highlight, message.payload.enabled);
      break;
    }
    case "setup.update-vertices":
    case "match.update-vertices": {
      setVertices(
        message.payload.availableSettlementVertices,
        message.payload.availableCityVertices,
        message.payload.highlight,
        message.payload.enabled,
      );
      break;
    }
    case "match.update-dice": {
      setDice({
        enabled: message.payload.enabled,
        value: message.payload.dice,
      });
      break;
    }
    case "match.update-hand": {
      setHand(message.payload.hand);
      break;
    }
    case "setup.update-round-player":
    case "match.update-round-player": {
      setCurrentRoundPlayer(message.payload.player);
      break;
    }
    case "match.update-resource-count": {
      setResourceCount(message.payload.resourceCount);
      break;
    }
    case "match.update-pass": {
      setPassAction(message.payload.enabled);
      break;
    }
    case "match.update-trade": {
      setTradeAction(message.payload.enabled);
      break;
    }
    case "match.update-trade-offers": {
      setActiveTradeOffers(message.payload.offers);
      break;
    }
    case "setup.hydrate": {
      setRoomStatus("setup"); // This should come from the API

      const edges = message.payload.edgeUpdate.payload;
      setEdges(edges.availableEdges, edges.highlight, edges.enabled);

      const vertices = message.payload.vertexUpdate.payload;
      setVertices(
        vertices.availableSettlementVertices,
        vertices.availableCityVertices,
        vertices.highlight,
        vertices.enabled,
      );

      setCurrentRoundPlayer(message.payload.roundPlayerUpdate.payload.player);

      setMapName(message.payload.mapName);
      setMap(message.payload.map);
      setPlayers(message.payload.players);
      setResourceCount(message.payload.resourceCount);
      setPorts(message.payload.ports);

      setRoads(message.payload.mapUpdate.payload.roads);
      setSettlements(message.payload.mapUpdate.payload.settlements);
      break;
    }
    case "match.hydrate": {
      setRoomStatus("match"); // This should come from the API

      const edges = message.payload.edgeUpdate.payload;
      setEdges(edges.availableEdges, edges.highlight, edges.enabled);
      const vertices = message.payload.vertexUpdate.payload;
      setVertices(
        vertices.availableSettlementVertices,
        vertices.availableCityVertices,
        vertices.highlight,
        vertices.enabled,
      );

      setCurrentRoundPlayer(message.payload.roundPlayerUpdate.payload.player);

      setMapName(message.payload.mapName);
      setMap(message.payload.map);
      setPlayers(message.payload.players);
      setResourceCount(message.payload.resourceCount);
      setPorts(message.payload.ports);

      setCities(message.payload.mapUpdate.payload.cities);
      setRoads(message.payload.mapUpdate.payload.roads);
      setSettlements(message.payload.mapUpdate.payload.settlements);

      setHand(message.payload.handUpdate.payload.hand);
      setDice({
        enabled: message.payload.diceUpdate.payload.enabled,
        value: message.payload.diceUpdate.payload.dice,
      });
      setPassAction(message.payload.passActionState.payload.enabled);
      setTradeAction(message.payload.tradeActionState.payload.enabled);
      setActiveTradeOffers(message.payload.tradeOffersUpdate.payload.offers);
      // TODO: set discard state
      // TODO: set robber movement state

      break;
    }
    default: {
      console.warn("Unknown websocket message type:", message.type);
    }
  }
}
