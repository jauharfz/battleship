import { defaultShips } from "./constants";
import { WireMessage } from "./types";

let ws: WebSocket;

export async function upgradeWebsocket(ws_token: string): Promise<WebSocket> {
  try {
    const protocol = window.location.protocol == "https:" ? "wss:" : "ws:";
    const ws_url = `${protocol}//${window.location.host}/ws?ticket=${ws_token}`;
    ws = new WebSocket(ws_url);
    return ws;
  } catch (error) {
    console.error("failed to upgrade to websocket connection");
    throw error;
  }
}
export function send(type: string, data?: any) {
  if (ws.readyState != WebSocket.OPEN) {
    console.warn("websocket connection isn't open yet, GET OUT!");
    return;
  }
  ws.send(JSON.stringify({ type, data }));
}

export function close() {
  ws.close(1000, "thanks for playing this game, nice to meet you !");
}

export async function attachWebSocketHandlers(
  onMessage: (msg: WireMessage) => void,
) {
  ws.onopen = () => {
    // WS.send("Hello, World!");
    // status!.textContent = "Connected";
    console.log("websocket connected");
  };

  ws.onmessage = (event) => {
    try {
      const msg: WireMessage = JSON.parse(event.data);
      onMessage(msg);
    } catch (err) {
      console.warn("Received non-JSON frame:", event.data, err);
    }
  };
  ws.onerror = (event) => {
    console.error(`SOCKET_ERROR: ${event}`);
  };

  ws.onclose = (event) => {
    console.log(`SOCKET_CLOSE: CODE=${event.code}, REASON=${event.reason}`);
  };
}
