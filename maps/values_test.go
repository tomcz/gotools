package maps

import (
	"crypto/rand"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestGetString(t *testing.T) {
	m := map[string]any{"key": "wibble"}
	assert.Equal(t, GetString(m, "key", "wobble"), "wibble")
	assert.Equal(t, GetString(m, rand.Text(), "foo"), "foo")
}

func TestGetBool(t *testing.T) {
	m := map[string]any{
		"k1": true,
		"k2": "true",
		"k3": "bad",
	}
	assert.Equal(t, GetBool(m, "k1", false), true)
	assert.Equal(t, GetBool(m, "k2", false), true)
	assert.Equal(t, GetBool(m, "k3", false), false)
	assert.Equal(t, GetBool(m, rand.Text(), true), true)
}

func TestGetInt(t *testing.T) {
	m := map[string]any{
		"k1": 2,
		"k2": float64(7),
		"k3": "10",
		"k4": "bad",
	}
	assert.Equal(t, GetInt(m, "k1", 5), 2)
	assert.Equal(t, GetInt(m, "k2", 5), 7)
	assert.Equal(t, GetInt(m, "k3", 5), 10)
	assert.Equal(t, GetInt(m, "k4", 5), 5)
	assert.Equal(t, GetInt(m, rand.Text(), 5), 5)
}

func TestGetInt64(t *testing.T) {
	m := map[string]any{
		"k1": int64(2),
		"k2": 6,
		"k3": float64(7),
		"k4": "10",
		"k5": "bad",
	}
	assert.Equal(t, GetInt64(m, "k1", 5), int64(2))
	assert.Equal(t, GetInt64(m, "k2", 5), int64(6))
	assert.Equal(t, GetInt64(m, "k3", 5), int64(7))
	assert.Equal(t, GetInt64(m, "k4", 5), int64(10))
	assert.Equal(t, GetInt64(m, "k5", 5), int64(5))
	assert.Equal(t, GetInt64(m, rand.Text(), 5), int64(5))
}

func TestGetUint(t *testing.T) {
	m := map[string]any{
		"k1": uint(2),
		"k2": 6,
		"k3": float64(7),
		"k4": "10",
		"k5": "bad",
	}
	assert.Equal(t, GetUint(m, "k1", 5), uint(2))
	assert.Equal(t, GetUint(m, "k2", 5), uint(6))
	assert.Equal(t, GetUint(m, "k3", 5), uint(7))
	assert.Equal(t, GetUint(m, "k4", 5), uint(10))
	assert.Equal(t, GetUint(m, "k5", 5), uint(5))
	assert.Equal(t, GetUint(m, rand.Text(), 5), uint(5))
}

func TestGetUint64(t *testing.T) {
	m := map[string]any{
		"k1": uint64(2),
		"k2": 6,
		"k3": float64(7),
		"k4": "10",
		"k5": "bad",
	}
	assert.Equal(t, GetUint64(m, "k1", 5), uint64(2))
	assert.Equal(t, GetUint64(m, "k2", 5), uint64(6))
	assert.Equal(t, GetUint64(m, "k3", 5), uint64(7))
	assert.Equal(t, GetUint64(m, "k4", 5), uint64(10))
	assert.Equal(t, GetUint64(m, "k5", 5), uint64(5))
	assert.Equal(t, GetUint64(m, rand.Text(), 5), uint64(5))
}

func TestGetFloat64(t *testing.T) {
	m := map[string]any{
		"k1": 2.2,
		"k2": 7.7,
		"k3": "10.3",
		"k4": "bad",
	}
	assert.Equal(t, GetFloat64(m, "k1", 5.1), 2.2)
	assert.Equal(t, GetFloat64(m, "k2", 5.1), 7.7)
	assert.Equal(t, GetFloat64(m, "k3", 5.1), 10.3)
	assert.Equal(t, GetFloat64(m, "k4", 5.1), 5.1)
	assert.Equal(t, GetFloat64(m, rand.Text(), 5.1), 5.1)
}

func TestGetDuration(t *testing.T) {
	m := map[string]any{
		"k1": time.Hour,
		"k2": "1m",
		"k3": "bad",
	}
	assert.Equal(t, GetDuration(m, "k1", time.Second), time.Hour)
	assert.Equal(t, GetDuration(m, "k2", time.Second), time.Minute)
	assert.Equal(t, GetDuration(m, "k3", time.Second), time.Second)
	assert.Equal(t, GetDuration(m, rand.Text(), time.Second), time.Second)
}
