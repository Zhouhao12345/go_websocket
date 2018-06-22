package views

import (
	"net/http"
	"go_ws/tools"
)

func ServeHome(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	signed, _ := tools.SingleSign(r)
	if signed == false {
		http.Redirect(w,r,"/login", http.StatusFound)
		return
	}
	http.ServeFile(w, r, "C:/Users/hao.zhou/go/src/go_ws/template/home.html")
}
