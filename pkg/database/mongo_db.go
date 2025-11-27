package database

import (
	"context"
	"fmt"
	"time"

	"github.com/ArdyJunata/go-realtime-market-data/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func InitMongo() (*mongo.Client, error) {
	credential := options.Credential{
		Username: config.GetString(config.CFG_MONGO_USERNAME),
		Password: config.GetString(config.CFG_MONGO_PASSWORD),
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	uri := config.GetString(config.CFG_MONGO_URI)
	if uri == "" {
		return nil, fmt.Errorf("uri is empty")
	}

	clientOpts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI).SetAuth(credential)

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
