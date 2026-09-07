package room

type Message struct {
	Sender *Client
	Data   string
}
