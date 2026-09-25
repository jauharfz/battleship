var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
export function upgradeWebsocket(ws_token) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            const protocol = window.location.protocol == "https:" ? "wss:" : "ws:";
            const ws_url = `${protocol}//${window.location.host}/ws?ticket=${ws_token}`;
            const ws = new WebSocket(ws_url);
            return ws;
        }
        catch (error) {
            console.error("failed to upgrade to websocket connection");
            throw error;
        }
    });
}
