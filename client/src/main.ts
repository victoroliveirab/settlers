import WebSocketConnection from "./websocket";
import GameState from "./state";

window.onload = () => {
  const pregameRoot = document.getElementById("wrapper-pregame");
  const root = document.getElementById("wrapper");
  const cookies = new URLSearchParams(document.cookie.replace(/; /g, "&"));
  const userID = cookies.get("settlersucookie");

  if (!pregameRoot || !root || !userID) return;

  const roomID = window.location.pathname.split("/").at(-1);

  new WebSocketConnection(
    `http://localhost:8080/ws?room=${roomID}`,
    new GameState(pregameRoot, root, userID),
  );
};
