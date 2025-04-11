import { createContext, useCallback, useEffect, useState } from "react";

import { safeParse } from "./state/helpers";
import { reducer } from "./state/reducer";
import { SettlersOutgoingMessages } from "./state/messages";
import { setUsername } from "./state/player";

type WebSocketStates = "connecting" | "ready" | "error" | "closed";

export const WebSocketContext = createContext(
  {} as {
    instance: WebSocket;
    sendMessage: (message: SettlersOutgoingMessages) => void;
    state: WebSocketStates;
  },
);

const roomID = window.location.pathname.split("/").at(-1);
const wsURL = import.meta.env.VITE_WS_URL;
const wsEndpoint = `${wsURL}?room=${roomID}`;

export const WebSocketProvider = (props: React.PropsWithChildren) => {
  const [state, setState] = useState<WebSocketStates>("connecting");
  const [ws] = useState(() => new WebSocket(wsEndpoint));

  useEffect(() => {
    ws.addEventListener("open", (e) => {
      console.debug("==============WEBSOCKET_OPEN==============");
      console.debug(e);
      console.debug("==========================================");

      const cookies = new URLSearchParams(document.cookie.replace(/; /g, "&"));
      const userID = cookies.get("settlersucookie");
      setUsername(userID);
      setState("ready");
    });

    ws.addEventListener("close", (e) => {
      console.debug("==============WEBSOCKET_CLOSE==============");
      console.debug(e);
      console.debug("===========================================");
      setState("closed");
    });

    ws.addEventListener("error", (e) => {
      console.error("ERROR EVENT", e);
      setState("error");
    });

    ws.addEventListener("message", (event) => {
      const message = safeParse(event.data);
      if (!message) return;
      console.debug("==============WEBSOCKET_MESSAGE==============");
      console.debug(message);
      console.debug("=============================================");
      reducer(message);
    });

    return () => {
      ws.close();
    };
  }, [ws]);

  const sendMessage = useCallback(
    (message: SettlersOutgoingMessages) => {
      ws.send(JSON.stringify(message));
    },
    [ws],
  );

  return (
    <WebSocketContext.Provider
      value={{
        instance: ws,
        sendMessage,
        state,
      }}
    >
      {props.children}
    </WebSocketContext.Provider>
  );
};
