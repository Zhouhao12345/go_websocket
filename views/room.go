package views

import (
	"net/http"
	"go_ws/models"
	"go_ws/tools"
	"encoding/json"
	"log"
)

func APIRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	signed, userId := tools.SingleSign(r)
	if signed == false {
		http.Error(w, "Please sign in firstly!", http.StatusOK)
		return
	}
	m := &models.Models{}

	//todo fixme improve it
	roomRows, err := m.SelectQuery(
		"select room.id as rid , room.desc as des, count(message.id) as unread from web_chatroom as room " +
			"inner join web_chatroom_users as chroom on room.id = chroom.chatroom_id " +
				"left join web_chatmessage as message on message.room_id = chroom.chatroom_id and " +
					"message.unread = 1 and message.user_id != " + userId +
						" where chroom.ggacuser_id = " + userId + " group by room.id")

	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "DB ERROR", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roomRows)
	return
}
