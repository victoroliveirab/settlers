import { createRoot } from "react-dom/client";
import { StrictMode } from "react";
import App from "./App";
import { WebSocketProvider } from "./context";
import "./index.css";

// Render the Root component
createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <WebSocketProvider>
      <App />
    </WebSocketProvider>
  </StrictMode>,
);
