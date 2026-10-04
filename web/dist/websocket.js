var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
let ws;
export function upgradeWebsocket(ws_token) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            const protocol = window.location.protocol == "https:" ? "wss:" : "ws:";
            const ws_url = `${protocol}//${window.location.host}/ws?ticket=${ws_token}`;
            ws = new WebSocket(ws_url);
            return ws;
        }
        catch (error) {
            console.error("failed to upgrade to websocket connection");
            throw error;
        }
    });
}
export function send(type, data) {
    if (ws.readyState != WebSocket.OPEN) {
        console.warn("websocket connection isn't open yet, GET OUT!");
        return;
    }
    ws.send(JSON.stringify({ type, data }));
}
export function close() {
    ws.close(1000, "thanks for playing this game, nice to meet you !");
}
export function attachWebSocketHandlers(onMessage) {
    return __awaiter(this, void 0, void 0, function* () {
        ws.onopen = () => {
            // WS.send("Hello, World!");
            // status!.textContent = "Connected";
            console.log("websocket connected");
        };
        ws.onmessage = (event) => {
            try {
                const msg = JSON.parse(event.data);
                onMessage(msg);
            }
            catch (err) {
                console.warn("Received non-JSON frame:", event.data, err);
            }
            ws.onerror = (event) => {
                console.error(`SOCKET_ERROR: ${event}`);
            };
            ws.onclose = (event) => {
                console.log(`SOCKET_CLOSE: CODE=${event.code}, REASON=${event.reason}`);
            };
        };
    });
}
