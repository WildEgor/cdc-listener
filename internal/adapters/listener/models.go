package listener

import (
	"context"
	"github.com/WildEgor/cdc-listener/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (k OperationType) string() string {
	return string(k)
}

func (k OperationType) IsPublishable() bool {
	return k == insertOperationType || k == updateOperationType || k == replaceOperationType || k == deleteOperationType
}

func (k OperationType) IsInvalid() bool {
	return k == invalidateOperationType
}

func (k OperationType) IsInsert() bool {
	return k == insertOperationType
}

type ChangeEventRaw struct {
	ID                        string
	Db                        string
	Coll                      string
	Kind                      OperationType
	FullDocumentBeforeChanges bson.D
	FullDocument              bson.D
}

type ChangeEventHandler func(ctx context.Context, event *ChangeEventRaw) error

type WatchCollectionOptions struct {
	WatchedDb          string
	WatchedColl        string
	ChangeEventHandler ChangeEventHandler
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
	GetResumeToken(db, coll string) string
	SaveResumeToken(token *models.ResumeTokenState) error
	GetWatchStream(watch *WatchCollectionOptions, opts *options.ChangeStreamOptions) (*mongo.ChangeStream, error)
	IsAlive() error
	Close() error
}

type ITokenSaver interface {
	Run()
	Stop()
	GetResumeToken(db, coll string) string
	SaveResumeToken(token *models.ResumeTokenState) error
}

type DbEventData struct {
	CurrentResumeToken       resumeTokenKey `bson:"_id"`
	FullDocument             bson.D         `bson:"fullDocument"`
	FullDocumentBeforeChange bson.D         `bson:"fullDocumentBeforeChange"`
	DocumentKey              documentKey    `bson:"documentKey"`
	OperationType            OperationType  `bson:"operationType"`
}

type documentKey struct {
	ID primitive.ObjectID `bson:"_id"`
}

type resumeTokenKey struct {
	Data string `bson:"_data"`
}
