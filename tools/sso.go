package tools

import (
	"net/http"
	"encoding/base64"
	"go_ws/models"
	"strings"
	"encoding/json"
)

func SingleSign(r *http.Request) (signed bool, userId string) {
	sessionId, err := r.Cookie("ggsessionid")
	if err != nil {
		return false, "0"
	}

	//todo fixme improve it
	sessionRow, err := models.SelectQuery(
		"select session_data from django_session where session_key = ?", sessionId.Value)
	if err != nil {
		return false, "0"
	}
	if len(sessionRow) == 0 {
		return false, "0"
	}
	sessionData := sessionRow[0]["session_data"]
	decodeBytes, err := base64.StdEncoding.DecodeString(sessionData)
	if err != nil {
		return false, "0"
	}
	sessionJson := strings.SplitN(string(decodeBytes), ":", 2)[1]
	var sessionMap map[string]interface{}
	if err := json.Unmarshal([]byte(sessionJson), &sessionMap); err !=nil {
		return false, "0"
	}
	if err != nil {
		return false, "0"
	}
	if user_id, ok := sessionMap["_auth_user_id"].(string); ok {
		return true, user_id
	}
	return false, "0"
}
