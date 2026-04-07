package envs

import (
	"os"
	"strconv"
	"time"
)

// GetString from an environment variable or the default value if not present.
func GetString(envName string, defaultValue string) string {
	value := os.Getenv(envName)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetBool from an environment variable or the default value if not present or not a bool.
func GetBool(envName string, defaultValue bool) bool {
	value := os.Getenv(envName)
	if value == "" {
		return defaultValue
	}
	res, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return res
}

// GetInt from an environment variable or the default value if not present or not an int.
func GetInt(envName string, defaultValue int) int {
	value := os.Getenv(envName)
	if value == "" {
		return defaultValue
	}
	res, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return res
}

// GetInt64 from an environment variable or the default value if not present or not an int64.
func GetInt64(envName string, defaultValue int64) int64 {
	value := os.Getenv(envName)
	if value == "" {
		return defaultValue
	}
	res, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}
	return res
}

// GetUint from an environment variable or the default value if not present or not a uint.
func GetUint(envName string, defaultValue uint) uint {
	value := os.Getenv(envName)
	if value == "" {
		return defaultValue
	}
	res, err := strconv.ParseUint(value, 10, 0)
	if err != nil {
		return defaultValue
	}
	return uint(res)
}

// GetUint64 from an environment variable or the default value if not present or not a uint64.
func GetUint64(envName string, defaultValue uint64) uint64 {
	value := os.Getenv(envName)
	if value == "" {
		return defaultValue
	}
	res, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return defaultValue
	}
	return res
}

// GetFloat64 from an environment variable or the default value if not present or not a float64.
func GetFloat64(envName string, defaultValue float64) float64 {
	value := os.Getenv(envName)
	if value == "" {
		return defaultValue
	}
	res, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}
	return res
}

// GetDuration from an environment variable or the default value if not present or not a Duration.
func GetDuration(envName string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(envName)
	if value == "" {
		return defaultValue
	}
	res, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return res
}
