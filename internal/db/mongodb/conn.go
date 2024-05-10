package mongodb

import (
	"context"
	"github.com/WildEgor/cdc-listener/internal/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log/slog"
	"time"
)

type MongoConnection struct {
	client *mongo.Client
	cfg    *configs.MongoConfig
}

func NewMongoConnection(cfg *configs.MongoConfig) *MongoConnection {
	conn := &MongoConnection{
		cfg: cfg,
	}

	conn.Connect(context.TODO())

	return conn
}

func (mc *MongoConnection) Connect(ctx context.Context) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mc.cfg.URI), options.Client().SetConnectTimeout(5*time.Second))
	if err != nil {
		slog.Error("fail connect to mongo", err)
		panic(err)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		slog.Error("fail connect to mongo", err)
		panic(err)
	}

	slog.Info("success connect to mongoDB")

	mc.client = client
}

func (mc *MongoConnection) IsAlive() bool {
	if err := mc.client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return false
	}
	return true
}

func (mc *MongoConnection) Disconnect(ctx context.Context) {
	if mc.client == nil {
		return
	}

	if err := mc.client.Disconnect(ctx); err != nil {
		slog.Error("fail disconnect to mongo", err)
		panic(err)
	}

	slog.Info("connection to mongo closed success")
}

func (mc *MongoConnection) DB(name string) *mongo.Database {
	return mc.client.Database(name)
}

// TODO
func (mc *MongoConnection) ResumeTokenColl() *mongo.Collection {
	return mc.client.Database("test").Collection("tokens")
}

func (mc *MongoConnection) DbColl(db, coll string) *mongo.Collection {
	return mc.client.Database(db).Collection(coll)
}
