package maps

import (
	"strconv"
	"time"
)

// GetString from a map entry or the default value if not present.
func GetString(m map[string]any, key string, defaultValue string) string {
	value, found := m[key]
	if !found {
		return defaultValue
	}
	if x, ok := value.(string); ok {
		return x
	}
	return defaultValue
}

// GetBool from a map entry or the default value if not present or not a bool.
func GetBool(m map[string]any, key string, defaultValue bool) bool {
	value, found := m[key]
	if !found {
		return defaultValue
	}
	if x, ok := value.(bool); ok {
		return x
	}
	if x, ok := value.(string); ok {
		res, err := strconv.ParseBool(x)
		if err == nil {
			return res
		}
	}
	return defaultValue
}

// GetInt from a map entry or the default value if not present or not an int.
func GetInt(m map[string]any, key string, defaultValue int) int {
	value, found := m[key]
	if !found {
		return defaultValue
	}
	if x, ok := value.(int); ok {
		return x
	}
	if x, ok := value.(float64); ok {
		return int(x)
	}
	if x, ok := value.(string); ok {
		res, err := strconv.Atoi(x)
		if err == nil {
			return res
		}
	}
	return defaultValue
}

// GetInt64 from a map entry or the default value if not present or not an int64.
func GetInt64(m map[string]any, key string, defaultValue int64) int64 {
	value, found := m[key]
	if !found {
		return defaultValue
	}
	if x, ok := value.(int64); ok {
		return x
	}
	if x, ok := value.(int); ok {
		return int64(x)
	}
	if x, ok := value.(float64); ok {
		return int64(x)
	}
	if x, ok := value.(string); ok {
		res, err := strconv.ParseInt(x, 10, 64)
		if err == nil {
			return res
		}
	}
	return defaultValue
}

// GetUint from a map entry or the default value if not present or not a uint.
func GetUint(m map[string]any, key string, defaultValue uint) uint {
	value, found := m[key]
	if !found {
		return defaultValue
	}
	if x, ok := value.(uint); ok {
		return x
	}
	if x, ok := value.(int); ok {
		return uint(x)
	}
	if x, ok := value.(float64); ok {
		return uint(x)
	}
	if x, ok := value.(string); ok {
		res, err := strconv.ParseUint(x, 10, 0)
		if err == nil {
			return uint(res)
		}
	}
	return defaultValue
}

// GetUint64 from a map entry or the default value if not present or not a uint64.
func GetUint64(m map[string]any, key string, defaultValue uint64) uint64 {
	value, found := m[key]
	if !found {
		return defaultValue
	}
	if x, ok := value.(uint64); ok {
		return x
	}
	if x, ok := value.(int); ok {
		return uint64(x)
	}
	if x, ok := value.(float64); ok {
		return uint64(x)
	}
	if x, ok := value.(string); ok {
		res, err := strconv.ParseUint(x, 10, 64)
		if err == nil {
			return res
		}
	}
	return defaultValue
}

// GetFloat64 from a map entry or the default value if not present or not a float64.
func GetFloat64(m map[string]any, key string, defaultValue float64) float64 {
	value, found := m[key]
	if !found {
		return defaultValue
	}
	if x, ok := value.(float64); ok {
		return x
	}
	if x, ok := value.(int); ok {
		return float64(x)
	}
	if x, ok := value.(string); ok {
		res, err := strconv.ParseFloat(x, 64)
		if err == nil {
			return res
		}
	}
	return defaultValue
}

// GetDuration from a map entry or the default value if not present or not a Duration.
func GetDuration(m map[string]any, key string, defaultValue time.Duration) time.Duration {
	value, found := m[key]
	if !found {
		return defaultValue
	}
	if x, ok := value.(time.Duration); ok {
		return x
	}
	if x, ok := value.(string); ok {
		res, err := time.ParseDuration(x)
		if err == nil {
			return res
		}
	}
	return defaultValue
}
