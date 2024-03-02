package utils

import "github.com/go-jose/go-jose/v3/json"

func ToJsonString(v any) string {
	buff, _ := json.MarshalIndent(v, "", "\t")
	return string(buff)
}

func StringToBytes(str string) []byte {
	return []byte(str)
}
