package env

import (
	"log"

	"github.com/joho/godotenv"
)

type Environment struct {
	PORT                      string `env:"PORT"`
	LOGGER_MODE               string `env:"LOGGER_MODE"`
	LINE_CHANNEL_SECRET       string `env:"LINE_CHANNEL_SECRET"`
	LINE_CHANNEL_ACCESS_TOKEN string `env:"LINE_CHANNEL_ACCESS_TOKEN"`
}

var V = Environment{
	PORT:        "8081",
	LOGGER_MODE: "LOCAL",
}

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println(err)
	}

	Parse(&V)

}
