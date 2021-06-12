package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCfg struct {
	User  string `mapstructure:"test_user_name"`
	Age   int    `mapstructure:"test_user_age"`
	Admin bool   `mapstructure:"test_is_admin"`
	Port  int    `mapstructure:"port"`
}

func TestPopulateFromEnv(t *testing.T) {
	os.Setenv("test_user_name", "Homer")
	os.Setenv("test_user_age", "42")
	os.Setenv("test_is_admin", "true")
	cfg := &testCfg{Port: 8080}
	err := PopulateFromEnv(cfg)
	if assert.NoError(t, err) {
		assert.Equal(t, "Homer", cfg.User)
		assert.Equal(t, 42, cfg.Age)
		assert.Equal(t, true, cfg.Admin)
		assert.Equal(t, 8080, cfg.Port)
	}
}
