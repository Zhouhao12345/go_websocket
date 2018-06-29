package tools

import (
	"encoding/json"
)

func StringToJson(message string, goObject *map[string]interface{}) error{
	if err := json.Unmarshal([]byte(string(message)), &goObject); err != nil {
		return err
	}
	return nil
}
