package views

import (
	"net/http"
	"time"
	"go_ws/tools"
)

func ServeHome(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	signed, userId := tools.SingleSign(r)
	if signed == false {
		http.ServeFile(w, r, "C:/Users/hao.zhou/go/src/go_ws/template/login.html")
		return
	}
	cookieAge := time.Hour * 24 / time.Second
	userCookie:=&http.Cookie{
		Name:   "user_id",
		Value:    userId,
		Path:     "/",
		HttpOnly: false,
		MaxAge:  int(cookieAge),
	}
	http.SetCookie(w, userCookie)
	http.ServeFile(w, r, "C:/Users/hao.zhou/go/src/go_ws/template/home.html")
}
