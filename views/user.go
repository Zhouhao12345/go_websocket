package views

import (
	"net/http"
	"go_ws/tools"
	"encoding/json"
	"go_ws/models"
	"log"
)

func APIUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	signed, userId := tools.SingleSign(r)
	if signed == false {
		http.Error(w, "Please sign in firstly!", http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userId)
	return
}

func APIUserFocused(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	signed, userId := tools.SingleSign(r)
	key := r.FormValue("key")
	if signed == false {
		http.Error(w, "Please sign in firstly!", http.StatusOK)
		return
	}

	m := &models.Models{}
	//todo fixme improve it
	userRaws, err := m.SelectQuery(
		"select users.id, users.username, gusers.avatar_image from auth_user as users " +
			"inner join web_focus as focus on focus.focus_user_id = users.id and disable = 0 " +
			"inner join web_ggacuser as gusers on gusers.user_ptr_id = users.id " +
				"where focus.user_id = ? and users.username like ?", userId, "%"+key+"%")

	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "DB ERROR", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userRaws)
	return
}
