package app

import (
	"context"
	"github.com/emamulandalib/airbringr-auth/config"
	"github.com/emamulandalib/airbringr-auth/repository"
	log "github.com/sirupsen/logrus"
)

func (app *App) Bootstrap() {
	mongo := repository.Mongo{
		ConnURI: config.Params.MongoDbConnURI,
		DBName:  config.Params.MongoDbName,
	}
	app.Mongo = &mongo

	err := mongo.Connect()

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	ctx := context.Background()
	authRepo := repository.Auth{ctx}
	err = authRepo.CreateUserIndex()

	if err != nil {
		log.Fatal(err.Error())
	}

	return
}
