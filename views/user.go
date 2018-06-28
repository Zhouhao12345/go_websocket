package views

import (
	"net/http"
	"go_ws/tools"
	"encoding/json"
	"go_ws/models"
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

func APILogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	param_userName, found1 := r.Form["username"]
	param_password, found2 := r.Form["password"]
	if !(found1 && found2) {
		http.Error(w, "Please provide username and password", http.StatusBadRequest)
	}
	user := &models.User{}
	logind := user.Auth(param_password[0], param_userName[0], DB)
	if !logind {
		http.Error(w, "Login Failed!", http.StatusBadRequest)
	}
	return
}


func APIRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	param_userName, found1 := r.Form["username"]
	param_password, found2 := r.Form["password"]

	if !(found1 && found2) {
		http.Error(w, "Please provide username and password", http.StatusBadRequest)
	}
	user := &models.User{}
	user.Register(param_password[0], param_userName[0], DB)
	return
}