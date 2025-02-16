import WebSocketConnection from "./websocket";
import { getCookie } from "./utils/cookie";

const name = getCookie("settlersucookie");
console.log({ name })

const conn = new WebSocketConnection("http://localhost:8080/ws");
