const PROTOCOL = window.location.protocol == "https:"?"wss:":"ws:";
const WS_URL = `${PROTOCOL}//${window.location.host}/ws`
const WS = new WebSocket(WS_URL);

const STATUS = document.getElementById("status");
const MESSAGE_INPUT = document.getElementById("message-input") as HTMLInputElement;
const LOG = document.getElementById("log");
const SEND_BUTTON = document.getElementById("send-button");

WS.onopen = ()=>{
	WS.send("Hello, World!");
	STATUS!.textContent = "Connected";
};

WS.onmessage = ((event)=>{
	try{
		const li = document.createElement("li");
		const msg:WireMessage = JSON.parse(event.data);
		li.textContent = HandleMessage(msg);
		LOG?.append(li);
	}catch(err){
		console.warn("Received non-JSON frame:",event.data,err)
	}
})

WS.onerror = (event)=>{
	console.error(`SOCKET_ERROR: ${event}`);
	STATUS!.textContent = "Error";
}

WS.onclose = (event)=>{
	console.log(`SOCKET_CLOSE: CODE=${event.code}, REASON=${event.reason}`);
	STATUS!.textContent = "Disconnected"; 
}

SEND_BUTTON!.addEventListener("click",()=>{
	if (WS.readyState != WebSocket.OPEN){
		throw new DOMException();
	}
	WS.send(MESSAGE_INPUT.value);
	MESSAGE_INPUT.textContent = "";
	MESSAGE_INPUT.value = "";
})

interface WireMessage {
	type: string;
	payload?: any;
}

function HandleMessage(msg:WireMessage):string{
	switch (msg.type){
		case "EVENT_BATTLE_START": 
			return "Opponent found! Battle Start";
		case "EVENT_GAME_OVER":
			return "Opponent disconnected. Game Over";
		case "EVENT_CHAT":
			return `Message: ${msg.payload.message}`
		default:
			return `Incoming message: ${msg}`;
	}
}
