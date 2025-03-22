package utils

import "encoding/json"

func StructToString(str any) string {
	data, _ := json.Marshal(str)
	return string(data)
}
