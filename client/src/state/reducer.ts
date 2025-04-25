import {
  setActiveTradeOffers,
  setBlockedTiles,
  setBuyDevCardAction,
  setCities,
  setCurrentRoundPlayer,
  setDevHand,
  setDevHandCount,
  setDevHandPermissions,
  setDice,
  setDiscard,
  setEdges,
  setHand,
  setKnightUsages,
  setLogs,
  setLongestRoadSizes,
  setMap,
  setMapName,
  setMonopoly,
  setPassAction,
  setPlayerPorts,
  setPlayers,
  setPoints,
  setPorts,
  setResourceCount,
  setRoads,
  setRobbablePlayers,
  setRobber,
  setSettlements,
  setTradeAction,
  setVertices,
  setYearOfPlenty,
} from "./match";
import {
  setEndDatetime,
  setPointsDistribution,
  setRoomName,
  setRoundsPlayed,
  setStartDatetime,
  setStatistics,
} from "./match-report";
import type { SettlersIncomingMessage } from "./messages";
import { setRoom, setRoomCapacity, setRoomParams, setRoomStatus } from "./room";

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
      setRoomName(message.payload.room.id);
      setRoomParams(message.payload.params);
      setRoomCapacity(message.payload.minMaxPlayers);
      break;
    }
    case "room.new-update": {
      setRoom(message.payload.room);
      setRoomName(message.payload.room.id);
      setRoomParams(message.payload.params);
      break;
    }
    case "room.update-capacity.success": {
      setRoom(message.payload.room);
      setRoomName(message.payload.room.id);
      setRoomParams(message.payload.params);
      break;
    }
    case "room.toggle-ready.success": {
      setRoom(message.payload.room);
      setRoomName(message.payload.room.id);
      setRoomParams(message.payload.params);
      break;
    }
    case "room.update-param.success": {
      setRoom(message.payload.room);
      setRoomName(message.payload.room.id);
      setRoomParams(message.payload.params);
      break;
    }
    case "room.start-game.success": {
      setRoomStatus(message.payload.roomStatus);
      setMapName(message.payload.mapName);
      setMap(message.payload.map);
      setPlayers(message.payload.players);
      setPorts(message.payload.ports);
      break;
    }
    case "setup.update-map":
    case "match.update-map": {
      setCities(message.payload.cities);
      setRoads(message.payload.roads);
      setSettlements(message.payload.settlements);
      setBlockedTiles(message.payload.blockedTiles);
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
      setCurrentRoundPlayer(
        message.payload.player,
        message.payload.deadline,
        message.payload.subDeadline,
        message.payload.serverNow,
      );
      break;
    }
    case "match.update-resource-count": {
      setResourceCount(message.payload.resourceCount);
      break;
    }
    case "match.update-dev-hand-count": {
      setDevHandCount(message.payload.devHandCount);
      break;
    }
    case "setup.update-points":
    case "match.update-points": {
      setPoints(message.payload.points);
      break;
    }
    case "setup.update-longest-road-size":
    case "match.update-longest-road-size": {
      setLongestRoadSizes(message.payload.longestRoadSizeByPlayer);
      break;
    }
    case "match.update-knight-usage": {
      setKnightUsages(message.payload.knightUsesByPlayer);
      break;
    }
    case "setup.update-ports":
    case "match.update-ports": {
      setPlayerPorts(message.payload.ports);
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
    case "match.update-buy-dev-card": {
      setBuyDevCardAction(message.payload.enabled);
      break;
    }
    case "match.update-dev-hand": {
      setDevHand(message.payload.devHand);
      break;
    }
    case "match.update-dev-hand-permissions": {
      setDevHandPermissions(message.payload.devHandPermissions);
      break;
    }
    case "match.update-discard-phase": {
      setDiscard(message.payload.discardAmounts, message.payload.enabled);
      break;
    }
    case "match.update-robber-movement": {
      setRobber(message.payload);
      break;
    }
    case "match.update-pick-robbed": {
      setRobbablePlayers(message.payload);
      break;
    }
    case "match.update-monopoly": {
      setMonopoly(message.payload.enabled);
      break;
    }
    case "match.update-year-of-plenty": {
      setYearOfPlenty(message.payload.enabled);
      break;
    }
    case "match.report.success": {
      setStatistics(message.payload.report.statistics);
      break;
    }
    case "setup.hydrate": {
      setRoomStatus(message.payload.roomStatus);

      const edges = message.payload.edgeUpdate.payload;
      setEdges(edges.availableEdges, edges.highlight, edges.enabled);

      const vertices = message.payload.vertexUpdate.payload;
      setVertices(
        vertices.availableSettlementVertices,
        vertices.availableCityVertices,
        vertices.highlight,
        vertices.enabled,
      );

      setCurrentRoundPlayer(
        message.payload.roundPlayerUpdate.payload.player,
        message.payload.roundPlayerUpdate.payload.deadline,
        message.payload.roundPlayerUpdate.payload.subDeadline,
        message.payload.roundPlayerUpdate.payload.serverNow,
      );

      setMapName(message.payload.mapName);
      setMap(message.payload.map);
      setPlayers(message.payload.players);
      setResourceCount(message.payload.resourceCount);
      setDevHandCount(message.payload.devHandCount);
      setPorts(message.payload.ports);

      setRoads(message.payload.mapUpdate.payload.roads);
      setSettlements(message.payload.mapUpdate.payload.settlements);
      setBlockedTiles(message.payload.mapUpdate.payload.blockedTiles);
      break;
    }
    case "match.hydrate": {
      setRoomStatus(message.payload.roomStatus);

      const edges = message.payload.edgeUpdate.payload;
      setEdges(edges.availableEdges, edges.highlight, edges.enabled);
      const vertices = message.payload.vertexUpdate.payload;
      setVertices(
        vertices.availableSettlementVertices,
        vertices.availableCityVertices,
        vertices.highlight,
        vertices.enabled,
      );

      setCurrentRoundPlayer(
        message.payload.roundPlayerUpdate.payload.player,
        message.payload.roundPlayerUpdate.payload.deadline,
        message.payload.roundPlayerUpdate.payload.subDeadline,
        message.payload.roundPlayerUpdate.payload.serverNow,
      );

      setMapName(message.payload.mapName);
      setMap(message.payload.map);
      setPlayers(message.payload.players);
      setPorts(message.payload.ports);

      setCities(message.payload.mapUpdate.payload.cities);
      setRoads(message.payload.mapUpdate.payload.roads);
      setSettlements(message.payload.mapUpdate.payload.settlements);
      setBlockedTiles(message.payload.mapUpdate.payload.blockedTiles);
      setPlayerPorts(message.payload.portsUpdate.payload.ports);

      setHand(message.payload.handUpdate.payload.hand);
      setDevHand(message.payload.devHandUpdate.payload.devHand);
      setDevHandPermissions(message.payload.devHandPermissionsUpdate.payload.devHandPermissions);
      setDice({
        enabled: message.payload.diceUpdate.payload.enabled,
        value: message.payload.diceUpdate.payload.dice,
      });
      setPassAction(message.payload.passActionState.payload.enabled);
      setTradeAction(message.payload.tradeActionState.payload.enabled);
      setActiveTradeOffers(message.payload.tradeOffersUpdate.payload.offers);
      setDiscard(
        message.payload.discardUpdate.payload.discardAmounts,
        message.payload.discardUpdate.payload.enabled,
      );
      setRobber(message.payload.robberMovementUpdate.payload);
      setRobbablePlayers(message.payload.robbablePlayersUpdate.payload);
      setBuyDevCardAction(message.payload.buyDevCardUpdate.payload.enabled);
      setRobber(message.payload.robberMovementUpdate.payload);

      setResourceCount(message.payload.resourceCount);
      setDevHandCount(message.payload.devHandCount);
      setPoints(message.payload.pointsUpdate.payload.points);
      setLongestRoadSizes(message.payload.longestRoadUpdate.payload.longestRoadSizeByPlayer);
      setKnightUsages(message.payload.knightsUsageUpdate.payload.knightUsesByPlayer);

      setYearOfPlenty(message.payload.yearOfPlentyUpdate.payload.enabled);

      break;
    }
    case "over.data": {
      setRoomStatus(message.payload.roomStatus);
      setStatistics(message.payload.report.statistics);
      setPointsDistribution(message.payload.report.pointsDistribution);
      setStartDatetime(message.payload.startDatetime);
      setEndDatetime(message.payload.endDatetime);
      setRoundsPlayed(message.payload.roundsPlayed);
      break;
    }
    case "over.hydrate": {
      setPlayers(message.payload.players);
      setMapName(message.payload.mapName);
      setRoomStatus(message.payload.roomStatus);
      setRoomName(message.payload.roomName);
      setStatistics(message.payload.report.statistics);
      setPointsDistribution(message.payload.report.pointsDistribution);
      setStartDatetime(message.payload.startDatetime);
      setEndDatetime(message.payload.endDatetime);
      setRoundsPlayed(message.payload.roundsPlayed);
      break;
    }
    default: {
      console.warn("Unknown websocket message type:", (message as any).type);
    }
  }
}
