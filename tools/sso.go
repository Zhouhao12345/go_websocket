package tools

import (
	"net/http"
	"go_ws/config"
	"go_ws/cache"
)

func SingleSign(r *http.Request) (signed bool, userId string) {
	sessionId, err := r.Cookie(config.SESSION_COOKIE_KEY)
	if err != nil {
		return false, "0"
	}

	//todo fixme improve it
	sessionValue, err1 := cache.Client.HGet("session:" + sessionId.Value, "id").Result()
	if err1 != nil {
		return false, "0"
	}
	return true, sessionValue
}
