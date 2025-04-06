import { Game } from "./pages/game";
import { useWebSocket } from "./hooks/useWebSocket";

export default function App() {
  const ws = useWebSocket();

  // Render a loading screen until the connection is established
  if (ws.state === "connecting") {
    return <div>Connecting to WebSocket...</div>;
  }

  // Once connected, render the App component
  return <Game />;
}
