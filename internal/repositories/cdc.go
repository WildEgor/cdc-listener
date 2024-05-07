package repositories

import (
	"context"
	"errors"
	"github.com/WildEgor/cdc-listener/internal/adapters/listener"
	"github.com/WildEgor/cdc-listener/internal/db/mongodb"
	appErrors "github.com/WildEgor/cdc-listener/internal/errors"
	"github.com/WildEgor/cdc-listener/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
)

var _ listener.ICDCRepository = (*CDCRepository)(nil)

type CDCRepository struct {
	db *mongodb.MongoConnection
}

func NewCDCRepository(db *mongodb.MongoConnection) *CDCRepository {
	return &CDCRepository{db: db}
}

func (r *CDCRepository) CreateCollection(ctx context.Context, opts *listener.CreateCollectionOptions) error {
	db := r.db.DB(opts.DbName)

	colls, err := db.ListCollectionNames(ctx, bson.D{{Key: "name", Value: opts.CollName}})
	if err != nil {
		return appErrors.ErrFailListCollections
	}

	if len(colls) == 0 {
		mongoOpt := options.CreateCollection()
		if opts.Capped {
			mongoOpt.SetCapped(true).SetSizeInBytes(opts.SizeInBytes)
		}

		if err := db.CreateCollection(ctx, opts.CollName, mongoOpt); err != nil {
			return appErrors.ErrFailCreateCollection
		}

		slog.Debug("created mongodb collection")
	}

	if opts.ChangeStreamPreAndPostImages {
		enablePreAndPostImages := bson.D{{Key: "collMod", Value: opts.CollName},
			{Key: "changeStreamPreAndPostImages", Value: bson.D{{Key: "enabled", Value: true}}}}
		if err = db.RunCommand(ctx, enablePreAndPostImages).Err(); err != nil {
			slog.Warn("could not enable changeStreamPreAndPostImages, is your MongoDB version at least 6.0?")
		}
	}

	return nil
}

func (r *CDCRepository) GetResumeToken(collCapped bool) (string, error) {
	findOneOpts := options.FindOne()
	if collCapped {
		// use natural sort for capped collections to get the last inserted resume token
		findOneOpts.SetSort(bson.D{{Key: "$natural", Value: -1}})
	} else {
		// cannot rely on natural sort for uncapped collections, sort by id instead
		findOneOpts.SetSort(bson.D{{Key: "_id", Value: -1}})
	}

	lastResumeToken := &models.ResumeToken{}
	err := r.db.ResumeTokenColl().FindOne(context.TODO(), bson.D{}, findOneOpts).Decode(lastResumeToken)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return "", appErrors.ErrFailFindResumeToken
	}

	return lastResumeToken.Value, nil
}

func (r *CDCRepository) GetWatchStream(opts *options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	cs, err := r.db.Coll("").Watch(context.TODO(), mongo.Pipeline{}, opts)
	if err != nil {
		return nil, appErrors.ErrFailFindChangeStream
	}

	return cs, nil
}

func (r *CDCRepository) SaveResumeToken(token string) error {
	_, err := r.db.ResumeTokenColl().InsertOne(context.TODO(), &models.ResumeToken{
		Value: token,
	})

	return err
}

// TODO
// IsAlive
func (r *CDCRepository) IsAlive() error {
	is := r.db.IsAlive()
	if !is {
		return errors.New("")
	}

	return nil
}

func (r *CDCRepository) Close() error {
	r.db.Disconnect(context.TODO())

	return nil
}
