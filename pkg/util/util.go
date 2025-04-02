package util

import (
	"encoding/json"
	"os"
	"strconv"

	. "github.com/jboelensns/openstack-cni/pkg/logging"
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

func FileExists(path string) (bool, error) {
	log := Log().With().Str("path", path).Logger()
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Error().Err(err).Msg("error is not exist in FileExists")
			return false, nil
		}
		if os.IsPermission(err) {
			log.Error().Err(err).Msg("error is permission FileExists")
			return true, nil
		}
		log.Error().Err(err).Msg("error other FileExists")
		return false, err
	}
	isDir := info.IsDir()
	log.Info().Bool("is_dir", isDir).Msg("is_dir in FileExists")
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
