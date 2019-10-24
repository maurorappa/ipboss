package main

import (
	toml "github.com/BurntSushi/toml"
	"log"
)

type config struct {
	Verbose      bool   `toml:"verbose"`
	EniId        string `toml:"eniId"`
	Interface    string `toml:"interface"`
	EniPrivateIp string `toml:"eniPrivateIp"`
	PublicIp     string `toml:"PublicIp"`
	AwsRegion    string `toml:"AwsRegion"`
	EcsCluster   string `toml:"EcsCluster"`
}

func loadConfig(path string) *config {
	conf := &config{}
	_, err := toml.DecodeFile(path, conf)
	if err != nil {
		log.Fatal(err)
	}

	return conf
}
