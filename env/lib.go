package env

import (
	"os"
	"reflect"
	"strconv"
	"time"
)

const (
	tagPrefix = "env"
)

func Parse(i interface{}) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	k := t.Kind()

	if k == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	num := t.NumField()
	for i := 0; i < num; i++ {
		c := t.Field(i).Tag
		envKey := c.Get(tagPrefix)
		val := os.Getenv(envKey)
		if val == "" {
			continue
		}

		field := reflect.New(reflect.TypeOf(val))
		field.Elem().Set(reflect.ValueOf(val))
		value := v.FieldByName(t.Field(i).Name)
		switch value.Interface().(type) {
		case time.Duration:
			d, err := time.ParseDuration(val)
			if err != nil {
				continue
			}
			value.Set(reflect.ValueOf(d))
		case bool:
			if val == "true" {
				value.SetBool(true)
			} else {
				value.SetBool(false)
			}
		case string:
			value.Set(field.Elem())
		case int:
			i, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			value.SetInt(int64(i))
		}
	}
}
