package app

import (
	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/emamulandalib/airbringr-auth/repository"
	log "github.com/sirupsen/logrus"
)

func (app *App) Bootstrap() {
	mongo := repository.Mongo{
		ConnURI: config.Params.MongoDbConnURI,
		DBName:  config.Params.MongoDbName,
	}

	err := mongo.Connect()

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	authRepo := repository.Auth{}
	err = authRepo.CreateUserIndex()

	if err != nil {
		log.Fatal(err.Error())
	}

	return
}
