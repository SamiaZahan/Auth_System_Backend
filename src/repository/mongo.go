package repository

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Mongo struct {
	ConnURI string
	DBName  string
}

func (m *Mongo) Connect() error {
	clientOptions := options.Client().ApplyURI(m.ConnURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	defer func() {
		err = client.Disconnect(ctx)
		log.Fatal(err.Error())
	}()

	DB = client.Database(m.DBName)
	return nil
}
