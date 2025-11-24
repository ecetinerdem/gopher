package env

import (
	"os"
	"strconv"
	"time"
)

//direnv allow .

func GetString(key, fallback string) string {
	value, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}

	return value
}

func GetInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}

	valueInt, err := strconv.Atoi(value)

	if err != nil {
		return fallback
	}

	return valueInt
}

func GetDuration(key string, fallback string) time.Duration {
	value, ok := os.LookupEnv(key)

	if !ok {
		duration, err := time.ParseDuration(fallback)
		if err != nil {
			return 15 * time.Minute
		}
		return duration
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 15 * time.Minute
	}
	return duration

}

func GetBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)

	if !ok {
		return fallback
	}
	boolValue, err := strconv.ParseBool(value)

	if err != nil {
		return fallback
	}

	return boolValue
}
