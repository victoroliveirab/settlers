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
      new MatchStateManager(ws, root, userID, {
        map: data.payload.map,
        resourceCount: data.payload.resourceCount,
        firstPlayer: data.payload.currentRoundPlayer,
        players: data.payload.players,
        logs: [],
        mapName: "base4",
      });
    }
  };
};
