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
console.log(event);
	let li = document.createElement("li");
	li.textContent = event.data;
	LOG?.append(li);
})

WS.onerror = (event)=>{
	console.error(`SOCKET_ERROR: ${event}`);
	STATUS!.textContent = "Error";
}

WS.onclose = (event)=>{
	console.log(`SOCKET_CLOSE: CODE=${event.code}, REASON=${event.reason}`);
	STATUS!.textContent = "Disconnected"; 
	WS.close(event.code,event.reason);
}

SEND_BUTTON!.addEventListener("click",()=>{
	if (WS.readyState != WebSocket.OPEN){
		throw new DOMException();
	}
	WS.send(MESSAGE_INPUT.value);
	MESSAGE_INPUT.textContent = "";
})
