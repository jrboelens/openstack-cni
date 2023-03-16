package util

import (
	"encoding/json"
	"os"
)

func FromJson(data []byte, i any) error {
	return json.Unmarshal(data, i)
}

func ToJson(i any) ([]byte, error) {
	return json.MarshalIndent(i, "", "  ")
}

func Getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
