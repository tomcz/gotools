package env

import (
	"errors"
	"os"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
)

// PopulateFromEnv populates the target using environment variables and
// mapstructure (https://github.com/go-viper/mapstructure/v2) mappings.
// Please note that the target must be a pointer to a struct.
func PopulateFromEnv(target any) error {
	r := reflect.TypeOf(target)
	if r.Kind() != reflect.Ptr {
		return errors.New("target must be a pointer")
	}
	r = r.Elem()
	if r.Kind() != reflect.Struct {
		return errors.New("target must be a pointer to a struct")
	}
	count := r.NumField()
	keys := make([]string, count)
	for i := 0; i < count; i++ {
		k := r.Field(i).Tag.Get("mapstructure")
		keys[i] = k
	}
	data := map[string]string{}
	for _, key := range keys {
		if val, ok := os.LookupEnv(key); ok {
			data[key] = val
		}
	}
	return mapstructure.WeakDecode(data, target)
}
