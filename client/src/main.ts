import MatchStateManager from "./phases/match/state";
import PreMatchStateManager from "./phases/pre-match/state";

window.onload = () => {
  const pregameRoot = document.getElementById("wrapper-pregame");
  const root = document.getElementById("wrapper");
  const cookies = new URLSearchParams(document.cookie.replace(/; /g, "&"));
  const userID = cookies.get("settlersucookie");

  if (!pregameRoot || !root || !userID) return;

  const roomID = window.location.pathname.split("/").at(-1);

  const wsURL = `http://localhost:8080/ws?room=${roomID}`;
  const ws = new WebSocket(wsURL);
  ws.onopen = (e) => {
    console.log("websocket connection opened", e);
  };
  ws.onmessage = (e) => {
    // TODO: try-catch this
    const data = JSON.parse(e.data) as { type: string; payload: any };
    console.log(data);
    if (data.type.startsWith("room.")) {
      new PreMatchStateManager(ws, e, pregameRoot, root, userID, roomID!);
    } else if (data.type === "setup.hydrate" || data.type === "game.hydrate") {
      pregameRoot.remove();
      const {
        availableEdges,
        availableVertices,
        currentRoundPlayer,
        dice,
        map,
        players,
        resourceCount,
        roads,
        roundType,
        settlements,
      } = data.payload.state;
      const state = new MatchStateManager(ws, root, userID, "base4", map, players);
      if (dice) state.setDice(dice[0], dice[1]);
      if (resourceCount) state.setResourcesCounts(resourceCount);
      if (roads) Object.values<any>(roads).forEach(state.addRoad);
      if (settlements) Object.values<any>(settlements).forEach(state.addSettlement);
      if (availableEdges) state.setAvailableEdges(availableEdges);
      if (availableVertices) state.setAvailableVertices(availableVertices);

      state.setRoundPlayer(currentRoundPlayer);
      state.setRoundType(roundType);
      // state.set
      state.updateUI();
    }
  };
};
