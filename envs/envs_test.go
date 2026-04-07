package envs

import (
	"crypto/rand"
	"os"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestGetString(t *testing.T) {
	key := rand.Text()
	assert.NilError(t, os.Setenv(key, "wibble"))
	assert.Equal(t, GetString(key, "wobble"), "wibble")

	assert.Equal(t, GetString(rand.Text(), "foo"), "foo")
}

func TestGetBool(t *testing.T) {
	key := rand.Text()
	assert.NilError(t, os.Setenv(key, "true"))
	assert.Equal(t, GetBool(key, false), true)

	assert.Equal(t, GetBool(rand.Text(), true), true)

	key2 := rand.Text()
	assert.NilError(t, os.Setenv(key2, "bad"))
	assert.Equal(t, GetBool(key2, false), false)
}

func TestGetInt(t *testing.T) {
	key := rand.Text()
	assert.NilError(t, os.Setenv(key, "10"))
	assert.Equal(t, GetInt(key, 5), 10)

	assert.Equal(t, GetInt(rand.Text(), 5), 5)

	key2 := rand.Text()
	assert.NilError(t, os.Setenv(key2, "bad"))
	assert.Equal(t, GetInt(key2, 5), 5)
}

func TestGetInt64(t *testing.T) {
	key := rand.Text()
	assert.NilError(t, os.Setenv(key, "10"))
	assert.Equal(t, GetInt64(key, 5), int64(10))

	assert.Equal(t, GetInt64(rand.Text(), 5), int64(5))

	key2 := rand.Text()
	assert.NilError(t, os.Setenv(key2, "bad"))
	assert.Equal(t, GetInt64(key2, 5), int64(5))
}

func TestGetUint(t *testing.T) {
	key := rand.Text()
	assert.NilError(t, os.Setenv(key, "10"))
	assert.Equal(t, GetUint(key, 5), uint(10))

	assert.Equal(t, GetUint(rand.Text(), 5), uint(5))

	key2 := rand.Text()
	assert.NilError(t, os.Setenv(key2, "bad"))
	assert.Equal(t, GetUint(key2, 5), uint(5))
}

func TestGetUint64(t *testing.T) {
	key := rand.Text()
	assert.NilError(t, os.Setenv(key, "10"))
	assert.Equal(t, GetUint64(key, 5), uint64(10))

	assert.Equal(t, GetUint64(rand.Text(), 5), uint64(5))

	key2 := rand.Text()
	assert.NilError(t, os.Setenv(key2, "bad"))
	assert.Equal(t, GetUint64(key2, 5), uint64(5))
}

func TestGetFloat64(t *testing.T) {
	key := rand.Text()
	assert.NilError(t, os.Setenv(key, "10.1"))
	assert.Equal(t, GetFloat64(key, 5.2), 10.1)

	assert.Equal(t, GetFloat64(rand.Text(), 5.2), 5.2)

	key2 := rand.Text()
	assert.NilError(t, os.Setenv(key2, "bad"))
	assert.Equal(t, GetFloat64(key2, 5.1), 5.1)
}

func TestGetDuration(t *testing.T) {
	key := rand.Text()
	assert.NilError(t, os.Setenv(key, "1m"))
	assert.Equal(t, GetDuration(key, time.Hour), time.Minute)

	assert.Equal(t, GetDuration(rand.Text(), time.Second), time.Second)

	key2 := rand.Text()
	assert.NilError(t, os.Setenv(key2, "bad"))
	assert.Equal(t, GetDuration(key2, time.Minute), time.Minute)
}
