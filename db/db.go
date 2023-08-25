package db

import (
	"context"
	"log"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type dbInjectionKey int

const injectionKey dbInjectionKey = iota

func WithMongo(ctx context.Context) (context.Context, func()) {
	dbAddr := viper.GetString("dbAddr")
	if dbAddr == "" {
		log.Fatalln("MongoDB address not set")
	}
	mongo, err := mongo.Connect(ctx, options.Client().ApplyURI(dbAddr))
	if err != nil {
		log.Fatalln(err)
	}
	return context.WithValue(ctx, injectionKey, mongo), func() {
		if err := mongo.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}
}

func LoadMongo(ctx context.Context) *mongo.Client {
	return ctx.Value(injectionKey).(*mongo.Client)
}

func Database(ctx context.Context) *mongo.Database {
	mongo := LoadMongo(ctx)
	cs, _ := connstring.ParseAndValidate(viper.GetString("dbAddr"))
	return mongo.Database(cs.Database)
}

func Collection(ctx context.Context, name string) *mongo.Collection {
	return Database(ctx).Collection(name)
}
