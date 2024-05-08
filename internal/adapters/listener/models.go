package listener

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OperationType string

const (
	insertOperationType     = "insert"
	updateOperationType     = "update"
	replaceOperationType    = "replace"
	deleteOperationType     = "delete"
	invalidateOperationType = "invalidate"
)

var publishableOperationTypes = map[string]struct{}{
	insertOperationType:  {},
	updateOperationType:  {},
	replaceOperationType: {},
	deleteOperationType:  {},
}

func (k OperationType) string() string {
	return string(k)
}

type ChangeEventHandler func(ctx context.Context, subj, msgId string, opType OperationType, data bson.M) error

type WatchCollectionOptions struct {
	WatchedSubj            string
	ResumeTokensCollCapped bool
	StreamName             string
	ChangeEventHandler     ChangeEventHandler
}

type IListener interface {
	WatchCollection(ctx context.Context, opts *WatchCollectionOptions) error
	Run(ctx context.Context) error
	Stop() error
}

type CreateCollectionOptions struct {
	DbName                       string
	CollName                     string
	Capped                       bool
	SizeInBytes                  int64
	ChangeStreamPreAndPostImages bool
}

type ICDCRepository interface {
	CreateCollection(ctx context.Context, opts *CreateCollectionOptions) error
	GetResumeToken(collCapped bool) (string, error)
	SaveResumeToken(token string) error
	GetWatchStream(db, coll string, opts *options.ChangeStreamOptions) (*mongo.ChangeStream, error)
	IsAlive() error
	Close() error
}
