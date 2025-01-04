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
func PopulateFromEnv(target any, tag ...string) error {
	tagName := "mapstructure"
	for _, name := range tag {
		tagName = name
	}
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
		k := r.Field(i).Tag.Get(tagName)
		keys[i] = k
	}
	envVars := map[string]string{}
	for _, key := range keys {
		if val, ok := os.LookupEnv(key); ok {
			envVars[key] = val
		}
	}
	config := &mapstructure.DecoderConfig{
		TagName:          tagName,
		Result:           target,
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(envVars)
}
