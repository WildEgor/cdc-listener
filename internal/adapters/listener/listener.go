package listener

import (
	"context"
	"github.com/WildEgor/cdc-listener/internal/adapters/publisher"
	"github.com/WildEgor/cdc-listener/internal/configs"
	"github.com/WildEgor/cdc-listener/internal/errors"
	"github.com/WildEgor/cdc-listener/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"sync"
	"time"
)

var _ IListener = (*Listener)(nil)

// Listener wrapper for collections listen logic
type Listener struct {
	eventPubAdapter *publisher.EventPublisherAdapter
	repo            ICDCRepository
	store           *CDCData
	cfg             *configs.ListenerConfig
	saver           ITokenSaver
}

// NewListener create new listener
func NewListener(pub *publisher.EventPublisherAdapter, repo ICDCRepository, cfg *configs.ListenerConfig, saver ITokenSaver) *Listener {
	pool := &sync.Pool{
		New: func() any {
			return &publisher.Event{}
		},
	}

	return &Listener{
		eventPubAdapter: pub,
		repo:            repo,
		saver:           saver,
		store:           NewCDCData(pool),
		cfg:             cfg,
	}
}

// WatchCollection get updates from collections
func (l *Listener) WatchCollection(ctx context.Context, opts *WatchCollectionOptions) error {
	resume := true
	for resume {
		time.Sleep(time.Second)

		csOpts := options.ChangeStream().
			SetFullDocumentBeforeChange(options.WhenAvailable).
			SetFullDocument(options.UpdateLookup)

		token := l.saver.GetResumeToken(opts.WatchedDb, opts.WatchedColl)
		if len(token) != 0 {
			csOpts.SetResumeAfter(bson.M{"_data": token})
		}

		cs, err := l.repo.GetWatchStream(opts, csOpts)
		if err != nil {
			return err
		}

		for cs.Next(ctx) {
			dbEvent := &DbEventData{}
			if err := cs.Decode(dbEvent); err != nil {
				return err
			}

			if dbEvent.OperationType.IsInvalid() {
				resume = false
				slog.Warn("invalid operation type")
				break
			}

			if !dbEvent.OperationType.IsPublishable() {
				continue
			}

			rawEvent := &ChangeEventRaw{
				ID:                        dbEvent.DocumentKey.ID.Hex(),
				Db:                        opts.WatchedDb,
				Coll:                      opts.WatchedColl,
				Kind:                      dbEvent.OperationType,
				FullDocumentBeforeChanges: dbEvent.FullDocumentBeforeChange,
				FullDocument:              dbEvent.FullDocument,
			}

			if err = opts.ChangeEventHandler(ctx, rawEvent); err != nil {
				slog.Error("could not publish change event", "err", err)
				break
			}

			resumeToken := cs.ResumeToken()
			if resumeToken != nil {
				var res models.MongoResumeToken
				err := bson.Unmarshal(resumeToken, &res)
				if err != nil {
					slog.Error("error unmarshalling resume token: " + err.Error())
					continue
				}

				l.saver.SaveResumeToken(&models.ResumeTokenState{
					Db:                     opts.WatchedDb,
					Coll:                   opts.WatchedColl,
					LastMongoResumeToken:   res.Token,
					LastMongoProcessedTime: time.Now(),
				})
			}
		}

		slog.Info("stopped watching mongodb collection")

		if err = cs.Close(context.Background()); err != nil {
			return errors.ErrFailCloseStream
		}
	}

	return nil
}

// Run start listen collections from db
func (l *Listener) Run(ctx context.Context) error {
	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		l.saver.Run()
		return nil
	})

	for subj, _ := range l.cfg.MappedFilter {
		group.Go(func() error {
			meta := l.cfg.GetDbCollBySubject(subj)

			watchCollOpts := &WatchCollectionOptions{
				WatchedDb:   meta.Db,
				WatchedColl: meta.Coll,
				ChangeEventHandler: func(ctx context.Context, rawEvent *ChangeEventRaw) error {
					_, err := l.store.Assert(rawEvent)
					if err != nil {
						return err
					}

					event := l.store.FilterEvent(groupCtx, l.cfg.MappedFilter)
					if event == nil {
						return nil
					}

					topic := l.cfg.GetTopicBySubject(subj)

					return l.eventPubAdapter.Publisher.Publish(groupCtx, topic, event)
				},
			}

			return l.WatchCollection(groupCtx, watchCollOpts)
		})
	}

	return group.Wait()
}

// Stop stop listener and disconnect db
func (l *Listener) Stop() error {

	l.saver.Stop()

	return l.repo.Close()
}
