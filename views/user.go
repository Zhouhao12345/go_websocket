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
	userRaw, err := models.SelectQuery(
		"select users.id, users.username, gusers.avatar_image_small from auth_user as users " +
			"inner join web_ggacuser as gusers on gusers.user_ptr_id = users.id " +
			"where users.id = ?", userId)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "DB ERROR", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tools.ApiJsonNormalization(userRaw, 0, "success"))
	return
}

func APIUserDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	signed, _ := tools.SingleSign(r)
	if signed == false {
		http.Error(w, "Please sign in firstly!", http.StatusOK)
		return
	}
	focuseuserId := r.FormValue("user_id")
	userRaws, err := models.SelectQuery(
		"select users.id, users.username, gusers.avatar_image_small, " +
			"IFNULL(k.num_focus, 0) as num_focus, IFNULL(n.num_focused,0) as num_focused " +
			"from auth_user as users " +
			"inner join web_ggacuser as gusers on gusers.user_ptr_id = users.id " +
			"left join ( select user_id, count(id) as num_focus " +
				"from web_focus where disable = 0 group by user_id) as k on k.user_id = users.id " +
			"left join ( select focus_user_id, count(id) as num_focused " +
				"from web_focus where disable = 0 group by focus_user_id) as n on n.focus_user_id = users.id " +
			"where users.id = ? group by users.id", focuseuserId)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "DB ERROR", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tools.ApiJsonNormalization(userRaws, 0, "success"))
	return
}

func APIUserFocusedCancel(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	signed, userId := tools.SingleSign(r)
	if signed == false {
		http.Error(w, "Please sign in firstly!", http.StatusOK)
		return
	}

	r.ParseForm()
	userId_focus := r.Form["user_id"][0]
	err1 := models.UpdateQuery("update web_focus set disable = 1 " +
		"where user_id = ? and focus_user_id = ?", userId, userId_focus)
	if err1 != nil {
		log.Printf("error: %v", err1)
		http.Error(w, "DB ERROR", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tools.ApiJsonNormalization(make([]map[string]string, 0), 0, "success"))
	w.Header().Set("Content-Type", "application/json")
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

	//todo fixme improve it
	userRaws, err := models.SelectQuery(
		"select users.id, users.username, gusers.avatar_image_small from auth_user as users " +
			"inner join web_focus as focus on focus.focus_user_id = users.id and disable = 0 " +
			"inner join web_ggacuser as gusers on gusers.user_ptr_id = users.id " +
				"where focus.user_id = ? and users.username like ?", userId, "%"+key+"%")

	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "DB ERROR", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tools.ApiJsonNormalization(userRaws, 0, "success"))
	return
}

func APIUserFocusedAgree(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	signed, userId := tools.SingleSign(r)
	if signed == false {
		http.Error(w, "Please sign in firstly!", http.StatusOK)
		return
	}

	r.ParseForm()
	userId_focus := r.Form["user_id"][0]

	focusRaws, err := models.SelectQuery("select id from web_focus " +
		"where user_id = ? and focus_user_id = ?", userId, userId_focus)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "DB ERROR", http.StatusInternalServerError)
		return
	}
	if len(focusRaws) != 0 {
		err1 := models.UpdateQuery("update web_focus set disable = 0 " +
			"where user_id = ? and focus_user_id = ?", userId, userId_focus)
		if err1 != nil {
			log.Printf("error: %v", err1)
			http.Error(w, "DB ERROR", http.StatusInternalServerError)
			return
		}
	} else {
		current_time := tools.Now().Format("2006-01-02 15:04:05")
		_ ,err2 := models.InsertQuery("insert into web_focus  VALUES (?,?,?,?,?,?,?,?)",
			0,userId, current_time, userId, current_time, 0, userId_focus , userId)
		if err2 != nil {
			log.Printf("error: %v", err2)
			http.Error(w, "DB ERROR", http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(tools.ApiJsonNormalization(make([]map[string]string, 0), 0, "success"))
	w.Header().Set("Content-Type", "application/json")
	return
}
