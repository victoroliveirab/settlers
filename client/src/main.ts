import WebSocketConnection from "./websocket";
import PreGameRenderer from "./renderer/pre-game";

window.onload = () => {
  const root = document.getElementById("wrapper");
  const cookies = new URLSearchParams(document.cookie.replace(/; /g, "&"));
  const userID = cookies.get("settlersucookie");

  if (!root || !userID) return;

  const preGameRenderer = new PreGameRenderer(root, userID);
  const roomID = window.location.pathname.split("/").at(-1);

  new WebSocketConnection(`http://localhost:8080/ws?room=${roomID}`, preGameRenderer);
};
