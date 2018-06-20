package views

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Theatre struct {
	hubs map[string]*Hub

	// Register requests from the clients.
	register chan *Hub

	// Unregister requests from clients.
	unregister chan *Hub
}

func NewTheatre() *Theatre {
	return &Theatre{
		hubs: make(map[string]*Hub),
		register:   make(chan *Hub),
		unregister: make(chan *Hub),
	}
}

func (t *Theatre) Run() {
	for {
		select {
			case hub :=<- t.register:
				t.hubs[hub.room_id] = hub
				go hub.run()
			case hub :=<- t.unregister:
				if _, ok := t.hubs[hub.room_id]; ok {
					delete(t.hubs, hub.room_id)
					close(hub.unregister)
					close(hub.register)
					close(hub.broadcast)
				}
		}
	}
}