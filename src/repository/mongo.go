package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Mongo struct {
	ConnURI string
	DBName  string
}

func (m *Mongo) Connect() (MongoClient *mongo.Client, err error) {
	clientOptions := options.Client().ApplyURI(m.ConnURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		return
	}

	MongoClient = client
	DB = client.Database(m.DBName)
	return MongoClient, err
}

func (m *Mongo) Disconnect() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = MongoClient.Disconnect(ctx)
	return
}
