# env

Populate structs from environment variables using [mapstructure](https://github.com/go-viper/mapstructure).

```go
import (
	"os"
	"testing"
	"github.com/tomcz/gotools/env"
	"gotest.tools/v3/assert"
)

type testCfg struct {
	User  string `env:"test_user_name"`
	Age   int    `env:"test_user_age"`
	Admin bool   `env:"test_is_admin"`
	Port  int    `env:"port"`
}

func TestEnv(t *testing.T) {
	os.Setenv("test_user_name", "Homer")
	os.Setenv("test_user_age", "42")
	os.Setenv("test_is_admin", "true")

	cfg := &testCfg{Port: 8080}

	err := env.PopulateFromEnv(cfg, "env")

	assert.NilError(t, err)
	assert.Equal(t, "Homer", cfg.User)
	assert.Equal(t, 42, cfg.Age)
	assert.Equal(t, true, cfg.Admin)
	assert.Equal(t, 8080, cfg.Port)
}
```
