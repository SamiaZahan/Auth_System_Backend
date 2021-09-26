package config

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

var (
	Params *Config
)

type Config struct {
	Port                  int    `required:"true"`
	MongoDbConnURI        string `required:"true"`
	MongoDbName           string `required:"true"`
	NotificationSvcDomain string `required:"true"`
}

func New() {
	config := Config{}
	err := envconfig.Process("auth", &config)

	if err != nil {
		log.Fatal(err)
	}

	Params = &config
}
