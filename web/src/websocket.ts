export async function upgradeWebsocket(ws_token: string): Promise<WebSocket> {
  try {
    const protocol = window.location.protocol == "https:" ? "wss:" : "ws:";
    const ws_url = `${protocol}//${window.location.host}/ws?ticket=${ws_token}`;
    const ws = new WebSocket(ws_url);
    return ws;
  } catch (error) {
    console.error("failed to upgrade to websocket connection");
    throw error;
  }
}
