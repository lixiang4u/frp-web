package utils

import (
	"github.com/go-jose/go-jose/v3/json"
	"strconv"
)

func ToJsonString(v any) string {
	buff, _ := json.MarshalIndent(v, "", "\t")
	return string(buff)
}

func StringToBytes(str string) []byte {
	return []byte(str)
}

func StringToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}
