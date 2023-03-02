package data

// struct representing  client connections
type Broker struct {
	// Registered Clients.
	Clients map[*Client]bool

	// messages from the Clients.
	Notification chan []byte

	// Register requests from the Clients.
	join chan *Client

	// Unregister requests from Clients.
	leave chan *Client

	RoomID int
}

func newBroker(ID int) *Broker {
	return &Broker{
		Notification: make(chan []byte),
		join:         make(chan *Client),
		leave:        make(chan *Client),
		Clients:      make(map[*Client]bool),
		RoomID:       ID,
	}
}
