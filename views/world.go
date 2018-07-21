package views

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type World struct {
	maps map[*Map]bool
	members map[string]*Member

	// Register requests from the clients.
	register chan *Map

	// Unregister requests from clients.
	unregister chan *Map

	registerMember chan *Member
	unregisterMember chan *Member

}

func NewWorld() *World {
	return &World{
		maps: make(map[*Map]bool),
		members: make(map[string]*Member),
		register:   make(chan *Map),
		unregister: make(chan *Map),
		registerMember: make(chan *Member),
		unregisterMember: make(chan *Member),
	}
}

func (t *World) Run() {
	for {
		select {
			case mp:= <-t.register:
				t.maps[mp] = true
				go mp.run()
			case mp :=<- t.unregister:
				if _, ok := t.maps[mp]; ok {
					delete(t.maps, mp)
				}
			case member :=<- t.registerMember:
				t.members[member.user] = member
			case member :=<- t.unregisterMember:
				if _, ok := t.members[member.user]; ok {
					delete(t.members, member.user)
					//close(member.move)
					//close(member.mapEnter)
					//close(member.test_connect)
					//close(member.receive_error)
				}
		}
	}
}