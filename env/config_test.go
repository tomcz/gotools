package env

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

type msTestCfg struct {
	User  string `mapstructure:"test_user_name"`
	Age   int    `mapstructure:"test_user_age"`
	Admin bool   `mapstructure:"test_is_admin"`
	Port  int    `mapstructure:"port"`
	Ratio float64
}

func TestPopulateFromEnv_mapstructure(t *testing.T) {
	os.Setenv("test_user_name", "Homer")
	os.Setenv("test_user_age", "42")
	os.Setenv("test_is_admin", "true")
	cfg := &msTestCfg{Port: 8080, Ratio: 1.1}
	err := PopulateFromEnv(cfg)
	assert.NilError(t, err)
	assert.Equal(t, "Homer", cfg.User)
	assert.Equal(t, 42, cfg.Age)
	assert.Equal(t, true, cfg.Admin)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, 1.1, cfg.Ratio)
}

type envTestCfg struct {
	User  string `env:"test_user_name"`
	Age   int    `env:"test_user_age"`
	Admin bool   `env:"test_is_admin"`
	Port  int    `env:"port"`
	Ratio float64
}

func TestPopulateFromEnv_env(t *testing.T) {
	os.Setenv("test_user_name", "Homer")
	os.Setenv("test_user_age", "42")
	os.Setenv("test_is_admin", "true")
	cfg := &envTestCfg{Port: 8080, Ratio: 1.1}
	err := PopulateFromEnv(cfg, "env")
	assert.NilError(t, err)
	assert.Equal(t, "Homer", cfg.User)
	assert.Equal(t, 42, cfg.Age)
	assert.Equal(t, true, cfg.Admin)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, 1.1, cfg.Ratio)
}
