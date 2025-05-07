package util

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
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

func GetenvAsBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	r, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return r
}

func GetenvAsDuration(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	r, err := time.ParseDuration(fmt.Sprintf("%sms", v))
	if err != nil {
		return def
	}
	return r
}

func FileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		if os.IsPermission(err) {
			return true, nil
		}
		return false, err
	}
	isDir := info.IsDir()
	return !isDir, nil
}

func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		if os.IsPermission(err) {
			return true, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}
