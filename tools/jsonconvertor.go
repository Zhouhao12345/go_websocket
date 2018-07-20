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

func ApiJsonNormalization(data []map[string]string, return_code int, result string) map[string]interface{}{
	if len(data) == 0 {
		return map[string]interface{}{
			"data": make([]int64, 0),
			"return_code": return_code,
			"result": result,
		}
	} else {
		return map[string]interface{}{
			"data": data,
			"return_code": return_code,
			"result": result,
		}
	}
}