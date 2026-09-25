import { getGuestJWT } from "./auth.js";
import { enterMatchmaking, poolMatchmaking } from "./matchmaking.js";
import { upgradeWebsocket } from "./websocket.js";

let ws: WebSocket;
let jwt_token: string;

const STATUS = document.getElementById("status");
const MESSAGE_INPUT = document.getElementById(
  "message-input",
) as HTMLInputElement;
const LOG = document.getElementById("log");
const SEND_BUTTON = document.getElementById("send-button");
const MATCH_BUTTON = document.getElementById("match-button");

let playerID: string;

async function mainApp(): Promise<void> {
  try {
    jwt_token = await getGuestJWT();
  } catch (error) {
    console.error("unexpected error: ", error);
    throw error;
  }
}

MATCH_BUTTON!.addEventListener("click", async () => {
  if (!jwt_token || jwt_token == "") {
    console.warn("invalid token");
    return;
  }
  await enterMatchmaking(jwt_token);
  const ws_token = await poolMatchmaking(jwt_token);
  ws = await upgradeWebsocket(ws_token);
  attachWebSocketHandlers(ws);
});

function attachWebSocketHandlers(socket: WebSocket) {
  socket.onopen = () => {
    // WS.send("Hello, World!");
    STATUS!.textContent = "Connected";
  };

  socket.onmessage = (event) => {
    if (!LOG) return;
    try {
      const li = document.createElement("li");
      const msg: WireMessage = JSON.parse(event.data);
      li.textContent = HandleMessage(msg);
      LOG?.append(li);
    } catch (err) {
      console.warn("Received non-JSON frame:", event.data, err);
    }
  };

  socket.onerror = (event) => {
    console.error(`SOCKET_ERROR: ${event}`);
    STATUS!.textContent = "Error";
  };

  socket.onclose = (event) => {
    console.log(`SOCKET_CLOSE: CODE=${event.code}, REASON=${event.reason}`);
    STATUS!.textContent = "Disconnected";
  };
}

mainApp();

SEND_BUTTON!.addEventListener("click", () => {
  if (ws.readyState != WebSocket.OPEN) {
    throw new DOMException();
  }
  ws.send(MESSAGE_INPUT.value);
  MESSAGE_INPUT.textContent = "";
  MESSAGE_INPUT.value = "";
});

interface WireMessage {
  type: string;
  data?: any;
}

//TODO: add try-catch
function HandleMessage(msg: WireMessage): string {
  switch (msg.type) {
    case "EVENT_ROOM_CREATED":
      playerID = msg.data.player_id;
      return "Room Created";
    case "EVENT_PLACE_SHIPMENT_RESULT":
      return msg.data.is_valid
        ? "Ships Succesfully Placed"
        : `Invalid Ship Placements, Reason: ${msg.data.reason}`;

    case "EVENT_PLAYER_READY":
      return `${msg.data.player_id} is${msg.data.is_player_ready ? " Ready" : "n't Ready Yet"}`;
    case "EVENT_BATTLE_START":
      return "The Wheel of Fate is Turning!";
    case "EVENT_FIRE_RESULT":
      // TODO: shloud display/render to player
      if (msg.data.hit) {
        return `${msg.data.target} fired by ${msg.data.fired_by} at [${msg.data.x},${msg.data.y}]`;
      } else {
        return `${msg.data.fired_by} fire, miss hit at $[${msg.data.x},${msg.data.y}]`;
      }
    case "EVENT_ERROR":
      return `${msg.data}`;
    case "EVENT_GAME_OVER":
      if (msg.data.status == "finished") {
        return `${msg.data.winner_id} win, because another player disconnected`;
      }
      return `match is aborted, because another player disconnected`;
    default:
      return `invalid type msg: ${msg}`;
  }
}

//save the client id -> if == unlock, otherwise lock

(window as any).game = {
  placeShips: (ships: Ship[] = defaultShips) => {
    send("CMD_PLACE_SHIPS", ships);
  },
  cancelReady: () => {
    send("CMD_CANCEL_READY", {});
  },
  fire: (x: number = 8, y: number = 9) => {
    send("CMD_FIRE", { x, y });
  },
  surrender: () => {
    send("CMD_SURRENDER", {});
  },
};

type Ship = {
  name: string;
  type: string;
  x: number;
  y: number;
  direction: string;
};

const defaultShips: Ship[] = [
  {
    name: "IJN_Taihou",
    type: "aircraft_carrier",
    x: 5,
    y: 9,
    direction: "Horizontal",
  },
  {
    name: "HMS_Belfast",
    type: "light_cruiser",
    x: 0,
    y: 0,
    direction: "Vertical",
  },
  {
    name: "KMS_Bismarck",
    type: "battleship",
    x: 6,
    y: 0,
    direction: "Horizontal",
  },
  {
    name: "KMS_Prinz_Eugen",
    type: "heavy_cruiser",
    x: 0,
    y: 7,
    direction: "Vertical",
  },
  {
    name: "USS_Laffey",
    type: "destroyer",
    x: 5,
    y: 5,
    direction: "Horizontal",
  },
];

function send(type: string, data?: any) {
  if (ws.readyState != WebSocket.OPEN) {
    console.warn("websocket connection isn't open yet, get out");
    return;
  }
  ws.send(JSON.stringify({ type, data }));
}
