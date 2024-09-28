package mongoDb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hydr0g3nz/e-commerce/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func DBConn(cfg *config.Config) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	dsn := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(dsn),
		options.Client().SetConnectTimeout(time.Second*10),
		options.Client().SetTimeout(time.Second*60))

	if err != nil {
		log.Fatalf("connect to db -> %s failed: %v", dsn, err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("pinging to db -> %s failed: %v", dsn, err)
	}
	log.Printf("connected to db -> %s", dsn)
	return client
}
