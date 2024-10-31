package util

import (
	"encoding/json"
	"os"
	"strconv"
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

func fileEntryExists(path string, isDir bool) (bool, error) {
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

	if isDir {
		return info.IsDir(), nil
	}
	return !info.IsDir(), nil
}

func FileExists(file string) (bool, error) {
	return fileEntryExists(file, false)
}

func DirExists(dir string) (bool, error) {
	return fileEntryExists(dir, true)
}
