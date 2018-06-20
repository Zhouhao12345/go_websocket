package views

import (
	"net/http"
	"go_ws/models"
	"go_ws/tools"
	"encoding/json"
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
	var roomRows []map[string]string = m.SelectQuery(
		"select room.id as rid , room.desc as des from web_chatroom as room inner join web_chatroom_users as " +
			"chroom on room.id = chroom.chatroom_id where chroom.ggacuser_id = " + userId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roomRows)
	return
}
