package utils

import "encoding/json"

func ToJson(v interface{}) map[string]interface{} {
	b, _ := json.Marshal(v)
	var result map[string]interface{}
	json.Unmarshal(b, &result)
	return result
}
