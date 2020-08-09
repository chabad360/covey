package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func settings(file string) (*config, error) { // TODO: Revamp config system (probably using configuration)
	var exists bool

	if err := godotenv.Load(file); err != nil {
		return nil, err
	}

	conf := config{}
	if conf.AgentID, exists = os.LookupEnv("AGENT_ID"); !exists || conf.AgentID == "" {
		return nil, fmt.Errorf("missing AGENT_ID")
	}

	if conf.LogLevel, exists = os.LookupEnv("LOG_LEVEL"); !exists || conf.LogLevel == "" {
		conf.LogLevel = "INFO"
	}

	var host string
	if host, exists = os.LookupEnv("AGENT_HOST"); !exists || host == "" {
		return nil, fmt.Errorf("missing AGENT_HOST")
	}

	var port string
	if port, exists = os.LookupEnv("AGENT_HOST_POST"); !exists || port == "" {
		port = "8080"
	}

	conf.AgentPath = fmt.Sprintf("http://%s:%s/agent/%s", host, port, conf.AgentID)

	return &conf, nil
}
