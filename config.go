package main

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port string
}

var Cfg *Config

func loadConfig() {
	Cfg = &Config{}

	err := godotenv.Load()
	if err == nil {
		Cfg.Port = os.Getenv("PORT")
	} else {
		Cfg.Port = "8080"
	}
}
