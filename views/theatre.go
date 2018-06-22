package views

import (
	"go_ws/models"
	"log"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Theatre struct {
	hubs map[string]*Hub
	members map[string]*Member

	// Register requests from the clients.
	register chan *Hub

	// Unregister requests from clients.
	unregister chan *Hub

	registerMember chan *Member
	unregisterMember chan *Member

	wakeHub chan *Hub
	local *time.Location
}

func NewTheatre(local *time.Location) *Theatre {
	return &Theatre{
		hubs: make(map[string]*Hub),
		members: make(map[string]*Member),
		register:   make(chan *Hub),
		unregister: make(chan *Hub),
		registerMember: make(chan *Member),
		unregisterMember: make(chan *Member),
		wakeHub: make(chan *Hub),
		local: local,
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
			case member :=<- t.registerMember:
				t.members[member.user] = member
			case member :=<- t.unregisterMember:
				if _, ok := t.members[member.user]; ok {
					delete(t.members, member.user)
					close(member.room)
				}
			case hub :=<- t.wakeHub:
				var inString string = "("
				for client := range hub.clients {
					inString = inString + client.user + ","
				}
				inString = inString + "0)"
				m:=&models.Models{}
				userRaws, err := m.SelectQuery(
					"select ggacuser_id from web_chatroom_users " +
						"where ggacuser_id not in " + inString + " and chatroom_id = " + hub.room_id)
				if err != nil {
					log.Printf("error: %v", err)
				}
				for _, user := range userRaws {
					userId := user["ggacuser_id"]
					if member,ok := t.members[userId]; ok{
						member.room <- []byte(hub.room_id)
					}
				}
		}
	}
}