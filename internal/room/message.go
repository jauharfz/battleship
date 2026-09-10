package room

import "encoding/json"

type InboundPayload struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type Message struct {
	Sender  *Client
	Payload InboundPayload
}

type OutboundPayload struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}
