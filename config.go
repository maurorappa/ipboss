package main

import (
	toml "github.com/BurntSushi/toml"
	"log"
)

type config struct {
	Verbose       bool   `toml:"verbose"`
	Aws           bool   `toml:"aws_instance"`
	Interface     string `toml:"interface"`
	Poll_interval int    `toml:"poll_interval"`
	Timeout       int    `toml:"timeout"`
	Rules         map[string]rule
}

type rule struct {
	IP      string
	Port    int
	Process string
}

func loadConfig(path string) *config {
	conf := &config{}
	_, err := toml.DecodeFile(path, conf)
	if err != nil {
		log.Fatal(err)
	}

	return conf
}
