package views

import (
	"net/http"
	"go_ws/config"
)

func ServeLogin(w http.ResponseWriter, r *http.Request)  {
	http.ServeFile(w, r, config.LOGIN_TEMPLATE)
}