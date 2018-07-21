package views

import (
	"net/http"
	"go_ws/tools"
	"encoding/json"
	"go_ws/models"
	"log"
	"github.com/satori/go.uuid"
	"go_ws/config"
	"go_ws/cache"
	"encoding/base64"
	"time"
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
		"select users.id, users.username, users.avatar_image from users where users.id = ?", userId)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "DB ERROR", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tools.ApiJsonNormalization(userRaw, 0, "success"))
	return
}

func APILogin(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	result , err := models.SelectQuery(
		"select id , password from users where username = ?", username)
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, "DB ERROR", http.StatusServiceUnavailable)
		return
	}

	if len(result) == 0 {
		json.NewEncoder(w).Encode(
			tools.ApiJsonNormalization(
				make([]map[string]string, 0), -1, "帐号不存在"))
		return
	}

	current_password , err := base64.StdEncoding.DecodeString(result[0]["password"])
	if password != string(current_password) {
		json.NewEncoder(w).Encode(
			tools.ApiJsonNormalization(
				make([]map[string]string, 0), -2, "密码错误"))
		return
	}

	sessionKey := uuid.Must(uuid.NewV4()).String()
	err1 := cache.Client.HSet("session:"+ sessionKey, "id", result[0]["id"]).Err()
	if err1 != nil {
		log.Printf(err.Error())
		http.Error(w, "Cache ERROR", http.StatusServiceUnavailable)
		return
	}
	err2 := cache.Client.Expire("session:"+ sessionKey, config.SESSION_MAX_AGE * time.Second).Err()
	if err2 != nil {
		log.Printf(err.Error())
		http.Error(w, "Cache ERROR", http.StatusServiceUnavailable)
		return
	}

	cookie := http.Cookie{
		Name: config.SESSION_COOKIE_KEY,
		Value: sessionKey,
		Domain: config.DOMAIN ,
		Path: "/",
		HttpOnly: true,
		MaxAge: config.SESSION_MAX_AGE}
	http.SetCookie(w, &cookie)
	json.NewEncoder(w).Encode(tools.ApiJsonNormalization(make([]map[string]string, 0), 0, "success"))
	return
}

func APIRegister(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	password_base64 := base64.StdEncoding.EncodeToString([]byte(password))
	_ , err := models.InsertQuery(
		"insert into users (username , password ) VALUES (? , ?)", username, password_base64)
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, "DB ERROR", http.StatusServiceUnavailable)
		return
	}
	json.NewEncoder(w).Encode(tools.ApiJsonNormalization(make([]map[string]string, 0), 0, "success"))
	return
}

func APILogout(w http.ResponseWriter, r *http.Request, world *World)  {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	signed, userID := tools.SingleSign(r)
	if signed == false {
		http.Error(w, "Please sign in firstly!", http.StatusOK)
		return
	}
	member := world.members[userID]
	world.unregisterMember <- member
	member.mp.unregister <- member

	cookieSession , err := r.Cookie(config.SESSION_COOKIE_KEY)
	if err != nil {
		http.Error(w, "Cookie Lost", http.StatusServiceUnavailable)
		return
	}
	err2 := cache.Client.Del("session:"+ cookieSession.Value).Err()
	if err2 != nil {
		http.Error(w, "Session Lost", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tools.ApiJsonNormalization(make([]map[string]string, 0), 0, "success"))
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
	//_ := r.FormValue("user_id")
	userRaws, err := models.SelectQuery(
		"select users.id, users.username, users.avatar_image from users")
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
