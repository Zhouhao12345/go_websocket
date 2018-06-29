package views

import (
	"go_ws/models"
	"log"
	"time"
	"go_ws/tools"
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
	deleteHub chan []byte

	local *time.Location
}

func NewTheatre() *Theatre {
	return &Theatre{
		hubs: make(map[string]*Hub),
		members: make(map[string]*Member),
		register:   make(chan *Hub),
		unregister: make(chan *Hub),
		registerMember: make(chan *Member),
		unregisterMember: make(chan *Member),
		wakeHub: make(chan *Hub),
		deleteHub:make(chan []byte),
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
					for member := range hub.members {
						hub.unregister <- member
					}
					delete(t.hubs, hub.room_id)
					hub = newHub(t)
				}
			case member :=<- t.registerMember:
				t.members[member.user] = member
			case member :=<- t.unregisterMember:
				if _, ok := t.members[member.user]; ok {
					delete(t.members, member.user)
					close(member.room)
				}
			case room_id := <- t.deleteHub:
				room_id_str := string(room_id)
				m := &models.Models{}
				userRaws, err := m.SelectQuery(
					"select ggacuser_id from web_chatroom_users " +
						"where chatroom_id = ?", room_id_str)
				if err != nil {
					log.Fatalln(err)
				}
				// delete room
				err2 := m.DeleteQuery(
					"delete from web_chatmessage where room_id = ?", room_id_str)
				if err2 != nil {
					log.Println(err2)
				}
				err3 := m.DeleteQuery("delete from web_chatroom_users where chatroom_id = ?", room_id_str)
				if err3 != nil {
					log.Println(err3)
				}
				err4 := m.DeleteQuery("delete from web_chatroom where id = ?", room_id_str)
				if err4 != nil {
					log.Println(err4)
				}
				for _,user := range userRaws {
					user_id := user["ggacuser_id"]
					if member, ok := t.members[user_id]; ok{
						member.room_deleted <- room_id
					}
				}
			case hub :=<- t.wakeHub:
				var inString string = "("
				for member := range hub.members {
					inString = inString + member.user + ","
				}
				inString = inString + "0)"
				m:=&models.Models{}
				userRaws, err := m.SelectQuery(
					"select ggacuser_id from web_chatroom_users " +
						"where ggacuser_id not in "+ inString +" and chatroom_id = ?", hub.room_id)
				if err != nil {
					log.Printf("error: %v", err)
				}
				for _, user := range userRaws {
					userId := user["ggacuser_id"]
					if member,ok := t.members[userId]; ok{
						roomRows, err := m.SelectQuery(
							"select room.id as rid , k.user_names as des, k.images as images , count(message.id) as unread from web_chatroom as room " +
								"inner join web_chatroom_users as chroom on room.id = chroom.chatroom_id " +
								"left join web_chatmessage as message on message.room_id = chroom.chatroom_id and " +
								"message.unread = 1 and message.user_id != ?"+
								" inner join (SELECT GROUP_CONCAT(users.username) as user_names, " +
									"GROUP_CONCAT(gusers.avatar_image) AS images, " +
								"chroom.chatroom_id FROM auth_user AS users " +
								"INNER JOIN web_chatroom_users AS chroom ON chroom.ggacuser_id = users.id " +
								"INNER JOIN web_ggacuser AS gusers ON gusers.user_ptr_id = users.id " +
								"WHERE users.id != ? GROUP BY chroom.chatroom_id) as k on " +
								"k.chatroom_id = chroom.chatroom_id " +
								" where room.id = ? group by room.id", userId, userId, hub.room_id)
						if err != nil {
							log.Printf("error: %v", err)
						}
						messageRows, err1 := m.SelectQuery(
							"select message.content, message.create_date from web_chatmessage as message " +
								"where message.room_id = ? order by create_date desc limit 1", roomRows[0]["rid"])
						if err1 != nil {
							log.Printf("error: %v", err1)
						}
						if len(messageRows) > 0 {
							roomRows[0]["create_date"] = messageRows[0]["create_date"]
							roomRows[0]["content"] = tools.RemoveHtmlTags(messageRows[0]["content"])
						} else {
							roomRows[0]["create_date"] = ""
							roomRows[0]["content"] = ""
						}
						member.room_created <- roomRows[0]
						member.room <- []byte(string(hub.room_id)+"&"+roomRows[0]["create_date"]+"&"+roomRows[0]["content"])
					}
				}
		}
	}
}