package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"log"
)

type config struct {
	ID        string `env:"AGENT_ID"`
	LogLevel  string `env:"LOG_LEVEL"`
	Port      int    `env:"AGENT_PORT" envDefault:"8080"`
	Host      string `env:"AGENT_HOST"`
	AgentPath string
}

func settings(conf *config) error {
	if err := env.Parse(conf); err != nil {
		log.Println(err)
	}

	conf.AgentPath = fmt.Sprintf("http://%s:%d/agent/%s", conf.Host, conf.Port, conf.ID)
	return nil
}
