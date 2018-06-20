package views

import (
	"net/http"
)

func ServeLogin(w http.ResponseWriter, r *http.Request)  {
	http.ServeFile(w, r, "C:/Users/hao.zhou/go/src/go_ws/template/login.html")
}