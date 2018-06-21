package views

import (
	"net/http"
	"go_ws/models"
	"go_ws/tools"
	"encoding/json"
	"fmt"
)

func APIMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	room := r.FormValue("room_id")
	signed, userId := tools.SingleSign(r)
	if signed == false {
		http.Error(w, "Please sign in firstly!", http.StatusOK)
		return
	}
	fmt.Println(userId)
	m := &models.Models{}

	//todo fixme improve it
	messageRows, err := m.SelectQuery(
		"select message.id as mid , message.content as content, " +
			"message.room_id as rid, message.user_id as from_uid, " +
			"user.username as from_name, guser.avatar_image as image " +
			"from web_chatmessage as message " +
			"inner join auth_user as user on user.id = message.user_id and user.is_active = 1 " +
			"inner join web_ggacuser as guser on guser.user_ptr_id = user.id" +
			" where message.room_id = " + room + " order by message.create_date")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messageRows)
	return
}