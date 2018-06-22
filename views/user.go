package views

import (
	"net/http"
	"go_ws/tools"
	"encoding/json"
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
