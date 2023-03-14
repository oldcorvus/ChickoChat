package data


// struct representing  client connections
type Broker struct {
	// Registered Clients.
	Clients map[*Client]bool

	// messages from the Clients.
	Notification chan *ChatEvent

	// Register requests from the Clients.
	Join chan *Client

	// Unregister requests from Clients.
	Leave chan *Client

	Room *ChatRoom
}

func NewBroker(room *ChatRoom) *Broker {
	return &Broker{
		Notification: make(chan *ChatEvent),
		Join:         make(chan *Client),
		Leave:        make(chan *Client),
		Clients:      make(map[*Client]bool),
		Room:         room,
	}
}

