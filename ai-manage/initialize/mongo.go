package initialize

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
	"time"
	"vland.live/app/global"
)

func initMongo(ctx context.Context) {
	config := global.Config.Mongo
	param := fmt.Sprintf("mongodb://%s:%s@%s",
		config.Username,
		config.Password,
		config.Path,
	)
	clientOptions := options.Client().ApplyURI(param).
		SetMinPoolSize(5).
		SetMaxPoolSize(config.MaxPoolSize)

	if global.Config.App.Us() {
		journal := true
		writeConcern := &writeconcern.WriteConcern{
			W:        "majority",
			Journal:  &journal,
			WTimeout: 1000 * time.Second,
		}
		clientOptions.SetWriteConcern(writeConcern)
	}

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalln("mongo connect failed", err)
		return
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalln("mongo ping failed", err)
	}
	global.Mongo = client
	log.Println("Mongo Client init success")
}
