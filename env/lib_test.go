package env

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	_ = os.Setenv("ENV", "develop")
	_ = os.Setenv("PORT", "8080")
	_ = os.Setenv("TTL", "1m")
	_ = os.Setenv("DAYS", "1")

	var s struct {
		ENV  string        `env:"ENV"`
		Port string        `env:"PORT"`
		TTL  time.Duration `env:"TTL"`
		DAYS int           `env:"DAYS"`
	}

	expected := struct {
		ENV  string        `env:"ENV"`
		Port string        `env:"PORT"`
		TTL  time.Duration `env:"TTL"`
		DAYS int           `env:"DAYS"`
	}{
		ENV:  "develop",
		Port: "8080",
		TTL:  1 * time.Minute,
		DAYS: 1,
	}

	Parse(&s)

	if !reflect.DeepEqual(expected, s) {
		t.Errorf("%#v is expected but got %#v\n", expected, s)
	}
}
